package cmd

import (
	"fmt"

	"github.com/brunojuliao/go-clappie/internal/background"
	"github.com/brunojuliao/go-clappie/internal/displays"
	"github.com/brunojuliao/go-clappie/internal/platform"
	"github.com/brunojuliao/go-clappie/internal/sidekicks"
	"github.com/brunojuliao/go-clappie/internal/skills"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list [category]",
	Short: "List available resources (skills, displays, background, sidekicks)",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		category := ""
		if len(args) > 0 {
			category = args[0]
		}

		switch category {
		case "displays":
			return listDisplays()
		case "skills", "background", "sidekicks":
			root, err := platform.ProjectRoot()
			if err != nil {
				return err
			}
			switch category {
			case "skills":
				return listSkills(root)
			case "background":
				return listBackground(root)
			case "sidekicks":
				return listSidekicks(root)
			}
		case "":
			root, _ := platform.ProjectRoot()
			return listAll(root)
		default:
			return fmt.Errorf("unknown category: %s (available: skills, displays, background, sidekicks)", category)
		}
		return nil
	},
}

func listSkills(root string) error {
	skillList, err := skills.Discover(root)
	if err != nil {
		return err
	}
	if len(skillList) == 0 {
		fmt.Println("No skills found.")
		return nil
	}
	fmt.Println("Skills:")
	for _, s := range skillList {
		fmt.Printf("  %s\n", s.Name)
	}
	return nil
}

func listDisplays() error {
	viewNames := displays.ListRegistered()
	if len(viewNames) == 0 {
		fmt.Println("No displays registered.")
		return nil
	}
	fmt.Println("Displays:")
	for _, name := range viewNames {
		fmt.Printf("  %s\n", name)
	}
	return nil
}

func listBackground(root string) error {
	apps, err := background.Discover(root)
	if err != nil {
		return err
	}
	if len(apps) == 0 {
		fmt.Println("No background apps found.")
		return nil
	}
	fmt.Println("Background Apps:")
	for _, app := range apps {
		fmt.Printf("  %s\n", app.Name)
	}
	return nil
}

func listSidekicks(root string) error {
	sks, err := sidekicks.ListActive(root)
	if err != nil {
		return err
	}
	if len(sks) == 0 {
		fmt.Println("No active sidekicks.")
		return nil
	}
	fmt.Println("Active Sidekicks:")
	for _, sk := range sks {
		fmt.Printf("  %s: %s\n", sk.ID, sk.Prompt)
	}
	return nil
}

func listAll(root string) error {
	if root != "" {
		listSkills(root)
		fmt.Println()
	}
	listDisplays()
	if root != "" {
		fmt.Println()
		listBackground(root)
		fmt.Println()
		listSidekicks(root)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(listCmd)
}
