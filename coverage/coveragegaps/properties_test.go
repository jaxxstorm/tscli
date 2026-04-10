package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExpandTypePropertyPaths(t *testing.T) {
	paths := expandTypePropertyPaths("", propertyTypeRegistry["apitype.DeviceListResponse"])

	mustContain := []string{
		"devices",
		"devices[]",
		"devices[].advertisedRoutes",
		"devices[].multipleConnections",
		"devices[].postureIdentity",
		"devices[].postureIdentity.serialNumbers",
		"devices[].postureIdentity.serialNumbers[]",
	}
	for _, item := range mustContain {
		if !contains(paths, item) {
			t.Fatalf("expected expanded paths to include %q, got %#v", item, paths)
		}
	}
}

func TestCollectOperationPropertyPaths(t *testing.T) {
	dir := t.TempDir()
	schemaPath := filepath.Join(dir, "openapi.yaml")
	schema := `
paths:
  /tailnet/{tailnet}/devices:
    get:
      responses:
        "200":
          content:
            application/json:
              schema:
                type: object
                properties:
                  devices:
                    type: array
                    items:
                      $ref: '#/components/schemas/Device'
components:
  schemas:
    Device:
      type: object
      properties:
        postureIdentity:
          type: object
          properties:
            serialNumbers:
              type: array
              items:
                type: string
`
	if err := os.WriteFile(schemaPath, []byte(schema), 0o644); err != nil {
		t.Fatalf("write schema: %v", err)
	}

	_, err := loadOperations(schemaPath)
	if err != nil {
		t.Fatalf("load operations: %v", err)
	}

	got, err := collectOperationPropertyPaths(&opsDoc, "get /tailnet/{tailnet}/devices", "response")
	if err != nil {
		t.Fatalf("collect operation properties: %v", err)
	}

	for _, want := range []string{
		"devices",
		"devices[]",
		"devices[].postureIdentity",
		"devices[].postureIdentity.serialNumbers",
		"devices[].postureIdentity.serialNumbers[]",
	} {
		if !contains(got, want) {
			t.Fatalf("expected property path %q in %#v", want, got)
		}
	}
}

func TestDerivePropertyCoverageAppliesDefaultExclusions(t *testing.T) {
	dir := t.TempDir()
	schemaPath := filepath.Join(dir, "openapi.yaml")
	schema := `
paths:
  /tailnet/{tailnet}/devices:
    get:
      responses:
        "200":
          content:
            application/json:
              schema:
                type: object
                properties:
                  devices:
                    type: array
                    items:
                      $ref: '#/components/schemas/Device'
components:
  schemas:
    Device:
      type: object
      properties:
        postureIdentity:
          type: object
          properties:
            serialNumbers:
              type: array
              items:
                type: string
`
	if err := os.WriteFile(schemaPath, []byte(schema), 0o644); err != nil {
		t.Fatalf("write schema: %v", err)
	}
	if _, err := loadOperations(schemaPath); err != nil {
		t.Fatalf("load operations: %v", err)
	}

	mapping := &commandMap{Commands: map[string][]string{
		"list devices": {"get /tailnet/{tailnet}/devices"},
	}}
	inventory, err := derivePropertyCoverage(
		&opsDoc,
		mapping,
		map[string]struct{}{},
		&propertyCoverageManifest{Operations: map[string]propertyCoverageOperation{}},
		&propertyExclusions{
			Default:    propertySideExclusion{Response: "audit pending"},
			Properties: map[string]string{},
		},
	)
	if err != nil {
		t.Fatalf("derive property coverage: %v", err)
	}
	if len(inventory.Uncovered) != 0 {
		t.Fatalf("expected no uncovered properties with default exclusions, got %#v", inventory.Uncovered)
	}
	if len(inventory.Excluded) == 0 {
		t.Fatalf("expected excluded properties, got none")
	}
}

func TestDerivePropertyCoverageSkipsUnknownMappedOperations(t *testing.T) {
	dir := t.TempDir()
	schemaPath := filepath.Join(dir, "openapi.yaml")
	schema := `
paths:
  /tailnet/{tailnet}/devices:
    get:
      responses:
        "200":
          content:
            application/json:
              schema:
                type: object
                properties:
                  devices:
                    type: array
                    items:
                      type: object
                      properties:
                        hostname:
                          type: string
`
	if err := os.WriteFile(schemaPath, []byte(schema), 0o644); err != nil {
		t.Fatalf("write schema: %v", err)
	}
	if _, err := loadOperations(schemaPath); err != nil {
		t.Fatalf("load operations: %v", err)
	}

	mapping := &commandMap{Commands: map[string][]string{
		"list devices": {"get /tailnet/{tailnet}/devices", "get /not-in-schema"},
	}}
	_, err := derivePropertyCoverage(
		&opsDoc,
		mapping,
		map[string]struct{}{},
		&propertyCoverageManifest{Operations: map[string]propertyCoverageOperation{
			"get /tailnet/{tailnet}/devices": {
				Response: []propertyCoverageEvidence{{Type: "apitype.DeviceListResponse"}},
			},
		}},
		&propertyExclusions{Properties: map[string]string{}},
	)
	if err != nil {
		t.Fatalf("derive property coverage should ignore unknown mapped operations, got %v", err)
	}
}

func contains(items []string, want string) bool {
	for _, item := range items {
		if item == want {
			return true
		}
	}
	return false
}
