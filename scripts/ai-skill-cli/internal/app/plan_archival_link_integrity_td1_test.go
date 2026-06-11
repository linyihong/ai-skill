package app

import (
	"strings"
	"testing"
)

// TD-1 Resolution Gate fixture (plans/active/2026-06-11-1100-plan-archival-link-integrity.md).
//
// These tests exercise staged-vs-worktree divergence to verify the
// validator reports against the commit candidate (staged blob), not
// the worktree. Per the Resolution Gate procedure:
//
//	1. Run divergence fixture (these tests).
//	2. Observe behaviour.
//	3. If false-pass or false-block is observed -> promote TD-1 to
//	   active scope, insert Phase 2.5 (readFileForScan reads staged).
//	4. Otherwise keep as documented limitation.
//
// Gate outcome (2026-06-11): both directions of divergence were
// observed with the worktree-only implementation. TD-1 promoted; Phase
// 2.5 landed in the same commit (readFileForScan now reads staged
// blob first, worktree only as fallback). These tests assert the
// post-fix behaviour and must stay green.

func TestValidatePlanArchivalLinkIntegrity_TD1_StagedHasFixWorktreeBroken(t *testing.T) {
	// Staged blob has the fix; worktree shows the broken old link.
	// The validator MUST report clean — the commit will land with the
	// fixed link, regardless of what the worktree shows.
	root := setupArchivalRepo(t)
	writeRepoFile(t, root, "plans/active/foo.md", "# foo\n")
	writeRepoFile(t, root, "docs/index.md", "See [foo](../plans/active/foo.md).\n")
	runGitFixture(t, root, "add", ".")
	runGitFixture(t, root, "commit", "-qm", "init", "--no-verify")
	archivePlan(t, root, "plans/active/foo.md", "plans/archived/foo.md")

	writeRepoFile(t, root, "docs/index.md", "See [foo](../plans/archived/foo.md).\n")
	runGitFixture(t, root, "add", "docs/index.md")
	writeRepoFile(t, root, "docs/index.md", "See [foo](../plans/active/foo.md).\n")

	got := validatePlanArchivalLinkIntegrity("", nil, root)
	if got != "" {
		t.Errorf("staged blob has fix; expected clean pass (commit-candidate semantics); got finding:\n%s", got)
	}
}

func TestValidatePlanArchivalLinkIntegrity_TD1_StagedBrokenWorktreeFixed(t *testing.T) {
	// Reverse: staged blob still references the old path (no staged
	// edit on the referrer), but the worktree was hand-fixed. The
	// commit will land with the broken link, so the validator MUST
	// report it.
	root := setupArchivalRepo(t)
	writeRepoFile(t, root, "plans/active/foo.md", "# foo\n")
	writeRepoFile(t, root, "docs/index.md", "See [foo](../plans/active/foo.md).\n")
	runGitFixture(t, root, "add", ".")
	runGitFixture(t, root, "commit", "-qm", "init", "--no-verify")
	archivePlan(t, root, "plans/active/foo.md", "plans/archived/foo.md")

	writeRepoFile(t, root, "docs/index.md", "See [foo](../plans/archived/foo.md).\n")

	got := validatePlanArchivalLinkIntegrity("", nil, root)
	if got == "" {
		t.Fatalf("staged blob still references old path; expected block finding (commit-candidate semantics); got clean output")
	}
	if !strings.Contains(got, "broken_inbound_link") {
		t.Errorf("expected broken_inbound_link category; got: %s", got)
	}
	if !strings.Contains(got, "docs/index.md") {
		t.Errorf("expected referrer path in output; got: %s", got)
	}
}

func TestValidatePlanArchivalLinkIntegrity_TD1_TextualStagedFixWorktreeBroken(t *testing.T) {
	// Textual-scan variant of staged-has-fix. Staged removed the bare
	// path mention; worktree still has it. No warning expected.
	root := setupArchivalRepo(t)
	writeRepoFile(t, root, "plans/active/foo.md", "# foo\n")
	writeRepoFile(t, root, "docs/notes.md", "See plans/active/foo.md for context.\n")
	runGitFixture(t, root, "add", ".")
	runGitFixture(t, root, "commit", "-qm", "init", "--no-verify")
	archivePlan(t, root, "plans/active/foo.md", "plans/archived/foo.md")

	writeRepoFile(t, root, "docs/notes.md", "See plans/archived/foo.md for context.\n")
	runGitFixture(t, root, "add", "docs/notes.md")
	writeRepoFile(t, root, "docs/notes.md", "See plans/active/foo.md for context.\n")

	got := validatePlanArchivalLinkIntegrity("", nil, root)
	if got != "" {
		t.Errorf("staged blob already updated; expected clean pass; got:\n%s", got)
	}
}
