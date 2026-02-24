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
	)

	cmd := &cobra.Command{
		Use:   "approval",
		Short: "Get service approval state for a device",
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			if serviceName == "" {
				return fmt.Errorf("--service is required")
			}
			if deviceID == "" {
				return fmt.Errorf("--device is required")
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
			var raw json.RawMessage
			if _, err := tscli.Do(
				context.Background(),
				client,
				http.MethodGet,
				path,
				nil,
				&raw,
			); err != nil {
				return fmt.Errorf("failed to get service approval: %w", err)
			}

			out, _ := json.MarshalIndent(raw, "", "  ")
			return output.Print(viper.GetString("output"), out)
		},
	}

	cmd.Flags().StringVar(&serviceName, "service", "", "Service name")
	cmd.Flags().StringVar(&deviceID, "device", "", "Device ID")
	return cmd
}
