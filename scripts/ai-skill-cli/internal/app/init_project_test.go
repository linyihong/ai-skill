package app

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInitProjectDryRunPlansFilesWithoutWriting(t *testing.T) {
	project := t.TempDir()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"init-project", "--project", project, "--tools", "roo,cursor", "--dry-run", "--json"}, &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected success, got %d; stderr=%s", code, stderr.String())
	}

	var result Result
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}
	if result.Command != "init-project" || result.Mode != "dry_run" {
		t.Fatalf("unexpected result identity: %#v", result)
	}
	if len(result.Mutations) != 0 {
		t.Fatalf("dry-run must not mutate, got %#v", result.Mutations)
	}
	if len(result.PlannedActions) != 4 {
		t.Fatalf("expected 4 planned actions, got %#v", result.PlannedActions)
	}
	if pathExists(filepath.Join(project, ".roomodes")) {
		t.Fatal("dry-run wrote .roomodes")
	}
	if pathExists(filepath.Join(project, ".agent-goals")) {
		t.Fatal("dry-run wrote .agent-goals")
	}
}

func TestInitProjectBlocksExistingFileWithoutForce(t *testing.T) {
	project := t.TempDir()
	if err := os.WriteFile(filepath.Join(project, ".roomodes"), []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"init-project", "--project", project, "--tools", "roo", "--dry-run", "--json"}, &stdout, &stderr)
	if code != ExitInvalidUsage {
		t.Fatalf("expected invalid usage for conflict, got %d; stderr=%s", code, stderr.String())
	}

	var result Result
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}
	if result.Error == nil || result.Error.Code != "target_exists" {
		t.Fatalf("expected target_exists error, got %#v", result.Error)
	}
}

func TestInitProjectForceAllowsExistingFileInDryRun(t *testing.T) {
	project := t.TempDir()
	if err := os.WriteFile(filepath.Join(project, ".roomodes"), []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"init-project", "--project", project, "--tools", "roo", "--dry-run", "--force", "--json"}, &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected success with force, got %d; stderr=%s", code, stderr.String())
	}
}

func TestInitProjectWriteModeWritesSelectedFiles(t *testing.T) {
	project := t.TempDir()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"init-project", "--project", project, "--tools", "claude,cursor", "--json"}, &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected write success, got %d; stderr=%s", code, stderr.String())
	}

	var result Result
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}
	if result.Mode != "write" || len(result.Mutations) != 4 {
		t.Fatalf("expected write mutations, got %#v", result)
	}
	if !pathExists(filepath.Join(project, "CLAUDE.md")) {
		t.Fatal("write mode did not create CLAUDE.md")
	}
	if !pathExists(filepath.Join(project, ".cursor", "rules", "ai-skill-bootstrap.mdc")) {
		t.Fatal("write mode did not create Cursor rule")
	}
	if pathExists(filepath.Join(project, ".roomodes")) {
		t.Fatal("selected tools unexpectedly wrote .roomodes")
	}
	goals, err := os.ReadFile(filepath.Join(project, ".agent-goals", "README.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(goals), "ai-skill goals") {
		t.Fatalf("expected Go CLI goals guidance, got %s", string(goals))
	}
}

func TestInitProjectWritesCodexBootstrap(t *testing.T) {
	project := t.TempDir()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"init-project", "--project", project, "--tools", "codex", "--json"}, &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected write success, got %d; stderr=%s", code, stderr.String())
	}

	content, err := os.ReadFile(filepath.Join(project, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(content), "ai-tools/agent/codex.md") {
		t.Fatalf("expected Codex adapter pointer, got %s", string(content))
	}
}

func TestInitProjectPlainOutputIncludesPlannedActions(t *testing.T) {
	project := t.TempDir()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"init-project", "--project", project, "--tools", "claude", "--dry-run"}, &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected success, got %d; stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "Planned actions:") {
		t.Fatalf("expected planned actions in plain output, got %q", stdout.String())
	}
	if pathExists(filepath.Join(project, "CLAUDE.md")) {
		t.Fatal("dry-run wrote CLAUDE.md")
	}
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
