package displays

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/brunojuliao/go-clappie/internal/engine"
	"github.com/brunojuliao/go-clappie/internal/platform"
	"github.com/brunojuliao/go-clappie/internal/uikit"
)

// NewProjectsView creates the projects workspace view.
func NewProjectsView(ctx *engine.Context) engine.View {
	view := uikit.NewView(ctx)

	ctx.SetTitle("Projects")
	ctx.SetDescription("Workspace manager")

	var projects []projectItem
	selectedIdx := 0

	load := func() {
		root, err := platform.ProjectRoot()
		if err != nil {
			return
		}

		projectsDir := platform.ProjectsDir(root)
		entries, err := os.ReadDir(projectsDir)
		if err != nil {
			return
		}

		projects = nil
		for _, e := range entries {
			if e.IsDir() {
				projects = append(projects, projectItem{
					Name: e.Name(),
					Path: filepath.Join(projectsDir, e.Name()),
				})
			}
		}
	}

	render := func() {
		var lines []string
		lines = append(lines, "")

		if len(projects) == 0 {
			lines = append(lines, "  No projects found.")
			lines = append(lines, "")
			lines = append(lines, engine.StyleDim("  Create directories in projects/ to get started."))
		} else {
			lines = append(lines, fmt.Sprintf("  %d projects", len(projects)))
			lines = append(lines, "")

			for i, p := range projects {
				prefix := "  "
				if i == selectedIdx {
					prefix = "▸ "
				}
				line := fmt.Sprintf("%s📁 %s", prefix, p.Name)
				if i == selectedIdx {
					line = engine.StyleBold(line)
				}
				lines = append(lines, line)
			}
		}

		ctx.Draw(lines)
	}

	return engine.View{
		Init: func() {
			load()
			render()
		},
		Render: render,
		OnKey: func(key string) bool {
			switch key {
			case "UP", "k":
				if selectedIdx > 0 {
					selectedIdx--
					render()
				}
				return true
			case "DOWN", "j":
				if selectedIdx < len(projects)-1 {
					selectedIdx++
					render()
				}
				return true
			}
			return view.HandleKey(key)
		},
	}
}

type projectItem struct {
	Name string
	Path string
}
