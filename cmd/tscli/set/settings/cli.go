// cmd/tscli/set/settings/cli.go
package settings

import (
	"fmt"
	"strings"

	"github.com/jaxxstorm/tscli/pkg/apitype"
	"github.com/jaxxstorm/tscli/pkg/output"
	"github.com/jaxxstorm/tscli/pkg/tscli"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var validJoin = map[string]struct{}{
	"none": {}, "admin": {}, "member": {},
}

func Command() *cobra.Command {
	var (
		devAppr, devAuto, usrAppr,
		netLog, regRoute, postureID bool
		keyDays  int
		joinRole string
	)

	cmd := &cobra.Command{
		Use:   "settings",
		Short: "Update tailnet settings",
		Long:  "Update tailnet settings and print the authoritative settings object returned by the API.",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			f := cmd.Flags()

			// require at least one flag
			changed := 0
			f.Visit(func(_ *pflag.Flag) { changed++ })
			if changed == 0 {
				return fmt.Errorf("at least one setting flag must be provided")
			}

			if f.Lookup("users-role-join").Changed {
				joinRole = strings.ToLower(joinRole)
				if _, ok := validJoin[joinRole]; !ok {
					return fmt.Errorf("invalid --users-role-join: %s (none|admin|member)", joinRole)
				}
			}
			if f.Lookup("devices-key-duration").Changed {
				if keyDays < 1 || keyDays > 180 {
					return fmt.Errorf("--devices-key-duration must be 1-180")
				}
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := tscli.New()
			if err != nil {
				return err
			}

			req := apitype.UpdateTailnetSettingsRequest{}
			f := cmd.Flags()

			if f.Lookup("devices-approval").Changed {
				req.DevicesApprovalOn = &devAppr
			}
			if f.Lookup("devices-auto-updates").Changed {
				req.DevicesAutoUpdatesOn = &devAuto
			}
			if f.Lookup("devices-key-duration").Changed {
				req.DevicesKeyDurationDays = &keyDays
			}
			if f.Lookup("users-approval").Changed {
				req.UsersApprovalOn = &usrAppr
			}
			if f.Lookup("users-role-join").Changed {
				req.UsersRoleAllowedToJoinExternalTailnets = &joinRole
			}
			if f.Lookup("network-flow-logging").Changed {
				req.NetworkFlowLoggingOn = &netLog
			}
			if f.Lookup("regional-routing").Changed {
				req.RegionalRoutingOn = &regRoute
			}
			if f.Lookup("posture-identity-collection").Changed {
				req.PostureIdentityCollectionOn = &postureID
			}

			raw, err := tscli.UpdateTailnetSettingsJSON(cmd.Context(), client, req)
			if err != nil {
				return fmt.Errorf("update failed: %w", err)
			}

			outputType := viper.GetString("output")
			return output.Print(outputType, raw)
		},
	}

	cmd.Flags().BoolVar(&devAppr, "devices-approval", false, "Enable/disable device approval")
	cmd.Flags().BoolVar(&devAuto, "devices-auto-updates", false, "Enable/disable device auto-updates")
	cmd.Flags().BoolVar(&usrAppr, "users-approval", false, "Enable/disable user approval")
	cmd.Flags().BoolVar(&netLog, "network-flow-logging", false, "Enable/disable network-flow logging")
	cmd.Flags().BoolVar(&regRoute, "regional-routing", false, "Enable/disable regional routing")
	cmd.Flags().BoolVar(&postureID, "posture-identity-collection", false, "Enable/disable posture identity collection")

	cmd.Flags().IntVar(&keyDays, "devices-key-duration", 0, "Device key expiry (1-180 days)")
	cmd.Flags().StringVar(&joinRole, "users-role-join", "", "Role allowed to join external tailnets (none|admin|member)")

	return cmd
}
