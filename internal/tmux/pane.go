package tmux

import (
	"fmt"
	"strconv"
	"strings"
)

// PaneInfo holds information about a tmux pane.
type PaneInfo struct {
	ID     string
	Width  int
	Height int
	Active bool
}

// SplitDirection indicates horizontal or vertical split.
type SplitDirection int

const (
	SplitHorizontal SplitDirection = iota // -h (side by side)
	SplitVertical                         // -v (top and bottom)
)

// SplitPane splits the current pane and returns the new pane ID.
func SplitPane(target string, dir SplitDirection, percentage int, cmd string) (string, error) {
	args := []string{"split-window"}
	if dir == SplitHorizontal {
		args = append(args, "-h")
	} else {
		args = append(args, "-v")
	}
	if percentage > 0 {
		args = append(args, "-p", strconv.Itoa(percentage))
	}
	if target != "" {
		args = append(args, "-t", target)
	}
	args = append(args, "-P", "-F", "#{pane_id}")
	if cmd != "" {
		args = append(args, cmd)
	}
	return Run(args...)
}

// SplitPaneBefore splits the pane with the new pane placed before (above/left).
func SplitPaneBefore(target string, dir SplitDirection, percentage int, cmd string) (string, error) {
	args := []string{"split-window", "-b"}
	if dir == SplitHorizontal {
		args = append(args, "-h")
	} else {
		args = append(args, "-v")
	}
	if percentage > 0 {
		args = append(args, "-p", strconv.Itoa(percentage))
	}
	if target != "" {
		args = append(args, "-t", target)
	}
	args = append(args, "-P", "-F", "#{pane_id}")
	if cmd != "" {
		args = append(args, cmd)
	}
	return Run(args...)
}

// KillPane kills a tmux pane.
func KillPane(paneID string) error {
	return RunSilent("kill-pane", "-t", paneID)
}

// FocusPane focuses a specific pane.
func FocusPane(paneID string) error {
	return RunSilent("select-pane", "-t", paneID)
}

// ZoomPane toggles zoom on a pane.
func ZoomPane(paneID string) error {
	return RunSilent("resize-pane", "-t", paneID, "-Z")
}

// GetPaneSize returns the width and height of a pane.
func GetPaneSize(paneID string) (int, int, error) {
	target := paneID
	if target == "" {
		target = ""
	}
	args := []string{"display-message"}
	if target != "" {
		args = append(args, "-t", target)
	}
	args = append(args, "-p", "#{pane_width} #{pane_height}")
	out, err := Run(args...)
	if err != nil {
		return 0, 0, err
	}
	parts := strings.SplitN(out, " ", 2)
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("unexpected pane size output: %q", out)
	}
	w, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("parse width: %w", err)
	}
	h, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("parse height: %w", err)
	}
	return w, h, nil
}

// PaneExists checks if a pane exists.
func PaneExists(paneID string) bool {
	_, err := Run("display-message", "-t", paneID, "-p", "")
	return err == nil
}

// ListPanes returns all panes in the current session.
func ListPanes() ([]PaneInfo, error) {
	out, err := Run("list-panes", "-F", "#{pane_id} #{pane_width} #{pane_height} #{pane_active}")
	if err != nil {
		return nil, err
	}
	var panes []PaneInfo
	for _, line := range strings.Split(out, "\n") {
		parts := strings.Fields(line)
		if len(parts) < 4 {
			continue
		}
		w, _ := strconv.Atoi(parts[1])
		h, _ := strconv.Atoi(parts[2])
		panes = append(panes, PaneInfo{
			ID:     parts[0],
			Width:  w,
			Height: h,
			Active: parts[3] == "1",
		})
	}
	return panes, nil
}

// CapturePane captures the visible content of a pane.
func CapturePane(paneID string) (string, error) {
	args := []string{"capture-pane", "-p"}
	if paneID != "" {
		args = append(args, "-t", paneID)
	}
	return Run(args...)
}
