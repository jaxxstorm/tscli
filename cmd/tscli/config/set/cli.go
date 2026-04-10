package set

import (
	"fmt"

	"github.com/jaxxstorm/tscli/pkg/config"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	return &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Persist a config value (~/.tscli.yaml)",
		Args:  cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			key, val := args[0], args[1]
			if err := config.SetPersistentValue(key, val); err != nil {
				return fmt.Errorf("write config: %w", err)
			}
			fmt.Printf("%s saved\n", key)
			return nil
		},
	}
}
