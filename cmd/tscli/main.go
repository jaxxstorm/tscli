// cmd/tscli/main.go

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/fang"
	configuration "github.com/jaxxstorm/tscli/cmd/tscli/config"
	"github.com/jaxxstorm/tscli/cmd/tscli/create"
	"github.com/jaxxstorm/tscli/cmd/tscli/delete"
	"github.com/jaxxstorm/tscli/cmd/tscli/get"
	"github.com/jaxxstorm/tscli/cmd/tscli/list"
	"github.com/jaxxstorm/tscli/cmd/tscli/set"
	"github.com/jaxxstorm/tscli/cmd/tscli/version"
	"github.com/jaxxstorm/tscli/pkg/config"
	"github.com/jaxxstorm/tscli/pkg/contract"
	"github.com/jaxxstorm/tscli/pkg/output"
	pkgversion "github.com/jaxxstorm/tscli/pkg/version"
	"github.com/spf13/cobra"
	viper "github.com/spf13/viper"
)

var (
	apiKey     string
	tailnet    string
	debug      bool
	outputType string
)

func configureCLI() *cobra.Command {
	config.Init()
	v := viper.GetViper() // use the global instance

	root := &cobra.Command{
		Use:  "tscli",
		Long: "A CLI tool for interacting with the Tailscale API.",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			// skip validation for "help", "version", and "config" commands
			if cmd.Name() == "help" || cmd.Name() == "version" || cmd.Name() == "config" ||
				(cmd.Parent() != nil && cmd.Parent().Name() == "config") ||
				(cmd.Parent() != nil && cmd.Parent().Parent() != nil && cmd.Parent().Parent().Name() == "config") {
				return nil
			}

			_ = v.BindPFlags(cmd.Flags())

			// Check for API key based on configuration mode
			var hasAPIKey bool
			if config.IsNewConfig() {
				tailnetConfig, err := config.GetActiveTailnetConfig()
				hasAPIKey = err == nil && tailnetConfig != nil && tailnetConfig.APIKey != ""
			} else {
				hasAPIKey = v.GetString("api-key") != ""
			}

			if !hasAPIKey {
				return fmt.Errorf("a Tailscale API key is required")
			}

			if v.GetString("tailnet") == "" {
				v.Set("tailnet", "-")
			}
			return nil
		},
	}

	root.AddCommand(
		get.Command(),
		list.Command(),
		delete.Command(),
		set.Command(),
		create.Command(),
		version.Command(),
		configuration.Command(),
	)

	root.PersistentFlags().StringVarP(&apiKey, "api-key", "k",
		"",
		"Tailscale API key")
	root.PersistentFlags().StringVarP(
		&outputType, "output", "o", "",
		fmt.Sprintf("Output: %v", output.Available()),
	)
	root.PersistentFlags().StringVarP(&tailnet, "tailnet", "n", v.GetString("tailnet"), "Tailscale tailnet")

	v.SetDefault("output", "json")

	v.AutomaticEnv()
	v.BindEnv("api-key", "TAILSCALE_API_KEY")
	v.BindEnv("tailnet", "TAILSCALE_TAILNET")
	v.BindEnv("output", "TSCLI_OUTPUT")
	v.BindPFlag("api-key", root.PersistentFlags().Lookup("api-key"))
	v.BindPFlag("tailnet", root.PersistentFlags().Lookup("tailnet"))
	v.BindPFlag("output", root.PersistentFlags().Lookup("output"))
	root.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Dump HTTP requests/responses")
	v.BindPFlag("debug", root.PersistentFlags().Lookup("debug"))
	v.BindEnv("debug", "TSCLI_DEBUG")

	return root
}

func main() {
	if err := fang.Execute(context.Background(), configureCLI(), fang.WithVersion(pkgversion.GetVersion())); err != nil {
		contract.IgnoreIoError(fmt.Fprintf(os.Stderr, "%v\n", err))
		os.Exit(1)
	}
}
