package update

import (
	"fmt"
	"strings"

	pkgagent "github.com/jaxxstorm/tscli/pkg/agent"
	"github.com/spf13/cobra"
)

func Command(root *cobra.Command) *cobra.Command {
	var dir string
	var force bool

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Refresh AI agent integrations for tscli",
		Long:  "Refresh the generated tscli agent bundle using the tool selection recorded in the target manifest. Without --dir, update refreshes the global user-level install.",
		RunE: func(_ *cobra.Command, _ []string) error {
			result, err := pkgagent.Update(root, pkgagent.UpdateOptions{
				RootDir: dir,
				Force:   force,
			})
			if err != nil {
				return err
			}

			fmt.Printf("tscli agent integrations updated (%s) in %s\n", result.InstallScope, result.RootDir)
			fmt.Printf("tools: %s\n", strings.Join(result.Tools, ", "))
			fmt.Printf("indexed leaf commands: %d\n", result.CommandCount)
			fmt.Printf("manifest: %s\n", result.ManifestPath)
			return nil
		},
	}

	cmd.Flags().StringVar(&dir, "dir", "", "Optional repository root containing a repo-local tscli agent manifest; omit to update the global install")
	cmd.Flags().BoolVar(&force, "force", false, "Overwrite unmanaged files at generated target paths")
	return cmd
}
