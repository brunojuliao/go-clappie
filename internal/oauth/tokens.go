package oauth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/brunojuliao/go-clappie/internal/platform"
)

func tokensDir(root string) string {
	return filepath.Join(root, "recall", "settings", "oauth")
}

func tokenPath(root, provider string) string {
	return filepath.Join(tokensDir(root), provider+".json")
}

func loadToken(root, provider string) (*TokenData, error) {
	path := tokenPath(root, provider)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var token TokenData
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, err
	}
	return &token, nil
}

func saveToken(root, provider string, token *TokenData) error {
	dir := tokensDir(root)
	if err := platform.EnsureDir(dir); err != nil {
		return err
	}
	data, err := json.MarshalIndent(token, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(tokenPath(root, provider), data, 0600)
}

func isExpired(token *TokenData) bool {
	if token.ExpiresAt.IsZero() {
		return false
	}
	return time.Now().After(token.ExpiresAt)
}

// GetToken returns the access token for a provider, refreshing if expired.
func GetToken(root, provider string) (string, error) {
	token, err := loadToken(root, provider)
	if err != nil {
		return "", fmt.Errorf("no token for %s: %w", provider, err)
	}

	if isExpired(token) && token.RefreshToken != "" {
		if err := Refresh(root, provider); err != nil {
			return "", fmt.Errorf("refresh failed: %w", err)
		}
		token, err = loadToken(root, provider)
		if err != nil {
			return "", err
		}
	}

	return token.AccessToken, nil
}

// Revoke removes the stored token for a provider.
func Revoke(root, provider string) error {
	path := tokenPath(root, provider)
	err := os.Remove(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	fmt.Printf("Token revoked for %s\n", provider)
	return nil
}

// Status prints the status of all OAuth providers.
func Status(root string) error {
	providers, err := ListProviders(root)
	if err != nil {
		return err
	}

	if len(providers) == 0 {
		fmt.Println("No OAuth providers configured.")
		return nil
	}

	for _, p := range providers {
		status := "not authenticated"
		if p.HasToken {
			status = "authenticated"
			if p.Expired {
				status = "expired"
			}
		}
		fmt.Printf("  %s: %s\n", p.Name, status)
	}
	return nil
}
