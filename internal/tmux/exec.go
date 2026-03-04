package tmux

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

const defaultTimeout = 5 * time.Second

// Run executes a tmux command with default timeout and returns trimmed output.
func Run(args ...string) (string, error) {
	return RunWithTimeout(defaultTimeout, args...)
}

// RunWithTimeout executes a tmux command with a specified timeout.
func RunWithTimeout(timeout time.Duration, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "tmux", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("tmux %s: timeout after %v", strings.Join(args, " "), timeout)
		}
		return "", fmt.Errorf("tmux %s: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(string(out)))
	}
	return strings.TrimSpace(string(out)), nil
}

// RunSilent executes a tmux command and ignores output.
func RunSilent(args ...string) error {
	_, err := Run(args...)
	return err
}
