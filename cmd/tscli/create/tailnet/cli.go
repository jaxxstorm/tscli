package tailnet

import (
	"encoding/json"
	"fmt"

	"github.com/jaxxstorm/tscli/pkg/config"
	"github.com/jaxxstorm/tscli/pkg/oauth"
	"github.com/jaxxstorm/tscli/pkg/output"
	"github.com/jaxxstorm/tscli/pkg/tscli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Command() *cobra.Command {
	var (
		displayName       string
		oauthClientID     string
		oauthClientSecret string
	)

	cmd := &cobra.Command{
		Use:   "tailnet",
		Short: "Create an API-driven tailnet",
		Long:  "Create an API-driven tailnet using organization-approved OAuth client credentials.",
		PreRunE: func(_ *cobra.Command, _ []string) error {
			if displayName == "" {
				return fmt.Errorf("--display-name is required")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			v := viper.GetViper()
			_ = v.BindPFlags(cmd.Flags())

			creds, err := config.ResolveOAuthRuntimeConfig(config.ChangedMap(cmd))
			if err != nil {
				return err
			}

			tokenResp, err := oauth.ExchangeClientCredentials(cmd.Context(), creds.ClientID, creds.ClientSecret)
			if err != nil {
				return fmt.Errorf("failed to exchange OAuth credentials: %w", err)
			}

			var resp map[string]any
			if _, err := tscli.DoBearer(cmd.Context(), "POST", "/organizations/-/tailnets", tokenResp.AccessToken, map[string]string{
				"displayName": displayName,
			}, &resp); err != nil {
				return fmt.Errorf("create tailnet: %w", err)
			}

			out, err := json.MarshalIndent(resp, "", "  ")
			if err != nil {
				return fmt.Errorf("marshal create tailnet response: %w", err)
			}

			return output.Print(v.GetString("output"), out)
		},
	}

	cmd.Flags().StringVar(&displayName, "display-name", "", "Unique display name for the API-driven tailnet")
	cmd.Flags().StringVar(&oauthClientID, "oauth-client-id", "", "OAuth client ID for organization tailnet lifecycle access")
	cmd.Flags().StringVar(&oauthClientSecret, "oauth-client-secret", "", "OAuth client secret for organization tailnet lifecycle access")

	return cmd
}
