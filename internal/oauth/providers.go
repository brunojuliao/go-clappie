package oauth

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/brunojuliao/go-clappie/internal/platform"
)

// ListProviders discovers all OAuth providers from skill directories.
func ListProviders(root string) ([]ProviderInfo, error) {
	skillsDir := platform.SkillsDir(root)
	entries, err := os.ReadDir(skillsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var providers []ProviderInfo
	seen := make(map[string]bool)

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		oauthPath := filepath.Join(skillsDir, e.Name(), "oauth.json")
		data, err := os.ReadFile(oauthPath)
		if err != nil {
			continue
		}

		var config OAuthConfig
		if err := json.Unmarshal(data, &config); err != nil {
			continue
		}

		for _, p := range config.Providers {
			if seen[p.Name] {
				continue
			}
			seen[p.Name] = true

			info := ProviderInfo{
				Name: p.Name,
				Path: oauthPath,
			}

			// Check if token exists
			token, err := loadToken(root, p.Name)
			if err == nil && token != nil {
				info.HasToken = true
				info.Expired = isExpired(token)
			}

			providers = append(providers, info)
		}
	}

	return providers, nil
}

// FindProvider finds a provider by name across all skill directories.
func FindProvider(root, name string) (*Provider, error) {
	skillsDir := platform.SkillsDir(root)
	entries, err := os.ReadDir(skillsDir)
	if err != nil {
		return nil, err
	}

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		oauthPath := filepath.Join(skillsDir, e.Name(), "oauth.json")
		data, err := os.ReadFile(oauthPath)
		if err != nil {
			continue
		}

		var config OAuthConfig
		if err := json.Unmarshal(data, &config); err != nil {
			continue
		}

		for _, p := range config.Providers {
			if p.Name == name {
				return &p, nil
			}
		}
	}

	return nil, nil
}
