package cli_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/jaxxstorm/tscli/internal/testutil/apimock"
)

func TestCreateKeyNodeAttributesPermissionError(t *testing.T) {
	mock := apimock.New(t)
	mock.AddJSON(http.MethodPost, "/keys", http.StatusForbidden, map[string]string{"message": "devices:posture_attributes permission required"})
	res := executeCLI(t, []string{"create", "key", "--node-attributes", `{"custom:managed":true}`}, map[string]string{"TSCLI_BASE_URL": mock.URL(), "TSCLI_OUTPUT": "json"})
	if res.err == nil || !strings.Contains(res.err.Error(), "create auth-key") || !strings.Contains(res.err.Error(), "devices:posture_attributes") {
		t.Fatalf("expected actionable API permission error, got %v", res.err)
	}
	if len(mock.Requests()) != 1 {
		t.Fatalf("expected one request, got %d", len(mock.Requests()))
	}
}
