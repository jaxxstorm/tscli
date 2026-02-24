package main

import (
	"encoding/json"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func assertOutputForMode(t *testing.T, mode, out string) {
	t.Helper()

	if strings.TrimSpace(out) == "" {
		t.Fatalf("expected non-empty output for mode %q", mode)
	}

	switch mode {
	case "json":
		var v any
		if err := json.Unmarshal([]byte(out), &v); err != nil {
			t.Fatalf("output is not valid json: %v\n%s", err, out)
		}
	case "yaml":
		var v any
		if err := yaml.Unmarshal([]byte(out), &v); err != nil {
			t.Fatalf("output is not valid yaml: %v\n%s", err, out)
		}
	case "pretty", "human":
		// These are presentation modes; non-empty output is sufficient.
	default:
		t.Fatalf("unknown output mode %q", mode)
	}
}
