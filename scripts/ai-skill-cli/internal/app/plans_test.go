package app

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

// buildTreeFromPlans is a test helper that scans frontmatter from a tmp dir.
func scanTmpFrontmatter(t *testing.T, tmp string) []PlanFrontmatter {
	t.Helper()
	return scanAllPlanFrontmatter(tmp)
}

// ---------------------------------------------------------------------------
// Tree builder
// ---------------------------------------------------------------------------
func TestBuildPlanTree_MainAndChildren(t *testing.T) {
	tmp := t.TempDir()
	makeMain(t, tmp, "plans/active/main.md", "main-x", "in-progress")
	makeSub(t, tmp, "plans/active/sub-a.md", "sub-a", "main-x", "completed", "r", true)
	makeSub(t, tmp, "plans/active/sub-b.md", "sub-b", "main-x", "in-progress", "r", true)
	roots := buildPlanTree(scanTmpFrontmatter(t, tmp), "all", false)
	if len(roots) != 1 || roots[0].ID != "main-x" {
		t.Fatalf("expected 1 root main-x, got %+v", roots)
	}
	if len(roots[0].Children) != 2 {
		t.Fatalf("expected 2 children, got %d", len(roots[0].Children))
	}
	if roots[0].BlockerCount != 1 {
		t.Fatalf("expected 1 blocker (sub-b in-progress required), got %d", roots[0].BlockerCount)
	}
	if roots[0].ArchiveReady {
		t.Fatalf("expected archive_ready=false (1 required pending)")
	}
}

func TestBuildPlanTree_ArchiveReadyMain(t *testing.T) {
	tmp := t.TempDir()
	makeMain(t, tmp, "plans/active/main.md", "main-y", "in-progress")
	makeSub(t, tmp, "plans/active/sub-done.md", "sub-done", "main-y", "completed", "r", true)
	makeSub(t, tmp, "plans/active/sub-optional.md", "sub-opt", "main-y", "in-progress", "r", false)
	roots := buildPlanTree(scanTmpFrontmatter(t, tmp), "all", false)
	if !roots[0].ArchiveReady {
		t.Fatalf("expected archive_ready=true (only optional pending)")
	}
}

func TestBuildPlanTree_OrphanExcludedByDefault(t *testing.T) {
	tmp := t.TempDir()
	makeSub(t, tmp, "plans/active/orphan.md", "orphan", "no-such-main", "draft", "r", true)
	roots := buildPlanTree(scanTmpFrontmatter(t, tmp), "all", false)
	if len(roots) != 0 {
		t.Fatalf("orphans should be excluded by default, got %+v", roots)
	}
}

func TestBuildPlanTree_OrphanIncluded(t *testing.T) {
	tmp := t.TempDir()
	makeSub(t, tmp, "plans/active/orphan.md", "orphan", "no-such-main", "draft", "r", true)
	roots := buildPlanTree(scanTmpFrontmatter(t, tmp), "all", true)
	if len(roots) != 1 || !roots[0].IsOrphan {
		t.Fatalf("expected 1 orphan included, got %+v", roots)
	}
	if roots[0].UnresolvedRef != "no-such-main" {
		t.Fatalf("expected UnresolvedRef=no-such-main, got %q", roots[0].UnresolvedRef)
	}
}

func TestBuildPlanTree_StateFilter(t *testing.T) {
	tmp := t.TempDir()
	makeMain(t, tmp, "plans/archived/main-old.md", "main-old", "completed")
	makeMain(t, tmp, "plans/active/main-new.md", "main-new", "in-progress")
	rootsAll := buildPlanTree(scanTmpFrontmatter(t, tmp), "all", false)
	if len(rootsAll) != 2 {
		t.Fatalf("expected 2 roots (all), got %d", len(rootsAll))
	}
	rootsAct := buildPlanTree(scanTmpFrontmatter(t, tmp), "active", false)
	if len(rootsAct) != 1 || rootsAct[0].ID != "main-new" {
		t.Fatalf("expected only main-new for active, got %+v", rootsAct)
	}
	rootsArc := buildPlanTree(scanTmpFrontmatter(t, tmp), "archived", false)
	if len(rootsArc) != 1 || rootsArc[0].ID != "main-old" {
		t.Fatalf("expected only main-old for archived, got %+v", rootsArc)
	}
}

func TestBuildPlanTree_ArchivedParentActiveChild(t *testing.T) {
	tmp := t.TempDir()
	makeMain(t, tmp, "plans/archived/main.md", "main-z", "completed")
	makeSub(t, tmp, "plans/active/sub.md", "sub-z", "main-z", "in-progress", "r", false)
	roots := buildPlanTree(scanTmpFrontmatter(t, tmp), "all", false)
	if len(roots) != 1 || roots[0].ID != "main-z" {
		t.Fatalf("expected 1 root, got %+v", roots)
	}
	if len(roots[0].Children) != 1 || roots[0].Children[0].ID != "sub-z" {
		t.Fatalf("expected sub-z child of main-z across active/archived boundary, got %+v", roots[0].Children)
	}
}

func TestBuildPlanTree_MultiRoot(t *testing.T) {
	tmp := t.TempDir()
	makeMain(t, tmp, "plans/active/main-a.md", "main-a", "in-progress")
	makeMain(t, tmp, "plans/active/main-b.md", "main-b", "in-progress")
	roots := buildPlanTree(scanTmpFrontmatter(t, tmp), "all", false)
	if len(roots) != 2 {
		t.Fatalf("expected 2 roots, got %d", len(roots))
	}
}

// ---------------------------------------------------------------------------
// Renderers
// ---------------------------------------------------------------------------
func TestRenderText_BasicTree(t *testing.T) {
	tmp := t.TempDir()
	makeMain(t, tmp, "plans/active/main.md", "main-x", "in-progress")
	makeSub(t, tmp, "plans/active/sub.md", "sub-x", "main-x", "completed", "r", true)
	roots := buildPlanTree(scanTmpFrontmatter(t, tmp), "all", false)
	out := renderPlanTreeText(roots, tmp)
	if !strings.Contains(out, "main-x") || !strings.Contains(out, "sub-x") {
		t.Fatalf("text output missing nodes:\n%s", out)
	}
	if !strings.Contains(out, "└──") {
		t.Fatalf("text output missing ASCII branch:\n%s", out)
	}
}

func TestRenderJSON_ValidStructure(t *testing.T) {
	tmp := t.TempDir()
	makeMain(t, tmp, "plans/active/main.md", "main-x", "in-progress")
	makeSub(t, tmp, "plans/active/sub.md", "sub-x", "main-x", "completed", "r", true)
	roots := buildPlanTree(scanTmpFrontmatter(t, tmp), "all", false)
	out := renderPlanTreeJSON(roots, tmp)
	var parsed map[string]any
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("invalid JSON: %v\noutput: %s", err, out)
	}
	if parsed["top_level"].(float64) != 1 {
		t.Fatalf("expected top_level=1, got %v", parsed["top_level"])
	}
}

func TestRenderMarkdown_BasicList(t *testing.T) {
	tmp := t.TempDir()
	makeMain(t, tmp, "plans/active/main.md", "main-x", "in-progress")
	makeSub(t, tmp, "plans/active/sub.md", "sub-x", "main-x", "completed", "r", true)
	roots := buildPlanTree(scanTmpFrontmatter(t, tmp), "all", false)
	out := renderPlanTreeMarkdown(roots, tmp)
	if !strings.Contains(out, "- `main-x`") {
		t.Fatalf("markdown missing main-x bullet:\n%s", out)
	}
	if !strings.Contains(out, "  - `sub-x`") {
		t.Fatalf("markdown missing indented sub-x:\n%s", out)
	}
}

// ---------------------------------------------------------------------------
// runPlansTree integration (flag parsing + dispatch)
// ---------------------------------------------------------------------------
func TestRunPlansTree_DefaultText(t *testing.T) {
	tmp := t.TempDir()
	makeMain(t, tmp, "plans/active/main.md", "main-int", "draft")
	var stdout, stderr bytes.Buffer
	rc := runPlansTree([]string{"--root", tmp}, &stdout, &stderr)
	if rc != ExitSuccess {
		t.Fatalf("rc=%d stderr=%s", rc, stderr.String())
	}
	if !strings.Contains(stdout.String(), "main-int") {
		t.Fatalf("expected main-int in output:\n%s", stdout.String())
	}
}

func TestRunPlansTree_InvalidFormatRejected(t *testing.T) {
	var stdout, stderr bytes.Buffer
	rc := runPlansTree([]string{"--root", t.TempDir(), "--format", "html"}, &stdout, &stderr)
	if rc != ExitInvalidUsage {
		t.Fatalf("expected ExitInvalidUsage, got %d", rc)
	}
}

func TestRunPlansTree_InvalidStateRejected(t *testing.T) {
	var stdout, stderr bytes.Buffer
	rc := runPlansTree([]string{"--root", t.TempDir(), "--state", "purple"}, &stdout, &stderr)
	if rc != ExitInvalidUsage {
		t.Fatalf("expected ExitInvalidUsage, got %d", rc)
	}
}

func TestRunPlans_HelpAndDispatch(t *testing.T) {
	var stdout, stderr bytes.Buffer
	rc := runPlans([]string{"help"}, &stdout, &stderr)
	if rc != ExitSuccess {
		t.Fatalf("rc=%d", rc)
	}
	if !strings.Contains(stdout.String(), "tree") {
		t.Fatalf("help missing tree mention:\n%s", stdout.String())
	}
	rc = runPlans([]string{"unknown"}, &stdout, &stderr)
	if rc != ExitInvalidUsage {
		t.Fatalf("expected ExitInvalidUsage for unknown subcommand")
	}
}
