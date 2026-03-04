package oauth

import "time"

// ProviderInfo holds information about an OAuth provider.
type ProviderInfo struct {
	Name     string
	HasToken bool
	Expired  bool
	Path     string
}

// Provider holds the configuration for an OAuth provider.
type Provider struct {
	Name         string `json:"name"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	AuthURL      string `json:"auth_url"`
	TokenURL     string `json:"token_url"`
	RedirectURL  string `json:"redirect_url"`
	Scopes       []string `json:"scopes"`
	PKCE         bool   `json:"pkce"`
}

// TokenData holds OAuth token data.
type TokenData struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	ExpiresAt    time.Time `json:"expires_at"`
	Scope        string    `json:"scope"`
}

// OAuthConfig holds the full oauth.json configuration.
type OAuthConfig struct {
	Providers []Provider `json:"providers"`
}
