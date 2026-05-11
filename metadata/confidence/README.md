# Metadata Confidence

`metadata/confidence/` defines how to label evidence strength for Knowledge Atoms and routing surfaces.

## Confidence Values

| Value | Meaning | Allowed use |
| --- | --- | --- |
| `low` | Plausible but not yet validated in use. | Candidate maps, early summaries, untested atom proposals. |
| `medium` | Supported by source review, link validation, or one successful use. | Navigation rows, candidate atoms, pilot maps. |
| `high` | Repeatedly used or reviewed, with clear validation evidence. | Validated atoms, stable routing, promoted guidance. |

## Status Relationship

| Lifecycle status | Minimum confidence | Notes |
| --- | --- | --- |
| `temporary` | `low` | Short-lived or project-local; do not index as durable. |
| `candidate` | `low` | Can be routed if labeled as candidate. |
| `validated` | `medium` | Requires a real use, review, or explicit validation record. |
| `stable` | `high` | Requires repeated use or strong review evidence. |
| `deprecated` | any | Must include replacement or reason. |

## Evidence Signals

Confidence can increase when:

- Markdown links and source paths resolve.
- The atom was used in a completed task.
- A reviewer or user accepted the guidance.
- A test, fixture, lint, link check, or close-loop validation passed.
- The atom remains aligned with old source-of-truth entrypoints after a skill update.

Confidence should stay low when:

- The atom is only a planning guess.
- The old skill is still changing and no synchronization rule exists.
- Validation is missing or blocked.
- The guidance depends on project-specific evidence that has not been generalized.

## Downgrade Conditions

Downgrade or mark stale when:

- Old `skills/` source changes and the atom has not been rechecked.
- Links break.
- A conflict appears with `shared-rules/` or source-of-truth skill behavior.
- A candidate path is mistaken for a replacement path.

## Validation

Every atom with `confidence: high` should include or link to:

- The source path.
- The validation method.
- The last known lifecycle state.
- Any compatibility or deprecation note.
