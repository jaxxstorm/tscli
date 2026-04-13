// cmd/tscli/get/logs/stream/cli.go
//
// Fetch log-streaming configuration or (with --status) the current stream
// status for configuration-audit or network-flow logs.
//
//	# current log-streaming config for configuration-audit logs
//	tscli get logs stream --type configuration
//
//	# current status for network-flow log stream
//	tscli get logs stream --type network --status
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

var allowed = map[string]struct{}{
	string(tsapi.LogTypeConfig):  {},
	string(tsapi.LogTypeNetwork): {},
}

func Command() *cobra.Command {
	var (
		logType string
		status  bool
	)

	cmd := &cobra.Command{
		Use:   "stream --type {configuration|network} [--status]",
		Short: "Get log-streaming config or status",
		Long: `Without flags this returns the log-streaming configuration.
Use --status to fetch the current stream status.`,

		Args: cobra.NoArgs,

		PreRunE: func(_ *cobra.Command, _ []string) error {
			if _, ok := allowed[logType]; !ok {
				return fmt.Errorf(`--type must be "configuration" or "network"`)
			}
			return nil
		},

		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := tscli.New()
			if err != nil {
				return err
			}

			path := fmt.Sprintf(
				"/tailnet/{tailnet}/logging/%s/stream",
				url.PathEscape(logType),
			)
			if status {
				path += "/status"
			}

			var raw json.RawMessage
			if _, err = tscli.Do(
				context.Background(),
				client,
				http.MethodGet,
				path,
				nil,
				&raw,
			); err != nil {
				return err
			}

			return output.Print(viper.GetString("output"), raw)
		},
	}

	cmd.Flags().StringVar(&logType, "type", "",
		`Log type: "configuration" or "network" (required)`)
	cmd.Flags().BoolVar(&status, "status", false,
		"Fetch /status instead of /stream (shows runtime status)")
	_ = cmd.MarkFlagRequired("type")

	return cmd
}
