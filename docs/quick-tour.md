# go-clappie Quick Tour

Run these from your shell pane (left/bottom). The first `display push` auto-starts the daemon in a new tmux pane.

## 1. Stack views (push/pop)

```bash
go-clappie display push chores
go-clappie display pop
```

## 2. Demo screen with all 14 UI components

```bash
go-clappie display push example-demo-screens/all-components
```

Renders buttons, toggles, text inputs, checkboxes, radio groups, progress bars, loaders, and alerts all at once. Use Tab/Shift+Tab to cycle focus.

## 3. Interactive dialogs

```bash
go-clappie display push utility/confirm -d '{"message":"Deploy to prod?"}'
go-clappie display push utility/list -d '{"title":"Pick a model","options":["Claude 4","Sonnet","Haiku"]}'
go-clappie display push utility/editor -d '{"value":"Edit me"}'
go-clappie display push utility/viewer -d '{"content":"Line 1\nLine 2\nLine 3"}'
```

## 4. Toast notifications

```bash
go-clappie display toast "Hello from Windows!" -t 5000
```

## 5. Close it all

```bash
go-clappie display close
```
