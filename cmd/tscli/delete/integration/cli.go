// cmd/tscli/delete/integration/cli.go
//
// `tscli delete posture-integration --id <integration-id>`
// Delete a single device-posture integration by its identifier.
package postureintegration

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
	var id string

	cmd := &cobra.Command{
		Use:   "posture-integration",
		Short: "Delete a posture integration",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := tscli.New()
			if err != nil {
				return err
			}

			var raw map[string]any
			if _, err := tscli.Do(
				context.Background(),
				client,
				http.MethodDelete,
				"/posture/integrations/"+id,
				nil,
				&raw,
			); err != nil {
				return fmt.Errorf("failed to delete posture integration %s: %w", id, err)
			}

			out, _ := json.MarshalIndent(raw, "", "  ")
			outputType := viper.GetString("output")
			output.Print(outputType, out)
			return nil
		},
	}

	cmd.Flags().StringVar(&id, "id", "", "Integration identifier")
	_ = cmd.MarkFlagRequired("id")

	return cmd
}
