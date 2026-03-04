package ipc

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/brunojuliao/go-clappie/internal/platform"
)

const (
	dialTimeout = 3 * time.Second
	readTimeout = 5 * time.Second
)

// SendCommand sends a command to the daemon and returns the response.
func SendCommand(socketPath string, cmd Command) (*Response, error) {
	conn, err := platform.Dial(socketPath)
	if err != nil {
		return nil, fmt.Errorf("connect to daemon: %w", err)
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(readTimeout))

	data, err := json.Marshal(cmd)
	if err != nil {
		return nil, fmt.Errorf("marshal command: %w", err)
	}

	// Write length-prefixed JSON
	data = append(data, '\n')
	if _, err := conn.Write(data); err != nil {
		return nil, fmt.Errorf("write command: %w", err)
	}

	// Read response
	buf, err := io.ReadAll(conn)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var resp Response
	if err := json.Unmarshal(buf, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	return &resp, nil
}

// Ping checks if the daemon is running.
func Ping(socketPath string) bool {
	conn, err := platform.Dial(socketPath)
	if err != nil {
		return false
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(2 * time.Second))

	cmd, _ := json.Marshal(Command{Action: ActionPing})
	cmd = append(cmd, '\n')
	if _, err := conn.Write(cmd); err != nil {
		return false
	}

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil && err != io.EOF {
		return false
	}

	var resp Response
	if err := json.Unmarshal(buf[:n], &resp); err != nil {
		return false
	}
	return resp.OK
}

// IsDaemonRunning checks whether a daemon is listening at the given socket.
func IsDaemonRunning(socketPath string) bool {
	conn, err := net.DialTimeout("unix", socketPath, dialTimeout)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}
