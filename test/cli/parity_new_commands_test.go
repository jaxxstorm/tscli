package cli_test

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/jaxxstorm/tscli/internal/testutil/apimock"
)

type parityCase struct {
	name        string
	args        []string
	method      string
	pathHint    string
	successBody any
}

func TestDeleteUsersCommandWithMockedAPI(t *testing.T) {
	t.Run("dry run lists candidates without deleting", func(t *testing.T) {
		mock := apimock.New(t)
		mock.AddJSON(http.MethodGet, "/users", http.StatusOK, apimock.UserListEnvelope())

		res := executeCLI(t, []string{"delete", "users", "--status", "suspended"}, map[string]string{"TSCLI_BASE_URL": mock.URL()})
		if res.err != nil {
			t.Fatalf("unexpected error: %v\nstderr:\n%s", res.err, res.stderr)
		}

		reqs := mock.Requests()
		if len(reqs) != 1 {
			t.Fatalf("expected one request, got %d (%v)", len(reqs), reqs)
		}
		if reqs[0].Method != http.MethodGet || !strings.Contains(reqs[0].Path, "/users") {
			t.Fatalf("unexpected request sequence: %v", reqs)
		}

		var summary map[string]any
		if err := json.Unmarshal([]byte(res.stdout), &summary); err != nil {
			t.Fatalf("unmarshal stdout: %v\nstdout:\n%s", err, res.stdout)
		}
		if summary["total"] != float64(1) {
			t.Fatalf("expected total=1, got %v", summary["total"])
		}
	})

	t.Run("confirm deletes matching users", func(t *testing.T) {
		mock := apimock.New(t)
		mock.AddJSON(http.MethodGet, "/users", http.StatusOK, apimock.UserListEnvelope())
		mock.AddJSON(http.MethodPost, "/users/user-2/delete", http.StatusOK, map[string]any{})

		res := executeCLI(t, []string{"delete", "users", "--status", "suspended", "--confirm"}, map[string]string{"TSCLI_BASE_URL": mock.URL()})
		if res.err != nil {
			t.Fatalf("unexpected error: %v\nstderr:\n%s", res.err, res.stderr)
		}

		reqs := mock.Requests()
		if len(reqs) != 2 {
			t.Fatalf("expected two requests, got %d (%v)", len(reqs), reqs)
		}
		if reqs[1].Method != http.MethodPost || !strings.Contains(reqs[1].Path, "/users/user-2/delete") {
			t.Fatalf("unexpected delete request: %v", reqs[1])
		}
	})
}

func TestParityCommandsWithMockedAPI(t *testing.T) {
	cases := []parityCase{
		{name: "get dns configuration", args: []string{"get", "dns", "configuration"}, method: http.MethodGet, pathHint: "/dns/configuration", successBody: map[string]any{"magicDNS": true}},
		{name: "set dns configuration", args: []string{"set", "dns", "configuration", "--body", `{"magicDNS":true}`}, method: http.MethodPost, pathHint: "/dns/configuration", successBody: map[string]any{"ok": true}},
		{name: "set key", args: []string{"set", "key", "--key", "k123", "--body", `{"description":"x"}`}, method: http.MethodPut, pathHint: "/keys/k123", successBody: map[string]any{"id": "k123"}},
		{name: "set logs stream", args: []string{"set", "logs", "stream", "--type", "network", "--body", `{"endpoint":"https://example"}`}, method: http.MethodPut, pathHint: "/logging/network/stream", successBody: map[string]any{"enabled": true}},
		{name: "delete logs stream", args: []string{"delete", "logs", "stream", "--type", "network"}, method: http.MethodDelete, pathHint: "/logging/network/stream", successBody: map[string]any{}},
		{name: "set device attributes", args: []string{"set", "device", "attributes", "--body", `{"nodes":{}}`}, method: http.MethodPatch, pathHint: "/device-attributes", successBody: map[string]any{"ok": true}},
		{name: "list services", args: []string{"list", "services"}, method: http.MethodGet, pathHint: "/services", successBody: map[string]any{"vipServices": []map[string]any{{"name": "svc"}}}},
		{name: "list services devices", args: []string{"list", "services", "devices", "--service", "svc"}, method: http.MethodGet, pathHint: "/services/svc/devices", successBody: []map[string]any{{"nodeId": "node-1"}}},
		{name: "get service", args: []string{"get", "service", "--service", "svc"}, method: http.MethodGet, pathHint: "/services/svc", successBody: map[string]any{"name": "svc"}},
		{name: "get service approval", args: []string{"get", "service", "approval", "--service", "svc", "--device", "node-1"}, method: http.MethodGet, pathHint: "/services/svc/device/node-1/approved", successBody: map[string]any{"approved": true}},
		{name: "set service", args: []string{"set", "service", "--service", "svc", "--body", `{"name":"svc"}`}, method: http.MethodPut, pathHint: "/services/svc", successBody: map[string]any{"name": "svc"}},
		{name: "set service approval", args: []string{"set", "service", "approval", "--service", "svc", "--device", "node-1", "--approved=true"}, method: http.MethodPost, pathHint: "/services/svc/device/node-1/approved", successBody: map[string]any{"approved": true}},
		{name: "delete service", args: []string{"delete", "service", "--service", "svc"}, method: http.MethodDelete, pathHint: "/services/svc", successBody: map[string]any{}},
		{name: "set webhook update", args: []string{"set", "webhook", "--id", "wh-1", "--subscription", "nodeCreated"}, method: http.MethodPatch, pathHint: "/webhooks/wh-1", successBody: map[string]any{"id": "wh-1"}},
		{name: "set webhook rotate", args: []string{"set", "webhook", "--id", "wh-1", "--rotate"}, method: http.MethodPost, pathHint: "/webhooks/wh-1/rotate", successBody: map[string]any{"id": "wh-1"}},
		{name: "set webhook test", args: []string{"set", "webhook", "test", "--id", "wh-1"}, method: http.MethodPost, pathHint: "/webhooks/wh-1/test", successBody: map[string]any{"ok": true}},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name+"_success", func(t *testing.T) {
			mock := apimock.New(t)
			mock.AddJSON(tc.method, tc.pathHint, http.StatusOK, tc.successBody)

			res := executeCLI(t, tc.args, map[string]string{"TSCLI_BASE_URL": mock.URL()})
			if res.err != nil {
				t.Fatalf("unexpected error: %v\nstderr:\n%s", res.err, res.stderr)
			}

			reqs := mock.Requests()
			if len(reqs) == 0 {
				t.Fatalf("expected request to mock API, got none")
			}
			if reqs[0].Method != tc.method {
				t.Fatalf("expected method %s, got %s", tc.method, reqs[0].Method)
			}
		})

		t.Run(tc.name+"_api_error", func(t *testing.T) {
			mock := apimock.New(t)
			mock.AddJSON(tc.method, tc.pathHint, http.StatusInternalServerError, apimock.Error("boom"))

			res := executeCLI(t, tc.args, map[string]string{"TSCLI_BASE_URL": mock.URL()})
			if res.err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !strings.Contains(strings.ToLower(res.err.Error()), "boom") {
				t.Fatalf("expected API error message in wrapped error, got %v", res.err)
			}
		})
	}
}

func TestParityCommandValidation(t *testing.T) {
	cases := []struct {
		name        string
		args        []string
		errContains string
	}{
		{name: "set dns configuration requires payload", args: []string{"set", "dns", "configuration"}, errContains: "required"},
		{name: "delete users requires filters", args: []string{"delete", "users"}, errContains: "at least one of --status, --last-seen, or --devices is required"},
		{name: "delete users rejects active status", args: []string{"delete", "users", "--status", "active"}, errContains: "inactive|suspended"},
		{name: "delete users rejects mixed activity filters", args: []string{"delete", "users", "--status", "suspended", "--last-seen", "24h"}, errContains: "mutually exclusive"},
		{name: "set key requires payload", args: []string{"set", "key", "--key", "k123"}, errContains: "required"},
		{name: "set service approval requires approved", args: []string{"set", "service", "approval", "--service", "svc", "--device", "node-1"}, errContains: "approved"},
		{name: "set webhook requires action", args: []string{"set", "webhook", "--id", "wh-1"}, errContains: "subscription"},
		{name: "set device attributes requires payload", args: []string{"set", "device", "attributes"}, errContains: "required"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			res := executeCLI(t, tc.args, nil)
			if res.err == nil {
				t.Fatalf("expected validation error")
			}
			if !strings.Contains(strings.ToLower(res.err.Error()), strings.ToLower(tc.errContains)) {
				t.Fatalf("expected error containing %q, got %v", tc.errContains, res.err)
			}
		})
	}
}
