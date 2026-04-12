package init

import (
	"fmt"
	"strings"

	pkgagent "github.com/jaxxstorm/tscli/pkg/agent"
	"github.com/spf13/cobra"
)

func Command(root *cobra.Command) *cobra.Command {
	var dir string
	var tools []string
	var force bool

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize AI agent integrations for tscli",
		Long:  "Generate tscli-backed AI agent instructions, command catalogs, and native prompt, skill, or command surfaces for supported tools.",
		RunE: func(_ *cobra.Command, _ []string) error {
			result, err := pkgagent.Init(root, pkgagent.InstallOptions{
				RootDir: dir,
				Tools:   tools,
				Force:   force,
			})
			if err != nil {
				return err
			}

			fmt.Printf("tscli agent integrations initialized (%s) in %s\n", result.InstallScope, result.RootDir)
			fmt.Printf("tools: %s\n", strings.Join(result.Tools, ", "))
			fmt.Printf("indexed leaf commands: %d\n", result.CommandCount)
			fmt.Printf("manifest: %s\n", result.ManifestPath)
			if result.InstallScope == pkgagent.ScopeLocal {
				fmt.Printf("refresh with: tscli agent update --dir %s\n", result.RootDir)
			} else {
				fmt.Printf("refresh with: tscli agent update\n")
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&dir, "dir", "", "Optional repository root for a repo-local install; omit to install global user-level integrations")
	cmd.Flags().StringSliceVar(&tools, "tool", nil, fmt.Sprintf("Tool integrations to install (supported: %s; availability depends on global vs local install target)", strings.Join(pkgagent.SupportedTools(), ", ")))
	cmd.Flags().BoolVar(&force, "force", false, "Overwrite unmanaged files at generated target paths")
	return cmd
}
