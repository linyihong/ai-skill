package planvalidate

import (
	"reflect"
	"strings"
	"testing"
)

// Phase 2.2 / Gate D: the engine integration test is the FIRST consumer (not the
// CLI — Gate D). Tests are grouped by BEHAVIOR, not by validator (Gate D.3), so
// swapping/adding consumers in later phases does not force a test rewrite.

func b(v bool) *bool { return &v }

// ruleIDs is a tiny helper to assert which rules fired regardless of order.
func ruleIDs(fs []Finding) map[string]int {
	m := map[string]int{}
	for _, f := range fs {
		m[f.RuleID]++
	}
	return m
}

// --- behavior: valid -------------------------------------------------------

func TestEngine_Valid_NoFindings(t *testing.T) {
	plans := []NormalizedPlanModel{
		{Path: "plans/active/_plan.md", ID: "main", PlanKind: "main", Status: "draft"},
		{Path: "plans/active/01.md", ID: "sub", PlanKind: "sub", Status: "draft", Parent: "main", RequiredForCompletion: b(true), SubPlanReason: "r"},
	}
	if fs := Validate(ValidationContext{ExecutionMode: ModeManual}, plans); len(fs) != 0 {
		t.Fatalf("valid tree produced findings: %+v", fs)
	}
}

// --- behavior: missing (incomplete sub-plan frontmatter) -------------------

func TestEngine_Missing_FrontmatterFieldsFlagged(t *testing.T) {
	plans := []NormalizedPlanModel{
		{Path: "plans/active/_plan.md", ID: "main", PlanKind: "main"},
		{Path: "plans/active/01.md", ID: "sub", PlanKind: "sub", Parent: "main"}, // no sub_plan_reason, no required_for_completion
	}
	got := ruleIDs(Validate(ValidationContext{}, plans))
	if got["plan_tree.frontmatter"] != 2 {
		t.Fatalf("expected 2 frontmatter findings (reason + required), got %+v", got)
	}
}

// --- behavior: broken_link (unresolved parent) -----------------------------

func TestEngine_BrokenLink_UnresolvedParent(t *testing.T) {
	plans := []NormalizedPlanModel{
		{Path: "plans/active/01.md", ID: "sub", PlanKind: "sub", Parent: "ghost", RequiredForCompletion: b(true), SubPlanReason: "r"},
	}
	if ruleIDs(Validate(ValidationContext{}, plans))["plan_tree.parent_reference"] != 1 {
		t.Fatalf("expected unresolved-parent finding")
	}
}

// --- behavior: archive_required (archived main, incomplete required sub) ----

func TestEngine_ArchiveRequired_IncompleteSubBlocks(t *testing.T) {
	plans := []NormalizedPlanModel{
		{Path: "plans/archived/_plan.md", ID: "main", PlanKind: "main", Location: "archived", Status: "completed"},
		{Path: "plans/archived/01.md", ID: "sub", PlanKind: "sub", Location: "archived", Status: "draft", Parent: "main", RequiredForCompletion: b(true), SubPlanReason: "r"},
	}
	if ruleIDs(Validate(ValidationContext{}, plans))["plan_tree.archive_order"] != 1 {
		t.Fatalf("expected archive-order finding for incomplete required sub")
	}
}

// --- behavior: opt_out (consumer applies effective policy, engine does not) -

func TestEngine_OptOut_IsConsumerSideNotEngineSide(t *testing.T) {
	plans := []NormalizedPlanModel{
		{Path: "plans/active/01.md", ID: "sub", PlanKind: "sub", Parent: "ghost", RequiredForCompletion: b(true), SubPlanReason: "r"},
	}
	all := Validate(ValidationContext{}, plans)
	if len(all) == 0 {
		t.Fatalf("engine should still emit the finding regardless of opt-out")
	}
	// The CONSUMER (here, the test) resolves effective policy and filters. The
	// engine stays policy-free: opt-out never suppresses inside Validate.
	optedOut := map[string]bool{"plan_tree.parent_reference": true}
	var effective []Finding
	for _, f := range all {
		if !optedOut[f.RuleID] {
			effective = append(effective, f)
		}
	}
	if len(effective) != 0 {
		t.Fatalf("consumer-side opt-out should suppress the finding, got %+v", effective)
	}
}

// --- behavior: mixed (collect-all, no fail-fast — Gate D.2) -----------------

func TestEngine_Mixed_CollectsAllFindings(t *testing.T) {
	plans := []NormalizedPlanModel{
		{Path: "plans/active/a.md", ID: "dup", PlanKind: "sub", Parent: "ghost"},        // unique(later) + parent + frontmatter
		{Path: "plans/active/b.md", ID: "dup", PlanKind: "main"},                         // duplicate id
	}
	got := ruleIDs(Validate(ValidationContext{}, plans))
	// Must see multiple distinct rule classes in one pass (no fail-fast).
	if got["plan_tree.unique_id"] == 0 || got["plan_tree.parent_reference"] == 0 || got["plan_tree.frontmatter"] == 0 {
		t.Fatalf("collect-all expected multiple rule classes, got %+v", got)
	}
}

// --- Gate D.4: negative evidence -------------------------------------------

// The portable boundary must hold by CONSTRUCTION: the engine cannot express an
// excluded validator because NormalizedPlanModel carries none of their inputs.
// runtime-trigger-wiring needs routing-registry/runtime.db data; checkbox-sync /
// status-sync need the commit message. If any such field appeared on the model,
// an excluded validator would become expressible and the boundary would be wrong.
func TestEngine_CannotExpressExcludedValidators(t *testing.T) {
	tp := reflect.TypeOf(NormalizedPlanModel{})
	forbidden := []string{"route", "registry", "commit", "message", "runtime", "discovery", "diff"}
	for i := 0; i < tp.NumField(); i++ {
		name := strings.ToLower(tp.Field(i).Name)
		for _, bad := range forbidden {
			if strings.Contains(name, bad) {
				t.Fatalf("NormalizedPlanModel field %q could feed an EXCLUDED validator (%q) — portable boundary breached", tp.Field(i).Name, bad)
			}
		}
	}
}
