package engine

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/mattn/go-runewidth"
)

// ANSI escape code constants.
const (
	ClearScreen = "\x1b[2J\x1b[H"
	CursorHide  = "\x1b[?25l"
	CursorShow  = "\x1b[?25h"
	CursorHome  = "\x1b[H"
	Reset       = "\x1b[0m"
	Bold        = "\x1b[1m"
	Dim         = "\x1b[2m"
	Italic      = "\x1b[3m"
	Underline   = "\x1b[4m"
	Inverse     = "\x1b[7m"

	MouseEnable  = "\x1b[?1000h\x1b[?1002h\x1b[?1006h"
	MouseDisable = "\x1b[?1006l\x1b[?1002l\x1b[?1000l"

	AltScreenEnable  = "\x1b[?1049h"
	AltScreenDisable = "\x1b[?1049l"
)

var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]|\x1b\].*?\x1b\\|\x1b\][^\x07]*\x07`)

// VisualWidth returns the visual width of a string, accounting for
// wide characters (emoji, CJK) and stripping ANSI codes.
func VisualWidth(s string) int {
	stripped := StripANSI(s)
	return runewidth.StringWidth(stripped)
}

// StripANSI removes all ANSI escape sequences from a string.
func StripANSI(s string) string {
	return ansiRegex.ReplaceAllString(s, "")
}

// TruncateToWidth truncates a string to fit within maxWidth visual columns.
func TruncateToWidth(s string, maxWidth int, ellipsis string) string {
	if VisualWidth(s) <= maxWidth {
		return s
	}

	ellipsisWidth := VisualWidth(ellipsis)
	targetWidth := maxWidth - ellipsisWidth
	if targetWidth < 0 {
		targetWidth = 0
	}

	// Walk through the string tracking visual width
	stripped := StripANSI(s)
	width := 0
	var result strings.Builder
	for _, r := range stripped {
		rw := runewidth.RuneWidth(r)
		if width+rw > targetWidth {
			break
		}
		result.WriteRune(r)
		width += rw
	}

	return result.String() + ellipsis
}

// PadRight pads a string to the specified visual width.
func PadRight(s string, width int) string {
	vw := VisualWidth(s)
	if vw >= width {
		return s
	}
	return s + strings.Repeat(" ", width-vw)
}

// PadCenter centers a string within the specified visual width.
func PadCenter(s string, width int) string {
	vw := VisualWidth(s)
	if vw >= width {
		return s
	}
	leftPad := (width - vw) / 2
	rightPad := width - vw - leftPad
	return strings.Repeat(" ", leftPad) + s + strings.Repeat(" ", rightPad)
}

// Color returns a string with RGB foreground color.
func Color(r, g, b int, text string) string {
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm%s\x1b[0m", r, g, b, text)
}

// BgColor returns a string with RGB background color.
func BgColor(r, g, b int, text string) string {
	return fmt.Sprintf("\x1b[48;2;%d;%d;%dm%s\x1b[0m", r, g, b, text)
}

// FgBg returns a string with both foreground and background RGB colors.
func FgBg(fr, fg, fb, br, bg, bb int, text string) string {
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm\x1b[48;2;%d;%d;%dm%s\x1b[0m", fr, fg, fb, br, bg, bb, text)
}

// StyleBold wraps text in bold.
func StyleBold(text string) string {
	return Bold + text + Reset
}

// StyleDim wraps text in dim.
func StyleDim(text string) string {
	return Dim + text + Reset
}

// StyleItalic wraps text in italic.
func StyleItalic(text string) string {
	return Italic + text + Reset
}

// StyleUnderline wraps text in underline.
func StyleUnderline(text string) string {
	return Underline + text + Reset
}

// StyleInverse wraps text in inverse.
func StyleInverse(text string) string {
	return Inverse + text + Reset
}

// CursorTo moves the cursor to the given position (1-based).
func CursorTo(row, col int) string {
	return fmt.Sprintf("\x1b[%d;%dH", row, col)
}

// CursorUp moves the cursor up n lines.
func CursorUp(n int) string {
	return fmt.Sprintf("\x1b[%dA", n)
}

// CursorDown moves the cursor down n lines.
func CursorDown(n int) string {
	return fmt.Sprintf("\x1b[%dB", n)
}

// IsWideChar returns true if the rune is a wide character (2 columns).
func IsWideChar(r rune) bool {
	return runewidth.RuneWidth(r) == 2
}

// IsEmoji returns true if the rune is an emoji.
func IsEmoji(r rune) bool {
	// Emoji ranges
	if r >= 0x1F300 && r <= 0x1FAFF {
		return true
	}
	if r >= 0x2600 && r <= 0x27BF {
		return true
	}
	if r >= 0xFE00 && r <= 0xFE0F {
		return true // Variation selectors
	}
	return unicode.Is(unicode.So, r) // Symbol, other
}
