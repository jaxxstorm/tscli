package oauth

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"golang.org/x/oauth2"
)

func TestExchangeClientCredentialsDefaultsToPublicTokenEndpoint(t *testing.T) {
	var requestedURL string
	httpClient := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		requestedURL = req.URL.String()
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(`{"access_token":"token","token_type":"Bearer","expires_in":3600}`)),
			Request:    req,
		}, nil
	})}
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, httpClient)

	if _, err := ExchangeClientCredentials(ctx, "client-id", "client-secret", ""); err != nil {
		t.Fatalf("exchange client credentials: %v", err)
	}
	if requestedURL != publicTokenURL {
		t.Fatalf("expected token request to %q, got %q", publicTokenURL, requestedURL)
	}
}

func TestExchangeClientCredentialsHonorsValidatedTokenURLOverride(t *testing.T) {
	t.Setenv("TSCLI_OAUTH_TOKEN_URL", "http://127.0.0.1:1/oauth/token")

	var requestedURL string
	httpClient := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		requestedURL = req.URL.String()
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(`{"access_token":"token","token_type":"Bearer","expires_in":3600}`)),
			Request:    req,
		}, nil
	})}
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, httpClient)

	if _, err := ExchangeClientCredentials(ctx, "client-id", "client-secret", ""); err != nil {
		t.Fatalf("exchange client credentials: %v", err)
	}
	if requestedURL != "http://127.0.0.1:1/oauth/token" {
		t.Fatalf("expected token request to overridden loopback URL, got %q", requestedURL)
	}
}

func TestExchangeClientCredentialsRejectsUnsafeTokenURLOverride(t *testing.T) {
	t.Setenv("TSCLI_OAUTH_TOKEN_URL", "http://attacker.invalid/collect")

	var requested bool
	httpClient := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		requested = true
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader("{}")), Request: req}, nil
	})}
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, httpClient)

	if _, err := ExchangeClientCredentials(ctx, "client-id", "client-secret", ""); err == nil {
		t.Fatal("expected an error for an unvalidated http:// TSCLI_OAUTH_TOKEN_URL")
	}
	if requested {
		t.Fatal("expected no request to be sent to the rejected token URL")
	}
}

func TestValidateTokenURL(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		wantErr bool
	}{
		{name: "https", raw: "https://api.tailscale.com/api/v2/oauth/token"},
		{name: "http loopback IPv4", raw: "http://127.0.0.1:8080/oauth/token"},
		{name: "http loopback IPv6", raw: "http://[::1]:8080/oauth/token"},
		{name: "http localhost", raw: "http://localhost:8080/oauth/token"},
		{name: "http non-loopback rejected", raw: "http://attacker.invalid/collect", wantErr: true},
		{name: "unsupported scheme rejected", raw: "ftp://api.tailscale.com/token", wantErr: true},
		{name: "relative URL rejected", raw: "/oauth/token", wantErr: true},
		{name: "malformed URL rejected", raw: "://bad-url", wantErr: true},
		{name: "empty host rejected", raw: "https:///oauth/token", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ValidateTokenURL(tt.raw)
			if tt.wantErr && err == nil {
				t.Fatalf("expected an error for %q, got none", tt.raw)
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("expected no error for %q, got %v", tt.raw, err)
			}
		})
	}
}

func TestResolveTokenURL(t *testing.T) {
	t.Run("no override returns fallback", func(t *testing.T) {
		if got, err := ResolveTokenURL("https://fallback.example/token"); err != nil || got != "https://fallback.example/token" {
			t.Fatalf("got (%q, %v), want (%q, nil)", got, err, "https://fallback.example/token")
		}
	})
	t.Run("valid override replaces fallback", func(t *testing.T) {
		t.Setenv("TSCLI_OAUTH_TOKEN_URL", "https://override.example/token")
		if got, err := ResolveTokenURL("https://fallback.example/token"); err != nil || got != "https://override.example/token" {
			t.Fatalf("got (%q, %v), want (%q, nil)", got, err, "https://override.example/token")
		}
	})
	t.Run("unsafe override is rejected", func(t *testing.T) {
		t.Setenv("TSCLI_OAUTH_TOKEN_URL", "http://attacker.invalid/token")
		if _, err := ResolveTokenURL("https://fallback.example/token"); err == nil {
			t.Fatal("expected an error for an unsafe override")
		}
	})
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
