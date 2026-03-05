package displays

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/brunojuliao/go-clappie/internal/engine"
	"github.com/brunojuliao/go-clappie/internal/filestore"
	"github.com/brunojuliao/go-clappie/internal/platform"
)

type choreItem struct {
	Name, Path, Title, Body, Status, Icon, Summary string
}

type choresScreen struct {
	chores      []choreItem
	selectedIdx int
	styles      *engine.Styles
	claudePane  string
}

func NewChoresScreen(data map[string]interface{}, styles *engine.Styles, claudePane string) engine.ScreenModel {
	m := &choresScreen{styles: styles, claudePane: claudePane}
	m.loadChores()
	return m
}

func (m *choresScreen) loadChores() {
	root, err := platform.ProjectRoot()
	if err != nil {
		return
	}
	dir := platform.ChoresHumansDir(root)
	entries, err := filestore.List(dir)
	if err != nil {
		return
	}
	m.chores = nil
	for _, entry := range entries {
		body, blocks, err := filestore.ReadAndParse(entry.Path)
		if err != nil {
			continue
		}
		title := filestore.GetMetaField(blocks, "chore-meta", "title")
		if title == "" {
			title = entry.Name
		}
		m.chores = append(m.chores, choreItem{
			Name:    entry.Name,
			Path:    entry.Path,
			Title:   title,
			Body:    body,
			Status:  filestore.GetMetaField(blocks, "chore-meta", "status"),
			Icon:    filestore.GetMetaField(blocks, "chore-meta", "icon"),
			Summary: filestore.GetMetaField(blocks, "chore-meta", "summary"),
		})
	}
	if m.selectedIdx >= len(m.chores) {
		m.selectedIdx = max(0, len(m.chores)-1)
	}
}

func (m *choresScreen) Init() tea.Cmd { return nil }

func (m *choresScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.selectedIdx > 0 {
				m.selectedIdx--
			}
			return m, nil
		case "down", "j":
			if m.selectedIdx < len(m.chores)-1 {
				m.selectedIdx++
			}
			return m, nil
		case "a", "A":
			if m.selectedIdx < len(m.chores) {
				c := m.chores[m.selectedIdx]
				_, blocks, _ := filestore.ReadAndParse(c.Path)
				filestore.SetMetaField(&blocks, "chore-meta", "status", "approved")
				filestore.WriteWithMeta(c.Path, c.Body, blocks)
				m.loadChores()
				return m, engine.SubmitToClaudeCmd(fmt.Sprintf("[go-clappie] Chore approved → %s", c.Title))
			}
		case "x", "X":
			if m.selectedIdx < len(m.chores) {
				c := m.chores[m.selectedIdx]
				_, blocks, _ := filestore.ReadAndParse(c.Path)
				filestore.SetMetaField(&blocks, "chore-meta", "status", "rejected")
				filestore.WriteWithMeta(c.Path, c.Body, blocks)
				m.loadChores()
				return m, engine.SubmitToClaudeCmd(fmt.Sprintf("[go-clappie] Chore rejected → %s", c.Title))
			}
		}
	case tea.MouseMsg:
		if msg.Button == tea.MouseButtonWheelUp && m.selectedIdx > 0 {
			m.selectedIdx--
		} else if msg.Button == tea.MouseButtonWheelDown && m.selectedIdx < len(m.chores)-1 {
			m.selectedIdx++
		}
		return m, nil
	}
	return m, nil
}

func (m *choresScreen) View() string {
	dim := lipgloss.NewStyle().Faint(true)
	bold := lipgloss.NewStyle().Bold(true)

	var lines []string
	lines = append(lines, "")

	if len(m.chores) == 0 {
		lines = append(lines, "  No pending chores.")
		lines = append(lines, "")
		lines = append(lines, dim.Render("  All caught up!"))
	} else {
		lines = append(lines, fmt.Sprintf("  %d pending chores", len(m.chores)))
		lines = append(lines, "")

		for i, c := range m.chores {
			prefix := "  "
			if i == m.selectedIdx {
				prefix = "▸ "
			}
			icon := c.Icon
			if icon == "" {
				icon = "📋"
			}
			line := fmt.Sprintf("%s%s %s", prefix, icon, c.Title)
			if c.Summary != "" {
				line += dim.Render(" — " + c.Summary)
			}
			if i == m.selectedIdx {
				line = bold.Render(line)
			}
			lines = append(lines, line)
		}

		if m.selectedIdx < len(m.chores) {
			lines = append(lines, "")
			lines = append(lines, bold.Render("  ─── Details ───"))
			lines = append(lines, "")
			if m.chores[m.selectedIdx].Body != "" {
				lines = append(lines, "  "+m.chores[m.selectedIdx].Body)
			}
		}
	}

	return strings.Join(lines, "\n")
}

func (m *choresScreen) Name() string          { return "Chores" }
func (m *choresScreen) Layout() (string, int) { return "centered", 60 }
func (m *choresScreen) Shortcuts() []engine.ShortcutHint {
	return []engine.ShortcutHint{
		{Key: "A", Label: "Approve"},
		{Key: "X", Label: "Reject"},
	}
}
