package displays

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/brunojuliao/go-clappie/internal/engine"
)

type partiesStatusScreen struct {
	gameName string
	styles   *engine.Styles
}

func NewPartiesStatusScreen(data map[string]interface{}, styles *engine.Styles, claudePane string) engine.ScreenModel {
	gameName := ""
	if data != nil {
		if g, ok := data["game"].(string); ok {
			gameName = g
		}
	}
	return &partiesStatusScreen{gameName: gameName, styles: styles}
}

func (m *partiesStatusScreen) Init() tea.Cmd { return nil }

func (m *partiesStatusScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if keyMsg.String() == "l" || keyMsg.String() == "L" {
			return m, engine.ToastCmd("Launching simulation...", 0)
		}
	}
	return m, nil
}

func (m *partiesStatusScreen) View() string {
	bold := lipgloss.NewStyle().Bold(true)
	dim := lipgloss.NewStyle().Faint(true)

	lines := []string{
		"",
		"  Game: " + bold.Render(m.gameName),
		"",
		dim.Render("  Status view - simulation details will appear here."),
	}
	return strings.Join(lines, "\n")
}

func (m *partiesStatusScreen) Name() string          { return "Party Status" }
func (m *partiesStatusScreen) Layout() (string, int) { return "centered", 70 }
func (m *partiesStatusScreen) Shortcuts() []engine.ShortcutHint {
	return []engine.ShortcutHint{{Key: "L", Label: "Launch"}}
}
