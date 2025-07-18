package linuxsampler

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"sync"
	"time"

	repo "github.com/raspidrum-srv/internal/repo"
	"github.com/raspidrum-srv/internal/repo/dbus"
	lscp "github.com/raspidrum-srv/libs/liblscp-go"
)

// Engine - LinuxSampler engine: gig, sfz, sf2
type LinuxSampler struct {
	Client  lscp.Client
	Engine  string
	DataDir string // root dir for sfz-files, samples and presets
	// Systemd is used to control and check the state of the linuxsampler systemd service. Linux only
	Systemd dbus.SystemdManager

	//clientMu          sync.Mutex
	healthcheckCancel context.CancelFunc
	healthcheckWg     sync.WaitGroup
}

func InitLinuxSampler(samplesPath string) (*LinuxSampler, error) {

	// Initialize sampler
	sampler := LinuxSampler{
		Engine:  "sfz",
		DataDir: samplesPath,
	}

	if runtime.GOOS == "linux" {
		// startup linuxsampler service
		// TODO: move to config
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Initialize systemd manager
		systemd, err := dbus.NewDbusSystemdManager()
		if err != nil {
			return nil, fmt.Errorf("failed to connect to systemd: %w", err)
		}

		// Ensure linuxsampler service is running
		sampler.Systemd = systemd
		if err := sampler.EnsureLinuxSamplerRunning(ctx); err != nil {
			return nil, fmt.Errorf("failed to ensure linuxsampler service is running: %w", err)
		}
	}

	// Initialize LinuxSampler client
	lsClient := lscp.NewClient("localhost", "8888", "1s")
	err := lsClient.Connect()
	if err != nil {
		return nil, fmt.Errorf("Failed connect to LinuxSampler: %w", err)
	}

	// Initialize sampler
	sampler.Client = lsClient

	if runtime.GOOS == "linux" {
		sampler.StartHealthCheck(context.Background())
	}

	return &sampler, nil
}

// Connect
// params grouped by audio channels. Audio channel is key of map
func (l *LinuxSampler) ConnectAudioOutput(driver string, params map[int][]repo.Param[string]) (devId int, err error) {
	devId, err = l.Client.CreateAudioOutputDevice(driver)
	if err != nil {
		return
	}
	if len(params) != 0 {
		// key (k) - channelId
		// value (v) - array of channel params
		for k, v := range params {
			for _, p := range v {
				prm := lscp.Parameter[any]{
					Name:  p.Name,
					Value: p.Value,
				}
				err = l.Client.SetAudioOutputChannelParameter(devId, k, prm)
				if err != nil {
					return
				}
			}
		}
	}
	return
}

// Connect to MIDI port and optional set port parameters (i.e. bindings)
func (l *LinuxSampler) ConnectMidiInput(driver string, params []repo.Param[string]) (devId int, err error) {
	devId, err = l.Client.CreateMidiInputDevice(driver)
	if err != nil {
		return
	}
	if len(params) != 0 {
		for _, p := range params {
			prm := lscp.Parameter[any]{
				Name:  p.Name,
				Value: p.Value,
			}
			err = l.Client.SetMidiInputPortParameter(devId, 0, prm)
			if err != nil {
				return
			}
		}
	}

	return
}

func (l *LinuxSampler) CreateChannel(audioDevId, midiDevId int) (channelId int, err error) {
	channelId, err = l.Client.AddSamplerChannel()
	if err != nil {
		return
	}
	err = l.Client.SetChannelAudioOutputDevice(channelId, audioDevId)
	if err != nil {
		return
	}
	err = l.Client.SetChannelMidiInputDevice(channelId, midiDevId)
	if err != nil {
		return
	}
	err = l.Client.LoadSamplerEngine(l.Engine, channelId)
	if err != nil {
		return
	}
	return
}

func (l *LinuxSampler) LoadInstrument(instrumentFile string, instrIdx int, channelId int) error {
	return l.Client.LoadInstrument(instrumentFile, 0, channelId)
}

func (l *LinuxSampler) SetChannelVolume(samplerChn int, volume float32) error {
	return l.Client.SetChannelVolume(samplerChn, volume)
}

func (l *LinuxSampler) SendMidiCC(samplerChn int, cc int, value float32) error {
	return l.Client.SendChannelMidiData(samplerChn, "CC", cc, int(value))
}

func (l *LinuxSampler) SetGlobalVolume(volume float32) error {
	return l.Client.SetVolume(volume)
}

// EnsureLinuxSamplerRunning checks if the linuxsampler systemd service is running, starts it if needed, and waits for it to become active.
// It uses the provided context for cancellation and timeout.
func (l *LinuxSampler) EnsureLinuxSamplerRunning(ctx context.Context) error {
	const serviceName = "linuxsampler.service"
	active, err := l.Systemd.IsServiceActive(ctx, serviceName)
	if err != nil {
		return fmt.Errorf("systemd check failed: %w", err)
	}
	if active {
		return nil // Already running
	}
	if err := l.Systemd.StartService(ctx, serviceName); err != nil {
		return fmt.Errorf("failed to start linuxsampler service: %w", err)
	}
	// Wait up to 10 seconds for the service to become active
	if err := l.Systemd.WaitForServiceActive(ctx, serviceName, 10*time.Second); err != nil {
		return fmt.Errorf("linuxsampler service did not become active: %w", err)
	}
	return nil
}

// StartHealthCheck launches a background goroutine that checks the connection to LinuxSampler every 2 seconds.
// On connection loss, it attempts to restart the service and reconnect the client as needed.
// hc — клиент для healthcheck (может быть mock в тестах). Если nil, используется l.Client.
func (l *LinuxSampler) StartHealthCheck(ctx context.Context) {
	if l.healthcheckCancel != nil {
		// Already running
		return
	}
	ctx, cancel := context.WithCancel(ctx)
	l.healthcheckCancel = cancel
	l.healthcheckWg.Add(1)
	go func() {
		defer l.healthcheckWg.Done()
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				//l.clientMu.Lock()
				client := l.Client
				//l.clientMu.Unlock()
				err := client.Ping()
				if err == nil {
					continue
				}
				// Connection lost, try to recover
				recoverCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
				defer cancel()
				errRun := l.EnsureLinuxSamplerRunning(recoverCtx)
				if errRun != nil {
					slog.Error("[HealthCheck] Failed to ensure linuxsampler running", slog.Any("error", errRun))
					continue
				}
				// If service was restarted, need reconnect
				if err := client.Connect(); err != nil {
					slog.Error("[HealthCheck] Failed to reconnect to linuxsampler", slog.Any("error", err))
					continue
				}
				slog.Info("[HealthCheck] Reconnected to linuxsampler")
			}
		}
	}()
}

// StopHealthCheck stops the background healthcheck goroutine and waits for it to finish.
func (l *LinuxSampler) StopHealthCheck() {
	if l.healthcheckCancel != nil {
		l.healthcheckCancel()
		l.healthcheckCancel = nil
		l.healthcheckWg.Wait()
	}
}
