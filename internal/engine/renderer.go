package engine

import (
	"fmt"
	"strings"
)

// Renderer composes the full terminal output from view content.
type Renderer struct {
	daemon *Daemon
}

// NewRenderer creates a new renderer.
func NewRenderer(d *Daemon) *Renderer {
	return &Renderer{daemon: d}
}

// Compose builds the full screen output from a view instance.
func (r *Renderer) Compose(view *ViewInstance, width, height int) string {
	ctx := view.Context
	lines := ctx.GetLines()
	dims := r.calculateDimensions(width, height)

	var output strings.Builder

	// Header area
	r.renderHeader(&output, view, width, dims)

	// Content area
	r.renderContent(&output, lines, ctx, width, dims)

	// Footer / shortcuts
	r.renderFooter(&output, ctx, width, dims)

	return output.String()
}

// Dimensions holds calculated layout dimensions.
type Dimensions struct {
	HeaderHeight  int
	SkyHeight     int
	ContentHeight int
	FooterHeight  int
	ShortcutLines int
	ContentStart  int
	IsMobile      bool
}

func (r *Renderer) calculateDimensions(width, height int) Dimensions {
	isMobile := width < 120

	headerHeight := 2 // empty line + breadcrumbs
	headerGap := 1    // toast area
	skyHeight := 8
	skyGap := 1
	shortcutLines := 2
	footerHeight := 5 // garden

	if isMobile {
		skyHeight = 4
		skyGap = 0
	}

	contentHeight := height - headerHeight - headerGap - skyHeight - skyGap - shortcutLines - footerHeight
	if contentHeight < 1 {
		contentHeight = 1
	}

	contentStart := headerHeight + headerGap + skyHeight + skyGap

	return Dimensions{
		HeaderHeight:  headerHeight,
		SkyHeight:     skyHeight,
		ContentHeight: contentHeight,
		FooterHeight:  footerHeight,
		ShortcutLines: shortcutLines,
		ContentStart:  contentStart,
		IsMobile:      isMobile,
	}
}

func (r *Renderer) renderHeader(out *strings.Builder, view *ViewInstance, width int, dims Dimensions) {
	theme := r.daemon.theme

	// Empty line
	out.WriteString(theme.BG(""))
	out.WriteString(strings.Repeat(" ", width))
	out.WriteString("\x1b[0m\n")

	// Breadcrumb line
	breadcrumbs := r.buildBreadcrumbs()
	line := fmt.Sprintf(" %s", breadcrumbs)
	pad := width - VisualWidth(line)
	if pad < 0 {
		pad = 0
	}
	out.WriteString(theme.BG(theme.FG(line)))
	out.WriteString(strings.Repeat(" ", pad))
	out.WriteString("\x1b[0m\n")

	// Header gap (toast area)
	out.WriteString(theme.BG(""))
	out.WriteString(strings.Repeat(" ", width))
	out.WriteString("\x1b[0m\n")
}

func (r *Renderer) buildBreadcrumbs() string {
	var parts []string
	for _, v := range r.daemon.viewStack {
		parts = append(parts, v.Name)
	}
	if len(parts) > 3 {
		parts = append([]string{"..."}, parts[len(parts)-2:]...)
	}
	return strings.Join(parts, " > ")
}

func (r *Renderer) renderContent(out *strings.Builder, lines []string, ctx *Context, width int, dims Dimensions) {
	theme := r.daemon.theme

	// Calculate content width based on layout
	contentWidth := width
	if ctx.layout == "centered" && ctx.maxWidth > 0 {
		if ctx.maxWidth < contentWidth {
			contentWidth = ctx.maxWidth
		}
	}

	leftPad := 0
	if ctx.layout == "centered" && contentWidth < width {
		leftPad = (width - contentWidth) / 2
	}

	scrollTop := ctx.scrollTop
	if scrollTop > len(lines)-dims.ContentHeight {
		scrollTop = len(lines) - dims.ContentHeight
	}
	if scrollTop < 0 {
		scrollTop = 0
	}

	// Render visible lines
	for i := 0; i < dims.ContentHeight; i++ {
		lineIdx := scrollTop + i
		var lineContent string
		if lineIdx < len(lines) {
			lineContent = lines[lineIdx]
		}

		if leftPad > 0 {
			out.WriteString(theme.BG(""))
			out.WriteString(strings.Repeat(" ", leftPad))
		}

		lineWidth := VisualWidth(lineContent)
		remaining := width - leftPad - lineWidth
		if remaining < 0 {
			remaining = 0
		}

		out.WriteString(theme.BG(lineContent))
		out.WriteString(strings.Repeat(" ", remaining))
		out.WriteString("\x1b[0m\n")
	}

	// Scrollbar
	if len(lines) > dims.ContentHeight {
		// TODO: render scrollbar on right edge
	}
}

func (r *Renderer) renderFooter(out *strings.Builder, ctx *Context, width int, dims Dimensions) {
	theme := r.daemon.theme

	// Shortcuts bar
	shortcuts := r.buildShortcutBar(ctx)
	for _, line := range shortcuts {
		pad := width - VisualWidth(line)
		if pad < 0 {
			pad = 0
		}
		out.WriteString(theme.BG(line))
		out.WriteString(strings.Repeat(" ", pad))
		out.WriteString("\x1b[0m\n")
	}

	// Garden / footer area
	for i := 0; i < dims.FooterHeight; i++ {
		out.WriteString(theme.BG(""))
		out.WriteString(strings.Repeat(" ", width))
		out.WriteString("\x1b[0m\n")
	}
}

func (r *Renderer) buildShortcutBar(ctx *Context) []string {
	shortcuts := ctx.GetShortcuts()
	if len(shortcuts) == 0 {
		return []string{"", ""}
	}

	theme := r.daemon.theme
	var parts []string
	for _, sc := range shortcuts {
		part := fmt.Sprintf(" %s %s ", theme.Dim(sc.Key), sc.Label)
		parts = append(parts, part)
	}

	line := " " + strings.Join(parts, "  ")
	return []string{line, ""}
}
