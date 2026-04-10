// cmd/tscli/list/devices/cli.go
//
// `tscli list devices [--all]`
// Prints every device in the tailnet.
//
// * With no flag it returns the “standard” fields the public API shows by
//   default.
// * With `--all` it requests every possible field (`?fields=all`).

package devices

import (
	"fmt"

	"github.com/jaxxstorm/tscli/pkg/output"
	"github.com/jaxxstorm/tscli/pkg/tscli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Command() *cobra.Command {
	var showAll bool

	cmd := &cobra.Command{
		Use:   "devices",
		Short: "List devices",
		Long: `List every device registered in your tailnet.

By default only the common fields are returned.
Use --all to include advanced fields such as ClientConnectivity, AdvertisedRoutes, and EnabledRoutes.

Structured output prints the API response body directly so documented response fields are preserved.

Examples

  # Standard view
  tscli list devices

  # Full view (all fields)
  tscli list devices --all
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := tscli.New()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			raw, err := tscli.ListDevicesJSON(cmd.Context(), client, showAll)
			if err != nil {
				return fmt.Errorf("failed to list devices: %w", err)
			}

			outputType := viper.GetString("output")
			return output.Print(outputType, raw)
		},
	}

	cmd.Flags().BoolVar(
		&showAll,
		"all",
		false,
		"Include every field returned by the API (equivalent to '?fields=all').",
	)

	return cmd
}
