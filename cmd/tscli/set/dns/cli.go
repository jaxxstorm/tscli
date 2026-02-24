package dns

import (
	"github.com/jaxxstorm/tscli/cmd/tscli/set/dns/configuration"
	"github.com/jaxxstorm/tscli/cmd/tscli/set/dns/nameservers"
	"github.com/jaxxstorm/tscli/cmd/tscli/set/dns/preferences"
	"github.com/jaxxstorm/tscli/cmd/tscli/set/dns/searchpaths"
	"github.com/jaxxstorm/tscli/cmd/tscli/set/dns/split"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	command := &cobra.Command{
		Use:   "dns",
		Short: "Set DNS information",
		Long:  "Set information about DNS settings.",
	}

	command.AddCommand(configuration.Command())
	command.AddCommand(nameservers.Command())
	command.AddCommand(preferences.Command())
	command.AddCommand(searchpaths.Command())
	command.AddCommand(split.Command())

	return command
}
