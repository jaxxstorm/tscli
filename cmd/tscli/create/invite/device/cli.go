// cmd/tscli/create/invite/device/cli.go
//
// Create *share invites* for an **existing device**.
//
//	# invite alice@example.com
//	tscli create invite device --device node-abc --email alice@example.com
//
//	# generate a multi-use invite (share link) that allows exit-node
//	tscli create invite device --device 123456 --multi-use --allow-exit-node
package device

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/mail"

	"github.com/jaxxstorm/tscli/pkg/output"

	"github.com/jaxxstorm/tscli/pkg/tscli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type invite struct {
	MultiUse      bool   `json:"multiUse,omitempty"`
	AllowExitNode bool   `json:"allowExitNode,omitempty"`
	Email         string `json:"email,omitempty"`
}

// Command registers `tscli create invite device`.
func Command() *cobra.Command {
	var (
		deviceID      string
		emails        []string
		multiUse      bool
		allowExitNode bool
	)

	cmd := &cobra.Command{
		Use:   "device",
		Short: "Create share invite(s) for a device",
		PreRunE: func(*cobra.Command, []string) error {
			if deviceID == "" {
				return errors.New("--device is required")
			}

			// validate e-mails
			for _, e := range emails {
				if _, err := mail.ParseAddress(e); err != nil {
					return fmt.Errorf("invalid --email %q: %w", e, err)
				}
			}

			if len(emails) == 0 && !multiUse {
				return errors.New("use --email (one or more) or --multi-use")
			}
			return nil
		},

		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := tscli.New()
			if err != nil {
				return err
			}

			var payload []invite
			if len(emails) > 0 {
				for _, e := range emails {
					payload = append(payload, invite{
						Email:         e,
						AllowExitNode: allowExitNode,
						MultiUse:      multiUse, // if user set --multi-use together with e-mails
					})
				}
			} else { // purely multi-use invite
				payload = []invite{{
					MultiUse:      true,
					AllowExitNode: allowExitNode,
				}}
			}

			var resp json.RawMessage
			if _, err := tscli.Do(
				context.Background(),
				client,
				http.MethodPost,
				fmt.Sprintf("/device/%s/device-invites", deviceID),
				payload,
				&resp,
			); err != nil {
				return fmt.Errorf("invite failed: %w", err)
			}

			out, _ := json.MarshalIndent(resp, "", "  ")
			outputType := viper.GetString("output")
			output.Print(outputType, out)
			return nil
		},
	}

	cmd.Flags().StringVar(&deviceID, "device", "", "Device ID or nodeId to share (required)")
	cmd.Flags().StringSliceVarP(&emails, "email", "e", nil,
		"Recipient e-mail. Omit for a multi-use invite.")
	cmd.Flags().BoolVar(&multiUse, "multi-use", false, "Generate a multi-use share link")
	cmd.Flags().BoolVar(&allowExitNode, "allow-exit-node", false,
		"Let the shared device act as an exit-node")

	_ = cmd.MarkFlagRequired("device")
	return cmd
}
