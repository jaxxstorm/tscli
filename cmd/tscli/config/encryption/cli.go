package encryption

import (
	setup "github.com/jaxxstorm/tscli/cmd/tscli/config/encryption/setup"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "encryption",
		Short: "Manage config secret encryption",
		Long:  "Commands for setting up and managing AGE-based encryption for persisted config secrets.",
	}
	cmd.AddCommand(setup.Command())
	return cmd
}
