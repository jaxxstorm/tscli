package tscli

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	tsapi "tailscale.com/client/tailscale/v2"
)

func TestDoDecodesResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v2/tailnet/example.com/settings" {
			http.Error(w, `{"message":"unexpected path"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"devicesApprovalOn":"off"}`))
	}))
	defer srv.Close()

	base, _ := url.Parse(srv.URL)
	client := &tsapi.Client{
		Tailnet: "example.com",
		APIKey:  "tskey-test",
		BaseURL: base,
		HTTP:    srv.Client(),
	}

	var out map[string]any
	_, err := Do(context.Background(), client, http.MethodGet, "/tailnet/{tailnet}/settings", nil, &out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["devicesApprovalOn"] != "off" {
		t.Fatalf("unexpected decoded payload: %#v", out)
	}
}

func TestDoReturnsDecodeErrorForInvalidPayload(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"bad-json"`))
	}))
	defer srv.Close()

	base, _ := url.Parse(srv.URL)
	client := &tsapi.Client{
		Tailnet: "example.com",
		APIKey:  "tskey-test",
		BaseURL: base,
		HTTP:    srv.Client(),
	}

	var out map[string]any
	_, err := Do(context.Background(), client, http.MethodGet, "/tailnet/{tailnet}/settings", nil, &out)
	if err == nil {
		t.Fatalf("expected decode error")
	}
	if !strings.Contains(err.Error(), "decode response") {
		t.Fatalf("expected decode response error, got %v", err)
	}
}

func TestDoReturnsAPIErrorPayload(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, `{"message":"server exploded"}`, http.StatusInternalServerError)
	}))
	defer srv.Close()

	base, _ := url.Parse(srv.URL)
	client := &tsapi.Client{
		Tailnet: "example.com",
		APIKey:  "tskey-test",
		BaseURL: base,
		HTTP:    srv.Client(),
	}

	_, err := Do(context.Background(), client, http.MethodGet, "/tailnet/{tailnet}/settings", nil, nil)
	if err == nil {
		t.Fatalf("expected api error")
	}
	if !strings.Contains(err.Error(), "tailscale API") {
		t.Fatalf("expected wrapped API error, got %v", err)
	}
}
