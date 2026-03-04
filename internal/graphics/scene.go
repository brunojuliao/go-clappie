package graphics

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"
)

// SkyScene renders an animated sky with sun/moon and clouds.
type SkyScene struct {
	width  int
	height int // in pixel rows (each terminal row = 2 pixel rows)
	frame  int
	clouds []cloud
	stars  []star
}

type cloud struct {
	x, y   float64
	speed  float64
	width  int
	height int
}

type star struct {
	x, y      int
	brightness float64
	twinkle   float64
}

// NewSkyScene creates a new sky scene animation.
func NewSkyScene(width, height int) *SkyScene {
	s := &SkyScene{
		width:  width * 2, // pixel width
		height: height * 2, // pixel height
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generate clouds
	for i := 0; i < 3; i++ {
		s.clouds = append(s.clouds, cloud{
			x:      float64(rng.Intn(s.width)),
			y:      float64(2 + rng.Intn(s.height/2)),
			speed:  0.2 + rng.Float64()*0.3,
			width:  8 + rng.Intn(12),
			height: 2 + rng.Intn(2),
		})
	}

	// Generate stars for dark mode
	for i := 0; i < 20; i++ {
		s.stars = append(s.stars, star{
			x:          rng.Intn(s.width),
			y:          rng.Intn(s.height),
			brightness: 0.5 + rng.Float64()*0.5,
			twinkle:    rng.Float64() * math.Pi * 2,
		})
	}

	return s
}

// RenderLight renders the sky scene for light theme.
func (s *SkyScene) RenderLight(bgR, bgG, bgB int) []string {
	s.frame++
	canvas := NewCanvas(s.width, s.height)

	// Sky gradient
	for y := 0; y < s.height; y++ {
		ratio := float64(y) / float64(s.height)
		r := int(135 + ratio*40)
		g := int(195 + ratio*20)
		b := int(235 - ratio*20)
		for x := 0; x < s.width; x++ {
			canvas.Set(x, y, r, g, b)
		}
	}

	// Sun
	sunX := s.width * 3 / 4
	sunY := s.height / 3
	for dy := -4; dy <= 4; dy++ {
		for dx := -4; dx <= 4; dx++ {
			dist := math.Sqrt(float64(dx*dx + dy*dy))
			if dist <= 4 {
				brightness := 1.0 - dist/6
				r := int(255 * brightness)
				g := int(220 * brightness)
				b := int(100 * brightness)
				canvas.Set(sunX+dx, sunY+dy, r, g, b)
			}
		}
	}

	// Clouds
	for i := range s.clouds {
		s.clouds[i].x += s.clouds[i].speed
		if s.clouds[i].x > float64(s.width+20) {
			s.clouds[i].x = -float64(s.clouds[i].width)
		}
		s.drawCloud(canvas, s.clouds[i], 255, 255, 255)
	}

	return canvas.Render(bgR, bgG, bgB)
}

// RenderDark renders the sky scene for dark theme.
func (s *SkyScene) RenderDark(bgR, bgG, bgB int) []string {
	s.frame++
	canvas := NewCanvas(s.width, s.height)

	// Dark sky (mostly transparent/bg colored)
	// Stars
	t := float64(s.frame) * 0.05
	for _, st := range s.stars {
		twinkle := 0.5 + 0.5*math.Sin(t+st.twinkle)
		brightness := st.brightness * twinkle
		c := int(200 * brightness)
		canvas.Set(st.x, st.y, c, c, c+30)
	}

	// Moon
	moonX := s.width * 3 / 4
	moonY := s.height / 3
	for dy := -3; dy <= 3; dy++ {
		for dx := -3; dx <= 3; dx++ {
			dist := math.Sqrt(float64(dx*dx + dy*dy))
			if dist <= 3 {
				// Crescent: darken one side
				if dx < -1 {
					continue
				}
				brightness := 1.0 - dist/5
				c := int(220 * brightness)
				canvas.Set(moonX+dx, moonY+dy, c, c, int(float64(c)*1.1))
			}
		}
	}

	return canvas.Render(bgR, bgG, bgB)
}

func (s *SkyScene) drawCloud(canvas *Canvas, c cloud, r, g, b int) {
	cx := int(c.x)
	cy := int(c.y)
	for dy := 0; dy < c.height; dy++ {
		for dx := 0; dx < c.width; dx++ {
			// Elliptical cloud shape
			fx := float64(dx) / float64(c.width) * 2 - 1
			fy := float64(dy) / float64(c.height) * 2 - 1
			if fx*fx+fy*fy <= 1 {
				canvas.Set(cx+dx, cy+dy, r, g, b)
			}
		}
	}
}

// CrabGarden renders the footer garden scene with animated crabs.
type CrabGarden struct {
	width int
	frame int
	crabs []crab
	rng   *rand.Rand
}

type crab struct {
	x     float64
	dir   int // -1 left, 1 right
	speed float64
	state int // animation state
}

// NewCrabGarden creates a new crab garden animation.
func NewCrabGarden(width int) *CrabGarden {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	g := &CrabGarden{
		width: width,
		rng:   rng,
	}

	// Create a few crabs
	for i := 0; i < 2; i++ {
		dir := 1
		if rng.Intn(2) == 0 {
			dir = -1
		}
		g.crabs = append(g.crabs, crab{
			x:     float64(rng.Intn(width - 10)),
			dir:   dir,
			speed: 0.1 + rng.Float64()*0.2,
			state: rng.Intn(4),
		})
	}

	return g
}

// Render renders the garden footer.
func (g *CrabGarden) Render(theme string) []string {
	g.frame++
	lines := make([]string, 5)

	// Grass line
	var grassLine strings.Builder
	for i := 0; i < g.width; i++ {
		grass := []string{"╌", "╌", "─", "╌", "~"}[i%5]
		r, gg, b := 80, 160, 60
		if i%3 == 0 {
			r, gg, b = 60, 140, 45
		}
		grassLine.WriteString(fmt.Sprintf("\x1b[38;2;%d;%d;%dm%s", r, gg, b, grass))
	}
	grassLine.WriteString("\x1b[0m")
	lines[0] = grassLine.String()

	// Crab line
	var crabLine strings.Builder
	crabPositions := make(map[int]string)
	for i := range g.crabs {
		g.crabs[i].x += float64(g.crabs[i].dir) * g.crabs[i].speed
		if g.crabs[i].x < 0 {
			g.crabs[i].dir = 1
		} else if g.crabs[i].x > float64(g.width-4) {
			g.crabs[i].dir = -1
		}
		pos := int(g.crabs[i].x)
		g.crabs[i].state = (g.frame / 8) % 4
		crabChar := g.crabSprite(g.crabs[i])
		crabPositions[pos] = crabChar
	}

	for i := 0; i < g.width; i++ {
		if sprite, ok := crabPositions[i]; ok {
			crabLine.WriteString(fmt.Sprintf("\x1b[38;2;200;150;50m%s\x1b[0m", sprite))
			// Skip next chars covered by sprite
		} else {
			crabLine.WriteString(" ")
		}
	}
	lines[1] = crabLine.String()

	// Dirt lines
	for i := 2; i < 5; i++ {
		var dirt strings.Builder
		for j := 0; j < g.width; j++ {
			r, gg, b := 140, 110, 80
			if (i+j)%7 == 0 {
				r, gg, b = 180, 170, 160 // pebble
			}
			dirt.WriteString(fmt.Sprintf("\x1b[48;2;%d;%d;%dm ", r, gg, b))
		}
		dirt.WriteString("\x1b[0m")
		lines[i] = dirt.String()
	}

	return lines
}

func (g *CrabGarden) crabSprite(c crab) string {
	// Simple ASCII crab with animation frames
	frames := []string{"🦀", "🦀", "🦀", "🦀"}
	if c.state < len(frames) {
		return frames[c.state]
	}
	return "🦀"
}
