package tmux

import (
	"strings"
)

// SessionExists checks if a tmux session exists.
func SessionExists(name string) bool {
	_, err := Run("has-session", "-t", name)
	return err == nil
}

// ListSessions returns all tmux session names.
func ListSessions() ([]string, error) {
	out, err := Run("list-sessions", "-F", "#{session_name}")
	if err != nil {
		return nil, err
	}
	if out == "" {
		return nil, nil
	}
	return strings.Split(out, "\n"), nil
}

// NewSession creates a new tmux session.
func NewSession(name string, cmd string, args ...string) error {
	tmuxArgs := []string{"new-session", "-d", "-s", name}
	if cmd != "" {
		fullCmd := cmd
		for _, a := range args {
			fullCmd += " " + a
		}
		tmuxArgs = append(tmuxArgs, fullCmd)
	}
	return RunSilent(tmuxArgs...)
}

// KillSession kills a tmux session.
func KillSession(name string) error {
	return RunSilent("kill-session", "-t", name)
}
