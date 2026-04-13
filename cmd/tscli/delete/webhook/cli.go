// cmd/tscli/delete/webhook/delete.go
package webhook

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jaxxstorm/tscli/pkg/output"

	"github.com/jaxxstorm/tscli/pkg/tscli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Command() *cobra.Command {
	var hookID string

	cmd := &cobra.Command{
		Use:   "webhook",
		Short: "Delete a webhook",
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := tscli.New()
			if err != nil {
				return err
			}

			if err := client.Webhooks().Delete(context.Background(), hookID); err != nil {
				return err
			}

			out, _ := json.MarshalIndent(map[string]string{
				"result": fmt.Sprintf("webhook %s deleted", hookID),
			}, "", "  ")
			outputType := viper.GetString("output")
			output.Print(outputType, out)
			return nil
		},
	}

	cmd.Flags().StringVar(&hookID, "id", "", "Webhook ID to delete")
	_ = cmd.MarkFlagRequired("id")
	return cmd
}
