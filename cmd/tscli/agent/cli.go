package agent

import (
	agentinit "github.com/jaxxstorm/tscli/cmd/tscli/agent/init"
	agentupdate "github.com/jaxxstorm/tscli/cmd/tscli/agent/update"
	"github.com/spf13/cobra"
)

func Command(root *cobra.Command) *cobra.Command {
	command := &cobra.Command{
		Use:   "agent",
		Short: "Manage AI agent integrations for tscli",
		Long:  "Generate and refresh tscli-backed AI agent instructions, skills, prompts, and commands for either a repo-local checkout or global user-level tooling.",
	}

	command.AddCommand(agentinit.Command(root))
	command.AddCommand(agentupdate.Command(root))
	return command
}
