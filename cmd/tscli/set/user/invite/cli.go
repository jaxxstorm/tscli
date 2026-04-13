// cmd/tscli/set/user/invite/cli.go
//
// `tscli set user invite --id <invite-id> --resend`
// Resend a user invite via POST /api/v2/user-invites/{id}/resend

package invite

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
	var (
		inviteID string
		resend   bool
	)

	cmd := &cobra.Command{
		Use:   "invite",
		Short: "Set user invite status",
		Long:  "Resend a user invite",
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			if inviteID == "" {
				return fmt.Errorf("--id is required")
			}
			if !resend {
				return fmt.Errorf("--resend is required (currently the only supported action)")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := tscli.New()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			endpoint := fmt.Sprintf("/user-invites/%s/resend", inviteID)

			var response map[string]interface{}
			if _, err := tscli.Do(
				context.Background(),
				client,
				http.MethodPost,
				endpoint,
				nil,
				&response,
			); err != nil {
				return fmt.Errorf("failed to resend user invite: %w", err)
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

	cmd.Flags().StringVar(&inviteID, "id", "", "User invite ID")
	cmd.Flags().BoolVar(&resend, "resend", false, "Resend the invite")

	_ = cmd.MarkFlagRequired("id")
	_ = cmd.MarkFlagRequired("resend")

	return cmd
}
