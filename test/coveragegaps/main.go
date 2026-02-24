package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"gopkg.in/yaml.v3"
)

var methods = map[string]struct{}{
	"get":    {},
	"post":   {},
	"put":    {},
	"patch":  {},
	"delete": {},
}

type openapiDoc struct {
	Paths map[string]map[string]any `yaml:"paths"`
}

type commandMap struct {
	Commands map[string][]string `yaml:"commands"`
}

type report struct {
	OpenAPIOperations int               `json:"openapi_operations"`
	ManifestCommands  int               `json:"manifest_commands"`
	CoveredOperations []string          `json:"covered_operations"`
	UncoveredOps      []string          `json:"uncovered_operations"`
	UnknownMappedOps  []string          `json:"unknown_mapped_operations"`
	CoveredCommands   []string          `json:"covered_commands"`
	UnmappedCommands  []string          `json:"unmapped_commands"`
	UnknownCommands   []string          `json:"unknown_mapped_commands"`
	CommandMappings   map[string]string `json:"command_mappings"`
}

func main() {
	openapiPath := flag.String("openapi", "pkg/contract/openapi/tailscale-v2-openapi.yaml", "Path to pinned OpenAPI schema")
	mappingPath := flag.String("mapping", "pkg/contract/openapi/command-operation-map.yaml", "Path to command->operation map")
	manifestPath := flag.String("manifest", "cmd/tscli/testdata/leaf_commands.txt", "Path to leaf command manifest")
	jsonOut := flag.String("json-out", "coverage/coverage-gaps.json", "Path for machine-readable report")
	mdOut := flag.String("md-out", "coverage/coverage-gaps.md", "Path for markdown report")
	baselinePath := flag.String("baseline", "coverage/coverage-gaps-baseline.json", "Path to baseline report for diffing")
	diffOut := flag.String("diff-out", "coverage/coverage-gaps-diff.md", "Path for baseline diff report")
	failOnRegression := flag.Bool("fail-on-regression", false, "Exit non-zero if uncovered operations or unmapped commands regress vs baseline")
	flag.Parse()

	ops, err := loadOperations(*openapiPath)
	if err != nil {
		fatalf("load OpenAPI operations: %v", err)
	}
	mapping, err := loadCommandMap(*mappingPath)
	if err != nil {
		fatalf("load command map: %v", err)
	}
	manifest, err := loadManifest(*manifestPath)
	if err != nil {
		fatalf("load manifest: %v", err)
	}

	opSet := make(map[string]struct{}, len(ops))
	for _, op := range ops {
		opSet[op] = struct{}{}
	}
	manifestSet := make(map[string]struct{}, len(manifest))
	for _, cmd := range manifest {
		manifestSet[cmd] = struct{}{}
	}

	coveredOpSet := map[string]struct{}{}
	unknownMappedSet := map[string]struct{}{}
	coveredCommands := map[string]struct{}{}
	unknownCommandSet := map[string]struct{}{}
	commandMapping := map[string]string{}

	for cmd, mapped := range mapping.Commands {
		if len(mapped) == 0 {
			continue
		}
		if _, ok := manifestSet[cmd]; ok {
			coveredCommands[cmd] = struct{}{}
		} else {
			unknownCommandSet[cmd] = struct{}{}
		}
		commandMapping[cmd] = strings.Join(mapped, ", ")
		for _, op := range mapped {
			if _, ok := opSet[op]; ok {
				coveredOpSet[op] = struct{}{}
			} else {
				unknownMappedSet[op] = struct{}{}
			}
		}
	}

	unmappedCommands := make([]string, 0)
	for _, cmd := range manifest {
		if _, ok := coveredCommands[cmd]; !ok {
			unmappedCommands = append(unmappedCommands, cmd)
		}
	}
	slices.Sort(unmappedCommands)

	coveredOps := stringKeys(coveredOpSet)
	unknownMappedOps := stringKeys(unknownMappedSet)
	uncoveredOps := make([]string, 0, len(ops))
	for _, op := range ops {
		if _, ok := coveredOpSet[op]; !ok {
			uncoveredOps = append(uncoveredOps, op)
		}
	}

	rep := report{
		OpenAPIOperations: len(ops),
		ManifestCommands:  len(manifest),
		CoveredOperations: coveredOps,
		UncoveredOps:      uncoveredOps,
		UnknownMappedOps:  unknownMappedOps,
		CoveredCommands:   stringKeys(coveredCommands),
		UnmappedCommands:  unmappedCommands,
		UnknownCommands:   stringKeys(unknownCommandSet),
		CommandMappings:   commandMapping,
	}

	if err := os.MkdirAll(filepath.Dir(*jsonOut), 0o755); err != nil {
		fatalf("create output directory: %v", err)
	}
	if err := writeJSON(*jsonOut, rep); err != nil {
		fatalf("write json report: %v", err)
	}
	if err := writeMarkdown(*mdOut, rep); err != nil {
		fatalf("write markdown report: %v", err)
	}

	diff, err := diffAgainstBaseline(*baselinePath, rep)
	if err != nil {
		fatalf("diff baseline: %v", err)
	}
	if err := os.WriteFile(*diffOut, []byte(diff.markdown), 0o644); err != nil {
		fatalf("write diff report: %v", err)
	}
	if *failOnRegression && (len(diff.newUncoveredOps) > 0 || len(diff.newUnmappedCommands) > 0) {
		fatalf("coverage regression: %d new uncovered operations, %d new unmapped commands",
			len(diff.newUncoveredOps), len(diff.newUnmappedCommands))
	}
}

func loadOperations(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var doc openapiDoc
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return nil, err
	}

	ops := make([]string, 0)
	for p, verbs := range doc.Paths {
		for method := range verbs {
			method = strings.ToLower(method)
			if _, ok := methods[method]; !ok {
				continue
			}
			ops = append(ops, method+" "+p)
		}
	}
	slices.Sort(ops)
	return ops, nil
}

func loadCommandMap(path string) (*commandMap, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var m commandMap
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	if m.Commands == nil {
		m.Commands = map[string][]string{}
	}
	return &m, nil
}

func loadManifest(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(data), "\n")
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		out = append(out, line)
	}
	slices.Sort(out)
	return out, nil
}

func writeJSON(path string, rep report) error {
	data, err := json.MarshalIndent(rep, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0o644)
}

func writeMarkdown(path string, rep report) error {
	var b strings.Builder
	b.WriteString("# Coverage Gaps Report\n\n")
	b.WriteString(fmt.Sprintf("- OpenAPI operations: `%d`\n", rep.OpenAPIOperations))
	b.WriteString(fmt.Sprintf("- Manifest commands: `%d`\n", rep.ManifestCommands))
	b.WriteString(fmt.Sprintf("- Covered operations: `%d`\n", len(rep.CoveredOperations)))
	b.WriteString(fmt.Sprintf("- Uncovered operations: `%d`\n", len(rep.UncoveredOps)))
	b.WriteString(fmt.Sprintf("- Covered commands: `%d`\n", len(rep.CoveredCommands)))
	b.WriteString(fmt.Sprintf("- Unmapped commands: `%d`\n", len(rep.UnmappedCommands)))
	b.WriteString(fmt.Sprintf("- Unknown mapped commands: `%d`\n", len(rep.UnknownCommands)))
	b.WriteString("\n## Unmapped Commands\n\n")
	for _, cmd := range rep.UnmappedCommands {
		b.WriteString("- `" + cmd + "`\n")
	}
	b.WriteString("\n## Unknown Mapped Operations\n\n")
	if len(rep.UnknownMappedOps) == 0 {
		b.WriteString("- None\n")
	} else {
		for _, op := range rep.UnknownMappedOps {
			b.WriteString("- `" + op + "`\n")
		}
	}
	b.WriteString("\n## Unknown Mapped Commands\n\n")
	if len(rep.UnknownCommands) == 0 {
		b.WriteString("- None\n")
	} else {
		for _, cmd := range rep.UnknownCommands {
			b.WriteString("- `" + cmd + "`\n")
		}
	}
	return os.WriteFile(path, []byte(b.String()), 0o644)
}

type baselineDiff struct {
	newUncoveredOps     []string
	newUnmappedCommands []string
	markdown            string
}

func diffAgainstBaseline(path string, current report) (baselineDiff, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return baselineDiff{markdown: "# Coverage Baseline Diff\n\n- No baseline file found.\n"}, nil
		}
		return baselineDiff{}, err
	}

	var base report
	if err := json.Unmarshal(data, &base); err != nil {
		return baselineDiff{}, err
	}

	baseUncovered := make(map[string]struct{}, len(base.UncoveredOps))
	for _, op := range base.UncoveredOps {
		baseUncovered[op] = struct{}{}
	}
	baseUnmapped := make(map[string]struct{}, len(base.UnmappedCommands))
	for _, cmd := range base.UnmappedCommands {
		baseUnmapped[cmd] = struct{}{}
	}

	var newOps []string
	for _, op := range current.UncoveredOps {
		if _, ok := baseUncovered[op]; !ok {
			newOps = append(newOps, op)
		}
	}
	var newUnmapped []string
	for _, cmd := range current.UnmappedCommands {
		if _, ok := baseUnmapped[cmd]; !ok {
			newUnmapped = append(newUnmapped, cmd)
		}
	}
	slices.Sort(newOps)
	slices.Sort(newUnmapped)

	var b strings.Builder
	b.WriteString("# Coverage Baseline Diff\n\n")
	b.WriteString(fmt.Sprintf("- New uncovered operations: `%d`\n", len(newOps)))
	b.WriteString(fmt.Sprintf("- New unmapped commands: `%d`\n", len(newUnmapped)))
	b.WriteString("\n## New Uncovered Operations\n\n")
	if len(newOps) == 0 {
		b.WriteString("- None\n")
	} else {
		for _, op := range newOps {
			b.WriteString("- `" + op + "`\n")
		}
	}
	b.WriteString("\n## New Unmapped Commands\n\n")
	if len(newUnmapped) == 0 {
		b.WriteString("- None\n")
	} else {
		for _, cmd := range newUnmapped {
			b.WriteString("- `" + cmd + "`\n")
		}
	}

	return baselineDiff{
		newUncoveredOps:     newOps,
		newUnmappedCommands: newUnmapped,
		markdown:            b.String(),
	}, nil
}

func stringKeys(m map[string]struct{}) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	slices.Sort(out)
	return out
}

func fatalf(format string, args ...any) {
	_, _ = fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
