// cmd/tscli/set/nameservers/cli.go
//
// `tscli set nameservers --nameserver 1.1.1.1 --nameserver 8.8.8.8`
// Replace the tailnet-wide DNS nameserver list.
//
// If you pass an empty slice (`--nameserver ""`) the custom list is removed
// and Tailscale falls back to its defaults.
package nameservers

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/jaxxstorm/tscli/pkg/output"

	"github.com/jaxxstorm/tscli/pkg/tscli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Command() *cobra.Command {
	var ns []string

	cmd := &cobra.Command{
		Use:     "nameservers",
		Aliases: []string{"ns"},
		Short:   "Set the DNS nameservers for the tailnet",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			if len(ns) == 0 {
				return fmt.Errorf("at least one --nameserver is required")
			}
			for _, ip := range ns {
				if net.ParseIP(ip) == nil {
					return fmt.Errorf("invalid IP address: %s", ip)
				}
			}
			return nil
		},

		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := tscli.New()
			if err != nil {
				return err
			}

			body := map[string][]string{"dns": ns}

			var resp json.RawMessage // <- receives the body untouched
			if _, err := tscli.Do(
				context.Background(),
				client,
				http.MethodPost,
				"/tailnet/{tailnet}/dns/nameservers",
				body,
				&resp,
			); err != nil {
				return fmt.Errorf("update failed: %w", err)
			}
			pretty, _ := json.MarshalIndent(resp, "", "  ")
			outputType := viper.GetString("output")
			output.Print(outputType, pretty)
			return nil
		},
	}

	cmd.Flags().StringSliceVarP(
		&ns,
		"nameserver", "N",
		nil,
		"DNS nameserver IP (repeatable). Example: --nameserver 1.1.1.1 --nameserver 8.8.8.8",
	)
	_ = cmd.MarkFlagRequired("nameserver")

	return cmd
}
