package app

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

func TestCloseLoopCleanTempRepoDryRun(t *testing.T) {
	repo := initTempGitRepo(t)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"close-loop", "--repo", repo, "--dry-run", "--json"}, &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected success, got %d; stderr=%s", code, stderr.String())
	}

	var result Result
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}
	if !hasCheckStatus(result.Checks, "working_tree", "clean") {
		t.Fatalf("expected clean working tree, got %#v", result.Checks)
	}
	if len(result.Mutations) != 0 {
		t.Fatalf("dry-run must not mutate, got %#v", result.Mutations)
	}
}

func TestCloseLoopGroupsOwnedDirtyPaths(t *testing.T) {
	repo := initTempGitRepo(t)
	writeFile(t, filepath.Join(repo, "scripts", "tool.sh"), "echo ok\n")
	writeFile(t, filepath.Join(repo, "workflow", "note.md"), "# Note\n")
	writeFile(t, filepath.Join(repo, "plans", "active", "plan.md"), "# Plan\n")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"close-loop", "--repo", repo, "--dry-run", "--json"}, &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected success, got %d; stderr=%s", code, stderr.String())
	}

	var result Result
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}
	if !hasCheckMessage(result.Checks, "owner_group", "scripts: scripts/tool.sh") {
		t.Fatalf("expected scripts owner group, got %#v", result.Checks)
	}
	if !hasCheckMessage(result.Checks, "owner_group", "architecture: plans/active/plan.md, workflow/note.md") {
		t.Fatalf("expected plans to share architecture owner group, got %#v", result.Checks)
	}
}

func TestCloseLoopBlocksUnknownDirtyPath(t *testing.T) {
	repo := initTempGitRepo(t)
	writeFile(t, filepath.Join(repo, "scratch.txt"), "temporary\n")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"close-loop", "--repo", repo, "--dry-run", "--json"}, &stdout, &stderr)
	if code != ExitValidationFailed {
		t.Fatalf("expected validation failure, got %d; stderr=%s", code, stderr.String())
	}

	var result Result
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}
	if result.Error == nil || result.Error.Code != "unknown_owner_group" {
		t.Fatalf("expected unknown_owner_group, got %#v", result.Error)
	}
}

func TestCloseLoopBlocksMergeState(t *testing.T) {
	repo := initTempGitRepo(t)
	gitDir := gitOutput(t, repo, "rev-parse", "--git-dir")
	if !filepath.IsAbs(gitDir) {
		gitDir = filepath.Join(repo, gitDir)
	}
	writeFile(t, filepath.Join(gitDir, "MERGE_HEAD"), "deadbeef\n")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"close-loop", "--repo", repo, "--dry-run", "--json"}, &stdout, &stderr)
	if code != ExitUnsafeRepoState {
		t.Fatalf("expected unsafe repo state, got %d; stderr=%s", code, stderr.String())
	}
}

func TestCloseLoopBlocksActiveLock(t *testing.T) {
	repo := initTempGitRepo(t)
	lockDir := filepath.Join(repo, ".git", "ai-skill-agent.lock")
	if err := os.MkdirAll(lockDir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(lockDir, "owner"), "other-agent\n")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"close-loop", "--repo", repo, "--dry-run", "--json"}, &stdout, &stderr)
	if code != ExitUnsafeRepoState {
		t.Fatalf("expected unsafe repo state, got %d; stderr=%s", code, stderr.String())
	}

	var result Result
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}
	if result.Error == nil || result.Error.Code != "active_close_loop_lock" {
		t.Fatalf("expected active_close_loop_lock, got %#v", result.Error)
	}
}

func TestCloseLoopBlocksRebaseState(t *testing.T) {
	repo := initTempGitRepo(t)
	gitDir := gitOutput(t, repo, "rev-parse", "--git-dir")
	if !filepath.IsAbs(gitDir) {
		gitDir = filepath.Join(repo, gitDir)
	}
	if err := os.MkdirAll(filepath.Join(gitDir, "rebase-merge"), 0o755); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"close-loop", "--repo", repo, "--dry-run", "--json"}, &stdout, &stderr)
	if code != ExitUnsafeRepoState {
		t.Fatalf("expected unsafe repo state, got %d; stderr=%s", code, stderr.String())
	}
}

func TestCloseLoopMissingGitBlocks(t *testing.T) {
	t.Setenv("PATH", emptyPathDir(t))

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"close-loop", "--repo", t.TempDir(), "--dry-run", "--json"}, &stdout, &stderr)
	if code != ExitMissingDependency {
		t.Fatalf("expected missing dependency, got %d; stderr=%s", code, stderr.String())
	}

	var result Result
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}
	if result.Error == nil || result.Error.Code != "missing_git" {
		t.Fatalf("expected missing_git, got %#v", result.Error)
	}
}

func TestCloseLoopCommitMissingGitBlocksBeforeWriteMode(t *testing.T) {
	t.Setenv("PATH", emptyPathDir(t))

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"close-loop", "--repo", t.TempDir(), "--commit", "--json"}, &stdout, &stderr)
	if code != ExitMissingDependency {
		t.Fatalf("expected missing dependency, got %d; stderr=%s", code, stderr.String())
	}

	var result Result
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}
	if result.Error == nil || result.Error.Code != "missing_git" {
		t.Fatalf("expected missing_git before write-mode block, got %#v", result.Error)
	}
}

func TestCloseLoopCommitModeBlockedUntilParity(t *testing.T) {
	repo := initTempGitRepo(t)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"close-loop", "--repo", repo, "--commit", "--json"}, &stdout, &stderr)
	if code != ExitPartialCloseBlocked {
		t.Fatalf("expected write mode blocked, got %d; stderr=%s", code, stderr.String())
	}
}

func initTempGitRepo(t *testing.T) string {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is required for temp repo tests")
	}
	repo := t.TempDir()
	runGit(t, repo, "init")
	runGit(t, repo, "config", "user.email", "test@example.invalid")
	runGit(t, repo, "config", "user.name", "Test User")
	writeFile(t, filepath.Join(repo, "README.md"), "# Fixture\n")
	runGit(t, repo, "add", "README.md")
	runGit(t, repo, "commit", "-m", "initial")
	return repo
}

func runGit(t *testing.T, repo string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", repo}, args...)...)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, string(output))
	}
}

func gitOutput(t *testing.T, repo string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", repo}, args...)...)
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("git %v failed: %v", args, err)
	}
	return filepath.Clean(string(bytes.TrimSpace(output)))
}

func writeFile(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func hasCheckMessage(checks []Check, name string, message string) bool {
	for _, check := range checks {
		if check.Name == name && check.Message == message {
			return true
		}
	}
	return false
}

func TestCloseLoopPathGroupingMatchesWindowsSafePrefixes(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("prefix grouping uses Git porcelain slash paths, covered by other tests")
	}
	if closeLoopGroupForPath("scripts/tool.sh") != "scripts" {
		t.Fatal("expected scripts group")
	}
}
