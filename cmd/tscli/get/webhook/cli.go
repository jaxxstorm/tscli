// cmd/tscli/get/webhook/cli.go
package webhook

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jaxxstorm/tscli/cmd/tscli/get/webhook/test"
	"github.com/jaxxstorm/tscli/pkg/output"

	"github.com/jaxxstorm/tscli/pkg/tscli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tsapi "tailscale.com/client/tailscale/v2"
)

func Command() *cobra.Command {
	var hookID string

	runGetWebhook := func(id *string) func(*cobra.Command, []string) error {
		return func(cmd *cobra.Command, _ []string) error {
			client, err := tscli.New()
			if err != nil {
				return err
			}

			var hook *tsapi.Webhook
			hook, err = client.Webhooks().Get(context.Background(), *id)
			if err != nil {
				return fmt.Errorf("failed to get webhook %s: %w", *id, err)
			}

			out, err := json.MarshalIndent(hook, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal webhook: %w", err)
			}
			outputType := viper.GetString("output")
			output.Print(outputType, out)
			return nil
		}
	}

	cmd := &cobra.Command{
		Use:   "webhook",
		Short: "Get information about a webhook from the Tailscale API",
		RunE:  runGetWebhook(&hookID),
	}

	cmd.Flags().StringVar(&hookID, "id", "", "Webhook ID to retrieve")
	_ = cmd.MarkFlagRequired("id")

	// Legacy path was `get webhook webhook --id ...`.
	var legacyHookID string
	legacy := &cobra.Command{
		Use:    "webhook",
		Short:  "Legacy alias for `get webhook`",
		Hidden: true,
		RunE:   runGetWebhook(&legacyHookID),
	}
	legacy.Flags().StringVar(&legacyHookID, "id", "", "Webhook ID to retrieve")
	_ = legacy.MarkFlagRequired("id")
	cmd.AddCommand(legacy)

	legacyTest := test.Command()
	legacyTest.Hidden = true
	legacyTest.Short = "Legacy alias for `set webhook test`"
	cmd.AddCommand(legacyTest)

	return cmd
}
