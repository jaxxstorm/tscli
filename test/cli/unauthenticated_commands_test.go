package cli_test

import (
	"strings"
	"testing"
)

func TestCreateTokenDoesNotRequireAPIKeyConfig(t *testing.T) {
	res := executeCLINoDefaults(t, []string{"create", "token", "--client-id", "cid", "--client-secret", "secret"}, map[string]string{
		"TSCLI_OAUTH_TOKEN_URL": "http://127.0.0.1:1/api/v2/oauth/token",
	})

	if res.err == nil {
		t.Fatalf("expected token exchange to fail against the fake OAuth endpoint")
	}

	if strings.Contains(strings.ToLower(res.err.Error()), "api key") {
		t.Fatalf("expected create token to bypass api-key runtime config, got %v", res.err)
	}
	if !strings.Contains(strings.ToLower(res.err.Error()), "failed to exchange oauth credentials") {
		t.Fatalf("expected oauth exchange error, got %v", res.err)
	}
}
