// cmd/tscli/set/device/posture/cli.go

package posture

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/jaxxstorm/tscli/pkg/output"

	"github.com/jaxxstorm/tscli/pkg/tscli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tsapi "tailscale.com/client/tailscale/v2"
)

var keyRe = regexp.MustCompile(`^custom:[A-Za-z0-9_:]{1,43}$`) // 7 + 43 = 50
var strValRe = regexp.MustCompile(`^[A-Za-z0-9_.]{1,50}$`)

func coerceValue(in string) (any, error) {
	if b, err := strconv.ParseBool(in); err == nil {
		return b, nil
	}
	if i, err := strconv.ParseInt(in, 10, 64); err == nil {
		return i, nil
	}
	if !strValRe.MatchString(in) {
		return nil, errors.New(`string value may be 1-50 chars: letters, numbers, underscores, periods`)
	}
	return in, nil
}

func parseExpiry(exp string) (tsapi.Time, error) {
	if exp == "" {
		return tsapi.Time{}, nil
	}
	t, err := time.Parse(time.RFC3339, exp)
	if err != nil {
		return tsapi.Time{}, err
	}
	return tsapi.Time{Time: t.UTC()}, nil
}

func Command() *cobra.Command {
	var (
		deviceID string
		key      string
		value    string
		expiry   string
		comment  string
	)

	var (
		parsedValue any
		parsedExp   tsapi.Time
	)

	cmd := &cobra.Command{
		Use:   "posture",
		Short: "Set custom posture attributes on a device",
		Long: `Set one or more custom posture attributes on a specific device.

Examples:

  # Set a string attribute
  tscli set device posture --device node-abc123 --key custom:group --value production

  # Set a boolean attribute  
  tscli set device posture --device node-abc123 --key custom:compliant --value true

  # Set an integer attribute
  tscli set device posture --device node-abc123 --key custom:score --value 95

  # Set an attribute with expiry
  tscli set device posture --device node-abc123 --key custom:temp --value test --expiry 2024-12-31T23:59:59Z
`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if deviceID == "" {
				return fmt.Errorf("--device is required")
			}
			if key == "" {
				return fmt.Errorf("--key is required")
			}
			if value == "" {
				return fmt.Errorf("--value is required")
			}
			if !keyRe.MatchString(key) {
				return fmt.Errorf(`invalid key %q: must match "custom:..." pattern`, key)
			}

			// value
			val, err := coerceValue(value)
			if err != nil {
				return fmt.Errorf("invalid --value: %w", err)
			}
			parsedValue = val

			// expiry
			exp, err := parseExpiry(expiry)
			if err != nil {
				return fmt.Errorf("invalid --expiry: %w", err)
			}
			parsedExp = exp

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := tscli.New()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			req := tsapi.DevicePostureAttributeRequest{
				Value:   parsedValue,
				Expiry:  parsedExp,
				Comment: comment,
			}

			if err := client.Devices().SetPostureAttribute(
				context.Background(),
				deviceID,
				key,
				req,
			); err != nil {
				return fmt.Errorf("failed to set posture attribute: %w", err)
			}

			payload := map[string]string{
				"result": fmt.Sprintf("device %s: %s set to %v", deviceID, key, parsedValue),
			}
			out, _ := json.MarshalIndent(payload, "", "  ")
			outputType := viper.GetString("output")
			output.Print(outputType, out)
			return nil
		},
	}

	cmd.Flags().StringVar(&deviceID, "device", "", "Device ID")
	cmd.Flags().StringVar(&key, "key", "", "Posture attribute key (must start with 'custom:')")
	cmd.Flags().StringVar(&value, "value", "", "Posture attribute value")
	cmd.Flags().StringVar(&expiry, "expiry", "", "Expiry time in RFC3339 format (optional)")
	cmd.Flags().StringVar(&comment, "comment", "", "Optional audit-log comment")

	_ = cmd.MarkFlagRequired("device")
	_ = cmd.MarkFlagRequired("key")
	_ = cmd.MarkFlagRequired("value")

	return cmd
}
