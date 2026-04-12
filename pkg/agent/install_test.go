package agent

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func TestInitWritesManagedAssets(t *testing.T) {
	root := testRoot()
	repo := t.TempDir()

	result, err := Init(root, InstallOptions{
		RootDir: repo,
	})
	if err != nil {
		t.Fatalf("init agent bundle: %v", err)
	}

	if result.CommandCount != 3 {
		t.Fatalf("expected 3 leaf commands, got %d", result.CommandCount)
	}

	manifestPath := filepath.Join(repo, filepath.FromSlash(localManifestRelPath))
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatalf("read manifest: %v", err)
	}

	var manifest Manifest
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		t.Fatalf("decode manifest: %v", err)
	}
	if manifest.ManagedBy != managedBy {
		t.Fatalf("expected managed-by %q, got %q", managedBy, manifest.ManagedBy)
	}
	if manifest.BundleVersion != bundleVersion {
		t.Fatalf("expected bundle version %q, got %q", bundleVersion, manifest.BundleVersion)
	}
	if manifest.InstallScope != ScopeLocal {
		t.Fatalf("expected install scope %q, got %q", ScopeLocal, manifest.InstallScope)
	}

	catalog, err := os.ReadFile(filepath.Join(repo, filepath.FromSlash(localCommandCatalogRelPath)))
	if err != nil {
		t.Fatalf("read catalog: %v", err)
	}
	body := string(catalog)
	for _, command := range []string{"tscli config show", "tscli list devices", "tscli set device"} {
		if !strings.Contains(body, command) {
			t.Fatalf("expected command catalog to include %q, got:\n%s", command, body)
		}
	}

	agents, err := os.ReadFile(filepath.Join(repo, localAgentsRelPath))
	if err != nil {
		t.Fatalf("read AGENTS.md: %v", err)
	}
	if !strings.Contains(string(agents), "tscli agent update") {
		t.Fatalf("expected AGENTS.md to mention refresh command, got:\n%s", string(agents))
	}
}

func TestInitWritesGlobalManagedAssets(t *testing.T) {
	root := testRoot()
	home := t.TempDir()
	t.Setenv("HOME", home)

	result, err := Init(root, InstallOptions{})
	if err != nil {
		t.Fatalf("init global agent bundle: %v", err)
	}

	if result.InstallScope != ScopeGlobal {
		t.Fatalf("expected global install scope, got %q", result.InstallScope)
	}
	if got := strings.Join(result.Tools, ","); got != strings.Join(defaultGlobalTools, ",") {
		t.Fatalf("expected default global tools %q, got %q", strings.Join(defaultGlobalTools, ","), got)
	}

	for _, rel := range []string{
		globalManifestRelPath,
		globalCommandCatalogRelPath,
		globalCodexSkillRelPath,
		globalClaudeInspectRelPath,
		globalClaudeOperateRelPath,
		globalOpenCodeInspectRelPath,
		globalOpenCodeOperateRelPath,
	} {
		if _, err := os.Stat(filepath.Join(home, filepath.FromSlash(rel))); err != nil {
			t.Fatalf("expected %s to exist: %v", rel, err)
		}
	}
}

func TestUpdateRefreshesManagedFilesFromManifest(t *testing.T) {
	root := testRoot()
	repo := t.TempDir()

	if _, err := Init(root, InstallOptions{
		RootDir: repo,
		Tools:   []string{ToolCodex},
	}); err != nil {
		t.Fatalf("init agent bundle: %v", err)
	}

	codexSkillPath := filepath.Join(repo, filepath.FromSlash(localCodexSkillRelPath))
	if err := os.WriteFile(codexSkillPath, []byte("stale"), 0o644); err != nil {
		t.Fatalf("write stale codex skill: %v", err)
	}

	result, err := Update(root, UpdateOptions{RootDir: repo})
	if err != nil {
		t.Fatalf("update agent bundle: %v", err)
	}

	if got := strings.Join(result.Tools, ","); got != ToolCodex {
		t.Fatalf("expected codex-only update, got %q", got)
	}

	updated, err := os.ReadFile(codexSkillPath)
	if err != nil {
		t.Fatalf("read updated codex skill: %v", err)
	}
	if !strings.Contains(string(updated), managedMarker) {
		t.Fatalf("expected managed codex skill content after update, got:\n%s", string(updated))
	}

	if _, err := os.Stat(filepath.Join(repo, filepath.FromSlash(localGitHubSkillRelPath))); !os.IsNotExist(err) {
		t.Fatalf("did not expect copilot assets to be created during codex-only update")
	}
}

func TestInitRemovesStaleManagedFilesWhenToolSelectionNarrows(t *testing.T) {
	root := testRoot()
	repo := t.TempDir()

	if _, err := Init(root, InstallOptions{RootDir: repo}); err != nil {
		t.Fatalf("init full agent bundle: %v", err)
	}

	stalePaths := []string{
		localAgentsRelPath,
		localClaudeRelPath,
		localGitHubSkillRelPath,
		localGitHubInspectRelPath,
		localGitHubOperateRelPath,
	}
	for _, rel := range stalePaths {
		if _, err := os.Stat(filepath.Join(repo, filepath.FromSlash(rel))); err != nil {
			t.Fatalf("expected %s to exist before narrowing tools: %v", rel, err)
		}
	}

	result, err := Init(root, InstallOptions{RootDir: repo, Tools: []string{ToolCodex}})
	if err != nil {
		t.Fatalf("re-init narrowed bundle: %v", err)
	}
	if got := strings.Join(result.Tools, ","); got != ToolCodex {
		t.Fatalf("expected codex-only re-init, got %q", got)
	}

	for _, rel := range stalePaths {
		if _, err := os.Stat(filepath.Join(repo, filepath.FromSlash(rel))); !os.IsNotExist(err) {
			t.Fatalf("expected stale managed file %s to be removed, got err=%v", rel, err)
		}
	}

	if _, err := os.Stat(filepath.Join(repo, filepath.FromSlash(localCodexSkillRelPath))); err != nil {
		t.Fatalf("expected codex skill to remain after narrowing tools: %v", err)
	}
}

func TestLoadManifestRejectsUnsupportedSchemaVersion(t *testing.T) {
	repo := t.TempDir()
	target := installTarget{
		Scope:              ScopeLocal,
		RootDir:            repo,
		ManifestRelPath:    localManifestRelPath,
		CommandCatalogPath: localCommandCatalogRelPath,
	}

	if err := os.MkdirAll(filepath.Join(repo, filepath.FromSlash(".tscli/agent")), 0o755); err != nil {
		t.Fatalf("mkdir manifest dir: %v", err)
	}
	content := strings.Join([]string{
		"managed-by: tscli-agent",
		"schema-version: 99",
		"bundle-version: vFuture",
		"install-scope: local",
		"tools: [codex]",
		"files: [.tscli/agent/commands.md, .codex/skills/tscli/SKILL.md, .tscli/agent/manifest.yaml]",
		"command-count: 1",
		"command-catalog: .tscli/agent/commands.md",
		"",
	}, "\n")
	if err := os.WriteFile(filepath.Join(repo, filepath.FromSlash(localManifestRelPath)), []byte(content), 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	_, err := loadManifest(target)
	if err == nil || !strings.Contains(err.Error(), "unsupported schema-version 99") {
		t.Fatalf("expected unsupported schema version error, got %v", err)
	}
}

func TestLoadManifestAllowsLegacySchemaVersionZero(t *testing.T) {
	repo := t.TempDir()
	target := installTarget{
		Scope:              ScopeLocal,
		RootDir:            repo,
		ManifestRelPath:    localManifestRelPath,
		CommandCatalogPath: localCommandCatalogRelPath,
	}

	if err := os.MkdirAll(filepath.Join(repo, filepath.FromSlash(".tscli/agent")), 0o755); err != nil {
		t.Fatalf("mkdir manifest dir: %v", err)
	}
	content := strings.Join([]string{
		"managed-by: tscli-agent",
		"schema-version: 0",
		"bundle-version: v1",
		"tools: [codex]",
		"files: [.tscli/agent/commands.md, .codex/skills/tscli/SKILL.md, .tscli/agent/manifest.yaml]",
		"command-count: 1",
		"command-catalog: .tscli/agent/commands.md",
		"",
	}, "\n")
	if err := os.WriteFile(filepath.Join(repo, filepath.FromSlash(localManifestRelPath)), []byte(content), 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	manifest, err := loadManifest(target)
	if err != nil {
		t.Fatalf("load legacy manifest: %v", err)
	}
	if manifest.SchemaVersion != 0 {
		t.Fatalf("expected legacy schema version 0, got %d", manifest.SchemaVersion)
	}
	if manifest.InstallScope != ScopeLocal {
		t.Fatalf("expected legacy manifest scope to default to %q, got %q", ScopeLocal, manifest.InstallScope)
	}
}

func testRoot() *cobra.Command {
	root := &cobra.Command{Use: "tscli"}

	configCmd := &cobra.Command{Use: "config"}
	configCmd.AddCommand(&cobra.Command{
		Use:   "show",
		Short: "Show tscli configuration",
		RunE:  func(*cobra.Command, []string) error { return nil },
	})

	listCmd := &cobra.Command{Use: "list"}
	listCmd.AddCommand(&cobra.Command{
		Use:   "devices",
		Short: "List devices",
		RunE:  func(*cobra.Command, []string) error { return nil },
	})

	setCmd := &cobra.Command{Use: "set"}
	setCmd.AddCommand(&cobra.Command{
		Use:   "device",
		Short: "Update a device",
		RunE:  func(*cobra.Command, []string) error { return nil },
	})

	root.AddCommand(configCmd, listCmd, setCmd)
	return root
}
