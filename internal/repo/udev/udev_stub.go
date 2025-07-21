//go:build !linux

package udev

import (
	"context"
	"fmt"
)

// Event contains information about a udev event.
type Event struct {
	Action    string
	DevPath   string
	Subsystem string
	DevName   string
	DevType   string
	Env       map[string]string
}

// Monitor represents a udev event monitor.
type Monitor struct{}

// NewMonitor creates a new udev monitor.
func NewMonitor() (*Monitor, error) {
	return &Monitor{}, nil
}

// Close closes the monitor.
func (m *Monitor) Close() {}

// Start begins monitoring for udev events.
func (m *Monitor) Start(ctx context.Context) (<-chan *Event, error) {
	return nil, fmt.Errorf("udev monitoring is only supported on Linux")
}
