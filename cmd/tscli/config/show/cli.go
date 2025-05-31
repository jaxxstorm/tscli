package show

import (
	"encoding/json"

	"github.com/jaxxstorm/tscli/pkg/config"
	"github.com/jaxxstorm/tscli/pkg/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show tscli configuration",
		RunE: func(cmd *cobra.Command, _ []string) error {
			config.WarnIfOverrides(cmd.ErrOrStderr(), cmd)

			out, _ := json.MarshalIndent(viper.AllSettings(), "", "  ")
			return output.Print(viper.GetString("format"), out)
		},
	}
	return cmd
}

// helper for pretty printing flag names
func keys(m map[string]struct{}) []string {
	k := make([]string, 0, len(m))
	for n := range m {
		k = append(k, n)
	}
	return k
}
