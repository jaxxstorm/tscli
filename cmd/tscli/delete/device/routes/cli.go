package routes

import (
	"fmt"

	"github.com/jaxxstorm/tscli/pkg/output"
	"github.com/jaxxstorm/tscli/pkg/tscli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Command() *cobra.Command {
	return &cobra.Command{
		Use:   "routes <device-id>",
		Short: "Delete all enabled subnet routes for a device",
		Long: `Delete all enabled subnet routes for a device.

This does not stop the device from advertising routes; change the device's
Tailscale configuration to do that instead.

Example:
  tscli delete device routes node-abc123`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := tscli.New()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			raw, err := tscli.ClearDeviceRoutesJSON(cmd.Context(), client, args[0])
			if err != nil {
				return fmt.Errorf("failed to delete enabled subnet routes: %w", err)
			}

			return output.Print(viper.GetString("output"), raw)
		},
	}
}
