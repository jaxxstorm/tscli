package cli_test

import (
	"encoding/json"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

type jsonTopLevelKind string

const (
	jsonTopLevelObject jsonTopLevelKind = "object"
	jsonTopLevelArray  jsonTopLevelKind = "array"
)

type jsonShapeExpectation struct {
	TopLevel         jsonTopLevelKind
	ObjectKeys       []string
	ObjectAbsentKeys []string
	ArrayItemKeys    []string
	RequireNonEmpty  bool
}

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

func assertTextOutput(t *testing.T, out string, contains ...string) {
	t.Helper()

	if strings.TrimSpace(out) == "" {
		t.Fatalf("expected non-empty text output")
	}

	for _, part := range contains {
		if !strings.Contains(out, part) {
			t.Fatalf("expected output to contain %q, got:\n%s", part, out)
		}
	}
}

func assertJSONShape(t *testing.T, out string, want jsonShapeExpectation) {
	t.Helper()

	var payload any
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("output is not valid json: %v\n%s", err, out)
	}

	switch want.TopLevel {
	case jsonTopLevelObject:
		obj, ok := payload.(map[string]any)
		if !ok {
			t.Fatalf("expected top-level object, got %T\n%s", payload, out)
		}
		for _, key := range want.ObjectKeys {
			if _, ok := obj[key]; !ok {
				t.Fatalf("expected top-level key %q in %#v", key, obj)
			}
		}
		for _, key := range want.ObjectAbsentKeys {
			if _, ok := obj[key]; ok {
				t.Fatalf("did not expect top-level key %q in %#v", key, obj)
			}
		}
	case jsonTopLevelArray:
		items, ok := payload.([]any)
		if !ok {
			t.Fatalf("expected top-level array, got %T\n%s", payload, out)
		}
		if want.RequireNonEmpty && len(items) == 0 {
			t.Fatalf("expected non-empty array output")
		}
		if len(want.ArrayItemKeys) == 0 || len(items) == 0 {
			return
		}
		first, ok := items[0].(map[string]any)
		if !ok {
			t.Fatalf("expected first array item to be object, got %T", items[0])
		}
		for _, key := range want.ArrayItemKeys {
			if _, ok := first[key]; !ok {
				t.Fatalf("expected array item key %q in %#v", key, first)
			}
		}
	default:
		t.Fatalf("unknown top-level kind %q", want.TopLevel)
	}
}
