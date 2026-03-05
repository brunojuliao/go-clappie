package displays

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textarea"

	"github.com/brunojuliao/go-clappie/internal/engine"
)

type utilityEditorScreen struct {
	textarea textarea.Model
	styles   *engine.Styles
}

func NewUtilityEditorScreen(data map[string]interface{}, styles *engine.Styles, claudePane string) engine.ScreenModel {
	ta := textarea.New()
	ta.Placeholder = "Type here..."
	ta.SetWidth(60)
	ta.SetHeight(20)
	ta.Focus()

	if data != nil {
		if v, ok := data["value"].(string); ok {
			ta.SetValue(v)
		}
	}

	return &utilityEditorScreen{textarea: ta, styles: styles}
}

func (m *utilityEditorScreen) Init() tea.Cmd {
	return textarea.Blink
}

func (m *utilityEditorScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		// Ctrl+S to save
		if keyMsg.String() == "ctrl+s" {
			return m, tea.Batch(
				engine.SubmitToClaudeCmd("[go-clappie] Editor → "+m.textarea.Value()),
				engine.PopViewCmd(),
			)
		}
	}

	var cmd tea.Cmd
	m.textarea, cmd = m.textarea.Update(msg)
	return m, cmd
}

func (m *utilityEditorScreen) View() string {
	return m.textarea.View() + "\n\n  Ctrl+S to save"
}

func (m *utilityEditorScreen) Name() string          { return "Editor" }
func (m *utilityEditorScreen) Layout() (string, int) { return "full", 0 }
func (m *utilityEditorScreen) Shortcuts() []engine.ShortcutHint {
	return []engine.ShortcutHint{{Key: "Ctrl+S", Label: "Save"}}
}
