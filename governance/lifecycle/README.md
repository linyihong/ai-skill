# Knowledge Lifecycle

`governance/lifecycle/` defines how knowledge moves from existing source files into the new AI-native layers without breaking old skill entrypoints.

## Source Of Truth Rule

Until a migration is explicitly promoted and validated, existing `skills/`, `shared-rules/`, `ai-tools/`, and `scripts/` files remain the source of truth for executable behavior.

New layer files under `analysis/`, `workflow/`, `intelligence/`, `runtime/`, `memory/`, `feedback/`, `models/`, `governance/`, `knowledge/`, and `metadata/` may act as:

- Routing surfaces.
- Candidate maps.
- Promotion targets.
- Metadata or summary surfaces.
- Governance and runtime design.

They must not silently replace old skill behavior before promotion.

## Durable Goal Boundary

Long-term lifecycle states belong in durable planning files, not in `.agent-goals/`.

| Goal type | Durable location |
| --- | --- |
| Repository roadmap, phase, migration sequence | `architecture/` |
| Layer responsibility, candidate destinations, promotion targets | Layer README files |
| Knowledge lifecycle, validation and deprecation rules | `governance/` |
| Routing, metadata and atom discovery | `knowledge/`, `metadata/`, `runtime/` |
| Active implementation work for the current conversation | `.agent-goals/` only until completion |

Before deleting an active `.agent-goals/` entry, confirm that any remaining roadmap, lifecycle, migration, promotion, deprecation or follow-up state has been written to the durable location above.

## Lifecycle States

| State | Meaning | Allowed content | Not allowed |
| --- | --- | --- | --- |
| `source-of-truth` | Current canonical behavior lives here. | Existing skill/shared rule/tool/script files. | Treating a newer map as override. |
| `candidate-map` | A map from current sources to future layer destinations. | Ownership boundaries, source-to-target tables, compatibility notes. | Bulk content migration or behavior changes. |
| `candidate-atom` | A proposed Knowledge Atom or summary. | Metadata, summary, links to source, validation criteria. | Marking as stable without use. |
| `validated-atom` | A candidate used or reviewed successfully. | Routing metadata, summary, checklist, validation evidence. | Removing old entrypoint. |
| `promoted` | New layer becomes a supported reference path. | Old entrypoint links to promoted atom, index routes to both when needed. | Deleting compatibility path without deprecation. |
| `deprecated` | Old path is being retired with a replacement. | Deprecation note, replacement link, validation record. | Breaking existing links or tool loading. |

## Promotion Gates

A candidate can be promoted only when all gates pass:

1. The old source path remains reachable or has a redirect note.
2. `metadata/schema.md` metadata exists for the promoted atom or surface.
3. `knowledge/indexes/README.md` routes relevant task intents.
4. The owning layer README links the new path.
5. Validation is documented in `governance/validation/`.
6. Diff review confirms no project-specific evidence, secrets, local absolute paths, or tool mirror paths were introduced.
7. Any durable roadmap or lifecycle state has been updated outside `.agent-goals/`.
8. Commit, push, readback, and clean status have completed.

## Update Strategy While Skills Still Change

When an old skill is updated before migration:

1. Update the old `skills/<name>/` source first.
2. Check whether a candidate map or promoted atom references the changed section.
3. If yes, update the map, metadata, summary, or index in the same change.
4. If no, record that no linked update was needed in the final validation.
5. Do not copy new skill text into a new layer unless the change is an explicit atom promotion.

## Deletion Rule

Do not delete or move old skill files during candidate-map or candidate-atom phases. Deletion can be considered only after:

- The replacement path is promoted.
- Existing tool adapters can still load the skill or have documented replacements.
- Links, indexes, summaries, and metadata are updated.
- A deprecation note and rollback path exist.
