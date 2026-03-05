package uikit

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ProgressConfig configures a progress bar component.
type ProgressConfig struct {
	Label string
	Value float64
	Max   float64
	Width int
}

// Progress is a progress bar component.
type Progress struct {
	config ProgressConfig
	value  float64
	width  int
}

// NewProgress creates a new progress bar.
func NewProgress(cfg ProgressConfig) Progress {
	w := cfg.Width
	if w == 0 {
		w = 30
	}
	if cfg.Max == 0 {
		cfg.Max = 1.0
	}
	return Progress{config: cfg, value: cfg.Value, width: w}
}

func (p Progress) Init() tea.Cmd                         { return nil }
func (p Progress) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return p, nil }

func (p Progress) View() string {
	barWidth := p.width - 4

	ratio := p.value / p.config.Max
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}

	filledWidth := int(float64(barWidth) * ratio)
	emptyWidth := barWidth - filledWidth

	bar := strings.Repeat("█", filledWidth) + strings.Repeat("░", emptyWidth)
	pct := lipgloss.NewStyle().Faint(true).Render(fmt.Sprintf("%3.0f%%", ratio*100))

	var lines []string
	if p.config.Label != "" {
		lines = append(lines, "  "+p.config.Label)
	}
	lines = append(lines, fmt.Sprintf("  %s %s", bar, pct))
	return joinLines(lines)
}

func (p Progress) IsFocusable() bool { return false }
func (p Progress) Focused() bool     { return false }
func (p Progress) Focus() Component  { return p }
func (p Progress) Blur() Component   { return p }

// SetValue sets the progress value.
func (p *Progress) SetValue(v float64) { p.value = v }

// Value returns the current value.
func (p Progress) Value() float64 { return p.value }
