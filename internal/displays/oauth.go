package displays

import (
	"fmt"

	"github.com/brunojuliao/go-clappie/internal/engine"
	"github.com/brunojuliao/go-clappie/internal/oauth"
	"github.com/brunojuliao/go-clappie/internal/platform"
	"github.com/brunojuliao/go-clappie/internal/uikit"
)

// NewOAuthView creates the OAuth management view.
func NewOAuthView(ctx *engine.Context) engine.View {
	view := uikit.NewView(ctx)

	ctx.SetTitle("OAuth")
	ctx.SetDescription("Token management")

	var providers []oauth.ProviderInfo
	selectedIdx := 0

	load := func() {
		root, err := platform.ProjectRoot()
		if err != nil {
			return
		}
		providers, _ = oauth.ListProviders(root)
	}

	render := func() {
		var lines []string
		lines = append(lines, "")

		if len(providers) == 0 {
			lines = append(lines, "  No OAuth providers configured.")
			lines = append(lines, "")
			lines = append(lines, engine.StyleDim("  Add oauth.json to a skill directory."))
		} else {
			lines = append(lines, fmt.Sprintf("  %d providers", len(providers)))
			lines = append(lines, "")

			for i, p := range providers {
				prefix := "  "
				if i == selectedIdx {
					prefix = "▸ "
				}
				status := engine.StyleDim("not authenticated")
				if p.HasToken {
					status = engine.Color(40, 167, 69, "authenticated")
					if p.Expired {
						status = engine.Color(255, 193, 7, "expired")
					}
				}
				line := fmt.Sprintf("%s%s [%s]", prefix, p.Name, status)
				if i == selectedIdx {
					line = engine.StyleBold(line)
				}
				lines = append(lines, line)
			}
		}

		ctx.Draw(lines)
	}

	view.RegisterShortcut("A", "Auth", func() {
		if selectedIdx < len(providers) {
			root, _ := platform.ProjectRoot()
			oauth.Auth(root, providers[selectedIdx].Name)
			ctx.Toast("Auth flow started...")
		}
	})

	view.RegisterShortcut("R", "Refresh", func() {
		if selectedIdx < len(providers) {
			root, _ := platform.ProjectRoot()
			oauth.Refresh(root, providers[selectedIdx].Name)
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
				if selectedIdx < len(providers)-1 {
					selectedIdx++
					render()
				}
				return true
			}
			return view.HandleKey(key)
		},
	}
}
