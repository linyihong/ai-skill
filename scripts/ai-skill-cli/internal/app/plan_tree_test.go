package app

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// makePlan writes a plan file under tmp with frontmatter content.
func makePlan(t *testing.T, tmp, rel, frontmatter string) {
	t.Helper()
	abs := filepath.Join(tmp, rel)
	if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
		t.Fatal(err)
	}
	body := frontmatter + "\n\n# Body\n"
	if err := os.WriteFile(abs, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func makeMain(t *testing.T, tmp, rel, id, status string) {
	fm := "---\n" +
		"id: " + id + "\n" +
		"plan_kind: main\n" +
		"status: " + status + "\n" +
		"owner: test\n" +
		"created: 2026-06-03\n" +
		"parent: null\n" +
		"---"
	makePlan(t, tmp, rel, fm)
}

func makeSub(t *testing.T, tmp, rel, id, parent, status, reason string, required bool) {
	reqStr := "false"
	if required {
		reqStr = "true"
	}
	fm := "---\n" +
		"id: " + id + "\n" +
		"plan_kind: sub\n" +
		"status: " + status + "\n" +
		"owner: test\n" +
		"created: 2026-06-03\n" +
		"parent: " + parent + "\n" +
		"required_for_completion: " + reqStr + "\n" +
		"sub_plan_reason: " + reason + "\n" +
		"---"
	makePlan(t, tmp, rel, fm)
}

// ---------------------------------------------------------------------------
// validatePlanTreeFrontmatter
// ---------------------------------------------------------------------------
func TestPlanTreeFrontmatter_MissingParent(t *testing.T) {
	tmp := t.TempDir()
	makePlan(t, tmp, "plans/active/sub-bad.md", "---\nid: sub-bad\nplan_kind: sub\nstatus: draft\nowner: t\ncreated: 2026-06-03\nrequired_for_completion: true\nsub_plan_reason: x\n---")
	got := validatePlanTreeFrontmatter("commit body", []string{"plans/active/sub-bad.md"}, tmp)
	if !strings.Contains(got, "missing: parent") {
		t.Fatalf("expected missing parent, got: %s", got)
	}
}

func TestPlanTreeFrontmatter_EmptyReason(t *testing.T) {
	tmp := t.TempDir()
	makeSub(t, tmp, "plans/active/sub-empty.md", "sub-empty", "main-x", "draft", `""`, true)
	got := validatePlanTreeFrontmatter("", []string{"plans/active/sub-empty.md"}, tmp)
	if !strings.Contains(got, "sub_plan_reason") {
		t.Fatalf("expected missing sub_plan_reason, got: %s", got)
	}
}

func TestPlanTreeFrontmatter_MissingRequired(t *testing.T) {
	tmp := t.TempDir()
	makePlan(t, tmp, "plans/active/sub-no-req.md", "---\nid: sub-no-req\nplan_kind: sub\nstatus: draft\nowner: t\ncreated: 2026-06-03\nparent: main-x\nsub_plan_reason: x\n---")
	got := validatePlanTreeFrontmatter("", []string{"plans/active/sub-no-req.md"}, tmp)
	if !strings.Contains(got, "required_for_completion") {
		t.Fatalf("expected missing required_for_completion, got: %s", got)
	}
}

func TestPlanTreeFrontmatter_PassValidSub(t *testing.T) {
	tmp := t.TempDir()
	makeSub(t, tmp, "plans/active/sub-good.md", "sub-good", "main-x", "draft", "good reason", true)
	got := validatePlanTreeFrontmatter("", []string{"plans/active/sub-good.md"}, tmp)
	if got != "" {
		t.Fatalf("expected pass, got: %s", got)
	}
}

func TestPlanTreeFrontmatter_SkipsMainAndUntagged(t *testing.T) {
	tmp := t.TempDir()
	makeMain(t, tmp, "plans/active/main.md", "main-x", "draft")
	makePlan(t, tmp, "plans/active/legacy.md", "# legacy no frontmatter")
	got := validatePlanTreeFrontmatter("", []string{"plans/active/main.md", "plans/active/legacy.md"}, tmp)
	if got != "" {
		t.Fatalf("main/legacy should pass, got: %s", got)
	}
}

func TestPlanTreeFrontmatter_OptOut(t *testing.T) {
	tmp := t.TempDir()
	makePlan(t, tmp, "plans/active/sub-bad.md", "---\nid: sub-bad\nplan_kind: sub\nstatus: draft\n---")
	got := validatePlanTreeFrontmatter("body\n[skip-plan-tree-frontmatter]\n", []string{"plans/active/sub-bad.md"}, tmp)
	if got != "" {
		t.Fatalf("opt-out should suppress, got: %s", got)
	}
}

// ---------------------------------------------------------------------------
// validatePlanTreeArchiveOrder
// ---------------------------------------------------------------------------
func TestPlanTreeArchiveOrder_BlocksOnPendingChild(t *testing.T) {
	tmp := t.TempDir()
	makeMain(t, tmp, "plans/archived/main-x.md", "main-x", "completed")
	makeSub(t, tmp, "plans/active/sub-pending.md", "sub-pending", "main-x", "in-progress", "r", true)
	got := validatePlanTreeArchiveOrder("", []string{"plans/archived/main-x.md"}, tmp)
	if !strings.Contains(got, "main-x") || !strings.Contains(got, "sub-pending") {
		t.Fatalf("expected block citing main-x + sub-pending, got: %s", got)
	}
}

func TestPlanTreeArchiveOrder_PassesWhenAllRequiredComplete(t *testing.T) {
	tmp := t.TempDir()
	makeMain(t, tmp, "plans/archived/main-x.md", "main-x", "completed")
	makeSub(t, tmp, "plans/archived/sub-done.md", "sub-done", "main-x", "completed", "r", true)
	makeSub(t, tmp, "plans/active/sub-optional.md", "sub-optional", "main-x", "in-progress", "r", false)
	got := validatePlanTreeArchiveOrder("", []string{"plans/archived/main-x.md"}, tmp)
	if got != "" {
		t.Fatalf("expected pass (required done, optional ignored), got: %s", got)
	}
}

func TestPlanTreeArchiveOrder_NoArchiveStaged(t *testing.T) {
	tmp := t.TempDir()
	makeMain(t, tmp, "plans/active/main.md", "main-x", "in-progress")
	got := validatePlanTreeArchiveOrder("", []string{"plans/active/main.md"}, tmp)
	if got != "" {
		t.Fatalf("expected pass when no archive in stage, got: %s", got)
	}
}

// ---------------------------------------------------------------------------
// validatePlanTreeParentReference
// ---------------------------------------------------------------------------
func TestPlanTreeParentReference_DanglingPointer(t *testing.T) {
	tmp := t.TempDir()
	makeSub(t, tmp, "plans/active/sub.md", "sub", "no-such-main", "draft", "r", true)
	got := validatePlanTreeParentReference("", []string{"plans/active/sub.md"}, tmp)
	if !strings.Contains(got, "no-such-main") {
		t.Fatalf("expected dangling parent, got: %s", got)
	}
}

func TestPlanTreeParentReference_ResolvesToArchivedMain(t *testing.T) {
	tmp := t.TempDir()
	makeMain(t, tmp, "plans/archived/main.md", "main-x", "completed")
	makeSub(t, tmp, "plans/active/sub.md", "sub", "main-x", "draft", "r", true)
	got := validatePlanTreeParentReference("", []string{"plans/active/sub.md"}, tmp)
	if got != "" {
		t.Fatalf("expected pass (archived main resolves), got: %s", got)
	}
}

// ---------------------------------------------------------------------------
// validatePlanTreeUniqueID
// ---------------------------------------------------------------------------
func TestPlanTreeUniqueID_Duplicate(t *testing.T) {
	tmp := t.TempDir()
	makeMain(t, tmp, "plans/active/a.md", "dup-id", "draft")
	makeMain(t, tmp, "plans/active/b.md", "dup-id", "draft")
	got := validatePlanTreeUniqueID("", []string{"plans/active/b.md"}, tmp)
	if !strings.Contains(got, "dup-id") {
		t.Fatalf("expected duplicate id violation, got: %s", got)
	}
}

func TestPlanTreeUniqueID_NoDuplicate(t *testing.T) {
	tmp := t.TempDir()
	makeMain(t, tmp, "plans/active/a.md", "uniq-1", "draft")
	makeMain(t, tmp, "plans/active/b.md", "uniq-2", "draft")
	got := validatePlanTreeUniqueID("", []string{"plans/active/a.md", "plans/active/b.md"}, tmp)
	if got != "" {
		t.Fatalf("expected pass, got: %s", got)
	}
}

func TestPlanTreeUniqueID_FixturesExcluded(t *testing.T) {
	tmp := t.TempDir()
	makeMain(t, tmp, "plans/active/a.md", "dup-id", "draft")
	// Fixture with the same id should be ignored.
	makeMain(t, tmp, "plans/active/foo/fixtures/dup.md", "dup-id", "draft")
	got := validatePlanTreeUniqueID("", []string{"plans/active/a.md"}, tmp)
	if got != "" {
		t.Fatalf("expected fixtures-segment to be excluded, got: %s", got)
	}
}

// ---------------------------------------------------------------------------
// validatePlanTreeFolderConvention
// ---------------------------------------------------------------------------
func TestPlanTreeFolderConvention_DepthWarning(t *testing.T) {
	tmp := t.TempDir()
	makeMain(t, tmp, "plans/active/foo/bar/baz/deep.md", "deep", "draft")
	got := validatePlanTreeFolderConvention("", []string{"plans/active/foo/bar/baz/deep.md"}, tmp)
	if !strings.Contains(got, "nested depth") {
		t.Fatalf("expected depth warning, got: %s", got)
	}
}

func TestPlanTreeFolderConvention_BadFilename(t *testing.T) {
	tmp := t.TempDir()
	makeSub(t, tmp, "plans/active/cluster/bad-name.md", "bad", "main-x", "draft", "r", true)
	got := validatePlanTreeFolderConvention("", []string{"plans/active/cluster/bad-name.md"}, tmp)
	if !strings.Contains(got, "_plan.md") {
		t.Fatalf("expected filename advisory, got: %s", got)
	}
}

func TestPlanTreeFolderConvention_GoodLayout(t *testing.T) {
	tmp := t.TempDir()
	makeMain(t, tmp, "plans/active/cluster/_plan.md", "main", "draft")
	makeSub(t, tmp, "plans/active/cluster/01-schema.md", "sub", "main", "draft", "r", true)
	got := validatePlanTreeFolderConvention("",
		[]string{"plans/active/cluster/_plan.md", "plans/active/cluster/01-schema.md"}, tmp)
	if got != "" {
		t.Fatalf("expected clean, got: %s", got)
	}
}

// ---------------------------------------------------------------------------
// parsePlanFrontmatterFromBytes sanity checks
// ---------------------------------------------------------------------------
func TestParseFrontmatter_FoldedReason(t *testing.T) {
	body := "---\nid: x\nplan_kind: sub\nstatus: draft\nowner: t\ncreated: 2026-06-03\nparent: m\nrequired_for_completion: true\nsub_plan_reason: >\n  this reason\n  spans two lines\n---\n"
	pf := parsePlanFrontmatterFromBytes("test.md", []byte(body))
	if !pf.HasFrontmatter {
		t.Fatal("expected HasFrontmatter")
	}
	if !strings.Contains(pf.SubPlanReason, "spans two lines") {
		t.Fatalf("folded reason not parsed, got: %q", pf.SubPlanReason)
	}
}

func TestParseFrontmatter_NoFrontmatter(t *testing.T) {
	pf := parsePlanFrontmatterFromBytes("test.md", []byte("# Just markdown\nno yaml here\n"))
	if pf.HasFrontmatter {
		t.Fatal("expected HasFrontmatter=false")
	}
}
