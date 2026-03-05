package engine

// ViewFactory creates a ScreenModel for a given view.
type ViewFactory func(data map[string]interface{}, styles *Styles, claudePane string) ScreenModel

// ViewModuleBT describes a registered bubbletea display view.
type ViewModuleBT struct {
	Create   ViewFactory
	Layout   string // "centered" or "full"
	MaxWidth int    // for centered layout
}
