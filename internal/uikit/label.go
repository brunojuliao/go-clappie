package uikit

import "github.com/brunojuliao/go-clappie/internal/engine"

// LabelConfig configures a label component.
type LabelConfig struct {
	Text string
	Dim  bool
}

// Label is a static text display component.
type Label struct {
	ComponentBase
	config LabelConfig
}

// NewLabel creates a new label.
func NewLabel(cfg LabelConfig) *Label {
	return &Label{
		ComponentBase: ComponentBase{
			Focusable: false,
			Width:     engine.VisualWidth(cfg.Text) + 4,
		},
		config: cfg,
	}
}

// Render renders the label.
func (l *Label) Render(focused bool) []string {
	text := "  " + l.config.Text
	if l.config.Dim {
		text = engine.StyleDim(text)
	}
	return []string{text}
}

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
	ComponentBase
	config SectionHeadingConfig
}

// NewSectionHeading creates a new section heading.
func NewSectionHeading(cfg SectionHeadingConfig) *SectionHeading {
	return &SectionHeading{
		ComponentBase: ComponentBase{
			Focusable: false,
			Width:     engine.VisualWidth(cfg.Text) + 4,
		},
		config: cfg,
	}
}

// Render renders the section heading.
func (s *SectionHeading) Render(focused bool) []string {
	return []string{
		"",
		"  " + engine.StyleBold(s.config.Text),
		"",
	}
}
