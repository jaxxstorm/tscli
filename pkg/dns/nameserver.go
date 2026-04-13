package dns

import (
	"fmt"
	"net"
	"net/url"
	"strings"
)

// ValidateNameserver accepts either a literal IP address or a valid HTTPS DoH endpoint.
func ValidateNameserver(value string) error {
	if net.ParseIP(value) != nil {
		return nil
	}

	parsed, err := url.Parse(value)
	if err != nil {
		return fmt.Errorf("invalid nameserver: %s", value)
	}
	if parsed.Scheme != "https" {
		return fmt.Errorf("invalid nameserver: %s", value)
	}
	if parsed.Host == "" || parsed.Hostname() == "" {
		return fmt.Errorf("invalid nameserver: %s", value)
	}
	if strings.Contains(parsed.Host, " ") {
		return fmt.Errorf("invalid nameserver: %s", value)
	}

	return nil
}
