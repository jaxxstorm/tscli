package delete

import (
	"github.com/jaxxstorm/tscli/cmd/tscli/delete/device"
	"github.com/jaxxstorm/tscli/cmd/tscli/delete/devices"
	postureintegration "github.com/jaxxstorm/tscli/cmd/tscli/delete/integration"
	"github.com/jaxxstorm/tscli/cmd/tscli/delete/key"
	"github.com/jaxxstorm/tscli/cmd/tscli/delete/logs"
	"github.com/jaxxstorm/tscli/cmd/tscli/delete/service"
	"github.com/jaxxstorm/tscli/cmd/tscli/delete/tailnet"
	"github.com/jaxxstorm/tscli/cmd/tscli/delete/user"
	"github.com/jaxxstorm/tscli/cmd/tscli/delete/users"
	"github.com/jaxxstorm/tscli/cmd/tscli/delete/webhook"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	command := &cobra.Command{
		Use:   "delete",
		Short: "Delete commands",
		Long:  "Commands that delete from the Tailscale API",
	}

	command.AddCommand(device.Command())
	command.AddCommand(devices.Command())
	command.AddCommand(user.Command())
	command.AddCommand(users.Command())
	command.AddCommand(key.Command())
	command.AddCommand(service.Command())
	command.AddCommand(webhook.Command())
	command.AddCommand(postureintegration.Command())
	command.AddCommand(logs.Command())
	command.AddCommand(tailnet.Command())

	return command
}
