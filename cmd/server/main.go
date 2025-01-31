package main

import (
	"fmt"
	"io"
	"log"
	"net"

	pb "github.com/raspidrum-srv/internal/pkg/grpc"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

type Config struct {
	Host struct {
		Addr string `mapstructure:"addr"`
		Port int    `mapstructure:"port"`
	} `mapstructure:"host"`
}

var cfg Config

// Server реализует GreeterServer
type Server struct {
	pb.UnimplementedChannelControlServer
}

// TODO: вынести в пакет работы с ChannelControl
func (s *Server) SetValue(stream pb.ChannelControl_SetValueServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		out := &pb.ControlValue{
			Key:   in.Key,
			Seq:   in.Seq,
			Value: in.Value,
		}
		if err := stream.Send(out); err != nil {
			return err
		}

	}
}

func main() {
	cfg, err := loadConfig("./configs")
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Host.Addr, cfg.Host.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	s := grpc.NewServer(opts...)
	pb.RegisterChannelControlServer(s, &Server{})
	log.Printf("Server is running, port: %d", cfg.Host.Port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Server error: %v", err)
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
