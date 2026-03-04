package platform

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ProjectRoot finds the go-clappie project root by walking up from cwd
// looking for a directory containing recall/, chores/, or notifications/.
func ProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if isProjectRoot(dir) {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", fmt.Errorf("go-clappie project root not found (no recall/, chores/, or notifications/ directory found)")
}

func isProjectRoot(dir string) bool {
	markers := []string{"recall", "chores", "notifications"}
	for _, m := range markers {
		info, err := os.Stat(filepath.Join(dir, m))
		if err == nil && info.IsDir() {
			return true
		}
	}
	return false
}

// SocketPath returns the socket path for the current tmux pane.
func SocketPath() string {
	paneID := os.Getenv("TMUX_PANE")
	if paneID == "" {
		paneID = "default"
	}
	// Sanitize pane ID for use in filename
	paneID = strings.ReplaceAll(paneID, "%", "")
	return filepath.Join(os.TempDir(), fmt.Sprintf("go-clappie-%s.sock", paneID))
}

// EnsureDir creates a directory and all parents if they don't exist.
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

// ChoresDir returns the path to the chores directory.
func ChoresDir(root string) string {
	return filepath.Join(root, "chores")
}

// ChoresHumansDir returns the path to chores/humans/.
func ChoresHumansDir(root string) string {
	return filepath.Join(root, "chores", "humans")
}

// ChoresBotsDir returns the path to chores/bots/.
func ChoresBotsDir(root string) string {
	return filepath.Join(root, "chores", "bots")
}

// NotificationsDirtyDir returns the path to notifications/dirty/.
func NotificationsDirtyDir(root string) string {
	return filepath.Join(root, "notifications", "dirty")
}

// NotificationsCleanDir returns the path to notifications/clean/.
func NotificationsCleanDir(root string) string {
	return filepath.Join(root, "notifications", "clean")
}

// RecallDir returns the path to the recall directory.
func RecallDir(root string) string {
	return filepath.Join(root, "recall")
}

// RecallMemoryDir returns the path to recall/memory/.
func RecallMemoryDir(root string) string {
	return filepath.Join(root, "recall", "memory")
}

// RecallLogsDir returns the path to recall/logs/.
func RecallLogsDir(root string) string {
	return filepath.Join(root, "recall", "logs")
}

// RecallSettingsDir returns the path to recall/settings/.
func RecallSettingsDir(root string) string {
	return filepath.Join(root, "recall", "settings")
}

// ProjectsDir returns the path to the projects directory.
func ProjectsDir(root string) string {
	return filepath.Join(root, "projects")
}

// SkillsDir returns the path to the skills directory.
func SkillsDir(root string) string {
	return filepath.Join(root, ".claude", "skills")
}
