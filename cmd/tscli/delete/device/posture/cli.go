// cmd/tscli/delete/device/posture/cli.go
//
// `tscli delete device posture --device <id> --key custom:attr`
// Removes a custom posture attribute from a device.

package posture

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	"github.com/jaxxstorm/tscli/pkg/output"

	"github.com/jaxxstorm/tscli/pkg/tscli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// valid custom: key (same rules: prefix custom: and ≤50 chars total)
var keyRe = regexp.MustCompile(`^custom:[A-Za-z0-9_:]{1,43}$`) // 7+43 = 50

func Command() *cobra.Command {
	var (
		deviceID string
		attrKey  string
	)

	cmd := &cobra.Command{
		Use:   "posture",
		Short: "Delete a posture attribute",
		Long:  "Remove (unset) a custom:* posture attribute from the specified device.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if !keyRe.MatchString(attrKey) {
				return errors.New(`--key must start with "custom:" and contain only letters, numbers, underscores or colons (≤ 50 chars)`)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := tscli.New()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			if err := client.Devices().DeletePostureAttribute(
				context.Background(),
				deviceID,
				attrKey,
			); err != nil {
				return fmt.Errorf("failed to delete posture attribute: %w", err)
			}

			resp := map[string]string{
				"result": fmt.Sprintf("device %s: %s deleted", deviceID, attrKey),
			}
			out, _ := json.MarshalIndent(resp, "", "  ")
			outputType := viper.GetString("output")
			output.Print(outputType, out)
			return nil
		},
	}

	// ---------------- flags ----------------
	cmd.Flags().StringVar(
		&deviceID,
		"device",
		"",
		`Device ID whose attribute will be removed (nodeId "node-abc123" or numeric id). Example: --device=node-abcdef123456`,
	)
	cmd.Flags().StringVar(
		&attrKey,
		"key",
		"",
		`Posture attribute key to remove (must start with "custom:"). Example: --key=custom:buildNumber`,
	)

	_ = cmd.MarkFlagRequired("device")
	_ = cmd.MarkFlagRequired("key")

	return cmd
}
