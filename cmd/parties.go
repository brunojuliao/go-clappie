package cmd

import (
	"fmt"

	"github.com/brunojuliao/go-clappie/internal/parties"
	"github.com/brunojuliao/go-clappie/internal/platform"
	"github.com/spf13/cobra"
)

var partiesCmd = &cobra.Command{
	Use:   "parties",
	Short: "Manage party simulations",
}

var partiesInitCmd = &cobra.Command{
	Use:   "init [game]",
	Short: "Initialize a new simulation from a game definition",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		root, err := platform.ProjectRoot()
		if err != nil {
			return err
		}
		simID, err := parties.Init(root, args[0])
		if err != nil {
			return err
		}
		fmt.Printf("Initialized simulation: %s\n", simID)
		return nil
	},
}

var partiesLaunchCmd = &cobra.Command{
	Use:   "launch [sim]",
	Short: "Launch a simulation (spawn AI agents)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		root, err := platform.ProjectRoot()
		if err != nil {
			return err
		}
		return parties.Launch(root, args[0])
	},
}

var partiesShowCmd = &cobra.Command{
	Use:   "show [sim]",
	Short: "Show simulation status/ledger",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		root, err := platform.ProjectRoot()
		if err != nil {
			return err
		}
		return parties.Show(root, args[0])
	},
}

var partiesEndCmd = &cobra.Command{
	Use:   "end [sim]",
	Short: "End a simulation",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		root, err := platform.ProjectRoot()
		if err != nil {
			return err
		}
		return parties.End(root, args[0])
	},
}

var partiesSetCmd = &cobra.Command{
	Use:   "set [sim] [key] [value]",
	Short: "Set a value in the simulation ledger",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		root, err := platform.ProjectRoot()
		if err != nil {
			return err
		}
		return parties.Set(root, args[0], args[1], args[2])
	},
}

var partiesGetCmd = &cobra.Command{
	Use:   "get [sim] [key]",
	Short: "Get a value from the simulation ledger",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		root, err := platform.ProjectRoot()
		if err != nil {
			return err
		}
		value, err := parties.Get(root, args[0], args[1])
		if err != nil {
			return err
		}
		fmt.Println(value)
		return nil
	},
}

var partiesRollCmd = &cobra.Command{
	Use:   "roll [spec]",
	Short: "Roll dice (e.g., 2d6+3, coin, pick item1 item2)",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := parties.Roll(args)
		if err != nil {
			return err
		}
		fmt.Println(result)
		return nil
	},
}

var partiesIdentityCmd = &cobra.Command{
	Use:   "identity [subcommand]",
	Short: "Manage identities",
}

var partiesGamesCmd = &cobra.Command{
	Use:   "games",
	Short: "List available games",
	RunE: func(cmd *cobra.Command, args []string) error {
		root, err := platform.ProjectRoot()
		if err != nil {
			return err
		}
		games, err := parties.ListGames(root)
		if err != nil {
			return err
		}
		if len(games) == 0 {
			fmt.Println("No games found.")
			return nil
		}
		for _, g := range games {
			fmt.Printf("  %s: %s\n", g.Name, g.Description)
		}
		return nil
	},
}

var partiesRulesCmd = &cobra.Command{
	Use:   "rules [game]",
	Short: "Show game rules",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		root, err := platform.ProjectRoot()
		if err != nil {
			return err
		}
		rules, err := parties.Rules(root, args[0])
		if err != nil {
			return err
		}
		fmt.Println(rules)
		return nil
	},
}

func init() {
	partiesCmd.AddCommand(partiesInitCmd)
	partiesCmd.AddCommand(partiesLaunchCmd)
	partiesCmd.AddCommand(partiesShowCmd)
	partiesCmd.AddCommand(partiesEndCmd)
	partiesCmd.AddCommand(partiesSetCmd)
	partiesCmd.AddCommand(partiesGetCmd)
	partiesCmd.AddCommand(partiesRollCmd)
	partiesCmd.AddCommand(partiesIdentityCmd)
	partiesCmd.AddCommand(partiesGamesCmd)
	partiesCmd.AddCommand(partiesRulesCmd)
	rootCmd.AddCommand(partiesCmd)
}
