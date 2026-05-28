package app

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestValidatePlanCheckboxSyncOptOut(t *testing.T) {
	body := "feat(x): Phase 1 of plans/active/foo.md\n\n[skip-plan-checkbox-sync]\n"
	got := validatePlanCheckboxSync(body, []string{"scripts/ai-skill-cli/internal/app/hooks.go"}, "")
	if got != "" {
		t.Errorf("opt-out trailer should bypass, got %q", got)
	}
}

func TestValidatePlanCheckboxSyncNoPlanReference(t *testing.T) {
	body := "feat(x): generic change without plan reference"
	got := validatePlanCheckboxSync(body, []string{"scripts/ai-skill-cli/internal/app/x.go"}, "")
	if got != "" {
		t.Errorf("no plan reference should bypass, got %q", got)
	}
}

func TestValidatePlanCheckboxSyncDocOnlyBypass(t *testing.T) {
	body := "docs(plans): refer to plans/active/foo.md"
	got := validatePlanCheckboxSync(body, []string{"README.md", "docs/x.md"}, "")
	if got != "" {
		t.Errorf("doc-only staging should bypass (no code work heuristic match), got %q", got)
	}
}

func TestValidatePlanCheckboxSyncPlanNotStaged(t *testing.T) {
	repo := initTempGitRepo(t)
	body := "feat(audit): land code referencing plans/active/foo.md"
	staged := []string{"scripts/ai-skill-cli/internal/app/x.go"}
	got := validatePlanCheckboxSync(body, staged, repo)
	if !strings.Contains(got, "plans/active/foo.md") || !strings.Contains(got, "not staged") {
		t.Errorf("expected violation citing plan not staged, got %q", got)
	}
}

func TestValidatePlanCheckboxSyncStagedButNoCheckboxFlip(t *testing.T) {
	repo := initTempGitRepo(t)
	planRel := "plans/active/foo.md"
	planAbs := filepath.Join(repo, planRel)
	writeFile(t, planAbs, "# Plan\n\n- [ ] task A\n- [ ] task B\n")
	runGit(t, repo, "add", planRel)
	runGit(t, repo, "commit", "-m", "init plan")
	// Edit prose but don't flip a checkbox
	writeFile(t, planAbs, "# Plan\n\nUpdated prose.\n\n- [ ] task A\n- [ ] task B\n")
	writeFile(t, filepath.Join(repo, "scripts/ai-skill-cli/internal/app/x.go"), "package app\n")
	runGit(t, repo, "add", planRel, "scripts/ai-skill-cli/internal/app/x.go")

	body := "feat(audit): work cited in plans/active/foo.md"
	staged := []string{planRel, "scripts/ai-skill-cli/internal/app/x.go"}
	got := validatePlanCheckboxSync(body, staged, repo)
	if !strings.Contains(got, "no `[ ]` → `[x]` transition") {
		t.Errorf("expected violation citing missing checkbox flip, got %q", got)
	}
}

func TestValidatePlanCheckboxSyncStagedWithCheckboxFlipPasses(t *testing.T) {
	repo := initTempGitRepo(t)
	planRel := "plans/active/foo.md"
	planAbs := filepath.Join(repo, planRel)
	writeFile(t, planAbs, "# Plan\n\n- [ ] task A\n- [ ] task B\n")
	runGit(t, repo, "add", planRel)
	runGit(t, repo, "commit", "-m", "init plan")
	// Flip one checkbox
	writeFile(t, planAbs, "# Plan\n\n- [x] task A\n- [ ] task B\n")
	writeFile(t, filepath.Join(repo, "scripts/ai-skill-cli/internal/app/x.go"), "package app\n")
	runGit(t, repo, "add", planRel, "scripts/ai-skill-cli/internal/app/x.go")

	body := "feat(audit): Phase 1 done per plans/active/foo.md"
	staged := []string{planRel, "scripts/ai-skill-cli/internal/app/x.go"}
	got := validatePlanCheckboxSync(body, staged, repo)
	if got != "" {
		t.Errorf("checkbox flip should pass validator, got %q", got)
	}
}

func TestPlanDiffFlipsCheckboxDetectsIndentedFlip(t *testing.T) {
	diff := "diff --git a/x b/x\n--- a/x\n+++ b/x\n@@ -1 +1 @@\n-- [ ] foo\n+- [x] foo\n"
	if !planDiffFlipsCheckbox(diff) {
		t.Error("expected unindented flip to be detected")
	}
	diff2 := "diff --git a/x b/x\n--- a/x\n+++ b/x\n@@ -1 +1 @@\n-  - [ ] foo\n+  - [x] foo\n"
	if !planDiffFlipsCheckbox(diff2) {
		t.Error("expected indented flip to be detected")
	}
}

func TestPlanDiffFlipsCheckboxIgnoresAddedUnchecked(t *testing.T) {
	diff := "+++ b/x\n+- [ ] new task added\n"
	if planDiffFlipsCheckbox(diff) {
		t.Error("adding `- [ ]` should not count as a checkbox flip")
	}
}

func TestPlanDiffFlipsCheckboxAcceptsCapitalX(t *testing.T) {
	diff := "+++ b/x\n+- [X] done with capital X\n"
	if !planDiffFlipsCheckbox(diff) {
		t.Error("expected capital `[X]` to count as a flip")
	}
}
