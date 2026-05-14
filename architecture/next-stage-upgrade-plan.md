# AI-native Knowledge Operating System 下一階段升級規劃

本文件是下一階段架構升級規劃書。它承接 [`ai-native-knowledge-operating-system.md`](ai-native-knowledge-operating-system.md) 的 reference-first、goal ledger、failure learning、rule weight 與 close-loop 基礎，規劃如何從現有 skill-centered repository 演進成 AI Knowledge Runtime System。

> **⚠️ 成本優化補充規劃**：本文件專注於架構分層與遷移路徑。關於 **token 成本優化、Bootstrap 極小化、Rule lazy-load、Summary layer、Context TTL** 等立即省錢措施，請見獨立的 [`context-cost-optimization-plan.md`](context-cost-optimization-plan.md)。兩份文件互補，建議先讀成本優化規劃再讀本文件。

## 目前走到哪裡

已完成的基礎層：

- Root `README.md` 已是 AI-native Knowledge Operating System dashboard。
- `shared-rules/` 已建立 dependency reading、linked updates、conversation goal ledger、failure learning、rule weight、language consistency 等 operating rules。
- `architecture/ai-native-knowledge-operating-system.md` 已定義 reference-first、compatibility inventory、Phase 3 deprecation checklist。
- `analysis/`、`intelligence/`、`workflow/`、`runtime/`、`memory/`、`feedback/`、`models/`、`governance/`、`knowledge/`、`metadata/` 已建立第一版 README skeleton，定義責任邊界。
- `.agent-goals/` 已作為 project-local active goal ledger 使用，完成後刪除，不進 git。
- Cursor / Claude tool docs 已指向 central repository 與 shared-rule bootstrap。
- `knowledge/indexes/README.md` 已建立第一版 task intent routing table 與索引格式。
- `metadata/schema.md` 已建立 Knowledge Atom metadata schema v1，可套用到第一批 atom candidates。
- `apk-analysis` pilot migration map 已建立，並新增 `analysis/apk/`、`workflow/apk-analysis/`、`intelligence/engineering/analytical-reasoning/` 候選目的地。
- 新分層流程優先策略已建立：`governance/lifecycle/`、`governance/validation/`、metadata 子規則、`runtime/routing/`、`knowledge/summaries/`、`knowledge/graphs/`、`knowledge/runtime/`。
- `knowledge/runtime/routing-registry.yaml` 已建立第一版 machine-readable routing registry，包含 8 筆 sample routing records。
- `scripts/validate-knowledge-runtime.rb` 已建立 deterministic validation helper，檢查 routing registry、refresh policy、summaries 與 graph records。
- `scripts/generate-knowledge-runtime-report.rb` 已建立第一版 generated runtime report 工具，產出 `knowledge/runtime/runtime-report.md`。
- `scripts/generate-model-context-report.rb` 已建立第一版 model-aware context report 工具，產出 `knowledge/runtime/model-context-report.md`。
- `feedback/promotion/README.md` 已建立 feedback promotion pipeline surface。
- `knowledge/summaries/` 已建立第一批 6 個 Knowledge Atom summaries，覆蓋 root bootstrap、metadata schema、apk-analysis pilot、goal ledger boundary、APK highest-leverage route selection 與 feedback promotion pipeline。
- `intelligence/engineering/analytical-reasoning/highest-leverage-analysis-path.md` 已建立第一個實際 APK engineering intelligence atom。
- `knowledge/summaries/` 已新增 APK highest-leverage route selection summary。
- `knowledge/graphs/` 已建立 5 個 graph records：source-boundary、metadata-navigation、apk-analysis-pilot、apk-highest-leverage-analysis、feedback-promotion-pipeline。
- `knowledge/runtime/refresh-policy.yaml` 已建立 generated summaries / graphs / registry refresh 流程，定義 refresh、revalidate、downgrade 與 no update needed。
- `knowledge/runtime/sqlite/README.md` 與 `scripts/generate-runtime-sqlite-index.rb` / `query-runtime-index.rb` / `validate-runtime-sqlite-index.rb` 已建立 SQLite / FTS runtime index prototype，作為低 token 搜尋候選 source 的 generated cache，不作 source-of-truth。

### ✅ 已完成：Context Cost Optimization（Phase 1）

以下為 [`context-cost-optimization-plan.md`](context-cost-optimization-plan.md) 中已完成的項目：

- **Bootstrap 極小化**：`CORE_BOOTSTRAP.md` 建立（3 rules, ~800 tokens），取代舊 Default Bootstrap（12 rules, ~5000 tokens）。
- **README 拆分**：根 `README.md` 縮短為 ~80 行超短入口。
- **Rule Lazy-load 機制**：`shared-rules/README.md` 引入 Runtime Activation Model，15 條 lazy-load rules 定義觸發條件。
- **Skill Index**：`skills-index.yaml` 建立（13 skills，含 triggers、cost metadata、entrypoint、summary path）。
- **Runtime Router**：`runtime/router/activation-rules.yaml`（15 activation rules）、`runtime/router/README.md`（routing decision flow）。
- **Context TTL**：`runtime/context/ttl-policy.yaml`（20 context types）、`runtime/context/README.md`（prune strategy）。
- **Context Cost Metadata**：`metadata/schema.md` 的 `context_cost` 升級為 object（estimated_tokens、load_strategy、cacheable、ttl）。
- **Routing Registry 升級**：`knowledge/runtime/routing-registry.yaml` 升級 v2，所有 records 含 cost metadata。
- **Summary Layer 擴充**：新增 8 個 summaries（app-development-guidance、travel-planning、repo-governance、knowledge-navigation、runtime-operations、model-routing、memory-operations、context-cost-optimization），總數從 6 → 14。
- **Knowledge Index 更新**：`knowledge/indexes/README.md` 加入 `skills-index.yaml` 作為首要路由。

### ✅ 已完成：Runtime Quality & Safety（Phase 2）

以下為根據外部 review 建議，優先實作的 **Runtime Quality & Safety** 層：

- **Token Budget System**：[`runtime/budget/token-budget.yaml`](../runtime/budget/token-budget.yaml) — 120K default max_tokens、per-model budgets、per-layer budget allocation、70% warning / 90% hard stop thresholds。
- **Context Health Score**：[`runtime/health/context-health-score.yaml`](../runtime/health/context-health-score.yaml) — 4 維度（relevance 0.35、duplication 0.20、staleness 0.25、conflict 0.20）、composite score、healthy/warning/critical thresholds。
- **Circuit Breaker**：[`runtime/guards/circuit-breaker.yaml`](../runtime/guards/circuit-breaker.yaml) — 5 guards（recursive depth max 4、tool calls 20/task、context growth 30%/task、hallucination risk 4 factors、conflict rules）。
- **Context Pollution Detection**：[`runtime/guards/context-pollution.yaml`](../runtime/guards/context-pollution.yaml) — 5 signals（conversation length 50 turns、repetitive edits 5 edits、module count 20 modules、cross-reference depth 5 layers、token utilization 85%）。
- **Tool Metadata**：[`tools/metadata/README.md`](../tools/metadata/README.md) — 每個工具標註 cost（avg_input_tokens、avg_output_tokens）、risk、contexts、activation strategy、compression support。
- **Tool Lazy Activation**：[`tools/routing/README.md`](../tools/routing/README.md) — 5-step activation flow、tool explosion detection（recursive_search、repetitive_read、tool_chain_too_long、output_too_large）。
- **Tool Output Compression**：[`tools/compression/README.md`](../tools/compression/README.md) — 4 levels（raw 1.0x、summary 0.2-0.3x、structured 0.1-0.2x、minimal 0.05-0.1x）、per-output-type strategies。
- **Memory Architecture 子層**：[`memory/working/README.md`](../memory/working/README.md)、[`memory/summary/README.md`](../memory/summary/README.md)、[`memory/decision/README.md`](../memory/decision/README.md) — 3 子層（working session-local、summary ≤500 tokens、decision immutable numbered）。
- **Decision System（ADR）**：[`decisions/README.md`](../decisions/README.md) — ADR lifecycle（proposed → accepted → deprecated → superseded）、naming convention `ADR-{number}-{short-title}.md`。
- **Anti-patterns**：[`anti-patterns/README.md`](../anti-patterns/README.md) + 5 patterns（context-explosion、recursive-tool-loop、hallucination-loop、stale-summary、skill-pollution）。
- **Skills Metadata v2**：[`skills-index.yaml`](../skills-index.yaml) 升級至 v2，所有 13 skills 加入 weight、domains、dependencies、conflicts、priority.runtime。

### ✅ 已完成：Runtime Pipeline（Phase 3）

將所有 Runtime Quality & Safety 元件串接成可執行的 orchestration flow：

- **Pipeline 概覽**：[`runtime/pipeline/README.md`](../runtime/pipeline/README.md) — 元件接線圖（bootstrap → routing → execution → close-loop）、跨階段通訊表（10 個觸發事件，如 Token usage > 90% → Context Pollution auto-archive、Recursive depth > 4 → Force close-loop）
- **Session Lifecycle**：[`runtime/pipeline/session-lifecycle.yaml`](../runtime/pipeline/session-lifecycle.yaml) — 4 階段定義：Bootstrap（2000 tokens, 2 guards）、Routing（2500 tokens, 3 guards）、Execution（100000 tokens, 11 guards）、Close-loop（1000 tokens, 1 guard）
- **Progressive Context Expansion**：[`runtime/pipeline/context-flow.yaml`](../runtime/pipeline/context-flow.yaml) — 4 層級（summary ~500 → module summary ~1500 → detailed source ~4500 → raw source ~10000 tokens），每層有 cache policy（session/task TTL）、entry/exit conditions、no-skip-levels rule
- **Guard Chain**：[`runtime/pipeline/guard-chain.yaml`](../runtime/pipeline/guard-chain.yaml) — 每 stage 的 guard 執行順序（ordered by severity）、檢查頻率（per_tool_call / per_task / per_edit）、layered violation 行為（critical → halt, high/medium → warn）
- **Skill Relevance Engine**：[`runtime/pipeline/relevance-engine.yaml`](../runtime/pipeline/relevance-engine.yaml) — 3 維度 scoring（trigger_match 0.5 + domain_match 0.3 + weight 0.2）、threshold 0.5、conflict penalty ×0.5、dependency_missing penalty ×0.8、3 個 scoring examples

### ✅ 已完成：Feedback Promotion Pipeline（Phase 4）

將 feedback lesson 從 `skills/*/feedback_history/` 的原始觀察，透過機器可讀的 scoring、workflow 與 lifecycle automation，推進到 `workflow/`、`intelligence/`、`shared-rules/`、`memory/` 或 runtime surfaces：

- **Promotion Pipeline 概覽**：[`feedback/pipeline/README.md`](../feedback/pipeline/README.md) — pipeline 架構圖（feedback_history → Promotion Engine → Promotion Workflow → Target Layer）、與既有層的關係、使用方式
- **Promotion Engine**：[`feedback/pipeline/promotion-engine.yaml`](../feedback/pipeline/promotion-engine.yaml) — 5 維度 scoring（impact 0.30 + maturity 0.25 + frequency 0.20 + freshness 0.15 + urgency 0.10）、threshold 0.7 immediate / 0.5 backlog、5 種 promotion target decisions（shared-rules/intelligence/workflow/skill-doc/archive）、3 個 scoring examples（cross-skill validated lesson → 0.71 promote_to_skill_doc、single-technique experimental → 0.27 archive、cross-skill engineering intelligence → 0.74 promote_to_intelligence）
- **Promotion Workflow**：[`feedback/pipeline/promotion-workflow.yaml`](../feedback/pipeline/promotion-workflow.yaml) — 5 階段 workflow（assess-lesson → prepare-content → write-target → update-linked → validate-close-loop）、每階段有 entry/exit conditions、steps、output、rollback-on-validation-failure rule
- **Lifecycle Automation**：[`feedback/pipeline/lifecycle-automation.yaml`](../feedback/pipeline/lifecycle-automation.yaml) — 4 種 automation（auto-archive-cold 180 days no references score<0.4、auto-downgrade-stale 90 days no re-validation、periodic-promotion-check weekly recalculate score、cold-data-threshold-monitor 50 lessons per category trigger index）、完整 state machine（new → experimental → candidate → validated → promoted → archived）、6 條 automation rules

### 已完成：子目錄擴充（Phase 6-11）

以下為後續補齊的各層子目錄與內容：

- **Phase 6：Governance 子目錄** — `governance/cleanup/README.md`（5 種 duplicate 類型、splitting 規則、ownership boundary）、`governance/dependency/README.md`（graph 更新時機、edge type controlled vocabulary、validation checklist）
- **Phase 7：Analysis 子目錄** — `analysis/repo/README.md`（4 種分析方法）、`analysis/production/README.md`（5 種生產分析）、`analysis/issue/README.md`（5 種 issue 分類與 priority 計算）
- **Phase 8：Feedback 子目錄** — `feedback/replay/README.md`（5 種 trigger conditions、replay flow）、`feedback/extraction/README.md`（extraction threshold table、intelligence type mapping）、`feedback/refinement/README.md`（6 種 trigger、5 種 problem types、version management）
- **Phase 9：Workflow 子目錄** — `workflow/software-delivery/README.md`（5 種 review types）、`workflow/repo-analysis/README.md`（5 種 analysis types）、`workflow/travel-planning/README.md`（6 種 planning types）
- **Phase 10：Memory 子目錄** — `memory/episodic/README.md`（情境記憶）、`memory/project/README.md`（專案記憶）、`memory/failure/README.md`（失效記憶），memory/ 完整 6 子層
- **Phase 11：Decisions ADR** — 5 筆實際 ADR（Reference-First Migration Strategy、Intelligence vs Knowledge Separation、Three-Layer Architecture、Feedback Promotion Pipeline、Memory Architecture）

### 已完成：Phase 13-16

- **Phase 13：Intelligence atoms 填充** — 在全部 8 個子目錄（architecture、tradeoffs、failure、domain、anti-patterns、distributed-systems、business、travel）各建立 1 個 candidate intelligence atom，遵循 principle → rationale → when to apply → decision flow → common misuse → token impact 格式。
- **Phase 14：Task-specific prompt artifact generator** — 建立 `runtime/prompt-artifacts/` 層，含 7 個 task type templates（apk-analysis、app-development-guidance、repo-analysis、travel-planning、repo-governance、knowledge-navigation、feedback-promotion）與 4 個 composition rules（workflow-steps、intelligence-atoms、analysis-methods、knowledge-summary），支援 priority-based culling 與 conflict resolution。
- **Phase 15：Knowledge graph records** — 建立 17 個完整 graph records：10 個 intelligence atom graphs、2 個 runtime graphs（pipeline、prompt-artifacts）、5 個 layer subdirectory graphs（workflow、analysis、governance、feedback、decisions），含 edges、validation、metadata。
- **Phase 16：Skill content extraction** — 從 `skills/apk-analysis/` 提取以下內容到新分層：
  - `TOOLS.md` → `analysis/apk/tools-and-failures.md`：媒體驗證工具、自動化腳本安全邊界
  - `WORKFLOW.md` → `workflow/apk-analysis/execution-flow.md`：8 條 capture window 詳細規則（tab coverage、lazy-load、evidence validation、replay runner、window split、read-only override、classifier、smoke）
  - `DOCUMENTATION.md` → `workflow/apk-analysis/artifact-gates.md`：SDK live self-generation audit、identity material audit、UI architecture map template、API catalog detail requirements、sanitization rules、developer guidance notes、feedback lesson writing tips、backfill rules

### ✅ 已完成：Phase 17-19

- **Phase 17：App Development Guidance 內容提取** — 從 `skills/app-development-guidance/` 提取以下內容到新分層：
  - `WORKFLOW.md` → `workflow/software-delivery/execution-flow.md`：8 個 workflow sections（Start From Evidence + Change Intake、Docs-first BDD closure loop、SDK defect closure loop、Same-session closure、Performance test gate、Backfill rules、Validate、Feed Back Reusable Lessons）
  - `DOCUMENTATION.md` → `workflow/software-delivery/artifact-gates.md`：6 個 artifact sections（Reusable Note Structure、Keep Separate、Reusable Guidance Boundary、Required Linked Update Statement、Good Guidance、Avoid）
  - `WORKFLOW.md` §2-5 → `analysis/development-guidance/risk-translation.md`：5 個 analysis methods（Translate To Risk、Choose The Owner Layer、Define Controls、File The Guidance、Apply Required Linked Updates）
  - 建立 `analysis/development-guidance/README.md` 定義 scope 與 4 個 analysis methods
  - 更新 `workflow/software-delivery/README.md` 加入 extracted content table

- **Phase 18：Travel Planning 內容提取** — 從 `skills/travel-planning/` 提取以下內容到新分層：
  - `WORKFLOW.md` → `workflow/travel-planning/execution-flow.md`：17 個 workflow sections（Intake、Source Triage、Agency Benchmark、Location Verification、Stop Planning、Weather、Transport、Lodging、Route Shape、Country Checks、Feasibility、Schedule、Calendar Output、車中泊、Recommendation Pass、Final Verification）
  - `DOCUMENTATION.md` → `workflow/travel-planning/artifact-gates.md`：14 個 output templates（Itinerary Summary、Day Plan、Weather-Aware Options、Source Table、Calendar/App-Ready Table、Offline Checklist、Agency Benchmark、Stop Experience、Restaurant、Exact Location、Non-Driving Transport、Self-Drive Cost、車中泊 Quietness、Final Verification Checklist）
  - 更新 `workflow/travel-planning/README.md` 加入 extracted content table 與 4 個 workflow flows

- **Phase 19：APK Analysis Techniques 提取** — 從 `skills/apk-analysis/techniques/` 提取以下內容到 `analysis/apk/techniques/`：
  - `techniques/README.md` → `analysis/apk/techniques/README.md`：routing rules（4 categories）、category rules、migration notes
  - `techniques/flutter-dart-aot/README.md` → `analysis/apk/techniques/flutter-dart-aot.md`：When To Use、Core Guidance、Common Flow、Success Shape、Pitfalls
  - `techniques/http-api/README.md` → `analysis/apk/techniques/http-api.md`：When To Use、Core Output、API Catalog Shape、API Documentation Flow、Finish Gate、UI Automation
  - `techniques/local-proxy/README.md` → `analysis/apk/techniques/local-proxy.md`：When To Use、Core Guidance、Handler Hook Flow、Attribution
  - `techniques/media-hls/README.md` → `analysis/apk/techniques/media-hls.md`：When To Use、Core Guidance、Media Chain（5-layer model）

### ✅ 已完成：Phase 22 — Repo Analysis 內容提取

- **Phase 22：Repository Analysis 內容提取** — 從 `skills/app-development-guidance/process/README.md` 提取 repo-discovery 相關內容到新分層：
  - `process/README.md` §Existing Project Documentation Backfill → `analysis/repo/documentation-backfill.md`：8 種文件恢復規則、6 種 pipeline artifact 恢復方法、7 步恢復順序
  - `process/README.md` §Traceability Gate → `analysis/repo/traceability-gate.md`：5 種追溯連結、stable ID 類型、未實作行為標記
  - `process/README.md` §Contract Governance Gate → `analysis/repo/contract-governance.md`：6 級文件優先順序、5 種衝突處理規則
  - 更新 `analysis/repo/README.md`：加入 3 個新分析方法、已提取內容表、擴充產出格式
  - 更新 `workflow/repo-analysis/README.md`：加入已提取內容表、Documentation Backfill Flow
  - `techniques/README.md` → `analysis/apk/techniques/README.md`：routing rules（4 categories）、category rules、migration notes
  - `techniques/flutter-dart-aot/README.md` → `analysis/apk/techniques/flutter-dart-aot.md`：When To Use、Core Guidance、Common Flow、Success Shape、Pitfalls
  - `techniques/http-api/README.md` → `analysis/apk/techniques/http-api.md`：When To Use、Core Output、API Catalog Shape、API Documentation Flow、Finish Gate、UI Automation
  - `techniques/local-proxy/README.md` → `analysis/apk/techniques/local-proxy.md`：When To Use、Core Guidance、Handler Hook Flow、Attribution
  - `techniques/media-hls/README.md` → `analysis/apk/techniques/media-hls.md`：When To Use、Core Guidance、Media Chain（5-layer model）

### ✅ 已完成：Phase 24 — Intelligence Atoms for app-development-guidance + travel-planning + repo-analysis

為新提取的 3 個領域建立工程智慧 atoms：

| 領域 | 新 atoms | 位置 |
|------|----------|------|
| **app-development-guidance** | docs-first-bdd-closure、risk-translation-heuristic、contract-governance-heuristic | `intelligence/engineering/development/` |
| **travel-planning** | source-triage-heuristic、feasibility-build-heuristic | `intelligence/travel/` |
| **repo-analysis** | documentation-backfill-heuristic、traceability-heuristic | `intelligence/engineering/analytical-reasoning/` |

同時建立對應的 graph records（`intelligence-app-development-guidance.yaml`、`intelligence-repo-analysis.yaml`），並更新 `intelligence-travel.yaml` 加入新 atom edges。

### ✅ 已完成：Phase 25 — Extract RUNBOOK.md to runtime/onboarding/

從 `skills/apk-analysis/RUNBOOK.md` 提取內容到新 `runtime/onboarding/` 層：

| 新文件 | 提取內容 |
|--------|----------|
| [`runtime/onboarding/apk-analysis-setup.md`](runtime/onboarding/apk-analysis-setup.md) | 放置位置、開場提示詞、自動回饋提示 |
| [`runtime/onboarding/apk-analysis-completion.md`](runtime/onboarding/apk-analysis-completion.md) | 完成定義、第一輪分析順序 |
| [`runtime/onboarding/README.md`](runtime/onboarding/README.md) | 目錄索引與 scope 定義 |

### ✅ 已完成：Phase 26 — Extract app-development-guidance subdirectories

從 `skills/app-development-guidance/` 的 `implementation/`、`platforms/`、`languages/`、`controls/`、`checklists/` 子目錄提取內容到對應的新分層：

| 新文件 | 目標層 | 原始來源 | 說明 |
|--------|--------|----------|------|
| [`metadata/development-guidance/controls-catalog.md`](../metadata/development-guidance/controls-catalog.md) | `metadata/` | `skills/app-development-guidance/controls/` | 6 種跨平台安全控制（API Transport、Auth & Session、Local Storage、Logging & Telemetry、Anti-Tamper、Release Build） |
| [`metadata/development-guidance/implementation-catalog.md`](../metadata/development-guidance/implementation-catalog.md) | `metadata/` | `skills/app-development-guidance/implementation/` | 5 類實作模式（Backend、Mobile、Embedded、Tooling、Examples）與 contract-to-implementation 映射流程 |
| [`metadata/development-guidance/platforms-catalog.md`](../metadata/development-guidance/platforms-catalog.md) | `metadata/` | `skills/app-development-guidance/platforms/` | 4 類平台指引（Mobile、Web、Backend、Embedded） |
| [`metadata/development-guidance/languages-catalog.md`](../metadata/development-guidance/languages-catalog.md) | `metadata/` | `skills/app-development-guidance/languages/` | 4 類語言陷阱（Dart、Kotlin/Java、Swift、TypeScript） |
| [`workflow/software-delivery/review-checklists.md`](../workflow/software-delivery/review-checklists.md) | `workflow/` | `skills/app-development-guidance/checklists/` | 6 種審查 checklist（Mobile Design Review、Mobile PR Review、Mobile Release Review、API Security Review、Contract Governance Review、Embedded Firmware Review） |

> **後續修正**：4 個 catalog 文件原放在 `analysis/`，後發現其本質為分類索引而非分析方法，已於 2026-05-13 移至 `metadata/development-guidance/`。

### ✅ 已完成：Phase 27 — Extract process/ to workflow/software-delivery/

從 `skills/app-development-guidance/process/README.md` 提取 contract-first 開發流程到 `workflow/software-delivery/development-process.md`：

| 新文件 | 目標層 | 原始來源 | 說明 |
|--------|--------|----------|------|
| [`workflow/software-delivery/development-process.md`](../workflow/software-delivery/development-process.md) | `workflow/` | `skills/app-development-guidance/process/README.md` | Contract-first 開發流程：Default Flow、Required Contracts、Initial Documentation Pack、Product Brief Validation Gate、Change Intake Gate、Contract Governance Gate、Traceability Gate、BDD Execution Closure、Test Strategy Gate、Embedded/Hardware Flow、Missing Information Gate、Existing Project Documentation Backfill、Contract-First Rules、Definition of Ready/Done |

### ✅ 已完成：Phase 28 — Technique → Intelligence Pilot（flutter-dart-aot）

---

### ✅ 已完成：Phase 30 — Feedback History 提取

**目標**：將 `skills/apk-analysis/feedback_history/` 和 `skills/app-development-guidance/feedback_history/` 的 lessons 提取到 `feedback/` 層。

**實際執行結果**：
- 建立 `feedback/extraction/apk-analysis-index.md`（61 lessons，6 個目標層）
- 建立 `feedback/extraction/development-guidance-index.md`（40 lessons，6 個目標層）
- 101/101 個 lesson 檔案已加入 `# Extracted — See <target>` 標記
- 100% lessons 已對應到目標層（workflow/、intelligence/、analysis/、shared-rules/）

---

### ✅ 已完成：Phase 29 — 其餘 3 個 Techniques Decomposition

**目標**：將 http-api、local-proxy、media-hls 三個 technique 比照 flutter-dart-aot 模式拆解。

**實際執行結果**：

1. **http-api decomposition**：
   - Workflow → `analysis/apk/workflows/http-api-documentation-flow.md`（7 步驟操作流程：API Entry → Group Index → Per-API Detail → Coverage/Gap Matrix → SDK Mapping → Finish Gate → UI Automation）
   - Intelligence → `intelligence/engineering/analytical-reasoning/heuristics/api-documentation-completeness.md`（何時開始、何時完成、Field Confidence 判斷）
   - 舊檔案已標註：`skills/apk-analysis/techniques/http-api/README.md`、`analysis/apk/techniques/http-api.md`

2. **local-proxy decomposition**：
   - Workflow → `analysis/apk/workflows/local-proxy-hook-flow.md`（6 步驟操作流程：確認證據 → 識別 Handler → Hook → Cast Netty → 去敏 → 歸因）
   - Intelligence → `intelligence/engineering/analytical-reasoning/heuristics/local-proxy-routing-diagnosis.md`（Local Proxy vs TLS Pinning 判斷表）
   - Intelligence → `intelligence/engineering/analytical-reasoning/signals/local-proxy-detection.md`（主要/次要/排除信號表）
   - 舊檔案已標註：`skills/apk-analysis/techniques/local-proxy/README.md`、`analysis/apk/techniques/local-proxy.md`

3. **media-hls decomposition**：
   - Workflow → `analysis/apk/workflows/media-hls-analysis-flow.md`（7 步驟操作流程：分離控制面/資料面 → Playlist → Key → Segments → 合併 → 容器驗證）
   - Intelligence → `intelligence/engineering/analytical-reasoning/signals/media-type-detection.md`（Magic Bytes 參考表、靜態 vs 動畫判斷、Container Probe 指令）
   - 舊檔案已標註：`skills/apk-analysis/techniques/media-hls/README.md`、`analysis/apk/techniques/media-hls.md`

**成功驗證標準**：
- ✅ 每個 technique 至少有 1 個 workflow 檔案 + 1-2 個 intelligence atoms
- ✅ 舊 technique 檔案已標註 `# Intelligence Extracted`
- ✅ `analysis/apk/README.md` 和 `intelligence/engineering/analytical-reasoning/README.md` 待更新
- ⏳ `knowledge/indexes/README.md` 和 `knowledge/runtime/routing-registry.yaml` 待更新

**新增檔案**（7 個）：
- `analysis/apk/workflows/http-api-documentation-flow.md`
- `analysis/apk/workflows/local-proxy-hook-flow.md`
- `analysis/apk/workflows/media-hls-analysis-flow.md`
- `intelligence/engineering/analytical-reasoning/heuristics/api-documentation-completeness.md`
- `intelligence/engineering/analytical-reasoning/heuristics/local-proxy-routing-diagnosis.md`
- `intelligence/engineering/analytical-reasoning/signals/local-proxy-detection.md`
- `intelligence/engineering/analytical-reasoning/signals/media-type-detection.md`

**更新檔案**（6 個）：
- `skills/apk-analysis/techniques/http-api/README.md`
- `skills/apk-analysis/techniques/local-proxy/README.md`
- `skills/apk-analysis/techniques/media-hls/README.md`
- `analysis/apk/techniques/http-api.md`
- `analysis/apk/techniques/local-proxy.md`
- `analysis/apk/techniques/media-hls.md`

---


Phase 28 是策略轉折點：從「搬遷內容」轉為「拆解 techniques，workflow 進 analysis/，intelligence 進 intelligence/」。

#### 策略摘要

- **核心目標**：提升 AI decision quality，不是把分類變漂亮
- **Technique Decomposition**：拆解（不是搬遷），HOW TO DO 進 `analysis/`，HOW TO THINK 進 `intelligence/`
- **舊 techniques 保留**，只標註已提取
- **不做完整 meta architecture**，先讓 intelligence 活起來

#### 建立的新目錄

| 目錄 | 說明 |
|------|------|
| `intelligence/engineering/analytical-reasoning/heuristics/` | 啟發式判斷規則（何時該用哪個技術） |
| `intelligence/engineering/analytical-reasoning/anti-patterns/` | 可預防的錯誤模式 |
| `intelligence/engineering/analytical-reasoning/failure/` | 具體失敗模式與診斷 |
| `intelligence/engineering/analytical-reasoning/signals/` | 技術特徵辨識信號 |
| `analysis/apk/workflows/` | HOW TO DO 操作流程 |

#### 建立的新檔案

| 檔案 | 層 | 說明 |
|------|-----|------|
| [`analysis/apk/workflows/frida-hook-flow.md`](../analysis/apk/workflows/frida-hook-flow.md) | `analysis/` | Frida hook 操作流程（command、setup、步驟） |
| [`intelligence/engineering/analytical-reasoning/heuristics/hook-selection.md`](../intelligence/engineering/analytical-reasoning/heuristics/hook-selection.md) | `intelligence/` | Hook 策略選擇啟發式（決策表） |
| [`intelligence/engineering/analytical-reasoning/anti-patterns/early-hook-instability.md`](../intelligence/engineering/analytical-reasoning/anti-patterns/early-hook-instability.md) | `intelligence/` | 過早 hook 導致不穩定（症狀表） |
| [`intelligence/engineering/analytical-reasoning/failure/frida-spawn-race.md`](../intelligence/engineering/analytical-reasoning/failure/frida-spawn-race.md) | `intelligence/` | Frida spawn race condition（診斷與緩解） |
| [`intelligence/engineering/analytical-reasoning/signals/flutter-dart-aot-detection.md`](../intelligence/engineering/analytical-reasoning/signals/flutter-dart-aot-detection.md) | `intelligence/` | Flutter/Dart AOT 辨識信號（信號表） |
| [`notes/intelligence-extraction-observations.md`](../notes/intelligence-extraction-observations.md) | `notes/` | Extraction 過程記錄與觀察 |

#### 更新的既有檔案

| 檔案 | 變更 |
|------|------|
| `intelligence/engineering/analytical-reasoning/README.md` | 加入新子目錄結構與 scope |
| `analysis/apk/README.md` | 加入 workflows/ 目錄，更新 migration notes |
| `skills/apk-analysis/techniques/flutter-dart-aot/README.md` | 加入 `# Intelligence Extracted` 標註 |
| `analysis/apk/techniques/flutter-dart-aot.md` | 加入 `# Intelligence Extracted` 標註 |

#### 成功驗證標準

Pilot 成功 = AI 開始能做 decision routing：
- 以前：只會照流程 dump
- 現在：能根據 signal 改變策略

### 下一階段 Phase 規劃（29-33）

以下為 Phase 28 之後的具體執行階段：

| Phase | 優先級 | 目標 | 主要產出 | 依賴 |
|-------|--------|------|----------|------|
| **Phase 29** | P1 | ✅ 已完成 | 其餘 3 個 techniques decomposition（http-api、local-proxy、media-hls） | 每個 technique 拆出 workflow → `analysis/apk/workflows/` + intelligence atoms → `intelligence/engineering/analytical-reasoning/{heuristics,signals}/` | Phase 28（pilot 模式已驗證） |
| **Phase 30** | P1 | ✅ 已完成 | Feedback history 提取到 feedback/ 層 | `skills/apk-analysis/feedback_history/` 和 `skills/app-development-guidance/feedback_history/` 的 lessons 提取到 `feedback/extraction/`，建立 category index | Phase 29（techniques 完成後，feedback 可對應到已建立的 workflow/intelligence） |
| **Phase 31** | P2 | ✅ 已完成 | Intelligence Extraction Pipeline 抽象化 | [`governance/lifecycle/intelligence-extraction-pipeline.md`](governance/lifecycle/intelligence-extraction-pipeline.md) — 7-step pipeline（內容審計 → 類型判斷 → 拆解執行 → 格式轉換 → 標註來源 → 驗證 → 更新索引） | Phase 29（所有 techniques 完成後才有足夠經驗） |
| **Phase 32** | P2 | ✅ 已完成 | SKILL.md 分解 | 將各 skill 的 `SKILL.md` 中剩餘內容（Quick Start、Default Workflow、Output Style、Feedback Loop）提取到對應新層 | Phase 30（feedback 提取完成） |
| **Phase 33** | P4 | ✅ 已完成 | Skill-Specific Intelligence Extraction | 見 [`plans/skill-specific-extraction.md`](plans/skill-specific-extraction.md) | Phase 31（pipeline 驗證成功）+ Phase 32（SKILL.md 分解完成） |

---

#### Phase 29：其餘 3 個 Techniques Decomposition ✅ 已完成

**目標**：將 http-api、local-proxy、media-hls 三個 technique 比照 flutter-dart-aot 模式拆解。

**實際產出**：

| Technique | Workflow | Intelligence Atoms |
|-----------|----------|-------------------|
| http-api | `analysis/apk/workflows/http-api-documentation-flow.md` | `intelligence/engineering/analytical-reasoning/heuristics/api-documentation-completeness.md` |
| local-proxy | `analysis/apk/workflows/local-proxy-hook-flow.md` | `intelligence/engineering/analytical-reasoning/heuristics/local-proxy-routing-diagnosis.md` + `intelligence/engineering/analytical-reasoning/signals/local-proxy-detection.md` |
| media-hls | `analysis/apk/workflows/media-hls-analysis-flow.md` | `intelligence/engineering/analytical-reasoning/signals/media-type-detection.md` |

**注意**：http-api 和 media-hls 的決策智慧較少，未產生獨立的 failure atom；local-proxy 因涉及 routing 判斷，產生了 2 個 intelligence atoms（heuristic + signals）。

---

#### Phase 30：Feedback History 提取 ✅ 已完成

**目標**：將 `skills/apk-analysis/feedback_history/` 和 `skills/app-development-guidance/feedback_history/` 的 lessons 提取到 `feedback/` 層。

**實際數據**：
- `skills/apk-analysis/feedback_history/`：61 條 lessons（common/ 38 條、flutter-dart-aot/ 13 條、http-api/ 22 條、local-proxy/ 5 條、media-hls/ 3 條）
- `skills/app-development-guidance/feedback_history/`：40 條 lessons（common/ 20 條、controls/ 1 條）
- **總計 101 條 lessons，100% 已提取**

**實際產出**：

| 產出 | 路徑 | 說明 |
|------|------|------|
| APK Analysis Index | `feedback/extraction/apk-analysis-index.md` | 61 lessons 分類到 6 個目標層 |
| App Development Guidance Index | `feedback/extraction/development-guidance-index.md` | 40 lessons 分類到 6 個目標層 |
| 已標註 lessons | `skills/*/feedback_history/*/*.md` | 101/101 檔案加入 `# Extracted — See <target>` 標記 |

**提取策略**：
- Feedback history lessons 已是「已提取產品」（從真實經驗提煉的泛化 lesson），不需重新提取內容
- 策略是：**建立索引 + 標註提取狀態**，而非重新提取內容到目標層
- 每個 lesson 的 `# Extracted — See <target>` 指向其 `Promotion Target` 對應的新架構檔案

**目標層分佈**：

| 目標層 | apk-analysis | app-development-guidance |
|--------|:-----------:|:-----------------------:|
| `workflow/*/execution-flow.md` | 29 | 18 |
| `intelligence/engineering/analytical-reasoning/` | 9 | — |
| `analysis/*/workflows/` | 12 | — |
| `analysis/*/tools-and-failures.md` | 3 | — |
| `workflow/*/artifact-gates.md` | 4 | 3 |
| `shared-rules/` | 4 | 2 |
| `workflow/*/development-process.md` | — | 10 |
| `analysis/*/controls-catalog.md` | — | 1 |
| `analysis/*/implementation-catalog.md` | — | 6 |

---

#### Phase 31：Pilot 驗證 + Intelligence Extraction Pipeline ✅ 已完成

**目標**：從 4 個 techniques 的實際 extraction 經驗中抽象出可重複的 extraction pipeline。

**實際產出**：

| 產出 | 路徑 | 說明 |
|------|------|------|
| Intelligence Extraction Pipeline | [`governance/lifecycle/intelligence-extraction-pipeline.md`](governance/lifecycle/intelligence-extraction-pipeline.md) | 7-step pipeline：內容審計 → 類型判斷 → 拆解執行 → 格式轉換 → 標註來源 → 驗證 → 更新索引 |
| 路由索引更新 | `knowledge/indexes/README.md` | 加入「執行 Intelligence Extraction Pipeline」路由列 |

**Pipeline 核心設計**：

1. **內容審計（Content Audit）**：8 個審計維度（操作步驟、判斷決策、工具命令、失敗模式、信號偵測、反模式、通用建議、產出規範）
2. **類型判斷（Type Classification）**：10 種元素類型，各有對應目標層與格式模板
3. **拆解執行（Decomposition）**：三種模式 — Technique Decomposition（Phase 28-29）、Feedback History Extraction（Phase 30）、SKILL.md Decomposition（Phase 32 未來）
4. **格式轉換（Format Transformation）**：5 種格式模板（Heuristic、Anti-pattern、Failure、Signal、Workflow）
5. **標註來源（Source Annotation）**：`# Intelligence Extracted` 或 `# Extracted — See <target>` 標記
6. **驗證（Validation）**：完整性檢查 + 品質檢查 + 不強制 atomize 的情況
7. **更新索引（Index Update）**：7 個需更新的文件清單

**已驗證的模式**：
- Technique Decomposition（4 techniques，Phase 28-29）✅
- Feedback History Extraction（101 lessons，Phase 30）✅

**未完成項目**：
- **Pilot 驗證**（在實際 APK analysis session 中驗證 intelligence atoms 是否改善 AI 決策品質）— 需在真實工作階段中進行，無法在此環境執行
- `knowledge/runtime/routing-registry.yaml` 更新 — 待 pipeline 正式啟用後再加入 routing record

---

#### ✅ 已完成：Phase 32 — SKILL.md 分解

**目標**：將各 skill 的 `SKILL.md` 中剩餘內容提取到對應新層。

**執行摘要**：

| 子步驟 | 內容 | 狀態 |
|--------|------|------|
| 32a | 分析 3 個 SKILL.md 的剩餘內容 | ✅ |
| 32b | 建立 `workflow/travel-planning/artifact-gates.md`（travel-planning 缺少此文件） | ✅ |
| 32c | 提取 apk-analysis SKILL.md Quick Start 步驟 5+7 → `runtime/onboarding/apk-analysis-quickstart.md` | ✅ |
| 32d | 提取 app-development-guidance SKILL.md Quick Start → `runtime/onboarding/app-development-guidance-quickstart.md` | ✅（已刪除：內容已回歸 workflow/software-delivery/） |
| 32e | 提取 travel-planning SKILL.md Quick Start → `runtime/onboarding/travel-planning-quickstart.md` | ✅ |
| 32f | 瘦身 3 個 SKILL.md 為純 routing 文件 | ✅ |
| 32g | 更新 `runtime/onboarding/README.md` + `architecture/next-stage-upgrade-plan.md` | ✅ |
| 32h | 提交 + push | ✅ |

**建立的新檔案**（4 個）：

| 檔案 | 說明 |
|------|------|
| `workflow/travel-planning/artifact-gates.md` | Travel Planning 產出規範與品質門檻（19 項必備項目、4 項品質門檻） |
| `runtime/onboarding/apk-analysis-quickstart.md` | APK 分析 Quick Start 操作步驟（UI 架構地圖 + 持久化資產轉換） |
| `runtime/onboarding/travel-planning-quickstart.md` | Travel Planning Quick Start 18 步驟 |

**瘦身成果**：

| SKILL.md | 原始行數 | 瘦身後行數 | 減少比例 |
|----------|---------|-----------|---------|
| `skills/apk-analysis/SKILL.md` | 158 | ~55 | ~65% |
| `skills/app-development-guidance/SKILL.md` | 132 | ~65 | ~51% |
| `skills/travel-planning/SKILL.md` | 102 | ~55 | ~46% |

**瘦身策略**：
- Quick Start：保留 routing summary（1-2 行 per step），詳細操作步驟 → `runtime/onboarding/<skill>-quickstart.md`
- Output Style / Safety：→ `workflow/<domain>/artifact-gates.md`
- Feedback Loop：→ `shared-rules/feedback-lessons.md` + `feedback/` 層
- 保留：header metadata、Shared Policy、When To Use、Out Of Scope、Default Workflow（純 routing）

**關鍵發現**：
- `travel-planning` 缺少 `workflow/travel-planning/artifact-gates.md`，需先建立才能提取 Output Style
- `apk-analysis` 的 Quick Start 步驟 5（UI 架構地圖）和步驟 7（持久化資產）是主要操作內容，步驟 1-4 是 routing，步驟 6 是 heuristic，步驟 8 是 feedback
- `app-development-guidance` 的 Quick Start 15 步驟全部可提取，保留 routing summary 即可

---

#### ✅ 已完成：Phase 33 — Skill-Specific Intelligence Extraction（P4）

**目標**：為每個 skill 設計專屬的 extraction strategy，並執行剩餘未提取內容的提取。

**執行摘要**：

| 子步驟 | 內容 | 狀態 |
|--------|------|------|
| 33a | 提取 `skills/app-development-guidance/CHECKLIST.md` → `workflow/software-delivery/review-checklist.md` | ✅ |
| 33b | 提取 `skills/travel-planning/TOOLS.md` → `analysis/travel/sources-and-tools.md` | ✅ |
| 33c | 提取 `skills/travel-planning/README.md` → `analysis/travel/README.md` | ✅ |
| 33d | 更新 `plans/skill-specific-extraction.md` 為執行結果 | ✅ |
| 33e | 更新 `architecture/next-stage-upgrade-plan.md` | ✅ |
| 33f | 提交 + push | ✅ |

**建立的新檔案**（3 個）：

| 檔案 | 說明 |
|------|------|
| `workflow/software-delivery/review-checklist.md` | 完整審查檢查清單（15 個類別：Change Intake、Test Strategy、Performance Test、Product To Contract、Documentation Backfill、Contract Governance、Reusable Guidance Boundary、Embedded/Hardware Review、API And Transport、Auth/Tokens/Sessions、Local Storage、Flutter/Android Build、Logging/Telemetry、Anti-Tamper/Risk、Release Gate） |
| `analysis/travel/sources-and-tools.md` | 旅遊規劃來源與工具完整清單（Source Hierarchy、Japan Travel Sources、Agency Sources、Checks By Travel Type、Stop/Food Sources、Restaurant Sources、Lodging Sources、Route Shape Sources、Transport Sources、Country-Specific Driving、Exact Location Sources、Weather Sources、Road-Trip Stops、Schedule Feasibility、Calendar/App-Ready、車中泊 Quietness、Source Note Format、Red Flags） |
| `analysis/travel/README.md` | Travel Planning 分析方法（Core Goals、Scope、Usage Pattern、Relationships） |

**更新的既有檔案**（4 個）：

| 檔案 | 變更 |
|------|------|
| `skills/app-development-guidance/CHECKLIST.md` | 加入 `# Extracted — See workflow/software-delivery/review-checklist.md` |
| `skills/travel-planning/TOOLS.md` | 加入 `# Extracted — See analysis/travel/sources-and-tools.md` |
| `skills/travel-planning/README.md` | 加入 `# Extracted — See analysis/travel/README.md` |
| `plans/skill-specific-extraction.md` | 從 pending 更新為執行結果（含完整 extraction status 表、annotation status 表、validation 結果） |

**提取策略摘要**：

| Skill | 提取策略 | 已提取內容 |
|-------|---------|-----------|
| `apk-analysis` | Technique Decomposition（Phase 28-29）+ Feedback History Index（Phase 30）+ SKILL.md Decomposition（Phase 32） | 全部提取完成 |
| `app-development-guidance` | Catalog Extraction（Phase 26）+ Process Extraction（Phase 27）+ Feedback History Index（Phase 30）+ SKILL.md Decomposition（Phase 32）+ Checklist Extraction（Phase 33） | 全部提取完成 |
| `travel-planning` | Workflow/Artifact Extraction（Phase 18）+ SKILL.md Decomposition（Phase 32）+ Tools/README Extraction（Phase 33） | 全部提取完成 |

**關鍵發現**：
- 3 個 skills 的內容結構差異大，但 extraction 策略可歸納為 3 種模式：**Technique Decomposition**（apk-analysis）、**Catalog Extraction**（app-development-guidance）、**Direct Promotion**（travel-planning）
- `travel-planning` 的 `TOOLS.md`（298 行）和 `README.md`（72 行）是純參考資料，直接 promotion 到 `analysis/` 層即可，不需拆解
- `app-development-guidance/CHECKLIST.md`（169 行）是純審查清單，直接 promotion 到 `workflow/` 層即可，不需拆解
- 所有 3 個 skills 的內容現在已全部提取到新分層，舊檔案已標註提取狀態

---

### 尚未完成的下一階段

- **既有 `skills/` 仍同時承載 workflow、analysis 方法、工程智慧、templates 與 feedback lessons**（by reference-first design，舊入口維持 active）。所有內容已提取到新分層，舊檔案已標註提取狀態。舊 skills 的清理時間線定義於 [`governance/lifecycle/README.md`](governance/lifecycle/README.md) 的 Skills Deprecation Timeline（Phase A→D），目前處於 Phase C（techniques 已刪除，`skills/apk-analysis/techniques/` 和 `analysis/apk/techniques/` 已清空）。
- **Runtime surfaces 更新**：✅ 已完成（2026-05-12）。`knowledge/runtime/routing-registry.yaml` 新增 5 條路由（app-development-guidance、travel-planning、onboarding、apk-workflows、apk-intelligence-atoms）。`knowledge/indexes/README.md` 新增 3 條索引列（onboarding、workflows、intelligence atoms）。generated summaries、reports、SQLite index 已可一鍵重建。
- **Intelligence Extraction Pipeline 實戰驗證（P2）**：✅ 已完成（Phase 32-33 已實際執行 pipeline 的 3 種模式，pipeline 文件已更新狀態表，遺漏的 index 更新已補齊）
- **Skills Deprecation（Phase B）**：✅ 已完成（2026-05-12）。所有 8 個 techniques 檔案（4 skills/ + 4 analysis/）已從 `# Intelligence Extracted` 升級為 `# Deprecated — see <new path>`。詳細時間線見 [`governance/lifecycle/README.md`](governance/lifecycle/README.md)。
- **跨層一致性檢查**：✅ 已完成（2026-05-12）。`analysis/apk/workflows/README.md`、`intelligence/engineering/analytical-reasoning/{heuristics,anti-patterns,failure,signals}/README.md` 均已列出完整 atom 表格，無遺漏。
- **跨 skill intelligence promotion（apk-analysis → heuristics）**：✅ 已完成（2026-05-12）。從 `intelligence/engineering/analytical-reasoning/` 提取 2 個 proven generalized atoms 到 `intelligence/engineering/heuristics/`：
  - `api-documentation-completeness.md` → `field-confidence-judgment.md`（Field Confidence 判斷表：confirmed/candidate/needs capture/meaning unknown 等 8 種狀態）
  - `media-type-detection.md` → `magic-bytes-reference.md`（Magic Bytes 參考表：MP4/WebP/GIF/PNG/JPEG/MP3/FLAC/TS/WebM）
  - 原始 atoms 已加入 `# Cross-Domain Promotion` 標註
- **跨 skill intelligence promotion（app-development-guidance → heuristics）**：✅ 已完成（2026-05-12）。從 `intelligence/engineering/development/` 提取 1 個 proven generalized atom 到 `intelligence/engineering/heuristics/`：
  - `contract-governance-heuristic.md` → `document-priority-hierarchy.md`（文件優先順序階層：6 層優先級、衝突分類流程、常見誤用表）
  - 原始 atom 已加入 `# Cross-Domain Promotion` 標註
- **Skills Deprecation（Phase C）**：✅ 已完成（2026-05-12）。已刪除 10 個舊 technique 檔案（`skills/apk-analysis/techniques/` 的 4 個子目錄 + 1 個 README，`analysis/apk/techniques/` 的 4 個 .md + 1 個 README）。刪除前已確認 Phase C 檢查清單 7 項條件全部滿足。15+ 個引用舊路徑的檔案已更新為指向新路徑。

### 下一階段 Phase 規劃（29-33）

以下為 Phase 28 之後的具體執行階段：

| Phase | 優先級 | 目標 | 主要產出 | 依賴 |
|-------|--------|------|----------|------|
| **Phase 29** | P1 | ✅ 已完成 | 其餘 3 個 techniques decomposition（http-api、local-proxy、media-hls） | 每個 technique 拆出 workflow → `analysis/apk/workflows/` + intelligence atoms → `intelligence/engineering/analytical-reasoning/{heuristics,anti-patterns,failure,signals}/` | Phase 28（pilot 模式已驗證） |
| **Phase 30** | P1 | ✅ 已完成 | Feedback history 提取到 feedback/ 層 | `skills/apk-analysis/feedback_history/` 和 `skills/app-development-guidance/feedback_history/` 的 lessons 提取到 `feedback/extraction/`，建立 category index | Phase 29（techniques 完成後，feedback 可對應到已建立的 workflow/intelligence） |
| **Phase 31** | P2 | ✅ 已完成 | Intelligence Extraction Pipeline 抽象化 | [`governance/lifecycle/intelligence-extraction-pipeline.md`](governance/lifecycle/intelligence-extraction-pipeline.md) — 7-step pipeline（內容審計 → 類型判斷 → 拆解執行 → 格式轉換 → 標註來源 → 驗證 → 更新索引） | Phase 29（所有 techniques 完成後才有足夠經驗） |
| **Phase 32** | P2 | ✅ 已完成 | SKILL.md 分解 | 將各 skill 的 `SKILL.md` 中剩餘內容（Quick Start、Default Workflow、Output Style、Feedback Loop）提取到對應新層 | Phase 30（feedback 提取完成） |
| **Phase 33** | P4 | ✅ 已完成 | Skill-Specific Intelligence Extraction | 見 [`plans/skill-specific-extraction.md`](plans/skill-specific-extraction.md) | Phase 31（pipeline 驗證成功）+ Phase 32（SKILL.md 分解完成） |

## 核心問題

下一階段要回答的不只是「有哪些 prompts 或 skills」，而是：

- AI 如何工作。
- AI 如何學習。
- AI 如何沉澱知識。
- AI 如何找到正確知識。
- AI 如何演化知識。
- AI 如何多模型協作。
- AI 如何長期維護知識。

因此整體方向要從 **Skill Collection** 升級為 **AI Knowledge Runtime System**。

## 目標架構分層

下一階段建議正式拆分：

```text
analysis/
intelligence/
workflow/
runtime/
memory/
feedback/
models/
governance/
knowledge/
metadata/
```

這些目錄不是一次搬完所有內容，而是先建立責任邊界、metadata schema 與 navigation layer，再逐批遷移。

## 各層責任

### `analysis/`

負責「如何觀察與拆解」。

目前結構：

```text
analysis/
  apk/                     ← APK 分析方法（pilot migrated）
  development-guidance/    ← 開發指引分析方法（風險翻譯、控制層選擇）
  repo/                    ← Repository 分析與理解方法
  production/              ← Production 問題分析與根因追蹤
  issue/                   ← Issue 分類與優先級判斷
```

核心責任：

- reverse engineering。
- 流程拆解。
- 技術觀察。
- pattern extraction。
- 分析方法。
- Repository 結構觀察與心智模型建立。
- Production incident RCA 與效能診斷。
- Issue triage、優先級計算與重複偵測。

不應承載過多：

- trade-off。
- architecture lesson。
- anti-pattern conclusion。

這些應抽取到 `intelligence/`。

### `intelligence/`

負責「沉澱工程智慧與領域知識」。本層不是百科知識（什麼是 Redis、什麼是 CQRS），而是**經過經驗抽象化後的工程智慧**——AI 的「專家腦內模型」。

#### 與其他層的差異

| 層 | 偏 | 範例 |
|---|----|------|
| `knowledge/` | 事實 | Redis supports pub/sub |
| `skills/` | 執行流程 | How to debug Redis latency |
| `intelligence/` | 判斷力與經驗法則 | If Redis latency spikes suddenly, check connection lifecycle before scaling |

#### 核心內容

- **Heuristics（經驗法則）** — 資深工程師直覺
- **Tradeoffs（取捨）** — 「沒有銀彈」的理解
- **Pattern Recognition（模式辨識）** — 可重複的設計與反設計模式
- **Failure Recognition（災難辨識）** — 抽象化後的失敗模式
- **Decision Intelligence（決策智慧）** — 架構與技術選擇的判斷力
- **Contextual Thinking（情境思考）** — 何時適用、何時不適用的邊界條件

建議結構：

```text
intelligence/
  engineering/
    architecture/          # 架構思考模式（非教學）
    domain/                # DDD / 業務模型智慧
    failure/               # 工程災難智慧（抽象化失敗模式）
    heuristics/            # 經驗法則（intelligence 核心）
    anti-patterns/         # 常見錯誤設計
    tradeoffs/             # 技術取捨智慧
    distributed-systems/   # 分散式系統生存經驗
  business/                # 商業決策智慧
  travel/                  # 特定領域智慧（Personal Domain Intelligence）
```

#### 各子目錄說明

| 目錄 | 核心 | 範例內容 |
|------|------|---------|
| `engineering/architecture/` | 架構判斷力 | `modular-monolith-vs-microservices.md`、`event-driven-tradeoffs.md`、`cqrs-when-not-to-use.md` |
| `engineering/domain/` | 業務建模經驗 | `aggregate-boundary-heuristics.md`、`inventory-domain-patterns.md` |
| `engineering/failure/` | AI 的「危險雷達」 | `connection-leak-patterns.md`、`distributed-lock-failure.md` |
| `engineering/heuristics/` | 資深工程師直覺 | `premature-optimization.md`、`abstraction-threshold.md`、`retry-smell.md` |
| `engineering/anti-patterns/` | AI 自動避雷 | `generic-repository-overuse.md`、`shared-database-microservices.md`、`god-service-pattern.md` |
| `engineering/tradeoffs/` | 「沒有銀彈」的理解 | `postgres-vs-mongodb.md`、`websocket-vs-polling.md`、`sqlite-vs-postgres.md` |
| `engineering/distributed-systems/` | 大型系統生存經驗 | `eventual-consistency-patterns.md`、`event-ordering-risks.md`、`backpressure-signals.md` |
| `business/` | 商業判斷力 | `saas-pricing-heuristics.md`、`enterprise-sales-warning.md`、`pmf-signals.md` |
| `travel/` | 領域經驗法則 | `japan-roadtrip-fatigue.md`、`ski-trip-routing.md`、`golden-week-avoidance.md` |

#### 與根目錄 `anti-patterns/` 的邊界

| 位置 | 內容 |
|------|------|
| [`anti-patterns/`](../../anti-patterns/) | Agent 操作層的 anti-patterns（context explosion、recursive tool loop、hallucination loop 等） |
| `intelligence/engineering/anti-patterns/` | 工程領域的 anti-patterns（設計錯誤、架構錯誤等） |

`intelligence/` 是 Senior Engineer Brain。

### `workflow/`

負責「AI 如何執行工作」。

目前結構：

```text
workflow/
  apk-analysis/              ← APK 分析執行流程（pilot migrated）
  software-delivery/         ← 軟體交付執行流程
  repo-analysis/             ← Repository 分析執行流程
  travel-planning/           ← 旅遊規劃執行流程
```

核心責任：

- planning flow。
- task decomposition。
- review flow。
- orchestration flow。
- execution flow。
- App 開發審查流程（design review、code review、release review、security review）。
- Repository 分析流程（onboarding、deep analysis、migration impact、tech debt assessment）。
- 旅遊規劃流程（itinerary planning、transportation research、budget planning）。

`workflow/` 應 reference `intelligence/`，而不是內嵌大量知識。

### `runtime/`

負責「AI 系統如何運作」。

建議結構：

```text
runtime/
  scheduler/
  routing/
  orchestration/
  context/
  budget/                ← Token Budget System
  health/                ← Context Health Score
  guards/                ← Circuit Breaker + Context Pollution Detection
  pipeline/              ← Session lifecycle, guard chain, relevance engine
  prompt-artifacts/      ← 新增：Task-specific prompt artifact generator
```

核心責任：

- dynamic loading。
- context injection。
- orchestration。
- task routing。
- context pruning。
- agent coordination。
- **token budget management**（`runtime/budget/token-budget.yaml`）。
- **context health scoring**（`runtime/health/context-health-score.yaml`）。
- **circuit breaker & guards**（`runtime/guards/circuit-breaker.yaml`、`runtime/guards/context-pollution.yaml`）。
- **session lifecycle management**（`runtime/pipeline/session-lifecycle.yaml`）。
- **prompt artifact generation**（`runtime/prompt-artifacts/`）— 根據 task type 自動組合 prompt 結構，引用 workflow/ 的執行步驟、intelligence/ 的工程智慧、analysis/ 的分析方法，產出針對當前任務優化的 prompt artifact。

### `tools/`

負責「AI 工具如何被管理與優化」。

建議結構：

```text
tools/
  metadata/        ← 工具成本、風險、activation strategy
  routing/         ← 工具 lazy activation、explosion detection
  compression/     ← 工具輸出壓縮（4 levels）
```

核心責任：

- 工具成本標註（avg_input_tokens、avg_output_tokens、risk）。
- 工具 lazy activation（preload / lazy / on_demand）。
- 工具爆炸偵測（recursive_search、repetitive_read、tool_chain_too_long、output_too_large）。
- 工具輸出壓縮（raw / summary / structured / minimal）。

### `memory/`

負責「長期記憶」。

目前結構：

```text
memory/
  working/         ← Session-local, discardable 工作記憶
  summary/         ← 壓縮 session 歷史（≤500 tokens）
  decision/        ← 輕量 ADR（immutable, numbered）
  episodic/        ← 情境記憶（跨 session 經驗 recall）
  project/         ← 專案記憶（跨 session 專案脈絡）
  failure/         ← 失效記憶（抽象化失效模式）
```

核心責任：

- experience replay。
- long-term memory。
- historical context。
- **session-local working memory**（[`memory/working/README.md`](memory/working/README.md)）。
- **compressed session summaries**（[`memory/summary/README.md`](memory/summary/README.md)）。
- **architecture decision records**（[`memory/decision/README.md`](memory/decision/README.md)）。
- **episodic memory**（[`memory/episodic/README.md`](memory/episodic/README.md)）— 跨 session 情境經驗 recall。
- **project memory**（[`memory/project/README.md`](memory/project/README.md)）— 跨 session 專案脈絡保持。
- **failure memory**（[`memory/failure/README.md`](memory/failure/README.md)）— 抽象化失效模式記錄。

### `decisions/`

負責「架構決策記錄（ADR）」。

目前結構：

```text
decisions/
  README.md                                    ← ADR 系統說明
  ADR-001-reference-first-migration-strategy.md ← Reference-First 遷移策略
  ADR-002-intelligence-vs-knowledge-separation.md ← Intelligence vs Knowledge 分離
  ADR-003-three-layer-architecture.md          ← Three-Layer Architecture
  ADR-004-feedback-promotion-pipeline.md       ← Feedback Promotion Pipeline
  ADR-005-memory-architecture.md               ← Memory Architecture（6 子層）
```

核心責任：

- 記錄關鍵架構決策。
- ADR lifecycle（proposed → accepted → deprecated → superseded）。
- 避免重複討論相同決策。
- 提供決策歷史追溯。

現有 ADR：

| ADR | 決策 | 狀態 |
| --- | --- | --- |
| ADR-001 | Reference-First Migration Strategy — 保留 `skills/` 為 source of truth，新分層只建立 reference/summary/index | accepted |
| ADR-002 | Intelligence vs Knowledge Separation — `intelligence/` 獨立於 `knowledge/`，兩者為平行層 | accepted |
| ADR-003 | Three-Layer Architecture — knowledge（事實）/ skills（流程）/ intelligence（判斷）三層平行 | accepted |
| ADR-004 | Feedback Promotion Pipeline — replay → extraction → refinement → promotion 五階段 pipeline | accepted |
| ADR-005 | Memory Architecture — 6 子層記憶模型（working/summary/decision/episodic/project/failure） | accepted |

### `anti-patterns/`

負責「已知失效模式記錄」。

建議結構：

```text
anti-patterns/
  README.md                    ← 索引
  context-explosion.md         ← Context 無限制增長
  recursive-tool-loop.md       ← 工具反覆呼叫無進展
  hallucination-loop.md        ← 無 canonical source 時過度推理
  stale-summary.md             ← Summary 與 source 不同步
  skill-pollution.md           ← 不相關 skill 浪費 token
```

核心責任：

- 記錄已知失效模式。
- 提供 detection signal 與 prevention strategy。
- 讓 agent 快速辨識並避免。

### `feedback/`

負責「系統如何持續演化」。

目前結構：

```text
feedback/
  replay/        ← 經驗重播（從過往 session 與 failure 提取教訓）
  extraction/    ← 智慧抽取（從 feedback 提煉 intelligence atom）
  refinement/    ← 流程精煉（持續改進 workflow）
  promotion/     ← 推廣管道（lesson 推進到各層）
  pipeline/      ← 自動化 pipeline YAML（lifecycle、promotion engine、workflow）
```

核心責任：

- workflow refinement。
- intelligence extraction。
- lesson replay。
- knowledge evolution。
- 經驗重播的觸發條件、流程與結果格式。
- Intelligence atom 的抽取門檻、流程與必備欄位。
- Workflow 的持續改進流程、觸發條件與版本管理。

### `models/`

負責「不同模型如何協作」。

建議結構：

```text
models/
  claude/
  gpt/
  gemini/
  qwen/
  small-model/
```

核心責任：

- capability profile。
- reasoning strength。
- context limit。
- routing strategy。
- compression strategy。
- prompt adaptation。

### `governance/`

負責「知識治理與系統維護」。

建議結構：

```text
governance/
  cleanup/
  splitting/
  lifecycle/
  validation/
```

核心責任：

- duplicate cleanup。
- lifecycle management。
- validation。
- splitting rules。
- dependency maintenance。

### `knowledge/`

負責「知識導航與知識圖譜」。

建議結構：

```text
knowledge/
  atoms/
  indexes/
  summaries/
  graphs/
  runtime/
```

核心思想是 Atomic Knowledge。真正目標不是單純拆小文件，而是支援 Dynamic Context Composition。

不要讓系統變成 Knowledge Fragment Hell；每個 atom 都必須能被 index、summary、graph 與 runtime metadata 找到。

### `metadata/`

負責「知識控制系統」。

建議結構：

```text
metadata/
  rules/
  ranking/
  confidence/
  compatibility/
```

`metadata/` 是 Rule Metadata System 的核心。Metadata 不是只描述文件，而是控制 runtime 行為。

每個 Knowledge Atom 應包含：

```yaml
id:
type:
domain:
tags:
priority:
confidence:
stability:
complexity:
context_cost:
depends:
related:
conflicts:
models:
summary:
checklist:
```

Runtime 依賴 metadata 進行：

1. Context Routing：現在該載入哪些知識。
2. Priority Selection：哪些規則優先。
3. Conflict Resolution：規則衝突時如何仲裁。
4. Dynamic Loading：根據 task 載入知識。
5. Model-aware Compression：小模型只讀 checklist 或 compressed knowledge。
6. Knowledge Promotion：`candidate` → `validated` → `stable`。
7. Knowledge Cleanup：找出過期知識。
8. Dependency Graph Construction：自動建立 knowledge graph。

## Knowledge Navigation System

Atomic Knowledge 必須搭配 navigation + index system。

建議建立：

```text
knowledge/indexes/
knowledge/summaries/
knowledge/graphs/
knowledge/runtime/
```

真正重要的不是知識量，而是 AI 能否找到正確知識。

### SQLite / FTS Runtime Index（prototype）

SQLite 適合導入為 generated runtime lookup cache，用來降低 agent 在大型 Markdown / YAML repository 中的初始讀取成本。它不保存 canonical truth；只保存可重建的 index rows，讓 agent 先用 task intent、keyword、tag、layer、priority、confidence 或 context cost 找到少量 candidate sources。

原則：

- Canonical source 仍是 Markdown / YAML：`skills/*/feedback_history/`、`shared-rules/`、`knowledge/summaries/`、`knowledge/graphs/`、`knowledge/runtime/routing-registry.yaml`。
- SQLite DB 不預設 commit；commit schema、generator、query helper 與 validation tests。
- Query result 只回傳少量 `source_path`、summary、tags、score 與 validation signal；需要執行、修改、promotion 或高信心結論時仍讀全文。
- Feedback lessons 可以被 index，但 lesson 全文仍留在 skill-local `feedback_history/`，直到 migration policy 明確改變。
- SQLite schema 應由 `metadata/schema.md` 與 runtime registry 控制，避免另創一套欄位語意。

第一版 tooling：

```text
knowledge/runtime/sqlite/README.md
scripts/generate-runtime-sqlite-index.rb
scripts/query-runtime-index.rb
scripts/validate-runtime-sqlite-index.rb
```

預期效益：

- 用 SQLite / FTS 先縮小候選集，再讀 source，可降低 token 消耗。
- Feedback lessons、summaries、graphs、routing registry 可用同一套 lookup path 查詢。
- 因 DB 可重建，未來不會把 generated cache 和 canonical source 混在一起。

Cold-data archive 的觸發門檻放在 `governance/lifecycle/README.md`：當單一 skill feedback lesson 超過約 50 條、單一 category 超過約 20 條，或 agent 為了找 lesson 需要讀大量 history/index 時，應優先使用 generated summary / SQLite FTS 作候選查詢，再按需讀 canonical source。

## Intelligence Feedback Loop

系統應形成閉環：

```text
Analysis -> Extraction -> Intelligence -> Workflow -> Feedback
```

例：

```text
apk-analysis
  -> intelligence extraction
  -> realtime intelligence
  -> workflow reference
  -> future refinement
```

## Multi-model Runtime Architecture

未來模型一定是混用，因此 workflow 應 model-aware。

範例：

```yaml
small-model:
  use:
    - checklist
    - compressed knowledge

large-model:
  use:
    - full intelligence graph
```

## Knowledge Lifecycle System

知識一定會熵增，因此每個知識單元需要 lifecycle：

```text
temporary/
candidate/
validated/
stable/
deprecated/
```

## 遷移原則

1. 不一次搬完所有檔案。
2. 先建立 top-level directory README，定義責任邊界。
3. 先定義 metadata schema，再遷移 content。
4. 先選一個 skill 做示範遷移，再擴展到其他 skill。
5. 保留 `skills/` 與 `shared-rules/` 相容層，直到 workflow / intelligence / metadata / runtime 的 reference path 穩定。
6. 每次搬移都必須保留舊連結或提供 redirect / index。
7. 每次遷移都要經過 `.agent-goals`、linked updates、diff review、commit/push/readback、clean status。

## 建議遷移階段

### Phase 0：目前已完成的基礎

- OS dashboard。
- `reference-first`。
- `rule-weight`。
- goal ledger。
- failure learning。
- language consistency。
- compatibility inventory。
- Phase 3 deprecation checklist。

### Phase 1：建立新架構目錄

Status: top-level README skeletons 已建立；尚未開始 bulk content migration。

建立下列目錄與 README：

```text
analysis/
intelligence/
workflow/
runtime/
memory/
feedback/
models/
governance/
knowledge/
metadata/
```

每個 README 只定義：

- 該層責任。
- 放什麼。
- 不放什麼。
- 與現有 `skills/`、`shared-rules/`、`ai-tools/` 的關係。
- 第一批候選遷移來源。

### Phase 2：Metadata System

Status: `metadata/schema.md` 已建立 Knowledge Atom schema v1；`metadata/rules/`、`metadata/ranking/`、`metadata/confidence/`、`metadata/compatibility/` 已建立第一版操作規則。

新增：

```text
metadata/schema.md
metadata/rules/
metadata/ranking/
metadata/confidence/
metadata/compatibility/
```

定義 Knowledge Atom schema 與 required/optional 欄位。

### Phase 3：Knowledge Navigation

Status: `knowledge/indexes/README.md` 已建立第一版 routing index；`knowledge/summaries/`、`knowledge/graphs/`、`knowledge/runtime/` 已建立格式與責任邊界，尚未生成大量內容。

新增：

```text
knowledge/indexes/
knowledge/summaries/
knowledge/graphs/
knowledge/runtime/
```

先做 index 與 summary，不急著做完整 graph runtime。

### Phase 4：Workflow / Intelligence 分離

Status: `apk-analysis` pilot 已完成 content extraction（6 個新檔案）；`intelligence/` 已重建為 9 個子目錄的專家智慧層；舊 `skills/apk-analysis/` 已不再作為 active entrypoint。

已完成項目：

- `analysis/apk/`：已提取 traffic triage 與 tools/failures 兩份分析文件。
- `workflow/apk-analysis/`：已提取 execution flow 與 artifact gates 兩份工作流程文件。
- `intelligence/engineering/analytical-reasoning/`：已提取 evidence-first-routing 與 live-readiness-gates 兩個 validated-intelligence atoms。
- `intelligence/` 重建為 9 個子目錄（architecture、domain、failure、heuristics、anti-patterns、tradeoffs、distributed-systems、business、travel），每個有 README.md 定義 scope。
- `knowledge/indexes/README.md`：新增 4 條 routing entries 指向新提取檔案。
- `architecture/apk-analysis-pilot-migration.md`：狀態更新為 `content-extracted`。

尚未完成：

- 將 `skills/apk-analysis/` 中更多內容（techniques/、feedback_history/）逐步提取到對應新層。
- 將 `skills/app-development-guidance/` 的 implementation 與 controls 提取到 intelligence。
- 將 `shared-rules/failure-patterns/` 中偏工程判斷的 pattern 摘要提取到 intelligence。

### Phase 5：Runtime / Models

Status: ✅ **已完成**。所有子項目已實作完畢。

已實作項目：

- `runtime/routing/README.md` — context routing 流程（7 步驟路由決策表）。
- `knowledge/runtime/routing-registry.yaml` — machine-readable routing registry。
- `models/profiles/README.md` — small / large / specialized 三種 profile，含 routing rules 與 metadata mapping。
- `models/compression/README.md` — 5 層 compression（index-only / summary-first / checklist-first / source-backed / graph-assisted），含 profile defaults 與 escalation rules。
- `runtime/budget/token-budget.yaml` — Token Budget System（120K default、per-model budgets、per-layer allocation、70%/90% thresholds）。
- `runtime/health/context-health-score.yaml` — Context Health Score（4 維度 composite score、healthy/warning/critical thresholds）。
- `runtime/guards/circuit-breaker.yaml` — Circuit Breaker（5 guards：recursive depth、tool calls、context growth、hallucination risk、conflict rules）。
- `runtime/guards/context-pollution.yaml` — Context Pollution Detection（5 signals、composite pollution score、auto-archive on critical）。
- `tools/metadata/README.md` + `tools/routing/README.md` — Tool Metadata & Lazy Activation（tool cost/risk/activation schema、explosion detection）。
- `tools/compression/README.md` — Tool Output Compression（4-level compression、per-output-type strategies）。
- `memory/working/README.md` + `memory/summary/README.md` + `memory/decision/README.md` — Memory Architecture 3 子層。
- `decisions/README.md` — Decision System（ADR lifecycle、naming convention）。
- `anti-patterns/README.md` + 5 patterns — Anti-patterns。
- `skills-index.yaml` — Skills Metadata v2（weight/domains/dependencies/conflicts/priority.runtime）。
- `runtime/pipeline/README.md` — Runtime Pipeline（7 階段 pipeline 架構與元件間通訊）。
- `runtime/context/README.md` — Context TTL & Pruning（session/task/conversation TTL 類型與 prune 策略）。

定義：

- context routing。
- dynamic loading。
- context pruning。
- model capability profiles。
- small-model / large-model 使用策略。

### Phase 6：Lifecycle / Governance

Status: ✅ **已完成**。所有子項目已實作完畢。

已實作項目：

- `governance/lifecycle/README.md` — Knowledge Lifecycle（6 階段 lifecycle states：source-of-truth / candidate-map / candidate-atom / validated-atom / promoted / deprecated；promotion gates；cold data archive；update strategy；deletion rule）。
- `governance/validation/README.md` — Knowledge Validation Gates（必要 gates 表、migration validation checklist、generated refresh checklist、pass / block rules）。
- `governance/cleanup/README.md` — Duplicate Cleanup & Splitting（5 種 duplicate 類型、偵測流程、splitting 規則與門檻、ownership boundary 表、清理執行流程）。
- `governance/dependency/README.md` — Dependency Graph Maintenance（graph 更新時機表、graph record 維護流程、edge type controlled vocabulary、依賴變更連動更新表、graph validation 檢查項目）。
- `governance/README.md` — 整體 governance surface，已連結所有 4 個子目錄。

定義：

- knowledge lifecycle。
- duplicate cleanup。
- dependency graph maintenance。
- validation gates。
- deprecation / archive process。

## Durable Roadmap Goals

`.agent-goals/` 只追蹤目前對話的 active implementation task；長期 phase、未完成能力、migration / promotion / deprecation 狀態保留在本表與相關 layer 文件。當本表某項被拉進本輪工作時，才建立 `.agent-goals/` entry；完成驗證後刪除 active goal，並把 durable 狀態回寫到本表或對應文件。

| Priority | Status | Goal | Durable location | Next action | Completion criteria |
| --- | --- | --- | --- | --- | --- |
| P1 | done | 建立 next-stage upgrade plan | `architecture/next-stage-upgrade-plan.md` | 已完成 | 規劃書 commit/push/readback，root/architecture 入口可找到 |
| P1 | done | 建立 top-level architecture directories | `analysis/`, `intelligence/`, `workflow/`, `runtime/`, `memory/`, `feedback/`, `models/`, `governance/`, `knowledge/`, `metadata/` | 已完成 | 每個目錄責任邊界清楚，不搬移大量內容 |
| P2 | done | 設計 metadata schema | `metadata/schema.md` | 已完成 | Schema 可套用到第一批 Knowledge Atom |
| P2 | done | 建立 knowledge navigation index | `knowledge/indexes/README.md` | 已完成 | Agent 能從 index 找到 task-relevant knowledge |
| P2 | done | 遷移第一個 skill 作為示範 | `architecture/apk-analysis-pilot-migration.md`, `analysis/apk/`, `workflow/apk-analysis/`, `intelligence/engineering/analytical-reasoning/` | 已完成 content extraction（6 files），舊入口仍 active | 舊入口仍可用，新路徑可被 reference-first 找到 |
| P1 | done | 建立新分層運作流程 | `governance/`, `metadata/`, `runtime/routing/`, `knowledge/` | 已完成第一版流程與格式 | 舊 `skills/` 維持 source of truth，新分層可作 routing / promotion / validation surface |
| P1 | done | 規範 active goal 與 durable roadmap 邊界 | `shared-rules/conversation-goal-ledger.md`, `shared-rules/content-layering.md`, `governance/lifecycle/README.md` | 已完成 | `.agent-goals/` 不作長期 archive；刪除 active goal 前需回寫 durable planning |
| P1 | done | 建立 machine-readable routing registry | `knowledge/runtime/routing-registry.yaml`, `runtime/routing/README.md` | 已完成第一版 registry 與 8 筆 sample routing records | Runtime 可用結構化資料從 task intent 找到 primary source、dependencies、candidate summaries 與 validation signal |
| P1 | done | 建立第一批 Knowledge Atom summaries | `knowledge/summaries/` | 已完成 root bootstrap、metadata schema、apk-analysis pilot、goal ledger boundary summaries | Summaries 指向 canonical source，且不取代 source-of-truth 文件 |
| P2 | done | 建立初版 knowledge graph records | `knowledge/graphs/` | 已完成 source-boundary、metadata-navigation、apk-analysis-pilot 三個 graph records | Graph records 描述 depends / related / preserves_entrypoint，不使用 replacement semantics |
| P2 | done | 建立 model-aware routing / compression strategy | `models/profiles/`, `models/compression/`, `runtime/routing/README.md` | 已完成 small / large / specialized profiles 與 compression levels | Model profile 可被 runtime routing 與 summaries 引用 |
| P2 | done | 設計 generated summaries / graph refresh 流程 | `governance/validation/`, `knowledge/runtime/refresh-policy.yaml` | 已完成 refresh / revalidate / downgrade / no update needed 流程 | Source 變更時有明確 revalidation / downgrade path |
| P1 | done | 建立 registry / refresh validation helper | `scripts/validate-knowledge-runtime.rb`, `governance/validation/README.md`, `knowledge/runtime/README.md` | 已完成 deterministic helper | Helper 可檢查 registry、refresh policy、summaries、graphs 的必要欄位、YAML / Markdown 格式與 canonical paths |
| P1 | done | 建立 runtime report generator | `scripts/generate-knowledge-runtime-report.rb`, `knowledge/runtime/runtime-report.md` | 已完成 deterministic report generator | Report 可由 registry、refresh policy、summaries、graphs 重新產生，並通過 runtime validator 與 Markdown link check |
| P1 | done | 建立 model-aware context report generator | `scripts/generate-model-context-report.rb`, `knowledge/runtime/model-context-report.md`, `models/README.md` | 已完成 deterministic model context report | Report 可依 routing registry 的 model profile / compression level 重新產生，並通過 runtime validator 與 Markdown link check |
| P1 | done | 建立第一個 APK engineering intelligence atom | `intelligence/engineering/analytical-reasoning/highest-leverage-analysis-path.md`, `knowledge/summaries/apk-highest-leverage-analysis.md`, `knowledge/graphs/apk-highest-leverage-analysis.yaml` | 已完成最高收益路線 candidate intelligence atom | Old skill entrypoint remains active；runtime registry、summary、graph 與 knowledge index 可 route 到此 atom |
| P1 | done | 建立 feedback promotion pipeline surface | `feedback/promotion/README.md`, `knowledge/summaries/feedback-promotion-pipeline.md`, `knowledge/graphs/feedback-promotion-pipeline.yaml` | 已完成 promotion / downgrade design surface | Lesson source 保留於 `feedback_history/`；runtime registry、summary、graph 與 knowledge index 可 route 到 promotion pipeline |
| P1 | done | 建立 SQLite / FTS runtime index prototype | `knowledge/runtime/sqlite/README.md`, `scripts/generate-runtime-sqlite-index.rb`, `scripts/query-runtime-index.rb`, `scripts/validate-runtime-sqlite-index.rb` | 已完成本機 generated lookup cache、ranked query helper、filtering 與 stale checksum validator | SQLite 作為 generated lookup cache，不提交 DB binary；feedback lessons 只被索引，不搬離 `feedback_history/` |
| P1 | done | 定義 cold feedback lesson archive lifecycle | `governance/lifecycle/README.md`, `knowledge/runtime/sqlite/README.md`, `feedback/README.md`, `memory/README.md` | 已完成冷資料觸發門檻與 source-of-truth 邊界 | Lesson 超過門檻時先使用 generated summary / SQLite FTS 查候選；Markdown 仍是 canonical source |
| P1 | done | 建立 knowledge runtime refresh orchestrator | `scripts/refresh-knowledge-runtime.rb`, `knowledge/runtime/README.md`, `governance/validation/README.md` | 已完成一鍵重建 reports / SQLite index 並執行 validators | Generated runtime surfaces 可用單一命令重建與驗證，降低 stale cache 風險 |
| P1 | done | 建立 knowledge graph query helper | `scripts/query-knowledge-graph.rb`, `knowledge/graphs/README.md`, `knowledge/runtime/README.md` | 已完成 source / target / type / keyword graph edge 查詢 | Graph query 只回傳候選 edge list；修改或高信心判斷仍讀 graph YAML 與 canonical source |
| P1 | done | 建立 model checklist generator | `scripts/generate-model-checklists.rb`, `knowledge/runtime/model-checklists.md`, `models/README.md` | 已完成 per-model context-loading checklist artifact | Checklist 由 routing registry 生成；需要修改或高信心判斷仍讀 model docs 與 canonical source |
| P1 | done | 建立 Token Budget System | `runtime/budget/token-budget.yaml` | 已完成 120K default max_tokens、per-model budgets、per-layer allocation、70%/90% thresholds | Token 用量可預測，不再因深度 reasoning 爆 token |
| P1 | done | 建立 Context Health Score | `runtime/health/context-health-score.yaml` | 已完成 4 維度 composite score、healthy/warning/critical thresholds | Context 健康度可量化，在惡化前主動介入 |
| P1 | done | 建立 Circuit Breaker | `runtime/guards/circuit-breaker.yaml` | 已完成 5 guards（recursive depth、tool calls、context growth、hallucination risk、conflict rules） | Agent 不再陷入無限迴圈或工具爆炸 |
| P1 | done | 建立 Context Pollution Detection | `runtime/guards/context-pollution.yaml` | 已完成 5 signals、composite pollution score、auto-archive on critical | Context 污染可自動偵測與歸檔 |
| P1 | done | 建立 Tool Metadata & Lazy Activation | `tools/metadata/README.md`, `tools/routing/README.md` | 已完成 tool cost/risk/activation schema、explosion detection | 工具層級 token 消耗可預測與控制 |
| P1 | done | 建立 Tool Output Compression | `tools/compression/README.md` | 已完成 4-level compression、per-output-type strategies | 工具輸出 token 減少 50-95% |
| P1 | done | 建立 Memory Architecture 子層 | `memory/working/README.md`, `memory/summary/README.md`, `memory/decision/README.md`, `memory/episodic/README.md`, `memory/project/README.md`, `memory/failure/README.md` | 已完成 6 子層（working/summary/decision/episodic/project/failure） | 記憶管理精準，不再單一 memory 層 |
| P1 | done | 建立 Decision System（ADR） | `decisions/README.md`, `decisions/ADR-001-reference-first-migration-strategy.md`, `decisions/ADR-002-intelligence-vs-knowledge-separation.md`, `decisions/ADR-003-three-layer-architecture.md`, `decisions/ADR-004-feedback-promotion-pipeline.md`, `decisions/ADR-005-memory-architecture.md` | 已完成 ADR lifecycle、naming convention、5 筆實際 ADR | 架構決策有記錄，避免重複討論 |
| P1 | done | 建立 Anti-patterns | `anti-patterns/README.md` + 5 patterns | 已完成 5 個 anti-pattern 文件 | 失效模式可主動辨識與避免 |
| P1 | done | 升級 Skills Metadata v2 | `skills-index.yaml` | 已完成所有 13 skills 加入 weight/domains/dependencies/conflicts/priority.runtime | Skill relevance scoring 與 conflict detection 可運作 |
| P1 | done | 重建 intelligence/ 為專家智慧層 | `intelligence/README.md`, `intelligence/engineering/{architecture,domain,failure,heuristics,anti-patterns,tradeoffs,distributed-systems}/`, `intelligence/business/`, `intelligence/travel/` | 已完成 9 子目錄結構與 scope 定義，尚未填充實際 atoms | 每個子目錄有 README.md 定義核心、範例內容、與其他層的關係；與根 `anti-patterns/` 邊界已明確定義 |
| P1 | done | Phase 29：其餘 3 個 techniques decomposition | `analysis/apk/workflows/`, `intelligence/engineering/analytical-reasoning/{heuristics,signals}/` | 已完成 http-api、local-proxy、media-hls 拆解（3 workflows + 4 intelligence atoms） | 每個 technique 有 workflow + 1-2 intelligence atoms，舊檔案已標註 |
| P1 | done | Phase 30：Feedback history 提取 | `feedback/extraction/apk-analysis-index.md`, `feedback/extraction/development-guidance-index.md` | 已完成 101 lessons 索引建立，100% lessons 已對應到目標層 | Category index 已建立，101/101 lessons 已對應到目標層 |
| P2 | done | Phase 31：Intelligence Extraction Pipeline | `governance/lifecycle/intelligence-extraction-pipeline.md` | 已完成 7-step pipeline（內容審計 → 類型判斷 → 拆解執行 → 格式轉換 → 標註來源 → 驗證 → 更新索引） | Pipeline 文件已建立，已驗證 Technique Decomposition 和 Feedback History Extraction 兩種模式 |
| P2 | done | Phase 32：SKILL.md 分解 | `runtime/onboarding/apk-analysis-quickstart.md`, `runtime/onboarding/app-development-guidance-quickstart.md`, `runtime/onboarding/travel-planning-quickstart.md`, `workflow/travel-planning/artifact-gates.md` | 已完成 3 個 SKILL.md 瘦身（46-65% 減少），4 個新檔案建立 | SKILL.md 中所有內容已對應到目標層，舊檔案已標註 |
| P4 | done | Phase 33：Skill-Specific Intelligence Extraction（P4） | `plans/skill-specific-extraction.md`, `workflow/software-delivery/review-checklist.md`, `analysis/travel/sources-and-tools.md`, `analysis/travel/README.md` | 已完成 3 個 skills 的剩餘內容提取（CHECKLIST.md、TOOLS.md、README.md），所有 skills 內容已全部提取到新分層 | 每個 skill 有專屬 extraction strategy，內容結構分析 → 拆解 → 標註 → 驗證的完整流程已執行 |
| P2 | done | Knowledge Graph 補齊（Phase 28-33） | `knowledge/graphs/workflow-software-delivery.yaml`, `knowledge/graphs/intelligence-apk-analysis-atoms.yaml`, `knowledge/graphs/workflow-layers.yaml`, `knowledge/graphs/analysis-layers.yaml`, `knowledge/graphs/runtime-onboarding.yaml` | 已完成 5 個 graph records 更新，新增 1 個 graph record（runtime-onboarding） | Phase 28-33 所有新檔案在 graph records 中有對應 edge |
| P2 | done | Intelligence Extraction Pipeline 實戰驗證（P2） | `governance/lifecycle/intelligence-extraction-pipeline.md`, `analysis/apk/workflows/README.md`, `intelligence/engineering/analytical-reasoning/{heuristics,anti-patterns,failure,signals}/README.md` | 已完成 pipeline 狀態表更新，5 個 README 的 placeholder 已替換為實際 atom 表格 | Pipeline 的 3 種模式（Technique Decomposition、Feedback History Extraction、SKILL.md Decomposition）均已實際執行 |
| P2 | done | Skills Deprecation（Phase B） | `governance/lifecycle/README.md`, `skills/apk-analysis/techniques/`, `analysis/apk/techniques/` | 已完成 8 個 techniques 檔案從 `# Intelligence Extracted` 升級為 `# Deprecated — see <new path>` | Phase B 條件（techniques decomposition 完成 + pilot 驗證通過）已滿足 |
| P2 | done | Runtime surfaces 更新 | `knowledge/runtime/routing-registry.yaml`, `knowledge/indexes/README.md` | 已完成 routing-registry 新增 5 條路由（app-development-guidance、travel-planning、onboarding、apk-workflows、apk-intelligence-atoms），indexes 新增 3 條索引列 | 所有新層級在 routing-registry 和 indexes 中有對應路由 |
| P2 | done | 跨層一致性檢查 | `analysis/apk/workflows/README.md`, `intelligence/engineering/analytical-reasoning/{heuristics,anti-patterns,failure,signals}/README.md` | 已完成所有 README 的 atom/workflow 表格檢查，無遺漏 | 每個 README 的表格與實際檔案一致 |
| P2 | done | 跨 skill intelligence promotion（apk-analysis → heuristics） | `intelligence/engineering/heuristics/field-confidence-judgment.md`, `intelligence/engineering/heuristics/magic-bytes-reference.md` | 已完成 2 個 proven generalized atoms 從 apk-analysis 提升到跨領域 heuristics 層 | 原始 atoms 有 `# Cross-Domain Promotion` 標註，新 atoms 有 `# Source` 指向原始 domain atom |
| P2 | done | 跨 skill intelligence promotion（app-development-guidance → heuristics） | `intelligence/engineering/heuristics/document-priority-hierarchy.md` | 已完成 1 個 proven generalized atom 從 app-development-guidance 提升到跨領域 heuristics 層 | 原始 atom 有 `# Cross-Domain Promotion` 標註，新 atom 有 `# Source` 指向原始 domain atom |
| P2 | done | Skills Deprecation（Phase C） | `governance/lifecycle/README.md`, `skills/apk-analysis/techniques/`, `analysis/apk/techniques/` | 已完成 10 個舊 technique 檔案刪除（4 個子目錄 + 1 個 README 在 skills/，4 個 .md + 1 個 README 在 analysis/），更新 15+ 個引用舊路徑的檔案 | Phase C 檢查清單 7 項條件全部滿足，刪除後可 git revert rollback |
| P2 | done | Bootstrap flow 修正 — primary_entrypoint 優先 | `skills-index.yaml`, `CORE_BOOTSTRAP.md`, `shared-rules/dependency-reading.md` | 已完成 3 個 skill 新增 primary_entrypoint 指向新分層路徑，CORE_BOOTSTRAP.md 啟動流程加入 primary_entrypoint 優先規則，dependency-reading.md 加入 primary_entrypoint 優先規則 | AI 啟動時優先讀 primary_entrypoint 指向的新分層路徑，不再被舊 entrypoint 導向舊 skills/ |
| P1 | done | AI Decision Contract Testing 框架 | `validation/README.md`, `validation/scenarios/apk-analysis/`, `validation/rules/heuristics/`, `validation/traces/template.yaml`, `validation/evaluations/template.yaml`, `plans/ai-decision-contract-testing.md` | 已完成 validation/ 目錄建立（scenarios/rules/traces/evaluations），3 個 Scenario（flutter-aot-hooking、local-proxy-vs-pinning、early-hook-prevention），3 個 Rule 定義，Trace/Evaluation template | AI 的 decision path、rule obedience、routing stability 可透過 stateless scenario 驗證 |

## 最終目標

AI-native Knowledge Operating System 的最終目標不只是讓 AI 產生內容，而是建立：

- AI-native Engineering System。
- Knowledge Graph Runtime。
- Multi-model Orchestration。
- Engineering Intelligence Platform。
- Long-term AI Learning System。

未來真正瓶頸不會只是模型強度，而是知識是否能被正確管理、導航、組合與演化。

這是本 repository 下一階段的核心方向。
