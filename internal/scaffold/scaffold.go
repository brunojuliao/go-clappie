package scaffold

import (
	"os"
	"path/filepath"
	"strings"
)

// Result holds the outcome of a scaffold operation.
type Result struct {
	Created []string
	Skipped []string
	Errors  []error
}

// Run scaffolds go-clappie into the given project root directory.
// It is idempotent — safe to run multiple times.
func Run(root string) Result {
	var r Result

	// Skill file
	skillPath := filepath.Join(root, ".claude", "skills", "go-clappie", "SKILL.md")
	writeIfMissing(skillPath, SkillMD, &r)

	// CLAUDE.md — skip if exists, append section if exists without go-clappie reference
	claudePath := filepath.Join(root, "CLAUDE.md")
	handleClaudeMD(claudePath, &r)

	// Data directories
	dirs := []string{
		"chores/humans",
		"chores/bots",
		"notifications/dirty",
		"notifications/clean",
		"recall/memory",
		"recall/logs",
		"recall/settings",
		"recall/sidekicks",
		"recall/parties",
		"projects",
	}
	for _, d := range dirs {
		dirPath := filepath.Join(root, d)
		createDir(dirPath, &r)
	}

	return r
}

func writeIfMissing(path string, content []byte, r *Result) {
	if _, err := os.Stat(path); err == nil {
		r.Skipped = append(r.Skipped, path)
		return
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		r.Errors = append(r.Errors, err)
		return
	}

	if err := os.WriteFile(path, content, 0644); err != nil {
		r.Errors = append(r.Errors, err)
		return
	}

	r.Created = append(r.Created, path)
}

func handleClaudeMD(path string, r *Result) {
	existing, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// No CLAUDE.md — create it
			if err := os.WriteFile(path, ClaudeMD, 0644); err != nil {
				r.Errors = append(r.Errors, err)
				return
			}
			r.Created = append(r.Created, path)
			return
		}
		r.Errors = append(r.Errors, err)
		return
	}

	// CLAUDE.md exists — check if it already mentions go-clappie
	if strings.Contains(string(existing), "go-clappie") {
		r.Skipped = append(r.Skipped, path)
		return
	}

	// Append go-clappie section
	appended := string(existing)
	if !strings.HasSuffix(appended, "\n") {
		appended += "\n"
	}
	appended += "\n" + string(ClaudeMD)

	if err := os.WriteFile(path, []byte(appended), 0644); err != nil {
		r.Errors = append(r.Errors, err)
		return
	}
	r.Created = append(r.Created, path+" (appended)")
}

func createDir(path string, r *Result) {
	if info, err := os.Stat(path); err == nil && info.IsDir() {
		r.Skipped = append(r.Skipped, path)
		return
	}

	if err := os.MkdirAll(path, 0755); err != nil {
		r.Errors = append(r.Errors, err)
		return
	}

	r.Created = append(r.Created, path)
}
