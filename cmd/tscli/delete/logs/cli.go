package logs

import (
	"github.com/jaxxstorm/tscli/cmd/tscli/delete/logs/stream"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	command := &cobra.Command{
		Use:   "logs",
		Short: "Delete log configuration",
		Long:  "Commands that remove Tailscale logging configuration.",
	}

	command.AddCommand(stream.Command())
	return command
}
