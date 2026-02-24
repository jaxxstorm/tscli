package openapi

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

type snapshotDoc struct {
	OpenAPI string                    `yaml:"openapi"`
	Info    map[string]any            `yaml:"info"`
	Paths   map[string]map[string]any `yaml:"paths"`
}

type snapshotMeta struct {
	SourceURL    string `yaml:"source_url"`
	FetchedAtUTC string `yaml:"fetched_at_utc"`
	OpenAPI      string `yaml:"openapi_version"`
	APIVersion   string `yaml:"api_version"`
	SHA256       string `yaml:"sha256"`
	SchemaFile   string `yaml:"schema_file"`
}

type opMap struct {
	Commands map[string][]string `yaml:"commands"`
}

func TestSnapshotMetadataMatchesFile(t *testing.T) {
	data, err := os.ReadFile("tailscale-v2-openapi.yaml")
	if err != nil {
		t.Fatalf("read schema: %v", err)
	}
	metaBytes, err := os.ReadFile("snapshot-metadata.yaml")
	if err != nil {
		t.Fatalf("read metadata: %v", err)
	}

	var meta snapshotMeta
	if err := yaml.Unmarshal(metaBytes, &meta); err != nil {
		t.Fatalf("unmarshal metadata: %v", err)
	}

	sum := sha256.Sum256(data)
	got := hex.EncodeToString(sum[:])
	if meta.SHA256 != got {
		t.Fatalf("sha256 mismatch: metadata=%s actual=%s", meta.SHA256, got)
	}
	if meta.SourceURL == "" || meta.FetchedAtUTC == "" {
		t.Fatalf("snapshot metadata missing source_url/fetched_at_utc")
	}
}

func TestSnapshotHasExpectedShape(t *testing.T) {
	data, err := os.ReadFile("tailscale-v2-openapi.yaml")
	if err != nil {
		t.Fatalf("read schema: %v", err)
	}

	var doc snapshotDoc
	if err := yaml.Unmarshal(data, &doc); err != nil {
		t.Fatalf("unmarshal schema: %v", err)
	}
	if doc.OpenAPI == "" {
		t.Fatalf("missing openapi version")
	}
	if len(doc.Paths) == 0 {
		t.Fatalf("expected non-empty paths")
	}

	requiredPaths := []string{
		"/tailnet/{tailnet}/devices",
		"/device/{deviceId}",
		"/tailnet/{tailnet}/keys",
	}
	for _, p := range requiredPaths {
		if _, ok := doc.Paths[p]; !ok {
			t.Fatalf("required path %q missing from pinned schema", p)
		}
	}
}

func TestMappedOperationsExistInPinnedSchema(t *testing.T) {
	schemaBytes, err := os.ReadFile("tailscale-v2-openapi.yaml")
	if err != nil {
		t.Fatalf("read schema: %v", err)
	}
	mapBytes, err := os.ReadFile("command-operation-map.yaml")
	if err != nil {
		t.Fatalf("read mapping: %v", err)
	}

	var doc snapshotDoc
	if err := yaml.Unmarshal(schemaBytes, &doc); err != nil {
		t.Fatalf("unmarshal schema: %v", err)
	}
	var mapping opMap
	if err := yaml.Unmarshal(mapBytes, &mapping); err != nil {
		t.Fatalf("unmarshal mapping: %v", err)
	}

	for cmd, ops := range mapping.Commands {
		for _, op := range ops {
			parts := strings.SplitN(op, " ", 2)
			if len(parts) != 2 {
				t.Fatalf("invalid operation mapping %q for %q", op, cmd)
			}
			method := parts[0]
			path := parts[1]

			verbs, ok := doc.Paths[path]
			if !ok {
				t.Fatalf("mapped path %q for %q not in schema", path, cmd)
			}
			if _, ok := verbs[strings.ToLower(method)]; !ok {
				t.Fatalf("mapped method %q for %q not in schema path %q", method, cmd, path)
			}
		}
	}
}
