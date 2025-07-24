package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	fileName = ".tscli"
	fileType = "yaml"
)

// TailnetConfig represents a single tailnet configuration
type TailnetConfig struct {
	Name   string `yaml:"name" json:"name"`
	APIKey string `yaml:"api-key" json:"api-key"`
}

// Config represents the overall configuration structure
type Config struct {
	// Legacy fields for backward compatibility
	APIKey string `yaml:"api-key,omitempty" json:"api-key,omitempty"`
	Debug  bool   `yaml:"debug" json:"debug"`
	Output string `yaml:"output" json:"output"`
	Help   bool   `yaml:"help" json:"help"`

	// New multi-tailnet support
	Tailnets      []TailnetConfig `yaml:"tailnets,omitempty" json:"tailnets,omitempty"`
	ActiveTailnet string          `yaml:"active-tailnet,omitempty" json:"active-tailnet,omitempty"`
}

func Init() {
	v := viper.GetViper()

	// Search order: cwd ⇒ $HOME
	v.SetConfigName(fileName)
	v.SetConfigType(fileType)
	v.AddConfigPath(".")
	if home, _ := os.UserHomeDir(); home != "" {
		v.AddConfigPath(home)
	}

	_ = v.ReadInConfig() // ignore "not found"
	v.SetDefault("output", "json")

	// NO automatic migration - keep both old and new logic working
}

func Save() error {
	v := viper.GetViper()
	path := v.ConfigFileUsed()
	if path == "" {
		home, _ := os.UserHomeDir()
		path = filepath.Join(home, fileName+"."+fileType)
	}
	return v.WriteConfigAs(path)
}

// GetActiveTailnetConfig returns the configuration for the currently active tailnet
// Returns nil if using legacy configuration
func GetActiveTailnetConfig() (*TailnetConfig, error) {
	v := viper.GetViper()

	// Check if we're using the new multi-tailnet configuration
	if v.IsSet("tailnets") && len(getTailnets()) > 0 {
		activeName := v.GetString("active-tailnet")
		if activeName == "" {
			return nil, fmt.Errorf("no active tailnet configured")
		}

		tailnets := getTailnets()
		for _, tailnet := range tailnets {
			if tailnet.Name == activeName {
				return &tailnet, nil
			}
		}

		return nil, fmt.Errorf("active tailnet %q not found", activeName)
	}

	// Using legacy configuration - return nil to indicate legacy mode
	return nil, nil
}

// IsLegacyConfig returns true if using the old single api-key configuration
func IsLegacyConfig() bool {
	v := viper.GetViper()
	// Legacy if we have api-key set but no tailnets configured
	return v.IsSet("api-key") && v.GetString("api-key") != "" && !v.IsSet("tailnets")
}

// IsNewConfig returns true if using the new multi-tailnet configuration
func IsNewConfig() bool {
	v := viper.GetViper()
	return v.IsSet("tailnets") && len(getTailnets()) > 0
}

// getTailnets returns all configured tailnets
func getTailnets() []TailnetConfig {
	v := viper.GetViper()

	// Get the raw tailnets data the same way config show does
	allSettings := v.AllSettings()
	tailnetsData, exists := allSettings["tailnets"]
	if !exists {
		return []TailnetConfig{}
	}

	var tailnets []TailnetConfig

	// Convert the interface{} data to our struct
	if slice, ok := tailnetsData.([]interface{}); ok {
		for _, item := range slice {
			if m, ok := item.(map[string]interface{}); ok {
				name, _ := m["name"].(string)
				apiKey, _ := m["api-key"].(string)
				if name != "" {
					tailnets = append(tailnets, TailnetConfig{
						Name:   name,
						APIKey: apiKey,
					})
				}
			}
		}
	}

	return tailnets
}

// AddTailnet adds a new tailnet configuration
func AddTailnet(name, apiKey string) error {
	v := viper.GetViper()

	tailnets := getTailnets()

	// Check if tailnet already exists
	for _, tailnet := range tailnets {
		if tailnet.Name == name {
			return fmt.Errorf("tailnet %q already exists", name)
		}
	}

	// Add new tailnet
	newTailnet := TailnetConfig{
		Name:   name,
		APIKey: apiKey,
	}

	tailnets = append(tailnets, newTailnet)
	v.Set("tailnets", tailnets)

	// If this is the first tailnet, make it active
	if len(tailnets) == 1 {
		v.Set("active-tailnet", name)
	}

	return Save()
}

// RemoveTailnet removes a tailnet configuration
func RemoveTailnet(name string) error {
	v := viper.GetViper()

	tailnets := getTailnets()
	var updatedTailnets []TailnetConfig
	found := false

	for _, tailnet := range tailnets {
		if tailnet.Name != name {
			updatedTailnets = append(updatedTailnets, tailnet)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("tailnet %q not found", name)
	}

	v.Set("tailnets", updatedTailnets)

	// If we removed the active tailnet, clear the active setting
	if v.GetString("active-tailnet") == name {
		if len(updatedTailnets) > 0 {
			v.Set("active-tailnet", updatedTailnets[0].Name)
		} else {
			v.Set("active-tailnet", "")
		}
	}

	return Save()
}

// SetActiveTailnet switches to a different tailnet
func SetActiveTailnet(name string) error {
	v := viper.GetViper()

	tailnets := getTailnets()
	found := false

	for _, tailnet := range tailnets {
		if tailnet.Name == name {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("tailnet %q not found", name)
	}

	v.Set("active-tailnet", name)
	return Save()
}

// ListTailnets returns all configured tailnets with indication of which is active
func ListTailnets() ([]TailnetConfig, string, error) {
	tailnets := getTailnets()
	active := viper.GetString("active-tailnet")
	return tailnets, active, nil
}
