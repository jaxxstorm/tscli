// cmd/tscli/set/device/key/cli.go
//
// `tscli set device key --device <device-id> --disable-expiry`
// Updates device key settings via POST /api/v2/device/{id}/key

package key

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jaxxstorm/tscli/pkg/output"

	"github.com/jaxxstorm/tscli/pkg/tscli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type keyRequest struct {
	KeyExpiryDisabled bool `json:"keyExpiryDisabled"`
}

func Command() *cobra.Command {
	var (
		deviceID      string
		disableExpiry bool
		enableExpiry  bool
	)

	cmd := &cobra.Command{
		Use:   "key",
		Short: "Set device key settings",
		Long:  "Update device key settings such as expiry behavior",
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			if deviceID == "" {
				return fmt.Errorf("--device is required")
			}
			if disableExpiry && enableExpiry {
				return fmt.Errorf("cannot use both --disable-expiry and --enable-expiry")
			}
			if !disableExpiry && !enableExpiry {
				return fmt.Errorf("must specify either --disable-expiry or --enable-expiry")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := tscli.New()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			req := keyRequest{
				KeyExpiryDisabled: disableExpiry,
			}

			endpoint := fmt.Sprintf("/device/%s/key", deviceID)

			var response map[string]interface{}
			if _, err := tscli.Do(
				context.Background(),
				client,
				http.MethodPost,
				endpoint,
				req,
				&response,
			); err != nil {
				return fmt.Errorf("failed to update device key settings for %s: %w", deviceID, err)
			}

			// If the API returns data, show it, otherwise show confirmation
			if len(response) > 0 {
				out, _ := json.MarshalIndent(response, "", "  ")
				outputType := viper.GetString("output")
				output.Print(outputType, out)
			} else {
				// Fallback confirmation message
				action := "enabled"
				if disableExpiry {
					action = "disabled"
				}
				payload := map[string]string{
					"result":            fmt.Sprintf("device %s key expiry %s", deviceID, action),
					"device":            deviceID,
					"keyExpiryDisabled": fmt.Sprintf("%t", disableExpiry),
				}
				out, _ := json.MarshalIndent(payload, "", "  ")
				outputType := viper.GetString("output")
				output.Print(outputType, out)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&deviceID, "device", "", "Device ID to update key settings for")
	cmd.Flags().BoolVar(&disableExpiry, "disable-expiry", false, "Disable key expiry for this device")
	cmd.Flags().BoolVar(&enableExpiry, "enable-expiry", false, "Enable key expiry for this device")
	_ = cmd.MarkFlagRequired("device")

	return cmd
}
