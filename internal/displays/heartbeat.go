package displays

import (
	"fmt"
	"time"

	"github.com/brunojuliao/go-clappie/internal/engine"
	"github.com/brunojuliao/go-clappie/internal/filestore"
	"github.com/brunojuliao/go-clappie/internal/platform"
	"github.com/brunojuliao/go-clappie/internal/uikit"
)

// NewHeartbeatView creates the heartbeat dashboard view.
func NewHeartbeatView(ctx *engine.Context) engine.View {
	view := uikit.NewView(ctx)

	ctx.SetTitle("Heartbeat")
	ctx.SetDescription("AI-powered cron scheduler")

	var checks []heartbeatCheck

	loadChecks := func() {
		root, err := platform.ProjectRoot()
		if err != nil {
			return
		}

		checksDir := platform.ChoresBotsDir(root)
		entries, err := filestore.List(checksDir)
		if err != nil {
			return
		}

		checks = nil
		for _, entry := range entries {
			body, blocks, err := filestore.ReadAndParse(entry.Path)
			if err != nil {
				continue
			}
			meta := filestore.GetMeta(blocks, "heartbeat-meta")
			interval := ""
			lastRun := ""
			status := ""
			if meta != nil {
				interval = meta.Fields["interval"]
				lastRun = meta.Fields["last_run"]
				status = meta.Fields["status"]
			}
			checks = append(checks, heartbeatCheck{
				Name:     entry.Name,
				Body:     body,
				Interval: interval,
				LastRun:  lastRun,
				Status:   status,
			})
		}
	}

	render := func() {
		var lines []string
		lines = append(lines, "")

		if len(checks) == 0 {
			lines = append(lines, "  No heartbeat checks configured.")
			lines = append(lines, "")
			lines = append(lines, "  Add check files to chores/bots/ to get started.")
		} else {
			lines = append(lines, fmt.Sprintf("  %d checks configured", len(checks)))
			lines = append(lines, "")
			for _, c := range checks {
				statusIcon := "○"
				if c.Status == "ok" {
					statusIcon = "●"
				} else if c.Status == "error" {
					statusIcon = "✗"
				}
				line := fmt.Sprintf("  %s %s", statusIcon, c.Name)
				if c.Interval != "" {
					line += fmt.Sprintf(" (every %s)", c.Interval)
				}
				if c.LastRun != "" {
					line += fmt.Sprintf(" — last: %s", c.LastRun)
				}
				lines = append(lines, line)
			}
		}

		lines = append(lines, "")
		lines = append(lines, fmt.Sprintf("  Last refreshed: %s", time.Now().Format("15:04:05")))

		ctx.Draw(lines)
	}

	// Refresh button
	view.Add(uikit.NewButton(uikit.ButtonConfig{
		Label:    "Refresh",
		Shortcut: "R",
		OnPress: func() {
			loadChecks()
			render()
		},
	}))

	view.RegisterShortcut("R", "Refresh", func() {
		loadChecks()
		render()
	})

	return engine.View{
		Init: func() {
			loadChecks()
			render()
		},
		Render: render,
		OnKey: func(key string) bool {
			return view.HandleKey(key)
		},
	}
}

type heartbeatCheck struct {
	Name     string
	Body     string
	Interval string
	LastRun  string
	Status   string
}
