package uikit

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// CheckboxConfig configures a checkbox component.
type CheckboxConfig struct {
	Label    string
	Checked  bool
	OnChange func(bool) tea.Cmd
}

// Checkbox is a boolean checkbox component.
type Checkbox struct {
	config  CheckboxConfig
	checked bool
	focused bool
}

// NewCheckbox creates a new checkbox.
func NewCheckbox(cfg CheckboxConfig) Checkbox {
	return Checkbox{config: cfg, checked: cfg.Checked}
}

func (c Checkbox) Init() tea.Cmd { return nil }

func (c Checkbox) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "enter", " ":
			c.checked = !c.checked
			if c.config.OnChange != nil {
				return c, c.config.OnChange(c.checked)
			}
			return c, nil
		}
	}
	return c, nil
}

func (c Checkbox) View() string {
	icon := "☐"
	if c.checked {
		icon = "☑"
	}
	line := fmt.Sprintf("  %s %s", icon, c.config.Label)
	if c.focused {
		return lipgloss.NewStyle().Bold(true).Render(line)
	}
	return line
}

func (c Checkbox) IsFocusable() bool { return true }
func (c Checkbox) Focused() bool     { return c.focused }

func (c Checkbox) Focus() Component {
	c.focused = true
	return c
}

func (c Checkbox) Blur() Component {
	c.focused = false
	return c
}

// Checked returns the current state.
func (c Checkbox) Checked() bool { return c.checked }

// SetChecked sets the checkbox state.
func (c *Checkbox) SetChecked(v bool) { c.checked = v }
