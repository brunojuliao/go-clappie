package displays

import (
	"github.com/brunojuliao/go-clappie/internal/engine"
	"github.com/brunojuliao/go-clappie/internal/uikit"
)

// NewUtilityConfirmView creates the utility confirm dialog view.
func NewUtilityConfirmView(ctx *engine.Context) engine.View {
	view := uikit.NewView(ctx)

	message := "Are you sure?"
	if ctx.Data != nil {
		if m, ok := ctx.Data["message"].(string); ok {
			message = m
		}
	}

	ctx.SetTitle("Confirm")

	render := func() {
		lines := []string{
			"",
			"  " + message,
			"",
		}
		ctx.Draw(lines)
	}

	view.Add(uikit.NewButton(uikit.ButtonConfig{
		Label:    "Yes",
		Shortcut: "Y",
		OnPress: func() {
			ctx.Submit("[clappie] Confirm → yes")
			ctx.Pop()
		},
	}))

	view.Add(uikit.NewButton(uikit.ButtonConfig{
		Label:    "No",
		Shortcut: "N",
		OnPress: func() {
			ctx.Submit("[clappie] Confirm → no")
			ctx.Pop()
		},
	}))

	return engine.View{
		Init:   render,
		Render: render,
		OnKey: func(key string) bool {
			return view.HandleKey(key)
		},
	}
}
