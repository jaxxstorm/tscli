package service

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
	var serviceName string

	cmd := &cobra.Command{
		Use:   "service",
		Short: "Delete a service",
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
			if _, err := tscli.Do(
				context.Background(),
				client,
				http.MethodDelete,
				path,
				nil,
				nil,
			); err != nil {
				return fmt.Errorf("failed to delete service %s: %w", serviceName, err)
			}

			out, _ := json.MarshalIndent(map[string]string{"result": fmt.Sprintf("service %s deleted", serviceName)}, "", "  ")
			return output.Print(viper.GetString("output"), out)
		},
	}

	cmd.Flags().StringVar(&serviceName, "service", "", "Service name")
	return cmd
}
