package engine

import (
	"strconv"
	"strings"
)

// MouseEvent represents a parsed mouse event.
type MouseEvent struct {
	X         int
	Y         int
	Button    int
	IsScroll  bool
	ScrollDir int // -1 up, 1 down
}

// ClickHandler handles click events at a specific screen region.
type ClickHandler struct {
	Handler   func(relX, relY int)
	Component interface{}
	ColStart  int
	ColEnd    int
}

// Pointer manages the click grid and mouse event parsing.
type Pointer struct {
	grid map[int]map[int]*ClickHandler // row -> col -> handler
}

// NewPointer creates a new pointer manager.
func NewPointer() *Pointer {
	return &Pointer{
		grid: make(map[int]map[int]*ClickHandler),
	}
}

// ClearGrid clears all registered click zones.
func (p *Pointer) ClearGrid() {
	p.grid = make(map[int]map[int]*ClickHandler)
}

// PaintClick registers a click handler for a rectangular region.
func (p *Pointer) PaintClick(row, colStart, colEnd int, handler func(relX, relY int)) {
	if _, ok := p.grid[row]; !ok {
		p.grid[row] = make(map[int]*ClickHandler)
	}
	for col := colStart; col <= colEnd; col++ {
		p.grid[row][col] = &ClickHandler{
			Handler:  handler,
			ColStart: colStart,
			ColEnd:   colEnd,
		}
	}
}

// RegisterClickZone registers a click handler directly at coordinates.
func (p *Pointer) RegisterClickZone(row, col int, handler func(relX, relY int)) {
	if _, ok := p.grid[row]; !ok {
		p.grid[row] = make(map[int]*ClickHandler)
	}
	p.grid[row][col] = &ClickHandler{Handler: handler}
}

// HandleClick dispatches a click event to the registered handler.
func (p *Pointer) HandleClick(x, y int) {
	row, ok := p.grid[y]
	if !ok {
		return
	}
	handler, ok := row[x]
	if !ok {
		return
	}
	if handler.Handler != nil {
		handler.Handler(x-handler.ColStart, 0)
	}
}

// ParseMouse parses raw terminal mouse input into a MouseEvent.
// Supports SGR mode (\x1b[<...M/m), urxvt mode, and X10 mode.
func (p *Pointer) ParseMouse(data []byte) *MouseEvent {
	s := string(data)

	// SGR mode: \x1b[<button;x;yM or \x1b[<button;x;ym
	if strings.HasPrefix(s, "\x1b[<") {
		return p.parseSGR(s[3:])
	}

	// X10 mode: \x1b[M + 3 bytes
	if len(data) >= 6 && data[0] == 27 && data[1] == '[' && data[2] == 'M' {
		btn := int(data[3]) - 32
		x := int(data[4]) - 32 - 1
		y := int(data[5]) - 32 - 1

		if btn == 64 || btn == 65 {
			return &MouseEvent{
				X: x, Y: y,
				IsScroll:  true,
				ScrollDir: map[int]int{64: -1, 65: 1}[btn],
			}
		}

		return &MouseEvent{X: x, Y: y, Button: btn}
	}

	return nil
}

func (p *Pointer) parseSGR(s string) *MouseEvent {
	// Format: button;x;yM or button;x;ym
	// Remove trailing M or m
	isRelease := false
	if strings.HasSuffix(s, "m") {
		isRelease = true
		s = s[:len(s)-1]
	} else if strings.HasSuffix(s, "M") {
		s = s[:len(s)-1]
	} else {
		return nil
	}

	parts := strings.Split(s, ";")
	if len(parts) != 3 {
		return nil
	}

	btn, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil
	}
	x, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil
	}
	y, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil
	}

	// Convert to 0-based
	x--
	y--

	// Scroll wheel: buttons 64 (up) and 65 (down)
	if btn == 64 || btn == 65 {
		return &MouseEvent{
			X: x, Y: y,
			IsScroll:  true,
			ScrollDir: map[int]int{64: -1, 65: 1}[btn],
		}
	}

	// Only process press events, not release
	if isRelease {
		return nil
	}

	return &MouseEvent{X: x, Y: y, Button: btn & 0x03}
}
