// cmd/tscli/create/key/cli.go
package key

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
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
}

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
	)

	cmd := &cobra.Command{
		Use:   "key",
		Short: "Create an auth-key or OAuth client",
		Long: "Create Tailscale auth-keys or OAuth clients.\n" +
			"Auth-key capability flags (--reusable, --ephemeral, --preauthorized) apply only to --type authkey.",
		Example: "  tscli create key --type authkey --description \"CI runner\" --expiry 720h --reusable --preauthorized\n" +
			"  tscli create key --type oauthclient --scope users:read --scope devices:read",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			kind = strings.ToLower(kind)
			if kind == "" {
				kind = "authkey"
			}
			if _, ok := keyKinds[kind]; !ok {
				return fmt.Errorf("--type must be authkey or oauthclient")
			}

			if kind == "oauthclient" {
				if len(scopes) == 0 {
					return fmt.Errorf("--scope is required for oauthclient")
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
			return nil
		},

		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := tscli.New()
			if err != nil {
				return err
			}
			ctx := context.Background()

			if kind == "authkey" {
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
				key, err := client.Keys().CreateAuthKey(ctx, req)
				if err != nil {
					return fmt.Errorf("create auth-key: %w", err)
				}
				b, _ := json.MarshalIndent(key, "", "  ")
				fmt.Fprintln(os.Stdout, string(b))
				return nil
			}

			req := tsapi.CreateOAuthClientRequest{
				Description: desc,
				Scopes:      scopes,
			}
			if len(tags) > 0 {
				req.Tags = tags
			}
			key, err := client.Keys().CreateOAuthClient(ctx, req)
			if err != nil {
				return fmt.Errorf("create oauth client: %w", err)
			}
			out, _ := json.MarshalIndent(key, "", "  ")
			outputType := viper.GetString("output")
			output.Print(outputType, out)
			return nil
		},
	}

	cmd.Flags().StringVar(&kind, "type", "authkey", "Key type: authkey|oauthclient")
	cmd.Flags().StringVar(&desc, "description", "", "Short description (≤50 chars)")
	cmd.Flags().DurationVar(&expiry, "expiry", 0, "Expiry duration (e.g. 720h) for auth-keys")
	cmd.Flags().BoolVar(&reusable, "reusable", false, "Auth-key only: allow the key to be used multiple times")
	cmd.Flags().BoolVar(&ephemeral, "ephemeral", false, "Auth-key only: mark devices authenticated with this key as ephemeral")
	cmd.Flags().BoolVar(&preauthorized, "preauthorized", false, "Auth-key only: create key in preauthorized state")
	cmd.Flags().StringSliceVar(&scopes, "scope", nil, "OAuth scopes (repeatable)")
	cmd.Flags().StringSliceVar(&tags, "tags", nil, "Allowed tags (repeatable) for OAuth client")

	return cmd
}
