package app

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPlanFolderizeCluster_DryRun(t *testing.T) {
	tmp := t.TempDir()
	base := "2026-06-29-1430-preparatory-refactoring-workflow"
	writePlanFile(t, tmp, "plans/active/"+base+".md", "# main\nlink to "+base+"-dogfood-evidence.md\n")
	writePlanFile(t, tmp, "plans/active/"+base+"-dogfood-evidence.md", "# evidence\nback to "+base+".md\n")

	clusters := scanFlatPlanClusters(tmp)
	if len(clusters) != 1 {
		t.Fatal("expected cluster")
	}
	log, err := planFolderizeCluster(tmp, clusters[0], true)
	if err != nil {
		t.Fatal(err)
	}
	if len(log) != 2 {
		t.Fatalf("expected 2 moves, got %v", log)
	}
	if _, err := os.Stat(filepath.Join(tmp, "plans/active", base, "_plan.md")); !os.IsNotExist(err) {
		t.Fatal("dry-run should not create folder")
	}
}

func TestPlanFolderizeCluster_Apply(t *testing.T) {
	tmp := t.TempDir()
	base := "2026-06-29-1430-preparatory-refactoring-workflow"
	writePlanFile(t, tmp, "plans/active/"+base+".md", "# main\n")
	writePlanFile(t, tmp, "plans/active/"+base+"-dogfood-evidence.md", "# evidence\n")

	clusters := scanFlatPlanClusters(tmp)
	log, err := planFolderizeCluster(tmp, clusters[0], false)
	if err != nil {
		t.Fatal(err)
	}
	if len(log) != 2 {
		t.Fatalf("expected 2 moves, got %v", log)
	}
	planPath := filepath.Join(tmp, "plans/active", base, "_plan.md")
	if _, err := os.Stat(planPath); err != nil {
		t.Fatalf("_plan.md missing: %v", err)
	}
	compPath := filepath.Join(tmp, "plans/active", base, "01-dogfood-evidence.md")
	if _, err := os.Stat(compPath); err != nil {
		t.Fatalf("companion missing: %v", err)
	}
	if _, err := os.Stat(filepath.Join(tmp, "plans/active", base+".md")); !os.IsNotExist(err) {
		t.Fatal("flat main should be removed")
	}
}

func TestRunPlansFolderize_DryRunCLI(t *testing.T) {
	tmp := t.TempDir()
	base := "2026-06-29-1430-preparatory-refactoring-workflow"
	writePlanFile(t, tmp, "plans/active/"+base+".md", "# main\n")
	writePlanFile(t, tmp, "plans/active/"+base+"-dogfood-evidence.md", "# evidence\n")

	var stdout strings.Builder
	rc := runPlansFolderize([]string{"--root", tmp, "--cluster", base, "--dry-run"}, &stdout, &strings.Builder{})
	if rc != 0 {
		t.Fatalf("exit=%d stdout=%s", rc, stdout.String())
	}
	if !strings.Contains(stdout.String(), "_plan.md") {
		t.Fatalf("expected dry-run output: %s", stdout.String())
	}
}

func TestRunPlansFolderize_UnknownCluster(t *testing.T) {
	tmp := t.TempDir()
	rc := runPlansFolderize([]string{"--root", tmp, "--cluster", "no-such", "--dry-run"}, &strings.Builder{}, &strings.Builder{})
	if rc != ExitValidationFailed {
		t.Fatalf("expected validation failed, got %d", rc)
	}
}
