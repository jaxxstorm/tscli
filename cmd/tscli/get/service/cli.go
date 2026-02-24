package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/jaxxstorm/tscli/cmd/tscli/get/service/approval"
	"github.com/jaxxstorm/tscli/pkg/output"
	"github.com/jaxxstorm/tscli/pkg/tscli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Command() *cobra.Command {
	var serviceName string

	cmd := &cobra.Command{
		Use:   "service",
		Short: "Get a service",
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			if serviceName == "" {
				return fmt.Errorf("--service is required")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := tscli.New()
			if err != nil {
				return err
			}

			path := fmt.Sprintf("/tailnet/{tailnet}/services/%s", url.PathEscape(serviceName))
			var raw json.RawMessage
			if _, err := tscli.Do(
				context.Background(),
				client,
				http.MethodGet,
				path,
				nil,
				&raw,
			); err != nil {
				return fmt.Errorf("failed to get service %s: %w", serviceName, err)
			}

			out, _ := json.MarshalIndent(raw, "", "  ")
			return output.Print(viper.GetString("output"), out)
		},
	}

	cmd.Flags().StringVar(&serviceName, "service", "", "Service name")
	cmd.AddCommand(approval.Command())
	return cmd
}
