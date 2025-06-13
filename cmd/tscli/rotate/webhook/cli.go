// cmd/tscli/rotate/webhook/cli.go
package webhook

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jaxxstorm/tscli/pkg/output"

	"github.com/jaxxstorm/tscli/pkg/tscli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tsapi "tailscale.com/client/tailscale/v2"
)

func Command() *cobra.Command {
	var hookID string

	cmd := &cobra.Command{
		Use:   "webhook",
		Short: "Rotate a webhook from the Tailscale API",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := tscli.New()
			if err != nil {
				return err
			}

			var hook *tsapi.Webhook
			hook, err = client.Webhooks().RotateSecret(context.Background(), hookID)
			if err != nil {
				return fmt.Errorf("failed to rotate secret for webhook %s: %w", hookID, err)
			}

			out, err := json.MarshalIndent(hook, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal webhook: %w", err)
			}
			format := viper.GetString("format")
			output.Print(format, out)
			return nil
		},
	}

	cmd.Flags().StringVar(&hookID, "id", "", "Webhook ID to rotate")
	_ = cmd.MarkFlagRequired("id")
	return cmd
}
