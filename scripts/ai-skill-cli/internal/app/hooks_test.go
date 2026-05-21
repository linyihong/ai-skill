package app

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestHooksInstallDryRunPlansWithoutWriting(t *testing.T) {
	repo := initTempGitRepo(t)
	writeFile(t, filepath.Join(repo, ".githooks", "pre-commit"), "#!/bin/sh\n")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"hooks", "install", "--repo", repo, "--dry-run", "--json"}, &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected success, got %d; stderr=%s", code, stderr.String())
	}

	var result Result
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}
	if result.Command != "hooks install" {
		t.Fatalf("unexpected command: %q", result.Command)
	}
	if len(result.PlannedActions) != 1 {
		t.Fatalf("expected one planned hook install, got %#v", result.PlannedActions)
	}
	if pathExists(filepath.Join(repo, ".git", "hooks", "pre-commit")) {
		t.Fatal("dry-run wrote hook target")
	}
}

func TestHooksInstallBlocksMissingSource(t *testing.T) {
	repo := initTempGitRepo(t)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"hooks", "install", "--repo", repo, "--dry-run", "--json"}, &stdout, &stderr)
	if code != ExitValidationFailed {
		t.Fatalf("expected validation failure, got %d; stderr=%s", code, stderr.String())
	}

	var result Result
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}
	if result.Error == nil || result.Error.Code != "missing_hook_source" {
		t.Fatalf("expected missing_hook_source, got %#v", result.Error)
	}
}

func TestHooksInstallBlocksExistingTargetWithoutForce(t *testing.T) {
	repo := initTempGitRepo(t)
	writeFile(t, filepath.Join(repo, ".githooks", "pre-commit"), "#!/bin/sh\n")
	writeFile(t, filepath.Join(repo, ".git", "hooks", "pre-commit"), "# existing\n")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"hooks", "install", "--repo", repo, "--dry-run", "--json"}, &stdout, &stderr)
	if code != ExitInvalidUsage {
		t.Fatalf("expected invalid usage for conflict, got %d; stderr=%s", code, stderr.String())
	}
}

func TestHooksInstallForceAllowsExistingTargetInDryRun(t *testing.T) {
	repo := initTempGitRepo(t)
	writeFile(t, filepath.Join(repo, ".githooks", "pre-commit"), "#!/bin/sh\n")
	writeFile(t, filepath.Join(repo, ".git", "hooks", "pre-commit"), "# existing\n")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"hooks", "install", "--repo", repo, "--dry-run", "--force", "--json"}, &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected success with force, got %d; stderr=%s", code, stderr.String())
	}
}

func TestHooksInstallReportsUnsafeGitStateButDoesNotBlockDryRun(t *testing.T) {
	repo := initTempGitRepo(t)
	writeFile(t, filepath.Join(repo, ".githooks", "pre-commit"), "#!/bin/sh\n")
	writeFile(t, filepath.Join(repo, ".git", "MERGE_HEAD"), "deadbeef\n")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"hooks", "install", "--repo", repo, "--dry-run", "--json"}, &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected dry-run success with warning, got %d; stderr=%s", code, stderr.String())
	}

	var result Result
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}
	if !hasCheckStatus(result.Checks, "git_operation", "warning") {
		t.Fatalf("expected git_operation warning, got %#v", result.Checks)
	}
}

func TestHooksInstallMissingGitBlocks(t *testing.T) {
	t.Setenv("PATH", emptyPathDir(t))

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"hooks", "install", "--repo", t.TempDir(), "--dry-run", "--json"}, &stdout, &stderr)
	if code != ExitMissingDependency {
		t.Fatalf("expected missing dependency, got %d; stderr=%s", code, stderr.String())
	}
}

func TestHooksInstallWriteModeBlockedUntilParity(t *testing.T) {
	repo := initTempGitRepo(t)
	writeFile(t, filepath.Join(repo, ".githooks", "pre-commit"), "#!/bin/sh\n")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"hooks", "install", "--repo", repo, "--json"}, &stdout, &stderr)
	if code != ExitPartialCloseBlocked {
		t.Fatalf("expected write mode blocked, got %d; stderr=%s", code, stderr.String())
	}
	if pathExists(filepath.Join(repo, ".git", "hooks", "pre-commit")) {
		t.Fatal("write-blocked mode wrote hook target")
	}
}

func TestHooksUnsupportedSubcommandReturnsInvalidUsage(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"hooks", "remove"}, &stdout, &stderr)
	if code != ExitInvalidUsage {
		t.Fatalf("expected invalid usage, got %d", code)
	}
}

func TestHooksInstallDoesNotWriteMutations(t *testing.T) {
	repo := initTempGitRepo(t)
	writeFile(t, filepath.Join(repo, ".githooks", "pre-commit"), "#!/bin/sh\n")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"hooks", "install", "--repo", repo, "--dry-run", "--json"}, &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected success, got %d; stderr=%s", code, stderr.String())
	}

	var result Result
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}
	if len(result.Mutations) != 0 {
		t.Fatalf("dry-run must not mutate, got %#v", result.Mutations)
	}
}

func TestListHookFilesIgnoresDirectories(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "pre-commit"), "#!/bin/sh\n")
	if err := os.Mkdir(filepath.Join(dir, "nested"), 0o755); err != nil {
		t.Fatal(err)
	}
	hooks := listHookFiles(dir)
	if len(hooks) != 1 || hooks[0] != "pre-commit" {
		t.Fatalf("expected only regular hook file, got %#v", hooks)
	}
}
