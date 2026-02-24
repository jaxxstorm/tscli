package upsert

import (
	"fmt"

	"github.com/jaxxstorm/tscli/pkg/config"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	var apiKey string

	cmd := &cobra.Command{
		Use:     "upsert <name>",
		Short:   "Create or update a tailnet profile",
		Long:    "Create or update a named tailnet profile with its API key and persist it in the config file.",
		Example: "tscli config profiles upsert _lbr_sandbox --api-key tskey-xxx",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			created, err := config.UpsertTailnetProfile(name, apiKey)
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

	cmd.Flags().StringVar(&apiKey, "api-key", "", "API key for the tailnet profile")
	_ = cmd.MarkFlagRequired("api-key")

	return cmd
}
