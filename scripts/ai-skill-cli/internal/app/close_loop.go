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

type closeLoopOptions struct {
	repoPath    string
	dryRun      bool
	commit      bool
	push        bool
	jsonOutput  bool
	plainOutput bool
}

func runCloseLoop(args []string, stdout io.Writer, stderr io.Writer) int {
	fs := newFlagSet("close-loop", stderr)
	opts := closeLoopOptions{}
	fs.StringVar(&opts.repoPath, "repo", ".", "repository path")
	fs.BoolVar(&opts.dryRun, "dry-run", false, "inspect without committing or pushing")
	fs.BoolVar(&opts.commit, "commit", false, "commit grouped changes")
	fs.BoolVar(&opts.push, "push", false, "push after commit")
	fs.BoolVar(&opts.jsonOutput, "json", false, "write machine-readable JSON output")
	fs.BoolVar(&opts.plainOutput, "plain", false, "write human-readable output")
	if err := fs.Parse(args); err != nil {
		return ExitInvalidUsage
	}
	if opts.jsonOutput && opts.plainOutput {
		_, _ = fmt.Fprintln(stderr, "--json and --plain are mutually exclusive")
		return ExitInvalidUsage
	}
	if opts.push && !opts.commit {
		_, _ = fmt.Fprintln(stderr, "--push requires --commit")
		return ExitInvalidUsage
	}

	result := buildCloseLoopResult(opts)
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

func buildCloseLoopResult(opts closeLoopOptions) Result {
	result := Result{
		Command:        "close-loop",
		Mode:           "dry_run",
		Status:         "success",
		ExitCode:       ExitSuccess,
		Checks:         []Check{},
		PlannedActions: []string{},
		Mutations:      []string{},
	}

	if opts.commit || opts.push {
		result.Status = "blocked"
		result.ExitCode = ExitPartialCloseBlocked
		result.Error = &CommandError{
			Code:        "write_mode_not_implemented",
			Message:     "close-loop currently supports dry-run inspection only.",
			Remediation: "Run with --dry-run until commit and push parity is fixture-backed.",
		}
	}

	gitCheck := checkGit()
	result.Checks = append(result.Checks, gitCheck)
	if gitCheck.Status != "ok" {
		result.Status = "blocked"
		result.ExitCode = ExitMissingDependency
		result.Error = &CommandError{
			Code:        "missing_git",
			Message:     "Git is required for close-loop but was not found in PATH.",
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

	if lock := closeLoopLockCheck(root); lock.Status != "ok" {
		result.Checks = append(result.Checks, lock)
		result.Status = "blocked"
		result.ExitCode = ExitUnsafeRepoState
		result.Error = &CommandError{Code: "active_close_loop_lock", Message: lock.Message, Remediation: "Ask the active agent or user to finish, or clear stale lock only after verifying it is safe."}
		return result
	}
	result.Checks = append(result.Checks, Check{Name: "close_loop_lock", Status: "ok", Message: "no active close-loop lock detected"})

	if unsafe := closeLoopUnsafeStateCheck(root); unsafe.Status != "ok" {
		result.Checks = append(result.Checks, unsafe)
		result.Status = "blocked"
		result.ExitCode = ExitUnsafeRepoState
		result.Error = &CommandError{Code: "unsafe_repo_state", Message: unsafe.Message, Remediation: "Finish or abort the active Git operation before running close-loop."}
		return result
	}
	result.Checks = append(result.Checks, Check{Name: "git_operation", Status: "ok", Message: "no merge, rebase, or cherry-pick state detected"})

	paths, err := closeLoopChangedPaths(root)
	if err != nil {
		result.Status = "blocked"
		result.ExitCode = ExitGeneralFailure
		result.Error = &CommandError{Code: "git_status_failed", Message: err.Error()}
		return result
	}
	if len(paths) == 0 {
		result.Checks = append(result.Checks, Check{Name: "working_tree", Status: "clean", Message: "no dirty paths"})
		return result
	}

	groups := map[string][]string{}
	for _, path := range paths {
		group := closeLoopGroupForPath(path)
		if group == "unknown" {
			result.Checks = append(result.Checks, Check{Name: "owner_group", Status: "failed", Message: path})
			result.Status = "blocked"
			result.ExitCode = ExitValidationFailed
			result.Error = &CommandError{Code: "unknown_owner_group", Message: "unrecognized dirty path: " + path, Remediation: "Add a grouping rule or handle the path manually."}
			return result
		}
		groups[group] = append(groups[group], path)
	}

	groupNames := make([]string, 0, len(groups))
	for group := range groups {
		groupNames = append(groupNames, group)
		sort.Strings(groups[group])
	}
	sort.Strings(groupNames)
	for _, group := range groupNames {
		result.Checks = append(result.Checks, Check{Name: "owner_group", Status: "ok", Message: fmt.Sprintf("%s: %s", group, strings.Join(groups[group], ", "))})
		result.PlannedActions = append(result.PlannedActions, fmt.Sprintf("would process %s group: %s", group, strings.Join(groups[group], ", ")))
	}
	return result
}

func closeLoopRepoRoot(repoPath string) (string, Check) {
	output, err := exec.Command("git", "-C", repoPath, "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", Check{Name: "repo_root", Status: "failed", Message: "not inside a Git work tree"}
	}
	root := strings.TrimSpace(string(output))
	normalized, err := pathutil.NormalizeForReport(root)
	if err != nil {
		return "", Check{Name: "repo_root", Status: "failed", Message: err.Error()}
	}
	return root, Check{Name: "repo_root", Status: "ok", Message: normalized}
}

func closeLoopLockCheck(root string) Check {
	lockDir := filepath.Join(root, ".git", "ai-skill-agent.lock")
	info, err := os.Stat(lockDir)
	if err != nil {
		return Check{Name: "close_loop_lock", Status: "ok"}
	}
	if !info.IsDir() {
		return Check{Name: "close_loop_lock", Status: "blocked", Message: "close-loop lock path exists but is not a directory"}
	}
	owner := "unknown"
	if content, err := os.ReadFile(filepath.Join(lockDir, "owner")); err == nil {
		owner = strings.TrimSpace(string(content))
	}
	return Check{Name: "close_loop_lock", Status: "blocked", Message: "active close-loop lock detected: " + owner}
}

func closeLoopUnsafeStateCheck(root string) Check {
	gitDirOutput, err := exec.Command("git", "-C", root, "rev-parse", "--git-dir").Output()
	if err != nil {
		return Check{Name: "git_operation", Status: "failed", Message: "cannot resolve git dir"}
	}
	gitDir := strings.TrimSpace(string(gitDirOutput))
	if !filepath.IsAbs(gitDir) {
		gitDir = filepath.Join(root, gitDir)
	}
	states := []struct {
		name string
		path string
		dir  bool
	}{
		{name: "merge in progress", path: "MERGE_HEAD"},
		{name: "rebase in progress", path: "rebase-merge", dir: true},
		{name: "rebase/apply in progress", path: "rebase-apply", dir: true},
		{name: "cherry-pick in progress", path: "CHERRY_PICK_HEAD"},
	}
	for _, state := range states {
		path := filepath.Join(gitDir, state.path)
		info, err := os.Stat(path)
		if err != nil {
			continue
		}
		if state.dir && info.IsDir() {
			return Check{Name: "git_operation", Status: "blocked", Message: state.name}
		}
		if !state.dir && !info.IsDir() {
			return Check{Name: "git_operation", Status: "blocked", Message: state.name}
		}
	}
	return Check{Name: "git_operation", Status: "ok"}
}

func closeLoopChangedPaths(root string) ([]string, error) {
	output, err := exec.Command("git", "-C", root, "status", "--porcelain=v1").Output()
	if err != nil {
		return nil, err
	}
	paths := []string{}
	for _, line := range strings.Split(string(output), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		if len(line) < 4 {
			continue
		}
		status := line[0:2]
		path := line[3:]
		if strings.Contains(path, " -> ") {
			parts := strings.Split(path, " -> ")
			path = parts[len(parts)-1]
		}
		if status == "??" && strings.HasSuffix(path, "/") {
			expanded, err := closeLoopUntrackedFiles(root, path)
			if err != nil {
				return nil, err
			}
			paths = append(paths, expanded...)
			continue
		}
		paths = append(paths, path)
	}
	sort.Strings(paths)
	return paths, nil
}

func closeLoopUntrackedFiles(root string, path string) ([]string, error) {
	output, err := exec.Command("git", "-C", root, "ls-files", "--others", "--exclude-standard", "--", path).Output()
	if err != nil {
		return nil, err
	}
	files := []string{}
	for _, line := range strings.Split(string(output), "\n") {
		if strings.TrimSpace(line) != "" {
			files = append(files, line)
		}
	}
	return files, nil
}

func closeLoopGroupForPath(path string) string {
	switch {
	case strings.HasPrefix(path, ".cursor/rules/"), strings.HasPrefix(path, "ai-tools/"):
		return "tooling"
	case strings.HasPrefix(path, "architecture/"),
		strings.HasPrefix(path, "analysis/"),
		strings.HasPrefix(path, "intelligence/"),
		strings.HasPrefix(path, "workflow/"),
		strings.HasPrefix(path, "runtime/"),
		strings.HasPrefix(path, "memory/"),
		strings.HasPrefix(path, "feedback/"),
		strings.HasPrefix(path, "models/"),
		strings.HasPrefix(path, "governance/"),
		strings.HasPrefix(path, "knowledge/"),
		strings.HasPrefix(path, "metadata/"),
		strings.HasPrefix(path, "plans/"),
		strings.HasPrefix(path, "skills/"):
		return "architecture"
	case strings.HasPrefix(path, "enforcement/"), path == "README.md", path == "CONTRIBUTING.md", path == ".gitignore":
		return "shared"
	case strings.HasPrefix(path, "scripts/"):
		return "scripts"
	default:
		return "unknown"
	}
}
