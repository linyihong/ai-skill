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
		result.Mode = "write"
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
	if opts.commit {
		result = closeLoopCommitGroups(result, root, groupNames, groups, opts.push)
	}
	return result
}

func closeLoopCommitGroups(result Result, root string, groupNames []string, groups map[string][]string, push bool) Result {
	release, err := acquireCloseLoopLock(root)
	if err != nil {
		result.Status = "blocked"
		result.ExitCode = ExitUnsafeRepoState
		result.Error = &CommandError{Code: "active_close_loop_lock", Message: err.Error()}
		return result
	}
	defer release()
	if check := closeLoopPlanCompletionCheck(root); check.Status == "blocked" {
		result.Checks = append(result.Checks, check)
		result.Status = "blocked"
		result.ExitCode = ExitValidationFailed
		result.Error = &CommandError{Code: "plan_closure_required", Message: check.Message}
		return result
	}
	result.Checks = append(result.Checks, Check{Name: "plan_closure", Status: "ok", Message: "no unarchived completed active plan detected"})
	if check := closeLoopRuntimeValidation(root); check.Status == "failed" {
		result.Checks = append(result.Checks, check)
		result.Status = "blocked"
		result.ExitCode = ExitValidationFailed
		result.Error = &CommandError{Code: "runtime_validation_failed", Message: check.Message}
		return result
	} else {
		result.Checks = append(result.Checks, check)
	}
	for _, group := range groupNames {
		paths := groups[group]
		if err := closeLoopScanPrivatePaths(root, paths); err != nil {
			result.Status = "blocked"
			result.ExitCode = ExitValidationFailed
			result.Error = &CommandError{Code: "private_path_detected", Message: err.Error()}
			return result
		}
		args := append([]string{"-C", root, "add", "--"}, paths...)
		if output, err := exec.Command("git", args...).CombinedOutput(); err != nil {
			result.Status = "blocked"
			result.ExitCode = ExitGeneralFailure
			result.Error = &CommandError{Code: "git_add_failed", Message: strings.TrimSpace(string(output))}
			return result
		}
		msg := closeLoopCommitMessageForGroup(group)
		if output, err := exec.Command("git", "-C", root, "commit", "-m", msg).CombinedOutput(); err != nil {
			result.Status = "blocked"
			result.ExitCode = ExitGeneralFailure
			result.Error = &CommandError{Code: "git_commit_failed", Message: strings.TrimSpace(string(output))}
			return result
		}
		result.Mutations = append(result.Mutations, "committed "+group+": "+msg)
	}
	if push {
		branchOutput, err := exec.Command("git", "-C", root, "branch", "--show-current").Output()
		if err != nil {
			result.Status = "blocked"
			result.ExitCode = ExitGeneralFailure
			result.Error = &CommandError{Code: "branch_lookup_failed", Message: err.Error()}
			return result
		}
		branch := strings.TrimSpace(string(branchOutput))
		if output, err := exec.Command("git", "-C", root, "push", "origin", branch).CombinedOutput(); err != nil {
			result.Status = "blocked"
			result.ExitCode = ExitGeneralFailure
			result.Error = &CommandError{Code: "git_push_failed", Message: strings.TrimSpace(string(output))}
			return result
		}
		result.Mutations = append(result.Mutations, "pushed branch: "+branch)
	}
	if paths, err := closeLoopChangedPaths(root); err == nil && len(paths) == 0 {
		result.Checks = append(result.Checks, Check{Name: "readback_status", Status: "ok", Message: "working tree clean after close-loop"})
	} else if err != nil {
		result.Status = "blocked"
		result.ExitCode = ExitGeneralFailure
		result.Error = &CommandError{Code: "readback_failed", Message: err.Error()}
	} else {
		result.Status = "blocked"
		result.ExitCode = ExitValidationFailed
		result.Error = &CommandError{Code: "dirty_after_commit", Message: strings.Join(paths, ", ")}
	}
	return result
}

func acquireCloseLoopLock(root string) (func(), error) {
	lockDir := filepath.Join(root, ".git", "ai-skill-agent.lock")
	if _, err := os.Stat(lockDir); err == nil {
		return nil, fmt.Errorf("active close-loop lock detected")
	}
	if err := os.Mkdir(lockDir, 0o755); err != nil {
		return nil, err
	}
	_ = os.WriteFile(filepath.Join(lockDir, "pid"), []byte(fmt.Sprintf("%d\n", os.Getpid())), 0o644)
	_ = os.WriteFile(filepath.Join(lockDir, "owner"), []byte(defaultGoalOwner()+"\n"), 0o644)
	_ = os.WriteFile(filepath.Join(lockDir, "startedAt"), []byte(nowUTC()+"\n"), 0o644)
	return func() { _ = os.RemoveAll(lockDir) }, nil
}

func closeLoopPlanCompletionCheck(root string) Check {
	active := filepath.Join(root, "plans", "active")
	entries, err := os.ReadDir(active)
	if err != nil {
		return Check{Name: "plan_closure", Status: "ok", Message: "no active plans directory"}
	}
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".md" {
			continue
		}
		path := filepath.Join(active, entry.Name())
		bytes, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		total, done := markdownTaskCounts(string(bytes))
		if total > 0 && total == done {
			if _, err := os.Stat(filepath.Join(root, "plans", "archived", entry.Name())); err != nil {
				return Check{Name: "plan_closure", Status: "blocked", Message: "completed active plan must be archived or marked exception: " + entry.Name()}
			}
		}
	}
	return Check{Name: "plan_closure", Status: "ok"}
}

func markdownTaskCounts(text string) (int, int) {
	total, done := 0, 0
	for _, line := range strings.Split(text, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "- [") || strings.HasPrefix(trimmed, "* [") {
			total++
			if strings.HasPrefix(trimmed, "- [x]") || strings.HasPrefix(trimmed, "- [X]") || strings.HasPrefix(trimmed, "* [x]") || strings.HasPrefix(trimmed, "* [X]") || strings.Contains(trimmed, "✅") {
				done++
			}
		}
		if strings.Contains(trimmed, ". ✅") {
			total++
			done++
		} else if strings.Contains(trimmed, ". ⏳") || strings.Contains(trimmed, ". [ ]") {
			total++
		}
	}
	return total, done
}

func closeLoopRuntimeValidation(root string) Check {
	if _, err := os.Stat(filepath.Join(root, "runtime")); err != nil {
		return Check{Name: "runtime_validation", Status: "skipped", Message: "no runtime directory in target repo"}
	}
	var stdout strings.Builder
	var stderr strings.Builder
	code := Run([]string{"runtime", "validate", "--repo", root, "--json"}, &stdout, &stderr)
	if code != ExitSuccess {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = compactSummary(stdout.String())
		}
		return Check{Name: "runtime_validation", Status: "failed", Message: msg}
	}
	return Check{Name: "runtime_validation", Status: "ok", Message: "ai-skill runtime validate passed"}
}

func closeLoopScanPrivatePaths(root string, paths []string) error {
	patterns := []string{"/Users/", "Authorization: Bearer ", "x-api-key:"}
	for _, rel := range paths {
		output, _ := exec.Command("git", "-C", root, "diff", "--", rel).Output()
		content := string(output)
		if _, err := os.Stat(filepath.Join(root, rel)); err == nil {
			if bytes, err := os.ReadFile(filepath.Join(root, rel)); err == nil {
				content += "\n" + string(bytes)
			}
		}
		for _, pattern := range patterns {
			if strings.Contains(content, pattern) {
				return fmt.Errorf("sensitive-looking content detected in %s", rel)
			}
		}
	}
	return nil
}

func closeLoopCommitMessageForGroup(group string) string {
	switch group {
	case "shared":
		return "docs(shared): close knowledge update loop"
	case "scripts":
		return "chore(scripts): update Go automation"
	case "tooling":
		return "docs(tools): update tool integration guidance"
	case "architecture":
		return "docs(architecture): close governance updates"
	default:
		return "docs(ai): close knowledge updates"
	}
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
