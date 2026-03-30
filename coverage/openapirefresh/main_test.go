package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestBuildSnapshotMetadata(t *testing.T) {
	schema := []byte(`openapi: 3.1.0
info:
  version: v2
paths:
  /devices:
    get: {}
    post: {}
  /devices/{id}:
    get: {}
    delete: {}
    parameters: []
`)

	got, err := buildSnapshotMetadata(schema, "https://example.test/openapi.yaml", "pkg/contract/openapi/tailscale-v2-openapi.yaml", time.Date(2026, time.March, 29, 12, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("build snapshot metadata: %v", err)
	}

	if got.OpenAPIVersion != "3.1.0" {
		t.Fatalf("unexpected openapi version: %q", got.OpenAPIVersion)
	}
	if got.APIVersion != "v2" {
		t.Fatalf("unexpected api version: %q", got.APIVersion)
	}
	if got.PathCount != 2 {
		t.Fatalf("unexpected path count: %d", got.PathCount)
	}
	if got.OperationCount != 4 {
		t.Fatalf("unexpected operation count: %d", got.OperationCount)
	}
	if got.SchemaFile != "tailscale-v2-openapi.yaml" {
		t.Fatalf("unexpected schema file: %q", got.SchemaFile)
	}
	if got.SHA256 == "" {
		t.Fatalf("expected sha256 to be populated")
	}
}

func TestReplaceFilesAtomically(t *testing.T) {
	dir := t.TempDir()
	schemaPath := filepath.Join(dir, "schema.yaml")
	metaPath := filepath.Join(dir, "metadata.yaml")

	if err := os.WriteFile(schemaPath, []byte("old-schema"), 0o644); err != nil {
		t.Fatalf("write schema: %v", err)
	}
	if err := os.WriteFile(metaPath, []byte("old-meta"), 0o644); err != nil {
		t.Fatalf("write metadata: %v", err)
	}

	err := replaceFilesAtomically(map[string][]byte{
		schemaPath: []byte("new-schema"),
		metaPath:   []byte("new-meta"),
	})
	if err != nil {
		t.Fatalf("replace files atomically: %v", err)
	}

	schemaBytes, err := os.ReadFile(schemaPath)
	if err != nil {
		t.Fatalf("read schema: %v", err)
	}
	metaBytes, err := os.ReadFile(metaPath)
	if err != nil {
		t.Fatalf("read metadata: %v", err)
	}

	if string(schemaBytes) != "new-schema" {
		t.Fatalf("unexpected schema contents: %q", string(schemaBytes))
	}
	if string(metaBytes) != "new-meta" {
		t.Fatalf("unexpected metadata contents: %q", string(metaBytes))
	}
}
