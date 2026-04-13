// cmd/tscli/set/dns/searchpaths/cli.go
//
// Replace the tailnet-wide DNS *search domains*.
//
//	# set two search paths
//	tscli set dns searchpaths --searchpath corp.example.com --searchpath svc.local
//
//	# clear (remove) custom search paths
//	tscli set searchpaths --searchpath ""
package searchpaths

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/jaxxstorm/tscli/pkg/output"

	"github.com/jaxxstorm/tscli/pkg/tscli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var domainRE = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9\.\-]*\.[a-zA-Z]{2,}$`)

// …imports unchanged…

func Command() *cobra.Command {
	var paths []string

	cmd := &cobra.Command{
		Use:     "searchpaths",
		Aliases: []string{"sp", "search"},
		Short:   "Set (replace) DNS search-paths for the tailnet",

		PreRunE: func(cmd *cobra.Command, _ []string) error {
			f := cmd.Flags().Lookup("searchpath")
			if !f.Changed {
				return fmt.Errorf("--searchpath is required (use --searchpath \"\" to clear)")
			}

			// No values at all → user passed `--searchpath ""` or `-s ""`
			if len(paths) == 0 {
				return nil // means “clear”
			}
			// Validate each domain
			for _, p := range paths {
				if p == "" { // explicit empty element also means clear
					continue
				}
				if !domainRE.MatchString(p) {
					return fmt.Errorf("invalid domain: %q", p)
				}
			}
			return nil
		},

		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := tscli.New()
			if err != nil {
				return err
			}

			// If the flag was given but zero/empty domains, clear the list.
			if len(paths) == 0 || (len(paths) == 1 && paths[0] == "") {
				paths = []string{}
			}

			body := map[string][]string{"searchPaths": paths}

			var resp json.RawMessage
			if _, err := tscli.Do(
				context.Background(),
				client,
				http.MethodPost,
				"/tailnet/{tailnet}/dns/searchpaths",
				body,
				&resp,
			); err != nil {
				return fmt.Errorf("update failed: %w", err)
			}

			out, _ := json.MarshalIndent(resp, "", "  ")
			outputType := viper.GetString("output")
			output.Print(outputType, out)
			return nil
		},
	}

	cmd.Flags().StringSliceVarP(
		&paths,
		"searchpath", "s",
		nil,
		"DNS search domain (repeatable). Use an empty value to clear, e.g. --searchpath \"\"",
	)

	return cmd
}
