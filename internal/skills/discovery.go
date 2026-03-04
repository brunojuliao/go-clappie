package skills

import (
	"os"
	"path/filepath"

	"github.com/brunojuliao/go-clappie/internal/platform"
)

// Discover scans for skills in .claude/skills/*/
func Discover(root string) ([]SkillInfo, error) {
	skillsDir := platform.SkillsDir(root)
	entries, err := os.ReadDir(skillsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var skillsList []SkillInfo
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}

		skillPath := filepath.Join(skillsDir, e.Name())
		info := SkillInfo{
			Name: e.Name(),
			Path: skillPath,
		}

		// Check for webhook.json
		webhookPath := filepath.Join(skillPath, "webhook.json")
		if _, err := os.Stat(webhookPath); err == nil {
			info.HasWebhook = true
		}

		skillsList = append(skillsList, info)
	}

	return skillsList, nil
}
