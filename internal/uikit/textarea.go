package uikit

import (
	"strings"

	"github.com/brunojuliao/go-clappie/internal/engine"
)

// TextareaConfig configures a textarea component.
type TextareaConfig struct {
	Placeholder string
	Width       int
	Height      int
	Value       string
	OnChange    func(string)
	OnSubmit    func(string)
}

// Textarea is a multi-line text input component.
type Textarea struct {
	ComponentBase
	config    TextareaConfig
	lines     []string
	cursorRow int
	cursorCol int
	scrollTop int
	height    int
}

// NewTextarea creates a new textarea.
func NewTextarea(cfg TextareaConfig) *Textarea {
	w := cfg.Width
	if w == 0 {
		w = 40
	}
	h := cfg.Height
	if h == 0 {
		h = 5
	}

	lines := []string{""}
	if cfg.Value != "" {
		lines = strings.Split(cfg.Value, "\n")
	}

	return &Textarea{
		ComponentBase: ComponentBase{Focusable: true, Width: w},
		config:        cfg,
		lines:         lines,
		height:        h,
	}
}

// Render renders the textarea.
func (ta *Textarea) Render(focused bool) []string {
	w := ta.GetWidth()
	innerWidth := w - 4

	border := "─"
	if focused {
		border = "━"
	}

	var output []string
	output = append(output, "┌"+strings.Repeat(border, w-2)+"┐")

	for i := 0; i < ta.height; i++ {
		lineIdx := ta.scrollTop + i
		var content string
		if lineIdx < len(ta.lines) {
			content = ta.lines[lineIdx]
		} else if i == 0 && len(ta.lines) == 1 && ta.lines[0] == "" && !focused {
			content = engine.StyleDim(ta.config.Placeholder)
		}

		// Truncate if needed
		stripped := engine.StripANSI(content)
		if engine.VisualWidth(stripped) > innerWidth {
			content = engine.TruncateToWidth(stripped, innerWidth, "")
		}

		padded := engine.PadRight(" "+content, w-2)
		output = append(output, "│"+padded+"│")
	}

	output = append(output, "└"+strings.Repeat(border, w-2)+"┘")
	return output
}

// OnKey handles key events for textarea.
func (ta *Textarea) OnKey(key string) bool {
	switch key {
	case "ENTER":
		// Check for Ctrl+Enter as submit
		// Regular enter inserts newline
		rest := ""
		if ta.cursorCol < len([]rune(ta.lines[ta.cursorRow])) {
			runes := []rune(ta.lines[ta.cursorRow])
			rest = string(runes[ta.cursorCol:])
			ta.lines[ta.cursorRow] = string(runes[:ta.cursorCol])
		}
		ta.cursorRow++
		ta.cursorCol = 0
		// Insert new line
		newLines := make([]string, len(ta.lines)+1)
		copy(newLines, ta.lines[:ta.cursorRow])
		newLines[ta.cursorRow] = rest
		copy(newLines[ta.cursorRow+1:], ta.lines[ta.cursorRow:])
		ta.lines = newLines
		ta.adjustScroll()
		ta.notifyChange()
		return true

	case "BACKSPACE":
		if ta.cursorCol > 0 {
			runes := []rune(ta.lines[ta.cursorRow])
			ta.lines[ta.cursorRow] = string(runes[:ta.cursorCol-1]) + string(runes[ta.cursorCol:])
			ta.cursorCol--
			ta.notifyChange()
		} else if ta.cursorRow > 0 {
			prevLen := len([]rune(ta.lines[ta.cursorRow-1]))
			ta.lines[ta.cursorRow-1] += ta.lines[ta.cursorRow]
			ta.lines = append(ta.lines[:ta.cursorRow], ta.lines[ta.cursorRow+1:]...)
			ta.cursorRow--
			ta.cursorCol = prevLen
			ta.adjustScroll()
			ta.notifyChange()
		}
		return true

	case "UP":
		if ta.cursorRow > 0 {
			ta.cursorRow--
			runes := []rune(ta.lines[ta.cursorRow])
			if ta.cursorCol > len(runes) {
				ta.cursorCol = len(runes)
			}
			ta.adjustScroll()
		}
		return true

	case "DOWN":
		if ta.cursorRow < len(ta.lines)-1 {
			ta.cursorRow++
			runes := []rune(ta.lines[ta.cursorRow])
			if ta.cursorCol > len(runes) {
				ta.cursorCol = len(runes)
			}
			ta.adjustScroll()
		}
		return true

	case "LEFT":
		if ta.cursorCol > 0 {
			ta.cursorCol--
		} else if ta.cursorRow > 0 {
			ta.cursorRow--
			ta.cursorCol = len([]rune(ta.lines[ta.cursorRow]))
			ta.adjustScroll()
		}
		return true

	case "RIGHT":
		runes := []rune(ta.lines[ta.cursorRow])
		if ta.cursorCol < len(runes) {
			ta.cursorCol++
		} else if ta.cursorRow < len(ta.lines)-1 {
			ta.cursorRow++
			ta.cursorCol = 0
			ta.adjustScroll()
		}
		return true

	case "HOME", "CTRL_A":
		ta.cursorCol = 0
		return true

	case "END", "CTRL_E":
		ta.cursorCol = len([]rune(ta.lines[ta.cursorRow]))
		return true

	default:
		// Printable character
		r := []rune(key)
		if len(r) == 1 && r[0] >= 32 {
			runes := []rune(ta.lines[ta.cursorRow])
			ta.lines[ta.cursorRow] = string(runes[:ta.cursorCol]) + key + string(runes[ta.cursorCol:])
			ta.cursorCol++
			ta.notifyChange()
			return true
		}
	}
	return false
}

func (ta *Textarea) adjustScroll() {
	if ta.cursorRow < ta.scrollTop {
		ta.scrollTop = ta.cursorRow
	}
	if ta.cursorRow >= ta.scrollTop+ta.height {
		ta.scrollTop = ta.cursorRow - ta.height + 1
	}
}

func (ta *Textarea) notifyChange() {
	if ta.config.OnChange != nil {
		ta.config.OnChange(ta.Value())
	}
}

// Value returns the current textarea content.
func (ta *Textarea) Value() string {
	return strings.Join(ta.lines, "\n")
}

// SetValue sets the textarea content.
func (ta *Textarea) SetValue(v string) {
	ta.lines = strings.Split(v, "\n")
	ta.cursorRow = 0
	ta.cursorCol = 0
	ta.scrollTop = 0
}
