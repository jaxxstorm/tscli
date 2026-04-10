// cmd/tscli/get/settings/cli.go
//
// `tscli get settings`
// Fetch the tailnet-wide settings object and print it as JSON.
package settings

import (
	"fmt"

	"github.com/jaxxstorm/tscli/pkg/output"
	"github.com/jaxxstorm/tscli/pkg/tscli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Command() *cobra.Command {
	return &cobra.Command{
		Use:   "settings",
		Short: "Get tailnet-wide settings",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := tscli.New()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			raw, err := tscli.GetTailnetSettingsJSON(cmd.Context(), client)
			if err != nil {
				return fmt.Errorf("failed to retrieve settings: %w", err)
			}

			outputType := viper.GetString("output")
			return output.Print(outputType, raw)
		},
	}
}
