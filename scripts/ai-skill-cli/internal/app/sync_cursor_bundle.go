package app

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/linyihong/Ai-skill/scripts/ai-skill-cli/internal/pathutil"
)

type syncCursorBundleOptions struct {
	repoPath    string
	targetPath  string
	dryRun      bool
	jsonOutput  bool
	plainOutput bool
}

func runSyncCursorBundle(args []string, stdout io.Writer, stderr io.Writer) int {
	fs := newFlagSet("sync-cursor-bundle", stderr)
	opts := syncCursorBundleOptions{}
	fs.StringVar(&opts.repoPath, "repo", ".", "Ai-skill repository path")
	fs.StringVar(&opts.targetPath, "target", "", "Cursor root target path, such as <fake-home>/.cursor")
	fs.BoolVar(&opts.dryRun, "dry-run", false, "preview sync without writing")
	fs.BoolVar(&opts.jsonOutput, "json", false, "write machine-readable JSON output")
	fs.BoolVar(&opts.plainOutput, "plain", false, "write human-readable output")
	if err := fs.Parse(args); err != nil {
		return ExitInvalidUsage
	}
	if opts.jsonOutput && opts.plainOutput {
		_, _ = fmt.Fprintln(stderr, "--json and --plain are mutually exclusive")
		return ExitInvalidUsage
	}

	result := buildSyncCursorBundleResult(opts)
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

func buildSyncCursorBundleResult(opts syncCursorBundleOptions) Result {
	result := Result{
		Command:        "sync-cursor-bundle",
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
			Message:     "sync-cursor-bundle currently supports dry-run planning only.",
			Remediation: "Run with --dry-run until symlink/copy fallback parity is fixture-backed.",
		}
	}

	repo, repoCheck := resolveExistingDir("repo", opts.repoPath)
	result.Checks = append(result.Checks, repoCheck)
	if repoCheck.Status != "ok" {
		result.Status = "blocked"
		result.ExitCode = ExitInvalidUsage
		result.Error = &CommandError{Code: "invalid_repo", Message: repoCheck.Message, Remediation: "Pass --repo with the Ai-skill repository root."}
		return result
	}

	sourceCheck := syncSourceCheck(repo)
	result.Checks = append(result.Checks, sourceCheck)
	if sourceCheck.Status != "ok" {
		result.Status = "blocked"
		result.ExitCode = ExitValidationFailed
		result.Error = &CommandError{Code: "invalid_sync_source", Message: sourceCheck.Message, Remediation: "Run from a complete Ai-skill checkout with enforcement/ present."}
		return result
	}

	target, targetCheck := resolveTargetRoot(opts.targetPath)
	result.Checks = append(result.Checks, targetCheck)
	if targetCheck.Status != "ok" {
		result.Status = "blocked"
		result.ExitCode = ExitInvalidUsage
		result.Error = &CommandError{Code: "invalid_target", Message: targetCheck.Message, Remediation: "Pass --target with an explicit fake or real Cursor root outside the repository."}
		return result
	}

	if pathInside(target, repo) {
		result.Checks = append(result.Checks, Check{Name: "target_boundary", Status: "failed", Message: "target is inside repository"})
		result.Status = "blocked"
		result.ExitCode = ExitValidationFailed
		result.Error = &CommandError{Code: "target_inside_repo", Message: "bundle target must live outside the repository", Remediation: "Choose a target outside the Ai-skill checkout."}
		return result
	}
	result.Checks = append(result.Checks, Check{Name: "target_boundary", Status: "ok", Message: "target is outside repository"})

	strategy := syncStrategy()
	result.Checks = append(result.Checks, Check{Name: "mirror_strategy", Status: "ok", Message: strategy})
	result.PlannedActions = append(result.PlannedActions,
		"ensure directory: "+filepath.Join(target, "bundles"),
		"ensure directory: "+filepath.Join(target, "bundles", "ai-skill"),
		"ensure directory: "+filepath.Join(target, "skills"),
		"remove legacy mirror path if present: "+filepath.Join(target, "bundles", "ai-skill", "shared-rules"),
		fmt.Sprintf("%s mirror: %s -> %s", strategy, filepath.Join(repo, "enforcement"), filepath.Join(target, "bundles", "shared-rules")),
		fmt.Sprintf("%s mirror: %s -> %s", strategy, filepath.Join(target, "bundles", "shared-rules"), filepath.Join(target, "shared-rules")),
	)

	skills := syncSkillDirs(filepath.Join(repo, "skills"))
	result.Checks = append(result.Checks, Check{Name: "skills", Status: "ok", Message: fmt.Sprintf("%d syncable skills", len(skills))})
	for _, skill := range skills {
		source := filepath.Join(repo, "skills", skill)
		bundle := filepath.Join(target, "bundles", "ai-skill", skill)
		active := filepath.Join(target, "skills", skill)
		result.PlannedActions = append(result.PlannedActions,
			fmt.Sprintf("%s mirror: %s -> %s", strategy, source, bundle),
			fmt.Sprintf("%s mirror: %s -> %s", strategy, bundle, active),
		)
	}

	return result
}

func resolveExistingDir(name string, path string) (string, Check) {
	if strings.TrimSpace(path) == "" {
		return "", Check{Name: name, Status: "failed", Message: "--" + name + " is required"}
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", Check{Name: name, Status: "failed", Message: err.Error()}
	}
	info, err := os.Stat(abs)
	if err != nil {
		return "", Check{Name: name, Status: "failed", Message: err.Error()}
	}
	if !info.IsDir() {
		return "", Check{Name: name, Status: "failed", Message: "path is not a directory"}
	}
	normalized, err := pathutil.NormalizeForReport(abs)
	if err != nil {
		return "", Check{Name: name, Status: "failed", Message: err.Error()}
	}
	return abs, Check{Name: name, Status: "ok", Message: normalized}
}

func resolveTargetRoot(path string) (string, Check) {
	if strings.TrimSpace(path) == "" {
		return "", Check{Name: "target", Status: "failed", Message: "--target is required"}
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", Check{Name: "target", Status: "failed", Message: err.Error()}
	}
	normalized, err := pathutil.NormalizeForReport(abs)
	if err != nil {
		return "", Check{Name: "target", Status: "failed", Message: err.Error()}
	}
	return abs, Check{Name: "target", Status: "ok", Message: normalized}
}

func syncSourceCheck(repo string) Check {
	enforcement := filepath.Join(repo, "enforcement")
	info, err := os.Stat(enforcement)
	if err != nil {
		return Check{Name: "sync_source", Status: "failed", Message: "missing enforcement/ source"}
	}
	if !info.IsDir() {
		return Check{Name: "sync_source", Status: "failed", Message: "enforcement source is not a directory"}
	}
	return Check{Name: "sync_source", Status: "ok", Message: "enforcement source present"}
}

func syncSkillDirs(skillsRoot string) []string {
	entries, err := os.ReadDir(skillsRoot)
	if err != nil {
		return nil
	}
	skills := []string{}
	for _, entry := range entries {
		if !entry.IsDir() || entry.Name() == "_template" {
			continue
		}
		if _, err := os.Stat(filepath.Join(skillsRoot, entry.Name(), "SKILL.md")); err == nil {
			skills = append(skills, entry.Name())
		}
	}
	sort.Strings(skills)
	return skills
}

func pathInside(path string, parent string) bool {
	rel, err := filepath.Rel(parent, path)
	if err != nil {
		return false
	}
	return rel == "." || (!strings.HasPrefix(rel, ".."+string(filepath.Separator)) && rel != "..")
}

func syncStrategy() string {
	return "copy-fallback"
}
