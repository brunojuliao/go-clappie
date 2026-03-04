package uikit

import (
	"fmt"

	"github.com/brunojuliao/go-clappie/internal/engine"
)

// CheckboxConfig configures a checkbox component.
type CheckboxConfig struct {
	Label    string
	Checked  bool
	OnChange func(bool)
}

// Checkbox is a boolean checkbox component.
type Checkbox struct {
	ComponentBase
	config  CheckboxConfig
	checked bool
}

// NewCheckbox creates a new checkbox.
func NewCheckbox(cfg CheckboxConfig) *Checkbox {
	return &Checkbox{
		ComponentBase: ComponentBase{
			Focusable: true,
			Width:     engine.VisualWidth(cfg.Label) + 6,
		},
		config:  cfg,
		checked: cfg.Checked,
	}
}

// Render renders the checkbox.
func (c *Checkbox) Render(focused bool) []string {
	icon := "☐"
	if c.checked {
		icon = "☑"
	}

	line := fmt.Sprintf("  %s %s", icon, c.config.Label)
	if focused {
		line = engine.StyleBold(line)
	}
	return []string{line}
}

// OnKey handles key events.
func (c *Checkbox) OnKey(key string) bool {
	if key == "ENTER" || key == "SPACE" {
		c.checked = !c.checked
		if c.config.OnChange != nil {
			c.config.OnChange(c.checked)
		}
		return true
	}
	return false
}

// OnClick handles click events.
func (c *Checkbox) OnClick(lineIdx, col int) bool {
	c.checked = !c.checked
	if c.config.OnChange != nil {
		c.config.OnChange(c.checked)
	}
	return true
}

// Checked returns the current state.
func (c *Checkbox) Checked() bool {
	return c.checked
}

// SetChecked sets the checkbox state.
func (c *Checkbox) SetChecked(v bool) {
	c.checked = v
}
