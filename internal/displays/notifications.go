package displays

import (
	"fmt"

	"github.com/brunojuliao/go-clappie/internal/engine"
	"github.com/brunojuliao/go-clappie/internal/filestore"
	"github.com/brunojuliao/go-clappie/internal/platform"
	"github.com/brunojuliao/go-clappie/internal/uikit"
)

// NewNotificationsView creates the notifications inbox view.
func NewNotificationsView(ctx *engine.Context) engine.View {
	view := uikit.NewView(ctx)

	ctx.SetTitle("Notifications")
	ctx.SetDescription("Notification inbox")

	var items []notifItem
	selectedIdx := 0

	loadNotifs := func() {
		root, err := platform.ProjectRoot()
		if err != nil {
			return
		}

		dir := platform.NotificationsCleanDir(root)
		entries, err := filestore.List(dir)
		if err != nil {
			return
		}

		items = nil
		for _, entry := range entries {
			body, blocks, err := filestore.ReadAndParse(entry.Path)
			if err != nil {
				continue
			}
			sourceID := filestore.GetMetaField(blocks, "meta", "source_id")
			context := filestore.GetMetaField(blocks, "meta", "context")
			created := filestore.GetMetaField(blocks, "meta", "created")

			items = append(items, notifItem{
				Name:     entry.Name,
				Path:     entry.Path,
				Body:     body,
				SourceID: sourceID,
				Context:  context,
				Created:  created,
			})
		}
	}

	render := func() {
		var lines []string
		lines = append(lines, "")

		if len(items) == 0 {
			lines = append(lines, "  No notifications.")
			lines = append(lines, "")
			lines = append(lines, engine.StyleDim("  Inbox zero!"))
		} else {
			lines = append(lines, fmt.Sprintf("  %d notifications", len(items)))
			lines = append(lines, "")

			for i, item := range items {
				prefix := "  "
				if i == selectedIdx {
					prefix = "▸ "
				}
				line := fmt.Sprintf("%s%s", prefix, item.Body)
				if len(line) > 55 {
					line = line[:52] + "..."
				}
				if item.Context != "" {
					line += engine.StyleDim(fmt.Sprintf(" [%s]", item.Context))
				}
				if i == selectedIdx {
					line = engine.StyleBold(line)
				}
				lines = append(lines, line)
			}
		}

		ctx.Draw(lines)
	}

	view.RegisterShortcut("D", "Dismiss", func() {
		if selectedIdx < len(items) {
			filestore.DeleteFile(items[selectedIdx].Path)
			loadNotifs()
			render()
		}
	})

	view.RegisterShortcut("C", "Clear All", func() {
		for _, item := range items {
			filestore.DeleteFile(item.Path)
		}
		items = nil
		selectedIdx = 0
		render()
	})

	return engine.View{
		Init: func() {
			loadNotifs()
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
				if selectedIdx < len(items)-1 {
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
			} else if dir > 0 && selectedIdx < len(items)-1 {
				selectedIdx++
			}
			render()
		},
	}
}

type notifItem struct {
	Name     string
	Path     string
	Body     string
	SourceID string
	Context  string
	Created  string
}
