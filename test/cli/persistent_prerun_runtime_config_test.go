package cli_test

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"filippo.io/age"
	"github.com/jaxxstorm/tscli/internal/testutil/apimock"
)

func TestCommandsWithPreRunUseActiveProfileRuntimeConfig(t *testing.T) {
	home := t.TempDir()
	configFile := filepath.Join(home, ".tscli.yaml")
	cfg := strings.Join([]string{
		"output: json",
		"active-tailnet: profile-tailnet",
		"tailnets:",
		"  - name: profile-tailnet",
		"    api-key: tskey-profile",
		"",
	}, "\n")
	if err := os.WriteFile(configFile, []byte(cfg), 0o600); err != nil {
		t.Fatalf("write config file: %v", err)
	}

	cases := []struct {
		name     string
		args     []string
		method   string
		pathHint string
		body     any
	}{
		{
			name:     "get service",
			args:     []string{"get", "service", "--service", "svc:portal"},
			method:   http.MethodGet,
			pathHint: "/tailnet/profile-tailnet/services/svc:portal",
			body:     apimock.Service(),
		},
		{
			name:     "get service approval",
			args:     []string{"get", "service", "approval", "--service", "svc:portal", "--device", "node-123"},
			method:   http.MethodGet,
			pathHint: "/tailnet/profile-tailnet/services/svc:portal/device/node-123/approved",
			body:     apimock.ServiceApproval(),
		},
		{
			name:     "list services devices",
			args:     []string{"list", "services", "devices", "--service", "svc:portal"},
			method:   http.MethodGet,
			pathHint: "/tailnet/profile-tailnet/services/svc:portal/devices",
			body:     apimock.ServiceDevices(),
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mock := apimock.New(t)
			mock.AddJSON(tc.method, tc.pathHint, http.StatusOK, tc.body)

			res := executeCLINoDefaults(t, tc.args, map[string]string{
				"HOME":           home,
				"TSCLI_BASE_URL": mock.URL(),
			})
			if res.err != nil {
				t.Fatalf("unexpected error: %v\nstderr:\n%s", res.err, res.stderr)
			}

			reqs := mock.Requests()
			if len(reqs) == 0 {
				t.Fatalf("expected request to mock API, got none")
			}
			if !strings.Contains(reqs[0].Path, tc.pathHint) {
				t.Fatalf("expected request path %q, got %q", tc.pathHint, reqs[0].Path)
			}
		})
	}
}

func TestCommandsWithPreRunUseActiveOAuthProfileRuntimeConfig(t *testing.T) {
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
	mock.AddJSON(http.MethodGet, "/tailnet/org-admin/devices", http.StatusOK, apimock.DeviceList())

	res := executeCLINoDefaults(t, []string{"list", "devices"}, map[string]string{
		"HOME":           home,
		"TSCLI_BASE_URL": mock.URL(),
	})
	if res.err != nil {
		t.Fatalf("unexpected error: %v\nstderr:\n%s", res.err, res.stderr)
	}

	reqs := mock.Requests()
	if len(reqs) < 2 {
		t.Fatalf("expected token exchange and api request, got %+v", reqs)
	}
	if got := reqs[1].Header.Get("Authorization"); got != "Bearer tok-profile" {
		t.Fatalf("expected bearer auth on general API request, got %q", got)
	}
}

func TestCommandsWithPreRunDecryptEncryptedActiveProfile(t *testing.T) {
	identity, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatalf("generate identity: %v", err)
	}

	home := t.TempDir()
	res := executeCLINoDefaults(t, []string{"config", "encryption", "setup", "--public-key", identity.Recipient().String(), "--private-key-source", "env"}, map[string]string{
		"HOME": home,
	})
	if res.err != nil {
		t.Fatalf("config encryption setup: %v\nstderr:\n%s", res.err, res.stderr)
	}

	res = executeCLINoDefaults(t, []string{"config", "profiles", "set", "sandbox", "--api-key", "tskey-encrypted"}, map[string]string{
		"HOME": home,
	})
	if res.err != nil {
		t.Fatalf("upsert encrypted api profile: %v\nstderr:\n%s", res.err, res.stderr)
	}

	mock := apimock.New(t)
	mock.AddJSON(http.MethodGet, "/tailnet/sandbox/devices", http.StatusOK, apimock.DeviceList())

	res = executeCLINoDefaults(t, []string{"list", "devices"}, map[string]string{
		"HOME":                  home,
		"TSCLI_BASE_URL":        mock.URL(),
		"TSCLI_AGE_PRIVATE_KEY": identity.String(),
	})
	if res.err != nil {
		t.Fatalf("unexpected error: %v\nstderr:\n%s", res.err, res.stderr)
	}

	reqs := mock.Requests()
	if len(reqs) == 0 || !strings.Contains(reqs[0].Path, "/tailnet/sandbox/devices") {
		t.Fatalf("expected encrypted profile request path, got %+v", reqs)
	}
}
