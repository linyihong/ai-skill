# Model Context Report

本檔由 `ai-skill runtime refresh` 產生，依 `knowledge/runtime/routing-registry.yaml` 的 model 欄位整理 model-aware context loading view。

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

### `medium`

| Route | Primary source | Compression | Reason |
| --- | --- | --- | --- |
| `route.governance.cognitive-state-evidence` | `governance/ai-runtime-governance/cognitive-state-governance.md` | `summary-first` | 先讀治理主文與 evidence hierarchy；只有修改 runtime guard、metadata evidence 或 validation scenario 時才展開相關 source。 |

### `small`

| Route | Primary source | Compression | Reason |
| --- | --- | --- | --- |
| `route.bootstrap.ai-skill` | `CORE_BOOTSTRAP.md` | `summary-first` | Bootstrap 階段只需要 Core Bootstrap 3 rules 與 OS layout，不需要完整 source。 |
| `route.runtime.phase-machine` | `runtime/phases/phase-machine.yaml` | `source-backed` | Phase Machine 是執行階段的核心狀態機，每個 session 啟動後必須載入。 |
| `route.runtime.obligation-ledger` | `runtime/obligations/obligation-ledger.yaml` | `source-backed` | Obligation Ledger 是執行階段的義務檢查清單，每個 phase entry 時必須載入。 |
| `route.runtime.blocking-gates` | `runtime/gates/blocking-gates.yaml` | `source-backed` | Blocking Gates 是執行階段的阻斷檢查，每個 phase transition 前必須載入。 |
| `route.runtime.recovery` | `runtime/recovery/recovery-strategies.yaml` | `source-backed` | State Repair System 在 blocking gate 阻斷時才需要載入，不佔用 bootstrap 預算。 |
| `route.runtime.scheduler` | `runtime/scheduler/priority-scheduler.yaml` | `source-backed` | Scheduler 在 phase entry 時載入，決定 execution queue 的優先順序。 |
| `route.runtime.transactions` | `runtime/transactions/transaction-machine.yaml` | `source-backed` | Transaction Runtime 在需要寫入 Ai-skill 檔案時才載入，不佔用 bootstrap 預算。 |
| `route.skill.discovery` | `knowledge/runtime/routing-registry.yaml` | `index-only` | Routing registry 輕量（~300 tokens），可在整個對話中 cache。 |
| `route.runtime.activation-rules` | `runtime/router/activation-rules.yaml` | `index-only` | Activation rules 輕量（~500 tokens），可在整個對話中 cache。 |
| `route.runtime.context-ttl` | `governance/ai-runtime-governance/context-attention-governance.md` | `summary-first` | 多數 context drift 先讀治理 gate 與 source intelligence；修改 runtime policy 時再展開 ttl-policy.yaml。 |
| `route.runtime.prompt-cache-alignment` | `runtime/context/prompt-cache-playbook.md` | `summary-first` | 先讀 prompt-cache summary 與 playbook；修改 metadata 或 enforcement 時再讀 full source。 |
| `route.governance.ai-runtime-five-step` | `governance/ai-runtime-governance/five-step-ai-governance.md` | `summary-first` | 只有 governance / architecture / automation 擴張決策需要讀取完整 five-step governance；一般任務不常駐載入。 |
| `route.runtime.context-loading` | `governance/lifecycle/routing-philosophy.md` | `summary-first` | routing decision 可先用 registry、index、summary；修改 source 時再升級。 |
| `route.governance.routing-signal` | `governance/ai-runtime-governance/routing-signal-governance.md` | `summary-first` | 多數路由疑義先讀治理 gate 與 source intelligence；只有修改 workflow/routing surface 時再展開 registry 和 activation table。 |
| `route.governance.validation-scenario` | `governance/ai-runtime-governance/validation-scenario-governance.md` | `summary-first` | 先讀 scenario governance 與 source intelligence；建立或修改 scenario 時再展開 validation format 與 failure-learning source。 |
| `route.runtime.router-flow` | `runtime/router/` | `index-only` | Router flow 是設計文件，需要 routing 決策時才讀。 |
| `route.intelligence.engineering.heuristics` | `intelligence/engineering/heuristics/README.md` | `index-only` | 通用 heuristics 在需要對應判斷時才讀取完整 atom。 |
| `route.runtime.context-ttl-doc` | `governance/lifecycle/context-ttl-philosophy.md` | `index-only` | TTL doc 是設計文件，需要 prune context 時才讀。 |
| `route.workflow.documentation-ai-native` | `workflow/documentation/README.md` | `summary-first` | 多數任務只需 workflow README 與 execution-flow；policy 全文留在 enforcement。 |
| `route.intelligence.requirements-cognition` | `intelligence/engineering/requirements/README.md` | `summary-first` | 多數需求任務先讀 requirements README 與 summary；高 ambiguity 時再展開 behavior / specification / validation docs。 |
| `route.intelligence.architectural-fit` | `intelligence/engineering/architecture/architectural-fit/README.md` | `summary-first` | 先讀 architectural-fit README 與 metadata，只有在 high-complexity 或 DDD 訊號強時展開完整 DDD tactical docs。 |
| `route.intelligence.engineering.agent-architecture` | `intelligence/engineering/agent-architecture/README.md` | `index-only` | agent-architecture atoms 在需要理解 AI 行為模式時才讀取完整內容。 |
| `route.feedback.history` | `feedback/history/README.md` | `index-only` | feedback/history 在需要查詢或新增 lesson 時才讀取。 |
| `route.runtime.decision-recording` | `runtime/decisions/decision-recording.yaml` | `source-backed` | 避免決策只留在 plan 或對話；查問題時依 lookup.by_symptom 回溯。 |
| `route.decisions.adr` | `decisions/README.md` | `summary-first` | decisions/ 在需要查詢歷史決策或建立新 ADR 時才讀取。 |
| `route.architecture.permanent-docs` | `architecture/README.md` | `summary-first` | architecture/ 在需要查詢架構定義或設計原則時才讀取。 |
| `route.evaluations.scenario-results` | `evaluations/README.md` | `index-only` | evaluations/ 在需要查詢或記錄 scenario 執行結果時才讀取。 |
| `route.tools.metadata-routing` | `tools/README.md` | `index-only` | tools/ 在需要查詢 tool 設定或 routing 規則時才讀取。 |
| `route.traces.decision-traces` | `traces/README.md` | `index-only` | traces/ 在需要查詢或記錄 decision trace 時才讀取。 |
| `route.anti-patterns.runtime-patterns` | `anti-patterns/README.md` | `summary-first` | anti-patterns/ 在需要查詢 runtime 失效模式或預防方式時才讀取。 |
| `route.runtime.compiler` | `runtime/compiler/compiler-rules.yaml` | `source-backed` | Compiler 在 workflow/enforcement 文件變更後才需要執行，不佔用 bootstrap 預算。 |
| `route.runtime.intelligence-routing` | `runtime/intelligence/intelligence-routing.yaml` | `index-only` | Intelligence routing 在需要領域知識時才載入，不佔用 bootstrap 預算。 |
| `route.runtime.output-governance` | `runtime/output-governance/` | `source-backed` | Output Governance 在 validation 與 finalize phase 時才需要載入，不佔用 bootstrap 預算。 |
| `route.runtime.distributed` | `runtime/distributed/` | `source-backed` | Distributed Runtime 在 multi-agent 協作或 async job 執行時才需要載入，不佔用 bootstrap 預算。 |
| `route.governance.knowledge-update-flow` | `governance/lifecycle/knowledge-update-flow.md` | `source-backed` | Knowledge Update Flow 在 checkpoint 或 finalize phase 需要執行知識更新時才載入，不佔用 bootstrap 預算。 |

### `specialized`

| Route | Primary source | Compression | Reason |
| --- | --- | --- | --- |
| `route.workflow.apk-analysis` | `workflow/apk-analysis/execution-flow.md` | `source-backed` | APK analysis 需要 workflow、analysis methods 與 domain-specific intelligence routing。 |
| `route.intelligence.apk-highest-leverage-path` | `intelligence/engineering/analytical-reasoning/highest-leverage-analysis-path.md` | `source-backed` | APK route selection 需要 domain workflow、feedback source 與 intelligence judgment 一起判斷。 |
| `route.workflow.software-delivery` | `workflow/software-delivery/execution-flow.md` | `source-backed` | App development guidance 需要 workflow、analysis methods 與 domain-specific controls/checklists routing。triggers 已擴充以支援中途切換。 |
| `route.workflow.greenfield` | `workflow/greenfield/execution-flow.md` | `source-backed` | Greenfield 需要 workflow、templates 與 software-delivery 的 BDD Closure 整合。 |
| `route.workflow.travel-planning` | `workflow/travel-planning/execution-flow.md` | `source-backed` | Travel planning 需要 workflow、analysis methods 與 domain-specific intelligence routing。 |
| `route.runtime.onboarding` | `runtime/onboarding/README.md` | `summary-first` | Onboarding 文件在需要執行對應 workflow 的完整流程時才讀取。 |
| `route.analysis.apk.workflows` | `analysis/apk/workflows/README.md` | `summary-first` | 需要執行對應操作流程時才讀取完整 workflow。 |
| `route.analysis.web` | `analysis/web/README.md` | `summary-first` | Web Scraping analysis 在需要從網頁提取資料時才讀取，不佔用 bootstrap 預算。 |
| `route.intelligence.apk-analysis.atoms` | `intelligence/engineering/apk-analysis/README.md` | `summary-first` | 需要對應決策智慧時才讀取完整 atom。 |
| `route.validation.ai-decision-contract` | `validation/README.md` | `summary-first` | 需要驗證 AI 決策品質時才讀取完整 scenario 與 rule 定義。 |

## Compression View

| Compression level | Routes | Escalation note |
| --- | --- | --- |
| `index-only` | `route.skill.discovery`, `route.runtime.activation-rules`, `route.runtime.router-flow`, `route.intelligence.engineering.heuristics`, `route.runtime.context-ttl-doc`, `route.intelligence.engineering.agent-architecture`, `route.feedback.history`, `route.evaluations.scenario-results`, `route.tools.metadata-routing`, `route.traces.decision-traces`, `route.runtime.intelligence-routing` | 依 `models/compression/README.md` 的 escalation rules 判斷。 |
| `source-backed` | `route.runtime.phase-machine`, `route.runtime.obligation-ledger`, `route.runtime.blocking-gates`, `route.runtime.recovery`, `route.runtime.scheduler`, `route.runtime.transactions`, `route.governance.durable-goal-boundary`, `route.metadata.knowledge-atom-schema`, `route.workflow.apk-analysis`, `route.intelligence.apk-highest-leverage-path`, `route.feedback.promotion-pipeline`, `route.models.model-aware-routing`, `route.workflow.software-delivery`, `route.workflow.greenfield`, `route.workflow.travel-planning`, `route.runtime.decision-recording`, `route.runtime.compiler`, `route.runtime.output-governance`, `route.runtime.distributed`, `route.governance.knowledge-update-flow` | 需要 primary source 與 required dependencies；適合 writeback、migration 或 domain work。 |
| `summary-first` | `route.bootstrap.ai-skill`, `route.runtime.context-ttl`, `route.runtime.prompt-cache-alignment`, `route.governance.ai-runtime-five-step`, `route.runtime.context-loading`, `route.governance.routing-signal`, `route.governance.validation-scenario`, `route.governance.cognitive-state-evidence`, `route.workflow.documentation-ai-native`, `route.runtime.onboarding`, `route.analysis.apk.workflows`, `route.analysis.web`, `route.intelligence.apk-analysis.atoms`, `route.intelligence.requirements-cognition`, `route.intelligence.architectural-fit`, `route.validation.ai-decision-contract`, `route.decisions.adr`, `route.architecture.permanent-docs`, `route.anti-patterns.runtime-patterns` | 適合先用 registry / summary 判斷 relevance；修改 source 時升級。 |

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

- 產生前應先確認 `routing-registry.yaml` 可通過 `ai-skill runtime validate`。
- 產生後應重新執行 `ai-skill runtime validate`，檢查本 report links。
- 本報告是 generated view，不取代 `models/profiles/README.md`、`models/compression/README.md` 或 routing registry。
