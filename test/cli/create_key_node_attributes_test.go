package cli_test

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	"github.com/jaxxstorm/tscli/internal/testutil/apimock"
)

func TestCreateKeyNodeAttributesPayload(t *testing.T) {
	mock := apimock.New(t)
	mock.AddJSON(http.MethodPost, "/keys", http.StatusOK, apimock.KeyResponse())
	res := executeCLI(t, []string{"create", "key", "--description", "posture", "--expiry", "24h", "--reusable", "--ephemeral", "--preauthorized", "--tags", "tag:ci", "--node-attributes", `{"custom:environment":"production,west","custom:rack":12,"custom:enabled":false,"custom:large":9007199254740993}`}, map[string]string{"TSCLI_BASE_URL": mock.URL(), "TSCLI_OUTPUT": "json"})
	if res.err != nil {
		t.Fatalf("command failed: %v\n%s", res.err, res.stderr)
	}
	reqs := mock.Requests()
	if len(reqs) != 1 {
		t.Fatalf("expected one request, got %d", len(reqs))
	}
	var body struct {
		Description   string `json:"description"`
		ExpirySeconds int64  `json:"expirySeconds"`
		Capabilities  struct {
			Devices struct {
				Create struct {
					Reusable       bool                       `json:"reusable"`
					Ephemeral      bool                       `json:"ephemeral"`
					Preauthorized  bool                       `json:"preauthorized"`
					Tags           []string                   `json:"tags"`
					NodeAttributes map[string]json.RawMessage `json:"nodeAttributes"`
				} `json:"create"`
			} `json:"devices"`
		} `json:"capabilities"`
	}
	if err := json.Unmarshal([]byte(reqs[0].Body), &body); err != nil {
		t.Fatal(err)
	}
	create := body.Capabilities.Devices.Create
	want := map[string]json.RawMessage{"custom:environment": json.RawMessage(`"production,west"`), "custom:rack": json.RawMessage(`12`), "custom:enabled": json.RawMessage(`false`), "custom:large": json.RawMessage(`9007199254740993`)}
	if !reflect.DeepEqual(create.NodeAttributes, want) {
		t.Fatalf("attributes = %s, want %v", reqs[0].Body, want)
	}
	if body.Description != "posture" || body.ExpirySeconds != 86400 || !create.Reusable || !create.Ephemeral || !create.Preauthorized || !reflect.DeepEqual(create.Tags, []string{"tag:ci"}) {
		t.Fatalf("existing fields lost: %s", reqs[0].Body)
	}
}
