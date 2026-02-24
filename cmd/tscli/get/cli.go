package get

import (
	"github.com/jaxxstorm/tscli/cmd/tscli/get/contacts"
	"github.com/jaxxstorm/tscli/cmd/tscli/get/device"
	"github.com/jaxxstorm/tscli/cmd/tscli/get/dns"
	postureintegration "github.com/jaxxstorm/tscli/cmd/tscli/get/integration"
	"github.com/jaxxstorm/tscli/cmd/tscli/get/key"
	"github.com/jaxxstorm/tscli/cmd/tscli/get/logs"
	"github.com/jaxxstorm/tscli/cmd/tscli/get/policy"
	"github.com/jaxxstorm/tscli/cmd/tscli/get/service"
	"github.com/jaxxstorm/tscli/cmd/tscli/get/settings"
	"github.com/jaxxstorm/tscli/cmd/tscli/get/user"
	"github.com/jaxxstorm/tscli/cmd/tscli/get/webhook"

	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	command := &cobra.Command{
		Use:   "get",
		Short: "Get commands",
		Long:  "Commands that get information from the Tailscale API",
	}

	command.AddCommand(device.Command())
	command.AddCommand(policy.Command())
	command.AddCommand(key.Command())
	command.AddCommand(postureintegration.Command())
	command.AddCommand(user.Command())
	command.AddCommand(webhook.Command())
	command.AddCommand(settings.Command())
	command.AddCommand(service.Command())
	command.AddCommand(contacts.Command())
	command.AddCommand(dns.Command())
	command.AddCommand(logs.Command())

	return command
}
