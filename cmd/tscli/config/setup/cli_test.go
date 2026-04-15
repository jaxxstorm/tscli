package setup

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jaxxstorm/tscli/pkg/config"
	"github.com/spf13/viper"
)

func TestNewModelStartsAtExpectedStep(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)

	t.Run("no profiles starts at encryption choice", func(t *testing.T) {
		home := t.TempDir()
		t.Setenv("HOME", home)

		m, err := newModel()
		if err != nil {
			t.Fatalf("newModel: %v", err)
		}
		if m.step != stepEncryptionChoice {
			t.Fatalf("expected encryption choice step, got %q", m.step)
		}
	})

	t.Run("existing profiles start at action choice", func(t *testing.T) {
		viper.Reset()
		home := t.TempDir()
		t.Setenv("HOME", home)
		viper.Set("active-tailnet", "sandbox")
		viper.Set("tailnets", []map[string]any{{
			"name":    "sandbox",
			"api-key": "tskey-sandbox",
		}})

		m, err := newModel()
		if err != nil {
			t.Fatalf("newModel: %v", err)
		}
		if m.step != stepActionChoice {
			t.Fatalf("expected action choice step, got %q", m.step)
		}
	})
}

func TestResolveAgeKeyPath(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	path, err := resolveAgeKeyPath("")
	if err != nil {
		t.Fatalf("resolve default path: %v", err)
	}
	if want := filepath.Join(home, ".tscli", "age.txt"); path != want {
		t.Fatalf("expected %q, got %q", want, path)
	}

	custom, err := resolveAgeKeyPath("~/.config/tscli/age.txt")
	if err != nil {
		t.Fatalf("resolve custom path: %v", err)
	}
	if want := filepath.Join(home, ".config", "tscli", "age.txt"); custom != want {
		t.Fatalf("expected %q, got %q", want, custom)
	}
}

func TestGenerateAndPersistAgeConfigWritesKeyFile(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)

	home := t.TempDir()
	t.Setenv("HOME", home)

	path := filepath.Join(home, ".tscli", "age.txt")
	cfg, err := generateAndPersistAgeConfig(path)
	if err != nil {
		t.Fatalf("generateAndPersistAgeConfig: %v", err)
	}
	if cfg.PrivateKeyPath != path {
		t.Fatalf("expected private key path %q, got %q", path, cfg.PrivateKeyPath)
	}

	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read key file: %v", err)
	}
	if !strings.Contains(string(body), "AGE-SECRET-KEY-") {
		t.Fatalf("expected private key in key file, got:\n%s", string(body))
	}
	if !strings.Contains(string(body), "public-key:") {
		t.Fatalf("expected public key comment in key file, got:\n%s", string(body))
	}
}

func TestAddAnotherNoTransitionsToDone(t *testing.T) {
	m := model{step: stepAddAnother, message: "Tailnet profile sandbox created."}
	m.input = "no"

	updatedModel, cmd := m.submit()
	updated := updatedModel.(model)
	if updated.step != stepDone {
		t.Fatalf("expected done step, got %q", updated.step)
	}
	if !strings.Contains(updated.message, "Setup complete") {
		t.Fatalf("expected setup complete message, got %q", updated.message)
	}
	if cmd == nil {
		t.Fatalf("expected quit command")
	}
	if msg := cmd(); msg != tea.Quit() {
		t.Fatalf("expected tea quit command, got %v", msg)
	}
}

func TestSelectProfileUsesCurrentSelection(t *testing.T) {
	m := model{
		step: stepSelectProfile,
		profiles: config.TailnetProfilesState{Tailnets: []config.TailnetProfile{{
			Name:              "org-admin",
			Tailnet:           "example.ts.net",
			OAuthClientID:     "cid",
			OAuthClientSecret: "secret",
		}}},
	}

	updatedModel, cmd := m.submit()
	updated := updatedModel.(model)
	if cmd != nil {
		t.Fatalf("expected no quit command")
	}
	if updated.step != stepAuthType {
		t.Fatalf("expected auth type step, got %q", updated.step)
	}
	if !updated.editing {
		t.Fatalf("expected editing mode to be enabled")
	}
	if updated.profile.Name != "org-admin" {
		t.Fatalf("expected selected profile to load, got %+v", updated.profile)
	}
	if updated.authType != "oauth" {
		t.Fatalf("expected oauth auth type, got %q", updated.authType)
	}
	if updated.choiceIndex != 1 {
		t.Fatalf("expected oauth selection index, got %d", updated.choiceIndex)
	}
}

func TestDeleteProfileChoicesIncludeExistingProfiles(t *testing.T) {
	m := model{step: stepDeleteProfile, profiles: config.TailnetProfilesState{
		ActiveTailnet: "sandbox",
		Tailnets:      []config.TailnetProfile{{Name: "sandbox"}, {Name: "prod"}},
	}}

	choices, ok := m.currentChoices()
	if !ok {
		t.Fatalf("expected delete step to expose choices")
	}
	if len(choices) != 2 {
		t.Fatalf("expected two profile choices, got %d", len(choices))
	}
	if choices[0].value != "sandbox" || choices[0].label != "sandbox (active)" {
		t.Fatalf("unexpected first choice: %+v", choices[0])
	}
	if choices[1].value != "prod" || choices[1].label != "prod" {
		t.Fatalf("unexpected second choice: %+v", choices[1])
	}
}
