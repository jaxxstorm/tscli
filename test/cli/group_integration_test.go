package cli_test

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/jaxxstorm/tscli/internal/testutil/apimock"
)

type groupCase struct {
	name        string
	args        []string
	method      string
	pathHint    string
	successBody any
}

func TestGroupCommandsWithMockedAPI(t *testing.T) {
	cases := []groupCase{
		{
			name:        "get",
			args:        []string{"get", "device", "--device", "node-123"},
			method:      http.MethodGet,
			pathHint:    "/device/",
			successBody: apimock.Device(),
		},
		{
			name:        "list",
			args:        []string{"list", "devices"},
			method:      http.MethodGet,
			pathHint:    "/devices",
			successBody: apimock.DeviceList(),
		},
		{
			name:        "create",
			args:        []string{"create", "key"},
			method:      http.MethodPost,
			pathHint:    "/keys",
			successBody: apimock.KeyResponse(),
		},
		{
			name:        "set",
			args:        []string{"set", "device", "name", "--device", "node-123", "--name", "new-name"},
			method:      http.MethodPost,
			pathHint:    "/name",
			successBody: map[string]any{},
		},
		{
			name:        "delete",
			args:        []string{"delete", "device", "--device", "node-123"},
			method:      http.MethodDelete,
			pathHint:    "/device/",
			successBody: map[string]any{},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name+"_success", func(t *testing.T) {
			mock := apimock.New(t)
			mock.AddJSON(tc.method, tc.pathHint, http.StatusOK, tc.successBody)

			res := executeCLI(t, tc.args, map[string]string{
				"TSCLI_BASE_URL": mock.URL(),
			})
			if res.err != nil {
				t.Fatalf("unexpected error: %v\nstderr:\n%s", res.err, res.stderr)
			}

			reqs := mock.Requests()
			if len(reqs) == 0 {
				t.Fatalf("expected request to mock API, got none")
			}
		})

		t.Run(tc.name+"_api_error", func(t *testing.T) {
			mock := apimock.New(t)
			mock.AddJSON(tc.method, tc.pathHint, http.StatusInternalServerError, apimock.Error("boom"))

			res := executeCLI(t, tc.args, map[string]string{
				"TSCLI_BASE_URL": mock.URL(),
			})
			if res.err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !strings.Contains(strings.ToLower(res.err.Error()), "boom") {
				t.Fatalf("expected API error message in wrapped error, got %v", res.err)
			}
		})
	}
}

func TestListDevicesOutputModes(t *testing.T) {
	for _, mode := range []string{"json", "yaml", "pretty", "human"} {
		mode := mode
		t.Run(mode, func(t *testing.T) {
			mock := apimock.New(t)
			mock.AddJSON(http.MethodGet, "/devices", http.StatusOK, apimock.DeviceList())

			res := executeCLI(t, []string{"list", "devices"}, map[string]string{
				"TSCLI_BASE_URL": mock.URL(),
				"TSCLI_OUTPUT":   mode,
			})
			if res.err != nil {
				t.Fatalf("unexpected error: %v", res.err)
			}
			assertOutputForMode(t, mode, res.stdout)
		})
	}
}

func TestListDevicesAllFlag(t *testing.T) {
	mock := apimock.New(t)
	mock.AddJSON(http.MethodGet, "/devices", http.StatusOK, apimock.DeviceList())

	res := executeCLI(t, []string{"list", "devices", "--all"}, map[string]string{
		"TSCLI_BASE_URL": mock.URL(),
		"TSCLI_OUTPUT":   "json",
	})
	if res.err != nil {
		t.Fatalf("unexpected error: %v", res.err)
	}
	assertOutputForMode(t, "json", res.stdout)
}

func TestListDevicesPropertyCoverage(t *testing.T) {
	mock := apimock.New(t)
	mock.AddJSON(http.MethodGet, "/devices", http.StatusOK, apimock.DeviceList())

	res := executeCLI(t, []string{"list", "devices", "--all"}, map[string]string{
		"TSCLI_BASE_URL": mock.URL(),
		"TSCLI_OUTPUT":   "json",
	})
	if res.err != nil {
		t.Fatalf("unexpected error: %v", res.err)
	}
	if !strings.Contains(res.stdout, `"advertisedRoutes"`) {
		t.Fatalf("expected advertisedRoutes in JSON output, got %s", res.stdout)
	}
	if !strings.Contains(res.stdout, `"multipleConnections"`) {
		t.Fatalf("expected multipleConnections in JSON output, got %s", res.stdout)
	}
	if !strings.Contains(res.stdout, `"postureIdentity"`) {
		t.Fatalf("expected postureIdentity in JSON output, got %s", res.stdout)
	}
	if !strings.Contains(res.stdout, `"serialNumbers"`) {
		t.Fatalf("expected postureIdentity.serialNumbers in JSON output, got %s", res.stdout)
	}
	var payload []map[string]any
	if err := json.Unmarshal([]byte(res.stdout), &payload); err != nil {
		t.Fatalf("expected top-level device array output, got %v\n%s", err, res.stdout)
	}
	if len(payload) == 0 {
		t.Fatalf("expected non-empty device array output")
	}
}

func TestGetDevicePropertyCoverage(t *testing.T) {
	mock := apimock.New(t)
	mock.AddJSON(http.MethodGet, "/device/", http.StatusOK, apimock.Device())

	res := executeCLI(t, []string{"get", "device", "--device", "node-123", "--all"}, map[string]string{
		"TSCLI_BASE_URL": mock.URL(),
		"TSCLI_OUTPUT":   "json",
	})
	if res.err != nil {
		t.Fatalf("unexpected error: %v", res.err)
	}
	for _, want := range []string{`"advertisedRoutes"`, `"multipleConnections"`, `"postureIdentity"`} {
		if !strings.Contains(res.stdout, want) {
			t.Fatalf("expected %s in JSON output, got %s", want, res.stdout)
		}
	}
}

func TestListRoutesPropertyCoverage(t *testing.T) {
	mock := apimock.New(t)
	mock.AddJSON(http.MethodGet, "/routes", http.StatusOK, apimock.DeviceRoutes())

	res := executeCLI(t, []string{"list", "routes", "--device", "node-123"}, map[string]string{
		"TSCLI_BASE_URL": mock.URL(),
		"TSCLI_OUTPUT":   "json",
	})
	if res.err != nil {
		t.Fatalf("unexpected error: %v", res.err)
	}
	for _, want := range []string{`"advertisedRoutes"`, `"enabledRoutes"`} {
		if !strings.Contains(res.stdout, want) {
			t.Fatalf("expected %s in JSON output, got %s", want, res.stdout)
		}
	}
}

func TestGetSettingsPropertyCoverage(t *testing.T) {
	mock := apimock.New(t)
	mock.AddJSON(http.MethodGet, "/settings", http.StatusOK, apimock.TailnetSettings())

	res := executeCLI(t, []string{"get", "settings"}, map[string]string{
		"TSCLI_BASE_URL": mock.URL(),
		"TSCLI_OUTPUT":   "json",
	})
	if res.err != nil {
		t.Fatalf("unexpected error: %v", res.err)
	}
	if !strings.Contains(res.stdout, `"postureIdentityCollectionOn"`) {
		t.Fatalf("expected postureIdentityCollectionOn in JSON output, got %s", res.stdout)
	}
}

func TestSetDeviceRoutesPropertyCoverage(t *testing.T) {
	mock := apimock.New(t)
	mock.AddJSON(http.MethodPost, "/routes", http.StatusOK, apimock.DeviceRoutes())

	res := executeCLI(t, []string{"set", "device", "routes", "--device", "node-123", "--route", "10.0.0.0/24"}, map[string]string{
		"TSCLI_BASE_URL": mock.URL(),
		"TSCLI_OUTPUT":   "json",
	})
	if res.err != nil {
		t.Fatalf("unexpected error: %v", res.err)
	}
	reqs := mock.Requests()
	if len(reqs) == 0 {
		t.Fatalf("expected request to mock API, got none")
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(reqs[0].Body), &body); err != nil {
		t.Fatalf("unmarshal request body: %v", err)
	}
	if _, ok := body["routes"]; !ok {
		t.Fatalf("expected routes request property, got %#v", body)
	}
	if !strings.Contains(res.stdout, `"advertisedRoutes"`) || !strings.Contains(res.stdout, `"enabledRoutes"`) {
		t.Fatalf("expected route response fields in output, got %s", res.stdout)
	}
	if strings.Contains(res.stdout, `"result"`) {
		t.Fatalf("expected authoritative API response, got synthetic summary %s", res.stdout)
	}
}

func TestSetSettingsPropertyCoverage(t *testing.T) {
	mock := apimock.New(t)
	mock.AddJSON(http.MethodPatch, "/settings", http.StatusOK, apimock.TailnetSettings())

	res := executeCLI(t, []string{"set", "settings", "--devices-approval", "--posture-identity-collection"}, map[string]string{
		"TSCLI_BASE_URL": mock.URL(),
		"TSCLI_OUTPUT":   "json",
	})
	if res.err != nil {
		t.Fatalf("unexpected error: %v", res.err)
	}
	reqs := mock.Requests()
	if len(reqs) == 0 {
		t.Fatalf("expected request to mock API, got none")
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(reqs[0].Body), &body); err != nil {
		t.Fatalf("unmarshal request body: %v", err)
	}
	for _, want := range []string{"devicesApprovalOn", "postureIdentityCollectionOn"} {
		if _, ok := body[want]; !ok {
			t.Fatalf("expected %s in request body, got %#v", want, body)
		}
	}
	if !strings.Contains(res.stdout, `"postureIdentityCollectionOn"`) {
		t.Fatalf("expected authoritative settings response in output, got %s", res.stdout)
	}
}

func TestCreateKeyValidation(t *testing.T) {
	res := executeCLI(t, []string{"create", "key", "--type", "oauthclient"}, nil)
	if res.err == nil {
		t.Fatalf("expected validation error for oauthclient without scopes")
	}
	if !strings.Contains(res.err.Error(), "--scope is required") {
		t.Fatalf("unexpected validation error: %v", res.err)
	}
}

func TestDeleteDeviceRequiresFlag(t *testing.T) {
	res := executeCLI(t, []string{"delete", "device"}, nil)
	if res.err == nil {
		t.Fatalf("expected required-flag validation error")
	}
	if !strings.Contains(strings.ToLower(res.err.Error()), "required") {
		t.Fatalf("unexpected validation error: %v", res.err)
	}
}

func TestIntegrationFailsWithoutMockServer(t *testing.T) {
	res := executeCLI(t, []string{"get", "device", "--device", "node-123"}, map[string]string{
		"TSCLI_BASE_URL": "http://127.0.0.1:1",
	})
	if res.err == nil {
		t.Fatalf("expected error when base URL is not a running mock server")
	}
}
