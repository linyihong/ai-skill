# Implementation Plan: <feature>

## Architecture Compatibility Preflight
- [ ] Candidate files still exist or are marked `not applicable` / `source missing`.
- [ ] Source-of-truth and generated surfaces are identified.
- [ ] Layer responsibility is consistent with current architecture.
- [ ] Compiler / validator impact is known before implementation starts.
- [ ] Open questions reconciled: each is marked `resolved` (with preflight evidence) / `still-open` / `deferred`, and resolved ones are written back to the plan — not left answered only in working notes.

## Pre-build Interrogation
- **Goal**: <what user/system outcome this plan must achieve>
- **Scope**: <affected behavior, files, workflow, runtime surface, tool adapter, or generated artifact>
- **Non-goals**: <explicitly out of scope>
- **Acceptance / validation target**: <test, runtime validate, scenario, query, review, link check>
- **Framework discovery**: <canonical source, owner layer, projection, mirror/cache/generated output, linked updates>
- **Duplication risk**: <none | remove duplicate | deprecate old path | explicit precedence>
- **Open questions**: <blocker_question | safe_assumption | scoped_out | invalidated>
- **Decision**: <proceed | ask_user | revise_plan | blocked>

## Intake routing (from change brief)
- **change_kind**: <feature|bugfix|replacement|migration|internal_refactor>
- **blocked_by_structure**: <true|false|n/a>
- **execution_mode**: <direct_change|preparatory_refactoring> <!-- default direct_change when omitted -->

```yaml
# Routed from intake.md dual-axis; copy into plan frontmatter when useful
change_kind: feature
blocked_by_structure: true
execution_mode: preparatory_refactoring
compatibility:
  missing_execution_mode: direct_change
  missing_intent_on_steps: allowed   # opt-in; required when preparatory_refactoring
```

## Refactor / Replacement Parity Preflight
- [ ] Not applicable, or parity inventory exists before implementation.
- [ ] Legacy surfaces, inputs, outputs / side effects and external dependencies are mapped to target surfaces.
- [ ] Deferred / not planned / tool-specific items include a non-blocking reason.
- [ ] High-risk side effects have dry-run, fake-root, fixture or equivalent validation evidence.
- [ ] **Not** using parity inventory for `feature` + `preparatory_refactoring` (orthogonal — see [`intake.md`](../intake.md)).

## Execution steps (Change Intent Lock)

When `execution_mode: preparatory_refactoring` (or plan explicitly uses intent on steps), each step declares `intent` and `behavior_change`. Canonical rules: [`implementation/execution-modes.md`](../implementation/execution-modes.md).

```yaml
steps:
  - id: prep-01
    intent: structure
    behavior_change:
      allowed: false
    action: <extract seam / isolate logic>
    checkpoint:
      observable_equivalence:
        required: true
        evidence: <fixture | regression | mutation | parity note>

  - id: feat-01
    intent: feature
    behavior_change:
      allowed: true
    action: <land feature behavior>
    validation: <new acceptance / BDD / contract proof>
```

**Intent transition** (state machine):
- `structure → feature` requires `observable_equivalence_passed`
- `feature → structure` requires `explicit_reopen_reason`

## Stop condition checklist

Structure mode — exit when **any** holds (`exit_when`):
- [ ] Target change becomes local
- [ ] Target test becomes expressible
- [ ] New abstraction created (and referenced by next feature step)

Force exit when **any** holds (`force_exit_when`):
- [ ] `no_structure_progress_detected` (feature still not local; acceptance not expressible; seam unused)
- [ ] `abstraction_not_used_by_next_feature`

**Avoid** (even in structure intent): broad_cleanup, style_only, debt_harvesting, opportunistic_refactor.

On force exit: no more structure steps — enter feature mode, shrink scope / `direct_change`, or return to intake to re-route.

## Task Breakdown
### Task 1: <title>
- **File(s)**: <path>
- **Description**: <what to do>
- **Acceptance**: <how to verify>
- **intent** *(if preparatory)*: <structure|feature>

### Task 2: <title>
- **File(s)**: <path>
- **Description**: <what to do>
- **Acceptance**: <how to verify>
- **intent** *(if preparatory)*: <structure|feature>

## Dependencies
- <prerequisite task or external dependency>

## Risk Assessment
| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| <risk> | <H/M/L> | <H/M/L> | <plan> |

## Traceability
- **Upstream**: <link to change-brief | contract | bdd-scenario>
- **Downstream**: → Review Report
