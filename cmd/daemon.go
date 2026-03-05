package cmd

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/brunojuliao/go-clappie/internal/displays"
	"github.com/brunojuliao/go-clappie/internal/engine"
)

var (
	daemonSocket     string
	daemonView       string
	daemonData       string
	daemonClaudePane string
)

var daemonCmd = &cobra.Command{
	Use:    "__daemon",
	Short:  "Run the display daemon (internal)",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Flags take precedence over env vars
		socketPath := daemonSocket
		if socketPath == "" {
			socketPath = os.Getenv("GO_CLAPPIE_SOCKET_PATH")
		}
		if socketPath == "" {
			return fmt.Errorf("GO_CLAPPIE_SOCKET_PATH not set")
		}

		initialView := daemonView
		if initialView == "" {
			initialView = os.Getenv("GO_CLAPPIE_INITIAL_VIEW")
		}

		initialData := daemonData
		if initialData == "" {
			initialData = os.Getenv("GO_CLAPPIE_INITIAL_DATA")
		}

		claudePane := daemonClaudePane
		if claudePane == "" {
			claudePane = os.Getenv("GO_CLAPPIE_CLAUDE_PANE")
		}

		// Log to file for debugging (stderr may not be visible in tmux pane).
		logFile, err := os.CreateTemp("", "go-clappie-daemon-*.log")
		if err == nil {
			log.SetOutput(logFile)
			defer logFile.Close()
			log.Printf("daemon starting: socket=%s view=%s claude=%s", socketPath, initialView, claudePane)
		}

		app := engine.NewApp(engine.AppConfig{
			SocketPath:  socketPath,
			InitialView: initialView,
			InitialData: initialData,
			ClaudePane:  claudePane,
			Registry:    displays.Registry,
		})

		p := tea.NewProgram(app, tea.WithAltScreen(), tea.WithMouseCellMotion())
		app.SetProgram(p)

		if err := app.StartIPCServer(); err != nil {
			log.Printf("IPC server failed: %v", err)
			return fmt.Errorf("start IPC server: %w", err)
		}
		defer app.Shutdown()

		log.Printf("starting bubbletea program")
		_, err = p.Run()
		if err != nil {
			log.Printf("bubbletea exited with error: %v", err)
		} else {
			log.Printf("bubbletea exited cleanly")
		}
		return err
	},
}

func init() {
	daemonCmd.Flags().StringVar(&daemonSocket, "socket", "", "IPC socket path")
	daemonCmd.Flags().StringVar(&daemonView, "view", "", "Initial view to display")
	daemonCmd.Flags().StringVar(&daemonData, "data", "", "Initial view data (JSON)")
	daemonCmd.Flags().StringVar(&daemonClaudePane, "claude-pane", "", "Claude's tmux pane ID")
	rootCmd.AddCommand(daemonCmd)
}
