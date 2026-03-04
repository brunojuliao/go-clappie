package displays

import (
	"fmt"

	"github.com/brunojuliao/go-clappie/internal/background"
	"github.com/brunojuliao/go-clappie/internal/engine"
	"github.com/brunojuliao/go-clappie/internal/platform"
	"github.com/brunojuliao/go-clappie/internal/uikit"
)

// NewBackgroundView creates the background apps dashboard view.
func NewBackgroundView(ctx *engine.Context) engine.View {
	view := uikit.NewView(ctx)

	ctx.SetTitle("Background")
	ctx.SetDescription("Long-running app management")

	var apps []background.App
	selectedIdx := 0

	load := func() {
		root, err := platform.ProjectRoot()
		if err != nil {
			return
		}
		apps, _ = background.List(root)
	}

	render := func() {
		var lines []string
		lines = append(lines, "")

		if len(apps) == 0 {
			lines = append(lines, "  No background apps found.")
		} else {
			lines = append(lines, fmt.Sprintf("  %d background apps", len(apps)))
			lines = append(lines, "")

			for i, app := range apps {
				prefix := "  "
				if i == selectedIdx {
					prefix = "▸ "
				}
				status := engine.StyleDim("stopped")
				if app.Running {
					status = engine.Color(40, 167, 69, "running")
				}
				line := fmt.Sprintf("%s%s [%s]", prefix, app.Name, status)
				if i == selectedIdx {
					line = engine.StyleBold(line)
				}
				lines = append(lines, line)
			}
		}

		ctx.Draw(lines)
	}

	view.RegisterShortcut("S", "Start/Stop", func() {
		if selectedIdx < len(apps) {
			root, _ := platform.ProjectRoot()
			app := apps[selectedIdx]
			if app.Running {
				background.Stop(app.Name)
			} else {
				background.Start(root, app.Name)
			}
			load()
			render()
		}
	})

	view.RegisterShortcut("R", "Refresh", func() {
		load()
		render()
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
				if selectedIdx < len(apps)-1 {
					selectedIdx++
					render()
				}
				return true
			}
			return view.HandleKey(key)
		},
	}
}
