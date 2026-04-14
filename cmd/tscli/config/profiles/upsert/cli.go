package upsert

import (
	"fmt"

	"github.com/jaxxstorm/tscli/pkg/config"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	var (
		apiKey            string
		tailnet           string
		oauthClientID     string
		oauthClientSecret string
	)

	cmd := &cobra.Command{
		Use:   "upsert <name>",
		Short: "Create or update a tailnet profile",
		Long:  "Create or update a named tailnet profile with either an API key or OAuth client credentials and persist it in the config file.",
		Example: "tscli config profiles upsert _lbr_sandbox --api-key tskey-xxx\n" +
			"tscli config profiles upsert org-admin --oauth-client-id cid --oauth-client-secret secret\n" +
			"tscli config profiles upsert sandbox --api-key tskey-xxx --profile-tailnet example.ts.net",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			created, err := config.UpsertTailnetProfile(config.TailnetProfile{
				Name:              name,
				Tailnet:           tailnet,
				APIKey:            apiKey,
				OAuthClientID:     oauthClientID,
				OAuthClientSecret: oauthClientSecret,
			})
			if err != nil {
				return err
			}

			if created {
				fmt.Fprintf(cmd.OutOrStdout(), "tailnet profile %s created\n", name)
				return nil
			}

			fmt.Fprintf(cmd.OutOrStdout(), "tailnet profile %s updated\n", name)
			return nil
		},
	}

	cmd.Flags().StringVar(&tailnet, "profile-tailnet", "", "Explicit effective tailnet value for this profile (defaults to the profile name)")
	cmd.Flags().StringVar(&apiKey, "api-key", "", "API key for the tailnet profile")
	cmd.Flags().StringVar(&oauthClientID, "oauth-client-id", "", "OAuth client ID for the tailnet profile")
	cmd.Flags().StringVar(&oauthClientSecret, "oauth-client-secret", "", "OAuth client secret for the tailnet profile")

	return cmd
}
