package sidekicks

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/brunojuliao/go-clappie/internal/platform"
)

// WebhookConfig holds webhook configuration for a skill.
type WebhookConfig struct {
	Port     int               `json:"port"`
	Routes   []WebhookRoute    `json:"routes"`
	Secret   string            `json:"secret"`
	Settings map[string]string `json:"settings"`
}

// WebhookRoute defines a webhook route.
type WebhookRoute struct {
	Path    string `json:"path"`
	Method  string `json:"method"`
	Handler string `json:"handler"`
}

// DiscoverWebhooks finds all webhook configurations from skills.
func DiscoverWebhooks(root string) (map[string]WebhookConfig, error) {
	skillsDir := platform.SkillsDir(root)
	entries, err := os.ReadDir(skillsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	configs := make(map[string]WebhookConfig)
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		webhookPath := filepath.Join(skillsDir, e.Name(), "webhook.json")
		data, err := os.ReadFile(webhookPath)
		if err != nil {
			continue
		}
		var cfg WebhookConfig
		if err := json.Unmarshal(data, &cfg); err != nil {
			continue
		}
		configs[e.Name()] = cfg
	}

	return configs, nil
}
