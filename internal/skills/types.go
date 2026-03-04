package skills

// SkillInfo holds information about a discovered skill.
type SkillInfo struct {
	Name      string
	Path      string
	HasWebhook bool
}

// WebhookConfig holds webhook configuration for a skill.
type WebhookConfig struct {
	Port   int    `json:"port"`
	Secret string `json:"secret"`
}
