package displays

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/brunojuliao/go-clappie/internal/engine"
	"github.com/brunojuliao/go-clappie/internal/platform"
)

type projectItem struct {
	Name, Path string
}

type projectsScreen struct {
	projects    []projectItem
	selectedIdx int
	styles      *engine.Styles
}

func NewProjectsScreen(data map[string]interface{}, styles *engine.Styles, claudePane string) engine.ScreenModel {
	m := &projectsScreen{styles: styles}
	m.load()
	return m
}

func (m *projectsScreen) load() {
	root, err := platform.ProjectRoot()
	if err != nil {
		return
	}
	projectsDir := platform.ProjectsDir(root)
	entries, err := os.ReadDir(projectsDir)
	if err != nil {
		return
	}
	m.projects = nil
	for _, e := range entries {
		if e.IsDir() {
			m.projects = append(m.projects, projectItem{
				Name: e.Name(),
				Path: filepath.Join(projectsDir, e.Name()),
			})
		}
	}
}

func (m *projectsScreen) Init() tea.Cmd { return nil }

func (m *projectsScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "up", "k":
			if m.selectedIdx > 0 {
				m.selectedIdx--
			}
		case "down", "j":
			if m.selectedIdx < len(m.projects)-1 {
				m.selectedIdx++
			}
		}
	}
	return m, nil
}

func (m *projectsScreen) View() string {
	dim := lipgloss.NewStyle().Faint(true)
	bold := lipgloss.NewStyle().Bold(true)

	var lines []string
	lines = append(lines, "")

	if len(m.projects) == 0 {
		lines = append(lines, "  No projects found.")
		lines = append(lines, "")
		lines = append(lines, dim.Render("  Create directories in projects/ to get started."))
	} else {
		lines = append(lines, fmt.Sprintf("  %d projects", len(m.projects)))
		lines = append(lines, "")
		for i, p := range m.projects {
			prefix := "  "
			if i == m.selectedIdx {
				prefix = "▸ "
			}
			line := fmt.Sprintf("%s📁 %s", prefix, p.Name)
			if i == m.selectedIdx {
				line = bold.Render(line)
			}
			lines = append(lines, line)
		}
	}

	return strings.Join(lines, "\n")
}

func (m *projectsScreen) Name() string                        { return "Projects" }
func (m *projectsScreen) Layout() (string, int)               { return "centered", 60 }
func (m *projectsScreen) Shortcuts() []engine.ShortcutHint { return nil }
