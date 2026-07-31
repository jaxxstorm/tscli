package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jaxxstorm/tscli/tools/internal/coveragegaps"
)

func main() {
	opts := coveragegaps.DefaultOptions()

	flag.StringVar(&opts.OpenAPIPath, "openapi", opts.OpenAPIPath, "Path to pinned OpenAPI schema")
	flag.StringVar(&opts.MappingPath, "mapping", opts.MappingPath, "Path to command->operation map")
	flag.StringVar(&opts.ManifestPath, "manifest", opts.ManifestPath, "Path to command manifest")
	flag.StringVar(&opts.ExclusionsPath, "exclusions", opts.ExclusionsPath, "Path to exclusions policy")
	flag.StringVar(&opts.JSONOut, "json-out", opts.JSONOut, "Path for machine-readable report")
	flag.StringVar(&opts.MarkdownOut, "md-out", opts.MarkdownOut, "Path for markdown report")
	flag.StringVar(&opts.BaselinePath, "baseline", opts.BaselinePath, "Path to baseline report for diffing")
	flag.StringVar(&opts.DiffOut, "diff-out", opts.DiffOut, "Path for baseline diff report")
	flag.StringVar(&opts.PropertyCoveragePath, "property-coverage", opts.PropertyCoveragePath, "Path to property coverage manifest")
	flag.StringVar(&opts.PropertyExclusionsPath, "property-exclusions", opts.PropertyExclusionsPath, "Path to property exclusion policy")
	flag.BoolVar(&opts.FailOnRegression, "fail-on-regression", opts.FailOnRegression, "Exit non-zero if uncovered operations or unmapped commands regress vs baseline")
	flag.BoolVar(&opts.FailOnGaps, "fail-on-gaps", opts.FailOnGaps, "Exit non-zero if uncovered in-scope operations, unmapped commands, unknown mapped operations, or unknown mapped commands remain")
	flag.Parse()

	if err := coveragegaps.Run(opts); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
