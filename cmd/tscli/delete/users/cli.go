package users

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/jaxxstorm/tscli/pkg/output"
	"github.com/jaxxstorm/tscli/pkg/tscli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tsapi "tailscale.com/client/tailscale/v2"
)

var newClient = tscli.New

var validStatuses = map[string]struct{}{
	"inactive":  {},
	"suspended": {},
}

var protectedRoles = map[string]struct{}{
	"owner":         {},
	"admin":         {},
	"it-admin":      {},
	"network-admin": {},
	"billing-admin": {},
}

type DeletionResult struct {
	UserID      string `json:"userId"`
	LoginName   string `json:"loginName"`
	DisplayName string `json:"displayName,omitempty"`
	Success     bool   `json:"success"`
	Reason      string `json:"reason,omitempty"`
}

type DeletionSummary struct {
	Total        int              `json:"total"`
	Successful   int              `json:"successful"`
	Failed       int              `json:"failed"`
	Skipped      int              `json:"skipped"`
	Results      []DeletionResult `json:"results"`
	FailedUsers  []string         `json:"failedUsers,omitempty"`
	SkippedUsers []string         `json:"skippedUsers,omitempty"`
}

type deleteUserFilters struct {
	status         string
	lastSeen       time.Duration
	lastSeenSet    bool
	deviceCount    int
	deviceCountSet bool
	includeAdmins  bool
	confirm        bool
}

func Command() *cobra.Command {
	var (
		status        string
		lastSeenInput string
		deviceCount   int
		includeAdmins bool
		confirm       bool
	)

	cmd := &cobra.Command{
		Use:   "users",
		Args:  cobra.NoArgs,
		Short: "Delete multiple tailnet users",
		Long: `Delete multiple Tailscale users based on status, inactivity, and device count.

This command evaluates users returned by the list users API and deletes matching users.
By default, it performs a dry run and reports what would be deleted. Pass --confirm to
actually delete users. Privileged users are excluded unless --admins=true is provided.

Examples:

  # Show suspended users that would be deleted
  tscli delete users --status suspended

  # Show users last seen more than 24 hours ago
  tscli delete users --last-seen 24h

  # Show users with no devices
  tscli delete users --devices 0

  # Delete inactive users with no devices
  tscli delete users --status inactive --devices 0 --confirm

  # Include admin users explicitly
  tscli delete users --last-seen 24h --admins=true --confirm
`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			filters, err := buildFilters(
				status,
				lastSeenInput,
				cmd.Flags().Lookup("last-seen").Changed,
				deviceCount,
				cmd.Flags().Lookup("devices").Changed,
				includeAdmins,
				confirm,
			)
			if err != nil {
				return err
			}

			client, err := newClient()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			summary, err := deleteUsers(cmd.Context(), client, filters)
			if err != nil {
				return fmt.Errorf("failed to delete users: %w", err)
			}

			out, err := json.MarshalIndent(summary, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal summary: %w", err)
			}

			return output.Print(viper.GetString("output"), out)
		},
	}

	cmd.Flags().StringVar(&status, "status", "", "Delete users by status: inactive|suspended")
	cmd.Flags().StringVar(&lastSeenInput, "last-seen", "", "Delete users last seen longer than this duration (e.g. 24h, 30m)")
	cmd.Flags().IntVar(&deviceCount, "devices", 0, "Only delete users with this device count")
	cmd.Flags().BoolVar(&includeAdmins, "admins", false, "Include privileged users in deletion candidates")
	cmd.Flags().BoolVar(&confirm, "confirm", false, "Actually delete users (default is a dry run)")

	return cmd
}

func buildFilters(status, lastSeenInput string, lastSeenSet bool, deviceCount int, deviceCountSet, includeAdmins, confirm bool) (deleteUserFilters, error) {
	filters := deleteUserFilters{
		status:         strings.ToLower(strings.TrimSpace(status)),
		deviceCount:    deviceCount,
		deviceCountSet: deviceCountSet,
		includeAdmins:  includeAdmins,
		confirm:        confirm,
	}

	if filters.status != "" {
		if _, ok := validStatuses[filters.status]; !ok {
			return deleteUserFilters{}, fmt.Errorf("invalid --status value: %s (supported: inactive|suspended)", status)
		}
	}

	if filters.status != "" && lastSeenSet {
		return deleteUserFilters{}, fmt.Errorf("--status and --last-seen are mutually exclusive; use one or the other")
	}

	if lastSeenSet {
		lastSeenDuration, err := parseLastSeen(lastSeenInput)
		if err != nil {
			return deleteUserFilters{}, err
		}
		filters.lastSeen = lastSeenDuration
		filters.lastSeenSet = true
	}

	if deviceCountSet && deviceCount < 0 {
		return deleteUserFilters{}, fmt.Errorf("--devices must be greater than or equal to 0")
	}

	if filters.status == "" && !filters.lastSeenSet && !filters.deviceCountSet {
		return deleteUserFilters{}, fmt.Errorf("at least one of --status, --last-seen, or --devices is required")
	}

	return filters, nil
}

func parseLastSeen(input string) (time.Duration, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return 0, fmt.Errorf("--last-seen cannot be empty")
	}

	duration, err := time.ParseDuration(trimmed)
	if err != nil {
		return 0, fmt.Errorf("invalid --last-seen value %q: use a duration like 24h", input)
	}
	if duration < 0 {
		return 0, fmt.Errorf("--last-seen must be greater than or equal to 0")
	}
	return duration, nil
}

func deleteUsers(ctx context.Context, client *tsapi.Client, filters deleteUserFilters) (*DeletionSummary, error) {
	users, err := client.Users().List(ctx, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	candidates, skippedUsers, err := filterUsers(users, filters, time.Now())
	if err != nil {
		return nil, err
	}

	summary := &DeletionSummary{
		Total:        len(candidates),
		Skipped:      len(skippedUsers),
		SkippedUsers: skippedUsers,
	}

	if len(candidates) == 0 {
		return summary, nil
	}

	if !filters.confirm {
		for _, user := range candidates {
			summary.Results = append(summary.Results, DeletionResult{
				UserID:      user.ID,
				LoginName:   user.LoginName,
				DisplayName: user.DisplayName,
				Success:     true,
				Reason:      "would delete user",
			})
		}
		summary.Successful = len(candidates)
		return summary, nil
	}

	for _, user := range candidates {
		_, err := tscli.Do(ctx, client, http.MethodPost, "/users/"+user.ID+"/delete", nil, nil)
		result := DeletionResult{
			UserID:      user.ID,
			LoginName:   user.LoginName,
			DisplayName: user.DisplayName,
			Success:     err == nil,
		}
		if err != nil {
			result.Reason = err.Error()
			summary.Failed++
			summary.FailedUsers = append(summary.FailedUsers, fmt.Sprintf("%s (%s)", user.LoginName, err.Error()))
		} else {
			result.Reason = "deleted"
			summary.Successful++
		}
		summary.Results = append(summary.Results, result)
	}

	return summary, nil
}

func filterUsers(users []tsapi.User, filters deleteUserFilters, now time.Time) ([]tsapi.User, []string, error) {
	var candidates []tsapi.User
	var skipped []string

	for _, user := range users {
		if !filters.includeAdmins && isProtectedRole(user.Role) {
			skipped = append(skipped, fmt.Sprintf("%s (protected role excluded)", user.LoginName))
			continue
		}

		if filters.status != "" && !strings.EqualFold(string(user.Status), filters.status) {
			skipped = append(skipped, fmt.Sprintf("%s (status %s)", user.LoginName, user.Status))
			continue
		}

		if filters.lastSeenSet {
			if user.LastSeen.IsZero() {
				skipped = append(skipped, fmt.Sprintf("%s (missing lastSeen)", user.LoginName))
				continue
			}
			if now.Sub(user.LastSeen) <= filters.lastSeen {
				skipped = append(skipped, fmt.Sprintf("%s (recently active)", user.LoginName))
				continue
			}
		}

		if filters.deviceCountSet && user.DeviceCount != filters.deviceCount {
			skipped = append(skipped, fmt.Sprintf("%s (deviceCount %d)", user.LoginName, user.DeviceCount))
			continue
		}

		candidates = append(candidates, user)
	}

	return candidates, skipped, nil
}

func isProtectedRole(role tsapi.UserRole) bool {
	_, ok := protectedRoles[strings.ToLower(string(role))]
	return ok
}
