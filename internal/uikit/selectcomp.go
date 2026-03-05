package uikit

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	OnChange func(int, string) tea.Cmd
	Style    SelectStyle
	Width    int
}

// Select is a dropdown-like select component.
type Select struct {
	config   SelectConfig
	selected int
	expanded bool
	focused  bool
	width    int
}

// NewSelect creates a new select.
func NewSelect(cfg SelectConfig) Select {
	w := cfg.Width
	if w == 0 {
		w = 30
	}
	return Select{config: cfg, selected: cfg.Selected, width: w}
}

// NewSelectBlock creates a block-style select.
func NewSelectBlock(cfg SelectConfig) Select {
	cfg.Style = SelectStyleBlock
	return NewSelect(cfg)
}

func (s Select) Init() tea.Cmd { return nil }

func (s Select) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if !s.expanded {
			switch keyMsg.String() {
			case "enter", " ":
				s.expanded = true
				return s, nil
			}
			return s, nil
		}

		switch keyMsg.String() {
		case "up", "k":
			if s.selected > 0 {
				s.selected--
			}
			return s, nil
		case "down", "j":
			if s.selected < len(s.config.Options)-1 {
				s.selected++
			}
			return s, nil
		case "enter", " ":
			s.expanded = false
			return s, s.notifyChange()
		case "esc":
			s.expanded = false
			return s, nil
		}
	}
	return s, nil
}

func (s Select) View() string {
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
		if s.focused {
			indicator = "▾"
		}
		line := fmt.Sprintf("  %s%s %s", label, currentValue, indicator)

		if s.config.Style == SelectStyleBlock {
			style := lipgloss.NewStyle().Width(s.width)
			if s.focused {
				style = style.Reverse(true)
			}
			return style.Render(line)
		}
		if s.focused {
			return lipgloss.NewStyle().Bold(true).Render(line)
		}
		return line
	}

	// Expanded: show all options
	var lines []string
	if s.config.Label != "" {
		lines = append(lines, "  "+lipgloss.NewStyle().Bold(true).Render(s.config.Label))
	}
	for i, opt := range s.config.Options {
		icon := " "
		if i == s.selected {
			icon = "▸"
		}
		line := fmt.Sprintf("    %s %s", icon, opt)
		if i == s.selected {
			style := lipgloss.NewStyle().Bold(true)
			if s.config.Style == SelectStyleBlock {
				style = style.Width(s.width).Reverse(true)
			}
			line = style.Render(line)
		}
		lines = append(lines, line)
	}
	return joinLines(lines)
}

func (s Select) notifyChange() tea.Cmd {
	if s.config.OnChange != nil && s.selected < len(s.config.Options) {
		return s.config.OnChange(s.selected, s.config.Options[s.selected])
	}
	return nil
}

func (s Select) IsFocusable() bool { return true }
func (s Select) Focused() bool     { return s.focused }

func (s Select) Focus() Component {
	s.focused = true
	return s
}

func (s Select) Blur() Component {
	s.focused = false
	return s
}

// Selected returns the selected index.
func (s Select) Selected() int { return s.selected }

// SelectedOption returns the selected option string.
func (s Select) SelectedOption() string {
	if s.selected >= 0 && s.selected < len(s.config.Options) {
		return s.config.Options[s.selected]
	}
	return ""
}
