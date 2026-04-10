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
	deviceProps := loadDeviceSchemaProperties(t, doc)

	postureIdentity, ok := deviceProps["postureIdentity"].(map[string]any)
	if !ok {
		t.Fatalf("expected Device.postureIdentity in pinned schema")
	}
	postureProps := requireMap(t, postureIdentity, "Device.postureIdentity.properties", "properties")
	if _, ok := postureProps["serialNumbers"]; !ok {
		t.Fatalf("expected postureIdentity.serialNumbers in pinned schema")
	}
}

func TestPinnedSchemaIncludesDeviceAdvancedRouteProperties(t *testing.T) {
	doc := loadSnapshotDoc(t)
	deviceProps := loadDeviceSchemaProperties(t, doc)

	if _, ok := deviceProps["advertisedRoutes"]; !ok {
		t.Fatalf("expected Device.advertisedRoutes in pinned schema")
	}
	if _, ok := deviceProps["multipleConnections"]; !ok {
		t.Fatalf("expected Device.multipleConnections in pinned schema")
	}
}

func TestPinnedSchemaIncludesDeviceRoutesSchemaProperties(t *testing.T) {
	doc := loadSnapshotDoc(t)
	routesSchema := resolveSchemaRef(t, doc, "#/components/schemas/DeviceRoutes")
	props := requireMap(t, routesSchema, "DeviceRoutes", "properties")
	for _, key := range []string{"advertisedRoutes", "enabledRoutes"} {
		if _, ok := props[key]; !ok {
			t.Fatalf("expected DeviceRoutes.%s in pinned schema", key)
		}
	}
}

func TestPinnedSchemaIncludesTailnetSettingsPostureIdentityCollectionProperty(t *testing.T) {
	doc := loadSnapshotDoc(t)
	settingsSchema := resolveSchemaRef(t, doc, "#/components/schemas/TailnetSettings")
	props := requireMap(t, settingsSchema, "TailnetSettings", "properties")
	if _, ok := props["postureIdentityCollectionOn"]; !ok {
		t.Fatalf("expected TailnetSettings.postureIdentityCollectionOn in pinned schema")
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
	reqBody := requireMap(t, postOp, "keys POST operation", "requestBody")
	content := requireMap(t, reqBody, "keys POST requestBody", "content")
	appJSON, ok := content["application/json"].(map[string]any)
	if !ok {
		t.Fatalf("keys POST requestBody.content.application/json has unexpected shape")
	}
	schema, ok := appJSON["schema"].(map[string]any)
	if !ok {
		t.Fatalf("keys POST request schema has unexpected shape")
	}
	props := requireMap(t, schema, "keys POST request schema", "properties")
	capabilityRef, ok := props["capabilities"].(map[string]any)
	if !ok {
		t.Fatalf("keys POST request schema properties.capabilities has unexpected shape")
	}
	ref, ok := capabilityRef["$ref"].(string)
	if !ok {
		t.Fatalf("keys POST request schema properties.capabilities.$ref missing")
	}
	caps := resolveSchemaRef(t, doc, ref)
	devices := requireMap(t, requireMap(t, caps, "KeyCapabilities", "properties"), "KeyCapabilities.properties", "devices")
	create := requireMap(t, requireMap(t, devices, "KeyCapabilities.properties.devices", "properties"), "KeyCapabilities.properties.devices.properties", "create")
	if _, ok := requireMap(t, create, "KeyCapabilities.properties.devices.properties.create", "properties")["tags"]; !ok {
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

func loadDeviceSchemaProperties(t *testing.T, doc snapshotDoc) map[string]any {
	t.Helper()

	verbs, ok := doc.Paths["/tailnet/{tailnet}/devices"]
	if !ok {
		t.Fatalf("devices path missing from schema")
	}
	getOp, ok := verbs["get"].(map[string]any)
	if !ok {
		t.Fatalf("devices GET operation missing from schema")
	}
	responses := requireMap(t, getOp, "devices GET operation", "responses")
	okResp, ok := responses["200"].(map[string]any)
	if !ok {
		t.Fatalf("devices GET operation responses.200 has unexpected shape")
	}
	content := requireMap(t, okResp, "devices GET operation responses.200", "content")
	appJSON, ok := content["application/json"].(map[string]any)
	if !ok {
		t.Fatalf("devices GET operation responses.200.content.application/json has unexpected shape")
	}
	schema, ok := appJSON["schema"].(map[string]any)
	if !ok {
		t.Fatalf("devices GET operation response schema has unexpected shape")
	}
	properties := requireMap(t, schema, "devices GET response schema", "properties")
	devices, ok := properties["devices"].(map[string]any)
	if !ok {
		t.Fatalf("devices GET response schema properties.devices has unexpected shape")
	}
	items, ok := devices["items"].(map[string]any)
	if !ok {
		t.Fatalf("devices GET response schema properties.devices.items has unexpected shape")
	}
	deviceRef, ok := items["$ref"].(string)
	if !ok {
		t.Fatalf("devices GET response schema properties.devices.items.$ref missing")
	}
	deviceSchema := resolveSchemaRef(t, doc, deviceRef)
	return requireMap(t, deviceSchema, "Device schema", "properties")
}

func requireMap(t *testing.T, parent map[string]any, label, key string) map[string]any {
	t.Helper()
	raw, ok := parent[key]
	if !ok {
		t.Fatalf("%s missing %q", label, key)
	}
	out, ok := raw.(map[string]any)
	if !ok {
		t.Fatalf("%s.%s has unexpected shape", label, key)
	}
	return out
}
