package stream

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	f "github.com/jaxxstorm/tscli/pkg/file"
	"github.com/jaxxstorm/tscli/pkg/output"
	"github.com/jaxxstorm/tscli/pkg/tscli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tsapi "tailscale.com/client/tailscale/v2"
)

var validLogTypes = map[string]struct{}{
	string(tsapi.LogTypeConfig):  {},
	string(tsapi.LogTypeNetwork): {},
}

func Command() *cobra.Command {
	var (
		logType  string
		filePath string
		inline   string
	)

	cmd := &cobra.Command{
		Use:   "stream",
		Short: "Set log streaming configuration",
		PreRunE: func(_ *cobra.Command, _ []string) error {
			if _, ok := validLogTypes[logType]; !ok {
				return fmt.Errorf(`--type must be "configuration" or "network"`)
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

			path := fmt.Sprintf("/tailnet/{tailnet}/logging/%s/stream", url.PathEscape(logType))
			var resp json.RawMessage
			if _, err := tscli.Do(
				context.Background(),
				client,
				http.MethodPut,
				path,
				payload,
				&resp,
			); err != nil {
				return fmt.Errorf("failed to update stream configuration: %w", err)
			}

			out, _ := json.MarshalIndent(resp, "", "  ")
			return output.Print(viper.GetString("output"), out)
		},
	}

	cmd.Flags().StringVar(&logType, "type", "", `Log type: "configuration" or "network" (required)`)
	cmd.Flags().StringVar(&filePath, "file", "", "Path to a JSON payload file, file://path, or '-' for stdin")
	cmd.Flags().StringVar(&inline, "body", "", "Inline JSON payload")
	_ = cmd.MarkFlagRequired("type")
	return cmd
}
