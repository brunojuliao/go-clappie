package displays

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/brunojuliao/go-clappie/internal/engine"
	"github.com/brunojuliao/go-clappie/internal/uikit"
)

// --- Hello World ---

type helloWorldScreen struct {
	partyMode bool
	styles    *engine.Styles
}

func NewHelloWorldScreen(data map[string]interface{}, styles *engine.Styles, claudePane string) engine.ScreenModel {
	return &helloWorldScreen{styles: styles}
}

func (m *helloWorldScreen) Init() tea.Cmd { return nil }

func (m *helloWorldScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if keyMsg.String() == "p" || keyMsg.String() == "P" {
			m.partyMode = !m.partyMode
		}
	}
	return m, nil
}

func (m *helloWorldScreen) View() string {
	dim := lipgloss.NewStyle().Faint(true)

	var lines []string
	lines = append(lines, "")
	lines = append(lines, "  Hello, World!")
	lines = append(lines, "")

	if m.partyMode {
		lines = append(lines, "  🎉🎊🥳🎉🎊🥳🎉")
		lines = append(lines, "")
	}

	lines = append(lines, dim.Render("  Welcome to Go-Clappie!"))
	lines = append(lines, dim.Render("  This is a demo display."))
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("  Party mode: %v", m.partyMode))

	return strings.Join(lines, "\n")
}

func (m *helloWorldScreen) Name() string          { return "Hello" }
func (m *helloWorldScreen) Layout() (string, int) { return "centered", 50 }
func (m *helloWorldScreen) Shortcuts() []engine.ShortcutHint {
	return []engine.ShortcutHint{{Key: "P", Label: "Party"}}
}

// --- All Components ---

type allComponentsScreen struct {
	container uikit.ViewContainer
	styles    *engine.Styles
}

func NewAllComponentsScreen(data map[string]interface{}, styles *engine.Styles, claudePane string) engine.ScreenModel {
	m := &allComponentsScreen{
		container: uikit.NewViewContainer(),
		styles:    styles,
	}

	vc := &m.container

	vc.Add(uikit.NewSectionHeading(uikit.SectionHeadingConfig{Text: "Buttons"}))
	vc.Add(uikit.NewButton(uikit.ButtonConfig{
		Label:   "Default Button",
		OnPress: func() tea.Cmd { return engine.ToastCmd("Pressed!", 0) },
	}))
	vc.Add(uikit.NewButtonFilled(uikit.ButtonConfig{
		Label:   "Filled Button",
		OnPress: func() tea.Cmd { return engine.ToastCmd("Filled!", 0) },
	}))
	vc.Add(uikit.NewButtonGhost(uikit.ButtonConfig{
		Label:   "Ghost Button",
		OnPress: func() tea.Cmd { return engine.ToastCmd("Ghost!", 0) },
	}))

	vc.Add(uikit.NewSectionHeading(uikit.SectionHeadingConfig{Text: "Toggles"}))
	vc.Add(uikit.NewToggle(uikit.ToggleConfig{Label: "Dark Mode", Value: true}))
	vc.Add(uikit.NewToggle(uikit.ToggleConfig{Label: "Notifications", Value: false}))

	vc.Add(uikit.NewSectionHeading(uikit.SectionHeadingConfig{Text: "Text Input"}))
	vc.Add(uikit.NewTextInput(uikit.TextInputConfig{Placeholder: "Type something...", Width: 40}))

	vc.Add(uikit.NewSectionHeading(uikit.SectionHeadingConfig{Text: "Textarea"}))
	vc.Add(uikit.NewTextarea(uikit.TextareaConfig{Placeholder: "Multi-line input...", Width: 40, Height: 3}))

	vc.Add(uikit.NewSectionHeading(uikit.SectionHeadingConfig{Text: "Checkbox"}))
	vc.Add(uikit.NewCheckbox(uikit.CheckboxConfig{Label: "I agree to the terms"}))
	vc.Add(uikit.NewCheckbox(uikit.CheckboxConfig{Label: "Send notifications", Checked: true}))

	vc.Add(uikit.NewSectionHeading(uikit.SectionHeadingConfig{Text: "Radio"}))
	vc.Add(uikit.NewRadio(uikit.RadioConfig{
		Label:   "Theme",
		Options: []string{"Dark", "Light", "Auto"},
	}))

	vc.Add(uikit.NewSectionHeading(uikit.SectionHeadingConfig{Text: "Select"}))
	vc.Add(uikit.NewSelect(uikit.SelectConfig{
		Label:   "Model",
		Options: []string{"Claude 4", "Claude 3.5", "Claude 3"},
	}))

	vc.Add(uikit.NewSectionHeading(uikit.SectionHeadingConfig{Text: "Progress"}))
	vc.Add(uikit.NewProgress(uikit.ProgressConfig{Label: "Upload", Value: 0.7, Width: 40}))

	vc.Add(uikit.NewSectionHeading(uikit.SectionHeadingConfig{Text: "Loader"}))
	vc.Add(uikit.NewLoader(uikit.LoaderConfig{Label: "Processing..."}))

	vc.Add(uikit.NewDivider(uikit.DividerConfig{Width: 50}))

	vc.Add(uikit.NewSectionHeading(uikit.SectionHeadingConfig{Text: "Alerts"}))
	vc.Add(uikit.NewAlert(uikit.AlertConfig{Type: uikit.AlertInfo, Message: "Informational message", Width: 50}))
	vc.Add(uikit.NewAlert(uikit.AlertConfig{Type: uikit.AlertSuccess, Message: "Operation completed", Width: 50}))
	vc.Add(uikit.NewAlert(uikit.AlertConfig{Type: uikit.AlertWarning, Message: "Proceed with caution", Width: 50}))
	vc.Add(uikit.NewAlert(uikit.AlertConfig{Type: uikit.AlertError, Message: "Something went wrong", Width: 50}))

	vc.Add(uikit.NewLabel(uikit.LabelConfig{Text: "This is a label", Dim: true}))

	return m
}

func (m *allComponentsScreen) Init() tea.Cmd { return nil }

func (m *allComponentsScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.container, cmd = m.container.Update(msg)
	return m, cmd
}

func (m *allComponentsScreen) View() string {
	return m.container.View()
}

func (m *allComponentsScreen) Name() string          { return "Components" }
func (m *allComponentsScreen) Layout() (string, int) { return "centered", 60 }
func (m *allComponentsScreen) Shortcuts() []engine.ShortcutHint {
	return nil
}
