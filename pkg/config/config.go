package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

const (
	fileName = ".tscli"
	fileType = "yaml"
)

func Init() {
	v := viper.GetViper()

	// Search order: cwd ⇒ $HOME
	v.SetConfigName(fileName)
	v.SetConfigType(fileType)
	v.AddConfigPath(".")
	if home, _ := os.UserHomeDir(); home != "" {
		v.AddConfigPath(home)
	}

	_ = v.ReadInConfig() // ignore “not found”
	v.SetDefault("output", "json")
	v.SetDefault("tailnet", "-")
}

func Save() error {
	return save(viper.GetViper())
}

func save(v *viper.Viper) error {
	settings, err := loadPersistedSettings(v)
	if err != nil {
		return err
	}
	return writePersistedSettings(v, settings)
}

func SetPersistentValue(key, value string) error {
	v := viper.GetViper()

	settings, err := loadPersistedSettings(v)
	if err != nil {
		return err
	}
	settings[key] = value
	v.Set(key, value)

	return writePersistedSettings(v, settings)
}

func ShowSettings() (map[string]any, error) {
	return loadPersistedSettings(viper.GetViper())
}

func configFilePath(v *viper.Viper) string {
	path := v.ConfigFileUsed()
	if path == "" {
		home, _ := os.UserHomeDir()
		path = filepath.Join(home, fileName+"."+fileType)
	}
	return path
}

func loadPersistedSettings(v *viper.Viper) (map[string]any, error) {
	path := configFilePath(v)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]any{}, nil
		}
		return nil, err
	}

	if len(data) == 0 {
		return map[string]any{}, nil
	}

	var settings map[string]any
	if err := yaml.Unmarshal(data, &settings); err != nil {
		return nil, err
	}
	if settings == nil {
		settings = map[string]any{}
	}

	return canonicalizePersistedSettings(settings), nil
}

func writePersistedSettings(v *viper.Viper, settings map[string]any) error {
	settings = canonicalizePersistedSettings(settings)

	data, err := yaml.Marshal(settings)
	if err != nil {
		return err
	}

	path := configFilePath(v)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

func canonicalizePersistedSettings(settings map[string]any) map[string]any {
	canonical := make(map[string]any, len(settings))
	for key, value := range settings {
		switch key {
		case "help", "debug", "base-url":
			continue
		default:
			if isEmptyPersistedValue(value) {
				continue
			}
			canonical[key] = value
		}
	}

	if hasProfileState(canonical) {
		delete(canonical, "tailnet")
		delete(canonical, "api-key")
	}

	return canonical
}

func hasProfileState(settings map[string]any) bool {
	if active, ok := settings["active-tailnet"].(string); ok && active != "" {
		return true
	}
	tailnets, ok := settings["tailnets"]
	if !ok {
		return false
	}
	switch value := tailnets.(type) {
	case []any:
		return len(value) > 0
	case []map[string]any:
		return len(value) > 0
	case []TailnetProfile:
		return len(value) > 0
	default:
		return true
	}
}

func isEmptyPersistedValue(value any) bool {
	switch v := value.(type) {
	case nil:
		return true
	case string:
		return v == ""
	case []any:
		return len(v) == 0
	case []map[string]any:
		return len(v) == 0
	case []TailnetProfile:
		return len(v) == 0
	default:
		return false
	}
}
