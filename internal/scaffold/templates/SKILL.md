---
name: go-clappie
description: >
  Go-Clappie is a digital assistant layer that turns Claude Code into a full personal assistant
  with interactive terminal UIs. It manages everything through one CLI: `go-clappie`.

  This is the core skill for the entire project. Load it whenever the user asks for ANYTHING
  personal-assistant related — not just when they say "go-clappie". Examples: emails, texts,
  notifications, inbox, dashboard, sidekicks, chores, heartbeat, background apps, displays,
  parties, memory, messages, todos, approval queue, status, automation, or any `[go-clappie]` prefixed message.
  If the user says "open notifications" or "spawn a sidekick" — that's this skill.
  Don't guess at how these systems work — load this skill, it has the docs.

  CORE SYSTEM - Terminal display engine with push/pop view navigation, mouse click handling,
  keyboard shortcuts, and two-way communication. Views send structured messages back to Claude
  prefixed with [go-clappie].

  SUBSYSTEMS:

  - Sidekicks: Autonomous AI agents running in background tmux sessions. Spawn with prompts,
    chain tasks, report back. Used for async work and background processing.

  - Chores: Human approval queue. High-stakes actions get drafted as chore files in
    chores/humans/ for user review before execution.

  - Heartbeat: AI-powered cron scheduler. Periodic check files in chores/bots/ trigger
    subagent spawns at configurable intervals.

  - Notifications: Bidirectional sync system. External events dump to notifications/dirty/,
    AI processor creates curated items in notifications/clean/.

  - Background: Long-running apps managed via .background marker files.

  - Display Engine: Terminal UI framework with 14 components (Button, Toggle, TextInput,
    Textarea, Checkbox, Radio, Select, Progress, Loader, Alert, Label, Divider, etc).

  - Parties: Gamified AI swarm simulations. Define games with player roles and rules,
    spawn multiple AI agents that interact according to game mechanics.

  - OAuth: Shared token management across skills. Auth flows, auto-refresh, token storage.

  - Skills: Skill discovery from .claude/skills/*/ directories.

  COMMANDS: go-clappie list, go-clappie display push/pop/toast/close/kill,
  go-clappie background start/stop/list/kill, go-clappie sidekick spawn/send/complete/report/end/kill,
  go-clappie parties games/init/launch/end/show/rules/set/get/roll,
  go-clappie oauth auth/token/status/refresh/revoke,
  go-clappie kill.
---

# Go-Clappie

Turns Claude Code into a digital assistant with interactive terminal UIs. Push views, handle user input, communicate back to Claude.

## Be Proactive

**Don't tell the user to run commands - just run them yourself.**

When the user says "open notifications", "show me my emails", "view the dashboard", etc. - just run `go-clappie display push <name>`. Act autonomously.

**Display name lookup:** Run `go-clappie list displays` if you're not certain. Common mappings:

| User says | Display name |
|-----------|-------------|
| "background manager", "background apps" | `background` |
| "sidekicks", "active sidekicks" | `sidekicks` |
| "notifications", "inbox" | `notifications` |
| "chores", "todo", "approval queue" | `chores` |
| "heartbeat", "dashboard", "status" | `heartbeat` |
| "projects", "scratch pad" | `projects` |
| "oauth", "tokens" | `oauth` |
| "parties", "games" | `parties` |

**All registered views:**

| View Name | Layout | Description |
|-----------|--------|-------------|
| `heartbeat` | full | AI-powered cron dashboard |
| `chores` | centered | Human approval queue |
| `notifications` | centered | Curated notification feed |
| `sidekicks` | centered | Autonomous agent management |
| `background` | centered | Long-running app management |
| `oauth` | centered | Token management |
| `parties` | centered | Party game index |
| `parties/status` | centered | Active party status |
| `projects` | centered | Project workspace |
| `utility/list` | centered | List picker dialog |
| `utility/confirm` | centered | Yes/No confirmation dialog |
| `utility/editor` | full | Multi-line text editor |
| `utility/viewer` | full | Read-only text viewer |
| `example-demo-screens/hello-world` | centered | Hello world demo |
| `example-demo-screens/all-components` | centered | All UI components demo |

## Data Directories

| Directory | Purpose |
|-----------|---------|
| `chores/humans/` | Human approval queue (pending chores) |
| `chores/bots/` | Heartbeat check files |
| `notifications/dirty/` | Raw sync stream from integrations |
| `notifications/clean/` | Curated items for user review |
| `recall/memory/` | Persistent memory files |
| `recall/logs/` | All logs (flat subdirectories) |
| `recall/settings/` | Runtime settings per skill |
| `recall/sidekicks/` | Sidekick session records |
| `recall/parties/` | Party game data |
| `projects/` | Workspace for apps, sites, repos |

## Logging Rules

All logs live in `recall/logs/` with flat structure:

```
recall/logs/
├── chores/        # YYYY-MM-DD-HHMM-name.txt
├── heartbeat/     # YYYY-MM-DD.txt (daily, append [HH:MM] entries)
├── sidekicks/     # YYYY-MM-DD-HHMM-source-slug.txt
└── notifications/ # YYYY-MM-DD.txt (daily processing log)
```

**Rules:** No files at root of `recall/logs/`. No nested subdirectories. No `.log` extension — always `.txt`. No improvised paths.

## Commands

```bash
# Discovery
go-clappie list                            # List everything
go-clappie list skills                     # Just skills and commands
go-clappie list displays                   # Just displays + navigation

# Displays
go-clappie display push <view> [options]   # Push view onto stack
go-clappie display pop                     # Go back
go-clappie display toast "<msg>" [-t ms]   # Toast notification
go-clappie display close                   # Close display
go-clappie display list                    # List running instances
go-clappie display kill                    # Kill displays only

# Background
go-clappie background start [app]          # Start a background app
go-clappie background stop [app]           # Stop a background app
go-clappie background list                 # List apps + status
go-clappie background kill                 # Kill background only

# Sidekicks
go-clappie sidekick spawn "prompt"         # Spawn with task
go-clappie sidekick send "message"         # Send to active sidekick
go-clappie sidekick complete "summary"     # End sidekick with summary
go-clappie sidekick report "message"       # Report to main Claude
go-clappie sidekick end                    # End active sidekick
go-clappie sidekick message <id> "text"    # DM a specific sidekick
go-clappie sidekick broadcast "text"       # Message all sidekicks

# Parties
go-clappie parties init <game> [context]   # Create ledger
go-clappie parties set <key> <val>         # Set shared state
go-clappie parties set <who> <key> <val>   # Set player state
go-clappie parties get [who] [key]         # Read state
go-clappie parties roll <spec>             # 1d6, 2d20, coin, pick "a,b,c"

# OAuth
go-clappie oauth auth <provider>           # Start OAuth flow
go-clappie oauth token <provider>          # Get access token
go-clappie oauth status                    # Show all tokens
go-clappie oauth refresh <provider>        # Force refresh
go-clappie oauth revoke <provider>         # Delete tokens

# Kill everything
go-clappie kill                            # Displays + background + sidekicks
```

## Options

```bash
-f, --focus         # Switch focus to display pane (default: stay in chat)
-d key=value        # Pass data (repeatable)
-d body=@/tmp/f.txt # Pass file contents
-t, --timeout <ms>  # Toast duration
```

## View Naming

```bash
heartbeat                         # Built-in view
heartbeat/dashboard               # Nested view
example-demo-screens/hello-world  # Demo screen
```

- No slash = top-level view
- Has slash = nested view

## [go-clappie] Messages

Views communicate back to Claude by typing into the chat window via tmux. All messages are prefixed with `[go-clappie]`.

**ctx.Submit()** — types + Enter (hard submit):
```
[go-clappie] Counter → 5
[go-clappie] Toggle → yes
[go-clappie] Confirm → yes
[go-clappie] List → selected-item
[go-clappie] Editor → edited text content
[go-clappie] Chore approved → chore-title
[go-clappie] Chore rejected → chore-title
[go-clappie] Sidekick complete → summary
[go-clappie] Sidekick report → message
[go-clappie] State changed → key = value
```

**ctx.Send()** — types only, no Enter (user can review).

Use arrow format, not JSON. It's going to a chat window — keep it human-readable.

## Data Formats (-d flag)

```bash
-d key=value              # Simple string
-d key=@/path/to/file     # Read from file (for long content)
-d '{"json": true}'       # JSON (must start with { or [)
```

**For long text:** Write to temp file first, then `-d body=@/tmp/draft.txt`.

## Subsystem Summaries

### Sidekicks
Autonomous AI agents in background tmux sessions. Spawn with prompts, chain work. When a sidekick completes, it sends `[go-clappie] Sidekick complete → summary` back to Claude.

### Chores
Human approval queue for high-stakes actions. Files in `chores/humans/` use `.txt` format with `[meta]` blocks containing title, status, icon, and summary fields. When you receive `[go-clappie] Chore approved → <title>`, read the chore file and execute it.

### Notifications
Bidirectional sync: `dirty/` (raw dumps) → `clean/` (curated items). Aggressive TLDR consolidation — 4-6 clean items, not 20. `source_id` connects items for lifecycle tracking.

### Heartbeat
AI-powered cron. Check files in `chores/bots/`. Spawns subagents at intervals. When you receive `[go-clappie] Heartbeat initiated`, read check files, spawn agents, update metadata, log results.

### Display Engine
Terminal UI framework. Views rendered via `Context.Draw()`. Components: Button, Toggle, TextInput, Textarea, Checkbox, Radio, Select, Progress, Loader, Alert, Label, Divider, SectionHeading, View.

### Background
Long-running apps discovered via `.background` marker files in `.claude/skills/go-clappie/clapps/` directories.

### OAuth
Shared token management across skills. Auth flows via local callback server, auto-refresh, PKCE support.

### Parties
Gamified AI swarm simulations layered on sidekicks. Define games with player roles, rules, and state. Launch simulations with shared ledgers.

### Memory
Save to `recall/memory/` whenever you learn something about the user. One fact per line, terse, newest at bottom. Files: `personal.txt`, `preferences.txt`, `work.txt`, `people.txt`, `technical.txt`, or create new ones as needed.

## File Format

All state is stored as `.txt` files with optional metadata blocks:

```
Body content goes here

---
[meta-name]
key: value
status: pending
created: 2025-03-03 14:30
```

## Requirements

- **tmux** — Must run inside tmux
- **go-clappie** binary on PATH
