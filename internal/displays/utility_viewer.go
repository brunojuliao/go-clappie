package displays

import (
	"strings"

	"github.com/brunojuliao/go-clappie/internal/engine"
	"github.com/brunojuliao/go-clappie/internal/filestore"
	"github.com/brunojuliao/go-clappie/internal/uikit"
)

// NewUtilityViewerView creates the utility file viewer view.
func NewUtilityViewerView(ctx *engine.Context) engine.View {
	view := uikit.NewView(ctx)

	ctx.SetTitle("Viewer")
	ctx.SetLayout("full", 0)

	content := ""
	if ctx.Data != nil {
		if c, ok := ctx.Data["content"].(string); ok {
			content = c
		}
		if path, ok := ctx.Data["path"].(string); ok {
			data, err := filestore.ReadFile(path)
			if err == nil {
				content = data
			}
		}
	}

	lines := strings.Split(content, "\n")
	scrollTop := 0

	render := func() {
		var output []string
		output = append(output, "")
		for _, line := range lines {
			output = append(output, "  "+line)
		}
		ctx.SetScrollTop(scrollTop)
		ctx.Draw(output)
	}

	return engine.View{
		Init:   render,
		Render: render,
		OnKey: func(key string) bool {
			return view.HandleKey(key)
		},
		OnScroll: func(dir int) {
			scrollTop += dir * 3
			if scrollTop < 0 {
				scrollTop = 0
			}
			max := len(lines) - 10
			if max < 0 {
				max = 0
			}
			if scrollTop > max {
				scrollTop = max
			}
			render()
		},
	}
}
