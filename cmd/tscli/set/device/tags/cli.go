// cmd/tscli/set/device/tags/cli.go

package tags

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/jaxxstorm/tscli/pkg/output"

	"github.com/jaxxstorm/tscli/pkg/tscli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var tagRe = regexp.MustCompile(`^tag:[A-Za-z0-9_\-]+$`)

func Command() *cobra.Command {
	var (
		deviceID string
		tags     []string
	)

	cmd := &cobra.Command{
		Use:   "tags",
		Short: "Set a device's tags",
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			for _, t := range tags {
				if !tagRe.MatchString(t) {
					return fmt.Errorf("invalid tag %q: must match tag:<n> and contain only letters, numbers, dashes, or underscores", t)
				}
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := tscli.New()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			if err := client.Devices().SetTags(
				context.Background(),
				deviceID,
				tags,
			); err != nil {
				return fmt.Errorf("failed to set device tags: %w", err)
			}

			out, _ := json.MarshalIndent(map[string]any{
				"result": "tags updated",
				"device": deviceID,
				"tags":   tags,
			}, "", "  ")
			outputType := viper.GetString("output")
			output.Print(outputType, out)
			return nil
		},
	}

	cmd.Flags().StringVar(&deviceID, "device", "", "Device ID to retag")
	cmd.Flags().StringSliceVar(&tags, "tag", nil, "Tag to apply (repeatable), e.g. --tag tag:web --tag tag:prod")
	_ = cmd.MarkFlagRequired("device")
	_ = cmd.MarkFlagRequired("tag")

	return cmd
}
