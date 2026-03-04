package cmd

import (
	"fmt"
	"os"

	"github.com/brunojuliao/go-clappie/internal/scaffold"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize go-clappie in the current project",
	Long:  "Scaffolds .claude/skills/go-clappie/SKILL.md, CLAUDE.md, and data directories into the current directory.",
	RunE:  runInit,
}

func runInit(cmd *cobra.Command, args []string) error {
	cwd, err := cmd.Flags().GetString("root")
	if err != nil || cwd == "" {
		// Default to current working directory
		var osErr error
		cwd, osErr = os.Getwd()
		if osErr != nil {
			return fmt.Errorf("get working directory: %w", osErr)
		}
	}

	result := scaffold.Run(cwd)

	// Print results
	for _, path := range result.Created {
		fmt.Printf("  Created: %s\n", path)
	}
	for _, path := range result.Skipped {
		fmt.Printf("  Skipped: %s\n", path)
	}
	for _, err := range result.Errors {
		fmt.Printf("  Error:   %v\n", err)
	}

	fmt.Println()
	if len(result.Errors) > 0 {
		return fmt.Errorf("%d errors during init", len(result.Errors))
	}

	fmt.Printf("Initialized go-clappie (%d created, %d skipped)\n", len(result.Created), len(result.Skipped))
	fmt.Println()
	fmt.Println("Next step:")
	fmt.Println("  go-clappie display push heartbeat")

	return nil
}

func init() {
	initCmd.Flags().String("root", "", "Project root directory (default: current directory)")
	rootCmd.AddCommand(initCmd)
}
