package chores

import "time"

// Status constants for chore lifecycle.
const (
	StatusPending   = "pending"
	StatusApproved  = "approved"
	StatusCompleted = "completed"
	StatusRejected  = "rejected"
	StatusShelved   = "shelved"
)

// Chore represents a human approval task.
type Chore struct {
	Name    string
	Path    string
	Title   string
	Summary string
	Icon    string
	Context string
	Status  string
	Created time.Time
	Body    string
}
