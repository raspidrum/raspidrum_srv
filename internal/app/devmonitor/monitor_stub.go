//go:build !linux

package devmonitor

import (
	"context"
	"log/slog"
)

type MonitorService struct {
}

// NewMonitorService creates a new MonitorService.
func NewMonitorService() (*MonitorService, error) {
	return &MonitorService{}, nil
}

func (s *MonitorService) Start(ctx context.Context) error {
	slog.Warn("device monitoring is only supported on Linux")
	return nil
}
