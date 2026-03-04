package ipc

import (
	"path/filepath"
	"testing"
)

func TestIPCRoundtrip(t *testing.T) {
	socketPath := filepath.Join(t.TempDir(), "test.sock")

	handler := func(cmd Command) Response {
		switch cmd.Action {
		case ActionPing:
			return Response{OK: true, Message: "pong"}
		default:
			return Response{OK: false, Error: "unknown"}
		}
	}

	server, err := NewServer(socketPath, handler)
	if err != nil {
		t.Fatalf("NewServer: %v", err)
	}

	go server.Serve()
	defer server.Close()

	// Test ping
	resp, err := SendCommand(socketPath, Command{Action: ActionPing})
	if err != nil {
		t.Fatalf("SendCommand: %v", err)
	}
	if !resp.OK {
		t.Error("expected OK response")
	}
	if resp.Message != "pong" {
		t.Errorf("message = %q, want pong", resp.Message)
	}

	// Test Ping helper
	if !Ping(socketPath) {
		t.Error("Ping should return true")
	}
}

func TestIPCUnknownAction(t *testing.T) {
	socketPath := filepath.Join(t.TempDir(), "test2.sock")

	handler := func(cmd Command) Response {
		return Response{OK: false, Error: "unknown action"}
	}

	server, err := NewServer(socketPath, handler)
	if err != nil {
		t.Fatalf("NewServer: %v", err)
	}

	go server.Serve()
	defer server.Close()

	resp, err := SendCommand(socketPath, Command{Action: "nonexistent"})
	if err != nil {
		t.Fatalf("SendCommand: %v", err)
	}
	if resp.OK {
		t.Error("expected not OK for unknown action")
	}
}

func TestPingNotRunning(t *testing.T) {
	socketPath := filepath.Join(t.TempDir(), "nonexistent.sock")
	if Ping(socketPath) {
		t.Error("Ping should return false for non-running server")
	}
}
