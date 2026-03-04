package uikit

// Component is the interface that all UI components implement.
type Component interface {
	// Render returns the rendered lines for this component.
	Render(focused bool) []string
	// IsFocusable returns true if this component can receive focus.
	IsFocusable() bool
	// GetWidth returns the visual width of this component.
	GetWidth() int
	// OnKey handles a key press. Returns true if handled.
	OnKey(key string) bool
	// OnClick handles a click at the given relative position.
	OnClick(lineIdx, col int) bool
}

// ComponentBase provides default implementations for Component.
type ComponentBase struct {
	Focusable bool
	Width     int
}

func (c *ComponentBase) IsFocusable() bool    { return c.Focusable }
func (c *ComponentBase) GetWidth() int         { return c.Width }
func (c *ComponentBase) OnKey(key string) bool { return false }
func (c *ComponentBase) OnClick(lineIdx, col int) bool { return false }
