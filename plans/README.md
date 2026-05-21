# Plans（計畫目錄）

## 目錄規則

| 子目錄 | 用途 | 生命週期 |
|--------|------|---------|
| [`active/`](active/) | 進行中或待審閱的計畫（draft / in-progress） | 完成後搬移至 `archived/` |
| [`archived/`](archived/) | 已執行完成的計畫（執行結果記錄） | 永久保留，作為決策記錄 |

## 原則

1. **`active/` 只放尚未開始或正在執行的計畫** — 一旦計畫執行完畢，立即搬移至 `archived/`
2. **`archived/` 的計畫不刪除** — 作為歷史決策記錄，可供日後查閱
3. **計畫檔案命名規則**：`YYYY-MM-DD-HHMM-<slug>.md`，日期時間前綴讓檔案按時間排序（精確到分鐘），slug 需能反映計畫核心目標
4. **每個計畫必須在檔頭標註狀態**：`draft` / `in-progress` / `completed`
5. **計畫完成後，若從中提煉出可重用的系統經驗，應建立對應的 intelligence atom**

## Plan 執行前架構相容性檢查（Architecture Compatibility Preflight）

開始執行任何 `active/` plan 前，agent **必須**先確認 plan 與現行架構相容。此檢查是 blocking gate；未完成前不得進入 implementation phase。

### 檢查清單

| # | 檢查項目 | 說明 |
|---|---------|------|
| 1 | **Candidate files 存在性** | plan 列出的 source、generated surface、runtime table、workflow / metadata path 是否仍存在；缺檔需標 `not applicable` 或 `source missing` |
| 2 | **Source-of-truth 一致性** | 確認應修改的是 canonical source、YAML source、embedded source、compiler source 或 generated DB；不得只改不生效的 mirror / generated output |
| 3 | **Layer responsibility** | plan 是否把 policy、runtime state、workflow、metadata、analysis、intelligence 放在正確 layer |
| 4 | **Compiler / generated surface** | 涉及 `runtime/`、`knowledge/`、`metadata/`、`validation/` 時，確認 compiler / validator 會讀到該 source，並列出需要重新生成的 artifact |
| 5 | **Linked updates** | 依 [`enforcement/linked-updates.md`](../enforcement/linked-updates.md) 確認相關 README、metadata、activation rules、templates、runtime DB 或 validators 是否要同步 |
| 6 | **Execution decision** | 若發現架構衝突，先暫停執行並更新 plan / 詢問使用者；不得邊實作邊假設 plan 仍正確 |

### 最低記錄格式

每次 preflight 至少要在工作筆記、plan Phase 0、或回覆中留下：

| 欄位 | 必填內容 |
| --- | --- |
| Trigger | 要開始執行哪個 plan / phase |
| Checked sources | 讀過哪些 current architecture sources |
| Conflicts | 無衝突，或列出 candidate path / source-of-truth / compiler / layer 衝突 |
| Decision | proceed / revise plan first / ask user / blocked |
| Validation | 用什麼方式確認（diff、runtime query、validator、link check、readback） |

### 強制執行規則

1. **任何 active plan 的 Phase 1 或 implementation phase 開始前，都必須先完成 Architecture Compatibility Preflight。**
2. 若 plan 已有 Phase 0，Phase 0 必須包含此檢查；若沒有，agent 必須先補做 preflight，再決定是否需要更新 plan。
3. 若 preflight 發現 plan 與 current architecture 衝突，必須先修正 plan 或取得使用者確認，不得直接繼續執行。
4. 涉及 `runtime.db`、generated reports、SQLite index 或 compiler outputs 時，preflight 必須確認「source 變更是否真的進入 generated surface」。

## Plan 完成閉環（Plan Completion Closure）

當一個 plan 的所有項目都標記為完成（`✅`）時，agent **必須**執行以下閉環檢查：

### 檢查清單

| # | 檢查項目 | 說明 |
|---|---------|------|
| 1 | **確認所有項目已完成** | 檢查 plan 中所有 task 是否都標記為 `✅`，無遺漏項目 |
| 2 | **執行 validator** | 若 plan 涉及 `knowledge/`、`validation/`、`intelligence/` 等層，執行 `ai-skill runtime refresh` |
| 3 | **檢查連動更新** | 依 [`enforcement/linked-updates.md`](../enforcement/linked-updates.md) 檢查 plan 改動是否需要同步其他檔案 |
| 4 | **更新 plans/README.md 狀態** | 將本 plan 在[目前狀態](#目前狀態)表格中的狀態改為 `✅ completed` |
| 5 | **搬移至 archived/** | 將 plan 檔案從 `active/` 搬移至 `archived/`，檔名與內容不變 |
| 6 | **Commit & push** | 提交搬移與狀態更新，並推送 |
| 7 | **最終確認** | 執行 `git status --short --branch` 確認工作樹乾淨 |

### 強制執行規則

1. **最後一個 Phase 完成後，agent 必須立即執行閉環檢查清單**，不得直接結束或進行 commit & push。
2. 若 plan 有多個 Phase，最後一個 Phase 的完成條件中必須包含「執行 Plan Completion Closure」。
3. 違反此規則的 commit 應被視為閉環不完整，需依 [`enforcement/linked-updates.md`](../enforcement/linked-updates.md) 的「閉環不完整時的強制補救」處理。

### 不搬移的例外情況

若 plan 符合以下任一條件，可留在 `active/` 但標註 `✅ completed`：

- Plan 是**持續生效的基礎建設**（如 validation gate、pre-commit hook），未來可能擴充新 Phase
- Plan 的 scope 是 ongoing 的維護性任務，沒有明確的「完成」邊界

例外情況必須在 plan 檔頭或 `plans/README.md` 表格中說明原因。

## 目前狀態

| 檔案 | 狀態 | 說明 |
|------|------|------|
| [`archived/2026-05-11-1112-next-stage-upgrade-plan.md`](archived/2026-05-11-1112-next-stage-upgrade-plan.md) | ✅ completed | 全局升級路線圖（所有 Phase 1-33 已執行完畢） |
| [`archived/2026-05-11-1129-apk-analysis-pilot-migration.md`](archived/2026-05-11-1129-apk-analysis-pilot-migration.md) | ✅ completed | APK Analysis Pilot Migration 狀態圖（原 architecture/） |
| [`archived/2026-05-12-1101-context-cost-optimization.md`](archived/2026-05-12-1101-context-cost-optimization.md) | ✅ completed | Phase 1：Context Cost Optimization 執行計畫（原 architecture/） |
| [`archived/2026-05-12-1458-technique-intelligence-pilot.md`](archived/2026-05-12-1458-technique-intelligence-pilot.md) | ✅ completed | Phase 28：Technique → Intelligence Pilot（flutter-dart-aot） |
| [`archived/2026-05-12-1506-skill-specific-extraction.md`](archived/2026-05-12-1506-skill-specific-extraction.md) | ✅ completed | Phase 33：Skill-Specific Intelligence Extraction |
| [`archived/2026-05-13-0954-cognitive-boundary-system.md`](archived/2026-05-13-0954-cognitive-boundary-system.md) | ✅ completed | Cognitive Boundary System 整合計畫，所有 Phase 1-8 已執行完畢 |
| [`archived/2026-05-13-1331-knowledge-runtime-validation-gate.md`](archived/2026-05-13-1331-knowledge-runtime-validation-gate.md) | ✅ completed | Part 1: Validation Gate 已完成；Part 2: UI Operation Intelligence Extraction 已完成 |
| [`archived/2026-05-13-0837-ai-decision-contract-testing.md`](archived/2026-05-13-0837-ai-decision-contract-testing.md) | ✅ completed | AI Decision Contract Testing 框架設計與實作 |
| [`archived/2026-05-14-1035-enforcement-layer-enhancement.md`](archived/2026-05-14-1035-enforcement-layer-enhancement.md) | ✅ completed | enforcement/ 後續強化計畫：Metadata Spec、Rule Graph、Activation Engine、Conflict Matrix、Deprecation Lifecycle（5 方向全完成） |
| [`archived/2026-05-14-1028-shared-rules-to-enforcement-migration.md`](archived/2026-05-14-1028-shared-rules-to-enforcement-migration.md) | ✅ completed | shared-rules/ → enforcement/ 搬遷計畫，含 Layer Responsibility Contract |
| [`archived/2026-05-18-scrapling-knowledge-integration-plan.md`](archived/2026-05-18-scrapling-knowledge-integration-plan.md) | ✅ completed | Scrapling 知識整合計畫：analysis/web/ + 6 份 intelligence 文件 + sanitization 強化 + routing 註冊，3 個 Phase 全完成 |
| [`archived/2026-05-18-0155-software-delivery-output-templates.md`](archived/2026-05-18-0155-software-delivery-output-templates.md) | ✅ completed | Software Delivery Output Templates — 建立 5 個輸出模板 + Greenfield 標準化流程 + Slash Command 模式 + 模板 Traceability 整合 |
| [`archived/2026-05-15-0920-runtime-execution-layer-upgrade-analysis.md`](archived/2026-05-15-0920-runtime-execution-layer-upgrade-analysis.md) | ✅ completed / archived | AI-native Cognitive Execution System 升級比對分析已完成；P0/P1/P2 execution runtime 缺口已由 `runtime/runtime.db`、`runtime/compiler/embedded_data.rb`、recovery、output governance、distributed runtime 與 cognitive governance plan 吸收，Agent VM 留作遠期方向 |
| [`archived/2026-05-15-0949-workflow-activation-contract-migration.md`](archived/2026-05-15-0949-workflow-activation-contract-migration.md) | ✅ superseded / archived | Per-workflow `activation-contract.yaml` 方案已被 ADR-006 registry-first workflow activation 取代；現行 source 是 activation #27、`route.workflow.*.activation_triggers` 與 `workflow/workflow-routing.md` |
| [`archived/2026-05-20-1039-runtime-recovery-escalation-system.md`](archived/2026-05-20-1039-runtime-recovery-escalation-system.md) | ✅ completed | Runtime Recovery & Escalation System — escalation policy、runtime guard、recovery procedure、metadata policy、workflow hooks 與 validation scenarios 全完成 |
| [`archived/2026-05-20-1307-ai-runtime-governance-five-step-integration.md`](archived/2026-05-20-1307-ai-runtime-governance-five-step-integration.md) | ✅ completed | AI Runtime Governance Five-Step Integration — Musk Five-Step source philosophy 與 AI runtime governance 轉譯層已完成 |
| [`archived/2026-05-20-1501-cognitive-state-evidence-governance.md`](archived/2026-05-20-1501-cognitive-state-evidence-governance.md) | ✅ completed / archived | Cognitive State & Evidence Governance — governance translation、evidence hierarchy enforcement、runtime-lite compressed guard、metadata evidence policy、validation scenarios 與 generated runtime surfaces 已完成 |
| [`active/2026-05-20-1745-memory-retrieval-activation-governance.md`](active/2026-05-20-1745-memory-retrieval-activation-governance.md) | 📝 draft | Memory Retrieval & Activation Governance — 將 memory/ 從 storage taxonomy 升級為 selective cognitive replay system，補足 retrieval、activation、replay cost、freshness、contamination boundary、working-memory buffer 與 promotion pipeline |
| [`active/2026-05-20-1802-model-aware-execution-routing.md`](active/2026-05-20-1802-model-aware-execution-routing.md) | 📝 draft | Model-Aware Execution Routing — 將 models/ 從 profile / compression documentation 升級為 execution strategy layer，定義 task complexity、cognitive state、autonomy mode、context budget 與 tool capability 的 routing contract |
| [`archived/2026-05-21-0834-cross-platform-go-script-runtime.md`](archived/2026-05-21-0834-cross-platform-go-script-runtime.md) | ✅ completed / archived | Cross-Platform Go Script Runtime — Windows、macOS、Linux repo-local binaries、native runtime refresh/validate/compile/query、CI artifacts、binary guards、mobile out-of-scope decision、legacy script disposition 已完成；持續生效 policy 轉由 `scripts/ai-skill-cli/docs/` 維護 |
| [`archived/2026-05-20-1601-ddd-intelligence-software-delivery-governance.md`](archived/2026-05-20-1601-ddd-intelligence-software-delivery-governance.md) | ✅ completed / archived | DDD Integration Plan — DDD domain intelligence、architecture selection、software-delivery architecture governance、metadata heuristics、validation scenarios、routing registry 與 generated runtime surfaces 已完成；DDD 維持 selectable architecture strategy，不 promotion 成 runtime invariant |
| [`archived/2026-05-20-1635-bdd-ddd-cognition-aligned-reframe.md`](archived/2026-05-20-1635-bdd-ddd-cognition-aligned-reframe.md) | ✅ completed / archived | BDD + DDD Cognition-Aligned Reframe — BDD 歸入 requirements cognition，DDD 歸入 domain architecture cognition，workflow 拆成 delivery stages，runtime 僅接收 metadata-only runtime-lite signal；routing、graphs、metadata、validation 與 generated runtime surfaces 已更新 |

## 誰會參考這裡（Inbound References）

- [`route.governance.durable-goal-boundary`](../knowledge/runtime/routing-registry.yaml:129) — candidate_sources 引用 `scripts/README.md`
- [`enforcement/conversation-goal-ledger.md`](../enforcement/conversation-goal-ledger.md) — 定義 active goal 與 durable planning 的邊界
- [`enforcement/linked-updates.md`](../enforcement/linked-updates.md) — 計畫完成後需執行連動更新檢查

## 與其他層的關係

- [`plans/archived/2026-05-11-1112-next-stage-upgrade-plan.md`](archived/2026-05-11-1112-next-stage-upgrade-plan.md) — 已完成的全局升級路線圖（所有 Phase 1-33 已執行完畢）
- [`governance/lifecycle/README.md`](../governance/lifecycle/README.md) — Skills Deprecation Timeline 等生命週期規則
- [`intelligence/engineering/agent-architecture/`](../intelligence/engineering/agent-architecture/) — 從已完成計畫中提煉的系統經驗結晶
