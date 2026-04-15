package deleteprofile

import (
	"fmt"

	"github.com/jaxxstorm/tscli/pkg/config"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	return &cobra.Command{
		Use:     "delete <name>",
		Short:   "Delete a tailnet profile",
		Long:    "Delete a tailnet profile from configuration.",
		Example: "tscli config profiles delete _lbr_sandbox",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			if err := config.RemoveTailnetProfile(name); err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "tailnet profile %s removed\n", name)
			return nil
		},
	}
}
