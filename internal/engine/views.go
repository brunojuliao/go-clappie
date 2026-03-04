package engine

// ViewModule describes a registered display view.
type ViewModule struct {
	Create   func(ctx *Context) View
	Layout   string // "centered" or "full"
	MaxWidth int    // for centered layout
}

// ViewInstance is a running instance of a view.
type ViewInstance struct {
	Name     string
	Module   ViewModule
	Context  *Context
	Instance View
}

// View is the interface that display views implement.
type View struct {
	Init     func()
	Render   func()
	OnKey    func(key string) bool
	OnClick  func(x, y int)
	OnScroll func(direction int)
	Cleanup  func()
}
