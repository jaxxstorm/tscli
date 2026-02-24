package cli_test

import (
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
