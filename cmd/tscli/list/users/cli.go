// cmd/tscli/list/users/cli.go
package users

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jaxxstorm/tscli/pkg/output"

	"github.com/jaxxstorm/tscli/pkg/tscli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tsapi "tailscale.com/client/tailscale/v2"
)

var (
	validTypes = map[string]tsapi.UserType{
		"member": tsapi.UserType("member"),
		"shared": tsapi.UserType("shared"),
		"all":    tsapi.UserType("all"),
	}
	validRoles = map[string]tsapi.UserRole{
		"owner":         tsapi.UserRole("owner"),
		"member":        tsapi.UserRole("member"),
		"admin":         tsapi.UserRole("admin"),
		"it-admin":      tsapi.UserRole("it-admin"),
		"network-admin": tsapi.UserRole("network-admin"),
		"billing-admin": tsapi.UserRole("billing-admin"),
		"auditor":       tsapi.UserRole("auditor"),
		"all":           tsapi.UserRole("all"),
	}
)

func Command() *cobra.Command {
	var (
		typeStr string
		roleStr string
		ut      *tsapi.UserType
		ur      *tsapi.UserRole
	)

	cmd := &cobra.Command{
		Use:   "users",
		Short: "List tailnet users",
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			if typeStr != "" {
				if t, ok := validTypes[strings.ToLower(typeStr)]; ok {
					ut = &t
				} else {
					return fmt.Errorf("invalid --type value: %s", typeStr)
				}
			}
			if roleStr != "" {
				if r, ok := validRoles[strings.ToLower(roleStr)]; ok {
					ur = &r
				} else {
					return fmt.Errorf("invalid --role value: %s", roleStr)
				}
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := tscli.New()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			users, err := client.Users().List(context.Background(), ut, ur)
			if err != nil {
				return fmt.Errorf("failed to list users: %w", err)
			}

			out, err := json.MarshalIndent(users, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal users: %w", err)
			}
			outputType := viper.GetString("output")
			output.Print(outputType, out)
			return nil
		},
	}

	cmd.Flags().StringVar(&typeStr, "type", "", "Filter by user type: member|shared|all")
	cmd.Flags().StringVar(&roleStr, "role", "", "Filter by user role: owner|member|admin|it-admin|network-admin|billing-admin|auditor|all")
	return cmd
}
