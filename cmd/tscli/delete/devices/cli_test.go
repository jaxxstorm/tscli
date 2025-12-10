package devices

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
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
			LastSeen: tsapi.Time{Time: oldTime},
		},
	}

	stubClient := func() (*tsapi.Client, error) {
		base, _ := url.Parse("http://fake")
		return &tsapi.Client{
			BaseURL: base,
			HTTP:    &http.Client{Transport: &dummyRT{devices: fakeDevices}},
		}, nil
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
