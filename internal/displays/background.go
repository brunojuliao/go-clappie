package displays

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/brunojuliao/go-clappie/internal/background"
	"github.com/brunojuliao/go-clappie/internal/engine"
	"github.com/brunojuliao/go-clappie/internal/platform"
)

type backgroundScreen struct {
	apps        []background.App
	selectedIdx int
	styles      *engine.Styles
}

func NewBackgroundScreen(data map[string]interface{}, styles *engine.Styles, claudePane string) engine.ScreenModel {
	m := &backgroundScreen{styles: styles}
	m.load()
	return m
}

func (m *backgroundScreen) load() {
	root, err := platform.ProjectRoot()
	if err != nil {
		return
	}
	m.apps, _ = background.List(root)
}

func (m *backgroundScreen) Init() tea.Cmd { return nil }

func (m *backgroundScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "up", "k":
			if m.selectedIdx > 0 {
				m.selectedIdx--
			}
		case "down", "j":
			if m.selectedIdx < len(m.apps)-1 {
				m.selectedIdx++
			}
		case "s", "S":
			if m.selectedIdx < len(m.apps) {
				root, _ := platform.ProjectRoot()
				app := m.apps[m.selectedIdx]
				if app.Running {
					background.Stop(app.Name)
				} else {
					background.Start(root, app.Name)
				}
				m.load()
			}
		case "r", "R":
			m.load()
		}
	}
	return m, nil
}

func (m *backgroundScreen) View() string {
	dim := lipgloss.NewStyle().Faint(true)
	bold := lipgloss.NewStyle().Bold(true)
	green := lipgloss.NewStyle().Foreground(lipgloss.Color("#28a745"))

	var lines []string
	lines = append(lines, "")

	if len(m.apps) == 0 {
		lines = append(lines, "  No background apps found.")
	} else {
		lines = append(lines, fmt.Sprintf("  %d background apps", len(m.apps)))
		lines = append(lines, "")
		for i, app := range m.apps {
			prefix := "  "
			if i == m.selectedIdx {
				prefix = "▸ "
			}
			status := dim.Render("stopped")
			if app.Running {
				status = green.Render("running")
			}
			line := fmt.Sprintf("%s%s [%s]", prefix, app.Name, status)
			if i == m.selectedIdx {
				line = bold.Render(line)
			}
			lines = append(lines, line)
		}
	}

	return strings.Join(lines, "\n")
}

func (m *backgroundScreen) Name() string          { return "Background" }
func (m *backgroundScreen) Layout() (string, int) { return "centered", 60 }
func (m *backgroundScreen) Shortcuts() []engine.ShortcutHint {
	return []engine.ShortcutHint{
		{Key: "S", Label: "Start/Stop"},
		{Key: "R", Label: "Refresh"},
	}
}
