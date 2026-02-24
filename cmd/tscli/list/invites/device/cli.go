// cmd/tscli/list/invites/device/cli.go
package device

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/jaxxstorm/tscli/pkg/output"

	"github.com/jaxxstorm/tscli/pkg/tscli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var validState = map[string]struct{}{
	"pending":  {},
	"accepted": {},
	"all":      {},
}

func Command() *cobra.Command {
	var (
		deviceID string
		state    string
	)

	cmd := &cobra.Command{
		Use:   "device",
		Short: "List device invites for a device in the tailnet",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			if state != "" {
				if _, ok := validState[strings.ToLower(state)]; !ok {
					return fmt.Errorf("invalid --state %q (pending|accepted|all)", state)
				}
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := tscli.New()
			if err != nil {
				return err
			}

			path := "/device/" + url.PathEscape(deviceID) + "/device-invites"
			if state != "" && strings.ToLower(state) != "all" {
				q := url.Values{}
				q.Set("state", strings.ToLower(state))
				path = path + "?" + q.Encode()
			}

			var raw []byte
			if _, err := tscli.Do(
				context.Background(),
				client,
				http.MethodGet,
				path,
				nil,
				&raw,
			); err != nil {
				return err
			}

			out, _ := json.MarshalIndent(json.RawMessage(raw), "", "  ")
			outputType := viper.GetString("output")
			output.Print(outputType, out)
			return nil
		},
	}

	cmd.Flags().StringVar(&deviceID, "device", "", "Device ID (nodeId or numeric id)")
	cmd.Flags().StringVar(&state, "state", "", "Filter by state: pending|accepted|all")
	_ = cmd.MarkFlagRequired("device")
	return cmd
}
