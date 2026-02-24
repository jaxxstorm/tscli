// cmd/tscli/delete/user/cli.go
//
// `tscli delete user --user <id-or-email>`
// Removes a user from the tailnet via POST /api/v2/users/{id}/delete

package user

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jaxxstorm/tscli/cmd/tscli/delete/user/invite"
	"github.com/jaxxstorm/tscli/pkg/output"

	"github.com/jaxxstorm/tscli/pkg/tscli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Command() *cobra.Command {
	var userID string

	cmd := &cobra.Command{
		Use:   "user",
		Short: "Delete a tailnet user",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := tscli.New()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			if _, err := tscli.Do(
				context.Background(),
				client,
				http.MethodPost,
				"/users/"+userID+"/delete",
				nil,
				nil,
			); err != nil {
				return fmt.Errorf("failed to delete user %s: %w", userID, err)
			}

			resp := map[string]string{
				"result": fmt.Sprintf("user %s deleted", userID),
			}
			out, _ := json.MarshalIndent(resp, "", "  ")
			outputType := viper.GetString("output")
			output.Print(outputType, out)
			return nil
		},
	}

	cmd.Flags().StringVar(&userID, "user", "", "User ID to delete")
	_ = cmd.MarkFlagRequired("user")

	// Add subcommands
	cmd.AddCommand(invite.Command())

	return cmd
}
