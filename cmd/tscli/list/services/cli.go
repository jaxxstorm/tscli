package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jaxxstorm/tscli/cmd/tscli/list/services/devices"
	"github.com/jaxxstorm/tscli/pkg/output"
	"github.com/jaxxstorm/tscli/pkg/tscli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "services",
		Short: "List tailnet services",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := tscli.New()
			if err != nil {
				return err
			}

			var raw json.RawMessage
			if _, err := tscli.Do(
				context.Background(),
				client,
				http.MethodGet,
				"/tailnet/{tailnet}/services",
				nil,
				&raw,
			); err != nil {
				return fmt.Errorf("failed to list services: %w", err)
			}

			out, _ := json.MarshalIndent(raw, "", "  ")
			return output.Print(viper.GetString("output"), out)
		},
	}

	cmd.AddCommand(devices.Command())
	return cmd
}
