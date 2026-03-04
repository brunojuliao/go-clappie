package engine

import (
	"strings"
	"sync"
	"time"
)

// Toast manages temporary notification overlays.
type Toast struct {
	message string
	timer   *time.Timer
	mu      sync.Mutex
}

// NewToast creates a new toast manager.
func NewToast() *Toast {
	return &Toast{}
}

// Show displays a toast message for the given duration.
// If duration is 0, default of 3 seconds is used.
func (t *Toast) Show(message string, duration time.Duration) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if duration == 0 {
		duration = 3 * time.Second
	}

	t.message = message

	if t.timer != nil {
		t.timer.Stop()
	}

	t.timer = time.AfterFunc(duration, func() {
		t.mu.Lock()
		t.message = ""
		t.mu.Unlock()
	})
}

// GetMessage returns the current toast message, or empty string if none.
func (t *Toast) GetMessage() string {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.message
}

// Clear removes the current toast.
func (t *Toast) Clear() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.message = ""
	if t.timer != nil {
		t.timer.Stop()
		t.timer = nil
	}
}

// Apply overlays the toast on the rendered output at row 2 (below breadcrumbs).
func (t *Toast) Apply(output string, width int) string {
	msg := t.GetMessage()
	if msg == "" {
		return output
	}

	lines := strings.Split(output, "\n")
	if len(lines) < 3 {
		return output
	}

	// Build toast line with colored background
	toastText := " " + msg + " "
	vw := VisualWidth(toastText)
	pad := width - vw
	if pad < 0 {
		pad = 0
	}

	// Use highlight color for toast background (amber-ish)
	toastLine := BgColor(200, 150, 50, Color(0, 0, 0, toastText+strings.Repeat(" ", pad)))

	// Replace row 2 (index 2, the header gap line)
	if len(lines) > 2 {
		lines[2] = toastLine
	}

	return strings.Join(lines, "\n")
}
