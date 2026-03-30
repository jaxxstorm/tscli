package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

var operationMethods = map[string]struct{}{
	"get":     {},
	"post":    {},
	"put":     {},
	"patch":   {},
	"delete":  {},
	"options": {},
	"head":    {},
	"trace":   {},
}

type openapiDoc struct {
	OpenAPI string                    `yaml:"openapi"`
	Info    openapiInfo               `yaml:"info"`
	Paths   map[string]map[string]any `yaml:"paths"`
}

type openapiInfo struct {
	Version string `yaml:"version"`
}

type snapshotMetadata struct {
	SourceURL      string `yaml:"source_url"`
	FetchedAtUTC   string `yaml:"fetched_at_utc"`
	OpenAPIVersion string `yaml:"openapi_version"`
	APIVersion     string `yaml:"api_version"`
	PathCount      int    `yaml:"path_count"`
	OperationCount int    `yaml:"operation_count"`
	SHA256         string `yaml:"sha256"`
	SchemaFile     string `yaml:"schema_file"`
}

type fileReplacement struct {
	path   string
	tmp    string
	backup string
	exists bool
}

func main() {
	sourceURL := flag.String("source-url", "", "Canonical OpenAPI source URL")
	schemaOut := flag.String("schema-out", "pkg/contract/openapi/tailscale-v2-openapi.yaml", "Path for pinned OpenAPI schema")
	metadataOut := flag.String("metadata-out", "pkg/contract/openapi/snapshot-metadata.yaml", "Path for snapshot metadata")
	flag.Parse()

	if *sourceURL == "" {
		fatalf("source-url is required")
	}

	schema, err := fetchSchema(*sourceURL)
	if err != nil {
		fatalf("fetch schema: %v", err)
	}

	fetchedAt := time.Now().UTC()
	metadata, err := buildSnapshotMetadata(schema, *sourceURL, *schemaOut, fetchedAt)
	if err != nil {
		fatalf("build snapshot metadata: %v", err)
	}

	metadataBytes, err := yaml.Marshal(metadata)
	if err != nil {
		fatalf("marshal snapshot metadata: %v", err)
	}

	if err := replaceFilesAtomically(map[string][]byte{
		*schemaOut:   schema,
		*metadataOut: metadataBytes,
	}); err != nil {
		fatalf("persist refreshed snapshot: %v", err)
	}
}

func fetchSchema(sourceURL string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, sourceURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return nil, fmt.Errorf("unexpected status %s: %s", resp.Status, strings.TrimSpace(string(body)))
	}

	return io.ReadAll(resp.Body)
}

func buildSnapshotMetadata(schema []byte, sourceURL, schemaPath string, fetchedAt time.Time) (snapshotMetadata, error) {
	var doc openapiDoc
	if err := yaml.Unmarshal(schema, &doc); err != nil {
		return snapshotMetadata{}, err
	}
	if doc.OpenAPI == "" {
		return snapshotMetadata{}, fmt.Errorf("schema missing openapi version")
	}
	if len(doc.Paths) == 0 {
		return snapshotMetadata{}, fmt.Errorf("schema missing paths")
	}

	sum := sha256.Sum256(schema)
	return snapshotMetadata{
		SourceURL:      sourceURL,
		FetchedAtUTC:   fetchedAt.Format(time.RFC3339),
		OpenAPIVersion: doc.OpenAPI,
		APIVersion:     doc.Info.Version,
		PathCount:      len(doc.Paths),
		OperationCount: countOperations(doc.Paths),
		SHA256:         hex.EncodeToString(sum[:]),
		SchemaFile:     filepath.Base(schemaPath),
	}, nil
}

func countOperations(paths map[string]map[string]any) int {
	total := 0
	for _, verbs := range paths {
		for method := range verbs {
			if _, ok := operationMethods[strings.ToLower(method)]; ok {
				total++
			}
		}
	}
	return total
}

func replaceFilesAtomically(files map[string][]byte) error {
	keys := make([]string, 0, len(files))
	for path := range files {
		keys = append(keys, path)
	}
	slices.Sort(keys)

	replacements := make([]fileReplacement, 0, len(keys))
	for i, path := range keys {
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}

		tmp := filepath.Join(dir, fmt.Sprintf(".%s.tmp-%d-%d", filepath.Base(path), os.Getpid(), i))
		if err := os.WriteFile(tmp, files[path], 0o644); err != nil {
			cleanupTemps(replacements)
			return err
		}

		_, err := os.Stat(path)
		exists := err == nil
		if err != nil && !os.IsNotExist(err) {
			_ = os.Remove(tmp)
			cleanupTemps(replacements)
			return err
		}

		replacements = append(replacements, fileReplacement{
			path:   path,
			tmp:    tmp,
			backup: filepath.Join(dir, fmt.Sprintf(".%s.bak-%d-%d", filepath.Base(path), os.Getpid(), i)),
			exists: exists,
		})
	}

	replaced := 0
	for i := range replacements {
		r := &replacements[i]
		if r.exists {
			if err := os.Rename(r.path, r.backup); err != nil {
				cleanupTemps(replacements)
				return err
			}
		}
		if err := os.Rename(r.tmp, r.path); err != nil {
			if rollbackErr := rollbackReplacements(replacements, replaced); rollbackErr != nil {
				return fmt.Errorf("%w (rollback failed: %v)", err, rollbackErr)
			}
			return err
		}
		replaced++
	}

	for _, r := range replacements {
		if r.exists {
			if err := os.Remove(r.backup); err != nil && !os.IsNotExist(err) {
				return err
			}
		}
	}

	return nil
}

func cleanupTemps(replacements []fileReplacement) {
	for _, r := range replacements {
		_ = os.Remove(r.tmp)
	}
}

func rollbackReplacements(replacements []fileReplacement, replaced int) error {
	var errs []string

	for i := replaced - 1; i >= 0; i-- {
		r := replacements[i]
		if err := os.Remove(r.path); err != nil && !os.IsNotExist(err) {
			errs = append(errs, err.Error())
			continue
		}
		if r.exists {
			if err := os.Rename(r.backup, r.path); err != nil {
				errs = append(errs, err.Error())
			}
		}
	}

	for _, r := range replacements[replaced:] {
		_ = os.Remove(r.tmp)
		if r.exists {
			if err := os.Rename(r.backup, r.path); err != nil && !os.IsNotExist(err) {
				errs = append(errs, err.Error())
			}
		}
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}
	return nil
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
