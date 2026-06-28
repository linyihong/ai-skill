package app

import (
	"reflect"
	"strings"
	"testing"
)

// Phase 3.3a — manual (engine/CLI) ↔ hook (legacy) consumer equivalence, compared
// via the Canonical Observation Record (E.1). Equivalence is on the observation
// boundary (Findings + OptOutEffect + DiscoveryScope), never the execution trace.

func validTree(t *testing.T) string {
	tmp := t.TempDir()
	makePlan(t, tmp, "plans/active/_plan.md",
		"---\nid: m\nplan_kind: main\nstatus: draft\nowner: t\ncreated: 2026-06-25\nparent: null\n---")
	makePlan(t, tmp, "plans/active/01.md",
		"---\nid: s\nplan_kind: sub\nstatus: draft\nowner: t\ncreated: 2026-06-25\nparent: m\nrequired_for_completion: true\nsub_plan_reason: x\n---")
	return tmp
}

func violationTree(t *testing.T) string {
	tmp := t.TempDir()
	makePlan(t, tmp, "plans/active/_plan.md",
		"---\nid: m\nplan_kind: main\nstatus: draft\nowner: t\ncreated: 2026-06-25\nparent: null\n---")
	makePlan(t, tmp, "plans/active/01.md",
		"---\nid: s\nplan_kind: sub\nstatus: draft\nowner: t\ncreated: 2026-06-25\nparent: ghost\nrequired_for_completion: true\nsub_plan_reason: x\n---")
	return tmp
}

func assertCOREqual(t *testing.T, label, text string, staged []string, root string) {
	t.Helper()
	legacy := legacyObservation(text, staged, root)
	engine := engineObservation(text, root)
	if !reflect.DeepEqual(legacy, engine) {
		t.Fatalf("%s: manual!=hook COR\n legacy=%+v\n engine=%+v", label, legacy, engine)
	}
}

func TestConsumerEquivalence_ManualVsHook_Valid(t *testing.T) {
	root := validTree(t)
	assertCOREqual(t, "valid", "feat: x\n", []string{"plans/active/01.md"}, root)
}

func TestConsumerEquivalence_ManualVsHook_Violation(t *testing.T) {
	root := violationTree(t)
	assertCOREqual(t, "violation", "feat: x\n", []string{"plans/active/01.md"}, root)
	// And the observation must actually carry the parent_reference finding.
	if !legacyObservation("feat: x\n", []string{"plans/active/01.md"}, root).Findings["plan_tree.parent_reference"] {
		t.Fatalf("expected parent_reference in observation")
	}
}

func TestConsumerEquivalence_ManualVsHook_OptOut(t *testing.T) {
	root := violationTree(t)
	text := "feat: x\n\n[skip-plan-tree-parent-reference]\n"
	assertCOREqual(t, "opt-out", text, []string{"plans/active/01.md"}, root)
	// Opt-out must move parent_reference from Findings to OptOutEffect in both.
	rec := engineObservation(text, root)
	if rec.Findings["plan_tree.parent_reference"] || !rec.OptOutEffect["plan_tree.parent_reference"] {
		t.Fatalf("opt-out should suppress finding into OptOutEffect, got %+v", rec)
	}
}

// E.1 guard: the COR must not expose any transport dimension, so transport can
// never leak into equivalence.
func TestObservationRecord_ExcludesTransport(t *testing.T) {
	tp := reflect.TypeOf(ObservationRecord{})
	for i := 0; i < tp.NumField(); i++ {
		name := strings.ToLower(tp.Field(i).Name)
		for _, bad := range []string{"exit", "mode", "snapshot", "timing", "message"} {
			if strings.Contains(name, bad) {
				t.Fatalf("ObservationRecord field %q exposes transport dimension %q", tp.Field(i).Name, bad)
			}
		}
	}
}
