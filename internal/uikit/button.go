package uikit

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ButtonStyle represents different button visual styles.
type ButtonStyle int

const (
	ButtonStyleDefault ButtonStyle = iota
	ButtonStyleFilled
	ButtonStyleGhost
	ButtonStyleInline
	ButtonStyleFullWidth
)

// ButtonConfig configures a button component.
type ButtonConfig struct {
	Label    string
	Shortcut string
	OnPress  func() tea.Cmd
	Style    ButtonStyle
	Width    int
}

// Button is a clickable button component.
type Button struct {
	config  ButtonConfig
	focused bool
	width   int
}

// NewButton creates a new button.
func NewButton(cfg ButtonConfig) Button {
	w := cfg.Width
	if w == 0 {
		w = len(cfg.Label) + 6
	}
	return Button{config: cfg, width: w}
}

// NewButtonFilled creates a filled-style button.
func NewButtonFilled(cfg ButtonConfig) Button {
	cfg.Style = ButtonStyleFilled
	return NewButton(cfg)
}

// NewButtonGhost creates a ghost-style button.
func NewButtonGhost(cfg ButtonConfig) Button {
	cfg.Style = ButtonStyleGhost
	return NewButton(cfg)
}

func (b Button) Init() tea.Cmd { return nil }

func (b Button) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "enter", " ":
			if b.config.OnPress != nil {
				return b, b.config.OnPress()
			}
		}
	}
	return b, nil
}

func (b Button) View() string {
	label := b.config.Label
	if b.config.Shortcut != "" {
		label = fmt.Sprintf("[%s] %s", b.config.Shortcut, label)
	}

	switch b.config.Style {
	case ButtonStyleFilled:
		return b.viewFilled(label)
	case ButtonStyleGhost:
		return b.viewGhost(label)
	case ButtonStyleInline:
		return b.viewInline(label)
	case ButtonStyleFullWidth:
		return b.viewFullWidth(label)
	default:
		return b.viewDefault(label)
	}
}

func (b Button) viewDefault(label string) string {
	w := b.width
	inner := w - 2
	padded := lipgloss.PlaceHorizontal(inner, lipgloss.Center, label)

	border := lipgloss.NormalBorder()
	style := lipgloss.NewStyle().
		Border(border).
		Width(inner)

	if b.focused {
		style = style.Bold(true).BorderForeground(lipgloss.Color("15"))
	}
	return style.Render(padded)
}

func (b Button) viewFilled(label string) string {
	style := lipgloss.NewStyle().
		Reverse(true).
		Padding(0, 1)
	if b.focused {
		style = style.Bold(true)
	}
	return style.Render(label)
}

func (b Button) viewGhost(label string) string {
	if b.focused {
		return lipgloss.NewStyle().Bold(true).Padding(0, 1).Render(label)
	}
	return lipgloss.NewStyle().Faint(true).Padding(0, 1).Render(label)
}

func (b Button) viewInline(label string) string {
	if b.focused {
		return lipgloss.NewStyle().Underline(true).Render(label)
	}
	return label
}

func (b Button) viewFullWidth(label string) string {
	padded := "  " + label + strings.Repeat(" ", max(0, b.width-len(label)-2))
	style := lipgloss.NewStyle().Reverse(true)
	if b.focused {
		style = style.Bold(true)
	}
	return style.Render(padded)
}

func (b Button) IsFocusable() bool { return true }
func (b Button) Focused() bool     { return b.focused }

func (b Button) Focus() Component {
	b.focused = true
	return b
}

func (b Button) Blur() Component {
	b.focused = false
	return b
}
