# Knowledge Runtime Report

本檔由 `ai-skill runtime refresh` 產生，彙整 runtime registry、summaries、graphs 與 refresh policy 的目前狀態。

## Source Surfaces

| Surface | Path | Count / Status |
| --- | --- | --- |
| Routing registry | [`routing-registry.yaml`](routing-registry.yaml) | 51 records |
| Refresh policy | [`refresh-policy.yaml`](refresh-policy.yaml) | candidate |
| Model context report | [`model-context-report.md`](model-context-report.md) | generated view |
| Model checklists | [`model-checklists.md`](model-checklists.md) | generated view |
| SQLite runtime index | [`sqlite/`](sqlite/) | generated lookup cache prototype |
| Summaries | [`../summaries/`](../summaries/) | 20 files |
| Graph records | [`../graphs/`](../graphs/) | 33 files |

## Routing Records

| ID | Primary source | Model | Compression | Validation signal |
| --- | --- | --- | --- | --- |
| `route.bootstrap.ai-skill` | `CORE_BOOTSTRAP.md` | `small` | `summary-first` | Core Bootstrap 3 rules 已讀，git status 已檢查。 |
| `route.runtime.phase-machine` | `runtime/runtime.db` | `small` | `source-backed` | 目前 phase 已確認，allowed_actions 與 forbidden_actions 已載入，phase transition rules 已檢查。 |
| `route.runtime.obligation-ledger` | `runtime/runtime.db` | `small` | `source-backed` | 本 phase 的 obligations 已確認，pending obligations 已記錄，blocking gates 已檢查。 |
| `route.runtime.blocking-gates` | `runtime/runtime.db` | `small` | `source-backed` | 本 phase 的 blocking gates 已檢查，無未通過的 critical/high gates，phase transition 可進行。 |
| `route.runtime.recovery` | `runtime/runtime.db` | `small` | `source-backed` | 阻斷 gate 已對應到 strategy；retry limit、strategy change、domain policy、required reload set、execution graph rebuild 與 recovery validation 已檢查。 |
| `route.runtime.scheduler` | `runtime/runtime.db` | `small` | `source-backed` | Queue 已依 priority 排序，blocking gates 優先於 obligations，dependencies 已解析。 |
| `route.runtime.transactions` | `runtime/runtime.db` | `small` | `source-backed` | Transaction state 正確，所有 rules 已檢查，templates 已套用。 |
| `route.skill.discovery` | `knowledge/runtime/routing-registry.yaml` | `small` | `index-only` | Task intent 已對應到 routing-registry.yaml 的 triggers，entrypoint 與 summary path 可解析。 |
| `route.runtime.activation-rules` | `runtime/runtime.db` | `small` | `index-only` | 目前 task 已比對 activation-rules.yaml，符合條件的 rules 已載入，不符合的已 deferred。 |
| `route.runtime.context-ttl` | `governance/ai-runtime-governance/context-attention-governance.md` | `small` | `summary-first` | Summary-first、attention budget、decision externalization、recap checkpoint 與 task-boundary prune 已檢查；必要時再讀 TTL policy。 |
| `route.runtime.prompt-cache-alignment` | `runtime/context/prompt-cache-playbook.md` | `small` | `summary-first` | Prompt cache playbook、enforcement rule、metadata provider_cache 欄位與 activation rule 已同步。 |
| `route.governance.durable-goal-boundary` | `enforcement/conversation-goal-ledger.md` | `large` | `source-backed` | 長期狀態已落到 durable planning 文件，且 active goal 完成後才刪除。 |
| `route.governance.ai-runtime-five-step` | `governance/ai-runtime-governance/five-step-ai-governance.md` | `small` | `summary-first` | 新增項目已通過 necessity、deletion、simplification、cycle-time 與 automation-last 檢查，且沒有加入 always-load context。 |
| `route.runtime.context-loading` | `governance/lifecycle/routing-philosophy.md` | `small` | `summary-first` | Primary source、deferred sources、source-of-truth gate 與 validation signal 已可被記錄。 |
| `route.governance.routing-signal` | `governance/ai-runtime-governance/routing-signal-governance.md` | `small` | `summary-first` | Task intent 已先於 path signal 確認；primary_source、negative signals、multi-route disambiguation 與 route validation signal 已檢查。 |
| `route.governance.validation-scenario` | `governance/ai-runtime-governance/validation-scenario-governance.md` | `small` | `summary-first` | Stateless reproduction、no answer leakage、failure class、expected/forbidden behavior 與 traceability gates 已檢查。 |
| `route.governance.cognitive-state-evidence` | `governance/ai-runtime-governance/cognitive-state-governance.md` | `medium` | `summary-first` | Claim scope、evidence owner、confidence integrity、contradiction propagation、runtime primitive boundary 與 scenario coverage 已檢查。 |
| `route.metadata.knowledge-atom-schema` | `metadata/schema.md` | `large` | `source-backed` | 欄位可套用到第一批 atom candidates，且 Markdown links 可解析。 |
| `route.workflow.apk-analysis` | `workflow/apk-analysis/execution-flow.md` | `specialized` | `source-backed` | 新分層路徑可讀取，workflow 與 analysis 內容已分離。 |
| `route.intelligence.apk-highest-leverage-path` | `intelligence/engineering/analytical-reasoning/highest-leverage-analysis-path.md` | `specialized` | `source-backed` | 已比較可用 routes、選定 evidence-to-cost ratio 最高路線、記錄 fallback 與 attribution 回補條件。 |
| `route.feedback.promotion-pipeline` | `feedback/promotion/README.md` | `large` | `source-backed` | 原 lesson source 保留、promotion target 明確、runtime surfaces 與 close-loop validation 已同步。 |
| `route.models.model-aware-routing` | `models/README.md` | `large` | `source-backed` | Strategy、capability dimensions、compression、fallback behavior、workflow shape 與 validation target 可被記錄，且 validation/scenarios/models 覆蓋主要 routing cases。 |
| `route.memory.retrieval-activation` | `memory/README.md` | `large` | `source-backed` | Memory activation decision 能記錄 trigger、candidate memory、qualification、replay budget、current-source revalidation 與 discard / promotion decision。 |
| `route.runtime.router-flow` | `runtime/README.md` | `small` | `index-only` | Routing flow 已理解，activation rules 與 TTL policy 已對應到對應階段。 |
| `route.intelligence.engineering.heuristics` | `intelligence/engineering/heuristics/README.md` | `small` | `index-only` | 各 atom 有明確原則與決策表，可反查驗證。 |
| `route.runtime.context-ttl-doc` | `governance/lifecycle/context-ttl-philosophy.md` | `small` | `index-only` | TTL policy 已套用，prune strategy 已理解。 |
| `route.workflow.software-delivery` | `workflow/software-delivery/execution-flow.md` | `specialized` | `source-backed` | Change intake、requirements cognition、contract precedence、BDD closure、artifact completeness、test strategy、performance evidence、same-session closure 與 architecture fit analysis 已檢查；workflow、analysis、intelligence、metadata 與 governance 內容已分離。 |
| `route.workflow.greenfield` | `workflow/greenfield/execution-flow.md` | `specialized` | `source-backed` | 新分層路徑可讀取，workflow 與 templates 內容完整。 |
| `route.workflow.travel-planning` | `workflow/travel-planning/execution-flow.md` | `specialized` | `source-backed` | 新分層路徑可讀取，workflow 與 analysis 內容已分離。 |
| `route.workflow.documentation-ai-native` | `workflow/documentation/README.md` | `small` | `summary-first` | 新文件或目錄具 index-first 導航、分類維度（kind/audience/stability）已標註； README-as-router、停止條件與單一真相已符合 documentation context governance； 長文已按 document-sizing 拆分；語言與工具敘述依專案自訂 policy（本 route 不預設 tool-neutral）。  |
| `route.runtime.onboarding` | `runtime/onboarding/apk-analysis-setup.md` | `specialized` | `summary-first` | 各 quickstart 的步驟可依序執行，且與對應 workflow 的內容一致。 |
| `route.analysis.apk.workflows` | `analysis/apk/workflows/README.md` | `specialized` | `summary-first` | 各 workflow 有明確步驟與產出格式，可依序執行。 |
| `route.analysis.web` | `analysis/web/README.md` | `specialized` | `summary-first` | 目標網站已評估（技術棧、JS 需求、anti-bot 保護），工具已選擇（HTTP/Dynamic/Stealth）， 提取策略已設計（selector/adaptive parsing），風險已評估（legal/technical/data quality）。  |
| `route.intelligence.apk-analysis.atoms` | `intelligence/engineering/apk-analysis/README.md` | `specialized` | `summary-first` | 各 atom 有明確決策表或信號表，可反查驗證。 |
| `route.intelligence.requirements-cognition` | `intelligence/engineering/requirements/README.md` | `small` | `summary-first` | Impact / journey evidence → requirement → behavior contract → acceptance criteria → validation target → execution artifact is traceable; no unconfirmed feature is written as requirement; high coverage is not treated as test effectiveness without targeted proof when risk is high. |
| `route.intelligence.architectural-fit` | `intelligence/engineering/architecture/architectural-fit/README.md` | `small` | `summary-first` | Chosen strategy、rejected lighter option、rejected heavier option、fit evidence 與 upgrade/downgrade trigger 已輸出；沒有預設套用 DDD/CQRS/event sourcing。 |
| `route.validation.ai-decision-contract` | `validation/README.md` | `specialized` | `summary-first` | Scenario 的 expected_route 與 actual trace 可比對；forbidden_routes 未被使用。 |
| `route.intelligence.engineering.agent-architecture` | `intelligence/engineering/agent-architecture/README.md` | `small` | `index-only` | 各 atom 有明確原則、症狀表與預防方式，可反查驗證。 |
| `route.feedback.history` | `feedback/history/README.md` | `small` | `index-only` | Lesson 已寫入 feedback/history/<domain>/ 對應分類，且 feedback/history/<domain>/README.md 已更新索引。 |
| `route.runtime.decision-recording` | `runtime/runtime.db` | `small` | `source-backed` | 本輪若有鎖定決策，對應 tier 的檔案與 README 索引已更新。 |
| `route.decisions.adr` | `decisions/README.md` | `small` | `summary-first` | decisions/README.md 已讀取，ADR 清單已查詢，graph edge 已確認。 |
| `route.architecture.permanent-docs` | `architecture/README.md` | `small` | `summary-first` | architecture/README.md 已讀取，架構文件清單已查詢。 |
| `route.evaluations.scenario-results` | `evaluations/README.md` | `small` | `index-only` | Evaluation 記錄已依 scenario 分類存放，格式與 template 一致。 |
| `route.tools.metadata-routing` | `tools/README.md` | `small` | `index-only` | Tool metadata 已定義，compression 策略與 routing 規則已建立。 |
| `route.traces.decision-traces` | `traces/README.md` | `small` | `index-only` | Trace 記錄已依 scenario 分類存放，格式與 template 一致。 |
| `route.anti-patterns.runtime-patterns` | `anti-patterns/README.md` | `small` | `summary-first` | Anti-pattern 已依格式記錄，症狀、預防與恢復方式已定義。 |
| `route.runtime.compiler` | `runtime/runtime.db` | `small` | `source-backed` | 所有 modified sources 已編譯，runtime.db 的 generated_surfaces 表包含最新記錄，`ai-skill runtime validate` 回傳 exit 0。 |
| `route.runtime.intelligence-routing` | `runtime/runtime.db` | `small` | `index-only` | Task intent 已對應到 intelligence-routing.yaml 的 domain，applicable_phases 已檢查，domain README 已載入。 |
| `route.runtime.output-governance` | `runtime/README.md` | `small` | `source-backed` | Language consistency 已確認，sanitization 已通過，tool neutrality 已檢查， format compliance 已驗證，governance gates 全部通過。  |
| `route.runtime.distributed` | `runtime/README.md` | `small` | `source-backed` | Distributed locks 已正確 acquire/release，multi-agent coordination rules 已遵守， async job lifecycle 狀態轉換正確，無 deadlock 或 stale state。  |
| `route.governance.knowledge-update-flow` | `governance/lifecycle/knowledge-update-flow.md` | `small` | `source-backed` | 目前 knowledge update 的步驟已確認，entry conditions 已滿足，exit conditions 已檢查，reference sources 已載入；linked update completeness gates 已通過， 下一步驟已決定。  |

## Summary Records

| Atom ID | Lifecycle | File | Summary |
| --- | --- | --- | --- |
| `architecture.apk-analysis-pilot` | `new-layer-promoted` | [`apk-analysis-pilot.md`](../summaries/apk-analysis-pilot.md) | `apk-analysis` 作為第一個 Workflow / Analysis / Intelligence 分離 pilot 的 migration map。新分層已 promoted：`workflow/apk-analysis/` 是端到端執行入口，`analysis/apk/` 保存可重用觀察、拆解與證據取得方法，`intelligence/engineering/analytical-reasoning/` 保存 reusable decision intelligence。 |
| `intelligence.apk-highest-leverage-analysis` | `candidate` | [`apk-highest-leverage-analysis.md`](../summaries/apk-highest-leverage-analysis.md) | APK 分析 checkpoint 應先界定未知，再依 time-to-evidence、語意距離、安全性與 validation clarity 選擇最高收益路線。 |
| `intelligence.architectural-fit` | `candidate` | [`architectural-fit.md`](../summaries/architectural-fit.md) | Architecture selection 先評估 business complexity、invariant density、integration pressure、lifecycle 與 team boundary，再選 CRUD、DDD Lite、Full DDD、event-driven 或 microservices。DDD 是 selectable architecture strategy，不是 runtime invariant。 |
| `architecture.context-cost-optimization` | `validated` | [`context-cost-optimization.md`](../summaries/context-cost-optimization.md) | Token 成本優化規劃。Phase 1（立即省錢）：Bootstrap 極小化（~800 tokens）、README 拆分、Rule lazy-load、Summary layer。Phase 2（架構升級）：Runtime Context Router、Context Cost Metadata、Skill Index、Context TTL。Phase 2.5（規範層已實作）：Provider Prompt Cache Alignment，定義 stable prefix / volatile suffix 與 provider cache metadata。Phase 3（長期）：Semantic Retrieval、Episodic Memory、Multi-model Routing。 |
| `workflow.software-delivery` | `validated` | [`development-guidance.md`](../summaries/development-guidance.md) | 將授權 App/API/Embedded/Firmware 觀察轉成開發 guidance、實作模式、控制項、檢查清單。涵蓋 mobile（Android/iOS/Flutter/React Native）、backend API、embedded firmware 的安全控制、實作模式與 release gate。原 `skills/app-development-guidance/` 已刪除，所有內容已遷移至新分層。提供 5 個標準化輸出模板（change-brief / contract / bdd-scenario / implementation-plan / review-report），位於 `workflow/software-delivery/templates/`。另提供 Greenfield 標準化流程（`workflow/greenfield/`）與 Slash Command 模式（`ai-tools/slash-commands.md`）。 |
| `analysis.dual-token-audit` | `candidate` | [`dual-token-audit.md`](../summaries/dual-token-audit.md) | 系統內同時存在兩套以上 token 機制（JWT + JWE、HMAC + 對稱加密、平台 token + 廠商回調 token）時的審計方法。觀察點：代碼結構、key 管理、algorithm 宣告、token 流向、validation 一致性、replay 防護、log 外洩、error 訊息。Audit 五步：列 token universe → 畫 flow → key/alg matrix → 接縫盤點 → failure mode 對照。 |
| `feedback.promotion.pipeline` | `candidate` | [`feedback-promotion-pipeline.md`](../summaries/feedback-promotion-pipeline.md) | 定義 feedback lesson 從 skill-local history 推進到 workflow、intelligence、enforcement、memory 或 runtime surfaces 的 promotion / downgrade gate。 |
| `governance.goal-ledger-boundary` | `validated` | [`goal-ledger-boundary.md`](../summaries/goal-ledger-boundary.md) | `.agent-goals/` 只保存 active conversation goals；長期 roadmap、phase、migration、promotion、deprecation 與治理狀態必須落到 durable planning 文件。 |
| `knowledge.navigation` | `validated` | [`knowledge-navigation.md`](../summaries/knowledge-navigation.md) | 知識導航系統：indexes（任務路由）、summaries（300-500 token 摘要）、graphs（知識圖譜邊）、runtime（routing registry、refresh policy、SQLite lookup cache）。讓 agent 用最小 token 成本找到正確知識。 |
| `memory.operations` | `candidate` | [`memory-operations.md`](../summaries/memory-operations.md) | Memory 是 selective replay system：working buffer、summary、episodic、project、failure、decision 與 retrieval-governance。Replay 需要 trigger、qualification、budget、freshness/scope check 與 current source revalidation。 |
| `metadata.schema.knowledge-atom` | `validated` | [`metadata-schema.md`](../summaries/metadata-schema.md) | Knowledge Atom metadata schema v1，定義 atom 的必填欄位、選填欄位、受控值、YAML 範本、驗證規則與 provider prompt cache hints。 |
| `intelligence.migration-seeder-anti-patterns` | `candidate` | [`migration-seeder-anti-patterns.md`](../summaries/migration-seeder-anti-patterns.md) | 把大量業務資料（廠商目錄、商品/SKU、遊戲清單、權限矩陣）以巨型 `INSERT` 包進 schema migration，使資料 lifecycle 與 schema lifecycle 被強制綁定。訊號：單檔 >50KB、檔名含 dataSeeder、業務人員想改資料要工程師寫 migration。替代方案依資料性質：列舉留 migration、目錄走 application seeder/admin、大量參考資料用外部 CSV + bulk loader。 |
| `models.routing` | `candidate` | [`model-routing.md`](../summaries/model-routing.md) | Model-aware execution strategy：profiles、capabilities、routing、workflow adaptation、governance、runtime primitives 與 compression。用於選擇 behavior shape，不宣稱 provider model 已切換。 |
| `runtime.prompt-cache-alignment` | `candidate` | [`prompt-cache-alignment.md`](../summaries/prompt-cache-alignment.md) | Provider prompt cache 對齊規範。定義 stable prefix、semi-stable middle、volatile suffix 的 context layout，並說明 `cacheable` 與 `provider_cache_candidate` 的邊界。 |
| `governance.repo-maintenance` | `validated` | [`repo-governance.md`](../summaries/repo-governance.md) | 本系統的維護、升級、遷移與治理。涵蓋 lifecycle management、validation、cleanup、splitting rules、dependency maintenance。 |
| `intelligence.requirements-cognition` | `candidate` | [`requirements-cognition.md`](../summaries/requirements-cognition.md) | Requirements cognition 先用 Impact Map × Customer Journey Map 對齊 product impact，再用 BDD-lite 處理 ambiguity、actor intent、behavior boundary、acceptance criteria、traceability、validation target 與 test effectiveness，而不是 Gherkin everywhere。 |
| `root.bootstrap.ai-skill` | `validated` | [`root-bootstrap.md`](../summaries/root-bootstrap.md) | Ai-skill 工作的 bootstrap 入口。Root README 定義 OS layout 與 cost-aware 啟動流程；CORE_BOOTSTRAP.md 定義 3 條核心規則（~800 tokens）；enforcement README 定義 Runtime Activation Model 與 lazy-load rules。 |
| `runtime.operations` | `validated` | [`runtime-operations.md`](../summaries/runtime-operations.md) | Runtime 層負責 context routing、dynamic loading、context pruning、agent coordination 與 orchestration。包含 router（activation rules、cost budget）、context（TTL policy、prune strategy、prompt cache layout）。 |
| `workflow.travel-planning` | `candidate` | [`travel-planning.md`](../summaries/travel-planning.md) | 依目的地、日期、交通與玩法規劃行程，包含營業時間查證、交通比較、住宿與備案。支援 itinerary 結構化輸出與可行性檢查。日本自駕含 Mapcode 粒度規則（沿線景點 2km+ 需各停車點獨立一行）與查詢工具鏈。 |
| `intelligence.vendor-integration-architecture` | `candidate` | [`vendor-integration-architecture.md`](../summaries/vendor-integration-architecture.md) | 整合超過 3 個外部廠商（支付聚合、社群登入、IM、博弈聚合、廣告聯播等）時的整合策略選型。五種策略：A. Adapter/Strategy（單模組多實作）/ B. Compile-time submodule per vendor / C. Plugin SPI（runtime 載入）/ D. Out-of-process service / E. Hybrid 分層。N ≥ 10 必須跳出 compile-time module per vendor 模式，否則編譯時間、IDE、升級成本爆炸。 |

## Graph Records

| ID | Source | Status | Edges | File |
| --- | --- | --- | --- | --- |
| `graph.analysis-layers` | `analysis/README.md` | `candidate` | 22 | [`analysis-layers.yaml`](../graphs/analysis-layers.yaml) |
| `graph.analysis-repo-methods` | `analysis/repo/README.md` | `candidate` | 6 | [`analysis-repo-methods.yaml`](../graphs/analysis-repo-methods.yaml) |
| `graph.apk-analysis-pilot` | `plans/archived/2026-05-11-1129-apk-analysis-pilot-migration.md` | `new-layer-promoted` | 5 | [`apk-analysis-pilot.yaml`](../graphs/apk-analysis-pilot.yaml) |
| `graph.apk-highest-leverage-analysis` | `intelligence/engineering/analytical-reasoning/highest-leverage-analysis-path.md` | `candidate` | 4 | [`apk-highest-leverage-analysis.yaml`](../graphs/apk-highest-leverage-analysis.yaml) |
| `graph.ddd-architecture-governance` | `intelligence/engineering/architecture/domain-modeling/README.md` | `candidate` | 7 | [`ddd-architecture-governance.yaml`](../graphs/ddd-architecture-governance.yaml) |
| `graph.decisions-adr` | `decisions/README.md` | `candidate` | 11 | [`decisions-adr.yaml`](../graphs/decisions-adr.yaml) |
| `graph.feedback-layers` | `feedback/README.md` | `candidate` | 9 | [`feedback-layers.yaml`](../graphs/feedback-layers.yaml) |
| `graph.feedback-promotion-pipeline` | `feedback/promotion/README.md` | `candidate` | 6 | [`feedback-promotion-pipeline.yaml`](../graphs/feedback-promotion-pipeline.yaml) |
| `graph.governance-layers` | `governance/README.md` | `candidate` | 22 | [`governance-layers.yaml`](../graphs/governance-layers.yaml) |
| `graph.intelligence-agent-architecture` | `intelligence/engineering/agent-architecture/README.md` | `candidate` | 20 | [`intelligence-agent-architecture.yaml`](../graphs/intelligence-agent-architecture.yaml) |
| `graph.intelligence-anti-patterns` | `intelligence/engineering/anti-patterns/generic-repository-overuse.md` | `candidate` | 9 | [`intelligence-anti-patterns.yaml`](../graphs/intelligence-anti-patterns.yaml) |
| `graph.intelligence-apk-analysis-atoms` | `intelligence/engineering/apk-analysis/README.md` | `candidate` | 20 | [`intelligence-apk-analysis-atoms.yaml`](../graphs/intelligence-apk-analysis-atoms.yaml) |
| `graph.intelligence-architecture` | `intelligence/engineering/architecture/modular-monolith-vs-microservices.md` | `candidate` | 7 | [`intelligence-architecture.yaml`](../graphs/intelligence-architecture.yaml) |
| `graph.intelligence-business` | `intelligence/business/saas-pricing-heuristics.md` | `candidate` | 3 | [`intelligence-business.yaml`](../graphs/intelligence-business.yaml) |
| `graph.intelligence-app-development-guidance` | `intelligence/engineering/development/README.md` | `candidate` | 8 | [`intelligence-development-guidance.yaml`](../graphs/intelligence-development-guidance.yaml) |
| `graph.intelligence-distributed-systems` | `intelligence/engineering/distributed-systems/eventual-consistency-patterns.md` | `candidate` | 5 | [`intelligence-distributed-systems.yaml`](../graphs/intelligence-distributed-systems.yaml) |
| `graph.intelligence-domain` | `intelligence/engineering/domain/aggregate-boundary-heuristics.md` | `candidate` | 7 | [`intelligence-domain.yaml`](../graphs/intelligence-domain.yaml) |
| `graph.intelligence-failure` | `intelligence/engineering/failure/connection-leak-patterns.md` | `candidate` | 5 | [`intelligence-failure.yaml`](../graphs/intelligence-failure.yaml) |
| `graph.intelligence-heuristics` | `intelligence/engineering/heuristics/README.md` | `candidate` | 11 | [`intelligence-heuristics.yaml`](../graphs/intelligence-heuristics.yaml) |
| `graph.intelligence-repo-analysis` | `intelligence/engineering/analytical-reasoning/README.md` | `candidate` | 8 | [`intelligence-repo-analysis.yaml`](../graphs/intelligence-repo-analysis.yaml) |
| `graph.intelligence-tradeoffs` | `intelligence/engineering/tradeoffs/postgres-vs-mongodb.md` | `candidate` | 7 | [`intelligence-tradeoffs.yaml`](../graphs/intelligence-tradeoffs.yaml) |
| `graph.intelligence-travel` | `intelligence/travel/README.md` | `candidate` | 7 | [`intelligence-travel.yaml`](../graphs/intelligence-travel.yaml) |
| `memory-operations` | `memory/README.md` | `candidate` | 6 | [`memory-operations.yaml`](../graphs/memory-operations.yaml) |
| `graph.metadata-navigation` | `metadata/schema.md` | `candidate` | 8 | [`metadata-navigation.yaml`](../graphs/metadata-navigation.yaml) |
| `graph.requirements-cognition` | `intelligence/engineering/requirements/README.md` | `candidate` | 11 | [`requirements-cognition.yaml`](../graphs/requirements-cognition.yaml) |
| `graph.runtime-onboarding` | `runtime/onboarding/apk-analysis-setup.md` | `candidate` | 4 | [`runtime-onboarding.yaml`](../graphs/runtime-onboarding.yaml) |
| `graph.runtime-pipeline` | `runtime/README.md` | `candidate` | 12 | [`runtime-pipeline.yaml`](../graphs/runtime-pipeline.yaml) |
| `graph.runtime-prompt-artifacts` | `runtime/README.md` | `candidate` | 8 | [`runtime-prompt-artifacts.yaml`](../graphs/runtime-prompt-artifacts.yaml) |
| `graph.source-boundary` | `governance/lifecycle/README.md` | `candidate` | 6 | [`source-boundary.yaml`](../graphs/source-boundary.yaml) |
| `graph.workflow-layers` | `workflow/README.md` | `candidate` | 20 | [`workflow-layers.yaml`](../graphs/workflow-layers.yaml) |
| `graph.workflow-software-delivery-architecture` | `workflow/software-delivery/architecture/README.md` | `candidate` | 7 | [`workflow-software-delivery-architecture.yaml`](../graphs/workflow-software-delivery-architecture.yaml) |
| `graph.workflow-app-development-guidance` | `workflow/software-delivery/README.md` | `candidate` | 11 | [`workflow-software-delivery.yaml`](../graphs/workflow-software-delivery.yaml) |
| `graph.workflow-travel-planning` | `workflow/travel-planning/README.md` | `candidate` | 6 | [`workflow-travel-planning.yaml`](../graphs/workflow-travel-planning.yaml) |

## Refresh Decisions

| Decision value | Meaning |
| --- | --- |
| `refresh_now` | 由 `refresh-policy.yaml` 定義的 generated surface decision。 |
| `revalidate_only` | 由 `refresh-policy.yaml` 定義的 generated surface decision。 |
| `downgrade_confidence` | 由 `refresh-policy.yaml` 定義的 generated surface decision。 |
| `no_update_needed` | 由 `refresh-policy.yaml` 定義的 generated surface decision。 |

## Validation

- 產生前應先執行 `ai-skill runtime validate`。
- 產生後應執行 Markdown link check、lints、close-loop dry run、commit / push / readback。
- 本報告是 generated view，不取代 `routing-registry.yaml`、`refresh-policy.yaml`、summary 或 graph source files。
