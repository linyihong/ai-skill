package app

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// validateRuntimeIndexFreshness is a commit-msg validator that enforces the
// invariant: staged markdown files that are present in the runtime index source
// inventory agree (by SHA-256 checksum) with the rows recorded in the staged
// knowledge/runtime/sqlite/runtime-index.sqlite. The runtime index builder is
// the canonical coverage boundary; this validator must not infer broader
// coverage from sibling directories.
//
// Opt-out: standalone "[skip-runtime-index-freshness]" trailer.
//
// Trigger rules (intentionally narrow to avoid false positives):
//
//   - Skip when the only staged file is
//     knowledge/runtime/sqlite/runtime-index.sqlite itself (an
//     index-only commit is a drift-repair commit; no markdown to verify).
//   - Skip when no staged path is a markdown file.
//
// Content semantics: comparisons use the STAGED blob (`git show :<path>`),
// not the worktree. This is the same snapshot that will land in the commit,
// so partial-stage scenarios (`git add -p`) do not produce false positives
// and post-stage worktree edits do not cause spurious blocks. See the
// ContentResolver abstraction below — runtime validate keeps its
// worktree-direct semantics; commit-time enforcement explicitly opts into
// staged.
func validateRuntimeIndexFreshness(text string, staged []string, root string) string {
	for _, line := range strings.Split(text, "\n") {
		if strings.TrimSpace(line) == "[skip-runtime-index-freshness]" {
			return ""
		}
	}
	if len(staged) == 0 {
		return ""
	}

	// Trigger gate: must touch at least one .md, and not just the index file.
	hasMarkdown := false
	for _, s := range staged {
		if strings.HasSuffix(strings.ToLower(s), ".md") {
			hasMarkdown = true
			break
		}
	}
	if !hasMarkdown {
		return ""
	}

	indexRel := filepath.ToSlash(filepath.Join("knowledge", "runtime", "sqlite", "runtime-index.sqlite"))
	indexStaged := false
	for _, s := range staged {
		if filepath.ToSlash(s) == indexRel {
			indexStaged = true
			break
		}
	}

	// Resolve the index DB path: if staged, materialize the staged blob to a
	// tempfile; otherwise read from disk. We never modify the file.
	dbPath, cleanup, err := resolveStagedIndexPath(root, indexRel, indexStaged)
	if err != nil {
		return fmt.Sprintf("runtime-index-freshness: cannot resolve staged runtime-index.sqlite: %v", err)
	}
	if cleanup != nil {
		defer cleanup()
	}
	if dbPath == "" {
		// Index file does not exist on disk and is not staged. Treat as "no
		// tracked source set known" → cannot enforce; skip.
		return ""
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Sprintf("runtime-index-freshness: cannot open runtime-index.sqlite: %v", err)
	}
	defer db.Close()

	// Build the source inventory defined by runtime-index.sqlite.
	rows, err := db.Query("SELECT source_path, checksum FROM sources WHERE checksum IS NOT NULL AND checksum != ''")
	if err != nil {
		// sources table may not exist in degenerate fixtures; skip.
		return ""
	}
	defer rows.Close()
	indexed := map[string]string{}
	for rows.Next() {
		var sp, ck string
		if err := rows.Scan(&sp, &ck); err != nil {
			return fmt.Sprintf("runtime-index-freshness: scan sources row failed: %v", err)
		}
		sp = filepath.ToSlash(sp)
		indexed[sp] = ck
	}
	if err := rows.Err(); err != nil {
		return fmt.Sprintf("runtime-index-freshness: iterate sources failed: %v", err)
	}

	resolver := newStagedBlobResolver(root)

	var staleViolations []string

	// Deterministic ordering for stable error messages.
	stagedSorted := append([]string(nil), staged...)
	sort.Strings(stagedSorted)

	for _, s := range stagedSorted {
		rel := filepath.ToSlash(s)
		if !strings.HasSuffix(strings.ToLower(rel), ".md") {
			continue
		}
		// Skip files that were deleted in the staged tree (no blob to read).
		// gitShowStagedBlob will error with non-zero exit on a delete.
		content, err := resolver.Read(rel)
		if err != nil {
			// Likely a staged delete; the source row should also be removed
			// in the same commit. If a row remains for a deleted source,
			// nativeRuntimeIndexSourceReferencesCheck will surface it at
			// runtime validate time. Don't double-block here.
			continue
		}

		expected, tracked := indexed[rel]
		if tracked {
			actual := fmt.Sprintf("%x", sha256.Sum256(content))
			if actual != expected {
				staleViolations = append(staleViolations, rel)
			}
		}
	}

	if len(staleViolations) == 0 {
		return ""
	}

	var parts []string
	if len(staleViolations) > 0 {
		parts = append(parts, "stale checksum(s) — staged markdown content disagrees with runtime-index.sqlite:\n    - "+
			strings.Join(staleViolations, "\n    - "))
	}

	return "runtime-index-freshness:\n  " + strings.Join(parts, "\n  ") +
		"\n  Remediation: run `ai-skill runtime refresh`, stage the regenerated knowledge/runtime/sqlite/runtime-index.sqlite, and retry the commit." +
		"\n  Opt-out: add a standalone `[skip-runtime-index-freshness]` trailer for genuine refresh-only or test-fixture commits."
}

// ContentResolver abstracts reading file content. Different commit-validation
// vs runtime-validation contexts need different semantics:
//
//   - WorktreeResolver: reads from the on-disk worktree. Used by
//     `runtime validate` to confirm the current working state is consistent.
//   - StagedBlobResolver: reads the staged blob via `git show :<path>`. Used
//     by the commit-msg validator to confirm the snapshot that will land in
//     the commit is consistent. Resists `git add -p` and post-stage edits.
type ContentResolver interface {
	Read(relPath string) ([]byte, error)
}

type worktreeResolver struct {
	root string
}

func newWorktreeResolver(root string) ContentResolver {
	return &worktreeResolver{root: root}
}

func (r *worktreeResolver) Read(relPath string) ([]byte, error) {
	return os.ReadFile(filepath.Join(r.root, filepath.FromSlash(relPath)))
}

type stagedBlobResolver struct {
	root string
}

func newStagedBlobResolver(root string) ContentResolver {
	return &stagedBlobResolver{root: root}
}

func (r *stagedBlobResolver) Read(relPath string) ([]byte, error) {
	cmd := exec.Command("git", "-C", r.root, "show", ":"+filepath.ToSlash(relPath))
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return out, nil
}

// resolveStagedIndexPath returns a filesystem path the caller can `sql.Open`
// against. If the runtime-index.sqlite blob is staged, we materialize it via
// `git show :<path>` to a tempfile so the staged snapshot is what we
// validate against. Otherwise we point at the on-disk copy. Returns
// ("", nil, nil) when neither source is available.
func resolveStagedIndexPath(root string, indexRel string, indexStaged bool) (string, func(), error) {
	if indexStaged {
		cmd := exec.Command("git", "-C", root, "show", ":"+indexRel)
		out, err := cmd.Output()
		if err != nil {
			return "", nil, fmt.Errorf("git show :%s: %w", indexRel, err)
		}
		tmp, err := os.CreateTemp("", "runtime-index-staged-*.sqlite")
		if err != nil {
			return "", nil, err
		}
		if _, err := tmp.Write(out); err != nil {
			tmp.Close()
			os.Remove(tmp.Name())
			return "", nil, err
		}
		if err := tmp.Close(); err != nil {
			os.Remove(tmp.Name())
			return "", nil, err
		}
		cleanup := func() { os.Remove(tmp.Name()) }
		return tmp.Name(), cleanup, nil
	}
	abs := filepath.Join(root, filepath.FromSlash(indexRel))
	if _, err := os.Stat(abs); err != nil {
		return "", nil, nil
	}
	return abs, nil, nil
}
