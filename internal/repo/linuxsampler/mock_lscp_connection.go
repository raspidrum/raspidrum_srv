package linuxsampler

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/raspidrum-srv/libs/liblscp-go"
)

var errNoError = errors.New("no error")

// mockLscpDriver implements LscpDriver interface
// (GetServerInfo, Connect, Host, Port, Timeout)
type mockLscpDriver struct {
	pingErr    atomic.Value // error
	connectErr atomic.Value // error
	connected  atomic.Bool
	conn       net.Conn
}

func (m *mockLscpDriver) Ping() error {
	err, _ := m.pingErr.Load().(error)
	if err == errNoError {
		return nil
	}
	return err
}

func (m *mockLscpDriver) Connect() error {
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

func (m *mockLscpDriver) Disconnect() error {
	return nil
}

func (m *mockLscpDriver) RetrieveInfo(lscpCmd string, isMultiResult bool) (liblscp.ResultSet, error) {
	if m.conn == nil {
		return liblscp.ResultSet{}, fmt.Errorf("not connected")
	}
	cmd := strings.Trim(lscpCmd, " ")
	_, err := fmt.Fprintf(m.conn, "%s\r\n", cmd)
	if err != nil {
		return liblscp.ResultSet{}, fmt.Errorf("failed lscp command: %s : %w", lscpCmd, err)
	}
	return liblscp.GetResultSetFromConn(m.conn, isMultiResult)
}

// MockPipeServer управляет in-memory сервером
type MockPipeServer struct {
	conn             net.Conn
	receivedMessages []string
	mu               sync.Mutex
	done             chan struct{}
	channelIdx       int
}

func startMockPipeServer(conn net.Conn) *MockPipeServer {
	server := &MockPipeServer{
		conn: conn,
		done: make(chan struct{}),
	}

	go server.listen()
	return server
}

func (m *MockPipeServer) listen() {
	//defer m.conn.Close()

	scanner := bufio.NewScanner(m.conn)
	for {
		select {
		case <-m.done:
			return
		default:
			if !scanner.Scan() {
				return
			}

			m.mu.Lock()
			req := scanner.Text()
			m.mu.Unlock()
			m.receivedMessages = append(m.receivedMessages, req)
			var resp string
			switch {
			case req == "ADD CHANNEL":
				resp = fmt.Sprintf("OK[%d]\n\r", m.channelIdx)
				m.channelIdx++
			default:
				// SET CHANNEL ...
				// LOAD LOAD ENGINE ...
				// LOAD INSTRUMENT ...
				resp = "OK\n\r"
			}
			m.conn.Write([]byte(resp))
		}
	}
}

func (m *MockPipeServer) stop() {
	if !IsClosed(m.done) {
		close(m.done)
	}
	m.channelIdx = 0
}

func (m *MockPipeServer) getMessages() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.receivedMessages
}

func IsClosed(ch <-chan struct{}) bool {
	select {
	case <-ch:
		return true
	default:
	}
	return false
}
