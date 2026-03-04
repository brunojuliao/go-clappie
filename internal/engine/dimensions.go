package engine

// LayoutDimensions calculates all layout dimensions for the display.
type LayoutDimensions struct {
	DisplayWidth  int
	DisplayHeight int

	HeaderHeight int
	HeaderGap    int
	SkyHeight    int
	SkyGap       int
	ShortcutRows int
	FooterHeight int

	ContentHeight int
	ContentStart  int
	ContentWidth  int

	IsMobile bool
}

// CalculateLayout computes layout dimensions from terminal size and settings.
func CalculateLayout(width, height int, showScene bool) LayoutDimensions {
	isMobile := width < 120

	d := LayoutDimensions{
		DisplayWidth:  width,
		DisplayHeight: height,
		HeaderHeight:  2,  // empty line + breadcrumbs
		HeaderGap:     1,  // toast area
		ShortcutRows:  2,  // shortcut bar
		FooterHeight:  5,  // garden / footer
		IsMobile:      isMobile,
	}

	if showScene {
		if isMobile {
			d.SkyHeight = 4
			d.SkyGap = 0
		} else {
			d.SkyHeight = 8
			d.SkyGap = 1
		}
	}

	d.ContentStart = d.HeaderHeight + d.HeaderGap + d.SkyHeight + d.SkyGap
	d.ContentHeight = height - d.ContentStart - d.ShortcutRows - d.FooterHeight
	if d.ContentHeight < 1 {
		d.ContentHeight = 1
	}
	d.ContentWidth = width

	return d
}
