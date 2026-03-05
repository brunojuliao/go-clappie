package engine

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// Styles provides lipgloss style presets derived from the Theme.
type Styles struct {
	theme *Theme

	// Base styles
	Background lipgloss.Style
	Text       lipgloss.Style
	TextMuted  lipgloss.Style
	Bold       lipgloss.Style

	// Semantic colors
	Success lipgloss.Style
	Error   lipgloss.Style
	Warning lipgloss.Style
	Info    lipgloss.Style

	// Layout elements
	HeaderStyle     lipgloss.Style
	ToastStyle      lipgloss.Style
	BorderFocused   lipgloss.Style
	BorderNormal    lipgloss.Style
	ShortcutKey     lipgloss.Style
	ShortcutLabel   lipgloss.Style
	SelectedItem    lipgloss.Style
	BreadcrumbStyle lipgloss.Style
}

// rgbToColor converts an RGB to a lipgloss.Color hex string.
func rgbToColor(c RGB) lipgloss.Color {
	return lipgloss.Color(fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B))
}

// NewStyles creates a Styles from a Theme.
func NewStyles(theme *Theme) *Styles {
	s := &Styles{theme: theme}
	s.Refresh()
	return s
}

// Refresh rebuilds all styles from the current theme colors.
func (s *Styles) Refresh() {
	t := s.theme

	bg := rgbToColor(t.C("background"))
	fg := rgbToColor(t.C("text"))
	muted := rgbToColor(t.C("textMuted"))

	s.Background = lipgloss.NewStyle().Background(bg)
	s.Text = lipgloss.NewStyle().Foreground(fg)
	s.TextMuted = lipgloss.NewStyle().Foreground(muted)
	s.Bold = lipgloss.NewStyle().Foreground(fg).Bold(true)

	s.Success = lipgloss.NewStyle().Foreground(rgbToColor(t.C("success")))
	s.Error = lipgloss.NewStyle().Foreground(rgbToColor(t.C("error")))
	s.Warning = lipgloss.NewStyle().Foreground(rgbToColor(t.C("warning")))
	s.Info = lipgloss.NewStyle().Foreground(rgbToColor(t.C("info")))

	s.HeaderStyle = lipgloss.NewStyle().
		Foreground(fg).
		Background(bg)

	s.ToastStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#000000")).
		Background(lipgloss.Color("#ffcc00")).
		Bold(true).
		Padding(0, 1)

	s.BorderFocused = lipgloss.NewStyle().
		BorderForeground(rgbToColor(t.C("primary")))

	s.BorderNormal = lipgloss.NewStyle().
		BorderForeground(rgbToColor(t.C("border")))

	s.ShortcutKey = lipgloss.NewStyle().
		Foreground(muted)

	s.ShortcutLabel = lipgloss.NewStyle().
		Foreground(fg)

	s.SelectedItem = lipgloss.NewStyle().
		Foreground(fg).
		Bold(true)

	s.BreadcrumbStyle = lipgloss.NewStyle().
		Foreground(muted)
}

// Theme returns the underlying Theme.
func (s *Styles) Theme() *Theme {
	return s.theme
}
