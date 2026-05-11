# Knowledge Runtime Report

本檔由 `ruby scripts/generate-knowledge-runtime-report.rb --write` 產生，彙整 runtime registry、summaries、graphs 與 refresh policy 的目前狀態。

## Source Surfaces

| Surface | Path | Count / Status |
| --- | --- | --- |
| Routing registry | [`routing-registry.yaml`](routing-registry.yaml) | 6 records |
| Refresh policy | [`refresh-policy.yaml`](refresh-policy.yaml) | candidate |
| Summaries | [`../summaries/`](../summaries/) | 4 files |
| Graph records | [`../graphs/`](../graphs/) | 3 files |

## Routing Records

| ID | Primary source | Model | Compression | Validation signal |
| --- | --- | --- | --- | --- |
| `route.bootstrap.ai-skill` | `README.md` | `large` | `source-backed` | Bootstrap required set 已讀，且 git status 已檢查。 |
| `route.governance.durable-goal-boundary` | `shared-rules/conversation-goal-ledger.md` | `large` | `source-backed` | 長期狀態已落到 durable planning 文件，且 active goal 完成後才刪除。 |
| `route.runtime.context-loading` | `runtime/routing/README.md` | `small` | `summary-first` | Primary source、deferred sources、source-of-truth gate 與 validation signal 已可被記錄。 |
| `route.metadata.knowledge-atom-schema` | `metadata/schema.md` | `large` | `source-backed` | 欄位可套用到第一批 atom candidates，且 Markdown links 可解析。 |
| `route.skill.apk-analysis` | `skills/apk-analysis/SKILL.md` | `specialized` | `source-backed` | 舊入口仍可讀，新 reference-first paths 可找到，且未 bulk migrate skill content。 |
| `route.models.model-aware-routing` | `models/profiles/README.md` | `large` | `source-backed` | Profile、compression level、primary source、deferred sources 與 validation signal 可被記錄。 |

## Summary Records

| Atom ID | Lifecycle | File | Summary |
| --- | --- | --- | --- |
| `architecture.apk-analysis-pilot` | `candidate` | [`apk-analysis-pilot.md`](../summaries/apk-analysis-pilot.md) | `apk-analysis` 作為第一個 Workflow / Analysis / Intelligence 分離 pilot 的 migration map。它建立新 reference-first 候選目的地，但保留 `skills/apk-analysis/SKILL.md` 作為 active skill entrypoint。 |
| `governance.goal-ledger-boundary` | `validated` | [`goal-ledger-boundary.md`](../summaries/goal-ledger-boundary.md) | `.agent-goals/` 只保存 active conversation goals；長期 roadmap、phase、migration、promotion、deprecation 與治理狀態必須落到 durable planning 文件。 |
| `metadata.schema.knowledge-atom` | `validated` | [`metadata-schema.md`](../summaries/metadata-schema.md) | Knowledge Atom metadata schema v1，定義 atom 的必填欄位、選填欄位、受控值、YAML 範本與驗證規則。 |
| `root.bootstrap.ai-skill` | `validated` | [`root-bootstrap.md`](../summaries/root-bootstrap.md) | Ai-skill 工作的 bootstrap 入口。Root README 定義 repository layout 與 reference-first 工作流；shared-rules README 定義 Default Bootstrap 與依任務補讀規則。 |

## Graph Records

| ID | Source | Status | Edges | File |
| --- | --- | --- | --- | --- |
| `graph.apk-analysis-pilot` | `architecture/apk-analysis-pilot-migration.md` | `candidate` | 6 | [`apk-analysis-pilot.yaml`](../graphs/apk-analysis-pilot.yaml) |
| `graph.metadata-navigation` | `metadata/schema.md` | `candidate` | 7 | [`metadata-navigation.yaml`](../graphs/metadata-navigation.yaml) |
| `graph.source-boundary` | `governance/lifecycle/README.md` | `candidate` | 6 | [`source-boundary.yaml`](../graphs/source-boundary.yaml) |

## Refresh Decisions

| Decision value | Meaning |
| --- | --- |
| `refresh_now` | 由 `refresh-policy.yaml` 定義的 generated surface decision。 |
| `revalidate_only` | 由 `refresh-policy.yaml` 定義的 generated surface decision。 |
| `downgrade_confidence` | 由 `refresh-policy.yaml` 定義的 generated surface decision。 |
| `no_update_needed` | 由 `refresh-policy.yaml` 定義的 generated surface decision。 |

## Validation

- 產生前應先執行 `ruby scripts/validate-knowledge-runtime.rb`。
- 產生後應執行 Markdown link check、lints、close-loop dry run、commit / push / readback。
- 本報告是 generated view，不取代 `routing-registry.yaml`、`refresh-policy.yaml`、summary 或 graph source files。
