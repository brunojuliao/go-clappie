package uikit

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textarea"
)

// TextareaConfig configures a textarea component.
type TextareaConfig struct {
	Placeholder string
	Width       int
	Height      int
	Value       string
	OnChange    func(string) tea.Cmd
	OnSubmit    func(string) tea.Cmd
}

// Textarea wraps bubbles/textarea.Model.
type Textarea struct {
	config TextareaConfig
	model  textarea.Model
}

// NewTextarea creates a new textarea.
func NewTextarea(cfg TextareaConfig) Textarea {
	m := textarea.New()
	m.Placeholder = cfg.Placeholder
	if cfg.Width > 0 {
		m.SetWidth(cfg.Width)
	}
	if cfg.Height > 0 {
		m.SetHeight(cfg.Height)
	}
	if cfg.Value != "" {
		m.SetValue(cfg.Value)
	}
	return Textarea{config: cfg, model: m}
}

func (ta Textarea) Init() tea.Cmd { return nil }

func (ta Textarea) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	oldVal := ta.model.Value()
	ta.model, cmd = ta.model.Update(msg)

	if ta.model.Value() != oldVal && ta.config.OnChange != nil {
		changeCm := ta.config.OnChange(ta.model.Value())
		return ta, tea.Batch(cmd, changeCm)
	}

	return ta, cmd
}

func (ta Textarea) View() string {
	return ta.model.View()
}

func (ta Textarea) IsFocusable() bool { return true }
func (ta Textarea) Focused() bool     { return ta.model.Focused() }

func (ta Textarea) Focus() Component {
	ta.model.Focus()
	return ta
}

func (ta Textarea) Blur() Component {
	ta.model.Blur()
	return ta
}

// Value returns the current textarea content.
func (ta Textarea) Value() string {
	return ta.model.Value()
}

// SetValue sets the textarea content.
func (ta *Textarea) SetValue(v string) {
	ta.model.SetValue(v)
}
