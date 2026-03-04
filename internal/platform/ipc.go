package platform

import (
	"fmt"
	"hash/fnv"
	"net"
	"os"
	"runtime"
)

// Listen creates a net.Listener using Unix sockets on Unix systems
// and TCP localhost on Windows. Returns the listener, address string, and error.
func Listen(socketPath string) (net.Listener, string, error) {
	if runtime.GOOS == "windows" {
		port := derivePort(socketPath)
		addr := fmt.Sprintf("127.0.0.1:%d", port)
		ln, err := net.Listen("tcp", addr)
		return ln, addr, err
	}
	// Remove stale socket file
	os.Remove(socketPath)
	ln, err := net.Listen("unix", socketPath)
	return ln, socketPath, err
}

// Dial connects to a listener at the given socket path.
// Uses Unix sockets on Unix systems and TCP on Windows.
func Dial(socketPath string) (net.Conn, error) {
	if runtime.GOOS == "windows" {
		port := derivePort(socketPath)
		return net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	}
	return net.Dial("unix", socketPath)
}

// derivePort hashes a socket path to a deterministic port number in the
// ephemeral range (49152-65535).
func derivePort(socketPath string) int {
	h := fnv.New32a()
	h.Write([]byte(socketPath))
	return 49152 + int(h.Sum32()%16383)
}
