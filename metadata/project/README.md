# Project Metadata Schemas

`metadata/project/` defines the canonical schemas that **downstream projects** use to declare project-local metadata consumed by Ai-skill framework executors.

This sub-layer is distinct from the rest of `metadata/`:

- `metadata/rules/`, `metadata/architecture/`, `metadata/recovery/` etc. declare **framework-level** controlled vocabulary, rule metadata, and routing/control surfaces — they are read by Ai-skill itself.
- `metadata/project/` declares schemas that **downstream projects** populate. Each schema describes the SHAPE of a file the project author creates at their own `<PROJECT_ROOT>/...` path.

## Current schemas

| Schema | Project file | Consumer | Status |
| --- | --- | --- | --- |
| [`ai-skill-project-schema.yaml`](ai-skill-project-schema.yaml) | `<PROJECT_ROOT>/.ai-skill-project.yaml` | sanitization mechanical enforcement (Phase 1C/1D, in design) | candidate |

## Boundary rules

- Schemas in this directory describe project-author-facing surfaces. They are NOT projection surfaces themselves.
- The framework parser implementation lives in `scripts/ai-skill-cli/internal/app/` (e.g. `project_metadata.go` for the `.ai-skill-project.yaml` schema).
- Migration trajectory of any consumer of these schemas belongs in `migration-notes.md`, not inline in the schema declarations. Schema YAML stays focused on the data shape.

## When to add a new schema here

Add a new schema YAML under `metadata/project/` when:

1. The framework introduces a new executor that consumes project-local declarative metadata
2. The metadata format is intended for downstream project authors to write (not for the framework to generate)
3. The framework needs a parser + validation layer separate from the executor

Do NOT add schemas here that describe internal Ai-skill runtime surfaces. Those belong in `runtime/` or `governance/lifecycle/`.

## Cross-links

- [`migration-notes.md`](migration-notes.md) — current adoption status of each schema in the framework runtime
- [`example-ai-skill-project.yaml`](example-ai-skill-project.yaml) — illustrative downstream project metadata file
- [`../README.md`](../README.md) — parent `metadata/` index

← [Back to metadata](../README.md)
