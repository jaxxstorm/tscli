package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func TestLoadTailnetProfilesStateValidation(t *testing.T) {
	t.Run("duplicate names fail", func(t *testing.T) {
		v := viper.New()
		v.Set("tailnets", []map[string]any{
			{"name": "sandbox", "api-key": "tskey-a"},
			{"name": "sandbox", "api-key": "tskey-b"},
		})

		_, err := loadTailnetProfilesState(v)
		if err == nil || !strings.Contains(err.Error(), "duplicate tailnet profile") {
			t.Fatalf("expected duplicate profile validation error, got %v", err)
		}
	})

	t.Run("missing api-key fails", func(t *testing.T) {
		v := viper.New()
		v.Set("tailnets", []map[string]any{
			{"name": "sandbox"},
		})

		_, err := loadTailnetProfilesState(v)
		if err == nil || !strings.Contains(err.Error(), "missing api-key") {
			t.Fatalf("expected missing api-key validation error, got %v", err)
		}
	})

	t.Run("missing active profile reference fails", func(t *testing.T) {
		v := viper.New()
		v.Set("active-tailnet", "unknown")
		v.Set("tailnets", []map[string]any{
			{"name": "sandbox", "api-key": "tskey-a"},
		})

		_, err := loadTailnetProfilesState(v)
		if err == nil || !strings.Contains(err.Error(), "active-tailnet") {
			t.Fatalf("expected active-tailnet validation error, got %v", err)
		}
	})

	t.Run("profiles are normalized and sorted", func(t *testing.T) {
		v := viper.New()
		v.Set("tailnets", []map[string]any{
			{"name": " zeta ", "api-key": " tskey-z "},
			{"name": "alpha", "api-key": "tskey-a"},
		})

		state, err := loadTailnetProfilesState(v)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(state.Tailnets) != 2 {
			t.Fatalf("expected 2 profiles, got %d", len(state.Tailnets))
		}
		if state.Tailnets[0].Name != "alpha" || state.Tailnets[1].Name != "zeta" {
			t.Fatalf("expected sorted profile names, got %+v", state.Tailnets)
		}
		if state.Tailnets[1].APIKey != "tskey-z" {
			t.Fatalf("expected trimmed api-key, got %q", state.Tailnets[1].APIKey)
		}
	})
}

func TestResolveRuntimeConfigPrecedence(t *testing.T) {
	t.Run("flags override env profile and legacy", func(t *testing.T) {
		v := viper.New()
		v.Set("active-tailnet", "profile-tailnet")
		v.Set("tailnets", []map[string]any{
			{"name": "profile-tailnet", "api-key": "profile-key"},
		})
		v.Set("tailnet", "legacy-tailnet")
		v.Set("api-key", "legacy-key")
		v.Set("tailnet", "flag-tailnet")
		v.Set("api-key", "flag-key")
		t.Setenv("TAILSCALE_TAILNET", "env-tailnet")
		t.Setenv("TAILSCALE_API_KEY", "env-key")

		resolved, err := resolveRuntimeConfig(v, map[string]struct{}{
			"tailnet": {},
			"api-key": {},
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resolved.Tailnet != "flag-tailnet" || resolved.APIKey != "flag-key" {
			t.Fatalf("expected flag values, got %+v", resolved)
		}
	})

	t.Run("env overrides profile and legacy", func(t *testing.T) {
		v := viper.New()
		v.Set("active-tailnet", "profile-tailnet")
		v.Set("tailnets", []map[string]any{
			{"name": "profile-tailnet", "api-key": "profile-key"},
		})
		v.Set("tailnet", "legacy-tailnet")
		v.Set("api-key", "legacy-key")
		t.Setenv("TAILSCALE_TAILNET", "env-tailnet")
		t.Setenv("TAILSCALE_API_KEY", "env-key")

		resolved, err := resolveRuntimeConfig(v, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resolved.Tailnet != "env-tailnet" || resolved.APIKey != "env-key" {
			t.Fatalf("expected env values, got %+v", resolved)
		}
	})

	t.Run("active profile overrides legacy", func(t *testing.T) {
		v := viper.New()
		v.Set("active-tailnet", "profile-tailnet")
		v.Set("tailnets", []map[string]any{
			{"name": "profile-tailnet", "api-key": "profile-key"},
		})
		v.Set("tailnet", "legacy-tailnet")
		v.Set("api-key", "legacy-key")

		resolved, err := resolveRuntimeConfig(v, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resolved.Tailnet != "profile-tailnet" || resolved.APIKey != "profile-key" {
			t.Fatalf("expected profile values, got %+v", resolved)
		}
	})

	t.Run("legacy config works without profiles", func(t *testing.T) {
		v := viper.New()
		v.Set("tailnet", "legacy-tailnet")
		v.Set("api-key", "legacy-key")

		resolved, err := resolveRuntimeConfig(v, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resolved.Tailnet != "legacy-tailnet" || resolved.APIKey != "legacy-key" {
			t.Fatalf("expected legacy values, got %+v", resolved)
		}
	})

	t.Run("missing tailnet defaults to dash", func(t *testing.T) {
		v := viper.New()
		v.Set("api-key", "legacy-key")

		resolved, err := resolveRuntimeConfig(v, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resolved.Tailnet != "-" {
			t.Fatalf("expected default tailnet '-', got %q", resolved.Tailnet)
		}
	})

	t.Run("missing api-key fails", func(t *testing.T) {
		v := viper.New()

		_, err := resolveRuntimeConfig(v, nil)
		if err == nil || !strings.Contains(err.Error(), "API key is required") {
			t.Fatalf("expected missing api-key error, got %v", err)
		}
	})
}

func TestTailnetProfilePersistenceHelpers(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)

	home := t.TempDir()
	t.Setenv("HOME", home)

	created, err := UpsertTailnetProfile("sandbox", "tskey-sandbox")
	if err != nil {
		t.Fatalf("upsert sandbox: %v", err)
	}
	if !created {
		t.Fatalf("expected first upsert to create profile")
	}

	created, err = UpsertTailnetProfile("prod", "tskey-prod")
	if err != nil {
		t.Fatalf("upsert prod: %v", err)
	}
	if !created {
		t.Fatalf("expected second upsert to create profile")
	}

	if err := SetActiveTailnet("prod"); err != nil {
		t.Fatalf("set active profile: %v", err)
	}

	state, err := ListTailnetProfiles()
	if err != nil {
		t.Fatalf("list profiles: %v", err)
	}
	if state.ActiveTailnet != "prod" {
		t.Fatalf("expected active profile prod, got %q", state.ActiveTailnet)
	}

	resolved, err := ResolveRuntimeConfig(nil)
	if err != nil {
		t.Fatalf("resolve runtime config: %v", err)
	}
	if resolved.Tailnet != "prod" || resolved.APIKey != "tskey-prod" {
		t.Fatalf("expected resolved active profile credentials, got %+v", resolved)
	}

	cfg, err := os.ReadFile(filepath.Join(home, ".tscli.yaml"))
	if err != nil {
		t.Fatalf("read persisted config: %v", err)
	}
	body := string(cfg)
	if !strings.Contains(body, "active-tailnet: prod") {
		t.Fatalf("expected active-tailnet in config, got:\n%s", body)
	}
	if !strings.Contains(body, "tailnets:") {
		t.Fatalf("expected tailnets block in config, got:\n%s", body)
	}

	var persisted map[string]any
	if err := yaml.Unmarshal(cfg, &persisted); err != nil {
		t.Fatalf("unmarshal persisted config: %v", err)
	}
	if _, ok := persisted["tailnet"]; ok {
		t.Fatalf("did not expect duplicated top-level tailnet in config, got:\n%s", body)
	}
	if _, ok := persisted["api-key"]; ok {
		t.Fatalf("did not expect duplicated top-level api-key in profile-backed config, got:\n%s", body)
	}

	if err := RemoveTailnetProfile("prod"); err == nil {
		t.Fatalf("expected deleting active profile to fail")
	}

	if err := RemoveTailnetProfile("sandbox"); err != nil {
		t.Fatalf("remove non-active profile: %v", err)
	}
}

func TestProfilePersistenceNormalizesMixedLegacyConfig(t *testing.T) {
	v := viper.New()

	home := t.TempDir()
	path := filepath.Join(home, ".tscli.yaml")
	data := strings.Join([]string{
		"output: pretty",
		"debug: false",
		"help: false",
		"active-tailnet: prod",
		"tailnet: prod",
		"api-key: tskey-prod",
		"tailnets:",
		"  - name: prod",
		"    api-key: tskey-prod",
		"  - name: sandbox",
		"    api-key: tskey-sandbox",
		"",
	}, "\n")
	if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
		t.Fatalf("write config file: %v", err)
	}

	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err != nil {
		t.Fatalf("read config file: %v", err)
	}

	if err := persistTailnetProfilesState(v, TailnetProfilesState{
		ActiveTailnet: "sandbox",
		Tailnets: []TailnetProfile{
			{Name: "prod", APIKey: "tskey-prod"},
			{Name: "sandbox", APIKey: "tskey-sandbox"},
		},
	}); err != nil {
		t.Fatalf("persist profiles: %v", err)
	}

	cfg, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read normalized config: %v", err)
	}
	body := string(cfg)

	var persisted map[string]any
	if err := yaml.Unmarshal(cfg, &persisted); err != nil {
		t.Fatalf("unmarshal normalized config: %v", err)
	}
	for _, unwanted := range []string{"tailnet", "api-key", "debug", "help"} {
		if _, ok := persisted[unwanted]; ok {
			t.Fatalf("did not expect %q in normalized config:\n%s", unwanted, body)
		}
	}
	if !strings.Contains(body, "output: pretty") {
		t.Fatalf("expected unrelated persisted keys to remain, got:\n%s", body)
	}
}
