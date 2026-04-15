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
			PrivateKey:        identity.String(),
			PrivateKeyCommand: "op read secret",
		})
		if err == nil || !strings.Contains(err.Error(), "not both") {
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
	v.Set("encryption.age.private-key", identity.String())

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
	v.Set("encryption.age.private-key", identity.String())

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
