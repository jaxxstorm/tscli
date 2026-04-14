package create

import (
	"github.com/jaxxstorm/tscli/cmd/tscli/create/integration"
	"github.com/jaxxstorm/tscli/cmd/tscli/create/invite"
	"github.com/jaxxstorm/tscli/cmd/tscli/create/key"
	"github.com/jaxxstorm/tscli/cmd/tscli/create/tailnet"
	"github.com/jaxxstorm/tscli/cmd/tscli/create/token"
	"github.com/jaxxstorm/tscli/cmd/tscli/create/webhook"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	command := &cobra.Command{
		Use:   "create",
		Short: "Create commands",
		Long:  "Commands that create in the Tailscale API",
	}

	command.AddCommand(postureintegration.Command())
	command.AddCommand(key.Command())
	command.AddCommand(webhook.Command())
	command.AddCommand(invite.Command())
	command.AddCommand(token.Command())
	command.AddCommand(tailnet.Command())

	return command
}
