package cli_test

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"filippo.io/age"
	"github.com/jaxxstorm/tscli/internal/testutil/apimock"
)

func TestConfigProfilesCommandFlow(t *testing.T) {
	home := t.TempDir()

	res := executeCLINoDefaults(t, []string{"config", "profiles", "set", "sandbox", "--api-key", "tskey-sandbox"}, map[string]string{
		"HOME": home,
	})
	if res.err != nil {
		t.Fatalf("upsert sandbox: %v\nstderr:\n%s", res.err, res.stderr)
	}
	if !strings.Contains(res.stdout, "created") {
		t.Fatalf("expected created message, got %q", res.stdout)
	}

	res = executeCLINoDefaults(t, []string{"config", "profiles", "set", "prod", "--api-key", "tskey-prod"}, map[string]string{
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

	configFile := filepath.Join(home, ".tscli.yaml")
	cfg, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatalf("read config file: %v", err)
	}
	if !strings.Contains(string(cfg), "active-tailnet: prod") {
		t.Fatalf("expected active-tailnet to remain prod after failed delete, got:\n%s", string(cfg))
	}
	if !strings.Contains(string(cfg), "tailnets:") || !strings.Contains(string(cfg), "name: prod") {
		t.Fatalf("expected remaining prod profile to stay persisted, got:\n%s", string(cfg))
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

	res := executeCLINoDefaults(t, []string{"config", "profiles", "set", "org-admin", "--oauth-client-id", "cid-org", "--oauth-client-secret", "secret-org"}, map[string]string{
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

func TestConfigEncryptionSetupPersistsAgeConfig(t *testing.T) {
	identity, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatalf("generate identity: %v", err)
	}

	home := t.TempDir()
	privateKeyPath := filepath.Join(home, "age.txt")
	if err := os.WriteFile(privateKeyPath, []byte(identity.String()+"\n"), 0o600); err != nil {
		t.Fatalf("write private key file: %v", err)
	}

	res := executeCLINoDefaults(t, []string{"config", "encryption", "setup", "--public-key", identity.Recipient().String(), "--private-key-source", "path", "--private-key-path", privateKeyPath}, map[string]string{
		"HOME": home,
	})
	if res.err != nil {
		t.Fatalf("config encryption setup: %v\nstderr:\n%s", res.err, res.stderr)
	}

	configFile := filepath.Join(home, ".tscli.yaml")
	cfg, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatalf("read config file: %v", err)
	}
	if !strings.Contains(string(cfg), "public-key: "+identity.Recipient().String()) {
		t.Fatalf("expected persisted age public key, got:\n%s", string(cfg))
	}
	if !strings.Contains(string(cfg), "private-key-path: "+privateKeyPath) {
		t.Fatalf("expected persisted age private key path, got:\n%s", string(cfg))
	}
}

func TestConfigEncryptionSetupPromptsForSupportedSources(t *testing.T) {
	home := t.TempDir()

	res := executeCLINoDefaultsWithInput(t, []string{"config", "encryption", "setup"}, map[string]string{
		"HOME": home,
	}, "age1invalid\n\n")
	if res.err == nil {
		t.Fatalf("expected setup to fail for missing private key source")
	}
	if !strings.Contains(res.stdout, "Private key source [path|env|command]:") {
		t.Fatalf("expected updated private key source prompt, got:\n%s", res.stdout)
	}
}

func TestConfigEncryptionSetupReusesExistingPathIdentity(t *testing.T) {
	identity, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatalf("generate identity: %v", err)
	}

	home := t.TempDir()
	privateKeyPath := filepath.Join(home, "age.txt")
	if err := os.WriteFile(privateKeyPath, []byte(identity.String()+"\n"), 0o600); err != nil {
		t.Fatalf("write private key file: %v", err)
	}

	res := executeCLINoDefaults(t, []string{"config", "encryption", "setup", "--private-key-source", "path", "--private-key-path", privateKeyPath}, map[string]string{
		"HOME": home,
	})
	if res.err != nil {
		t.Fatalf("config encryption setup reuse existing path: %v\nstderr:\n%s", res.err, res.stderr)
	}
	if !strings.Contains(res.stdout, "Reusing existing AGE identity") {
		t.Fatalf("expected reuse message, got:\n%s", res.stdout)
	}

	configFile := filepath.Join(home, ".tscli.yaml")
	cfg, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatalf("read config file: %v", err)
	}
	body := string(cfg)
	if !strings.Contains(body, "public-key: "+identity.Recipient().String()) {
		t.Fatalf("expected derived public key, got:\n%s", body)
	}
	if !strings.Contains(body, "private-key-path: "+privateKeyPath) {
		t.Fatalf("expected persisted private key path, got:\n%s", body)
	}
}

func TestConfigProfilesUpsertPromptsForOAuthCredentials(t *testing.T) {
	home := t.TempDir()

	res := executeCLINoDefaultsWithInput(t, []string{"config", "profiles", "set", "org-admin"}, map[string]string{
		"HOME": home,
	}, "oauth\ncid\nsecret\n")
	if res.err != nil {
		t.Fatalf("interactive oauth upsert: %v\nstderr:\n%s", res.err, res.stderr)
	}
	if !strings.Contains(res.stdout, "Auth type [api-key|oauth]:") || !strings.Contains(res.stdout, "OAuth client ID:") || !strings.Contains(res.stdout, "OAuth client secret:") {
		t.Fatalf("expected interactive oauth prompts, got:\n%s", res.stdout)
	}

	configFile := filepath.Join(home, ".tscli.yaml")
	cfg, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatalf("read config file: %v", err)
	}
	body := string(cfg)
	if !strings.Contains(body, "oauth-client-id: cid") || !strings.Contains(body, "oauth-client-secret: secret") {
		t.Fatalf("expected oauth credentials in config file, got:\n%s", body)
	}
}

func TestConfigSetupCreatesPlaintextProfile(t *testing.T) {
	home := t.TempDir()

	res := executeCLINoDefaultsWithInput(t, []string{"config", "setup"}, map[string]string{
		"HOME": home,
	}, "no\napi-key\nsandbox\n\ntskey-sandbox\nno\njson\nno\n")
	if res.err != nil {
		t.Fatalf("config setup plaintext: %v\nstderr:\n%s", res.err, res.stderr)
	}
	if !strings.Contains(res.stdout, "Encrypt your credentials?") || !strings.Contains(res.stdout, "Use an API key (it will expire) or OAuth credentials?") {
		t.Fatalf("expected setup prompts, got:\n%s", res.stdout)
	}
	if !strings.Contains(res.stdout, "Choose your default output format") || !strings.Contains(res.stdout, "Enable debug HTTP request/response logging by default?") {
		t.Fatalf("expected output/debug prompts, got:\n%s", res.stdout)
	}

	configFile := filepath.Join(home, ".tscli.yaml")
	cfg, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatalf("read config file: %v", err)
	}
	body := string(cfg)
	if !strings.Contains(body, "active-tailnet: sandbox") || !strings.Contains(body, "api-key: tskey-sandbox") {
		t.Fatalf("expected plaintext profile in config file, got:\n%s", body)
	}
	if !strings.Contains(body, "output: json") || !strings.Contains(body, "debug: false") {
		t.Fatalf("expected persisted output/debug defaults in config file, got:\n%s", body)
	}
	if strings.Contains(body, "api-key-encrypted:") {
		t.Fatalf("did not expect encrypted api key in plaintext setup, got:\n%s", body)
	}
}

func TestConfigSetupCreatesEncryptedProfileAndKeyFile(t *testing.T) {
	home := t.TempDir()

	res := executeCLINoDefaultsWithInput(t, []string{"config", "setup"}, map[string]string{
		"HOME": home,
	}, "yes\n\napi-key\nsandbox\n\ntskey-sandbox\nno\npretty\nyes\n")
	if res.err != nil {
		t.Fatalf("config setup encrypted: %v\nstderr:\n%s", res.err, res.stderr)
	}

	keyFile := filepath.Join(home, ".tscli", "age.txt")
	data, err := os.ReadFile(keyFile)
	if err != nil {
		t.Fatalf("read generated key file: %v", err)
	}
	if !strings.Contains(string(data), "AGE-SECRET-KEY-") {
		t.Fatalf("expected generated private key file, got:\n%s", string(data))
	}

	configFile := filepath.Join(home, ".tscli.yaml")
	cfg, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatalf("read config file: %v", err)
	}
	body := string(cfg)
	if !strings.Contains(body, "public-key:") || !strings.Contains(body, "private-key-path: "+keyFile) {
		t.Fatalf("expected persisted encryption config, got:\n%s", body)
	}
	if !strings.Contains(body, "output: pretty") || !strings.Contains(body, "debug: true") {
		t.Fatalf("expected persisted output/debug defaults, got:\n%s", body)
	}
	if !strings.Contains(body, "api-key-encrypted:") {
		t.Fatalf("expected encrypted api key in config file, got:\n%s", body)
	}
	if strings.Contains(body, "api-key: tskey-sandbox") {
		t.Fatalf("did not expect plaintext api key in config file, got:\n%s", body)
	}
}

func TestConfigSetupReusesExistingAgeIdentityFile(t *testing.T) {
	identity, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatalf("generate identity: %v", err)
	}

	home := t.TempDir()
	keyFile := filepath.Join(home, ".tscli", "age.txt")
	if err := os.MkdirAll(filepath.Dir(keyFile), 0o755); err != nil {
		t.Fatalf("create key dir: %v", err)
	}
	originalBody := "# existing key\n# public-key: " + identity.Recipient().String() + "\n" + identity.String() + "\n"
	if err := os.WriteFile(keyFile, []byte(originalBody), 0o600); err != nil {
		t.Fatalf("write age key file: %v", err)
	}

	res := executeCLINoDefaultsWithInput(t, []string{"config", "setup"}, map[string]string{
		"HOME": home,
	}, "yes\n\nyes\napi-key\nsandbox\n\ntskey-sandbox\nno\njson\nno\n")
	if res.err != nil {
		t.Fatalf("config setup reuse existing key: %v\nstderr:\n%s", res.err, res.stderr)
	}
	if !strings.Contains(res.stdout, "Reuse the existing age identity?") {
		t.Fatalf("expected reuse prompt, got:\n%s", res.stdout)
	}
	if !strings.Contains(res.stdout, "Use arrow keys to choose") {
		t.Fatalf("expected structured choice output, got:\n%s", res.stdout)
	}

	data, err := os.ReadFile(keyFile)
	if err != nil {
		t.Fatalf("read key file: %v", err)
	}
	if string(data) != originalBody {
		t.Fatalf("expected existing key file to remain unchanged, got:\n%s", string(data))
	}

	configFile := filepath.Join(home, ".tscli.yaml")
	cfg, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatalf("read config file: %v", err)
	}
	body := string(cfg)
	if !strings.Contains(body, "public-key: "+identity.Recipient().String()) {
		t.Fatalf("expected reused public key in config, got:\n%s", body)
	}
	if !strings.Contains(body, "private-key-path: "+keyFile) {
		t.Fatalf("expected reused private key path in config, got:\n%s", body)
	}
	if !strings.Contains(body, "api-key-encrypted:") {
		t.Fatalf("expected encrypted api key in config, got:\n%s", body)
	}
}

func TestConfigSetupCanReplaceExistingAgeIdentityFile(t *testing.T) {
	identity, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatalf("generate identity: %v", err)
	}

	home := t.TempDir()
	keyFile := filepath.Join(home, ".tscli", "age.txt")
	if err := os.MkdirAll(filepath.Dir(keyFile), 0o755); err != nil {
		t.Fatalf("create key dir: %v", err)
	}
	originalBody := identity.String() + "\n"
	if err := os.WriteFile(keyFile, []byte(originalBody), 0o600); err != nil {
		t.Fatalf("write age key file: %v", err)
	}

	res := executeCLINoDefaultsWithInput(t, []string{"config", "setup"}, map[string]string{
		"HOME": home,
	}, "yes\n\nno\napi-key\nsandbox\n\ntskey-sandbox\nno\njson\nno\n")
	if res.err != nil {
		t.Fatalf("config setup replace existing key: %v\nstderr:\n%s", res.err, res.stderr)
	}

	data, err := os.ReadFile(keyFile)
	if err != nil {
		t.Fatalf("read replaced key file: %v", err)
	}
	if string(data) == originalBody {
		t.Fatalf("expected key file to be replaced")
	}
	if !strings.Contains(string(data), "AGE-SECRET-KEY-") {
		t.Fatalf("expected generated key file contents, got:\n%s", string(data))
	}
}

func TestConfigSetupInvalidExistingAgeIdentityFallsBackToGeneration(t *testing.T) {
	home := t.TempDir()
	keyFile := filepath.Join(home, ".tscli", "age.txt")
	if err := os.MkdirAll(filepath.Dir(keyFile), 0o755); err != nil {
		t.Fatalf("create key dir: %v", err)
	}
	if err := os.WriteFile(keyFile, []byte("not-an-age-key\n"), 0o600); err != nil {
		t.Fatalf("write invalid key file: %v", err)
	}

	res := executeCLINoDefaultsWithInput(t, []string{"config", "setup"}, map[string]string{
		"HOME": home,
	}, "yes\n\napi-key\nsandbox\n\ntskey-sandbox\nno\njson\nno\n")
	if res.err != nil {
		t.Fatalf("config setup invalid existing key fallback: %v\nstderr:\n%s", res.err, res.stderr)
	}
	if !strings.Contains(res.stdout, "could not be reused") {
		t.Fatalf("expected invalid key fallback message, got:\n%s", res.stdout)
	}

	data, err := os.ReadFile(keyFile)
	if err != nil {
		t.Fatalf("read generated key file: %v", err)
	}
	if !strings.Contains(string(data), "AGE-SECRET-KEY-") {
		t.Fatalf("expected generated key file after fallback, got:\n%s", string(data))
	}
}

func TestConfigSetupRerunCanAddAndDeleteProfiles(t *testing.T) {
	home := t.TempDir()

	res := executeCLINoDefaultsWithInput(t, []string{"config", "setup"}, map[string]string{
		"HOME": home,
	}, "yes\n\napi-key\nsandbox\n\ntskey-sandbox\nyes\noauth\norg-admin\n\ncid\nsecret\nno\njson\nno\n")
	if res.err != nil {
		t.Fatalf("config setup initial run: %v\nstderr:\n%s", res.err, res.stderr)
	}

	res = executeCLINoDefaultsWithInput(t, []string{"config", "setup"}, map[string]string{
		"HOME": home,
	}, "delete\norg-admin\nquit\n")
	if res.err != nil {
		t.Fatalf("config setup rerun delete: %v\nstderr:\n%s", res.err, res.stderr)
	}
	if !strings.Contains(res.stdout, "Add, modify, or delete profiles?") {
		t.Fatalf("expected rerun management prompt, got:\n%s", res.stdout)
	}
	if strings.Contains(res.stdout, "Choose your default output format") || strings.Contains(res.stdout, "Enable debug HTTP request/response logging by default?") {
		t.Fatalf("did not expect initial setup preference prompts during rerun, got:\n%s", res.stdout)
	}

	configFile := filepath.Join(home, ".tscli.yaml")
	cfg, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatalf("read config file: %v", err)
	}
	body := string(cfg)
	if !strings.Contains(body, "api-key-encrypted:") {
		t.Fatalf("expected encrypted api key to remain, got:\n%s", body)
	}
	if !strings.Contains(body, "active-tailnet: sandbox") {
		t.Fatalf("expected sandbox to remain active, got:\n%s", body)
	}
	if strings.Contains(body, "name: org-admin") || strings.Contains(body, "oauth-client-id: cid") || strings.Contains(body, "oauth-client-secret-encrypted:") {
		t.Fatalf("expected deleted oauth profile to be removed, got:\n%s", body)
	}
}

func TestConfigSetupRerunCanModifySelectedProfile(t *testing.T) {
	home := t.TempDir()

	res := executeCLINoDefaultsWithInput(t, []string{"config", "setup"}, map[string]string{
		"HOME": home,
	}, "no\noauth\norg-admin\nexample.ts.net\ncid\nsecret\nno\njson\nno\n")
	if res.err != nil {
		t.Fatalf("config setup initial run: %v\nstderr:\n%s", res.err, res.stderr)
	}

	res = executeCLINoDefaultsWithInput(t, []string{"config", "setup"}, map[string]string{
		"HOME": home,
	}, "modify\norg-admin\n\n\n\nupdated-secret\nquit\n")
	if res.err != nil {
		t.Fatalf("config setup rerun modify: %v\nstderr:\n%s", res.err, res.stderr)
	}
	if !strings.Contains(res.stdout, "Select a profile to modify") {
		t.Fatalf("expected modify selection prompt, got:\n%s", res.stdout)
	}
	if strings.Contains(res.stdout, "Choose your default output format") || strings.Contains(res.stdout, "Enable debug HTTP request/response logging by default?") {
		t.Fatalf("did not expect initial setup preference prompts during rerun modify, got:\n%s", res.stdout)
	}

	configFile := filepath.Join(home, ".tscli.yaml")
	cfg, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatalf("read config file: %v", err)
	}
	body := string(cfg)
	if !strings.Contains(body, "name: org-admin") {
		t.Fatalf("expected updated profile to remain present, got:\n%s", body)
	}
	if !strings.Contains(body, "tailnet: example.ts.net") {
		t.Fatalf("expected existing tailnet override to be preserved, got:\n%s", body)
	}
	if !strings.Contains(body, "oauth-client-id: cid") {
		t.Fatalf("expected existing oauth client id to be preserved, got:\n%s", body)
	}
	if !strings.Contains(body, "oauth-client-secret: updated-secret") {
		t.Fatalf("expected oauth client secret to be updated, got:\n%s", body)
	}
	if strings.Contains(body, "oauth-client-secret: secret") {
		t.Fatalf("expected old oauth client secret to be replaced, got:\n%s", body)
	}
}

func TestConfigProfilesEncryptSecretsWhenEnabled(t *testing.T) {
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

	res = executeCLINoDefaults(t, []string{"config", "profiles", "set", "sandbox", "--api-key", "tskey-sandbox"}, map[string]string{
		"HOME": home,
	})
	if res.err != nil {
		t.Fatalf("upsert encrypted api-key profile: %v\nstderr:\n%s", res.err, res.stderr)
	}

	res = executeCLINoDefaults(t, []string{"config", "profiles", "set", "org-admin", "--oauth-client-id", "cid", "--oauth-client-secret", "secret"}, map[string]string{
		"HOME": home,
	})
	if res.err != nil {
		t.Fatalf("upsert encrypted oauth profile: %v\nstderr:\n%s", res.err, res.stderr)
	}

	configFile := filepath.Join(home, ".tscli.yaml")
	cfg, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatalf("read config file: %v", err)
	}
	body := string(cfg)
	if strings.Contains(body, "api-key: tskey-sandbox") {
		t.Fatalf("did not expect plaintext api-key in encrypted config:\n%s", body)
	}
	if strings.Contains(body, "oauth-client-secret: secret") {
		t.Fatalf("did not expect plaintext oauth-client-secret in encrypted config:\n%s", body)
	}
	if !strings.Contains(body, "api-key-encrypted:") || !strings.Contains(body, "oauth-client-secret-encrypted:") {
		t.Fatalf("expected encrypted sibling fields, got:\n%s", body)
	}

	res = executeCLINoDefaults(t, []string{"config", "profiles", "list"}, map[string]string{
		"HOME": home,
	})
	if res.err != nil {
		t.Fatalf("list encrypted profiles: %v\nstderr:\n%s", res.err, res.stderr)
	}
}

func TestConfigProfilesUpsertUsesProfileTailnetFlag(t *testing.T) {
	home := t.TempDir()

	res := executeCLINoDefaults(t, []string{"config", "profiles", "set", "sandbox", "--api-key", "tskey-sandbox", "--profile-tailnet", "example.ts.net"}, map[string]string{
		"HOME": home,
	})
	if res.err != nil {
		t.Fatalf("upsert profile tailnet: %v\nstderr:\n%s", res.err, res.stderr)
	}

	configFile := filepath.Join(home, ".tscli.yaml")
	cfg, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatalf("read config file: %v", err)
	}
	if !strings.Contains(string(cfg), "tailnet: example.ts.net") {
		t.Fatalf("expected persisted profile tailnet, got:\n%s", string(cfg))
	}
	if strings.Contains(string(cfg), "\nprofile-tailnet:") {
		t.Fatalf("did not expect profile-tailnet key to persist, got:\n%s", string(cfg))
	}
}

func TestConfigProfilesUpsertAcceptsDeprecatedTailnetAlias(t *testing.T) {
	home := t.TempDir()

	res := executeCLINoDefaults(t, []string{"config", "profiles", "set", "sandbox", "--api-key", "tskey-sandbox", "--tailnet", "example.ts.net"}, map[string]string{
		"HOME": home,
	})
	if res.err != nil {
		t.Fatalf("upsert profile tailnet alias: %v\nstderr:\n%s", res.err, res.stderr)
	}

	configFile := filepath.Join(home, ".tscli.yaml")
	cfg, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatalf("read config file: %v", err)
	}
	if !strings.Contains(string(cfg), "tailnet: example.ts.net") {
		t.Fatalf("expected persisted profile tailnet from deprecated alias, got:\n%s", string(cfg))
	}
}

func TestConfigProfilesUpsertRejectsMixedAuthShapes(t *testing.T) {
	res := executeCLINoDefaults(t, []string{"config", "profiles", "set", "mixed", "--api-key", "tskey-mixed", "--oauth-client-id", "cid", "--oauth-client-secret", "secret"}, nil)
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
		{"config", "profiles", "set", "sandbox", "--api-key", "tskey-sandbox"},
		{"config", "profiles", "set", "prod", "--api-key", "tskey-prod"},
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
	for _, unwanted := range []string{"tailnet", "api-key", "help"} {
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
	if got, ok := shown["debug"].(bool); !ok || got {
		t.Fatalf("expected persisted debug false in output, got %s", res.stdout)
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
