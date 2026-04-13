// cmd/tscli/set/device/ip/cli.go
//
// `tscli set device ip --device <id> --ip 100.64.0.42`
package ip

import (
	"context"
	"encoding/json"

	"fmt"
	"net"

	"github.com/jaxxstorm/tscli/pkg/output"

	"github.com/jaxxstorm/tscli/pkg/tscli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Command() *cobra.Command {
	var (
		deviceID string
		ipv4     string
	)

	cmd := &cobra.Command{
		Use:   "ip",
		Short: "Set a device's Tailscale IPv4 address",
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			ip := net.ParseIP(ipv4)
			if ip == nil || ip.To4() == nil {
				return fmt.Errorf("invalid IPv4 address: %s", ipv4)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := tscli.New()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			if err := client.Devices().SetIPv4Address(
				context.Background(),
				deviceID,
				ipv4,
			); err != nil {
				return fmt.Errorf("failed to set device IPv4: %w", err)
			}

			out, _ := json.MarshalIndent(map[string]string{
				"result": fmt.Sprintf("device %s: IPv4 set to %s", deviceID, ipv4),
			}, "", "  ")

			outputType := viper.GetString("output")
			output.Print(outputType, out)
			return nil
		},
	}

	cmd.Flags().StringVar(&deviceID, "device", "", "Device ID")
	cmd.Flags().StringVar(&ipv4, "ip", "", "IPv4 address to assign")

	_ = cmd.MarkFlagRequired("device")
	_ = cmd.MarkFlagRequired("ip")

	return cmd
}
