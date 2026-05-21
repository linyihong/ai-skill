package app

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/linyihong/Ai-skill/scripts/ai-skill-cli/internal/pathutil"
)

type hooksOptions struct {
	command     string
	repoPath    string
	dryRun      bool
	force       bool
	jsonOutput  bool
	plainOutput bool
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
			_, _ = fmt.Fprintln(stderr, "usage: ai-skill hooks run <pre-commit|post-commit> [flags]")
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

	hooks := []string{"pre-commit", "post-commit"}
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
	case "run post-commit":
		if os.Getenv("AI_SKILL_SYNC_CURSOR_BUNDLE") == "1" {
			result.Checks = append(result.Checks, Check{Name: "cursor_bundle_sync", Status: "skipped", Message: "Go sync-cursor-bundle write mode is not enabled"})
		} else {
			result.Checks = append(result.Checks, Check{Name: "cursor_bundle_sync", Status: "skipped", Message: "reference-only default"})
		}
		return result
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
exec "$BIN" hooks run %s --repo "$ROOT"
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
