package app

import (
	"strings"
	"testing"
)

// Phase 2.3b (Shadow Confidence Window) Stage-2 fixture replay.
//
// These exercise the shadow ACCOUNTING PATHS (transport + parity-on-violation)
// directly through planValidateShadowCheck, rather than waiting for natural
// commits. Per the closure rule, 2.3b validates that the divergence buckets are
// exercised, not that real-world violations were observed. We deliberately do NOT
// manufacture missing/extra divergence (that would be an engine-challenge test,
// not shadow confidence).

// Replay A — a real violation (sub parent does not resolve) must fire in BOTH
// legacy and engine, landing in the `same` bucket: parity on a violation, with no
// genuine missing/extra and no exit-code effect (Gate C.1).
func TestPlanValidateShadow_ViolationParity(t *testing.T) {
	tmp := t.TempDir()
	makePlan(t, tmp, "plans/active/_plan.md",
		"---\nid: m\nplan_kind: main\nstatus: draft\nowner: t\ncreated: 2026-06-24\nparent: null\n---")
	makePlan(t, tmp, "plans/active/01-x.md",
		"---\nid: s\nplan_kind: sub\nstatus: draft\nowner: t\ncreated: 2026-06-24\nparent: ghost\nrequired_for_completion: true\nsub_plan_reason: x\n---")

	chk := planValidateShadowCheck(commitMsgCtx{
		text:   "feat: trigger parent_reference\n",
		staged: []string{"plans/active/01-x.md"},
		root:   tmp,
	})
	if chk.Status != "ok" {
		t.Fatalf("violation parity should converge (ok), got status=%q msg=%q", chk.Status, chk.Message)
	}
	if !strings.Contains(chk.Message, "same=plan_tree.parent_reference") {
		t.Fatalf("violation should land in same bucket (both fire), got %q", chk.Message)
	}
	if !strings.Contains(chk.Message, "missing=-") || !strings.Contains(chk.Message, "extra=-") {
		t.Fatalf("violation parity must have no genuine missing/extra, got %q", chk.Message)
	}
}

// Replay B — the same violation with an opt-out trailer: legacy suppresses, the
// policy-free engine still emits, so the difference is reclassified to the
// `transport` bucket (benign, converged), never `extra`.
func TestPlanValidateShadow_OptOutTransport(t *testing.T) {
	tmp := t.TempDir()
	makePlan(t, tmp, "plans/active/_plan.md",
		"---\nid: m\nplan_kind: main\nstatus: draft\nowner: t\ncreated: 2026-06-24\nparent: null\n---")
	makePlan(t, tmp, "plans/active/01-x.md",
		"---\nid: s\nplan_kind: sub\nstatus: draft\nowner: t\ncreated: 2026-06-24\nparent: ghost\nrequired_for_completion: true\nsub_plan_reason: x\n---")

	chk := planValidateShadowCheck(commitMsgCtx{
		text:   "feat: opt out parent check\n\n[skip-plan-tree-parent-reference]\n",
		staged: []string{"plans/active/01-x.md"},
		root:   tmp,
	})
	if chk.Status != "ok" {
		t.Fatalf("opt-out transport is benign and should converge (ok), got status=%q msg=%q", chk.Status, chk.Message)
	}
	if !strings.Contains(chk.Message, "transport=plan_tree.parent_reference") {
		t.Fatalf("opt-out should reclassify to transport, got %q", chk.Message)
	}
	if !strings.Contains(chk.Message, "extra=-") {
		t.Fatalf("opt-out must not be a genuine extra, got %q", chk.Message)
	}
}
