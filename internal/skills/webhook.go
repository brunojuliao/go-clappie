package skills

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// LoadWebhookConfig loads the webhook.json from a skill directory.
func LoadWebhookConfig(skillPath string) (*WebhookConfig, error) {
	path := filepath.Join(skillPath, "webhook.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config WebhookConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
