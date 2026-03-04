package background

import (
	"os"
	"path/filepath"

	"github.com/brunojuliao/go-clappie/internal/platform"
)

// Discover scans for background-capable apps by looking for .background marker files.
func Discover(root string) ([]App, error) {
	skillsDir := platform.SkillsDir(root)
	clappieSkillDir := filepath.Join(skillsDir, "go-clappie", "clapps")

	var apps []App

	// Walk the clapps directory looking for .background files
	entries, err := os.ReadDir(clappieSkillDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		bgMarker := filepath.Join(clappieSkillDir, e.Name(), ".background")
		if _, err := os.Stat(bgMarker); err == nil {
			apps = append(apps, App{
				Name: e.Name(),
				Path: filepath.Join(clappieSkillDir, e.Name()),
			})
		}
	}

	return apps, nil
}
