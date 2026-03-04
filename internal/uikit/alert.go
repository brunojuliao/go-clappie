package uikit

import (
	"fmt"
	"strings"

	"github.com/brunojuliao/go-clappie/internal/engine"
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
	ComponentBase
	config AlertConfig
}

// NewAlert creates a new alert.
func NewAlert(cfg AlertConfig) *Alert {
	w := cfg.Width
	if w == 0 {
		w = 50
	}
	return &Alert{
		ComponentBase: ComponentBase{Focusable: false, Width: w},
		config:        cfg,
	}
}

// Render renders the alert.
func (a *Alert) Render(focused bool) []string {
	w := a.GetWidth()

	icon := a.icon()
	styleFunc := a.styleFunc()

	topBorder := "╭" + strings.Repeat("─", w-2) + "╮"
	botBorder := "╰" + strings.Repeat("─", w-2) + "╯"

	content := fmt.Sprintf(" %s %s", icon, a.config.Message)
	padded := engine.PadRight(content, w-2)

	return []string{
		styleFunc(topBorder),
		styleFunc("│" + padded + "│"),
		styleFunc(botBorder),
	}
}

func (a *Alert) icon() string {
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

func (a *Alert) styleFunc() func(string) string {
	switch a.config.Type {
	case AlertSuccess:
		return func(s string) string { return engine.Color(40, 167, 69, s) }
	case AlertWarning:
		return func(s string) string { return engine.Color(255, 193, 7, s) }
	case AlertError:
		return func(s string) string { return engine.Color(220, 53, 69, s) }
	default:
		return func(s string) string { return engine.Color(23, 162, 184, s) }
	}
}

// SetMessage changes the alert message.
func (a *Alert) SetMessage(msg string) {
	a.config.Message = msg
}

// SetType changes the alert type.
func (a *Alert) SetType(t AlertType) {
	a.config.Type = t
}
