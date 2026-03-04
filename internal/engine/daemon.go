package engine

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/brunojuliao/go-clappie/internal/ipc"
	"github.com/brunojuliao/go-clappie/internal/tmux"
	"golang.org/x/term"
)

// DaemonConfig holds configuration for the display daemon.
type DaemonConfig struct {
	SocketPath  string
	InitialView string
	InitialData string
	ClaudePane  string
	Registry    map[string]ViewModule // injected view registry
}

// Daemon is the display engine daemon.
type Daemon struct {
	config    DaemonConfig
	server    *ipc.Server
	viewStack []*ViewInstance
	width     int
	height    int
	theme     *Theme
	toast     *Toast
	keyboard  *Keyboard
	pointer   *Pointer
	renderer  *Renderer
	mu        sync.Mutex
	quit      chan struct{}
	oldState  *term.State
}

// NewDaemon creates a new display daemon.
func NewDaemon(config DaemonConfig) (*Daemon, error) {
	d := &Daemon{
		config: config,
		theme:  NewTheme(),
		toast:  NewToast(),
		quit:   make(chan struct{}),
	}

	d.keyboard = NewKeyboard()
	d.pointer = NewPointer()
	d.renderer = NewRenderer(d)

	// Create IPC server
	server, err := ipc.NewServer(config.SocketPath, d.handleCommand)
	if err != nil {
		return nil, fmt.Errorf("create IPC server: %w", err)
	}
	d.server = server

	return d, nil
}

// Run starts the daemon main loop.
func (d *Daemon) Run() error {
	// Setup terminal
	if err := d.setupTerminal(); err != nil {
		return fmt.Errorf("setup terminal: %w", err)
	}
	defer d.restoreTerminal()

	// Get initial dimensions
	d.updateDimensions()

	// Start IPC server in background
	go func() {
		if err := d.server.Serve(); err != nil {
			log.Printf("IPC server error: %v", err)
		}
	}()

	// Push initial view
	if d.config.InitialView != "" {
		var data map[string]interface{}
		if d.config.InitialData != "" {
			json.Unmarshal([]byte(d.config.InitialData), &data)
		}
		d.pushView(d.config.InitialView, data)
	}

	// Start heartbeat (check if parent pane still exists)
	go d.heartbeatLoop()

	// Main event loop
	d.eventLoop()

	return nil
}

// Shutdown gracefully shuts down the daemon.
func (d *Daemon) Shutdown() {
	select {
	case <-d.quit:
		return // Already shutting down
	default:
		close(d.quit)
	}
	d.server.Close()
}

func (d *Daemon) setupTerminal() error {
	// Set stdin to raw mode
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	d.oldState = oldState

	// Hide cursor
	fmt.Fprint(os.Stdout, "\x1b[?25l")
	// Enable mouse tracking (SGR mode)
	fmt.Fprint(os.Stdout, "\x1b[?1000h\x1b[?1002h\x1b[?1006h")
	// Enable alternate screen buffer
	fmt.Fprint(os.Stdout, "\x1b[?1049h")
	// Clear screen
	fmt.Fprint(os.Stdout, "\x1b[2J\x1b[H")

	return nil
}

func (d *Daemon) restoreTerminal() {
	// Disable mouse tracking
	fmt.Fprint(os.Stdout, "\x1b[?1006l\x1b[?1002l\x1b[?1000l")
	// Show cursor
	fmt.Fprint(os.Stdout, "\x1b[?25h")
	// Disable alternate screen buffer
	fmt.Fprint(os.Stdout, "\x1b[?1049l")

	if d.oldState != nil {
		term.Restore(int(os.Stdin.Fd()), d.oldState)
	}
}

func (d *Daemon) updateDimensions() {
	w, h, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		// Fallback to tmux dimensions
		w, h, err = tmux.GetPaneSize("")
		if err != nil {
			w, h = 80, 24
		}
	}
	d.mu.Lock()
	d.width = w
	d.height = h
	d.mu.Unlock()
}

func (d *Daemon) eventLoop() {
	buf := make([]byte, 256)
	stdinCh := make(chan []byte, 16)

	// Read stdin in goroutine
	go func() {
		for {
			n, err := os.Stdin.Read(buf)
			if err != nil {
				return
			}
			data := make([]byte, n)
			copy(data, buf[:n])
			stdinCh <- data
		}
	}()

	// Ticker for periodic refresh
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-d.quit:
			return
		case data := <-stdinCh:
			d.handleInput(data)
		case <-ticker.C:
			d.render()
		}
	}
}

func (d *Daemon) handleInput(data []byte) {
	// Check for mouse input first
	if mouseEvent := d.pointer.ParseMouse(data); mouseEvent != nil {
		if mouseEvent.IsScroll {
			d.handleScroll(mouseEvent.ScrollDir)
		} else {
			d.handleClick(mouseEvent.X, mouseEvent.Y)
		}
		return
	}

	// Parse keyboard input
	key := d.keyboard.Parse(data)
	if key == "" {
		return
	}

	// Global keys
	switch key {
	case "CTRL_C":
		d.Shutdown()
		return
	case "ESC":
		if len(d.viewStack) > 1 {
			d.popView()
			d.render()
		} else {
			d.Shutdown()
		}
		return
	}

	// Forward to current view
	d.mu.Lock()
	view := d.currentView()
	d.mu.Unlock()

	if view != nil && view.Instance.OnKey != nil {
		if view.Instance.OnKey(key) {
			return // handled
		}
	}
}

func (d *Daemon) handleClick(x, y int) {
	d.pointer.HandleClick(x, y)
	d.render()
}

func (d *Daemon) handleScroll(direction int) {
	d.mu.Lock()
	view := d.currentView()
	d.mu.Unlock()

	if view != nil && view.Instance.OnScroll != nil {
		view.Instance.OnScroll(direction)
		d.render()
	}
}

func (d *Daemon) heartbeatLoop() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-d.quit:
			return
		case <-ticker.C:
			if d.config.ClaudePane != "" && !tmux.PaneExists(d.config.ClaudePane) {
				d.Shutdown()
				return
			}
		}
	}
}

func (d *Daemon) handleCommand(cmd ipc.Command) ipc.Response {
	switch cmd.Action {
	case ipc.ActionPing:
		return ipc.Response{OK: true, Message: "pong"}

	case ipc.ActionPushView:
		var data map[string]interface{}
		if cmd.Data != nil {
			json.Unmarshal(cmd.Data, &data)
		}
		if err := d.pushView(cmd.View, data); err != nil {
			return ipc.Response{OK: false, Error: err.Error()}
		}
		d.render()
		return ipc.Response{OK: true}

	case ipc.ActionPopView:
		d.popView()
		d.render()
		return ipc.Response{OK: true}

	case ipc.ActionToast:
		duration := cmd.Duration
		if duration == 0 {
			duration = 3000
		}
		d.toast.Show(cmd.Message, time.Duration(duration)*time.Millisecond)
		d.render()
		return ipc.Response{OK: true}

	case ipc.ActionClose:
		d.Shutdown()
		return ipc.Response{OK: true}

	case ipc.ActionListViews:
		names := make([]string, len(d.viewStack))
		for i, v := range d.viewStack {
			names[i] = v.Name
		}
		data, _ := json.Marshal(names)
		return ipc.Response{OK: true, Data: data}

	case ipc.ActionKill:
		d.Shutdown()
		return ipc.Response{OK: true}

	default:
		return ipc.Response{OK: false, Error: fmt.Sprintf("unknown action: %s", cmd.Action)}
	}
}

func (d *Daemon) pushView(name string, data map[string]interface{}) error {
	module, ok := d.config.Registry[name]
	if !ok {
		return fmt.Errorf("unknown view: %s", name)
	}

	ctx := d.createContext(name, data)
	instance := module.Create(ctx)

	vi := &ViewInstance{
		Name:     name,
		Module:   module,
		Context:  ctx,
		Instance: instance,
	}

	d.mu.Lock()
	d.viewStack = append(d.viewStack, vi)
	d.mu.Unlock()

	if instance.Init != nil {
		instance.Init()
	}

	return nil
}

func (d *Daemon) popView() {
	d.mu.Lock()
	defer d.mu.Unlock()

	if len(d.viewStack) == 0 {
		return
	}

	top := d.viewStack[len(d.viewStack)-1]
	if top.Instance.Cleanup != nil {
		top.Instance.Cleanup()
	}
	d.viewStack = d.viewStack[:len(d.viewStack)-1]
}

func (d *Daemon) currentView() *ViewInstance {
	if len(d.viewStack) == 0 {
		return nil
	}
	return d.viewStack[len(d.viewStack)-1]
}

func (d *Daemon) render() {
	d.mu.Lock()
	view := d.currentView()
	w, h := d.width, d.height
	d.mu.Unlock()

	if view == nil || w == 0 || h == 0 {
		return
	}

	// Render view content
	if view.Instance.Render != nil {
		view.Instance.Render()
	}

	// Compose full screen
	output := d.renderer.Compose(view, w, h)

	// Apply toast overlay
	output = d.toast.Apply(output, w)

	// Write to terminal
	fmt.Fprint(os.Stdout, "\x1b[H") // Move cursor to top-left
	fmt.Fprint(os.Stdout, output)
}

func (d *Daemon) createContext(name string, data map[string]interface{}) *Context {
	return &Context{
		Name:       name,
		Data:       data,
		daemon:     d,
		lines:      nil,
		title:      "",
		desc:       "",
		layout:     "centered",
		maxWidth:   60,
		scrollTop:  0,
		shortcuts:  make(map[string]Shortcut),
	}
}
