package dns

import "testing"

func TestValidateNameserver(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value string
		valid bool
	}{
		{name: "ipv4", value: "1.1.1.1", valid: true},
		{name: "ipv6", value: "2606:4700:4700::1111", valid: true},
		{name: "empty clears list", value: "", valid: true},
		{name: "doh base", value: "https://dns.google/dns-query", valid: true},
		{name: "doh with path and query", value: "https://dns.nextdns.io/abcd12?device=test", valid: true},
		{name: "http rejected", value: "http://dns.google/dns-query", valid: false},
		{name: "invalid port rejected", value: "https://dns.google:65536/dns-query", valid: false},
		{name: "missing host rejected", value: "https:///dns-query", valid: false},
		{name: "nonsense rejected", value: "not-a-nameserver", valid: false},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateNameserver(tc.value)
			if tc.valid && err != nil {
				t.Fatalf("expected %q to be valid, got %v", tc.value, err)
			}
			if !tc.valid && err == nil {
				t.Fatalf("expected %q to be invalid", tc.value)
			}
		})
	}
}
