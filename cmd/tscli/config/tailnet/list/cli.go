package list

import (
	"encoding/json"
	"fmt"

	"github.com/jaxxstorm/tscli/pkg/config"
	"github.com/jaxxstorm/tscli/pkg/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Command() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all tailnet configurations",
		RunE: func(_ *cobra.Command, _ []string) error {
			// Check if using legacy config
			if config.IsLegacyConfig() {
				fmt.Println("Warning: You are using legacy configuration (single api-key).")
				fmt.Println("No tailnet configurations found. Use 'tscli config tailnet add' to start using multi-tailnet configuration.")
				return nil
			}

			tailnets, active, err := config.ListTailnets()
			if err != nil {
				return fmt.Errorf("failed to list tailnets: %w", err)
			}

			if len(tailnets) == 0 {
				fmt.Println("No tailnets configured")
				return nil
			}

			outputType := viper.GetString("output")
			if outputType == "json" || outputType == "yaml" {
				type tailnetDisplay struct {
					Name     string `json:"name" yaml:"name"`
					APIKey   string `json:"api-key" yaml:"api-key"`
					IsActive bool   `json:"is-active" yaml:"is-active"`
				}

				var displayTailnets []tailnetDisplay
				for _, tailnet := range tailnets {
					displayTailnets = append(displayTailnets, tailnetDisplay{
						Name:     tailnet.Name,
						APIKey:   tailnet.APIKey,
						IsActive: tailnet.Name == active,
					})
				}

				out, err := json.MarshalIndent(displayTailnets, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal output: %w", err)
				}

				return output.Print(outputType, out)
			}

			// Human-readable output
			fmt.Println("Configured tailnets:")
			for _, tailnet := range tailnets {
				marker := " "
				if tailnet.Name == active {
					marker = "*"
				}

				// Safely truncate API key
				apiKeyDisplay := tailnet.APIKey
				if len(apiKeyDisplay) > 8 {
					apiKeyDisplay = apiKeyDisplay[:8] + "..."
				} else if len(apiKeyDisplay) > 0 {
					apiKeyDisplay = apiKeyDisplay + "..."
				}

				fmt.Printf("%s %s (API Key: %s)\n", marker, tailnet.Name, apiKeyDisplay)
			}

			if active != "" {
				fmt.Printf("\nActive tailnet: %s\n", active)
			}

			return nil
		},
	}
}
