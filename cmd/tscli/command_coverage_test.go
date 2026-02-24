package main

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestLeafCommandManifest(t *testing.T) {
	expected, err := loadLeafManifest()
	if err != nil {
		t.Fatalf("load manifest: %v", err)
	}

	actual := leafCommands(configureCLI())

	var missing []string
	for _, cmd := range expected {
		if !slices.Contains(actual, cmd) {
			missing = append(missing, cmd)
		}
	}

	var extra []string
	for _, cmd := range actual {
		if !slices.Contains(expected, cmd) {
			extra = append(extra, cmd)
		}
	}

	if len(missing) > 0 || len(extra) > 0 {
		t.Fatalf("leaf command manifest out of date\nmissing: %v\nextra: %v", missing, extra)
	}
}

func TestLeafCommandsHelpSmoke(t *testing.T) {
	commands, err := loadLeafManifest()
	if err != nil {
		t.Fatalf("load manifest: %v", err)
	}

	for _, command := range commands {
		t.Run(command, func(t *testing.T) {
			args := append(strings.Fields(command), "--help")
			res := executeCLI(t, args, nil)
			if res.err != nil {
				t.Fatalf("unexpected error for %q --help: %v\nstderr:\n%s", command, res.err, res.stderr)
			}
			if !strings.Contains(strings.ToLower(res.stdout), "usage") {
				t.Fatalf("expected usage in help output for %q, got:\n%s", command, res.stdout)
			}
		})
	}
}

func leafCommands(root *cobra.Command) []string {
	var out []string
	var walk func(*cobra.Command, []string)

	walk = func(cmd *cobra.Command, path []string) {
		children := nonBuiltinChildren(cmd)
		if len(children) == 0 {
			out = append(out, strings.Join(path, " "))
			return
		}
		for _, child := range children {
			walk(child, append(path, child.Name()))
		}
	}

	for _, child := range nonBuiltinChildren(root) {
		walk(child, []string{child.Name()})
	}

	slices.Sort(out)
	return out
}

func nonBuiltinChildren(cmd *cobra.Command) []*cobra.Command {
	children := cmd.Commands()
	filtered := make([]*cobra.Command, 0, len(children))
	for _, child := range children {
		switch child.Name() {
		case "help", "completion":
			continue
		default:
			filtered = append(filtered, child)
		}
	}
	return filtered
}

func loadLeafManifest() ([]string, error) {
	path := filepath.Join("testdata", "leaf_commands.txt")
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(b), "\n")
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		out = append(out, line)
	}
	return out, nil
}
