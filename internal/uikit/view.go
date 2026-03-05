package uikit

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/brunojuliao/go-clappie/internal/engine"
)

// ViewContainer manages a list of components with focus cycling.
type ViewContainer struct {
	components []Component
	focusIndex int
}

// NewViewContainer creates a new view container.
func NewViewContainer() ViewContainer {
	return ViewContainer{
		focusIndex: -1,
	}
}

// Add adds a component to the container.
func (vc *ViewContainer) Add(c Component) {
	vc.components = append(vc.components, c)
	// Auto-focus first focusable component
	if vc.focusIndex == -1 && c.IsFocusable() {
		vc.focusIndex = len(vc.components) - 1
		vc.components[vc.focusIndex] = c.Focus()
	}
}

// Update handles key messages and delegates to focused component.
func (vc ViewContainer) Update(msg tea.Msg) (ViewContainer, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "tab":
			vc.focusNext()
			return vc, nil
		case "shift+tab":
			vc.focusPrev()
			return vc, nil
		}
	}

	// Forward to focused component
	if vc.focusIndex >= 0 && vc.focusIndex < len(vc.components) {
		updated, cmd := vc.components[vc.focusIndex].Update(msg)
		vc.components[vc.focusIndex] = updated.(Component)
		return vc, cmd
	}

	return vc, nil
}

// View renders all components joined vertically.
func (vc ViewContainer) View() string {
	var views []string
	for _, c := range vc.components {
		views = append(views, c.View())
	}
	return lipgloss.JoinVertical(lipgloss.Left, views...)
}

func (vc *ViewContainer) focusNext() {
	if len(vc.components) == 0 {
		return
	}
	// Blur current
	if vc.focusIndex >= 0 && vc.focusIndex < len(vc.components) {
		vc.components[vc.focusIndex] = vc.components[vc.focusIndex].Blur()
	}
	start := vc.focusIndex + 1
	for i := 0; i < len(vc.components); i++ {
		idx := (start + i) % len(vc.components)
		if vc.components[idx].IsFocusable() {
			vc.focusIndex = idx
			vc.components[idx] = vc.components[idx].Focus()
			return
		}
	}
}

func (vc *ViewContainer) focusPrev() {
	if len(vc.components) == 0 {
		return
	}
	// Blur current
	if vc.focusIndex >= 0 && vc.focusIndex < len(vc.components) {
		vc.components[vc.focusIndex] = vc.components[vc.focusIndex].Blur()
	}
	start := vc.focusIndex - 1
	if start < 0 {
		start = len(vc.components) - 1
	}
	for i := 0; i < len(vc.components); i++ {
		idx := (start - i + len(vc.components)) % len(vc.components)
		if vc.components[idx].IsFocusable() {
			vc.focusIndex = idx
			vc.components[idx] = vc.components[idx].Focus()
			return
		}
	}
}

// Shortcuts returns shortcut hints from all shortcut-bearing components.
func (vc ViewContainer) Shortcuts() []engine.ShortcutHint {
	return nil
}
