package list

import (
	"github.com/jaxxstorm/tscli/cmd/tscli/list/devices"
	"github.com/jaxxstorm/tscli/cmd/tscli/list/integration"
	"github.com/jaxxstorm/tscli/cmd/tscli/list/invites"
	"github.com/jaxxstorm/tscli/cmd/tscli/list/keys"
	"github.com/jaxxstorm/tscli/cmd/tscli/list/logs"
	"github.com/jaxxstorm/tscli/cmd/tscli/list/nameservers"
	"github.com/jaxxstorm/tscli/cmd/tscli/list/routes"
	"github.com/jaxxstorm/tscli/cmd/tscli/list/services"
	"github.com/jaxxstorm/tscli/cmd/tscli/list/tailnets"
	"github.com/jaxxstorm/tscli/cmd/tscli/list/users"
	"github.com/jaxxstorm/tscli/cmd/tscli/list/webhooks"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	command := &cobra.Command{
		Use:   "list",
		Short: "List commands",
		Long:  "Commands that list information from the Tailscale API",
	}

	command.AddCommand(devices.Command())
	command.AddCommand(routes.Command())
	command.AddCommand(keys.Command())
	command.AddCommand(users.Command())
	command.AddCommand(invites.Command())
	command.AddCommand(webhooks.Command())
	command.AddCommand(postureintegration.Command())
	command.AddCommand(nameservers.Command())
	command.AddCommand(logs.Command())
	command.AddCommand(services.Command())
	command.AddCommand(tailnets.Command())

	return command
}
