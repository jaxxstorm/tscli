// cmd/tscli/create/key/cli.go
package key

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/jaxxstorm/tscli/pkg/output"

	"github.com/jaxxstorm/tscli/pkg/tscli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tsapi "tailscale.com/client/tailscale/v2"
)

var keyKinds = map[string]struct{}{
	"authkey":     {},
	"oauthclient": {},
	"federated":   {},
}

var newClient = tscli.New
var doRequest = tscli.Do

var scopeEnum = map[string]struct{}{
	"devices:core":  {},
	"devices:read":  {},
	"devices:write": {},
	"dns:read":      {},
	"dns:write":     {},
	"logging:read":  {},
	"logging:write": {},
	"tailnet:read":  {},
	"tailnet:write": {},
	"users:read":    {},
	"users:write":   {},
	"auth_keys":     {},
}

func scopesNeedTags(sc []string) bool {
	for _, s := range sc {
		if s == "devices:core" || s == "auth_keys" {
			return true
		}
	}
	return false
}

type federatedKeyRequest struct {
	KeyType          string            `json:"keyType"`
	Description      string            `json:"description,omitempty"`
	Scopes           []string          `json:"scopes"`
	Tags             []string          `json:"tags,omitempty"`
	Issuer           string            `json:"issuer"`
	Subject          string            `json:"subject"`
	Audience         string            `json:"audience,omitempty"`
	CustomClaimRules map[string]string `json:"customClaimRules,omitempty"`
}

func parseClaimRules(claims []string) (map[string]string, error) {
	if len(claims) == 0 {
		return nil, nil
	}
	rules := make(map[string]string, len(claims))
	for _, claim := range claims {
		parts := strings.SplitN(claim, "=", 2)
		if len(parts) != 2 || parts[0] == "" {
			return nil, fmt.Errorf("--claim must be key=value")
		}
		rules[parts[0]] = parts[1]
	}
	return rules, nil
}

func printKey(format string, key *tsapi.Key) error {
	if key == nil {
		return fmt.Errorf("no key data returned")
	}
	b, err := json.MarshalIndent(key, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal key: %w", err)
	}
	return output.Print(format, b)
}

func Command() *cobra.Command {
	var (
		kind          string
		desc          string
		expiry        time.Duration
		scopes        []string
		tags          []string
		reusable      bool
		ephemeral     bool
		preauthorized bool
		issuer        string
		subject       string
		audience      string
		claims        []string
	)

	cmd := &cobra.Command{
		Use:   "key",
		Short: "Create an auth-key, OAuth client, or federated credential",
		Long: "Create Tailscale auth-keys, OAuth clients, or federated identities.\n" +
			"Auth-key capability flags (--reusable, --ephemeral, --preauthorized) apply only to --type authkey.",
		Example: "  tscli create key --type authkey --description \"CI runner\" --expiry 720h --reusable --preauthorized\n" +
			"  tscli create key --type oauthclient --scope users:read --scope devices:read\n" +
			"  tscli create key --type federated --scope users:read --issuer https://example.com --subject example-*",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			kind = strings.ToLower(kind)
			if kind == "" {
				kind = "authkey"
			}
			if _, ok := keyKinds[kind]; !ok {
				return fmt.Errorf("--type must be authkey, oauthclient, or federated")
			}

			if kind == "oauthclient" || kind == "federated" {
				if len(scopes) == 0 {
					return fmt.Errorf("--scope is required for %s", kind)
				}
				for _, s := range scopes {
					if _, ok := scopeEnum[s]; !ok {
						return fmt.Errorf("invalid scope %q", s)
					}
				}
				if scopesNeedTags(scopes) && len(tags) == 0 {
					return fmt.Errorf("--tags required when scope includes devices:core or auth_keys")
				}
			}
			if kind == "federated" {
				if issuer == "" {
					return fmt.Errorf("--issuer is required for federated credentials")
				}
				parsed, err := url.Parse(issuer)
				if err != nil {
					return fmt.Errorf("--issuer must be a valid URL: %w", err)
				}
				if parsed.Scheme != "https" || parsed.Host == "" {
					return fmt.Errorf("--issuer must be an https URL")
				}
				if subject == "" {
					return fmt.Errorf("--subject is required for federated credentials")
				}
			}
			return nil
		},

		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := newClient()
			if err != nil {
				return err
			}
			ctx := context.Background()
			outputType := viper.GetString("output")

			var key *tsapi.Key
			switch kind {
			case "authkey":
				req := tsapi.CreateKeyRequest{
					Description:  desc,
					Capabilities: tsapi.KeyCapabilities{},
				}
				req.Capabilities.Devices.Create.Reusable = reusable
				req.Capabilities.Devices.Create.Ephemeral = ephemeral
				req.Capabilities.Devices.Create.Preauthorized = preauthorized
				if cmd.Flags().Lookup("expiry").Changed {
					req.ExpirySeconds = int64(expiry.Seconds())
				}
				key, err = client.Keys().CreateAuthKey(ctx, req)
				if err != nil {
					return fmt.Errorf("create auth-key: %w", err)
				}
			case "oauthclient":
				req := tsapi.CreateOAuthClientRequest{
					Description: desc,
					Scopes:      scopes,
				}
				if len(tags) > 0 {
					req.Tags = tags
				}
				key, err = client.Keys().CreateOAuthClient(ctx, req)
				if err != nil {
					return fmt.Errorf("create oauth client: %w", err)
				}
			case "federated":
				ruleMap, err := parseClaimRules(claims)
				if err != nil {
					return err
				}
				req := federatedKeyRequest{
					KeyType:          "federated",
					Description:      desc,
					Scopes:           scopes,
					Tags:             tags,
					Issuer:           issuer,
					Subject:          subject,
					Audience:         audience,
					CustomClaimRules: ruleMap,
				}
				key = &tsapi.Key{}
				if _, err := doRequest(ctx, client, http.MethodPost, "/tailnet/{tailnet}/keys", req, key); err != nil {
					return fmt.Errorf("create federated credential: %w", err)
				}
			default:
				return fmt.Errorf("unsupported key type %q", kind)
			}

			return printKey(outputType, key)
		},
	}

	cmd.Flags().StringVar(&kind, "type", "authkey", "Key type: authkey|oauthclient|federated")
	cmd.Flags().StringVar(&desc, "description", "", "Short description (≤50 chars)")
	cmd.Flags().DurationVar(&expiry, "expiry", 0, "Expiry duration (e.g. 720h) for auth-keys")
	cmd.Flags().BoolVar(&reusable, "reusable", false, "Auth-key only: allow the key to be used multiple times")
	cmd.Flags().BoolVar(&ephemeral, "ephemeral", false, "Auth-key only: mark devices authenticated with this key as ephemeral")
	cmd.Flags().BoolVar(&preauthorized, "preauthorized", false, "Auth-key only: create key in preauthorized state")
	cmd.Flags().StringSliceVar(&scopes, "scope", nil, "OAuth/federated scopes (repeatable)")
	cmd.Flags().StringSliceVar(&tags, "tags", nil, "Allowed tags (repeatable) for OAuth client and federated credentials")
	cmd.Flags().StringVar(&issuer, "issuer", "", "Federated only: issuer HTTPS URL")
	cmd.Flags().StringVar(&subject, "subject", "", "Federated only: subject pattern")
	cmd.Flags().StringVar(&audience, "audience", "", "Federated only: audience hint")
	cmd.Flags().StringSliceVar(&claims, "claim", nil, "Federated only: custom claim rule key=value (repeatable)")

	return cmd
}
