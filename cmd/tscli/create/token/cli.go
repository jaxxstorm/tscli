// cmd/tscli/create/token/cli.go
//
// Exchange an OAuth client-id/secret for an access-token.
//
//	tscli create token --client-id XXX --client-secret YYY
package token

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jaxxstorm/tscli/pkg/oauth"
	"github.com/jaxxstorm/tscli/pkg/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Command() *cobra.Command {
	var (
		clientID     string
		clientSecret string
	)

	cmd := &cobra.Command{
		Use:   "token",
		Short: "Create an OAuth access-token from a client_id / client_secret pair",
		Long:  "POSTs to /api/v2/oauth/token. No API-key auth is required.",
		Example: `  tscli create token \
      --client-id     "$OAUTH_CLIENT_ID" \
      --client-secret "$OAUTH_CLIENT_SECRET"`,
		PreRunE: func(_ *cobra.Command, _ []string) error {
			if clientID == "" || clientSecret == "" {
				return errors.New("--client-id and --client-secret are both required")
			}
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			// Use the OAuth library for token exchange
			tokenResp, err := oauth.ExchangeClientCredentials(cmd.Context(), clientID, clientSecret)
			if err != nil {
				return fmt.Errorf("failed to exchange OAuth credentials: %w", err)
			}

			tokenBytes, err := json.MarshalIndent(tokenResp, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal token response: %w", err)
			}

			outputType := viper.GetString("output")
			output.Print(outputType, tokenBytes)
			return nil
		},
	}

	cmd.Flags().StringVar(&clientID, "client-id", "", "OAuth client ID (required)")
	cmd.Flags().StringVar(&clientSecret, "client-secret", "", "OAuth client secret (required)")
	_ = cmd.MarkFlagRequired("client-id")
	_ = cmd.MarkFlagRequired("client-secret")

	return cmd
}
