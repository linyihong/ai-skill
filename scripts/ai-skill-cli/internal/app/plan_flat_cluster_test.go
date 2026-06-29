package app

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writePlanFile(t *testing.T, tmp, rel, body string) {
	t.Helper()
	abs := filepath.Join(tmp, rel)
	if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(abs, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestScanFlatPlanClusters_PreparatoryRefactoringFixture(t *testing.T) {
	tmp := t.TempDir()
	base := "2026-06-29-1430-preparatory-refactoring-workflow"
	writePlanFile(t, tmp, "plans/active/"+base+".md", "---\nid: x\n---\n# main\n")
	writePlanFile(t, tmp, "plans/active/"+base+"-dogfood-evidence.md", "# evidence\n")

	clusters := scanFlatPlanClusters(tmp)
	if len(clusters) != 1 {
		t.Fatalf("expected 1 cluster, got %d", len(clusters))
	}
	c := clusters[0]
	if c.Base != base {
		t.Fatalf("base=%q want %q", c.Base, base)
	}
	if c.memberCount() != 2 {
		t.Fatalf("memberCount=%d want 2", c.memberCount())
	}
	if len(c.Companions) != 1 || c.Companions[0].Suffix != "dogfood-evidence" {
		t.Fatalf("companions=%+v", c.Companions)
	}
}

func TestScanFlatPlanClusters_SingleFileIgnored(t *testing.T) {
	tmp := t.TempDir()
	writePlanFile(t, tmp, "plans/active/2026-06-29-1430-standalone.md", "# solo\n")
	if len(scanFlatPlanClusters(tmp)) != 0 {
		t.Fatal("single top-level plan should not form a cluster")
	}
}

func TestFlatClusterWarningsForStaged(t *testing.T) {
	tmp := t.TempDir()
	base := "2026-06-29-1430-preparatory-refactoring-workflow"
	writePlanFile(t, tmp, "plans/active/"+base+".md", "# main\n")
	writePlanFile(t, tmp, "plans/active/"+base+"-dogfood-evidence.md", "# evidence\n")

	w := flatClusterWarningsForStaged([]string{"plans/active/" + base + "-dogfood-evidence.md"}, tmp)
	if len(w) != 1 {
		t.Fatalf("expected 1 warning, got %d: %v", len(w), w)
	}
	if !strings.Contains(w[0], "folderize") {
		t.Fatalf("expected folderize hint, got: %s", w[0])
	}
}

func TestValidatePlanTreeFolderConvention_FlatClusterWarning(t *testing.T) {
	tmp := t.TempDir()
	base := "2026-06-29-1430-preparatory-refactoring-workflow"
	writePlanFile(t, tmp, "plans/active/"+base+".md", "# main\n")
	writePlanFile(t, tmp, "plans/active/"+base+"-dogfood-evidence.md", "# evidence\n")

	got := validatePlanTreeFolderConvention("", []string{"plans/active/" + base + ".md"}, tmp)
	if !strings.Contains(got, "flat multi-file cluster") {
		t.Fatalf("expected flat cluster warning, got: %s", got)
	}
}

func TestRewriteClusterLinks(t *testing.T) {
	c := flatPlanCluster{
		Base: "2026-06-29-1430-preparatory-refactoring-workflow",
		Companions: []flatPlanCompanion{{
			Suffix: "dogfood-evidence",
		}},
	}
	in := "See [evidence](2026-06-29-1430-preparatory-refactoring-workflow-dogfood-evidence.md) and [main](2026-06-29-1430-preparatory-refactoring-workflow.md)"
	out := rewriteClusterLinks(in, c)
	if !strings.Contains(out, "01-dogfood-evidence.md") || !strings.Contains(out, "_plan.md") {
		t.Fatalf("rewrite failed: %s", out)
	}
}

func TestCompanionTargetName(t *testing.T) {
	if got := companionTargetName(1, "dogfood-evidence"); got != "01-dogfood-evidence.md" {
		t.Fatalf("got %q", got)
	}
}
