package displays

import (
	"fmt"

	"github.com/brunojuliao/go-clappie/internal/engine"
	"github.com/brunojuliao/go-clappie/internal/parties"
	"github.com/brunojuliao/go-clappie/internal/platform"
	"github.com/brunojuliao/go-clappie/internal/uikit"
)

// NewPartiesIndexView creates the parties game listing view.
func NewPartiesIndexView(ctx *engine.Context) engine.View {
	view := uikit.NewView(ctx)

	ctx.SetTitle("Parties")
	ctx.SetDescription("AI swarm simulations")

	var games []parties.GameInfo
	selectedIdx := 0

	load := func() {
		root, err := platform.ProjectRoot()
		if err != nil {
			return
		}
		games, _ = parties.ListGames(root)
	}

	render := func() {
		var lines []string
		lines = append(lines, "")

		if len(games) == 0 {
			lines = append(lines, "  No games found.")
			lines = append(lines, "")
			lines = append(lines, engine.StyleDim("  Create game files to get started."))
		} else {
			lines = append(lines, fmt.Sprintf("  %d games available", len(games)))
			lines = append(lines, "")

			for i, g := range games {
				prefix := "  "
				if i == selectedIdx {
					prefix = "▸ "
				}
				line := fmt.Sprintf("%s%s", prefix, g.Name)
				if g.Description != "" {
					line += engine.StyleDim(" — "+g.Description)
				}
				if i == selectedIdx {
					line = engine.StyleBold(line)
				}
				lines = append(lines, line)
			}
		}

		ctx.Draw(lines)
	}

	view.RegisterShortcut("I", "Init", func() {
		if selectedIdx < len(games) {
			root, _ := platform.ProjectRoot()
			simID, err := parties.Init(root, games[selectedIdx].Name)
			if err != nil {
				ctx.Toast(fmt.Sprintf("Error: %v", err))
			} else {
				ctx.Toast(fmt.Sprintf("Initialized: %s", simID))
			}
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
				if selectedIdx < len(games)-1 {
					selectedIdx++
					render()
				}
				return true
			case "ENTER":
				if selectedIdx < len(games) {
					ctx.Push("parties/status", map[string]interface{}{
						"game": games[selectedIdx].Name,
					})
				}
				return true
			}
			return view.HandleKey(key)
		},
	}
}
