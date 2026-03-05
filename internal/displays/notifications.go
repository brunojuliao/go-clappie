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

type notifItem struct {
	Name, Path, Body, SourceID, Context, Created string
}

type notificationsScreen struct {
	items       []notifItem
	selectedIdx int
	styles      *engine.Styles
}

func NewNotificationsScreen(data map[string]interface{}, styles *engine.Styles, claudePane string) engine.ScreenModel {
	m := &notificationsScreen{styles: styles}
	m.loadNotifs()
	return m
}

func (m *notificationsScreen) loadNotifs() {
	root, err := platform.ProjectRoot()
	if err != nil {
		return
	}
	dir := platform.NotificationsCleanDir(root)
	entries, err := filestore.List(dir)
	if err != nil {
		return
	}
	m.items = nil
	for _, entry := range entries {
		body, blocks, err := filestore.ReadAndParse(entry.Path)
		if err != nil {
			continue
		}
		m.items = append(m.items, notifItem{
			Name:     entry.Name,
			Path:     entry.Path,
			Body:     body,
			SourceID: filestore.GetMetaField(blocks, "meta", "source_id"),
			Context:  filestore.GetMetaField(blocks, "meta", "context"),
			Created:  filestore.GetMetaField(blocks, "meta", "created"),
		})
	}
}

func (m *notificationsScreen) Init() tea.Cmd { return nil }

func (m *notificationsScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.selectedIdx > 0 {
				m.selectedIdx--
			}
		case "down", "j":
			if m.selectedIdx < len(m.items)-1 {
				m.selectedIdx++
			}
		case "d", "D":
			if m.selectedIdx < len(m.items) {
				filestore.DeleteFile(m.items[m.selectedIdx].Path)
				m.loadNotifs()
			}
		case "c", "C":
			for _, item := range m.items {
				filestore.DeleteFile(item.Path)
			}
			m.items = nil
			m.selectedIdx = 0
		}
	case tea.MouseMsg:
		if msg.Button == tea.MouseButtonWheelUp && m.selectedIdx > 0 {
			m.selectedIdx--
		} else if msg.Button == tea.MouseButtonWheelDown && m.selectedIdx < len(m.items)-1 {
			m.selectedIdx++
		}
	}
	return m, nil
}

func (m *notificationsScreen) View() string {
	dim := lipgloss.NewStyle().Faint(true)
	bold := lipgloss.NewStyle().Bold(true)

	var lines []string
	lines = append(lines, "")

	if len(m.items) == 0 {
		lines = append(lines, "  No notifications.")
		lines = append(lines, "")
		lines = append(lines, dim.Render("  Inbox zero!"))
	} else {
		lines = append(lines, fmt.Sprintf("  %d notifications", len(m.items)))
		lines = append(lines, "")

		for i, item := range m.items {
			prefix := "  "
			if i == m.selectedIdx {
				prefix = "▸ "
			}
			line := fmt.Sprintf("%s%s", prefix, item.Body)
			if len(line) > 55 {
				line = line[:52] + "..."
			}
			if item.Context != "" {
				line += dim.Render(fmt.Sprintf(" [%s]", item.Context))
			}
			if i == m.selectedIdx {
				line = bold.Render(line)
			}
			lines = append(lines, line)
		}
	}

	return strings.Join(lines, "\n")
}

func (m *notificationsScreen) Name() string          { return "Notifications" }
func (m *notificationsScreen) Layout() (string, int) { return "centered", 60 }
func (m *notificationsScreen) Shortcuts() []engine.ShortcutHint {
	return []engine.ShortcutHint{
		{Key: "D", Label: "Dismiss"},
		{Key: "C", Label: "Clear All"},
	}
}
