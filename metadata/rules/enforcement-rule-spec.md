# Enforcement Rule Metadata Spec

本文件定義 Enforcement Rule 專屬的 metadata 規格，繼承 [`metadata/schema.md`](../schema.md) 的 Knowledge Atom schema，並新增 enforcement-specific 欄位。

## Schema 狀態

| 欄位 | 值 |
| --- | --- |
| Schema version | `enforcement-rule/v1` |
| Lifecycle | `candidate` |
| Owner layer | `metadata/rules/` |
| Primary consumers | `runtime/router/`, `knowledge/graphs/rules/`, `governance/lifecycle/` |

## 繼承欄位（來自 Knowledge Atom schema）

Enforcement Rule metadata **必須**包含 Knowledge Atom schema 的所有必填欄位：

| 欄位 | 型別 | 用途 |
| --- | --- | --- |
| `id` | string | 穩定識別碼，格式：`enforcement.<rule-short-name>`（例如 `enforcement.rule-weight`） |
| `title` | string | 給人讀的標題 |
| `type` | enum | 固定為 `rule` |
| `layer` | enum | 固定為 `enforcement` |
| `source_path` | string | 指向 `enforcement/<filename>.md` |
| `summary` | string | 一到兩句說明此 rule 提供 agent 什麼能力或判斷 |
| `domains` | string array | 適用的 domain |
| `tags` | string array | 檢索標籤 |
| `status` | enum | lifecycle 狀態 |
| `priority` | enum | 載入優先序（P0/P1/P2/P3） |
| `confidence` | enum | 證據信心 |
| `stability` | enum | 預期變動速度 |
| `context_cost` | object | 閱讀成本與載入策略 |
| `when_to_read` | string | 載入此 rule 的觸發條件 |
| `validation` | string | 如何確認此 rule 仍然最新且可安全使用 |

## 新增欄位（Enforcement Rule 專屬）

| 欄位 | 型別 | 必填 | 用途 |
| --- | --- | --- | --- |
| `activation_conditions` | object | yes | 結構化觸發條件（見下方詳細定義） |
| `always_apply` | boolean | yes | 是否每個 session 都必須載入（Core Bootstrap） |
| `scope` | enum | yes | 適用範圍：`agent-behavior`、`content-governance`、`workflow`、`tool-usage`、`cross-skill` |
| `deprecated_by` | string | no | 被哪條 rule 取代（rule ID），僅 `status: deprecated` 時填寫 |
| `rule_id` | string | yes | 與 `id` 相同，但明確標示為 rule identifier（供程式化查詢） |

### `activation_conditions` 詳細定義

`activation_conditions` 是 object，包含以下子欄位：

| 欄位 | 型別 | 必填 | 用途 |
| --- | --- | --- | --- |
| `when` | string array | yes | 觸發此 rule 的具體情境描述（對應 `runtime/router/activation-rules.yaml` 的 `activation.when`） |
| `load_strategy` | enum | yes | `preload`（Core Bootstrap）、`lazy`（依條件 activate）、`on_demand`（使用者要求才讀） |
| `priority` | enum | yes | 載入優先序（P0/P1/P2/P3），與頂層 `priority` 一致 |
| `estimated_tokens` | integer | yes | 精確 token 估算 |

## YAML 範本

```yaml
# Enforcement Rule Metadata
# 繼承 Knowledge Atom schema 必填欄位 + Enforcement Rule 專屬欄位
id: enforcement.<rule-short-name>
title: <Rule Title>
schema_version: enforcement-rule/v1
type: rule
layer: enforcement
source_path: enforcement/<filename>.md
summary: <Summary>
domains:
  - repo-governance
tags:
  - <tag1>
  - <tag2>
status: validated
priority: P2
confidence: high
stability: stable
context_cost:
  estimated_tokens: <number>
  load_strategy: lazy
  cacheable: true
  ttl:
    task: 1
when_to_read: <Trigger condition description>
validation: <Validation description>
complexity: low
depends:
  - <dependency rule ID or path>
related:
  - <related rule ID or path>
conflicts: []
replaces: []
models:
  small:
  large:
  specialized:
checklist: []
runtime_notes:
governance_notes:

# Enforcement Rule 專屬欄位
activation_conditions:
  when:
    - <situation 1>
    - <situation 2>
  load_strategy: lazy
  priority: P2
  estimated_tokens: <number>
always_apply: false
scope: agent-behavior
deprecated_by:
rule_id: enforcement.<rule-short-name>
```

## 範例：rule-weight

```yaml
id: enforcement.rule-weight
title: 規則權重與衝突優先序
schema_version: enforcement-rule/v1
type: rule
layer: enforcement
source_path: enforcement/rule-weight.md
summary: 當 enforcement rules、skill workflow、tool adapter、使用者目標或效率規則看似衝突時，依安全/source/validation/user-goal/tool adapter/效率的權重排序處理。
domains:
  - repo-governance
tags:
  - rule-weight
  - conflict-resolution
  - priority
status: validated
priority: P0
confidence: high
stability: stable
context_cost:
  estimated_tokens: 300
  load_strategy: preload
  cacheable: true
  ttl:
    session: 1
when_to_read: 每個 session 啟動時必讀（Core Bootstrap）
validation: 規則權重順序與 enforcement/README.md 的 Runtime Activation Model 一致
complexity: low
depends:
  - enforcement/dependency-reading.md
related:
  - metadata/rules/conflict-matrix.yaml
conflicts: []
replaces: []
models:
  small: 只使用權重順序表
  large: 完整讀取所有範例
checklist:
  - 依 rule-weight.md 的權重順序判斷衝突
  - 不確定時依 rule-weight.md 的處理流程
runtime_notes: Core Bootstrap 規則，每個 session 預載入
governance_notes: 修改權重順序時需同步更新 enforcement/README.md 的索引

activation_conditions:
  when:
    - 任何規則看似衝突時
    - session 啟動時（Core Bootstrap）
  load_strategy: preload
  priority: P0
  estimated_tokens: 300
always_apply: true
scope: agent-behavior
deprecated_by:
rule_id: enforcement.rule-weight
```

## 範例：linked-updates

```yaml
id: enforcement.linked-updates
title: 連動更新
schema_version: enforcement-rule/v1
type: rule
layer: enforcement
source_path: enforcement/linked-updates.md
summary: 全庫必須連動更新規則：改一處影響多處時，相關文件必須同步更新或明確檢查。
domains:
  - repo-governance
tags:
  - linked-updates
  - consistency
  - multi-file-change
status: validated
priority: P1
confidence: high
stability: stable
context_cost:
  estimated_tokens: 800
  load_strategy: lazy
  cacheable: true
  ttl:
    task: 1
when_to_read: 當進行 multi-file change 或 architecture update 時
validation: linked-updates.md 中的依賴表與 enforcement/ 實際檔案結構一致
complexity: medium
depends:
  - enforcement/dependency-reading.md
related:
  - enforcement/rule-weight.md
  - knowledge/graphs/rules/
conflicts: []
replaces: []
models:
  small: 使用常見連動關係表
  large: 完整讀取所有連動規則
checklist:
  - 修改前檢查 linked-updates.md 的依賴表
  - 修改後執行連動更新
runtime_notes: Lazy-load，P1 優先權
governance_notes: 新增或刪除 enforcement rule 時需同步更新 linked-updates.md

activation_conditions:
  when:
    - 修改會影響多個檔案
    - 架構重構或目錄搬遷
    - 新增或刪除 enforcement rule
  load_strategy: lazy
  priority: P1
  estimated_tokens: 800
always_apply: false
scope: workflow
deprecated_by:
rule_id: enforcement.linked-updates
```

## 驗證規則

1. 每個 `.yaml` 檔案的 `source_path` 必須指向 `enforcement/` 下的對應檔案。
2. `id` 格式必須為 `enforcement.<rule-short-name>`（小寫 kebab-case）。
3. `layer` 必須為 `enforcement`。
4. `type` 必須為 `rule`。
5. `always_apply: true` 的 rule 必須有 `load_strategy: preload`。
6. `deprecated_by` 僅在 `status: deprecated` 時可填寫。
7. `activation_conditions.when` 中的描述應與 `runtime/router/activation-rules.yaml` 的 `activation.when` 一致。
