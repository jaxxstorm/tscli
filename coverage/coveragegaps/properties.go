package main

import (
	"fmt"
	"os"
	"path"
	"reflect"
	"slices"
	"strings"

	"gopkg.in/yaml.v3"
	tsapi "tailscale.com/client/tailscale/v2"
)

type propertyCoverageManifest struct {
	Operations map[string]propertyCoverageOperation `yaml:"operations"`
}

type propertyCoverageOperation struct {
	Request  []propertyCoverageEvidence `yaml:"request"`
	Response []propertyCoverageEvidence `yaml:"response"`
}

type propertyCoverageEvidence struct {
	Root     string         `yaml:"root"`
	Type     string         `yaml:"type"`
	Evidence map[string]any `yaml:"evidence,omitempty"`
}

type propertyExclusions struct {
	Default    propertySideExclusion `yaml:"default"`
	Properties map[string]string     `yaml:"properties"`
}

type propertySideExclusion struct {
	Request  string `yaml:"request"`
	Response string `yaml:"response"`
}

type propertyInventory struct {
	Covered              []string
	Excluded             []string
	Uncovered            []string
	UncoveredByOperation map[string][]string
	ExcludedByOperation  map[string][]string
	CoveredByOperation   map[string][]string
}

type propertyReference struct {
	Operation string
	Side      string
	Path      string
}

type openapiComponents struct {
	Schemas       map[string]any `yaml:"schemas"`
	Responses     map[string]any `yaml:"responses"`
	RequestBodies map[string]any `yaml:"requestBodies"`
}

type deviceListResponse struct {
	Devices []tsapi.Device `json:"devices"`
}

type createKeyOperationRequest struct {
	KeyType          string                `json:"keyType"`
	Capabilities     tsapi.KeyCapabilities `json:"capabilities"`
	ExpirySeconds    int64                 `json:"expirySeconds"`
	Description      string                `json:"description"`
	Scopes           []string              `json:"scopes"`
	Tags             []string              `json:"tags"`
	Issuer           string                `json:"issuer"`
	Subject          string                `json:"subject"`
	Audience         string                `json:"audience"`
	CustomClaimRules map[string]string     `json:"customClaimRules"`
}

var propertyTypeRegistry = map[string]reflect.Type{
	"tsapi.CreateKeyRequest":                reflect.TypeOf(tsapi.CreateKeyRequest{}),
	"tsapi.Device":                          reflect.TypeOf(tsapi.Device{}),
	"tsapi.Key":                             reflect.TypeOf(tsapi.Key{}),
	"tsapi.TailnetSettings":                 reflect.TypeOf(tsapi.TailnetSettings{}),
	"tsapi.UpdateTailnetSettingsRequest":    reflect.TypeOf(tsapi.UpdateTailnetSettingsRequest{}),
	"coverage.device_list_response":         reflect.TypeOf(deviceListResponse{}),
	"coverage.create_key_operation_request": reflect.TypeOf(createKeyOperationRequest{}),
}

func loadPropertyCoverage(path string) (*propertyCoverageManifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &propertyCoverageManifest{Operations: map[string]propertyCoverageOperation{}}, nil
		}
		return nil, err
	}
	var manifest propertyCoverageManifest
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		return nil, err
	}
	if manifest.Operations == nil {
		manifest.Operations = map[string]propertyCoverageOperation{}
	}
	return &manifest, nil
}

func loadPropertyExclusions(path string) (*propertyExclusions, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &propertyExclusions{Properties: map[string]string{}}, nil
		}
		return nil, err
	}
	var exclusions propertyExclusions
	if err := yaml.Unmarshal(data, &exclusions); err != nil {
		return nil, err
	}
	if exclusions.Properties == nil {
		exclusions.Properties = map[string]string{}
	}
	return &exclusions, nil
}

func derivePropertyCoverage(
	doc *openapiDoc,
	mapping *commandMap,
	excludedOps map[string]struct{},
	manifest *propertyCoverageManifest,
	exclusions *propertyExclusions,
) (propertyInventory, error) {
	ops := uniqueMappedOperations(mapping, excludedOps)

	coveredSet := map[string]struct{}{}
	excludedSet := map[string]struct{}{}
	uncoveredSet := map[string]struct{}{}

	for _, op := range ops {
		reqProps, err := collectOperationPropertyPaths(doc, op, "request")
		if err != nil {
			return propertyInventory{}, err
		}
		respProps, err := collectOperationPropertyPaths(doc, op, "response")
		if err != nil {
			return propertyInventory{}, err
		}

		reqCovered, err := expandCoverage(manifest.Operations[op].Request)
		if err != nil {
			return propertyInventory{}, fmt.Errorf("%s request coverage: %w", op, err)
		}
		respCovered, err := expandCoverage(manifest.Operations[op].Response)
		if err != nil {
			return propertyInventory{}, fmt.Errorf("%s response coverage: %w", op, err)
		}

		classifyProperties(op, "request", reqProps, reqCovered, exclusions, excludedSet, coveredSet, uncoveredSet)
		classifyProperties(op, "response", respProps, respCovered, exclusions, excludedSet, coveredSet, uncoveredSet)
	}

	return propertyInventory{
		Covered:              stringKeys(coveredSet),
		Excluded:             stringKeys(excludedSet),
		Uncovered:            stringKeys(uncoveredSet),
		CoveredByOperation:   groupPropertiesByOperation(coveredSet),
		ExcludedByOperation:  groupPropertiesByOperation(excludedSet),
		UncoveredByOperation: groupPropertiesByOperation(uncoveredSet),
	}, nil
}

func uniqueMappedOperations(mapping *commandMap, excludedOps map[string]struct{}) []string {
	seen := map[string]struct{}{}
	for _, ops := range mapping.Commands {
		for _, op := range ops {
			if _, excluded := excludedOps[op]; excluded {
				continue
			}
			seen[op] = struct{}{}
		}
	}
	return stringKeys(seen)
}

func expandCoverage(specs []propertyCoverageEvidence) (map[string]struct{}, error) {
	covered := map[string]struct{}{}
	for _, spec := range specs {
		t, ok := propertyTypeRegistry[spec.Type]
		if !ok {
			return nil, fmt.Errorf("unknown type %q", spec.Type)
		}
		for _, p := range expandTypePropertyPaths(spec.Root, t) {
			covered[p] = struct{}{}
		}
	}
	return covered, nil
}

func classifyProperties(
	operation, side string,
	props []string,
	covered map[string]struct{},
	exclusions *propertyExclusions,
	excludedSet, coveredSet, uncoveredSet map[string]struct{},
) {
	defaultReason := ""
	if side == "request" {
		defaultReason = exclusions.Default.Request
	} else {
		defaultReason = exclusions.Default.Response
	}

	hasCoverage := len(covered) > 0

	for _, prop := range props {
		ref := operation + " " + side + " " + prop
		if _, ok := covered[prop]; ok {
			coveredSet[ref] = struct{}{}
			continue
		}
		if reason, ok := matchPropertyExclusion(ref, exclusions.Properties); ok && reason != "" {
			excludedSet[ref] = struct{}{}
			continue
		}
		if !hasCoverage && defaultReason != "" {
			excludedSet[ref] = struct{}{}
			continue
		}
		uncoveredSet[ref] = struct{}{}
	}
}

func matchPropertyExclusion(ref string, exclusions map[string]string) (string, bool) {
	for pattern, reason := range exclusions {
		if pattern == ref {
			return reason, true
		}
		matched, err := path.Match(pattern, ref)
		if err == nil && matched {
			return reason, true
		}
	}
	return "", false
}

func groupPropertiesByOperation(items map[string]struct{}) map[string][]string {
	grouped := map[string][]string{}
	for ref := range items {
		op, rest, ok := strings.Cut(ref, " request ")
		if ok {
			grouped[op+" request"] = append(grouped[op+" request"], rest)
			continue
		}
		op, rest, ok = strings.Cut(ref, " response ")
		if ok {
			grouped[op+" response"] = append(grouped[op+" response"], rest)
		}
	}
	for key := range grouped {
		slices.Sort(grouped[key])
	}
	return grouped
}

func expandTypePropertyPaths(root string, t reflect.Type) []string {
	out := map[string]struct{}{}
	seen := map[reflect.Type]bool{}

	var walk func(prefix string, current reflect.Type)
	walk = func(prefix string, current reflect.Type) {
		current = indirectType(current)
		if current == nil {
			return
		}
		if prefix != "" {
			out[prefix] = struct{}{}
		}

		switch current.Kind() {
		case reflect.Struct:
			if seen[current] {
				return
			}
			seen[current] = true
			defer delete(seen, current)
			for i := 0; i < current.NumField(); i++ {
				field := current.Field(i)
				if !field.IsExported() {
					continue
				}
				name, ok := jsonFieldName(field)
				if !ok {
					continue
				}
				child := joinProperty(prefix, name)
				walk(child, field.Type)
			}
		case reflect.Slice, reflect.Array:
			elem := indirectType(current.Elem())
			arrayPrefix := prefix + "[]"
			out[arrayPrefix] = struct{}{}
			walk(arrayPrefix, elem)
		case reflect.Map:
			return
		default:
			return
		}
	}

	walk(root, t)
	return stringKeys(out)
}

func indirectType(t reflect.Type) reflect.Type {
	for t != nil && t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return t
}

func jsonFieldName(field reflect.StructField) (string, bool) {
	tag := field.Tag.Get("json")
	if tag == "-" {
		return "", false
	}
	if tag == "" {
		return "", false
	}
	name := strings.Split(tag, ",")[0]
	if name == "" {
		return "", false
	}
	return name, true
}

func joinProperty(prefix, name string) string {
	if prefix == "" {
		return name
	}
	if strings.HasSuffix(prefix, "[]") {
		return prefix + "." + name
	}
	return prefix + "." + name
}

func collectOperationPropertyPaths(doc *openapiDoc, opKey, side string) ([]string, error) {
	method, apiPath, ok := strings.Cut(opKey, " ")
	if !ok {
		return nil, fmt.Errorf("invalid operation key %q", opKey)
	}
	verbs, ok := doc.Paths[apiPath]
	if !ok {
		return nil, fmt.Errorf("path %q not found in schema", apiPath)
	}
	opRaw, ok := verbs[strings.ToLower(method)]
	if !ok {
		return nil, fmt.Errorf("operation %q not found in schema", opKey)
	}
	op, ok := opRaw.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("operation %q has unexpected shape", opKey)
	}

	var schema any
	var err error
	if side == "request" {
		schema, err = requestSchema(doc, op)
	} else {
		schema, err = responseSchema(doc, op)
	}
	if err != nil || schema == nil {
		return nil, err
	}

	paths := map[string]struct{}{}
	if err := collectSchemaPropertyPaths(doc, schema, "", paths); err != nil {
		return nil, err
	}
	return stringKeys(paths), nil
}

func requestSchema(doc *openapiDoc, op map[string]any) (any, error) {
	raw, ok := op["requestBody"]
	if !ok {
		return nil, nil
	}
	raw, err := resolveMaybeRef(doc, raw)
	if err != nil {
		return nil, err
	}
	body, ok := raw.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("requestBody has unexpected shape")
	}
	return contentSchema(body)
}

func responseSchema(doc *openapiDoc, op map[string]any) (any, error) {
	raw, ok := op["responses"]
	if !ok {
		return nil, nil
	}
	responses, ok := raw.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("responses have unexpected shape")
	}
	var codes []string
	for code := range responses {
		if strings.HasPrefix(code, "2") {
			codes = append(codes, code)
		}
	}
	slices.Sort(codes)
	for _, code := range codes {
		resolved, err := resolveMaybeRef(doc, responses[code])
		if err != nil {
			return nil, err
		}
		response, ok := resolved.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("response %s has unexpected shape", code)
		}
		schema, err := contentSchema(response)
		if err != nil {
			return nil, err
		}
		if schema != nil {
			return schema, nil
		}
	}
	return nil, nil
}

func contentSchema(container map[string]any) (any, error) {
	contentRaw, ok := container["content"]
	if !ok {
		return nil, nil
	}
	content, ok := contentRaw.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("content has unexpected shape")
	}
	for _, mediaType := range []string{"application/json", "application/merge-patch+json"} {
		item, ok := content[mediaType]
		if !ok {
			continue
		}
		typed, ok := item.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("%s content has unexpected shape", mediaType)
		}
		return typed["schema"], nil
	}
	return nil, nil
}

func resolveMaybeRef(doc *openapiDoc, raw any) (any, error) {
	node, ok := raw.(map[string]any)
	if !ok {
		return raw, nil
	}
	ref, ok := node["$ref"].(string)
	if !ok {
		return raw, nil
	}
	return resolveRef(doc, ref)
}

func resolveRef(doc *openapiDoc, ref string) (any, error) {
	const prefix = "#/components/"
	if !strings.HasPrefix(ref, prefix) {
		return nil, fmt.Errorf("unsupported ref %q", ref)
	}
	parts := strings.Split(strings.TrimPrefix(ref, prefix), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("unsupported ref %q", ref)
	}
	section, name := parts[0], parts[1]
	switch section {
	case "schemas":
		raw, ok := doc.Components.Schemas[name]
		if !ok {
			return nil, fmt.Errorf("schema ref %q not found", ref)
		}
		return raw, nil
	case "responses":
		raw, ok := doc.Components.Responses[name]
		if !ok {
			return nil, fmt.Errorf("response ref %q not found", ref)
		}
		return raw, nil
	case "requestBodies":
		raw, ok := doc.Components.RequestBodies[name]
		if !ok {
			return nil, fmt.Errorf("requestBody ref %q not found", ref)
		}
		return raw, nil
	default:
		return nil, fmt.Errorf("unsupported ref section %q", ref)
	}
}

func collectSchemaPropertyPaths(doc *openapiDoc, raw any, prefix string, out map[string]struct{}) error {
	resolved, err := resolveMaybeRef(doc, raw)
	if err != nil {
		return err
	}
	node, ok := resolved.(map[string]any)
	if !ok {
		return nil
	}

	for _, keyword := range []string{"allOf", "oneOf", "anyOf"} {
		if items, ok := node[keyword].([]any); ok {
			for _, item := range items {
				if err := collectSchemaPropertyPaths(doc, item, prefix, out); err != nil {
					return err
				}
			}
		}
	}

	if items, ok := node["items"]; ok {
		arrayPrefix := prefix + "[]"
		if prefix != "" {
			out[arrayPrefix] = struct{}{}
		}
		return collectSchemaPropertyPaths(doc, items, arrayPrefix, out)
	}

	propsRaw, hasProps := node["properties"]
	if !hasProps {
		return nil
	}
	props, ok := propsRaw.(map[string]any)
	if !ok {
		return nil
	}
	names := make([]string, 0, len(props))
	for name := range props {
		names = append(names, name)
	}
	slices.Sort(names)
	for _, name := range names {
		child := joinProperty(prefix, name)
		out[child] = struct{}{}
		if err := collectSchemaPropertyPaths(doc, props[name], child, out); err != nil {
			return err
		}
	}
	return nil
}
