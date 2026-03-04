package oauth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Auth starts the OAuth authorization flow for a provider.
func Auth(root, providerName string) error {
	provider, err := FindProvider(root, providerName)
	if err != nil || provider == nil {
		return fmt.Errorf("provider %q not found", providerName)
	}

	// Generate PKCE challenge if needed
	var codeVerifier, codeChallenge string
	if provider.PKCE {
		verifier := make([]byte, 32)
		rand.Read(verifier)
		codeVerifier = base64.RawURLEncoding.EncodeToString(verifier)

		h := sha256.Sum256([]byte(codeVerifier))
		codeChallenge = base64.RawURLEncoding.EncodeToString(h[:])
	}

	// Generate state parameter
	stateBytes := make([]byte, 16)
	rand.Read(stateBytes)
	state := base64.RawURLEncoding.EncodeToString(stateBytes)

	// Build authorization URL
	params := url.Values{
		"client_id":     {provider.ClientID},
		"redirect_uri":  {provider.RedirectURL},
		"response_type": {"code"},
		"scope":         {strings.Join(provider.Scopes, " ")},
		"state":         {state},
	}
	if provider.PKCE {
		params.Set("code_challenge", codeChallenge)
		params.Set("code_challenge_method", "S256")
	}

	authURL := provider.AuthURL + "?" + params.Encode()
	fmt.Printf("Open this URL to authenticate:\n%s\n\n", authURL)

	// Start callback server
	return startCallbackServer(root, provider, state, codeVerifier)
}

func startCallbackServer(root string, provider *Provider, expectedState, codeVerifier string) error {
	codeCh := make(chan string, 1)
	errCh := make(chan error, 1)

	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		state := r.URL.Query().Get("state")
		if state != expectedState {
			http.Error(w, "invalid state", http.StatusBadRequest)
			errCh <- fmt.Errorf("state mismatch")
			return
		}

		code := r.URL.Query().Get("code")
		if code == "" {
			errorDesc := r.URL.Query().Get("error_description")
			http.Error(w, "no code", http.StatusBadRequest)
			errCh <- fmt.Errorf("no code: %s", errorDesc)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, "<html><body><h1>Authentication successful!</h1><p>You can close this window.</p></body></html>")
		codeCh <- code
	})

	server := &http.Server{
		Addr:    ":8085",
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(ctx)
	}()

	fmt.Println("Waiting for callback on :8085...")

	select {
	case code := <-codeCh:
		return exchangeCode(root, provider, code, codeVerifier)
	case err := <-errCh:
		return err
	case <-time.After(5 * time.Minute):
		return fmt.Errorf("auth timeout after 5 minutes")
	}
}

func exchangeCode(root string, provider *Provider, code, codeVerifier string) error {
	data := url.Values{
		"grant_type":   {"authorization_code"},
		"code":         {code},
		"redirect_uri": {provider.RedirectURL},
		"client_id":    {provider.ClientID},
	}
	if provider.ClientSecret != "" {
		data.Set("client_secret", provider.ClientSecret)
	}
	if codeVerifier != "" {
		data.Set("code_verifier", codeVerifier)
	}

	resp, err := http.PostForm(provider.TokenURL, data)
	if err != nil {
		return fmt.Errorf("token exchange: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("token exchange failed (%d): %s", resp.StatusCode, string(body))
	}

	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
		Scope        string `json:"scope"`
	}
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return fmt.Errorf("parse token: %w", err)
	}

	token := &TokenData{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		TokenType:    tokenResp.TokenType,
		Scope:        tokenResp.Scope,
	}
	if tokenResp.ExpiresIn > 0 {
		token.ExpiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	}

	if err := saveToken(root, provider.Name, token); err != nil {
		return fmt.Errorf("save token: %w", err)
	}

	fmt.Printf("Authentication successful for %s\n", provider.Name)
	return nil
}

// Refresh refreshes an OAuth token.
func Refresh(root, providerName string) error {
	provider, err := FindProvider(root, providerName)
	if err != nil || provider == nil {
		return fmt.Errorf("provider %q not found", providerName)
	}

	token, err := loadToken(root, providerName)
	if err != nil {
		return fmt.Errorf("no token to refresh: %w", err)
	}

	if token.RefreshToken == "" {
		return fmt.Errorf("no refresh token available")
	}

	data := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {token.RefreshToken},
		"client_id":     {provider.ClientID},
	}
	if provider.ClientSecret != "" {
		data.Set("client_secret", provider.ClientSecret)
	}

	resp, err := http.PostForm(provider.TokenURL, data)
	if err != nil {
		return fmt.Errorf("refresh request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("refresh failed (%d): %s", resp.StatusCode, string(body))
	}

	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
	}
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return err
	}

	newToken := &TokenData{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		TokenType:    tokenResp.TokenType,
	}
	if newToken.RefreshToken == "" {
		newToken.RefreshToken = token.RefreshToken
	}
	if tokenResp.ExpiresIn > 0 {
		newToken.ExpiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	}

	return saveToken(root, providerName, newToken)
}
