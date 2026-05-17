# Knowledge Atom Metadata Schema

本文件定義第一版 Knowledge Atom metadata。目標是讓 `knowledge/indexes/`、未來 summaries / graphs / runtime routing，以及各層候選知識能用同一組欄位描述「這份知識是什麼、何時讀、可信度如何、依賴什麼、適合哪些模型」。

本 schema 先服務 navigation 與 migration planning，不要求立即把既有 `skills/` 或 `enforcement/` 全部轉成 atom。

## Schema 狀態

| 欄位 | 值 |
| --- | --- |
| Schema version | `knowledge-atom/v1` |
| Lifecycle | `candidate` |
| Owner layer | `metadata/` |
| Primary consumers | `knowledge/indexes/`, `runtime/`, `governance/` |

## 必填欄位

每個 Knowledge Atom candidate 至少要能填寫下列欄位。

| 欄位 | 型別 | 必填 | 用途 |
| --- | --- | --- | --- |
| `id` | string | yes | 穩定識別碼，使用小寫 kebab-case；優先採 `<layer>.<domain>.<short-name>`。 |
| `title` | string | yes | 給人讀的標題。 |
| `type` | enum | yes | Atom 類型：`rule`、`workflow`、`analysis-method`、`intelligence`、`template`、`checklist`、`index`、`schema`、`tool-adapter`、`failure-pattern`、`reference`。 |
| `layer` | enum | yes | 主要歸屬層：`analysis`、`intelligence`、`workflow`、`runtime`、`memory`、`feedback`、`models`、`governance`、`knowledge`、`metadata`、`enforcement`、`skills`、`ai-tools`、`scripts`、`architecture`。 |
| `source_path` | string | yes | canonical repository-relative source 檔案或資料夾路徑。 |
| `summary` | string | yes | 一到兩句說明此 atom 提供 agent 什麼能力或判斷。 |
| `domains` | string array | yes | 此 atom 適用的 domain 或 capability，例如 `apk-analysis`、`app-development`、`travel-planning`、`repo-governance`。 |
| `tags` | string array | yes | 檢索標籤；使用小寫 kebab-case。 |
| `status` | enum | yes | lifecycle 狀態：`temporary`、`candidate`、`validated`、`stable`、`deprecated`。 |
| `priority` | enum | yes | 載入優先序：`P0`、`P1`、`P2`、`P3`，沿用 goal ledger vocabulary。 |
| `confidence` | enum | yes | 證據信心：`low`、`medium`、`high`。 |
| `stability` | enum | yes | 預期變動速度：`experimental`、`evolving`、`stable`、`legacy`。 |
| `context_cost` | object | yes | 閱讀成本與載入策略（見下方詳細定義）。 |
| `when_to_read` | string | yes | 載入此 atom 的觸發條件。 |
| `validation` | string | yes | agent 如何確認此 atom 仍然最新且可安全使用。 |

## `context_cost` 詳細定義

`context_cost` 是 object，包含以下子欄位：

| 欄位 | 型別 | 必填 | 用途 |
| --- | --- | --- | --- |
| `estimated_tokens` | integer | yes | 精確 token 估算（整數）。 |
| `load_strategy` | enum | yes | 載入策略：`preload`（每個 session 必讀）、`lazy`（依條件 activate）、`on_demand`（使用者要求才讀）。 |
| `cacheable` | boolean | yes | 是否可在 session/conversation 內 cache。 |
| `ttl` | object | no | Context TTL 設定（見下方）。 |
| `breakdown` | object | no | 成本細項，例如 `header_and_navigation`、`core_content`、`examples`、`validation`。 |

### `context_cost.ttl` 子欄位

| 欄位 | 型別 | 用途 |
| --- | --- | --- |
| `session` | integer | 活幾個 session（預設 1）。 |
| `task` | integer | 活幾個 task。 |
| `conversation` | boolean | 是否活整個對話。 |

### `context_cost` 範例

```yaml
context_cost:
  estimated_tokens: 1200
  load_strategy: lazy
  cacheable: true
  ttl:
    task: 1
  breakdown:
    header_and_navigation: 200
    core_content: 700
    examples: 200
    validation: 100
```

## 選填欄位

當欄位能改善 routing、衝突處理、model-aware loading 或 cost-aware routing 時才填寫。

| 欄位 | 型別 | 用途 |
| --- | --- | --- |
| `complexity` | enum | `low`、`medium`、`high`；支援 model routing 與 compression strategy。 |
| `depends` | string array | 使用此 atom 前必須先讀的 atom ID 或路徑。 |
| `related` | string array | 可能有幫助但非必讀的 atom ID 或路徑。 |
| `conflicts` | string array | 可能衝突的 atom ID、路徑或規則類別。 |
| `replaces` | string array | 已 deprecated 或被取代的 atom ID / 路徑。 |
| `models` | object | 模型適用性備註，例如 `small`、`large`、`specialized`。 |
| `checklist` | string array | 低 context 或小模型使用的短檢查清單。 |
| `runtime_notes` | string | dynamic loading、compression、orchestration 或 TTL 備註。 |
| `governance_notes` | string | lifecycle、review cadence、ownership 或 deprecation 備註。 |

## 受控值

### `type`

- `rule`：可執行政策或操作規則。
- `workflow`：任務執行流程。
- `analysis-method`：觀察、拆解或 extraction method。
- `intelligence`：工程判斷、trade-off、anti-pattern 或可重用 domain knowledge。
- `template`：可重用文件或 prompt template。
- `checklist`：驗證或 review checklist。
- `index`：navigation 或 routing index。
- `schema`：metadata 或 contract schema。
- `tool-adapter`：工具專屬執行差異。
- `failure-pattern`：已知 agent failure 的可重用 prevention pattern。
- `reference`：roadmap、architecture note 或背景 reference。

### `priority`

- `P0`：safety、secrets、source-of-truth、data-loss 或 destructive-action control。
- `P1`：active goal closure、required bootstrap、canonical writeback 或 validation gate。
- `P2`：task-relevant workflow、domain intelligence 或 migration guidance。
- `P3`：optional optimization、cleanup、examples 或 background context。

### `context_cost`

- `low`：快速 index、checklist 或短規則。
- `medium`：聚焦 workflow 或單一用途 reference。
- `high`：寬泛文件、多段 workflow，或只有高度相關時才應讀取的 source。

## YAML 範本

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

## Atom 範例

```yaml
id: knowledge.indexes.task-routing
title: Knowledge navigation task routing index
schema_version: knowledge-atom/v1
type: index
layer: knowledge
source_path: knowledge/indexes/README.md
summary: 將 task intent 導向 agent 應先讀取的 canonical source，並提供 related sources 與 validation signals。
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
when_to_read: 當 agent 需要找到 task-relevant Ai-skill knowledge，但不應載入所有 skill 或 enforcement rule 時使用。
validation: Links 可解析、primary sources 仍為 canonical，且 roadmap status 符合目前 repository structure。
complexity: low
depends:
  - README.md
  - enforcement/dependency-reading.md
related:
  - metadata/schema.md
  - plans/archived/2026-05-11-1112-next-stage-upgrade-plan.md
conflicts: []
models:
  small: 只使用 task routing table 與 validation signal。
  large: 任務跨多個 layer 時讀取 related sources。
checklist:
  - 將 task intent 對應到索引列。
  - 先讀 primary source。
  - 只在需要時載入 related sources。
runtime_notes: 適合作為 deeper context loading 前的低成本 routing atom。
governance_notes: 新增 top-level layers、skills 或 canonical entrypoints 時同步更新。
```

## 驗證規則

- `source_path` 必須指向 canonical repository path，不可指向 local tool mirror。
- `depends`、`related`、`conflicts`、`replaces` 有 atom ID 時優先使用 atom ID；migration 期間可使用 repository-relative path。
- `summary`、`when_to_read`、`validation` 必須足夠具體，讓 agent 能判斷是否需要載入此 atom。
- 不可用 metadata 覆蓋可執行 enforcement rules；若規則衝突，依 `enforcement/rule-weight.md`。
- Atom 至少經過一次真實使用或 review，且有 validation evidence 後，才可標記為 `stable`。
