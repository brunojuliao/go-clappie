package uikit

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"
)

// TextInputConfig configures a text input component.
type TextInputConfig struct {
	Placeholder string
	Width       int
	Value       string
	OnChange    func(string) tea.Cmd
	OnSubmit    func(string) tea.Cmd
}

// TextInput wraps bubbles/textinput.Model.
type TextInput struct {
	config TextInputConfig
	model  textinput.Model
}

// NewTextInput creates a new text input.
func NewTextInput(cfg TextInputConfig) TextInput {
	m := textinput.New()
	m.Placeholder = cfg.Placeholder
	if cfg.Width > 0 {
		m.Width = cfg.Width - 2 // account for border padding
	}
	if cfg.Value != "" {
		m.SetValue(cfg.Value)
	}
	return TextInput{config: cfg, model: m}
}

func (ti TextInput) Init() tea.Cmd { return nil }

func (ti TextInput) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if keyMsg.String() == "enter" && ti.config.OnSubmit != nil {
			return ti, ti.config.OnSubmit(ti.model.Value())
		}
	}

	var cmd tea.Cmd
	oldVal := ti.model.Value()
	ti.model, cmd = ti.model.Update(msg)

	if ti.model.Value() != oldVal && ti.config.OnChange != nil {
		changeCm := ti.config.OnChange(ti.model.Value())
		return ti, tea.Batch(cmd, changeCm)
	}

	return ti, cmd
}

func (ti TextInput) View() string {
	return ti.model.View()
}

func (ti TextInput) IsFocusable() bool { return true }
func (ti TextInput) Focused() bool     { return ti.model.Focused() }

func (ti TextInput) Focus() Component {
	ti.model.Focus()
	return ti
}

func (ti TextInput) Blur() Component {
	ti.model.Blur()
	return ti
}

// Value returns the current input value.
func (ti TextInput) Value() string {
	return ti.model.Value()
}

// SetValue sets the input value.
func (ti *TextInput) SetValue(v string) {
	ti.model.SetValue(v)
}
