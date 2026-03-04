package filestore

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// FileLock provides file-based locking for concurrent access.
type FileLock struct {
	path string
	file *os.File
	mu   sync.Mutex
}

// NewFileLock creates a new file lock.
func NewFileLock(path string) *FileLock {
	return &FileLock{
		path: path + ".lock",
	}
}

// Lock acquires the file lock, blocking until available.
func (fl *FileLock) Lock() error {
	fl.mu.Lock()

	if err := os.MkdirAll(filepath.Dir(fl.path), 0755); err != nil {
		fl.mu.Unlock()
		return fmt.Errorf("create lock dir: %w", err)
	}

	f, err := os.OpenFile(fl.path, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		fl.mu.Unlock()
		return fmt.Errorf("open lock file: %w", err)
	}
	fl.file = f

	if err := lockFile(f); err != nil {
		f.Close()
		fl.mu.Unlock()
		return fmt.Errorf("lock file: %w", err)
	}
	return nil
}

// TryLock attempts to acquire the lock without blocking.
func (fl *FileLock) TryLock() (bool, error) {
	if !fl.mu.TryLock() {
		return false, nil
	}

	if err := os.MkdirAll(filepath.Dir(fl.path), 0755); err != nil {
		fl.mu.Unlock()
		return false, fmt.Errorf("create lock dir: %w", err)
	}

	f, err := os.OpenFile(fl.path, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		fl.mu.Unlock()
		return false, fmt.Errorf("open lock file: %w", err)
	}
	fl.file = f

	if !tryLockFile(f) {
		f.Close()
		fl.mu.Unlock()
		return false, nil
	}
	return true, nil
}

// Unlock releases the file lock.
func (fl *FileLock) Unlock() error {
	if fl.file != nil {
		unlockFile(fl.file)
		fl.file.Close()
		os.Remove(fl.path)
		fl.file = nil
	}
	fl.mu.Unlock()
	return nil
}

// WithLock runs a function while holding the lock.
func (fl *FileLock) WithLock(fn func() error) error {
	if err := fl.Lock(); err != nil {
		return err
	}
	defer fl.Unlock()
	return fn()
}

// WithTimeout tries to acquire the lock within the timeout duration.
func (fl *FileLock) WithTimeout(timeout time.Duration, fn func() error) error {
	deadline := time.Now().Add(timeout)
	for {
		ok, err := fl.TryLock()
		if err != nil {
			return err
		}
		if ok {
			defer fl.Unlock()
			return fn()
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("lock timeout after %v", timeout)
		}
		time.Sleep(10 * time.Millisecond)
	}
}
