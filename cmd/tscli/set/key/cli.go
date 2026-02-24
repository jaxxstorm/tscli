package key

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
		keyID    string
		filePath string
		inline   string
	)

	cmd := &cobra.Command{
		Use:   "key",
		Short: "Update a tailnet key",
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			if keyID == "" {
				return fmt.Errorf("--key is required")
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

			var resp json.RawMessage
			if _, err := tscli.Do(
				context.Background(),
				client,
				http.MethodPut,
				"/tailnet/{tailnet}/keys/"+keyID,
				payload,
				&resp,
			); err != nil {
				return fmt.Errorf("failed to update key %s: %w", keyID, err)
			}

			out, _ := json.MarshalIndent(resp, "", "  ")
			return output.Print(viper.GetString("output"), out)
		},
	}

	cmd.Flags().StringVar(&keyID, "key", "", "Key ID to update")
	cmd.Flags().StringVar(&filePath, "file", "", "Path to a JSON payload file, file://path, or '-' for stdin")
	cmd.Flags().StringVar(&inline, "body", "", "Inline JSON payload")
	return cmd
}
