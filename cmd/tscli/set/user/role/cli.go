// cmd/tscli/set/user/role/cli.go
package role

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

var validRoles = map[string]struct{}{
	"owner": {}, "member": {}, "admin": {}, "it-admin": {},
	"network-admin": {}, "billing-admin": {}, "auditor": {}, "all": {},
}

func Command() *cobra.Command {
	var (
		userID string
		role   string
	)

	cmd := &cobra.Command{
		Use:   "role",
		Short: "Update a user's role",
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			if role == "" {
				return fmt.Errorf("--role is required")
			}
			if _, ok := validRoles[strings.ToLower(role)]; !ok {
				return fmt.Errorf("invalid role %q", role)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := tscli.New()
			if err != nil {
				return err
			}

			_, err = tscli.Do(
				context.Background(),
				client,
				http.MethodPost,
				"/users/"+userID+"/role",
				map[string]string{"role": role},
				nil,
			)
			if err != nil {
				return err
			}

			out, _ := json.MarshalIndent(map[string]string{
				"result": fmt.Sprintf("user %s role set to %s", userID, role),
			}, "", "  ")
			outputType := viper.GetString("output")
			output.Print(outputType, out)
			return nil
		},
	}

	cmd.Flags().StringVar(&userID, "user", "", "User ID or email")
	cmd.Flags().StringVar(&role, "role", "", "New role (owner|member|admin|it-admin|network-admin|billing-admin|auditor|all)")
	_ = cmd.MarkFlagRequired("user")
	_ = cmd.MarkFlagRequired("role")
	return cmd
}
