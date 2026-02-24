package stream

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
	tsapi "tailscale.com/client/tailscale/v2"
)

var validLogTypes = map[string]struct{}{
	string(tsapi.LogTypeConfig):  {},
	string(tsapi.LogTypeNetwork): {},
}

func Command() *cobra.Command {
	var logType string

	cmd := &cobra.Command{
		Use:   "stream",
		Short: "Disable log streaming",
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			if _, ok := validLogTypes[logType]; !ok {
				return fmt.Errorf(`--type must be "configuration" or "network"`)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := tscli.New()
			if err != nil {
				return err
			}

			path := fmt.Sprintf("/tailnet/{tailnet}/logging/%s/stream", url.PathEscape(logType))
			if _, err := tscli.Do(
				context.Background(),
				client,
				http.MethodDelete,
				path,
				nil,
				nil,
			); err != nil {
				return fmt.Errorf("failed to disable log streaming: %w", err)
			}

			out, _ := json.MarshalIndent(map[string]string{"result": fmt.Sprintf("%s log streaming disabled", logType)}, "", "  ")
			return output.Print(viper.GetString("output"), out)
		},
	}

	cmd.Flags().StringVar(&logType, "type", "", `Log type: "configuration" or "network" (required)`)
	_ = cmd.MarkFlagRequired("type")
	return cmd
}
