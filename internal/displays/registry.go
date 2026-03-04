package displays

import (
	"sort"

	"github.com/brunojuliao/go-clappie/internal/engine"
)

// Registry maps view names to their module definitions.
var Registry = map[string]engine.ViewModule{}

// Register adds a view module to the registry.
func Register(name string, module engine.ViewModule) {
	Registry[name] = module
}

// ListRegistered returns sorted names of all registered views.
func ListRegistered() []string {
	names := make([]string, 0, len(Registry))
	for name := range Registry {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func init() {
	// Register built-in views
	Register("heartbeat", engine.ViewModule{
		Create: NewHeartbeatView,
		Layout: "full",
	})
	Register("chores", engine.ViewModule{
		Create:   NewChoresView,
		Layout:   "centered",
		MaxWidth: 60,
	})
	Register("notifications", engine.ViewModule{
		Create:   NewNotificationsView,
		Layout:   "centered",
		MaxWidth: 60,
	})
	Register("sidekicks", engine.ViewModule{
		Create:   NewSidekicksView,
		Layout:   "centered",
		MaxWidth: 70,
	})
	Register("background", engine.ViewModule{
		Create:   NewBackgroundView,
		Layout:   "centered",
		MaxWidth: 60,
	})
	Register("oauth", engine.ViewModule{
		Create:   NewOAuthView,
		Layout:   "centered",
		MaxWidth: 60,
	})
	Register("parties", engine.ViewModule{
		Create:   NewPartiesIndexView,
		Layout:   "centered",
		MaxWidth: 60,
	})
	Register("parties/status", engine.ViewModule{
		Create:   NewPartiesStatusView,
		Layout:   "centered",
		MaxWidth: 70,
	})
	Register("projects", engine.ViewModule{
		Create:   NewProjectsView,
		Layout:   "centered",
		MaxWidth: 60,
	})

	// Utility displays
	Register("utility/list", engine.ViewModule{
		Create:   NewUtilityListView,
		Layout:   "centered",
		MaxWidth: 50,
	})
	Register("utility/confirm", engine.ViewModule{
		Create:   NewUtilityConfirmView,
		Layout:   "centered",
		MaxWidth: 50,
	})
	Register("utility/editor", engine.ViewModule{
		Create: NewUtilityEditorView,
		Layout: "full",
	})
	Register("utility/viewer", engine.ViewModule{
		Create: NewUtilityViewerView,
		Layout: "full",
	})

	// Demo screens
	Register("example-demo-screens/hello-world", engine.ViewModule{
		Create:   NewHelloWorldView,
		Layout:   "centered",
		MaxWidth: 50,
	})
	Register("example-demo-screens/all-components", engine.ViewModule{
		Create:   NewAllComponentsView,
		Layout:   "centered",
		MaxWidth: 60,
	})
}
