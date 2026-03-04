package notifications

import "time"

// DirtyItem represents a raw notification from an integration.
type DirtyItem struct {
	Name     string
	Path     string
	Body     string
	Source   string
	SourceID string
	Created  time.Time
}

// CleanItem represents a curated notification.
type CleanItem struct {
	Name     string
	Path     string
	Body     string
	SourceID string
	Context  string
	Created  time.Time
}
