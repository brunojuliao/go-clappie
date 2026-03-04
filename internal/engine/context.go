package engine

import (
	"github.com/brunojuliao/go-clappie/internal/tmux"
)

// Shortcut represents a keyboard shortcut.
type Shortcut struct {
	Key     string
	Label   string
	Handler func()
}

// Context provides the API for display views to interact with the daemon.
type Context struct {
	Name       string
	Data       map[string]interface{}
	Width      int
	Height     int
	StackIndex int
	StackDepth int

	daemon    *Daemon
	lines     []string
	title     string
	desc      string
	layout    string
	maxWidth  int
	scrollTop int
	shortcuts map[string]Shortcut
}

// Draw sets the rendered lines for this view.
func (c *Context) Draw(lines []string) {
	c.lines = lines
	c.daemon.render()
}

// SetTitle sets the view title.
func (c *Context) SetTitle(title string) {
	c.title = title
}

// GetTitle returns the view title.
func (c *Context) GetTitle() string {
	return c.title
}

// SetDescription sets the view description/subtitle.
func (c *Context) SetDescription(desc string) {
	c.desc = desc
}

// GetDescription returns the view description.
func (c *Context) GetDescription() string {
	return c.desc
}

// SetLayout sets the layout mode and optional max width.
func (c *Context) SetLayout(layout string, maxWidth int) {
	c.layout = layout
	if maxWidth > 0 {
		c.maxWidth = maxWidth
	}
}

// GetLayout returns the current layout mode.
func (c *Context) GetLayout() string {
	return c.layout
}

// Push pushes a new view onto the stack.
func (c *Context) Push(viewName string, data map[string]interface{}) error {
	return c.daemon.pushView(viewName, data)
}

// Pop pops the current view from the stack.
func (c *Context) Pop() {
	c.daemon.popView()
	c.daemon.render()
}

// Submit types a [clappie] message into Claude's pane and presses Enter.
func (c *Context) Submit(message string) {
	tmux.SubmitToClaudePane(c.daemon.config.ClaudePane, message)
}

// Send types a [clappie] message into Claude's pane without pressing Enter.
func (c *Context) Send(message string) {
	tmux.SendToClaudePane(c.daemon.config.ClaudePane, message)
}

// Toast shows a toast notification.
func (c *Context) Toast(message string) {
	c.daemon.toast.Show(message, 0)
	c.daemon.render()
}

// RegisterShortcut registers a keyboard shortcut for this view.
func (c *Context) RegisterShortcut(key, label string, handler func()) {
	c.shortcuts[key] = Shortcut{
		Key:     key,
		Label:   label,
		Handler: handler,
	}
}

// GetLines returns the currently drawn lines.
func (c *Context) GetLines() []string {
	return c.lines
}

// GetShortcuts returns the registered shortcuts.
func (c *Context) GetShortcuts() map[string]Shortcut {
	return c.shortcuts
}

// ScrollTop returns the current scroll position.
func (c *Context) ScrollTop() int {
	return c.scrollTop
}

// SetScrollTop sets the scroll position.
func (c *Context) SetScrollTop(top int) {
	c.scrollTop = top
}

// SetContentLines hints the number of content lines for vertical centering.
func (c *Context) SetContentLines(count int) {
	// Used by renderer for vertical centering
}
