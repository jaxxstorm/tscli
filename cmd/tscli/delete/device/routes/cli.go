package routes

import (
	"fmt"

	"github.com/jaxxstorm/tscli/pkg/output"
	"github.com/jaxxstorm/tscli/pkg/tscli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Command() *cobra.Command {
	var deviceID string

	cmd := &cobra.Command{
		Use:   "routes",
		Short: "Delete all enabled subnet routes for a device",
		Long: `Delete all enabled subnet routes for a device.

This does not stop the device from advertising routes; change the device's
Tailscale configuration to do that instead.

Examples

  tscli delete device routes --device node-abc123`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := tscli.New()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			raw, err := tscli.ClearDeviceRoutesJSON(cmd.Context(), client, deviceID)
			if err != nil {
				return fmt.Errorf("failed to delete enabled subnet routes: %w", err)
			}

			return output.Print(viper.GetString("output"), raw)
		},
	}

	cmd.Flags().StringVar(&deviceID, "device", "", `Device ID whose enabled routes will be deleted (nodeId "node-abc123" or numeric id).`)
	_ = cmd.MarkFlagRequired("device")

	return cmd
}
