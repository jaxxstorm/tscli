// cmd/tscli/get/dns/preferences/cli.go
//
// `tscli get prefs`
// Fetch the tailnet-wide *preferences* object (Auto-Update, MagicDNS,
// exit-node defaults, etc.) in JSON form.
package preferences

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

// Command registers `tscli get dns preferences`.
func Command() *cobra.Command {
	return &cobra.Command{
		Use:     "preferences",
		Aliases: []string{"prefs"},
		Short:   "Get tailnet DNS preferences",
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
				"/tailnet/{tailnet}/dns/preferences",
				nil,
				&raw,
			); err != nil {
				return fmt.Errorf("failed to fetch preferences: %w", err)
			}

			out, _ := json.MarshalIndent(raw, "", "  ")
			outputType := viper.GetString("output")
			output.Print(outputType, out)
			return nil
		},
	}
}
