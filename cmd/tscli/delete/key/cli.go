package key

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
	var keyID string

	cmd := &cobra.Command{
		Use:   "key",
		Short: "Delete a tailnet auth key",
		Long:  "Remove an auth key from the tailnet",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := tscli.New()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			if _, err := tscli.Do(
				context.Background(),
				client,
				http.MethodDelete,
				"/tailnet/{tailnet}/keys/"+keyID,
				nil,
				nil,
			); err != nil {
				return fmt.Errorf("failed to delete key %s: %w", keyID, err)
			}

			resp := map[string]string{
				"result": fmt.Sprintf("key %s deleted", keyID),
			}
			out, _ := json.MarshalIndent(resp, "", "  ")
			outputType := viper.GetString("output")
			output.Print(outputType, out)
			return nil
		},
	}

	cmd.Flags().StringVar(&keyID, "key", "", "Key ID to delete")
	_ = cmd.MarkFlagRequired("key")

	return cmd
}
