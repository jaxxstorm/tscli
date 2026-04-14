package cli_test

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jaxxstorm/tscli/internal/testutil/apimock"
)

func TestTailnetLifecycleCommands(t *testing.T) {
	mock := apimock.New(t)
	mock.AddRaw(http.MethodPost, "/api/v2/oauth/token", http.StatusOK, `{"access_token":"tok-123","token_type":"Bearer","expires_in":3600}`)
	mock.AddRaw(http.MethodPost, "/api/v2/organizations/-/tailnets", http.StatusOK, `{"id":"T123","displayName":"Sandbox","orgId":"o123","dnsName":"tail123.ts.net","createdAt":"2025-01-01T12:00:00Z","oauthClient":{"id":"k123","secret":"tskey-client-secret"}}`)
	mock.AddRaw(http.MethodGet, "/api/v2/organizations/-/tailnets", http.StatusOK, `{"tailnets":[{"id":"T123","displayName":"Sandbox","orgId":"o123","createdAt":"2025-01-01T12:00:00Z"}]}`)
	mock.AddRaw(http.MethodDelete, "/api/v2/tailnet/T123", http.StatusOK, `{}`)

	env := map[string]string{
		"TSCLI_BASE_URL":        mock.URL(),
		"TSCLI_OAUTH_TOKEN_URL": mock.URL() + "/api/v2/oauth/token",
	}

	res := executeCLINoDefaults(t, []string{"create", "tailnet", "--display-name", "Sandbox", "--oauth-client-id", "cid", "--oauth-client-secret", "secret"}, env)
	if res.err != nil {
		t.Fatalf("create tailnet: %v\nstderr:\n%s", res.err, res.stderr)
	}
	var created map[string]any
	if err := json.Unmarshal([]byte(res.stdout), &created); err != nil {
		t.Fatalf("unmarshal create tailnet output: %v\noutput:\n%s", err, res.stdout)
	}
	if created["displayName"] != "Sandbox" {
		t.Fatalf("expected authoritative create response, got %s", res.stdout)
	}

	res = executeCLINoDefaults(t, []string{"list", "tailnets", "--oauth-client-id", "cid", "--oauth-client-secret", "secret"}, env)
	if res.err != nil {
		t.Fatalf("list tailnets: %v\nstderr:\n%s", res.err, res.stderr)
	}
	var listed map[string]any
	if err := json.Unmarshal([]byte(res.stdout), &listed); err != nil {
		t.Fatalf("unmarshal list tailnets output: %v\noutput:\n%s", err, res.stdout)
	}
	tailnets, _ := listed["tailnets"].([]any)
	if len(tailnets) != 1 {
		t.Fatalf("expected authoritative list response, got %s", res.stdout)
	}

	res = executeCLINoDefaults(t, []string{"delete", "tailnet", "--id", "T123", "--oauth-client-id", "cid", "--oauth-client-secret", "secret"}, env)
	if res.err != nil {
		t.Fatalf("delete tailnet: %v\nstderr:\n%s", res.err, res.stderr)
	}
	var deleted map[string]any
	if err := json.Unmarshal([]byte(res.stdout), &deleted); err != nil {
		t.Fatalf("unmarshal delete tailnet output: %v\noutput:\n%s", err, res.stdout)
	}
	if deleted["result"] != "tailnet deleted" {
		t.Fatalf("expected delete summary response, got %s", res.stdout)
	}

	reqs := mock.Requests()
	if len(reqs) < 6 {
		t.Fatalf("expected token and lifecycle requests, got %+v", reqs)
	}
	if !strings.Contains(reqs[1].Body, `"displayName":"Sandbox"`) {
		t.Fatalf("expected create request body to include displayName, got %s", reqs[1].Body)
	}
	for _, req := range reqs {
		if strings.Contains(req.Path, "/organizations/-/tailnets") || strings.Contains(req.Path, "/tailnet/T123") {
			if got := req.Header.Get("Authorization"); got != "Bearer tok-123" {
				t.Fatalf("expected bearer auth on lifecycle request, got %q for %s", got, req.Path)
			}
		}
	}
}

func TestTailnetLifecycleCommandsUseActiveOAuthProfile(t *testing.T) {
	home := t.TempDir()
	configFile := filepath.Join(home, ".tscli.yaml")
	cfg := strings.Join([]string{
		"output: json",
		"active-tailnet: org-admin",
		"tailnets:",
		"  - name: org-admin",
		"    oauth-client-id: cid-profile",
		"    oauth-client-secret: secret-profile",
		"",
	}, "\n")
	if err := os.WriteFile(configFile, []byte(cfg), 0o600); err != nil {
		t.Fatalf("write config file: %v", err)
	}

	mock := apimock.New(t)
	mock.AddRaw(http.MethodPost, "/api/v2/oauth/token", http.StatusOK, `{"access_token":"tok-profile","token_type":"Bearer","expires_in":3600}`)
	mock.AddRaw(http.MethodGet, "/api/v2/organizations/-/tailnets", http.StatusOK, `{"tailnets":[]}`)

	res := executeCLINoDefaults(t, []string{"list", "tailnets"}, map[string]string{
		"HOME":                  home,
		"TSCLI_BASE_URL":        mock.URL(),
		"TSCLI_OAUTH_TOKEN_URL": mock.URL() + "/api/v2/oauth/token",
	})
	if res.err != nil {
		t.Fatalf("list tailnets with oauth profile: %v\nstderr:\n%s", res.err, res.stderr)
	}

	reqs := mock.Requests()
	if len(reqs) < 2 {
		t.Fatalf("expected token exchange and lifecycle request, got %+v", reqs)
	}
	if got := reqs[0].Header.Get("Authorization"); got != "Basic "+base64.StdEncoding.EncodeToString([]byte("cid-profile:secret-profile")) {
		t.Fatalf("expected basic auth header for oauth token exchange, got %q", got)
	}
	if got := reqs[1].Header.Get("Authorization"); got != "Bearer tok-profile" {
		t.Fatalf("expected bearer auth from profile credentials, got %q", got)
	}
}

func TestTailnetLifecycleCommandErrorsAreActionable(t *testing.T) {
	t.Run("create tailnet requires display-name", func(t *testing.T) {
		res := executeCLINoDefaults(t, []string{"create", "tailnet", "--oauth-client-id", "cid", "--oauth-client-secret", "secret"}, nil)
		if res.err == nil || !strings.Contains(res.err.Error(), "--display-name is required") {
			t.Fatalf("expected display-name validation error, got %v", res.err)
		}
	})

	t.Run("delete tailnet requires id", func(t *testing.T) {
		res := executeCLINoDefaults(t, []string{"delete", "tailnet", "--oauth-client-id", "cid", "--oauth-client-secret", "secret"}, nil)
		if res.err == nil || !strings.Contains(res.err.Error(), "required flag(s) \"id\" not set") {
			t.Fatalf("expected id validation error, got %v", res.err)
		}
	})

	t.Run("lifecycle commands bypass api-key pre-run", func(t *testing.T) {
		res := executeCLINoDefaults(t, []string{"list", "tailnets", "--oauth-client-id", "cid", "--oauth-client-secret", "secret"}, map[string]string{
			"TSCLI_OAUTH_TOKEN_URL": "http://127.0.0.1:1/api/v2/oauth/token",
		})
		if res.err == nil {
			t.Fatalf("expected oauth exchange to fail")
		}
		if strings.Contains(strings.ToLower(res.err.Error()), "api key") {
			t.Fatalf("expected oauth exchange error instead of api-key pre-run error, got %v", res.err)
		}
	})

	t.Run("api errors are wrapped with lifecycle context", func(t *testing.T) {
		mock := apimock.New(t)
		mock.AddRaw(http.MethodPost, "/api/v2/oauth/token", http.StatusOK, `{"access_token":"tok-123","token_type":"Bearer","expires_in":3600}`)
		mock.AddRaw(http.MethodGet, "/api/v2/organizations/-/tailnets", http.StatusForbidden, `{"message":"forbidden"}`)

		res := executeCLINoDefaults(t, []string{"list", "tailnets", "--oauth-client-id", "cid", "--oauth-client-secret", "secret"}, map[string]string{
			"TSCLI_BASE_URL":        mock.URL(),
			"TSCLI_OAUTH_TOKEN_URL": mock.URL() + "/api/v2/oauth/token",
		})
		if res.err == nil {
			t.Fatalf("expected lifecycle API error")
		}
		if !strings.Contains(res.err.Error(), "list tailnets") || !strings.Contains(res.err.Error(), "forbidden") {
			t.Fatalf("expected wrapped lifecycle API error, got %v", res.err)
		}
	})
}
