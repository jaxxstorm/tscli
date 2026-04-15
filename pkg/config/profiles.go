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
	Name                       string `mapstructure:"name" json:"name" yaml:"name"`
	Tailnet                    string `mapstructure:"tailnet" json:"tailnet,omitempty" yaml:"tailnet,omitempty"`
	APIKey                     string `mapstructure:"api-key" json:"api-key,omitempty" yaml:"api-key,omitempty"`
	APIKeyEncrypted            string `mapstructure:"api-key-encrypted" json:"api-key-encrypted,omitempty" yaml:"api-key-encrypted,omitempty"`
	OAuthClientID              string `mapstructure:"oauth-client-id" json:"oauth-client-id,omitempty" yaml:"oauth-client-id,omitempty"`
	OAuthClientSecret          string `mapstructure:"oauth-client-secret" json:"oauth-client-secret,omitempty" yaml:"oauth-client-secret,omitempty"`
	OAuthClientSecretEncrypted string `mapstructure:"oauth-client-secret-encrypted" json:"oauth-client-secret-encrypted,omitempty" yaml:"oauth-client-secret-encrypted,omitempty"`
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

type ResolvedCommandAuth struct {
	Tailnet   string
	APIKey    string
	OAuth     ResolvedOAuthConfig
	UsesOAuth bool
	Source    string
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
	if state.ActiveTailnet == name {
		return fmt.Errorf("tailnet profile %q is active; set a different active profile before deleting it", name)
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

func ResolveCommandAuthConfig(flagOverrides map[string]struct{}) (ResolvedCommandAuth, error) {
	return resolveCommandAuthConfig(viper.GetViper(), flagOverrides)
}

func resolveRuntimeConfig(v *viper.Viper, flagOverrides map[string]struct{}) (ResolvedRuntimeConfig, error) {
	state, err := loadTailnetProfilesState(v)
	if err != nil {
		return ResolvedRuntimeConfig{}, err
	}

	legacyAPIKey := strings.TrimSpace(v.GetString("api-key"))
	legacyTailnet := strings.TrimSpace(v.GetString("tailnet"))

	apiKey, apiKeySet, err := resolveWithPrecedence(
		v,
		"api-key",
		"TAILSCALE_API_KEY",
		legacyAPIKey,
		state,
		flagOverrides,
		func(profile TailnetProfile) (string, error) {
			return profile.ResolveAPIKey(v)
		},
	)
	if err != nil {
		return ResolvedRuntimeConfig{}, err
	}
	tailnet, tailnetSet, err := resolveWithPrecedence(
		v,
		"tailnet",
		"TAILSCALE_TAILNET",
		legacyTailnet,
		state,
		flagOverrides,
		func(profile TailnetProfile) (string, error) { return profile.EffectiveTailnet(), nil },
	)
	if err != nil {
		return ResolvedRuntimeConfig{}, err
	}

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

	clientID, clientIDSet, err := resolveWithPrecedence(
		v,
		"oauth-client-id",
		"TSCLI_OAUTH_CLIENT_ID",
		strings.TrimSpace(v.GetString("oauth-client-id")),
		state,
		flagOverrides,
		func(profile TailnetProfile) (string, error) { return profile.OAuthClientID, nil },
	)
	if err != nil {
		return ResolvedOAuthConfig{}, err
	}
	clientSecret, clientSecretSet, err := resolveWithPrecedence(
		v,
		"oauth-client-secret",
		"TSCLI_OAUTH_CLIENT_SECRET",
		strings.TrimSpace(v.GetString("oauth-client-secret")),
		state,
		flagOverrides,
		func(profile TailnetProfile) (string, error) {
			return profile.ResolveOAuthClientSecret(v)
		},
	)
	if err != nil {
		return ResolvedOAuthConfig{}, err
	}

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

func resolveCommandAuthConfig(v *viper.Viper, flagOverrides map[string]struct{}) (ResolvedCommandAuth, error) {
	state, err := loadTailnetProfilesState(v)
	if err != nil {
		return ResolvedCommandAuth{}, err
	}

	legacyTailnet := strings.TrimSpace(v.GetString("tailnet"))
	tailnet, tailnetSet, err := resolveWithPrecedence(
		v,
		"tailnet",
		"TAILSCALE_TAILNET",
		legacyTailnet,
		state,
		flagOverrides,
		func(profile TailnetProfile) (string, error) { return profile.EffectiveTailnet(), nil },
	)
	if err != nil {
		return ResolvedCommandAuth{}, err
	}
	if !tailnetSet || strings.TrimSpace(tailnet) == "" {
		tailnet = "-"
	}

	if auth, ok, err := resolveProfileCommandAuth(v, state, flagOverrides); err != nil {
		return ResolvedCommandAuth{}, err
	} else if ok {
		auth.Tailnet = strings.TrimSpace(tailnet)
		setResolvedCommandAuth(v, auth)
		return auth, nil
	}

	legacyAPIKey := strings.TrimSpace(v.GetString("api-key"))
	apiKey, apiKeySet, err := resolveWithPrecedence(
		v,
		"api-key",
		"TAILSCALE_API_KEY",
		legacyAPIKey,
		TailnetProfilesState{},
		flagOverrides,
		nil,
	)
	if err != nil {
		return ResolvedCommandAuth{}, err
	}
	if apiKeySet {
		if strings.TrimSpace(apiKey) == "" {
			return ResolvedCommandAuth{}, fmt.Errorf("a Tailscale API key is required")
		}
		auth := ResolvedCommandAuth{
			Tailnet: strings.TrimSpace(tailnet),
			APIKey:  strings.TrimSpace(apiKey),
			Source:  "api-key",
		}
		setResolvedCommandAuth(v, auth)
		return auth, nil
	}

	creds, err := resolveOAuthRuntimeConfig(v, flagOverrides)
	if err != nil {
		return ResolvedCommandAuth{}, err
	}
	auth := ResolvedCommandAuth{
		Tailnet:   strings.TrimSpace(tailnet),
		OAuth:     creds,
		UsesOAuth: true,
		Source:    "oauth",
	}
	setResolvedCommandAuth(v, auth)
	return auth, nil
}

func resolveProfileCommandAuth(v *viper.Viper, state TailnetProfilesState, flagOverrides map[string]struct{}) (ResolvedCommandAuth, bool, error) {
	if _, ok := flagOverrides["api-key"]; ok {
		return ResolvedCommandAuth{}, false, nil
	}
	if _, ok := flagOverrides["oauth-client-id"]; ok {
		return ResolvedCommandAuth{}, false, nil
	}
	if _, ok := flagOverrides["oauth-client-secret"]; ok {
		return ResolvedCommandAuth{}, false, nil
	}
	if _, ok := os.LookupEnv("TAILSCALE_API_KEY"); ok {
		return ResolvedCommandAuth{}, false, nil
	}
	if _, ok := os.LookupEnv("TSCLI_OAUTH_CLIENT_ID"); ok {
		return ResolvedCommandAuth{}, false, nil
	}
	if _, ok := os.LookupEnv("TSCLI_OAUTH_CLIENT_SECRET"); ok {
		return ResolvedCommandAuth{}, false, nil
	}
	if state.ActiveTailnet == "" {
		return ResolvedCommandAuth{}, false, nil
	}

	profile, found := findTailnetProfile(state.Tailnets, state.ActiveTailnet)
	if !found {
		return ResolvedCommandAuth{}, false, nil
	}

	switch profile.AuthType() {
	case "api-key":
		apiKey, err := profile.ResolveAPIKey(v)
		if err != nil {
			return ResolvedCommandAuth{}, false, err
		}
		apiKey = strings.TrimSpace(apiKey)
		if apiKey == "" {
			return ResolvedCommandAuth{}, false, fmt.Errorf("a Tailscale API key is required")
		}
		return ResolvedCommandAuth{APIKey: apiKey, Source: "api-key"}, true, nil
	case "oauth":
		secret, err := profile.ResolveOAuthClientSecret(v)
		if err != nil {
			return ResolvedCommandAuth{}, false, err
		}
		creds := ResolvedOAuthConfig{
			ClientID:     strings.TrimSpace(profile.OAuthClientID),
			ClientSecret: strings.TrimSpace(secret),
		}
		if creds.ClientID == "" || creds.ClientSecret == "" {
			return ResolvedCommandAuth{}, false, fmt.Errorf("OAuth client credentials are required; provide --oauth-client-id and --oauth-client-secret, set TSCLI_OAUTH_CLIENT_ID and TSCLI_OAUTH_CLIENT_SECRET, or configure them on the active profile")
		}
		return ResolvedCommandAuth{OAuth: creds, UsesOAuth: true, Source: "oauth"}, true, nil
	default:
		return ResolvedCommandAuth{}, false, nil
	}
}

func setResolvedCommandAuth(v *viper.Viper, auth ResolvedCommandAuth) {
	v.Set("tailnet", strings.TrimSpace(auth.Tailnet))
	v.Set("api-key", strings.TrimSpace(auth.APIKey))
	if auth.UsesOAuth {
		v.Set("oauth-client-id", auth.OAuth.ClientID)
		v.Set("oauth-client-secret", auth.OAuth.ClientSecret)
		v.Set("api-key", "")
		return
	}
	v.Set("oauth-client-id", "")
	v.Set("oauth-client-secret", "")
}

func resolveWithPrecedence(
	v *viper.Viper,
	key string,
	envKey string,
	legacyValue string,
	state TailnetProfilesState,
	flagOverrides map[string]struct{},
	profileValue func(TailnetProfile) (string, error),
) (string, bool, error) {
	if _, ok := flagOverrides[key]; ok {
		value := strings.TrimSpace(v.GetString(key))
		return value, true, nil
	}

	if envValue, ok := os.LookupEnv(envKey); ok {
		envValue = strings.TrimSpace(envValue)
		return envValue, true, nil
	}

	if state.ActiveTailnet != "" {
		if profile, found := findTailnetProfile(state.Tailnets, state.ActiveTailnet); found {
			if profileValue != nil {
				value, err := profileValue(profile)
				if err != nil {
					return "", false, err
				}
				value = strings.TrimSpace(value)
				if value != "" {
					return value, true, nil
				}
			}
		}
	}

	legacyValue = strings.TrimSpace(legacyValue)
	if legacyValue != "" {
		return legacyValue, true, nil
	}

	return "", false, nil
}

func loadTailnetProfilesState(v *viper.Viper) (TailnetProfilesState, error) {
	if err := validateAgeEncryptionConfig(loadAgeEncryptionConfig(v)); err != nil {
		return TailnetProfilesState{}, err
	}

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
		Name:                       strings.TrimSpace(profile.Name),
		Tailnet:                    strings.TrimSpace(profile.Tailnet),
		APIKey:                     strings.TrimSpace(profile.APIKey),
		APIKeyEncrypted:            strings.TrimSpace(profile.APIKeyEncrypted),
		OAuthClientID:              strings.TrimSpace(profile.OAuthClientID),
		OAuthClientSecret:          strings.TrimSpace(profile.OAuthClientSecret),
		OAuthClientSecretEncrypted: strings.TrimSpace(profile.OAuthClientSecretEncrypted),
	}
}

func validateProfileAuthShape(profile TailnetProfile) error {
	hasAPIKey := profile.APIKey != ""
	hasAPIKeyEncrypted := profile.APIKeyEncrypted != ""
	hasOAuthID := profile.OAuthClientID != ""
	hasOAuthSecret := profile.OAuthClientSecret != ""
	hasOAuthSecretEncrypted := profile.OAuthClientSecretEncrypted != ""
	hasOAuth := hasOAuthID || hasOAuthSecret || hasOAuthSecretEncrypted

	switch {
	case hasAPIKey && hasAPIKeyEncrypted:
		return fmt.Errorf("tailnet profile %q must not include both api-key and api-key-encrypted", profile.Name)
	case hasOAuthSecret && hasOAuthSecretEncrypted:
		return fmt.Errorf("tailnet profile %q must not include both oauth-client-secret and oauth-client-secret-encrypted", profile.Name)
	case (hasAPIKey || hasAPIKeyEncrypted) && hasOAuth:
		return fmt.Errorf("tailnet profile %q must use either api-key auth or oauth-client-id/oauth-client-secret auth, not both", profile.Name)
	case hasAPIKey || hasAPIKeyEncrypted:
		return nil
	case hasOAuthID && (hasOAuthSecret || hasOAuthSecretEncrypted):
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
	if p.APIKey != "" || p.APIKeyEncrypted != "" {
		return "api-key"
	}
	if p.OAuthClientID != "" || p.OAuthClientSecret != "" || p.OAuthClientSecretEncrypted != "" {
		return "oauth"
	}
	return "unknown"
}

func (p TailnetProfile) ResolveAPIKey(v *viper.Viper) (string, error) {
	if p.APIKey != "" {
		return p.APIKey, nil
	}
	if p.APIKeyEncrypted != "" {
		return decryptSecret(v, p.APIKeyEncrypted)
	}
	return "", nil
}

func (p TailnetProfile) ResolveOAuthClientSecret(v *viper.Viper) (string, error) {
	if p.OAuthClientSecret != "" {
		return p.OAuthClientSecret, nil
	}
	if p.OAuthClientSecretEncrypted != "" {
		return decryptSecret(v, p.OAuthClientSecretEncrypted)
	}
	return "", nil
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

	encryptedProfiles, err := encryptProfilesForPersistence(v, state.Tailnets)
	if err != nil {
		return err
	}
	state.Tailnets = encryptedProfiles

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
