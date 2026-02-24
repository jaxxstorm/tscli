// cmd/tscli/set/device/invite/cli.go
//
// `tscli set device invite --id <invite-id|invite-url-or-code> --status <resend|accept>`
// Resend or accept a device invite.

package invite

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/jaxxstorm/tscli/pkg/output"

	"github.com/jaxxstorm/tscli/pkg/tscli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var validStatuses = map[string]string{
	"resend": "resend",
	"accept": "accept",
}

func Command() *cobra.Command {
	var (
		inviteID string
		status   string
	)

	cmd := &cobra.Command{
		Use:   "invite",
		Short: "Set device invite status",
		Long:  "Resend or accept a device invite. Valid statuses: resend, accept",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			if inviteID == "" {
				return fmt.Errorf("--id is required")
			}
			if status == "" {
				return fmt.Errorf("--status is required")
			}
			if _, ok := validStatuses[strings.ToLower(status)]; !ok {
				return fmt.Errorf("invalid --status: %s (resend|accept)", status)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := tscli.New()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			action := validStatuses[strings.ToLower(status)]
			endpoint := fmt.Sprintf("/device-invites/%s/resend", inviteID)
			var body any
			if action == "accept" {
				endpoint = "/device-invites/-/accept"
				body = map[string]string{"invite": inviteID}
			}

			var response map[string]interface{}
			if _, err := tscli.Do(
				context.Background(),
				client,
				http.MethodPost,
				endpoint,
				body,
				&response,
			); err != nil {
				return fmt.Errorf("failed to %s device invite: %w", action, err)
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

	cmd.Flags().StringVar(&inviteID, "id", "", "Invite ID, URL, or code (required)")
	cmd.Flags().StringVar(&status, "status", "", "Action to perform: resend or accept")

	_ = cmd.MarkFlagRequired("id")
	_ = cmd.MarkFlagRequired("status")

	return cmd
}
