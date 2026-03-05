# Go-Clappie

A Go port of [Clappie](https://github.com/whatnickcodes/clappie) — the personal assistant framework that runs inside tmux and communicates with Claude Code.

Single binary. No runtime dependencies. Cross-platform.

## Overview

Go-Clappie is a complete rewrite of the original JavaScript/Bun-based Clappie framework in Go. It preserves the same architecture, file formats, and tmux-based workflow while compiling down to a single ~11MB binary.

### What It Does

- Runs a TUI display engine inside tmux panes
- Communicates with Claude Code via `[go-clappie]` messages typed into tmux
- File-based state management (no database) using `.txt` files with `[meta]` blocks
- CLI sends commands to a daemon process via Unix sockets (TCP fallback on Windows)

### Subsystems

| Subsystem | Description |
|-----------|-------------|
| **Display Engine** | Bubbletea-powered TUI with view stack, keyboard/mouse input, themes |
| **UI Kit** | 14 components: buttons, toggles, text inputs, textareas, checkboxes, radios, selects, progress bars, loaders, labels, dividers, alerts |
| **Chores** | Human approval queue for AI-proposed actions |
| **Notifications** | Bidirectional sync pipeline with aggressive TLDR |
| **Heartbeat** | AI-powered cron scheduler with interval-based checks |
| **Background** | Long-running app management via tmux sessions |
| **Sidekicks** | Autonomous Claude sessions with webhook routing |
| **Parties** | Gamified AI swarm simulations with dice, ledgers, identities |
| **OAuth** | Shared token management with PKCE auth flows |
| **Skills** | Skill discovery from `.claude/skills/*/` directories |
| **Graphics** | Unicode quarter-block pixel art and animated scenes |
| **Telegram** | Bidirectional messaging via [telegram-bridge](https://github.com/brunojuliao/go-clappie-telegram-bridge) (separate binary) |

## Requirements

- **Go 1.25+** (for building from source)
- **tmux** (required at runtime)
- **winpty** (required on Windows/MSYS2 for TUI apps inside tmux)
- A terminal emulator with ANSI color support

## Installation

### From Source

```bash
git clone https://github.com/brunojuliao/go-clappie.git
cd go-clappie
go build -o go-clappie .
```

Or install directly:

```bash
go install github.com/brunojuliao/go-clappie@latest
```

### Pre-built Binaries

Download from [Releases](https://github.com/brunojuliao/go-clappie/releases) for your platform.

## Building

### Single Target

```bash
make build        # Build for current OS/arch
make install      # go install
```

### All Targets

```bash
make build-all    # Cross-compile all 5 targets
```

This produces binaries in `dist/`:

| Target | File | OS | Architecture |
|--------|------|----|-------------|
| Linux x86_64 | `go-clappie-linux-amd64` | Linux | amd64 |
| Linux ARM | `go-clappie-linux-arm64` | Linux | arm64 |
| macOS Intel | `go-clappie-darwin-amd64` | macOS | amd64 |
| macOS Apple Silicon | `go-clappie-darwin-arm64` | macOS | arm64 |
| Windows | `go-clappie-windows-amd64.exe` | Windows | amd64 |

Rename the binary for your platform to `go-clappie` (or `go-clappie.exe` on Windows) and place it somewhere in your `PATH`.

### Other Commands

```bash
make test         # Run all tests
make vet          # Static analysis
make fmt          # Format code
make lint         # Run golangci-lint
make clean        # Remove build artifacts
```

## Supported Platforms

### Linux

Works out of the box. tmux is available in all major package managers:

```bash
# Debian/Ubuntu
sudo apt install tmux

# Arch
sudo pacman -S tmux

# Fedora
sudo dnf install tmux
```

### macOS

Works out of the box. Install tmux via Homebrew:

```bash
brew install tmux
```

### Windows

Go-Clappie on Windows requires **Git for Windows** (Git Bash) with **tmux** and **winpty** installed manually. Neither is included with Git Bash by default.

#### Option A: Via MSYS2 (Recommended)

1. **Install [MSYS2](https://www.msys2.org/)** if you don't have it already.

2. **Install tmux and winpty in MSYS2:**
   ```bash
   pacman -S tmux winpty
   ```

3. **Copy the required files** from MSYS2 to Git for Windows:

   From `C:\msys64\usr\bin\` copy these files to `C:\Program Files\Git\usr\bin\`:
   - `tmux.exe`
   - `msys-event-2-1-*.dll` (e.g., `msys-event-2-1-7.dll`)
   - `msys-event_core-2-1-*.dll` (if present)
   - `winpty.exe`
   - `winpty-agent.exe`

4. **Restart Git Bash** and verify:
   ```bash
   tmux -V
   winpty --version
   ```

#### Option B: Direct Package Download

Download packages from the [MSYS2 repository](https://repo.msys2.org/msys/x86_64/):

1. **tmux:**
   - `tmux-*.pkg.tar.zst`
   - `libevent-*.pkg.tar.xz` (dependency)

2. **winpty:**
   - `winpty-*.pkg.tar.zst`

3. **Extract and copy** the binaries to `C:\Program Files\Git\usr\bin\`:
   - `tmux.exe` + `msys-event*.dll` (from tmux and libevent packages)
   - `winpty.exe` + `winpty-agent.exe` (from winpty package)

#### Why winpty?

Native Windows binaries (like Claude Code and go-clappie's display daemon) can't interact with MSYS2 PTYs directly. `winpty` bridges the gap by creating a real Windows console. Go-clappie auto-detects winpty and uses it when available.

#### Important Notes for Windows

- tmux **only works with MinTTY** (`git-bash.exe`). It will not work with `cmd.exe`, PowerShell, or Windows Terminal running `bash.exe` directly.
- Make sure Git for Windows and MSYS2 are both **64-bit** installations.
- IPC uses TCP localhost (port derived from socket path hash) instead of Unix sockets on Windows.

## Usage

Go-Clappie runs inside a tmux session. Basic commands:

```bash
# Display management
go-clappie display push heartbeat       # Open the heartbeat dashboard
go-clappie display push chores          # Open the chore approval queue
go-clappie display pop                  # Close current view
go-clappie display list                 # List open displays

# Background apps
go-clappie background start <app>       # Start a background app
go-clappie background list              # List running apps

# Sidekicks
go-clappie sidekick spawn <name>        # Spawn an autonomous Claude session
go-clappie sidekick list                # List active sidekicks

# Parties (AI simulations)
go-clappie parties init <game>          # Initialize a game
go-clappie parties roll 2d6+3           # Roll dice

# OAuth
go-clappie oauth auth <provider>        # Start OAuth flow
go-clappie oauth status                 # Show token status

# General
go-clappie list displays                # List all 15 built-in views
go-clappie list skills                  # List discovered skills
go-clappie --help                       # Full command reference
```

## Project Structure

```
go-clappie/
  main.go                     # Entry point
  Makefile                    # Build targets
  cmd/                        # Cobra CLI commands (9 files)
  internal/
    platform/                 # Cross-platform abstractions (IPC, paths, tmux)
    ipc/                      # JSON command/response protocol over sockets
    tmux/                     # tmux session/pane/sendkeys operations
    engine/                   # Display engine daemon (bubbletea app, view stack, styles, theme)
    uikit/                    # 14 UI components
    graphics/                 # Quarter-block pixel art and animations
    filestore/                # File-based state with [meta] block format
    chores/                   # Human approval queue
    notifications/            # Bidirectional sync pipeline
    heartbeat/                # AI-powered cron scheduler
    background/               # Long-running app management
    sidekicks/                # Autonomous Claude sessions + webhook server
    parties/                  # AI swarm simulations (dice, ledgers, identities)
    oauth/                    # Token management with PKCE
    skills/                   # Skill discovery
    displays/                 # 15 built-in display views
```

## Dependencies

| Package | Purpose |
|---------|---------|
| [github.com/spf13/cobra](https://github.com/spf13/cobra) | CLI framework |
| [github.com/charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea) | TUI framework (Elm architecture) |
| [github.com/charmbracelet/lipgloss](https://github.com/charmbracelet/lipgloss) | Terminal styling and layout |
| [github.com/charmbracelet/bubbles](https://github.com/charmbracelet/bubbles) | TUI components (text input, viewport, spinner) |
| [github.com/mattn/go-runewidth](https://github.com/mattn/go-runewidth) | Visual width calculation (emoji, CJK) |

Everything else uses the Go standard library.

## Testing

```bash
go test ./...
```

Test suites cover: filestore meta parsing, IPC round-trip, ANSI width calculation, dice rolling, heartbeat scheduling, and file CRUD operations.

## Architecture

### Adaptive IPC

- **Linux/macOS:** Unix domain sockets at `/tmp/go-clappie-{TMUX_PANE}.sock`
- **Windows:** TCP on `127.0.0.1:{port}` where port is derived from an FNV hash of the socket path

### Display Daemon

The CLI spawns itself as a daemon process (`go-clappie __daemon`). The daemon:
1. Creates a bubbletea `Program` with alt-screen and mouse support
2. Runs an IPC server on a goroutine
3. Bubbletea manages the event loop (keyboard, mouse, window resize, custom messages)
4. IPC mutations are thread-safe via `program.Send()` into the Elm update cycle

### File Format

All state is stored as `.txt` files with optional metadata blocks:

```
Body content goes here

---
[meta]
key: value
status: pending
created: 2025-03-03 14:30
```

### Import Cycle Resolution

The `engine` package cannot import `displays` (circular dependency). Solution: the view registry is injected via `AppConfig.Registry` at startup from `cmd/daemon.go`.

## Related Projects

| Project | Description |
|---------|-------------|
| [go-clappie-telegram-bridge](https://github.com/brunojuliao/go-clappie-telegram-bridge) | Lightweight Telegram bot that plugs into go-clappie's file-based ecosystem. Incoming messages land in `notifications/dirty/`, Claude replies via `notifications/outbox/`. Runs as a companion binary in a tmux window alongside Claude. |
| [Clappie](https://github.com/whatnickcodes/clappie) | The original JavaScript/Bun framework that go-clappie is a full Go port of. |

## Credits

Based on the original [Clappie](https://github.com/whatnickcodes/clappie) by [@whatnickcodes](https://github.com/whatnickcodes).

Co-created by [Bruno Juliao](https://github.com/brunojuliao) and [Claude](https://claude.ai) (Anthropic's AI assistant). The entire Go codebase — 102 source files across 13 subsystems — was pair-programmed from the ground up with Claude Code.

## License

See [LICENSE](LICENSE) for details.
