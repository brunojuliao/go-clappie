package heartbeat

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/brunojuliao/go-clappie/internal/filestore"
	"github.com/brunojuliao/go-clappie/internal/platform"
)

// DiscoverChecks finds all heartbeat check files.
func DiscoverChecks(root string) ([]Check, error) {
	dir := platform.ChoresBotsDir(root)
	entries, err := filestore.List(dir)
	if err != nil {
		return nil, err
	}

	var checks []Check
	for _, entry := range entries {
		body, blocks, err := filestore.ReadAndParse(entry.Path)
		if err != nil {
			continue
		}

		meta := filestore.GetMeta(blocks, "heartbeat-meta")

		check := Check{
			Name: entry.Name,
			Path: entry.Path,
			Body: body,
		}

		if meta != nil {
			if interval, ok := meta.Fields["interval"]; ok {
				check.Interval = ParseInterval(interval)
			}
			if lastRun, ok := meta.Fields["last_run"]; ok {
				if t, err := time.Parse("2006-01-02 15:04:05", lastRun); err == nil {
					check.LastRun = t
				}
			}
			check.Status = meta.Fields["status"]
		}

		checks = append(checks, check)
	}

	return checks, nil
}

// IsDue returns true if a check is due to run.
func IsDue(check Check) bool {
	if check.Interval == 0 {
		return false
	}
	if check.LastRun.IsZero() {
		return true
	}
	return time.Since(check.LastRun) >= check.Interval
}

// MarkRun updates the last_run timestamp for a check.
func MarkRun(check Check, status string) error {
	body, blocks, err := filestore.ReadAndParse(check.Path)
	if err != nil {
		return err
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	filestore.SetMetaField(&blocks, "heartbeat-meta", "last_run", now)
	filestore.SetMetaField(&blocks, "heartbeat-meta", "status", status)

	return filestore.WriteWithMeta(check.Path, body, blocks)
}

// ParseInterval parses interval strings like "5m", "1h", "30s", "1d".
func ParseInterval(s string) time.Duration {
	s = strings.TrimSpace(strings.ToLower(s))
	if s == "" {
		return 0
	}

	// Try standard Go duration first
	if d, err := time.ParseDuration(s); err == nil {
		return d
	}

	// Handle "d" for days
	if strings.HasSuffix(s, "d") {
		numStr := s[:len(s)-1]
		if n, err := strconv.Atoi(numStr); err == nil {
			return time.Duration(n) * 24 * time.Hour
		}
	}

	return 0
}

// FormatLogEntry formats a heartbeat run result.
func FormatLogEntry(results []CheckResult) string {
	now := time.Now().Format("15:04:05")
	var parts []string
	for _, r := range results {
		icon := "✓"
		if !r.OK {
			icon = "✗"
		}
		part := fmt.Sprintf("%s %s", r.Name, icon)
		if r.Note != "" {
			part += fmt.Sprintf(" (%s)", r.Note)
		}
		parts = append(parts, part)
	}
	return fmt.Sprintf("[%s] %s", now, strings.Join(parts, " | "))
}

// CheckResult holds the outcome of a single heartbeat check.
type CheckResult struct {
	Name string
	OK   bool
	Note string
}
