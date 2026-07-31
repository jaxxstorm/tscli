package openapirefresh

import (
	refresh "github.com/jaxxstorm/tscli/internal/maintenance/openapirefresh"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	opts := refresh.DefaultOptions()

	command := &cobra.Command{
		Use:   "openapi-refresh",
		Short: "Refresh the pinned Tailscale OpenAPI snapshot",
		Long:  "Fetch the Tailscale OpenAPI schema from the configured source URL and atomically refresh the pinned schema and snapshot metadata files.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return refresh.Run(opts)
		},
	}

	command.Flags().StringVar(&opts.SourceURL, "source-url", opts.SourceURL, "Canonical OpenAPI source URL")
	command.Flags().StringVar(&opts.SchemaOut, "schema-out", opts.SchemaOut, "Path for pinned OpenAPI schema")
	command.Flags().StringVar(&opts.MetadataOut, "metadata-out", opts.MetadataOut, "Path for snapshot metadata")

	return command
}
