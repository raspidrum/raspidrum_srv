//go:build linux
// +build linux

package alsa

import (
	"fmt"
	"log/slog"
	"regexp"
	"strconv"

	"github.com/raspidrum-srv/internal/app/devmonitor"
	"github.com/raspidrum-srv/internal/app/midi"
	libalsa "github.com/raspidrum-srv/libs/libalsa"
)

const (
	uDevMidiSubsystem = "sound"
	uDevDevPath       = `/card(\d+)/midiC(\d+)D(\d+)`
)

type AlsaMidiProvider struct {
	deviceListener func(devId string, state midi.MIDIPortState)
	devPathRegexp  *regexp.Regexp
}

func NewAlsaMidiProvider(monitor *devmonitor.MonitorService) (*AlsaMidiProvider, error) {
	p := AlsaMidiProvider{
		deviceListener: nil,
	}
	var err error
	p.devPathRegexp, err = regexp.Compile(`sound/card(\d+)/midiC(\d+)D(\d+)`)
	if err != nil {
		return nil, fmt.Errorf("failed to compile udev device path regex: %w", err)
	}

	monitor.AddListener("addMidi", "add", uDevMidiSubsystem, uDevDevPath, p.HandleUDevEvent)
	monitor.AddListener("removeMidi", "remove", uDevMidiSubsystem, uDevDevPath, p.HandleUDevEvent)
	return &p, nil
}

// Return only hardware ports
func (p *AlsaMidiProvider) GetMIDIPorts() ([]midi.MIDIPortInfo, error) {
	ports, err := libalsa.ListMidiPorts()
	if err != nil {
		return nil, fmt.Errorf("failed getting MIDI ports from ALSA: %w", err)
	}

	res := make([]midi.MIDIPortInfo, 0)
	for _, port := range ports {
		// TODO: filter port: only hardware needed
		mport := midi.MIDIPortInfo{
			Driver:     "ALSA",
			DevId:      p.getDevIdFromMidiPort(port),
			PortId:     fmt.Sprintf("%d:%d", port.ClientId, port.PortId),
			Name:       port.ClientName,
			State:      midi.MIDIPortStateConnected,
			DeviceType: midi.MIDIPortTypeHardware,
		}
		res = append(res, mport)
	}
	return res, nil
}

func (p *AlsaMidiProvider) getDevIdFromMidiPort(port libalsa.MidiPortInfo) string {
	return fmt.Sprintf("%d:%d", port.CardId, port.PortId)
}

func (p *AlsaMidiProvider) SubscribeDeviceState(listener func(devId string, state midi.MIDIPortState)) error {
	if listener == nil {
		return fmt.Errorf("listener cannot be nil")
	}
	p.deviceListener = listener
	return nil
}

func (p *AlsaMidiProvider) UnsubscribeDeviceState() error {
	p.deviceListener = nil
	return nil
}

func (p *AlsaMidiProvider) notifyDeviceState(devId string, state midi.MIDIPortState) {
	if p.deviceListener != nil {
		p.deviceListener(devId, state)
	}
}

func (p *AlsaMidiProvider) HandleUDevEvent(devId string, subsystem string, action string, details map[string]string) {
	if subsystem != uDevMidiSubsystem {
		return
	}
	did, ok := getDevIdFromUDevEvent(p.devPathRegexp, devId)
	if !ok {
		return
	}
	switch action {
	case "add":
		p.notifyDeviceState(did, midi.MIDIPortStateConnected)
	case "remove":
		p.notifyDeviceState(did, midi.MIDIPortStateDisconnected)
	default:
		slog.Info("Unhandled udev action: %s for device %s\n", action, devId)
	}
}

func atoi(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}

func getDevIdFromUDevEvent(rg *regexp.Regexp, devId string) (string, bool) {
	matches := rg.FindStringSubmatch(devId)
	if len(matches) == 4 {
		card := atoi(matches[1])
		port := atoi(matches[3])
		return fmt.Sprintf("%d:%d", card, port), true
	}
	return "", false
}
