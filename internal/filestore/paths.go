package filestore

import (
	"path/filepath"
	"time"
)

// LogPath returns the path for a log file in a specific log category.
func LogPath(root, category string) string {
	return filepath.Join(root, "recall", "logs", category)
}

// ChoreLogPath returns the path for a chore log file.
func ChoreLogPath(root string) string {
	return LogPath(root, "chores")
}

// HeartbeatLogPath returns the path for a heartbeat log file.
func HeartbeatLogPath(root string) string {
	return LogPath(root, "heartbeat")
}

// SidekickLogPath returns the path for a sidekick log file.
func SidekickLogPath(root string) string {
	return LogPath(root, "sidekicks")
}

// NotificationLogPath returns the path for a notification log file.
func NotificationLogPath(root string) string {
	return LogPath(root, "notifications")
}

// TimestampedName returns a filename with timestamp prefix.
func TimestampedName(name string) string {
	ts := time.Now().Format("2006-01-02-150405")
	return ts + "-" + name
}

// SettingsPath returns the path for a settings file.
func SettingsPath(root, name string) string {
	return filepath.Join(root, "recall", "settings", name)
}
