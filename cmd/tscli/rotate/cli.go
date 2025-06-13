package rotate

import (
	"github.com/jaxxstorm/tscli/cmd/tscli/rotate/webhook"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	command := &cobra.Command{
		Use:   "rotate",
		Short: "Rotate commands",
		Long:  "Commands that rotate information from the Tailscale API",
	}

	command.AddCommand(webhook.Command())

	return command
}
