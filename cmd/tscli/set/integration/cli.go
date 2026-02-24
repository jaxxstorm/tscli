// cmd/tscli/set/posture-integration/cli.go
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
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var validProv = map[string]struct{}{
	"falcon": {}, "intune": {}, "jamfpro": {}, "kandji": {},
	"kolide": {}, "sentinelone": {},
}

func Command() *cobra.Command {
	var (
		id, provider      string
		cloudID, clientID string
		tenantID, secret  string
	)

	cmd := &cobra.Command{
		Use:   "posture-integration",
		Short: "Create or update a posture integration",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			// ensure at least one field besides --id was supplied
			changed := 0
			cmd.Flags().VisitAll(func(f *pflag.Flag) {
				if f.Name != "id" && f.Changed {
					changed++
				}
			})
			if changed == 0 {
				return fmt.Errorf("provide at least one field flag to update")
			}

			if cmd.Flags().Lookup("provider").Changed {
				provider = strings.ToLower(provider)
				if _, ok := validProv[provider]; !ok {
					return fmt.Errorf("invalid --provider: %s", provider)
				}
			}
			return nil
		},

		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := tscli.New()
			if err != nil {
				return err
			}

			// Build PATCH body with only changed flags
			body := make(map[string]any)
			f := cmd.Flags()

			if f.Lookup("provider").Changed {
				body["provider"] = provider
			}
			if f.Lookup("cloud-id").Changed {
				body["cloudId"] = cloudID
			}
			if f.Lookup("client-id").Changed {
				body["clientId"] = clientID
			}
			if f.Lookup("tenant-id").Changed {
				body["tenantId"] = tenantID
			}
			if f.Lookup("client-secret").Changed {
				body["clientSecret"] = secret
			}

			_, err = tscli.Do(
				context.Background(),
				client,
				http.MethodPatch,
				"/posture/integrations/"+id,
				body,
				nil,
			)
			if err != nil {
				return fmt.Errorf("update failed: %w", err)
			}

			out, _ := json.MarshalIndent(map[string]any{
				"result": "integration updated",
				"id":     id,
				"fields": body,
			}, "", "  ")
			outputType := viper.GetString("output")
			output.Print(outputType, out)
			return nil
		},
	}

	cmd.Flags().StringVar(&id, "id", "", "Integration identifier")
	cmd.Flags().StringVar(&provider, "provider", "", "Provider (falcon|intune|jamfpro|kandji|kolide|sentinelone)")
	cmd.Flags().StringVar(&cloudID, "cloud-id", "", "Provider cloud/region ID")
	cmd.Flags().StringVar(&clientID, "client-id", "", "Client/application ID")
	cmd.Flags().StringVar(&tenantID, "tenant-id", "", "Microsoft Intune tenant ID")
	cmd.Flags().StringVar(&secret, "client-secret", "", "Client secret / API token")

	_ = cmd.MarkFlagRequired("id")
	return cmd
}
