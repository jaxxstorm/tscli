package update

import (
	"fmt"

	"github.com/jaxxstorm/tscli/pkg/config"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	return &cobra.Command{
		Use:   "update <name> <new-name> [new-api-key]",
		Short: "Update an existing tailnet configuration",
		Args:  cobra.RangeArgs(2, 3),
		RunE: func(_ *cobra.Command, args []string) error {
			oldName := args[0]
			newName := args[1]
			var newAPIKey string

			// Get current tailnets
			tailnets, active, err := config.ListTailnets()
			if err != nil {
				return fmt.Errorf("failed to list tailnets: %w", err)
			}

			// Find the tailnet to update
			var currentTailnet *config.TailnetConfig
			for _, tailnet := range tailnets {
				if tailnet.Name == oldName {
					currentTailnet = &tailnet
					break
				}
			}

			if currentTailnet == nil {
				return fmt.Errorf("tailnet %q not found", oldName)
			}

			// Use existing API key if not provided
			if len(args) == 3 {
				newAPIKey = args[2]
			} else {
				newAPIKey = currentTailnet.APIKey
			}

			// Remove old tailnet
			if err := config.RemoveTailnet(oldName); err != nil {
				return fmt.Errorf("failed to remove old tailnet: %w", err)
			}

			// Add new tailnet
			if err := config.AddTailnet(newName, newAPIKey); err != nil {
				return fmt.Errorf("failed to add updated tailnet: %w", err)
			}

			// If this was the active tailnet, make the new one active
			if active == oldName {
				if err := config.SetActiveTailnet(newName); err != nil {
					return fmt.Errorf("failed to set active tailnet: %w", err)
				}
			}

			fmt.Printf("Tailnet updated from %q to %q\n", oldName, newName)
			return nil
		},
	}
}
