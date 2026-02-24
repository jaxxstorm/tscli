// cmd/tscli/set/device/cli.go

package device

import (
	"github.com/jaxxstorm/tscli/cmd/tscli/set/device/attributes"
	"github.com/jaxxstorm/tscli/cmd/tscli/set/device/authorization"
	"github.com/jaxxstorm/tscli/cmd/tscli/set/device/expiry"
	"github.com/jaxxstorm/tscli/cmd/tscli/set/device/invite"
	"github.com/jaxxstorm/tscli/cmd/tscli/set/device/ip"
	"github.com/jaxxstorm/tscli/cmd/tscli/set/device/key"
	"github.com/jaxxstorm/tscli/cmd/tscli/set/device/name"
	"github.com/jaxxstorm/tscli/cmd/tscli/set/device/posture"
	"github.com/jaxxstorm/tscli/cmd/tscli/set/device/routes"
	"github.com/jaxxstorm/tscli/cmd/tscli/set/device/tags"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	command := &cobra.Command{
		Use:   "device",
		Short: "Set device commands",
		Long:  "Commands that set device information on the Tailscale API",
	}

	command.AddCommand(attributes.Command())
	command.AddCommand(authorization.Command())
	command.AddCommand(expiry.Command())
	command.AddCommand(invite.Command())
	command.AddCommand(ip.Command())
	command.AddCommand(key.Command())
	command.AddCommand(name.Command())
	command.AddCommand(posture.Command())
	command.AddCommand(routes.Command())
	command.AddCommand(tags.Command())

	return command
}
