package app

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

// planRename is one staged archive event detected via
// `git diff --cached --find-renames=90% --name-status`.
type planRename struct {
	OldPath string
	NewPath string
}

// linkFinding is one broken or stale reference discovered during a plan
// archive event. Severity / Category drive output formatting; the
// SuggestedReplacement carries the rewrite hint when known.
type linkFinding struct {
	Severity             string
	Category             string
	File                 string
	Line                 int
	Column               int
	Target               string
	SuggestedReplacement string
}

// validatePlanArchivalLinkIntegrity emits a block finding when a staged
// commit moves one or more plan files between plans/active/ and
// plans/archived/ and any markdown link — in the moved file (outbound)
// or any other repo .md file (inbound) — would be broken by the move.
//
// Opt-out: standalone "[skip-plan-archival-link-integrity]" trailer.
//
// Trigger: `git diff --cached --find-renames=90% --name-status` shows
// at least one R* rename between plans/active/ and plans/archived/.
//
// Multi-archive in same commit is handled by building the full rename
// map before any resolution; cross-references between simultaneously
// archived plans resolve correctly.
//
// Bare textual references (non-link prose mentions of old paths) are
// deferred to a follow-up step and not reported here.
//
// Plan: plans/active/2026-06-11-1100-plan-archival-link-integrity.md
// Phase: 1 (Implementation — outbound + inbound markdown link checks).
func validatePlanArchivalLinkIntegrity(text string, staged []string, root string) string {
	for _, line := range strings.Split(text, "\n") {
		if strings.TrimSpace(line) == "[skip-plan-archival-link-integrity]" {
			return ""
		}
	}

	renames := detectPlanArchivalRenames(root)
	if len(renames) == 0 {
		return ""
	}

	renameMap := make(map[string]string, len(renames))
	movedNew := make(map[string]bool, len(renames))
	for _, r := range renames {
		renameMap[r.OldPath] = r.NewPath
		movedNew[r.NewPath] = true
	}

	var findings []linkFinding

	for _, r := range renames {
		content, err := readStagedFileContent(root, r.NewPath)
		if err != nil {
			continue
		}
		for _, link := range extractMarkdownLinks(content) {
			f, ok := classifyLink(r.NewPath, link, renameMap, root, "broken_outbound_link")
			if ok {
				findings = append(findings, f)
			}
		}
	}

	for _, mdPath := range listRepoMarkdown(root) {
		if movedNew[mdPath] {
			continue
		}
		content, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(mdPath)))
		if err != nil {
			continue
		}
		for _, link := range extractMarkdownLinks(content) {
			f, ok := classifyLink(mdPath, link, renameMap, root, "broken_inbound_link")
			if ok && f.SuggestedReplacement != "" {
				findings = append(findings, f)
			}
		}
	}

	if len(findings) == 0 {
		return ""
	}
	return formatLinkFindings(findings)
}

// classifyLink resolves a link target against the current rename map
// and disk state. Returns (finding, true) when the link is broken;
// otherwise (zero, false). The "kind" arg distinguishes outbound vs
// inbound only for the Category label.
func classifyLink(fromFile string, link Link, renameMap map[string]string, root, kind string) (linkFinding, bool) {
	target := stripLinkFragment(link.Target)
	if target == "" {
		return linkFinding{}, false
	}
	resolved := resolveRepoPath(fromFile, target)
	if resolved == "" {
		return linkFinding{}, false
	}
	if newPath, ok := renameMap[resolved]; ok {
		return linkFinding{
			Severity:             "block",
			Category:             kind,
			File:                 fromFile,
			Line:                 link.Line,
			Column:               link.Column,
			Target:               link.Target,
			SuggestedReplacement: suggestReplacement(fromFile, newPath, link.Target),
		}, true
	}
	if pathExistsInRepo(root, resolved) {
		return linkFinding{}, false
	}
	return linkFinding{
		Severity: "block",
		Category: kind,
		File:     fromFile,
		Line:     link.Line,
		Column:   link.Column,
		Target:   link.Target,
	}, true
}

// detectPlanArchivalRenames runs `git diff --cached --find-renames=90%
// --name-status` and returns plan archive rename events
// (plans/active/ ↔ plans/archived/, either direction).
func detectPlanArchivalRenames(root string) []planRename {
	cmd := exec.Command("git", "-C", root, "diff", "--cached", "--find-renames=90%", "--name-status")
	out, err := cmd.Output()
	if err != nil {
		return nil
	}
	return parsePlanRenames(string(out))
}

// parsePlanRenames is split out so it can be unit-tested without a git
// repo. Each line is the `--name-status` tab-separated form for renames:
//
//	R100\tplans/active/foo.md\tplans/archived/foo.md
func parsePlanRenames(diffOut string) []planRename {
	var renames []planRename
	for _, line := range strings.Split(diffOut, "\n") {
		if !strings.HasPrefix(line, "R") {
			continue
		}
		fields := strings.Split(line, "\t")
		if len(fields) < 3 {
			continue
		}
		oldPath := strings.TrimSpace(fields[1])
		newPath := strings.TrimSpace(fields[2])
		if isPlanArchiveMove(oldPath, newPath) {
			renames = append(renames, planRename{OldPath: oldPath, NewPath: newPath})
		}
	}
	return renames
}

func isPlanArchiveMove(oldPath, newPath string) bool {
	active := "plans/active/"
	archived := "plans/archived/"
	return (strings.HasPrefix(oldPath, active) && strings.HasPrefix(newPath, archived)) ||
		(strings.HasPrefix(oldPath, archived) && strings.HasPrefix(newPath, active))
}

func readStagedFileContent(root, repoPath string) ([]byte, error) {
	cmd := exec.Command("git", "-C", root, "show", ":"+repoPath)
	return cmd.Output()
}

func stripLinkFragment(target string) string {
	if idx := strings.Index(target, "#"); idx >= 0 {
		return target[:idx]
	}
	return target
}

// resolveRepoPath joins a markdown link target against the source file's
// directory and cleans the result as a POSIX-style repo path. Empty
// targets and targets that resolve outside the repo (leading "..") are
// not handled specially here; the existence check filters them out.
func resolveRepoPath(fromFile, linkTarget string) string {
	if linkTarget == "" {
		return ""
	}
	dir := path.Dir(fromFile)
	return path.Clean(path.Join(dir, linkTarget))
}

func pathExistsInRepo(root, repoPath string) bool {
	_, err := os.Stat(filepath.Join(root, filepath.FromSlash(repoPath)))
	return err == nil
}

func listRepoMarkdown(root string) []string {
	cmd := exec.Command("git", "-C", root, "ls-files", "*.md")
	out, err := cmd.Output()
	if err != nil {
		return nil
	}
	var files []string
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line != "" {
			files = append(files, line)
		}
	}
	return files
}

// suggestReplacement computes a rewrite hint for a broken link. When the
// original target looks repo-rooted ("plans/active/..." or
// "plans/archived/..."), give the new repo-rooted path. Otherwise give a
// repo-relative path computed POSIX-style from the source file's
// directory.
func suggestReplacement(fromFile, newPath, originalTarget string) string {
	if (strings.HasPrefix(originalTarget, "plans/active/") || strings.HasPrefix(originalTarget, "plans/archived/")) &&
		!strings.HasPrefix(originalTarget, "./") &&
		!strings.HasPrefix(originalTarget, "../") {
		return newPath
	}
	return posixRel(path.Dir(fromFile), newPath)
}

// posixRel is a POSIX-path equivalent of filepath.Rel; produces a path
// using '/' separators regardless of host OS, so suggestions are stable
// on Windows agents.
func posixRel(fromDir, toPath string) string {
	if fromDir == "." || fromDir == "" {
		return toPath
	}
	fromParts := strings.Split(fromDir, "/")
	toParts := strings.Split(toPath, "/")
	common := 0
	for common < len(fromParts) && common < len(toParts) && fromParts[common] == toParts[common] {
		common++
	}
	var parts []string
	for i := common; i < len(fromParts); i++ {
		parts = append(parts, "..")
	}
	parts = append(parts, toParts[common:]...)
	if len(parts) == 0 {
		return "."
	}
	return strings.Join(parts, "/")
}

func formatLinkFindings(findings []linkFinding) string {
	var bullets []string
	for _, f := range findings {
		loc := fmt.Sprintf("%s:%d:%d", f.File, f.Line, f.Column)
		msg := fmt.Sprintf("%s [%s] target=%q", loc, f.Category, f.Target)
		if f.SuggestedReplacement != "" {
			msg += fmt.Sprintf("  → suggested: %q", f.SuggestedReplacement)
		}
		bullets = append(bullets, msg)
	}
	return "plan-archival-link-integrity: archive event leaves broken markdown link(s):\n    - " +
		strings.Join(bullets, "\n    - ") +
		"\n  Update each broken link to its suggested target, or add a standalone `[skip-plan-archival-link-integrity]` trailer for emergency archives."
}
