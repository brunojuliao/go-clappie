package engine

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/brunojuliao/go-clappie/internal/ipc"
	"github.com/brunojuliao/go-clappie/internal/tmux"
)

// AppConfig holds configuration for the bubbletea app.
type AppConfig struct {
	SocketPath  string
	InitialView string
	InitialData string
	ClaudePane  string
	Registry    map[string]ViewModuleBT
}

// AppModel is the root bubbletea model that manages the view stack and layout.
type AppModel struct {
	config    AppConfig
	viewStack []ScreenModel
	theme     *Theme
	styles    *Styles
	toast     string
	width     int
	height    int
	program   *tea.Program
	server    *ipc.Server
	ipcCmds   chan IPCCommandMsg
}

// NewApp creates a new AppModel.
func NewApp(config AppConfig) *AppModel {
	theme := NewTheme()
	return &AppModel{
		config:  config,
		theme:   theme,
		styles:  NewStyles(theme),
		ipcCmds: make(chan IPCCommandMsg, 16),
	}
}

// SetProgram sets the tea.Program reference (needed for program.Send from IPC).
func (a *AppModel) SetProgram(p *tea.Program) {
	a.program = p
}

// StartIPCServer creates and starts the IPC server in the background.
func (a *AppModel) StartIPCServer() error {
	server, err := ipc.NewServer(a.config.SocketPath, a.handleIPCCommand)
	if err != nil {
		return fmt.Errorf("create IPC server: %w", err)
	}
	a.server = server
	go func() {
		a.server.Serve()
	}()
	return nil
}

// Shutdown cleans up resources.
func (a *AppModel) Shutdown() {
	if a.server != nil {
		a.server.Close()
	}
}

// handleIPCCommand is called from the IPC server goroutine.
// For commands that need to modify app state, we send via program.Send.
// For read-only commands like ping/list-views, we respond directly.
func (a *AppModel) handleIPCCommand(cmd ipc.Command) ipc.Response {
	switch cmd.Action {
	case ipc.ActionPing:
		return ipc.Response{OK: true, Message: "pong"}

	case ipc.ActionListViews:
		names := make([]string, len(a.viewStack))
		for i, v := range a.viewStack {
			names[i] = v.Name()
		}
		data, _ := json.Marshal(names)
		return ipc.Response{OK: true, Data: data}

	default:
		// Commands that mutate state go through the tea.Program event loop
		replyCh := make(chan ipc.Response, 1)
		a.program.Send(IPCCommandMsg{Cmd: cmd, ReplyCh: replyCh})
		select {
		case resp := <-replyCh:
			return resp
		case <-time.After(5 * time.Second):
			return ipc.Response{OK: false, Error: "timeout"}
		}
	}
}

// Init implements tea.Model.
func (a *AppModel) Init() tea.Cmd {
	var cmds []tea.Cmd

	// Push initial view
	if a.config.InitialView != "" {
		var data map[string]interface{}
		if a.config.InitialData != "" {
			json.Unmarshal([]byte(a.config.InitialData), &data)
		}
		cmds = append(cmds, PushViewCmd(a.config.InitialView, data))
	}

	// Start tick for animations
	cmds = append(cmds, TickCmd(500*time.Millisecond))

	// Start heartbeat
	cmds = append(cmds, HeartbeatCmd(5*time.Second))

	return tea.Batch(cmds...)
}

// Update implements tea.Model.
func (a *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		log.Printf("WindowSizeMsg: %dx%d", msg.Width, msg.Height)
		a.width = msg.Width
		a.height = msg.Height
		// Forward to current view (needed by viewport-based views like utility/viewer)
		if len(a.viewStack) > 0 {
			return a.updateCurrentView(msg)
		}
		return a, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return a, tea.Quit
		case "esc":
			if len(a.viewStack) > 1 {
				return a, PopViewCmd()
			}
			return a, tea.Quit
		default:
			// Forward to current view
			if len(a.viewStack) > 0 {
				return a.updateCurrentView(msg)
			}
		}
		return a, nil

	case tea.MouseMsg:
		if len(a.viewStack) > 0 {
			return a.updateCurrentView(msg)
		}
		return a, nil

	case PushViewMsg:
		return a.pushView(msg.Name, msg.Data)

	case PopViewMsg:
		return a.popView()

	case ToastMsg:
		a.toast = msg.Message
		duration := msg.Duration
		if duration == 0 {
			duration = 3 * time.Second
		}
		return a, tea.Tick(duration, func(time.Time) tea.Msg {
			return ToastExpiredMsg{}
		})

	case ToastExpiredMsg:
		a.toast = ""
		return a, nil

	case SubmitToClaudeMsg:
		tmux.SubmitToClaudePane(a.config.ClaudePane, msg.Message)
		return a, nil

	case SendToClaudeMsg:
		tmux.SendToClaudePane(a.config.ClaudePane, msg.Message)
		return a, nil

	case TickMsg:
		// Forward tick to current view for animations
		var cmd tea.Cmd
		if len(a.viewStack) > 0 {
			_, cmd = a.updateCurrentView(msg)
		}
		return a, tea.Batch(cmd, TickCmd(500*time.Millisecond))

	case HeartbeatCheckMsg:
		if a.config.ClaudePane != "" && !tmux.PaneExists(a.config.ClaudePane) {
			return a, tea.Quit
		}
		return a, HeartbeatCmd(5 * time.Second)

	case IPCCommandMsg:
		return a.handleIPCMsg(msg)

	default:
		// Forward unknown messages to current view
		if len(a.viewStack) > 0 {
			return a.updateCurrentView(msg)
		}
		return a, nil
	}
}

// View implements tea.Model.
func (a *AppModel) View() string {
	// Use fallback dimensions if WindowSizeMsg hasn't arrived yet
	// (common on Windows/MSYS2/MinTTY where terminal size queries may fail).
	if a.width == 0 {
		a.width = 80
	}
	if a.height == 0 {
		a.height = 24
	}

	var sections []string

	// Header: breadcrumbs
	sections = append(sections, a.renderHeader())

	// Toast (or empty line)
	sections = append(sections, a.renderToast())

	// Content
	sections = append(sections, a.renderContent())

	// Shortcuts
	sections = append(sections, a.renderShortcuts())

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (a *AppModel) renderHeader() string {
	breadcrumbs := a.buildBreadcrumbs()
	line := " " + a.styles.BreadcrumbStyle.Render(breadcrumbs)
	return a.styles.HeaderStyle.Width(a.width).Render(line)
}

func (a *AppModel) buildBreadcrumbs() string {
	var parts []string
	for _, v := range a.viewStack {
		parts = append(parts, v.Name())
	}
	if len(parts) > 3 {
		parts = append([]string{"..."}, parts[len(parts)-2:]...)
	}
	return strings.Join(parts, " > ")
}

func (a *AppModel) renderToast() string {
	if a.toast != "" {
		toast := a.styles.ToastStyle.Render(" " + a.toast + " ")
		return lipgloss.PlaceHorizontal(a.width, lipgloss.Center, toast)
	}
	return strings.Repeat(" ", a.width)
}

func (a *AppModel) renderContent() string {
	if len(a.viewStack) == 0 {
		return ""
	}

	current := a.viewStack[len(a.viewStack)-1]
	content := current.View()

	// Calculate available height for content
	// header(1) + toast(1) + shortcuts(2) = 4 lines overhead
	contentHeight := a.height - 4
	if contentHeight < 1 {
		contentHeight = 1
	}

	mode, maxWidth := current.Layout()

	contentStyle := lipgloss.NewStyle().Height(contentHeight)

	if mode == "centered" && maxWidth > 0 && maxWidth < a.width {
		contentStyle = contentStyle.Width(maxWidth)
		rendered := contentStyle.Render(content)
		return lipgloss.PlaceHorizontal(a.width, lipgloss.Center, rendered)
	}

	contentStyle = contentStyle.Width(a.width)
	return contentStyle.Render(content)
}

func (a *AppModel) renderShortcuts() string {
	if len(a.viewStack) == 0 {
		return ""
	}

	current := a.viewStack[len(a.viewStack)-1]
	shortcuts := current.Shortcuts()

	if len(shortcuts) == 0 {
		return lipgloss.NewStyle().Width(a.width).Height(2).Render("")
	}

	var parts []string
	for _, sc := range shortcuts {
		key := a.styles.ShortcutKey.Render(sc.Key)
		label := a.styles.ShortcutLabel.Render(sc.Label)
		parts = append(parts, fmt.Sprintf(" %s %s ", key, label))
	}
	line := " " + strings.Join(parts, "  ")

	return lipgloss.NewStyle().Width(a.width).Height(2).Render(line)
}

// pushView creates and pushes a new view onto the stack.
func (a *AppModel) pushView(name string, data map[string]interface{}) (tea.Model, tea.Cmd) {
	module, ok := a.config.Registry[name]
	if !ok {
		return a, nil
	}

	screen := module.Create(data, a.styles, a.config.ClaudePane)
	a.viewStack = append(a.viewStack, screen)
	log.Printf("pushView: %s (stack size: %d)", name, len(a.viewStack))

	initCmd := screen.Init()

	// Send current dimensions to the new view so viewport-based views
	// (like utility/viewer) can initialize even if WindowSizeMsg already fired.
	var sizeCmd tea.Cmd
	if a.width > 0 && a.height > 0 {
		idx := len(a.viewStack) - 1
		updated, cmd := a.viewStack[idx].Update(tea.WindowSizeMsg{Width: a.width, Height: a.height})
		a.viewStack[idx] = updated.(ScreenModel)
		sizeCmd = cmd
	}

	return a, tea.Batch(initCmd, sizeCmd)
}

// popView removes the top view from the stack.
func (a *AppModel) popView() (tea.Model, tea.Cmd) {
	if len(a.viewStack) == 0 {
		return a, nil
	}
	a.viewStack = a.viewStack[:len(a.viewStack)-1]
	return a, nil
}

// updateCurrentView forwards a message to the current view and updates the stack.
func (a *AppModel) updateCurrentView(msg tea.Msg) (tea.Model, tea.Cmd) {
	idx := len(a.viewStack) - 1
	updated, cmd := a.viewStack[idx].Update(msg)
	a.viewStack[idx] = updated.(ScreenModel)
	return a, cmd
}

// handleIPCMsg processes IPC commands that mutate state.
func (a *AppModel) handleIPCMsg(msg IPCCommandMsg) (tea.Model, tea.Cmd) {
	cmd := msg.Cmd
	var resp ipc.Response
	var teaCmd tea.Cmd

	switch cmd.Action {
	case ipc.ActionPushView:
		var data map[string]interface{}
		if cmd.Data != nil {
			json.Unmarshal(cmd.Data, &data)
		}
		_, teaCmd = a.pushView(cmd.View, data)
		resp = ipc.Response{OK: true}

	case ipc.ActionPopView:
		a.popView()
		resp = ipc.Response{OK: true}

	case ipc.ActionToast:
		duration := cmd.Duration
		if duration == 0 {
			duration = 3000
		}
		a.toast = cmd.Message
		log.Printf("toast: %q (duration=%dms)", cmd.Message, duration)
		teaCmd = tea.Tick(time.Duration(duration)*time.Millisecond, func(time.Time) tea.Msg {
			return ToastExpiredMsg{}
		})
		resp = ipc.Response{OK: true}

	case ipc.ActionClose, ipc.ActionKill:
		resp = ipc.Response{OK: true}
		if msg.ReplyCh != nil {
			msg.ReplyCh <- resp
		}
		return a, tea.Quit

	case ipc.ActionGetTheme:
		resp = ipc.Response{OK: true, Message: a.theme.GetMode()}

	case ipc.ActionSetTheme:
		a.theme.SetMode(cmd.Message)
		a.styles.Refresh()
		resp = ipc.Response{OK: true}

	default:
		resp = ipc.Response{OK: false, Error: fmt.Sprintf("unknown action: %s", cmd.Action)}
	}

	if msg.ReplyCh != nil {
		msg.ReplyCh <- resp
	}
	return a, teaCmd
}
