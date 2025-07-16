package dbus

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/godbus/dbus/v5"
)

// SystemdManager abstracts systemd operations for service control and status checks.
type SystemdManager interface {
	// IsServiceActive checks if the given systemd service is active (running).
	IsServiceActive(ctx context.Context, name string) (bool, error)
	// StartService starts the given systemd service unit.
	StartService(ctx context.Context, name string) error
	// WaitForServiceActive waits until the service is active or timeout/context is done.
	WaitForServiceActive(ctx context.Context, name string, timeout time.Duration) error
}

// DbusSystemdManager implements SystemdManager using D-Bus API.
type DbusSystemdManager struct {
	conn *dbus.Conn
}

// NewDbusSystemdManager creates a new DbusSystemdManager and connects to the system bus.
func NewDbusSystemdManager() (*DbusSystemdManager, error) {
	conn, err := dbus.SystemBus()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to system bus: %w", err)
	}
	return &DbusSystemdManager{conn: conn}, nil
}

// IsServiceActive checks if the given systemd service is active (running).
func (m *DbusSystemdManager) IsServiceActive(ctx context.Context, name string) (bool, error) {
	unitPath, err := m.getUnitPath(ctx, name)
	if err != nil {
		return false, err
	}
	obj := m.conn.Object("org.freedesktop.systemd1", unitPath)
	variant := dbus.Variant{}
	err = obj.CallWithContext(ctx, "org.freedesktop.DBus.Properties.Get", 0, "org.freedesktop.systemd1.Unit", "ActiveState").Store(&variant)
	if err != nil {
		return false, fmt.Errorf("failed to get ActiveState: %w", err)
	}
	state, ok := variant.Value().(string)
	if !ok {
		return false, fmt.Errorf("unexpected type for ActiveState: %T", variant.Value())
	}
	return state == "active", nil
}

// StartService starts the given systemd service unit.
func (m *DbusSystemdManager) StartService(ctx context.Context, name string) error {
	obj := m.conn.Object("org.freedesktop.systemd1", "/org/freedesktop/systemd1")
	var jobPath dbus.ObjectPath
	call := obj.CallWithContext(ctx, "org.freedesktop.systemd1.Manager.StartUnit", 0, name, "replace")
	if call.Err != nil {
		return fmt.Errorf("failed to start unit %s: %w", name, call.Err)
	}
	if err := call.Store(&jobPath); err != nil {
		return fmt.Errorf("failed to parse StartUnit response: %w", err)
	}
	return nil
}

// WaitForServiceActive waits until the service is active or timeout/context is done.
func (m *DbusSystemdManager) WaitForServiceActive(ctx context.Context, name string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		active, err := m.IsServiceActive(ctx, name)
		if err != nil {
			return err
		}
		if active {
			return nil
		}
		if time.Now().After(deadline) {
			return errors.New("timeout waiting for service to become active")
		}
		time.Sleep(200 * time.Millisecond)
	}
}

// getUnitPath resolves the D-Bus object path for a given unit name.
func (m *DbusSystemdManager) getUnitPath(ctx context.Context, name string) (dbus.ObjectPath, error) {
	obj := m.conn.Object("org.freedesktop.systemd1", "/org/freedesktop/systemd1")
	var unitPath dbus.ObjectPath
	call := obj.CallWithContext(ctx, "org.freedesktop.systemd1.Manager.GetUnit", 0, name)
	if call.Err != nil {
		return "", fmt.Errorf("failed to get unit path for %s: %w", name, call.Err)
	}
	if err := call.Store(&unitPath); err != nil {
		return "", fmt.Errorf("failed to parse GetUnit response: %w", err)
	}
	return unitPath, nil
} 