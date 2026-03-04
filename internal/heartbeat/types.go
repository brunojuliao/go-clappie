package heartbeat

import "time"

// Check represents a heartbeat check configuration.
type Check struct {
	Name     string
	Path     string
	Body     string
	Interval time.Duration
	LastRun  time.Time
	Status   string // "ok", "error", ""
}
