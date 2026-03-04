package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/brunojuliao/go-clappie/internal/displays"
	"github.com/brunojuliao/go-clappie/internal/engine"
	"github.com/spf13/cobra"
)

var daemonCmd = &cobra.Command{
	Use:    "__daemon",
	Short:  "Run the display daemon (internal)",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		socketPath := os.Getenv("CLAPPIE_SOCKET_PATH")
		if socketPath == "" {
			return fmt.Errorf("CLAPPIE_SOCKET_PATH not set")
		}

		initialView := os.Getenv("CLAPPIE_INITIAL_VIEW")
		initialData := os.Getenv("CLAPPIE_INITIAL_DATA")
		claudePane := os.Getenv("CLAPPIE_CLAUDE_PANE")

		d, err := engine.NewDaemon(engine.DaemonConfig{
			SocketPath:  socketPath,
			InitialView: initialView,
			InitialData: initialData,
			ClaudePane:  claudePane,
			Registry:    displays.Registry,
		})
		if err != nil {
			return fmt.Errorf("create daemon: %w", err)
		}

		// Handle signals
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
		go func() {
			<-sigCh
			d.Shutdown()
		}()

		return d.Run()
	},
}

func init() {
	rootCmd.AddCommand(daemonCmd)
}
