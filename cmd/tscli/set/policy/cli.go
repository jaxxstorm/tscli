// cmd/tscli/set/policy/cli.go
package policy

import (
	"errors"
	"fmt"
	"os"

	f "github.com/jaxxstorm/tscli/pkg/file"
	"github.com/jaxxstorm/tscli/pkg/tscli"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	var (
		filePath string
		inline   string
		ifMatch  string
	)

	cmd := &cobra.Command{
		Use:   "policy",
		Short: "Replace the tailnet ACL / policy file",
		Long: `Accepts JSON or HuJSON.

  --file  file://path.hujson  (or "-": stdin)
  --body  '{"ACLs":[...]}'   (inline)

  Optional optimistic lock:
  --if-match ts-default  OR  value from previous GET /acl ETag`,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			if filePath == "" && inline == "" {
				return errors.New("one of --file or --body is required")
			}
			if filePath != "" && inline != "" {
				return errors.New("--file and --body are mutually exclusive")
			}
			return nil
		},

		RunE: func(cmd *cobra.Command, _ []string) error {
			raw, err := f.ReadInput(filePath, inline)
			if err != nil {
				return err
			}

			if err := f.ValidatePolicy(raw); err != nil {
				return fmt.Errorf("invalid ACL: %w", err)
			}

			client, err := tscli.New()
			if err != nil {
				return err
			}

			if err := client.PolicyFile().Set(cmd.Context(), string(raw), ifMatch); err != nil {
				return fmt.Errorf("update failed: %w", err)
			}

			updated, err := client.PolicyFile().Raw(cmd.Context())
			if err != nil {
				return fmt.Errorf("fetch updated file: %w", err)
			}
			fmt.Fprintln(os.Stdout, updated.HuJSON) // unchanged formatting
			return nil
		},
	}

	cmd.Flags().StringVar(&filePath, "file", "", "file://path, local path, or '-' for stdin")
	cmd.Flags().StringVar(&inline, "body", "", "Raw ACL JSON/HuJSON string")
	cmd.Flags().StringVar(&ifMatch, "if-match", "", "Value for If-Match header (etag or ts-default)")

	return cmd
}
