package displays

import (
	"fmt"

	"github.com/brunojuliao/go-clappie/internal/engine"
	"github.com/brunojuliao/go-clappie/internal/filestore"
	"github.com/brunojuliao/go-clappie/internal/platform"
	"github.com/brunojuliao/go-clappie/internal/uikit"
)

// NewChoresView creates the chores approval queue view.
func NewChoresView(ctx *engine.Context) engine.View {
	view := uikit.NewView(ctx)

	ctx.SetTitle("Chores")
	ctx.SetDescription("Human approval queue")

	var choresList []choreItem
	selectedIdx := 0

	loadChores := func() {
		root, err := platform.ProjectRoot()
		if err != nil {
			return
		}

		dir := platform.ChoresHumansDir(root)
		entries, err := filestore.List(dir)
		if err != nil {
			return
		}

		choresList = nil
		for _, entry := range entries {
			body, blocks, err := filestore.ReadAndParse(entry.Path)
			if err != nil {
				continue
			}

			title := filestore.GetMetaField(blocks, "chore-meta", "title")
			if title == "" {
				title = entry.Name
			}
			status := filestore.GetMetaField(blocks, "chore-meta", "status")
			icon := filestore.GetMetaField(blocks, "chore-meta", "icon")
			summary := filestore.GetMetaField(blocks, "chore-meta", "summary")

			choresList = append(choresList, choreItem{
				Name:    entry.Name,
				Path:    entry.Path,
				Title:   title,
				Body:    body,
				Status:  status,
				Icon:    icon,
				Summary: summary,
			})
		}

		if selectedIdx >= len(choresList) {
			selectedIdx = len(choresList) - 1
		}
		if selectedIdx < 0 {
			selectedIdx = 0
		}
	}

	render := func() {
		var lines []string
		lines = append(lines, "")

		if len(choresList) == 0 {
			lines = append(lines, "  No pending chores.")
			lines = append(lines, "")
			lines = append(lines, engine.StyleDim("  All caught up!"))
		} else {
			lines = append(lines, fmt.Sprintf("  %d pending chores", len(choresList)))
			lines = append(lines, "")

			for i, c := range choresList {
				prefix := "  "
				if i == selectedIdx {
					prefix = "▸ "
				}
				icon := c.Icon
				if icon == "" {
					icon = "📋"
				}
				line := fmt.Sprintf("%s%s %s", prefix, icon, c.Title)
				if c.Summary != "" {
					line += engine.StyleDim(" — "+c.Summary)
				}
				if i == selectedIdx {
					line = engine.StyleBold(line)
				}
				lines = append(lines, line)
			}

			// Show selected chore details
			if selectedIdx < len(choresList) {
				lines = append(lines, "")
				lines = append(lines, engine.StyleBold("  ─── Details ───"))
				lines = append(lines, "")
				if choresList[selectedIdx].Body != "" {
					lines = append(lines, "  "+choresList[selectedIdx].Body)
				}
			}
		}

		ctx.Draw(lines)
	}

	view.RegisterShortcut("A", "Approve", func() {
		if selectedIdx < len(choresList) {
			c := choresList[selectedIdx]
			ctx.Submit(fmt.Sprintf("[clappie] Chore approved → %s", c.Title))
			// Update status
			_, blocks, _ := filestore.ReadAndParse(c.Path)
			filestore.SetMetaField(&blocks, "chore-meta", "status", "approved")
			filestore.WriteWithMeta(c.Path, c.Body, blocks)
			loadChores()
			render()
		}
	})

	view.RegisterShortcut("X", "Reject", func() {
		if selectedIdx < len(choresList) {
			c := choresList[selectedIdx]
			ctx.Submit(fmt.Sprintf("[clappie] Chore rejected → %s", c.Title))
			_, blocks, _ := filestore.ReadAndParse(c.Path)
			filestore.SetMetaField(&blocks, "chore-meta", "status", "rejected")
			filestore.WriteWithMeta(c.Path, c.Body, blocks)
			loadChores()
			render()
		}
	})

	return engine.View{
		Init: func() {
			loadChores()
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
				if selectedIdx < len(choresList)-1 {
					selectedIdx++
					render()
				}
				return true
			}
			return view.HandleKey(key)
		},
		OnScroll: func(dir int) {
			if dir < 0 && selectedIdx > 0 {
				selectedIdx--
			} else if dir > 0 && selectedIdx < len(choresList)-1 {
				selectedIdx++
			}
			render()
		},
	}
}

type choreItem struct {
	Name    string
	Path    string
	Title   string
	Body    string
	Status  string
	Icon    string
	Summary string
}
