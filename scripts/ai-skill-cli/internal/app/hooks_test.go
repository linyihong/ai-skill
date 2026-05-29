package app

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
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

func TestFormatDirtyGitRepoReportCombinesNestedRepos(t *testing.T) {
	workspace := t.TempDir()
	repoA := filepath.Join(workspace, "TATA")
	repoB := filepath.Join(workspace, "mr.HS")
	os.MkdirAll(repoA, 0o755)
	os.MkdirAll(repoB, 0o755)
	runGit(t, repoA, "init")
	runGit(t, repoA, "config", "user.email", "test@example.invalid")
	runGit(t, repoA, "config", "user.name", "Test User")
	writeFile(t, filepath.Join(repoA, "README.md"), "# A\n")
	runGit(t, repoA, "add", "README.md")
	runGit(t, repoA, "commit", "-m", "initial")
	runGit(t, repoB, "init")
	runGit(t, repoB, "config", "user.email", "test@example.invalid")
	runGit(t, repoB, "config", "user.name", "Test User")
	writeFile(t, filepath.Join(repoB, "README.md"), "# B\n")
	runGit(t, repoB, "add", "README.md")
	runGit(t, repoB, "commit", "-m", "initial")

	writeFile(t, filepath.Join(repoA, "dirty.txt"), "dirty\n")
	writeFile(t, filepath.Join(repoB, "dirty.txt"), "dirty\n")

	report := formatDirtyGitRepoReport(workspace)
	if !strings.Contains(report, "TATA") || !strings.Contains(report, "mr.HS") {
		t.Fatalf("expected combined nested repo report, got:\n%s", report)
	}
	if !strings.Contains(report, "### Project Git Report") {
		t.Fatalf("expected final response report instruction, got:\n%s", report)
	}
}

func TestRunUserPromptSubmitHookIncludesNestedGitReport(t *testing.T) {
	workspace := t.TempDir()
	writeFile(t, filepath.Join(workspace, "CORE_BOOTSTRAP.md"), "# Bootstrap\n")
	repo := filepath.Join(workspace, "nested")
	os.MkdirAll(repo, 0o755)
	runGit(t, repo, "init")
	runGit(t, repo, "config", "user.email", "test@example.invalid")
	runGit(t, repo, "config", "user.name", "Test User")
	writeFile(t, filepath.Join(repo, "README.md"), "# Nested\n")
	runGit(t, repo, "add", "README.md")
	runGit(t, repo, "commit", "-m", "initial")
	writeFile(t, filepath.Join(repo, "dirty.txt"), "dirty\n")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if code := runUserPromptSubmitHook(workspace, &stdout, &stderr); code != ExitSuccess {
		t.Fatalf("expected success, got %d; stderr=%s", code, stderr.String())
	}
	var output map[string]map[string]string
	if err := json.Unmarshal(stdout.Bytes(), &output); err != nil {
		t.Fatalf("decode hook output: %v\n%s", err, stdout.String())
	}
	ctx := output["hookSpecificOutput"]["additionalContext"]
	if !strings.Contains(ctx, "nested") || !strings.Contains(ctx, "### Project Git Report") {
		t.Fatalf("expected nested git report in context, got:\n%s", ctx)
	}
}

func TestRunStopHookBlocksCursorPayloadWithoutCognitive(t *testing.T) {
	setHookStdin(t, `{"assistant_response":"Done. Tests passed."}`)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runStopHook(t.TempDir(), &stdout, &stderr)
	if code != ExitValidationFailed {
		t.Fatalf("expected missing Cognitive block to fail, got %d; stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stderr.String(), "Missing obligation") {
		t.Fatalf("expected missing obligation message, got %s", stderr.String())
	}
}

func TestRunStopHookBlocksCursorOkOnlyPayload(t *testing.T) {
	setHookStdin(t, `{"assistant_response":"OK"}`)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runStopHook(t.TempDir(), &stdout, &stderr)
	if code != ExitValidationFailed {
		t.Fatalf("expected OK-only response to fail, got %d; stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stderr.String(), "Missing obligation") {
		t.Fatalf("expected missing obligation message, got %s", stderr.String())
	}
}

func TestRunStopHookBlocksMissingAssistantText(t *testing.T) {
	setHookStdin(t, `{"hook_event_name":"afterAgentResponse"}`)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runStopHook(t.TempDir(), &stdout, &stderr)
	if code != ExitValidationFailed {
		t.Fatalf("expected missing assistant text to fail closed, got %d; stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stderr.String(), "BLOCK_NO_ASSISTANT_TEXT") {
		t.Fatalf("expected missing assistant text diagnostic, got %s", stderr.String())
	}
}

func TestRunStopHookAllowsCursorPayloadWithCompactCognitive(t *testing.T) {
	setHookStdin(t, `{"assistant_response":"Done.\n\nCognitive: NORMAL·SUMMARY_FIRST·STANDARD·NONE / V:CHECKLIST / Cost:LOW / Sig:user_keyword_fast"}`)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runStopHook(t.TempDir(), &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected Cognitive block to pass, got %d; stderr=%s", code, stderr.String())
	}
}

func setHookStdin(t *testing.T, input string) {
	t.Helper()
	previous := os.Stdin
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := writer.WriteString(input); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	os.Stdin = reader
	t.Cleanup(func() {
		os.Stdin = previous
		_ = reader.Close()
	})
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

func TestValidateTokenBudget(t *testing.T) {
	// No Token Estimate trailer → no-op
	v := validateTokenBudget(map[string]string{"execution_mode": "DEEP"}, "feat: x")
	if v != "" {
		t.Fatalf("expected no-op without estimate, got %q", v)
	}
	// Within budget
	modes := map[string]string{"execution_mode": "NORMAL", "context_mode": "SUMMARY_FIRST", "governance_mode": "STANDARD", "memory_mode": "EPISODIC"}
	v = validateTokenBudget(modes, "feat: x\n\nToken Estimate: 3000\n")
	if v != "" {
		t.Fatalf("expected no violation within budget, got %q", v)
	}
	// Exceeds tuple budget
	v = validateTokenBudget(modes, "feat: x\n\nToken Estimate: 9999\n")
	if v == "" {
		t.Fatal("expected violation when over tuple budget")
	}
	// Exceeds execution_mode default (no exact tuple match)
	modes2 := map[string]string{"execution_mode": "FAST", "context_mode": "SUMMARY_FIRST", "governance_mode": "LIGHT", "memory_mode": "NONE"}
	v = validateTokenBudget(modes2, "feat: x\n\nToken Estimate: 1500\n")
	if v == "" {
		t.Fatal("expected violation when over FAST default budget")
	}
	// Opt-out trailer skips
	v = validateTokenBudget(modes, "feat: x\n\nToken Estimate: 999999\n\n[skip-token-budget]\n")
	if v != "" {
		t.Fatalf("expected no violation with opt-out, got %q", v)
	}
	// Unknown execution_mode → no enforcement
	v = validateTokenBudget(map[string]string{"execution_mode": "WEIRD"}, "feat: x\n\nToken Estimate: 99999\n")
	if v != "" {
		t.Fatalf("expected no enforcement for unknown mode, got %q", v)
	}
}

func TestDeriveCognitiveCost(t *testing.T) {
	cases := []struct {
		exec string
		ctx  string
		want string
	}{
		// cost lookup: FAST × INDEX_ONLY -> LOW
		{"FAST", "INDEX_ONLY", "LOW"},
		// cost lookup: FAST × other context -> MEDIUM
		{"FAST", "SUMMARY_FIRST", "MEDIUM"},
		// cost lookup: NORMAL × INDEX_ONLY -> LOW
		{"NORMAL", "INDEX_ONLY", "LOW"},
		// cost lookup: NORMAL × SUMMARY_FIRST -> LOW
		{"NORMAL", "SUMMARY_FIRST", "LOW"},
		// cost lookup: NORMAL × CHECKLIST_FIRST -> MEDIUM
		{"NORMAL", "CHECKLIST_FIRST", "MEDIUM"},
		// cost lookup: NORMAL × SOURCE_BACKED -> MEDIUM
		{"NORMAL", "SOURCE_BACKED", "MEDIUM"},
		// cost lookup: DEEP × any context -> HIGH
		{"DEEP", "SOURCE_BACKED", "HIGH"},
		// cost lookup: FORENSIC × any context -> VERY_HIGH
		{"FORENSIC", "GRAPH_ASSISTED", "VERY_HIGH"},
		// cost lookup: RECOVERY × any context -> VERY_HIGH
		{"RECOVERY", "CHECKLIST_FIRST", "VERY_HIGH"},
	}
	for _, tt := range cases {
		if got := deriveCognitiveCost(tt.exec, tt.ctx); got != tt.want {
			t.Fatalf("deriveCognitiveCost(%s, %s) = %s, want %s", tt.exec, tt.ctx, got, tt.want)
		}
	}
}

func TestValidateCognitiveCost(t *testing.T) {
	valid := map[string]string{"execution_mode": "NORMAL", "context_mode": "SUMMARY_FIRST", "cognitive_cost": "LOW"}
	if v := validateCognitiveCost(valid); v != "" {
		t.Fatalf("expected valid cost, got %q", v)
	}
	mismatch := map[string]string{"execution_mode": "DEEP", "context_mode": "SOURCE_BACKED", "cognitive_cost": "LOW"}
	if v := validateCognitiveCost(mismatch); v == "" {
		t.Fatal("expected cost mismatch BLOCK for DEEP + LOW")
	}
	missing := map[string]string{"execution_mode": "NORMAL", "context_mode": "SUMMARY_FIRST"}
	if v := validateCognitiveCost(missing); v == "" {
		t.Fatal("expected missing cognitive_cost violation")
	}
}

func TestValidateActivationSignals(t *testing.T) {
	root := t.TempDir()
	os.MkdirAll(filepath.Join(root, "runtime"), 0o755)
	discovery := "signals:\n  - name: file_diff_runtime_schema\n  - name: user_keyword_deep\n"
	os.WriteFile(filepath.Join(root, "runtime", "cognitive-modes-discovery.yaml"), []byte(discovery), 0o644)

	validText := "feat: x\n\n### Cognitive Mode 報告\n\nactivation_reason:\n  - file_diff_runtime_schema\n"
	if v := validateActivationSignals(commitMsgCtx{text: validText, root: root, modes: map[string]string{}}); v != "" {
		t.Fatalf("expected known activation signal to pass, got %q", v)
	}
	unknownText := "feat: x\n\n### Cognitive Mode 報告\n\nactivation_reason:\n  - made_up_signal\n"
	if v := validateActivationSignals(commitMsgCtx{text: unknownText, root: root, modes: map[string]string{}}); v == "" {
		t.Fatal("expected unknown signal to BLOCK")
	}
	emptyText := "feat: x\n\n### Cognitive Mode 報告\n\n| execution_mode | NORMAL | reason |\n"
	if v := validateActivationSignals(commitMsgCtx{text: emptyText, root: root, modes: map[string]string{}}); v == "" {
		t.Fatal("expected empty activation_reason to BLOCK")
	}
	compact := map[string]string{"activation_signal": "user_keyword_deep"}
	if v := validateActivationSignals(commitMsgCtx{text: "feat: x", root: root, modes: compact}); v != "" {
		t.Fatalf("expected compact Sig to pass, got %q", v)
	}
}

func TestValidateCapabilitySnippet(t *testing.T) {
	highRisk := map[string]string{"execution_mode": "DEEP", "governance_mode": "STRICT"}
	if v := validateCapabilitySnippet(highRisk, "feat: x\n\n### Cognitive Mode 報告\n"); v == "" {
		t.Fatal("expected high-risk mode without Capability summary to BLOCK")
	}
	withSnippet := "feat: x\n\n### Cognitive Mode 報告\n\nCapability summary:\n  execution_mode=DEEP -> source-backed reads\n"
	if v := validateCapabilitySnippet(highRisk, withSnippet); v != "" {
		t.Fatalf("expected high-risk mode with Capability summary to pass, got %q", v)
	}
	lowRisk := map[string]string{"execution_mode": "NORMAL", "governance_mode": "STANDARD"}
	if v := validateCapabilitySnippet(lowRisk, "feat: x\n\n### Cognitive Mode 報告\n"); v != "" {
		t.Fatalf("expected low-risk mode without snippet to pass, got %q", v)
	}
}

func TestInflatedRejection(t *testing.T) {
	costMismatch := map[string]string{"execution_mode": "DEEP", "context_mode": "SUMMARY_FIRST", "cognitive_cost": "LOW"}
	if v := validateCognitiveCost(costMismatch); v == "" {
		t.Fatal("expected typo DEEP + LOW mismatch BLOCK")
	}
	root := t.TempDir()
	os.MkdirAll(filepath.Join(root, "runtime"), 0o755)
	os.WriteFile(filepath.Join(root, "runtime", "cognitive-modes-discovery.yaml"), []byte("signals:\n  - name: file_diff_notes_ephemeral\n"), 0o644)
	unknown := "feat: typo\n\n### Cognitive Mode 報告\n\nactivation_reason:\n  - made_up_signal\n"
	if v := validateActivationSignals(commitMsgCtx{text: unknown, root: root, modes: map[string]string{}}); v == "" {
		t.Fatal("expected unknown signal BLOCK")
	}
	valid := "feat: docs\n\n### Cognitive Mode 報告\n\nactivation_reason:\n  - file_diff_notes_ephemeral\n"
	if v := validateActivationSignals(commitMsgCtx{text: valid, root: root, modes: map[string]string{}}); v != "" {
		t.Fatalf("expected legitimate known signal PASS, got %q", v)
	}
}

func TestValidateAdaptiveTriggers(t *testing.T) {
	// Case 1: contradiction_risk fires when keyword + ≥2 distinct sources
	modes := map[string]string{"execution_mode": "NORMAL", "context_mode": "SUMMARY_FIRST", "governance_mode": "STANDARD", "memory_mode": "EPISODIC"}
	body := "feat: reconcile contradict plans/active/a.md vs constitution/ADR-001.md"
	v := validateAdaptiveTriggers(modes, body)
	if v == "" {
		t.Fatal("expected contradiction_risk violation")
	}
	// Same body but with STRICT + SOURCE_BACKED → ok
	modes2 := map[string]string{"execution_mode": "DEEP", "context_mode": "SOURCE_BACKED", "governance_mode": "STRICT", "memory_mode": "DECISION_REPLAY"}
	v = validateAdaptiveTriggers(modes2, body)
	if v != "" {
		t.Fatalf("expected no violation with STRICT+SOURCE_BACKED, got %q", v)
	}

	// Case 2: repeated_failure — 2 failure-pattern refs
	body2 := "fix: address enforcement/failure-patterns/foo.md and enforcement/failure-patterns/bar.md"
	v = validateAdaptiveTriggers(modes, body2)
	if v == "" {
		t.Fatal("expected repeated_failure violation")
	}
	// Same body with RECOVERY + FAILURE_REPLAY → ok
	modes3 := map[string]string{"execution_mode": "RECOVERY", "context_mode": "CHECKLIST_FIRST", "governance_mode": "STRICT", "memory_mode": "FAILURE_REPLAY"}
	v = validateAdaptiveTriggers(modes3, body2)
	if v != "" {
		t.Fatalf("expected no violation with RECOVERY+FAILURE_REPLAY, got %q", v)
	}

	// Case 3: budget_near_ceiling — Token Estimate at 90% of NORMAL budget (5000)
	body3 := "feat: medium work\n\nToken Estimate: 4500\n"
	v = validateAdaptiveTriggers(modes, body3)
	if v == "" {
		t.Fatal("expected near-ceiling warning")
	}
	if !contains(v, "≥80%") {
		t.Fatalf("expected near-ceiling warning text, got %q", v)
	}

	// Opt-out marker skips all
	v = validateAdaptiveTriggers(modes, body+"\n\n[skip-adaptive]\n")
	if v != "" {
		t.Fatalf("expected no violation with opt-out, got %q", v)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || (len(s) > 0 && (func() bool {
		for i := 0; i+len(sub) <= len(s); i++ {
			if s[i:i+len(sub)] == sub {
				return true
			}
		}
		return false
	}())))
}

func TestValidateBootstrapEntryThinness(t *testing.T) {
	// Build a temp repo with a thin CLAUDE.md → no violation
	tmp := t.TempDir()
	thin := "# entry\n\nRead CORE_BOOTSTRAP.md. That's the canonical source.\n"
	if err := os.WriteFile(filepath.Join(tmp, "CLAUDE.md"), []byte(thin), 0o644); err != nil {
		t.Fatal(err)
	}
	v := validateBootstrapEntryThinness("feat: x", []string{"CLAUDE.md"}, tmp)
	if v != "" {
		t.Fatalf("expected no violation for thin file, got %q", v)
	}

	// Add forbidden enum content → violation
	bloated := thin + "\nMode values: FAST/NORMAL/DEEP/FORENSIC/RECOVERY\n"
	if err := os.WriteFile(filepath.Join(tmp, "CLAUDE.md"), []byte(bloated), 0o644); err != nil {
		t.Fatal(err)
	}
	v = validateBootstrapEntryThinness("feat: x", []string{"CLAUDE.md"}, tmp)
	if v == "" {
		t.Fatal("expected violation for enum content")
	}

	// Over line count → violation
	manyLines := ""
	for i := 0; i < 40; i++ {
		manyLines += "line\n"
	}
	if err := os.WriteFile(filepath.Join(tmp, "CLAUDE.md"), []byte(manyLines), 0o644); err != nil {
		t.Fatal(err)
	}
	v = validateBootstrapEntryThinness("feat: x", []string{"CLAUDE.md"}, tmp)
	if v == "" {
		t.Fatal("expected violation for >30 lines")
	}

	// Opt-out skips even bloated file
	v = validateBootstrapEntryThinness("feat: x\n\n[skip-bootstrap-thinness]\n", []string{"CLAUDE.md"}, tmp)
	if v != "" {
		t.Fatalf("expected no violation with opt-out, got %q", v)
	}

	// Non-entry file staged → no enforcement
	v = validateBootstrapEntryThinness("feat: x", []string{"README.md"}, tmp)
	if v != "" {
		t.Fatalf("expected no enforcement on non-entry file, got %q", v)
	}
}

func TestValidateRuntimeYamlProjects(t *testing.T) {
	tmp := t.TempDir()
	os.MkdirAll(filepath.Join(tmp, "runtime"), 0o755)
	goodYAML := "runtime_projection:\n  enabled: true\n  target_key: runtime.test.contract\n"
	badYAML := "runtime_projection:\n  enabled: false\n"
	os.WriteFile(filepath.Join(tmp, "runtime", "good.yaml"), []byte(goodYAML), 0o644)
	os.WriteFile(filepath.Join(tmp, "runtime", "bad.yaml"), []byte(badYAML), 0o644)

	// Good only → pass
	v := validateRuntimeYamlProjects("feat: x", []string{"runtime/good.yaml"}, tmp)
	if v != "" {
		t.Fatalf("expected pass for projected yaml, got %q", v)
	}
	// Bad → block
	v = validateRuntimeYamlProjects("feat: x", []string{"runtime/bad.yaml"}, tmp)
	if v == "" {
		t.Fatal("expected violation for non-projected yaml")
	}
	// Opt-out
	v = validateRuntimeYamlProjects("feat: x\n\n[skip-runtime-yaml-projection]\n", []string{"runtime/bad.yaml"}, tmp)
	if v != "" {
		t.Fatalf("expected opt-out to bypass, got %q", v)
	}
}

func TestValidateMarkdownYamlSync(t *testing.T) {
	tmp := t.TempDir()
	os.MkdirAll(filepath.Join(tmp, "governance"), 0o755)
	// Create a paired md+yaml; staging md alone should violate
	os.WriteFile(filepath.Join(tmp, "governance", "foo.md"), []byte("# foo"), 0o644)
	os.WriteFile(filepath.Join(tmp, "governance", "foo.yaml"), []byte("id: foo"), 0o644)
	// Create md without yaml sibling; should not violate
	os.WriteFile(filepath.Join(tmp, "governance", "orphan.md"), []byte("# orphan"), 0o644)

	// md alone (with sibling existing) → block
	v := validateMarkdownYamlSync("feat: x", []string{"governance/foo.md"}, tmp)
	if v == "" {
		t.Fatal("expected violation when sibling yaml exists but not staged")
	}
	// both staged → pass
	v = validateMarkdownYamlSync("feat: x", []string{"governance/foo.md", "governance/foo.yaml"}, tmp)
	if v != "" {
		t.Fatalf("expected pass when both staged, got %q", v)
	}
	// md without sibling → pass
	v = validateMarkdownYamlSync("feat: x", []string{"governance/orphan.md"}, tmp)
	if v != "" {
		t.Fatalf("expected pass for orphan md, got %q", v)
	}
	// opt-out
	v = validateMarkdownYamlSync("feat: x\n\n[skip-markdown-yaml-sync]\n", []string{"governance/foo.md"}, tmp)
	if v != "" {
		t.Fatalf("expected opt-out to bypass, got %q", v)
	}
}

func TestValidateGlossaryRetroOwn(t *testing.T) {
	tmp := t.TempDir()

	// Happy path: framework surface + glossary both staged
	v := validateGlossaryRetroOwn(
		"feat: add new cognitive signal",
		[]string{"runtime/cognitive-modes-discovery.yaml", "knowledge/glossary/ai-skill.md"},
		tmp,
	)
	if v != "" {
		t.Fatalf("expected pass when glossary also staged, got %q", v)
	}

	// Block: cognitive-modes staged but no glossary
	v = validateGlossaryRetroOwn(
		"feat: add new cognitive signal",
		[]string{"runtime/cognitive-modes-discovery.yaml"},
		tmp,
	)
	if v == "" {
		t.Fatal("expected block when cognitive-modes staged without glossary")
	}

	// Block: runtime/economics surface staged but no glossary
	v = validateGlossaryRetroOwn(
		"feat: add economics contract",
		[]string{"runtime/economics/token-costs.yaml"},
		tmp,
	)
	if v == "" {
		t.Fatal("expected block for runtime/economics/ without glossary")
	}

	// Block: ecosystem/ surface staged but no glossary
	v = validateGlossaryRetroOwn(
		"feat: add ecosystem adaptation",
		[]string{"ecosystem/cognition/adaptation.yaml"},
		tmp,
	)
	if v == "" {
		t.Fatal("expected block for ecosystem/ without glossary")
	}

	// Opt-out: skip marker bypasses block
	v = validateGlossaryRetroOwn(
		"chore: typo fix\n\n[skip-glossary-retro-own]\n",
		[]string{"runtime/cognitive-modes-discovery.yaml"},
		tmp,
	)
	if v != "" {
		t.Fatalf("expected opt-out to bypass, got %q", v)
	}

	// No framework surface staged → no enforcement
	v = validateGlossaryRetroOwn(
		"feat: unrelated doc change",
		[]string{"README.md"},
		tmp,
	)
	if v != "" {
		t.Fatalf("expected no enforcement on non-framework staging, got %q", v)
	}

	// Other runtime/*.yaml (not cognitive-modes) should not trigger
	v = validateGlossaryRetroOwn(
		"feat: tweak unrelated runtime yaml",
		[]string{"runtime/cli-modification-policy.yaml"},
		tmp,
	)
	if v != "" {
		t.Fatalf("expected no enforcement on non-cognitive-modes runtime yaml, got %q", v)
	}
}
