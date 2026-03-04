package displays

import (
	"github.com/brunojuliao/go-clappie/internal/engine"
	"github.com/brunojuliao/go-clappie/internal/uikit"
)

// NewUtilityEditorView creates the utility text editor view.
func NewUtilityEditorView(ctx *engine.Context) engine.View {
	view := uikit.NewView(ctx)

	ctx.SetTitle("Editor")
	ctx.SetLayout("full", 0)

	initialValue := ""
	if ctx.Data != nil {
		if v, ok := ctx.Data["value"].(string); ok {
			initialValue = v
		}
	}

	textarea := uikit.NewTextarea(uikit.TextareaConfig{
		Placeholder: "Type here...",
		Width:       60,
		Height:      20,
		Value:       initialValue,
	})
	view.Add(textarea)

	view.Add(uikit.NewButton(uikit.ButtonConfig{
		Label:    "Save",
		Shortcut: "S",
		OnPress: func() {
			ctx.Submit("[clappie] Editor → " + textarea.Value())
			ctx.Pop()
		},
		Style: uikit.ButtonStyleFilled,
	}))

	render := func() {
		lines := view.Render()
		ctx.Draw(lines)
	}

	return engine.View{
		Init:   render,
		Render: render,
		OnKey: func(key string) bool {
			if view.HandleKey(key) {
				render()
				return true
			}
			return false
		},
	}
}
