package config

import (
	"strings"
	"testing"

	"github.com/spf13/viper"
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
	if got := viper.GetString("tailnet"); got != "prod" {
		t.Fatalf("expected legacy tailnet mirrored to prod, got %q", got)
	}
	if got := viper.GetString("api-key"); got != "tskey-prod" {
		t.Fatalf("expected legacy api-key mirrored to active profile, got %q", got)
	}

	if err := RemoveTailnetProfile("prod"); err == nil {
		t.Fatalf("expected deleting active profile to fail")
	}

	if err := RemoveTailnetProfile("sandbox"); err != nil {
		t.Fatalf("remove non-active profile: %v", err)
	}
}
