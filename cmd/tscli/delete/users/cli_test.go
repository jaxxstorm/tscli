package users

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	tsapi "tailscale.com/client/tailscale/v2"
)

type stubRoundTripper struct {
	users       []tsapi.User
	deleteError map[string]int
	requests    []string
}

func (s *stubRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	s.requests = append(s.requests, fmt.Sprintf("%s %s", req.Method, req.URL.Path))

	switch {
	case req.Method == http.MethodGet && strings.Contains(req.URL.Path, "/users"):
		body, _ := json.Marshal(map[string][]tsapi.User{"users": s.users})
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(body)),
			Header:     make(http.Header),
		}, nil
	case req.Method == http.MethodPost && strings.Contains(req.URL.Path, "/users/") && strings.HasSuffix(req.URL.Path, "/delete"):
		for userID, status := range s.deleteError {
			if strings.Contains(req.URL.Path, "/users/"+userID+"/delete") {
				body, _ := json.Marshal(map[string]string{"message": "boom"})
				return &http.Response{
					StatusCode: status,
					Body:       io.NopCloser(bytes.NewReader(body)),
					Header:     make(http.Header),
				}, nil
			}
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{}`)),
			Header:     make(http.Header),
		}, nil
	default:
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       io.NopCloser(strings.NewReader(`{"message":"not found"}`)),
			Header:     make(http.Header),
		}, nil
	}
}

func newStubClientWithUsers(users []tsapi.User, deleteError map[string]int) (*tsapi.Client, *stubRoundTripper, error) {
	base, _ := url.Parse("http://fake")
	rt := &stubRoundTripper{users: users, deleteError: deleteError}
	return &tsapi.Client{
		BaseURL: base,
		HTTP:    &http.Client{Transport: rt},
	}, rt, nil
}

func TestBuildFilters(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		status         string
		lastSeen       string
		lastSeenSet    bool
		devices        int
		deviceCountSet bool
		wantErr        string
	}{
		{name: "status ok", status: "suspended"},
		{name: "last-seen duration ok", lastSeen: "24h", lastSeenSet: true},
		{name: "devices only ok", devices: 0, deviceCountSet: true},
		{name: "status and devices ok", status: "inactive", devices: 0, deviceCountSet: true},
		{name: "status invalid", status: "active", wantErr: "supported: inactive|suspended"},
		{name: "mutually exclusive", status: "inactive", lastSeen: "24h", lastSeenSet: true, wantErr: "mutually exclusive"},
		{name: "no filters", wantErr: "at least one of --status, --last-seen, or --devices is required"},
		{name: "negative devices", devices: -1, deviceCountSet: true, wantErr: "greater than or equal to 0"},
		{name: "invalid last-seen", lastSeen: "yesterday", lastSeenSet: true, wantErr: "invalid --last-seen"},
		{name: "empty last-seen when set", lastSeen: "", lastSeenSet: true, wantErr: "--last-seen cannot be empty"},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			_, err := buildFilters(tc.status, tc.lastSeen, tc.lastSeenSet, tc.devices, tc.deviceCountSet, false, false)
			if tc.wantErr == "" && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tc.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error containing %q", tc.wantErr)
				}
				if !strings.Contains(strings.ToLower(err.Error()), strings.ToLower(tc.wantErr)) {
					t.Fatalf("error %q does not contain %q", err.Error(), tc.wantErr)
				}
			}
		})
	}
}

func TestCommandRejectsPositionalArgs(t *testing.T) {
	t.Parallel()

	cmd := Command()
	cmd.SetArgs([]string{"devices"})
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)

	err := cmd.ExecuteContext(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestFilterUsers(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, time.April, 15, 12, 0, 0, 0, time.UTC)
	users := []tsapi.User{
		{ID: "u1", LoginName: "member-suspended@example.com", Role: tsapi.UserRole("member"), Status: tsapi.UserStatus("suspended"), DeviceCount: 0, LastSeen: now.Add(-48 * time.Hour)},
		{ID: "u2", LoginName: "member-inactive@example.com", Role: tsapi.UserRole("member"), Status: tsapi.UserStatus("inactive"), DeviceCount: 1, LastSeen: now.Add(-72 * time.Hour)},
		{ID: "u3", LoginName: "admin-suspended@example.com", Role: tsapi.UserRole("admin"), Status: tsapi.UserStatus("suspended"), DeviceCount: 0, LastSeen: now.Add(-72 * time.Hour)},
		{ID: "u4", LoginName: "owner-suspended@example.com", Role: tsapi.UserRole("owner"), Status: tsapi.UserStatus("suspended"), DeviceCount: 0, LastSeen: now.Add(-96 * time.Hour)},
		{ID: "u5", LoginName: "it-admin-suspended@example.com", Role: tsapi.UserRole("it-admin"), Status: tsapi.UserStatus("suspended"), DeviceCount: 0, LastSeen: now.Add(-96 * time.Hour)},
		{ID: "u6", LoginName: "member-active@example.com", Role: tsapi.UserRole("member"), Status: tsapi.UserStatus("active"), DeviceCount: 0, LastSeen: now.Add(-2 * time.Hour)},
		{ID: "u7", LoginName: "member-missing-lastseen@example.com", Role: tsapi.UserRole("member"), Status: tsapi.UserStatus("inactive"), DeviceCount: 0},
	}

	tests := []struct {
		name           string
		filters        deleteUserFilters
		wantCandidates []string
		wantSkipped    string
	}{
		{
			name:           "status filter excludes protected roles by default",
			filters:        deleteUserFilters{status: "suspended"},
			wantCandidates: []string{"member-suspended@example.com"},
			wantSkipped:    "protected role excluded",
		},
		{
			name:           "last seen and device count",
			filters:        deleteUserFilters{lastSeen: 24 * time.Hour, lastSeenSet: true, deviceCount: 0, deviceCountSet: true},
			wantCandidates: []string{"member-suspended@example.com"},
			wantSkipped:    "recently active",
		},
		{
			name:           "include admins explicitly",
			filters:        deleteUserFilters{status: "suspended", includeAdmins: true},
			wantCandidates: []string{"member-suspended@example.com", "admin-suspended@example.com", "owner-suspended@example.com", "it-admin-suspended@example.com"},
		},
		{
			name:           "missing last seen is skipped",
			filters:        deleteUserFilters{lastSeen: 24 * time.Hour, lastSeenSet: true, deviceCount: 0, deviceCountSet: true, includeAdmins: true},
			wantCandidates: []string{"member-suspended@example.com", "admin-suspended@example.com", "owner-suspended@example.com", "it-admin-suspended@example.com"},
			wantSkipped:    "missing lastSeen",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			candidates, skipped, err := filterUsers(users, tc.filters, now)
			if err != nil {
				t.Fatalf("filterUsers: %v", err)
			}

			var got []string
			for _, user := range candidates {
				got = append(got, user.LoginName)
			}

			if strings.Join(got, ",") != strings.Join(tc.wantCandidates, ",") {
				t.Fatalf("candidates=%v, want %v", got, tc.wantCandidates)
			}

			if tc.wantSkipped != "" && !strings.Contains(strings.Join(skipped, " | "), tc.wantSkipped) {
				t.Fatalf("skipped=%v, want substring %q", skipped, tc.wantSkipped)
			}
		})
	}
}

func TestDeleteUsersDryRunAndConfirm(t *testing.T) {
	t.Parallel()

	users := []tsapi.User{
		{ID: "u1", LoginName: "inactive@example.com", DisplayName: "Inactive User", Role: tsapi.UserRole("member"), Status: tsapi.UserStatus("inactive"), DeviceCount: 0, LastSeen: time.Now().Add(-48 * time.Hour)},
		{ID: "u2", LoginName: "admin@example.com", DisplayName: "Admin User", Role: tsapi.UserRole("admin"), Status: tsapi.UserStatus("inactive"), DeviceCount: 0, LastSeen: time.Now().Add(-48 * time.Hour)},
	}

	t.Run("dry run does not post deletes", func(t *testing.T) {
		client, rt, err := newStubClientWithUsers(users, nil)
		if err != nil {
			t.Fatalf("newStubClientWithUsers: %v", err)
		}

		summary, err := deleteUsers(context.Background(), client, deleteUserFilters{status: "inactive"})
		if err != nil {
			t.Fatalf("deleteUsers: %v", err)
		}
		if summary.Total != 1 || summary.Successful != 1 || summary.Skipped != 1 {
			t.Fatalf("unexpected summary: %+v", summary)
		}
		if len(rt.requests) != 1 {
			t.Fatalf("requests=%v, want only list request", rt.requests)
		}
	})

	t.Run("confirm posts deletes and records failures", func(t *testing.T) {
		client, rt, err := newStubClientWithUsers(users, map[string]int{"u1": http.StatusInternalServerError})
		if err != nil {
			t.Fatalf("newStubClientWithUsers: %v", err)
		}

		summary, err := deleteUsers(context.Background(), client, deleteUserFilters{status: "inactive", includeAdmins: true, confirm: true})
		if err != nil {
			t.Fatalf("deleteUsers: %v", err)
		}
		if summary.Total != 2 || summary.Failed != 1 || summary.Successful != 1 {
			t.Fatalf("unexpected summary: %+v", summary)
		}
		if len(rt.requests) != 3 {
			t.Fatalf("requests=%v, want list plus two delete requests", rt.requests)
		}
	})
}

func TestCommandValidation(t *testing.T) {
	fakeUsers := []tsapi.User{{ID: "u1", LoginName: "inactive@example.com", Role: tsapi.UserRole("member"), Status: tsapi.UserStatus("inactive"), DeviceCount: 0, LastSeen: time.Now().Add(-48 * time.Hour)}}
	stubClient := func() (*tsapi.Client, error) {
		client, _, err := newStubClientWithUsers(fakeUsers, nil)
		return client, err
	}

	tests := []struct {
		name           string
		args           []string
		wantErr        bool
		wantErrContain string
	}{
		{name: "status ok", args: []string{"--status", "inactive"}},
		{name: "last-seen ok", args: []string{"--last-seen", "24h"}},
		{name: "devices only ok", args: []string{"--devices", "0"}},
		{name: "confirm ok", args: []string{"--status", "inactive", "--confirm"}},
		{name: "invalid status", args: []string{"--status", "active"}, wantErr: true},
		{name: "mutually exclusive", args: []string{"--status", "inactive", "--last-seen", "24h"}, wantErr: true},
		{name: "empty last-seen", args: []string{"--last-seen="}, wantErr: true, wantErrContain: "--last-seen cannot be empty"},
		{name: "missing filters", args: []string{}, wantErr: true},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			save := newClient
			newClient = stubClient
			defer func() { newClient = save }()

			cmd := Command()
			cmd.SetArgs(tc.args)
			cmd.SetOut(io.Discard)
			cmd.SetErr(io.Discard)

			origStdout := os.Stdout
			origStderr := os.Stderr
			nullFile, _ := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
			os.Stdout = nullFile
			os.Stderr = nullFile
			defer func() {
				os.Stdout = origStdout
				os.Stderr = origStderr
				_ = nullFile.Close()
			}()

			err := cmd.ExecuteContext(context.Background())
			if tc.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tc.wantErrContain != "" && (err == nil || !strings.Contains(err.Error(), tc.wantErrContain)) {
				t.Fatalf("error = %v, want substring %q", err, tc.wantErrContain)
			}
		})
	}
}
