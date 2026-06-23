package planvalidate

import "fmt"

// Phase 2.2 (Gate D): the validation engine.
//
// The engine is pure: it consumes a ValidationContext and the normalized plan
// set and returns ALL findings (no fail-fast — Gate D.2). It implements only the
// plan_profile.core rules — the portable, dependency-free plan-structure checks.
//
// Negative evidence (Gate D.4): the engine deliberately CANNOT express the
// excluded validators. NormalizedPlanModel carries no commit message (so
// checkbox-sync / status-sync are un-implementable here) and no routing-registry
// / runtime.db data (so runtime-trigger-wiring is un-implementable here). The
// boundary therefore holds by construction, not by convention. See
// engine_test.go TestEngine_CannotExpressExcludedValidators.
//
// Finding transport policy is DEFERRED. A Finding is minimal (a single Blocking
// flag), NOT a severity taxonomy. How a consumer maps Blocking to its transport
// (hook -> block, CI -> fail, manual -> warning) and how opt-out is applied are
// CONSUMER concerns resolved outside the engine ("engine receives effective
// policy / consumer resolves policy source"). Do NOT add Warning/Info/Fatal
// severities here until Q7 matures.

// Finding is a single validation result. Intentionally minimal (Gate D.1):
// a rule id, a human-readable message, and one Blocking flag — no severity enum.
type Finding struct {
	RuleID   string
	Message  string
	Blocking bool
}

// Validate runs the plan_profile.core rules over the normalized plan set and
// returns every finding discovered (Gate D.2: collect-all, never fail-fast).
//
// ctx is part of the stable contract. The current core rules are pure over the
// model set and do not consult ctx; archival rules added later will read
// ctx.ExecutionMode to decide staged-blob vs working-tree transport.
func Validate(_ ValidationContext, plans []NormalizedPlanModel) []Finding {
	var findings []Finding
	findings = append(findings, checkSubPlanFrontmatter(plans)...)
	findings = append(findings, checkUniqueID(plans)...)
	findings = append(findings, checkParentReference(plans)...)
	findings = append(findings, checkArchiveOrder(plans)...)
	return findings
}

// checkSubPlanFrontmatter: a sub plan must declare parent, a non-empty
// sub_plan_reason, and required_for_completion. (validatePlanTreeFrontmatter)
func checkSubPlanFrontmatter(plans []NormalizedPlanModel) []Finding {
	var out []Finding
	for _, p := range plans {
		if p.PlanKind != "sub" {
			continue
		}
		if p.Parent == "" {
			out = append(out, Finding{"plan_tree.frontmatter", fmt.Sprintf("%s: sub-plan missing parent", p.Path), true})
		}
		if p.SubPlanReason == "" {
			out = append(out, Finding{"plan_tree.frontmatter", fmt.Sprintf("%s: sub-plan missing non-empty sub_plan_reason", p.Path), true})
		}
		if p.RequiredForCompletion == nil {
			out = append(out, Finding{"plan_tree.frontmatter", fmt.Sprintf("%s: sub-plan missing required_for_completion", p.Path), true})
		}
	}
	return out
}

// checkUniqueID: every plan id must be unique across the set.
// (validatePlanTreeUniqueID)
func checkUniqueID(plans []NormalizedPlanModel) []Finding {
	var out []Finding
	seen := map[string][]string{}
	for _, p := range plans {
		if p.ID == "" {
			continue
		}
		seen[p.ID] = append(seen[p.ID], p.Path)
	}
	for id, paths := range seen {
		if len(paths) > 1 {
			out = append(out, Finding{"plan_tree.unique_id", fmt.Sprintf("duplicate plan id %q across: %v", id, paths), true})
		}
	}
	return out
}

// checkParentReference: every parent pointer must resolve to an existing id in
// the set (defends against orphan nodes). (validatePlanTreeParentReference)
func checkParentReference(plans []NormalizedPlanModel) []Finding {
	ids := map[string]bool{}
	for _, p := range plans {
		if p.ID != "" {
			ids[p.ID] = true
		}
	}
	var out []Finding
	for _, p := range plans {
		if p.Parent == "" {
			continue
		}
		if !ids[p.Parent] {
			out = append(out, Finding{"plan_tree.parent_reference", fmt.Sprintf("%s: parent %q does not resolve to any plan", p.Path, p.Parent), true})
		}
	}
	return out
}

// checkArchiveOrder: when a main plan is archived, every required sub-plan must
// be completed. (validatePlanTreeArchiveOrder)
func checkArchiveOrder(plans []NormalizedPlanModel) []Finding {
	childrenOf := map[string][]NormalizedPlanModel{}
	for _, p := range plans {
		if p.Parent != "" {
			childrenOf[p.Parent] = append(childrenOf[p.Parent], p)
		}
	}
	var out []Finding
	for _, p := range plans {
		isMain := p.PlanKind == "main" || (p.PlanKind == "" && p.Parent == "")
		if !isMain || p.Location != "archived" {
			continue
		}
		for _, c := range childrenOf[p.ID] {
			if c.RequiredForCompletion != nil && *c.RequiredForCompletion && c.Status != "completed" {
				out = append(out, Finding{"plan_tree.archive_order", fmt.Sprintf("%s: required sub-plan %q is %q, not completed, but parent is archived", p.Path, c.ID, c.Status), true})
			}
		}
	}
	return out
}
