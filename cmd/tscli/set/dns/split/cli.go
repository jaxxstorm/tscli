// cmd/tscli/set/splitdns/cli.go
//
// Replace *or* patch the split-DNS mapping.
//
//	# add two nameservers for example.com, one DoH endpoint for other.com
//	tscli set splitdns \
//	   --entry example.com=1.1.1.1 \
//	   --entry example.com=8.8.8.8 \
//	   --entry other.com=https://dns.google/dns-query
//
//	# clear a single domain (entry with empty RHS)
//	tscli set splitdns --entry stale.com=
//
//	# replace the whole mapping (PUT, drop everything else)
//	tscli set splitdns --replace --entry corp.local=10.0.0.53
package split

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	cldns "github.com/jaxxstorm/tscli/pkg/dns"
	"github.com/jaxxstorm/tscli/pkg/output"

	"github.com/jaxxstorm/tscli/pkg/tscli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var domainRE = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9\.\-]*\.[a-zA-Z]{2,}$`)

// Command registers `tscli set dns split-dns`.
func Command() *cobra.Command {
	var (
		entries []string // domain=ip (repeatable, one per IP)
		replace bool
	)

	cmd := &cobra.Command{
		Use:     "split-dns",
		Aliases: []string{"splitdns", "split", "sd"},
		Short:   "Set or patch split-DNS domains → nameservers",

		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			if !cmd.Flags().Lookup("entry").Changed {
				return fmt.Errorf("at least one --entry is required (domain=ip)")
			}
			return nil
		},

		RunE: func(cmd *cobra.Command, _ []string) error {
			/* ---------------------------------------------------------- */
			// Build payload map[string]any
			payload := make(map[string]any)
			tmp := make(map[string][]string)

			for _, e := range entries {
				parts := strings.SplitN(e, "=", 2)
				if len(parts) != 2 {
					return fmt.Errorf("invalid --entry %q (expect domain=ip)", e)
				}
				dom, ip := strings.ToLower(parts[0]), parts[1]

				if !domainRE.MatchString(dom) {
					return fmt.Errorf("invalid domain: %q", dom)
				}

				// empty RHS → clear this domain
				if ip == "" {
					payload[dom] = nil
					continue
				}

				if err := cldns.ValidateNameserver(ip); err != nil {
					return fmt.Errorf("invalid nameserver %q for %s", ip, dom)
				}
				tmp[dom] = append(tmp[dom], ip)
			}

			// merge temp slices unless domain already set to nil
			for d, ips := range tmp {
				if _, cleared := payload[d]; cleared {
					continue // explicit nil wins
				}
				payload[d] = ips
			}

			client, err := tscli.New()
			if err != nil {
				return err
			}

			method := http.MethodPatch
			if replace {
				method = http.MethodPut
			}

			var resp json.RawMessage
			if _, err := tscli.Do(
				context.Background(),
				client,
				method,
				"/tailnet/{tailnet}/dns/split-dns",
				payload,
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

	cmd.Flags().StringArrayVarP(
		&entries,
		"entry", "e",
		nil,
		`Mapping "domain=nameserver". Repeat --entry for multiple IP/DoH values or domains. Set an empty value to clear.`,
	)
	cmd.Flags().BoolVar(&replace, "replace", false,
		"Replace the entire mapping (PUT) instead of patching (PATCH).")

	return cmd
}
