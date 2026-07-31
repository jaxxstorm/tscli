package maintenance

import (
	"github.com/jaxxstorm/tscli/cmd/tscli/maintenance/coveragegaps"
	"github.com/jaxxstorm/tscli/cmd/tscli/maintenance/openapirefresh"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	command := &cobra.Command{
		Use:   "maintenance",
		Short: "Run repository maintenance workflows",
		Long:  "Commands for running tscli repository maintenance workflows such as coverage-gap analysis and OpenAPI snapshot refresh.",
	}

	command.AddCommand(coveragegaps.Command())
	command.AddCommand(openapirefresh.Command())
	return command
}
