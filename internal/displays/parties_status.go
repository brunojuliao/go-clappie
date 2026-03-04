package displays

import (
	"github.com/brunojuliao/go-clappie/internal/engine"
	"github.com/brunojuliao/go-clappie/internal/uikit"
)

// NewPartiesStatusView creates the parties simulation status view.
func NewPartiesStatusView(ctx *engine.Context) engine.View {
	view := uikit.NewView(ctx)

	ctx.SetTitle("Party Status")
	ctx.SetDescription("Simulation status")

	gameName := ""
	if ctx.Data != nil {
		if g, ok := ctx.Data["game"].(string); ok {
			gameName = g
		}
	}

	render := func() {
		var lines []string
		lines = append(lines, "")
		lines = append(lines, "  Game: "+engine.StyleBold(gameName))
		lines = append(lines, "")
		lines = append(lines, engine.StyleDim("  Status view - simulation details will appear here."))
		ctx.Draw(lines)
	}

	view.RegisterShortcut("L", "Launch", func() {
		ctx.Toast("Launching simulation...")
	})

	return engine.View{
		Init:   render,
		Render: render,
		OnKey: func(key string) bool {
			return view.HandleKey(key)
		},
	}
}
