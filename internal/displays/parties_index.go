package displays

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/brunojuliao/go-clappie/internal/engine"
	"github.com/brunojuliao/go-clappie/internal/parties"
	"github.com/brunojuliao/go-clappie/internal/platform"
)

type partiesIndexScreen struct {
	games       []parties.GameInfo
	selectedIdx int
	styles      *engine.Styles
}

func NewPartiesIndexScreen(data map[string]interface{}, styles *engine.Styles, claudePane string) engine.ScreenModel {
	m := &partiesIndexScreen{styles: styles}
	m.load()
	return m
}

func (m *partiesIndexScreen) load() {
	root, err := platform.ProjectRoot()
	if err != nil {
		return
	}
	m.games, _ = parties.ListGames(root)
}

func (m *partiesIndexScreen) Init() tea.Cmd { return nil }

func (m *partiesIndexScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "up", "k":
			if m.selectedIdx > 0 {
				m.selectedIdx--
			}
		case "down", "j":
			if m.selectedIdx < len(m.games)-1 {
				m.selectedIdx++
			}
		case "enter":
			if m.selectedIdx < len(m.games) {
				return m, engine.PushViewCmd("parties/status", map[string]interface{}{
					"game": m.games[m.selectedIdx].Name,
				})
			}
		case "i", "I":
			if m.selectedIdx < len(m.games) {
				root, _ := platform.ProjectRoot()
				simID, err := parties.Init(root, m.games[m.selectedIdx].Name)
				if err != nil {
					return m, engine.ToastCmd(fmt.Sprintf("Error: %v", err), 0)
				}
				return m, engine.ToastCmd(fmt.Sprintf("Initialized: %s", simID), 0)
			}
		}
	}
	return m, nil
}

func (m *partiesIndexScreen) View() string {
	dim := lipgloss.NewStyle().Faint(true)
	bold := lipgloss.NewStyle().Bold(true)

	var lines []string
	lines = append(lines, "")

	if len(m.games) == 0 {
		lines = append(lines, "  No games found.")
		lines = append(lines, "")
		lines = append(lines, dim.Render("  Create game files to get started."))
	} else {
		lines = append(lines, fmt.Sprintf("  %d games available", len(m.games)))
		lines = append(lines, "")
		for i, g := range m.games {
			prefix := "  "
			if i == m.selectedIdx {
				prefix = "▸ "
			}
			line := fmt.Sprintf("%s%s", prefix, g.Name)
			if g.Description != "" {
				line += dim.Render(" — " + g.Description)
			}
			if i == m.selectedIdx {
				line = bold.Render(line)
			}
			lines = append(lines, line)
		}
	}

	return strings.Join(lines, "\n")
}

func (m *partiesIndexScreen) Name() string          { return "Parties" }
func (m *partiesIndexScreen) Layout() (string, int) { return "centered", 60 }
func (m *partiesIndexScreen) Shortcuts() []engine.ShortcutHint {
	return []engine.ShortcutHint{{Key: "I", Label: "Init"}}
}
