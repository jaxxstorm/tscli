package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/jaxxstorm/tscli/cmd/tscli/set/service/approval"
	f "github.com/jaxxstorm/tscli/pkg/file"
	"github.com/jaxxstorm/tscli/pkg/output"
	"github.com/jaxxstorm/tscli/pkg/tscli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Command() *cobra.Command {
	var (
		serviceName string
		filePath    string
		inline      string
	)

	cmd := &cobra.Command{
		Use:   "service",
		Short: "Update a service",
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			if serviceName == "" {
				return fmt.Errorf("--service is required")
			}
			if filePath == "" && inline == "" {
				return fmt.Errorf("one of --file or --body is required")
			}
			if filePath != "" && inline != "" {
				return fmt.Errorf("--file and --body are mutually exclusive")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := tscli.New()
			if err != nil {
				return err
			}

			raw, err := f.ReadInput(filePath, inline)
			if err != nil {
				return err
			}

			var payload any
			if err := json.Unmarshal(raw, &payload); err != nil {
				return fmt.Errorf("invalid JSON payload: %w", err)
			}

			path := fmt.Sprintf("/tailnet/{tailnet}/services/%s", url.PathEscape(serviceName))
			var resp json.RawMessage
			if _, err := tscli.Do(
				context.Background(),
				client,
				http.MethodPut,
				path,
				payload,
				&resp,
			); err != nil {
				return fmt.Errorf("failed to update service %s: %w", serviceName, err)
			}

			out, _ := json.MarshalIndent(resp, "", "  ")
			return output.Print(viper.GetString("output"), out)
		},
	}

	cmd.Flags().StringVar(&serviceName, "service", "", "Service name")
	cmd.Flags().StringVar(&filePath, "file", "", "Path to a JSON payload file, file://path, or '-' for stdin")
	cmd.Flags().StringVar(&inline, "body", "", "Inline JSON payload")
	cmd.AddCommand(approval.Command())
	return cmd
}
