package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/brunojuliao/go-clappie/internal/ipc"
	"github.com/brunojuliao/go-clappie/internal/platform"
	ttmux "github.com/brunojuliao/go-clappie/internal/tmux"
	"github.com/spf13/cobra"
)

var displayCmd = &cobra.Command{
	Use:   "display",
	Short: "Manage display views",
}

var displayPushCmd = &cobra.Command{
	Use:   "push [view]",
	Short: "Push a view onto the display stack",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		viewName := args[0]
		socketPath := platform.SocketPath()

		data, err := parseData()
		if err != nil {
			return err
		}

		// Check if daemon is running
		if !ipc.Ping(socketPath) {
			// Start daemon with this view
			return startDaemon(viewName, data)
		}

		// Send push command to running daemon
		command := ipc.Command{
			Action:  ipc.ActionPushView,
			View:    viewName,
			Data:    data,
			NoFocus: !focusFlag,
		}
		resp, err := ipc.SendCommand(socketPath, command)
		if err != nil {
			return fmt.Errorf("send push command: %w", err)
		}
		if !resp.OK {
			return fmt.Errorf("push failed: %s", resp.Error)
		}
		return nil
	},
}

var displayPopCmd = &cobra.Command{
	Use:   "pop",
	Short: "Pop the current view from the display stack",
	RunE: func(cmd *cobra.Command, args []string) error {
		socketPath := platform.SocketPath()
		resp, err := ipc.SendCommand(socketPath, ipc.Command{Action: ipc.ActionPopView})
		if err != nil {
			return fmt.Errorf("send pop command: %w", err)
		}
		if !resp.OK {
			return fmt.Errorf("pop failed: %s", resp.Error)
		}
		return nil
	},
}

var displayToastCmd = &cobra.Command{
	Use:   "toast [message]",
	Short: "Show a toast notification",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		socketPath := platform.SocketPath()
		command := ipc.Command{
			Action:   ipc.ActionToast,
			Message:  args[0],
			Duration: timeoutFlag,
		}
		resp, err := ipc.SendCommand(socketPath, command)
		if err != nil {
			return fmt.Errorf("send toast command: %w", err)
		}
		if !resp.OK {
			return fmt.Errorf("toast failed: %s", resp.Error)
		}
		return nil
	},
}

var displayCloseCmd = &cobra.Command{
	Use:   "close",
	Short: "Close the display daemon",
	RunE: func(cmd *cobra.Command, args []string) error {
		socketPath := platform.SocketPath()
		resp, err := ipc.SendCommand(socketPath, ipc.Command{Action: ipc.ActionClose})
		if err != nil {
			return fmt.Errorf("send close command: %w", err)
		}
		if !resp.OK {
			return fmt.Errorf("close failed: %s", resp.Error)
		}
		return nil
	},
}

var displayListCmd = &cobra.Command{
	Use:   "list",
	Short: "List views on the display stack",
	RunE: func(cmd *cobra.Command, args []string) error {
		socketPath := platform.SocketPath()
		resp, err := ipc.SendCommand(socketPath, ipc.Command{Action: ipc.ActionListViews})
		if err != nil {
			return fmt.Errorf("send list command: %w", err)
		}
		if !resp.OK {
			return fmt.Errorf("list failed: %s", resp.Error)
		}
		if resp.Data != nil {
			fmt.Println(string(resp.Data))
		}
		return nil
	},
}

var displayKillCmd = &cobra.Command{
	Use:   "kill",
	Short: "Kill the display daemon",
	RunE: func(cmd *cobra.Command, args []string) error {
		socketPath := platform.SocketPath()
		_, _ = ipc.SendCommand(socketPath, ipc.Command{Action: ipc.ActionKill})
		return nil
	},
}

func init() {
	displayCmd.AddCommand(displayPushCmd)
	displayCmd.AddCommand(displayPopCmd)
	displayCmd.AddCommand(displayToastCmd)
	displayCmd.AddCommand(displayCloseCmd)
	displayCmd.AddCommand(displayListCmd)
	displayCmd.AddCommand(displayKillCmd)
	rootCmd.AddCommand(displayCmd)
}

// startDaemon spawns the daemon process in a new tmux pane.
func startDaemon(initialView string, initialData json.RawMessage) error {
	if !platform.InTmux() {
		return fmt.Errorf("go-clappie must be run inside tmux")
	}

	claudePane := platform.TmuxPaneID()
	socketPath := platform.SocketPath()

	// Get pane dimensions to determine mobile vs desktop layout
	w, h, err := ttmux.GetPaneSize("")
	if err != nil {
		return fmt.Errorf("get pane size: %w", err)
	}
	isMobile := h > w || w < 120

	// Build the daemon command
	binary, err := os.Executable()
	if err != nil {
		return fmt.Errorf("get executable path: %w", err)
	}

	// Prepare environment
	env := fmt.Sprintf("GO_CLAPPIE_CLAUDE_PANE=%s GO_CLAPPIE_SOCKET_PATH=%s GO_CLAPPIE_INITIAL_VIEW=%s",
		claudePane, socketPath, initialView)
	if initialData != nil {
		env += fmt.Sprintf(" GO_CLAPPIE_INITIAL_DATA='%s'", string(initialData))
	}
	if os.Getenv("GO_CLAPPIE_ALLOW_SENDKEYS") != "" {
		env += " GO_CLAPPIE_ALLOW_SENDKEYS=1"
	}

	daemonCmd := fmt.Sprintf("%s %s __daemon", env, binary)

	var paneID string
	if isMobile {
		// Mobile: split vertically, put UI above (-b), 70%
		paneID, err = ttmux.SplitPaneBefore("", ttmux.SplitVertical, 70, daemonCmd)
		if err != nil {
			return fmt.Errorf("split pane (mobile): %w", err)
		}
		// Zoom the new pane for mobile
		ttmux.ZoomPane(paneID)
	} else {
		// Desktop: split horizontally (right), 70%
		paneID, err = ttmux.SplitPane("", ttmux.SplitHorizontal, 70, daemonCmd)
		if err != nil {
			return fmt.Errorf("split pane (desktop): %w", err)
		}
	}

	_ = paneID

	// Focus back to Claude pane unless -f flag
	if !focusFlag {
		ttmux.FocusPane(claudePane)
	}

	return nil
}

// startDaemonDirect starts the daemon process directly (for the __daemon command).
func startDaemonDirect() error {
	binary, err := exec.LookPath(os.Args[0])
	if err != nil {
		return err
	}

	socketPath := os.Getenv("GO_CLAPPIE_SOCKET_PATH")
	if socketPath == "" {
		socketPath = platform.SocketPath()
	}

	_ = binary
	_ = socketPath
	return nil
}

// isMobileLayout detects mobile layout based on terminal dimensions.
func isMobileLayout() bool {
	w, err := strconv.Atoi(os.Getenv("COLUMNS"))
	if err != nil {
		return false
	}
	h, err := strconv.Atoi(os.Getenv("LINES"))
	if err != nil {
		return false
	}
	return h > w || w < 120
}

// formatViewName converts slash-separated view names for display.
func formatViewName(name string) string {
	return strings.ReplaceAll(name, "/", " > ")
}
