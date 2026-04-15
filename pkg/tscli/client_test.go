package tscli

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

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
