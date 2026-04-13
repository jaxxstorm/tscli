// cmd/tscli/get/policy-preview/cli.go
//
// Preview which ACL rules match a given user *or* dst IP:port without
// saving the policy file.
//
//	# current tailnet policy against alice@example.com
//	tscli get policy-preview --type user   --value alice@example.com --current
//
//	# draft file against 10.0.0.10:443
//	tscli get policy-preview --type ipport --value 10.0.0.10:443 --file ./draft.hujson
package preview

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jaxxstorm/tscli/pkg/output"
	"github.com/spf13/viper"
	"net/http"
	"net/url"
	"strings"

	f "github.com/jaxxstorm/tscli/pkg/file"
	"github.com/jaxxstorm/tscli/pkg/tscli"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	var (
		kind    string
		value   string
		file    string
		body    string
		useLive bool
	)

	cmd := &cobra.Command{
		Use:   "preview",
		Short: "Preview which ACL rules match a user or IP:port",
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			kind = strings.ToLower(kind)
			if kind != "user" && kind != "ipport" {
				return errors.New(`--type must be "user" or "ipport"`)
			}
			if value == "" {
				return errors.New("--value is required")
			}

			// exactly one source
			sources := 0
			if useLive {
				sources++
			}
			if file != "" {
				sources++
			}
			if body != "" {
				sources++
			}
			if sources == 0 {
				return errors.New("provide ACL via --current, --file or --body")
			}
			if sources > 1 {
				return errors.New("--current, --file and --body are mutually exclusive")
			}
			return nil
		},

		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := tscli.New()
			if err != nil {
				return err
			}

			var raw []byte
			switch {
			case useLive:
				pf, err := client.PolicyFile().Raw(cmd.Context())
				if err != nil {
					return fmt.Errorf("fetch current policy: %w", err)
				}
				raw = []byte(pf.HuJSON)

			case file != "":
				raw, err = f.ReadInput(file, "")
				if err != nil {
					return err
				}

			default: // body
				raw = []byte(body)
			}

			// validate only for local inputs
			if !useLive {
				if err := f.ValidatePolicy(raw); err != nil {
					return fmt.Errorf("invalid policy file: %w", err)
				}
			}

			path := fmt.Sprintf(
				"/tailnet/{tailnet}/acl/preview?type=%s&previewFor=%s",
				kind, urlQueryEscape(value),
			)

			var resp json.RawMessage
			if _, err := tscli.Do(
				context.Background(),
				client,
				http.MethodPost,
				path,
				raw,
				&resp,
			); err != nil {
				return err
			}

			out, _ := json.MarshalIndent(resp, "", "  ")
			outputType := viper.GetString("output")
			output.Print(outputType, out)
			return nil
		},
	}

	cmd.Flags().StringVar(&kind, "type", "", "user | ipport (required)")
	cmd.Flags().StringVar(&value, "value", "", "email or \"ip:port\" (required)")
	cmd.Flags().StringVar(&file, "file", "", "path, file://path or '-' for stdin")
	cmd.Flags().StringVar(&body, "body", "", "inline ACL JSON/HuJSON")
	cmd.Flags().BoolVar(&useLive, "current", false, "use current tailnet ACL")

	_ = cmd.MarkFlagRequired("type")
	_ = cmd.MarkFlagRequired("value")

	return cmd
}

func urlQueryEscape(s string) string {
	return strings.ReplaceAll(url.QueryEscape(s), "+", "%20")
}
