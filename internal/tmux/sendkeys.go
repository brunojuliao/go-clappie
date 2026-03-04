package tmux

import (
	"os"
	"strings"
)

// SendKeys sends text to a tmux pane via send-keys.
func SendKeys(paneID string, text string, enter bool) error {
	args := []string{"send-keys", "-t", paneID}

	// Escape special tmux characters
	escaped := escapeSendKeys(text)
	args = append(args, escaped)

	if enter {
		args = append(args, "Enter")
	}
	return RunSilent(args...)
}

// SendKeysLiteral sends text literally using load-buffer + paste-buffer.
// This avoids tmux's send-keys escaping issues for long/complex text.
func SendKeysLiteral(paneID string, text string) error {
	// Write to temp file, load into tmux buffer, paste
	f, err := os.CreateTemp("", "clappie-sendkeys-*")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())

	if _, err := f.WriteString(text); err != nil {
		f.Close()
		return err
	}
	f.Close()

	if err := RunSilent("load-buffer", f.Name()); err != nil {
		return err
	}
	return RunSilent("paste-buffer", "-t", paneID, "-d")
}

// escapeSendKeys escapes special tmux send-keys characters.
func escapeSendKeys(s string) string {
	// Escape semicolons which tmux interprets as command separator
	s = strings.ReplaceAll(s, ";", "\\;")
	// Escape dollar signs
	s = strings.ReplaceAll(s, "$", "\\$")
	// Escape tilde which tmux interprets specially
	s = strings.ReplaceAll(s, "~", "\\~")
	return s
}

// SubmitToClaudePane sends a [clappie] message to Claude's pane and presses Enter.
func SubmitToClaudePane(claudePaneID string, message string) error {
	if claudePaneID == "" {
		claudePaneID = os.Getenv("CLAPPIE_CLAUDE_PANE")
	}
	if claudePaneID == "" {
		return nil // No claude pane configured
	}
	return SendKeysLiteral(claudePaneID, message+"\n")
}

// SendToClaudePane sends a [clappie] message to Claude's pane without pressing Enter.
func SendToClaudePane(claudePaneID string, message string) error {
	if claudePaneID == "" {
		claudePaneID = os.Getenv("CLAPPIE_CLAUDE_PANE")
	}
	if claudePaneID == "" {
		return nil
	}
	return SendKeysLiteral(claudePaneID, message)
}
