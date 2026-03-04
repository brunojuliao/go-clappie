package background

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/brunojuliao/go-clappie/internal/tmux"
)

const sessionPrefix = "go-clappie-bg-"

// Start starts a background app in a new tmux session.
func Start(root, name string) error {
	apps, err := Discover(root)
	if err != nil {
		return err
	}

	var app *App
	for i := range apps {
		if apps[i].Name == name {
			app = &apps[i]
			break
		}
	}
	if app == nil {
		return fmt.Errorf("background app %q not found", name)
	}

	sessionName := fmt.Sprintf("%s%s-%d", sessionPrefix, name, time.Now().Unix())

	binary, err := os.Executable()
	if err != nil {
		return fmt.Errorf("get executable: %w", err)
	}

	// Start the app in a new tmux session
	return tmux.NewSession(sessionName, binary, "__daemon")
}

// Stop stops a background app by killing its tmux session.
func Stop(name string) error {
	sessions, err := tmux.ListSessions()
	if err != nil {
		return err
	}

	prefix := sessionPrefix + name
	for _, s := range sessions {
		if strings.HasPrefix(s, prefix) {
			return tmux.KillSession(s)
		}
	}
	return fmt.Errorf("background app %q is not running", name)
}

// List returns all background apps with their running status.
func List(root string) ([]App, error) {
	apps, err := Discover(root)
	if err != nil {
		return nil, err
	}

	sessions, _ := tmux.ListSessions()
	sessionSet := make(map[string]bool)
	for _, s := range sessions {
		sessionSet[s] = true
	}

	for i := range apps {
		prefix := sessionPrefix + apps[i].Name
		for s := range sessionSet {
			if strings.HasPrefix(s, prefix) {
				apps[i].Running = true
				apps[i].Session = s
				break
			}
		}
	}

	return apps, nil
}

// KillAll kills all background app sessions.
func KillAll() error {
	sessions, err := tmux.ListSessions()
	if err != nil {
		return err
	}

	for _, s := range sessions {
		if strings.HasPrefix(s, sessionPrefix) {
			tmux.KillSession(s)
		}
	}
	return nil
}
