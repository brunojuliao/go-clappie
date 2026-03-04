package cmd

import (
	"fmt"

	"github.com/brunojuliao/go-clappie/internal/background"
	"github.com/brunojuliao/go-clappie/internal/platform"
	"github.com/spf13/cobra"
)

var backgroundCmd = &cobra.Command{
	Use:   "background",
	Short: "Manage background apps",
}

var backgroundStartCmd = &cobra.Command{
	Use:   "start [app]",
	Short: "Start a background app",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		root, err := platform.ProjectRoot()
		if err != nil {
			return err
		}
		return background.Start(root, args[0])
	},
}

var backgroundStopCmd = &cobra.Command{
	Use:   "stop [app]",
	Short: "Stop a background app",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return background.Stop(args[0])
	},
}

var backgroundListCmd = &cobra.Command{
	Use:   "list",
	Short: "List background apps",
	RunE: func(cmd *cobra.Command, args []string) error {
		root, err := platform.ProjectRoot()
		if err != nil {
			return err
		}
		apps, err := background.List(root)
		if err != nil {
			return err
		}
		if len(apps) == 0 {
			fmt.Println("No background apps found.")
			return nil
		}
		for _, app := range apps {
			status := "stopped"
			if app.Running {
				status = "running"
			}
			fmt.Printf("  %s [%s]\n", app.Name, status)
		}
		return nil
	},
}

var backgroundKillCmd = &cobra.Command{
	Use:   "kill",
	Short: "Kill all background apps",
	RunE: func(cmd *cobra.Command, args []string) error {
		return background.KillAll()
	},
}

func init() {
	backgroundCmd.AddCommand(backgroundStartCmd)
	backgroundCmd.AddCommand(backgroundStopCmd)
	backgroundCmd.AddCommand(backgroundListCmd)
	backgroundCmd.AddCommand(backgroundKillCmd)
	rootCmd.AddCommand(backgroundCmd)
}
