// cmd/tscli/list/routes/cli.go
//
// Lists the subnet routes that a device advertises **and** those currently
// enabled for it.
//
// Example:
//
//	tscli list routes --device node-abcdef123456
package routes

import (
	"encoding/json"
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
		Short: "List a device's subnet routes",
		Long:  "Show both advertised and enabled subnet routes for a device.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if deviceID == "" {
				return fmt.Errorf("--device is required")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := tscli.New()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			routes, err := client.Devices().SubnetRoutes(cmd.Context(), deviceID)
			if err != nil {
				return fmt.Errorf("failed to list routes for device %s: %w", deviceID, err)
			}

			out, err := json.MarshalIndent(routes, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal routes into JSON: %w", err)
			}
			outputType := viper.GetString("output")
			output.Print(outputType, out)
			return nil
		},
	}

	// Flags -----------------------------------------------------------------
	cmd.Flags().StringVar(
		&deviceID,
		"device",
		"",
		`Device ID to inspect (nodeId "node-abc123" or numeric id). Example: --device node-abcdef123456`,
	)
	_ = cmd.MarkFlagRequired("device")

	return cmd
}
