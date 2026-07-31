package coveragegaps

import (
	gaps "github.com/jaxxstorm/tscli/internal/maintenance/coveragegaps"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	opts := gaps.DefaultOptions()
	check := false

	command := &cobra.Command{
		Use:   "coverage-gaps",
		Short: "Generate OpenAPI command coverage-gap reports",
		Long:  "Generate tscli OpenAPI command and property coverage-gap reports from the pinned schema, command map, and leaf command manifest.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if check {
				opts.FailOnRegression = true
				opts.FailOnGaps = true
			}
			return gaps.Run(opts)
		},
	}

	command.Flags().StringVar(&opts.OpenAPIPath, "openapi", opts.OpenAPIPath, "Path to pinned OpenAPI schema")
	command.Flags().StringVar(&opts.MappingPath, "mapping", opts.MappingPath, "Path to command-to-operation map")
	command.Flags().StringVar(&opts.ManifestPath, "manifest", opts.ManifestPath, "Path to leaf command manifest")
	command.Flags().StringVar(&opts.ExclusionsPath, "exclusions", opts.ExclusionsPath, "Path to operation and command exclusions policy")
	command.Flags().StringVar(&opts.JSONOut, "json-out", opts.JSONOut, "Path for machine-readable report")
	command.Flags().StringVar(&opts.MarkdownOut, "md-out", opts.MarkdownOut, "Path for markdown report")
	command.Flags().StringVar(&opts.BaselinePath, "baseline", opts.BaselinePath, "Path to baseline report for diffing")
	command.Flags().StringVar(&opts.DiffOut, "diff-out", opts.DiffOut, "Path for baseline diff report")
	command.Flags().StringVar(&opts.PropertyCoveragePath, "property-coverage", opts.PropertyCoveragePath, "Path to property coverage manifest")
	command.Flags().StringVar(&opts.PropertyExclusionsPath, "property-exclusions", opts.PropertyExclusionsPath, "Path to property exclusion policy")
	command.Flags().BoolVar(&opts.FailOnRegression, "fail-on-regression", opts.FailOnRegression, "Exit non-zero if uncovered operations or unmapped commands regress vs baseline")
	command.Flags().BoolVar(&opts.FailOnGaps, "fail-on-gaps", opts.FailOnGaps, "Exit non-zero if uncovered in-scope operations, unmapped commands, unknown mapped operations, or unknown mapped commands remain")
	command.Flags().BoolVar(&check, "check", check, "Enable strict baseline regression and gap failure checks")

	return command
}
