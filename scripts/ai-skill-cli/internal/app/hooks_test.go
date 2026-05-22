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
	if len(result.PlannedActions) != 3 {
		t.Fatalf("expected three planned hook installs, got %#v", result.PlannedActions)
	}
	if pathExists(filepath.Join(repo, ".git", "hooks", "pre-commit")) {
		t.Fatal("dry-run wrote hook target")
	}
}

func TestHooksInstallWriteModeInstallsAdapters(t *testing.T) {
	repo := initTempGitRepo(t)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"hooks", "install", "--repo", repo, "--json"}, &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected install success, got %d; stderr=%s stdout=%s", code, stderr.String(), stdout.String())
	}
	content, err := os.ReadFile(filepath.Join(repo, ".git", "hooks", "pre-commit"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(content, []byte("hooks run pre-commit")) {
		t.Fatalf("expected Go hook runner adapter, got %s", string(content))
	}
	prePushContent, err := os.ReadFile(filepath.Join(repo, ".git", "hooks", "pre-push"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(prePushContent, []byte("hooks run pre-push")) {
		t.Fatalf("expected Go pre-push hook runner adapter, got %s", string(prePushContent))
	}
}

func TestHooksInstallBlocksExistingTargetWithoutForce(t *testing.T) {
	repo := initTempGitRepo(t)
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

func TestHooksRunPostCommitReferenceOnly(t *testing.T) {
	repo := initTempGitRepo(t)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"hooks", "run", "post-commit", "--repo", repo, "--json"}, &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected post-commit success, got %d; stderr=%s stdout=%s", code, stderr.String(), stdout.String())
	}
}

func TestHooksRunPrePushSkipsWithoutCLIChanges(t *testing.T) {
	repo := initTempGitRepo(t)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"hooks", "run", "pre-push", "--repo", repo, "--json"}, &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected pre-push success, got %d; stderr=%s stdout=%s", code, stderr.String(), stdout.String())
	}
	var result Result
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}
	if !hasCheckStatus(result.Checks, "cli_ci_preflight", "skipped") {
		t.Fatalf("expected skipped cli_ci_preflight, got %#v", result.Checks)
	}
}

func TestHasCLICIPreflightChange(t *testing.T) {
	if !hasCLICIPreflightChange([]string{"scripts/ai-skill-cli/internal/app/hooks.go"}) {
		t.Fatal("expected CLI source to trigger preflight")
	}
	if !hasCLICIPreflightChange([]string{".github/workflows/ai-skill-cli.yml"}) {
		t.Fatal("expected workflow to trigger preflight")
	}
	if hasCLICIPreflightChange([]string{"scripts/README.md"}) {
		t.Fatal("scripts README alone should not trigger Go preflight")
	}
}

func TestParseGitHubRemote(t *testing.T) {
	cases := []string{
		"https://github.com/linyihong/Ai-skill.git",
		"git@github.com:linyihong/Ai-skill.git",
		"ssh://git@github.com/linyihong/Ai-skill.git",
	}
	for _, input := range cases {
		owner, repo, ok := parseGitHubRemote(input)
		if !ok || owner != "linyihong" || repo != "Ai-skill" {
			t.Fatalf("parseGitHubRemote(%q) = %q %q %v", input, owner, repo, ok)
		}
	}
	if _, _, ok := parseGitHubRemote("https://example.com/linyihong/Ai-skill.git"); ok {
		t.Fatal("non-GitHub remote should not parse")
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
