package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"

	pb "github.com/raspidrum-srv/internal/pkg/grpc"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	channelcontrol "github.com/raspidrum-srv/internal/app/channel_control"

	lscp "github.com/raspidrum-srv/libs/liblscp-go"
)

type Config struct {
	Host struct {
		Addr string `mapstructure:"addr"`
		Port int    `mapstructure:"port"`
	} `mapstructure:"host"`
}

var cfg Config

func main() {
	cfg, err := loadConfig("./configs")
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	lsClient := lscp.NewClient("localhost", "8888", "1s")
	err = lsClient.Connect()
	if err != nil {
		slog.Error(fmt.Sprintln(fmt.Errorf("Failed connect to LinuxSampler: %w", err)))
		os.Exit(1)
	}

	// start GRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Host.Addr, cfg.Host.Port))
	if err != nil {
		slog.Error(fmt.Sprintln(fmt.Errorf("Failed to listen: %w", err)))
		os.Exit(1)
	}
	var opts []grpc.ServerOption
	s := grpc.NewServer(opts...)
	server := channelcontrol.NewChannelControlServer()
	pb.RegisterChannelControlServer(s, server)
	slog.Info("Server is running", slog.Int("port:", cfg.Host.Port))
	if err := s.Serve(lis); err != nil {
		slog.Error(fmt.Sprintln(fmt.Errorf("Server error: %w", err)))
		os.Exit(1)
	}
}

func loadConfig(configPath string) (*Config, error) {
	v := viper.New()
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
