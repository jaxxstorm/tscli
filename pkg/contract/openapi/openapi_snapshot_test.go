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
	OpenAPI    string                    `yaml:"openapi"`
	Info       map[string]any            `yaml:"info"`
	Paths      map[string]map[string]any `yaml:"paths"`
	Components struct {
		Schemas map[string]any `yaml:"schemas"`
	} `yaml:"components"`
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

func TestPinnedSchemaIncludesDevicePostureIdentityProperty(t *testing.T) {
	doc := loadSnapshotDoc(t)

	verbs, ok := doc.Paths["/tailnet/{tailnet}/devices"]
	if !ok {
		t.Fatalf("devices path missing from schema")
	}
	getOp, ok := verbs["get"].(map[string]any)
	if !ok {
		t.Fatalf("devices GET operation missing from schema")
	}
	responses := getOp["responses"].(map[string]any)
	okResp := responses["200"].(map[string]any)
	content := okResp["content"].(map[string]any)
	appJSON := content["application/json"].(map[string]any)
	schema := appJSON["schema"].(map[string]any)
	properties := schema["properties"].(map[string]any)
	devices := properties["devices"].(map[string]any)
	items := devices["items"].(map[string]any)
	deviceRef := items["$ref"].(string)
	deviceSchema := resolveSchemaRef(t, doc, deviceRef)
	deviceProps := deviceSchema["properties"].(map[string]any)

	postureIdentity, ok := deviceProps["postureIdentity"].(map[string]any)
	if !ok {
		t.Fatalf("expected Device.postureIdentity in pinned schema")
	}
	postureProps := postureIdentity["properties"].(map[string]any)
	if _, ok := postureProps["serialNumbers"]; !ok {
		t.Fatalf("expected postureIdentity.serialNumbers in pinned schema")
	}
}

func TestPinnedSchemaIncludesCreateKeyTagsProperty(t *testing.T) {
	doc := loadSnapshotDoc(t)

	verbs, ok := doc.Paths["/tailnet/{tailnet}/keys"]
	if !ok {
		t.Fatalf("keys path missing from schema")
	}
	postOp, ok := verbs["post"].(map[string]any)
	if !ok {
		t.Fatalf("keys POST operation missing from schema")
	}
	reqBody := postOp["requestBody"].(map[string]any)
	content := reqBody["content"].(map[string]any)
	appJSON := content["application/json"].(map[string]any)
	schema := appJSON["schema"].(map[string]any)
	props := schema["properties"].(map[string]any)
	caps := resolveSchemaRef(t, doc, props["capabilities"].(map[string]any)["$ref"].(string))
	devices := caps["properties"].(map[string]any)["devices"].(map[string]any)
	create := devices["properties"].(map[string]any)["create"].(map[string]any)
	if _, ok := create["properties"].(map[string]any)["tags"]; !ok {
		t.Fatalf("expected capabilities.devices.create.tags in pinned schema")
	}
}

func loadSnapshotDoc(t *testing.T) snapshotDoc {
	t.Helper()
	data, err := os.ReadFile("tailscale-v2-openapi.yaml")
	if err != nil {
		t.Fatalf("read schema: %v", err)
	}
	var doc snapshotDoc
	if err := yaml.Unmarshal(data, &doc); err != nil {
		t.Fatalf("unmarshal schema: %v", err)
	}
	return doc
}

func resolveSchemaRef(t *testing.T, doc snapshotDoc, ref string) map[string]any {
	t.Helper()
	const prefix = "#/components/schemas/"
	if !strings.HasPrefix(ref, prefix) {
		t.Fatalf("unsupported schema ref %q", ref)
	}
	name := strings.TrimPrefix(ref, prefix)
	raw, ok := doc.Components.Schemas[name]
	if !ok {
		t.Fatalf("schema ref %q not found", ref)
	}
	schema, ok := raw.(map[string]any)
	if !ok {
		t.Fatalf("schema ref %q has unexpected shape", ref)
	}
	return schema
}
