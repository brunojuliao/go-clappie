package ipc

import "encoding/json"

// Command represents an IPC command sent from CLI to daemon.
type Command struct {
	Action string          `json:"action"`
	View   string          `json:"view,omitempty"`
	Data   json.RawMessage `json:"data,omitempty"`

	// Display commands
	Message  string `json:"message,omitempty"`
	Duration int    `json:"duration,omitempty"` // milliseconds
	NoFocus  bool   `json:"noFocus,omitempty"`

	// Layout
	Layout   string `json:"layout,omitempty"`
	MaxWidth int    `json:"maxWidth,omitempty"`
}

// Response represents the daemon's response to a command.
type Response struct {
	OK      bool            `json:"ok"`
	Error   string          `json:"error,omitempty"`
	Data    json.RawMessage `json:"data,omitempty"`
	Message string          `json:"message,omitempty"`
}

// Known action constants.
const (
	ActionPing      = "ping"
	ActionPushView  = "push_view"
	ActionPopView   = "pop_view"
	ActionToast     = "toast"
	ActionClose     = "close"
	ActionListViews = "list_views"
	ActionKill      = "kill"
	ActionGetTheme  = "get_theme"
	ActionSetTheme  = "set_theme"
)
