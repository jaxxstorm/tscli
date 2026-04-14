package cli

import (
	"fmt"

	agentcmd "github.com/jaxxstorm/tscli/cmd/tscli/agent"
	configuration "github.com/jaxxstorm/tscli/cmd/tscli/config"
	"github.com/jaxxstorm/tscli/cmd/tscli/create"
	"github.com/jaxxstorm/tscli/cmd/tscli/delete"
	"github.com/jaxxstorm/tscli/cmd/tscli/get"
	"github.com/jaxxstorm/tscli/cmd/tscli/list"
	"github.com/jaxxstorm/tscli/cmd/tscli/set"
	"github.com/jaxxstorm/tscli/cmd/tscli/version"
	"github.com/jaxxstorm/tscli/pkg/config"
	"github.com/jaxxstorm/tscli/pkg/output"
	"github.com/spf13/cobra"
	viper "github.com/spf13/viper"
)

var (
	apiKey     string
	tailnet    string
	debug      bool
	outputType string
)

func Configure() *cobra.Command {
	config.Init()
	v := viper.GetViper()

	root := &cobra.Command{
		Use:  "tscli",
		Long: "A CLI tool for interacting with the Tailscale API.",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			if cmd.Name() == "help" || cmd.Name() == "version" || cmd.Name() == "completion" {
				return nil
			}
			if isLocalCommand(cmd) {
				return nil
			}
			if skipsAPIKeyPreRun(cmd) {
				return nil
			}

			_ = v.BindPFlags(cmd.Flags())

			if _, err := config.ResolveRuntimeConfig(config.ChangedMap(cmd)); err != nil {
				return err
			}
			return nil
		},
	}

	root.AddCommand(
		agentcmd.Command(root),
		get.Command(),
		list.Command(),
		delete.Command(),
		set.Command(),
		create.Command(),
		version.Command(),
		configuration.Command(),
	)

	root.PersistentFlags().StringVarP(&apiKey, "api-key", "k", "", "Tailscale API key")
	root.PersistentFlags().StringVarP(&outputType, "output", "o", "", fmt.Sprintf("Output: %v", output.Available()))
	root.PersistentFlags().StringVarP(&tailnet, "tailnet", "n", "-", "Tailscale tailnet")

	v.SetDefault("output", "json")

	v.AutomaticEnv()
	v.BindEnv("api-key", "TAILSCALE_API_KEY")
	v.BindEnv("tailnet", "TAILSCALE_TAILNET")
	v.BindEnv("output", "TSCLI_OUTPUT")
	v.BindEnv("base-url", "TSCLI_BASE_URL")
	v.BindPFlag("api-key", root.PersistentFlags().Lookup("api-key"))
	v.BindPFlag("tailnet", root.PersistentFlags().Lookup("tailnet"))
	v.BindPFlag("output", root.PersistentFlags().Lookup("output"))
	root.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Dump HTTP requests/responses")
	v.BindPFlag("debug", root.PersistentFlags().Lookup("debug"))
	v.BindEnv("debug", "TSCLI_DEBUG")

	return root
}

func isLocalCommand(cmd *cobra.Command) bool {
	for current := cmd; current != nil; current = current.Parent() {
		switch current.Name() {
		case "config", "agent":
			return true
		}
	}
	return false
}

func skipsAPIKeyPreRun(cmd *cobra.Command) bool {
	segments := make([]string, 0, 4)
	for current := cmd; current != nil; current = current.Parent() {
		if current.Name() == "tscli" {
			break
		}
		segments = append([]string{current.Name()}, segments...)
	}

	if len(segments) == 2 {
		return (segments[0] == "create" && segments[1] == "token") ||
			(segments[0] == "create" && segments[1] == "tailnet") ||
			(segments[0] == "list" && segments[1] == "tailnets") ||
			(segments[0] == "delete" && segments[1] == "tailnet")
	}

	return false
}
