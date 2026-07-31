package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jaxxstorm/tscli/internal/maintenance/openapirefresh"
)

func main() {
	opts := openapirefresh.DefaultOptions()

	flag.StringVar(&opts.SourceURL, "source-url", opts.SourceURL, "Canonical OpenAPI source URL")
	flag.StringVar(&opts.SchemaOut, "schema-out", opts.SchemaOut, "Path for pinned OpenAPI schema")
	flag.StringVar(&opts.MetadataOut, "metadata-out", opts.MetadataOut, "Path for snapshot metadata")
	flag.Parse()

	if err := openapirefresh.Run(opts); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
