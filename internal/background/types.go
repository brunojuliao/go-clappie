package background

// App represents a background application.
type App struct {
	Name    string
	Path    string // path to the clapp directory
	Running bool
	Session string // tmux session name
}
