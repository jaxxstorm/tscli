package configuration

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

func Command() *cobra.Command {
	return &cobra.Command{
		Use:   "configuration",
		Short: "Get tailnet DNS configuration",
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
				"/tailnet/{tailnet}/dns/configuration",
				nil,
				&raw,
			); err != nil {
				return fmt.Errorf("failed to fetch DNS configuration: %w", err)
			}

			out, _ := json.MarshalIndent(raw, "", "  ")
			return output.Print(viper.GetString("output"), out)
		},
	}
}
