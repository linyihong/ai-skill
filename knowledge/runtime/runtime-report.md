# Knowledge Runtime Report

本檔由 `ruby scripts/generate-knowledge-runtime-report.rb --write` 產生，彙整 runtime registry、summaries、graphs 與 refresh policy 的目前狀態。

## Source Surfaces

| Surface | Path | Count / Status |
| --- | --- | --- |
| Routing registry | [`routing-registry.yaml`](routing-registry.yaml) | 22 records |
| Refresh policy | [`refresh-policy.yaml`](refresh-policy.yaml) | candidate |
| Model context report | [`model-context-report.md`](model-context-report.md) | generated view |
| Model checklists | [`model-checklists.md`](model-checklists.md) | generated view |
| SQLite runtime index | [`sqlite/`](sqlite/) | generated lookup cache prototype |
| Summaries | [`../summaries/`](../summaries/) | 14 files |
| Graph records | [`../graphs/`](../graphs/) | 29 files |

## Routing Records

| ID | Primary source | Model | Compression | Validation signal |
| --- | --- | --- | --- | --- |
| `route.bootstrap.ai-skill` | `CORE_BOOTSTRAP.md` | `small` | `summary-first` | Core Bootstrap 3 rules 已讀，skills-index.yaml 已查詢，git status 已檢查。 |
| `route.skill.discovery` | `skills-index.yaml` | `small` | `index-only` | Task intent 已對應到 skills-index.yaml 的 triggers，entrypoint 與 summary path 可解析。 |
| `route.runtime.activation-rules` | `runtime/router/activation-rules.yaml` | `small` | `index-only` | 目前 task 已比對 activation-rules.yaml，符合條件的 rules 已載入，不符合的已 deferred。 |
| `route.runtime.context-ttl` | `runtime/context/ttl-policy.yaml` | `small` | `index-only` | Context TTL 已套用，過期 context 已 prune，token 使用率在預算內。 |
| `route.governance.durable-goal-boundary` | `shared-rules/conversation-goal-ledger.md` | `large` | `source-backed` | 長期狀態已落到 durable planning 文件，且 active goal 完成後才刪除。 |
| `route.runtime.context-loading` | `runtime/routing/README.md` | `small` | `summary-first` | Primary source、deferred sources、source-of-truth gate 與 validation signal 已可被記錄。 |
| `route.metadata.knowledge-atom-schema` | `metadata/schema.md` | `large` | `source-backed` | 欄位可套用到第一批 atom candidates，且 Markdown links 可解析。 |
| `route.skill.apk-analysis` | `skills/apk-analysis/SKILL.md` | `specialized` | `source-backed` | 舊入口仍可讀，新 reference-first paths 可找到，且未 bulk migrate skill content。 |
| `route.intelligence.apk-highest-leverage-path` | `intelligence/engineering/analytical-reasoning/highest-leverage-analysis-path.md` | `specialized` | `source-backed` | 已比較可用 routes、選定 evidence-to-cost ratio 最高路線、記錄 fallback 與 attribution 回補條件。 |
| `route.feedback.promotion-pipeline` | `feedback/promotion/README.md` | `large` | `source-backed` | 原 lesson source 保留、promotion target 明確、runtime surfaces 與 close-loop validation 已同步。 |
| `route.models.model-aware-routing` | `models/profiles/README.md` | `large` | `source-backed` | Profile、compression level、primary source、deferred sources 與 validation signal 可被記錄。 |
| `route.runtime.router-flow` | `runtime/router/README.md` | `small` | `index-only` | Routing flow 已理解，activation rules 與 TTL policy 已對應到對應階段。 |
| `route.intelligence.engineering.heuristics` | `intelligence/engineering/heuristics/README.md` | `small` | `index-only` | 各 atom 有明確原則與決策表，可反查驗證。 |
| `route.runtime.context-ttl-doc` | `runtime/context/README.md` | `small` | `index-only` | TTL policy 已套用，prune strategy 已理解。 |
| `route.skill.app-development-guidance` | `skills/app-development-guidance/SKILL.md` | `specialized` | `source-backed` | 舊入口仍可讀，新 reference-first paths 可找到，且未 bulk migrate skill content。 |
| `route.skill.travel-planning` | `skills/travel-planning/SKILL.md` | `specialized` | `source-backed` | 舊入口仍可讀，新 reference-first paths 可找到，且未 bulk migrate skill content。 |
| `route.runtime.onboarding` | `runtime/onboarding/README.md` | `specialized` | `summary-first` | 各 quickstart 的步驟可依序執行，且與對應 SKILL.md 的 Quick Start 摘要一致。 |
| `route.analysis.apk.workflows` | `analysis/apk/workflows/README.md` | `specialized` | `summary-first` | 各 workflow 有明確步驟與產出格式，可依序執行。 |
| `route.intelligence.apk-analysis.atoms` | `intelligence/engineering/analytical-reasoning/README.md` | `specialized` | `summary-first` | 各 atom 有明確決策表或信號表，可反查驗證。 |
| `route.validation.ai-decision-contract` | `validation/README.md` | `specialized` | `summary-first` | Scenario 的 expected_route 與 actual trace 可比對；forbidden_routes 未被使用。 |
| `route.intelligence.engineering.agent-architecture` | `intelligence/engineering/agent-architecture/README.md` | `small` | `index-only` | 各 atom 有明確原則、症狀表與預防方式，可反查驗證。 |
| `route.feedback.history` | `feedback/history/README.md` | `small` | `index-only` | Lesson 已寫入 feedback/history/<domain>/ 對應分類，且 feedback/history/<domain>/README.md 已更新索引。 |

## Summary Records

| Atom ID | Lifecycle | File | Summary |
| --- | --- | --- | --- |
| `architecture.apk-analysis-pilot` | `candidate` | [`apk-analysis-pilot.md`](../summaries/apk-analysis-pilot.md) | `apk-analysis` 作為第一個 Workflow / Analysis / Intelligence 分離 pilot 的 migration map。它建立新 reference-first 候選目的地，但保留 `skills/apk-analysis/SKILL.md` 作為 active skill entrypoint。 |
| `intelligence.apk-highest-leverage-analysis` | `candidate` | [`apk-highest-leverage-analysis.md`](../summaries/apk-highest-leverage-analysis.md) | APK 分析 checkpoint 應先界定未知，再依 time-to-evidence、語意距離、安全性與 validation clarity 選擇最高收益路線。 |
| `skill.app-development-guidance` | `validated` | [`app-development-guidance.md`](../summaries/app-development-guidance.md) | 將授權 App/API/Embedded/Firmware 觀察轉成開發 guidance、實作模式、控制項、檢查清單。涵蓋 mobile（Android/iOS/Flutter/React Native）、backend API、embedded firmware 的安全控制、實作模式與 release gate。 |
| `architecture.context-cost-optimization` | `validated` | [`context-cost-optimization.md`](../summaries/context-cost-optimization.md) | Token 成本優化規劃。Phase 1（立即省錢）：Bootstrap 極小化（~800 tokens）、README 拆分、Rule lazy-load、Summary layer。Phase 2（架構升級）：Runtime Context Router、Context Cost Metadata、Skill Index、Context TTL。Phase 3（長期）：Semantic Retrieval、Episodic Memory、Multi-model Routing。 |
| `feedback.promotion.pipeline` | `candidate` | [`feedback-promotion-pipeline.md`](../summaries/feedback-promotion-pipeline.md) | 定義 feedback lesson 從 skill-local history 推進到 workflow、intelligence、shared-rules、memory 或 runtime surfaces 的 promotion / downgrade gate。 |
| `governance.goal-ledger-boundary` | `validated` | [`goal-ledger-boundary.md`](../summaries/goal-ledger-boundary.md) | `.agent-goals/` 只保存 active conversation goals；長期 roadmap、phase、migration、promotion、deprecation 與治理狀態必須落到 durable planning 文件。 |
| `knowledge.navigation` | `validated` | [`knowledge-navigation.md`](../summaries/knowledge-navigation.md) | 知識導航系統：indexes（任務路由）、summaries（300-500 token 摘要）、graphs（知識圖譜邊）、runtime（routing registry、refresh policy、SQLite lookup cache）。讓 agent 用最小 token 成本找到正確知識。 |
| `memory.operations` | `candidate` | [`memory-operations.md`](../summaries/memory-operations.md) | 長期記憶層：short-term（目前 task context）、episodic（過去 task 關鍵決策與結果）、project（專案歷史脈絡）、failure（反覆失效模式）。支援 similarity-based retrieval。 |
| `metadata.schema.knowledge-atom` | `validated` | [`metadata-schema.md`](../summaries/metadata-schema.md) | Knowledge Atom metadata schema v1，定義 atom 的必填欄位、選填欄位、受控值、YAML 範本與驗證規則。 |
| `models.routing` | `candidate` | [`model-routing.md`](../summaries/model-routing.md) | 多模型協作架構：capability profile（small/large/specialized）、compression strategy（checklist/compressed/full）、model-aware context report。根據 task 複雜度選擇模型與 context 策略。 |
| `governance.repo-maintenance` | `validated` | [`repo-governance.md`](../summaries/repo-governance.md) | AI-native Knowledge Operating System 本身的維護、升級、遷移與治理。涵蓋 lifecycle management、validation、cleanup、splitting rules、dependency maintenance。 |
| `root.bootstrap.ai-skill` | `validated` | [`root-bootstrap.md`](../summaries/root-bootstrap.md) | Ai-skill 工作的 bootstrap 入口。Root README 定義 OS layout 與 cost-aware 啟動流程；CORE_BOOTSTRAP.md 定義 3 條核心規則（~800 tokens）；shared-rules README 定義 Runtime Activation Model 與 lazy-load rules。 |
| `runtime.operations` | `validated` | [`runtime-operations.md`](../summaries/runtime-operations.md) | Runtime 層負責 context routing、dynamic loading、context pruning、agent coordination 與 orchestration。包含 router（activation rules、cost budget）、context（TTL policy、prune strategy）。 |
| `skill.travel-planning` | `candidate` | [`travel-planning.md`](../summaries/travel-planning.md) | 依目的地、日期、交通與玩法規劃行程，包含營業時間查證、交通比較、住宿與備案。支援 itinerary 結構化輸出與可行性檢查。 |

## Graph Records

| ID | Source | Status | Edges | File |
| --- | --- | --- | --- | --- |
| `graph.analysis-layers` | `analysis/README.md` | `candidate` | 17 | [`analysis-layers.yaml`](../graphs/analysis-layers.yaml) |
| `graph.analysis-repo-methods` | `analysis/repo/README.md` | `candidate` | 8 | [`analysis-repo-methods.yaml`](../graphs/analysis-repo-methods.yaml) |
| `graph.apk-analysis-pilot` | `architecture/apk-analysis-pilot-migration.md` | `candidate` | 6 | [`apk-analysis-pilot.yaml`](../graphs/apk-analysis-pilot.yaml) |
| `graph.apk-highest-leverage-analysis` | `intelligence/engineering/analytical-reasoning/highest-leverage-analysis-path.md` | `candidate` | 4 | [`apk-highest-leverage-analysis.yaml`](../graphs/apk-highest-leverage-analysis.yaml) |
| `graph.decisions-adr` | `decisions/README.md` | `candidate` | 9 | [`decisions-adr.yaml`](../graphs/decisions-adr.yaml) |
| `graph.feedback-layers` | `feedback/README.md` | `candidate` | 9 | [`feedback-layers.yaml`](../graphs/feedback-layers.yaml) |
| `graph.feedback-promotion-pipeline` | `feedback/promotion/README.md` | `candidate` | 6 | [`feedback-promotion-pipeline.yaml`](../graphs/feedback-promotion-pipeline.yaml) |
| `graph.governance-layers` | `governance/README.md` | `candidate` | 8 | [`governance-layers.yaml`](../graphs/governance-layers.yaml) |
| `graph.intelligence-agent-architecture` | `intelligence/engineering/agent-architecture/README.md` | `candidate` | 20 | [`intelligence-agent-architecture.yaml`](../graphs/intelligence-agent-architecture.yaml) |
| `graph.intelligence-anti-patterns` | `intelligence/engineering/anti-patterns/generic-repository-overuse.md` | `candidate` | 5 | [`intelligence-anti-patterns.yaml`](../graphs/intelligence-anti-patterns.yaml) |
| `graph.intelligence-apk-analysis-atoms` | `intelligence/engineering/analytical-reasoning/README.md` | `candidate` | 18 | [`intelligence-apk-analysis-atoms.yaml`](../graphs/intelligence-apk-analysis-atoms.yaml) |
| `graph.intelligence-app-development-guidance` | `intelligence/engineering/development/README.md` | `candidate` | 9 | [`intelligence-app-development-guidance.yaml`](../graphs/intelligence-app-development-guidance.yaml) |
| `graph.intelligence-architecture` | `intelligence/engineering/architecture/modular-monolith-vs-microservices.md` | `candidate` | 5 | [`intelligence-architecture.yaml`](../graphs/intelligence-architecture.yaml) |
| `graph.intelligence-business` | `intelligence/business/saas-pricing-heuristics.md` | `candidate` | 3 | [`intelligence-business.yaml`](../graphs/intelligence-business.yaml) |
| `graph.intelligence-distributed-systems` | `intelligence/engineering/distributed-systems/eventual-consistency-patterns.md` | `candidate` | 5 | [`intelligence-distributed-systems.yaml`](../graphs/intelligence-distributed-systems.yaml) |
| `graph.intelligence-domain` | `intelligence/engineering/domain/aggregate-boundary-heuristics.md` | `candidate` | 5 | [`intelligence-domain.yaml`](../graphs/intelligence-domain.yaml) |
| `graph.intelligence-failure` | `intelligence/engineering/failure/connection-leak-patterns.md` | `candidate` | 5 | [`intelligence-failure.yaml`](../graphs/intelligence-failure.yaml) |
| `graph.intelligence-heuristics` | `intelligence/engineering/heuristics/README.md` | `candidate` | 10 | [`intelligence-heuristics.yaml`](../graphs/intelligence-heuristics.yaml) |
| `graph.intelligence-repo-analysis` | `intelligence/engineering/analytical-reasoning/README.md` | `candidate` | 9 | [`intelligence-repo-analysis.yaml`](../graphs/intelligence-repo-analysis.yaml) |
| `graph.intelligence-tradeoffs` | `intelligence/engineering/tradeoffs/postgres-vs-mongodb.md` | `candidate` | 5 | [`intelligence-tradeoffs.yaml`](../graphs/intelligence-tradeoffs.yaml) |
| `graph.intelligence-travel` | `intelligence/travel/README.md` | `candidate` | 7 | [`intelligence-travel.yaml`](../graphs/intelligence-travel.yaml) |
| `graph.metadata-navigation` | `metadata/schema.md` | `candidate` | 7 | [`metadata-navigation.yaml`](../graphs/metadata-navigation.yaml) |
| `graph.runtime-onboarding` | `runtime/onboarding/README.md` | `candidate` | 13 | [`runtime-onboarding.yaml`](../graphs/runtime-onboarding.yaml) |
| `graph.runtime-pipeline` | `runtime/pipeline/README.md` | `candidate` | 12 | [`runtime-pipeline.yaml`](../graphs/runtime-pipeline.yaml) |
| `graph.runtime-prompt-artifacts` | `runtime/prompt-artifacts/README.md` | `candidate` | 8 | [`runtime-prompt-artifacts.yaml`](../graphs/runtime-prompt-artifacts.yaml) |
| `graph.source-boundary` | `governance/lifecycle/README.md` | `candidate` | 6 | [`source-boundary.yaml`](../graphs/source-boundary.yaml) |
| `graph.workflow-app-development-guidance` | `workflow/app-development-guidance/README.md` | `candidate` | 11 | [`workflow-app-development-guidance.yaml`](../graphs/workflow-app-development-guidance.yaml) |
| `graph.workflow-layers` | `workflow/README.md` | `candidate` | 16 | [`workflow-layers.yaml`](../graphs/workflow-layers.yaml) |
| `graph.workflow-travel-planning` | `workflow/travel-planning/README.md` | `candidate` | 9 | [`workflow-travel-planning.yaml`](../graphs/workflow-travel-planning.yaml) |

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
