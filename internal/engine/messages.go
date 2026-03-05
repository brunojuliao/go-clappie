package engine

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/brunojuliao/go-clappie/internal/ipc"
)

// PushViewMsg requests pushing a new view onto the stack.
type PushViewMsg struct {
	Name string
	Data map[string]interface{}
}

// PopViewMsg requests popping the current view from the stack.
type PopViewMsg struct{}

// ToastMsg shows a toast notification.
type ToastMsg struct {
	Message  string
	Duration time.Duration
}

// ToastExpiredMsg indicates the toast timer has expired.
type ToastExpiredMsg struct{}

// IPCCommandMsg wraps an IPC command received from external process.
type IPCCommandMsg struct {
	Cmd      ipc.Command
	ReplyCh  chan ipc.Response
}

// SubmitToClaudeMsg types a message into Claude's pane and presses Enter.
type SubmitToClaudeMsg struct {
	Message string
}

// SendToClaudeMsg types a message into Claude's pane without pressing Enter.
type SendToClaudeMsg struct {
	Message string
}

// TickMsg is sent periodically for animations and refresh.
type TickMsg struct{}

// HeartbeatCheckMsg triggers a heartbeat check.
type HeartbeatCheckMsg struct{}

// RefreshMsg triggers a re-render.
type RefreshMsg struct{}

// --- Command constructors ---

// PushViewCmd returns a command that pushes a view.
func PushViewCmd(name string, data map[string]interface{}) tea.Cmd {
	return func() tea.Msg {
		return PushViewMsg{Name: name, Data: data}
	}
}

// PopViewCmd returns a command that pops the current view.
func PopViewCmd() tea.Cmd {
	return func() tea.Msg {
		return PopViewMsg{}
	}
}

// ToastCmd returns a command that shows a toast.
func ToastCmd(message string, duration time.Duration) tea.Cmd {
	if duration == 0 {
		duration = 3 * time.Second
	}
	return func() tea.Msg {
		return ToastMsg{Message: message, Duration: duration}
	}
}

// SubmitToClaudeCmd returns a command that submits a message to Claude's pane.
func SubmitToClaudeCmd(message string) tea.Cmd {
	return func() tea.Msg {
		return SubmitToClaudeMsg{Message: message}
	}
}

// SendToClaudeCmd returns a command that sends a message to Claude's pane (no Enter).
func SendToClaudeCmd(message string) tea.Cmd {
	return func() tea.Msg {
		return SendToClaudeMsg{Message: message}
	}
}

// TickCmd returns a command that sends a tick after a delay.
func TickCmd(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(time.Time) tea.Msg {
		return TickMsg{}
	})
}

// HeartbeatCmd returns a command that sends a heartbeat check after a delay.
func HeartbeatCmd(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(time.Time) tea.Msg {
		return HeartbeatCheckMsg{}
	})
}
