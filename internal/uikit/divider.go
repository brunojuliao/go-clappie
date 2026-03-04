package uikit

import "strings"

// DividerVariant represents the divider style.
type DividerVariant int

const (
	DividerThin DividerVariant = iota
	DividerThick
	DividerDashed
	DividerDotted
	DividerSpace
)

// DividerConfig configures a divider component.
type DividerConfig struct {
	Variant DividerVariant
	Width   int
}

// Divider is a horizontal line separator component.
type Divider struct {
	ComponentBase
	config DividerConfig
}

// NewDivider creates a new divider.
func NewDivider(cfg DividerConfig) *Divider {
	w := cfg.Width
	if w == 0 {
		w = 40
	}
	return &Divider{
		ComponentBase: ComponentBase{Focusable: false, Width: w},
		config:        cfg,
	}
}

// Render renders the divider.
func (d *Divider) Render(focused bool) []string {
	w := d.GetWidth() - 4 // padding

	var char string
	switch d.config.Variant {
	case DividerThick:
		char = "━"
	case DividerDashed:
		char = "╌"
	case DividerDotted:
		char = "·"
	case DividerSpace:
		return []string{""}
	default:
		char = "─"
	}

	return []string{"  " + strings.Repeat(char, w)}
}
