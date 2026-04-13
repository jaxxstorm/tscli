// cmd/tscli/create/posture-integration/cli.go
//
// Create a brand-new device-posture integration.
//
//	tscli create posture-integration \
//	    --provider falcon \
//	    --cloud-id us-1 \
//	    --client-id ABC123 \
//	    --client-secret superSecret
package postureintegration

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/jaxxstorm/tscli/pkg/output"

	"github.com/jaxxstorm/tscli/pkg/tscli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var providerEnum = map[string]struct{}{
	"falcon": {}, "intune": {}, "jamfpro": {}, "kandji": {},
	"kolide": {}, "sentinelone": {},
}

func Command() *cobra.Command {
	var (
		provider          string
		cloudID, clientID string
		tenantID, secret  string
	)

	cmd := &cobra.Command{
		Use:   "posture-integration",
		Short: "Create a new posture integration",
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			provider = strings.ToLower(provider)
			if _, ok := providerEnum[provider]; !ok {
				return fmt.Errorf("invalid --provider: %s", provider)
			}
			if secret == "" {
				return fmt.Errorf("--client-secret is required when creating an integration")
			}
			// ensure at least provider+secret plus one extra field? not necessary
			return nil
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := tscli.New()
			if err != nil {
				return err
			}

			// Build POST body with supplied flags only
			body := map[string]any{
				"provider":     provider,
				"clientSecret": secret,
			}
			f := cmd.Flags()
			if f.Lookup("cloud-id").Changed {
				body["cloudId"] = cloudID
			}
			if f.Lookup("client-id").Changed {
				body["clientId"] = clientID
			}
			if f.Lookup("tenant-id").Changed {
				body["tenantId"] = tenantID
			}

			// Send POST request
			var resp map[string]any
			if _, err := tscli.Do(
				context.Background(),
				client,
				http.MethodPost,
				"/tailnet/{tailnet}/posture/integrations",
				body,
				&resp,
			); err != nil {
				return fmt.Errorf("creation failed: %w", err)
			}

			out, _ := json.MarshalIndent(resp, "", "  ")
			outputType := viper.GetString("output")
			output.Print(outputType, out)
			return nil
		},
	}

	// Flags
	cmd.Flags().StringVar(&provider, "provider", "", "Provider (falcon|intune|jamfpro|kandji|kolide|sentinelone)")
	cmd.Flags().StringVar(&cloudID, "cloud-id", "", "Provider cloud/region ID")
	cmd.Flags().StringVar(&clientID, "client-id", "", "Client/application ID")
	cmd.Flags().StringVar(&tenantID, "tenant-id", "", "Microsoft Intune tenant ID")
	cmd.Flags().StringVar(&secret, "client-secret", "", "Client secret / API token (required)")

	_ = cmd.MarkFlagRequired("provider")
	_ = cmd.MarkFlagRequired("client-secret")

	return cmd
}
