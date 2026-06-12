package app

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// Integration tests for validatePlanArchivalLinkIntegrity. Each test
// builds a real git repo in t.TempDir(), stages an archive event, and
// invokes the validator end-to-end (parser + rename detection + resolver
// + textual scan + formatter).
//
// Plan: plans/archived/2026-06-11-1100-plan-archival-link-integrity.md
// Phase: 2 (Tests — integration / fixture-based).
//
// TD-1 staged-vs-worktree divergence fixture is intentionally NOT here;
// it lives in a separate file so the Resolution Gate's evidence run is
// isolated.

func setupArchivalRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	runGitFixture(t, dir, "init", "-q", "-b", "main")
	runGitFixture(t, dir, "config", "user.email", "test@example.com")
	runGitFixture(t, dir, "config", "user.name", "Test User")
	runGitFixture(t, dir, "config", "commit.gpgsign", "false")
	return dir
}

func runGitFixture(t *testing.T, root string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, out)
	}
}

// archivePlan stages a rename of oldRel -> newRel by doing a real file
// move + `git add -A`, so the staged diff carries the rename even on
// platforms where `git mv` does not auto-create the destination dir.
func archivePlan(t *testing.T, root, oldRel, newRel string) {
	t.Helper()
	oldFull := filepath.Join(root, filepath.FromSlash(oldRel))
	newFull := filepath.Join(root, filepath.FromSlash(newRel))
	if err := os.MkdirAll(filepath.Dir(newFull), 0o755); err != nil {
		t.Fatalf("mkdir archive dir: %v", err)
	}
	if err := os.Rename(oldFull, newFull); err != nil {
		t.Fatalf("rename %s -> %s: %v", oldRel, newRel, err)
	}
	runGitFixture(t, root, "add", "-A")
}

func writeRepoFile(t *testing.T, root, relPath, content string) {
	t.Helper()
	full := filepath.Join(root, filepath.FromSlash(relPath))
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", relPath, err)
	}
}

func TestValidatePlanArchivalLinkIntegrity_OutboundBroken(t *testing.T) {
	// foo.md uses a bare relative link `./sibling.md` that resolves
	// against plans/active/ at write time. After archive, the same link
	// resolves against plans/archived/sibling.md which does not exist —
	// the move breaks the link without anyone editing it.
	root := setupArchivalRepo(t)
	writeRepoFile(t, root, "plans/active/sibling.md", "# sibling\n")
	writeRepoFile(t, root, "plans/active/foo.md", "# foo\n\nSee [parent](./sibling.md) for context.\n")
	runGitFixture(t, root, "add", "plans/active/sibling.md", "plans/active/foo.md")
	runGitFixture(t, root, "commit", "-qm", "init", "--no-verify")
	archivePlan(t, root, "plans/active/foo.md", "plans/archived/foo.md")

	got := validatePlanArchivalLinkIntegrity("", nil, root)
	if got == "" {
		t.Fatalf("expected block finding for broken outbound link; got empty")
	}
	if !strings.Contains(got, "broken_outbound_link") {
		t.Errorf("expected broken_outbound_link category; got: %s", got)
	}
	if !strings.Contains(got, "plans/archived/foo.md") {
		t.Errorf("expected file path of moved plan; got: %s", got)
	}
}

func TestValidatePlanArchivalLinkIntegrity_InboundBroken(t *testing.T) {
	root := setupArchivalRepo(t)
	writeRepoFile(t, root, "plans/active/foo.md", "# foo\n")
	writeRepoFile(t, root, "docs/index.md", "# index\n\nSee [foo](../plans/active/foo.md) for the plan.\n")
	runGitFixture(t, root, "add", ".")
	runGitFixture(t, root, "commit", "-qm", "init", "--no-verify")
	archivePlan(t, root, "plans/active/foo.md", "plans/archived/foo.md")

	got := validatePlanArchivalLinkIntegrity("", nil, root)
	if got == "" {
		t.Fatalf("expected block finding for broken inbound link; got empty")
	}
	if !strings.Contains(got, "broken_inbound_link") {
		t.Errorf("expected broken_inbound_link category; got: %s", got)
	}
	if !strings.Contains(got, "docs/index.md") {
		t.Errorf("expected referrer path in output; got: %s", got)
	}
	if !strings.Contains(got, `suggested:`) {
		t.Errorf("expected suggested replacement in output; got: %s", got)
	}
}

func TestValidatePlanArchivalLinkIntegrity_TextualWarning(t *testing.T) {
	// Acceptance Contract: a stale bare textual mention is WARNING severity —
	// it must NOT block the commit. The block validator returns "" (so the
	// dispatcher does not block); the advisory surface returns the warning.
	root := setupArchivalRepo(t)
	writeRepoFile(t, root, "plans/active/foo.md", "# foo\n")
	writeRepoFile(t, root, "docs/notes.md", "Note: see plans/active/foo.md when you need context.\n")
	runGitFixture(t, root, "add", ".")
	runGitFixture(t, root, "commit", "-qm", "init", "--no-verify")
	archivePlan(t, root, "plans/active/foo.md", "plans/archived/foo.md")

	if got := validatePlanArchivalLinkIntegrity("", nil, root); got != "" {
		t.Fatalf("warning-only archive must NOT block (validator must return empty); got: %s", got)
	}
	warn := warnPlanArchivalLinkIntegrity("", nil, root)
	if warn == "" {
		t.Fatalf("expected advisory warning for bare textual mention; got empty")
	}
	if !strings.Contains(warn, "stale_textual_reference") {
		t.Errorf("expected stale_textual_reference category; got: %s", warn)
	}
	if !strings.Contains(warn, "advisory") {
		t.Errorf("expected advisory framing; got: %s", warn)
	}
	if strings.Contains(warn, "historical_provenance_reference") {
		t.Errorf("did not expect info category in output; got: %s", warn)
	}
}

// The three dispatcher-semantics E2E cases requested in the Finding-1 review:
// the gate blocks iff there is a block-severity (broken-link) finding;
// warning-only archives pass while still surfacing the advisory.

func TestPlanArchivalLinkIntegrity_OnlyBlock_BlocksNoWarn(t *testing.T) {
	root := setupArchivalRepo(t)
	writeRepoFile(t, root, "plans/active/sibling.md", "# sibling\n")
	writeRepoFile(t, root, "plans/active/foo.md", "# foo\n\nSee [parent](./sibling.md).\n")
	runGitFixture(t, root, "add", ".")
	runGitFixture(t, root, "commit", "-qm", "init", "--no-verify")
	archivePlan(t, root, "plans/active/foo.md", "plans/archived/foo.md")

	if got := validatePlanArchivalLinkIntegrity("", nil, root); got == "" {
		t.Fatalf("broken outbound link must block (non-empty validator return)")
	}
	if warn := warnPlanArchivalLinkIntegrity("", nil, root); warn != "" {
		t.Errorf("no textual reference present, warn channel must be empty; got: %s", warn)
	}
}

func TestPlanArchivalLinkIntegrity_OnlyWarning_AllowedSurfaced(t *testing.T) {
	root := setupArchivalRepo(t)
	writeRepoFile(t, root, "plans/active/foo.md", "# foo\n")
	writeRepoFile(t, root, "docs/notes.md", "Background: plans/active/foo.md described it.\n")
	runGitFixture(t, root, "add", ".")
	runGitFixture(t, root, "commit", "-qm", "init", "--no-verify")
	archivePlan(t, root, "plans/active/foo.md", "plans/archived/foo.md")

	if got := validatePlanArchivalLinkIntegrity("", nil, root); got != "" {
		t.Fatalf("warning-only archive must be allowed (validator must return empty); got: %s", got)
	}
	if warn := warnPlanArchivalLinkIntegrity("", nil, root); warn == "" {
		t.Errorf("advisory must still surface the stale textual reference")
	}
}

func TestPlanArchivalLinkIntegrity_BlockPlusWarning_BlocksAndWarns(t *testing.T) {
	root := setupArchivalRepo(t)
	writeRepoFile(t, root, "plans/active/sibling.md", "# sibling\n")
	writeRepoFile(t, root, "plans/active/foo.md", "# foo\n\nSee [parent](./sibling.md).\n")
	writeRepoFile(t, root, "docs/notes.md", "Background: plans/active/foo.md described it.\n")
	runGitFixture(t, root, "add", ".")
	runGitFixture(t, root, "commit", "-qm", "init", "--no-verify")
	archivePlan(t, root, "plans/active/foo.md", "plans/archived/foo.md")

	if got := validatePlanArchivalLinkIntegrity("", nil, root); got == "" {
		t.Fatalf("broken link present, validator must block")
	}
	if warn := warnPlanArchivalLinkIntegrity("", nil, root); warn == "" {
		t.Errorf("textual reference present, advisory must also surface")
	}
}

func TestValidatePlanArchivalLinkIntegrity_TextualWithProvenanceSuppressed(t *testing.T) {
	root := setupArchivalRepo(t)
	writeRepoFile(t, root, "plans/active/foo.md", "# foo\n")
	writeRepoFile(t, root, "docs/notes.md",
		"<!-- archival-provenance -->\nOriginally lived at plans/active/foo.md.\n")
	runGitFixture(t, root, "add", ".")
	runGitFixture(t, root, "commit", "-qm", "init", "--no-verify")
	archivePlan(t, root, "plans/active/foo.md", "plans/archived/foo.md")

	got := validatePlanArchivalLinkIntegrity("", nil, root)
	if got != "" {
		t.Errorf("expected empty output when only provenance-marked info finding; got: %s", got)
	}
}

func TestValidatePlanArchivalLinkIntegrity_CleanArchivePass(t *testing.T) {
	root := setupArchivalRepo(t)
	writeRepoFile(t, root, "plans/active/foo.md", "# foo\n\nSee [parent](sibling.md).\n")
	writeRepoFile(t, root, "plans/active/sibling.md", "# sibling\n")
	runGitFixture(t, root, "add", ".")
	runGitFixture(t, root, "commit", "-qm", "init", "--no-verify")
	archivePlan(t, root, "plans/active/foo.md", "plans/archived/foo.md")
	writeRepoFile(t, root, "plans/archived/foo.md", "# foo\n\nSee [parent](../active/sibling.md).\n")
	runGitFixture(t, root, "add", "plans/archived/foo.md")

	got := validatePlanArchivalLinkIntegrity("", nil, root)
	if got != "" {
		t.Errorf("expected clean pass; got: %s", got)
	}
}

func TestValidatePlanArchivalLinkIntegrity_BareIdProvenanceNoFalsePositive(t *testing.T) {
	root := setupArchivalRepo(t)
	writeRepoFile(t, root, "plans/active/foo.md", "# foo\n")
	writeRepoFile(t, root, "docs/notes.md", "Historical note: project foo (id `foo`) shipped last quarter.\n")
	runGitFixture(t, root, "add", ".")
	runGitFixture(t, root, "commit", "-qm", "init", "--no-verify")
	archivePlan(t, root, "plans/active/foo.md", "plans/archived/foo.md")

	got := validatePlanArchivalLinkIntegrity("", nil, root)
	if got != "" {
		t.Errorf("bare-id mention without path should not false-positive; got: %s", got)
	}
}

func TestValidatePlanArchivalLinkIntegrity_EscapedParensPath(t *testing.T) {
	root := setupArchivalRepo(t)
	writeRepoFile(t, root, "plans/active/sib(ling).md", "# sib(ling)\n")
	writeRepoFile(t, root, "plans/active/foo.md",
		"# foo\n\nSee [parent](sib\\(ling\\).md) for context.\n")
	runGitFixture(t, root, "add", ".")
	runGitFixture(t, root, "commit", "-qm", "init", "--no-verify")
	archivePlan(t, root, "plans/active/foo.md", "plans/archived/foo.md")
	writeRepoFile(t, root, "plans/archived/foo.md",
		"# foo\n\nSee [parent](../active/sib\\(ling\\).md) for context.\n")
	runGitFixture(t, root, "add", "plans/archived/foo.md")

	got := validatePlanArchivalLinkIntegrity("", nil, root)
	if got != "" {
		t.Errorf("escaped parens path should resolve correctly; got: %s", got)
	}
}

func TestValidatePlanArchivalLinkIntegrity_MultiArchiveCrossRefBroken(t *testing.T) {
	// A originally used the (deliberately explicit) link
	// `[b](../active/b.md)` while in active/. After A and B both archive,
	// A's content is unchanged (so the rename is detected at 100%
	// similarity) but the link still points at active/b.md, which no
	// longer exists. The batch rename map should let the suggestion
	// land as `b.md` (same-dir archived sibling), even though B is
	// being archived in the *same commit*.
	root := setupArchivalRepo(t)
	writeRepoFile(t, root, "plans/active/a.md", "# a\n\nLinks to [b](../active/b.md).\n")
	writeRepoFile(t, root, "plans/active/b.md", "# b\n")
	runGitFixture(t, root, "add", ".")
	runGitFixture(t, root, "commit", "-qm", "init", "--no-verify")
	archivePlan(t, root, "plans/active/a.md", "plans/archived/a.md")
	archivePlan(t, root, "plans/active/b.md", "plans/archived/b.md")

	got := validatePlanArchivalLinkIntegrity("", nil, root)
	if got == "" {
		t.Fatalf("expected block finding for cross-archive broken link; got empty")
	}
	if !strings.Contains(got, "broken_outbound_link") {
		t.Errorf("expected broken_outbound_link category; got: %s", got)
	}
	if !strings.Contains(got, `suggested: "b.md"`) {
		t.Errorf("expected suggested rewrite to b.md from same-dir archive; got: %s", got)
	}
}

func TestValidatePlanArchivalLinkIntegrity_MultiArchiveCrossRefResolved(t *testing.T) {
	root := setupArchivalRepo(t)
	writeRepoFile(t, root, "plans/active/a.md", "# a\n\nLinks to [b](b.md).\n")
	writeRepoFile(t, root, "plans/active/b.md", "# b\n")
	runGitFixture(t, root, "add", ".")
	runGitFixture(t, root, "commit", "-qm", "init", "--no-verify")
	archivePlan(t, root, "plans/active/a.md", "plans/archived/a.md")
	archivePlan(t, root, "plans/active/b.md", "plans/archived/b.md")
	// A's link `[b](b.md)` already resolves correctly because B is now in
	// the same archived/ directory. No retarget needed.

	got := validatePlanArchivalLinkIntegrity("", nil, root)
	if got != "" {
		t.Errorf("expected clean pass for cross-archive same-dir link; got: %s", got)
	}
}

func TestValidatePlanArchivalLinkIntegrity_OptOutTrailer(t *testing.T) {
	root := setupArchivalRepo(t)
	writeRepoFile(t, root, "plans/active/foo.md", "# foo\n\nSee [x](../active/missing.md).\n")
	runGitFixture(t, root, "add", ".")
	runGitFixture(t, root, "commit", "-qm", "init", "--no-verify")
	archivePlan(t, root, "plans/active/foo.md", "plans/archived/foo.md")

	commitMsg := "emergency archive\n\n[skip-plan-archival-link-integrity]\n"
	got := validatePlanArchivalLinkIntegrity(commitMsg, nil, root)
	if got != "" {
		t.Errorf("opt-out trailer should suppress findings; got: %s", got)
	}
}

func TestValidatePlanArchivalLinkIntegrity_NoArchiveNoOp(t *testing.T) {
	root := setupArchivalRepo(t)
	writeRepoFile(t, root, "docs/x.md", "# x\n")
	runGitFixture(t, root, "add", ".")
	runGitFixture(t, root, "commit", "-qm", "init", "--no-verify")
	writeRepoFile(t, root, "docs/x.md", "# x updated\n")
	runGitFixture(t, root, "add", "docs/x.md")

	got := validatePlanArchivalLinkIntegrity("", nil, root)
	if got != "" {
		t.Errorf("non-archive commit should be no-op; got: %s", got)
	}
}
