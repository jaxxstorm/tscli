package config

import (
	"github.com/jaxxstorm/tscli/cmd/tscli/config/encryption"
	"github.com/jaxxstorm/tscli/cmd/tscli/config/get"
	"github.com/jaxxstorm/tscli/cmd/tscli/config/profiles"
	"github.com/jaxxstorm/tscli/cmd/tscli/config/set"
	"github.com/jaxxstorm/tscli/cmd/tscli/config/show"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	command := &cobra.Command{
		Use:   "config",
		Short: "Config commands",
		Long:  "Commands that show tscli's configuration",
	}

	command.AddCommand(show.Command())
	command.AddCommand(set.Command())
	command.AddCommand(get.Command())
	command.AddCommand(encryption.Command())
	command.AddCommand(profiles.Command())
	return command
}
