# Experience Runtime (Cross-Cutting)

Immersive client surfaces (player, editor, onboarding) cut across **runtime state**, **journey specification**, **validation execution**, and **UI contracts**. This directory holds **cross-cutting templates** — not a `software-delivery` slice.

## L3 reminder

```text
Contract → Behavior → Validation Capability → Evidence
```

Experience runtime describes **client state transitions** and **which validation dimensions** each transition needs. It does not replace journey BDD or integration tests.

## Boundary table

| Concern | Owner | This cross-cutting doc | Do not duplicate |
| --- | --- | --- | --- |
| Journey specification | Project BDD / Gherkin | Reference states only | Full scenario text |
| Journey execution / evidence | `workflow/software-delivery/validation/` | Point to `evidence:*` tokens | Integration runner code |
| UI visual contract | Project `ui-style-reference/specs/` | Map states → spec sections | Pixel values |
| Evidence types | `validation/evidence-types/` | Map transitions → types | Collection methods in gate `requires:` |
| Failure evolution | `validation/failure-evolution-catalog.md` | Link incidents → states | Ad-hoc incident logs |
| Experience-validation pipeline | `plans/active/2026-06-09-1040-*` | Complement; taxonomy split with evidence-types | Collapse collection_method into evidence_type |

## Pilot templates

| Surface | File | Downstream consumer |
| --- | --- | --- |
| Immersive player | [`player.yaml`](player.yaml) | Project `.ai-skill/project/rules/player-client-patterns.md`, screen mapping |

## State → evidence mapping (pattern)

| Transition class | Typical evidence_type |
| --- | --- |
| User-visible overlay / CTA | `user_visible` |
| Route / href / history | `navigation` |
| Time boundary / poll / debounce | `temporal_behavior` |
| sessionStorage / cookie readback | `state_persistence` |
| Video load / pause / seek | `media_playback` |
| Static contract string | `source_contract` |

Gate `requires:` lists **only** types relevant to the changed transition — not all six.

## Slice promotion conditions

Promote to `software-delivery` slice **only when**:

1. **Player** pilot has stable state machine + claim registry + integration envelope (✅ downstream pilot)
2. **Editor** (or equivalent authoring surface) documents the same pattern
3. **Onboarding** (or equivalent first-run surface) documents the same pattern

Until then: **no** `route.workflow.experience-runtime`, **no** `sd-experience-runtime` artifact gates.

## Related

- [`../software-delivery/validation/evidence-gate-vocabulary.md`](../software-delivery/validation/evidence-gate-vocabulary.md)
- [`../software-delivery/validation/authority-decision-table.md`](../software-delivery/validation/authority-decision-table.md)
- Plan: [`2026-06-12-1600-projection-break-evidence-taxonomy-experience-runtime.md`](../../plans/active/2026-06-12-1600-projection-break-evidence-taxonomy-experience-runtime.md) §C.3
