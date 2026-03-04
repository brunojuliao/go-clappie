package engine

import "testing"

func TestKeyboardParse(t *testing.T) {
	kb := NewKeyboard()

	tests := []struct {
		name string
		data []byte
		want string
	}{
		{"Enter", []byte{13}, "ENTER"},
		{"Tab", []byte{9}, "TAB"},
		{"Escape", []byte{27}, "ESC"},
		{"Backspace", []byte{127}, "BACKSPACE"},
		{"Ctrl+C", []byte{3}, "CTRL_C"},
		{"Ctrl+D", []byte{4}, "CTRL_D"},
		{"Space", []byte{32}, "SPACE"},
		{"Letter a", []byte{'a'}, "a"},
		{"Letter Z", []byte{'Z'}, "Z"},
		{"Digit 5", []byte{'5'}, "5"},
		{"Arrow Up", []byte{27, '[', 'A'}, "UP"},
		{"Arrow Down", []byte{27, '[', 'B'}, "DOWN"},
		{"Arrow Right", []byte{27, '[', 'C'}, "RIGHT"},
		{"Arrow Left", []byte{27, '[', 'D'}, "LEFT"},
		{"Home", []byte{27, '[', 'H'}, "HOME"},
		{"End", []byte{27, '[', 'F'}, "END"},
		{"Shift+Tab", []byte{27, '[', 'Z'}, "SHIFT_TAB"},
		{"Delete", []byte{27, '[', '3', '~'}, "DELETE"},
		{"PageUp", []byte{27, '[', '5', '~'}, "PAGEUP"},
		{"PageDown", []byte{27, '[', '6', '~'}, "PAGEDOWN"},
		{"F1", []byte{27, 'O', 'P'}, "F1"},
		{"F2", []byte{27, 'O', 'Q'}, "F2"},
		{"Alt+a", []byte{27, 'a'}, "ALT_a"},
		{"Empty", []byte{}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := kb.Parse(tt.data)
			if got != tt.want {
				t.Errorf("Parse(%v) = %q, want %q", tt.data, got, tt.want)
			}
		})
	}
}
