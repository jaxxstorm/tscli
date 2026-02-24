package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadOperations(t *testing.T) {
	dir := t.TempDir()
	schema := filepath.Join(dir, "openapi.yaml")
	content := []byte("paths:\n  /foo:\n    get: {}\n    post: {}\n  /bar:\n    patch: {}\n    options: {}\n")
	if err := os.WriteFile(schema, content, 0o644); err != nil {
		t.Fatalf("write schema: %v", err)
	}

	ops, err := loadOperations(schema)
	if err != nil {
		t.Fatalf("load operations: %v", err)
	}
	if len(ops) != 3 {
		t.Fatalf("expected 3 operations, got %d (%v)", len(ops), ops)
	}
	if ops[0].Key == "" || ops[0].Domain == "" {
		t.Fatalf("expected operation key and domain to be populated, got %#v", ops[0])
	}
}

func TestLoadManifest(t *testing.T) {
	dir := t.TempDir()
	manifest := filepath.Join(dir, "manifest.txt")
	content := []byte("# comment\n\nlist devices\nget device\n")
	if err := os.WriteFile(manifest, content, 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	items, err := loadManifest(manifest)
	if err != nil {
		t.Fatalf("load manifest: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}
}

func TestLoadExclusionsMissingFile(t *testing.T) {
	p, err := loadExclusions(filepath.Join(t.TempDir(), "does-not-exist.yaml"))
	if err != nil {
		t.Fatalf("load exclusions: %v", err)
	}
	if len(p.Commands) != 0 || len(p.Operations) != 0 {
		t.Fatalf("expected empty exclusions, got %#v", p)
	}
}
