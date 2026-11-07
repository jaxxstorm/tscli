// cmd/tscli/get/key/cli.go
//
// `tscli get key --key <auth-key-id>`
// Fetch details for a single reusable auth-key.

package key

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/jaxxstorm/tscli/pkg/output"

	"github.com/jaxxstorm/tscli/pkg/apitype"
	"github.com/jaxxstorm/tscli/pkg/tscli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Command() *cobra.Command {
	var keyID string

	cmd := &cobra.Command{
		Use:   "key",
		Short: "Get a single tailnet auth key",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := tscli.New()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			var raw apitype.Key
			_, err = tscli.Do(context.Background(), client, http.MethodGet, "/tailnet/{tailnet}/keys/"+url.PathEscape(keyID), nil, &raw)
			if err != nil {
				return fmt.Errorf("failed to get key %s: %w", keyID, err)
			}

			out, err := json.MarshalIndent(raw, "", "  ")
			if err != nil {
				return fmt.Errorf("marshal key: %w", err)
			}
			outputType := viper.GetString("output")
			return output.Print(outputType, out)
		},
	}

	cmd.Flags().StringVar(&keyID, "key", "", "Key ID to retrieve")
	_ = cmd.MarkFlagRequired("key")

	return cmd
}
