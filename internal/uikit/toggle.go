package uikit

import (
	"fmt"

	"github.com/brunojuliao/go-clappie/internal/engine"
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
	OnChange func(bool)
	Style    ToggleStyle
	Width    int
}

// Toggle is a binary on/off switch component.
type Toggle struct {
	ComponentBase
	config ToggleConfig
	value  bool
}

// NewToggle creates a new toggle.
func NewToggle(cfg ToggleConfig) *Toggle {
	w := cfg.Width
	if w == 0 {
		w = engine.VisualWidth(cfg.Label) + 10
	}
	return &Toggle{
		ComponentBase: ComponentBase{Focusable: true, Width: w},
		config:        cfg,
		value:         cfg.Value,
	}
}

// NewToggleBlock creates a block-style toggle.
func NewToggleBlock(cfg ToggleConfig) *Toggle {
	cfg.Style = ToggleStyleBlock
	return NewToggle(cfg)
}

// Render renders the toggle.
func (t *Toggle) Render(focused bool) []string {
	indicator := "○"
	if t.value {
		indicator = "●"
	}

	label := t.config.Label
	if t.config.Shortcut != "" {
		label = fmt.Sprintf("[%s] %s", t.config.Shortcut, label)
	}

	line := fmt.Sprintf("  %s %s", indicator, label)

	if focused {
		line = engine.StyleBold(line)
	}

	if t.config.Style == ToggleStyleBlock {
		w := t.GetWidth()
		line = engine.PadRight(line, w)
		if focused {
			line = engine.StyleInverse(line)
		}
	}

	return []string{line}
}

// OnKey handles key events.
func (t *Toggle) OnKey(key string) bool {
	if key == "ENTER" || key == "SPACE" {
		t.value = !t.value
		if t.config.OnChange != nil {
			t.config.OnChange(t.value)
		}
		return true
	}
	return false
}

// OnClick handles click events.
func (t *Toggle) OnClick(lineIdx, col int) bool {
	t.value = !t.value
	if t.config.OnChange != nil {
		t.config.OnChange(t.value)
	}
	return true
}

// Value returns the current toggle state.
func (t *Toggle) Value() bool {
	return t.value
}

// SetValue sets the toggle state.
func (t *Toggle) SetValue(v bool) {
	t.value = v
}
