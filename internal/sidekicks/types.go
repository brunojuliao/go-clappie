package sidekicks

import "time"

// SidekickInfo holds information about a sidekick session.
type SidekickInfo struct {
	ID        string
	Prompt    string
	Model     string
	Squad     string
	Status    string
	PaneID    string
	CreatedAt time.Time
}

// SidekickMeta holds metadata for a sidekick session file.
type SidekickMeta struct {
	ID        string
	Prompt    string
	Model     string
	Squad     string
	Status    string
	PaneID    string
	CreatedAt string
}

// SpawnConfig holds configuration for spawning a sidekick.
type SpawnConfig struct {
	Prompt string
	Model  string
	Squad  string
	Skill  string
}
