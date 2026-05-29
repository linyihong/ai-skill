package app

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
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
	if len(result.PlannedActions) != 6 {
		t.Fatalf("expected 6 planned actions, got %#v", result.PlannedActions)
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
	if result.Mode != "write" || len(result.Mutations) != 7 {
		t.Fatalf("expected write mutations, got %#v", result)
	}
	claudePath := filepath.Join(project, "CLAUDE.md")
	if !pathExists(claudePath) {
		t.Fatal("write mode did not create CLAUDE.md")
	}
	repo, repoCheck := resolveInitProjectAiSkillRepo()
	if repoCheck.Status != "ok" {
		t.Fatalf("expected repo resolution ok, got %#v", repoCheck)
	}
	claudeContent, err := os.ReadFile(claudePath)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(claudeContent), filepath.ToSlash(repo)) {
		t.Fatalf("CLAUDE.md must not contain local absolute Ai-skill path, got %s", string(claudeContent))
	}
	if !strings.Contains(string(claudeContent), "<AI_SKILL_REPO>/CORE_BOOTSTRAP.md") {
		t.Fatalf("expected portable placeholder in CLAUDE.md, got %s", string(claudeContent))
	}
	if !strings.Contains(string(claudeContent), "AI_SKILL_REPO") || !strings.Contains(string(claudeContent), "Windows PowerShell") {
		t.Fatalf("expected cross-platform AI_SKILL_REPO setup guidance in CLAUDE.md, got %s", string(claudeContent))
	}
	if !strings.Contains(string(claudeContent), "Pointer 展開規則") || !strings.Contains(string(claudeContent), "final close-out Cognitive Mode") {
		t.Fatalf("expected explicit pointer expansion guidance in CLAUDE.md, got %s", string(claudeContent))
	}
	claudeSettingsPath := filepath.Join(project, ".claude", "settings.json")
	if !pathExists(claudeSettingsPath) {
		t.Fatal("write mode did not create Claude Code hook settings")
	}
	claudeSettings, err := os.ReadFile(claudeSettingsPath)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(claudeSettings), filepath.ToSlash(repo)) {
		t.Fatalf("Claude settings must not contain local absolute Ai-skill path, got %s", string(claudeSettings))
	}
	if !strings.Contains(string(claudeSettings), "user-prompt-submit") || !strings.Contains(string(claudeSettings), "AI_SKILL_REPO") {
		t.Fatalf("expected Claude hooks to call Ai-skill Go runner, got %s", string(claudeSettings))
	}
	if !strings.Contains(string(claudeSettings), ".ai-skill/local.env") {
		t.Fatalf("expected Claude hooks to source project-local env, got %s", string(claudeSettings))
	}
	if !pathExists(filepath.Join(project, ".cursor", "rules", "ai-skill-bootstrap.mdc")) {
		t.Fatal("write mode did not create Cursor rule")
	}
	cursorHooksPath := filepath.Join(project, ".cursor", "hooks.json")
	cursorHooks, err := os.ReadFile(cursorHooksPath)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(cursorHooks), filepath.ToSlash(repo)) {
		t.Fatalf("Cursor hooks must not contain local absolute Ai-skill path, got %s", string(cursorHooks))
	}
	if !strings.Contains(string(cursorHooks), `"stop"`) ||
		!strings.Contains(string(cursorHooks), "hooks run stop") ||
		!strings.Contains(string(cursorHooks), `"failClosed": true`) ||
		!strings.Contains(string(cursorHooks), ".ai-skill/local.env") {
		t.Fatalf("expected Cursor hooks to enforce final Cognitive close-out through Ai-skill CLI, got %s", string(cursorHooks))
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
	localIgnore, err := os.ReadFile(filepath.Join(project, ".ai-skill", ".gitignore"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(localIgnore), "local.env") && !strings.Contains(string(localIgnore), "*") {
		t.Fatalf("expected .ai-skill/.gitignore to ignore local files, got %s", string(localIgnore))
	}
	localEnvInfo, err := os.Stat(filepath.Join(project, ".ai-skill", "local.env"))
	if err != nil {
		t.Fatal(err)
	}
	// Windows does not honor Unix-style permission bits on os.WriteFile;
	// Mode().Perm() returns 0o666 / 0o444 regardless of the requested mode.
	// The 0o600 security guarantee is still made by initProjectWriteFile;
	// only the assertion is Unix-specific.
	if runtime.GOOS != "windows" {
		if localEnvInfo.Mode().Perm() != 0o600 {
			t.Fatalf("expected local.env mode 0600, got %o", localEnvInfo.Mode().Perm())
		}
	}
	localEnv, err := os.ReadFile(filepath.Join(project, ".ai-skill", "local.env"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(localEnv), "export AI_SKILL_REPO=") || !strings.Contains(string(localEnv), repo) {
		t.Fatalf("expected local.env to set current Ai-skill repo, got %s", string(localEnv))
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
	// AGENTS.md is now generic agent entry that routes via ai-tools/README.md
	// (not direct to ai-tools/agent/codex.md). Codex users reach the adapter
	// through the routing hub. See bootstrap-contract-yaml-migration plan.
	if !strings.Contains(string(content), "ai-tools/README.md") {
		t.Fatalf("expected AGENTS.md to route via ai-tools/README.md, got %s", string(content))
	}
	if strings.Contains(string(content), "ai-tools/agent/codex.md") {
		t.Fatalf("AGENTS.md must NOT direct-link to codex.md (generic entry locks would break multi-agent use), got %s", string(content))
	}
	repo, repoCheck := resolveInitProjectAiSkillRepo()
	if repoCheck.Status != "ok" {
		t.Fatalf("expected repo resolution ok, got %#v", repoCheck)
	}
	if strings.Contains(string(content), filepath.ToSlash(repo)) {
		t.Fatalf("AGENTS.md must not contain local absolute Ai-skill path, got %s", string(content))
	}
	if !strings.Contains(string(content), "AI_SKILL_REPO") || !strings.Contains(string(content), "Windows PowerShell") {
		t.Fatalf("expected cross-platform AI_SKILL_REPO setup guidance in AGENTS.md, got %s", string(content))
	}
	if !strings.Contains(string(content), "Pointer 展開規則") || !strings.Contains(string(content), "final close-out Cognitive Mode") {
		t.Fatalf("expected explicit pointer expansion guidance in AGENTS.md, got %s", string(content))
	}
}

func TestInitProjectWritesCopilotBootstrap(t *testing.T) {
	project := t.TempDir()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"init-project", "--project", project, "--tools", "copilot", "--json"}, &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected write success, got %d; stderr=%s", code, stderr.String())
	}

	projectInstructionsPath := filepath.Join(project, ".github", "copilot-instructions.md")
	scopedInstructionsPath := filepath.Join(project, ".github", "instructions", "ai-skill.instructions.md")
	copilotReadmePath := filepath.Join(project, ".copilot", "README.md")
	copilotPromptPath := filepath.Join(project, ".copilot", "bootstrap-prompt.md")
	copilotScriptPath := filepath.Join(project, ".copilot", "start-copilot.sh")
	projectInstructions, err := os.ReadFile(projectInstructionsPath)
	if err != nil {
		t.Fatal(err)
	}
	scopedInstructions, err := os.ReadFile(scopedInstructionsPath)
	if err != nil {
		t.Fatal(err)
	}
	copilotReadme, err := os.ReadFile(copilotReadmePath)
	if err != nil {
		t.Fatal(err)
	}
	copilotPrompt, err := os.ReadFile(copilotPromptPath)
	if err != nil {
		t.Fatal(err)
	}
	copilotScript, err := os.ReadFile(copilotScriptPath)
	if err != nil {
		t.Fatal(err)
	}
	repo, repoCheck := resolveInitProjectAiSkillRepo()
	if repoCheck.Status != "ok" {
		t.Fatalf("expected repo resolution ok, got %#v", repoCheck)
	}
	for path, content := range map[string]string{
		projectInstructionsPath: string(projectInstructions),
		scopedInstructionsPath:  string(scopedInstructions),
		copilotReadmePath:       string(copilotReadme),
		copilotPromptPath:       string(copilotPrompt),
		copilotScriptPath:       string(copilotScript),
	} {
		if strings.Contains(content, filepath.ToSlash(repo)) {
			t.Fatalf("%s must not contain local absolute Ai-skill path, got %s", path, content)
		}
		if strings.Contains(content, "phase.bootstrap obligations=") {
			t.Fatalf("Copilot instructions must not copy Bootstrap Receipt examples, got %s", content)
		}
	}
	if !strings.Contains(string(projectInstructions), "在回覆任何使用者請求前") ||
		!strings.Contains(string(projectInstructions), "不得因為請求看似只是列檔") {
		t.Fatalf("expected project instructions to make bootstrap mandatory before any request, got %s", string(projectInstructions))
	}
	if !strings.Contains(string(scopedInstructions), "Before answering any user request") ||
		!strings.Contains(string(scopedInstructions), "simple file listings") {
		t.Fatalf("expected scoped instructions to reject simple-query bootstrap exemptions, got %s", string(scopedInstructions))
	}
	for path, content := range map[string]string{
		projectInstructionsPath: string(projectInstructions),
		scopedInstructionsPath:  string(scopedInstructions),
		copilotReadmePath:       string(copilotReadme),
		copilotPromptPath:       string(copilotPrompt),
	} {
		if !strings.Contains(content, "<AI_SKILL_REPO>/CORE_BOOTSTRAP.md") ||
			!strings.Contains(content, "<AI_SKILL_REPO>/runtime/core-bootstrap.yaml") ||
			!strings.Contains(content, "ai-tools/agent/copilot.md") {
			t.Fatalf("expected %s to point to canonical sources, got %s", path, content)
		}
	}
	if !strings.Contains(string(copilotScript), "Deletion condition") ||
		!strings.Contains(string(copilotScript), "repo-local ai-skill binary") ||
		!strings.Contains(string(copilotScript), "copilot start --project") {
		t.Fatalf("expected Copilot script to be a temporary repo-local CLI wrapper, got %s", string(copilotScript))
	}
	if runtime.GOOS != "windows" {
		info, err := os.Stat(copilotScriptPath)
		if err != nil {
			t.Fatal(err)
		}
		if info.Mode().Perm() != 0o755 {
			t.Fatalf("expected start-copilot.sh mode 0755, got %o", info.Mode().Perm())
		}
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
