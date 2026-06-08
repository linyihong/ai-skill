package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
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
	repoA := filepath.Join(workspace, "repo-alpha")
	repoB := filepath.Join(workspace, "repo-beta")
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
	if !strings.Contains(report, "repo-alpha") || !strings.Contains(report, "repo-beta") {
		t.Fatalf("expected combined nested repo report, got:\n%s", report)
	}
	if !strings.Contains(report, "### Project Git Report") {
		t.Fatalf("expected final response report instruction, got:\n%s", report)
	}
}

func TestFormatDirtyGitRepoReportSkipsAiSkillRepoRoot(t *testing.T) {
	workspace := t.TempDir()
	writeFile(t, filepath.Join(workspace, "CORE_BOOTSTRAP.md"), "# Bootstrap\n")
	writeFile(t, filepath.Join(workspace, "runtime", "core-bootstrap.yaml"), "schema_version: 1\n")
	if err := os.MkdirAll(filepath.Join(workspace, "scripts", "ai-skill-cli"), 0o755); err != nil {
		t.Fatalf("mkdir ai-skill cli dir: %v", err)
	}
	runGit(t, workspace, "init")
	runGit(t, workspace, "config", "user.email", "test@example.invalid")
	runGit(t, workspace, "config", "user.name", "Test User")
	writeFile(t, filepath.Join(workspace, "README.md"), "# Ai-skill\n")
	runGit(t, workspace, "add", ".")
	runGit(t, workspace, "commit", "-m", "initial")
	writeFile(t, filepath.Join(workspace, "dirty.txt"), "dirty\n")

	if report := formatDirtyGitRepoReport(workspace); report != "" {
		t.Fatalf("expected Ai-skill repo root to skip Project Git Report, got:\n%s", report)
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

	// ADR-011: hook now drains stdin to look up transcript_path. Empty
	// payload → transcript not resolvable → conservative path (inject full
	// bootstrap). Pipe an empty JSON object so the read does not block.
	setHookStdin(t, `{}`)

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

// TestRunUserPromptSubmitHookFirstTurnInjectsFullBootstrap is ADR-011 case A:
// transcript has no Bootstrap Receipt yet → hook must inject the full
// CORE_BOOTSTRAP.md + the "Receipt not yet observed" prompt alongside the
// MUST close-out block.
func TestRunUserPromptSubmitHookFirstTurnInjectsFullBootstrap(t *testing.T) {
	workspace := t.TempDir()
	// Force AI_SKILL_REPO so resolveClaudeAiSkillRepo doesn't escape to the
	// real repo's CORE_BOOTSTRAP.md and let us assert on a controllable marker.
	bootstrapMarker := "# Bootstrap stub for ADR-011 first-turn test\n"
	writeFile(t, filepath.Join(workspace, "CORE_BOOTSTRAP.md"), bootstrapMarker)
	writeFile(t, filepath.Join(workspace, "runtime", "core-bootstrap.yaml"), "schema_version: 1\n")
	t.Setenv("AI_SKILL_REPO", workspace)

	// Transcript with only a non-Receipt assistant turn — acknowledgment scan
	// must return false.
	tr := writeBootstrapTranscript(t, workspace,
		"Hello, working on the task. No Receipt yet.", nil)
	setHookStdin(t, fmt.Sprintf(`{"transcript_path":%q}`, tr))

	var stdout, stderr bytes.Buffer
	if code := runUserPromptSubmitHook(workspace, &stdout, &stderr); code != ExitSuccess {
		t.Fatalf("expected success, got %d; stderr=%s", code, stderr.String())
	}
	var output map[string]map[string]string
	if err := json.Unmarshal(stdout.Bytes(), &output); err != nil {
		t.Fatalf("decode hook output: %v\n%s", err, stdout.String())
	}
	ctx := output["hookSpecificOutput"]["additionalContext"]
	if !strings.Contains(ctx, "final close-out obligation") {
		t.Errorf("MUST block missing from first-turn context:\n%s", ctx)
	}
	if !strings.Contains(ctx, "Bootstrap Receipt not yet observed") {
		t.Errorf("conditional bootstrap prompt missing from first-turn context:\n%s", ctx)
	}
	if !strings.Contains(ctx, bootstrapMarker) {
		t.Errorf("full CORE_BOOTSTRAP.md not injected on first turn; got:\n%s", ctx)
	}
}

// TestRunUserPromptSubmitHookSubsequentTurnSkipsBootstrap is ADR-011 case B:
// transcript already contains a Bootstrap Receipt line in a prior assistant
// text turn → hook must inject ONLY the MUST close-out block; CORE_BOOTSTRAP.md
// must NOT be re-injected (this is the ~2-3K token saving).
func TestRunUserPromptSubmitHookSubsequentTurnSkipsBootstrap(t *testing.T) {
	workspace := t.TempDir()
	bootstrapMarker := "# Bootstrap stub for ADR-011 subsequent-turn test\n"
	writeFile(t, filepath.Join(workspace, "CORE_BOOTSTRAP.md"), bootstrapMarker)

	receiptLine := "Bootstrap: rules=✓ phase=phase.bootstrap obligations=23 gates=25\n" +
		"Active per-turn obligations: obligation.cognitive.mode_report, obligation.finality.close_loop_check\n" +
		"\nDoing work..."
	tr := writeBootstrapTranscript(t, workspace, receiptLine, nil)
	setHookStdin(t, fmt.Sprintf(`{"transcript_path":%q}`, tr))

	var stdout, stderr bytes.Buffer
	if code := runUserPromptSubmitHook(workspace, &stdout, &stderr); code != ExitSuccess {
		t.Fatalf("expected success, got %d; stderr=%s", code, stderr.String())
	}
	var output map[string]map[string]string
	if err := json.Unmarshal(stdout.Bytes(), &output); err != nil {
		t.Fatalf("decode hook output: %v\n%s", err, stdout.String())
	}
	ctx := output["hookSpecificOutput"]["additionalContext"]
	if !strings.Contains(ctx, "final close-out obligation") {
		t.Errorf("MUST block missing from subsequent-turn context:\n%s", ctx)
	}
	if strings.Contains(ctx, "Bootstrap Receipt not yet observed") {
		t.Errorf("conditional prompt should be omitted when Receipt already in transcript; got:\n%s", ctx)
	}
	if strings.Contains(ctx, bootstrapMarker) {
		t.Errorf("CORE_BOOTSTRAP.md must NOT be re-injected after acknowledgment; got:\n%s", ctx)
	}
}

// No hook_event_name => Claude path (cursorStop=false). A Claude Stop block MUST
// be exit 0 + top-level {"decision":"block","reason":...} (the decision:block JSON
// does the blocking, NOT exit code 30, which Claude treats as a non-blocking error
// and stops anyway). Assert the JSON decision field, not just the exit code.
func TestRunStopHookBlocksClaudePayloadWithoutCognitive(t *testing.T) {
	setHookStdin(t, `{"assistant_response":"Done. Tests passed."}`)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runStopHook(t.TempDir(), &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected Claude Stop block via decision JSON to exit 0, got %d; stderr=%s", code, stderr.String())
	}
	var output map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &output); err != nil {
		t.Fatalf("decode Claude stop decision: %v\n%s", err, stdout.String())
	}
	if output["decision"] != "block" {
		t.Fatalf("expected decision:block, got %#v", output)
	}
	reason, _ := output["reason"].(string)
	if !strings.Contains(reason, "Missing obligation") {
		t.Fatalf("expected missing obligation reason, got %#v", output)
	}
}

func TestRunStopHookBlocksClaudeOkOnlyPayload(t *testing.T) {
	setHookStdin(t, `{"assistant_response":"OK"}`)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runStopHook(t.TempDir(), &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected Claude Stop block via decision JSON to exit 0, got %d; stderr=%s", code, stderr.String())
	}
	var output map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &output); err != nil {
		t.Fatalf("decode Claude stop decision: %v\n%s", err, stdout.String())
	}
	if output["decision"] != "block" {
		t.Fatalf("expected decision:block, got %#v", output)
	}
	reason, _ := output["reason"].(string)
	if !strings.Contains(reason, "Missing obligation") {
		t.Fatalf("expected missing obligation reason, got %#v", output)
	}
}

func TestRunStopHookLoopsCursorStopOkOnlyPayload(t *testing.T) {
	setHookStdin(t, `{"hook_event_name":"stop","assistant_response":"OK"}`)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runStopHook(t.TempDir(), &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected Cursor stop to loop with success exit, got %d; stderr=%s", code, stderr.String())
	}
	var output map[string]string
	if err := json.Unmarshal(stdout.Bytes(), &output); err != nil {
		t.Fatalf("decode Cursor stop output: %v\n%s", err, stdout.String())
	}
	if !strings.Contains(output["followup_message"], "Cognitive Mode block") {
		t.Fatalf("expected followup_message to request Cognitive block, got %#v", output)
	}
	if !strings.Contains(output["followup_message"], "Bootstrap Receipt") {
		t.Fatalf("expected followup_message to request Bootstrap Receipt, got %#v", output)
	}
}

func TestRunStopHookAllowsCursorUserAbortedStop(t *testing.T) {
	setHookStdin(t, `{"hook_event_name":"stop","status":"aborted","assistant_response":"partial progress without close-out"}`)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runStopHook(t.TempDir(), &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected user-aborted Cursor stop to pass, got %d; stderr=%s", code, stderr.String())
	}
	if stdout.String() != "" {
		t.Fatalf("expected no followup loop for user-aborted Cursor stop, got %s", stdout.String())
	}
	if !strings.Contains(stderr.String(), "ALLOW_CURSOR_USER_ABORT") {
		t.Fatalf("expected user-abort diagnostic, got %s", stderr.String())
	}
}

func TestRunStopHookAllowsAfterAgentResponseAuditOnly(t *testing.T) {
	setHookStdin(t, `{"hook_event_name":"afterAgentResponse"}`)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runStopHook(t.TempDir(), &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected afterAgentResponse audit-only hook to pass, got %d; stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stderr.String(), "ALLOW_AFTER_AGENT_RESPONSE_AUDIT_ONLY") {
		t.Fatalf("expected audit-only diagnostic, got %s", stderr.String())
	}
}

func TestRunStopHookAllowsCursorStopMissingAssistantText(t *testing.T) {
	setHookStdin(t, `{"hook_event_name":"stop"}`)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runStopHook(t.TempDir(), &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected Cursor stop missing assistant text to fail open, got %d; stderr=%s", code, stderr.String())
	}
	if stdout.String() != "" {
		t.Fatalf("expected no followup loop without assistant final text, got %s", stdout.String())
	}
	if !strings.Contains(stderr.String(), "ALLOW_CURSOR_NO_ASSISTANT_TEXT") {
		t.Fatalf("expected no-assistant-text fail-open diagnostic, got %s", stderr.String())
	}
}

func TestRunStopHookAllowsCursorPayloadWithCompactCognitive(t *testing.T) {
	setHookStdin(t, fmt.Sprintf(`{"assistant_response":%q}`, validStopCloseOutText()))

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runStopHook(t.TempDir(), &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected Cognitive block to pass, got %d; stderr=%s", code, stderr.String())
	}
}

func TestRunStopHookAllowsExplicitBootstrapAcknowledgementWithoutCanonicalReceipt(t *testing.T) {
	text := "已讀 CORE_BOOTSTRAP.md 與 runtime/core-bootstrap.yaml。Active per-turn obligations: obligation.cognitive.mode_report, obligation.feedback.learning_report, obligation.finality.close_loop_check\n\nDone.\n\n" +
		compactFeedbackReport("NONE", "LOCAL", "N/A", "") +
		"\nCognitive: NORMAL·SUMMARY_FIRST·STANDARD·NONE / V:CHECKLIST / Cost:LOW / Sig:repair"
	setHookStdin(t, fmt.Sprintf(`{"assistant_response":%q}`, text))

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runStopHook(t.TempDir(), &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected explicit bootstrap acknowledgement to pass, got %d; stderr=%s", code, stderr.String())
	}
}

func TestRunStopHookBlocksMissingFeedbackLearningReport(t *testing.T) {
	text := "Bootstrap: rules=✓ phase=phase.bootstrap obligations=23 gates=25\n" +
		"Active per-turn obligations: obligation.cognitive.mode_report, obligation.feedback.learning_report, obligation.finality.close_loop_check\n\n" +
		"Done.\n\n" +
		"Cognitive: NORMAL·SUMMARY_FIRST·STANDARD·NONE / V:CHECKLIST / Cost:LOW / Sig:user_keyword_fast"
	setHookStdin(t, fmt.Sprintf(`{"hook_event_name":"stop","assistant_response":%q}`, text))

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runStopHook(t.TempDir(), &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected Cursor stop to loop with success exit, got %d; stderr=%s", code, stderr.String())
	}
	var output map[string]string
	if err := json.Unmarshal(stdout.Bytes(), &output); err != nil {
		t.Fatalf("decode Cursor stop output: %v\n%s", err, stdout.String())
	}
	if !strings.Contains(output["followup_message"], "Feedback / Learning Report") {
		t.Fatalf("expected followup_message to request Feedback / Learning Report, got %#v", output)
	}
}

func TestRunStopHookAllowsFullFeedbackLearningReport(t *testing.T) {
	text := "Bootstrap: rules=✓ phase=phase.bootstrap obligations=23 gates=25\n" +
		"Active per-turn obligations: obligation.cognitive.mode_report, obligation.feedback.learning_report, obligation.finality.close_loop_check\n\n" +
		"Done.\n\n" +
		"### Feedback / Learning Report\n\n" +
		"| Field | Value |\n" +
		"| --- | --- |\n" +
		"| feedback_decision | NEEDED |\n" +
		"| repo_context | LOCAL |\n" +
		"| writeback_status | COMPLETED |\n" +
		"| target | workflow |\n\n" +
		"Cognitive: NORMAL·SUMMARY_FIRST·STANDARD·NONE / V:CHECKLIST / Cost:LOW / Sig:user_keyword_fast"
	setHookStdin(t, fmt.Sprintf(`{"assistant_response":%q}`, text))

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runStopHook(t.TempDir(), &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected full Feedback / Learning Report to pass, got %d; stderr=%s", code, stderr.String())
	}
	if stdout.String() != "" {
		t.Fatalf("expected no stop output, got %s", stdout.String())
	}
}

func TestRunStopHookBlocksInvalidFeedbackLearningEnum(t *testing.T) {
	text := "Bootstrap: rules=✓ phase=phase.bootstrap obligations=23 gates=25\n" +
		"Active per-turn obligations: obligation.cognitive.mode_report, obligation.feedback.learning_report, obligation.finality.close_loop_check\n\n" +
		"Done.\n\n" +
		compactFeedbackReport("MAYBE", "LOCAL", "N/A", "") +
		"\nCognitive: NORMAL·SUMMARY_FIRST·STANDARD·NONE / V:CHECKLIST / Cost:LOW / Sig:user_keyword_fast"
	setHookStdin(t, fmt.Sprintf(`{"assistant_response":%q}`, text))

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runStopHook(t.TempDir(), &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected Claude stop block via decision JSON to exit 0, got %d; stderr=%s", code, stderr.String())
	}
	var output map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &output); err != nil {
		t.Fatalf("decode Claude stop decision: %v\n%s", err, stdout.String())
	}
	reason, _ := output["reason"].(string)
	if !strings.Contains(reason, "FeedbackDecision") {
		t.Fatalf("expected invalid FeedbackDecision reason, got %#v", output)
	}
}

func TestRunStopHookBlocksFeedbackNeededWithoutTarget(t *testing.T) {
	text := "Bootstrap: rules=✓ phase=phase.bootstrap obligations=23 gates=25\n" +
		"Active per-turn obligations: obligation.cognitive.mode_report, obligation.feedback.learning_report, obligation.finality.close_loop_check\n\n" +
		"Done.\n\n" +
		compactFeedbackReport("NEEDED", "LOCAL", "COMPLETED", "") +
		"\nCognitive: NORMAL·SUMMARY_FIRST·STANDARD·NONE / V:CHECKLIST / Cost:LOW / Sig:user_keyword_fast"
	setHookStdin(t, fmt.Sprintf(`{"assistant_response":%q}`, text))

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runStopHook(t.TempDir(), &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected Claude stop block via decision JSON to exit 0, got %d; stderr=%s", code, stderr.String())
	}
	var output map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &output); err != nil {
		t.Fatalf("decode Claude stop decision: %v\n%s", err, stdout.String())
	}
	reason, _ := output["reason"].(string)
	if !strings.Contains(reason, "requires a non-`none` `Target`") {
		t.Fatalf("expected missing target reason, got %#v", output)
	}
}

func TestRunStopHookRepairPromptSaysCorrectedFinalIsAccepted(t *testing.T) {
	setHookStdin(t, `{"hook_event_name":"stop","assistant_response":"OK"}`)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runStopHook(t.TempDir(), &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected Cursor stop to loop with success exit, got %d; stderr=%s", code, stderr.String())
	}
	var output map[string]string
	if err := json.Unmarshal(stdout.Bytes(), &output); err != nil {
		t.Fatalf("decode Cursor stop output: %v\n%s", err, stdout.String())
	}
	if !strings.Contains(output["followup_message"], "Repair is allowed") {
		t.Fatalf("expected repairable bootstrap wording, got %#v", output)
	}
	if !strings.Contains(output["followup_message"], "corrected final response is accepted as repair") {
		t.Fatalf("expected corrected final repair wording, got %#v", output)
	}
}

func TestRunStopHookAllowsCursorPlanToolResponseWithoutCloseOutLoop(t *testing.T) {
	setHookStdin(t, `{"hook_event_name":"stop","assistant_response":"Plan file created at: /tmp/project_overlay_rules.plan.md\n\nYou can read the plan contents from this file. The provided to-dos have been added to the file as well."}`)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runStopHook(t.TempDir(), &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected Cursor stop to allow non-final plan tool response, got %d; stderr=%s", code, stderr.String())
	}
	if stdout.String() != "" {
		t.Fatalf("expected no followup loop for non-final plan tool response, got %s", stdout.String())
	}
	if !strings.Contains(stderr.String(), "ALLOW_CURSOR_NON_FINAL_TOOL_RESPONSE") {
		t.Fatalf("expected non-final tool response diagnostic, got %s", stderr.String())
	}
}

func TestRunStopHookAllowsCursorTodoToolResponseWithoutCloseOutLoop(t *testing.T) {
	setHookStdin(t, `{"hook_event_name":"stop","assistant_response":"Successfully updated TODOs. Make sure to follow and update your TODO list as you make progress."}`)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runStopHook(t.TempDir(), &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected Cursor stop to allow non-final todo tool response, got %d; stderr=%s", code, stderr.String())
	}
	if stdout.String() != "" {
		t.Fatalf("expected no followup loop for non-final todo tool response, got %s", stdout.String())
	}
	if !strings.Contains(stderr.String(), "ALLOW_CURSOR_NON_FINAL_TOOL_RESPONSE") {
		t.Fatalf("expected non-final tool response diagnostic, got %s", stderr.String())
	}
}

func TestRunStopHookAllowsCursorSwitchModeResponseWithoutCloseOutLoop(t *testing.T) {
	cases := []string{
		"Switched composer mode from agent to plan",
		"Switched composer mode from plan to agent",
		"Switched composer mode from plan to build",
		"Switched from Plan to Build mode.",
		"Switched to Agent mode",
		"Switched to Plan mode",
		"Switched to Build mode",
		"Switched to Ask mode",
		"Switched to Debug mode",
		"You are now in Agent mode.",
		"You are now in Plan mode.",
		"You are now in Build mode.",
		"You are now in Ask mode.",
		"You are now in Debug mode.",
		"Successfully switched to Plan mode.",
		"Successfully switched to Build mode.",
		"Mode switched to plan.",
		"Mode switched to build.",
	}
	for _, assistantResponse := range cases {
		t.Run(assistantResponse, func(t *testing.T) {
			setHookStdin(t, fmt.Sprintf(`{"hook_event_name":"stop","assistant_response":%q}`, assistantResponse))

			var stdout bytes.Buffer
			var stderr bytes.Buffer
			code := runStopHook(t.TempDir(), &stdout, &stderr)
			if code != ExitSuccess {
				t.Fatalf("expected Cursor stop to allow non-final switch-mode response, got %d; stderr=%s", code, stderr.String())
			}
			if stdout.String() != "" {
				t.Fatalf("expected no followup loop for non-final switch-mode response, got %s", stdout.String())
			}
			if !strings.Contains(stderr.String(), "ALLOW_CURSOR_NON_FINAL_TOOL_RESPONSE") {
				t.Fatalf("expected non-final tool response diagnostic, got %s", stderr.String())
			}
		})
	}
}

func TestRunStopHookAllowsCursorSwitchModeTranscriptWithoutCloseOutLoop(t *testing.T) {
	dir := t.TempDir()
	transcriptPath := writeBootstrapTranscript(t, dir, "Successfully switched to Plan mode.", nil)
	setHookStdin(t, fmt.Sprintf(`{"hook_event_name":"stop","transcript_path":%q}`, transcriptPath))

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runStopHook(dir, &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected Cursor transcript stop to allow non-final switch-mode response, got %d; stderr=%s", code, stderr.String())
	}
	if stdout.String() != "" {
		t.Fatalf("expected no followup loop for non-final switch-mode transcript, got %s", stdout.String())
	}
	if !strings.Contains(stderr.String(), "ALLOW_CURSOR_NON_FINAL_TOOL_RESPONSE") {
		t.Fatalf("expected non-final tool response diagnostic, got %s", stderr.String())
	}
}

func TestRunStopHookDoesNotTreatFinalMentioningSwitchModeAsToolStatus(t *testing.T) {
	setHookStdin(t, `{"hook_event_name":"stop","assistant_response":"The earlier Cursor message said Switched to Plan mode, but this is my final answer without a close-out block."}`)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runStopHook(t.TempDir(), &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected Cursor stop to loop with success exit, got %d; stderr=%s", code, stderr.String())
	}
	var output map[string]string
	if err := json.Unmarshal(stdout.Bytes(), &output); err != nil {
		t.Fatalf("decode Cursor stop output: %v\n%s", err, stdout.String())
	}
	if !strings.Contains(output["followup_message"], "Cognitive Mode block") {
		t.Fatalf("expected close-out followup for final response, got %#v", output)
	}
	if strings.Contains(stderr.String(), "ALLOW_CURSOR_NON_FINAL_TOOL_RESPONSE") {
		t.Fatalf("final response mentioning mode switch must not be treated as non-final: %s", stderr.String())
	}
}

func validStopCloseOutText() string {
	return "Bootstrap: rules=✓ phase=phase.bootstrap obligations=23 gates=25\n" +
		"Active per-turn obligations: obligation.cognitive.mode_report, obligation.feedback.learning_report, obligation.finality.close_loop_check\n\n" +
		"Done.\n\n" +
		compactFeedbackReport("NONE", "LOCAL", "N/A", "") +
		"\nCognitive: NORMAL·SUMMARY_FIRST·STANDARD·NONE / V:CHECKLIST / Cost:LOW / Sig:user_keyword_fast"
}

func compactFeedbackReport(decision, repoContext, writeback, target string) string {
	report := "FeedbackDecision: " + decision + "\n" +
		"RepoContext: " + repoContext + "\n" +
		"Writeback: " + writeback
	if strings.TrimSpace(target) != "" {
		report += "\nTarget: " + target
	}
	return report + "\n"
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

// writeBootstrapTranscript builds a minimal Claude Code JSONL transcript with
// the supplied assistant text turn and (optionally) Read tool_use blocks for
// the given file paths. It returns the transcript path.
func writeBootstrapTranscript(t *testing.T, dir string, assistantText string, readPaths []string) string {
	t.Helper()
	path := filepath.Join(dir, "transcript.jsonl")
	var lines []string
	if len(readPaths) > 0 {
		var items []map[string]any
		for i, p := range readPaths {
			items = append(items, map[string]any{
				"type":  "tool_use",
				"id":    fmt.Sprintf("tu_%d", i),
				"name":  "Read",
				"input": map[string]any{"file_path": p},
			})
		}
		entry := map[string]any{
			"type": "assistant",
			"message": map[string]any{
				"content": items,
			},
		}
		buf, err := json.Marshal(entry)
		if err != nil {
			t.Fatalf("marshal read entry: %v", err)
		}
		lines = append(lines, string(buf))
	}
	if assistantText != "" {
		entry := map[string]any{
			"type": "assistant",
			"message": map[string]any{
				"content": []map[string]any{
					{"type": "text", "text": assistantText},
				},
			},
		}
		buf, err := json.Marshal(entry)
		if err != nil {
			t.Fatalf("marshal text entry: %v", err)
		}
		lines = append(lines, string(buf))
	}
	writeFile(t, path, strings.Join(lines, "\n")+"\n")
	return path
}

func TestTranscriptHasRequiredBootstrapReadsDetectsBothReads(t *testing.T) {
	dir := t.TempDir()
	tr := writeBootstrapTranscript(t, dir, "", []string{
		"/repo/CORE_BOOTSTRAP.md",
		"/repo/runtime/core-bootstrap.yaml",
	})
	ok, missing := transcriptHasRequiredBootstrapReads(tr, bootstrapRequiredReadSuffixes)
	if !ok || len(missing) != 0 {
		t.Fatalf("expected ok=true missing=[]; got ok=%v missing=%v", ok, missing)
	}
}

func TestTranscriptHasRequiredBootstrapReadsReportsMissing(t *testing.T) {
	dir := t.TempDir()
	tr := writeBootstrapTranscript(t, dir, "", []string{"/repo/CORE_BOOTSTRAP.md"})
	ok, missing := transcriptHasRequiredBootstrapReads(tr, bootstrapRequiredReadSuffixes)
	if ok {
		t.Fatalf("expected ok=false when one required file missing")
	}
	if len(missing) != 1 || missing[0] != "runtime/core-bootstrap.yaml" {
		t.Fatalf("expected missing=[runtime/core-bootstrap.yaml]; got %v", missing)
	}
}

func TestTranscriptHasRequiredBootstrapReadsAcceptsWindowsPaths(t *testing.T) {
	dir := t.TempDir()
	tr := writeBootstrapTranscript(t, dir, "", []string{
		`C:\yiHong\Programs\Ai-skill\CORE_BOOTSTRAP.md`,
		`C:\yiHong\Programs\Ai-skill\runtime\core-bootstrap.yaml`,
	})
	ok, missing := transcriptHasRequiredBootstrapReads(tr, bootstrapRequiredReadSuffixes)
	if !ok {
		t.Fatalf("expected ok=true for backslash paths; missing=%v", missing)
	}
}

func TestTranscriptHasRequiredBootstrapReadsAcceptsCursorReadFileTools(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "transcript.jsonl")
	entries := []map[string]any{
		{
			"type": "assistant",
			"message": map[string]any{
				"content": []map[string]any{
					{"type": "tool_use", "id": "tu_1", "name": "ReadFile", "input": map[string]any{"path": "/repo/CORE_BOOTSTRAP.md"}},
					{"type": "tool_use", "id": "tu_2", "name": "functions.ReadFile", "input": map[string]any{"path": "/repo/runtime/core-bootstrap.yaml"}},
				},
			},
		},
	}
	var lines []string
	for _, entry := range entries {
		buf, err := json.Marshal(entry)
		if err != nil {
			t.Fatalf("marshal transcript entry: %v", err)
		}
		lines = append(lines, string(buf))
	}
	writeFile(t, path, strings.Join(lines, "\n")+"\n")

	ok, missing := transcriptHasRequiredBootstrapReads(path, bootstrapRequiredReadSuffixes)
	if !ok || len(missing) != 0 {
		t.Fatalf("expected Cursor ReadFile tool_use entries to satisfy bootstrap reads; ok=%v missing=%v", ok, missing)
	}
}

func TestPreToolUseHookBlocksReceiptWithoutReads(t *testing.T) {
	dir := t.TempDir()
	// Transcript contains a Bootstrap Receipt line but NO Read tool_use entries —
	// this models the failure case where the agent copies the format from the
	// SessionStart hook reminder without actually dereferencing the canonical
	// files. The strengthened gate must block.
	tr := writeBootstrapTranscript(t, dir,
		"Bootstrap: rules=✓ phase=phase.bootstrap obligations=2 gates=2\nDone.",
		nil)
	payload := fmt.Sprintf(`{"tool_name":"Bash","transcript_path":%q}`, tr)
	setHookStdin(t, payload)

	var stdout, stderr bytes.Buffer
	code := runPreToolUseHook(dir, &stdout, &stderr)
	// Claude Code blocks PreToolUse via exit 0 + permissionDecision:deny (or exit
	// 2). It does NOT block on other non-zero codes. So a real block MUST be
	// exit 0 with the deny JSON on stdout — NOT ExitValidationFailed(30), which
	// silently let the tool run.
	if code != ExitSuccess {
		t.Fatalf("expected ExitSuccess (deny carried by JSON); got %d; stderr=%s", code, stderr.String())
	}
	assertPreToolUseDeny(t, stdout.String())
	if !strings.Contains(stderr.String(), "BLOCK_RECEIPT_WITHOUT_READS") {
		t.Fatalf("expected BLOCK_RECEIPT_WITHOUT_READS in stderr; got:\n%s", stderr.String())
	}
}

// assertPreToolUseDeny verifies stdout carries a Claude PreToolUse deny decision
// (the only thing that actually blocks besides exit 2).
func assertPreToolUseDeny(t *testing.T, stdout string) {
	t.Helper()
	var out struct {
		HookSpecificOutput struct {
			HookEventName      string `json:"hookEventName"`
			PermissionDecision string `json:"permissionDecision"`
			Reason             string `json:"permissionDecisionReason"`
		} `json:"hookSpecificOutput"`
	}
	if err := json.Unmarshal([]byte(strings.TrimSpace(stdout)), &out); err != nil {
		t.Fatalf("stdout is not a hook decision JSON: %v\n%s", err, stdout)
	}
	if out.HookSpecificOutput.HookEventName != "PreToolUse" {
		t.Fatalf("hookEventName = %q, want PreToolUse", out.HookSpecificOutput.HookEventName)
	}
	if out.HookSpecificOutput.PermissionDecision != "deny" {
		t.Fatalf("permissionDecision = %q, want deny", out.HookSpecificOutput.PermissionDecision)
	}
	if strings.TrimSpace(out.HookSpecificOutput.Reason) == "" {
		t.Fatalf("deny must carry permissionDecisionReason; got empty")
	}
}

// assertCursorPreToolUseDeny verifies stdout carries a Cursor preToolUse deny
// decision (native {permission:"deny"} — the Claude permissionDecision compat
// shim is OFF by default in Cursor, so we must emit the native form).
func assertCursorPreToolUseDeny(t *testing.T, stdout string) {
	t.Helper()
	var out struct {
		Permission   string `json:"permission"`
		UserMessage  string `json:"user_message"`
		AgentMessage string `json:"agent_message"`
	}
	if err := json.Unmarshal([]byte(strings.TrimSpace(stdout)), &out); err != nil {
		t.Fatalf("stdout is not a Cursor permission JSON: %v\n%s", err, stdout)
	}
	if out.Permission != "deny" {
		t.Fatalf("permission = %q, want deny", out.Permission)
	}
	if strings.TrimSpace(out.UserMessage) == "" || strings.TrimSpace(out.AgentMessage) == "" {
		t.Fatalf("deny must carry user_message + agent_message; got %#v", out)
	}
}

func TestRenderCursorPreToolUseDecision_Deny(t *testing.T) {
	var stdout bytes.Buffer
	code := renderCursorPreToolUseDecision(&stdout, hookDecision{Deny: true, Reason: "read primary_source first"})
	if code != ExitSuccess {
		t.Fatalf("deny must return exit 0 (permission JSON carries the block), got %d", code)
	}
	assertCursorPreToolUseDeny(t, stdout.String())
	if !strings.Contains(stdout.String(), "read primary_source first") {
		t.Fatalf("reason must be in user_message/agent_message; got %s", stdout.String())
	}
}

func TestRenderCursorPreToolUseDecision_Allow(t *testing.T) {
	var stdout bytes.Buffer
	code := renderCursorPreToolUseDecision(&stdout, hookDecision{Deny: false})
	if code != ExitSuccess {
		t.Fatalf("allow must be exit 0, got %d", code)
	}
	// Empty stdout == allow in Cursor (proven by the shipped stop hook). We must
	// NOT emit non-empty invalid output (Cursor errors on stdout that is not JSON).
	if strings.TrimSpace(stdout.String()) != "" {
		t.Fatalf("allow must emit NO decision JSON (empty stdout = allow); got %q", stdout.String())
	}
}

func TestDetectPreToolUseHost(t *testing.T) {
	cursor := mustRawPayload(t, `{"hook_event_name":"preToolUse","cursor_version":"3.4.17","tool_name":"run_terminal_cmd"}`)
	if got := detectPreToolUseHost(cursor); got != hostCursor {
		t.Fatalf("cursor_version present => hostCursor, got %d", got)
	}
	// Claude payloads also carry hook_event_name, so that field must NOT trigger
	// cursor detection — only cursor_version does.
	claude := mustRawPayload(t, `{"hook_event_name":"PreToolUse","tool_name":"Bash"}`)
	if got := detectPreToolUseHost(claude); got != hostClaude {
		t.Fatalf("no cursor_version => hostClaude, got %d", got)
	}
}

func TestPreToolUseReadAllowed(t *testing.T) {
	for _, tool := range []string{
		"Read",
		"read_file",
		"ReadFile",
		"functions.ReadFile",
		"list_dir",
		"grep",
		"glob_file_search",
		"codebase_search",
		"Glob",
		"functions.Glob",
		"rg",
		"functions.rg",
		"SemanticSearch",
		"functions.SemanticSearch",
	} {
		if !preToolUseReadAllowed(hostCursor, tool) {
			t.Fatalf("cursor read tool %q must be allowed", tool)
		}
	}
	if preToolUseReadAllowed(hostCursor, "run_terminal_cmd") {
		t.Fatalf("cursor non-read tool must NOT be allowed")
	}
	if !preToolUseReadAllowed(hostClaude, "Read") {
		t.Fatalf("claude Read must be allowed")
	}
	if preToolUseReadAllowed(hostClaude, "read_file") {
		t.Fatalf("claude host must not broaden beyond Read")
	}
}

func mustRawPayload(t *testing.T, s string) map[string]json.RawMessage {
	t.Helper()
	var p map[string]json.RawMessage
	if err := json.Unmarshal([]byte(s), &p); err != nil {
		t.Fatalf("bad payload json: %v", err)
	}
	return p
}

// Cursor preToolUse with a non-read tool and a transcript that lacks a Bootstrap
// Receipt must block via native {permission:"deny"} — the same gate logic as
// Claude, rendered to Cursor's transport.
func TestRunPreToolUseHookCursorBlocksWithoutReceipt(t *testing.T) {
	dir := t.TempDir()
	tr := writeBootstrapTranscript(t, dir, "Working on it (no receipt).", nil)
	payload := fmt.Sprintf(`{"hook_event_name":"preToolUse","cursor_version":"3.4.17","tool_name":"run_terminal_cmd","transcript_path":%q}`, tr)
	setHookStdin(t, payload)

	var stdout, stderr bytes.Buffer
	code := runPreToolUseHook(dir, &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected ExitSuccess (deny carried by permission JSON); got %d; stderr=%s", code, stderr.String())
	}
	assertCursorPreToolUseDeny(t, stdout.String())
	if !strings.Contains(stderr.String(), "BLOCK_NO_RECEIPT") {
		t.Fatalf("expected BLOCK_NO_RECEIPT in stderr; got:\n%s", stderr.String())
	}
}

func TestRunPreToolUseHookCursorAllowsReadTool(t *testing.T) {
	for _, toolName := range []string{"read_file", "ReadFile", "functions.ReadFile"} {
		t.Run(toolName, func(t *testing.T) {
			dir := t.TempDir()
			tr := writeBootstrapTranscript(t, dir, "no receipt yet", nil)
			payload := fmt.Sprintf(`{"hook_event_name":"preToolUse","cursor_version":"3.4.17","tool_name":%q,"transcript_path":%q}`, toolName, tr)
			setHookStdin(t, payload)

			var stdout, stderr bytes.Buffer
			code := runPreToolUseHook(dir, &stdout, &stderr)
			if code != ExitSuccess {
				t.Fatalf("expected %s allowed, got %d; stderr=%s", toolName, code, stderr.String())
			}
			if strings.TrimSpace(stdout.String()) != "" {
				t.Fatalf("read allow must emit empty stdout; got %q", stdout.String())
			}
			if !strings.Contains(stderr.String(), "ALLOW_READ_TOOL") {
				t.Fatalf("expected ALLOW_READ_TOOL diagnostic; got:\n%s", stderr.String())
			}
		})
	}
}

func TestRunPreToolUseHookCursorFailOpenWithoutTranscript(t *testing.T) {
	dir := t.TempDir()
	// No transcript_path: the gate cannot verify and MUST fail open (allow),
	// matching the Claude path and the plan's Phase 1 contract.
	payload := `{"hook_event_name":"preToolUse","cursor_version":"3.4.17","tool_name":"run_terminal_cmd"}`
	setHookStdin(t, payload)

	var stdout, stderr bytes.Buffer
	code := runPreToolUseHook(dir, &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected fail-open allow, got %d; stderr=%s", code, stderr.String())
	}
	if strings.TrimSpace(stdout.String()) != "" {
		t.Fatalf("fail-open allow must emit empty stdout; got %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "ALLOW_NO_TRANSCRIPT") {
		t.Fatalf("expected ALLOW_NO_TRANSCRIPT diagnostic; got:\n%s", stderr.String())
	}
}

// writeCursorE2EAiSkillRepo builds a fake AI_SKILL_REPO directory that satisfies
// resolveClaudeAiSkillRepo (CORE_BOOTSTRAP.md + runtime/core-bootstrap.yaml) and
// carries a minimal routing-registry.yaml with two routes:
//   - route.workflow.web-scrape   (auto-detect; user_signal "web scraping";
//     primary_source workflow/web-scrape/README.md)
//   - route.workflow.ddd          (auto-detect; user_signal "DDD"; used for the
//     miss-vs-conflict matrix — exists so a transcript without web-scrape signals
//     still has a non-degenerate registry, not so it activates by default)
//
// The helper sets AI_SKILL_REPO via t.Setenv so resolveClaudeAiSkillRepo picks up
// the fake repo instead of walking out to the real Ai-skill checkout (which
// would pull in the production routing-registry and make assertions registry-
// dependent).
func writeCursorE2EAiSkillRepo(t *testing.T) string {
	t.Helper()
	repo := t.TempDir()
	if err := os.WriteFile(filepath.Join(repo, "CORE_BOOTSTRAP.md"), []byte("# fake\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(repo, "runtime"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repo, "runtime", "core-bootstrap.yaml"), []byte("schema_version: \"1.0\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(repo, "knowledge", "runtime"), 0o755); err != nil {
		t.Fatal(err)
	}
	registry := `registry_version: knowledge-routing/v2
records:
  - id: route.workflow.web-scrape
    route_type: workflow
    activation_mode: auto-detect
    task_intent: web scraping
    activation_triggers:
      activation_any_of:
        user_signals:
          - web scraping
    primary_source: workflow/web-scrape/README.md
  - id: route.workflow.ddd
    route_type: workflow
    activation_mode: auto-detect
    task_intent: domain modeling
    activation_triggers:
      activation_any_of:
        user_signals:
          - DDD architecture review
    primary_source: workflow/ddd/README.md
`
	if err := os.WriteFile(filepath.Join(repo, "knowledge", "runtime", "routing-registry.yaml"), []byte(registry), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("AI_SKILL_REPO", repo)
	return repo
}

// writeCursorE2ETranscript emits a Cursor-flavored JSONL transcript that:
//   - satisfies the bootstrap read-log gate (Cursor ReadFile tool_use entries for
//     CORE_BOOTSTRAP.md + runtime/core-bootstrap.yaml),
//   - satisfies the receipt-present gate (assistant text containing "Bootstrap:"),
//   - injects userText as the latest user turn (drives the detector),
//   - optionally records additional Read tool_use paths (e.g. workflow primary_source).
func writeCursorE2ETranscript(t *testing.T, dir, userText string, extraReads []string) string {
	t.Helper()
	path := filepath.Join(dir, "transcript.jsonl")
	var lines []string

	// Cursor-style bootstrap reads — exercise the ReadFile tool name that the
	// Bugfix Addendum (plans/active/2026-06-05-0200) added to the read-log gate.
	bootstrapReads := []map[string]any{
		{"type": "tool_use", "id": "tu_b1", "name": "ReadFile", "input": map[string]any{"path": "/repo/CORE_BOOTSTRAP.md"}},
		{"type": "tool_use", "id": "tu_b2", "name": "functions.ReadFile", "input": map[string]any{"path": "/repo/runtime/core-bootstrap.yaml"}},
	}
	for i, p := range extraReads {
		bootstrapReads = append(bootstrapReads, map[string]any{
			"type":  "tool_use",
			"id":    fmt.Sprintf("tu_x%d", i),
			"name":  "ReadFile",
			"input": map[string]any{"path": p},
		})
	}
	bootstrapEntry := map[string]any{
		"type":    "assistant",
		"message": map[string]any{"content": bootstrapReads},
	}
	buf, err := json.Marshal(bootstrapEntry)
	if err != nil {
		t.Fatalf("marshal bootstrap reads: %v", err)
	}
	lines = append(lines, string(buf))

	// Receipt + canonical numbers — drives the receipt-present branch in
	// runPreToolUseHook so the workflow gate (finishPreToolUse) actually runs.
	receiptEntry := map[string]any{
		"type": "assistant",
		"message": map[string]any{
			"content": []map[string]any{
				{"type": "text", "text": "Bootstrap: rules=✓ phase=phase.bootstrap obligations=23 gates=25\nWorking on it."},
			},
		},
	}
	buf, err = json.Marshal(receiptEntry)
	if err != nil {
		t.Fatalf("marshal receipt: %v", err)
	}
	lines = append(lines, string(buf))

	// User turn last — BuildRuntimeContext scans for the latest user role to
	// decide substantive + detector input.
	userEntry := map[string]any{
		"type": "user",
		"message": map[string]any{
			"role":    "user",
			"content": userText,
		},
	}
	buf, err = json.Marshal(userEntry)
	if err != nil {
		t.Fatalf("marshal user: %v", err)
	}
	lines = append(lines, string(buf))

	writeFile(t, path, strings.Join(lines, "\n")+"\n")
	return path
}

// TestRunPreToolUseHookCursorWorkflowGate_LockedUnreadDenies is the Phase 4 e2e
// case for "locked active_route + primary_source NOT yet Read → Cursor deny".
// It composes Cursor payload + Cursor transcript format + workflow gate so a
// regression in any link (host detection, bootstrap read-log Cursor compat,
// detector activation, primary_source lookup, Cursor deny render) fails one
// assertion instead of leaking past all three.
func TestRunPreToolUseHookCursorWorkflowGate_LockedUnreadDenies(t *testing.T) {
	dir := t.TempDir()
	writeCursorE2EAiSkillRepo(t)
	tr := writeCursorE2ETranscript(t, dir, "幫我做 web scraping 抓網站", nil)
	payload := fmt.Sprintf(`{"hook_event_name":"preToolUse","cursor_version":"3.4.17","tool_name":"run_terminal_cmd","transcript_path":%q}`, tr)
	setHookStdin(t, payload)

	var stdout, stderr bytes.Buffer
	code := runPreToolUseHook(dir, &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected ExitSuccess (deny carried by permission JSON); got %d; stderr=%s", code, stderr.String())
	}
	assertCursorPreToolUseDeny(t, stdout.String())
	if !strings.Contains(stderr.String(), "BLOCK_WORKFLOW_PRIMARY_SOURCE") {
		t.Fatalf("expected BLOCK_WORKFLOW_PRIMARY_SOURCE diagnostic; got:\n%s", stderr.String())
	}
	if !strings.Contains(stderr.String(), "route.workflow.web-scrape") {
		t.Fatalf("expected locked route id in diagnostic; got:\n%s", stderr.String())
	}
}

// TestRunPreToolUseHookCursorWorkflowGate_LockedAndReadAllows: same setup as the
// deny case but the transcript ALSO records a Cursor ReadFile of the route's
// primary_source. The workflow gate must self-clear → empty stdout allow.
func TestRunPreToolUseHookCursorWorkflowGate_LockedAndReadAllows(t *testing.T) {
	dir := t.TempDir()
	writeCursorE2EAiSkillRepo(t)
	tr := writeCursorE2ETranscript(t, dir, "幫我做 web scraping 抓網站",
		[]string{"/repo/workflow/web-scrape/README.md"})
	payload := fmt.Sprintf(`{"hook_event_name":"preToolUse","cursor_version":"3.4.17","tool_name":"run_terminal_cmd","transcript_path":%q}`, tr)
	setHookStdin(t, payload)

	var stdout, stderr bytes.Buffer
	code := runPreToolUseHook(dir, &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected allow once primary_source Read; got %d; stderr=%s", code, stderr.String())
	}
	if strings.TrimSpace(stdout.String()) != "" {
		t.Fatalf("workflow-gate allow must emit empty stdout; got %q", stdout.String())
	}
	if strings.Contains(stderr.String(), "BLOCK_WORKFLOW_PRIMARY_SOURCE") {
		t.Fatalf("must NOT block once primary_source Read; stderr=%s", stderr.String())
	}
}

// TestRunPreToolUseHookCursorWorkflowGate_MissAllows: bootstrap satisfied but
// user request matches no route → detector miss → workflow gate fails open.
// Guards against the failure mode the §SAFETY block in workflowPrimarySourceGate
// was written for ("never wedge unrelated tool calls on a detector miss").
func TestRunPreToolUseHookCursorWorkflowGate_MissAllows(t *testing.T) {
	dir := t.TempDir()
	writeCursorE2EAiSkillRepo(t)
	tr := writeCursorE2ETranscript(t, dir, "hi 早安", nil)
	payload := fmt.Sprintf(`{"hook_event_name":"preToolUse","cursor_version":"3.4.17","tool_name":"run_terminal_cmd","transcript_path":%q}`, tr)
	setHookStdin(t, payload)

	var stdout, stderr bytes.Buffer
	code := runPreToolUseHook(dir, &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected fail-open allow on detector miss; got %d; stderr=%s", code, stderr.String())
	}
	if strings.TrimSpace(stdout.String()) != "" {
		t.Fatalf("detector-miss allow must emit empty stdout; got %q", stdout.String())
	}
	if strings.Contains(stderr.String(), "BLOCK_WORKFLOW_PRIMARY_SOURCE") {
		t.Fatalf("detector miss must never block workflow gate; stderr=%s", stderr.String())
	}
}

func TestRenderClaudePreToolUseDecision_Deny(t *testing.T) {
	var stdout bytes.Buffer
	code := renderClaudePreToolUseDecision(&stdout, hookDecision{Deny: true, Reason: "because policy"})
	// MUST be exit 0 — the block is carried by the JSON, not the exit code.
	// Returning a non-zero like 30 is the bug: Claude treats it as non-blocking.
	if code != ExitSuccess {
		t.Fatalf("deny must return exit 0 (JSON carries the block), got %d", code)
	}
	assertPreToolUseDeny(t, stdout.String())
	if !strings.Contains(stdout.String(), "because policy") {
		t.Fatalf("reason must be in permissionDecisionReason; got %s", stdout.String())
	}
}

func TestRenderClaudePreToolUseDecision_Allow(t *testing.T) {
	var stdout bytes.Buffer
	code := renderClaudePreToolUseDecision(&stdout, hookDecision{Deny: false})
	if code != ExitSuccess {
		t.Fatalf("allow must be exit 0, got %d", code)
	}
	if strings.TrimSpace(stdout.String()) != "" {
		t.Fatalf("allow must emit NO decision JSON (normal permission flow); got %q", stdout.String())
	}
}

func TestRenderClaudeStopDecision_Deny(t *testing.T) {
	var stdout bytes.Buffer
	code := renderClaudeStopDecision(&stdout, hookDecision{Deny: true, Reason: "because close-out missing"})
	// MUST be exit 0 — the block is carried by top-level decision:block JSON, not
	// the exit code. Returning a non-zero like 30 is the bug: Claude treats it as
	// a non-blocking Stop error and stops anyway.
	if code != ExitSuccess {
		t.Fatalf("deny must return exit 0 (JSON carries the block), got %d", code)
	}
	var output map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &output); err != nil {
		t.Fatalf("decode stop decision: %v\n%s", err, stdout.String())
	}
	if output["decision"] != "block" {
		t.Fatalf("Stop block must use top-level decision:block (NOT permissionDecision); got %#v", output)
	}
	if reason, _ := output["reason"].(string); !strings.Contains(reason, "because close-out missing") {
		t.Fatalf("reason must be in the reason field; got %#v", output)
	}
}

func TestRenderClaudeStopDecision_Allow(t *testing.T) {
	var stdout bytes.Buffer
	code := renderClaudeStopDecision(&stdout, hookDecision{Deny: false})
	if code != ExitSuccess {
		t.Fatalf("allow must be exit 0, got %d", code)
	}
	if strings.TrimSpace(stdout.String()) != "" {
		t.Fatalf("allow must emit NO decision JSON (normal stop flow); got %q", stdout.String())
	}
}

func TestSessionStartHookEmitsPlaceholderReceiptNotResolvedNumbers(t *testing.T) {
	// Verifies proposal 2: SessionStart MUST NOT inject a pre-rendered Receipt
	// with concrete numbers into the agent prompt. Doing so lets the agent copy
	// the Receipt verbatim without ever Reading canonical files, defeating
	// gate.bootstrap.receipt_present even after the read-log strengthening.
	workspace := t.TempDir()
	writeFile(t, filepath.Join(workspace, "CORE_BOOTSTRAP.md"), "# Bootstrap stub\n")

	var stdout, stderr bytes.Buffer
	if code := runSessionStartHook(workspace, &stdout, &stderr); code != ExitSuccess {
		t.Fatalf("expected success, got %d; stderr=%s", code, stderr.String())
	}
	var output map[string]map[string]string
	if err := json.Unmarshal(stdout.Bytes(), &output); err != nil {
		t.Fatalf("decode: %v\n%s", err, stdout.String())
	}
	ctx := output["hookSpecificOutput"]["additionalContext"]

	// Must contain the placeholder template form.
	for _, want := range []string{
		"<phase_id from phase_machine_init>",
		"<COUNT(*) FROM obligations>",
		"<COUNT(*) FROM gates>",
		"intentionally does NOT print",
	} {
		if !strings.Contains(ctx, want) {
			t.Errorf("expected SessionStart context to contain %q; got:\n%s", want, ctx)
		}
	}

	// Must NOT contain a fully-resolved Receipt — assert no
	// "Bootstrap: rules=✓ ... obligations=<digits> gates=<digits>" pattern.
	resolvedRe := regexp.MustCompile(`Bootstrap: rules=✓ phase=\S+ obligations=\d+ gates=\d+`)
	if resolvedRe.MatchString(ctx) {
		t.Errorf("SessionStart MUST NOT pre-render a Receipt with concrete digits; matched:\n%s",
			resolvedRe.FindString(ctx))
	}
}

func TestPreToolUseHookAllowsReceiptWithVerifiedReads(t *testing.T) {
	dir := t.TempDir()
	tr := writeBootstrapTranscript(t, dir,
		"Bootstrap: rules=✓ phase=phase.bootstrap obligations=2 gates=2\nDone.",
		[]string{
			"/repo/CORE_BOOTSTRAP.md",
			"/repo/runtime/core-bootstrap.yaml",
		})
	payload := fmt.Sprintf(`{"tool_name":"Bash","transcript_path":%q}`, tr)
	setHookStdin(t, payload)

	var stdout, stderr bytes.Buffer
	code := runPreToolUseHook(dir, &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected ExitSuccess when reads and receipt both present; got %d; stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stderr.String(), "ALLOW_RECEIPT_FOUND_WITH_READS") {
		t.Fatalf("expected ALLOW_RECEIPT_FOUND_WITH_READS in stderr; got:\n%s", stderr.String())
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

func TestCompactCognitiveContractRunsStagedValidators(t *testing.T) {
	repo := initTempGitRepo(t)
	writeFile(t, filepath.Join(repo, "runtime", "cognitive-modes-discovery.yaml"), "signals:\n  - name: user_keyword_fast\n  - name: file_diff_runtime_schema\n")
	writeFile(t, filepath.Join(repo, "runtime", "runtime.db"), "not-a-real-db-for-hook-fallback\n")
	runGit(t, repo, "add", "runtime/runtime.db")

	msg := filepath.Join(repo, "COMMIT_EDITMSG")
	writeFile(t, msg, "docs: touch runtime\n\nCognitive: NORMAL·SUMMARY_FIRST·STANDARD·NONE / V:CHECKLIST / Cost:LOW / Sig:user_keyword_fast\n")
	result := runCommitMsgHook(Result{}, repo, []string{msg})
	if result.Status != "blocked" || result.ExitCode != ExitValidationFailed {
		t.Fatalf("expected compact form to block on staged runtime file, got status=%q exit=%d error=%#v", result.Status, result.ExitCode, result.Error)
	}
	if result.Error == nil || !strings.Contains(result.Error.Message, "execution_mode=NORMAL insufficient") {
		t.Fatalf("expected staged validator violation in error, got %#v", result.Error)
	}
}

func TestFullCognitiveContractAllowsRuntimeDiffWithDeepStrict(t *testing.T) {
	repo := initTempGitRepo(t)
	writeFile(t, filepath.Join(repo, "runtime", "cognitive-modes-discovery.yaml"), "signals:\n  - name: file_diff_runtime_schema\n")
	writeFile(t, filepath.Join(repo, "runtime", "runtime.db"), "not-a-real-db-for-hook-fallback\n")
	runGit(t, repo, "add", "runtime/runtime.db")

	msg := filepath.Join(repo, "COMMIT_EDITMSG")
	writeFile(t, msg, `docs: touch runtime

### Cognitive Mode 報告

| 維度 | 值 | 理由 |
|---|---|---|
| execution_mode | DEEP | runtime diff |
| context_mode | SOURCE_BACKED | runtime source |
| governance_mode | STRICT | generated surface risk |
| memory_mode | NONE | no replay |
| validation_mode | SOURCE_BACKED | runtime validation |
| cognitive_cost | HIGH | derived |

activation_reason:
  - file_diff_runtime_schema

Capability summary:
  execution_mode=DEEP -> source-backed reads and linked-update validation.
`)
	result := runCommitMsgHook(Result{}, repo, []string{msg})
	if result.Status == "blocked" || result.ExitCode != 0 {
		t.Fatalf("expected full DEEP/STRICT form to pass, got status=%q exit=%d error=%#v", result.Status, result.ExitCode, result.Error)
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
