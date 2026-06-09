# Repository Topology — Migration Notes

> **Purpose**: a small living record of where `runtime/repository-topology.yaml` sits between the v1 schema (currently in production) and the v2 schema (designed but not yet adopted in projection). Reviewers three months from now should be able to read this file alone and know whether the live YAML's v1 shape is deliberate (waiting for Phase 1C) or accidental (forgotten cleanup).
>
> Update this file every time a consumer or producer migrates between schema versions, or when a new consumer is added to the runtime.

## `runtime/repository-topology.yaml` (v2 in production as of Phase 1C₁)

### Phase 1C₁ migration landed 2026-06-09

The live `runtime/repository-topology.yaml` was upgraded from v1 to v2 in commit landing this section. `runtime_compiler.go` line 339's tuple-format projection rule was removed; the file is now compiled via `compileRepositoryTopology` in `repository_topology_compile.go`, which writes `repository_topology` rows with JSON content carrying BOTH v1 keys (`subtree`, `shared`) AND v2 keys (`path`, `shared_layer`, `owner`, `purpose`). The legacy `sanitization_scan.go::repositoryTopologyRow` continues to read the table unchanged; Phase 1D will retire its v1 dependency.

The v1 / v2 sections below are kept as **historical reference** so future readers can see what the schema looked like before and after the upgrade. They are no longer prescriptive — the live state is v2.

### v1 schema (historical; production through Phase 1B)

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

### Current framework runtime state (as of Phase 1C₁ landing)

| Surface | Schema in use | Reads via |
| --- | --- | --- |
| `runtime/repository-topology.yaml` (on disk) | **v2** | — |
| `runtime_compiler.go` projection wiring | calls `compileRepositoryTopology` (custom function) | — |
| `repository_topology_compile.go::compileRepositoryTopology` | reads v2 via `LoadRepositoryTopology`; writes dual-shape JSON | direct YAML I/O via Phase 1B loader |
| `runtime.db.repository_topology` (projected table) | rows have JSON with BOTH v1 keys (`subtree`, `shared`) AND v2 keys (`path`, `shared_layer`, `owner`, `purpose`) | populated by `compileRepositoryTopology` |
| `sanitization_scan.go::repositoryTopologyRow` (legacy consumer) | reads `subtree` + `shared` from JSON content (still works via backward-compat dual-shape) | unchanged; will migrate in Phase 1D |

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
| Schema spec (this file) | v2 documented | v2 documented | Phase 1B (landed) | framework maintainer |
| Loader `repository_topology.go` | v1 read + v2 read/write | v1 read + v2 read/write | Phase 1B (landed) | framework maintainer |
| Live YAML `runtime/repository-topology.yaml` | **v2** | v2 | **Phase 1C₁ (landed)** | framework maintainer |
| Projection rule (compiler) | calls `compileRepositoryTopology` | calls `compileRepositoryTopology` | **Phase 1C₁ (landed)** | framework maintainer |
| Live consumer `sanitization_scan.go` | reads `subtree` + `shared` from SQLite via backward-compat dual-shape JSON | reads v2 fields (`owner`/`purpose`) directly | Phase 1D | framework maintainer |
| Registry / scenarios | — | enforcement-registry binding | Phase 4 | framework maintainer |

### History — why Phase 1B kept the v1 file on disk

Phase 1B (commits `b09359a` + `acf2693`) introduced the canonical loader (`LoadRepositoryTopology`) but explicitly did NOT upgrade the live YAML. The reason: `runtime_compiler.go` line 339 was a tuple-format projection rule hard-coded against v1 field names (`shared_layer_classification`, `subtree`). Upgrading the YAML in isolation would have broken projection.

Phase 1C₁ (this section's landing commit) executes the upgrade ATOMICALLY:

1. Live YAML upgraded v1 → v2 (in-place)
2. `runtime_compiler.go` line 339 tuple entry removed
3. New `compileRepositoryTopology` function added to the compile pipeline
4. JSON content carries BOTH v1 and v2 keys so legacy scanner reads correctly

The atomic discipline matches the migration-notes §Schema version precedence (mixed-shape files) rule: partial-shape transitions are silent-drop hazards; the migration must be one commit, not several.

Until Phase 1D lands, the v1 keys in JSON content remain — they are the contract that lets `sanitization_scan.go::repositoryTopologyRow` keep reading without modification. Phase 1D will:

- Either drop the v1 keys (and migrate the scanner's `repositoryTopologyRow` to read v2 keys directly)
- Or keep the v1 keys but route the scanner through `LoadRepositoryTopology` for stronger typing

The choice is deferred to Phase 1D's design pass.

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
