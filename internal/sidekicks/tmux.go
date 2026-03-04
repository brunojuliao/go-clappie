package sidekicks

import (
	"fmt"
	"os"

	ttmux "github.com/brunojuliao/go-clappie/internal/tmux"
)

// spawnTmuxPane creates a new tmux pane for a sidekick.
func spawnTmuxPane(root string, config SpawnConfig) (string, error) {
	binary, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("get executable: %w", err)
	}

	// Build the command for the sidekick pane
	// The sidekick runs claude with the prompt
	cmd := fmt.Sprintf("claude --dangerously-skip-permissions '%s'", escapePrompt(config.Prompt))
	if config.Model != "" {
		cmd = fmt.Sprintf("claude --dangerously-skip-permissions --model '%s' '%s'", config.Model, escapePrompt(config.Prompt))
	}

	_ = binary // Not spawning self for sidekick, spawning claude directly

	// Split pane horizontally
	paneID, err := ttmux.SplitPane("", ttmux.SplitHorizontal, 50, cmd)
	if err != nil {
		return "", err
	}

	return paneID, nil
}

func escapePrompt(s string) string {
	// Escape single quotes for shell
	result := ""
	for _, c := range s {
		if c == '\'' {
			result += "'\"'\"'"
		} else {
			result += string(c)
		}
	}
	return result
}
