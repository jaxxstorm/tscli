package profiles

import (
	"github.com/jaxxstorm/tscli/cmd/tscli/config/profiles/deleteprofile"
	"github.com/jaxxstorm/tscli/cmd/tscli/config/profiles/list"
	"github.com/jaxxstorm/tscli/cmd/tscli/config/profiles/setactive"
	"github.com/jaxxstorm/tscli/cmd/tscli/config/profiles/upsert"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	command := &cobra.Command{
		Use:   "profiles",
		Short: "Manage tailnet profiles",
		Long:  "Commands to list, switch, upsert, and delete multi-tailnet configuration profiles.",
	}

	command.AddCommand(
		list.Command(),
		setactive.Command(),
		upsert.Command(),
		deleteprofile.Command(),
	)

	return command
}
