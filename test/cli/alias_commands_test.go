package cli_test

import (
	"strings"
	"testing"
)

func TestLegacyAliasesRemainUsable(t *testing.T) {
	cases := []struct {
		name string
		args []string
	}{
		{name: "get dns ns", args: []string{"get", "dns", "ns", "--help"}},
		{name: "get dns prefs", args: []string{"get", "dns", "prefs", "--help"}},
		{name: "get dns split", args: []string{"get", "dns", "split", "--help"}},
		{name: "set dns ns", args: []string{"set", "dns", "ns", "--help"}},
		{name: "set dns prefs", args: []string{"set", "dns", "prefs", "--help"}},
		{name: "set dns splitdns", args: []string{"set", "dns", "splitdns", "--help"}},
		{name: "get webhook webhook", args: []string{"get", "webhook", "webhook", "--help"}},
		{name: "get webhook test (legacy)", args: []string{"get", "webhook", "test", "--help"}},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			res := executeCLI(t, tc.args, nil)
			if res.err != nil {
				t.Fatalf("unexpected error: %v\nstderr:\n%s", res.err, res.stderr)
			}
			if !strings.Contains(strings.ToLower(res.stdout), "usage") {
				t.Fatalf("expected usage output, got:\n%s", res.stdout)
			}
		})
	}
}
