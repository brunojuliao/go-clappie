package uikit

import (
	"fmt"
	"strings"

	"github.com/brunojuliao/go-clappie/internal/engine"
)

// ButtonStyle represents different button visual styles.
type ButtonStyle int

const (
	ButtonStyleDefault ButtonStyle = iota
	ButtonStyleFilled
	ButtonStyleGhost
	ButtonStyleInline
	ButtonStyleFullWidth
)

// ButtonConfig configures a button component.
type ButtonConfig struct {
	Label    string
	Shortcut string
	OnPress  func()
	Style    ButtonStyle
	Width    int
}

// Button is a clickable button component.
type Button struct {
	ComponentBase
	config ButtonConfig
}

// NewButton creates a new button.
func NewButton(cfg ButtonConfig) *Button {
	w := cfg.Width
	if w == 0 {
		w = engine.VisualWidth(cfg.Label) + 6
	}
	return &Button{
		ComponentBase: ComponentBase{Focusable: true, Width: w},
		config:        cfg,
	}
}

// NewButtonFilled creates a filled-style button.
func NewButtonFilled(cfg ButtonConfig) *Button {
	cfg.Style = ButtonStyleFilled
	return NewButton(cfg)
}

// NewButtonGhost creates a ghost-style button.
func NewButtonGhost(cfg ButtonConfig) *Button {
	cfg.Style = ButtonStyleGhost
	return NewButton(cfg)
}

// NewButtonInline creates an inline-style button.
func NewButtonInline(cfg ButtonConfig) *Button {
	cfg.Style = ButtonStyleInline
	return NewButton(cfg)
}

// NewButtonFullWidth creates a full-width button.
func NewButtonFullWidth(cfg ButtonConfig) *Button {
	cfg.Style = ButtonStyleFullWidth
	return NewButton(cfg)
}

// Render renders the button.
func (b *Button) Render(focused bool) []string {
	label := b.config.Label
	if b.config.Shortcut != "" {
		label = fmt.Sprintf("[%s] %s", b.config.Shortcut, label)
	}

	switch b.config.Style {
	case ButtonStyleFilled:
		return b.renderFilled(label, focused)
	case ButtonStyleGhost:
		return b.renderGhost(label, focused)
	case ButtonStyleInline:
		return b.renderInline(label, focused)
	case ButtonStyleFullWidth:
		return b.renderFullWidth(label, focused)
	default:
		return b.renderDefault(label, focused)
	}
}

func (b *Button) renderDefault(label string, focused bool) []string {
	w := b.GetWidth()
	padded := engine.PadCenter(label, w-2)

	if focused {
		top := "┌" + strings.Repeat("─", w-2) + "┐"
		mid := "│" + engine.StyleBold(padded) + "│"
		bot := "└" + strings.Repeat("─", w-2) + "┘"
		return []string{top, mid, bot}
	}

	top := "┌" + strings.Repeat("─", w-2) + "┐"
	mid := "│" + padded + "│"
	bot := "└" + strings.Repeat("─", w-2) + "┘"
	return []string{top, mid, bot}
}

func (b *Button) renderFilled(label string, focused bool) []string {
	w := b.GetWidth()
	padded := engine.PadCenter(label, w)
	if focused {
		return []string{engine.StyleInverse(engine.StyleBold(padded))}
	}
	return []string{engine.StyleInverse(padded)}
}

func (b *Button) renderGhost(label string, focused bool) []string {
	if focused {
		return []string{engine.StyleBold("  " + label + "  ")}
	}
	return []string{engine.StyleDim("  " + label + "  ")}
}

func (b *Button) renderInline(label string, focused bool) []string {
	if focused {
		return []string{engine.StyleUnderline(label)}
	}
	return []string{label}
}

func (b *Button) renderFullWidth(label string, focused bool) []string {
	w := b.GetWidth()
	padded := engine.PadRight("  "+label, w)
	if focused {
		return []string{engine.StyleInverse(engine.StyleBold(padded))}
	}
	return []string{engine.StyleInverse(padded)}
}

// OnKey handles key events for the button.
func (b *Button) OnKey(key string) bool {
	if key == "ENTER" || key == "SPACE" {
		if b.config.OnPress != nil {
			b.config.OnPress()
		}
		return true
	}
	return false
}

// OnClick handles click events.
func (b *Button) OnClick(lineIdx, col int) bool {
	if b.config.OnPress != nil {
		b.config.OnPress()
	}
	return true
}

// GetShortcut returns the button's shortcut key.
func (b *Button) GetShortcut() string {
	return b.config.Shortcut
}
