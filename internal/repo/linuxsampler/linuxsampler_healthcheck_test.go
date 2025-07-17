package linuxsampler

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/raspidrum-srv/libs/liblscp-go"
)

// mockSystemdManager implements dbus.SystemdManager
// (IsServiceActive, StartService, WaitForServiceActive)
type mockSystemdManager struct {
	ensureErr  atomic.Value // error
	called     atomic.Bool
	failActive atomic.Bool // если true, IsServiceActive всегда возвращает false и ошибку
}

func (m *mockSystemdManager) IsServiceActive(ctx context.Context, name string) (bool, error) {
	if m.failActive.Load() {
		return false, errors.New("systemd fail")
	}
	return true, nil
}
func (m *mockSystemdManager) StartService(ctx context.Context, name string) error {
	return nil
}
func (m *mockSystemdManager) WaitForServiceActive(ctx context.Context, name string, timeout time.Duration) error {
	return nil
}

func TestHealthCheck_Success(t *testing.T) {
	lscpDrv := &mockLscpDriver{}
	lscpDrv.pingErr.Store(errNoError)
	client := liblscp.NewClientWithDriver(lscpDrv)
	s := &LinuxSampler{
		Client:  client,
		Systemd: &mockSystemdManager{},
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s.StartHealthCheck(ctx)
	time.Sleep(2500 * time.Millisecond)
	s.StopHealthCheck()
	// Should not call EnsureLinuxSamplerRunning
	// (мы не проверяем вызовы SystemdManager в этом тесте)
}

func TestHealthCheck_ReconnectOnFailure(t *testing.T) {
	lscpDrv := &mockLscpDriver{}
	lscpDrv.pingErr.Store(errors.New("fail"))
	lscpDrv.connectErr.Store(errNoError)
	client := liblscp.NewClientWithDriver(lscpDrv)

	ms := &mockSystemdManager{}
	ms.ensureErr.Store(errNoError)

	s := &LinuxSampler{
		Client:  client,
		Systemd: ms,
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s.StartHealthCheck(ctx)
	time.Sleep(2500 * time.Millisecond)
	s.StopHealthCheck()
	if !lscpDrv.connected.Load() {
		t.Error("Client should be reconnected on failure")
	}
}

func TestHealthCheck_EnsureFail(t *testing.T) {
	lscpDrv := &mockLscpDriver{}
	lscpDrv.pingErr.Store(errors.New("fail"))
	lscpDrv.connectErr.Store(errNoError)
	client := liblscp.NewClientWithDriver(lscpDrv)

	ms := &mockSystemdManager{}
	ms.ensureErr.Store(errors.New("systemd fail"))
	ms.failActive.Store(true)

	s := &LinuxSampler{
		Client:  client,
		Systemd: ms,
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s.StartHealthCheck(ctx)
	time.Sleep(2500 * time.Millisecond)
	s.StopHealthCheck()
	if lscpDrv.connected.Load() {
		t.Error("Client should not be reconnected if EnsureLinuxSamplerRunning fails")
	}
}
