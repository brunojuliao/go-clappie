package engine

import "testing"

func TestVisualWidth(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"hello", 5},
		{"", 0},
		{"abc def", 7},
		// ANSI codes should not count
		{"\x1b[31mred\x1b[0m", 3},
		{"\x1b[1m\x1b[38;2;255;0;0mbold red\x1b[0m", 8},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := VisualWidth(tt.input)
			if got != tt.want {
				t.Errorf("VisualWidth(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

func TestStripANSI(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"hello", "hello"},
		{"\x1b[31mred\x1b[0m", "red"},
		{"\x1b[1m\x1b[38;2;255;0;0mbold\x1b[0m text", "bold text"},
		{"no codes here", "no codes here"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := StripANSI(tt.input)
			if got != tt.want {
				t.Errorf("StripANSI(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestTruncateToWidth(t *testing.T) {
	tests := []struct {
		input    string
		maxWidth int
		want     string
	}{
		{"hello", 10, "hello"},           // no truncation needed
		{"hello world", 8, "hello..."},    // truncated
		{"hi", 2, "hi"},                   // exact fit
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := TruncateToWidth(tt.input, tt.maxWidth, "...")
			gotWidth := VisualWidth(StripANSI(got))
			if gotWidth > tt.maxWidth {
				t.Errorf("TruncateToWidth(%q, %d) visual width = %d, exceeds max", tt.input, tt.maxWidth, gotWidth)
			}
		})
	}
}

func TestPadRight(t *testing.T) {
	got := PadRight("hi", 10)
	if VisualWidth(got) != 10 {
		t.Errorf("PadRight(\"hi\", 10) width = %d, want 10", VisualWidth(got))
	}
}

func TestPadCenter(t *testing.T) {
	got := PadCenter("hi", 10)
	if VisualWidth(got) != 10 {
		t.Errorf("PadCenter(\"hi\", 10) width = %d, want 10", VisualWidth(got))
	}
}
