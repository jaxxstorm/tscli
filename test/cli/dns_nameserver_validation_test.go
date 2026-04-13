package cli_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/jaxxstorm/tscli/internal/testutil/apimock"
)

func TestSetDNSNameserversValidation(t *testing.T) {
	cases := []struct {
		name        string
		args        []string
		errContains string
	}{
		{
			name: "accepts ipv4 nameserver",
			args: []string{"set", "dns", "nameservers", "--nameserver", "1.1.1.1"},
		},
		{
			name: "accepts doh nameserver",
			args: []string{"set", "dns", "nameservers", "--nameserver", "https://dns.google/dns-query"},
		},
		{
			name:        "rejects malformed nameserver",
			args:        []string{"set", "dns", "nameservers", "--nameserver", "dns.google"},
			errContains: "invalid nameserver",
		},
		{
			name:        "rejects non-https doh nameserver",
			args:        []string{"set", "dns", "nameservers", "--nameserver", "http://dns.google/dns-query"},
			errContains: "invalid nameserver",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mock := apimock.New(t)
			mock.AddJSON(http.MethodPost, "/dns/nameservers", http.StatusOK, map[string]any{"dns": []string{"ok"}})

			res := executeCLI(t, tc.args, map[string]string{"TSCLI_BASE_URL": mock.URL()})
			if tc.errContains == "" {
				if res.err != nil {
					t.Fatalf("unexpected error: %v", res.err)
				}
				return
			}

			if res.err == nil {
				t.Fatalf("expected validation error")
			}
			if !strings.Contains(strings.ToLower(res.err.Error()), strings.ToLower(tc.errContains)) {
				t.Fatalf("expected error containing %q, got %v", tc.errContains, res.err)
			}
		})
	}
}

func TestSetDNSSplitValidation(t *testing.T) {
	cases := []struct {
		name        string
		args        []string
		errContains string
	}{
		{
			name: "accepts doh entry",
			args: []string{"set", "dns", "split-dns", "--entry", "corp.example.com=https://dns.google/dns-query"},
		},
		{
			name:        "rejects invalid nameserver entry",
			args:        []string{"set", "dns", "split-dns", "--entry", "corp.example.com=dns.google"},
			errContains: "invalid nameserver",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mock := apimock.New(t)
			mock.AddJSON(http.MethodPatch, "/dns/split-dns", http.StatusOK, apimock.DNSSplitConfig())

			res := executeCLI(t, tc.args, map[string]string{"TSCLI_BASE_URL": mock.URL()})
			if tc.errContains == "" {
				if res.err != nil {
					t.Fatalf("unexpected error: %v", res.err)
				}
				return
			}

			if res.err == nil {
				t.Fatalf("expected validation error")
			}
			if !strings.Contains(strings.ToLower(res.err.Error()), strings.ToLower(tc.errContains)) {
				t.Fatalf("expected error containing %q, got %v", tc.errContains, res.err)
			}
		})
	}
}

func TestSetDNSNameserverDoHValuesReachAPI(t *testing.T) {
	const dohEndpoint = "https://dns.google/dns-query"

	t.Run("nameservers", func(t *testing.T) {
		mock := apimock.New(t)
		mock.AddJSON(http.MethodPost, "/dns/nameservers", http.StatusOK, map[string]any{"dns": []string{dohEndpoint}})

		res := executeCLI(t, []string{"set", "dns", "nameservers", "--nameserver", dohEndpoint}, map[string]string{"TSCLI_BASE_URL": mock.URL()})
		if res.err != nil {
			t.Fatalf("unexpected error: %v\nstderr:\n%s", res.err, res.stderr)
		}

		reqs := mock.Requests()
		if len(reqs) == 0 {
			t.Fatalf("expected request to mock API, got none")
		}
		if !strings.Contains(reqs[0].Body, dohEndpoint) {
			t.Fatalf("expected request body to include DoH endpoint, got %s", reqs[0].Body)
		}
	})

	t.Run("split dns", func(t *testing.T) {
		mock := apimock.New(t)
		mock.AddJSON(http.MethodPatch, "/dns/split-dns", http.StatusOK, apimock.DNSSplitConfig())

		res := executeCLI(t, []string{"set", "dns", "split-dns", "--entry", "corp.example.com=" + dohEndpoint}, map[string]string{"TSCLI_BASE_URL": mock.URL()})
		if res.err != nil {
			t.Fatalf("unexpected error: %v\nstderr:\n%s", res.err, res.stderr)
		}

		reqs := mock.Requests()
		if len(reqs) == 0 {
			t.Fatalf("expected request to mock API, got none")
		}
		if !strings.Contains(reqs[0].Body, dohEndpoint) {
			t.Fatalf("expected request body to include DoH endpoint, got %s", reqs[0].Body)
		}
	})
}
