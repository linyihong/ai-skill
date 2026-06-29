package app

import "github.com/linyihong/Ai-skill/scripts/ai-skill-cli/internal/planvalidate"

// Phase 3.3 (Consumer Equivalence) — Canonical Observation Record (E.1).
//
// Consumer equivalence compares the OBSERVATION BOUNDARY, not the execution
// trace. Every consumer (manual CLI, commit-msg hook, CI) reduces to the same
// ObservationRecord before comparison, so transport-layer differences cannot
// leak into the equivalence check. The record deliberately contains ONLY the
// MUST-equal dimensions; it has no field for ExitCode / ExecutionMode /
// SnapshotOrigin / Timing / Message (those MAY differ or are IGNORE).
type ObservationRecord struct {
	// Findings: RuleID -> Blocking (presence-level; the MUST-equal observation).
	Findings map[string]bool
	// OptOutEffect: RuleID -> suppressed-by-opt-out (effective policy result).
	OptOutEffect map[string]bool
	// DiscoveryScope: the dirs scanned; MUST be equal across consumers.
	DiscoveryScope []string
}

func newObservationRecord() ObservationRecord {
	return ObservationRecord{
		Findings:       map[string]bool{},
		OptOutEffect:   map[string]bool{},
		DiscoveryScope: planTreeDiscoveryScope(),
	}
}

// planTreeDiscoveryScope is the fixed discovery scope (scope A). All consumers
// share it, so it is equal by construction — proving "discovery scope MUST equal"
// without coupling consumers.
func planTreeDiscoveryScope() []string {
	return []string{"plans/active", "plans/archived"}
}

// ruleApplicability reports, per core RuleID, whether the rule had any input to
// evaluate over the model set. Phase 3.3c ONLY (R.3): it distinguishes "rule
// applicable and passed" from "rule silently inapplicable" — two states that both
// yield zero findings but are NOT observation-preserving. Deliberately NOT part
// of the COR (does not pollute the 3.3a/3.3b equivalence surface).
func ruleApplicability(models []planvalidate.NormalizedPlanModel) map[string]bool {
	ap := map[string]bool{
		"plan_tree.frontmatter":      false,
		"plan_tree.unique_id":        false,
		"plan_tree.parent_reference": false,
		"plan_tree.archive_order":    false,
	}
	for _, m := range models {
		if m.PlanKind == "sub" {
			ap["plan_tree.frontmatter"] = true
		}
		if m.ID != "" {
			ap["plan_tree.unique_id"] = true
		}
		if m.Parent != "" {
			ap["plan_tree.parent_reference"] = true
		}
		isMain := m.PlanKind == "main" || (m.PlanKind == "" && m.Parent == "")
		if isMain && m.Location == "archived" {
			ap["plan_tree.archive_order"] = true
		}
	}
	return ap
}

// legacyObservation builds the COR for the commit-msg hook consumer (the legacy
// plan-tree validators). Opt-out is resolved from the commit message (the hook's
// transport); a suppressed validator records OptOutEffect, not a Finding.
func legacyObservation(text string, staged []string, root string) ObservationRecord {
	rec := newObservationRecord()
	for _, m := range planTreeLegacyMirror {
		if hasOptOutTrailer(text, m.optOut) {
			rec.OptOutEffect[m.ruleID] = true
		}
		if m.runLegacy(text, staged, root) != "" {
			rec.Findings[m.ruleID] = true
		}
	}
	return rec
}

// engineObservation builds the COR for an engine-backed consumer (manual CLI /
// CI). The engine is policy-free: the consumer resolves opt-out from its own
// transport (here, the same commit text) and applies it — the engine never sees
// it. This is what keeps equivalence from pushing context into the engine.
func engineObservation(text string, root string) ObservationRecord {
	optedOut := map[string]bool{}
	for _, m := range planTreeLegacyMirror {
		if hasOptOutTrailer(text, m.optOut) {
			optedOut[m.ruleID] = true
		}
	}
	models, compat := normalizedPlansFromRoot(root)
	fs := planvalidate.Validate(planvalidate.ValidationContext{Root: root, ExecutionMode: planvalidate.ModeManual}, models)
	fs = append(fs, compat...)
	rec := newObservationRecord()
	for _, f := range fs {
		if optedOut[f.RuleID] {
			rec.OptOutEffect[f.RuleID] = true
			continue
		}
		rec.Findings[f.RuleID] = f.Blocking
	}
	return rec
}
