package cmd

import (
	"fmt"

	"github.com/brunojuliao/go-clappie/internal/background"
	"github.com/brunojuliao/go-clappie/internal/ipc"
	"github.com/brunojuliao/go-clappie/internal/platform"
	"github.com/brunojuliao/go-clappie/internal/sidekicks"
	"github.com/spf13/cobra"
)

var killCmd = &cobra.Command{
	Use:   "kill",
	Short: "Kill all go-clappie processes (display, background, sidekicks)",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Kill display daemon
		socketPath := platform.SocketPath()
		if ipc.Ping(socketPath) {
			ipc.SendCommand(socketPath, ipc.Command{Action: ipc.ActionKill})
			fmt.Println("Display daemon killed.")
		}

		// Kill background apps
		if err := background.KillAll(); err != nil {
			fmt.Printf("Warning: killing background apps: %v\n", err)
		} else {
			fmt.Println("Background apps killed.")
		}

		// Kill sidekicks
		if err := sidekicks.KillAll(); err != nil {
			fmt.Printf("Warning: killing sidekicks: %v\n", err)
		} else {
			fmt.Println("Sidekicks killed.")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(killCmd)
}
