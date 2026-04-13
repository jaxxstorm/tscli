package cli_test

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jaxxstorm/tscli/internal/testutil/apimock"
)

func TestCommandsWithPreRunUseActiveProfileRuntimeConfig(t *testing.T) {
	home := t.TempDir()
	configFile := filepath.Join(home, ".tscli.yaml")
	cfg := strings.Join([]string{
		"output: json",
		"active-tailnet: profile-tailnet",
		"tailnets:",
		"  - name: profile-tailnet",
		"    api-key: tskey-profile",
		"",
	}, "\n")
	if err := os.WriteFile(configFile, []byte(cfg), 0o600); err != nil {
		t.Fatalf("write config file: %v", err)
	}

	cases := []struct {
		name     string
		args     []string
		method   string
		pathHint string
		body     any
	}{
		{
			name:     "get service",
			args:     []string{"get", "service", "--service", "svc:portal"},
			method:   http.MethodGet,
			pathHint: "/tailnet/profile-tailnet/services/svc:portal",
			body:     apimock.Service(),
		},
		{
			name:     "get service approval",
			args:     []string{"get", "service", "approval", "--service", "svc:portal", "--device", "node-123"},
			method:   http.MethodGet,
			pathHint: "/tailnet/profile-tailnet/services/svc:portal/device/node-123/approved",
			body:     apimock.ServiceApproval(),
		},
		{
			name:     "list services devices",
			args:     []string{"list", "services", "devices", "--service", "svc:portal"},
			method:   http.MethodGet,
			pathHint: "/tailnet/profile-tailnet/services/svc:portal/devices",
			body:     apimock.ServiceDevices(),
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mock := apimock.New(t)
			mock.AddJSON(tc.method, tc.pathHint, http.StatusOK, tc.body)

			res := executeCLINoDefaults(t, tc.args, map[string]string{
				"HOME":           home,
				"TSCLI_BASE_URL": mock.URL(),
			})
			if res.err != nil {
				t.Fatalf("unexpected error: %v\nstderr:\n%s", res.err, res.stderr)
			}

			reqs := mock.Requests()
			if len(reqs) == 0 {
				t.Fatalf("expected request to mock API, got none")
			}
			if !strings.Contains(reqs[0].Path, tc.pathHint) {
				t.Fatalf("expected request path %q, got %q", tc.pathHint, reqs[0].Path)
			}
		})
	}
}
