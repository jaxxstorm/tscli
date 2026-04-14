package cli_test

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jaxxstorm/tscli/internal/testutil/apimock"
)

func TestConfigProfilesCommandFlow(t *testing.T) {
	home := t.TempDir()

	res := executeCLINoDefaults(t, []string{"config", "profiles", "upsert", "sandbox", "--api-key", "tskey-sandbox"}, map[string]string{
		"HOME": home,
	})
	if res.err != nil {
		t.Fatalf("upsert sandbox: %v\nstderr:\n%s", res.err, res.stderr)
	}
	if !strings.Contains(res.stdout, "created") {
		t.Fatalf("expected created message, got %q", res.stdout)
	}

	res = executeCLINoDefaults(t, []string{"config", "profiles", "upsert", "prod", "--api-key", "tskey-prod"}, map[string]string{
		"HOME": home,
	})
	if res.err != nil {
		t.Fatalf("upsert prod: %v\nstderr:\n%s", res.err, res.stderr)
	}

	res = executeCLINoDefaults(t, []string{"config", "profiles", "set-active", "prod"}, map[string]string{
		"HOME": home,
	})
	if res.err != nil {
		t.Fatalf("set-active: %v\nstderr:\n%s", res.err, res.stderr)
	}

	res = executeCLINoDefaults(t, []string{"config", "profiles", "list"}, map[string]string{
		"HOME":         home,
		"TSCLI_OUTPUT": "json",
	})
	if res.err != nil {
		t.Fatalf("list profiles: %v\nstderr:\n%s", res.err, res.stderr)
	}
	if !strings.Contains(res.stdout, `"active-tailnet"`) || !strings.Contains(res.stdout, `"prod"`) {
		t.Fatalf("expected active-tailnet in output, got %s", res.stdout)
	}
	if !strings.Contains(res.stdout, "sandbox") || !strings.Contains(res.stdout, "prod") {
		t.Fatalf("expected both profile names in output, got %s", res.stdout)
	}
	var listed map[string]any
	if err := json.Unmarshal([]byte(res.stdout), &listed); err != nil {
		t.Fatalf("unmarshal profile list output: %v\noutput:\n%s", err, res.stdout)
	}
	tailnets, _ := listed["tailnets"].([]any)
	if len(tailnets) == 0 {
		t.Fatalf("expected tailnets in output, got %s", res.stdout)
	}
	first, _ := tailnets[0].(map[string]any)
	if first["auth-type"] != "api-key" {
		t.Fatalf("expected auth-type api-key, got %v", first["auth-type"])
	}

	res = executeCLINoDefaults(t, []string{"config", "profiles", "delete", "sandbox"}, map[string]string{
		"HOME": home,
	})
	if res.err != nil {
		t.Fatalf("delete non-active profile: %v\nstderr:\n%s", res.err, res.stderr)
	}

	res = executeCLINoDefaults(t, []string{"config", "profiles", "delete", "prod"}, map[string]string{
		"HOME": home,
	})
	if res.err == nil {
		t.Fatalf("expected deleting active profile to fail")
	}
	if !strings.Contains(strings.ToLower(res.err.Error()), "active") {
		t.Fatalf("expected active-profile guidance error, got %v", res.err)
	}

	configFile := filepath.Join(home, ".tscli.yaml")
	cfg, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatalf("read config file: %v", err)
	}
	if !strings.Contains(string(cfg), "active-tailnet: prod") {
		t.Fatalf("expected persisted active-tailnet in config file, got:\n%s", string(cfg))
	}
	if strings.Contains(string(cfg), "\ntailnet:") {
		t.Fatalf("did not expect duplicated top-level tailnet in config file, got:\n%s", string(cfg))
	}
	if strings.Contains(string(cfg), "\napi-key:") {
		t.Fatalf("did not expect duplicated top-level api-key in config file, got:\n%s", string(cfg))
	}
}

func TestConfigProfilesSupportOAuthCredentials(t *testing.T) {
	home := t.TempDir()

	res := executeCLINoDefaults(t, []string{"config", "profiles", "upsert", "org-admin", "--oauth-client-id", "cid-org", "--oauth-client-secret", "secret-org"}, map[string]string{
		"HOME": home,
	})
	if res.err != nil {
		t.Fatalf("upsert oauth profile: %v\nstderr:\n%s", res.err, res.stderr)
	}
	if !strings.Contains(res.stdout, "created") {
		t.Fatalf("expected created message, got %q", res.stdout)
	}

	res = executeCLINoDefaults(t, []string{"config", "profiles", "list"}, map[string]string{
		"HOME":         home,
		"TSCLI_OUTPUT": "json",
	})
	if res.err != nil {
		t.Fatalf("list oauth profiles: %v\nstderr:\n%s", res.err, res.stderr)
	}
	var listed map[string]any
	if err := json.Unmarshal([]byte(res.stdout), &listed); err != nil {
		t.Fatalf("unmarshal oauth profile list output: %v\noutput:\n%s", err, res.stdout)
	}
	tailnets, _ := listed["tailnets"].([]any)
	if len(tailnets) != 1 {
		t.Fatalf("expected one oauth profile in output, got %s", res.stdout)
	}
	first, _ := tailnets[0].(map[string]any)
	if first["auth-type"] != "oauth" {
		t.Fatalf("expected oauth auth type in list output, got %v", first["auth-type"])
	}

	configFile := filepath.Join(home, ".tscli.yaml")
	cfg, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatalf("read config file: %v", err)
	}
	if !strings.Contains(string(cfg), "oauth-client-id: cid-org") || !strings.Contains(string(cfg), "oauth-client-secret: secret-org") {
		t.Fatalf("expected persisted oauth profile credentials, got:\n%s", string(cfg))
	}
	if strings.Contains(string(cfg), "\noauth-client-id:") || strings.Contains(string(cfg), "\noauth-client-secret:") {
		// top-level keys would appear at column 0; nested profile keys are indented
		for _, line := range strings.Split(string(cfg), "\n") {
			if strings.HasPrefix(line, "oauth-client-id:") || strings.HasPrefix(line, "oauth-client-secret:") {
				t.Fatalf("did not expect duplicated top-level oauth keys in config file, got:\n%s", string(cfg))
			}
		}
	}

	res = executeCLINoDefaults(t, []string{"config", "profiles", "set-active", "org-admin"}, map[string]string{"HOME": home})
	if res.err != nil {
		t.Fatalf("set active oauth profile: %v\nstderr:\n%s", res.err, res.stderr)
	}

	res = executeCLINoDefaults(t, []string{"config", "profiles", "delete", "org-admin"}, map[string]string{"HOME": home})
	if res.err == nil {
		t.Fatalf("expected deleting active oauth profile to fail")
	}
}

func TestConfigProfilesUpsertRejectsMixedAuthShapes(t *testing.T) {
	res := executeCLINoDefaults(t, []string{"config", "profiles", "upsert", "mixed", "--api-key", "tskey-mixed", "--oauth-client-id", "cid", "--oauth-client-secret", "secret"}, nil)
	if res.err == nil {
		t.Fatalf("expected mixed auth shape to fail")
	}
	if !strings.Contains(strings.ToLower(res.err.Error()), "either api-key auth or oauth-client-id/oauth-client-secret auth") {
		t.Fatalf("expected mixed auth shape error, got %v", res.err)
	}
}

func TestRuntimeUsesActiveProfileWithoutEnvOrFlags(t *testing.T) {
	home := t.TempDir()
	configFile := filepath.Join(home, ".tscli.yaml")
	cfg := strings.Join([]string{
		"output: json",
		"active-tailnet: profile-tailnet",
		"tailnets:",
		"  - name: profile-tailnet",
		"    api-key: tskey-profile",
		"tailnet: legacy-tailnet",
		"api-key: legacy-key",
		"",
	}, "\n")
	if err := os.WriteFile(configFile, []byte(cfg), 0o600); err != nil {
		t.Fatalf("write config file: %v", err)
	}

	mock := apimock.New(t)
	mock.AddJSON(http.MethodGet, "/tailnet/profile-tailnet/devices", http.StatusOK, apimock.DeviceList())

	res := executeCLINoDefaults(t, []string{"list", "devices"}, map[string]string{
		"HOME":           home,
		"TSCLI_BASE_URL": mock.URL(),
	})
	if res.err != nil {
		t.Fatalf("unexpected error: %v\nstderr:\n%s", res.err, res.stderr)
	}

	reqs := mock.Requests()
	if len(reqs) == 0 || !strings.Contains(reqs[0].Path, "/tailnet/profile-tailnet/devices") {
		t.Fatalf("expected request path to use active profile tailnet, got %+v", reqs)
	}
}

func TestSwitchingActiveProfileChangesRuntimeTailnet(t *testing.T) {
	home := t.TempDir()

	for _, args := range [][]string{
		{"config", "profiles", "upsert", "sandbox", "--api-key", "tskey-sandbox"},
		{"config", "profiles", "upsert", "prod", "--api-key", "tskey-prod"},
	} {
		res := executeCLINoDefaults(t, args, map[string]string{"HOME": home})
		if res.err != nil {
			t.Fatalf("setup %v: %v\nstderr:\n%s", args, res.err, res.stderr)
		}
	}

	mock := apimock.New(t)
	mock.AddJSON(http.MethodGet, "/tailnet/prod/devices", http.StatusOK, apimock.DeviceList())
	mock.AddJSON(http.MethodGet, "/tailnet/sandbox/devices", http.StatusOK, apimock.DeviceList())

	res := executeCLINoDefaults(t, []string{"config", "profiles", "set-active", "prod"}, map[string]string{"HOME": home})
	if res.err != nil {
		t.Fatalf("set active prod: %v\nstderr:\n%s", res.err, res.stderr)
	}

	res = executeCLINoDefaults(t, []string{"list", "devices"}, map[string]string{
		"HOME":           home,
		"TSCLI_BASE_URL": mock.URL(),
	})
	if res.err != nil {
		t.Fatalf("list devices with prod active: %v\nstderr:\n%s", res.err, res.stderr)
	}

	res = executeCLINoDefaults(t, []string{"config", "profiles", "set-active", "sandbox"}, map[string]string{"HOME": home})
	if res.err != nil {
		t.Fatalf("set active sandbox: %v\nstderr:\n%s", res.err, res.stderr)
	}

	res = executeCLINoDefaults(t, []string{"list", "devices"}, map[string]string{
		"HOME":           home,
		"TSCLI_BASE_URL": mock.URL(),
	})
	if res.err != nil {
		t.Fatalf("list devices with sandbox active: %v\nstderr:\n%s", res.err, res.stderr)
	}

	reqs := mock.Requests()
	if len(reqs) < 2 {
		t.Fatalf("expected requests for both active profiles, got %+v", reqs)
	}
	if !strings.Contains(reqs[0].Path, "/tailnet/prod/devices") {
		t.Fatalf("expected first request to use prod tailnet, got %+v", reqs[0])
	}
	if !strings.Contains(reqs[1].Path, "/tailnet/sandbox/devices") {
		t.Fatalf("expected second request to use sandbox tailnet, got %+v", reqs[1])
	}
}

func TestConfigShowNormalizesProfileBackedConfig(t *testing.T) {
	home := t.TempDir()
	configFile := filepath.Join(home, ".tscli.yaml")
	cfg := strings.Join([]string{
		"output: pretty",
		"debug: false",
		"help: false",
		"active-tailnet: prod",
		"tailnet: prod",
		"api-key: tskey-prod",
		"tailnets:",
		"  - name: prod",
		"    api-key: tskey-prod",
		"",
	}, "\n")
	if err := os.WriteFile(configFile, []byte(cfg), 0o600); err != nil {
		t.Fatalf("write config file: %v", err)
	}

	res := executeCLINoDefaults(t, []string{"config", "show"}, map[string]string{
		"HOME":         home,
		"TSCLI_OUTPUT": "json",
	})
	if res.err != nil {
		t.Fatalf("config show: %v\nstderr:\n%s", res.err, res.stderr)
	}

	var shown map[string]any
	if err := json.Unmarshal([]byte(res.stdout), &shown); err != nil {
		t.Fatalf("unmarshal config show output: %v\noutput:\n%s", err, res.stdout)
	}
	for _, unwanted := range []string{"tailnet", "api-key", "debug", "help"} {
		if _, ok := shown[unwanted]; ok {
			t.Fatalf("did not expect top-level %q in config show output: %s", unwanted, res.stdout)
		}
	}
	if _, ok := shown["active-tailnet"]; !ok {
		t.Fatalf("expected active-tailnet in output, got %s", res.stdout)
	}
	if _, ok := shown["tailnets"]; !ok {
		t.Fatalf("expected canonical profile keys in output, got %s", res.stdout)
	}
}

func TestRuntimeOverridePrecedenceForProfileConfig(t *testing.T) {
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

	t.Run("env overrides profile", func(t *testing.T) {
		mock := apimock.New(t)
		mock.AddJSON(http.MethodGet, "/tailnet/env-tailnet/devices", http.StatusOK, apimock.DeviceList())

		res := executeCLINoDefaults(t, []string{"list", "devices"}, map[string]string{
			"HOME":              home,
			"TSCLI_BASE_URL":    mock.URL(),
			"TAILSCALE_TAILNET": "env-tailnet",
			"TAILSCALE_API_KEY": "tskey-env",
		})
		if res.err != nil {
			t.Fatalf("unexpected error: %v\nstderr:\n%s", res.err, res.stderr)
		}

		reqs := mock.Requests()
		if len(reqs) == 0 || !strings.Contains(reqs[0].Path, "/tailnet/env-tailnet/devices") {
			t.Fatalf("expected env tailnet path, got %+v", reqs)
		}
	})

	t.Run("flags override env", func(t *testing.T) {
		mock := apimock.New(t)
		mock.AddJSON(http.MethodGet, "/tailnet/flag-tailnet/devices", http.StatusOK, apimock.DeviceList())

		res := executeCLINoDefaults(t, []string{"--tailnet", "flag-tailnet", "--api-key", "tskey-flag", "list", "devices"}, map[string]string{
			"HOME":              home,
			"TSCLI_BASE_URL":    mock.URL(),
			"TAILSCALE_TAILNET": "env-tailnet",
			"TAILSCALE_API_KEY": "tskey-env",
		})
		if res.err != nil {
			t.Fatalf("unexpected error: %v\nstderr:\n%s", res.err, res.stderr)
		}

		reqs := mock.Requests()
		if len(reqs) == 0 || !strings.Contains(reqs[0].Path, "/tailnet/flag-tailnet/devices") {
			t.Fatalf("expected flag tailnet path, got %+v", reqs)
		}
	})
}

func TestRuntimeFailsOnInvalidActiveTailnetReference(t *testing.T) {
	home := t.TempDir()
	configFile := filepath.Join(home, ".tscli.yaml")
	cfg := strings.Join([]string{
		"output: json",
		"active-tailnet: missing-tailnet",
		"tailnets:",
		"  - name: profile-tailnet",
		"    api-key: tskey-profile",
		"tailnet: legacy-tailnet",
		"api-key: legacy-key",
		"",
	}, "\n")
	if err := os.WriteFile(configFile, []byte(cfg), 0o600); err != nil {
		t.Fatalf("write config file: %v", err)
	}

	mock := apimock.New(t)
	mock.AddJSON(http.MethodGet, "/tailnet/", http.StatusOK, apimock.DeviceList())

	res := executeCLINoDefaults(t, []string{"list", "devices"}, map[string]string{
		"HOME":           home,
		"TSCLI_BASE_URL": mock.URL(),
	})
	if res.err == nil {
		t.Fatalf("expected invalid active-tailnet error")
	}
	if !strings.Contains(strings.ToLower(res.err.Error()), "active-tailnet") {
		t.Fatalf("expected active-tailnet error, got %v", res.err)
	}
}
