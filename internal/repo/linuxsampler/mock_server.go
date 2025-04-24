package linuxsampler

import (
	"bufio"
	"net"
	"sync"
)

// MockPipeServer управляет in-memory сервером
type MockPipeServer struct {
	conn             net.Conn
	receivedMessages []string
	mu               sync.Mutex
	done             chan struct{}
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
	defer m.conn.Close()

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
			m.receivedMessages = append(m.receivedMessages, scanner.Text())
			m.mu.Unlock()
		}
	}
}

func (m *MockPipeServer) stop() {
	close(m.done)
}

func (m *MockPipeServer) getMessages() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.receivedMessages
}
