# Knowledge Runtime Report

本檔由 `ruby scripts/generate-knowledge-runtime-report.rb --write` 產生，彙整 runtime registry、summaries、graphs 與 refresh policy 的目前狀態。

## Source Surfaces

| Surface | Path | Count / Status |
| --- | --- | --- |
| Routing registry | [`routing-registry.yaml`](routing-registry.yaml) | 41 records |
| Refresh policy | [`refresh-policy.yaml`](refresh-policy.yaml) | candidate |
| Model context report | [`model-context-report.md`](model-context-report.md) | generated view |
| Model checklists | [`model-checklists.md`](model-checklists.md) | generated view |
| SQLite runtime index | [`sqlite/`](sqlite/) | generated lookup cache prototype |
| Summaries | [`../summaries/`](../summaries/) | 14 files |
| Graph records | [`../graphs/`](../graphs/) | 29 files |

## Routing Records

| ID | Primary source | Model | Compression | Validation signal |
| --- | --- | --- | --- | --- |
| `route.bootstrap.ai-skill` | `CORE_BOOTSTRAP.md` | `small` | `summary-first` | Core Bootstrap 3 rules 已讀，git status 已檢查。 |
| `route.runtime.phase-machine` | `runtime/compiler/embedded_data.rb` | `small` | `source-backed` | 目前 phase 已確認，allowed_actions 與 forbidden_actions 已載入，phase transition rules 已檢查。 |
| `route.runtime.obligation-ledger` | `runtime/compiler/embedded_data.rb` | `small` | `source-backed` | 本 phase 的 obligations 已確認，pending obligations 已記錄，blocking gates 已檢查。 |
| `route.runtime.blocking-gates` | `runtime/compiler/embedded_data.rb` | `small` | `source-backed` | 本 phase 的 blocking gates 已檢查，無未通過的 critical/high gates，phase transition 可進行。 |
| `route.runtime.recovery` | `runtime/compiler/embedded_data.rb` | `small` | `source-backed` | 阻斷 gate 已對應到 strategy，repair steps 已執行，verify steps 已通過。 |
| `route.runtime.scheduler` | `runtime/compiler/embedded_data.rb` | `small` | `source-backed` | Queue 已依 priority 排序，blocking gates 優先於 obligations，dependencies 已解析。 |
| `route.runtime.transactions` | `runtime/compiler/embedded_data.rb` | `small` | `source-backed` | Transaction state 正確，所有 rules 已檢查，templates 已套用。 |
| `route.skill.discovery` | `knowledge/runtime/routing-registry.yaml` | `small` | `index-only` | Task intent 已對應到 routing-registry.yaml 的 triggers，entrypoint 與 summary path 可解析。 |
| `route.runtime.activation-rules` | `runtime/compiler/embedded_data.rb` | `small` | `index-only` | 目前 task 已比對 activation-rules.yaml，符合條件的 rules 已載入，不符合的已 deferred。 |
| `route.runtime.context-ttl` | `runtime/compiler/embedded_data.rb` | `small` | `index-only` | Context TTL 已套用，過期 context 已 prune，token 使用率在預算內。 |
| `route.governance.durable-goal-boundary` | `enforcement/conversation-goal-ledger.md` | `large` | `source-backed` | 長期狀態已落到 durable planning 文件，且 active goal 完成後才刪除。 |
| `route.runtime.context-loading` | `runtime/routing/` | `small` | `summary-first` | Primary source、deferred sources、source-of-truth gate 與 validation signal 已可被記錄。 |
| `route.metadata.knowledge-atom-schema` | `metadata/schema.md` | `large` | `source-backed` | 欄位可套用到第一批 atom candidates，且 Markdown links 可解析。 |
| `route.workflow.apk-analysis` | `workflow/apk-analysis/execution-flow.md` | `specialized` | `source-backed` | 新分層路徑可讀取，workflow 與 analysis 內容已分離。 |
| `route.intelligence.apk-highest-leverage-path` | `intelligence/engineering/analytical-reasoning/highest-leverage-analysis-path.md` | `specialized` | `source-backed` | 已比較可用 routes、選定 evidence-to-cost ratio 最高路線、記錄 fallback 與 attribution 回補條件。 |
| `route.feedback.promotion-pipeline` | `feedback/promotion/README.md` | `large` | `source-backed` | 原 lesson source 保留、promotion target 明確、runtime surfaces 與 close-loop validation 已同步。 |
| `route.models.model-aware-routing` | `models/profiles/README.md` | `large` | `source-backed` | Profile、compression level、primary source、deferred sources 與 validation signal 可被記錄。 |
| `route.runtime.router-flow` | `runtime/router/` | `small` | `index-only` | Routing flow 已理解，activation rules 與 TTL policy 已對應到對應階段。 |
| `route.intelligence.engineering.heuristics` | `intelligence/engineering/heuristics/README.md` | `small` | `index-only` | 各 atom 有明確原則與決策表，可反查驗證。 |
| `route.runtime.context-ttl-doc` | `runtime/context/` | `small` | `index-only` | TTL policy 已套用，prune strategy 已理解。 |
| `route.workflow.software-delivery` | `workflow/software-delivery/execution-flow.md` | `specialized` | `source-backed` | 新分層路徑可讀取，workflow 與 analysis 內容已分離。 |
| `route.workflow.travel-planning` | `workflow/travel-planning/execution-flow.md` | `specialized` | `source-backed` | 新分層路徑可讀取，workflow 與 analysis 內容已分離。 |
| `route.workflow.documentation-ai-native` | `workflow/documentation/README.md` | `small` | `summary-first` | 新文件或目錄具 index-first 導航、分類維度（kind/audience/stability）已標註； 長文已按 document-sizing 拆分；語言與工具敘述依專案自訂 policy（本 route 不預設 tool-neutral）。  |
| `route.runtime.onboarding` | `runtime/onboarding/README.md` | `specialized` | `summary-first` | 各 quickstart 的步驟可依序執行，且與對應 workflow 的內容一致。 |
| `route.analysis.apk.workflows` | `analysis/apk/workflows/README.md` | `specialized` | `summary-first` | 各 workflow 有明確步驟與產出格式，可依序執行。 |
| `route.analysis.web` | `analysis/web/README.md` | `specialized` | `summary-first` | 目標網站已評估（技術棧、JS 需求、anti-bot 保護），工具已選擇（HTTP/Dynamic/Stealth）， 提取策略已設計（selector/adaptive parsing），風險已評估（legal/technical/data quality）。  |
| `route.intelligence.apk-analysis.atoms` | `intelligence/engineering/analytical-reasoning/README.md` | `specialized` | `summary-first` | 各 atom 有明確決策表或信號表，可反查驗證。 |
| `route.validation.ai-decision-contract` | `validation/README.md` | `specialized` | `summary-first` | Scenario 的 expected_route 與 actual trace 可比對；forbidden_routes 未被使用。 |
| `route.intelligence.engineering.agent-architecture` | `intelligence/engineering/agent-architecture/README.md` | `small` | `index-only` | 各 atom 有明確原則、症狀表與預防方式，可反查驗證。 |
| `route.feedback.history` | `feedback/history/README.md` | `small` | `index-only` | Lesson 已寫入 feedback/history/<domain>/ 對應分類，且 feedback/history/<domain>/README.md 已更新索引。 |
| `route.decisions.adr` | `decisions/README.md` | `small` | `summary-first` | decisions/README.md 已讀取，ADR 清單已查詢，graph edge 已確認。 |
| `route.architecture.permanent-docs` | `architecture/README.md` | `small` | `summary-first` | architecture/README.md 已讀取，架構文件清單已查詢。 |
| `route.evaluations.scenario-results` | `evaluations/README.md` | `small` | `index-only` | Evaluation 記錄已依 scenario 分類存放，格式與 template 一致。 |
| `route.tools.metadata-routing` | `tools/README.md` | `small` | `index-only` | Tool metadata 已定義，compression 策略與 routing 規則已建立。 |
| `route.traces.decision-traces` | `traces/README.md` | `small` | `index-only` | Trace 記錄已依 scenario 分類存放，格式與 template 一致。 |
| `route.anti-patterns.runtime-patterns` | `anti-patterns/README.md` | `small` | `summary-first` | Anti-pattern 已依格式記錄，症狀、預防與恢復方式已定義。 |
| `route.runtime.compiler` | `runtime/compiler/embedded_data.rb` | `small` | `source-backed` | 所有 modified sources 已編譯，runtime.db 的 generated_surfaces 表包含最新記錄，validate-runtime-db.rb 回傳 exit 0。 |
| `route.runtime.intelligence-routing` | `runtime/compiler/embedded_data.rb` | `small` | `index-only` | Task intent 已對應到 intelligence-routing.yaml 的 domain，applicable_phases 已檢查，domain README 已載入。 |
| `route.runtime.output-governance` | `runtime/output-governance/` | `small` | `source-backed` | Language consistency 已確認，sanitization 已通過，tool neutrality 已檢查， format compliance 已驗證，governance gates 全部通過。  |
| `route.runtime.distributed` | `runtime/distributed/` | `small` | `source-backed` | Distributed locks 已正確 acquire/release，multi-agent coordination rules 已遵守， async job lifecycle 狀態轉換正確，無 deadlock 或 stale state。  |
| `route.governance.knowledge-update-flow` | `governance/lifecycle/knowledge-update-flow.md` | `small` | `source-backed` | 目前 knowledge update 的步驟已確認，entry conditions 已滿足，exit conditions 已檢查，reference sources 已載入，下一步驟已決定。  |

## Summary Records

| Atom ID | Lifecycle | File | Summary |
| --- | --- | --- | --- |
| `architecture.apk-analysis-pilot` | `candidate` | [`apk-analysis-pilot.md`](../summaries/apk-analysis-pilot.md) | `apk-analysis` 作為第一個 Workflow / Analysis / Intelligence 分離 pilot 的 migration map。它建立新 reference-first 候選目的地，舊 `skills/apk-analysis/` 已不再作為 active entrypoint。 |
| `intelligence.apk-highest-leverage-analysis` | `candidate` | [`apk-highest-leverage-analysis.md`](../summaries/apk-highest-leverage-analysis.md) | APK 分析 checkpoint 應先界定未知，再依 time-to-evidence、語意距離、安全性與 validation clarity 選擇最高收益路線。 |
| `architecture.context-cost-optimization` | `validated` | [`context-cost-optimization.md`](../summaries/context-cost-optimization.md) | Token 成本優化規劃。Phase 1（立即省錢）：Bootstrap 極小化（~800 tokens）、README 拆分、Rule lazy-load、Summary layer。Phase 2（架構升級）：Runtime Context Router、Context Cost Metadata、Skill Index、Context TTL。Phase 3（長期）：Semantic Retrieval、Episodic Memory、Multi-model Routing。 |
| `workflow.software-delivery` | `validated` | [`development-guidance.md`](../summaries/development-guidance.md) | 將授權 App/API/Embedded/Firmware 觀察轉成開發 guidance、實作模式、控制項、檢查清單。涵蓋 mobile（Android/iOS/Flutter/React Native）、backend API、embedded firmware 的安全控制、實作模式與 release gate。原 `skills/app-development-guidance/` 已刪除，所有內容已遷移至新分層。 |
| `feedback.promotion.pipeline` | `candidate` | [`feedback-promotion-pipeline.md`](../summaries/feedback-promotion-pipeline.md) | 定義 feedback lesson 從 skill-local history 推進到 workflow、intelligence、enforcement、memory 或 runtime surfaces 的 promotion / downgrade gate。 |
| `governance.goal-ledger-boundary` | `validated` | [`goal-ledger-boundary.md`](../summaries/goal-ledger-boundary.md) | `.agent-goals/` 只保存 active conversation goals；長期 roadmap、phase、migration、promotion、deprecation 與治理狀態必須落到 durable planning 文件。 |
| `knowledge.navigation` | `validated` | [`knowledge-navigation.md`](../summaries/knowledge-navigation.md) | 知識導航系統：indexes（任務路由）、summaries（300-500 token 摘要）、graphs（知識圖譜邊）、runtime（routing registry、refresh policy、SQLite lookup cache）。讓 agent 用最小 token 成本找到正確知識。 |
| `memory.operations` | `candidate` | [`memory-operations.md`](../summaries/memory-operations.md) | 長期記憶層：short-term（目前 task context）、episodic（過去 task 關鍵決策與結果）、project（專案歷史脈絡）、failure（反覆失效模式）。支援 similarity-based retrieval。 |
| `metadata.schema.knowledge-atom` | `validated` | [`metadata-schema.md`](../summaries/metadata-schema.md) | Knowledge Atom metadata schema v1，定義 atom 的必填欄位、選填欄位、受控值、YAML 範本與驗證規則。 |
| `models.routing` | `candidate` | [`model-routing.md`](../summaries/model-routing.md) | 多模型協作架構：capability profile（small/large/specialized）、compression strategy（checklist/compressed/full）、model-aware context report。根據 task 複雜度選擇模型與 context 策略。 |
| `governance.repo-maintenance` | `validated` | [`repo-governance.md`](../summaries/repo-governance.md) | AI-native Knowledge Operating System 本身的維護、升級、遷移與治理。涵蓋 lifecycle management、validation、cleanup、splitting rules、dependency maintenance。 |
| `root.bootstrap.ai-skill` | `validated` | [`root-bootstrap.md`](../summaries/root-bootstrap.md) | Ai-skill 工作的 bootstrap 入口。Root README 定義 OS layout 與 cost-aware 啟動流程；CORE_BOOTSTRAP.md 定義 3 條核心規則（~800 tokens）；enforcement README 定義 Runtime Activation Model 與 lazy-load rules。 |
| `runtime.operations` | `validated` | [`runtime-operations.md`](../summaries/runtime-operations.md) | Runtime 層負責 context routing、dynamic loading、context pruning、agent coordination 與 orchestration。包含 router（activation rules、cost budget）、context（TTL policy、prune strategy）。 |
| `workflow.travel-planning` | `candidate` | [`travel-planning.md`](../summaries/travel-planning.md) | 依目的地、日期、交通與玩法規劃行程，包含營業時間查證、交通比較、住宿與備案。支援 itinerary 結構化輸出與可行性檢查。 |

## Graph Records

| ID | Source | Status | Edges | File |
| --- | --- | --- | --- | --- |
| `graph.analysis-layers` | `analysis/README.md` | `candidate` | 22 | [`analysis-layers.yaml`](../graphs/analysis-layers.yaml) |
| `graph.analysis-repo-methods` | `analysis/repo/README.md` | `candidate` | 6 | [`analysis-repo-methods.yaml`](../graphs/analysis-repo-methods.yaml) |
| `graph.apk-analysis-pilot` | `plans/archived/2026-05-11-1129-apk-analysis-pilot-migration.md` | `candidate` | 5 | [`apk-analysis-pilot.yaml`](../graphs/apk-analysis-pilot.yaml) |
| `graph.apk-highest-leverage-analysis` | `intelligence/engineering/analytical-reasoning/highest-leverage-analysis-path.md` | `candidate` | 4 | [`apk-highest-leverage-analysis.yaml`](../graphs/apk-highest-leverage-analysis.yaml) |
| `graph.decisions-adr` | `decisions/README.md` | `candidate` | 9 | [`decisions-adr.yaml`](../graphs/decisions-adr.yaml) |
| `graph.feedback-layers` | `feedback/README.md` | `candidate` | 9 | [`feedback-layers.yaml`](../graphs/feedback-layers.yaml) |
| `graph.feedback-promotion-pipeline` | `feedback/promotion/README.md` | `candidate` | 6 | [`feedback-promotion-pipeline.yaml`](../graphs/feedback-promotion-pipeline.yaml) |
| `graph.governance-layers` | `governance/README.md` | `candidate` | 22 | [`governance-layers.yaml`](../graphs/governance-layers.yaml) |
| `graph.intelligence-agent-architecture` | `intelligence/engineering/agent-architecture/README.md` | `candidate` | 20 | [`intelligence-agent-architecture.yaml`](../graphs/intelligence-agent-architecture.yaml) |
| `graph.intelligence-anti-patterns` | `intelligence/engineering/anti-patterns/generic-repository-overuse.md` | `candidate` | 5 | [`intelligence-anti-patterns.yaml`](../graphs/intelligence-anti-patterns.yaml) |
| `graph.intelligence-apk-analysis-atoms` | `intelligence/engineering/analytical-reasoning/README.md` | `candidate` | 17 | [`intelligence-apk-analysis-atoms.yaml`](../graphs/intelligence-apk-analysis-atoms.yaml) |
| `graph.intelligence-architecture` | `intelligence/engineering/architecture/modular-monolith-vs-microservices.md` | `candidate` | 5 | [`intelligence-architecture.yaml`](../graphs/intelligence-architecture.yaml) |
| `graph.intelligence-business` | `intelligence/business/saas-pricing-heuristics.md` | `candidate` | 3 | [`intelligence-business.yaml`](../graphs/intelligence-business.yaml) |
| `graph.intelligence-app-development-guidance` | `intelligence/engineering/development/README.md` | `candidate` | 8 | [`intelligence-development-guidance.yaml`](../graphs/intelligence-development-guidance.yaml) |
| `graph.intelligence-distributed-systems` | `intelligence/engineering/distributed-systems/eventual-consistency-patterns.md` | `candidate` | 5 | [`intelligence-distributed-systems.yaml`](../graphs/intelligence-distributed-systems.yaml) |
| `graph.intelligence-domain` | `intelligence/engineering/domain/aggregate-boundary-heuristics.md` | `candidate` | 5 | [`intelligence-domain.yaml`](../graphs/intelligence-domain.yaml) |
| `graph.intelligence-failure` | `intelligence/engineering/failure/connection-leak-patterns.md` | `candidate` | 5 | [`intelligence-failure.yaml`](../graphs/intelligence-failure.yaml) |
| `graph.intelligence-heuristics` | `intelligence/engineering/heuristics/README.md` | `candidate` | 10 | [`intelligence-heuristics.yaml`](../graphs/intelligence-heuristics.yaml) |
| `graph.intelligence-repo-analysis` | `intelligence/engineering/analytical-reasoning/README.md` | `candidate` | 8 | [`intelligence-repo-analysis.yaml`](../graphs/intelligence-repo-analysis.yaml) |
| `graph.intelligence-tradeoffs` | `intelligence/engineering/tradeoffs/postgres-vs-mongodb.md` | `candidate` | 5 | [`intelligence-tradeoffs.yaml`](../graphs/intelligence-tradeoffs.yaml) |
| `graph.intelligence-travel` | `intelligence/travel/README.md` | `candidate` | 7 | [`intelligence-travel.yaml`](../graphs/intelligence-travel.yaml) |
| `graph.metadata-navigation` | `metadata/schema.md` | `candidate` | 7 | [`metadata-navigation.yaml`](../graphs/metadata-navigation.yaml) |
| `graph.runtime-onboarding` | `runtime/onboarding/README.md` | `candidate` | 4 | [`runtime-onboarding.yaml`](../graphs/runtime-onboarding.yaml) |
| `graph.runtime-pipeline` | `runtime/pipeline/` | `candidate` | 12 | [`runtime-pipeline.yaml`](../graphs/runtime-pipeline.yaml) |
| `graph.runtime-prompt-artifacts` | `runtime/prompt-artifacts/` | `candidate` | 8 | [`runtime-prompt-artifacts.yaml`](../graphs/runtime-prompt-artifacts.yaml) |
| `graph.source-boundary` | `governance/lifecycle/README.md` | `candidate` | 6 | [`source-boundary.yaml`](../graphs/source-boundary.yaml) |
| `graph.workflow-layers` | `workflow/README.md` | `candidate` | 19 | [`workflow-layers.yaml`](../graphs/workflow-layers.yaml) |
| `graph.workflow-app-development-guidance` | `workflow/software-delivery/README.md` | `candidate` | 9 | [`workflow-software-delivery.yaml`](../graphs/workflow-software-delivery.yaml) |
| `graph.workflow-travel-planning` | `workflow/travel-planning/README.md` | `candidate` | 6 | [`workflow-travel-planning.yaml`](../graphs/workflow-travel-planning.yaml) |

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
