package key

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/spf13/viper"
	tsapi "tailscale.com/client/tailscale/v2"
)

func TestCreateKeyFederated(t *testing.T) {
	t.Helper()
	origNewClient := newClient
	origDoRequest := doRequest
	t.Cleanup(func() {
		newClient = origNewClient
		doRequest = origDoRequest
	})

	newClient = func() (*tsapi.Client, error) {
		return &tsapi.Client{}, nil
	}

	doCalled := false
	doRequest = func(ctx context.Context, client *tsapi.Client, method, path string, body any, out any) (http.Header, error) {
		doCalled = true
		if method != http.MethodPost {
			t.Fatalf("expected POST, got %s", method)
		}
		if path != "/tailnet/{tailnet}/keys" {
			t.Fatalf("unexpected path %s", path)
		}
		req, ok := body.(federatedKeyRequest)
		if !ok {
			t.Fatalf("unexpected body type %T", body)
		}
		if req.KeyType != "federated" {
			t.Fatalf("unexpected key type %s", req.KeyType)
		}
		outKey, _ := out.(*tsapi.Key)
		if outKey != nil {
			*outKey = tsapi.Key{ID: "federated", KeyType: "federated"}
		}
		return nil, nil
	}

	cmd := Command()
	cmd.SetArgs([]string{
		"--type", "federated",
		"--scope", "users:read",
		"--issuer", "https://issuer.example",
		"--subject", "example-*",
	})
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)

	if err := cmd.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("command failed: %v", err)
	}
	if !doCalled {
		t.Fatalf("expected federated request to be sent")
	}
}

func TestCommandValidationErrors(t *testing.T) {
	origNewClient := newClient
	origDoRequest := doRequest
	t.Cleanup(func() {
		newClient = origNewClient
		doRequest = origDoRequest
	})

	newClient = func() (*tsapi.Client, error) {
		return &tsapi.Client{}, nil
	}
	doRequest = func(ctx context.Context, client *tsapi.Client, method, path string, body any, out any) (http.Header, error) {
		return nil, nil
	}

	viper.Set("output", "json")

	cases := []struct {
		name string
		args []string
		want string
	}{
		{"invalid type", []string{"--type", "bogus"}, "must be authkey"},
		{"oauthclient missing scope", []string{"--type", "oauthclient"}, "--scope is required"},
		{"federated missing issuer", []string{"--type", "federated", "--scope", "users:read"}, "--issuer is required"},
		{"federated missing subject", []string{"--type", "federated", "--scope", "users:read", "--issuer", "https://issuer.example"}, "--subject is required"},
		{"federated non-https issuer", []string{"--type", "federated", "--scope", "users:read", "--issuer", "http://issuer", "--subject", "subject"}, "https URL"},
		{"federated invalid claim", []string{"--type", "federated", "--scope", "users:read", "--issuer", "https://issuer.example", "--subject", "subject", "--claim", "invalid"}, "--claim must be"},
	}

	for _, tt := range cases {
		tt := tt
		cmd := Command()
		cmd.SetArgs(tt.args)
		cmd.SetOut(io.Discard)
		cmd.SetErr(io.Discard)
		err := cmd.ExecuteContext(context.Background())
		if err == nil || !strings.Contains(err.Error(), tt.want) {
			t.Fatalf("%s: expected error containing %q, got %v", tt.name, tt.want, err)
		}
	}
}
