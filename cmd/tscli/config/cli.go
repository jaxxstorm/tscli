package config

import (
	"github.com/jaxxstorm/tscli/cmd/tscli/config/get"
	"github.com/jaxxstorm/tscli/cmd/tscli/config/set"
	"github.com/jaxxstorm/tscli/cmd/tscli/config/show"
	"github.com/jaxxstorm/tscli/cmd/tscli/config/tailnet"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	command := &cobra.Command{
		Use:   "config",
		Short: "Config commands",
		Long:  "Commands that show tscli's configuration",
	}

	command.AddCommand(show.Command())
	command.AddCommand(set.Command())
	command.AddCommand(get.Command())
	command.AddCommand(tailnet.Command())
	return command
}
