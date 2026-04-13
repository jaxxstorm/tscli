package dns

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
)

// ValidateNameserver accepts either a literal IP address or a valid HTTPS DoH endpoint.
func ValidateNameserver(value string) error {
	if value == "" {
		return nil
	}

	if net.ParseIP(value) != nil {
		return nil
	}

	parsed, err := url.Parse(value)
	if err != nil {
		return invalidNameserverError(value)
	}
	if parsed.Scheme != "https" {
		return invalidNameserverError(value)
	}
	if parsed.Host == "" || parsed.Hostname() == "" {
		return invalidNameserverError(value)
	}
	if strings.Contains(parsed.Host, " ") {
		return invalidNameserverError(value)
	}

	if parsed.Port() != "" {
		port, err := strconv.ParseUint(parsed.Port(), 10, 16)
		if err != nil || port > 65535 {
			return invalidNameserverError(value)
		}
	}

	return nil
}

func invalidNameserverError(value string) error {
	return fmt.Errorf("invalid nameserver %q", value)
}
