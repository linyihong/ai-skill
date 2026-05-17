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
| `route.governance.durable-goal-boundary` | `enforcement/conversation-goal-ledger.md` | `source-backed` | Goal 刪除與 durable planning gate 會影響 long-term state，需讀 shared rule source。 |
| `route.metadata.knowledge-atom-schema` | `metadata/schema.md` | `source-backed` | 建立或修改 atom metadata 時需讀 schema 與子規則全文。 |
| `route.feedback.promotion-pipeline` | `feedback/promotion/README.md` | `source-backed` | Promotion / downgrade 會影響 enforcement rules、skill history、intelligence 與 runtime surfaces，需要讀 full source 與 validation gates。 |
| `route.models.model-aware-routing` | `models/profiles/README.md` | `source-backed` | 修改 model strategy 會影響 runtime routing，需讀 profiles、compression 與 routing source。 |

### `small`

| Route | Primary source | Compression | Reason |
| --- | --- | --- | --- |
| `route.bootstrap.ai-skill` | `CORE_BOOTSTRAP.md` | `summary-first` | Bootstrap 階段只需要 Core Bootstrap 3 rules 與 OS layout，不需要完整 source。 |
| `route.skill.discovery` | `skills-index.yaml` | `index-only` | Skills index 輕量（~300 tokens），可在整個對話中 cache。 |
| `route.runtime.activation-rules` | `runtime/compiler/embedded_data.rb` | `index-only` | Activation rules 輕量（~500 tokens），可在整個對話中 cache。 |
| `route.runtime.context-ttl` | `runtime/compiler/embedded_data.rb` | `index-only` | TTL policy 輕量，需要 prune 時才讀。 |
| `route.runtime.context-loading` | `runtime/routing/` | `summary-first` | routing decision 可先用 registry、index、summary；修改 source 時再升級。 |
| `route.runtime.router-flow` | `runtime/router/` | `index-only` | Router flow 是設計文件，需要 routing 決策時才讀。 |
| `route.intelligence.engineering.heuristics` | `intelligence/engineering/heuristics/README.md` | `index-only` | 通用 heuristics 在需要對應判斷時才讀取完整 atom。 |
| `route.runtime.context-ttl-doc` | `runtime/context/` | `index-only` | TTL doc 是設計文件，需要 prune context 時才讀。 |
| `route.workflow.documentation-ai-native` | `workflow/documentation/README.md` | `summary-first` | 多數任務只需 workflow README 與 execution-flow；policy 全文留在 enforcement。 |
| `route.intelligence.engineering.agent-architecture` | `intelligence/engineering/agent-architecture/README.md` | `index-only` | agent-architecture atoms 在需要理解 AI 行為模式時才讀取完整內容。 |
| `route.feedback.history` | `feedback/history/README.md` | `index-only` | feedback/history 在需要查詢或新增 lesson 時才讀取。 |
| `route.decisions.adr` | `decisions/README.md` | `summary-first` | decisions/ 在需要查詢歷史決策或建立新 ADR 時才讀取。 |
| `route.architecture.permanent-docs` | `architecture/README.md` | `summary-first` | architecture/ 在需要查詢架構定義或設計原則時才讀取。 |
| `route.evaluations.scenario-results` | `evaluations/README.md` | `index-only` | evaluations/ 在需要查詢或記錄 scenario 執行結果時才讀取。 |
| `route.tools.metadata-routing` | `tools/README.md` | `index-only` | tools/ 在需要查詢 tool 設定或 routing 規則時才讀取。 |
| `route.traces.decision-traces` | `traces/README.md` | `index-only` | traces/ 在需要查詢或記錄 decision trace 時才讀取。 |
| `route.anti-patterns.runtime-patterns` | `anti-patterns/README.md` | `summary-first` | anti-patterns/ 在需要查詢 runtime 失效模式或預防方式時才讀取。 |

### `specialized`

| Route | Primary source | Compression | Reason |
| --- | --- | --- | --- |
| `route.workflow.apk-analysis` | `workflow/apk-analysis/execution-flow.md` | `source-backed` | APK analysis 需要 workflow、analysis methods 與 domain-specific intelligence routing。 |
| `route.intelligence.apk-highest-leverage-path` | `intelligence/engineering/analytical-reasoning/highest-leverage-analysis-path.md` | `source-backed` | APK route selection 需要 domain workflow、feedback source 與 intelligence judgment 一起判斷。 |
| `route.workflow.software-delivery` | `workflow/software-delivery/execution-flow.md` | `source-backed` | App development guidance 需要 workflow、analysis methods 與 domain-specific controls/checklists routing。 |
| `route.workflow.travel-planning` | `workflow/travel-planning/execution-flow.md` | `source-backed` | Travel planning 需要 workflow、analysis methods 與 domain-specific intelligence routing。 |
| `route.runtime.onboarding` | `runtime/onboarding/README.md` | `summary-first` | Onboarding 文件在需要執行對應 workflow 的完整流程時才讀取。 |
| `route.analysis.apk.workflows` | `analysis/apk/workflows/README.md` | `summary-first` | 需要執行對應操作流程時才讀取完整 workflow。 |
| `route.intelligence.apk-analysis.atoms` | `intelligence/engineering/analytical-reasoning/README.md` | `summary-first` | 需要對應決策智慧時才讀取完整 atom。 |
| `route.validation.ai-decision-contract` | `validation/README.md` | `summary-first` | 需要驗證 AI 決策品質時才讀取完整 scenario 與 rule 定義。 |

## Compression View

| Compression level | Routes | Escalation note |
| --- | --- | --- |
| `index-only` | `route.skill.discovery`, `route.runtime.activation-rules`, `route.runtime.context-ttl`, `route.runtime.router-flow`, `route.intelligence.engineering.heuristics`, `route.runtime.context-ttl-doc`, `route.intelligence.engineering.agent-architecture`, `route.feedback.history`, `route.evaluations.scenario-results`, `route.tools.metadata-routing`, `route.traces.decision-traces` | 依 `models/compression/README.md` 的 escalation rules 判斷。 |
| `source-backed` | `route.governance.durable-goal-boundary`, `route.metadata.knowledge-atom-schema`, `route.workflow.apk-analysis`, `route.intelligence.apk-highest-leverage-path`, `route.feedback.promotion-pipeline`, `route.models.model-aware-routing`, `route.workflow.software-delivery`, `route.workflow.travel-planning` | 需要 primary source 與 required dependencies；適合 writeback、migration 或 domain work。 |
| `summary-first` | `route.bootstrap.ai-skill`, `route.runtime.context-loading`, `route.workflow.documentation-ai-native`, `route.runtime.onboarding`, `route.analysis.apk.workflows`, `route.intelligence.apk-analysis.atoms`, `route.validation.ai-decision-contract`, `route.decisions.adr`, `route.architecture.permanent-docs`, `route.anti-patterns.runtime-patterns` | 適合先用 registry / summary 判斷 relevance；修改 source 時升級。 |

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
