package displays

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/brunojuliao/go-clappie/internal/engine"
)

type utilityListScreen struct {
	title       string
	options     []string
	selectedIdx int
	styles      *engine.Styles
}

func NewUtilityListScreen(data map[string]interface{}, styles *engine.Styles, claudePane string) engine.ScreenModel {
	title := "Select"
	if data != nil {
		if t, ok := data["title"].(string); ok {
			title = t
		}
	}
	var options []string
	if data != nil {
		if opts, ok := data["options"].([]interface{}); ok {
			for _, o := range opts {
				if s, ok := o.(string); ok {
					options = append(options, s)
				}
			}
		}
	}
	return &utilityListScreen{title: title, options: options, styles: styles}
}

func (m *utilityListScreen) Init() tea.Cmd { return nil }

func (m *utilityListScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "up", "k":
			if m.selectedIdx > 0 {
				m.selectedIdx--
			}
		case "down", "j":
			if m.selectedIdx < len(m.options)-1 {
				m.selectedIdx++
			}
		case "enter":
			if m.selectedIdx < len(m.options) {
				return m, tea.Batch(
					engine.SubmitToClaudeCmd(fmt.Sprintf("[go-clappie] List → %s", m.options[m.selectedIdx])),
					engine.PopViewCmd(),
				)
			}
		}
	}
	return m, nil
}

func (m *utilityListScreen) View() string {
	bold := lipgloss.NewStyle().Bold(true)

	var lines []string
	lines = append(lines, "")
	for i, opt := range m.options {
		prefix := "  "
		if i == m.selectedIdx {
			prefix = "▸ "
		}
		line := fmt.Sprintf("%s%s", prefix, opt)
		if i == m.selectedIdx {
			line = bold.Render(line)
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func (m *utilityListScreen) Name() string                        { return m.title }
func (m *utilityListScreen) Layout() (string, int)               { return "centered", 50 }
func (m *utilityListScreen) Shortcuts() []engine.ShortcutHint { return nil }
