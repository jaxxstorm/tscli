package add

import (
	"fmt"

	"github.com/jaxxstorm/tscli/pkg/config"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	return &cobra.Command{
		Use:   "add <name> <api-key>",
		Short: "Add a new tailnet configuration",
		Args:  cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			// Check if using legacy config and warn
			if config.IsLegacyConfig() {
				fmt.Println("Warning: You are currently using legacy configuration (single api-key).")
				fmt.Println("Adding a tailnet will switch you to the new multi-tailnet configuration.")
				fmt.Println("Your existing api-key will remain accessible via the legacy config until you migrate.")
				fmt.Println()
			}

			name, apiKey := args[0], args[1]

			if err := config.AddTailnet(name, apiKey); err != nil {
				return fmt.Errorf("failed to add tailnet: %w", err)
			}

			fmt.Printf("Tailnet %q added successfully\n", name)
			return nil
		},
	}
}
