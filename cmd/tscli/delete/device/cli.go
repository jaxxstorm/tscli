package device

import (
	"encoding/json"
	"fmt"

	"github.com/jaxxstorm/tscli/cmd/tscli/delete/device/invite"
	"github.com/jaxxstorm/tscli/cmd/tscli/delete/device/posture"
	"github.com/jaxxstorm/tscli/cmd/tscli/delete/device/routes"
	"github.com/jaxxstorm/tscli/pkg/output"

	"github.com/jaxxstorm/tscli/pkg/tscli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Command() *cobra.Command {
	command := &cobra.Command{
		Use:   "device",
		Short: "Delete device commands",
		Long:  "Delete devices from the Tailscale API",
		RunE: func(cmd *cobra.Command, args []string) error {

			client, err := tscli.New()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			deviceID, err := cmd.Flags().GetString("device")
			if err != nil {
				return fmt.Errorf("failed to get device flag: %w", err)
			}

			if err := client.Devices().Delete(cmd.Context(), deviceID); err != nil {
				return fmt.Errorf("failed to delete: %w", err)
			}

			payload := map[string]string{"result": fmt.Sprintf("device %s is deleted", deviceID)}

			out, err := json.MarshalIndent(payload, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal payload into JSON: %w", err)
			}
			outputType := viper.GetString("output")
			output.Print(outputType, out)
			return nil

		},
	}

	command.Flags().String("device", "", `Device ID to get (nodeId "node-abc123" or numeric id). Example: --device=node-abcdef123456`)
	_ = command.MarkFlagRequired("device")

	// Add subcommands
	command.AddCommand(invite.Command())
	command.AddCommand(posture.Command())
	command.AddCommand(routes.Command())

	return command
}
