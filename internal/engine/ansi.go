package engine

import (
	"os"
	"regexp"
	"strings"

	"github.com/mattn/go-runewidth"
)

func init() {
	// MSYS2/MinTTY renders Unicode ambiguous-width characters (box-drawing,
	// geometric shapes, bullets, block elements) as 2 cells wide.
	// Configure go-runewidth to match the terminal's actual rendering.
	if os.Getenv("MSYSTEM") != "" || os.Getenv("TERM_PROGRAM") == "mintty" {
		runewidth.DefaultCondition.EastAsianWidth = true
	}
	// Manual override
	switch os.Getenv("GO_CLAPPIE_EAST_ASIAN_WIDTH") {
	case "1":
		runewidth.DefaultCondition.EastAsianWidth = true
	case "0":
		runewidth.DefaultCondition.EastAsianWidth = false
	}
}

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

// RepeatToWidth repeats a string to fill the given visual width.
// This handles characters that may be 1 or 2 cells wide depending on terminal.
func RepeatToWidth(ch string, width int) string {
	chWidth := VisualWidth(ch)
	if chWidth <= 0 || width <= 0 {
		return ""
	}
	count := width / chWidth
	return strings.Repeat(ch, count)
}
