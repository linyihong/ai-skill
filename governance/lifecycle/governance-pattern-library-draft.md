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

### Per-step counts (N=5)

| Step | Filled cells | Notes |
|---|---|---|
| Observation | 5/5 | Universal so far |
| Rule | 4/5 | #4 has "(implicit)" — structural-invariant hole |
| Registry | 5/5 | Universal so far |
| Projection | 4/5 | **#5 (F19) empty** — compile-time-direct-read hole (new this pass) |
| Executor | 5/5 | Universal so far |
| Validation | 5/5 | Universal so far |

**Two distinct optional-step signals now, on two different steps:**

1. **Rule** (sample #4): no explicit *Rule* layer. The rule is encoded directly in the executor (`nativeRuntimeIndexChecksumsCheck`) without a yaml/md surface declaring "checksum drift is forbidden." Signal: *Rule* is optional when the invariant is **structural** (matches a checksum) rather than **editorial** (matches a pattern).
2. **Projection** (sample #5 / F19): no projected intermediate surface. The executor `LintValidationScenarios` calls `loadRegistrySnapshot` + `os.ReadFile` to read `enforcement-registry.yaml` and the referenced scenario yaml **directly** at compile time; there is no `runtime.db` table or other pre-digested surface between the canonical source and the executor.

**Why F19 has no Projection — a load-bearing hypothesis (Projection ⇔ hot path):**

Projection appears to exist precisely when the executor runs in a **hot path** and cannot afford to re-read+parse canonical yaml on every invocation:

| # | Executor runs at | Projection present? | Surface |
|---|---|---|---|
| 1 Workflow | PreToolUse (every non-Read tool call) | ✅ | `runtime.db` routes |
| 2 Discovery | tool-call / advisory time | ✅ | `runtime.discovery.config` |
| 3 Sanitization | per-commit | ✅ | `derived_match_tokens` |
| 4 Runtime Index | compile + commit | ✅ (but Registry *is* the db: `sources` table) | `sources` table |
| 5 F19 | `runtime compile` only | ❌ | reads canonical yaml directly |

The pattern: **hot-path executors need a projection; a compile-time executor reading the canonical surface once per compile does not.** Sample #4 is the boundary case — it runs at compile time too, but its registry *is already* a SQLite db, so "reading the registry" and "reading a projection" collapse into the same `sources`-table read. F19 is the first sample where the registry is plain yaml and the executor is purely compile-time, so the Projection step has nothing to do and is genuinely absent — not skipped for convenience.

This is the inverse relationship to the Rule signal: Rule drops out for **structural** invariants; Projection drops out for **compile-time-direct-read** executors. The two optional steps are governed by different conditions, which is stronger evidence than two holes on the *same* step would have been.

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
| Total samples | ≥ 5 | 5 (F19 added 2026-06-12) | ✅ |
| At least one non-fitting sample analysed | ≥ 1 | 2 (#4 Rule-missing; #5/F19 Projection-missing) | ✅ |
| ≥ 3 samples per step | ≥ 3 each | 4-5 each | ✅ |

**All three gate criteria are now met.** But meeting the gate does **not** mean "promote the 6-step chain as-is." Per plan Phase 1, three branches exist, and the evidence points to the **middle** one:

> - If 5+ samples AND ≥3 per step AND 6-step shape **consistent** → Phase 2 as-is
> - If samples found but **step coverage uneven** → **revise template shape (some steps optional, document variants)** ← we are here
> - If samples diverge significantly → abandon promotion

Step coverage is uneven in a *principled* way: Observation / Registry / Executor / Validation are 5/5 (invariant core), while **Rule (4/5) and Projection (4/5) are each conditionally optional under different conditions** (structural-invariant; compile-time-direct-read). So the promotable artifact is not a rigid 6-step checklist but a **4-step invariant core + 2 conditional steps**:

```
Observation → [Rule?] → Registry → [Projection?] → Executor → Validation
              structural-           hot-path-
              invariant skips       executor adds
```

This is a better outcome than a clean 5th confirmation would have been: a clean fit would have left confirmation-bias unrefuted, whereas two holes on two different steps, each with a falsifiable condition, is what lets the eventual template state *when* a step may be omitted instead of demanding all six unconditionally.

Recommended Phase 1 disposition (for the plan author to confirm): proceed toward Phase 2 extraction, but author the template as **core + conditional steps**, and verify each conditional rule against the 6th sample before freezing the variant wording. Until the template lands, do not link this draft from `governance/lifecycle/README.md` and do not reference it as a normative pattern.

## Why this draft matters even if the gate never passes

Even if the 6-step hypothesis is disproved, the analysis produces durable knowledge:

- Concrete evidence of which steps are *truly* universal across governance subsystems
- A list of variants (e.g. "structural invariants skip the Rule layer")
- A record of which subsystems were *forced* into the shape vs. which fit naturally — that itself is a confirmation-bias signal worth capturing in `enforcement/failure-patterns/` if reproduced

---

## Sequential vs. branching shape

> Resolves Phase 0 checklist item "Check whether the pattern is sequential (always Observation → ... → Validation) or has parallel/optional branches".

**Current evidence (N=5) is sequential, but with four nuances**:

1. **Two optional steps observed, on two different steps under two different conditions** — sample #4 (Runtime Index Freshness) has no explicit *Rule* layer (rule encoded structurally in `nativeRuntimeIndexChecksumsCheck`, a sha256 equality check); sample #5 (F19) has no *Projection* layer (executor reads canonical yaml directly at compile time). The governing conditions are independent: *Rule* drops out for **structural** invariants (mechanical equality, not editorial pattern); *Projection* drops out for **compile-time-direct-read** executors (no hot path → no need to pre-digest the canonical surface). See §"Per-step counts" for the Projection ⇔ hot-path table.
2. **Within-step parallelism in Executor** — samples #2 and #3 have executor pairs: Discovery Bridge has `discovery.go` core + advisory injector hook integration as two co-equal entry points; Sanitization Phase 1 (planned) has scanner core + commit-msg validator. F19 also shows the within-step structure but at a *single* placement: its executor is `runtime compile`-only, with commit-transaction dual-placement explicitly deferred (child plan Q6 + "Close and Observe" decision). The single "Executor" cell hides a sub-DAG whose *cardinality* itself varies (1 placement for F19, 2 for #2/#3); the template should say "Executor = core + integration points (placement count is a per-subsystem decision, see `validation-coverage-gap-executor-placement.md`)".
3. **No observed re-entry / loop** — every sample so far is one-pass from Observation to Validation. No sample has a "Validation discovers gap → re-enter Rule" cycle yet, though Phase D of Discovery Bridge (3-week empirical) is essentially that. The shape may turn out to have a feedback loop annotation for subsystems with empirical iteration gates.
4. **Observation can be a named failure pattern, not just an ad-hoc note** — F19's Observation cell points at a *promoted* failure pattern (`coverage-evidence-dangling-reference.md`), where samples #1–#4 had inline observations. This suggests the mature form of the Observation step is "a failure pattern in `enforcement/failure-patterns/`," and the 6-step (positive) family is the constructive dual of the failure-pattern (anti-pattern) family — consistent with the plan's "this captures the positive shape the anti-pattern violates against."

**Working interpretation (updated N=5)**: the shape is **a 4-step invariant core (Observation → Registry → Executor → Validation) with two conditionally-optional steps (Rule, Projection) inserted under distinct conditions, plus within-step branching at Executor**. The two optional steps being governed by *different* falsifiable conditions — rather than two holes on the same step — is what upgrades this from "soft signal" to "document as variants" (plan Phase 1 middle branch). A 6th sample should be used to test the two conditions, not just to add another tally mark.

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
