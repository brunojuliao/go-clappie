package uikit

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// AlertType represents the alert severity.
type AlertType int

const (
	AlertInfo AlertType = iota
	AlertSuccess
	AlertWarning
	AlertError
)

// AlertConfig configures an alert component.
type AlertConfig struct {
	Type    AlertType
	Message string
	Width   int
}

// Alert is a notification/message box component.
type Alert struct {
	config AlertConfig
	width  int
}

// NewAlert creates a new alert.
func NewAlert(cfg AlertConfig) Alert {
	w := cfg.Width
	if w == 0 {
		w = 50
	}
	return Alert{config: cfg, width: w}
}

func (a Alert) Init() tea.Cmd                         { return nil }
func (a Alert) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return a, nil }

func (a Alert) View() string {
	icon := a.icon()
	color := a.color()

	content := fmt.Sprintf(" %s %s", icon, a.config.Message)

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(color).
		Foreground(color).
		Width(a.width - 4)

	return style.Render(content)
}

func (a Alert) icon() string {
	switch a.config.Type {
	case AlertSuccess:
		return "✓"
	case AlertWarning:
		return "⚠"
	case AlertError:
		return "✗"
	default:
		return "ℹ"
	}
}

func (a Alert) color() lipgloss.Color {
	switch a.config.Type {
	case AlertSuccess:
		return lipgloss.Color("#28a745")
	case AlertWarning:
		return lipgloss.Color("#ffc107")
	case AlertError:
		return lipgloss.Color("#dc3545")
	default:
		return lipgloss.Color("#17a2b8")
	}
}

func (a Alert) IsFocusable() bool { return false }
func (a Alert) Focused() bool     { return false }
func (a Alert) Focus() Component  { return a }
func (a Alert) Blur() Component   { return a }

// SetMessage changes the alert message.
func (a *Alert) SetMessage(msg string) { a.config.Message = msg }

// SetType changes the alert type.
func (a *Alert) SetType(t AlertType) { a.config.Type = t }
