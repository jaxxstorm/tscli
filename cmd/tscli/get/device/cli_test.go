package device

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	tsapi "tailscale.com/client/tailscale/v2"
)

type dummyRT struct{ list, one []byte }

func (d *dummyRT) RoundTrip(req *http.Request) (*http.Response, error) {
	var body []byte
	switch {
	case strings.Contains(req.URL.Path, "/devices"):
		body = d.list
	default:
		body = d.one
	}
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewReader(body)),
		Header:     make(http.Header),
	}, nil
}

func TestGetDeviceFlagValidation(t *testing.T) {
	t.Parallel()

	fake := tsapi.Device{
		ID:        "123",
		NodeID:    "node-123",
		Hostname:  "ok",
		OS:        "linux",
		Addresses: []string{"100.64.0.1"},
	}
	one, _ := json.Marshal(fake)
	list, _ := json.Marshal(map[string][]tsapi.Device{"devices": {fake}})

	stubClient := func() (*tsapi.Client, error) {
		base, _ := url.Parse("http://fake")
		return &tsapi.Client{
			BaseURL: base,
			HTTP:    &http.Client{Transport: &dummyRT{list: list, one: one}},
		}, nil
	}

	cases := []struct {
		name    string
		args    []string
		useStub bool
		wantErr bool
	}{
		{"missing all flags", []string{}, false, true},
		{"unknown flag", []string{"--bogus"}, false, true},
		{"device id ok", []string{"--device", "123"}, true, false},
		{"ip ok", []string{"--ip", "100.64.0.1"}, true, false},
		{"hostname ok", []string{"--name", "ok"}, true, false},
		{"mutually exclusive", []string{"--device", "123", "--ip", "100.64.0.1"}, false, true},
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
