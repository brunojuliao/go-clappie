package uikit

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// RadioConfig configures a radio component.
type RadioConfig struct {
	Label    string
	Options  []string
	Selected int
	OnChange func(int, string) tea.Cmd
}

// Radio is a single-select radio group component.
type Radio struct {
	config   RadioConfig
	selected int
	focused  bool
}

// NewRadio creates a new radio group.
func NewRadio(cfg RadioConfig) Radio {
	return Radio{config: cfg, selected: cfg.Selected}
}

func (r Radio) Init() tea.Cmd { return nil }

func (r Radio) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "up", "k":
			if r.selected > 0 {
				r.selected--
				return r, r.notifyChange()
			}
			return r, nil
		case "down", "j":
			if r.selected < len(r.config.Options)-1 {
				r.selected++
				return r, r.notifyChange()
			}
			return r, nil
		case "enter", " ":
			return r, r.notifyChange()
		}
	}
	return r, nil
}

func (r Radio) View() string {
	var lines []string

	if r.config.Label != "" {
		label := r.config.Label
		if r.focused {
			label = lipgloss.NewStyle().Bold(true).Render(label)
		}
		lines = append(lines, "  "+label)
	}

	for i, opt := range r.config.Options {
		icon := "○"
		if i == r.selected {
			icon = "●"
		}
		line := fmt.Sprintf("    %s %s", icon, opt)
		if r.focused && i == r.selected {
			line = lipgloss.NewStyle().Bold(true).Render(line)
		}
		lines = append(lines, line)
	}

	return joinLines(lines)
}

func (r Radio) notifyChange() tea.Cmd {
	if r.config.OnChange != nil && r.selected < len(r.config.Options) {
		return r.config.OnChange(r.selected, r.config.Options[r.selected])
	}
	return nil
}

func (r Radio) IsFocusable() bool { return true }
func (r Radio) Focused() bool     { return r.focused }

func (r Radio) Focus() Component {
	r.focused = true
	return r
}

func (r Radio) Blur() Component {
	r.focused = false
	return r
}

// Selected returns the selected index.
func (r Radio) Selected() int { return r.selected }

// SelectedOption returns the selected option string.
func (r Radio) SelectedOption() string {
	if r.selected < len(r.config.Options) {
		return r.config.Options[r.selected]
	}
	return ""
}
