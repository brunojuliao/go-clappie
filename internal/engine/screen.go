package engine

import tea "github.com/charmbracelet/bubbletea"

// ShortcutHint represents a keyboard shortcut displayed in the footer.
// Unlike Shortcut (which includes a Handler), this is display-only;
// key handling is done in the ScreenModel's Update method.
type ShortcutHint struct {
	Key   string
	Label string
}

// ScreenModel is the interface that display views implement.
// It extends tea.Model with metadata for the layout engine.
type ScreenModel interface {
	tea.Model

	// Name returns the view name for breadcrumbs.
	Name() string

	// Layout returns the layout mode ("centered" or "full") and max width.
	Layout() (mode string, maxWidth int)

	// Shortcuts returns the keyboard shortcuts to display in the footer.
	Shortcuts() []ShortcutHint
}
