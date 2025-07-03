// cmd/tscli/get/user/invite/cli.go
//
// `tscli get user invite --id <invite-id>`
// Fetch a single user invite by its id.
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
	var id string

	cmd := &cobra.Command{
		Use:   "invite",
		Short: "Get a user invite",
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
				http.MethodGet,
				"/tailnet/{tailnet}/user-invites/"+id,
				nil,
				&raw,
			); err != nil {
				return fmt.Errorf("failed to get user invite %s: %w", id, err)
			}

			out, _ := json.MarshalIndent(raw, "", "  ")
			outputType := viper.GetString("output")
			output.Print(outputType, out)
			return nil
		},
	}

	cmd.Flags().StringVar(&id, "id", "", "User invite ID to retrieve")
	_ = cmd.MarkFlagRequired("id")

	return cmd
}
