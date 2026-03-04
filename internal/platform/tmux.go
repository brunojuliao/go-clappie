package platform

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// TmuxAvailable checks if tmux is installed and accessible.
func TmuxAvailable() bool {
	_, err := exec.LookPath("tmux")
	return err == nil
}

// InTmux returns true if we're running inside a tmux session.
func InTmux() bool {
	return os.Getenv("TMUX") != ""
}

// TmuxPaneID returns the current tmux pane ID from env.
func TmuxPaneID() string {
	return os.Getenv("TMUX_PANE")
}

// TmuxExec runs a tmux command and returns its output.
func TmuxExec(args ...string) (string, error) {
	cmd := exec.Command("tmux", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("tmux %s: %w: %s", strings.Join(args, " "), err, string(out))
	}
	return strings.TrimSpace(string(out)), nil
}
