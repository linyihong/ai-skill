# Executable Contract Schema

本文件定義 owner-layer executable YAML contract 的共同欄位。它用來區分「metadata YAML」與「agent 可直接執行的流程契約」。

YAML 放置政策的 canonical source 是 [`../governance/lifecycle/executable-contract-boundary.md`](../governance/lifecycle/executable-contract-boundary.md)。本檔只定義欄位 schema，不決定 contract 應放在 `enforcement/`、`workflow/`、`governance/`、`ai-tools/`、`metadata/` 或 runtime projection table。

## 適用範圍

當 owner-layer YAML 會影響 agent 執行時，必須設定：

```yaml
runtime_projection:
  enabled: true
```

設定後，runtime compiler 會將該 YAML 投影到 `runtime/runtime.db` 的 `generated_surfaces`。一般 metadata、routing、graph 或背景說明 YAML 不應設定此欄位。`metadata/rules/*.yaml` 預設仍是 metadata-only；只有符合本 schema 並明確啟用 runtime projection 時，才是 executable contract 的例外放置。

## 必填欄位

| 欄位 | 用途 |
| --- | --- |
| `schema_version` | 新 contract 使用 `executable-contract/v1`；舊 contract 可在遷移期間保留舊版本，但不得新增不完整 contract。 |
| `id` | 穩定 contract ID，例如 `workflow.documentation.execution_flow`。 |
| `title` | 人類可讀標題。 |
| `owner_layer` | 擁有層：`governance`、`enforcement`、`workflow`、`ai-tools` 或 `metadata`。 |
| `source_markdown` | companion Markdown source。 |
| `status` | `active`、`draft`、`deprecated`。 |
| `contract_type` | `policy_gate`、`workflow_flow`、`onboarding_flow`、`promotion_gate`、`schema_contract`。 |
| `blocking_level` | `blocking`、`advisory`、`informational`。 |
| `runtime_projection.target_key` | `generated_surfaces` 的穩定 key。 |
| `activation` | 何時必須讀取或執行此 contract。 |

## 執行欄位

新 contract 至少需要包含一組可執行欄位：

- `steps`
- `gates`
- `required_sources`
- `required_evidence`
- `success_criteria`
- `failure_modes`
- `final_status_report`

若某欄位不適用，應在鄰近欄位寫明原因，不要以空 contract 通過 validation。

## YAML 化判斷

| 文件形態 | 決策 |
| --- | --- |
| 有 ordered steps、activation、required reads、blocking gates、failure actions、required evidence 或 final report | 建立 executable YAML contract。 |
| 只有哲學、背景、tradeoff、設計理由或人類導讀索引 | 保持 Markdown-only。 |
| 只有 rule metadata、routing summary 或 front-matter | 不算 executable contract，除非補齊本 schema 並啟用 runtime projection。 |
| 其他專案的 docs / runbook / ADR / workflow 定義 agent 要反覆執行的 gate 或流程 | 建立 project-local companion YAML 或等價 structured contract；欄位語意應映射到本 schema。 |

## 驗證

`ai-skill runtime validate` 必須確認：

1. 啟用 runtime projection 的 YAML 可解析。
2. `runtime_projection.target_key` 在 `runtime/runtime.db generated_surfaces` 中存在且狀態為 `synced`。
3. `generated_surfaces.data` 包含來源 YAML 的 `id`、`runtime_projection` 與至少一組 execution-bearing 欄位。
4. 新增 `schema_version: executable-contract/v1` 的 contract 時，必填欄位齊全。

← [回到 Metadata](README.md)
