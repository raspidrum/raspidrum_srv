//go:build linux

package udev

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"syscall"
	"unsafe"
)

const (
	NETLINK_KOBJECT_UEVENT = 15
	UEVENT_BUFFER_SIZE     = 2048
)

// NetlinkSockaddr represents sockaddr_nl structure.
type NetlinkSockaddr struct {
	Family uint16
	Pad    uint16
	Pid    uint32
	Groups uint32
}

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
type Monitor struct {
	socket    int
	stopCh    chan struct{}
	closeOnce sync.Once
}

// NewMonitor creates a new udev monitor.
func NewMonitor() (*Monitor, error) {
	socket, err := syscall.Socket(syscall.AF_NETLINK, syscall.SOCK_RAW, NETLINK_KOBJECT_UEVENT)
	if err != nil {
		return nil, fmt.Errorf("failed to create netlink socket: %v", err)
	}

	// Set a read timeout on the socket. This is crucial for allowing graceful shutdown.
	// Without a timeout, the syscall.Read call would block indefinitely, preventing
	// the monitoring goroutine from checking the stop channel.
	tv := syscall.Timeval{Sec: 1, Usec: 0} // 1-second timeout
	if err := syscall.SetsockoptTimeval(socket, syscall.SOL_SOCKET, syscall.SO_RCVTIMEO, &tv); err != nil {
		syscall.Close(socket)
		return nil, fmt.Errorf("failed to set socket read timeout: %v", err)
	}

	addr := &NetlinkSockaddr{
		Family: syscall.AF_NETLINK,
		Pid:    uint32(syscall.Getpid()),
		Groups: 1, // Group for kernel uevents.
	}

	if err := syscall.Bind(socket, (*syscall.SockaddrNetlink)(unsafe.Pointer(addr))); err != nil {
		syscall.Close(socket)
		return nil, fmt.Errorf("failed to bind netlink socket: %v", err)
	}

	return &Monitor{
		socket: socket,
		stopCh: make(chan struct{}),
	}, nil
}

// Close closes the monitor.
func (m *Monitor) Close() {
	m.closeOnce.Do(func() {
		close(m.stopCh)
		// Closing the socket will make any blocking Read call return an error.
		syscall.Close(m.socket)
	})
}

func parseUevent(buffer []byte) *Event {
	env := make(map[string]string)
	event := &Event{Env: env}

	lines := strings.Split(string(buffer), "\x00")
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key, value := parts[0], parts[1]
		env[key] = value

		switch key {
		case "ACTION":
			event.Action = value
		case "DEVPATH":
			event.DevPath = value
		case "SUBSYSTEM":
			event.Subsystem = value
		case "DEVNAME":
			event.DevName = value
		case "DEVTYPE":
			event.DevType = value
		}
	}
	return event
}

func (m *Monitor) receiveEvent() (*Event, error) {
	buffer := make([]byte, UEVENT_BUFFER_SIZE)
	n, err := syscall.Read(m.socket, buffer)
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, nil
	}
	return parseUevent(buffer[:n]), nil
}

// Start begins monitoring for udev events.
func (m *Monitor) Start(ctx context.Context) (<-chan *Event, error) {
	eventsCh := make(chan *Event)

	go func() {
		defer close(eventsCh)
		for {
			// This select ensures we check for shutdown signal before blocking on read.
			// However, the read itself can still block, which is why a socket timeout is needed.
			select {
			case <-m.stopCh:
				return
			default:
			}

			event, err := m.receiveEvent()
			if err != nil {
				// EAGAIN is the error code for a timeout. This is expected.
				if err == syscall.EAGAIN || err == syscall.EWOULDBLOCK {
					continue // Timeout occurred, just loop again to check stopCh.
				}

				// For any other error, log it if it wasn't a clean shutdown, then exit.
				select {
				case <-m.stopCh:
					// Error occurred because the socket was closed. This is a clean shutdown.
				default:
					slog.Error("Udev receive error", "error", err)
				}
				return
			}

			if event != nil {
				select {
				case eventsCh <- event:
				case <-m.stopCh:
					return
				}
			}
		}
	}()

	go func() {
		<-ctx.Done()
		m.Close()
	}()

	return eventsCh, nil
}
