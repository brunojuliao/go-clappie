# go-clappie Quick Tour

Run these from your shell pane (left/bottom) while the daemon is running.

## 1. Stack views (push/pop)

```bash
go-clappie display push chores
go-clappie display pop
```

## 2. Demo screen with all 14 UI components

```bash
go-clappie display push example-demo-screens/all-components
```

Renders buttons, toggles, text inputs, checkboxes, radio groups, progress bars, loaders, and alerts all at once.

## 3. Interactive dialogs

```bash
go-clappie display push utility/confirm -d message="Deploy to prod?"
go-clappie display push utility/list -d title="Pick a model" -d '{"options":["Claude 4","Sonnet","Haiku"]}'
```

## 4. Toast notifications

```bash
go-clappie display toast "Hello from Windows!" -t 3000
```

## 5. Close it all

```bash
go-clappie display close
```
