# Project Metadata Schema — Migration Notes

> **Purpose**: a small living record of which framework surfaces consume each project metadata schema and where each consumer sits on its adoption trajectory. Reviewers three months from now should be able to read this file alone and know whether legacy code paths are deliberate (waiting for a planned phase) or accidental (forgotten cleanup).
>
> Update this file every time a consumer migrates from one schema version to another, or when a new consumer is added.

## `ai-skill-project-schema.yaml` (v1, candidate)

### Schema (canonical surface, 2026-06-08)

Declared at: [`ai-skill-project-schema.yaml`](ai-skill-project-schema.yaml).

```
project:
  id: <kebab-case>
  visibility: private | public
  private_entities:
    - name: <canonical entity name>      # governance layer ID
      kind: codename | client | product | individual | other
      match_tokens: [<aliases>]           # execution layer matching surface
      case_variants: auto | [<list>]
```

Governance layer (entity identity) and execution layer (match tokens) are **explicitly separated**. Reviewer discipline: any code touching this schema must respect that boundary — entity name flows into governance / audit / debug output; match_tokens flow into the scanner's matching pipeline. Mixing the two is the dual-source-of-truth trap this schema was designed to avoid.

### Current framework runtime state (as of Phase 1A landing; Phase 1C migration pending)

The pre-existing scanner implementation in `scripts/ai-skill-cli/internal/app/sanitization_scan.go` (landed in commits `2ff3a01`, `1e97bcf`, `97ea413`, before Phase 1A was designed) consumes a **legacy flat schema**:

```go
// sanitization_scan.go (legacy, lines 21-26)
type aiSkillProjectMetadata struct {
    ID              string
    Visibility      string
    PrivateTokens   []string   // flat list, no entity grouping
    PrivateEntities []string   // flat list, NOT structured entity
}
```

Projection emits to a single SQLite table: `derived_forbidden_tokens` (with columns `token`, `canonical_token`, `owning_project_id`, `source_metadata_path`, `suggested_placeholder`).

The legacy reader treats `private_tokens` and `private_entities` identically — both are concatenated into a single flat token list before projection. There is no `kind` field, no per-entity grouping, and no separation between governance and execution layers.

### Target state (Phase 1C + Phase 1D, in design)

After Phase 1C lands, the projection rule will read project metadata via the **canonical parser** at `scripts/ai-skill-cli/internal/app/project_metadata.go::LoadProjectMetadata` and emit to **two** SQLite tables:

- `derived_private_entities` — governance layer. One row per `(project_id, entity_name, kind, source_metadata_path)`. Used by governance lints, audit, and finding messages.
- `derived_match_tokens` — execution layer. One row per `(matched_token, owning_entity_name, owning_project_id)` after `case_variants` expansion. Used by the scanner's hot path.

After Phase 1D lands, the scanner in `sanitization_scan.go` will be renamed to a legacy reader (or replaced outright) and the canonical scanner will query `derived_match_tokens` instead of `derived_forbidden_tokens`.

### Migration ownership

| Surface | Current schema | Target schema | Migration phase | Owner |
| --- | --- | --- | --- | --- |
| Schema declaration (this file) | — | v1 canonical | **Phase 1A (this commit)** | framework maintainer |
| Parser `project_metadata.go` | — | structured loader | **Phase 1A (this commit)** | framework maintainer |
| Topology surface | — | `runtime/repository-topology.yaml` v2 | Phase 1B | framework maintainer |
| Projection rule | flat → `derived_forbidden_tokens` | structured → `derived_private_entities` + `derived_match_tokens` | Phase 1C | framework maintainer |
| Scanner | reads `derived_forbidden_tokens` | reads `derived_match_tokens` (with entity attribution from join) | Phase 1D | framework maintainer |
| Registry / failure pattern / scenarios | — | enforcement-registry rule_class executor binding | Phase 4 | framework maintainer |

### Why the legacy reader still exists after Phase 1A

The pre-existing `sanitization_scan.go::readProjectMetadata` (legacy flat schema reader) was committed before the Phase 1A canonical schema was designed. Phase 1A by deliberate scope discipline **does not modify `sanitization_scan.go`** — even one-line edits would blur the boundaries between phases 1A / 1C / 1D that the plan rewrite of 2026-06-08 specifically established.

Until Phase 1C lands, both readers co-exist in the binary:

- `LoadProjectMetadata` (canonical, in `project_metadata.go`) — parses the v1 structured schema; **not yet wired into any consumer**. Exists for parser correctness + test coverage.
- `readProjectMetadata` (legacy, in `sanitization_scan.go`) — parses the flat schema; **wired into the live projection rule** that emits `derived_forbidden_tokens`.

Phase 1C will introduce the new projection rule that calls `LoadProjectMetadata`, then either retire `readProjectMetadata` or rename it `readProjectMetadataLegacy` for a transition window. Phase 1D will switch the scanner to the new tables.

### What this is NOT

This file is **not** an architecture decision record. The architectural decision (governance/execution layer separation) was made in plan `2026-06-06-1800-sanitization-mechanical-enforcement.md` review on 2026-06-08 and is captured there. This file only tracks **adoption status** of that decision across the framework runtime.

Do not promote this file to ADR or canonical philosophy. When all consumers have migrated to the target state, delete the "Current framework runtime state" section and replace with a one-line "fully migrated as of `<commit>`" note. When the entire schema is retired (e.g. v2 supersedes v1), this file becomes a historical artifact and may move to `archived/` or be deleted.

### Cross-links

- Schema: [`ai-skill-project-schema.yaml`](ai-skill-project-schema.yaml)
- Example: [`example-ai-skill-project.yaml`](example-ai-skill-project.yaml)
- Parent plan: [`plans/archived/2026-06-06-1800-sanitization-mechanical-enforcement.md`](../../plans/archived/2026-06-06-1800-sanitization-mechanical-enforcement.md)
- Legacy reader source: `scripts/ai-skill-cli/internal/app/sanitization_scan.go::readProjectMetadata` (DO NOT modify until Phase 1C)

← [Back to project metadata index](README.md)
