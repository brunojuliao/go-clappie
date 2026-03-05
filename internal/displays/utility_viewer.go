package displays

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/viewport"

	"github.com/brunojuliao/go-clappie/internal/engine"
	"github.com/brunojuliao/go-clappie/internal/filestore"
)

type utilityViewerScreen struct {
	viewport viewport.Model
	content  string
	styles   *engine.Styles
	ready    bool
}

func NewUtilityViewerScreen(data map[string]interface{}, styles *engine.Styles, claudePane string) engine.ScreenModel {
	content := ""
	if data != nil {
		if c, ok := data["content"].(string); ok {
			content = c
		}
		if path, ok := data["path"].(string); ok {
			d, err := filestore.ReadFile(path)
			if err == nil {
				content = d
			}
		}
	}

	// Prefix each line with indent
	var lines []string
	for _, line := range strings.Split(content, "\n") {
		lines = append(lines, "  "+line)
	}
	content = strings.Join(lines, "\n")

	return &utilityViewerScreen{content: content, styles: styles}
}

func (m *utilityViewerScreen) Init() tea.Cmd { return nil }

func (m *utilityViewerScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport = viewport.New(msg.Width, msg.Height-4)
		m.viewport.SetContent(m.content)
		m.ready = true
		return m, nil
	}

	if m.ready {
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m *utilityViewerScreen) View() string {
	if !m.ready {
		return "Loading..."
	}
	return m.viewport.View()
}

func (m *utilityViewerScreen) Name() string                        { return "Viewer" }
func (m *utilityViewerScreen) Layout() (string, int)               { return "full", 0 }
func (m *utilityViewerScreen) Shortcuts() []engine.ShortcutHint { return nil }
