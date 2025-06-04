package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"path"
	"runtime"

	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	"github.com/raspidrum-srv/internal/app/preset"
	pb "github.com/raspidrum-srv/internal/pkg/grpc"
	"github.com/raspidrum-srv/internal/repo/db"
	lsampler "github.com/raspidrum-srv/internal/repo/linuxsampler"
	lscp "github.com/raspidrum-srv/libs/liblscp-go"
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
}

var cfg Config

func main() {
	cfg, err := loadConfig("./configs")
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	_, projectPath, _, _ := runtime.Caller(0)
	projectPath = path.Join(path.Dir(projectPath), "../../")

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Initialize LinuxSampler client
	lsClient := lscp.NewClient("localhost", "8888", "1s")
	err = lsClient.Connect()
	if err != nil {
		slog.Error(fmt.Sprintln(fmt.Errorf("Failed connect to LinuxSampler: %w", err)))
		os.Exit(1)
	}

	// Initialize sampler
	sampler := &lsampler.LinuxSampler{
		Client:  lsClient,
		Engine:  "sfz",
		DataDir: path.Join(projectPath, cfg.Data.Sampler),
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
	var opts []grpc.ServerOption
	s := grpc.NewServer(opts...)

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

func loadConfig(configPath string) (*Config, error) {
	v := viper.New()
	// TODO: get config name from  env variable
	v.SetConfigName("dev")
	v.AddConfigPath(configPath)
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}
