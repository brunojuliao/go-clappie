package uikit

import (
	"github.com/brunojuliao/go-clappie/internal/engine"
)

// ShortcutEntry represents a keyboard shortcut in the view.
type ShortcutEntry struct {
	Key     string
	Label   string
	Handler func()
}

// View manages a list of components with focus and shortcuts.
type View struct {
	ctx        *engine.Context
	components []Component
	focusIndex int
	shortcuts  map[string]ShortcutEntry
	width      int
}

// NewView creates a new view.
func NewView(ctx *engine.Context) *View {
	return &View{
		ctx:        ctx,
		focusIndex: -1,
		shortcuts:  make(map[string]ShortcutEntry),
		width:      60,
	}
}

// Add adds a component to the view.
func (v *View) Add(c Component) {
	v.components = append(v.components, c)
	// Auto-focus first focusable component
	if v.focusIndex == -1 && c.IsFocusable() {
		v.focusIndex = len(v.components) - 1
	}
}

// SetWidth sets the container width.
func (v *View) SetWidth(w int) {
	v.width = w
}

// Render renders all components and returns their combined lines.
func (v *View) Render() []string {
	var allLines []string
	for i, c := range v.components {
		focused := i == v.focusIndex
		lines := c.Render(focused)
		allLines = append(allLines, lines...)
	}
	return allLines
}

// HandleKey routes key events through the view.
func (v *View) HandleKey(key string) bool {
	// Tab navigation
	if key == "TAB" {
		v.FocusNext()
		return true
	}
	if key == "SHIFT_TAB" {
		v.FocusPrev()
		return true
	}

	// Forward to focused component
	if v.focusIndex >= 0 && v.focusIndex < len(v.components) {
		if v.components[v.focusIndex].OnKey(key) {
			return true
		}
	}

	// Check shortcuts
	if sc, ok := v.shortcuts[key]; ok {
		if sc.Handler != nil {
			sc.Handler()
		}
		return true
	}
	// Case-insensitive shortcut match
	if len(key) == 1 {
		upper := key
		if key[0] >= 'a' && key[0] <= 'z' {
			upper = string(key[0] - 32)
		}
		if sc, ok := v.shortcuts[upper]; ok {
			if sc.Handler != nil {
				sc.Handler()
			}
			return true
		}
	}

	return false
}

// FocusNext moves focus to the next focusable component.
func (v *View) FocusNext() {
	if len(v.components) == 0 {
		return
	}
	start := v.focusIndex + 1
	for i := 0; i < len(v.components); i++ {
		idx := (start + i) % len(v.components)
		if v.components[idx].IsFocusable() {
			v.focusIndex = idx
			return
		}
	}
}

// FocusPrev moves focus to the previous focusable component.
func (v *View) FocusPrev() {
	if len(v.components) == 0 {
		return
	}
	start := v.focusIndex - 1
	if start < 0 {
		start = len(v.components) - 1
	}
	for i := 0; i < len(v.components); i++ {
		idx := (start - i + len(v.components)) % len(v.components)
		if v.components[idx].IsFocusable() {
			v.focusIndex = idx
			return
		}
	}
}

// RegisterShortcut registers a keyboard shortcut.
func (v *View) RegisterShortcut(key, label string, handler func()) {
	v.shortcuts[key] = ShortcutEntry{
		Key:     key,
		Label:   label,
		Handler: handler,
	}
	// Also register with context for footer display
	v.ctx.RegisterShortcut(key, label, handler)
}

// GetShortcuts returns all registered shortcuts.
func (v *View) GetShortcuts() map[string]ShortcutEntry {
	return v.shortcuts
}

// FocusedComponent returns the currently focused component, or nil.
func (v *View) FocusedComponent() Component {
	if v.focusIndex >= 0 && v.focusIndex < len(v.components) {
		return v.components[v.focusIndex]
	}
	return nil
}
