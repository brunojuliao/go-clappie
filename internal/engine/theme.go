package engine

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// RGB represents an RGB color.
type RGB struct {
	R, G, B int
}

// Theme manages the color scheme.
type Theme struct {
	mode      string            // "dark" or "light"
	colors    map[string]RGB
	overrides map[string]RGB
	root      string
}

// Default color palettes.
var lightColors = map[string]RGB{
	"background":    {245, 243, 238},
	"primary":       {60, 100, 200},
	"text":          {30, 30, 30},
	"textMuted":     {120, 120, 120},
	"border":        {200, 195, 185},
	"divider":       {220, 215, 205},
	"success":       {40, 167, 69},
	"error":         {220, 53, 69},
	"warning":       {255, 193, 7},
	"info":          {23, 162, 184},
	"highlight":     {255, 248, 220},
	"accent":        {106, 90, 205},
	"crab":          {200, 150, 50},

	// Sky
	"sky":           {135, 195, 235},
	"cloudBright":   {255, 255, 255},
	"cloudMid":      {230, 235, 240},
	"sunCore":       {255, 220, 100},
	"skyTitle":      {255, 255, 255},
	"skyLead":       {200, 220, 240},

	// Garden
	"grass1":        {80, 160, 60},
	"grass2":        {60, 140, 45},
	"grass3":        {90, 170, 70},
	"grass4":        {70, 150, 55},
	"stem":          {60, 130, 40},
	"dirt":          {140, 110, 80},
	"pebbles":       {180, 170, 160},

	// Flowers
	"flowerWhite":   {255, 255, 255},
	"flowerRed":     {220, 60, 80},
	"flowerPurple":  {160, 100, 200},
	"flowerBlue":    {100, 140, 220},
	"flowerPink":    {240, 140, 180},
	"flowerOrange":  {240, 160, 60},

	// Sprites
	"spriteEyes":    {30, 30, 30},
}

var darkColors = map[string]RGB{
	"background":    {28, 28, 32},
	"primary":       {100, 140, 230},
	"text":          {240, 240, 240},
	"textMuted":     {140, 140, 150},
	"border":        {60, 60, 70},
	"divider":       {50, 50, 60},
	"success":       {60, 187, 89},
	"error":         {240, 73, 89},
	"warning":       {255, 213, 27},
	"info":          {43, 182, 204},
	"highlight":     {50, 48, 40},
	"accent":        {136, 120, 225},
	"crab":          {200, 150, 50},

	// Sky (null in JS = transparent for tmux bg)
	"sky":           {28, 28, 32},
	"cloudBright":   {60, 60, 70},
	"cloudMid":      {50, 50, 60},
	"sunCore":       {200, 180, 80},
	"skyTitle":      {240, 240, 240},
	"skyLead":       {140, 140, 150},

	// Garden (same as light)
	"grass1":        {80, 160, 60},
	"grass2":        {60, 140, 45},
	"grass3":        {90, 170, 70},
	"grass4":        {70, 150, 55},
	"stem":          {60, 130, 40},
	"dirt":          {140, 110, 80},
	"pebbles":       {180, 170, 160},

	// Flowers (same as light)
	"flowerWhite":   {255, 255, 255},
	"flowerRed":     {220, 60, 80},
	"flowerPurple":  {160, 100, 200},
	"flowerBlue":    {100, 140, 220},
	"flowerPink":    {240, 140, 180},
	"flowerOrange":  {240, 160, 60},

	// Sprites
	"spriteEyes":    {240, 240, 240},
}

// NewTheme creates a theme with default settings.
func NewTheme() *Theme {
	t := &Theme{
		mode:      "dark",
		overrides: make(map[string]RGB),
	}
	t.updateColors()
	return t
}

// InitFromRoot initializes theme from project settings.
func (t *Theme) InitFromRoot(root string) {
	t.root = root

	// Load mode from file
	modePath := filepath.Join(root, "recall", "settings", "theme", "mode.txt")
	data, err := os.ReadFile(modePath)
	if err == nil {
		mode := strings.TrimSpace(string(data))
		if mode == "light" || mode == "dark" {
			t.mode = mode
		}
	}

	// Load color overrides
	overridePath := filepath.Join(root, "recall", "settings", "theme", "colors.txt")
	data, err = os.ReadFile(overridePath)
	if err == nil {
		for _, line := range strings.Split(string(data), "\n") {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				continue
			}
			name := strings.TrimSpace(parts[0])
			hex := strings.TrimSpace(parts[1])
			if rgb, ok := hexToRGB(hex); ok {
				t.overrides[name] = rgb
			}
		}
	}

	t.updateColors()
}

func (t *Theme) updateColors() {
	if t.mode == "dark" {
		t.colors = make(map[string]RGB, len(darkColors))
		for k, v := range darkColors {
			t.colors[k] = v
		}
	} else {
		t.colors = make(map[string]RGB, len(lightColors))
		for k, v := range lightColors {
			t.colors[k] = v
		}
	}
	// Apply overrides
	for k, v := range t.overrides {
		t.colors[k] = v
	}
}

// IsDark returns true if the current theme is dark.
func (t *Theme) IsDark() bool {
	return t.mode == "dark"
}

// GetMode returns the theme mode name.
func (t *Theme) GetMode() string {
	return t.mode
}

// SetMode sets the theme mode and saves it.
func (t *Theme) SetMode(mode string) {
	t.mode = mode
	t.updateColors()

	if t.root != "" {
		dir := filepath.Join(t.root, "recall", "settings", "theme")
		os.MkdirAll(dir, 0755)
		os.WriteFile(filepath.Join(dir, "mode.txt"), []byte(mode), 0644)
	}
}

// Toggle switches between dark and light mode.
func (t *Theme) Toggle() {
	if t.mode == "dark" {
		t.SetMode("light")
	} else {
		t.SetMode("dark")
	}
}

// C returns the RGB color for a given name.
func (t *Theme) C(name string) RGB {
	if c, ok := t.colors[name]; ok {
		return c
	}
	return RGB{128, 128, 128}
}

func hexToRGB(hex string) (RGB, bool) {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return RGB{}, false
	}
	r, err := strconv.ParseInt(hex[0:2], 16, 32)
	if err != nil {
		return RGB{}, false
	}
	g, err := strconv.ParseInt(hex[2:4], 16, 32)
	if err != nil {
		return RGB{}, false
	}
	b, err := strconv.ParseInt(hex[4:6], 16, 32)
	if err != nil {
		return RGB{}, false
	}
	return RGB{int(r), int(g), int(b)}, true
}

