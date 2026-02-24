// cmd/tscli/set/user/access/cli.go
//
// `tscli set user access --user <id> [--approve | --suspend | --restore]`
package access

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
		userID  string
		approve bool
		suspend bool
		restore bool
	)

	cmd := &cobra.Command{
		Use:   "access",
		Short: "Approve, suspend, or restore a user",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			actions := 0
			if approve {
				actions++
			}
			if suspend {
				actions++
			}
			if restore {
				actions++
			}
			if actions != 1 {
				return fmt.Errorf("exactly one of --approve, --suspend, or --restore must be set")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := tscli.New()
			if err != nil {
				return err
			}

			var path, msg string
			switch {
			case approve:
				path = "/users/" + userID + "/approve"
				msg = "approved"
			case suspend:
				path = "/users/" + userID + "/suspend"
				msg = "suspended"
			case restore:
				path = "/users/" + userID + "/restore"
				msg = "restored"
			}

			if _, err := tscli.Do(
				context.Background(),
				client,
				http.MethodPost,
				path,
				nil,
				nil,
			); err != nil {
				return err
			}

			out, _ := json.MarshalIndent(map[string]string{
				"result": fmt.Sprintf("user %s %s", userID, msg),
			}, "", "  ")
			outputType := viper.GetString("output")
			output.Print(outputType, out)
			return nil
		},
	}

	cmd.Flags().StringVar(&userID, "user", "", "User ID or email address")
	cmd.Flags().BoolVar(&approve, "approve", false, "Approve the user")
	cmd.Flags().BoolVar(&suspend, "suspend", false, "Suspend the user")
	cmd.Flags().BoolVar(&restore, "restore", false, "Restore the user")
	_ = cmd.MarkFlagRequired("user")

	return cmd
}
