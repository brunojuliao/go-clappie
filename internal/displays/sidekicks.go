package displays

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/brunojuliao/go-clappie/internal/engine"
	"github.com/brunojuliao/go-clappie/internal/platform"
	"github.com/brunojuliao/go-clappie/internal/sidekicks"
)

type sidekicksScreen struct {
	active      []sidekicks.SidekickInfo
	selectedIdx int
	styles      *engine.Styles
}

func NewSidekicksScreen(data map[string]interface{}, styles *engine.Styles, claudePane string) engine.ScreenModel {
	m := &sidekicksScreen{styles: styles}
	m.load()
	return m
}

func (m *sidekicksScreen) load() {
	root, err := platform.ProjectRoot()
	if err != nil {
		return
	}
	m.active, _ = sidekicks.ListActive(root)
}

func (m *sidekicksScreen) Init() tea.Cmd { return nil }

func (m *sidekicksScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "up", "k":
			if m.selectedIdx > 0 {
				m.selectedIdx--
			}
		case "down", "j":
			if m.selectedIdx < len(m.active)-1 {
				m.selectedIdx++
			}
		case "r", "R":
			m.load()
		case "K":
			if m.selectedIdx < len(m.active) {
				root, _ := platform.ProjectRoot()
				sidekicks.End(root)
				m.load()
			}
		}
	}
	return m, nil
}

func (m *sidekicksScreen) View() string {
	dim := lipgloss.NewStyle().Faint(true)
	bold := lipgloss.NewStyle().Bold(true)

	var lines []string
	lines = append(lines, "")

	if len(m.active) == 0 {
		lines = append(lines, "  No active sidekicks.")
		lines = append(lines, "")
		lines = append(lines, dim.Render("  Use 'go-clappie sidekick spawn \"prompt\"' to spawn one."))
	} else {
		lines = append(lines, fmt.Sprintf("  %d active sidekicks", len(m.active)))
		lines = append(lines, "")
		for i, sk := range m.active {
			prefix := "  "
			if i == m.selectedIdx {
				prefix = "▸ "
			}
			line := fmt.Sprintf("%s%s: %s", prefix, sk.ID, sk.Prompt)
			if len(line) > 65 {
				line = line[:62] + "..."
			}
			if i == m.selectedIdx {
				line = bold.Render(line)
			}
			lines = append(lines, line)
		}
	}

	return strings.Join(lines, "\n")
}

func (m *sidekicksScreen) Name() string          { return "Sidekicks" }
func (m *sidekicksScreen) Layout() (string, int) { return "centered", 70 }
func (m *sidekicksScreen) Shortcuts() []engine.ShortcutHint {
	return []engine.ShortcutHint{
		{Key: "R", Label: "Refresh"},
		{Key: "K", Label: "Kill Selected"},
	}
}
