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
	if len(result.PlannedActions) != 4 {
		t.Fatalf("expected four planned hook installs (pre-commit/commit-msg/post-commit/pre-push), got %#v", result.PlannedActions)
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

func TestParseCognitiveModeBlock(t *testing.T) {
	text := "feat: example\n\n### Cognitive Mode 報告\n\n| 維度 | 值 | 理由 |\n|------|------|------|\n| execution_mode | DEEP | rationale |\n| context_mode | SOURCE_BACKED | rationale |\n| governance_mode | STRICT | rationale |\n| memory_mode | EPISODIC | rationale |\n"
	modes := parseCognitiveModeBlock(text)
	if modes["execution_mode"] != "DEEP" || modes["context_mode"] != "SOURCE_BACKED" || modes["governance_mode"] != "STRICT" || modes["memory_mode"] != "EPISODIC" {
		t.Fatalf("parse mismatch: %#v", modes)
	}
}

func TestValidateExecutionModeFloors(t *testing.T) {
	// FAST forbidden when touching enforcement/
	v := validateExecutionModeFloors(map[string]string{"execution_mode": "FAST"}, []string{"enforcement/foo.md"})
	if v == "" {
		t.Fatal("expected FAST violation when touching enforcement/")
	}
	// DEEP without STRICT governance
	v = validateExecutionModeFloors(map[string]string{"execution_mode": "DEEP", "governance_mode": "STANDARD", "context_mode": "SOURCE_BACKED"}, nil)
	if v == "" {
		t.Fatal("expected DEEP violation without STRICT governance")
	}
	// DEEP + STRICT + SOURCE_BACKED → OK
	v = validateExecutionModeFloors(map[string]string{"execution_mode": "DEEP", "governance_mode": "STRICT", "context_mode": "SOURCE_BACKED"}, nil)
	if v != "" {
		t.Fatalf("expected no violation, got %q", v)
	}
	// RECOVERY requires FAILURE_REPLAY memory
	v = validateExecutionModeFloors(map[string]string{"execution_mode": "RECOVERY", "governance_mode": "STRICT", "context_mode": "CHECKLIST_FIRST", "memory_mode": "EPISODIC"}, nil)
	if v == "" {
		t.Fatal("expected RECOVERY violation without FAILURE_REPLAY memory")
	}
}

func TestValidateGovernanceModeConsistency(t *testing.T) {
	// LIGHT touching enforcement/
	v := validateGovernanceModeConsistency(map[string]string{"governance_mode": "LIGHT"}, []string{"enforcement/x.md"}, "feat: x")
	if v == "" {
		t.Fatal("expected LIGHT violation")
	}
	// LOCKDOWN without approval
	v = validateGovernanceModeConsistency(map[string]string{"governance_mode": "LOCKDOWN"}, nil, "feat: critical")
	if v == "" {
		t.Fatal("expected LOCKDOWN violation without approval")
	}
	// LOCKDOWN with approval trailer
	v = validateGovernanceModeConsistency(map[string]string{"governance_mode": "LOCKDOWN"}, nil, "feat: critical\n\n[approved-by: alice]\n")
	if v != "" {
		t.Fatalf("expected no violation with approval, got %q", v)
	}
}

func TestValidateMemoryModeSubdir(t *testing.T) {
	// NONE but touching memory/episodic/
	v := validateMemoryModeSubdir(map[string]string{"memory_mode": "NONE"}, []string{"memory/episodic/foo.md"})
	if v == "" {
		t.Fatal("expected NONE violation when touching memory/episodic/")
	}
	// EPISODIC touching memory/decision/ — wrong subdir
	v = validateMemoryModeSubdir(map[string]string{"memory_mode": "EPISODIC"}, []string{"memory/decision/foo.md"})
	if v == "" {
		t.Fatal("expected EPISODIC violation when touching memory/decision/")
	}
	// EPISODIC touching memory/episodic/ — OK
	v = validateMemoryModeSubdir(map[string]string{"memory_mode": "EPISODIC"}, []string{"memory/episodic/foo.md"})
	if v != "" {
		t.Fatalf("expected no violation, got %q", v)
	}
	// Layer doc exempt
	v = validateMemoryModeSubdir(map[string]string{"memory_mode": "NONE"}, []string{"memory/README.md", "memory/retrieval-governance/activation-thresholds.md"})
	if v != "" {
		t.Fatalf("expected no violation for layer docs, got %q", v)
	}
}

func TestValidatePlanStatusSync(t *testing.T) {
	// Trigger fires: completion + Phase + plan ref, but plan not staged → block
	body := "feat: Phase 3 完成\n\nSee plans/active/2026-05-22-1629-runtime-cognitive-modes-system.md"
	v := validatePlanStatusSync(body, []string{"scripts/ai-skill-cli/internal/app/hooks.go"})
	if v == "" {
		t.Fatal("expected violation when plan completion claimed but plan not staged")
	}
	// Same body but plan IS staged → ok
	v = validatePlanStatusSync(body, []string{"plans/active/2026-05-22-1629-runtime-cognitive-modes-system.md"})
	if v != "" {
		t.Fatalf("expected no violation when plan is staged, got %q", v)
	}
	// No completion vocabulary → no trigger
	v = validatePlanStatusSync("docs: see plans/active/foo.md for context\n\nPhase 3 context", nil)
	if v != "" {
		t.Fatalf("expected no violation without completion vocabulary, got %q", v)
	}
	// No phase mention → no trigger
	v = validatePlanStatusSync("feat: completed plans/active/foo.md feature\n", nil)
	if v != "" {
		t.Fatalf("expected no violation without Phase N mention, got %q", v)
	}
	// Opt-out trailer skips
	body2 := "post-mortem: Phase 3 完成 looking back at plans/active/foo.md\n\n[skip-plan-status-sync]\n"
	v = validatePlanStatusSync(body2, nil)
	if v != "" {
		t.Fatalf("expected no violation with opt-out marker, got %q", v)
	}
}
