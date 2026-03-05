package uikit

import "strings"

// joinLines joins string slices with newlines.
func joinLines(lines []string) string {
	return strings.Join(lines, "\n")
}
