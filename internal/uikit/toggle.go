package uikit

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ToggleStyle represents toggle visual styles.
type ToggleStyle int

const (
	ToggleStyleDefault ToggleStyle = iota
	ToggleStyleBlock
)

// ToggleConfig configures a toggle component.
type ToggleConfig struct {
	Label    string
	Shortcut string
	Value    bool
	OnChange func(bool) tea.Cmd
	Style    ToggleStyle
	Width    int
}

// Toggle is a binary on/off switch component.
type Toggle struct {
	config  ToggleConfig
	value   bool
	focused bool
	width   int
}

// NewToggle creates a new toggle.
func NewToggle(cfg ToggleConfig) Toggle {
	w := cfg.Width
	if w == 0 {
		w = len(cfg.Label) + 10
	}
	return Toggle{config: cfg, value: cfg.Value, width: w}
}

// NewToggleBlock creates a block-style toggle.
func NewToggleBlock(cfg ToggleConfig) Toggle {
	cfg.Style = ToggleStyleBlock
	return NewToggle(cfg)
}

func (t Toggle) Init() tea.Cmd { return nil }

func (t Toggle) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "enter", " ":
			t.value = !t.value
			if t.config.OnChange != nil {
				return t, t.config.OnChange(t.value)
			}
			return t, nil
		}
	}
	return t, nil
}

func (t Toggle) View() string {
	indicator := "○"
	if t.value {
		indicator = "●"
	}

	label := t.config.Label
	if t.config.Shortcut != "" {
		label = fmt.Sprintf("[%s] %s", t.config.Shortcut, label)
	}

	line := fmt.Sprintf("  %s %s", indicator, label)

	if t.config.Style == ToggleStyleBlock {
		style := lipgloss.NewStyle().Width(t.width)
		if t.focused {
			style = style.Bold(true).Reverse(true)
		}
		return style.Render(line)
	}

	if t.focused {
		return lipgloss.NewStyle().Bold(true).Render(line)
	}
	return line
}

func (t Toggle) IsFocusable() bool { return true }
func (t Toggle) Focused() bool     { return t.focused }

func (t Toggle) Focus() Component {
	t.focused = true
	return t
}

func (t Toggle) Blur() Component {
	t.focused = false
	return t
}

// Value returns the current toggle state.
func (t Toggle) Value() bool {
	return t.value
}

// SetValue sets the toggle state.
func (t *Toggle) SetValue(v bool) {
	t.value = v
}
