package sidekicks

import (
	"fmt"

	"github.com/brunojuliao/go-clappie/internal/tmux"
)

// Spawn creates and starts a new sidekick agent.
func Spawn(root string, config SpawnConfig) (string, error) {
	// Create tmux pane for the sidekick
	paneID, err := spawnTmuxPane(root, config)
	if err != nil {
		return "", fmt.Errorf("spawn pane: %w", err)
	}

	// Create session record
	id, err := createSession(root, config, paneID)
	if err != nil {
		// Clean up pane
		tmux.KillPane(paneID)
		return "", fmt.Errorf("create session: %w", err)
	}

	return id, nil
}

// Send sends a message to the active sidekick.
func Send(root, message string) error {
	active, err := ListActive(root)
	if err != nil {
		return err
	}
	if len(active) == 0 {
		return fmt.Errorf("no active sidekicks")
	}

	sk := active[len(active)-1] // Most recent
	return tmux.SendKeysLiteral(sk.PaneID, message+"\n")
}

// Complete completes the active sidekick with a summary.
func Complete(root, summary string) error {
	active, err := ListActive(root)
	if err != nil {
		return err
	}
	if len(active) == 0 {
		return fmt.Errorf("no active sidekicks")
	}

	sk := active[len(active)-1]

	// Send completion message
	tmux.SendKeysLiteral(sk.PaneID, fmt.Sprintf("[clappie] Sidekick complete → %s\n", summary))

	// Update status
	return updateStatus(root, sk.ID, "completed")
}

// Report sends a report back to the main Claude pane.
func Report(root, message string) error {
	tmux.SubmitToClaudePane("", fmt.Sprintf("[clappie] Sidekick report → %s", message))
	return nil
}

// End ends the active sidekick session.
func End(root string) error {
	active, err := ListActive(root)
	if err != nil {
		return err
	}
	if len(active) == 0 {
		return fmt.Errorf("no active sidekicks")
	}

	sk := active[len(active)-1]

	// Kill the pane
	if sk.PaneID != "" && tmux.PaneExists(sk.PaneID) {
		tmux.KillPane(sk.PaneID)
	}

	return updateStatus(root, sk.ID, "ended")
}

// Message sends a message to a specific sidekick by ID.
func Message(root, id, message string) error {
	sk, err := Get(root, id)
	if err != nil {
		return err
	}
	if sk.Status != "active" {
		return fmt.Errorf("sidekick %s is not active", id)
	}
	return tmux.SendKeysLiteral(sk.PaneID, message+"\n")
}

// Broadcast sends a message to all active sidekicks.
func Broadcast(root, message string) error {
	active, err := ListActive(root)
	if err != nil {
		return err
	}
	for _, sk := range active {
		tmux.SendKeysLiteral(sk.PaneID, message+"\n")
	}
	return nil
}

// KillAll kills all sidekick sessions and panes.
func KillAll() error {
	sessions, _ := tmux.ListSessions()
	for _, s := range sessions {
		if len(s) > 3 && s[:3] == "sk-" {
			tmux.KillSession(s)
		}
	}
	return nil
}
