# Failure Evolution Catalog

Structured entries for **Failure → Classification → Authority → Evolution → Writeback**. Not an incident log — each entry must include **counterfactual**.

Authority routing: [`authority-decision-table.md`](authority-decision-table.md).  
Evidence gates: [`evidence-gate-vocabulary.md`](evidence-gate-vocabulary.md).

## Required fields

| Field | Purpose |
| --- | --- |
| `failure` | Stable id |
| `classification` | Projection layer + missing evidence / symptom |
| `authority` | Who may change what |
| `evolution_target` | rule / gate / scenario / code / playbook |
| `writeback` | Paths touched |
| `counterfactual` | What would have blocked merge earlier |

Optional: `not_framework_invariant_because`, `scenario`, `experience_runtime_state`.

## Entries

### `player_preview_gate_projection_break`

```yaml
failure: player_preview_gate_projection_break
classification:
  layer: L2_behavior_wrong_dom_target
  projection_break: Contract_and_BDD_pass; Behavior_on_wrong_video; L3_caught_via_integration
  missing_evidence_before_fix: [user_visible, temporal_behavior]
  symptom: BDD source assert passed; preview overlay absent in browser
experience_runtime_state:
  expected: preview → gated
  actual: preview poll on adjacent preload video
authority:
  - domain_pattern
  - implementation_defect
evolution_target:
  - integration test with evidence envelope
  - player-client-patterns.md (playerStage DOM ownership)
  - gate.requires evidence tokens
  - claim_registry in project workflow
writeback:
  framework:
    - validation/evidence-types/*
    - workflow/software-delivery/validation/evidence-gate-vocabulary.md
    - workflow/cross-cutting/experience-runtime/player.yaml
  project:
    - tests/integration/player-preview-gate.integration.mjs
    - .ai-skill/project/rules/player-client-patterns.md
    - docs/frontend-contracts/screen-mapping/player-preview-gate.md
counterfactual: >
  If gate.short_drama.validation_complete required evidence:user_visible and
  integration envelope contained claim preview_overlay_shown with poll log
  scoped to playerStage video, wrong adjacent-video target would have blocked
  merge before deploy.
not_framework_invariant_because: >
  Taxonomy already sufficient; gap was missing L3 test and DOM ownership rule.
scenario: validation/scenarios/failure-derived/projection-break-missing-browser-evidence-v1.yaml
enforcement_promotion: deferred
enforcement_note: >
  Advisory pattern only; not mechanical rule_class until second cross-project instance.
```

## Adding entries

1. Classify projection break layer — do not skip to framework invariant.
2. Pick authority from [`authority-decision-table.md`](authority-decision-table.md).
3. Add failure-derived scenario when the pattern is statelessly reproducible.
4. Link experience-runtime state if immersive client involved.

## Related

- [`authority-decision-table.md`](authority-decision-table.md)
- [`../cross-cutting/experience-runtime/README.md`](../cross-cutting/experience-runtime/README.md)
- Plan Phase 4: [`2026-06-12-1600-projection-break-evidence-taxonomy-experience-runtime.md`](../../plans/active/2026-06-12-1600-projection-break-evidence-taxonomy-experience-runtime.md)
