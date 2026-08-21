package tscli

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	tsapi "tailscale.com/client/tailscale/v2"
)

func TestNewOAuthClientRejectsUnsafeBaseURL(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)
	viper.Set("tailnet", "example.com")
	viper.Set("oauth-client-id", "client-id")
	viper.Set("oauth-client-secret", "client-secret")
	viper.Set("base-url", "http://attacker.invalid")

	if _, err := New(); err == nil {
		t.Fatal("expected New() to reject a non-loopback http:// base-url")
	}
}

func TestNewAPIKeyClientRejectsUnsafeBaseURL(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)
	viper.Set("tailnet", "example.com")
	viper.Set("api-key", "tskey-test")
	viper.Set("base-url", "http://attacker.invalid")

	if _, err := New(); err == nil {
		t.Fatal("expected New() to reject a non-loopback http:// base-url on the api-key auth path")
	}
}

func TestNewOAuthClientRejectsUnsafeOAuthTokenURLOverride(t *testing.T) {
	t.Setenv("TSCLI_OAUTH_TOKEN_URL", "http://attacker.invalid/collect")
	viper.Reset()
	t.Cleanup(viper.Reset)
	viper.Set("tailnet", "example.com")
	viper.Set("oauth-client-id", "client-id")
	viper.Set("oauth-client-secret", "client-secret")
	viper.Set("base-url", "https://api.tailscale.com")

	if _, err := New(); err == nil {
		t.Fatal("expected New() to reject a non-loopback http:// TSCLI_OAUTH_TOKEN_URL")
	}
}

func TestNewOAuthClientHonorsValidatedOAuthTokenURLOverride(t *testing.T) {
	var tokenRequests atomic.Int32
	override := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		tokenRequests.Add(1)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"access_token":"overridden","token_type":"Bearer","expires_in":3600}`))
	}))
	defer override.Close()

	configured := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer overridden" {
			t.Errorf("expected overridden bearer token, got %q", got)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer configured.Close()

	t.Setenv("TSCLI_OAUTH_TOKEN_URL", override.URL+"/oauth/token")
	viper.Reset()
	t.Cleanup(viper.Reset)
	viper.Set("tailnet", "example.com")
	viper.Set("oauth-client-id", "client-id")
	viper.Set("oauth-client-secret", "client-secret")
	viper.Set("base-url", configured.URL)

	client, err := New()
	if err != nil {
		t.Fatalf("new client: %v", err)
	}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, configured.URL+"/resource", nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	resp, err := client.HTTP.Do(req)
	if err != nil {
		t.Fatalf("request with oauth client: %v", err)
	}
	_ = resp.Body.Close()

	if got := tokenRequests.Load(); got != 1 {
		t.Fatalf("expected one token request to the overridden endpoint, got %d", got)
	}
}

func TestOauthTokenURLPreservesBaseURLPathPrefix(t *testing.T) {
	tests := []struct {
		name    string
		baseURL string
		want    string
	}{
		{name: "bare host", baseURL: "https://api.tailscale.com", want: "https://api.tailscale.com/api/v2/oauth/token"},
		{name: "trailing slash", baseURL: "https://api.tailscale.com/", want: "https://api.tailscale.com/api/v2/oauth/token"},
		{name: "path prefix", baseURL: "https://proxy.example.com/tailscale", want: "https://proxy.example.com/tailscale/api/v2/oauth/token"},
		{name: "path prefix trailing slash", baseURL: "https://proxy.example.com/tailscale/", want: "https://proxy.example.com/tailscale/api/v2/oauth/token"},
		{name: "versioned API root", baseURL: "https://proxy.example.com/api/v2", want: "https://proxy.example.com/api/v2/oauth/token"},
		{name: "versioned API root trailing slash", baseURL: "https://proxy.example.com/api/v2/", want: "https://proxy.example.com/api/v2/oauth/token"},
		{name: "prefixed versioned API root", baseURL: "https://proxy.example.com/tailscale/api/v2", want: "https://proxy.example.com/tailscale/api/v2/oauth/token"},
		{name: "host with port", baseURL: "http://127.0.0.1:8080", want: "http://127.0.0.1:8080/api/v2/oauth/token"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := url.Parse(tt.baseURL)
			if err != nil {
				t.Fatalf("parse base url: %v", err)
			}
			got, err := oauthTokenURL(u)
			if err != nil {
				t.Fatalf("oauthTokenURL: %v", err)
			}
			if got != tt.want {
				t.Fatalf("oauthTokenURL(%q) = %q, want %q", tt.baseURL, got, tt.want)
			}
		})
	}
}

func TestOAuthBearerTransportRefreshesExpiredTokenAndRetries401(t *testing.T) {
	var (
		mu           sync.Mutex
		tokenCalls   int
		resourceHits int
	)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oauth/token":
			mu.Lock()
			tokenCalls++
			call := tokenCalls
			mu.Unlock()
			w.Header().Set("Content-Type", "application/json")
			_, _ = fmt.Fprintf(w, `{"access_token":"tok-%d","token_type":"Bearer","expires_in":3600}`, call)
		case "/resource":
			mu.Lock()
			resourceHits++
			hit := resourceHits
			mu.Unlock()
			if hit == 1 {
				if got := r.Header.Get("Authorization"); got != "Bearer tok-1" {
					t.Fatalf("expected first request to use tok-1, got %q", got)
				}
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			if got := r.Header.Get("Authorization"); got != "Bearer tok-2" {
				t.Fatalf("expected retried request to use tok-2, got %q", got)
			}
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	transport := &oauthBearerTransport{
		rt:           server.Client().Transport,
		clientID:     "cid",
		clientSecret: "secret",
		tokenURL:     server.URL + "/oauth/token",
		expiresAt:    time.Now().Add(-time.Minute),
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL+"/resource", nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}

	resp, err := transport.RoundTrip(req)
	if err != nil {
		t.Fatalf("round trip: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected final response 200, got %d", resp.StatusCode)
	}

	mu.Lock()
	defer mu.Unlock()
	if tokenCalls != 2 {
		t.Fatalf("expected two token exchanges, got %d", tokenCalls)
	}
	if resourceHits != 2 {
		t.Fatalf("expected two resource requests, got %d", resourceHits)
	}
}

func TestExchangeOAuthClientCredentialsRejectsUnsafeBaseURL(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)
	viper.Set("base-url", "http://attacker.invalid")

	var requested bool
	httpClient := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		requested = true
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader("{}")), Request: req}, nil
	})}
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, httpClient)

	_, err := ExchangeOAuthClientCredentials(ctx, "client-id", "client-secret")
	if err == nil {
		t.Fatal("expected ExchangeOAuthClientCredentials to reject a non-loopback http:// base-url")
	}
	if !strings.Contains(err.Error(), "invalid base-url") {
		t.Fatalf("expected an invalid base-url validation error, got %v", err)
	}
	if requested {
		t.Fatal("expected no request to be sent to the rejected base-url")
	}
}

func TestExchangeOAuthClientCredentialsSucceedsAgainstConfiguredBaseURL(t *testing.T) {
	var tokenRequests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v2/oauth/token" {
			http.NotFound(w, r)
			return
		}
		tokenRequests.Add(1)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"access_token":"tok","token_type":"Bearer","expires_in":3600}`))
	}))
	defer server.Close()

	viper.Reset()
	t.Cleanup(viper.Reset)
	viper.Set("base-url", server.URL)

	resp, err := ExchangeOAuthClientCredentials(context.Background(), "client-id", "client-secret")
	if err != nil {
		t.Fatalf("exchange oauth client credentials: %v", err)
	}
	if resp.AccessToken != "tok" {
		t.Fatalf("expected access token %q, got %q", "tok", resp.AccessToken)
	}
	if got := tokenRequests.Load(); got != 1 {
		t.Fatalf("expected one token request, got %d", got)
	}
}

func TestExchangeOAuthClientCredentialsAcceptsVersionedBaseURL(t *testing.T) {
	var gotPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"access_token":"tok","token_type":"Bearer","expires_in":3600}`))
	}))
	defer server.Close()

	viper.Reset()
	t.Cleanup(viper.Reset)
	viper.Set("base-url", server.URL+"/api/v2")

	if _, err := ExchangeOAuthClientCredentials(context.Background(), "client-id", "client-secret"); err != nil {
		t.Fatalf("exchange oauth client credentials: %v", err)
	}
	if want := "/api/v2/oauth/token"; gotPath != want {
		t.Fatalf("expected token request path %q, got %q", want, gotPath)
	}
}

func TestNewOAuthClientDoesNotFollowRedirectsAcrossHosts(t *testing.T) {
	var evilRequests atomic.Int32
	evil := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		evilRequests.Add(1)
		if got := r.Header.Get("Authorization"); got != "" {
			t.Errorf("expected no Authorization header sent to redirect target, got %q", got)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer evil.Close()

	configured := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v2/oauth/token":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"access_token":"secret-token","token_type":"Bearer","expires_in":3600}`))
		case "/resource":
			http.Redirect(w, r, evil.URL+"/resource", http.StatusFound)
		default:
			http.NotFound(w, r)
		}
	}))
	defer configured.Close()

	viper.Reset()
	t.Cleanup(viper.Reset)
	viper.Set("tailnet", "example.com")
	viper.Set("oauth-client-id", "client-id")
	viper.Set("oauth-client-secret", "client-secret")
	viper.Set("base-url", configured.URL)

	client, err := New()
	if err != nil {
		t.Fatalf("new client: %v", err)
	}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, configured.URL+"/resource", nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	resp, err := client.HTTP.Do(req)
	if err != nil {
		t.Fatalf("request with oauth client: %v", err)
	}
	_ = resp.Body.Close()

	if resp.StatusCode != http.StatusFound {
		t.Fatalf("expected the client to NOT follow the redirect (status 302 returned as-is), got %d", resp.StatusCode)
	}
	if got := evilRequests.Load(); got != 0 {
		t.Fatalf("expected no request to the redirect target, got %d", got)
	}
}

func TestNewUsesConfiguredUserAgentOverride(t *testing.T) {
	var gotUserAgent string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUserAgent = r.Header.Get("User-Agent")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	viper.Reset()
	t.Cleanup(viper.Reset)
	viper.Set("tailnet", "example.com")
	viper.Set("api-key", "tskey-test")
	viper.Set("base-url", server.URL)
	viper.Set("user-agent", "my-app/1.0")

	client, err := New()
	if err != nil {
		t.Fatalf("new client: %v", err)
	}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL+"/resource", nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	resp, err := client.HTTP.Do(req)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	_ = resp.Body.Close()

	if gotUserAgent != "my-app/1.0" {
		t.Fatalf("expected overridden User-Agent %q, got %q", "my-app/1.0", gotUserAgent)
	}
}

func TestNewDefaultsUserAgentWhenNotConfigured(t *testing.T) {
	var gotUserAgent string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUserAgent = r.Header.Get("User-Agent")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	viper.Reset()
	t.Cleanup(viper.Reset)
	viper.Set("tailnet", "example.com")
	viper.Set("api-key", "tskey-test")
	viper.Set("base-url", server.URL)

	client, err := New()
	if err != nil {
		t.Fatalf("new client: %v", err)
	}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL+"/resource", nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	resp, err := client.HTTP.Do(req)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	_ = resp.Body.Close()

	if gotUserAgent == "" || gotUserAgent == "my-app/1.0" {
		t.Fatalf("expected a non-empty default User-Agent, got %q", gotUserAgent)
	}
}

func TestDoPreservesBaseURLPathPrefix(t *testing.T) {
	var gotPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()

	base, err := url.Parse(server.URL + "/proxy-prefix")
	if err != nil {
		t.Fatalf("parse base url: %v", err)
	}
	client := &tsapi.Client{
		Tailnet:   "example.com",
		UserAgent: "test",
		BaseURL:   base,
		HTTP:      server.Client(),
		APIKey:    "tskey-test",
	}

	if _, err := Do(context.Background(), client, http.MethodGet, "/tailnet/{tailnet}/devices", nil, nil); err != nil {
		t.Fatalf("Do: %v", err)
	}

	want := "/proxy-prefix/api/v2/tailnet/example.com/devices"
	if gotPath != want {
		t.Fatalf("expected request path %q (preserving base-url prefix), got %q", want, gotPath)
	}
}

func TestDoAcceptsVersionedBaseURL(t *testing.T) {
	var gotPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	base, err := url.Parse(server.URL + "/api/v2")
	if err != nil {
		t.Fatalf("parse base url: %v", err)
	}
	client := &tsapi.Client{
		Tailnet:   "example.com",
		UserAgent: "test",
		BaseURL:   base,
		HTTP:      server.Client(),
		APIKey:    "tskey-test",
	}

	if _, err := Do(context.Background(), client, http.MethodGet, "/tailnet/{tailnet}/devices", nil, nil); err != nil {
		t.Fatalf("Do: %v", err)
	}
	if want := "/api/v2/tailnet/example.com/devices"; gotPath != want {
		t.Fatalf("expected request path %q, got %q", want, gotPath)
	}
}

func TestDoBearerHonorsUserAgentEnvironmentOverride(t *testing.T) {
	var gotUserAgent string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUserAgent = r.Header.Get("User-Agent")
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	t.Setenv("TSCLI_USER_AGENT", "my-bearer-app/1.0")
	viper.Reset()
	t.Cleanup(viper.Reset)
	if err := viper.BindEnv("user-agent", "TSCLI_USER_AGENT"); err != nil {
		t.Fatalf("bind user-agent environment variable: %v", err)
	}
	viper.Set("base-url", server.URL)

	if _, err := DoBearer(context.Background(), http.MethodGet, "/organizations/-/tailnets", "tok-123", nil, nil); err != nil {
		t.Fatalf("DoBearer: %v", err)
	}
	if want := "my-bearer-app/1.0"; gotUserAgent != want {
		t.Fatalf("expected overridden User-Agent %q, got %q", want, gotUserAgent)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
