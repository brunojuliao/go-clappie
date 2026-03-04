package uikit

import (
	"fmt"

	"github.com/brunojuliao/go-clappie/internal/engine"
)

// RadioConfig configures a radio component.
type RadioConfig struct {
	Label    string
	Options  []string
	Selected int
	OnChange func(int, string)
}

// Radio is a single-select radio group component.
type Radio struct {
	ComponentBase
	config   RadioConfig
	selected int
}

// NewRadio creates a new radio group.
func NewRadio(cfg RadioConfig) *Radio {
	maxWidth := engine.VisualWidth(cfg.Label) + 4
	for _, opt := range cfg.Options {
		w := engine.VisualWidth(opt) + 6
		if w > maxWidth {
			maxWidth = w
		}
	}
	return &Radio{
		ComponentBase: ComponentBase{Focusable: true, Width: maxWidth},
		config:        cfg,
		selected:      cfg.Selected,
	}
}

// Render renders the radio group.
func (r *Radio) Render(focused bool) []string {
	var lines []string

	if r.config.Label != "" {
		label := r.config.Label
		if focused {
			label = engine.StyleBold(label)
		}
		lines = append(lines, "  "+label)
	}

	for i, opt := range r.config.Options {
		icon := "○"
		if i == r.selected {
			icon = "●"
		}
		line := fmt.Sprintf("    %s %s", icon, opt)
		if focused && i == r.selected {
			line = engine.StyleBold(line)
		}
		lines = append(lines, line)
	}

	return lines
}

// OnKey handles key events.
func (r *Radio) OnKey(key string) bool {
	switch key {
	case "UP":
		if r.selected > 0 {
			r.selected--
			r.notifyChange()
		}
		return true
	case "DOWN":
		if r.selected < len(r.config.Options)-1 {
			r.selected++
			r.notifyChange()
		}
		return true
	case "ENTER", "SPACE":
		r.notifyChange()
		return true
	}
	return false
}

// OnClick handles click events.
func (r *Radio) OnClick(lineIdx, col int) bool {
	// Adjust for label line
	optIdx := lineIdx
	if r.config.Label != "" {
		optIdx--
	}
	if optIdx >= 0 && optIdx < len(r.config.Options) {
		r.selected = optIdx
		r.notifyChange()
		return true
	}
	return false
}

func (r *Radio) notifyChange() {
	if r.config.OnChange != nil && r.selected < len(r.config.Options) {
		r.config.OnChange(r.selected, r.config.Options[r.selected])
	}
}

// Selected returns the selected index.
func (r *Radio) Selected() int {
	return r.selected
}

// SelectedOption returns the selected option string.
func (r *Radio) SelectedOption() string {
	if r.selected < len(r.config.Options) {
		return r.config.Options[r.selected]
	}
	return ""
}
