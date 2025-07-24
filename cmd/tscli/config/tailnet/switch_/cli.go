package switch_

import (
	"fmt"

	"github.com/jaxxstorm/tscli/pkg/config"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	return &cobra.Command{
		Use:   "switch <name>",
		Short: "Switch to a different tailnet configuration",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			// Check if using legacy config
			if config.IsLegacyConfig() {
				return fmt.Errorf("cannot switch tailnets: you are using legacy configuration (single api-key). Use 'tscli config tailnet add' to start using multi-tailnet configuration")
			}

			name := args[0]

			if err := config.SetActiveTailnet(name); err != nil {
				return fmt.Errorf("failed to switch tailnet: %w", err)
			}

			fmt.Printf("Switched to tailnet %q\n", name)
			return nil
		},
	}
}
