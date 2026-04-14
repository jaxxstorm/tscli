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
	Name              string `mapstructure:"name" json:"name" yaml:"name"`
	Tailnet           string `mapstructure:"tailnet" json:"tailnet,omitempty" yaml:"tailnet,omitempty"`
	APIKey            string `mapstructure:"api-key" json:"api-key,omitempty" yaml:"api-key,omitempty"`
	OAuthClientID     string `mapstructure:"oauth-client-id" json:"oauth-client-id,omitempty" yaml:"oauth-client-id,omitempty"`
	OAuthClientSecret string `mapstructure:"oauth-client-secret" json:"oauth-client-secret,omitempty" yaml:"oauth-client-secret,omitempty"`
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

type ResolvedOAuthConfig struct {
	ClientID     string
	ClientSecret string
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

func UpsertTailnetProfile(profile TailnetProfile) (bool, error) {
	v := viper.GetViper()

	state, err := loadTailnetProfilesState(v)
	if err != nil {
		return false, err
	}

	profile = normalizeProfile(profile)

	if profile.Name == "" {
		return false, fmt.Errorf("tailnet profile name is required")
	}
	if err := validateProfileAuthShape(profile); err != nil {
		return false, err
	}

	created := true
	for i := range state.Tailnets {
		if state.Tailnets[i].Name == profile.Name {
			state.Tailnets[i] = profile
			created = false
			break
		}
	}
	if created {
		state.Tailnets = append(state.Tailnets, profile)
	}

	if strings.TrimSpace(state.ActiveTailnet) == "" {
		state.ActiveTailnet = profile.Name
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

func ResolveOAuthRuntimeConfig(flagOverrides map[string]struct{}) (ResolvedOAuthConfig, error) {
	return resolveOAuthRuntimeConfig(viper.GetViper(), flagOverrides)
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
		func(profile TailnetProfile) string { return profile.APIKey },
	)
	tailnet, tailnetSet := resolveWithPrecedence(
		v,
		"tailnet",
		"TAILSCALE_TAILNET",
		legacyTailnet,
		state,
		flagOverrides,
		func(profile TailnetProfile) string { return profile.EffectiveTailnet() },
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

func resolveOAuthRuntimeConfig(v *viper.Viper, flagOverrides map[string]struct{}) (ResolvedOAuthConfig, error) {
	state, err := loadTailnetProfilesState(v)
	if err != nil {
		return ResolvedOAuthConfig{}, err
	}

	clientID, clientIDSet := resolveWithPrecedence(
		v,
		"oauth-client-id",
		"TSCLI_OAUTH_CLIENT_ID",
		strings.TrimSpace(v.GetString("oauth-client-id")),
		state,
		flagOverrides,
		func(profile TailnetProfile) string { return profile.OAuthClientID },
	)
	clientSecret, clientSecretSet := resolveWithPrecedence(
		v,
		"oauth-client-secret",
		"TSCLI_OAUTH_CLIENT_SECRET",
		strings.TrimSpace(v.GetString("oauth-client-secret")),
		state,
		flagOverrides,
		func(profile TailnetProfile) string { return profile.OAuthClientSecret },
	)

	if !clientIDSet || clientID == "" || !clientSecretSet || clientSecret == "" {
		return ResolvedOAuthConfig{}, fmt.Errorf("OAuth client credentials are required; provide --oauth-client-id and --oauth-client-secret, set TSCLI_OAUTH_CLIENT_ID and TSCLI_OAUTH_CLIENT_SECRET, or configure them on the active profile")
	}

	resolved := ResolvedOAuthConfig{
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}

	v.Set("oauth-client-id", resolved.ClientID)
	v.Set("oauth-client-secret", resolved.ClientSecret)

	return resolved, nil
}

func resolveWithPrecedence(
	v *viper.Viper,
	key string,
	envKey string,
	legacyValue string,
	state TailnetProfilesState,
	flagOverrides map[string]struct{},
	profileValue func(TailnetProfile) string,
) (string, bool) {
	if _, ok := flagOverrides[key]; ok {
		value := strings.TrimSpace(v.GetString(key))
		return value, value != ""
	}

	if envValue, ok := os.LookupEnv(envKey); ok {
		envValue = strings.TrimSpace(envValue)
		return envValue, envValue != ""
	}

	if state.ActiveTailnet != "" {
		if profile, found := findTailnetProfile(state.Tailnets, state.ActiveTailnet); found {
			if profileValue != nil {
				value := strings.TrimSpace(profileValue(profile))
				if value != "" {
					return value, true
				}
			}
		}
	}

	legacyValue = strings.TrimSpace(legacyValue)
	if legacyValue != "" {
		return legacyValue, true
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
		normalized = append(normalized, normalizeProfile(profile))
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
		if err := validateProfileAuthShape(profile); err != nil {
			return err
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

func normalizeProfile(profile TailnetProfile) TailnetProfile {
	return TailnetProfile{
		Name:              strings.TrimSpace(profile.Name),
		Tailnet:           strings.TrimSpace(profile.Tailnet),
		APIKey:            strings.TrimSpace(profile.APIKey),
		OAuthClientID:     strings.TrimSpace(profile.OAuthClientID),
		OAuthClientSecret: strings.TrimSpace(profile.OAuthClientSecret),
	}
}

func validateProfileAuthShape(profile TailnetProfile) error {
	hasAPIKey := profile.APIKey != ""
	hasOAuthID := profile.OAuthClientID != ""
	hasOAuthSecret := profile.OAuthClientSecret != ""
	hasOAuth := hasOAuthID || hasOAuthSecret

	switch {
	case hasAPIKey && hasOAuth:
		return fmt.Errorf("tailnet profile %q must use either api-key auth or oauth-client-id/oauth-client-secret auth, not both", profile.Name)
	case hasAPIKey:
		return nil
	case hasOAuthID && hasOAuthSecret:
		return nil
	case hasOAuth:
		return fmt.Errorf("tailnet profile %q must include both oauth-client-id and oauth-client-secret", profile.Name)
	default:
		return fmt.Errorf("tailnet profile %q must include api-key or both oauth-client-id and oauth-client-secret", profile.Name)
	}
}

func (p TailnetProfile) EffectiveTailnet() string {
	if p.Tailnet != "" {
		return p.Tailnet
	}
	return p.Name
}

func (p TailnetProfile) AuthType() string {
	if p.APIKey != "" {
		return "api-key"
	}
	if p.OAuthClientID != "" || p.OAuthClientSecret != "" {
		return "oauth"
	}
	return "unknown"
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

	settings, err := loadPersistedSettings(v)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	settings["tailnets"] = state.Tailnets
	settings["active-tailnet"] = state.ActiveTailnet
	delete(settings, "tailnet")
	delete(settings, "api-key")
	delete(settings, "oauth-client-id")
	delete(settings, "oauth-client-secret")

	v.Set("tailnets", state.Tailnets)
	v.Set("active-tailnet", state.ActiveTailnet)

	if err := writePersistedSettings(v, settings); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	return nil
}
