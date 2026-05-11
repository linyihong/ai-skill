# Model Context Report

本檔由 `ruby scripts/generate-model-context-report.rb --write` 產生，依 `knowledge/runtime/routing-registry.yaml` 的 model 欄位整理 model-aware context loading view。

## Source Surfaces

| Surface | Path | Purpose |
| --- | --- | --- |
| Routing registry | [`routing-registry.yaml`](routing-registry.yaml) | 提供每條 route 的 model profile 與 compression level。 |
| Model profiles | [`../../models/profiles/README.md`](../../models/profiles/README.md) | 定義 `small`、`large`、`specialized` 的讀取深度與 guardrails。 |
| Compression strategy | [`../../models/compression/README.md`](../../models/compression/README.md) | 定義 `summary-first`、`source-backed`、`graph-assisted` 等壓縮層級。 |

## Profile View

### `large`

| Route | Primary source | Compression | Reason |
| --- | --- | --- | --- |
| `route.bootstrap.ai-skill` | `README.md` | `source-backed` | Ai-skill bootstrap 和 repository close-loop 需要完整 source 與 required dependencies。 |
| `route.governance.durable-goal-boundary` | `shared-rules/conversation-goal-ledger.md` | `source-backed` | Goal 刪除與 durable planning gate 會影響 long-term state，需讀 shared rule source。 |
| `route.metadata.knowledge-atom-schema` | `metadata/schema.md` | `source-backed` | 建立或修改 atom metadata 時需讀 schema 與子規則全文。 |
| `route.feedback.promotion-pipeline` | `feedback/promotion/README.md` | `source-backed` | Promotion / downgrade 會影響 shared rules、skill history、intelligence 與 runtime surfaces，需要讀 full source 與 validation gates。 |
| `route.models.model-aware-routing` | `models/profiles/README.md` | `source-backed` | 修改 model strategy 會影響 runtime routing，需讀 profiles、compression 與 routing source。 |

### `small`

| Route | Primary source | Compression | Reason |
| --- | --- | --- | --- |
| `route.runtime.context-loading` | `runtime/routing/README.md` | `summary-first` | routing decision 可先用 registry、index、summary；修改 source 時再升級。 |

### `specialized`

| Route | Primary source | Compression | Reason |
| --- | --- | --- | --- |
| `route.skill.apk-analysis` | `skills/apk-analysis/SKILL.md` | `source-backed` | APK analysis 需要 skill entrypoint、workflow 與 domain-specific technique routing。 |
| `route.intelligence.apk-highest-leverage-path` | `intelligence/engineering/apk-analysis/highest-leverage-analysis-path.md` | `source-backed` | APK route selection 需要 domain workflow、feedback source 與 intelligence judgment 一起判斷。 |

## Compression View

| Compression level | Routes | Escalation note |
| --- | --- | --- |
| `source-backed` | `route.bootstrap.ai-skill`, `route.governance.durable-goal-boundary`, `route.metadata.knowledge-atom-schema`, `route.skill.apk-analysis`, `route.intelligence.apk-highest-leverage-path`, `route.feedback.promotion-pipeline`, `route.models.model-aware-routing` | 需要 primary source 與 required dependencies；適合 writeback、migration 或 domain work。 |
| `summary-first` | `route.runtime.context-loading` | 適合先用 registry / summary 判斷 relevance；修改 source 時升級。 |

## Agent Output Shape

使用本 report 決定 model-aware loading 時，回報：

```text
Profile:
Compression level:
Primary source:
Summaries used:
Required full sources:
Deferred sources:
Escalation trigger:
Validation signal:
```

## Validation

- 產生前應先確認 `routing-registry.yaml` 可通過 `ruby scripts/validate-knowledge-runtime.rb`。
- 產生後應重新執行 `ruby scripts/validate-knowledge-runtime.rb`，檢查本 report links。
- 本報告是 generated view，不取代 `models/profiles/README.md`、`models/compression/README.md` 或 routing registry。
