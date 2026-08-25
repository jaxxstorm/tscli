package tailnets

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/jaxxstorm/tscli/pkg/apitype"
	"github.com/jaxxstorm/tscli/pkg/config"
	"github.com/jaxxstorm/tscli/pkg/output"
	"github.com/jaxxstorm/tscli/pkg/tscli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Command() *cobra.Command {
	var (
		oauthClientID     string
		oauthClientSecret string
		limit             int
		cursor            string
	)

	cmd := &cobra.Command{
		Use:   "tailnets",
		Short: "List organization tailnets",
		Long:  "List the organization tailnets visible to an OAuth client approved for tailnet lifecycle APIs.",
		PreRunE: func(*cobra.Command, []string) error {
			if limit < 1 || limit > 100 {
				return fmt.Errorf("--limit must be between 1 and 100")
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

			tokenResp, err := tscli.ExchangeOAuthClientCredentials(cmd.Context(), creds.ClientID, creds.ClientSecret)
			if err != nil {
				return fmt.Errorf("failed to exchange OAuth credentials: %w", err)
			}

			q := url.Values{}
			q.Set("limit", fmt.Sprint(limit))
			if cursor != "" {
				q.Set("cursor", cursor)
			}

			var raw json.RawMessage
			if _, err := tscli.DoBearer(cmd.Context(), http.MethodGet, "/organizations/-/tailnets?"+q.Encode(), tokenResp.AccessToken, nil, &raw); err != nil {
				return fmt.Errorf("list tailnets: %w", err)
			}

			out := raw
			if outputType := v.GetString("output"); outputType == "pretty" || outputType == "human" {
				var listResp apitype.OrganizationTailnetListResponse
				if err := json.Unmarshal(raw, &listResp); err != nil {
					return fmt.Errorf("decode list tailnets response: %w", err)
				}
				out, err = json.MarshalIndent(listResp.Tailnets, "", "  ")
			} else {
				out, err = json.MarshalIndent(raw, "", "  ")
			}
			if err != nil {
				return fmt.Errorf("marshal list tailnets response: %w", err)
			}

			return output.Print(v.GetString("output"), out)
		},
	}

	cmd.Flags().StringVar(&oauthClientID, "oauth-client-id", "", "OAuth client ID for organization tailnet lifecycle access")
	cmd.Flags().StringVar(&oauthClientSecret, "oauth-client-secret", "", "OAuth client secret for organization tailnet lifecycle access")
	cmd.Flags().IntVar(&limit, "limit", 100, "Maximum number of tailnets to return (1-100)")
	cmd.Flags().StringVar(&cursor, "cursor", "", "Pagination cursor from a previous response")

	return cmd
}
