package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	"github.com/raspidrum-srv/internal/app/audio"
	"github.com/raspidrum-srv/internal/app/devmonitor"
	"github.com/raspidrum-srv/internal/app/midi"
	"github.com/raspidrum-srv/internal/app/preset"
	pb "github.com/raspidrum-srv/internal/pkg/grpc"
	"github.com/raspidrum-srv/internal/repo/audioprovider"
	"github.com/raspidrum-srv/internal/repo/db"
	"github.com/raspidrum-srv/internal/repo/dbus"
	lsampler "github.com/raspidrum-srv/internal/repo/linuxsampler"
	"github.com/raspidrum-srv/internal/repo/midiprovider"
	"github.com/raspidrum-srv/util"
)

type Config struct {
	Host struct {
		Addr string `mapstructure:"addr"`
		Port int    `mapstructure:"port"`
	} `mapstructure:"host"`
	Data struct {
		DB      string `mapstructure:"dbRoot"`
		Sampler string `mapstructure:"samplerRoot"`
	} `mapstructure:"data"`
	Log struct {
		Level string `mapstructure:"level"`
	} `mapstructure:"log"`
	Audio struct {
		BlackList []string `mapstructure:"blackList"`
	} `mapstructure:"audio"`
}

var cfg Config

func main() {
	setLogging()

	var err error
	cfg, err = loadConfig("./configs")
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	// Initialize systemd manager
	var systemd dbus.SystemdManager
	if runtime.GOOS == "linux" {
		systemd, err = dbus.NewDbusSystemdManager()
		if err != nil {
			slog.Error(fmt.Sprintln("failed to connect to systemd: %w", err))
			os.Exit(1)
		}
	}

	projectPath := util.AbsPathify("", ".")

	samplerDataPath := util.AbsPathify(projectPath, cfg.Data.Sampler)
	slog.Info("Working dir: " + samplerDataPath)
	sampler, err := lsampler.InitLinuxSampler(samplerDataPath, systemd)
	if err != nil {
		slog.Error(fmt.Sprintln(err))
		os.Exit(1)
	}

	// Initialize database
	db, err := db.NewSqlite(cfg.Data.DB)
	if err != nil {
		slog.Error(fmt.Sprintln(fmt.Errorf("Failed to initialize database: %w", err)))
		os.Exit(1)
	}
	defer db.Close()

	// Initialize filesystem
	fs := afero.NewOsFs()

	// Initialize audio provider and get device
	var audioProvider audio.AudioDeviceProvider
	audioProvider, err = audioprovider.NewAudioProvider(cfg.Audio.BlackList)
	if err != nil {
		slog.Error("Failed to initialize audio provider", "error", err)
		os.Exit(1)
	}

	// Get audio device
	audioDev, err := audioProvider.GetAudioDev()
	if err != nil {
		slog.Error("Failed to get audio device", "error", err)
		os.Exit(1)
	}

	slog.Info("Audio device initialized",
		"name", audioDev.Name(),
		"driver", audioDev.Driver())

	// Initialize midi device. Get current and start monitoring for changes
	// Create a context that can be cancelled.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	midiDev := initMidi(ctx, cancel)

	// start gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Host.Addr, cfg.Host.Port))
	if err != nil {
		slog.Error(fmt.Sprintln(fmt.Errorf("Failed to listen: %w", err)))
		os.Exit(1)
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(grpcUnaryLoggingInterceptor), grpc.StreamInterceptor(grpcStreamLoggingInterceptor))

	//cleanup, err := admin.Register(s)
	//if err != nil {
	//	log.Fatalf("failed to register admin: %v", err)
	//}
	//defer cleanup()

	// Register gRPC services
	presetServer := preset.NewPresetServer(db, sampler, fs, midiDev, audioDev)
	pb.RegisterKitPresetServer(s, presetServer)
	pb.RegisterChannelControlServer(s, presetServer)

	slog.Info("Server is running", slog.Int("port:", cfg.Host.Port))
	go func() {
		if err := s.Serve(lis); err != nil {
			slog.Error(fmt.Sprintln(fmt.Errorf("Server error: %w", err)))
			os.Exit(1)
		}
	}()

	// Wait for termination signal before shutting down
	slog.Info("Press Ctrl+C to stop the server")
	sigChan := make(chan os.Signal, 1)
	// Handle SIGINT (Ctrl+C) and SIGTERM (termination) for graceful shutdown
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	slog.Info("Shutting down server...")
	s.GracefulStop()
	cancel() // broadcast the cancellation to all services
	slog.Info("Server gracefully stopped")
}

func loadConfig(configPath string) (Config, error) {
	v := viper.New()
	// get config name from  env variable. default: dev
	configName := os.Getenv("RDRUM_CONFIG")
	if configName == "" {
		configName = "dev"
	}
	v.SetConfigName(configName)
	v.AddConfigPath(configPath)
	v.SetConfigType("yaml")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.BindEnv("log.level", "SRV_LOG_LEVEL")

	if err := v.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return cfg, nil
}

// UnaryInterceptor for grpc logging
func grpcUnaryLoggingInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	// logging incoming request
	slog.Debug(
		"gRPC request",
		slog.String("method", info.FullMethod),
		slog.Any("request", req),
	)

	// handle request
	resp, err := handler(ctx, req)

	// logging response
	if err != nil {
		slog.Error("gRPC",
			slog.String("method", info.FullMethod),
			slog.Any("request", req),
			slog.Any("error", err),
		)
	} else {
		slog.Debug(
			"gRPC response",
			slog.String("method", info.FullMethod),
			slog.Any("response", resp),
		)
	}
	return resp, err
}

// StreamInterceptor for grpc logging
func grpcStreamLoggingInterceptor(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	// logging incoming request
	slog.Debug(
		"gRPC request",
		slog.String("method", info.FullMethod),
		slog.Any("request", srv),
	)

	err := handler(srv, ss)

	// logging response
	if err != nil {
		slog.Error("gRPC",
			slog.String("method", info.FullMethod),
			slog.Any("error", err),
		)
	} else {
		slog.Debug(
			"gRPC response",
			slog.String("method", info.FullMethod),
			slog.Any("response", srv),
		)
	}

	return err
}

func initMidi(ctx context.Context, cancel context.CancelFunc) midi.MIDIDevice {
	// Initialize and start USB monitor service
	devMon, err := devmonitor.NewMonitorService()
	if err != nil {
		slog.Error("Failed to initialize device monitor", "error", err)
		os.Exit(1)
	}
	// Initialize and start the ALSA MIDI provider
	midiPr, err := midiprovider.NewMidiProvider(devMon)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to initialize MIDI provider: %v", err))
		os.Exit(1)
	}
	// Initialize midi device
	midiDev, err := midi.NewMIDIDevice(midiPr)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to initialize MIDI device: %v", err))
		os.Exit(1)
	}

	// Start monitoring for device changes
	go func() {
		if err := devMon.Start(ctx); err != nil {
			slog.Error(fmt.Sprintf("Failed to start device monitor: %v", err))
			cancel()
			os.Exit(1)
		}
	}()

	return midiDev
}

func setLogging() {
	//logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	//slog.SetDefault(logger)
	var level slog.Level
	switch cfg.Log.Level {
	case "DEBUG":
		level = slog.LevelDebug
	case "INFO":
		level = slog.LevelInfo
	case "WARNING":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}
	slog.SetLogLoggerLevel(level)
}
