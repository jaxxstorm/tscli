// cmd/tscli/get/logs/aws/cli.go
//
// `tscli get logs aws external-id`
// Get or create AWS external ID for log streaming integration
package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jaxxstorm/tscli/cmd/tscli/get/logs/aws/validate"
	"github.com/jaxxstorm/tscli/pkg/output"
	"github.com/jaxxstorm/tscli/pkg/tscli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "aws",
		Short: "Get AWS external ID for log streaming",
		Long: `Get or create an AWS external ID for use in IAM role trust policies.
This external ID is required when setting up AWS log streaming integration.`,

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
				"/tailnet/{tailnet}/aws-external-id",
				nil,
				&response,
			); err != nil {
				return fmt.Errorf("failed to get or create AWS external ID: %w", err)
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

	// Add subcommands
	cmd.AddCommand(validate.Command())

	return cmd
}
