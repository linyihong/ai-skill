# Evidence Types

`evidence_type` answers **what was proven**. It is not `collection_method` (how acquired) or `artifact_shape` (output form).

## L3 model

```text
Contract → Behavior → Validation Capability → Evidence (artifact + proof shape)
```

Validation Capability **produces** Evidence. Evidence is not a synonym for Validation.

## Three-way taxonomy

| Layer | Gate `requires:` | Where defined |
| --- | --- | --- |
| `evidence_type` | ✅ `evidence:user_visible` | This directory |
| `collection_method` | ❌ | [`2026-06-08 evidence acquisition layer`](../archived/2026-06-08-1544-evidence-acquisition-layer.md) / `sd-validation` |
| `artifact_shape` | ❌ | Integration envelope / close-out report |

## Envelope shape (project execution)

```yaml
claim: preview_overlay_shown
evidence_type: user_visible
collection_method: browser_observation
artifact_shape: screenshot
```

Trace chain for completion claims:

```text
gate → claim → artifact
```

## OQ-5: token inheritance — **reject**

`evidence_type` is **not** a subtype system. Do **not** model:

```text
user_visible → browser_observation → screenshot
```

as inherited gate tokens. That becomes an ontology and inflates `requires:` lists.

Instead, each evidence type file lists:

- `supported_collection_methods`
- `supported_artifact_shapes`

Gate `requires:` lists **evidence_type only**.

## Catalog

| File | `evidence_type` | Proves |
| --- | --- | --- |
| [source-contract.md](source-contract.md) | `source_contract` | Static contract / string / schema alignment |
| [user-visible.md](user-visible.md) | `user_visible` | User-visible UI state |
| [navigation.md](navigation.md) | `navigation` | Route / href / history / basePath correctness |
| [state-persistence.md](state-persistence.md) | `state_persistence` | State survives navigation / session |
| [media-playback.md](media-playback.md) | `media_playback` | Media element / stream behavior |
| [temporal-behavior.md](temporal-behavior.md) | `temporal_behavior` | Time-bounded observable transition |

Naming rule: evidence types must answer **what was proven**, not **why it breaks** (failure classes belong in classification, not type names).

## Related plans

- [`2026-06-12-1600-projection-break-evidence-taxonomy-experience-runtime.md`](../../plans/archived/2026-06-12-1600-projection-break-evidence-taxonomy-experience-runtime.md)
