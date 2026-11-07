// cmd/tscli/list/keys/cli.go
//
// `tscli list keys`
// List every reusable auth key in the tailnet.

package keys

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jaxxstorm/tscli/pkg/output"

	"github.com/jaxxstorm/tscli/pkg/apitype"
	"github.com/jaxxstorm/tscli/pkg/tscli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Command() *cobra.Command {

	var showAll bool

	cmd := &cobra.Command{
		Use:   "keys",
		Short: "List tailnet auth keys",
		Long:  "Return all reusable authentication keys defined for the tailnet.",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := tscli.New()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			var response struct {
				Keys []apitype.Key `json:"keys"`
			}
			path := "/tailnet/{tailnet}/keys"
			if showAll {
				path += "?all=true"
			}
			_, err = tscli.Do(context.Background(), client, http.MethodGet, path, nil, &response)
			if err != nil {
				return fmt.Errorf("failed to list keys: %w", err)
			}

			out, err := json.MarshalIndent(response.Keys, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal keys into JSON: %w", err)
			}
			outputType := viper.GetString("output")
			return output.Print(outputType, out)
		},
	}

	cmd.Flags().BoolVar(
		&showAll,
		"all",
		false,
		"Include all keys.",
	)

	return cmd
}
