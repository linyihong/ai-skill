package app

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func mkParent(t *testing.T, p string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatalf("mkdir parent: %v", err)
	}
}

// --- helpers ---

// seedRuntimeIndex creates a minimal runtime-index.sqlite at the given path
// with a `sources` table populated from the supplied (path, content) map.
// checksums are computed automatically.
func seedRuntimeIndex(t *testing.T, dbPath string, sources map[string]string) {
	t.Helper()
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	defer db.Close()
	if _, err := db.Exec(`CREATE TABLE sources (source_path TEXT PRIMARY KEY, checksum TEXT)`); err != nil {
		t.Fatalf("create table: %v", err)
	}
	for p, content := range sources {
		ck := fmt.Sprintf("%x", sha256.Sum256([]byte(content)))
		if _, err := db.Exec("INSERT INTO sources (source_path, checksum) VALUES (?, ?)", filepath.ToSlash(p), ck); err != nil {
			t.Fatalf("insert: %v", err)
		}
	}
}

// stageAll runs `git add` for each path then verifies they are staged.
func stageAll(t *testing.T, repo string, paths ...string) {
	t.Helper()
	args := append([]string{"add", "--"}, paths...)
	runGit(t, repo, args...)
}

// editWorktree mutates the on-disk worktree AFTER staging, simulating a
// post-stage local edit. The staged blob is untouched.
func editWorktree(t *testing.T, repo string, rel string, content string) {
	t.Helper()
	writeFile(t, filepath.Join(repo, rel), content)
}

// --- trigger gate tests ---

func TestRuntimeIndexFreshnessOptOut(t *testing.T) {
	body := "chore: bypass freshness\n\n[skip-runtime-index-freshness]\n"
	got := validateRuntimeIndexFreshness(body, []string{"enforcement/rule.md"}, "")
	if got != "" {
		t.Errorf("opt-out should bypass validator, got %q", got)
	}
}

func TestRuntimeIndexFreshnessNoMarkdownNoTrigger(t *testing.T) {
	staged := []string{"scripts/ai-skill-cli/internal/app/hooks.go"}
	got := validateRuntimeIndexFreshness("feat: code only", staged, "")
	if got != "" {
		t.Errorf("non-md commit should skip, got %q", got)
	}
}

func TestRuntimeIndexFreshnessIndexOnlyCommitSkips(t *testing.T) {
	// Refresh-only commit (just the index file, no markdown) must not block.
	staged := []string{"knowledge/runtime/sqlite/runtime-index.sqlite"}
	got := validateRuntimeIndexFreshness("chore: refresh index", staged, "")
	if got != "" {
		t.Errorf("index-only commit should skip (no md staged), got %q", got)
	}
}

// --- fixture tests using a real git repo ---

func TestRuntimeIndexFreshnessStaleStagedBlobBlocked(t *testing.T) {
	repo := initTempGitRepo(t)

	// Tracked markdown with original content A.
	mdRel := "enforcement/foo.md"
	contentA := "# Foo\nbody A\n"
	writeFile(t, filepath.Join(repo, mdRel), contentA)
	stageAll(t, repo, mdRel)

	// Seed index with checksum of OLD content (so staged != indexed).
	indexAbs := filepath.Join(repo, "knowledge", "runtime", "sqlite", "runtime-index.sqlite")
	mkParent(t, indexAbs) // create parent dir
	seedRuntimeIndex(t, indexAbs, map[string]string{
		mdRel: "# Foo\nbody OLD\n", // checksum of OLD differs from staged contentA
	})

	body := "fix: update foo content"
	staged := []string{mdRel}
	got := validateRuntimeIndexFreshness(body, staged, repo)
	if got == "" {
		t.Fatalf("expected block for stale staged blob, got empty")
	}
	if !strings.Contains(got, "stale checksum") {
		t.Errorf("expected stale-checksum violation, got %q", got)
	}
	if !strings.Contains(got, mdRel) {
		t.Errorf("expected violation to cite %s, got %q", mdRel, got)
	}
}

func TestRuntimeIndexFreshnessFreshStagedBlobPasses(t *testing.T) {
	repo := initTempGitRepo(t)

	mdRel := "enforcement/foo.md"
	content := "# Foo\nfresh body\n"
	writeFile(t, filepath.Join(repo, mdRel), content)
	stageAll(t, repo, mdRel)

	indexAbs := filepath.Join(repo, "knowledge", "runtime", "sqlite", "runtime-index.sqlite")
	mkParent(t, indexAbs)
	seedRuntimeIndex(t, indexAbs, map[string]string{
		mdRel: content, // checksum matches staged
	})

	got := validateRuntimeIndexFreshness("fix: foo", []string{mdRel}, repo)
	if got != "" {
		t.Errorf("fresh staged blob should pass, got %q", got)
	}
}

// Critical regression: simulates `git add -p` style flow where the staged
// blob == index checksum but the worktree has post-stage edits.
// Worktree-based validation would false-positive; staged-blob validation
// must pass.
func TestRuntimeIndexFreshnessPartialStageNotFalsePositive(t *testing.T) {
	repo := initTempGitRepo(t)

	mdRel := "enforcement/foo.md"
	stagedContent := "# Foo\nstaged body\n"
	writeFile(t, filepath.Join(repo, mdRel), stagedContent)
	stageAll(t, repo, mdRel)

	// Post-stage worktree edit — staged blob is unaffected.
	editWorktree(t, repo, mdRel, "# Foo\nstaged body\nplus extra unstaged work\n")

	indexAbs := filepath.Join(repo, "knowledge", "runtime", "sqlite", "runtime-index.sqlite")
	mkParent(t, indexAbs)
	seedRuntimeIndex(t, indexAbs, map[string]string{
		mdRel: stagedContent, // matches STAGED blob, not worktree
	})

	got := validateRuntimeIndexFreshness("fix: partial stage", []string{mdRel}, repo)
	if got != "" {
		t.Errorf("staged-blob semantics must not be misled by post-stage worktree edit, got %q", got)
	}
}

func TestRuntimeIndexFreshnessUnindexedSiblingAllowed(t *testing.T) {
	repo := initTempGitRepo(t)

	existingRel := "enforcement/existing.md"
	newRel := "enforcement/new.md"
	writeFile(t, filepath.Join(repo, existingRel), "# Existing\n")
	writeFile(t, filepath.Join(repo, newRel), "# New\n")
	stageAll(t, repo, existingRel, newRel)

	indexAbs := filepath.Join(repo, "knowledge", "runtime", "sqlite", "runtime-index.sqlite")
	mkParent(t, indexAbs)
	// Only existing.md is in the source inventory. A sibling markdown without a
	// source row is outside the runtime index builder's canonical coverage.
	seedRuntimeIndex(t, indexAbs, map[string]string{
		existingRel: "# Existing\n",
	})

	got := validateRuntimeIndexFreshness("feat: add new rule", []string{existingRel, newRel}, repo)
	if got != "" {
		t.Errorf("unindexed sibling markdown should be outside freshness scope, got %q", got)
	}
}

func TestRuntimeIndexFreshnessMarkdownInUntrackedDirAllowed(t *testing.T) {
	repo := initTempGitRepo(t)

	// Markdown in a directory not covered by the index at all (e.g. plans/).
	mdRel := "plans/active/whatever.md"
	writeFile(t, filepath.Join(repo, mdRel), "# Plan\n")
	stageAll(t, repo, mdRel)

	indexAbs := filepath.Join(repo, "knowledge", "runtime", "sqlite", "runtime-index.sqlite")
	mkParent(t, indexAbs)
	// Index only tracks enforcement/; plans/ is not a tracked dir.
	seedRuntimeIndex(t, indexAbs, map[string]string{
		"enforcement/foo.md": "# Foo\n",
	})

	got := validateRuntimeIndexFreshness("docs: plan", []string{mdRel}, repo)
	if got != "" {
		t.Errorf("md in untracked dir should pass, got %q", got)
	}
}

func TestRuntimeIndexFreshnessIndexFileAbsentSkips(t *testing.T) {
	// No runtime-index.sqlite at all and not staged — validator cannot
	// enforce; skip to fail-open.
	repo := initTempGitRepo(t)
	mdRel := "enforcement/foo.md"
	writeFile(t, filepath.Join(repo, mdRel), "# Foo\n")
	stageAll(t, repo, mdRel)
	got := validateRuntimeIndexFreshness("docs: foo", []string{mdRel}, repo)
	if got != "" {
		t.Errorf("absent index should skip, got %q", got)
	}
}
