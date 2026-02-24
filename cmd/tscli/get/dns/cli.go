package dns

import (
	"github.com/jaxxstorm/tscli/cmd/tscli/get/dns/configuration"
	"github.com/jaxxstorm/tscli/cmd/tscli/get/dns/nameservers"
	"github.com/jaxxstorm/tscli/cmd/tscli/get/dns/preferences"
	"github.com/jaxxstorm/tscli/cmd/tscli/get/dns/searchpaths"
	"github.com/jaxxstorm/tscli/cmd/tscli/get/dns/split"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	command := &cobra.Command{
		Use:   "dns",
		Short: "Get DNS information",
		Long:  "Get information about DNS settings.",
	}

	command.AddCommand(configuration.Command())
	command.AddCommand(nameservers.Command())
	command.AddCommand(preferences.Command())
	command.AddCommand(searchpaths.Command())
	command.AddCommand(split.Command())

	return command
}
