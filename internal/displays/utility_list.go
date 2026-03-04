package displays

import (
	"fmt"

	"github.com/brunojuliao/go-clappie/internal/engine"
	"github.com/brunojuliao/go-clappie/internal/uikit"
)

// NewUtilityListView creates the utility list picker view.
func NewUtilityListView(ctx *engine.Context) engine.View {
	view := uikit.NewView(ctx)

	title := "Select"
	if ctx.Data != nil {
		if t, ok := ctx.Data["title"].(string); ok {
			title = t
		}
	}

	ctx.SetTitle(title)

	var options []string
	if ctx.Data != nil {
		if opts, ok := ctx.Data["options"].([]interface{}); ok {
			for _, o := range opts {
				if s, ok := o.(string); ok {
					options = append(options, s)
				}
			}
		}
	}

	selectedIdx := 0

	render := func() {
		var lines []string
		lines = append(lines, "")

		for i, opt := range options {
			prefix := "  "
			if i == selectedIdx {
				prefix = "▸ "
			}
			line := fmt.Sprintf("%s%s", prefix, opt)
			if i == selectedIdx {
				line = engine.StyleBold(line)
			}
			lines = append(lines, line)
		}

		ctx.Draw(lines)
	}

	return engine.View{
		Init:   render,
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
				if selectedIdx < len(options)-1 {
					selectedIdx++
					render()
				}
				return true
			case "ENTER":
				if selectedIdx < len(options) {
					ctx.Submit(fmt.Sprintf("[go-clappie] List → %s", options[selectedIdx]))
					ctx.Pop()
				}
				return true
			}
			return view.HandleKey(key)
		},
	}
}
