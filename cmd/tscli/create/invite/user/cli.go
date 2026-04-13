// cmd/tscli/create/invite/user/cli.go
//
// Send an *invite* e-mail so a new user can join your tailnet.
//
//	# invite alice as a member (default role)
//	tscli create invite user --email alice@example.com
//
//	# invite bob as an IT-admin
//	tscli create invite user --email bob@example.com --role it-admin
package user

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
	"strings"

	"github.com/jaxxstorm/tscli/pkg/output"

	"github.com/jaxxstorm/tscli/pkg/tscli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var allowedRole = map[string]struct{}{
	"owner":         {},
	"member":        {},
	"admin":         {},
	"it-admin":      {},
	"network-admin": {},
	"billing-admin": {},
	"auditor":       {},
}

// Command registers `tscli create invite user`.
func Command() *cobra.Command {
	var (
		email string
		role  string
	)

	cmd := &cobra.Command{
		Use:   "user",
		Short: "Create a user invite",
		Long:  "Invite a user to the current Tailscale tailnet with an optional role.",
		PreRunE: func(*cobra.Command, []string) error {
			if _, err := mail.ParseAddress(email); err != nil {
				return fmt.Errorf("invalid email: %w", err)
			}
			role = strings.ToLower(role)
			if role == "" {
				role = "member"
			}
			if _, ok := allowedRole[role]; !ok {
				return fmt.Errorf("invalid --role %q (owner|member|admin|it-admin|network-admin|billing-admin|auditor)", role)
			}
			return nil
		},

		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := tscli.New()
			if err != nil {
				return err
			}

			body := []map[string]any{{
				"email": email,
				"role":  role,
			}}

			var resp json.RawMessage
			if _, err := tscli.Do(
				context.Background(),
				client,
				http.MethodPost,
				"/tailnet/{tailnet}/user-invites",
				body,
				&resp,
			); err != nil {
				return fmt.Errorf("invite failed: %w", err)
			}

			out, _ := json.MarshalIndent(resp, "", "  ")
			outputType := viper.GetString("output")
			output.Print(outputType, out)
			return nil
		},
	}

	cmd.Flags().StringVar(&email, "email", "", "E-mail address to invite (required)")
	cmd.Flags().StringVar(&role, "role", "member",
		"User role: owner|member|admin|it-admin|network-admin|billing-admin|auditor")

	_ = cmd.MarkFlagRequired("email")
	return cmd
}
