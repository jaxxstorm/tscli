// cmd/tscli/get/webhook/test/cli.go
package test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jaxxstorm/tscli/pkg/output"
	"github.com/jaxxstorm/tscli/pkg/tscli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test",
		Short: "Test webhook delivery",
		Long: `Test webhook delivery by sending a test event to all configured webhooks.
This sends a test event to verify that webhook endpoints are reachable and configured correctly.`,

		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := tscli.New()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			var response map[string]interface{}
			if _, err := tscli.Do(
				context.Background(),
				client,
				http.MethodPost,
				"/webhooks/test",
				nil,
				&response,
			); err != nil {
				return fmt.Errorf("failed to test webhook: %w", err)
			}

			out, err := json.MarshalIndent(response, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal response: %w", err)
			}

			outputType := viper.GetString("output")
			output.Print(outputType, out)
			return nil
		},
	}

	return cmd
}
