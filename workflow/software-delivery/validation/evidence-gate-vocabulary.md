# Evidence Gate Vocabulary

Canonical vocabulary for **gate `requires:`** tokens. Only `evidence_type` may appear in `requires:`.

## Token format

```yaml
requires:
  - evidence:user_visible
  - evidence:navigation
  - evidence:temporal_behavior
```

| Token prefix | Allowed in `requires:` | Example |
| --- | --- | --- |
| `evidence:<evidence_type>` | ✅ | `evidence:user_visible` |
| `collection_method` | ❌ | `browser_observation` |
| `artifact_shape` | ❌ | `screenshot` |
| `browser_review` | ❌ (activity summary only) | `browser_review: passed` |

Catalog: [`validation/evidence-types/`](../../../validation/evidence-types/README.md)

## Three-way split (do not collapse)

| Layer | Role | Gate token? |
| --- | --- | --- |
| `evidence_type` | What was proven | ✅ `evidence:*` |
| `collection_method` | How acquired | ❌ envelope field |
| `artifact_shape` | Output form | ❌ envelope field |

OQ-5: **reject inheritance**. Do not model `user_visible → browser_observation → screenshot` as a subtype tree. Each evidence type lists `supported_collection_methods` / `supported_artifact_shapes` in its catalog file.

## Trace chain (required for completion claims)

```text
gate → claim → artifact
```

Forbidden:

- artifact without `claim.id`
- screenshot-only close-out without linked claim
- `requires:` listing `browser_observation` instead of `evidence:*`

## Gate vs envelope

**Gate** (`requires:`) — minimum evidence types that must be satisfied for a surface or task class.

**Envelope** (close-out / integration output) — per-claim proof:

```yaml
gate: gate.<project>.validation_complete
claims:
  - id: preview_overlay_shown
    evidence_type: user_visible
    collection_method: browser_observation
    artifact_shape: screenshot
    status: passed|blocked
    artifact: .tmp/evidence/player-preview-gate-overlay.png
browser_review: passed   # activity summary; not a gate token
```

## Claim registry (query shape — Phase 2 design target)

Stable claim ids bind gates to integration tests and artifacts:

```yaml
claims:
  - id: preview_overlay_shown
    evidence_type:
      - user_visible
    validates: "Preview limit mask and plan guide visible to guest"

  - id: subscribe_href_resolves
    evidence_type:
      - navigation
    validates: "Subscribe CTA href resolves without double basePath"

  - id: preview_limit_enforced
    evidence_type:
      - temporal_behavior
    validates: "Preview boundary reached before overlay; no auto-next during preview"
```

Rules:

- `id` is stable across BDD refs, integration test names, and envelope output.
- One claim may require one or more evidence types only when each type proves a distinct question.
- Do not invent claim ids per run; add new ids only when the observable assertion is new.

## Registered evidence types (Phase 1 catalog)

| Token | Question |
| --- | --- |
| `evidence:source_contract` | Does static contract / string / schema align? |
| `evidence:user_visible` | Did the user actually see it? |
| `evidence:navigation` | Did navigation land correctly? |
| `evidence:state_persistence` | Did state survive navigation / session? |
| `evidence:media_playback` | Did media element / stream behave correctly? |
| `evidence:temporal_behavior` | Did time-bounded behavior hold? |

## Surface-specific requires (downstream pattern)

Project gates list **only** evidence types relevant to the changed surface. Example (immersive player preview):

```yaml
gate.short_drama.validation_complete:
  requires:
    - evidence:user_visible
    - evidence:navigation
    - evidence:temporal_behavior
  when_changed_surface_includes:
    - preview_limit_overlay
    - subscribe_cta_navigation
    - time_bounded_client_transition
```

Do not copy the full six-type list into every gate.

## browser_review demotion

| Before | After |
| --- | --- |
| Pass/fail gate token | Human-readable activity summary |
| Proxy for UX complete | Must be backed by `evidence:*` + claim-linked artifacts |

Allowed close-out:

```yaml
browser_review: passed|blocked|deferred
evidence_envelope: [...]
```

## Delivery vs publication

| Field | Meaning |
| --- | --- |
| `delivery_status` | Local validation passed (integration / envelope) |
| `publication_status` | Changes reachable via remote / deploy |

**Integration pass ≠ delivered.** A gate may pass locally while the project repo is `not_published`.

## Related

- [`authority-decision-table.md`](authority-decision-table.md)
- [`validation/evidence-types/README.md`](../../../validation/evidence-types/README.md)
- [`README.md`](README.md)
- Scenario: [`validation/scenarios/software-delivery/evidence-type-projection-break-v1.yaml`](../../../validation/scenarios/software-delivery/evidence-type-projection-break-v1.yaml)
