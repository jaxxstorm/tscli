package cli_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/jaxxstorm/tscli/internal/testutil/apimock"
)

func TestGetSettingsAdminConsoleSessionTimeout(t *testing.T) {
	for _, minutes := range []int{0, 180, 43200} {
		t.Run(fmt.Sprint(minutes), func(t *testing.T) {
			mock := apimock.New(t)
			want := map[string]any{"adminConsoleSessionTimeoutMinutes": float64(minutes)}
			mock.AddJSON(http.MethodGet, "/settings", http.StatusOK, want)
			res := executeCLI(t, []string{"get", "settings"}, map[string]string{"TSCLI_BASE_URL": mock.URL(), "TSCLI_OUTPUT": "json"})
			if res.err != nil {
				t.Fatal(res.err)
			}
			var got map[string]any
			if err := json.Unmarshal([]byte(res.stdout), &got); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("got %v, want %v", got, want)
			}
		})
	}
}

func TestSetSettingsAdminConsoleSessionTimeoutValidation(t *testing.T) {
	for _, minutes := range []string{"-1", "43201"} {
		t.Run(minutes, func(t *testing.T) {
			mock := apimock.New(t)
			mock.AddJSON(http.MethodPatch, "/settings", http.StatusOK, map[string]any{})
			res := executeCLI(t, []string{"set", "settings", "--admin-console-session-timeout-minutes=" + minutes}, map[string]string{"TSCLI_BASE_URL": mock.URL()})
			if res.err == nil || !strings.Contains(res.err.Error(), "--admin-console-session-timeout-minutes must be 0-43200") {
				t.Errorf("expected range validation error, got %v", res.err)
			}
			if len(mock.Requests()) != 0 {
				t.Fatal("invalid timeout sent an API request")
			}
		})
	}
}

func TestSetSettingsAdminConsoleSessionTimeoutOmitted(t *testing.T) {
	mock := apimock.New(t)
	mock.AddJSON(http.MethodPatch, "/settings", http.StatusOK, map[string]any{})
	res := executeCLI(t, []string{"set", "settings", "--devices-approval"}, map[string]string{"TSCLI_BASE_URL": mock.URL()})
	if res.err != nil {
		t.Fatal(res.err)
	}
	reqs := mock.Requests()
	if len(reqs) != 1 {
		t.Fatalf("got %d requests, want 1", len(reqs))
	}
	var got map[string]any
	if err := json.Unmarshal([]byte(reqs[0].Body), &got); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, map[string]any{"devicesApprovalOn": true}) {
		t.Fatalf("unexpected patch body: %v", got)
	}
}

func TestSetSettingsAdminConsoleSessionTimeout(t *testing.T) {
	for _, minutes := range []int{0, 1, 180, 43200} {
		t.Run(fmt.Sprint(minutes), func(t *testing.T) {
			mock := apimock.New(t)
			want := map[string]any{"adminConsoleSessionTimeoutMinutes": float64(minutes)}
			mock.AddJSON(http.MethodPatch, "/settings", http.StatusOK, want)
			res := executeCLI(t, []string{"set", "settings", "--admin-console-session-timeout-minutes", fmt.Sprint(minutes)}, map[string]string{"TSCLI_BASE_URL": mock.URL(), "TSCLI_OUTPUT": "json"})
			if res.err != nil {
				t.Fatal(res.err)
			}
			reqs := mock.Requests()
			if len(reqs) != 1 {
				t.Fatalf("got %d requests, want 1", len(reqs))
			}
			for _, raw := range []string{reqs[0].Body, res.stdout} {
				var got map[string]any
				if err := json.Unmarshal([]byte(raw), &got); err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(got, want) {
					t.Fatalf("got %v, want %v", got, want)
				}
			}
		})
	}
}
