package tmux

import "fmt"

// SetPaneStyle sets the style of a tmux pane (background color, etc).
func SetPaneStyle(paneID string, style string) error {
	args := []string{"select-pane"}
	if paneID != "" {
		args = append(args, "-t", paneID)
	}
	args = append(args, "-P", style)
	return RunSilent(args...)
}

// SetPaneBG sets the background color of a pane using RGB values.
func SetPaneBG(paneID string, r, g, b int) error {
	style := fmt.Sprintf("bg=#%02x%02x%02x", r, g, b)
	return SetPaneStyle(paneID, style)
}

// SetSessionStyle sets the status bar style for a session.
func SetSessionStyle(session string, option string, value string) error {
	args := []string{"set-option"}
	if session != "" {
		args = append(args, "-t", session)
	}
	args = append(args, option, value)
	return RunSilent(args...)
}
