package cli_test

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/jaxxstorm/tscli/internal/testutil/apimock"
)

type exampleOutputCase struct {
	command       string
	args          []string
	argsFunc      func(*testing.T, map[string]string) []string
	shape         jsonShapeExpectation
	textContains  []string
	supportsModes bool
	setup         func(t *testing.T, mock *apimock.Server, env map[string]string)
}

func TestExampleCommandCoverageManifest(t *testing.T) {
	expected, err := loadLeafManifest()
	if err != nil {
		t.Fatalf("load manifest: %v", err)
	}

	var actual []string
	for _, tc := range exampleOutputCases() {
		actual = append(actual, tc.command)
	}
	slices.Sort(actual)

	var missing []string
	for _, command := range expected {
		if !slices.Contains(actual, command) {
			missing = append(missing, command)
		}
	}

	var extra []string
	for _, command := range actual {
		if !slices.Contains(expected, command) {
			extra = append(extra, command)
		}
	}

	if len(missing) > 0 || len(extra) > 0 {
		t.Fatalf("example output manifest out of date\nmissing: %v\nextra: %v", missing, extra)
	}
}

func TestExampleCommandOutputShapes(t *testing.T) {
	for _, tc := range exampleOutputCases() {
		tc := tc
		t.Run(tc.command, func(t *testing.T) {
			res := runExampleOutputCase(t, tc, "json")
			if tc.shape.TopLevel != "" {
				assertJSONShape(t, res.stdout, tc.shape)
				return
			}
			assertTextOutput(t, res.stdout, tc.textContains...)
		})
	}
}

func TestExampleCommandOutputModes(t *testing.T) {
	for _, tc := range exampleOutputCases() {
		if !tc.supportsModes {
			continue
		}
		tc := tc
		for _, mode := range []string{"json", "yaml", "pretty", "human"} {
			mode := mode
			t.Run(tc.command+"/"+mode, func(t *testing.T) {
				res := runExampleOutputCase(t, tc, mode)
				assertOutputForMode(t, mode, res.stdout)
			})
		}
	}
}

func TestServiceCommandRenderedOutput(t *testing.T) {
	cases := []struct {
		name     string
		args     []string
		mode     string
		contains []string
		absent   []string
		counts   map[string]int
	}{
		{
			name:     "list services pretty",
			args:     []string{"list", "services"},
			mode:     "pretty",
			contains: []string{"svc:demo-speedtest", "svc:demo-streamer", "addrs:", "annotations:", "tailscale.com/owner-references", "─"},
			absent:   []string{"vipServices:"},
			counts: map[string]int{
				"svc:demo-speedtest": 1,
				"svc:demo-streamer":  1,
			},
		},
		{
			name:     "list services human",
			args:     []string{"list", "services"},
			mode:     "human",
			contains: []string{"svc:demo-speedtest", "svc:demo-streamer", "annotations", "tailscale.com/owner-references", "─"},
			absent:   []string{"vipServices:"},
			counts: map[string]int{
				"svc:demo-speedtest": 1,
				"svc:demo-streamer":  1,
			},
		},
		{
			name:     "get service pretty",
			args:     []string{"get", "service", "--service", "svc"},
			mode:     "pretty",
			contains: []string{"svc:demo-speedtest", "addrs:", "annotations:", "tailscale.com/owner-references"},
			absent:   []string{"vipServices:"},
			counts: map[string]int{
				"svc:demo-speedtest": 1,
			},
		},
		{
			name:     "get service human",
			args:     []string{"get", "service", "--service", "svc"},
			mode:     "human",
			contains: []string{"svc:demo-speedtest", "annotations", "tailscale.com/owner-references"},
			absent:   []string{"vipServices:"},
			counts: map[string]int{
				"svc:demo-speedtest": 1,
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mock := apimock.New(t)
			env := map[string]string{"TSCLI_OUTPUT": tc.mode, "TSCLI_BASE_URL": mock.URL()}

			switch tc.args[0] {
			case "list":
				addJSONForMethods(mock, apimock.ServiceList(), http.MethodGet)
			default:
				addJSONForMethods(mock, apimock.Service(), http.MethodGet)
			}

			res := executeCLI(t, tc.args, env)
			if res.err != nil {
				t.Fatalf("unexpected error: %v\nstderr:\n%s", res.err, res.stderr)
			}
			assertTextOutput(t, res.stdout, tc.contains...)
			for _, part := range tc.absent {
				if strings.Contains(res.stdout, part) {
					t.Fatalf("did not expect output to contain %q, got:\n%s", part, res.stdout)
				}
			}
			for part, want := range tc.counts {
				if got := strings.Count(res.stdout, part); got != want {
					t.Fatalf("expected %q count %d, got %d\noutput:\n%s", part, want, got, res.stdout)
				}
			}
		})
	}
}

func TestTailnetListRenderedOutput(t *testing.T) {
	cases := []struct {
		name     string
		mode     string
		contains []string
		absent   []string
		counts   map[string]int
	}{
		{
			name:     "list tailnets pretty",
			mode:     "pretty",
			contains: []string{"Sandbox", "createdAt:", "id:", "orgId:"},
			absent:   []string{"tailnets:"},
			counts: map[string]int{
				"Sandbox": 1,
			},
		},
		{
			name:     "list tailnets human",
			mode:     "human",
			contains: []string{"Sandbox", "createdAt", "id", "orgId"},
			absent:   []string{"tailnets:"},
			counts: map[string]int{
				"Sandbox": 1,
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mock := apimock.New(t)
			env := map[string]string{
				"TSCLI_OUTPUT":          tc.mode,
				"TSCLI_BASE_URL":        mock.URL(),
				"TSCLI_OAUTH_TOKEN_URL": mock.URL() + "/api/v2/oauth/token",
			}

			mock.AddRaw(http.MethodPost, "/api/v2/oauth/token", http.StatusOK, `{"access_token":"tok-123","token_type":"Bearer","expires_in":3600}`)
			mock.AddJSON(http.MethodGet, "/api/v2/organizations/-/tailnets", http.StatusOK, map[string]any{
				"tailnets": []map[string]any{{
					"id":          "T123",
					"displayName": "Sandbox",
					"orgId":       "o123",
					"createdAt":   "2025-01-01T12:00:00Z",
				}},
			})

			res := executeCLI(t, []string{"list", "tailnets", "--oauth-client-id", "cid", "--oauth-client-secret", "secret"}, env)
			if res.err != nil {
				t.Fatalf("unexpected error: %v\nstderr:\n%s", res.err, res.stderr)
			}

			assertTextOutput(t, res.stdout, tc.contains...)
			for _, part := range tc.absent {
				if strings.Contains(res.stdout, part) {
					t.Fatalf("did not expect output to contain %q, got:\n%s", part, res.stdout)
				}
			}
			for part, want := range tc.counts {
				if got := strings.Count(res.stdout, part); got != want {
					t.Fatalf("expected %q count %d, got %d\noutput:\n%s", part, want, got, res.stdout)
				}
			}
		})
	}
}

func TestListServicesStructuredModesPreserveRawEnvelope(t *testing.T) {
	const body = `{"vipServices":[{"name":"svc:demo-speedtest"}],"extra":"kept","count":1}`

	for _, mode := range []string{"json", "yaml"} {
		t.Run(mode, func(t *testing.T) {
			mock := apimock.New(t)
			env := map[string]string{"TSCLI_OUTPUT": mode, "TSCLI_BASE_URL": mock.URL()}
			mock.AddRaw(http.MethodGet, "/services", http.StatusOK, body)

			res := executeCLI(t, []string{"list", "services"}, env)
			if res.err != nil {
				t.Fatalf("unexpected error: %v\nstderr:\n%s", res.err, res.stderr)
			}

			switch mode {
			case "json":
				var payload map[string]any
				if err := json.Unmarshal([]byte(res.stdout), &payload); err != nil {
					t.Fatalf("json output is invalid: %v\n%s", err, res.stdout)
				}
				if got := payload["extra"]; got != "kept" {
					t.Fatalf("expected extra field to be preserved, got %#v", payload)
				}
				if got := payload["count"]; got != float64(1) {
					t.Fatalf("expected count field to be preserved, got %#v", payload)
				}
			case "yaml":
				assertTextOutput(t, res.stdout, "extra: kept", "count: 1")
			}
		})
	}
}

func TestListServicesStructuredModesDoNotMaterializeMissingVIPServices(t *testing.T) {
	const body = `{"extra":"kept"}`

	for _, mode := range []string{"json", "yaml"} {
		t.Run(mode, func(t *testing.T) {
			mock := apimock.New(t)
			env := map[string]string{"TSCLI_OUTPUT": mode, "TSCLI_BASE_URL": mock.URL()}
			mock.AddRaw(http.MethodGet, "/services", http.StatusOK, body)

			res := executeCLI(t, []string{"list", "services"}, env)
			if res.err != nil {
				t.Fatalf("unexpected error: %v\nstderr:\n%s", res.err, res.stderr)
			}

			if strings.Contains(res.stdout, "vipServices") {
				t.Fatalf("did not expect vipServices to be materialized in %s output:\n%s", mode, res.stdout)
			}
		})
	}
}

func runExampleOutputCase(t *testing.T, tc exampleOutputCase, mode string) execResult {
	t.Helper()

	mock := apimock.New(t)
	env := map[string]string{
		"TSCLI_OUTPUT": mode,
	}
	args := tc.args
	if tc.argsFunc != nil {
		args = tc.argsFunc(t, env)
	}
	if tc.setup != nil {
		tc.setup(t, mock, env)
	}

	res := executeCLI(t, args, env)
	if res.err != nil {
		t.Fatalf("unexpected error for %q: %v\nstderr:\n%s", tc.command, res.err, res.stderr)
	}
	return res
}

func exampleOutputCases() []exampleOutputCase {
	return []exampleOutputCase{
		localTextCaseWithArgs("agent init", func(t *testing.T, _ map[string]string) []string {
			return []string{"agent", "init", "--dir", t.TempDir()}
		}, nil, "tscli agent integrations initialized"),
		localTextCaseWithArgs("agent update", func(t *testing.T, env map[string]string) []string {
			repo := t.TempDir()
			env["TSCLI_AGENT_TEST_DIR"] = repo
			return []string{"agent", "update", "--dir", repo}
		}, func(t *testing.T, _ *apimock.Server, env map[string]string) {
			repo := env["TSCLI_AGENT_TEST_DIR"]
			res := executeCLINoDefaults(t, []string{"agent", "init", "--dir", repo}, nil)
			if res.err != nil {
				t.Fatalf("prepare agent update example: %v\nstderr:\n%s", res.err, res.stderr)
			}
		}, "tscli agent integrations updated"),
		localTextCase("config get", []string{"config", "get", "output"}, nil, "json"),
		localTextCase("config profiles delete", []string{"config", "profiles", "delete", "sandbox"}, setupProfileHome, "tailnet profile sandbox removed"),
		localObjectCase("config profiles list", []string{"config", "profiles", "list"}, setupProfileHome, jsonShapeExpectation{
			TopLevel:   jsonTopLevelObject,
			ObjectKeys: []string{"active-tailnet", "tailnets"},
		}),
		localTextCase("config profiles set-active", []string{"config", "profiles", "set-active", "sandbox"}, setupProfileHome, "active tailnet set to sandbox"),
		localTextCase("config profiles upsert", []string{"config", "profiles", "upsert", "sandbox", "--api-key", "tskey-sandbox"}, nil, "tailnet profile sandbox created"),
		localTextCase("config profiles upsert", []string{"config", "profiles", "upsert", "org-admin", "--oauth-client-id", "cid", "--oauth-client-secret", "secret"}, nil, "tailnet profile org-admin created"),
		localTextCase("config set", []string{"config", "set", "output", "yaml"}, nil, "output saved"),
		localObjectCase("config show", []string{"config", "show"}, nil, jsonShapeExpectation{
			TopLevel:   jsonTopLevelObject,
			ObjectKeys: []string{"output"},
		}),
		apiArrayCase("create invite device", []string{"create", "invite", "device", "--device", "node-123", "--email", "user@example.com"}, apimock.InviteList(), []string{"id", "email"}),
		apiArrayCase("create invite user", []string{"create", "invite", "user", "--email", "user@example.com"}, apimock.InviteList(), []string{"id", "email"}),
		apiObjectCase("create key", []string{"create", "key"}, apimock.KeyResponse(), "id", "key"),
		apiObjectCase("create posture-integration", []string{"create", "posture-integration", "--provider", "falcon", "--client-secret", "secret"}, apimock.PostureIntegration(), "id", "provider"),
		lifecycleObjectCase("create tailnet", []string{"create", "tailnet", "--display-name", "Sandbox", "--oauth-client-id", "cid", "--oauth-client-secret", "secret"}, map[string]any{"id": "T123", "displayName": "Sandbox", "oauthClient": map[string]any{"id": "k123", "secret": "tskey-client-secret"}}, "id", "displayName", "oauthClient"),
		oauthObjectCase("create token", []string{"create", "token", "--client-id", "cid", "--client-secret", "secret"}, "access_token", "token_type"),
		apiObjectCase("create webhook", []string{"create", "webhook", "--url", "https://example.com/hook", "--subscription", "nodeCreated"}, apimock.Webhook(), "endpointUrl"),
		summaryObjectCase("delete device", []string{"delete", "device", "--device", "node-123"}, "result"),
		summaryObjectCase("delete device invite", []string{"delete", "device", "invite", "--id", "invite-1"}, "result"),
		summaryObjectCase("delete device posture", []string{"delete", "device", "posture", "--device", "node-123", "--key", "custom:group"}, "result"),
		customCase("delete devices", []string{"delete", "devices"}, jsonShapeExpectation{
			TopLevel:   jsonTopLevelObject,
			ObjectKeys: []string{"total", "results"},
		}, true, func(t *testing.T, mock *apimock.Server, env map[string]string) {
			env["TSCLI_BASE_URL"] = mock.URL()
			addJSONForMethods(mock, apimock.DeviceList(), http.MethodGet)
		}),
		summaryObjectCase("delete key", []string{"delete", "key", "--key", "k123"}, "result"),
		summaryObjectCase("delete logs stream", []string{"delete", "logs", "stream", "--type", "network"}, "result"),
		apiObjectCase("delete posture-integration", []string{"delete", "posture-integration", "--id", "pi-1"}, map[string]any{"id": "pi-1", "deleted": true}, "id", "deleted"),
		summaryLifecycleCase("delete tailnet", []string{"delete", "tailnet", "--oauth-client-id", "cid", "--oauth-client-secret", "secret"}, "result"),
		summaryObjectCase("delete service", []string{"delete", "service", "--service", "svc"}, "result"),
		summaryObjectCase("delete user", []string{"delete", "user", "--user", "user@example.com"}, "result"),
		summaryObjectCase("delete user invite", []string{"delete", "user", "invite", "--id", "invite-1"}, "result"),
		summaryObjectCase("delete webhook", []string{"delete", "webhook", "--id", "wh-1"}, "result"),
		apiObjectCase("get contacts", []string{"get", "contacts"}, apimock.Contacts(), "account", "security"),
		apiObjectCase("get device", []string{"get", "device", "--device", "node-123", "--all"}, apimock.Device(), "advertisedRoutes", "multipleConnections", "postureIdentity"),
		apiObjectCase("get device invite", []string{"get", "device", "invite", "--id", "invite-1"}, apimock.Invite(), "id", "email"),
		apiObjectCase("get device posture", []string{"get", "device", "posture", "--device", "node-123"}, apimock.DevicePostureResponse(), "attributes"),
		apiObjectCase("get dns configuration", []string{"get", "dns", "configuration"}, apimock.DNSConfiguration(), "magicDNS"),
		apiArrayCase("get dns nameservers", []string{"get", "dns", "nameservers"}, apimock.DNSNameservers(), nil),
		apiObjectCase("get dns preferences", []string{"get", "dns", "preferences"}, map[string]any{"magicDNS": true}, "magicDNS"),
		apiArrayCase("get dns searchpaths", []string{"get", "dns", "searchpaths"}, apimock.DNSSearchPaths(), nil),
		apiObjectCase("get dns split-dns", []string{"get", "dns", "split-dns"}, apimock.DNSSplitConfig(), "corp.example.com"),
		apiObjectCase("get key", []string{"get", "key", "--key", "k123"}, apimock.KeyResponse(), "id", "key"),
		apiObjectCase("get logs aws", []string{"get", "logs", "aws"}, apimock.AWSExternalID(), "externalId"),
		apiObjectCase("get logs aws validate", []string{"get", "logs", "aws", "validate", "--external-id", "ext-123", "--role-arn", "arn:aws:iam::123456789012:role/demo"}, apimock.AWSValidation(), "valid"),
		apiObjectCase("get logs stream", []string{"get", "logs", "stream", "--type", "network"}, apimock.LogsStream(), "enabled", "endpoint"),
		customCase("get policy", []string{"get", "policy", "--json"}, jsonShapeExpectation{
			TopLevel:   jsonTopLevelObject,
			ObjectKeys: []string{"acls"},
		}, false, func(t *testing.T, mock *apimock.Server, env map[string]string) {
			env["TSCLI_BASE_URL"] = mock.URL()
			mock.AddRaw(http.MethodGet, "", http.StatusOK, apimock.Policy())
		}),
		apiObjectCase("get policy preview", []string{"get", "policy", "preview", "--type", "user", "--value", "user@example.com", "--body", "{}"}, apimock.PolicyPreview(), "matches"),
		apiObjectCase("get policy validate", []string{"get", "policy", "validate", "--body", "{}"}, apimock.PolicyValidation(), "valid"),
		apiObjectCase("get posture-integration", []string{"get", "posture-integration", "--id", "pi-1"}, apimock.PostureIntegration(), "id", "provider"),
		apiObjectCase("get service", []string{"get", "service", "--service", "svc"}, apimock.Service(), "name", "addrs", "ports", "annotations"),
		apiObjectCase("get service approval", []string{"get", "service", "approval", "--service", "svc", "--device", "node-123"}, apimock.ServiceApproval(), "approved"),
		apiObjectCase("get settings", []string{"get", "settings"}, apimock.TailnetSettings(), "devicesApprovalOn", "postureIdentityCollectionOn"),
		apiObjectCase("get user", []string{"get", "user", "--user", "user@example.com"}, apimock.User(), "id", "loginName"),
		apiObjectCase("get user invite", []string{"get", "user", "invite", "--id", "invite-1"}, apimock.Invite(), "id", "email"),
		apiObjectCase("get webhook", []string{"get", "webhook", "--id", "wh-1"}, apimock.Webhook(), "endpointUrl"),
		apiArrayCase("list devices", []string{"list", "devices", "--all"}, apimock.DeviceList(), []string{"advertisedRoutes", "multipleConnections", "postureIdentity"}),
		customCase("list invites device", []string{"list", "invites", "device", "--device", "node-123"}, jsonShapeExpectation{
			TopLevel:        jsonTopLevelArray,
			ArrayItemKeys:   []string{"id", "email"},
			RequireNonEmpty: true,
		}, true, func(t *testing.T, mock *apimock.Server, env map[string]string) {
			env["TSCLI_BASE_URL"] = mock.URL()
			mock.AddRaw(http.MethodGet, "", http.StatusOK, `[{"id":"invite-1","email":"user@example.com","status":"pending"}]`)
		}),
		apiArrayCase("list invites user", []string{"list", "invites", "user"}, apimock.InviteList(), []string{"id", "email"}),
		apiArrayCase("list keys", []string{"list", "keys", "--all"}, apimock.KeyListEnvelope(), []string{"id", "key"}),
		apiArrayCase("list logs config", []string{"list", "logs", "config"}, apimock.LogsConfiguration(), []string{"id", "action"}),
		apiArrayCase("list logs network", []string{"list", "logs", "network"}, apimock.LogsNetwork(), []string{"id", "srcIP"}),
		apiArrayCase("list nameservers", []string{"list", "nameservers"}, apimock.DNSNameservers(), nil),
		apiObjectCase("list posture-integrations", []string{"list", "posture-integrations"}, apimock.PostureIntegrationList(), "integrations"),
		lifecycleObjectCase("list tailnets", []string{"list", "tailnets", "--oauth-client-id", "cid", "--oauth-client-secret", "secret"}, map[string]any{"tailnets": []map[string]any{{"id": "T123", "displayName": "Sandbox", "orgId": "o123", "createdAt": "2025-01-01T12:00:00Z"}}}, "tailnets"),
		apiObjectCase("list routes", []string{"list", "routes", "--device", "node-123"}, apimock.DeviceRoutes(), "advertisedRoutes", "enabledRoutes"),
		customCase("list services", []string{"list", "services"}, jsonShapeExpectation{
			TopLevel:   jsonTopLevelObject,
			ObjectKeys: []string{"vipServices"},
		}, true, func(t *testing.T, mock *apimock.Server, env map[string]string) {
			env["TSCLI_BASE_URL"] = mock.URL()
			addJSONForMethods(mock, apimock.ServiceList(), http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete)
		}),
		apiArrayCase("list services devices", []string{"list", "services", "devices", "--service", "svc"}, apimock.ServiceDevices(), []string{"nodeId"}),
		apiArrayCase("list users", []string{"list", "users"}, apimock.UserListEnvelope(), []string{"id", "loginName"}),
		apiArrayCase("list webhooks", []string{"list", "webhooks"}, apimock.WebhookListEnvelope(), []string{"endpointUrl"}),
		summaryObjectCase("set contact", []string{"set", "contact", "--type", "primary", "--email", "ops@example.com"}, "result", "type", "email"),
		apiObjectCase("set device attributes", []string{"set", "device", "attributes", "--body", `{"nodes":{}}`}, map[string]any{"ok": true}, "ok"),
		summaryObjectCase("set device authorization", []string{"set", "device", "authorization", "--device", "node-123"}, "result"),
		summaryObjectCase("set device expiry", []string{"set", "device", "expiry", "--device", "node-123"}, "result"),
		apiObjectCase("set device invite", []string{"set", "device", "invite", "--id", "invite-1", "--status", "resend"}, apimock.Invite(), "id", "email"),
		summaryObjectCase("set device ip", []string{"set", "device", "ip", "--device", "node-123", "--ip", "100.64.0.42"}, "result"),
		summaryObjectCase("set device key", []string{"set", "device", "key", "--device", "node-123", "--disable-expiry"}, "result", "device", "keyExpiryDisabled"),
		summaryObjectCase("set device name", []string{"set", "device", "name", "--device", "node-123", "--name", "new-name"}, "result"),
		summaryObjectCase("set device posture", []string{"set", "device", "posture", "--device", "node-123", "--key", "custom:group", "--value", "prod"}, "result"),
		apiObjectCase("set device routes", []string{"set", "device", "routes", "--device", "node-123", "--route", "10.0.0.0/24"}, apimock.DeviceRoutes(), "advertisedRoutes", "enabledRoutes"),
		summaryObjectCase("set device tags", []string{"set", "device", "tags", "--device", "node-123", "--tag", "tag:prod"}, "result", "device", "tags"),
		apiObjectCase("set dns configuration", []string{"set", "dns", "configuration", "--body", `{"magicDNS":true}`}, apimock.DNSConfiguration(), "magicDNS"),
		apiObjectCase("set dns nameservers", []string{"set", "dns", "nameservers", "--nameserver", "1.1.1.1"}, map[string]any{"dns": []string{"1.1.1.1"}}, "dns"),
		apiObjectCase("set dns preferences", []string{"set", "dns", "preferences", "--magicdns"}, map[string]any{"magicDNS": true}, "magicDNS"),
		apiObjectCase("set dns searchpaths", []string{"set", "dns", "searchpaths", "--searchpath", "corp.example.com"}, map[string]any{"searchPaths": []string{"corp.example.com"}}, "searchPaths"),
		apiObjectCase("set dns split-dns", []string{"set", "dns", "split-dns", "--entry", "corp.example.com=1.1.1.1"}, apimock.DNSSplitConfig(), "corp.example.com"),
		apiObjectCase("set key", []string{"set", "key", "--key", "k123", "--body", `{"description":"updated"}`}, apimock.KeyResponse(), "id", "key"),
		apiObjectCase("set logs stream", []string{"set", "logs", "stream", "--type", "network", "--body", `{"endpoint":"https://example.com"}`}, apimock.LogsStream(), "enabled", "endpoint"),
		customCase("set policy", []string{"set", "policy", "--body", `{"acls":[]}`}, jsonShapeExpectation{}, false, func(t *testing.T, mock *apimock.Server, env map[string]string) {
			env["TSCLI_BASE_URL"] = mock.URL()
			mock.AddRaw(http.MethodPost, "", http.StatusOK, "")
			mock.AddRaw(http.MethodGet, "", http.StatusOK, apimock.Policy())
		}, `"acls"`),
		summaryObjectCase("set posture-integration", []string{"set", "posture-integration", "--id", "pi-1", "--client-id", "client-1"}, "result", "id", "fields"),
		apiObjectCase("set service", []string{"set", "service", "--service", "svc", "--body", `{"name":"svc"}`}, apimock.Service(), "name"),
		apiObjectCase("set service approval", []string{"set", "service", "approval", "--service", "svc", "--device", "node-123", "--approved=true"}, apimock.ServiceApproval(), "approved"),
		apiObjectCase("set settings", []string{"set", "settings", "--devices-approval"}, apimock.TailnetSettings(), "devicesApprovalOn"),
		summaryObjectCase("set user access", []string{"set", "user", "access", "--user", "user@example.com", "--approve"}, "result"),
		apiObjectCase("set user invite", []string{"set", "user", "invite", "--id", "invite-1", "--resend"}, apimock.Invite(), "id", "email"),
		summaryObjectCase("set user role", []string{"set", "user", "role", "--user", "user@example.com", "--role", "member"}, "result"),
		apiObjectCase("set webhook", []string{"set", "webhook", "--id", "wh-1", "--subscription", "nodeCreated"}, apimock.Webhook(), "id", "endpointUrl"),
		apiObjectCase("set webhook test", []string{"set", "webhook", "test", "--id", "wh-1"}, map[string]any{"ok": true}, "ok"),
		localTextCase("version", []string{"version"}, nil),
	}
}

func apiObjectCase(command string, args []string, body any, keys ...string) exampleOutputCase {
	return customCase(command, args, jsonShapeExpectation{
		TopLevel:   jsonTopLevelObject,
		ObjectKeys: keys,
	}, true, func(t *testing.T, mock *apimock.Server, env map[string]string) {
		env["TSCLI_BASE_URL"] = mock.URL()
		addJSONForMethods(mock, body, http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete)
	})
}

func apiArrayCase(command string, args []string, body any, itemKeys []string) exampleOutputCase {
	return customCase(command, args, jsonShapeExpectation{
		TopLevel:        jsonTopLevelArray,
		ArrayItemKeys:   itemKeys,
		RequireNonEmpty: true,
	}, true, func(t *testing.T, mock *apimock.Server, env map[string]string) {
		env["TSCLI_BASE_URL"] = mock.URL()
		addJSONForMethods(mock, body, http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete)
	})
}

func summaryObjectCase(command string, args []string, keys ...string) exampleOutputCase {
	return customCase(command, args, jsonShapeExpectation{
		TopLevel:   jsonTopLevelObject,
		ObjectKeys: keys,
	}, true, func(t *testing.T, mock *apimock.Server, env map[string]string) {
		env["TSCLI_BASE_URL"] = mock.URL()
		addJSONForMethods(mock, map[string]any{}, http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete)
	})
}

func oauthObjectCase(command string, args []string, keys ...string) exampleOutputCase {
	return customCase(command, args, jsonShapeExpectation{
		TopLevel:   jsonTopLevelObject,
		ObjectKeys: keys,
	}, true, func(t *testing.T, mock *apimock.Server, env map[string]string) {
		env["TSCLI_OAUTH_TOKEN_URL"] = mock.URL() + "/api/v2/oauth/token"
		mock.AddRaw(http.MethodPost, "/api/v2/oauth/token", http.StatusOK, `{"access_token":"tok-123","token_type":"Bearer","expires_in":3600}`)
	})
}

func lifecycleObjectCase(command string, args []string, body any, keys ...string) exampleOutputCase {
	return customCase(command, args, jsonShapeExpectation{
		TopLevel:   jsonTopLevelObject,
		ObjectKeys: keys,
	}, true, func(t *testing.T, mock *apimock.Server, env map[string]string) {
		env["TSCLI_BASE_URL"] = mock.URL()
		env["TSCLI_OAUTH_TOKEN_URL"] = mock.URL() + "/api/v2/oauth/token"
		mock.AddRaw(http.MethodPost, "/api/v2/oauth/token", http.StatusOK, `{"access_token":"tok-123","token_type":"Bearer","expires_in":3600}`)
		mock.AddJSON(http.MethodGet, "/api/v2/organizations/-/tailnets", http.StatusOK, body)
		mock.AddJSON(http.MethodPost, "/api/v2/organizations/-/tailnets", http.StatusOK, body)
	})
}

func summaryLifecycleCase(command string, args []string, keys ...string) exampleOutputCase {
	return customCase(command, args, jsonShapeExpectation{
		TopLevel:   jsonTopLevelObject,
		ObjectKeys: keys,
	}, true, func(t *testing.T, mock *apimock.Server, env map[string]string) {
		env["TSCLI_BASE_URL"] = mock.URL()
		env["TSCLI_OAUTH_TOKEN_URL"] = mock.URL() + "/api/v2/oauth/token"
		mock.AddRaw(http.MethodPost, "/api/v2/oauth/token", http.StatusOK, `{"access_token":"tok-123","token_type":"Bearer","expires_in":3600}`)
		mock.AddRaw(http.MethodDelete, "/api/v2/tailnet/-", http.StatusOK, `{}`)
	})
}

func localObjectCase(command string, args []string, setup func(*testing.T, *apimock.Server, map[string]string), shape jsonShapeExpectation) exampleOutputCase {
	return customCase(command, args, shape, true, setup)
}

func localTextCase(command string, args []string, setup func(*testing.T, *apimock.Server, map[string]string), contains ...string) exampleOutputCase {
	return customCase(command, args, jsonShapeExpectation{}, false, setup, contains...)
}

func localTextCaseWithArgs(command string, argsFunc func(*testing.T, map[string]string) []string, setup func(*testing.T, *apimock.Server, map[string]string), contains ...string) exampleOutputCase {
	return exampleOutputCase{
		command:       command,
		argsFunc:      argsFunc,
		shape:         jsonShapeExpectation{},
		textContains:  contains,
		supportsModes: false,
		setup:         setup,
	}
}

func customCase(command string, args []string, shape jsonShapeExpectation, supportsModes bool, setup func(*testing.T, *apimock.Server, map[string]string), contains ...string) exampleOutputCase {
	return exampleOutputCase{
		command:       command,
		args:          args,
		shape:         shape,
		textContains:  contains,
		supportsModes: supportsModes,
		setup:         setup,
	}
}

func addJSONForMethods(mock *apimock.Server, body any, methods ...string) {
	for _, method := range methods {
		mock.AddJSON(method, "", http.StatusOK, body)
	}
}

func setupProfileHome(t *testing.T, _ *apimock.Server, env map[string]string) {
	t.Helper()

	home := t.TempDir()
	path := filepath.Join(home, ".tscli.yaml")
	if err := os.WriteFile(path, []byte(`output: json
active-tailnet: other
tailnets:
  - name: other
    api-key: tskey-other
  - name: sandbox
    api-key: tskey-sandbox
`), 0o600); err != nil {
		t.Fatalf("write profile config: %v", err)
	}
	env["HOME"] = home
}
