# Dogfood Evidence — Implementation execution mode (force_exit path)

**Plan**: [`2026-06-29-1430-preparatory-refactoring-workflow.md`](2026-06-29-1430-preparatory-refactoring-workflow.md) Phase 4  
**Date**: 2026-06-29  
**Task**: Advisory `illegal_transition` scan for implementation-plan `steps[]` YAML (dogfood-before-validator)

## Intake routing

```yaml
change_kind: feature
blocked_by_structure: true   # planvalidate engine pure; intent YAML buried in markdown fences
execution_mode: preparatory_refactoring
```

**Not** replacement parity — no old/new capability inventory.

## Implementation plan (as executed)

```yaml
execution_mode: preparatory_refactoring
steps:
  - id: prep-01
    intent: structure
    behavior_change:
      allowed: false
    action: extract ExtractFencedYAMLBlocks seam for ```yaml fence parsing
    checkpoint:
      observable_equivalence:
        required: true
        evidence: planvalidate.TestExtractFencedYAMLBlocks

  - id: prep-02
    intent: structure
    behavior_change:
      allowed: false
    action: draft full ImplementationPlanParser (execution_mode, compatibility, nested checkpoint schema)
  # force_exit: no_structure_progress_detected — parser does not make illegal-transition test expressible; abstraction_not_used_by_next_feature

  - id: scope-shrink
    intent: feature
    behavior_change:
      allowed: true
    action: drop full parser; ship minimal ImplementationStep + DetectIllegalIntentTransitions only

  - id: feat-01
    intent: feature
    behavior_change:
      allowed: true
    action: AdvisoryValidateImplementationIntent + tests; Blocking=false (no hook wiring)
    validation: go test ./internal/planvalidate/ -run ImplementationIntent
```

## Path classification

| Segment | Classification |
|---------|----------------|
| prep-01 → checkpoint | structure + observable equivalence (tests green) |
| prep-02 | structure inflation / abstraction orphan → **force_exit** |
| scope-shrink → feat-01 | shrink scope + feature (advisory validator surface) |

**Valid evidence path**: Stop 設計驗證 path（structure → structure → force_exit → 縮 scope → feature）

## Failure-mode review lens

| failure | Dogfood observation |
|---------|---------------------|
| intent oscillation | Avoided — no feature→structure reopen |
| structure inflation | prep-02 would have added nested parser without exit |
| fake equivalence | Mechanical check accepts any non-empty `evidence`; human review must require runnable proof |
| abstraction orphan | prep-02 parser unused after shrink — caught by force_exit |
| compatibility collapse | N/A — new code path; old plans without `steps` unaffected |
| illegal transition | `DetectIllegalIntentTransitions` — primary feature outcome |

## Artifacts

| Artifact | Path |
|----------|------|
| YAML fence seam | `scripts/ai-skill-cli/internal/planvalidate/markdown_yaml.go` |
| Advisory intent transitions | `scripts/ai-skill-cli/internal/planvalidate/implementation_intent.go` |
| Tests | `scripts/ai-skill-cli/internal/planvalidate/implementation_intent_test.go` |
| Failure modes §8 | `workflow/software-delivery/implementation/execution-modes.md` |

## Validator posture

- **This phase**: advisory only (`Blocking: false`); **not** commit-msg / stop-hook wired.
- **Next** (after evidence): evaluate illegal-transition maturity ladder; schema may move to planvalidate engine profile.

## Recorded fields (Phase 4 checklist)

| Field | Value |
|-------|-------|
| change_kind | feature |
| execution_mode | preparatory_refactoring |
| intent sequence | structure → structure → (force_exit) → feature → feature |
| transitions | prep-01→scope-shrink legal (equivalence); prep-02 abandoned before illegal transition |
| force_exit trigger | `no_structure_progress_detected` + `abstraction_not_used_by_next_feature` |
| checkpoint | TestExtractFencedYAMLBlocks, TestAdvisoryValidateImplementationIntent_forceExitDogfoodFixture |
