# Knowledge Runtime Surfaces

`knowledge/runtime/` 定義給 runtime 使用的 knowledge view：如何從 task intent 找到 primary source、required dependencies、candidate sources、source-of-truth gate 與 validation signal。本層目前仍是文件與 registry surface，不是自動化 runtime。

## 目前入口

- [`routing-registry.yaml`](routing-registry.yaml)：第一版 machine-readable routing registry，包含 8 筆 sample routing records。
- [`refresh-policy.yaml`](refresh-policy.yaml)：generated summaries、graphs 與 routing registry 的 refresh / revalidate / downgrade 規則。
- [`runtime-report.md`](runtime-report.md)：由 generator 產生的目前 routes、summaries、graphs 與 refresh decisions 概覽。
- [`model-context-report.md`](model-context-report.md)：由 generator 產生的 model profile / compression loading view。
- [`../../scripts/validate-knowledge-runtime.rb`](../../scripts/validate-knowledge-runtime.rb)：檢查 generated surfaces 的 YAML / Markdown 格式、必要欄位與 canonical path。
- [`../../scripts/generate-knowledge-runtime-report.rb`](../../scripts/generate-knowledge-runtime-report.rb)：從 runtime surfaces 產生 deterministic report。
- [`../../scripts/generate-model-context-report.rb`](../../scripts/generate-model-context-report.rb)：從 routing registry 的 model 欄位產生 context loading report。

## Runtime Inputs

| Input | Source |
| --- | --- |
| Task intent routing | `knowledge/indexes/README.md` |
| Machine-readable routing registry | `knowledge/runtime/routing-registry.yaml` |
| Generated refresh policy | `knowledge/runtime/refresh-policy.yaml` |
| Generated runtime report | `knowledge/runtime/runtime-report.md` |
| Generated model context report | `knowledge/runtime/model-context-report.md` |
| Atom metadata | `metadata/schema.md` |
| Ranking rules | `metadata/ranking/README.md` |
| Confidence rules | `metadata/confidence/README.md` |
| Compatibility rules | `metadata/compatibility/README.md` |
| Model profiles | `models/profiles/README.md` |
| Compression strategy | `models/compression/README.md` |
| Lifecycle and validation gates | `governance/lifecycle/README.md`, `governance/validation/README.md` |
| Runtime routing design | `runtime/routing/README.md` |

## Runtime View Format

Runtime view 應回答：

| Field | Purpose |
| --- | --- |
| `task_intent` | Agent 正在嘗試完成的任務意圖。 |
| `primary_source` | 第一個應讀的 canonical source。 |
| `required_dependencies` | 必讀的 shared rules、skill entrypoints 或 metadata。 |
| `candidate_sources` | 可選的 maps、summaries 或 atoms。 |
| `source_of_truth_gate` | 舊 entrypoint 是否仍勝過 candidate new-layer path。 |
| `ranking_reason` | 為什麼此 source 排第一。 |
| `model_profile` | 建議使用 `small`、`large` 或 `specialized` profile。 |
| `compression_level` | 建議使用 `index-only`、`summary-first`、`checklist-first`、`source-backed` 或 `graph-assisted`。 |
| `validation_signal` | 如何確認這條 route 可安全使用。 |

## Runtime Rules

- Runtime view 不得跳過 required shared-rule bootstrap。
- Lifecycle promotion gates 通過前，runtime view 不得取代舊 skill behavior。
- 低成本 summary 只有在連回 canonical source 時，才可優先使用。
- Tool mirrors 是 deployment surfaces，不是 source paths。
- 故意延後的 sources 必須記錄為 deferred，而不是假裝不需要。
- 任何會修改 canonical source 或需要 close-loop 的任務，都必須使用 `source-backed` 或更高讀取深度。
- Source 變更後，必須依 `refresh-policy.yaml` 判斷 summaries、graphs、routing registry 是 refresh、revalidate、downgrade 或 no update needed。
- 修改 registry、summaries、graphs 或 refresh policy 後，重新執行 `ruby scripts/generate-knowledge-runtime-report.rb --write`。
- 修改 registry model 欄位、model profiles 或 compression strategy 後，重新執行 `ruby scripts/generate-model-context-report.rb --write`。
- 修改 registry、refresh policy、summaries 或 graphs 後，執行 `ruby scripts/validate-knowledge-runtime.rb`，再做 lints、Markdown link check、close-loop dry run 與 commit / push / readback。

## 尚未實作

- Automatic graph construction。
- Generated summaries。
- Machine-readable registry 的自動生成工具。
- Generated refresh 的自動執行工具。
- Model-aware compression output 的更細粒度生成。

這些項目會在 governance、metadata 與 routing surfaces 穩定後再推進。
