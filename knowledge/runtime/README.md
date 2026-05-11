# Knowledge Runtime Surfaces

`knowledge/runtime/` 定義給 runtime 使用的 knowledge view：如何從 task intent 找到 primary source、required dependencies、candidate sources、source-of-truth gate 與 validation signal。本層目前仍是文件與 registry surface，不是自動化 runtime。

## 目前入口

- [`routing-registry.yaml`](routing-registry.yaml)：第一版 machine-readable routing registry，包含 6 筆 sample routing records。

## Runtime Inputs

| Input | Source |
| --- | --- |
| Task intent routing | `knowledge/indexes/README.md` |
| Machine-readable routing registry | `knowledge/runtime/routing-registry.yaml` |
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

## 尚未實作

- Automatic graph construction。
- Generated summaries。
- Machine-readable registry 的自動生成或驗證工具。
- Model-aware compression output 的自動生成。

這些項目會在 governance、metadata 與 routing surfaces 穩定後再推進。
