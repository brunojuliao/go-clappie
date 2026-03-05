package uikit

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

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
	config DividerConfig
	width  int
}

// NewDivider creates a new divider.
func NewDivider(cfg DividerConfig) Divider {
	w := cfg.Width
	if w == 0 {
		w = 40
	}
	return Divider{config: cfg, width: w}
}

func (d Divider) Init() tea.Cmd                         { return nil }
func (d Divider) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return d, nil }

func (d Divider) View() string {
	w := d.width - 4

	var ch string
	switch d.config.Variant {
	case DividerThick:
		ch = "━"
	case DividerDashed:
		ch = "╌"
	case DividerDotted:
		ch = "·"
	case DividerSpace:
		return ""
	default:
		ch = "─"
	}

	return "  " + strings.Repeat(ch, w)
}

func (d Divider) IsFocusable() bool { return false }
func (d Divider) Focused() bool     { return false }
func (d Divider) Focus() Component  { return d }
func (d Divider) Blur() Component   { return d }
