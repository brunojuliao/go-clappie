package uikit

import (
	"fmt"
	"time"

	"github.com/brunojuliao/go-clappie/internal/engine"
)

var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// LoaderConfig configures a loader component.
type LoaderConfig struct {
	Label string
}

// Loader is an animated loading spinner component.
type Loader struct {
	ComponentBase
	config LoaderConfig
	frame  int
}

// NewLoader creates a new loader.
func NewLoader(cfg LoaderConfig) *Loader {
	return &Loader{
		ComponentBase: ComponentBase{
			Focusable: false,
			Width:     engine.VisualWidth(cfg.Label) + 4,
		},
		config: cfg,
	}
}

// Render renders the loader with current animation frame.
func (l *Loader) Render(focused bool) []string {
	// Advance frame based on time
	l.frame = int(time.Now().UnixMilli()/100) % len(spinnerFrames)
	spinner := spinnerFrames[l.frame]
	return []string{fmt.Sprintf("  %s %s", spinner, engine.StyleDim(l.config.Label))}
}
