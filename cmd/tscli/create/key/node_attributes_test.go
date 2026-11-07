package key

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"testing"

	tsapi "tailscale.com/client/tailscale/v2"
)

func TestParseNodeAttributesBoundaries(t *testing.T) {
	for _, count := range []int{0, 100} {
		input := map[string]int{}
		for i := 0; i < count; i++ {
			input[fmt.Sprintf("custom:a%d", i)] = -i
		}
		raw, err := json.Marshal(input)
		if err != nil {
			t.Fatal(err)
		}
		got, err := parseNodeAttributes(string(raw))
		if err != nil || len(got) != count {
			t.Fatalf("count %d: got %v, %v", count, got, err)
		}
	}
	got, err := parseNodeAttributes(`{"custom:zero":0,"custom:negative":-42,"custom:empty":"","custom:true":true}`)
	if err != nil {
		t.Fatal(err)
	}
	for key, want := range map[string]string{"custom:zero": "0", "custom:negative": "-42", "custom:empty": `""`, "custom:true": "true"} {
		if string(got[key]) != want {
			t.Fatalf("%s = %s, want %s", key, got[key], want)
		}
	}
}

func TestNodeAttributesValidationBeforeClient(t *testing.T) {
	original := newClient
	t.Cleanup(func() { newClient = original })
	newClient = func() (*tsapi.Client, error) { return nil, fmt.Errorf("invalid input reached client initialization") }
	tooMany := map[string]int{}
	for i := 0; i < 101; i++ {
		tooMany[fmt.Sprintf("custom:a%d", i)] = i
	}
	oversized, err := json.Marshal(tooMany)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct{ name, value, kind, want string }{
		{"null object", `null`, "authkey", "JSON object"},
		{"array", `[]`, "authkey", "JSON object"},
		{"malformed", `{`, "authkey", "JSON object"},
		{"trailing value", `{} {}`, "authkey", "JSON object"},
		{"namespace", `{"rack":1}`, "authkey", "custom:"},
		{"empty name", `{"custom:":1}`, "authkey", "custom:"},
		{"null value", `{"custom:a":null}`, "authkey", "string, integer, or boolean"},
		{"object value", `{"custom:a":{}}`, "authkey", "string, integer, or boolean"},
		{"array value", `{"custom:a":[]}`, "authkey", "string, integer, or boolean"},
		{"fraction", `{"custom:a":1.5}`, "authkey", "string, integer, or boolean"},
		{"limit", string(oversized), "authkey", "100"},
		{"oauth", `{"custom:a":true}`, "oauthclient", "only with --type authkey"},
		{"federated", `{}`, "federated", "only with --type authkey"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cmd := Command()
			cmd.SetOut(io.Discard)
			cmd.SetErr(io.Discard)
			cmd.SetArgs([]string{"--type", tc.kind, "--scope", "users:read", "--issuer", "https://issuer.example", "--subject", "test", "--node-attributes", tc.value})
			err := cmd.Execute()
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("want error containing %q, got %v", tc.want, err)
			}
		})
	}
}
