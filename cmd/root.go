package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	// Global flags
	dataFlags   []string
	focusFlag   bool
	timeoutFlag int
)

// rootCmd represents the base command.
var rootCmd = &cobra.Command{
	Use:   "clappie",
	Short: "Clappie — personal assistant framework",
	Long:  "Clappie is a personal assistant framework that runs inside tmux, communicating with Claude Code via [clappie] messages.",
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringArrayVarP(&dataFlags, "data", "d", nil, "Data to pass (key=value, key=@file, or JSON)")
	rootCmd.PersistentFlags().BoolVarP(&focusFlag, "focus", "f", false, "Focus the display pane after push")
	rootCmd.PersistentFlags().IntVarP(&timeoutFlag, "timeout", "t", 0, "Timeout in milliseconds")
}

// parseData parses the -d flags into a JSON-encodable map.
func parseData() (json.RawMessage, error) {
	if len(dataFlags) == 0 {
		return nil, nil
	}

	// Check if the first arg is raw JSON
	if len(dataFlags) == 1 && strings.HasPrefix(strings.TrimSpace(dataFlags[0]), "{") {
		raw := json.RawMessage(dataFlags[0])
		if json.Valid(raw) {
			return raw, nil
		}
	}

	result := make(map[string]interface{})
	for _, d := range dataFlags {
		// key=@file — read file content as value
		if idx := strings.Index(d, "=@"); idx > 0 {
			key := d[:idx]
			filePath := d[idx+2:]
			content, err := os.ReadFile(filePath)
			if err != nil {
				return nil, fmt.Errorf("read data file %s: %w", filePath, err)
			}
			result[key] = strings.TrimSpace(string(content))
			continue
		}

		// key=value
		if idx := strings.Index(d, "="); idx > 0 {
			key := d[:idx]
			value := d[idx+1:]
			result[key] = value
			continue
		}

		return nil, fmt.Errorf("invalid data flag: %q (expected key=value or key=@file)", d)
	}

	data, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(data), nil
}
