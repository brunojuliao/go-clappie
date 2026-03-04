package uikit

import (
	"fmt"
	"strings"

	"github.com/brunojuliao/go-clappie/internal/engine"
)

// ProgressConfig configures a progress bar component.
type ProgressConfig struct {
	Label string
	Value float64 // 0.0 to 1.0
	Max   float64
	Width int
}

// Progress is a progress bar component.
type Progress struct {
	ComponentBase
	config ProgressConfig
	value  float64
}

// NewProgress creates a new progress bar.
func NewProgress(cfg ProgressConfig) *Progress {
	w := cfg.Width
	if w == 0 {
		w = 30
	}
	if cfg.Max == 0 {
		cfg.Max = 1.0
	}
	return &Progress{
		ComponentBase: ComponentBase{Focusable: false, Width: w},
		config:        cfg,
		value:         cfg.Value,
	}
}

// Render renders the progress bar.
func (p *Progress) Render(focused bool) []string {
	w := p.GetWidth()
	barWidth := w - 4

	ratio := p.value / p.config.Max
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}

	filled := int(float64(barWidth) * ratio)
	empty := barWidth - filled

	bar := strings.Repeat("█", filled) + strings.Repeat("░", empty)
	pct := fmt.Sprintf("%3.0f%%", ratio*100)

	var lines []string
	if p.config.Label != "" {
		lines = append(lines, "  "+p.config.Label)
	}
	lines = append(lines, fmt.Sprintf("  %s %s", bar, engine.StyleDim(pct)))
	return lines
}

// SetValue sets the progress value.
func (p *Progress) SetValue(v float64) {
	p.value = v
}

// Value returns the current value.
func (p *Progress) Value() float64 {
	return p.value
}
