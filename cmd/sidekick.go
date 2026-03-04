package cmd

import (
	"fmt"

	"github.com/brunojuliao/go-clappie/internal/platform"
	"github.com/brunojuliao/go-clappie/internal/sidekicks"
	"github.com/spf13/cobra"
)

var (
	sidekickModel string
	sidekickSquad string
)

var sidekickCmd = &cobra.Command{
	Use:   "sidekick",
	Short: "Manage sidekick agents",
}

var sidekickSpawnCmd = &cobra.Command{
	Use:   "spawn [prompt]",
	Short: "Spawn a new sidekick agent",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		root, err := platform.ProjectRoot()
		if err != nil {
			return err
		}
		id, err := sidekicks.Spawn(root, sidekicks.SpawnConfig{
			Prompt: args[0],
			Model:  sidekickModel,
			Squad:  sidekickSquad,
		})
		if err != nil {
			return err
		}
		fmt.Printf("Spawned sidekick: %s\n", id)
		return nil
	},
}

var sidekickSendCmd = &cobra.Command{
	Use:   "send [message]",
	Short: "Send a message to the active sidekick",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		root, err := platform.ProjectRoot()
		if err != nil {
			return err
		}
		return sidekicks.Send(root, args[0])
	},
}

var sidekickCompleteCmd = &cobra.Command{
	Use:   "complete [summary]",
	Short: "Complete the active sidekick",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		root, err := platform.ProjectRoot()
		if err != nil {
			return err
		}
		return sidekicks.Complete(root, args[0])
	},
}

var sidekickReportCmd = &cobra.Command{
	Use:   "report [message]",
	Short: "Report back to the main Claude pane",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		root, err := platform.ProjectRoot()
		if err != nil {
			return err
		}
		return sidekicks.Report(root, args[0])
	},
}

var sidekickEndCmd = &cobra.Command{
	Use:   "end",
	Short: "End the active sidekick session",
	RunE: func(cmd *cobra.Command, args []string) error {
		root, err := platform.ProjectRoot()
		if err != nil {
			return err
		}
		return sidekicks.End(root)
	},
}

var sidekickKillCmd = &cobra.Command{
	Use:   "kill",
	Short: "Kill all sidekick sessions",
	RunE: func(cmd *cobra.Command, args []string) error {
		return sidekicks.KillAll()
	},
}

var sidekickMessageCmd = &cobra.Command{
	Use:   "message [id] [message]",
	Short: "Send a message to a specific sidekick",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		root, err := platform.ProjectRoot()
		if err != nil {
			return err
		}
		return sidekicks.Message(root, args[0], args[1])
	},
}

var sidekickBroadcastCmd = &cobra.Command{
	Use:   "broadcast [message]",
	Short: "Broadcast a message to all sidekicks",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		root, err := platform.ProjectRoot()
		if err != nil {
			return err
		}
		return sidekicks.Broadcast(root, args[0])
	},
}

var sidekickListCmd = &cobra.Command{
	Use:   "list",
	Short: "List active sidekicks",
	RunE: func(cmd *cobra.Command, args []string) error {
		root, err := platform.ProjectRoot()
		if err != nil {
			return err
		}
		sks, err := sidekicks.ListActive(root)
		if err != nil {
			return err
		}
		if len(sks) == 0 {
			fmt.Println("No active sidekicks.")
			return nil
		}
		for _, sk := range sks {
			fmt.Printf("  %s: %s\n", sk.ID, sk.Prompt)
		}
		return nil
	},
}

func init() {
	sidekickSpawnCmd.Flags().StringVar(&sidekickModel, "model", "", "Model to use")
	sidekickSpawnCmd.Flags().StringVar(&sidekickSquad, "squad", "", "Squad/group name")

	sidekickCmd.AddCommand(sidekickSpawnCmd)
	sidekickCmd.AddCommand(sidekickSendCmd)
	sidekickCmd.AddCommand(sidekickCompleteCmd)
	sidekickCmd.AddCommand(sidekickReportCmd)
	sidekickCmd.AddCommand(sidekickEndCmd)
	sidekickCmd.AddCommand(sidekickKillCmd)
	sidekickCmd.AddCommand(sidekickMessageCmd)
	sidekickCmd.AddCommand(sidekickBroadcastCmd)
	sidekickCmd.AddCommand(sidekickListCmd)
	rootCmd.AddCommand(sidekickCmd)
}
