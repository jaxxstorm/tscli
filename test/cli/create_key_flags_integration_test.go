package cli_test

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/jaxxstorm/tscli/internal/testutil/apimock"
)

func TestCreateKeyAuthkeyCapabilityFlagsPayload(t *testing.T) {
	mock := apimock.New(t)
	mock.AddJSON(http.MethodPost, "/keys", http.StatusOK, apimock.KeyResponse())

	res := executeCLI(t, []string{
		"create", "key",
		"--type", "authkey",
		"--description", "ci-runner",
		"--expiry", "24h",
		"--reusable",
		"--ephemeral",
		"--preauthorized",
	}, map[string]string{
		"TSCLI_BASE_URL": mock.URL(),
		"TSCLI_OUTPUT":   "json",
	})
	if res.err != nil {
		t.Fatalf("unexpected error: %v\nstderr:\n%s", res.err, res.stderr)
	}

	reqs := mock.Requests()
	if len(reqs) == 0 {
		t.Fatalf("expected request to mock API, got none")
	}

	var body map[string]any
	if err := json.Unmarshal([]byte(reqs[0].Body), &body); err != nil {
		t.Fatalf("unmarshal request body: %v\nbody=%s", err, reqs[0].Body)
	}

	caps, ok := body["capabilities"].(map[string]any)
	if !ok {
		t.Fatalf("expected capabilities object in request body, got %#v", body["capabilities"])
	}
	devices, ok := caps["devices"].(map[string]any)
	if !ok {
		t.Fatalf("expected devices object in capabilities, got %#v", caps["devices"])
	}
	create, ok := devices["create"].(map[string]any)
	if !ok {
		t.Fatalf("expected create object in capabilities.devices, got %#v", devices["create"])
	}

	if got := create["reusable"]; got != true {
		t.Fatalf("expected reusable=true, got %#v", got)
	}
	if got := create["ephemeral"]; got != true {
		t.Fatalf("expected ephemeral=true, got %#v", got)
	}
	if got := create["preauthorized"]; got != true {
		t.Fatalf("expected preauthorized=true, got %#v", got)
	}

	if got := int(body["expirySeconds"].(float64)); got != 86400 {
		t.Fatalf("expected expirySeconds=86400, got %d", got)
	}
}

func TestCreateKeyAuthkeyTagsPayload(t *testing.T) {
	mock := apimock.New(t)
	mock.AddJSON(http.MethodPost, "/keys", http.StatusOK, apimock.KeyResponse())

	res := executeCLI(t, []string{
		"create", "key",
		"--type", "authkey",
		"--tags", "tag:tsdns",
	}, map[string]string{
		"TSCLI_BASE_URL": mock.URL(),
		"TSCLI_OUTPUT":   "json",
	})
	if res.err != nil {
		t.Fatalf("unexpected error: %v\nstderr:\n%s", res.err, res.stderr)
	}

	reqs := mock.Requests()
	if len(reqs) == 0 {
		t.Fatalf("expected request to mock API, got none")
	}

	var body map[string]any
	if err := json.Unmarshal([]byte(reqs[0].Body), &body); err != nil {
		t.Fatalf("unmarshal request body: %v\nbody=%s", err, reqs[0].Body)
	}

	caps := body["capabilities"].(map[string]any)
	devices := caps["devices"].(map[string]any)
	create := devices["create"].(map[string]any)

	tags, ok := create["tags"].([]any)
	if !ok {
		t.Fatalf("expected tags array in capabilities.devices.create, got %#v", create["tags"])
	}
	if len(tags) != 1 || tags[0] != "tag:tsdns" {
		t.Fatalf("expected tags [tag:tsdns], got %#v", tags)
	}
}

func TestCreateKeyAuthkeyCompatibilityWithoutFlags(t *testing.T) {
	mock := apimock.New(t)
	mock.AddJSON(http.MethodPost, "/keys", http.StatusOK, apimock.KeyResponse())

	res := executeCLI(t, []string{"create", "key"}, map[string]string{
		"TSCLI_BASE_URL": mock.URL(),
		"TSCLI_OUTPUT":   "json",
	})
	if res.err != nil {
		t.Fatalf("unexpected error: %v\nstderr:\n%s", res.err, res.stderr)
	}

	reqs := mock.Requests()
	if len(reqs) == 0 {
		t.Fatalf("expected request to mock API, got none")
	}

	var body map[string]any
	if err := json.Unmarshal([]byte(reqs[0].Body), &body); err != nil {
		t.Fatalf("unmarshal request body: %v\nbody=%s", err, reqs[0].Body)
	}

	caps := body["capabilities"].(map[string]any)
	devices := caps["devices"].(map[string]any)
	create := devices["create"].(map[string]any)

	if got := create["reusable"]; got != false {
		t.Fatalf("expected reusable=false default, got %#v", got)
	}
	if got := create["ephemeral"]; got != false {
		t.Fatalf("expected ephemeral=false default, got %#v", got)
	}
	if got := create["preauthorized"]; got != false {
		t.Fatalf("expected preauthorized=false default, got %#v", got)
	}
}

func TestCreateKeyOAuthClientUnchangedWithAuthkeyFlagsPresent(t *testing.T) {
	mock := apimock.New(t)
	mock.AddJSON(http.MethodPost, "/keys", http.StatusOK, apimock.KeyResponse())

	res := executeCLI(t, []string{
		"create", "key",
		"--type", "oauthclient",
		"--scope", "users:read",
		"--reusable",
		"--ephemeral",
		"--preauthorized",
	}, map[string]string{
		"TSCLI_BASE_URL": mock.URL(),
		"TSCLI_OUTPUT":   "json",
	})
	if res.err != nil {
		t.Fatalf("unexpected error: %v\nstderr:\n%s", res.err, res.stderr)
	}

	reqs := mock.Requests()
	if len(reqs) == 0 {
		t.Fatalf("expected request to mock API, got none")
	}

	var body map[string]any
	if err := json.Unmarshal([]byte(reqs[0].Body), &body); err != nil {
		t.Fatalf("unmarshal request body: %v\nbody=%s", err, reqs[0].Body)
	}

	if got := body["keyType"]; got != "client" {
		t.Fatalf("expected keyType=client for oauth flow, got %#v", got)
	}
	if _, hasCapabilities := body["capabilities"]; hasCapabilities {
		t.Fatalf("expected oauth request to omit authkey capabilities, got %#v", body["capabilities"])
	}
}

func TestCreateKeyOAuthClientValidationStillRequiresScope(t *testing.T) {
	res := executeCLI(t, []string{
		"create", "key",
		"--type", "oauthclient",
		"--reusable",
		"--ephemeral",
		"--preauthorized",
	}, nil)
	if res.err == nil {
		t.Fatalf("expected validation error for oauthclient without scopes")
	}
	if !strings.Contains(res.err.Error(), "--scope is required") {
		t.Fatalf("unexpected validation error: %v", res.err)
	}
}
