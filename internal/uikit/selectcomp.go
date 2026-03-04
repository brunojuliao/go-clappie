package uikit

import (
	"fmt"

	"github.com/brunojuliao/go-clappie/internal/engine"
)

// SelectStyle represents select visual styles.
type SelectStyle int

const (
	SelectStyleDefault SelectStyle = iota
	SelectStyleBlock
)

// SelectConfig configures a select component.
type SelectConfig struct {
	Label    string
	Options  []string
	Selected int
	OnChange func(int, string)
	Style    SelectStyle
	Width    int
}

// Select is a dropdown-like select component.
type Select struct {
	ComponentBase
	config   SelectConfig
	selected int
	expanded bool
}

// NewSelect creates a new select.
func NewSelect(cfg SelectConfig) *Select {
	w := cfg.Width
	if w == 0 {
		w = 30
	}
	return &Select{
		ComponentBase: ComponentBase{Focusable: true, Width: w},
		config:        cfg,
		selected:      cfg.Selected,
	}
}

// NewSelectBlock creates a block-style select.
func NewSelectBlock(cfg SelectConfig) *Select {
	cfg.Style = SelectStyleBlock
	return NewSelect(cfg)
}

// Render renders the select.
func (s *Select) Render(focused bool) []string {
	currentValue := ""
	if s.selected >= 0 && s.selected < len(s.config.Options) {
		currentValue = s.config.Options[s.selected]
	}

	if !s.expanded {
		label := s.config.Label
		if label != "" {
			label += ": "
		}
		indicator := "▸"
		if focused {
			indicator = "▾"
		}
		line := fmt.Sprintf("  %s%s %s", label, currentValue, indicator)

		if s.config.Style == SelectStyleBlock {
			line = engine.PadRight(line, s.GetWidth())
			if focused {
				line = engine.StyleInverse(line)
			}
		} else if focused {
			line = engine.StyleBold(line)
		}
		return []string{line}
	}

	// Expanded: show all options
	var lines []string
	if s.config.Label != "" {
		lines = append(lines, "  "+engine.StyleBold(s.config.Label))
	}
	for i, opt := range s.config.Options {
		icon := " "
		if i == s.selected {
			icon = "▸"
		}
		line := fmt.Sprintf("    %s %s", icon, opt)
		if i == s.selected {
			line = engine.StyleBold(line)
		}
		if s.config.Style == SelectStyleBlock {
			line = engine.PadRight(line, s.GetWidth())
			if i == s.selected {
				line = engine.StyleInverse(line)
			}
		}
		lines = append(lines, line)
	}
	return lines
}

// OnKey handles key events.
func (s *Select) OnKey(key string) bool {
	if !s.expanded {
		if key == "ENTER" || key == "SPACE" {
			s.expanded = true
			return true
		}
		return false
	}

	switch key {
	case "UP":
		if s.selected > 0 {
			s.selected--
		}
		return true
	case "DOWN":
		if s.selected < len(s.config.Options)-1 {
			s.selected++
		}
		return true
	case "ENTER", "SPACE":
		s.expanded = false
		s.notifyChange()
		return true
	case "ESC":
		s.expanded = false
		return true
	}
	return false
}

// OnClick handles click events.
func (s *Select) OnClick(lineIdx, col int) bool {
	if !s.expanded {
		s.expanded = true
		return true
	}
	optIdx := lineIdx
	if s.config.Label != "" {
		optIdx--
	}
	if optIdx >= 0 && optIdx < len(s.config.Options) {
		s.selected = optIdx
		s.expanded = false
		s.notifyChange()
		return true
	}
	return false
}

func (s *Select) notifyChange() {
	if s.config.OnChange != nil && s.selected < len(s.config.Options) {
		s.config.OnChange(s.selected, s.config.Options[s.selected])
	}
}

// Selected returns the selected index.
func (s *Select) Selected() int {
	return s.selected
}

// SelectedOption returns the selected option string.
func (s *Select) SelectedOption() string {
	if s.selected >= 0 && s.selected < len(s.config.Options) {
		return s.config.Options[s.selected]
	}
	return ""
}
