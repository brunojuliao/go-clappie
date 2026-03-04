package displays

import (
	"fmt"

	"github.com/brunojuliao/go-clappie/internal/engine"
	"github.com/brunojuliao/go-clappie/internal/platform"
	"github.com/brunojuliao/go-clappie/internal/sidekicks"
	"github.com/brunojuliao/go-clappie/internal/uikit"
)

// NewSidekicksView creates the sidekicks dashboard view.
func NewSidekicksView(ctx *engine.Context) engine.View {
	view := uikit.NewView(ctx)

	ctx.SetTitle("Sidekicks")
	ctx.SetDescription("Autonomous agent management")

	var activeSidekicks []sidekicks.SidekickInfo
	selectedIdx := 0

	load := func() {
		root, err := platform.ProjectRoot()
		if err != nil {
			return
		}
		activeSidekicks, _ = sidekicks.ListActive(root)
	}

	render := func() {
		var lines []string
		lines = append(lines, "")

		if len(activeSidekicks) == 0 {
			lines = append(lines, "  No active sidekicks.")
			lines = append(lines, "")
			lines = append(lines, engine.StyleDim("  Use 'go-clappie sidekick spawn \"prompt\"' to spawn one."))
		} else {
			lines = append(lines, fmt.Sprintf("  %d active sidekicks", len(activeSidekicks)))
			lines = append(lines, "")

			for i, sk := range activeSidekicks {
				prefix := "  "
				if i == selectedIdx {
					prefix = "▸ "
				}
				line := fmt.Sprintf("%s%s: %s", prefix, sk.ID, sk.Prompt)
				if len(line) > 65 {
					line = line[:62] + "..."
				}
				if i == selectedIdx {
					line = engine.StyleBold(line)
				}
				lines = append(lines, line)
			}
		}

		ctx.Draw(lines)
	}

	view.RegisterShortcut("R", "Refresh", func() {
		load()
		render()
	})

	view.RegisterShortcut("K", "Kill Selected", func() {
		if selectedIdx < len(activeSidekicks) {
			root, _ := platform.ProjectRoot()
			sidekicks.End(root)
			load()
			render()
		}
	})

	return engine.View{
		Init: func() {
			load()
			render()
		},
		Render: render,
		OnKey: func(key string) bool {
			switch key {
			case "UP", "k":
				if selectedIdx > 0 {
					selectedIdx--
					render()
				}
				return true
			case "DOWN", "j":
				if selectedIdx < len(activeSidekicks)-1 {
					selectedIdx++
					render()
				}
				return true
			}
			return view.HandleKey(key)
		},
	}
}
