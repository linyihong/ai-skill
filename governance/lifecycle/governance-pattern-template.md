# Governance Pattern Template — Mechanical Enforcement Subsystem

> **Status**: promoted template (positive pattern). Normative.
> Sibling of `enforcement/failure-patterns/` — failure patterns capture
> anti-patterns, this captures the positive recurring shape a new mechanical
> governance subsystem should take.
>
> Promoted 2026-06-25 from
> [`governance-pattern-library-draft.md`](governance-pattern-library-draft.md)
> under plan
> [`2026-06-08-2100-governance-pattern-library-extraction.md`](../../plans/active/2026-06-08-2100-governance-pattern-library-extraction.md)
> (Phase 2). The draft remains as the **evidence log** (7 samples, two
> falsifiable predicates) and is the source for *why* each step is or is not
> invariant. Read this file for the contract; read the draft for the proof.

## Thesis — a governance pattern is not a fixed pipeline

The naïve reading of these subsystems is a 6-step pipeline
(`Observation → Rule → Registry → Projection → Executor → Validation`). Seven
samples disprove that as an *invariant*: two of the six steps are conditionally
absent, each under its own falsifiable predicate.

**The real shape is an invariant core plus justified omissions:**

```
Observation ─┐
Registry     │  invariant core   (present in 7/7 samples)
Executor     │  — omitting any of these is a defect
Validation  ─┘

Rule        — near-universal (6/7); omittable ONLY for a pure structural invariant
Projection  — conditional   (4/7); omittable when the executor reads the source directly
```

A template that demanded all six steps unconditionally would force structural
invariants to invent an empty rule surface, and force direct-consumption
executors to invent a redundant projection. The value of this template is not
"do these six things" — it is **"do these four always, and these two only when
their predicate fires; an omission outside the predicate is a defect."**

## Applicability

Use this template when onboarding a **new mechanical governance capability** —
a rule whose violation can be detected by code at build-time, compile-time, or
commit-time (schema checks, reference resolvability, drift detection, coverage
gates). It is the **build-time** shape ("how a new enforced rule comes into
existence").

Do **not** use it for:

- Maintenance-time surface agreement → that is the *Reference Integrity* family
  (still observation-stage in the draft).
- Enforcement-time blast-radius / standing decisions → that is the *Failure
  Authority* family (still observation-stage in the draft).
- Editorial / advisory guidance with no mechanical executor — if there is no
  code that can detect the violation, this is not the right shape.

## The six steps

### Invariant core (omitting any of these is a defect)

| Step | Contract | Canonical example |
|---|---|---|
| **1. Observation** | A named failure / drift / miss that motivates the subsystem. The mature form is a promoted entry in `enforcement/failure-patterns/` — the positive subsystem is the constructive dual of an anti-pattern. | F19's observation points at [`coverage-evidence-dangling-reference.md`](../../enforcement/failure-patterns/coverage-evidence-dangling-reference.md). |
| **3. Registry** | The binding that makes the capability *declared and discoverable* — an `enforcement-registry.yaml` `rule_class`, a routing-registry entry, or a runtime-index table. Without it the rule exists only in prose and cannot be enforced or audited. | `enforcement-registry.yaml` `plan_tree_governance` rule_class with its `executors[]` block. |
| **5. Executor** | The code that detects the violation — validator(s), gate(s), or lint. May be a **sub-DAG** of 1..N placements (1 for a single compile lint, up to 5 for plan-tree's commit-msg validators). Cardinality is a per-subsystem decision; see [`validation-coverage-gap-executor-placement.md`](../../enforcement/failure-patterns/validation-coverage-gap-executor-placement.md). | `plan_tree.go` — 5 commit-msg validators. |
| **6. Validation** | Executable proof the executor detects the violation. Two **interchangeable, additive** sub-forms: a Go unit test **and/or** a scenario yaml. "Validation present" (core invariant) is distinct from "scenario-yaml coverage complete" (a separate governance metric). A subsystem may be Validation-complete in one sub-form while the other is independently tracked. | `runtime_trigger_wiring_test.go` **+** `orphan-routing-entry-v1.yaml` (both sub-forms). |

### Conditional steps (omittable only under their predicate)

| Step | Contract | Predicate (when REQUIRED) | When OMITTABLE |
|---|---|---|---|
| **2. Rule** | The human-authored policy surface defining what counts as a violation (a `*.yaml`, schema, or prose rule). | The invariant is **policy-derived** — a human authored the definition of the valid relation ("what counts as a leak", "what counts as a valid consumer"). This is the **near-universal** case. | **Only** when the invariant is a **pure structural (in)equality** — the violation *is* `A != B` and no human authorship defines the valid relation (checksum match, raw id collision). Attested by exactly **one** clean case (#4, runtime-index checksum). Reference-resolvability-style invariants usually *look* structural but smuggle in a policy rule — treat the exemption as rare. |
| **4. Projection** | A derived runtime surface materialized for the executor (`runtime.db` routes, `derived_match_tokens`). | The executor needs **indirect consumption**: a hot path that cannot re-parse the source each call, pre-digested data the source does not expose, or a unified query surface over heterogeneous sources. | When the executor consumes the **authoritative source directly** (`os.ReadFile` the staged plans, `yaml.Unmarshal` the registry). Stated at the *contract* level, not the implementation — the essence is "reads the authoritative source, needs no intermediate surface." Attested by **three** independent cases (#5/#6/#7); the best-supported predicate. |

## Decision tree

Walk this for each candidate subsystem. Each "no" on a core step is a defect to
fix; each conditional step is a deliberate, recorded choice.

```
1. Is there a named Observation (ideally a failure-pattern)?  ──no──▶ DEFECT: name the failure first
        │ yes
        ▼
2. Is the invariant policy-derived (a human defines the valid relation)?
        ├─ yes ─▶ Rule REQUIRED (authored policy surface)
        └─ no, it is a pure structural (in)equality (checksum / raw uniqueness)
                 ─▶ Rule omittable — RECORD the structural justification
        │
        ▼
3. Is the capability declared in a Registry (rule_class / routing / index)?  ──no──▶ DEFECT
        │ yes
        ▼
4. Does the executor need indirect consumption
   (hot path / pre-digested / unified query)?
        ├─ yes ─▶ Projection REQUIRED
        └─ no, it reads the authoritative source directly
                 ─▶ Projection omittable — RECORD the direct-consumption justification
        │
        ▼
5. Is there an Executor (≥1 validator/gate/lint)?  ──no──▶ DEFECT
        │ yes
        ▼
6. Is there Validation (Go test and/or scenario yaml)?  ──no──▶ DEFECT
        │ yes
        ▼
   Subsystem conforms to the template.
```

## Template Exit Criteria

A subsystem **conforms** when the invariant core is complete and every omitted
conditional step carries its recorded justification. An omission outside its
predicate is not "a lighter variant" — it is a defect.

**Acceptable omissions (with recorded justification):**

- **Rule absent** — only if the invariant is a *pure structural (in)equality*
  (checksum match, raw id collision) where no human authorship defines the valid
  relation. The justification must name the structural relation.
- **Projection absent** — only if the executor performs *direct authoritative
  consumption* (reads the source surface itself). The justification must state
  why indirect consumption is not needed.

**Forbidden omissions (always a defect):**

- **Observation absent** — building enforcement before naming the failure it
  prevents.
- **Registry absent** — a rule enforceable only by prose, undiscoverable and
  unauditable.
- **Executor absent** — a declared rule with no code to detect violation
  (`rule-without-executor`).
- **Validation absent** — an executor with no executable proof it detects the
  violation.

> The forbidden list is the load-bearing half of this template. The point of
> distilling seven samples was not to license omissions but to make the *legal*
> ones explicit so the *illegal* ones stand out.

## Relationship to failure patterns

This template is the positive (constructive) dual of two anti-patterns:

- Step 5 **Executor** ↔
  [`rule-without-executor.md`](../../enforcement/failure-patterns/rule-without-executor.md)
  (a Registry/Rule with no Executor) and
  [`validation-coverage-gap-executor-placement.md`](../../enforcement/failure-patterns/validation-coverage-gap-executor-placement.md)
  (executor placement variants / coverage gaps).
- Step 1 **Observation** ↔ the failure-pattern library as a whole: the mature
  Observation step *is* a promoted failure pattern.

## Provenance

- Evidence log (7 samples, per-step counts, the two falsifiable predicates):
  [`governance-pattern-library-draft.md`](governance-pattern-library-draft.md).
- Sibling families still in incubation (own N≥5 gates, not yet promoted):
  *Reference Integrity* and *Failure Authority* — both documented in the draft.
  If promoted they become siblings here, not children of this template.

← [Back to governance/lifecycle](README.md)
