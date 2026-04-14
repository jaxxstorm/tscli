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
		oauthClientID     string
		oauthClientSecret string
	)

	cmd := &cobra.Command{
		Use:   "tailnet",
		Short: "Delete an API-driven tailnet",
		Long:  "Delete the current API-driven tailnet using that tailnet's OAuth client credentials.",
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

			if _, err := tscli.DoBearer(cmd.Context(), "DELETE", "/tailnet/-", tokenResp.AccessToken, nil, nil); err != nil {
				return fmt.Errorf("delete tailnet: %w", err)
			}

			resp := map[string]string{"result": "tailnet deleted"}
			out, err := json.MarshalIndent(resp, "", "  ")
			if err != nil {
				return fmt.Errorf("marshal delete tailnet response: %w", err)
			}

			return output.Print(v.GetString("output"), out)
		},
	}

	cmd.Flags().StringVar(&oauthClientID, "oauth-client-id", "", "OAuth client ID for deleting the API-driven tailnet")
	cmd.Flags().StringVar(&oauthClientSecret, "oauth-client-secret", "", "OAuth client secret for deleting the API-driven tailnet")

	return cmd
}
