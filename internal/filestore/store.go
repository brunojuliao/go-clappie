package filestore

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Entry represents a file in the directory-as-database.
type Entry struct {
	Name    string    // filename without extension
	Path    string    // full path
	ModTime time.Time // last modification time
}

// List returns all .txt files in a directory, sorted by modification time (newest first).
func List(dir string) ([]Entry, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("list %s: %w", dir, err)
	}

	var result []Entry
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".txt") {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		result = append(result, Entry{
			Name:    strings.TrimSuffix(e.Name(), ".txt"),
			Path:    filepath.Join(dir, e.Name()),
			ModTime: info.ModTime(),
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].ModTime.After(result[j].ModTime)
	})

	return result, nil
}

// ReadFile reads a text file and returns its content.
func ReadFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read %s: %w", path, err)
	}
	return string(data), nil
}

// WriteFile writes content to a text file, creating directories as needed.
func WriteFile(path string, content string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("mkdir for %s: %w", path, err)
	}
	return os.WriteFile(path, []byte(content), 0644)
}

// DeleteFile removes a file.
func DeleteFile(path string) error {
	err := os.Remove(path)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("delete %s: %w", path, err)
	}
	return nil
}

// Exists checks if a file exists.
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// FilePath returns the path for a named entry in a directory.
func FilePath(dir, name string) string {
	if strings.HasSuffix(name, ".txt") {
		return filepath.Join(dir, name)
	}
	return filepath.Join(dir, name+".txt")
}

// Count returns the number of .txt files in a directory.
func Count(dir string) int {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0
	}
	count := 0
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".txt") {
			count++
		}
	}
	return count
}

// ReadAndParse reads a file and parses it into body and meta blocks.
func ReadAndParse(path string) (body string, blocks []MetaBlock, err error) {
	content, err := ReadFile(path)
	if err != nil {
		return "", nil, err
	}
	body, blocks = ParseFile(content)
	return body, blocks, nil
}

// WriteWithMeta writes a file with body and meta blocks.
func WriteWithMeta(path string, body string, blocks []MetaBlock) error {
	content := FormatFile(body, blocks)
	return WriteFile(path, content)
}
