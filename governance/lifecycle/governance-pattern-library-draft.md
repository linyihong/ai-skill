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
| 5 | Validation Scenario Governance Executor — F19 (plan 2026-06-01-0100, archived 2026-06-12) | declared coverage evidence unverified / dangling `coverage_evidence` refs (`coverage-evidence-dangling-reference.md`) | `coverage_evidence` schema in `enforcement-registry.yaml` (`validation_scenarios[]` / `regression_scenarios[]` / `coverage_target_pct`) | `enforcement-registry` F19 `validation_scenario_governance` rule_class + `executors[]` block | **(none — executor reads the registry yaml + scenario corpus directly at compile time)** | `scenario_lint.go` `LintValidationScenarios` wired into `runtime compile` | `scenario_lint_test.go` (5 tests, ≥1 fail+1 pass per check) + `scenario-lint-dangling-coverage-ref-regression-v1.yaml` |
| 6 | Plan-Tree Hierarchy Governance (plan 2026-06-02-1200, mechanical) | plan-tree drift (broken parent pointer, id collision, archive mis-order) | frontmatter-schema **policy** (`01-frontmatter-schema.md`) — *but bundles one **structural** sub-invariant: `validatePlanTreeUniqueID` (id collision = `A==B`)* | `enforcement-registry` `plan_tree_governance` rule_class + 5 `executors[]` | **(none — 5 validators read staged/worktree plan files directly via `os.ReadFile`)** | `plan_tree.go` — **5** commit-msg validators (frontmatter / archive-order / parent-ref / unique-id / folder-convention) | `plan_tree_test.go` Go unit tests *(scenario-yaml sub-form `scenario_exists: pending`)* |
| 7 | Runtime Trigger Wiring (`runtime_trigger_wiring`, mechanical) | orphan runtime surface (new route / target_key with no consumer) | **`system-upgrade-governance.yaml §define_runtime_trigger_flow` (policy — enumerates the valid wiring forms)** — Rule **present**, *not* structural | `enforcement-registry` `runtime_trigger_wiring` rule_class | **(none — `routeWiredInTree`/`targetKeyConsumedInTree` use `os.ReadFile` + `filepath.Walk` + `stagedDiff`, direct)** | `validateRuntimeTriggerWiring` (commit-msg validator) | `runtime_trigger_wiring_test.go` **+** `orphan-routing-entry-v1.yaml` / `orphan-projection-target-key-v1.yaml` (both sub-forms) |

### Per-step counts (N=7)

| Step | Filled cells | Notes |
|---|---|---|
| Observation | 7/7 | Universal so far |
| Rule | 6/7 | **Only #4 empty.** #7 was the predicate-1 candidate but turned out to *have* a policy Rule — see "predicate 1 NOT closed" below |
| Registry | 7/7 | Universal so far |
| Projection | 4/7 | **#5, #6, #7 empty** — Projection-optional now has **three** confirming direct-consumption cases |
| Executor | 7/7 | Universal so far |
| Validation | 7/7 | Universal so far; #7 is complete in *both* sub-forms (Go test + 2 scenario yaml) |

**The real result: the shape is not a 6-step invariant — it is a 4-step invariant core + 2 conditional steps.**

The acceptance gate was set up to test whether the 6-step chain is an invariant. The evidence (now N=7) answers that directly: **it is not** an invariant — but it is **not wrong** either. Four steps are invariant; two are conditional, each governed by its own falsifiable predicate.

```
Original hypothesis (6-step invariant)        Result (N=7)

Observation                                   Observation  ┐
   → Rule                                     Registry     │ invariant core
   → Registry                                 Executor     │ (7/7 each)
   → Projection                               Validation   ┘
   → Executor
   → Validation                               Rule        — conditional (6/7; near-universal)
                                              Projection  — conditional (4/7)
```

**Conditional predicate 1 — Rule** (the hole at sample #4, mature wording):

> *Rule is conditionally required for **policy-derived** governance. It is optional when the invariant is **structural**.*

A *structural* invariant needs no authored rule surface because the violation IS the structural inequality. "Checksum must match content" requires no `rule: { checksum_must_match: true }` — `A != B` is already the violation. A *policy-derived* invariant (what counts as a leak, what counts as a stale plan reference) is human-authored and therefore always carries a Rule surface. Sample #4 (`nativeRuntimeIndexChecksumsCheck`) is the clean structural case.

*Granularity refinement (from sample #6)*: the predicate operates at **rule_class granularity**, not per-validator. `plan_tree_governance` carries a Rule surface (the frontmatter-schema policy) because most of its five validators are policy-derived (parent-pointer semantics, folder convention) — yet one of them, `validatePlanTreeUniqueID`, is purely *structural* (id collision = `A==B`, no policy needed) and rides on the shared rule surface anyway. So the precise predicate is: **a rule_class carries a Rule surface iff it contains *any* policy-derived invariant; structural sub-invariants bundled into a policy rule_class inherit it.**

*Predicate 1 NOT closed — and the failed attempt is itself the finding (from sample #7)*: the strongest open test was "find an *entirely structural* rule_class with no Rule surface," and `runtime_trigger_wiring` was the candidate. It **failed to close it** — but informatively. Its invariant ("a new route/target_key must have a consumer") *looks* structural (reference resolvability), yet it carries a real policy Rule surface, `system-upgrade-governance.yaml §define_runtime_trigger_flow`, because **"what counts as a valid consumer" is human-authored** (the enumerated forms: discovery signal / Go consumer / `manual_activation` annotation). Reference-resolvability is structural in *mechanism* but policy in *what-relation-must-hold*. The lesson: a great many invariants that feel structural actually smuggle in a policy definition of the valid relation. **Rule-optional is therefore rare — a narrow carve-out, not a common variant.** After 7 samples it is attested by exactly **one** clean case (#4, checksum equality), where the violation is a pure structural (in)equality needing no authored relation. The template should present Rule as **near-universal**, with the structural exemption as a deliberately-justified exception (criterion: *the violation is a pure structural (in)equality and no human authorship defines the valid relation* — checksum match, raw id collision), not as a co-equal "this step is often optional."

**Conditional predicate 2 — Projection** (the hole at sample #5 / F19):

> *Projection is conditionally required for **indirect-consumption** executors. It is optional when the executor can consume the authoritative registry/source surface **directly**.*

This is deliberately framed at the contract level, **not** the implementation level. F19's executor happens to call `loadRegistrySnapshot` + `os.ReadFile` today; tomorrow it could be `yaml.Unmarshal`, the day after `registry.Load()`. The essence is invariant under all three: the executor reads the *authoritative source* and needs no intermediate surface. Projection only materialises when consumption must be **indirect** — i.e. the executor needs one of:

- **high-frequency reads** (a hot path that cannot re-parse the source each call), or
- **pre-digested data** (a derived form the source does not expose), or
- **a unified query surface** (one shape over many heterogeneous sources).

| # | Executor consumes source… | Projection present? | Surface |
|---|---|---|---|
| 1 Workflow | indirectly (high-frequency: every non-Read PreToolUse call) | ✅ | `runtime.db` routes |
| 2 Discovery | indirectly (pre-digested capability map) | ✅ | `runtime.discovery.config` |
| 3 Sanitization | indirectly (pre-digested match tokens) | ✅ | `derived_match_tokens` |
| 4 Runtime Index | boundary — Registry *is* the db, so source-read and projection-read collapse into the `sources`-table read | ✅ (degenerate) | `sources` table |
| 5 F19 | **directly** (reads the authoritative registry + scenario yaml) | ❌ | — |
| 6 Plan-Tree | **directly** (5 commit-msg validators `os.ReadFile` the staged/worktree plans) | ❌ | — |
| 7 Runtime Trigger Wiring | **directly** (`os.ReadFile` + `filepath.Walk` tree scan + `stagedDiff`) | ❌ | — |

So Projection is reclassified from a *governance-process step* to an **executor optimization layer**: present only when indirect consumption is required, absent when the executor can read the authoritative surface directly. Sample #4 is the boundary case — its registry already *is* a database, so the two collapse; #5 (F19), #6 (Plan-Tree), and #7 (Runtime Trigger Wiring) are three clean cases where the source is plain yaml/markdown/tree and the executor reads it directly, leaving Projection genuinely empty (not skipped for convenience, not documentation debt). **Predicate 2 now has three independent confirming samples (#5, #6, #7)**, spanning two executor kinds (compile-time lint + commit-msg validator) and three source shapes (registry yaml, plan frontmatter, routing-registry + tree walk). It is the better-attested of the two predicates by a wide margin.

**Why two holes on two different steps is the stronger result:** Rule and Projection are governed by *independent* predicates (policy-vs-structural; indirect-vs-direct consumption). Two holes on the *same* step would only have weakened that step; two holes on two steps, each with its own falsifiable condition, is what distinguishes "invariant core + conditional stages" from "6-step with some noise."

**Validation sub-form nuance (from sample #6, reinforced by #7):** Validation is 7/7 (universal core), but #6 shows the step has *two interchangeable sub-forms* — **Go unit tests** (`plan_tree_test.go`) and **scenario-yaml** (`validation/scenarios/`). `plan_tree_governance` satisfies Validation via Go tests while its registry `scenario_exists` is still `pending`; the two are not the same axis. #7 (Runtime Trigger Wiring) is the *complete* case — it has **both** sub-forms (`runtime_trigger_wiring_test.go` + `orphan-routing-entry-v1.yaml` / `orphan-projection-target-key-v1.yaml`), confirming the sub-forms are additive, not exclusive. The template's Validation step should read "executable proof (unit test **and/or** scenario yaml) that the executor detects the violation," and a subsystem can be Validation-complete in one sub-form while the other is independently tracked. "Validation present" (core invariant) and "scenario-yaml coverage complete" (a separate governance metric, the very thing F19 itself enforces) must not be conflated.

**Cross-family note (sample #7 ↔ Reference Integrity):** `runtime_trigger_wiring`'s invariant — "a route/target_key reference must resolve to a consumer" — is also a candidate **Reference Integrity** sample (source surface = routing-registry route / runtime target_key; target surface = a consumer; drift = orphan). It surfaced here while testing the 6-step Rule predicate; it is logged as a Reference Integrity candidate in that family's section rather than expanding this one.

## Counter-sample candidates (to inventory next)

These are listed in the plan's Phase 0 checklist. They have NOT been analysed yet. The goal is to find samples that do not cleanly fit the 6-step shape, because a negative case is more informative than a fifth positive.

- [x] `runtime-trigger-wiring` validator — **inventoried as sample #7** (2026-06-12). Reads source directly (no projection). Has a Registry step (`runtime_trigger_wiring` rule_class) AND a policy Rule (`system-upgrade-governance.yaml §define_runtime_trigger_flow`). Findings: did **not** close predicate 1 (it is policy-derived, not structural — "what counts as wired" is human-authored); instead it became the **3rd** Projection-optional confirmation and reframed Rule-optional as a *rare* carve-out (still only #4 clean). Also a Reference Integrity candidate. See sample row 7 + predicate-1 finding.
- [x] `plan-tree-hierarchy` plan — **inventoried as sample #6** (2026-06-12). Has a Registry step (`plan_tree_governance` rule_class), not pure schema. Findings: confirms Projection-optional (2nd case), refines Rule predicate to rule_class granularity, surfaces Validation sub-form nuance. See sample inventory row 6 + predicate sections.
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
| Total samples | ≥ 5 | 7 (F19 #5 + Plan-Tree #6 + Runtime-Trigger-Wiring #7) | ✅ |
| At least one non-fitting sample analysed | ≥ 1 | 4 (#4 Rule-missing; #5/#6/#7 Projection-missing) | ✅ |
| ≥ 3 samples per step | ≥ 3 each | 4-7 each | ✅ |

**All three gate criteria are met, and the Phase 1 decision is closed.** The gate's own decision tree resolves deterministically:

> - If 5+ samples AND ≥3 per step AND 6-step shape **consistent** → Phase 2 as-is
> - If samples found but **step coverage uneven** → **revise template shape (some steps optional, document variants)** ← resolved here
> - If samples diverge significantly → abandon promotion

Coverage is **uneven but principled**, so the middle branch fires, and the revision is already performed: the promotable artifact is a **4-step invariant core** (Observation → Registry → Executor → Validation), with **Rule** conditionally required for policy-derived governance and **Projection** conditionally required for indirect-consumption executors. That makes Phase 1 a *closed decision*, not a pending one — what remains is Phase 2 **wording refinement**, not a branch choice.

This is a stronger outcome than a clean 5th confirmation would have been: a clean fit would have left confirmation bias unrefuted, whereas reducing the hypothesis to a 4-step core + 2 falsifiable predicates is genuine new knowledge — it tells the eventual template *when* a step may be omitted instead of demanding all six unconditionally.

**Phase 2 precondition — partially discharged (2026-06-12).** The recommendation was to find a 6th sample chosen to *challenge* the two predicates rather than tally. Sample #6 (Plan-Tree Hierarchy Governance) did that:

- **Projection-optional — strengthened.** #6 is a second, independent direct-consumption-no-projection case, from a different executor kind (commit-msg validator vs F19's compile-time lint) over a different source (plan frontmatter vs registry yaml). Predicate 2 now rests on two clean samples, not one.
- **Rule-optional — refined, not yet independently re-confirmed.** #6's `validatePlanTreeUniqueID` is structural but rides on a policy rule_class's rule surface, which sharpened the predicate to **rule_class granularity** but did *not* add a second clean standalone structural-no-rule case. Predicate 1 still rests on #4 alone.
- **New nuance — Validation sub-form.** #6 exposed that Validation has interchangeable Go-test / scenario-yaml sub-forms tracked on separate axes.

**Remaining open test before freezing template wording:** an *entirely structural* rule_class confirmed to carry **no** Rule surface (clean second case for predicate 1). Candidate left: `runtime-trigger-wiring`. Until that lands and the template wording is frozen, do not link this draft from `governance/lifecycle/README.md` and do not reference it as a normative pattern.

## Why this draft matters even if the gate never passes

Even if the 6-step hypothesis is disproved, the analysis produces durable knowledge:

- Concrete evidence of which steps are *truly* universal across governance subsystems
- A list of variants (e.g. "structural invariants skip the Rule layer")
- A record of which subsystems were *forced* into the shape vs. which fit naturally — that itself is a confirmation-bias signal worth capturing in `enforcement/failure-patterns/` if reproduced

---

## Sequential vs. branching shape

> Resolves Phase 0 checklist item "Check whether the pattern is sequential (always Observation → ... → Validation) or has parallel/optional branches".

**Current evidence (N=6) is sequential, but with five nuances**:

1. **Two optional steps observed, on two different steps under two independent conditions** — sample #4 (Runtime Index Freshness) has no explicit *Rule* layer (rule encoded structurally in `nativeRuntimeIndexChecksumsCheck`, a sha256 equality check); samples #5 (F19) and #6 (Plan-Tree) have no *Projection* layer (executor consumes the authoritative source directly). The governing predicates are independent: *Rule* is conditionally required for **policy-derived** governance and optional for **structural** invariants (at rule_class granularity — see Rule predicate granularity note); *Projection* is conditionally required for **indirect-consumption** executors and optional when the executor consumes the authoritative source **directly** (predicate stated at contract level, not implementation — `os.ReadFile` today, `registry.Load()` tomorrow, same essence). See §"Per-step counts" for the direct-vs-indirect consumption table.
2. **Within-step parallelism in Executor, with varying cardinality** — the single "Executor" cell hides a sub-DAG whose *cardinality* varies widely: 1 placement (#5 F19, `runtime compile`-only), 2 (#2 Discovery `discovery.go` + advisory injector; #3 Sanitization scanner + commit-msg), up to **5** (#6 Plan-Tree: frontmatter / archive-order / parent-ref / unique-id / folder-convention commit-msg validators). The template should say "Executor = core + integration points (count is a per-subsystem decision, see `validation-coverage-gap-executor-placement.md`)".
3. **No observed re-entry / loop** — every sample so far is one-pass from Observation to Validation. No sample has a "Validation discovers gap → re-enter Rule" cycle yet, though Phase D of Discovery Bridge (3-week empirical) is essentially that. The shape may turn out to have a feedback loop annotation for subsystems with empirical iteration gates.
4. **Observation can be a named failure pattern, not just an ad-hoc note** — F19's Observation cell points at a *promoted* failure pattern (`coverage-evidence-dangling-reference.md`), where samples #1–#4 had inline observations. This suggests the mature form of the Observation step is "a failure pattern in `enforcement/failure-patterns/`," and the 6-step (positive) family is the constructive dual of the failure-pattern (anti-pattern) family — consistent with the plan's "this captures the positive shape the anti-pattern violates against."
5. **Validation has interchangeable sub-forms** — sample #6 satisfies Validation via Go unit tests (`plan_tree_test.go`) while its scenario-yaml coverage is still `pending`. Validation-the-core-step ("executable proof the executor detects the violation") is distinct from scenario-yaml-coverage (a separate governance metric). The two must not be conflated; a subsystem can be Validation-complete in the test sub-form with scenario-yaml lagging.

**Working interpretation (updated N=7)**: the shape is **a 4-step invariant core (Observation → Registry → Executor → Validation), with Rule near-universal (one rare structural exemption) and Projection conditionally absent for direct-consumption executors, plus within-step branching at Executor and interchangeable sub-forms at Validation**. The two conditional steps are *asymmetric* in how well-attested they are: **Projection-optional has three independent confirmations (#5/#6/#7)** and is solid; **Rule-optional has exactly one clean case (#4)** after a deliberate attempt (#7) to find a second failed by revealing the candidate was policy-derived. So the template should treat Projection as a genuine conditional stage and Rule as near-universal with a narrow, criterion-gated structural carve-out. Further sampling for predicate 1 should target pure structural (in)equality invariants (checksum / raw uniqueness), since reference-resolvability-style invariants keep turning out to be policy-backed.

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

### Reference Integrity sample inventory (N=5, observation-stage)

| # | Subsystem | Source Surface | Target Surface | Drift Surface | Phase that exposes drift |
|---|---|---|---|---|---|
| 1 | Sanitization | `enforcement/sanitization.md` (prose rule) | scanner / allowlist executor | rule ↔ executor (rule changes, executor does not pick up new private token) | pre-commit / commit-msg |
| 2 | Workflow Activation | detector (`detector.go`) | Discovery fallback (`discovery.go`) | detector miss ↔ fallback presence (detector miss with no Discovery wiring) | PreToolUse runtime |
| 3 | Runtime Index | source files (e.g. `CORE_BOOTSTRAP.md`) | stored sha256 in `runtime-index.sqlite` | content ↔ checksum | `runtime validate` / commit-msg `validateRuntimeIndexFreshness` |
| 4 | Registry Reference | plan path in `plans/active/*.md` | `child_plan` pointer in `enforcement-registry.yaml` | path rename ↔ registry pointer | `runtime compile` (only triggered when workflow `paths:` filter hits, hence the silent-drift incident) |
| 5 | Runtime Trigger Wiring | new route / target_key in `routing-registry.yaml` / `runtime/*.yaml` | a consumer (discovery signal / Go consumer / routing-registry reference) | reference ↔ consumer (orphan: reference added, consumer absent) | commit-msg `validateRuntimeTriggerWiring` |

> #5 surfaced 2026-06-12 while testing the 6-step Rule predicate (it is 6-step sample #7). Cross-family appearance is consistent with #2 (also a 6-step sample). It does **not** advance this family past its gate — the non-fitting criterion below is still unmet.

### Reference Integrity acceptance gate (mirrors 6-step gate)

| Criterion | Threshold | Current | Met? |
|---|---|---|---|
| Total samples | ≥ 5 | 5 (Runtime Trigger Wiring added 2026-06-12) | ✅ |
| At least one non-fitting sample analysed | ≥ 1 | 0 | ❌ |
| Drift-surface variety | ≥ 3 distinct drift surfaces | 5 distinct (rule↔executor, detector↔fallback, content↔checksum, path↔pointer, reference↔consumer) | ✅ |

**Gate still NOT met** despite N≥5: no non-fitting sample has been analysed for this family (unlike the 6-step family, where #4/#5/#6/#7 produced the variant structure). Reference Integrity needs a deliberately-chosen non-fitting case — e.g. one of the pre-emptive falsification candidates below — before its shape can be revised and promoted. Adding a 5th *confirming* sample does not substitute for that.

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

---

## Parallel observation: Failure Authority family

> Surfaced in the 2026-06-10 Phase 1D review of the sanitization plan
> (2026-06-06-1800). Captured here in the **same incubator** under the same
> gate / falsification discipline as the other two families. It is a **third
> independent family**, not a view of either above.

### The question it asks

The 6-step family asks *"how does a new mechanical rule get onboarded?"*.
Reference Integrity asks *"do these two surfaces still agree?"*. This family
asks a different question that kept recurring across Discovery, Runtime Index,
and Sanitization:

> **Who has standing to block compile/commit — and on the authority of which source?**

The failure being named is the conflation of two separate things:

```
Metadata Presence  ⇒  Compile Authority
```

i.e. "this file exists and is invalid" silently becomes "therefore the whole
pipeline halts" — regardless of whether the file is an authoritative input to
`runtime.db` or an ephemeral, git-ignored, project-local scratch file.

### The invariant (observation-stage; not yet normative)

> **Failure Authority** — Only inputs originating from *compile-authoritative*
> sources may block runtime compilation (or a commit gate). *Non-authoritative*
> sources may emit warnings, but must not prevent runtime generation.

Canonical decision shape (the recurring control flow, abstracted from any one subsystem):

```
event failure (invalid / drift / miss)
        │
        ▼
classify authority of the SOURCE
        │
   ┌────┴────┐
authoritative   non-authoritative
   │                │
 FAIL             WARN
```

Illustrative structured form (NOT a projected/canonical surface yet — promotion
to a machine-readable invariant is gated; this is the shape it would take):

```yaml
# illustrative — lives in the draft until the gate passes
invariant:
  name: failure-authority
  classify: source            # not: presence, not: validity
  rule:
    authoritative:     fail   # may block compile / commit
    non_authoritative: warn   # may warn, must not block
```

The pivot is *classify the source's authority, not the input's validity*. A file
can be 100% invalid and still have no standing to halt the pipeline.

### Conceptual core: Standing (the deepest framing)

The family's real subject is not "failure" — it is **standing**, in the legal
sense: *who has the right to halt a proceeding.* The operative distinction:

```
Validity   ≠   Authority (Standing)
```

> Many things can be wrong. **Not every wrong thing has standing to stop the
> process.** A non-authoritative source raising a 100%-correct error still has no
> right to block compile — its remedy is a warning, not a halt.

Finding A is the proof case: a metadata file being *100% invalid* does **not**
entail it has *standing to halt compile*. "Failure Authority" is the mechanism;
**Standing is the principle.** Every executor below must answer the standing
question (*does this subject have the right to block?*) before the validity
question (*is this subject correct?*).

### Authority is per-SUBJECT, not per-path

A tempting shortcut is to make authority a function of file path (topology). But
path is only the signal for *one* kind of subject. The three samples already use
**three different authority sources**:

| Sample | Subject kind | Authority signal (NOT universally "path") |
|---|---|---|
| Workflow / Discovery | `route` / `discovery-provider` | route-lock + in-repo scope; provider is advisory-by-contract |
| Runtime Index | `runtime-index-row` | presence of a `sources` row |
| Project Metadata (Finding A) | `metadata-file` | `repository-topology.yaml` → `shared_layer` / `owner` |

Topology answers the standing question *for the metadata-file kind only*. Welding
`authority := path-classification` would make the classifier silently wrong for
the other two kinds — and would degrade into `if runtimeIndexRow… if discovery…`
branching that no longer matches its own name. Hence the contract below is over a
**subject**, not a path.

Corollary that motivates promoting this rather than just documenting it: if, for
the metadata-file kind, `shared_layer: true` and `shared_layer: false` both
ultimately `compile fail`, the topology v2 classification work (owner / purpose /
shared_layer) loses most of its governance value at the one moment it should
matter most.

### Sample inventory (N=4, observation-stage)

| # | Subsystem | Authoritative source (may block) | Non-authoritative source (must not block) | Invariant currently holds? |
|---|---|---|---|---|
| 1 | Workflow Activation / Discovery Bridge | a single locked `active_route` inside the repo | detector miss / multi-route conflict / routing registry unresolvable / running outside the repo → gate **fails open**; Discovery fallback is advisory-only and never blocks | ✅ holds — non-authoritative failure must not halt the system |
| 2 | Runtime Index Freshness | files with a `sources` row in `runtime-index.sqlite` | files *without* a source row ("outside this validator's freshness scope" — core-bootstrap.yaml runtime_index_freshness rationale) | ✅ holds — not every source row has standing to fail compile; only registered ones |
| 3 | Project Metadata Compile (Sanitization) | malformed shared-layer `.ai-skill-project.yaml` | malformed `.agent-goals/…` (`shared_layer:false`) project-local metadata | ✅ holds (since 2026-06-10) — now classifier-scoped via `ClassifyFailureAuthority` (Executor #1) |
| 4 | Authority Classifier itself | — | — | ✅ the shared `ClassifyFailureAuthority(subject)` impl all three above will route through |

Sample #3 *was* the live counter-instance (the Phase 1D hard-fail blocked compile
repo-wide, including non-authoritative `.agent-goals/`); it is now brought into
line — `compileProjectMetadataDerived` builds a `metadata-file` subject from
topology and only hard-fails when the subject is authoritative. Samples #1 and #2
already honoured the invariant, so #3 was the outlier corrected, not a new
behaviour invented.

### Acceptance gate (mirrors the other two families)

| Criterion | Threshold | Current | Met? |
|---|---|---|---|
| Total samples | ≥ 5 | 4 | ❌ |
| At least one non-fitting sample analysed | ≥ 1 | 1 (sample #3 was a violation — analysed and now corrected) | ✅ |
| Distinct authority signals | ≥ 3 | 3 (source-row / route-lock+repo-scope / shared_layer) | ✅ |

### Pre-emptive falsification questions

- [ ] **Commit-msg validators** (e.g. cognitive-mode block, plan-status-sync): do they ever block on staged files outside an authoritative scope, or is every staged file authoritative-by-definition (it's being committed)? If the latter, commit gates may be a degenerate case where Class B does not exist.
- [ ] **Bootstrap receipt gate**: is "no receipt" a Failure Authority decision (the session is authoritative) or a different family entirely?
- [ ] Does Runtime Index sample #1 *really* warn (vs silently ignore) non-authoritative drift? If it silently ignores rather than warns, the invariant's "warn only" half is unproven and may need softening to "must not block (warning optional)".

### Authority Classification Contract (the missing layer — subject-based, docs-only)

Inserting a layer the earlier sketch skipped. The chain is **not**
`invariant → classifier → executor`; it is:

```
Failure Authority invariant          (principle: Standing; Validity ≠ Authority)
        ↓
Authority Classification Contract     (THIS section — subject-based, language-neutral)
        ↓
Authority Classifier                  (one implementation of the contract)
        ↓
Executor                              (a caller that obeys the classifier's verdict)
```

The contract is what lets *five future executors share one definition of standing*
instead of each re-deciding it. It is written here as a specification first; no Go
is implied yet.

**Subject, not path.** Authority is resolved over an abstract *subject*:

```
AuthoritySubject:
  kind:  discovery-provider | runtime-index-row | metadata-file | generated-surface | …
  # kind-specific attributes the classifier may read:
  path:         string   # when the kind is path-bound (metadata-file)
  owner:        string   # e.g. project-local
  shared_layer: bool     # from repository-topology (metadata-file kind)
  registered:   bool     # e.g. has a runtime-index `sources` row
  advisory:     bool     # e.g. discovery-provider is advisory-by-contract
```

**Standing rule, keyed by signal (not by path).** A subject is *authoritative*
(may block) iff any authoritative signal holds and no override demotes it:

```yaml
# illustrative — draft contract, not a projected surface
authority_classification:
  authoritative:        # has standing to FAIL the process
    - shared_layer: true
    - registered: true            # runtime-index sources row, tracked runtime surface
    - tracked_runtime_surface
  non_authoritative:    # may WARN only, must not block
    - shared_layer: false
    - owner: project-local
    - advisory: true              # discovery-provider, advisory hooks
    - ephemeral / untracked
```

**Contract API shape (illustrative; the deliberate non-goal is `isCompileAuthoritative(path)`):**

```go
type AuthorityLevel int
const (
    NonAuthoritative AuthorityLevel = iota   // may warn, must not block
    Authoritative                            // may fail
)

// The contract surface every executor calls:
func ClassifyFailureAuthority(subject AuthoritySubject) AuthorityLevel
```

`isCompileAuthoritative(path)` is then **not** the contract — it is at most a thin
convenience wrapper that constructs a `metadata-file` subject and delegates to
`ClassifyFailureAuthority`. If it ever grows `if runtimeIndexRow… if discovery…`
branches, that is the signal the path-shaped wrapper was mistaken for the contract.

**Resolution precedence & safe defaults** (pinned by the first implementation,
`scripts/ai-skill-cli/internal/app/failure_authority.go`):

- **Demotion wins.** A non-authoritative signal (`owner: project-local`,
  `shared_layer: false`) overrides any authoritative signal — a project-local
  file has no standing even if it sits in a shared subtree.
- **`shared_layer` is tri-state** (`unknown | true | false`), not a bool. On a
  topology **miss** (`unknown`) a metadata-file *keeps* standing (fail-safe
  toward protection): silently demoting a possibly-shared file would let a real
  leak pass as a mere warning. A bare bool would default that case to the
  dangerous direction.
- **Unknown kind must earn standing.** A subject kind the contract has not
  explicitly granted standing defaults to *non-authoritative*. Adding a new
  *blocking* kind requires an explicit case — a grant is always conscious, never
  accidental.

> Build constraint (honoured): the **contract was specified before any classifier
> code**, so the first implementation (`ClassifyFailureAuthority`, landed with
> unit tests across metadata-file / runtime-index-row / discovery-provider /
> unknown kinds) conformed to a shared definition of standing rather than
> retrofitting one from Finding A's path check.

### Dependency inversion — Finding A is an executor, not the cause

The direction of derivation matters more than the fix itself. The wrong order
treats the classifier as a patch for one bug:

```
Finding A  →  Authority Classifier        (bug-driven; classifier is a special-case)
```

The right order derives everything from the invariant, with the Classification
Contract as the shared layer, and Finding A as the invariant's *first executor*:

```
Failure Authority invariant
        ↓
Authority Classification Contract       (subject-based; shared by all executors)
        ↓
Authority Classifier                    ClassifyFailureAuthority(subject)
        ↓
Executor #1: Project Metadata Compile   (the Finding A fix, metadata-file subject)
Executor #2..n: Discovery / Runtime Index / future surfaces reuse the SAME contract
```

So the materialised fix is not "special-case `.agent-goals/`", it is a call into
the shared contract (via a metadata-file subject):

```go
subject := AuthoritySubject{Kind: "metadata-file", Path: rel, SharedLayer: shared, Owner: owner}
if ClassifyFailureAuthority(subject) == NonAuthoritative {
    warn(...)
    continue          // non-authoritative: may warn, must not block
}
return error          // authoritative: may fail
```

Consequence: Discovery, Runtime Index, and any future compiler surface inherit
the *same* classifier rather than re-deciding "who may block" ad hoc. That is the
ROI argument for building the invariant first — the alternative leaves three
near-identical authority decisions scattered and divergent.

**Sequencing (agreed build order):**

1. ✅ Complete this family observation (the invariant + Standing core + 3 samples).
2. ✅ Specify the **Authority Classification Contract** (← above; subject-based,
   docs-only, no Go). The shared definition of standing.
3. ✅ Implement the Authority Classifier as the contract's first implementation
   (`ClassifyFailureAuthority(subject)` in `scripts/ai-skill-cli/internal/app/failure_authority.go`),
   *from the contract*, not from Finding A. Unit-tested across all four subject
   kinds — proving it is genuinely multi-kind, not a path-welded special-case.
   No caller wired yet.
4. ✅ Land Finding A as Executor #1 (`compileProjectMetadataDerived` builds a
   `metadata-file` subject from topology and calls `ClassifyFailureAuthority`;
   non-authoritative → warn+skip, authoritative → hard-fail). Sample #3 moved
   ❌→✅; family now N=4 toward the N≥5 gate.

See plan 2026-06-06-1800 §"Phase 1D review — Finding A" for the deferred fix and
plan 2026-06-08-2100 for the incubator gate.

### Cross-link forward

If promoted, this becomes a third sibling in
`governance/lifecycle/governance-pattern-library/` (folder), alongside the
6-step and Reference Integrity families. Anti-pattern side (to author if a
second violation reproduces): `enforcement/failure-patterns/` —
"non-authoritative-source-blocks-pipeline" (not yet created; Finding A is the
first instance).

← [Back to governance/lifecycle](README.md)
