# Repository Topology — Migration Notes

> **Purpose**: a small living record of where `runtime/repository-topology.yaml` sits between the v1 schema (currently in production) and the v2 schema (designed but not yet adopted in projection). Reviewers three months from now should be able to read this file alone and know whether the live YAML's v1 shape is deliberate (waiting for Phase 1C) or accidental (forgotten cleanup).
>
> Update this file every time a consumer or producer migrates between schema versions, or when a new consumer is added to the runtime.

## `runtime/repository-topology.yaml` (currently v1, target v2)

### v1 schema (canonical surface as of Phase 1B landing; Phase 1C migration pending)

```yaml
schema_version: 1
status: active
owner_layer: runtime

runtime_projection:
  enabled: true
  target_key: runtime.repository_topology.config

shared_layer_classification:
  - subtree: plans/
    shared: true
  - subtree: scripts/
    shared: false
  # ...

expected_consumers:        # ← this field is what v2 removes
  - sanitization
  - workflow_activation
  - governance_lint
  - dependency_reading
```

Projection rule lives in `scripts/ai-skill-cli/internal/app/runtime_compiler.go` line 339:

```go
{"runtime/repository-topology.yaml", "repository_topology", "subtree",
 "shared_layer_classification", []string{"subtree"}}
```

This rule is **hard-coded against the v1 field names** (`shared_layer_classification`, `subtree`). The live YAML cannot be upgraded to v2 in isolation; doing so would break projection. The upgrade is paired with the projection rewrite in Phase 1C.

### v2 schema (target; designed but NOT yet on disk)

```yaml
schema_version: 2
status: active
owner_layer: runtime

runtime_projection:
  enabled: true
  target_key: runtime.repository_topology.config

# Frozen governance decision: consumer list is DERIVED from code-reference,
# not manually maintained. This block is intentionally permanent.
# Do NOT re-add a manual `expected_consumers:` field — manual consumer
# lists go stale (new subsystem reads topology but no one updates the list),
# which is the exact failure pattern this v2 was designed to remove.
consumer_tracking:
  strategy: code_reference
  rationale: |
    Manual consumer lists go stale (new subsystem reads topology but no
    one updates the list). Code-reference derivation is the durable
    governance decision. This block is intentionally permanent — do not
    re-add a manual `expected_consumers:` field.

subtrees:
  - path: plans/
    shared_layer: true
    owner: framework-maintainer
    purpose: "Plan tracking; referenced by enforcement-registry child_plan"
  - path: workflow/
    shared_layer: true
    owner: framework-maintainer
    purpose: "Cross-skill workflow contracts and execution flows"
  # ...
  - path: scripts/
    shared_layer: false
    owner: tooling-maintainer
    purpose: "CLI / runtime implementation; not consumed as reusable knowledge"

invariants:
  - Topology is a source-of-truth surface, not a detector heuristic.
  - Shared-layer classification must not infer project privacy by absence.
  - Project-local metadata remains the source of private token declarations.
  # Note: `expected_consumers:` is intentionally absent. See consumer_tracking
  # block above. Adding it back would re-introduce the stale-reference
  # failure pattern that v2 was designed to eliminate.
```

### Schema version precedence (mixed-shape files)

**Explicit `schema_version:` wins over shape inference.** If a YAML declares `schema_version: 1` but also contains v2-only fields (`subtrees:`, `consumer_tracking:`), the loader treats the file as v1 and **silently ignores the v2 fields**. The reverse is also true: `schema_version: 2` with `shared_layer_classification:` ignores the v1 fields.

This is intentional. The author signalled their intent via the explicit version, and partial migration files (where both shapes co-exist) would otherwise produce undefined behaviour. The Phase 1C migration of `runtime/repository-topology.yaml` must therefore be **atomic**: bump `schema_version: 1` → `2` and add the v2 fields in the **same edit**. Splitting the migration into "add subtrees first, bump version later" would silently lose the new fields.

If `schema_version:` is absent, the loader falls back to shape inference:

1. Presence of `subtrees:` OR `consumer_tracking:` → v2
2. Presence of `shared_layer_classification:` OR `expected_consumers:` → v1
3. Empty file → v2 (writer-friendly default)

Locked in by `TestLoadRepositoryTopology_ExplicitV1IgnoresV2Fields`.

### Schema deltas

| Aspect | v1 | v2 |
| --- | --- | --- |
| Subtree list field | `shared_layer_classification` | `subtrees` |
| Per-subtree path field | `subtree` | `path` |
| Per-subtree shared field | `shared` | `shared_layer` |
| Per-subtree owner | absent | `owner` (required) |
| Per-subtree purpose | absent | `purpose` (required) |
| Consumer list | `expected_consumers: [name, name, ...]` (manual) | `consumer_tracking: { strategy: code_reference, rationale: ... }` (frozen) |

### Current framework runtime state

| Surface | Schema in use | Reads via |
| --- | --- | --- |
| `runtime/repository-topology.yaml` (on disk) | v1 | — |
| `runtime_compiler.go` line 339 projection rule | v1 (field names hard-coded) | direct YAML loader inside compiler |
| `runtime.db.repository_topology` (projected table) | v1 row shape (subtree, content JSON) | populated by projection rule |
| `sanitization_scan.go::repositoryTopologyRow` (legacy consumer) | v1 row shape | reads SQLite table, not YAML |
| `scripts/ai-skill-cli/internal/app/repository_topology.go::LoadRepositoryTopology` (Phase 1B canonical) | reads BOTH v1 + v2; writes ONLY v2 | direct YAML I/O; **not yet wired into projection** |

### Target state (Phase 1C + Phase 1D)

After Phase 1C lands:

1. `runtime/repository-topology.yaml` upgrades to v2 on disk (subtrees with owner+purpose, expected_consumers replaced with consumer_tracking)
2. `runtime_compiler.go` projection rule rewrites to call `LoadRepositoryTopology` rather than hard-coding the v1 field names
3. `runtime.db.repository_topology` row content includes the v2 fields (owner, purpose) so downstream consumers can query them

After Phase 1D lands:

4. `sanitization_scan.go::repositoryTopologyRow` is either renamed legacy or extended to read the v2 fields it cares about (owner / purpose flow into finding messages for governance attribution)

### Migration ownership

| Surface | Current schema | Target schema | Migration phase | Owner |
| --- | --- | --- | --- | --- |
| Schema spec (this file) | — | v2 documented | **Phase 1B (this commit)** | framework maintainer |
| Loader `repository_topology.go` | — | v1 read + v2 read/write | **Phase 1B (this commit)** | framework maintainer |
| Live YAML `runtime/repository-topology.yaml` | v1 | v2 | Phase 1C | framework maintainer |
| Projection rule (compiler) | v1 hard-coded | calls `LoadRepositoryTopology` | Phase 1C | framework maintainer |
| Live consumer `sanitization_scan.go` | reads `subtree` + `shared` from SQLite | reads v2 fields via join | Phase 1D | framework maintainer |
| Registry / scenarios | — | enforcement-registry binding | Phase 4 | framework maintainer |

### Why the v1 file still exists after Phase 1B

The pre-existing `runtime/repository-topology.yaml` (v1) was committed in the same wave as the original sanitization scanner (commits `2ff3a01`, `1e97bcf`, `97ea413`). Phase 1B by deliberate scope discipline **does not upgrade the live YAML file** — even a one-field edit to add `owner` would either:

- Force the projection rule to also know about the new field (Phase 1C territory), OR
- Break projection (the compiler iterates `shared_layer_classification` which would no longer exist)

So Phase 1B introduces the canonical Go loader (`LoadRepositoryTopology` in `repository_topology.go`) that can read BOTH schemas. Phase 1C then upgrades the live YAML to v2 atomically with the projection rewrite.

Until Phase 1C lands, both shapes co-exist conceptually:

- v1: the actual content of `runtime/repository-topology.yaml`, consumed by the live projection rule
- v2: the target shape, exercised by the new loader's tests against in-memory fixtures and by future migrations

### What this is NOT

This file is **not** an architecture decision record. The architectural decision (consumer_tracking freeze + owner/purpose addition) was made in plan `2026-06-06-1800-sanitization-mechanical-enforcement.md` review on 2026-06-08. This file only tracks **adoption status** of that decision.

Do not promote this file to ADR. When all surfaces have migrated to v2, replace the "Current framework runtime state" table with a one-line "fully migrated as of `<commit>`" note. When v3 supersedes v2 (e.g. severity-stratified subtrees), this file rotates: the v2 section becomes the new "current state" and a v3 section gets added.

### Cross-links

- Live YAML: `runtime/repository-topology.yaml` (still v1; do not edit shape until Phase 1C)
- Canonical loader: `scripts/ai-skill-cli/internal/app/repository_topology.go`
- Loader tests: `scripts/ai-skill-cli/internal/app/repository_topology_test.go`
- Parent plan: [`plans/active/2026-06-06-1800-sanitization-mechanical-enforcement.md`](../plans/active/2026-06-06-1800-sanitization-mechanical-enforcement.md) §Phase 1B
- Sibling pattern: [`metadata/project/migration-notes.md`](../metadata/project/migration-notes.md) (Phase 1A migration notes, same trajectory shape applied to the project metadata schema)
- Legacy consumer source: `scripts/ai-skill-cli/internal/app/sanitization_scan.go::repositoryTopologyRow` (DO NOT modify until Phase 1D)
