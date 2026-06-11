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
// also reported. They default to a warning (category
// stale_textual_reference); when a same-line or preceding-line
// <!-- archival-provenance --> marker is present the finding is
// downgraded to info (category historical_provenance_reference) and
// suppressed from validator output.
//
// Plan: plans/archived/2026-06-11-1100-plan-archival-link-integrity.md
// Phase: 1 (Implementation — outbound + inbound + bare textual scan).
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

	mdFiles := listRepoMarkdown(root)
	inSet := make(map[string]bool, len(mdFiles))
	for _, m := range mdFiles {
		inSet[m] = true
	}
	for newPath := range movedNew {
		if !inSet[newPath] {
			mdFiles = append(mdFiles, newPath)
			inSet[newPath] = true
		}
	}

	var findings []linkFinding
	for _, mdPath := range mdFiles {
		content, err := readFileForScan(root, mdPath, movedNew)
		if err != nil {
			continue
		}
		isMoved := movedNew[mdPath]
		kind := "broken_inbound_link"
		if isMoved {
			kind = "broken_outbound_link"
		}
		for _, link := range extractMarkdownLinks(content) {
			f, ok := classifyLink(mdPath, link, renameMap, root, kind)
			if !ok {
				continue
			}
			if !isMoved && f.SuggestedReplacement == "" {
				continue
			}
			findings = append(findings, f)
		}
		findings = append(findings, scanBareTextualReferences(mdPath, content, renames)...)
	}

	return formatLinkFindings(findings)
}

// readFileForScan reads content for the inbound / textual scan. The
// staged blob is the canonical commit candidate, so we read via
// `git show :<path>` first; if the file is not in the index (untracked
// .md, or git not available) we fall back to the worktree.
//
// Resolves TD-1 (staged vs worktree drift) per the Resolution Gate
// recorded in plans/archived/2026-06-11-1100-plan-archival-link-integrity.md.
// The fixture run observed both directions of divergence (staged-has-fix
// and staged-broken-worktree-fixed); the validator now reports against
// what will actually be committed.
//
// movedNew is retained on the signature for callsite readability but no
// longer affects routing — the staged-first path covers renamed files
// uniformly with non-renamed staged-modified files.
func readFileForScan(root, mdPath string, _ map[string]bool) ([]byte, error) {
	if data, err := readStagedFileContent(root, mdPath); err == nil {
		return data, nil
	}
	return os.ReadFile(filepath.Join(root, filepath.FromSlash(mdPath)))
}

// archivalProvenanceMarker is the explicit opt-in line that downgrades
// a bare textual reference from warning to info. Authors intentionally
// keeping historical paths in prose write this on the same line or the
// line immediately above.
const archivalProvenanceMarker = "<!-- archival-provenance -->"

// scanBareTextualReferences finds plain-text mentions of any renamed
// plan's old path that are NOT inside a markdown link target (those
// are already covered by the inbound link check). Provenance-marked
// occurrences are categorised as info.
func scanBareTextualReferences(mdPath string, content []byte, renames []planRename) []linkFinding {
	s := string(content)
	hit := false
	for _, r := range renames {
		if strings.Contains(s, r.OldPath) {
			hit = true
			break
		}
	}
	if !hit {
		return nil
	}
	lines := strings.Split(s, "\n")
	var findings []linkFinding
	for lineIdx, line := range lines {
		for _, r := range renames {
			old := r.OldPath
			start := 0
			for {
				pos := strings.Index(line[start:], old)
				if pos < 0 {
					break
				}
				abs := start + pos
				start = abs + len(old)
				if isLinkTargetContext(line, abs) {
					continue
				}
				provenance := strings.Contains(line, archivalProvenanceMarker) ||
					(lineIdx > 0 && strings.Contains(lines[lineIdx-1], archivalProvenanceMarker))
				severity := "warning"
				category := "stale_textual_reference"
				if provenance {
					severity = "info"
					category = "historical_provenance_reference"
				}
				findings = append(findings, linkFinding{
					Severity:             severity,
					Category:             category,
					File:                 mdPath,
					Line:                 lineIdx + 1,
					Column:               abs + 1,
					Target:               old,
					SuggestedReplacement: r.NewPath,
				})
			}
		}
	}
	return findings
}

// isLinkTargetContext returns true when pos sits at the start of a
// markdown link target, i.e. the preceding non-space characters are
// `](`. Such hits belong to the markdown link path and are reported by
// the inbound scan instead of as bare textual references.
func isLinkTargetContext(line string, pos int) bool {
	i := pos - 1
	for i >= 0 && (line[i] == ' ' || line[i] == '\t') {
		i--
	}
	if i < 0 || line[i] != '(' {
		return false
	}
	i--
	if i < 0 || line[i] != ']' {
		return false
	}
	return true
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

// formatLinkFindings groups findings by severity and emits block + warning
// sections. Info findings are intentionally suppressed from output — the
// provenance marker exists precisely to silence them. Returns "" when
// only info findings are present.
func formatLinkFindings(findings []linkFinding) string {
	var blocks, warnings []linkFinding
	for _, f := range findings {
		switch f.Severity {
		case "block":
			blocks = append(blocks, f)
		case "warning":
			warnings = append(warnings, f)
		}
	}
	if len(blocks) == 0 && len(warnings) == 0 {
		return ""
	}
	var out strings.Builder
	out.WriteString("plan-archival-link-integrity: archive event leaves references requiring attention:")
	if len(blocks) > 0 {
		out.WriteString("\n  blocking:")
		for _, f := range blocks {
			out.WriteString("\n    - " + formatFindingLine(f))
		}
	}
	if len(warnings) > 0 {
		out.WriteString("\n  warnings:")
		for _, f := range warnings {
			out.WriteString("\n    - " + formatFindingLine(f))
		}
	}
	out.WriteString("\n  Update each reference, or add a standalone `[skip-plan-archival-link-integrity]` trailer for emergency archives.")
	return out.String()
}

func formatFindingLine(f linkFinding) string {
	loc := fmt.Sprintf("%s:%d:%d", f.File, f.Line, f.Column)
	msg := fmt.Sprintf("%s [%s] target=%q", loc, f.Category, f.Target)
	if f.SuggestedReplacement != "" {
		msg += fmt.Sprintf("  → suggested: %q", f.SuggestedReplacement)
	}
	return msg
}
