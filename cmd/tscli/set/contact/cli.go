// cmd/tscli/set/contacts/cli.go
package contact

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
	tsapi "tailscale.com/client/tailscale/v2"
)

var validTypes = map[string]tsapi.ContactType{
	"primary":  tsapi.ContactType("primary"),
	"billing":  tsapi.ContactType("billing"),
	"security": tsapi.ContactType("security"),
}

func Command() *cobra.Command {
	var (
		typeStr string
		email   string
		resend  bool
	)

	cmd := &cobra.Command{
		Use:   "contact",
		Short: "Update a contact record or resend verification email",
		Long:  "Update contact info or resend verification email. Allowed types: primary | billing | security. Use --resend to send verification email for a contact type.",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			if !resend && email == "" {
				return fmt.Errorf("--email is required when not using --resend")
			}
			if _, ok := validTypes[strings.ToLower(typeStr)]; !ok {
				return fmt.Errorf("invalid --type: %s (primary|billing|security)", typeStr)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := tscli.New()
			if err != nil {
				return err
			}

			if resend {
				// Resend verification email
				ctype := validTypes[strings.ToLower(typeStr)]
				endpoint := fmt.Sprintf("/tailnet/{tailnet}/contacts/%s/resend-verification-email", string(ctype))

				if _, err := tscli.Do(
					context.Background(),
					client,
					http.MethodPost,
					endpoint,
					nil, // no request body
					nil, // no response body expected
				); err != nil {
					return fmt.Errorf("failed to resend verification email for %s contact: %w", typeStr, err)
				}

				out, _ := json.MarshalIndent(map[string]string{
					"result": "verification email sent",
					"type":   strings.ToLower(typeStr),
				}, "", "  ")
				outputType := viper.GetString("output")
				output.Print(outputType, out)
				return nil
			}

			// Update contact email
			req := tsapi.UpdateContactRequest{
				Email: tsapi.PointerTo(email),
			}

			ctype := validTypes[strings.ToLower(typeStr)]
			if err := client.Contacts().Update(
				context.Background(),
				ctype,
				req,
			); err != nil {
				return fmt.Errorf("update failed: %w", err)
			}

			out, _ := json.MarshalIndent(map[string]string{
				"result": "contact updated",
				"type":   strings.ToLower(typeStr),
				"email":  email,
			}, "", "  ")
			outputType := viper.GetString("output")
			output.Print(outputType, out)
			return nil
		},
	}

	cmd.Flags().StringVar(&typeStr, "type", "primary", "Contact type to update (primary|billing|security)")
	cmd.Flags().StringVar(&email, "email", "", "New email address")
	cmd.Flags().BoolVar(&resend, "resend", false, "Resend verification email for the contact type")

	return cmd
}
