package cli_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/jaxxstorm/tscli/internal/testutil/apimock"
	"gopkg.in/yaml.v3"
)

func TestKeyResponseFailures(t *testing.T) {
	for _, tc := range []struct {
		name         string
		args         []string
		method, path string
	}{
		{"create", []string{"create", "key"}, http.MethodPost, "/keys"},
		{"get", []string{"get", "key", "--key", "k123"}, http.MethodGet, "/keys/k123"},
		{"list", []string{"list", "keys"}, http.MethodGet, "/keys"},
	} {
		for _, status := range []int{http.StatusForbidden, http.StatusOK} {
			t.Run(fmt.Sprintf("%s_%d", tc.name, status), func(t *testing.T) {
				mock := apimock.New(t)
				body := `{"message":"permission denied"}`
				if status == http.StatusOK {
					body = `{"invalid`
				}
				mock.AddRaw(tc.method, tc.path, status, body)
				res := executeCLI(t, tc.args, map[string]string{"TSCLI_BASE_URL": mock.URL(), "TSCLI_OUTPUT": "json"})
				if res.err == nil {
					t.Fatal("expected API/decoding error")
				}
				if json.Valid([]byte(res.stdout)) {
					t.Fatalf("unexpected success JSON: %s", res.stdout)
				}
				if status == http.StatusForbidden && !strings.Contains(res.err.Error(), "permission denied") {
					t.Fatalf("lost API error: %v", res.err)
				}
			})
		}
	}
}

func TestListKeysEmptyResponse(t *testing.T) {
	mock := apimock.New(t)
	mock.AddJSON(http.MethodGet, "/keys", http.StatusOK, map[string]any{"keys": []any{}})
	res := executeCLI(t, []string{"list", "keys"}, map[string]string{"TSCLI_BASE_URL": mock.URL(), "TSCLI_OUTPUT": "json"})
	if res.err != nil {
		t.Fatal(res.err)
	}
	if strings.TrimSpace(res.stdout) != "[]" {
		t.Fatalf("empty array changed: %s", res.stdout)
	}
}

func TestCreateKeyNodeAttributesResponse(t *testing.T) {
	testKeyNodeAttributesResponse(t, []string{"create", "key"}, http.MethodPost, "/keys", false)
}

func TestGetKeyNodeAttributesResponse(t *testing.T) {
	testKeyNodeAttributesResponse(t, []string{"get", "key", "--key", "k123"}, http.MethodGet, "/keys/k123", false)
}

func TestListKeysNodeAttributesResponse(t *testing.T) {
	for _, args := range [][]string{{"list", "keys"}, {"list", "keys", "--all"}} {
		t.Run(strings.Join(args, "_"), func(t *testing.T) { testKeyNodeAttributesResponse(t, args, http.MethodGet, "/keys", true) })
	}
}

func TestKeyNodeAttributesYAMLResponse(t *testing.T) {
	testKeyNodeAttributesResponse(t, []string{"create", "key"}, http.MethodPost, "/keys", false, "yaml")
}

func testKeyNodeAttributesResponse(t *testing.T, args []string, method, path string, list bool, formats ...string) {
	t.Helper()
	attrs := map[string]any{"custom:environment": "production", "custom:enabled": false, "custom:large": json.Number("9007199254740993")}
	caps := map[string]any{"devices": map[string]any{"create": map[string]any{"reusable": true, "ephemeral": false, "preauthorized": true, "tags": []string{"tag:ci"}, "nodeAttributes": attrs}}}
	key := apimock.KeyResponse()
	key["capabilities"] = caps
	key["created"] = "2026-11-07T12:00:00Z"
	key["invalid"] = true
	var response any = key
	if list {
		response = map[string]any{"keys": []any{key}}
	}
	mock := apimock.New(t)
	mock.AddJSON(method, path, http.StatusOK, response)
	format := "json"
	if len(formats) > 0 {
		format = formats[0]
	}
	res := executeCLI(t, args, map[string]string{"TSCLI_BASE_URL": mock.URL(), "TSCLI_OUTPUT": format})
	if res.err != nil {
		t.Fatalf("command: %v\n%s", res.err, res.stderr)
	}
	text := res.stdout
	if format == "yaml" {
		var value any
		if err := yaml.Unmarshal([]byte(text), &value); err != nil {
			t.Fatal(err)
		}
		b, err := json.Marshal(value)
		if err != nil {
			t.Fatal(err)
		}
		text = string(b)
	}
	dec := json.NewDecoder(strings.NewReader(text))
	dec.UseNumber()
	var got map[string]any
	if list {
		var keys []map[string]any
		if err := dec.Decode(&keys); err != nil {
			t.Fatal(err)
		}
		if len(keys) != 1 {
			t.Fatalf("keys: %v", keys)
		}
		got = keys[0]
	} else if err := dec.Decode(&got); err != nil {
		t.Fatal(err)
	}
	wantBytes, _ := json.Marshal(caps)
	wantDecoder := json.NewDecoder(strings.NewReader(string(wantBytes)))
	wantDecoder.UseNumber()
	var want map[string]any
	if err := wantDecoder.Decode(&want); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got["capabilities"], want) {
		t.Fatalf("capabilities lost: %s", res.stdout)
	}
	for _, field := range []string{"id", "key", "description", "created", "invalid"} {
		if !reflect.DeepEqual(got[field], key[field]) {
			t.Errorf("%s: got %v want %v", field, got[field], key[field])
		}
	}
	reqs := mock.Requests()
	if len(reqs) != 1 {
		t.Fatalf("requests: %v", reqs)
	}
	wantQuery := ""
	if slices.Contains(args, "--all") {
		wantQuery = "all=true"
	}
	if reqs[0].Query != wantQuery || reqs[0].Method != method || !strings.HasSuffix(reqs[0].Path, path) {
		t.Fatalf("unexpected request: %+v", reqs[0])
	}
}
