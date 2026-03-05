package displays

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/brunojuliao/go-clappie/internal/engine"
	"github.com/brunojuliao/go-clappie/internal/filestore"
	"github.com/brunojuliao/go-clappie/internal/platform"
)

type heartbeatCheck struct {
	Name     string
	Body     string
	Interval string
	LastRun  string
	Status   string
}

type heartbeatScreen struct {
	checks  []heartbeatCheck
	styles  *engine.Styles
}

func NewHeartbeatScreen(data map[string]interface{}, styles *engine.Styles, claudePane string) engine.ScreenModel {
	m := &heartbeatScreen{styles: styles}
	m.loadChecks()
	return m
}

func (m *heartbeatScreen) loadChecks() {
	root, err := platform.ProjectRoot()
	if err != nil {
		return
	}
	checksDir := platform.ChoresBotsDir(root)
	entries, err := filestore.List(checksDir)
	if err != nil {
		return
	}
	m.checks = nil
	for _, entry := range entries {
		body, blocks, err := filestore.ReadAndParse(entry.Path)
		if err != nil {
			continue
		}
		meta := filestore.GetMeta(blocks, "heartbeat-meta")
		interval, lastRun, status := "", "", ""
		if meta != nil {
			interval = meta.Fields["interval"]
			lastRun = meta.Fields["last_run"]
			status = meta.Fields["status"]
		}
		m.checks = append(m.checks, heartbeatCheck{
			Name: entry.Name, Body: body,
			Interval: interval, LastRun: lastRun, Status: status,
		})
	}
}

func (m *heartbeatScreen) Init() tea.Cmd { return nil }

func (m *heartbeatScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "r", "R":
			m.loadChecks()
			return m, nil
		}
	}
	return m, nil
}

func (m *heartbeatScreen) View() string {
	var lines []string
	lines = append(lines, "")

	if len(m.checks) == 0 {
		lines = append(lines, "  No heartbeat checks configured.")
		lines = append(lines, "")
		lines = append(lines, "  Add check files to chores/bots/ to get started.")
	} else {
		lines = append(lines, fmt.Sprintf("  %d checks configured", len(m.checks)))
		lines = append(lines, "")
		for _, c := range m.checks {
			statusIcon := "○"
			if c.Status == "ok" {
				statusIcon = lipgloss.NewStyle().Foreground(lipgloss.Color("#28a745")).Render("●")
			} else if c.Status == "error" {
				statusIcon = lipgloss.NewStyle().Foreground(lipgloss.Color("#dc3545")).Render("✗")
			}
			line := fmt.Sprintf("  %s %s", statusIcon, c.Name)
			if c.Interval != "" {
				line += fmt.Sprintf(" (every %s)", c.Interval)
			}
			if c.LastRun != "" {
				line += lipgloss.NewStyle().Faint(true).Render(fmt.Sprintf(" — last: %s", c.LastRun))
			}
			lines = append(lines, line)
		}
	}

	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("  Last refreshed: %s", time.Now().Format("15:04:05")))

	return strings.Join(lines, "\n")
}

func (m *heartbeatScreen) Name() string                       { return "Heartbeat" }
func (m *heartbeatScreen) Layout() (string, int)              { return "full", 0 }
func (m *heartbeatScreen) Shortcuts() []engine.ShortcutHint {
	return []engine.ShortcutHint{{Key: "R", Label: "Refresh"}}
}
