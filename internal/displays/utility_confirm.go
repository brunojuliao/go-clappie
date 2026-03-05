package displays

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/brunojuliao/go-clappie/internal/engine"
)

type utilityConfirmScreen struct {
	message string
	styles  *engine.Styles
}

func NewUtilityConfirmScreen(data map[string]interface{}, styles *engine.Styles, claudePane string) engine.ScreenModel {
	message := "Are you sure?"
	if data != nil {
		if m, ok := data["message"].(string); ok {
			message = m
		}
	}
	return &utilityConfirmScreen{message: message, styles: styles}
}

func (m *utilityConfirmScreen) Init() tea.Cmd { return nil }

func (m *utilityConfirmScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "y", "Y":
			return m, tea.Batch(
				engine.SubmitToClaudeCmd("[go-clappie] Confirm → yes"),
				engine.PopViewCmd(),
			)
		case "n", "N":
			return m, tea.Batch(
				engine.SubmitToClaudeCmd("[go-clappie] Confirm → no"),
				engine.PopViewCmd(),
			)
		}
	}
	return m, nil
}

func (m *utilityConfirmScreen) View() string {
	lines := []string{
		"",
		"  " + m.message,
		"",
		"  [Y] Yes    [N] No",
	}
	return strings.Join(lines, "\n")
}

func (m *utilityConfirmScreen) Name() string          { return "Confirm" }
func (m *utilityConfirmScreen) Layout() (string, int) { return "centered", 50 }
func (m *utilityConfirmScreen) Shortcuts() []engine.ShortcutHint {
	return []engine.ShortcutHint{
		{Key: "Y", Label: "Yes"},
		{Key: "N", Label: "No"},
	}
}
