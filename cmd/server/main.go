package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	"github.com/raspidrum-srv/internal/app/preset"
	pb "github.com/raspidrum-srv/internal/pkg/grpc"
	"github.com/raspidrum-srv/internal/repo/db"
	"github.com/raspidrum-srv/internal/repo/dbus"
	lsampler "github.com/raspidrum-srv/internal/repo/linuxsampler"
	lscp "github.com/raspidrum-srv/libs/liblscp-go"
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
}

var cfg Config

func main() {
	setLogging()
	
	var err error
	cfg, err = loadConfig("./configs")
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	projectPath := util.AbsPathify(".")

	sampler, err := connectLinuxSamper(path.Join(projectPath, cfg.Data.Sampler))
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

	// start GRPC server
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

	// Register services
	presetServer := preset.NewPresetServer(db, sampler, fs)
	pb.RegisterKitPresetServer(s, presetServer)
	pb.RegisterChannelControlServer(s, presetServer)

	slog.Info("Server is running", slog.Int("port:", cfg.Host.Port))
	if err := s.Serve(lis); err != nil {
		slog.Error(fmt.Sprintln(fmt.Errorf("Server error: %w", err)))
		os.Exit(1)
	}
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

func setLogging() {
	//logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	//slog.SetDefault(logger)
	var level slog.Level
	switch cfg.Log.Level {
	case "DEBUG": level = slog.LevelDebug
	case "INFO": level = slog.LevelInfo
	case "WARNING": level = slog.LevelWarn 
	case "ERROR": level = slog.LevelError
	default: level = slog.LevelInfo
	}
	slog.SetLogLoggerLevel(level)
}


func connectLinuxSamper(samplesPath string) (*lsampler.LinuxSampler, error) {
	// Initialize sampler
	sampler := lsampler.LinuxSampler{
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
	sampler.Client =  lsClient

	return &sampler, nil
}