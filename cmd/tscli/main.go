// cmd/tscli/main.go

package main

import (
	"fmt"
	"os"

	configuration "github.com/jaxxstorm/tscli/cmd/tscli/config"
	"github.com/jaxxstorm/tscli/cmd/tscli/create"
	"github.com/jaxxstorm/tscli/cmd/tscli/delete"
	"github.com/jaxxstorm/tscli/cmd/tscli/get"
	"github.com/jaxxstorm/tscli/cmd/tscli/list"
	"github.com/jaxxstorm/tscli/cmd/tscli/rotate"
	"github.com/jaxxstorm/tscli/cmd/tscli/set"
	"github.com/jaxxstorm/tscli/cmd/tscli/version"
	"github.com/jaxxstorm/tscli/pkg/config"
	"github.com/jaxxstorm/tscli/pkg/contract"
	"github.com/jaxxstorm/tscli/pkg/output"
	"github.com/spf13/cobra"
	viper "github.com/spf13/viper"
)

var (
	apiKey  string
	tailnet string
	debug   bool
	format  string
)

func configureCLI() *cobra.Command {
	config.Init()
	v := viper.GetViper() // use the global instance

	root := &cobra.Command{
		Use:  "tscli",
		Long: "A CLI tool for interacting with the Tailscale API.",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			// skip validation for "help" and "version" commands
			if cmd.Name() == "help" || cmd.Name() == "version" {
				return nil
			}

			_ = v.BindPFlags(cmd.Flags())
			if v.GetString("api-key") == "" {
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
		rotate.Command(),
	)

	root.PersistentFlags().StringVarP(&apiKey, "api-key", "k",
		"",
		"Tailscale API key")
	root.PersistentFlags().StringVarP(
		&format, "format", "f", "",
		fmt.Sprintf("Output format: %v", output.Available()),
	)
	root.PersistentFlags().StringVarP(&tailnet, "tailnet", "n", v.GetString("tailnet"), "Tailscale tailnet")

	v.SetDefault("format", "json")

	v.AutomaticEnv()
	v.BindEnv("api-key", "TAILSCALE_API_KEY")
	v.BindEnv("tailnet", "TAILSCALE_TAILNET")
	v.BindEnv("format", "TSCLI_FORMAT")
	v.BindPFlag("api-key", root.PersistentFlags().Lookup("api-key"))
	v.BindPFlag("tailnet", root.PersistentFlags().Lookup("tailnet"))
	v.BindPFlag("format", root.PersistentFlags().Lookup("format"))
	root.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Dump HTTP requests/responses")
	v.BindPFlag("debug", root.PersistentFlags().Lookup("debug"))
	v.BindEnv("debug", "TSCLI_DEBUG")

	return root
}

func main() {
	if err := configureCLI().Execute(); err != nil {
		contract.IgnoreIoError(fmt.Fprintf(os.Stderr, "%v\n", err))
		os.Exit(1)
	}
}
