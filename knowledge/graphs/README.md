# Knowledge Graphs

`knowledge/graphs/` will describe relationships between Knowledge Atoms, source files, skills, shared rules, and runtime routing surfaces. During the current phase, this directory defines graph format only; it does not generate a full graph.

## Graph Purpose

Graphs help agents understand:

- Required dependencies.
- Related sources.
- Conflicts.
- Replacement and deprecation paths.
- Promotion flow from old skills into new layers.

## Edge Types

Use these edge labels for future graph records:

| Edge | Meaning |
| --- | --- |
| `depends_on` | Source must be read before the atom can be used. |
| `related_to` | Source may be useful but is not required. |
| `conflicts_with` | Source may conflict and requires rule-weight or governance resolution. |
| `replaces` | New atom supersedes an older source after promotion. |
| `preserves_entrypoint` | New layer path keeps old source reachable. |
| `promotes_from` | Atom was extracted or promoted from an old skill/shared rule. |
| `routes_to` | Index or runtime routing points to the source. |

## Graph Record Format

```yaml
id:
source:
edges:
  - type:
    target:
    reason:
    validation:
status: candidate
```

## Compatibility Rules

- Use canonical repository-relative paths or atom IDs.
- Do not model tool mirror paths as canonical sources.
- If a graph says a new path `replaces` an old path, the lifecycle state must already be promoted or deprecated.
- Candidate maps should use `preserves_entrypoint`, not `replaces`.

## Initial Graph Candidates

| Graph candidate | Purpose | Status |
| --- | --- | --- |
| `apk-analysis-pilot` | Connect `skills/apk-analysis/` to `analysis/apk/`, `workflow/apk-analysis/`, and `intelligence/engineering/apk-analysis/`. | candidate |
| `metadata-navigation` | Connect `knowledge/indexes/README.md` to `metadata/schema.md` and metadata subrules. | candidate |
| `source-boundary` | Connect governance lifecycle to old skill entrypoints and validation gates. | candidate |
