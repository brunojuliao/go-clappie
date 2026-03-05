package uikit

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// LabelConfig configures a label component.
type LabelConfig struct {
	Text string
	Dim  bool
}

// Label is a static text display component.
type Label struct {
	config LabelConfig
}

// NewLabel creates a new label.
func NewLabel(cfg LabelConfig) Label {
	return Label{config: cfg}
}

func (l Label) Init() tea.Cmd                         { return nil }
func (l Label) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return l, nil }

func (l Label) View() string {
	text := "  " + l.config.Text
	if l.config.Dim {
		return lipgloss.NewStyle().Faint(true).Render(text)
	}
	return text
}

func (l Label) IsFocusable() bool { return false }
func (l Label) Focused() bool     { return false }
func (l Label) Focus() Component  { return l }
func (l Label) Blur() Component   { return l }

// SetText changes the label text.
func (l *Label) SetText(text string) {
	l.config.Text = text
}

// SectionHeadingConfig configures a section heading.
type SectionHeadingConfig struct {
	Text string
}

// SectionHeading is a bold section heading component.
type SectionHeading struct {
	config SectionHeadingConfig
}

// NewSectionHeading creates a new section heading.
func NewSectionHeading(cfg SectionHeadingConfig) SectionHeading {
	return SectionHeading{config: cfg}
}

func (s SectionHeading) Init() tea.Cmd                         { return nil }
func (s SectionHeading) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return s, nil }

func (s SectionHeading) View() string {
	bold := lipgloss.NewStyle().Bold(true).Render("  " + s.config.Text)
	return "\n" + bold + "\n"
}

func (s SectionHeading) IsFocusable() bool { return false }
func (s SectionHeading) Focused() bool     { return false }
func (s SectionHeading) Focus() Component  { return s }
func (s SectionHeading) Blur() Component   { return s }
