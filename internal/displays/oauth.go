package displays

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/brunojuliao/go-clappie/internal/engine"
	"github.com/brunojuliao/go-clappie/internal/oauth"
	"github.com/brunojuliao/go-clappie/internal/platform"
)

type oauthScreen struct {
	providers   []oauth.ProviderInfo
	selectedIdx int
	styles      *engine.Styles
}

func NewOAuthScreen(data map[string]interface{}, styles *engine.Styles, claudePane string) engine.ScreenModel {
	m := &oauthScreen{styles: styles}
	m.load()
	return m
}

func (m *oauthScreen) load() {
	root, err := platform.ProjectRoot()
	if err != nil {
		return
	}
	m.providers, _ = oauth.ListProviders(root)
}

func (m *oauthScreen) Init() tea.Cmd { return nil }

func (m *oauthScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "up", "k":
			if m.selectedIdx > 0 {
				m.selectedIdx--
			}
		case "down", "j":
			if m.selectedIdx < len(m.providers)-1 {
				m.selectedIdx++
			}
		case "a", "A":
			if m.selectedIdx < len(m.providers) {
				root, _ := platform.ProjectRoot()
				oauth.Auth(root, m.providers[m.selectedIdx].Name)
				return m, engine.ToastCmd("Auth flow started...", 0)
			}
		case "r", "R":
			if m.selectedIdx < len(m.providers) {
				root, _ := platform.ProjectRoot()
				oauth.Refresh(root, m.providers[m.selectedIdx].Name)
				m.load()
			}
		}
	}
	return m, nil
}

func (m *oauthScreen) View() string {
	dim := lipgloss.NewStyle().Faint(true)
	bold := lipgloss.NewStyle().Bold(true)
	green := lipgloss.NewStyle().Foreground(lipgloss.Color("#28a745"))
	yellow := lipgloss.NewStyle().Foreground(lipgloss.Color("#ffc107"))

	var lines []string
	lines = append(lines, "")

	if len(m.providers) == 0 {
		lines = append(lines, "  No OAuth providers configured.")
		lines = append(lines, "")
		lines = append(lines, dim.Render("  Add oauth.json to a skill directory."))
	} else {
		lines = append(lines, fmt.Sprintf("  %d providers", len(m.providers)))
		lines = append(lines, "")
		for i, p := range m.providers {
			prefix := "  "
			if i == m.selectedIdx {
				prefix = "▸ "
			}
			status := dim.Render("not authenticated")
			if p.HasToken {
				status = green.Render("authenticated")
				if p.Expired {
					status = yellow.Render("expired")
				}
			}
			line := fmt.Sprintf("%s%s [%s]", prefix, p.Name, status)
			if i == m.selectedIdx {
				line = bold.Render(line)
			}
			lines = append(lines, line)
		}
	}

	return strings.Join(lines, "\n")
}

func (m *oauthScreen) Name() string          { return "OAuth" }
func (m *oauthScreen) Layout() (string, int) { return "centered", 60 }
func (m *oauthScreen) Shortcuts() []engine.ShortcutHint {
	return []engine.ShortcutHint{
		{Key: "A", Label: "Auth"},
		{Key: "R", Label: "Refresh"},
	}
}
