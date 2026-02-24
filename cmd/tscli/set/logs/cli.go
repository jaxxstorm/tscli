package logs

import (
	"github.com/jaxxstorm/tscli/cmd/tscli/set/logs/stream"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	command := &cobra.Command{
		Use:   "logs",
		Short: "Set log configuration",
		Long:  "Commands that mutate Tailscale logging settings.",
	}

	command.AddCommand(stream.Command())
	return command
}
