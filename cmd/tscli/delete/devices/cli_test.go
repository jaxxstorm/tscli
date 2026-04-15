package devices

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	tsapi "tailscale.com/client/tailscale/v2"
)

type dummyRT struct {
	devices []tsapi.Device
}

func (d *dummyRT) RoundTrip(req *http.Request) (*http.Response, error) {
	resp := map[string][]tsapi.Device{"devices": d.devices}
	body, _ := json.Marshal(resp)
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewReader(body)),
		Header:     make(http.Header),
	}, nil
}

func newStubClientWithDevices(devices []tsapi.Device) (*tsapi.Client, error) {
	base, _ := url.Parse("http://fake")
	return &tsapi.Client{
		BaseURL: base,
		HTTP:    &http.Client{Transport: &dummyRT{devices: devices}},
	}, nil
}

func TestDeleteDevicesFlagValidation(t *testing.T) {
	t.Parallel()

	// Create a fake device that was last seen long ago
	oldTime := time.Now().Add(-24 * time.Hour)
	fakeDevices := []tsapi.Device{
		{
			ID:       "123",
			NodeID:   "node-123",
			Hostname: "test-device",
			Name:     "test-device.tail123.ts.net",
			OS:       "linux",
			LastSeen: &tsapi.Time{Time: oldTime},
		},
	}

	stubClient := func() (*tsapi.Client, error) {
		return newStubClientWithDevices(fakeDevices)
	}

	cases := []struct {
		name    string
		args    []string
		useStub bool
		wantErr bool
	}{
		{"default dry run ok", []string{}, true, false},
		{"unknown flag", []string{"--bogus"}, false, true},
		{"exclude ok", []string{"--exclude", "prod"}, true, false},
		{"include ok", []string{"--include", "test"}, true, false},
		{"mutually exclusive", []string{"--exclude", "prod", "--include", "test"}, false, true},
		{"multiple excludes ok", []string{"--exclude", "prod", "--exclude", "server"}, true, false},
		{"multiple includes ok", []string{"--include", "dev", "--include", "test"}, true, false},
		{"ephemeral flag ok", []string{"--ephemeral"}, true, false},
		{"last-seen flag ok", []string{"--last-seen", "1h"}, true, false},
		{"confirm flag ok", []string{"--confirm"}, true, false},
		{"extra positional arg rejected", []string{"users"}, false, true},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			save := newClient
			if tc.useStub {
				newClient = stubClient
			}
			defer func() { newClient = save }()

			cmd := Command()
			cmd.SetArgs(tc.args)
			cmd.SetOut(io.Discard)
			cmd.SetErr(io.Discard)

			origStdout := os.Stdout
			origStderr := os.Stderr
			nullFile, _ := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
			os.Stdout = nullFile
			os.Stderr = nullFile
			defer func() {
				os.Stdout = origStdout
				os.Stderr = origStderr
				_ = nullFile.Close()
			}()

			err := cmd.ExecuteContext(context.Background())
			if tc.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestIsIncluded(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name        string
		deviceName  string
		includeList []string
		want        bool
	}{
		{"empty list", "test-device", []string{}, false},
		{"exact match", "test-device", []string{"test-device"}, true},
		{"partial match", "test-device", []string{"test"}, true},
		{"no match", "prod-server", []string{"test", "dev"}, false},
		{"multiple patterns one match", "dev-machine", []string{"test", "dev"}, true},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := isIncluded(tc.deviceName, tc.includeList)
			if got != tc.want {
				t.Fatalf("isIncluded(%q, %v) = %v, want %v", tc.deviceName, tc.includeList, got, tc.want)
			}
		})
	}
}

func TestIsExcluded(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name        string
		deviceName  string
		excludeList []string
		want        bool
	}{
		{"empty list", "test-device", []string{}, false},
		{"exact match", "test-device", []string{"test-device"}, true},
		{"partial match", "test-device", []string{"test"}, true},
		{"no match", "prod-server", []string{"test", "dev"}, false},
		{"multiple patterns one match", "dev-machine", []string{"test", "dev"}, true},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := isExcluded(tc.deviceName, tc.excludeList)
			if got != tc.want {
				t.Fatalf("isExcluded(%q, %v) = %v, want %v", tc.deviceName, tc.excludeList, got, tc.want)
			}
		})
	}
}

func TestDeleteDisconnectedDevicesFilters(t *testing.T) {
	t.Parallel()

	oldTime := time.Now().Add(-2 * time.Hour)
	devices := []tsapi.Device{
		{
			ID:       "prod-1",
			Name:     "prod-server",
			LastSeen: &tsapi.Time{Time: oldTime},
		},
		{
			ID:       "dev-1",
			Name:     "dev-machine",
			LastSeen: &tsapi.Time{Time: oldTime},
		},
		{
			ID:       "test-1",
			Name:     "qa-device",
			LastSeen: &tsapi.Time{Time: oldTime},
		},
	}

	client, err := newStubClientWithDevices(devices)
	if err != nil {
		t.Fatalf("create stub client: %v", err)
	}

	cases := []struct {
		name                   string
		include                []string
		exclude                []string
		expectedTotal          int
		expectSkippedContains  string
		expectedSkippedDetails int
	}{
		{
			name:                   "include filter",
			include:                []string{"dev"},
			expectedTotal:          2,
			expectSkippedContains:  "prod-server (not included)",
			expectedSkippedDetails: 1,
		},
		{
			name:                   "exclude filter",
			exclude:                []string{"prod"},
			expectedTotal:          2,
			expectSkippedContains:  "prod-server (excluded)",
			expectedSkippedDetails: 1,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			summary, err := deleteDisconnectedDevices(
				context.Background(),
				client,
				time.Minute,
				tc.exclude,
				tc.include,
				false,
				false,
			)
			if err != nil {
				t.Fatalf("deleteDisconnectedDevices: %v", err)
			}
			if summary.Total != tc.expectedTotal {
				t.Fatalf("total=%d, want %d", summary.Total, tc.expectedTotal)
			}
			if len(summary.SkippedDevices) != tc.expectedSkippedDetails {
				t.Fatalf("skippedDevices=%d, want %d", len(summary.SkippedDevices), tc.expectedSkippedDetails)
			}
			if tc.expectSkippedContains != "" {
				joined := strings.Join(summary.SkippedDevices, " | ")
				if !strings.Contains(joined, tc.expectSkippedContains) {
					t.Fatalf("skippedDevices missing %q: %s", tc.expectSkippedContains, joined)
				}
			}
		})
	}
}
