package approval

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/jaxxstorm/tscli/pkg/output"
	"github.com/jaxxstorm/tscli/pkg/tscli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Command() *cobra.Command {
	var (
		serviceName string
		deviceID    string
		approved    bool
	)

	cmd := &cobra.Command{
		Use:   "approval",
		Short: "Set service approval state for a device",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			if serviceName == "" {
				return fmt.Errorf("--service is required")
			}
			if deviceID == "" {
				return fmt.Errorf("--device is required")
			}
			if !cmd.Flags().Lookup("approved").Changed {
				return fmt.Errorf("--approved is required")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := tscli.New()
			if err != nil {
				return err
			}

			path := fmt.Sprintf(
				"/tailnet/{tailnet}/services/%s/device/%s/approved",
				url.PathEscape(serviceName),
				url.PathEscape(deviceID),
			)
			payload := map[string]bool{"approved": approved}

			var resp json.RawMessage
			if _, err := tscli.Do(
				context.Background(),
				client,
				http.MethodPost,
				path,
				payload,
				&resp,
			); err != nil {
				return fmt.Errorf("failed to set service approval: %w", err)
			}

			out, _ := json.MarshalIndent(resp, "", "  ")
			return output.Print(viper.GetString("output"), out)
		},
	}

	cmd.Flags().StringVar(&serviceName, "service", "", "Service name")
	cmd.Flags().StringVar(&deviceID, "device", "", "Device ID")
	cmd.Flags().BoolVar(&approved, "approved", false, "Whether the service is approved for this device")
	return cmd
}
