package linuxsampler

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	lscp "github.com/raspidrum-srv/libs/liblscp-go"
)

// mockClient implements LinuxSamplerClient interface
// (GetServerInfo, Connect, Host, Port, Timeout)
type mockClient struct {
	getServerInfoErr atomic.Value // error
	connectErr       atomic.Value // error
	connected        atomic.Bool
}

func (m *mockClient) GetServerInfo() (lscp.ServerInfo, error) {
	err, _ := m.getServerInfoErr.Load().(error)
	if err == errNoError {
		return lscp.ServerInfo{}, nil
	}
	return lscp.ServerInfo{}, err
}
func (m *mockClient) Connect() error {
	err, _ := m.connectErr.Load().(error)
	if err == errNoError {
		m.connected.Store(true)
		return nil
	}
	if err == nil {
		m.connected.Store(true)
	}
	return err
}
//func (m *mockClient) Host() string   { return "localhost" }
//func (m *mockClient) Port() string   { return "1234" }
//func (m *mockClient) Timeout() string { return "1s" }

// mockSystemdManager implements dbus.SystemdManager
// (IsServiceActive, StartService, WaitForServiceActive)
type mockSystemdManager struct {
	ensureErr atomic.Value // error
	called    atomic.Bool
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

var errNoError = errors.New("no error")

func TestHealthCheck_Success(t *testing.T) {
	client := &mockClient{}
	client.getServerInfoErr.Store(errNoError)
	s := &LinuxSampler{
		Client:  lscp.Client{}, // не используется
		Systemd: &mockSystemdManager{},
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s.StartHealthCheck(ctx, client)
	time.Sleep(2500 * time.Millisecond)
	s.StopHealthCheck()
	// Should not call EnsureLinuxSamplerRunning
	// (мы не проверяем вызовы SystemdManager в этом тесте)
}

func TestHealthCheck_ReconnectOnFailure(t *testing.T) {
	client := &mockClient{}
	client.getServerInfoErr.Store(errors.New("fail"))
	client.connectErr.Store(errNoError)
	ms := &mockSystemdManager{}
	ms.ensureErr.Store(errNoError)

	s := &LinuxSampler{
		Client:  client,
		Systemd: ms,
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s.StartHealthCheck(ctx, client)
	time.Sleep(2500 * time.Millisecond)
	s.StopHealthCheck()
	if !client.connected.Load() {
		t.Error("Client should be reconnected on failure")
	}
}

func TestHealthCheck_EnsureFail(t *testing.T) {
	client := &mockClient{}
	client.getServerInfoErr.Store(errors.New("fail"))
	client.connectErr.Store(errNoError)
	ms := &mockSystemdManager{}
	ms.ensureErr.Store(errors.New("systemd fail"))
	ms.failActive.Store(true)

	s := &LinuxSampler{
		Client:  lscp.Client{}, // не используется
		Systemd: ms,
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s.StartHealthCheck(ctx, client)
	time.Sleep(2500 * time.Millisecond)
	s.StopHealthCheck()
	if client.connected.Load() {
		t.Error("Client should not be reconnected if EnsureLinuxSamplerRunning fails")
	}
} 