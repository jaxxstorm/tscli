// cmd/tscli/delete/devices/cli.go
//
// `tscli delete devices [flags]`
// Bulk delete disconnected Tailscale devices with concurrent deletion for improved performance.

package devices

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/jaxxstorm/tscli/pkg/output"
	"github.com/jaxxstorm/tscli/pkg/tscli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tsapi "tailscale.com/client/tailscale/v2"
)

var newClient = tscli.New

// DeletionResult represents the result of a device deletion operation
type DeletionResult struct {
	DeviceID   string
	DeviceName string
	Success    bool
	Error      error
}

// DeletionSummary represents the overall summary of the deletion operation
type DeletionSummary struct {
	Total          int              `json:"total"`
	Successful     int              `json:"successful"`
	Failed         int              `json:"failed"`
	Skipped        int              `json:"skipped"`
	Results        []DeletionResult `json:"results"`
	FailedDevices  []string         `json:"failedDevices,omitempty"`
	SkippedDevices []string         `json:"skippedDevices,omitempty"`
}

func Command() *cobra.Command {
	var (
		lastSeenDuration time.Duration
		exclude          []string
		include          []string
		confirm          bool
		ephemeral        bool
	)

	cmd := &cobra.Command{
		Use:   "devices",
		Short: "Delete multiple disconnected devices",
		Long: `Delete multiple disconnected Tailscale devices based on last-seen duration.

This command allows you to bulk delete devices that haven't been seen for a specified duration.
Each deletion runs in its own goroutine for improved performance.

By default, this command shows what devices would be deleted without actually deleting them.
Use the --confirm flag to actually perform the deletions.

Examples:

  # Show devices that would be deleted (default behavior)
  tscli delete devices --last-seen 15m

  # Delete only ephemeral devices not seen for 1 hour
  tscli delete devices --last-seen 1h --ephemeral --confirm

  # Delete devices not seen for 24 hours, excluding specific patterns
  tscli delete devices --last-seen 24h --exclude server --exclude prod --confirm

  # Delete only devices matching specific patterns
  tscli delete devices --last-seen 24h --include dev --include test --confirm

  # Actually delete devices (requires --confirm)
  tscli delete devices --last-seen 15m --confirm
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Validate that exclude and include are mutually exclusive
			if len(exclude) > 0 && len(include) > 0 {
				return fmt.Errorf("--exclude and --include are mutually exclusive; use one or the other")
			}

			client, err := newClient()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			summary, err := deleteDisconnectedDevices(cmd.Context(), client, lastSeenDuration, exclude, include, ephemeral, confirm)
			if err != nil {
				return fmt.Errorf("failed to delete devices: %w", err)
			}

			out, err := json.MarshalIndent(summary, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal summary: %w", err)
			}

			outputType := viper.GetString("output")
			output.Print(outputType, out)
			return nil
		},
	}

	cmd.Flags().DurationVar(&lastSeenDuration, "last-seen", 15*time.Minute,
		"Duration to consider a device disconnected (e.g., 15m, 1h, 24h)")
	cmd.Flags().StringSliceVar(&exclude, "exclude", nil,
		"Device names to exclude by partial match (can be specified multiple times)")
	cmd.Flags().StringSliceVar(&include, "include", nil,
		"Device names to include by partial match (can be specified multiple times)")
	cmd.Flags().BoolVar(&ephemeral, "ephemeral", false,
		"Only delete ephemeral devices")
	cmd.Flags().BoolVar(&confirm, "confirm", false,
		"Actually delete devices (default is a dry run)")

	return cmd
}

func deleteDisconnectedDevices(ctx context.Context, client *tsapi.Client, lastSeenTimeout time.Duration, excludedDevices []string, includedDevices []string, ephemeralOnly bool, confirm bool) (*DeletionSummary, error) {
	// List all devices with full details to get ephemeral status
	devices, err := client.Devices().ListWithAllFields(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list devices: %w", err)
	}

	now := time.Now()
	var candidateDevices []tsapi.Device
	var skippedDevices []string

	// Filter devices based on criteria
	for _, device := range devices {
		// Check if the device matches the inclusion list (if specified)
		if len(includedDevices) > 0 && !isIncluded(device.Name, includedDevices) {
			skippedDevices = append(skippedDevices, fmt.Sprintf("%s (not included)", device.Name))
			continue
		}

		// Check if the device is in the exclusion list by partial match
		if isExcluded(device.Name, excludedDevices) {
			skippedDevices = append(skippedDevices, fmt.Sprintf("%s (excluded)", device.Name))
			continue
		}

		// Check if we should only process ephemeral devices
		if ephemeralOnly && !device.IsEphemeral {
			skippedDevices = append(skippedDevices, fmt.Sprintf("%s (not ephemeral)", device.Name))
			continue
		}

		// Calculate time since the device was last seen
		timeSinceLastSeen := now.Sub(device.LastSeen.Time)
		if timeSinceLastSeen > lastSeenTimeout {
			candidateDevices = append(candidateDevices, device)
		} else {
			skippedDevices = append(skippedDevices, fmt.Sprintf("%s (recently active)", device.Name))
		}
	}

	summary := &DeletionSummary{
		Total:          len(candidateDevices),
		SkippedDevices: skippedDevices,
		Skipped:        len(skippedDevices),
	}

	if len(candidateDevices) == 0 {
		return summary, nil
	}

	if !confirm {
		// Default behavior: show what would be deleted without actually deleting
		for _, device := range candidateDevices {
			timeSinceLastSeen := now.Sub(device.LastSeen.Time)
			summary.Results = append(summary.Results, DeletionResult{
				DeviceID:   device.ID,
				DeviceName: device.Name,
				Success:    true,
				Error:      fmt.Errorf("would delete device %s (last seen %v ago)", device.Name, timeSinceLastSeen),
			})
		}
		summary.Successful = len(candidateDevices)
		return summary, nil
	}

	// Perform concurrent deletions
	resultsChan := make(chan DeletionResult, len(candidateDevices))
	var wg sync.WaitGroup

	for _, device := range candidateDevices {
		wg.Add(1)
		go func(dev tsapi.Device) {
			defer wg.Done()
			err := client.Devices().Delete(ctx, dev.ID)
			result := DeletionResult{
				DeviceID:   dev.ID,
				DeviceName: dev.Name,
				Success:    err == nil,
				Error:      err,
			}
			resultsChan <- result
		}(device)
	}

	wg.Wait()
	close(resultsChan)

	// Collect results
	for result := range resultsChan {
		summary.Results = append(summary.Results, result)
		if result.Success {
			summary.Successful++
		} else {
			summary.Failed++
			summary.FailedDevices = append(summary.FailedDevices,
				fmt.Sprintf("%s (%s)", result.DeviceName, result.Error.Error()))
		}
	}

	return summary, nil
}

func isExcluded(deviceName string, excludedList []string) bool {
	for _, exclude := range excludedList {
		if strings.Contains(deviceName, exclude) {
			return true
		}
	}
	return false
}

func isIncluded(deviceName string, includedList []string) bool {
	for _, include := range includedList {
		if strings.Contains(deviceName, include) {
			return true
		}
	}
	return false
}
