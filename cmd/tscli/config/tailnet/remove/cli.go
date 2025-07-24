package remove

import (
	"fmt"

	"github.com/jaxxstorm/tscli/pkg/config"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	return &cobra.Command{
		Use:   "remove <name>",
		Short: "Remove a tailnet configuration",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			// Check if using legacy config
			if config.IsLegacyConfig() {
				return fmt.Errorf("cannot remove tailnets: you are using legacy configuration (single api-key). Use 'tscli config tailnet add' to start using multi-tailnet configuration")
			}

			name := args[0]

			if err := config.RemoveTailnet(name); err != nil {
				return fmt.Errorf("failed to remove tailnet: %w", err)
			}

			fmt.Printf("Tailnet %q removed successfully\n", name)
			return nil
		},
	}
}
