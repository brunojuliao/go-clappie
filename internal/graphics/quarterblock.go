package graphics

import (
	"fmt"
	"strings"
)

// Unicode quarter-block characters for pixel art rendering.
// Each character cell is divided into a 2x2 grid of quarter blocks.
//
// ┌──┐
// │▘▝│  Upper-left, upper-right
// │▖▗│  Lower-left, lower-right
// └──┘
//
// All 16 combinations of filled/empty quarter blocks:
const (
	QBEmpty       = ' '  // 0000
	QBLowerLeft   = '▖'  // 0001
	QBLowerRight  = '▗'  // 0010
	QBLowerHalf   = '▄'  // 0011
	QBUpperLeft   = '▘'  // 0100
	QBLeftHalf    = '▌'  // 0101
	QBDiagUp      = '▚'  // 0110
	QBNoLowerR    = '▙'  // 0111
	QBUpperRight  = '▝'  // 1000
	QBDiagDown    = '▞'  // 1001
	QBRightHalf   = '▐'  // 1010
	QBNoLowerL    = '▜'  // 1011
	QBUpperHalf   = '▀'  // 1100
	QBNoUpperR    = '▛'  // 1101
	QBNoUpperL    = '▟'  // 1110
	QBFull        = '█'  // 1111
)

// quarterBlocks maps the 4-bit pattern to the corresponding character.
var quarterBlocks = [16]rune{
	QBEmpty, QBLowerLeft, QBLowerRight, QBLowerHalf,
	QBUpperLeft, QBLeftHalf, QBDiagUp, QBNoLowerR,
	QBUpperRight, QBDiagDown, QBRightHalf, QBNoLowerL,
	QBUpperHalf, QBNoUpperR, QBNoUpperL, QBFull,
}

// Pixel represents a colored pixel in the quarter-block grid.
type Pixel struct {
	R, G, B int
	Filled  bool
}

// Canvas is a 2D grid of pixels that renders using quarter blocks.
// Each terminal cell represents a 2x2 block of pixels, so a canvas
// of width W and height H produces W/2 columns and H/2 rows.
type Canvas struct {
	Width  int // pixel width
	Height int // pixel height
	Pixels [][]Pixel
}

// NewCanvas creates a new pixel canvas.
func NewCanvas(width, height int) *Canvas {
	pixels := make([][]Pixel, height)
	for y := range pixels {
		pixels[y] = make([]Pixel, width)
	}
	return &Canvas{
		Width:  width,
		Height: height,
		Pixels: pixels,
	}
}

// Set sets a pixel at the given coordinates.
func (c *Canvas) Set(x, y, r, g, b int) {
	if x >= 0 && x < c.Width && y >= 0 && y < c.Height {
		c.Pixels[y][x] = Pixel{R: r, G: g, B: b, Filled: true}
	}
}

// Clear resets all pixels.
func (c *Canvas) Clear() {
	for y := range c.Pixels {
		for x := range c.Pixels[y] {
			c.Pixels[y][x] = Pixel{}
		}
	}
}

// Render renders the canvas to terminal lines using quarter-block characters.
// Each output cell represents a 2x2 pixel block. Uses foreground color for
// the dominant filled region.
func (c *Canvas) Render(bgR, bgG, bgB int) []string {
	rows := (c.Height + 1) / 2
	cols := (c.Width + 1) / 2

	lines := make([]string, rows)
	for row := 0; row < rows; row++ {
		var sb strings.Builder
		for col := 0; col < cols; col++ {
			py := row * 2
			px := col * 2

			// Get the 4 pixels in this cell
			ul := c.getPixel(px, py)
			ur := c.getPixel(px+1, py)
			ll := c.getPixel(px, py+1)
			lr := c.getPixel(px+1, py+1)

			// Determine pattern
			pattern := 0
			if ll.Filled {
				pattern |= 1
			}
			if lr.Filled {
				pattern |= 2
			}
			if ul.Filled {
				pattern |= 4
			}
			if ur.Filled {
				pattern |= 8
			}

			char := quarterBlocks[pattern]

			if pattern == 0 {
				// All empty - use background
				sb.WriteString(fmt.Sprintf("\x1b[48;2;%d;%d;%dm ", bgR, bgG, bgB))
			} else if pattern == 15 {
				// All filled - use foreground as background
				avg := c.averageColor(ul, ur, ll, lr)
				sb.WriteString(fmt.Sprintf("\x1b[48;2;%d;%d;%dm ", avg.R, avg.G, avg.B))
			} else {
				// Mixed - use foreground color on bg
				fg := c.averageFilledColor(ul, ur, ll, lr)
				sb.WriteString(fmt.Sprintf("\x1b[38;2;%d;%d;%dm\x1b[48;2;%d;%d;%dm%c",
					fg.R, fg.G, fg.B, bgR, bgG, bgB, char))
			}
		}
		sb.WriteString("\x1b[0m")
		lines[row] = sb.String()
	}

	return lines
}

func (c *Canvas) getPixel(x, y int) Pixel {
	if x >= 0 && x < c.Width && y >= 0 && y < c.Height {
		return c.Pixels[y][x]
	}
	return Pixel{}
}

func (c *Canvas) averageColor(pixels ...Pixel) Pixel {
	var r, g, b, count int
	for _, p := range pixels {
		if p.Filled {
			r += p.R
			g += p.G
			b += p.B
			count++
		}
	}
	if count == 0 {
		return Pixel{}
	}
	return Pixel{R: r / count, G: g / count, B: b / count, Filled: true}
}

func (c *Canvas) averageFilledColor(pixels ...Pixel) Pixel {
	return c.averageColor(pixels...)
}
