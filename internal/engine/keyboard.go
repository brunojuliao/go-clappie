package engine

// Keyboard parses raw terminal input into key names.
type Keyboard struct{}

// NewKeyboard creates a new keyboard parser.
func NewKeyboard() *Keyboard {
	return &Keyboard{}
}

// Parse converts raw bytes from terminal input into a key name string.
func (k *Keyboard) Parse(data []byte) string {
	if len(data) == 0 {
		return ""
	}

	// Single byte
	if len(data) == 1 {
		b := data[0]
		switch {
		case b == 3:
			return "CTRL_C"
		case b == 4:
			return "CTRL_D"
		case b == 9:
			return "TAB"
		case b == 10 || b == 13:
			return "ENTER"
		case b == 27:
			return "ESC"
		case b == 127 || b == 8:
			return "BACKSPACE"
		case b == 1:
			return "CTRL_A"
		case b == 2:
			return "CTRL_B"
		case b == 5:
			return "CTRL_E"
		case b == 6:
			return "CTRL_F"
		case b == 11:
			return "CTRL_K"
		case b == 12:
			return "CTRL_L"
		case b == 14:
			return "CTRL_N"
		case b == 16:
			return "CTRL_P"
		case b == 18:
			return "CTRL_R"
		case b == 21:
			return "CTRL_U"
		case b == 23:
			return "CTRL_W"
		case b == 26:
			return "CTRL_Z"
		case b == 32:
			return "SPACE"
		case b >= 33 && b <= 126:
			return string(b)
		}
		return ""
	}

	// Escape sequences
	if data[0] == 27 {
		if len(data) >= 3 && data[1] == '[' {
			return k.parseCSI(data[2:])
		}
		if len(data) >= 3 && data[1] == 'O' {
			return k.parseSS3(data[2])
		}
		// Alt+key
		if len(data) == 2 && data[1] >= 33 && data[1] <= 126 {
			return "ALT_" + string(data[1])
		}
	}

	// Multi-byte UTF-8 character
	s := string(data)
	if len([]rune(s)) == 1 {
		return s
	}

	return ""
}

func (k *Keyboard) parseCSI(data []byte) string {
	if len(data) == 0 {
		return ""
	}

	// Simple sequences: ESC [ X
	if len(data) == 1 {
		switch data[0] {
		case 'A':
			return "UP"
		case 'B':
			return "DOWN"
		case 'C':
			return "RIGHT"
		case 'D':
			return "LEFT"
		case 'H':
			return "HOME"
		case 'F':
			return "END"
		case 'Z':
			return "SHIFT_TAB"
		}
	}

	// Extended sequences: ESC [ N ~
	if len(data) >= 2 && data[len(data)-1] == '~' {
		switch string(data[:len(data)-1]) {
		case "1", "7":
			return "HOME"
		case "2":
			return "INSERT"
		case "3":
			return "DELETE"
		case "4", "8":
			return "END"
		case "5":
			return "PAGEUP"
		case "6":
			return "PAGEDOWN"
		case "15":
			return "F5"
		case "17":
			return "F6"
		case "18":
			return "F7"
		case "19":
			return "F8"
		case "20":
			return "F9"
		case "21":
			return "F10"
		case "23":
			return "F11"
		case "24":
			return "F12"
		}
	}

	// Modified arrow keys: ESC [ 1 ; mod X
	if len(data) >= 4 && data[0] == '1' && data[1] == ';' {
		mod := data[2]
		key := data[3]
		prefix := ""
		switch mod {
		case '2':
			prefix = "SHIFT_"
		case '3':
			prefix = "ALT_"
		case '5':
			prefix = "CTRL_"
		case '6':
			prefix = "CTRL_SHIFT_"
		}
		switch key {
		case 'A':
			return prefix + "UP"
		case 'B':
			return prefix + "DOWN"
		case 'C':
			return prefix + "RIGHT"
		case 'D':
			return prefix + "LEFT"
		}
	}

	return ""
}

func (k *Keyboard) parseSS3(b byte) string {
	switch b {
	case 'P':
		return "F1"
	case 'Q':
		return "F2"
	case 'R':
		return "F3"
	case 'S':
		return "F4"
	case 'H':
		return "HOME"
	case 'F':
		return "END"
	}
	return ""
}
