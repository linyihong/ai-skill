package app

import (
	"path/filepath"
	"strings"
	"testing"
)

// --- findArchivedPlans tests ---

func TestFindArchivedPlansDetectsMove(t *testing.T) {
	staged := []string{
		"plans/active/my-plan.md",
		"plans/archived/my-plan.md",
	}
	got := findArchivedPlans(staged)
	if len(got) != 1 || got[0] != "plans/archived/my-plan.md" {
		t.Errorf("expected archived path, got %v", got)
	}
}

func TestFindArchivedPlansIgnoresAddWithoutDelete(t *testing.T) {
	staged := []string{"plans/archived/my-plan.md"}
	got := findArchivedPlans(staged)
	if len(got) != 0 {
		t.Errorf("expected no match when active not deleted, got %v", got)
	}
}

func TestFindArchivedPlansIgnoresUnrelatedFiles(t *testing.T) {
	staged := []string{
		"scripts/ai-skill-cli/internal/app/hooks.go",
		"runtime/core-bootstrap.yaml",
	}
	got := findArchivedPlans(staged)
	if len(got) != 0 {
		t.Errorf("expected no archived plans, got %v", got)
	}
}

// --- bodyJustifiesUnchecked tests ---

func TestBodyJustifiesUncheckedDeferred(t *testing.T) {
	if !bodyJustifiesUnchecked("remaining items deferred to follow-up plan") {
		t.Error("expected 'deferred' to justify")
	}
}

func TestBodyJustifiesUncheckedNonGoal(t *testing.T) {
	if !bodyJustifiesUnchecked("these are non-goal items") {
		t.Error("expected 'non-goal' to justify")
	}
}

func TestBodyJustifiesUncheckedChinese(t *testing.T) {
	if !bodyJustifiesUnchecked("剩餘項目延後處理") {
		t.Error("expected '延後' to justify")
	}
}

func TestBodyJustifiesUncheckedNoKeyword(t *testing.T) {
	if bodyJustifiesUnchecked("archive old plan, all done") {
		t.Error("expected no justification without keywords")
	}
}

// --- validatePlanArchivalAudit fixture tests ---

func TestValidatePlanArchivalAuditOptOut(t *testing.T) {
	body := "chore: archive plan\n\n[skip-plan-archival-audit]\n"
	got := validatePlanArchivalAudit(body, []string{}, "")
	if got != "" {
		t.Errorf("opt-out should bypass validator, got %q", got)
	}
}

func TestValidatePlanArchivalAuditNoArchivalNoTrigger(t *testing.T) {
	body := "feat: add feature"
	staged := []string{"scripts/ai-skill-cli/internal/app/hooks.go"}
	got := validatePlanArchivalAudit(body, staged, "")
	if got != "" {
		t.Errorf("non-archival commit should not trigger, got %q", got)
	}
}

func TestValidatePlanArchivalAuditUnjustifiedCheckboxBlocked(t *testing.T) {
	repo := initTempGitRepo(t)

	planRel := "plans/archived/foo.md"
	planAbs := filepath.Join(repo, planRel)
	writeFile(t, planAbs, "# Plan\n\n- [ ] unfinished task\n- [x] done task\n")

	body := "chore: archive foo plan"
	staged := []string{"plans/active/foo.md", "plans/archived/foo.md"}
	got := validatePlanArchivalAudit(body, staged, repo)
	if !strings.Contains(got, "plan-archival-audit") {
		t.Errorf("expected block with plan-archival-audit, got %q", got)
	}
	if !strings.Contains(got, "foo.md") {
		t.Errorf("expected violation to cite file name, got %q", got)
	}
}

func TestValidatePlanArchivalAuditJustifiedCheckboxPasses(t *testing.T) {
	repo := initTempGitRepo(t)

	planRel := "plans/archived/foo.md"
	planAbs := filepath.Join(repo, planRel)
	writeFile(t, planAbs, "# Plan\n\n- [ ] unfinished task\n")

	body := "chore: archive foo plan\n\nRemaining items deferred to follow-up plan."
	staged := []string{"plans/active/foo.md", "plans/archived/foo.md"}
	got := validatePlanArchivalAudit(body, staged, repo)
	if got != "" {
		t.Errorf("justified unchecked should pass, got %q", got)
	}
}

func TestValidatePlanArchivalAuditAllCheckedPasses(t *testing.T) {
	repo := initTempGitRepo(t)

	planRel := "plans/archived/foo.md"
	planAbs := filepath.Join(repo, planRel)
	writeFile(t, planAbs, "# Plan\n\n- [x] all done\n- [x] also done\n")

	body := "chore: archive completed foo plan"
	staged := []string{"plans/active/foo.md", "plans/archived/foo.md"}
	got := validatePlanArchivalAudit(body, staged, repo)
	if got != "" {
		t.Errorf("all-checked plan should pass, got %q", got)
	}
}
