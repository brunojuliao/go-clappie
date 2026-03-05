package uikit

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
)

// LoaderConfig configures a loader component.
type LoaderConfig struct {
	Label string
}

// Loader wraps bubbles/spinner.Model.
type Loader struct {
	config  LoaderConfig
	spinner spinner.Model
}

// NewLoader creates a new loader.
func NewLoader(cfg LoaderConfig) Loader {
	s := spinner.New()
	s.Spinner = spinner.Dot
	return Loader{config: cfg, spinner: s}
}

func (l Loader) Init() tea.Cmd { return l.spinner.Tick }

func (l Loader) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	l.spinner, cmd = l.spinner.Update(msg)
	return l, cmd
}

func (l Loader) View() string {
	label := lipgloss.NewStyle().Faint(true).Render(l.config.Label)
	return fmt.Sprintf("  %s %s", l.spinner.View(), label)
}

func (l Loader) IsFocusable() bool { return false }
func (l Loader) Focused() bool     { return false }
func (l Loader) Focus() Component  { return l }
func (l Loader) Blur() Component   { return l }
