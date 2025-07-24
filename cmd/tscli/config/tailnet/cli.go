package tailnet

import (
	"github.com/jaxxstorm/tscli/cmd/tscli/config/tailnet/add"
	"github.com/jaxxstorm/tscli/cmd/tscli/config/tailnet/list"
	"github.com/jaxxstorm/tscli/cmd/tscli/config/tailnet/remove"
	"github.com/jaxxstorm/tscli/cmd/tscli/config/tailnet/switch_"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	command := &cobra.Command{
		Use:   "tailnet",
		Short: "Manage tailnet configurations",
		Long:  "Commands to add, remove, list, and switch between tailnet configurations",
	}

	command.AddCommand(add.Command())
	command.AddCommand(list.Command())
	command.AddCommand(remove.Command())
	command.AddCommand(switch_.Command())

	return command
}
