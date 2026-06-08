# Governance Pattern Library — Draft Inventory

> **Status**: draft / observation-stage. **Not** a promoted template.
> Tracked under plan [`2026-06-08-2100-governance-pattern-library-extraction.md`](../../plans/active/2026-06-08-2100-governance-pattern-library-extraction.md).
> Do not link to this file as a normative source. Promotion to
> `governance-pattern-template.md` is gated on the acceptance criteria in the plan.

## Hypothesised shape

```
Observation → Rule → Registry → Projection → Executor → Validation
```

The hypothesis: recent successful governance subsystems all decompose into these 6 steps. The purpose of this draft is to **test that hypothesis against more samples**, not to promote it.

## Sample inventory

| # | Sample | Observation | Rule | Registry | Projection | Executor | Validation |
|---|---|---|---|---|---|---|---|
| 1 | Workflow Activation Engine (plan 2026-05-31-1900, archived) | detector miss | `routing-registry.yaml` | `enforcement-registry` binding | `runtime.db` routes | `detector.go` + PreToolUse gate | regression scenarios |
| 2 | Discovery Bridge (plan 2026-06-06-1700, Phase A landed) | detector-miss-no-fallback | `discovery-bridge.yaml` | `capability_discovery` rule_class | `runtime.discovery.config` | `discovery.go` + advisory injector | `workflow-discovery-bridge-light-v1.yaml` scenario |
| 3 | Sanitization Mechanical Enforcement (plan 2026-06-06-1800, in design) | leak-on-canonical-write | `sanitization-patterns.yaml` + topology | `enforcement-registry` rule_class | `derived_match_tokens` | `sanitization_scan.go` | metadata-derived-fail/pass scenarios |
| 4 | Runtime Index Freshness (commit c5874a8, landed) | source-tree checksum drift | (implicit) | `runtime-index.sqlite` | `sources` table | `nativeRuntimeIndexChecksumsCheck` + commit-msg validator | runtime-index-freshness tests |

### Per-step counts (N=4)

| Step | Filled cells | Notes |
|---|---|---|
| Observation | 4/4 | Universal so far |
| Rule | 3/4 | #4 has "(implicit)" — first warning sign |
| Registry | 4/4 | Universal so far |
| Projection | 4/4 | Universal so far |
| Executor | 4/4 | Universal so far |
| Validation | 4/4 | Universal so far |

**Already a soft signal**: sample #4 has no explicit *Rule* layer. The rule is encoded directly in the executor (`nativeRuntimeIndexChecksumsCheck`) without a yaml/md surface declaring "checksum drift is forbidden." This may indicate the *Rule* step is optional when the invariant is structural (matches a checksum), not editorial (matches a pattern).

## Counter-sample candidates (to inventory next)

These are listed in the plan's Phase 0 checklist. They have NOT been analysed yet. The goal is to find samples that do not cleanly fit the 6-step shape, because a negative case is more informative than a fifth positive.

- [ ] `runtime-trigger-wiring` validator — staged file routing; rule may be declarative, but is there a registry/projection split, or does the validator read source directly?
- [ ] `plan-tree-hierarchy` plan — does plan frontmatter validation have a "Registry" step, or is it pure schema?
- [ ] `bootstrap-contract-yaml-migration` — the migration itself: where does it sit on the 6 steps? Was it Observation → Rule → Projection → Validation, skipping Registry + Executor?
- [ ] Commit-time registry reference consistency (spawn chip `task_a068faa6`) — currently a spawned task; once landed will be the 5th genuine sample.

## What we're looking for in counter-samples

When filling each candidate's row:

1. **Cleanly fits all 6 steps** → adds to confirmation count, but doesn't disprove confirmation bias.
2. **Missing a step but still successful** → strong evidence the step is *optional*, not *universal*. Document which step and why.
3. **Has a step the 6-step shape doesn't have** → strong evidence the hypothesis is wrong or incomplete. Document the new step.
4. **Cannot be mapped at all** → strongest signal that "6-step governance" is not a universal shape, just a recurring one for one class of problem (mechanical enforcement of editorial rules).

## Gate decision (will be filled at Phase 1)

| Criterion | Threshold | Current | Met? |
|---|---|---|---|
| Total samples | ≥ 5 | 4 | ❌ |
| At least one non-fitting sample analysed | ≥ 1 | 0 (sample #4 partial) | ❌ |
| ≥ 3 samples per step | ≥ 3 each | 3-4 each | ✅ |

When all rows are ✅, proceed to Phase 2 of the plan. Until then, do not link this draft from `governance/lifecycle/README.md` and do not reference it as a normative pattern.

## Why this draft matters even if the gate never passes

Even if the 6-step hypothesis is disproved, the analysis produces durable knowledge:

- Concrete evidence of which steps are *truly* universal across governance subsystems
- A list of variants (e.g. "structural invariants skip the Rule layer")
- A record of which subsystems were *forced* into the shape vs. which fit naturally — that itself is a confirmation-bias signal worth capturing in `enforcement/failure-patterns/` if reproduced

← [Back to governance/lifecycle](README.md)
