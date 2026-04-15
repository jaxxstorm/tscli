package upsert

import (
	"bufio"
	"fmt"
	"strings"

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
		Use:   "set <name>",
		Short: "Create or update a tailnet profile",
		Long:  "Create or update a named tailnet profile with either an API key or OAuth client credentials and persist it in the config file. Missing auth values are prompted for interactively.",
		Example: "tscli config profiles set _lbr_sandbox --api-key tskey-xxx\n" +
			"tscli config profiles set org-admin --oauth-client-id cid --oauth-client-secret secret\n" +
			"tscli config profiles set org-admin\n" +
			"tscli config profiles set sandbox --api-key tskey-xxx --profile-tailnet example.ts.net",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			reader := bufio.NewReader(cmd.InOrStdin())

			if err := promptForProfileAuth(cmd, reader, &apiKey, &oauthClientID, &oauthClientSecret); err != nil {
				return err
			}

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
	cmd.Flags().StringVar(&tailnet, "tailnet", "", "Deprecated alias for --profile-tailnet")
	_ = cmd.Flags().MarkDeprecated("tailnet", "use --profile-tailnet instead")
	cmd.Flags().StringVar(&apiKey, "api-key", "", "API key for the tailnet profile")
	cmd.Flags().StringVar(&oauthClientID, "oauth-client-id", "", "OAuth client ID for the tailnet profile")
	cmd.Flags().StringVar(&oauthClientSecret, "oauth-client-secret", "", "OAuth client secret for the tailnet profile")

	return cmd
}

func promptForProfileAuth(cmd *cobra.Command, reader *bufio.Reader, apiKey, oauthClientID, oauthClientSecret *string) error {
	hasAPIKey := strings.TrimSpace(*apiKey) != ""
	hasOAuthID := strings.TrimSpace(*oauthClientID) != ""
	hasOAuthSecret := strings.TrimSpace(*oauthClientSecret) != ""

	if !hasAPIKey && !hasOAuthID && !hasOAuthSecret {
		fmt.Fprint(cmd.OutOrStdout(), "Auth type [api-key|oauth]: ")
		value, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		switch strings.ToLower(strings.TrimSpace(value)) {
		case "api-key":
		case "oauth":
			if err := promptForValue(cmd, reader, "OAuth client ID: ", oauthClientID); err != nil {
				return err
			}
			if err := promptForValue(cmd, reader, "OAuth client secret: ", oauthClientSecret); err != nil {
				return err
			}
			return nil
		default:
			return fmt.Errorf("auth type must be one of: api-key, oauth")
		}
	}

	if strings.TrimSpace(*apiKey) == "" && strings.TrimSpace(*oauthClientID) == "" && strings.TrimSpace(*oauthClientSecret) == "" {
		return promptForValue(cmd, reader, "API key: ", apiKey)
	}
	if strings.TrimSpace(*apiKey) == "" && (strings.TrimSpace(*oauthClientID) == "" || strings.TrimSpace(*oauthClientSecret) == "") {
		if strings.TrimSpace(*oauthClientID) == "" {
			if err := promptForValue(cmd, reader, "OAuth client ID: ", oauthClientID); err != nil {
				return err
			}
		}
		if strings.TrimSpace(*oauthClientSecret) == "" {
			if err := promptForValue(cmd, reader, "OAuth client secret: ", oauthClientSecret); err != nil {
				return err
			}
		}
	}

	return nil
}

func promptForValue(cmd *cobra.Command, reader *bufio.Reader, prompt string, target *string) error {
	if strings.TrimSpace(*target) != "" {
		return nil
	}
	fmt.Fprint(cmd.OutOrStdout(), prompt)
	value, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	*target = strings.TrimSpace(value)
	return nil
}
