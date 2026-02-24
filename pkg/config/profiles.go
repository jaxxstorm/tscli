package config

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/spf13/viper"
)

// TailnetProfile represents one named tailnet credential profile.
type TailnetProfile struct {
	Name   string `mapstructure:"name" json:"name" yaml:"name"`
	APIKey string `mapstructure:"api-key" json:"api-key" yaml:"api-key"`
}

// TailnetProfilesState represents profile-backed configuration.
type TailnetProfilesState struct {
	ActiveTailnet string           `mapstructure:"active-tailnet" json:"active-tailnet" yaml:"active-tailnet"`
	Tailnets      []TailnetProfile `mapstructure:"tailnets" json:"tailnets" yaml:"tailnets"`
}

// ResolvedRuntimeConfig is the effective auth context after applying precedence.
type ResolvedRuntimeConfig struct {
	APIKey  string
	Tailnet string
}

func ListTailnetProfiles() (TailnetProfilesState, error) {
	return loadTailnetProfilesState(viper.GetViper())
}

func SetActiveTailnet(name string) error {
	v := viper.GetViper()

	state, err := loadTailnetProfilesState(v)
	if err != nil {
		return err
	}

	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("tailnet profile name is required")
	}

	if _, found := findTailnetProfile(state.Tailnets, name); !found {
		return fmt.Errorf("tailnet profile %q not found", name)
	}

	state.ActiveTailnet = name
	return persistTailnetProfilesState(v, state)
}

func UpsertTailnetProfile(name, apiKey string) (bool, error) {
	v := viper.GetViper()

	state, err := loadTailnetProfilesState(v)
	if err != nil {
		return false, err
	}

	name = strings.TrimSpace(name)
	apiKey = strings.TrimSpace(apiKey)

	if name == "" {
		return false, fmt.Errorf("tailnet profile name is required")
	}
	if apiKey == "" {
		return false, fmt.Errorf("tailnet profile api-key is required")
	}

	created := true
	for i := range state.Tailnets {
		if state.Tailnets[i].Name == name {
			state.Tailnets[i].APIKey = apiKey
			created = false
			break
		}
	}
	if created {
		state.Tailnets = append(state.Tailnets, TailnetProfile{
			Name:   name,
			APIKey: apiKey,
		})
	}

	if strings.TrimSpace(state.ActiveTailnet) == "" {
		state.ActiveTailnet = name
	}

	if err := persistTailnetProfilesState(v, state); err != nil {
		return false, err
	}

	return created, nil
}

func RemoveTailnetProfile(name string) error {
	v := viper.GetViper()

	state, err := loadTailnetProfilesState(v)
	if err != nil {
		return err
	}

	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("tailnet profile name is required")
	}
	if state.ActiveTailnet == name {
		return fmt.Errorf("cannot remove active tailnet profile %q; set a different active profile first", name)
	}

	index := -1
	for i, profile := range state.Tailnets {
		if profile.Name == name {
			index = i
			break
		}
	}
	if index < 0 {
		return fmt.Errorf("tailnet profile %q not found", name)
	}

	state.Tailnets = append(state.Tailnets[:index], state.Tailnets[index+1:]...)
	return persistTailnetProfilesState(v, state)
}

func ResolveRuntimeConfig(flagOverrides map[string]struct{}) (ResolvedRuntimeConfig, error) {
	return resolveRuntimeConfig(viper.GetViper(), flagOverrides)
}

func resolveRuntimeConfig(v *viper.Viper, flagOverrides map[string]struct{}) (ResolvedRuntimeConfig, error) {
	state, err := loadTailnetProfilesState(v)
	if err != nil {
		return ResolvedRuntimeConfig{}, err
	}

	legacyAPIKey := strings.TrimSpace(v.GetString("api-key"))
	legacyTailnet := strings.TrimSpace(v.GetString("tailnet"))

	apiKey, apiKeySet := resolveWithPrecedence(
		v,
		"api-key",
		"TAILSCALE_API_KEY",
		legacyAPIKey,
		state,
		flagOverrides,
	)
	tailnet, tailnetSet := resolveWithPrecedence(
		v,
		"tailnet",
		"TAILSCALE_TAILNET",
		legacyTailnet,
		state,
		flagOverrides,
	)

	if !tailnetSet || strings.TrimSpace(tailnet) == "" {
		tailnet = "-"
	}

	if !apiKeySet || strings.TrimSpace(apiKey) == "" {
		return ResolvedRuntimeConfig{}, fmt.Errorf("a Tailscale API key is required")
	}

	resolved := ResolvedRuntimeConfig{
		APIKey:  strings.TrimSpace(apiKey),
		Tailnet: strings.TrimSpace(tailnet),
	}

	v.Set("api-key", resolved.APIKey)
	v.Set("tailnet", resolved.Tailnet)

	return resolved, nil
}

func resolveWithPrecedence(
	v *viper.Viper,
	key string,
	envKey string,
	legacyValue string,
	state TailnetProfilesState,
	flagOverrides map[string]struct{},
) (string, bool) {
	if _, ok := flagOverrides[key]; ok {
		return strings.TrimSpace(v.GetString(key)), true
	}

	if envValue, ok := os.LookupEnv(envKey); ok {
		return strings.TrimSpace(envValue), true
	}

	if state.ActiveTailnet != "" {
		if profile, found := findTailnetProfile(state.Tailnets, state.ActiveTailnet); found {
			if key == "tailnet" {
				return profile.Name, true
			}
			if key == "api-key" {
				return profile.APIKey, true
			}
		}
	}

	if strings.TrimSpace(legacyValue) != "" {
		return strings.TrimSpace(legacyValue), true
	}

	return "", false
}

func loadTailnetProfilesState(v *viper.Viper) (TailnetProfilesState, error) {
	var profiles []TailnetProfile
	if err := v.UnmarshalKey("tailnets", &profiles); err != nil {
		return TailnetProfilesState{}, fmt.Errorf("decode tailnets: %w", err)
	}

	state := TailnetProfilesState{
		ActiveTailnet: strings.TrimSpace(v.GetString("active-tailnet")),
		Tailnets:      normalizeProfiles(profiles),
	}

	if err := validateTailnetProfilesState(state); err != nil {
		return TailnetProfilesState{}, err
	}

	return state, nil
}

func normalizeProfiles(profiles []TailnetProfile) []TailnetProfile {
	normalized := make([]TailnetProfile, 0, len(profiles))
	for _, profile := range profiles {
		normalized = append(normalized, TailnetProfile{
			Name:   strings.TrimSpace(profile.Name),
			APIKey: strings.TrimSpace(profile.APIKey),
		})
	}

	slices.SortFunc(normalized, func(a, b TailnetProfile) int {
		return strings.Compare(a.Name, b.Name)
	})

	return normalized
}

func validateTailnetProfilesState(state TailnetProfilesState) error {
	seen := map[string]struct{}{}

	for i, profile := range state.Tailnets {
		if profile.Name == "" {
			return fmt.Errorf("tailnet profile at index %d is missing name", i)
		}
		if profile.APIKey == "" {
			return fmt.Errorf("tailnet profile %q is missing api-key", profile.Name)
		}
		if _, exists := seen[profile.Name]; exists {
			return fmt.Errorf("duplicate tailnet profile name %q", profile.Name)
		}
		seen[profile.Name] = struct{}{}
	}

	if state.ActiveTailnet != "" && len(state.Tailnets) > 0 {
		if _, exists := seen[state.ActiveTailnet]; !exists {
			return fmt.Errorf("active-tailnet %q does not match any configured tailnet profile", state.ActiveTailnet)
		}
	}

	return nil
}

func findTailnetProfile(profiles []TailnetProfile, name string) (TailnetProfile, bool) {
	for _, profile := range profiles {
		if profile.Name == name {
			return profile, true
		}
	}
	return TailnetProfile{}, false
}

func persistTailnetProfilesState(v *viper.Viper, state TailnetProfilesState) error {
	state.ActiveTailnet = strings.TrimSpace(state.ActiveTailnet)
	state.Tailnets = normalizeProfiles(state.Tailnets)

	if err := validateTailnetProfilesState(state); err != nil {
		return err
	}

	v.Set("tailnets", state.Tailnets)
	v.Set("active-tailnet", state.ActiveTailnet)

	if state.ActiveTailnet != "" {
		active, _ := findTailnetProfile(state.Tailnets, state.ActiveTailnet)
		v.Set("tailnet", active.Name)
		v.Set("api-key", active.APIKey)
	}

	if err := save(v); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	return nil
}
