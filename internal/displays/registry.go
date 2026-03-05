package displays

import (
	"sort"

	"github.com/brunojuliao/go-clappie/internal/engine"
)

// Registry maps view names to their bubbletea module definitions.
var Registry = map[string]engine.ViewModuleBT{}

// Register adds a view module to the registry.
func Register(name string, module engine.ViewModuleBT) {
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
	Register("heartbeat", engine.ViewModuleBT{
		Create: NewHeartbeatScreen,
		Layout: "full",
	})
	Register("chores", engine.ViewModuleBT{
		Create:   NewChoresScreen,
		Layout:   "centered",
		MaxWidth: 60,
	})
	Register("notifications", engine.ViewModuleBT{
		Create:   NewNotificationsScreen,
		Layout:   "centered",
		MaxWidth: 60,
	})
	Register("sidekicks", engine.ViewModuleBT{
		Create:   NewSidekicksScreen,
		Layout:   "centered",
		MaxWidth: 70,
	})
	Register("background", engine.ViewModuleBT{
		Create:   NewBackgroundScreen,
		Layout:   "centered",
		MaxWidth: 60,
	})
	Register("oauth", engine.ViewModuleBT{
		Create:   NewOAuthScreen,
		Layout:   "centered",
		MaxWidth: 60,
	})
	Register("parties", engine.ViewModuleBT{
		Create:   NewPartiesIndexScreen,
		Layout:   "centered",
		MaxWidth: 60,
	})
	Register("parties/status", engine.ViewModuleBT{
		Create:   NewPartiesStatusScreen,
		Layout:   "centered",
		MaxWidth: 70,
	})
	Register("projects", engine.ViewModuleBT{
		Create:   NewProjectsScreen,
		Layout:   "centered",
		MaxWidth: 60,
	})

	// Utility displays
	Register("utility/list", engine.ViewModuleBT{
		Create:   NewUtilityListScreen,
		Layout:   "centered",
		MaxWidth: 50,
	})
	Register("utility/confirm", engine.ViewModuleBT{
		Create:   NewUtilityConfirmScreen,
		Layout:   "centered",
		MaxWidth: 50,
	})
	Register("utility/editor", engine.ViewModuleBT{
		Create: NewUtilityEditorScreen,
		Layout: "full",
	})
	Register("utility/viewer", engine.ViewModuleBT{
		Create: NewUtilityViewerScreen,
		Layout: "full",
	})

	// Demo screens
	Register("example-demo-screens/hello-world", engine.ViewModuleBT{
		Create:   NewHelloWorldScreen,
		Layout:   "centered",
		MaxWidth: 50,
	})
	Register("example-demo-screens/all-components", engine.ViewModuleBT{
		Create:   NewAllComponentsScreen,
		Layout:   "centered",
		MaxWidth: 60,
	})
}
