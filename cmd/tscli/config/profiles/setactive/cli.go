package setactive

import (
	"fmt"

	"github.com/jaxxstorm/tscli/pkg/config"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	return &cobra.Command{
		Use:     "set-active <name>",
		Short:   "Set the active tailnet profile",
		Long:    "Set the active tailnet profile that runtime commands will use when no flag or environment override is provided.",
		Example: "tscli config profiles set-active _lbr_sandbox",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			if err := config.SetActiveTailnet(name); err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "active tailnet set to %s\n", name)
			return nil
		},
	}
}
