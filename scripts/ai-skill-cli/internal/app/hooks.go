package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/linyihong/Ai-skill/scripts/ai-skill-cli/internal/pathutil"
)

type hooksOptions struct {
	command     string
	repoPath    string
	dryRun      bool
	force       bool
	jsonOutput  bool
	plainOutput bool
	positional  []string
}

func runHooks(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) == 0 {
		_, _ = fmt.Fprintln(stderr, "usage: ai-skill hooks <install|run> [flags]")
		return ExitInvalidUsage
	}
	opts := hooksOptions{command: args[0]}
	if opts.command != "install" && opts.command != "run" {
		_, _ = fmt.Fprintf(stderr, "unsupported hooks command: %s\n", opts.command)
		return ExitInvalidUsage
	}

	name := "hooks " + opts.command
	if opts.command == "run" {
		if len(args) < 2 {
			_, _ = fmt.Fprintln(stderr, "usage: ai-skill hooks run <pre-commit|commit-msg|post-commit|pre-push> [flags]")
			return ExitInvalidUsage
		}
		opts.command = "run " + args[1]
		args = append([]string{"run"}, args[2:]...)
	}
	fs := newFlagSet(name, stderr)
	fs.StringVar(&opts.repoPath, "repo", ".", "repository path")
	fs.BoolVar(&opts.dryRun, "dry-run", false, "preview hook installation without writing")
	fs.BoolVar(&opts.force, "force", false, "allow overwriting existing hook targets")
	fs.BoolVar(&opts.jsonOutput, "json", false, "write machine-readable JSON output")
	fs.BoolVar(&opts.plainOutput, "plain", false, "write human-readable output")
	if err := fs.Parse(args[1:]); err != nil {
		return ExitInvalidUsage
	}
	opts.positional = fs.Args()
	if opts.jsonOutput && opts.plainOutput {
		_, _ = fmt.Fprintln(stderr, "--json and --plain are mutually exclusive")
		return ExitInvalidUsage
	}

	var result Result
	if strings.HasPrefix(opts.command, "run ") {
		result = buildHooksRunResult(opts)
	} else {
		result = buildHooksInstallResult(opts)
	}
	if opts.jsonOutput {
		if err := writeJSON(stdout, result); err != nil {
			_, _ = fmt.Fprintf(stderr, "write output: %v\n", err)
			return ExitGeneralFailure
		}
		return result.ExitCode
	}
	if err := writePlain(stdout, result); err != nil {
		_, _ = fmt.Fprintf(stderr, "write output: %v\n", err)
		return ExitGeneralFailure
	}
	return result.ExitCode
}

func buildHooksInstallResult(opts hooksOptions) Result {
	result := Result{
		Command:        "hooks install",
		Mode:           "dry_run",
		Status:         "success",
		ExitCode:       ExitSuccess,
		Checks:         []Check{},
		PlannedActions: []string{},
		Mutations:      []string{},
	}

	if !opts.dryRun {
		result.Mode = "write"
	}

	gitCheck := checkGit()
	result.Checks = append(result.Checks, gitCheck)
	if gitCheck.Status != "ok" {
		result.Status = "blocked"
		result.ExitCode = ExitMissingDependency
		result.Error = &CommandError{
			Code:        "missing_git",
			Message:     "Git is required for hooks install but was not found in PATH.",
			Remediation: "Install Git and ensure the git executable is available in PATH.",
		}
		return result
	}

	root, repoCheck := closeLoopRepoRoot(opts.repoPath)
	result.Checks = append(result.Checks, repoCheck)
	if repoCheck.Status != "ok" {
		result.Status = "blocked"
		result.ExitCode = ExitInvalidUsage
		result.Error = &CommandError{Code: "invalid_repo", Message: repoCheck.Message, Remediation: "Pass --repo with a valid Git working tree."}
		return result
	}

	operationCheck := closeLoopUnsafeStateCheck(root)
	if operationCheck.Status != "ok" {
		operationCheck.Status = "warning"
		operationCheck.Remediation = "Hook installation planning is allowed, but commit or push must not be triggered while Git has an active operation."
		result.Checks = append(result.Checks, operationCheck)
	} else {
		result.Checks = append(result.Checks, Check{Name: "git_operation", Status: "ok", Message: "no merge, rebase, or cherry-pick state detected"})
	}

	targetDir, targetCheck := gitHooksTargetDir(root)
	result.Checks = append(result.Checks, targetCheck)
	if targetCheck.Status != "ok" {
		result.Status = "blocked"
		result.ExitCode = ExitInvalidUsage
		result.Error = &CommandError{Code: "missing_hooks_target", Message: targetCheck.Message, Remediation: "Use a normal Git working tree with a .git/hooks directory."}
		return result
	}

	hooks := []string{"pre-commit", "commit-msg", "post-commit", "pre-push"}
	for _, hook := range hooks {
		result.PlannedActions = append(result.PlannedActions, fmt.Sprintf("install Go hook adapter: %s", filepath.Join(targetDir, hook)))
	}
	conflicts := existingHookTargets(targetDir, hooks)
	if len(conflicts) > 0 && !opts.force {
		result.Checks = append(result.Checks, Check{Name: "hook_conflicts", Status: "failed", Message: strings.Join(conflicts, ", ")})
		result.Status = "blocked"
		result.ExitCode = ExitInvalidUsage
		result.Error = &CommandError{Code: "target_exists", Message: "hook targets already exist: " + strings.Join(conflicts, ", "), Remediation: "Pass --force only after reviewing the planned overwrite list."}
		return result
	}
	result.Checks = append(result.Checks, Check{Name: "hook_conflicts", Status: "ok", Message: "no blocking hook target conflicts"})
	if !opts.dryRun {
		for _, hook := range hooks {
			target := filepath.Join(targetDir, hook)
			if err := os.WriteFile(target, []byte(hookAdapterContent(hook)), 0o755); err != nil {
				result.Status = "blocked"
				result.ExitCode = ExitGeneralFailure
				result.Error = &CommandError{Code: "hook_install_failed", Message: err.Error()}
				return result
			}
			result.Mutations = append(result.Mutations, "installed hook adapter: "+target)
		}
	}
	return result
}

func buildHooksRunResult(opts hooksOptions) Result {
	result := Result{
		Command:        "hooks " + opts.command,
		Mode:           "write",
		Status:         "success",
		ExitCode:       ExitSuccess,
		Checks:         []Check{},
		PlannedActions: []string{},
		Mutations:      []string{},
	}
	gitCheck := checkGit()
	result.Checks = append(result.Checks, gitCheck)
	if gitCheck.Status != "ok" {
		result.Status = "blocked"
		result.ExitCode = ExitMissingDependency
		result.Error = &CommandError{Code: "missing_git", Message: "Git is required for hook runner but was not found in PATH."}
		return result
	}
	root, repoCheck := closeLoopRepoRoot(opts.repoPath)
	result.Checks = append(result.Checks, repoCheck)
	if repoCheck.Status != "ok" {
		result.Status = "blocked"
		result.ExitCode = ExitInvalidUsage
		result.Error = &CommandError{Code: "invalid_repo", Message: repoCheck.Message}
		return result
	}
	switch opts.command {
	case "run pre-commit":
		return runPreCommitHook(result, root)
	case "run commit-msg":
		return runCommitMsgHook(result, root, opts.positional)
	case "run post-commit":
		if os.Getenv("AI_SKILL_SYNC_CURSOR_BUNDLE") == "1" {
			result.Checks = append(result.Checks, Check{Name: "cursor_bundle_sync", Status: "skipped", Message: "Go sync-cursor-bundle write mode is not enabled"})
		} else {
			result.Checks = append(result.Checks, Check{Name: "cursor_bundle_sync", Status: "skipped", Message: "reference-only default"})
		}
		return result
	case "run pre-push":
		return runPrePushHook(result, root)
	default:
		result.Status = "blocked"
		result.ExitCode = ExitInvalidUsage
		result.Error = &CommandError{Code: "invalid_hook", Message: "unsupported hook runner: " + opts.command}
		return result
	}
}

func runPreCommitHook(result Result, root string) Result {
	staged, err := gitLines(root, "diff", "--cached", "--name-only")
	if err != nil {
		result.Status = "blocked"
		result.ExitCode = ExitGeneralFailure
		result.Error = &CommandError{Code: "staged_lookup_failed", Message: err.Error()}
		return result
	}
	if hasRuntimeSourceChange(staged) {
		var stdout strings.Builder
		var stderr strings.Builder
		code := Run([]string{"runtime", "compile", "--repo", root, "--json"}, &stdout, &stderr)
		if code != ExitSuccess {
			result.Status = "blocked"
			result.ExitCode = ExitValidationFailed
			result.Error = &CommandError{Code: "runtime_compile_failed", Message: compactSummary(stdout.String() + stderr.String())}
			return result
		}
		if _, err := exec.Command("git", "-C", root, "add", filepath.Join(root, "runtime", "runtime.db")).CombinedOutput(); err != nil {
			result.Status = "blocked"
			result.ExitCode = ExitGeneralFailure
			result.Error = &CommandError{Code: "runtime_db_stage_failed", Message: err.Error()}
			return result
		}
		result.Mutations = append(result.Mutations, "compiled and staged runtime/runtime.db")
	}
	if hasRuntimeDBChange(staged) || hasKnowledgeValidationChange(staged) {
		var stdout strings.Builder
		var stderr strings.Builder
		code := Run([]string{"runtime", "validate", "--repo", root, "--json"}, &stdout, &stderr)
		if code != ExitSuccess {
			result.Status = "blocked"
			result.ExitCode = ExitValidationFailed
			result.Error = &CommandError{Code: "runtime_validate_failed", Message: compactSummary(stdout.String() + stderr.String())}
			return result
		}
		result.Checks = append(result.Checks, Check{Name: "runtime_validation", Status: "ok", Message: "staged runtime/knowledge/validation changes passed"})
	}
	if len(result.Mutations) == 0 {
		result.Checks = append(result.Checks, Check{Name: "pre_commit", Status: "ok", Message: "no runtime or knowledge hook action required"})
	}
	return result
}

// runCommitMsgHook enforces Phase 4 behavioral wiring of
// gate.execution.cognitive_mode_resolved. Commit message body must contain
// the literal '### Cognitive Mode 報告' block (template defined in
// models/cognitive-modes/README.md). Mechanical commits may opt out via
// '[skip-cognitive-mode]' in the body. Merge commits auto-skip.
func runCommitMsgHook(result Result, root string, positional []string) Result {
	if len(positional) == 0 {
		// Hook called without message file path; cannot enforce, fail open with warning.
		result.Checks = append(result.Checks, Check{Name: "cognitive_mode_block", Status: "warning", Message: "no commit message file passed; check skipped"})
		return result
	}
	msgPath := positional[0]
	if !filepath.IsAbs(msgPath) {
		msgPath = filepath.Join(root, msgPath)
	}
	body, err := os.ReadFile(msgPath)
	if err != nil {
		result.Status = "blocked"
		result.ExitCode = ExitGeneralFailure
		result.Error = &CommandError{Code: "commit_msg_read_failed", Message: err.Error()}
		return result
	}
	text := string(body)

	// Auto-skip merge commits (git auto-generated header)
	if strings.HasPrefix(strings.TrimLeft(text, " \t\n"), "Merge ") {
		result.Checks = append(result.Checks, Check{Name: "cognitive_mode_block", Status: "skipped", Message: "merge commit auto-skipped"})
		return result
	}

	// Primary path: Cognitive Mode 報告 block present → PASS.
	// Checked BEFORE opt-out marker to avoid false positives when commit body
	// documents/quotes the opt-out token (e.g. "Opt-out via '[skip-cognitive-mode]'").
	if strings.Contains(text, "### Cognitive Mode 報告") {
		result.Checks = append(result.Checks, Check{Name: "cognitive_mode_block", Status: "ok", Message: "Cognitive Mode 報告 present"})
		return result
	}

	// Fallback path: opt-out marker on its own line (require leading whitespace or BOL
	// to reduce false positives from prose mentions). Mechanical commits should
	// place the marker as a standalone trailer line.
	for _, line := range strings.Split(text, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "[skip-cognitive-mode]" {
			result.Checks = append(result.Checks, Check{Name: "cognitive_mode_block", Status: "skipped", Message: "[skip-cognitive-mode] opt-out marker present on its own line"})
			return result
		}
	}

	result.Status = "blocked"
	result.ExitCode = ExitValidationFailed
	result.Error = &CommandError{
		Code:        "cognitive_mode_block_missing",
		Message:     "Commit message body must include '### Cognitive Mode 報告' block (4-dim execution/context/governance/memory resolution).",
		Remediation: "Add the block per models/cognitive-modes/README.md template, or add a standalone '[skip-cognitive-mode]' trailer line for mechanical commits.",
	}
	return result
}

func runPrePushHook(result Result, root string) Result {
	changed, upstream, err := cliCIPrePushPaths(root)
	if err != nil {
		result.Checks = append(result.Checks, Check{Name: "cli_ci_scope", Status: "warning", Message: err.Error(), Remediation: "Running Go tests conservatively because changed paths could not be resolved."})
		changed = []string{"scripts/ai-skill-cli/"}
	} else {
		result.Checks = append(result.Checks, Check{Name: "cli_ci_scope", Status: "ok", Message: "compared against " + upstream})
	}
	if !hasCLICIPreflightChange(changed) {
		result.Checks = append(result.Checks, Check{Name: "cli_ci_preflight", Status: "skipped", Message: "no CLI, hook, or workflow changes since upstream"})
		return result
	}
	result.Checks = append(result.Checks, githubWorkflowHistoryCheck(root))
	cmd := exec.Command("go", "test", "./...")
	cmd.Dir = filepath.Join(root, "scripts", "ai-skill-cli")
	output, err := cmd.CombinedOutput()
	if err != nil {
		result.Status = "blocked"
		result.ExitCode = ExitValidationFailed
		result.Error = &CommandError{
			Code:        "cli_go_test_failed",
			Message:     compactSummary(string(output)),
			Remediation: "Run `cd scripts/ai-skill-cli && go test ./...`; if repo-local binaries are stale, rebuild `bin/` with `go run ./cmd/releasebuild --stable-names --version \"repo-$(git rev-parse --short HEAD)\" --commit \"$(git rev-parse --short HEAD)\" --dist bin`.",
		}
		return result
	}
	result.Checks = append(result.Checks, Check{Name: "cli_ci_preflight", Status: "ok", Message: "go test ./... passed"})
	return result
}

type githubWorkflowRunsResponse struct {
	WorkflowRuns []struct {
		Status     string `json:"status"`
		Conclusion string `json:"conclusion"`
		HTMLURL    string `json:"html_url"`
		HeadSHA    string `json:"head_sha"`
	} `json:"workflow_runs"`
}

func githubWorkflowHistoryCheck(root string) Check {
	remote, err := exec.Command("git", "-C", root, "config", "--get", "remote.origin.url").Output()
	if err != nil {
		return Check{Name: "github_workflow_history", Status: "skipped", Message: "remote.origin.url unavailable"}
	}
	owner, repo, ok := parseGitHubRemote(strings.TrimSpace(string(remote)))
	if !ok {
		return Check{Name: "github_workflow_history", Status: "skipped", Message: "origin is not a GitHub remote"}
	}
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/actions/workflows/ai-skill-cli.yml/runs?per_page=1&status=completed", owner, repo)
	client := http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return Check{Name: "github_workflow_history", Status: "skipped", Message: err.Error()}
	}
	req.Header.Set("User-Agent", "ai-skill-cli-pre-push")
	resp, err := client.Do(req)
	if err != nil {
		return Check{Name: "github_workflow_history", Status: "warning", Message: "could not query latest completed GitHub workflow: " + err.Error(), Remediation: "Continuing with local Go preflight; check GitHub Actions manually if needed."}
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return Check{Name: "github_workflow_history", Status: "warning", Message: fmt.Sprintf("GitHub workflow query returned HTTP %d", resp.StatusCode), Remediation: "Continuing with local Go preflight; check GitHub Actions manually if needed."}
	}
	var runs githubWorkflowRunsResponse
	if err := json.NewDecoder(resp.Body).Decode(&runs); err != nil {
		return Check{Name: "github_workflow_history", Status: "warning", Message: "could not decode GitHub workflow response: " + err.Error(), Remediation: "Continuing with local Go preflight; check GitHub Actions manually if needed."}
	}
	if len(runs.WorkflowRuns) == 0 {
		return Check{Name: "github_workflow_history", Status: "skipped", Message: "no completed ai-skill CLI workflow runs found"}
	}
	run := runs.WorkflowRuns[0]
	message := fmt.Sprintf("latest completed run %s at %s", run.Conclusion, run.HTMLURL)
	if run.Conclusion != "success" {
		return Check{Name: "github_workflow_history", Status: "warning", Message: message, Remediation: "Previous completed GitHub workflow is red; local go test must pass before pushing a fix."}
	}
	return Check{Name: "github_workflow_history", Status: "ok", Message: message}
}

func parseGitHubRemote(remote string) (string, string, bool) {
	remote = strings.TrimSuffix(strings.TrimSpace(remote), ".git")
	switch {
	case strings.HasPrefix(remote, "https://github.com/"):
		parts := strings.Split(strings.TrimPrefix(remote, "https://github.com/"), "/")
		if len(parts) >= 2 {
			return parts[0], parts[1], true
		}
	case strings.HasPrefix(remote, "git@github.com:"):
		parts := strings.Split(strings.TrimPrefix(remote, "git@github.com:"), "/")
		if len(parts) >= 2 {
			return parts[0], parts[1], true
		}
	case strings.HasPrefix(remote, "ssh://git@github.com/"):
		parts := strings.Split(strings.TrimPrefix(remote, "ssh://git@github.com/"), "/")
		if len(parts) >= 2 {
			return parts[0], parts[1], true
		}
	}
	return "", "", false
}

func cliCIPrePushPaths(root string) ([]string, string, error) {
	upstreamOutput, err := exec.Command("git", "-C", root, "rev-parse", "--abbrev-ref", "--symbolic-full-name", "@{u}").Output()
	if err != nil {
		changed, fallbackErr := gitLines(root, "diff", "--name-only", "HEAD^...HEAD")
		if fallbackErr != nil {
			return []string{}, "HEAD", nil
		}
		return changed, "HEAD^...HEAD", nil
	}
	upstream := strings.TrimSpace(string(upstreamOutput))
	changed, err := gitLines(root, "diff", "--name-only", upstream+"...HEAD")
	if err != nil {
		return nil, upstream, err
	}
	return changed, upstream, nil
}

func hasCLICIPreflightChange(paths []string) bool {
	for _, path := range paths {
		if strings.HasPrefix(path, "scripts/ai-skill-cli/") ||
			strings.HasPrefix(path, "scripts/git-hooks/") ||
			path == ".github/workflows/ai-skill-cli.yml" {
			return true
		}
	}
	return false
}

func gitLines(root string, args ...string) ([]string, error) {
	output, err := exec.Command("git", append([]string{"-C", root}, args...)...).Output()
	if err != nil {
		return nil, err
	}
	lines := []string{}
	for _, line := range strings.Split(string(output), "\n") {
		if strings.TrimSpace(line) != "" {
			lines = append(lines, strings.TrimSpace(line))
		}
	}
	return lines, nil
}

func hasRuntimeSourceChange(paths []string) bool {
	return false
}

func hasRuntimeDBChange(paths []string) bool {
	for _, path := range paths {
		if path == "runtime/runtime.db" {
			return true
		}
	}
	return false
}

func hasKnowledgeValidationChange(paths []string) bool {
	for _, path := range paths {
		if strings.HasPrefix(path, "knowledge/") || strings.HasPrefix(path, "validation/") || strings.HasPrefix(path, "scripts/validate-") {
			return true
		}
	}
	return false
}

func hookAdapterContent(hook string) string {
	return fmt.Sprintf(`#!/usr/bin/env sh
set -eu
ROOT="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
BIN="${AI_SKILL_CLI:-$ROOT/scripts/ai-skill-cli/bin/ai-skill}"
if [ ! -x "$BIN" ]; then
  case "$(uname -s 2>/dev/null | tr '[:upper:]' '[:lower:]')" in
    darwin) os=darwin ;;
    linux) os=linux ;;
    mingw*|msys*|cygwin*) os=windows ;;
    *) os=unknown ;;
  esac
  arch="$(uname -m 2>/dev/null || echo unknown)"
  case "$arch" in
    arm64|aarch64) arch=arm64 ;;
    x86_64|amd64) arch=amd64 ;;
  esac
  suffix=""
  [ "$os" = "windows" ] && suffix=".exe"
  BIN="$ROOT/scripts/ai-skill-cli/bin/ai-skill-$os-$arch$suffix"
fi
exec "$BIN" hooks run %s --repo "$ROOT" "$@"
`, hook)
}

func hookSourceCheck(sourceDir string) Check {
	normalized, err := pathutil.NormalizeForReport(sourceDir)
	if err != nil {
		normalized = sourceDir
	}
	info, err := os.Stat(sourceDir)
	if err != nil {
		return Check{Name: "hook_source", Status: "missing", Message: normalized}
	}
	if !info.IsDir() {
		return Check{Name: "hook_source", Status: "failed", Message: normalized + " is not a directory"}
	}
	hooks := listHookFiles(sourceDir)
	if len(hooks) == 0 {
		return Check{Name: "hook_source", Status: "failed", Message: normalized + " contains no hook files"}
	}
	return Check{Name: "hook_source", Status: "ok", Message: fmt.Sprintf("%s (%d hooks)", normalized, len(hooks))}
}

func gitHooksTargetDir(root string) (string, Check) {
	gitDirOutput, err := exec.Command("git", "-C", root, "rev-parse", "--git-dir").Output()
	if err != nil {
		return "", Check{Name: "hook_target", Status: "failed", Message: "cannot resolve git dir"}
	}
	gitDir := strings.TrimSpace(string(gitDirOutput))
	if !filepath.IsAbs(gitDir) {
		gitDir = filepath.Join(root, gitDir)
	}
	targetDir := filepath.Join(gitDir, "hooks")
	normalized, err := pathutil.NormalizeForReport(targetDir)
	if err != nil {
		normalized = targetDir
	}
	info, err := os.Stat(targetDir)
	if err != nil {
		return "", Check{Name: "hook_target", Status: "missing", Message: normalized}
	}
	if !info.IsDir() {
		return "", Check{Name: "hook_target", Status: "failed", Message: normalized + " is not a directory"}
	}
	return targetDir, Check{Name: "hook_target", Status: "ok", Message: normalized}
}

func listHookFiles(sourceDir string) []string {
	entries, err := os.ReadDir(sourceDir)
	if err != nil {
		return nil
	}
	hooks := []string{}
	for _, entry := range entries {
		if entry.Type().IsRegular() {
			hooks = append(hooks, entry.Name())
		}
	}
	sort.Strings(hooks)
	return hooks
}

func existingHookTargets(targetDir string, hooks []string) []string {
	conflicts := []string{}
	for _, hook := range hooks {
		target := filepath.Join(targetDir, hook)
		if _, err := os.Stat(target); err == nil {
			conflicts = append(conflicts, target)
		}
	}
	sort.Strings(conflicts)
	return conflicts
}
