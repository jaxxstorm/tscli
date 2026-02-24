package webhook

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jaxxstorm/tscli/cmd/tscli/set/webhook/test"
	"github.com/jaxxstorm/tscli/pkg/output"
	"github.com/jaxxstorm/tscli/pkg/tscli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Command() *cobra.Command {
	var (
		hookID        string
		subscriptions []string
		rotate        bool
	)

	cmd := &cobra.Command{
		Use:   "webhook",
		Short: "Update or rotate a webhook",
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			if hookID == "" {
				return fmt.Errorf("--id is required")
			}

			hasSubs := len(subscriptions) > 0
			if hasSubs == rotate {
				return fmt.Errorf("set either --subscription (one or more) or --rotate")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := tscli.New()
			if err != nil {
				return err
			}

			if rotate {
				var rotateResp json.RawMessage
				if _, err := tscli.Do(
					context.Background(),
					client,
					http.MethodPost,
					fmt.Sprintf("/webhooks/%s/rotate", hookID),
					nil,
					&rotateResp,
				); err != nil {
					return fmt.Errorf("failed to rotate webhook secret: %w", err)
				}
				out, _ := json.MarshalIndent(rotateResp, "", "  ")
				return output.Print(viper.GetString("output"), out)
			}

			payload := map[string]any{"subscriptions": subscriptions}
			var updateResp json.RawMessage
			if _, err := tscli.Do(
				context.Background(),
				client,
				http.MethodPatch,
				fmt.Sprintf("/webhooks/%s", hookID),
				payload,
				&updateResp,
			); err != nil {
				return fmt.Errorf("failed to update webhook subscriptions: %w", err)
			}

			out, _ := json.MarshalIndent(updateResp, "", "  ")
			return output.Print(viper.GetString("output"), out)
		},
	}

	cmd.Flags().StringVar(&hookID, "id", "", "Webhook ID to update")
	cmd.Flags().StringSliceVar(&subscriptions, "subscription", nil, "Webhook subscriptions to set (repeatable)")
	cmd.Flags().BoolVar(&rotate, "rotate", false, "Rotate the webhook secret")
	cmd.AddCommand(test.Command())
	return cmd
}
