// cmd/tscli/get/policy/validate/cli.go
//
// `tscli get policy validate --file <policy.json>`
// Validate a policy file using the Tailscale API
package validate

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	f "github.com/jaxxstorm/tscli/pkg/file"
	"github.com/jaxxstorm/tscli/pkg/output"
	"github.com/jaxxstorm/tscli/pkg/tscli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Command() *cobra.Command {
	var (
		file string
		body string
	)

	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate a policy file",
		Long: `Validate a policy file using the Tailscale API.
The policy can be provided via --file or --body flag.`,

		PreRunE: func(cmd *cobra.Command, _ []string) error {
			// exactly one source
			sources := 0
			if file != "" {
				sources++
			}
			if body != "" {
				sources++
			}
			if sources == 0 {
				return fmt.Errorf("provide policy via --file or --body")
			}
			if sources > 1 {
				return fmt.Errorf("--file and --body are mutually exclusive")
			}
			return nil
		},

		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := tscli.New()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			var raw []byte
			if file != "" {
				raw, err = f.ReadInput(file, "")
				if err != nil {
					return fmt.Errorf("failed to read file: %w", err)
				}
			} else {
				raw = []byte(body)
			}

			// Parse the policy file into the expected JSON structure
			var policyData any
			if err := json.Unmarshal(raw, &policyData); err != nil {
				return fmt.Errorf("failed to parse policy file: %w", err)
			}

			var response map[string]interface{}
			if _, err := tscli.Do(
				context.Background(),
				client,
				http.MethodPost,
				"/tailnet/{tailnet}/acl/validate",
				policyData,
				&response,
			); err != nil {
				return fmt.Errorf("policy validation failed: %w", err)
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

	cmd.Flags().StringVar(&file, "file", "", "Path to policy file, file://path or '-' for stdin")
	cmd.Flags().StringVar(&body, "body", "", "Inline policy JSON")

	return cmd
}
