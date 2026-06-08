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

---

## Sequential vs. branching shape

> Resolves Phase 0 checklist item "Check whether the pattern is sequential (always Observation → ... → Validation) or has parallel/optional branches".

**Current evidence (N=4) is sequential, but with three nuances**:

1. **Single optional step observed** — sample #4 (Runtime Index Freshness) has no explicit *Rule* layer; rule is encoded structurally in the executor (`nativeRuntimeIndexChecksumsCheck` matches sha256 against stored checksum). This is the first soft signal that *Rule* is conditionally optional: structural invariants (mechanical equality checks) may not need an editorial rule surface, while editorial invariants (what counts as a leak, what counts as a stale plan reference) always seem to.
2. **Within-step parallelism in Executor** — samples #2 and #3 have executor pairs: Discovery Bridge has `discovery.go` core + advisory injector hook integration as two co-equal entry points; Sanitization Phase 1 (planned) has scanner core + commit-msg validator. The single "Executor" cell hides a sub-DAG. If 5th sample also shows multi-entry executors, the template should explicitly say "Executor = core + integration points" rather than a single function.
3. **No observed re-entry / loop** — every sample so far is one-pass from Observation to Validation. No sample has a "Validation discovers gap → re-enter Rule" cycle yet, though Phase D of Discovery Bridge (3-week empirical) is essentially that. The shape may turn out to have a feedback loop annotation for subsystems with empirical iteration gates.

**Working interpretation**: the shape is **mostly sequential with a single optional step (Rule) and within-step branching at Executor**. Whether this holds against 5th + 6th samples will determine if the template captures the sequence as-is or with explicit branch / optional notation.

---

## Parallel observation: Reference Integrity family

> Surfaced in the same 2026-06-08 conversation that produced the 6-step hypothesis. Captured here in the **same incubator** because (a) it is also a cross-cutting governance shape across the same subsystems and (b) the gate / falsification discipline should be applied to both before either is promoted. They are **two independent families**, not two views of one.

### Shape

```
Knowledge Surface
    ↓
Reference (one surface points to another)
    ↓
Drift between surface and reference
    ↓
Phase that surfaces the drift (often: phase ≠ the one that caused it)
```

### Distinguishing from the 6-step shape

| Axis | 6-step (above) | Reference Integrity (this section) |
|---|---|---|
| Lifecycle phase | **Build-time** (how a new mechanical rule gets onboarded) | **Maintenance-time** (how existing surfaces stay consistent) |
| Trigger | "We need a new governance capability" | "These two surfaces should agree but might not" |
| Output | A new validator + projection + executor chain | A new drift detector / consistency check |
| Failure mode if absent | The rule cannot be enforced (rule-without-executor) | Two surfaces silently disagree (validation-coverage-gap, today's 2026-06-06 incident) |

These are independent families that happen to operate on the same governance subsystems. Promoting one does not imply or require promoting the other.

### Reference Integrity sample inventory (N=4, observation-stage)

| # | Subsystem | Source Surface | Target Surface | Drift Surface | Phase that exposes drift |
|---|---|---|---|---|---|
| 1 | Sanitization | `enforcement/sanitization.md` (prose rule) | scanner / allowlist executor | rule ↔ executor (rule changes, executor does not pick up new private token) | pre-commit / commit-msg |
| 2 | Workflow Activation | detector (`detector.go`) | Discovery fallback (`discovery.go`) | detector miss ↔ fallback presence (detector miss with no Discovery wiring) | PreToolUse runtime |
| 3 | Runtime Index | source files (e.g. `CORE_BOOTSTRAP.md`) | stored sha256 in `runtime-index.sqlite` | content ↔ checksum | `runtime validate` / commit-msg `validateRuntimeIndexFreshness` |
| 4 | Registry Reference | plan path in `plans/active/*.md` | `child_plan` pointer in `enforcement-registry.yaml` | path rename ↔ registry pointer | `runtime compile` (only triggered when workflow `paths:` filter hits, hence the silent-drift incident) |

### Reference Integrity acceptance gate (mirrors 6-step gate)

| Criterion | Threshold | Current | Met? |
|---|---|---|---|
| Total samples | ≥ 5 | 4 | ❌ |
| At least one non-fitting sample analysed | ≥ 1 | 0 | ❌ |
| Drift-surface variety | ≥ 3 distinct drift surfaces | 4 distinct (rule↔executor, detector↔fallback, content↔checksum, path↔pointer) | ✅ |

### Pre-emptive falsification questions

To avoid retrofitting future samples into this shape, the candidates to test against:

- [ ] **Conversation goal ledger** (`enforcement/conversation-goal-ledger.md`): source surface = `.agent-goals/` files; target surface = ... what? If nothing, this is the first sample that proves Reference Integrity is NOT universal across governance subsystems — drift requires two surfaces.
- [ ] **Cognitive Mode reporting** (per-turn obligation + commit-msg validator): is the per-turn report a *reference* to the canonical YAML schema, or just a transient declaration? If transient, Reference Integrity does not apply, and the family is narrower than the 6-step family.
- [ ] **Discovery Bridge advisory output** (this session's Phase A): the advisory text is rendered from `runtime.discovery.config` + per-task signal; no persistent target surface. Reference Integrity probably does not apply.

If 2+ of these confirm Reference Integrity does NOT apply, the family is narrower than the 6-step family — that is itself a useful learning (it means the 6-step shape covers more governance subsystems than the Reference Integrity shape, even though both surfaced at the same time).

### Cross-link forward

Both families may eventually become entries in `governance/lifecycle/governance-pattern-library/` (folder, not single file) if they both pass their gates. They will be siblings, not parent/child. The failure patterns that motivate them (the anti-pattern side) already live in `enforcement/failure-patterns/`:

- 6-step shape ↔ `rule-without-executor.md` + `validation-coverage-gap-executor-placement.md`
- Reference Integrity ↔ `markdown-yaml-sync-drift.md` + `validation-coverage-gap-executor-placement.md` (the latter spans both families)

← [Back to governance/lifecycle](README.md)
