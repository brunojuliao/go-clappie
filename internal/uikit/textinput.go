package uikit

import (
	"strings"

	"github.com/brunojuliao/go-clappie/internal/engine"
)

// TextInputConfig configures a text input component.
type TextInputConfig struct {
	Placeholder string
	Width       int
	Value       string
	OnChange    func(string)
	OnSubmit    func(string)
}

// TextInput is a single-line text input component.
type TextInput struct {
	ComponentBase
	config   TextInputConfig
	value    string
	cursor   int
}

// NewTextInput creates a new text input.
func NewTextInput(cfg TextInputConfig) *TextInput {
	w := cfg.Width
	if w == 0 {
		w = 30
	}
	return &TextInput{
		ComponentBase: ComponentBase{Focusable: true, Width: w},
		config:        cfg,
		value:         cfg.Value,
		cursor:        len(cfg.Value),
	}
}

// Render renders the text input.
func (ti *TextInput) Render(focused bool) []string {
	w := ti.GetWidth()
	content := ti.value
	if content == "" && !focused {
		content = engine.StyleDim(ti.config.Placeholder)
	}

	// Build the input line
	border := "─"
	if focused {
		border = "━"
	}

	topBorder := "┌" + strings.Repeat(border, w-2) + "┐"
	botBorder := "└" + strings.Repeat(border, w-2) + "┘"

	// Truncate or pad content
	displayContent := content
	visWidth := engine.VisualWidth(engine.StripANSI(displayContent))
	innerWidth := w - 4 // 2 for borders, 2 for padding
	if visWidth > innerWidth {
		displayContent = engine.TruncateToWidth(engine.StripANSI(displayContent), innerWidth, "")
	}

	padded := engine.PadRight(" "+displayContent, w-2)
	mid := "│" + padded + "│"

	if focused {
		// Show cursor
		mid = "│" + padded + "│"
	}

	return []string{topBorder, mid, botBorder}
}

// OnKey handles key events for text input.
func (ti *TextInput) OnKey(key string) bool {
	switch key {
	case "ENTER":
		if ti.config.OnSubmit != nil {
			ti.config.OnSubmit(ti.value)
		}
		return true
	case "BACKSPACE":
		if ti.cursor > 0 {
			runes := []rune(ti.value)
			ti.value = string(runes[:ti.cursor-1]) + string(runes[ti.cursor:])
			ti.cursor--
			ti.notifyChange()
		}
		return true
	case "DELETE":
		runes := []rune(ti.value)
		if ti.cursor < len(runes) {
			ti.value = string(runes[:ti.cursor]) + string(runes[ti.cursor+1:])
			ti.notifyChange()
		}
		return true
	case "LEFT":
		if ti.cursor > 0 {
			ti.cursor--
		}
		return true
	case "RIGHT":
		if ti.cursor < len([]rune(ti.value)) {
			ti.cursor++
		}
		return true
	case "HOME", "CTRL_A":
		ti.cursor = 0
		return true
	case "END", "CTRL_E":
		ti.cursor = len([]rune(ti.value))
		return true
	case "CTRL_U":
		ti.value = string([]rune(ti.value)[ti.cursor:])
		ti.cursor = 0
		ti.notifyChange()
		return true
	case "CTRL_K":
		ti.value = string([]rune(ti.value)[:ti.cursor])
		ti.notifyChange()
		return true
	default:
		// Printable character
		if len(key) == 1 && key[0] >= 32 && key[0] <= 126 {
			runes := []rune(ti.value)
			ti.value = string(runes[:ti.cursor]) + key + string(runes[ti.cursor:])
			ti.cursor++
			ti.notifyChange()
			return true
		}
		// Multi-byte UTF-8
		r := []rune(key)
		if len(r) == 1 {
			runes := []rune(ti.value)
			ti.value = string(runes[:ti.cursor]) + key + string(runes[ti.cursor:])
			ti.cursor++
			ti.notifyChange()
			return true
		}
	}
	return false
}

func (ti *TextInput) notifyChange() {
	if ti.config.OnChange != nil {
		ti.config.OnChange(ti.value)
	}
}

// Value returns the current input value.
func (ti *TextInput) Value() string {
	return ti.value
}

// SetValue sets the input value.
func (ti *TextInput) SetValue(v string) {
	ti.value = v
	ti.cursor = len([]rune(v))
}

// Clear clears the input value.
func (ti *TextInput) Clear() {
	ti.value = ""
	ti.cursor = 0
}
