package sidekicks

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

// RouteTable maps paths to skill handlers.
type RouteTable struct {
	routes map[string]routeEntry
}

type routeEntry struct {
	Skill   string
	Handler string
	Secret  string
}

// NewRouteTable creates a new route table from webhook configs.
func NewRouteTable(configs map[string]WebhookConfig) *RouteTable {
	rt := &RouteTable{
		routes: make(map[string]routeEntry),
	}

	for skill, cfg := range configs {
		for _, route := range cfg.Routes {
			path := route.Path
			if !strings.HasPrefix(path, "/") {
				path = "/" + path
			}
			rt.routes[path] = routeEntry{
				Skill:   skill,
				Handler: route.Handler,
				Secret:  cfg.Secret,
			}
		}
	}

	return rt
}

// Match finds the route entry for a given path.
func (rt *RouteTable) Match(path string) (string, string, bool) {
	entry, ok := rt.routes[path]
	if !ok {
		return "", "", false
	}
	return entry.Skill, entry.Handler, true
}

// VerifyHMAC verifies an HMAC signature.
func VerifyHMAC(message, signature, secret string) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(message))
	expected := hex.EncodeToString(mac.Sum(nil))

	// Compare with or without "sha256=" prefix
	sig := strings.TrimPrefix(signature, "sha256=")
	return hmac.Equal([]byte(expected), []byte(sig))
}

// GetSecret returns the secret for a given path.
func (rt *RouteTable) GetSecret(path string) string {
	entry, ok := rt.routes[path]
	if !ok {
		return ""
	}
	return entry.Secret
}

// ListRoutes returns all registered route paths.
func (rt *RouteTable) ListRoutes() []string {
	var paths []string
	for p := range rt.routes {
		paths = append(paths, p)
	}
	return paths
}

// FormatRouteInfo returns a debug string for route info.
func (rt *RouteTable) FormatRouteInfo() string {
	var sb strings.Builder
	for path, entry := range rt.routes {
		sb.WriteString(fmt.Sprintf("  %s → %s/%s\n", path, entry.Skill, entry.Handler))
	}
	return sb.String()
}
