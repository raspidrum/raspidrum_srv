//go:build linux
// +build linux

package devmonitor

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"time"

	"github.com/raspidrum-srv/internal/repo/udev"
)

const (
	ClassAudio       = 0x01
	ClassMassStorage = 0x08
)

type deviceListenerInfo struct {
	filterSubsystem string
	filterDevPath   *regexp.Regexp
	filterAction    string
	listener        func(devId string, subsystem string, action string, details map[string]string)
}

type DevListener func(devId string, subsystem string, action string, details map[string]string)

// MonitorService orchestrates device monitoring.
type MonitorService struct {
	listeners map[string]deviceListenerInfo
}

// NewMonitorService creates a new MonitorService.
func NewMonitorService() (*MonitorService, error) {
	return &MonitorService{}, nil
}

// Start begins the monitoring process.
func (s *MonitorService) Start(ctx context.Context) error {
	for {
		udevMon, err := udev.NewMonitor()
		if err != nil {
			return err
		}

		eventsCh, errCh, err := udevMon.Start(ctx)
		if err != nil {
			slog.Error("Failed to start monitoring", "error", err)
			return err
		}
		slog.Info("Starting udev device monitoring...")

		monitoringStopped := false

		for !monitoringStopped {
			select {
			case <-ctx.Done():
				slog.Info("Stopping monitor service.")
				udevMon.Close()
				return nil
			case event, ok := <-eventsCh:
				if !ok {
					slog.Info("Udev event channel closed. Restarting in 1s...")
					monitoringStopped = true
					break
				}
				// find listeners and notify them
				slog.Info("Received udev event", "event", event)
				for _, listener := range s.listeners {
					if (listener.filterSubsystem == "" || listener.filterSubsystem == event.Subsystem) &&
						(listener.filterDevPath == nil || listener.filterDevPath.MatchString(event.DevPath)) &&
						(listener.filterAction == "" || listener.filterAction == event.Action) {
						listener.listener(getDeviceIDsFromEvent(event), event.Subsystem, event.Action, event.Env)
					}
				}
			case err, ok := <-errCh:
				if ok && err != nil {
					slog.Warn(fmt.Sprintf("Monitor error: %v", err))
				}
				slog.Warn("Monitoring stopped due to error. Restarting in 1s...")
				monitoringStopped = true
			}
		}
		udevMon.Close()
		time.Sleep(2 * time.Second)
		slog.Info("Restarting udev monitor...")
	}
}

func getDeviceIDsFromEvent(event *udev.Event) string {
	return event.DevPath
}

func (s *MonitorService) AddListener(key, action, subsystem, devPath string, listener DevListener) error {
	if _, exists := s.listeners[key]; exists {
		slog.Warn("Listener already exists for this filter", "key", key)
		return nil
	}
	if s.listeners == nil {
		s.listeners = make(map[string]deviceListenerInfo)
	}
	devPathRegex, err := regexp.Compile(devPath)
	if err != nil {
		return fmt.Errorf("failed to compile regex for devPath: %w", err)
	}

	// Add the listener to the list
	s.listeners[key] = deviceListenerInfo{
		filterSubsystem: subsystem,
		filterDevPath:   devPathRegex,
		filterAction:    action,
		listener:        listener,
	}
	slog.Info("Added new device listener", "key", key)
	return nil
}

func (s *MonitorService) RemoveListener(key string) {
	if _, exists := s.listeners[key]; !exists {
		slog.Warn("Listener not found for this filter", "key", key)
		return
	}
	delete(s.listeners, key)
	slog.Info("Removed device listener", "key", key)
}
