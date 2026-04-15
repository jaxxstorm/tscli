package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"filippo.io/age"
	"github.com/spf13/viper"
)

func TestValidateAgeEncryptionConfig(t *testing.T) {
	identity, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatalf("generate identity: %v", err)
	}

	t.Run("requires public key when enabled", func(t *testing.T) {
		err := validateAgeEncryptionConfig(AgeEncryptionConfig{PrivateKey: identity.String()})
		if err == nil || !strings.Contains(err.Error(), "public-key") {
			t.Fatalf("expected public-key validation error, got %v", err)
		}
	})

	t.Run("rejects conflicting private key sources", func(t *testing.T) {
		err := validateAgeEncryptionConfig(AgeEncryptionConfig{
			PublicKey:         identity.Recipient().String(),
			PrivateKeyPath:    "/tmp/key.txt",
			PrivateKey:        identity.String(),
			PrivateKeyCommand: "op read secret",
		})
		if err == nil || !strings.Contains(err.Error(), "configure only one") {
			t.Fatalf("expected conflicting private key source error, got %v", err)
		}
	})

	t.Run("accepts env only configuration", func(t *testing.T) {
		err := validateAgeEncryptionConfig(AgeEncryptionConfig{PublicKey: identity.Recipient().String()})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestEncryptAndDecryptSecret(t *testing.T) {
	identity, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatalf("generate identity: %v", err)
	}

	v := viper.New()
	v.Set("encryption.age.public-key", identity.Recipient().String())
	privateKeyPath := filepath.Join(t.TempDir(), "age.txt")
	if err := os.WriteFile(privateKeyPath, []byte(identity.String()+"\n"), 0o600); err != nil {
		t.Fatalf("write private key file: %v", err)
	}
	v.Set("encryption.age.private-key-path", privateKeyPath)

	ciphertext, err := encryptSecret(v, "super-secret")
	if err != nil {
		t.Fatalf("encrypt secret: %v", err)
	}
	if strings.Contains(ciphertext, "super-secret") {
		t.Fatalf("expected ciphertext to omit plaintext, got %q", ciphertext)
	}

	plaintext, err := decryptSecret(v, ciphertext)
	if err != nil {
		t.Fatalf("decrypt secret: %v", err)
	}
	if plaintext != "super-secret" {
		t.Fatalf("expected decrypted plaintext, got %q", plaintext)
	}
}

func TestResolveRuntimeConfigDecryptsEncryptedProfileSecrets(t *testing.T) {
	identity, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatalf("generate identity: %v", err)
	}

	v := viper.New()
	v.Set("encryption.age.public-key", identity.Recipient().String())
	privateKeyPath := filepath.Join(t.TempDir(), "age.txt")
	if err := os.WriteFile(privateKeyPath, []byte(identity.String()+"\n"), 0o600); err != nil {
		t.Fatalf("write private key file: %v", err)
	}
	v.Set("encryption.age.private-key-path", privateKeyPath)

	ciphertext, err := encryptSecret(v, "tskey-encrypted")
	if err != nil {
		t.Fatalf("encrypt api key: %v", err)
	}

	v.Set("active-tailnet", "sandbox")
	v.Set("tailnets", []map[string]any{{
		"name":              "sandbox",
		"api-key-encrypted": ciphertext,
	}})

	resolved, err := resolveRuntimeConfig(v, nil)
	if err != nil {
		t.Fatalf("resolve runtime config: %v", err)
	}
	if resolved.APIKey != "tskey-encrypted" {
		t.Fatalf("expected decrypted api key, got %q", resolved.APIKey)
	}
}

func TestResolveAgeIdentityFromPrivateKeyPath(t *testing.T) {
	identity, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatalf("generate identity: %v", err)
	}

	privateKeyPath := filepath.Join(t.TempDir(), "keys.txt")
	if err := os.WriteFile(privateKeyPath, []byte("# created by test\n"+identity.String()+"\n"), 0o600); err != nil {
		t.Fatalf("write private key file: %v", err)
	}

	v := viper.New()
	v.Set("encryption.age.public-key", identity.Recipient().String())
	v.Set("encryption.age.private-key-path", privateKeyPath)

	resolved, err := resolveAgeIdentity(v)
	if err != nil {
		t.Fatalf("resolve age identity: %v", err)
	}
	if resolved == nil {
		t.Fatalf("expected resolved identity")
	}
}

func TestSetAgeEncryptionConfigPersistsSettings(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)

	identity, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatalf("generate identity: %v", err)
	}

	home := t.TempDir()
	t.Setenv("HOME", home)

	err = SetAgeEncryptionConfig(AgeEncryptionConfig{
		PublicKey: identity.Recipient().String(),
	})
	if err != nil {
		t.Fatalf("set age encryption config: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(home, ".tscli.yaml"))
	if err != nil {
		t.Fatalf("read config file: %v", err)
	}
	body := string(data)
	if !strings.Contains(body, "encryption:") || !strings.Contains(body, "public-key:") {
		t.Fatalf("expected encryption settings in config, got:\n%s", body)
	}
}

func TestSetAgeEncryptionConfigExpandsPrivateKeyPath(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)

	identity, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatalf("generate identity: %v", err)
	}

	home := t.TempDir()
	t.Setenv("HOME", home)

	privateKeyDir := filepath.Join(home, ".config", "age")
	if err := os.MkdirAll(privateKeyDir, 0o755); err != nil {
		t.Fatalf("create private key dir: %v", err)
	}
	privateKeyPath := filepath.Join(privateKeyDir, "key.txt")
	if err := os.WriteFile(privateKeyPath, []byte(identity.String()+"\n"), 0o600); err != nil {
		t.Fatalf("write private key file: %v", err)
	}

	err = SetAgeEncryptionConfig(AgeEncryptionConfig{
		PublicKey:      identity.Recipient().String(),
		PrivateKeyPath: "~/.config/age/key.txt",
	})
	if err != nil {
		t.Fatalf("set age encryption config: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(home, ".tscli.yaml"))
	if err != nil {
		t.Fatalf("read config file: %v", err)
	}
	if !strings.Contains(string(data), "private-key-path: "+privateKeyPath) {
		t.Fatalf("expected expanded private key path, got:\n%s", string(data))
	}
	if viper.GetString("encryption.age.private-key-path") != privateKeyPath {
		t.Fatalf("expected viper private key path %q, got %q", privateKeyPath, viper.GetString("encryption.age.private-key-path"))
	}
}

func TestInspectAgeIdentityFile(t *testing.T) {
	t.Run("loads reusable identity and derives public key", func(t *testing.T) {
		identity, err := age.GenerateX25519Identity()
		if err != nil {
			t.Fatalf("generate identity: %v", err)
		}

		path := filepath.Join(t.TempDir(), "age.txt")
		if err := os.WriteFile(path, []byte("# generated by test\n"+identity.String()+"\n"), 0o600); err != nil {
			t.Fatalf("write age identity file: %v", err)
		}

		inspected, err := InspectAgeIdentityFile(path)
		if err != nil {
			t.Fatalf("inspect age identity file: %v", err)
		}
		if inspected.Path != path {
			t.Fatalf("expected path %q, got %q", path, inspected.Path)
		}
		if inspected.PublicKey != identity.Recipient().String() {
			t.Fatalf("expected public key %q, got %q", identity.Recipient().String(), inspected.PublicKey)
		}
		if inspected.Identity == nil {
			t.Fatalf("expected parsed identity")
		}
	})

	t.Run("expands home path", func(t *testing.T) {
		identity, err := age.GenerateX25519Identity()
		if err != nil {
			t.Fatalf("generate identity: %v", err)
		}

		home := t.TempDir()
		t.Setenv("HOME", home)
		if err := os.MkdirAll(filepath.Join(home, ".tscli"), 0o755); err != nil {
			t.Fatalf("create home dir: %v", err)
		}
		if err := os.WriteFile(filepath.Join(home, ".tscli", "age.txt"), []byte(identity.String()+"\n"), 0o600); err != nil {
			t.Fatalf("write age identity file: %v", err)
		}

		inspected, err := InspectAgeIdentityFile("~/.tscli/age.txt")
		if err != nil {
			t.Fatalf("inspect home age identity file: %v", err)
		}
		if want := filepath.Join(home, ".tscli", "age.txt"); inspected.Path != want {
			t.Fatalf("expected path %q, got %q", want, inspected.Path)
		}
	})

	t.Run("rejects invalid identity file", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "invalid.txt")
		if err := os.WriteFile(path, []byte("not-an-age-key\n"), 0o600); err != nil {
			t.Fatalf("write invalid key file: %v", err)
		}

		if _, err := InspectAgeIdentityFile(path); err == nil {
			t.Fatalf("expected invalid identity error")
		}
	})

	t.Run("rejects empty file", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "empty.txt")
		if err := os.WriteFile(path, nil, 0o600); err != nil {
			t.Fatalf("write empty key file: %v", err)
		}

		if _, err := InspectAgeIdentityFile(path); err == nil || !strings.Contains(err.Error(), "empty") {
			t.Fatalf("expected empty file error, got %v", err)
		}
	})
}
