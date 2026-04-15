package config

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"filippo.io/age"
	"github.com/spf13/viper"
)

type AgeEncryptionConfig struct {
	PublicKey         string `mapstructure:"public-key" json:"public-key,omitempty" yaml:"public-key,omitempty"`
	PrivateKeyPath    string `mapstructure:"private-key-path" json:"private-key-path,omitempty" yaml:"private-key-path,omitempty"`
	PrivateKey        string `mapstructure:"private-key" json:"private-key,omitempty" yaml:"private-key,omitempty"`
	PrivateKeyCommand string `mapstructure:"private-key-command" json:"private-key-command,omitempty" yaml:"private-key-command,omitempty"`
}

func loadAgeEncryptionConfig(v *viper.Viper) AgeEncryptionConfig {
	return AgeEncryptionConfig{
		PublicKey:         strings.TrimSpace(v.GetString("encryption.age.public-key")),
		PrivateKeyPath:    strings.TrimSpace(v.GetString("encryption.age.private-key-path")),
		PrivateKey:        strings.TrimSpace(v.GetString("encryption.age.private-key")),
		PrivateKeyCommand: strings.TrimSpace(v.GetString("encryption.age.private-key-command")),
	}
}

func validateAgeEncryptionConfig(cfg AgeEncryptionConfig) error {
	hasPublicKey := cfg.PublicKey != ""
	hasPath := cfg.PrivateKeyPath != ""
	hasConfigPrivateKey := cfg.PrivateKey != ""
	hasCommand := cfg.PrivateKeyCommand != ""
	if !hasPublicKey && !hasPath && !hasConfigPrivateKey && !hasCommand {
		return nil
	}
	if !hasPublicKey {
		return fmt.Errorf("encryption.age.public-key is required when config encryption is enabled")
	}
	privateKeySources := 0
	if hasPath {
		privateKeySources++
	}
	if hasConfigPrivateKey {
		privateKeySources++
	}
	if hasCommand {
		privateKeySources++
	}
	if privateKeySources > 1 {
		return fmt.Errorf("configure only one of encryption.age.private-key-path, encryption.age.private-key-command, or encryption.age.private-key")
	}
	if _, err := age.ParseX25519Recipient(cfg.PublicKey); err != nil {
		return fmt.Errorf("invalid encryption.age.public-key: %w", err)
	}
	if hasPath {
		if _, err := os.Stat(cfg.PrivateKeyPath); err != nil {
			return fmt.Errorf("invalid encryption.age.private-key-path: %w", err)
		}
	}
	if hasConfigPrivateKey {
		if _, err := age.ParseX25519Identity(cfg.PrivateKey); err != nil {
			return fmt.Errorf("invalid encryption.age.private-key: %w", err)
		}
	}
	return nil
}

func encryptionEnabled(v *viper.Viper) bool {
	return loadAgeEncryptionConfig(v).PublicKey != ""
}

func encryptProfilesForPersistence(v *viper.Viper, profiles []TailnetProfile) ([]TailnetProfile, error) {
	if !encryptionEnabled(v) {
		return profiles, nil
	}

	out := make([]TailnetProfile, 0, len(profiles))
	for _, profile := range profiles {
		copy := profile
		if copy.APIKey != "" {
			ciphertext, err := encryptSecret(v, copy.APIKey)
			if err != nil {
				return nil, fmt.Errorf("encrypt api key for profile %q: %w", copy.Name, err)
			}
			copy.APIKeyEncrypted = ciphertext
			copy.APIKey = ""
		}
		if copy.OAuthClientSecret != "" {
			ciphertext, err := encryptSecret(v, copy.OAuthClientSecret)
			if err != nil {
				return nil, fmt.Errorf("encrypt oauth client secret for profile %q: %w", copy.Name, err)
			}
			copy.OAuthClientSecretEncrypted = ciphertext
			copy.OAuthClientSecret = ""
		}
		out = append(out, copy)
	}
	return out, nil
}

func encryptSecret(v *viper.Viper, plaintext string) (string, error) {
	cfg := loadAgeEncryptionConfig(v)
	if err := validateAgeEncryptionConfig(cfg); err != nil {
		return "", err
	}
	recipient, err := age.ParseX25519Recipient(cfg.PublicKey)
	if err != nil {
		return "", fmt.Errorf("parse age public key: %w", err)
	}

	var buf bytes.Buffer
	w, err := age.Encrypt(&buf, recipient)
	if err != nil {
		return "", fmt.Errorf("create age writer: %w", err)
	}
	if _, err := io.WriteString(w, plaintext); err != nil {
		return "", fmt.Errorf("write plaintext: %w", err)
	}
	if err := w.Close(); err != nil {
		return "", fmt.Errorf("close age writer: %w", err)
	}
	return strings.TrimSpace(buf.String()), nil
}

func decryptSecret(v *viper.Viper, ciphertext string) (string, error) {
	identity, err := resolveAgeIdentity(v)
	if err != nil {
		return "", err
	}
	r, err := age.Decrypt(strings.NewReader(ciphertext), identity)
	if err != nil {
		return "", fmt.Errorf("decrypt secret: %w", err)
	}
	plaintext, err := io.ReadAll(r)
	if err != nil {
		return "", fmt.Errorf("read decrypted secret: %w", err)
	}
	return strings.TrimSpace(string(plaintext)), nil
}

func resolveAgeIdentity(v *viper.Viper) (age.Identity, error) {
	if value, ok := os.LookupEnv("TSCLI_AGE_PRIVATE_KEY"); ok {
		value = strings.TrimSpace(value)
		if value == "" {
			return nil, fmt.Errorf("TSCLI_AGE_PRIVATE_KEY is set but empty")
		}
		identity, err := parseAgeIdentity(value)
		if err != nil {
			return nil, fmt.Errorf("invalid TSCLI_AGE_PRIVATE_KEY: %w", err)
		}
		return identity, nil
	}

	cfg := loadAgeEncryptionConfig(v)
	if err := validateAgeEncryptionConfig(cfg); err != nil {
		return nil, err
	}
	if cfg.PrivateKeyPath != "" {
		data, err := os.ReadFile(cfg.PrivateKeyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read encryption.age.private-key-path: %w", err)
		}
		value := strings.TrimSpace(string(data))
		if value == "" {
			return nil, fmt.Errorf("encryption.age.private-key-path points to an empty key file")
		}
		identity, err := parseAgeIdentity(value)
		if err != nil {
			return nil, fmt.Errorf("invalid key in encryption.age.private-key-path: %w", err)
		}
		return identity, nil
	}
	if cfg.PrivateKeyCommand != "" {
		out, err := exec.Command("sh", "-c", cfg.PrivateKeyCommand).Output()
		if err != nil {
			return nil, fmt.Errorf("failed to execute encryption.age.private-key-command: %w", err)
		}
		value := strings.TrimSpace(string(out))
		if value == "" {
			return nil, fmt.Errorf("encryption.age.private-key-command returned an empty key")
		}
		identity, err := parseAgeIdentity(value)
		if err != nil {
			return nil, fmt.Errorf("invalid private key returned by encryption.age.private-key-command: %w", err)
		}
		return identity, nil
	}
	if cfg.PrivateKey != "" {
		identity, err := age.ParseX25519Identity(cfg.PrivateKey)
		if err != nil {
			return nil, fmt.Errorf("invalid encryption.age.private-key: %w", err)
		}
		return identity, nil
	}

	return nil, fmt.Errorf("encrypted config requires TSCLI_AGE_PRIVATE_KEY, encryption.age.private-key-path, encryption.age.private-key-command, or encryption.age.private-key")
}

func parseAgeIdentity(value string) (age.Identity, error) {
	identities, err := age.ParseIdentities(strings.NewReader(strings.TrimSpace(value)))
	if err == nil && len(identities) > 0 {
		return identities[0], nil
	}

	identity, parseErr := age.ParseX25519Identity(strings.TrimSpace(value))
	if parseErr == nil {
		return identity, nil
	}
	if err != nil {
		return nil, err
	}
	return nil, parseErr
}

func SetAgeEncryptionConfig(cfg AgeEncryptionConfig) error {
	v := viper.GetViper()
	cfg.PublicKey = strings.TrimSpace(cfg.PublicKey)
	cfg.PrivateKeyPath = strings.TrimSpace(cfg.PrivateKeyPath)
	cfg.PrivateKey = strings.TrimSpace(cfg.PrivateKey)
	cfg.PrivateKeyCommand = strings.TrimSpace(cfg.PrivateKeyCommand)
	if err := validateAgeEncryptionConfig(cfg); err != nil {
		return err
	}

	settings, err := loadPersistedSettings(v)
	if err != nil {
		return err
	}
	if cfg.PublicKey == "" && cfg.PrivateKeyPath == "" && cfg.PrivateKey == "" && cfg.PrivateKeyCommand == "" {
		delete(settings, "encryption")
	} else {
		settings["encryption"] = map[string]any{
			"age": map[string]any{
				"public-key":          cfg.PublicKey,
				"private-key-path":    cfg.PrivateKeyPath,
				"private-key":         cfg.PrivateKey,
				"private-key-command": cfg.PrivateKeyCommand,
			},
		}
	}
	v.Set("encryption.age.public-key", cfg.PublicKey)
	v.Set("encryption.age.private-key-path", cfg.PrivateKeyPath)
	v.Set("encryption.age.private-key", cfg.PrivateKey)
	v.Set("encryption.age.private-key-command", cfg.PrivateKeyCommand)
	return writePersistedSettings(v, settings)
}
