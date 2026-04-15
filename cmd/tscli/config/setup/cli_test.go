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

func TestInitialSetupAddAnotherNoTransitionsToOutputChoice(t *testing.T) {
	m := model{step: stepAddAnother, message: "Tailnet profile sandbox created.", initialSetup: true}
	m.input = "no"

	updatedModel, cmd := m.submit()
	updated := updatedModel.(model)
	if cmd != nil {
		t.Fatalf("expected no quit command")
	}
	if updated.step != stepOutputChoice {
		t.Fatalf("expected output choice step, got %q", updated.step)
	}
	if updated.choiceIndex != 0 {
		t.Fatalf("expected json output choice by default, got %d", updated.choiceIndex)
	}
}

func TestInitialSetupDebugChoicePersistsPreferences(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)
	home := t.TempDir()
	t.Setenv("HOME", home)

	m := model{
		step:             stepDebugChoice,
		message:          "Tailnet profile sandbox created.",
		initialSetup:     true,
		outputPreference: "human",
	}
	m.input = "yes"

	updatedModel, cmd := m.submit()
	updated := updatedModel.(model)
	if updated.err != nil {
		t.Fatalf("expected no error, got %v", updated.err)
	}
	if updated.step != stepDone {
		t.Fatalf("expected done step, got %q", updated.step)
	}
	if !updated.debugPreference {
		t.Fatalf("expected debug preference to be enabled")
	}
	if cmd == nil {
		t.Fatalf("expected quit command")
	}
	if msg := cmd(); msg != tea.Quit() {
		t.Fatalf("expected tea quit command, got %v", msg)
	}

	cfg, err := os.ReadFile(filepath.Join(home, ".tscli.yaml"))
	if err != nil {
		t.Fatalf("read config file: %v", err)
	}
	body := string(cfg)
	if !strings.Contains(body, "output: human") {
		t.Fatalf("expected output preference in config, got:\n%s", body)
	}
	if !strings.Contains(body, "debug: true") {
		t.Fatalf("expected debug preference in config, got:\n%s", body)
	}
}

func TestActionChoiceSettingsTransitionsToOutputChoice(t *testing.T) {
	m := model{step: stepActionChoice, outputPreference: "human"}
	m.input = "settings"

	updatedModel, cmd := m.submit()
	updated := updatedModel.(model)
	if cmd != nil {
		t.Fatalf("expected no quit command")
	}
	if updated.step != stepOutputChoice {
		t.Fatalf("expected output choice step, got %q", updated.step)
	}
	if updated.choiceIndex != 2 {
		t.Fatalf("expected human output choice index, got %d", updated.choiceIndex)
	}
}

func TestActionChoiceIncludesPreferencesOption(t *testing.T) {
	m := model{step: stepActionChoice}

	choices, ok := m.currentChoices()
	if !ok {
		t.Fatalf("expected action choice options")
	}
	if len(choices) != 5 {
		t.Fatalf("expected five action choices, got %d", len(choices))
	}
	if choices[3].value != "settings" || choices[3].label != "Modify preferences" {
		t.Fatalf("unexpected preferences choice: %+v", choices[3])
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

func TestSelectProfileUsesEncryptedAuthShape(t *testing.T) {
	m := model{
		step: stepSelectProfile,
		profiles: config.TailnetProfilesState{Tailnets: []config.TailnetProfile{{
			Name:                       "org-admin",
			OAuthClientID:              "cid",
			OAuthClientSecretEncrypted: "ciphertext",
		}}},
	}

	updatedModel, _ := m.submit()
	updated := updatedModel.(model)
	if updated.authType != "oauth" {
		t.Fatalf("expected oauth auth type, got %q", updated.authType)
	}
	if updated.choiceIndex != 1 {
		t.Fatalf("expected oauth selection index, got %d", updated.choiceIndex)
	}
}

func TestModifyEncryptedProfileKeepsExistingSecrets(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)
	t.Setenv("HOME", t.TempDir())

	m := model{
		step:     stepAPIKey,
		editing:  true,
		authType: "api-key",
		profile: config.TailnetProfile{
			Name:            "sandbox",
			APIKeyEncrypted: "ciphertext",
		},
	}

	updatedModel, cmd := m.submit()
	updated := updatedModel.(model)
	if cmd != nil {
		t.Fatalf("expected no quit command")
	}
	if updated.err != nil {
		t.Fatalf("expected encrypted api-key to be treated as present, got %v", updated.err)
	}
	if updated.step != stepActionChoice {
		t.Fatalf("expected successful modify to return to action choice, got %q", updated.step)
	}

	m = model{
		step:     stepAPIKey,
		editing:  true,
		authType: "api-key",
		input:    "new-key",
		profile: config.TailnetProfile{
			Name:            "sandbox",
			APIKeyEncrypted: "ciphertext",
		},
	}

	updatedModel, cmd = m.submit()
	updated = updatedModel.(model)
	if cmd != nil {
		t.Fatalf("expected no quit command")
	}
	if updated.err != nil {
		t.Fatalf("expected encrypted api-key replacement to succeed, got %v", updated.err)
	}

	m = model{
		step:     stepOAuthClientSecret,
		editing:  true,
		authType: "oauth",
		profile: config.TailnetProfile{
			Name:                       "org-admin",
			OAuthClientID:              "cid",
			OAuthClientSecretEncrypted: "ciphertext",
		},
	}

	updatedModel, cmd = m.submit()
	updated = updatedModel.(model)
	if cmd != nil {
		t.Fatalf("expected no quit command")
	}
	if updated.err != nil {
		t.Fatalf("expected encrypted oauth secret to be treated as present, got %v", updated.err)
	}
	if updated.step != stepActionChoice {
		t.Fatalf("expected successful modify to return to action choice, got %q", updated.step)
	}

	m = model{
		step:     stepOAuthClientSecret,
		editing:  true,
		authType: "oauth",
		input:    "updated-secret",
		profile: config.TailnetProfile{
			Name:                       "org-admin",
			OAuthClientID:              "cid",
			OAuthClientSecretEncrypted: "ciphertext",
		},
	}

	updatedModel, cmd = m.submit()
	updated = updatedModel.(model)
	if cmd != nil {
		t.Fatalf("expected no quit command")
	}
	if updated.err != nil {
		t.Fatalf("expected encrypted oauth secret replacement to succeed, got %v", updated.err)
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
