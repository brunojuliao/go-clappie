package uikit

import tea "github.com/charmbracelet/bubbletea"

// Component is the interface that all UI components implement.
type Component interface {
	tea.Model

	// IsFocusable returns true if this component can receive focus.
	IsFocusable() bool

	// Focused() returns true if the component has focus.
	Focused() bool

	// Focus returns a copy with focus.
	Focus() Component

	// Blur returns a copy without focus.
	Blur() Component
}
