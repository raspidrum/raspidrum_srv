package linuxsampler

import (
	"bufio"
	"fmt"
	"net"
	"sync"
)

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
