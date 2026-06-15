# Authority Decision Table

Who may change what after a validation / projection failure. Prevents every incident from escalating to framework invariant changes.

## Loop

```text
Failure
  → Classification（projection break layer + missing evidence_type）
  → Authority Decision
  → Evolution Target
  → Writeback
```

## Authority types

| Authority | May change | Typical writeback |
| --- | --- | --- |
| `framework_invariant` | Ai-skill validation catalog, gate vocabulary, scenarios, workflow templates | Ai-skill repo |
| `domain_pattern` | Project overlay rules, screen mapping, player-spec, workflow gate consumption | `.ai-skill/project/rules/`, project docs |
| `implementation_defect` | Application code, integration tests, fixtures | `server_doc/framework-code/`, `tests/integration/` |
| `env_deploy_incident` | Deploy playbook, smoke checklist, staging URL config | project docs, gitignored env |

## Classification → authority (default routing)

| Classification | Primary authority | Secondary |
| --- | --- | --- |
| Missing L3 validation capability | `domain_pattern` + `framework_invariant` (if taxonomy gap) | — |
| L2 behavior on wrong runtime target | `implementation_defect` + `domain_pattern` | — |
| BDD pass but UX broken | `domain_pattern` (add integration + claim) | not framework_invariant alone |
| Gate vocabulary drift | `framework_invariant` | project workflow sync |
| Deploy / env mismatch | `env_deploy_incident` | — |
| Taxonomy needs new evidence_type | `framework_invariant` | watch-list; avoid type inflation |

## Pilot entry: player preview gate projection break

```yaml
failure: player_preview_gate_projection_break
classification:
  layer: L2_behavior_wrong_dom_target
  missing_evidence: [user_visible, temporal_behavior]
  symptom: BDD source assert passed; overlay never shown in browser
authority:
  - domain_pattern
  - implementation_defect
evolution_target:
  - integration test with evidence envelope
  - playerStage video selector ownership rule (deferred P2-4)
  - gate.requires evidence tokens
writeback:
  - tests/integration/player-preview-gate.integration.mjs
  - docs/frontend-contracts/ui-style-reference/specs/player-spec.md
  - validation/evidence-types/*
counterfactual: >
  If gate required evidence:user_visible and integration envelope contained
  claim preview_overlay_shown with poll log, wrong adjacent-video target
  would have blocked merge before production.
not_framework_invariant_because: >
  Incident was DOM ownership + missing L3 test, not a cross-project taxonomy error.
```

## Escalation rules

1. **Do not** promote to `framework_invariant` without cross-project or taxonomy-gap justification.
2. **Do** add project integration + claim when BDD proves contract strings but not runtime ownership.
3. **Do** update gate vocabulary when the same `requires:` ambiguity appears in two projects (defer until then).
4. Failure-derived scenario writeback → [`validation/scenarios/failure-derived/`](../../../validation/scenarios/failure-derived/) (Phase 4).

## Related

- [`evidence-gate-vocabulary.md`](evidence-gate-vocabulary.md)
- Plan: [`2026-06-12-1600-projection-break-evidence-taxonomy-experience-runtime.md`](../../../plans/active/2026-06-12-1600-projection-break-evidence-taxonomy-experience-runtime.md) §C.4
