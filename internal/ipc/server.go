package ipc

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/brunojuliao/go-clappie/internal/platform"
)

// Handler processes an IPC command and returns a response.
type Handler func(cmd Command) Response

// Server listens for IPC connections and dispatches commands.
type Server struct {
	listener net.Listener
	addr     string
	handler  Handler
	mu       sync.Mutex
	done     chan struct{}
}

// NewServer creates a new IPC server.
func NewServer(socketPath string, handler Handler) (*Server, error) {
	ln, addr, err := platform.Listen(socketPath)
	if err != nil {
		return nil, fmt.Errorf("listen: %w", err)
	}
	return &Server{
		listener: ln,
		addr:     addr,
		handler:  handler,
		done:     make(chan struct{}),
	}, nil
}

// Addr returns the address the server is listening on.
func (s *Server) Addr() string {
	return s.addr
}

// Serve starts accepting connections. Blocks until Close is called.
func (s *Server) Serve() error {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.done:
				return nil
			default:
				log.Printf("ipc: accept error: %v", err)
				continue
			}
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024) // 1MB max

	if !scanner.Scan() {
		return
	}

	var cmd Command
	if err := json.Unmarshal(scanner.Bytes(), &cmd); err != nil {
		resp := Response{OK: false, Error: fmt.Sprintf("invalid command: %v", err)}
		data, _ := json.Marshal(resp)
		conn.Write(data)
		return
	}

	resp := s.handler(cmd)
	data, err := json.Marshal(resp)
	if err != nil {
		resp = Response{OK: false, Error: fmt.Sprintf("marshal error: %v", err)}
		data, _ = json.Marshal(resp)
	}
	conn.Write(data)
}

// Close shuts down the server.
func (s *Server) Close() error {
	close(s.done)
	return s.listener.Close()
}
