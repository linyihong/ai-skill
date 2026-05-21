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
		_, _ = fmt.Fprintln(stderr, "usage: ai-skill hooks install [flags]")
		return ExitInvalidUsage
	}
	opts := hooksOptions{command: args[0]}
	if opts.command != "install" {
		_, _ = fmt.Fprintf(stderr, "unsupported hooks command: %s\n", opts.command)
		return ExitInvalidUsage
	}

	fs := newFlagSet("hooks install", stderr)
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

	result := buildHooksInstallResult(opts)
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
		result.Status = "blocked"
		result.ExitCode = ExitPartialCloseBlocked
		result.Error = &CommandError{
			Code:        "write_mode_not_implemented",
			Message:     "hooks install currently supports dry-run planning only.",
			Remediation: "Run with --dry-run until hook copy and chmod parity is fixture-backed.",
		}
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

	sourceDir := filepath.Join(root, ".githooks")
	sourceCheck := hookSourceCheck(sourceDir)
	result.Checks = append(result.Checks, sourceCheck)
	if sourceCheck.Status != "ok" {
		result.Status = "blocked"
		result.ExitCode = ExitValidationFailed
		result.Error = &CommandError{Code: "missing_hook_source", Message: sourceCheck.Message, Remediation: "Create .githooks/ with hook files before planning installation."}
		return result
	}

	targetDir, targetCheck := gitHooksTargetDir(root)
	result.Checks = append(result.Checks, targetCheck)
	if targetCheck.Status != "ok" {
		result.Status = "blocked"
		result.ExitCode = ExitInvalidUsage
		result.Error = &CommandError{Code: "missing_hooks_target", Message: targetCheck.Message, Remediation: "Use a normal Git working tree with a .git/hooks directory."}
		return result
	}

	hooks := listHookFiles(sourceDir)
	for _, hook := range hooks {
		result.PlannedActions = append(result.PlannedActions, fmt.Sprintf("install hook: %s -> %s", filepath.Join(sourceDir, hook), filepath.Join(targetDir, hook)))
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
	return result
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
