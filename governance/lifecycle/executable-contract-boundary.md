# Executable Contract Boundary

This document defines when a governance, enforcement, or workflow document needs a machine-readable YAML contract and how that contract reaches runtime.

## Core Rule

Source ownership stays with the layer that owns the concept. Runtime execution surfaces are projected into `runtime/runtime.db`.

```text
owner-layer Markdown / YAML contract
  -> runtime compiler
  -> runtime.db generated_surfaces / projection tables
```

Do not move governance, enforcement, or workflow source files into `runtime/` just because they affect execution. `runtime/` remains the runtime engine and SQLite registry boundary.

## When YAML Is Required

A document needs a YAML contract when agents must execute it as a workflow or gate.

Required signals:

- Ordered steps
- Trigger or activation conditions
- Required reads or dependencies
- `depends_on` relationships
- Exit conditions
- Blocking gates
- Required evidence
- Failure actions
- Final status/report requirements

If a document only explains philosophy, background, tradeoffs, or design rationale, keep it Markdown-only unless a later workflow extracts executable gates from it.

## Placement Rule

| Source type | YAML contract location | Runtime projection |
| --- | --- | --- |
| Governance lifecycle flow | `governance/**/*.yaml` | Required when execution-affecting |
| Enforcement policy contract | `enforcement/**/*.yaml` or `metadata/rules/*.yaml` | Required when execution-affecting |
| Workflow execution flow | `workflow/**/*.yaml` | Required when execution-affecting |
| Runtime internal config | `runtime.db` canonical documents | Already runtime-owned |
| Philosophy / rationale / ADR | Markdown only | Not required |

## Runtime Projection Rule

YAML contracts that affect agent execution must include:

```yaml
runtime_projection:
  enabled: true
  target_key: governance.example.contract
  surface: generated_surfaces
```

The compiler only projects contracts that opt in with `runtime_projection.enabled: true`. This prevents ordinary metadata, graph, and validation YAML from becoming runtime noise.

## Schema Rule

New executable contracts should follow [`../../metadata/executable-contract-schema.md`](../../metadata/executable-contract-schema.md). Metadata YAML is not an executable contract unless it defines contract fields such as `contract_type`, `blocking_level`, `activation`, execution-bearing fields, and `runtime_projection.enabled: true`.

## Agent Rule

When a Markdown file says a process must be run as a workflow, the agent must load the companion YAML contract first, then use the Markdown for explanation and maintenance context.

If no YAML contract exists but the document has executable signals, the agent must treat that as a linked-update gap and either create the contract or record why it is not applicable.

## Initial Contract Inventory

The initial inventory lives in [`executable-contract-boundary.yaml`](executable-contract-boundary.yaml). It marks current documents as:

- `contract_exists`: YAML already exists.
- `candidate`: needs YAML contract next.
- `markdown_only`: intentionally not executable.
- `runtime_db_only`: runtime-internal config already belongs in SQLite.

## Related

- [`knowledge-update-flow.yaml`](knowledge-update-flow.yaml)
- [`compiler-philosophy.md`](compiler-philosophy.md)
- [`../../runtime/README.md`](../../runtime/README.md)
- [`../../scripts/ai-skill-cli/internal/app/runtime_compiler.go`](../../scripts/ai-skill-cli/internal/app/runtime_compiler.go)
