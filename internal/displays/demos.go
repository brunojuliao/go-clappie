package displays

import (
	"fmt"

	"github.com/brunojuliao/go-clappie/internal/engine"
	"github.com/brunojuliao/go-clappie/internal/uikit"
)

// NewHelloWorldView creates the hello world demo view.
func NewHelloWorldView(ctx *engine.Context) engine.View {
	view := uikit.NewView(ctx)

	ctx.SetTitle("Hello")
	ctx.SetDescription("Hello World Demo")

	partyMode := false

	render := func() {
		var lines []string
		lines = append(lines, "")
		lines = append(lines, "  Hello, World!")
		lines = append(lines, "")

		if partyMode {
			lines = append(lines, "  🎉🎊🥳🎉🎊🥳🎉")
			lines = append(lines, "")
		}

		lines = append(lines, engine.StyleDim("  Welcome to Clappie!"))
		lines = append(lines, engine.StyleDim("  This is a demo display."))
		lines = append(lines, "")
		lines = append(lines, fmt.Sprintf("  Party mode: %v", partyMode))

		ctx.Draw(lines)
	}

	view.RegisterShortcut("P", "Party", func() {
		partyMode = !partyMode
		render()
	})

	return engine.View{
		Init:   render,
		Render: render,
		OnKey: func(key string) bool {
			return view.HandleKey(key)
		},
	}
}

// NewAllComponentsView creates the all-components demo view.
func NewAllComponentsView(ctx *engine.Context) engine.View {
	view := uikit.NewView(ctx)

	ctx.SetTitle("Components")
	ctx.SetDescription("All UI Components Demo")

	// Add one of each component type
	view.Add(uikit.NewSectionHeading(uikit.SectionHeadingConfig{Text: "Buttons"}))
	view.Add(uikit.NewButton(uikit.ButtonConfig{Label: "Default Button", OnPress: func() { ctx.Toast("Pressed!") }}))
	view.Add(uikit.NewButtonFilled(uikit.ButtonConfig{Label: "Filled Button", OnPress: func() { ctx.Toast("Filled!") }}))
	view.Add(uikit.NewButtonGhost(uikit.ButtonConfig{Label: "Ghost Button", OnPress: func() { ctx.Toast("Ghost!") }}))

	view.Add(uikit.NewSectionHeading(uikit.SectionHeadingConfig{Text: "Toggles"}))
	view.Add(uikit.NewToggle(uikit.ToggleConfig{Label: "Dark Mode", Value: true}))
	view.Add(uikit.NewToggle(uikit.ToggleConfig{Label: "Notifications", Value: false}))

	view.Add(uikit.NewSectionHeading(uikit.SectionHeadingConfig{Text: "Text Input"}))
	view.Add(uikit.NewTextInput(uikit.TextInputConfig{Placeholder: "Type something...", Width: 40}))

	view.Add(uikit.NewSectionHeading(uikit.SectionHeadingConfig{Text: "Textarea"}))
	view.Add(uikit.NewTextarea(uikit.TextareaConfig{Placeholder: "Multi-line input...", Width: 40, Height: 3}))

	view.Add(uikit.NewSectionHeading(uikit.SectionHeadingConfig{Text: "Checkbox"}))
	view.Add(uikit.NewCheckbox(uikit.CheckboxConfig{Label: "I agree to the terms"}))
	view.Add(uikit.NewCheckbox(uikit.CheckboxConfig{Label: "Send notifications", Checked: true}))

	view.Add(uikit.NewSectionHeading(uikit.SectionHeadingConfig{Text: "Radio"}))
	view.Add(uikit.NewRadio(uikit.RadioConfig{
		Label:   "Theme",
		Options: []string{"Dark", "Light", "Auto"},
	}))

	view.Add(uikit.NewSectionHeading(uikit.SectionHeadingConfig{Text: "Select"}))
	view.Add(uikit.NewSelect(uikit.SelectConfig{
		Label:   "Model",
		Options: []string{"Claude 4", "Claude 3.5", "Claude 3"},
	}))

	view.Add(uikit.NewSectionHeading(uikit.SectionHeadingConfig{Text: "Progress"}))
	view.Add(uikit.NewProgress(uikit.ProgressConfig{Label: "Upload", Value: 0.7, Width: 40}))

	view.Add(uikit.NewSectionHeading(uikit.SectionHeadingConfig{Text: "Loader"}))
	view.Add(uikit.NewLoader(uikit.LoaderConfig{Label: "Processing..."}))

	view.Add(uikit.NewDivider(uikit.DividerConfig{Width: 50}))

	view.Add(uikit.NewSectionHeading(uikit.SectionHeadingConfig{Text: "Alerts"}))
	view.Add(uikit.NewAlert(uikit.AlertConfig{Type: uikit.AlertInfo, Message: "Informational message", Width: 50}))
	view.Add(uikit.NewAlert(uikit.AlertConfig{Type: uikit.AlertSuccess, Message: "Operation completed", Width: 50}))
	view.Add(uikit.NewAlert(uikit.AlertConfig{Type: uikit.AlertWarning, Message: "Proceed with caution", Width: 50}))
	view.Add(uikit.NewAlert(uikit.AlertConfig{Type: uikit.AlertError, Message: "Something went wrong", Width: 50}))

	view.Add(uikit.NewLabel(uikit.LabelConfig{Text: "This is a label", Dim: true}))

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
