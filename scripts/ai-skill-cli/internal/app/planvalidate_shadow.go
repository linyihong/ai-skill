package app

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/linyihong/Ai-skill/scripts/ai-skill-cli/internal/planvalidate"
)

// Phase 2.3 (Gate C): non-blocking engine shadow.
//
// planValidateShadowCheck runs the new planvalidate engine ALONGSIDE the legacy
// commit-msg plan-tree validators and reports their divergence as an
// informational Check. It deliberately returns a Check and nothing else, so the
// caller cannot let it influence the exit code: the commit outcome stays
// "legacy only" until divergence converges (Gate C.1).
//
// Equality is by RuleID + Blocking, never message text (Gate C.2). Divergence is
// bucketed into same/missing/extra/transport/context (Gate C.3) so expected,
// benign differences (opt-out transport; execution context) are not reported as
// genuine mismatches.
//
// Scope note: the 4 legacy core validators mirrored here read the working tree
// via scanAllPlanFrontmatter, so the engine (fed from the same scan) should
// match. The known intentional divergence is opt-out: the engine is policy-free
// and still emits a finding the legacy path suppresses via a [skip-*] trailer,
// which is reclassified to the Transport bucket.

// legacyRuleID maps a legacy validator to the engine RuleID it mirrors.
var planTreeLegacyMirror = []struct {
	ruleID    string
	optOut    string
	runLegacy func(text string, staged []string, root string) string
}{
	{"plan_tree.frontmatter", "[skip-plan-tree-frontmatter]", validatePlanTreeFrontmatter},
	{"plan_tree.unique_id", "[skip-plan-tree-unique-id]", validatePlanTreeUniqueID},
	{"plan_tree.parent_reference", "[skip-plan-tree-parent-reference]", validatePlanTreeParentReference},
	{"plan_tree.archive_order", "[skip-plan-tree-archive-order]", validatePlanTreeArchiveOrder},
}

// normalizedPlansFromRoot is the shared loader: it scans <root>/plans/active and
// <root>/plans/archived via the existing scanAllPlanFrontmatter (no new traversal
// abstraction — Phase 2.4 scope (A)), maps each plan to the planvalidate compat
// input, and normalizes it. Both the commit-msg shadow and the `plans validate`
// CLI consumer use this so they feed the engine identical models.
//
// Scope note (2.4 / Q8 boundary): discovery is fixed to plans/active|archived.
// Custom plans dirs, external path conventions, and schema dialects are explicit
// non-goals here; they belong to Q8 / Phase 3, not the CLI transport surface.
// normalizedPlansFromRoot returns the normalized models plus any compat-layer
// findings (e.g. unsupported schema_version). A compat reject is surfaced as a
// BLOCKING finding keyed by the CompatError's ReasonClass so it is deterministic
// and diagnosable end-to-end (Phase 3.2) — never silently degraded.
func normalizedPlansFromRoot(root string) ([]planvalidate.NormalizedPlanModel, []planvalidate.Finding) {
	var models []planvalidate.NormalizedPlanModel
	var compat []planvalidate.Finding
	for _, fm := range scanAllPlanFrontmatter(root) {
		if !fm.HasFrontmatter {
			continue
		}
		raw := planvalidate.RawPlan{
			Path:          fm.Path,
			Location:      planLocation(fm.Path),
			SchemaVersion: fm.SchemaVersion, // end-to-end: loader carries declared version into compat layer
			Fields: map[string]string{
				"id":              fm.ID,
				"plan_kind":       fm.PlanKind,
				"status":          fm.Status,
				"parent":          fm.Parent,
				"sub_plan_reason": fm.SubPlanReason,
			},
		}
		if fm.RequiredForCompletion != nil {
			if *fm.RequiredForCompletion {
				raw.Fields["required_for_completion"] = "true"
			} else {
				raw.Fields["required_for_completion"] = "false"
			}
		}
		model, err := planvalidate.Normalize(raw)
		if err != nil {
			var ce *planvalidate.CompatError
			if errors.As(err, &ce) {
				compat = append(compat, planvalidate.Finding{
					RuleID:   "compat." + ce.ReasonClass,
					Message:  fm.Path + ": " + ce.Error(),
					Blocking: true,
				})
			}
			continue // do not emit a degraded model
		}
		models = append(models, model)
	}
	return models, compat
}

func stagedTouchesPlans(staged []string) bool {
	for _, s := range staged {
		if (strings.HasPrefix(s, "plans/active/") || strings.HasPrefix(s, "plans/archived/")) && strings.HasSuffix(s, ".md") {
			return true
		}
	}
	return false
}

func planValidateShadowCheck(ctx commitMsgCtx) Check {
	if !stagedTouchesPlans(ctx.staged) {
		return Check{Name: "planvalidate_shadow", Status: "skipped", Message: "no plan files staged"}
	}

	// Legacy findings: which mirrored validators fire (presence, not count).
	var legacy []planvalidate.Finding
	hints := planvalidate.CompareHints{OptedOut: map[string]bool{}}
	for _, m := range planTreeLegacyMirror {
		if hasOptOutTrailer(ctx.text, m.optOut) {
			hints.OptedOut[m.ruleID] = true
		}
		if m.runLegacy(ctx.text, ctx.staged, ctx.root) != "" {
			legacy = append(legacy, planvalidate.Finding{RuleID: m.ruleID, Blocking: true})
		}
	}

	// Engine findings: normalize the working-tree plan set and run the engine.
	models, compat := normalizedPlansFromRoot(ctx.root)
	engine := planvalidate.Validate(planvalidate.ValidationContext{
		Root:          ctx.root,
		ExecutionMode: planvalidate.ModeCommit,
	}, models)
	engine = append(engine, compat...)

	d := planvalidate.Compare(legacy, engine, hints)

	msg := fmt.Sprintf("same=%v missing=%v extra=%v transport=%v context=%v",
		compact(d.Same), compact(d.Missing), compact(d.Extra), compact(d.Transport), compact(d.Context))
	if d.Converged() {
		return Check{Name: "planvalidate_shadow", Status: "ok", Message: "engine parity (no genuine divergence) — " + msg}
	}
	return Check{Name: "planvalidate_shadow", Status: "warning", Message: "engine divergence (non-blocking, exit code unaffected) — " + msg}
}

func compact(s []string) string {
	if len(s) == 0 {
		return "-"
	}
	sort.Strings(s)
	return strings.Join(s, ",")
}
