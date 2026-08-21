package oauth

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"os"
	"time"

	"golang.org/x/oauth2/clientcredentials"
)

const publicTokenURL = "https://api.tailscale.com/api/v2/oauth/token"

// TokenResponse represents the response from the OAuth token exchange
type TokenResponse struct {
	AccessToken string    `json:"access_token"`
	TokenType   string    `json:"token_type"`
	ExpiresIn   int       `json:"expires_in"`
	ExpiresAt   time.Time `json:"expires_at"`
}

// ExchangeClientCredentials exchanges OAuth client credentials for an
// access token at defaultTokenURL — pass "" to use Tailscale's public
// token endpoint (https://api.tailscale.com) — unless overridden by the
// TSCLI_OAUTH_TOKEN_URL environment variable. Because the client secret is
// POSTed to whatever URL is used, both defaultTokenURL and any
// TSCLI_OAUTH_TOKEN_URL override are validated by ValidateTokenURL: each
// must be an absolute https:// URL, or an absolute http:// URL restricted
// to a loopback host (for local testing).
func ExchangeClientCredentials(ctx context.Context, clientID, clientSecret, defaultTokenURL string) (*TokenResponse, error) {
	if defaultTokenURL == "" {
		defaultTokenURL = publicTokenURL
	}
	if _, err := ValidateTokenURL(defaultTokenURL); err != nil {
		return nil, fmt.Errorf("defaultTokenURL: %w", err)
	}
	tokenURL, err := ResolveTokenURL(defaultTokenURL)
	if err != nil {
		return nil, err
	}
	return ExchangeClientCredentialsAtURL(ctx, clientID, clientSecret, tokenURL)
}

// ExchangeClientCredentialsAtURL exchanges OAuth client credentials at a
// caller-selected token endpoint. The caller is responsible for trusting the
// endpoint because the client credentials are sent to it.
func ExchangeClientCredentialsAtURL(ctx context.Context, clientID, clientSecret, tokenURL string) (*TokenResponse, error) {
	config := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     tokenURL,
	}

	token, err := config.Token(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange credentials: %w", err)
	}

	// Convert oauth2.Token to our TokenResponse format
	expiresIn := 0
	expiresAt := time.Time{}
	if !token.Expiry.IsZero() {
		expiresIn = int(time.Until(token.Expiry).Seconds())
		// Ensure we don't return negative values if token is already expired
		if expiresIn < 0 {
			expiresIn = 0
		}
		expiresAt = token.Expiry
	}

	return &TokenResponse{
		AccessToken: token.AccessToken,
		TokenType:   token.TokenType,
		ExpiresIn:   expiresIn,
		ExpiresAt:   expiresAt,
	}, nil
}

// ValidateTokenURL parses raw and requires it to be an absolute https://
// URL, or an absolute http:// URL restricted to a loopback host
// (127.0.0.0/8, ::1, or "localhost"). It is used to validate any
// caller- or environment-configured OAuth token/base URL, since OAuth
// client credentials (and, via pkg/tscli's base-url, API keys) are sent to
// that URL — an unrestricted http:// endpoint would send them in
// cleartext to whatever host is configured.
func ValidateTokenURL(raw string) (*url.URL, error) {
	u, err := url.Parse(raw)
	if err != nil || !u.IsAbs() || u.Host == "" || u.Opaque != "" {
		return nil, fmt.Errorf("must be an absolute URL with scheme and host: %q", raw)
	}
	switch u.Scheme {
	case "https":
		return u, nil
	case "http":
		if isLoopbackHost(u.Hostname()) {
			return u, nil
		}
		return nil, fmt.Errorf("http:// is only allowed for loopback hosts, got %q", u.Host)
	default:
		return nil, fmt.Errorf("scheme %q is not supported, use https://", u.Scheme)
	}
}

// ResolveTokenURL returns the TSCLI_OAUTH_TOKEN_URL environment variable,
// validated via ValidateTokenURL, if it is set and non-empty; otherwise it
// returns fallback unchanged. Centralizing this here means the override
// behavior (and its validation) is defined in exactly one place, reused by
// both ExchangeClientCredentials and pkg/tscli's OAuth token-URL
// derivation.
func ResolveTokenURL(fallback string) (string, error) {
	override := os.Getenv("TSCLI_OAUTH_TOKEN_URL")
	if override == "" {
		return fallback, nil
	}
	u, err := ValidateTokenURL(override)
	if err != nil {
		return "", fmt.Errorf("TSCLI_OAUTH_TOKEN_URL: %w", err)
	}
	return u.String(), nil
}

func isLoopbackHost(host string) bool {
	if host == "localhost" {
		return true
	}
	if ip := net.ParseIP(host); ip != nil {
		return ip.IsLoopback()
	}
	return false
}
