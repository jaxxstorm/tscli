package apitype

import (
	"encoding/json"
	"reflect"
	tsapi "tailscale.com/client/tailscale/v2"
	"testing"
)

func TestKeyWithoutNodeAttributesMatchesSDK(t *testing.T) {
	for _, body := range []string{`{"id":"k123"}`, `{"id":"k123","capabilities":{"devices":{"create":{"reusable":true,"ephemeral":true,"preauthorized":true,"tags":["tag:ci"]}}},"scopes":["auth_keys"],"tags":["tag:ci"],"expirySeconds":123,"created":"2026-11-07T12:00:00Z","description":"test","key":"test-material","invalid":true,"userId":"u123","keyType":"auth","issuer":"https://example.com","audience":"aud","subject":"sub","customClaimRules":{"claim":"value"}}`} {
		var local Key
		var sdk tsapi.Key
		if err := json.Unmarshal([]byte(body), &local); err != nil {
			t.Fatal(err)
		}
		if err := json.Unmarshal([]byte(body), &sdk); err != nil {
			t.Fatal(err)
		}
		a, err := json.Marshal(local)
		if err != nil {
			t.Fatal(err)
		}
		b, err := json.Marshal(sdk)
		if err != nil {
			t.Fatal(err)
		}
		var got, want any
		if err := json.Unmarshal(a, &got); err != nil {
			t.Fatal(err)
		}
		if err := json.Unmarshal(b, &want); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("output compatibility changed: %s != %s", a, b)
		}
	}
}
