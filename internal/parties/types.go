package parties

import "time"

// GameInfo holds basic game information.
type GameInfo struct {
	Name        string
	Description string
	Path        string
}

// Game represents a parsed game definition.
type Game struct {
	Name           string
	Description    string
	Rules          string
	PlayerCards    []PlayerCard
	SuggestedState map[string]string
	Path           string
}

// PlayerCard defines a role/character in the game.
type PlayerCard struct {
	Name        string
	Description string
	Abilities   []string
}

// Ledger holds the state of a simulation.
type Ledger struct {
	SimID      string
	GameName   string
	State      map[string]string
	Participants []Participant
	CreatedAt  time.Time
	Path       string
}

// Participant represents an AI agent in a simulation.
type Participant struct {
	Name     string
	Role     string
	PaneID   string
	Identity string
}

// Identity represents a persistent character identity.
type Identity struct {
	Name     string
	Template string
	Memories []string
}
