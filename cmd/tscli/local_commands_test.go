package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestVersionCommand(t *testing.T) {
	res := executeCLI(t, []string{"version"}, nil)
	if res.err != nil {
		t.Fatalf("unexpected error: %v", res.err)
	}
	if strings.TrimSpace(res.stdout) == "" {
		t.Fatalf("expected non-empty version output")
	}
}

func TestConfigPrecedenceFlagOverEnvOverConfig(t *testing.T) {
	home := t.TempDir()
	cfgPath := filepath.Join(home, ".tscli.yaml")
	cfg := "api-key: from-config\ntailnet: from-config\noutput: json\n"
	if err := os.WriteFile(cfgPath, []byte(cfg), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	t.Run("config value used", func(t *testing.T) {
		res := executeCLINoDefaults(t, []string{"config", "get", "api-key"}, map[string]string{
			"HOME": home,
		})
		if res.err != nil {
			t.Fatalf("unexpected error: %v", res.err)
		}
		if got := strings.TrimSpace(res.stdout); got != "from-config" {
			t.Fatalf("expected from-config, got %q", got)
		}
	})

	t.Run("env overrides config", func(t *testing.T) {
		res := executeCLINoDefaults(t, []string{"config", "get", "api-key"}, map[string]string{
			"HOME":              home,
			"TAILSCALE_API_KEY": "from-env",
		})
		if res.err != nil {
			t.Fatalf("unexpected error: %v", res.err)
		}
		if got := strings.TrimSpace(res.stdout); got != "from-env" {
			t.Fatalf("expected from-env, got %q", got)
		}
	})

	t.Run("flag overrides env", func(t *testing.T) {
		res := executeCLINoDefaults(t, []string{"--api-key", "from-flag", "config", "get", "api-key"}, map[string]string{
			"HOME":              home,
			"TAILSCALE_API_KEY": "from-env",
		})
		if res.err != nil {
			t.Fatalf("unexpected error: %v", res.err)
		}
		if got := strings.TrimSpace(res.stdout); got != "from-flag" {
			t.Fatalf("expected from-flag, got %q", got)
		}
	})
}
