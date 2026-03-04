package cmd

import (
	"fmt"

	"github.com/brunojuliao/go-clappie/internal/oauth"
	"github.com/brunojuliao/go-clappie/internal/platform"
	"github.com/spf13/cobra"
)

var oauthCmd = &cobra.Command{
	Use:   "oauth",
	Short: "Manage OAuth tokens",
}

var oauthAuthCmd = &cobra.Command{
	Use:   "auth [provider]",
	Short: "Start OAuth authorization flow",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		root, err := platform.ProjectRoot()
		if err != nil {
			return err
		}
		return oauth.Auth(root, args[0])
	},
}

var oauthTokenCmd = &cobra.Command{
	Use:   "token [provider]",
	Short: "Get the current token for a provider",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		root, err := platform.ProjectRoot()
		if err != nil {
			return err
		}
		token, err := oauth.GetToken(root, args[0])
		if err != nil {
			return err
		}
		fmt.Println(token)
		return nil
	},
}

var oauthStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show OAuth status for all providers",
	RunE: func(cmd *cobra.Command, args []string) error {
		root, err := platform.ProjectRoot()
		if err != nil {
			return err
		}
		return oauth.Status(root)
	},
}

var oauthRefreshCmd = &cobra.Command{
	Use:   "refresh [provider]",
	Short: "Refresh a token",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		root, err := platform.ProjectRoot()
		if err != nil {
			return err
		}
		return oauth.Refresh(root, args[0])
	},
}

var oauthRevokeCmd = &cobra.Command{
	Use:   "revoke [provider]",
	Short: "Revoke a token",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		root, err := platform.ProjectRoot()
		if err != nil {
			return err
		}
		return oauth.Revoke(root, args[0])
	},
}

var oauthProvidersCmd = &cobra.Command{
	Use:   "providers",
	Short: "List available OAuth providers",
	RunE: func(cmd *cobra.Command, args []string) error {
		root, err := platform.ProjectRoot()
		if err != nil {
			return err
		}
		providers, err := oauth.ListProviders(root)
		if err != nil {
			return err
		}
		if len(providers) == 0 {
			fmt.Println("No OAuth providers found.")
			return nil
		}
		for _, p := range providers {
			fmt.Printf("  %s\n", p.Name)
		}
		return nil
	},
}

func init() {
	oauthCmd.AddCommand(oauthAuthCmd)
	oauthCmd.AddCommand(oauthTokenCmd)
	oauthCmd.AddCommand(oauthStatusCmd)
	oauthCmd.AddCommand(oauthRefreshCmd)
	oauthCmd.AddCommand(oauthRevokeCmd)
	oauthCmd.AddCommand(oauthProvidersCmd)
	rootCmd.AddCommand(oauthCmd)
}
