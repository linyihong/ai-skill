# Knowledge Atom Metadata Schema

本文件定義第一版 Knowledge Atom metadata。目標是讓 `knowledge/indexes/`、未來 summaries / graphs / runtime routing，以及各層候選知識能用同一組欄位描述「這份知識是什麼、何時讀、可信度如何、依賴什麼、適合哪些模型」。

本 schema 先服務 navigation 與 migration planning，不要求立即把既有 `skills/` 或 `shared-rules/` 全部轉成 atom。

## Schema 狀態

| 欄位 | 值 |
| --- | --- |
| Schema version | `knowledge-atom/v1` |
| Lifecycle | `candidate` |
| Owner layer | `metadata/` |
| Primary consumers | `knowledge/indexes/`, `runtime/`, `governance/` |

## Required Fields

每個 Knowledge Atom candidate 至少要能填寫下列欄位。

| Field | Type | Required | Purpose |
| --- | --- | --- | --- |
| `id` | string | yes | Stable identifier, using lowercase kebab-case. Prefer `<layer>.<domain>.<short-name>`. |
| `title` | string | yes | Human-readable title. |
| `type` | enum | yes | Atom type: `rule`, `workflow`, `analysis-method`, `intelligence`, `template`, `checklist`, `index`, `schema`, `tool-adapter`, `failure-pattern`, `reference`. |
| `layer` | enum | yes | Primary layer: `analysis`, `intelligence`, `workflow`, `runtime`, `memory`, `feedback`, `models`, `governance`, `knowledge`, `metadata`, `shared-rules`, `skills`, `ai-tools`, `scripts`, `architecture`. |
| `source_path` | string | yes | Canonical repository-relative path to the source file or directory. |
| `summary` | string | yes | One or two sentences explaining what the atom gives an agent. |
| `domains` | string array | yes | Domains or capabilities this atom applies to, for example `apk-analysis`, `app-development`, `travel-planning`, `repo-governance`. |
| `tags` | string array | yes | Retrieval tags. Use lowercase kebab-case. |
| `status` | enum | yes | Lifecycle status: `temporary`, `candidate`, `validated`, `stable`, `deprecated`. |
| `priority` | enum | yes | Loading priority: `P0`, `P1`, `P2`, `P3`. Use the same vocabulary as goal ledger priorities. |
| `confidence` | enum | yes | Evidence confidence: `low`, `medium`, `high`. |
| `stability` | enum | yes | Expected change rate: `experimental`, `evolving`, `stable`, `legacy`. |
| `context_cost` | enum | yes | Approximate read cost: `low`, `medium`, `high`. |
| `when_to_read` | string | yes | Trigger condition for loading this atom. |
| `validation` | string | yes | How an agent can verify the atom is current and safe to use. |

## Optional Fields

Use optional fields when they improve routing, conflict handling, or model-aware loading.

| Field | Type | Purpose |
| --- | --- | --- |
| `complexity` | enum | `low`, `medium`, `high`; helps model routing and compression strategy. |
| `depends` | string array | Atom IDs or paths that must be read first. |
| `related` | string array | Atom IDs or paths that may be useful but are not required. |
| `conflicts` | string array | Atom IDs, paths, or rule categories that may conflict. |
| `replaces` | string array | Deprecated or superseded atom IDs / paths. |
| `models` | object | Model suitability notes, for example `small`, `large`, `specialized`. |
| `checklist` | string array | Short checklist for low-context or small-model usage. |
| `runtime_notes` | string | Notes for dynamic loading, compression, or orchestration. |
| `governance_notes` | string | Lifecycle, review cadence, ownership, or deprecation notes. |

## Controlled Values

### `type`

- `rule`: executable policy or operating rule.
- `workflow`: task execution flow.
- `analysis-method`: observation, decomposition, or extraction method.
- `intelligence`: engineering judgment, trade-off, anti-pattern, or reusable domain knowledge.
- `template`: reusable document or prompt template.
- `checklist`: validation or review checklist.
- `index`: navigation or routing index.
- `schema`: metadata or contract schema.
- `tool-adapter`: tool-specific execution guidance.
- `failure-pattern`: reusable prevention pattern for known agent failures.
- `reference`: roadmap, architecture note, or background reference.

### `priority`

- `P0`: safety, secrets, source-of-truth, data-loss, or destructive-action control.
- `P1`: active goal closure, required bootstrap, canonical writeback, or validation gate.
- `P2`: task-relevant workflow, domain intelligence, or migration guidance.
- `P3`: optional optimization, cleanup, examples, or background context.

### `context_cost`

- `low`: quick index, checklist, or short rule.
- `medium`: focused workflow or single-purpose reference.
- `high`: broad document, multi-section workflow, or source that should be read only when strongly relevant.

## YAML Template

```yaml
id:
title:
schema_version: knowledge-atom/v1
type:
layer:
source_path:
summary:
domains: []
tags: []
status: candidate
priority:
confidence:
stability:
context_cost:
when_to_read:
validation:
complexity:
depends: []
related: []
conflicts: []
replaces: []
models:
  small:
  large:
  specialized:
checklist: []
runtime_notes:
governance_notes:
```

## Example Atom

```yaml
id: knowledge.indexes.task-routing
title: Knowledge navigation task routing index
schema_version: knowledge-atom/v1
type: index
layer: knowledge
source_path: knowledge/indexes/README.md
summary: Routes task intents to the first canonical source an agent should read, with related sources and validation signals.
domains:
  - repo-governance
  - knowledge-navigation
tags:
  - routing
  - navigation
  - context-loading
status: candidate
priority: P2
confidence: medium
stability: evolving
context_cost: low
when_to_read: Use when an agent needs to find task-relevant Ai-skill knowledge without loading every skill or shared rule.
validation: Links resolve, primary sources remain canonical, and roadmap status matches current repository structure.
complexity: low
depends:
  - README.md
  - shared-rules/dependency-reading.md
related:
  - metadata/schema.md
  - architecture/next-stage-upgrade-plan.md
conflicts: []
models:
  small: Use the task routing table and validation signal only.
  large: Read related sources when the task spans multiple layers.
checklist:
  - Match task intent to a row.
  - Read primary source first.
  - Load related sources only when needed.
runtime_notes: Suitable as a low-cost routing atom before deeper context loading.
governance_notes: Update when new top-level layers, skills, or canonical entrypoints are added.
```

## Validation Rules

- `source_path` must point to a canonical repository path, not a local tool mirror.
- `depends`, `related`, `conflicts`, and `replaces` should use atom IDs when available; repository-relative paths are acceptable during migration.
- `summary`, `when_to_read`, and `validation` must be specific enough for an agent to decide whether to load the atom.
- Do not use metadata to override executable shared rules. If a rule conflict exists, follow `shared-rules/rule-weight.md`.
- Do not mark an atom `stable` until it has survived at least one real use or review with validation evidence.
