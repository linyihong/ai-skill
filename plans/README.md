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
4. **每個計畫必須在檔頭標註狀態**：4-way enum — `draft` / `in-progress` / `completed (auto-detected)` / `completed (doc-only / pre-2026-strengthened)`。後者只給 2026-05-28 §`define_runtime_trigger_flow` 強化規則生效前已 archive 的 plan，受 [`governance/lifecycle/system-upgrade-governance.yaml`](../governance/lifecycle/system-upgrade-governance.yaml) §`pre_2026_05_28_doc_only_completion` grandfather flag 保護（primary sunset 2026-08-31，conditional 延展至 2026-11-30）；deadline 後須升 auto-detected 或降 orphan 下架。
5. **計畫完成後，若從中提煉出可重用的系統經驗，應建立對應的 intelligence atom**
6. **架構提案在 plan 內處理，不寫 proposed ADR**：依 [`governance/lifecycle/decision-promotion-pipeline.md`](../governance/lifecycle/decision-promotion-pipeline.md) §No-Proposed-ADR Rule，constitution/ 只放 accepted ADRs；提案、討論、alternatives 評估全在本 plan 的 `Decision Rationale` section 進行；plan completed 且通過 `ADR Promotion Criteria` 後才建立 accepted ADR。

## Plan Tree Hierarchy（主計畫／子計畫樹狀治理）

當主計畫執行中需要拆出有獨立 acceptance、跨多 phase、需獨立 sign-off 或可獨立 archive 的支線時，用 **plan tree** 表達 main ↔ sub 階層，而不是把 Phase 6/7/8 塞爆主計畫或另開孤立 plan 用 prose 連結。完整治理規則見 [`governance/lifecycle/plan-tree-hierarchy.md`](../governance/lifecycle/plan-tree-hierarchy.md)；落地範例見 [`active/2026-06-02-1200-plan-tree-hierarchy-governance/_plan.md`](active/2026-06-02-1200-plan-tree-hierarchy-governance/_plan.md)。

**核心模型（minimal governance）**：

- **Single source of truth = frontmatter `parent` pointer**。Hierarchy 由 `parent` 決定，不維護 `children:`（runtime scan 推導）。
- **Folder + `_plan.md` + `NN-` 前綴是 UI convention，不是 truth**。folder 放錯時 `ai-skill plans tree` 仍能建出正確樹；folder shape 只發 warning，不 block commit。
- **Lifecycle 與 storage 分離**：`status` 是生命週期、`plans/active|archived/` 是儲存位置；archive gate 只看 `status`，不看 location。

**Frontmatter（sub-plan 必填）**：`id` / `plan_kind: sub` / `status` / `owner` / `created` / `parent: <main-id>` / `required_for_completion: bool` / `sub_plan_reason: <非空 free text>`。Spike 用 `plan_kind: spike`（建議 `required_for_completion: false`）。詞彙定義見 [`knowledge/glossary/ai-skill.md`](../knowledge/glossary/ai-skill.md)（`plan_kind` / `parent` / `plan_tree` / `required_for_completion` / `sub_plan_reason`）。

**Commit-msg 機械強制（4 block + 1 warning）**：

| Validator | Severity | 規則 |
|---|---|---|
| `validatePlanTreeFrontmatter` | block | sub-plan 缺 `parent` / `sub_plan_reason`（空字串視為缺）/ `required_for_completion` |
| `validatePlanTreeArchiveOrder` | block | 主計畫 archive 時，所有 `parent == <main>` 且 `required_for_completion: true` 的 sub-plan 必須 `status: completed` |
| `validatePlanTreeParentReference` | block | `parent` 指向的 id 必須存在於 active + archived 全集（防 orphan node） |
| `validatePlanTreeUniqueID` | block | 全 repo plan `id` 唯一（防 parent pointer 指錯） |
| `validatePlanTreeFolderConvention` | warning | folder 缺 `_plan.md` / 檔名不符 `NN-` 前綴 / 深度 ≥ 3 |

**何時開 sub-plan**：`sub_plan_reason` 非空為唯一強制條件；recommended triggers（independent sign-off / multi-phase own acceptance / independent runtime trigger / parallel owners / independent spike archive）只是參考，不強制。單一 phase 內 step 細分用 checkbox、< 1 session 工作 inline 寫主計畫、純文件補強直接 commit，**不開 sub-plan**。

**檢視**：`ai-skill plans tree`（`--state active|archived|all`，`--format text|json|markdown`）純讀 frontmatter 動態建樹，顯示 status / 進度 / blocker。

## Plan 模板必填章節

任何涉及架構變更、新流程、跨層改動的 plan 必須包含下列章節（簡單修補類 plan 可省略 Decision Rationale）：

| 章節 | 必填條件 | 用途 |
|------|---------|------|
| `Status` 標頭 | 全部 plan | 標 `draft` / `in-progress` / `completed` |
| `Decision Rationale` | 架構/流程/跨層 plan | 取代「proposed ADR」的提案內容；含 Problem & Why Now / Decision / Alternatives Considered / Why Not an ADR Yet / ADR Promotion Criteria / Consequences |
| `Runtime Execution Path` | 涉及 framework/runtime/governance/workflow/validation/scenario/metadata/compiler 改動 | 明列 runtime owner、trigger flow、trigger location、activation contract、generated surface、validation scenarios、test passing evidence；trigger flow 必須說明 event → detector / route / query → loaded contract/source → runtime action / blocker → evidence，不能只寫「routing 會處理」。**新增明文 forbidden（2026-05-28）**：(a) 加 `route.*` entry 到 routing-registry 但沒同時 wire 一個 discovery signal 或 commit-msg validator → 不算 runtime integration；(b) 把資料 project 進 `runtime.db generated_surfaces` 或 SQLite index 但沒宣告 consumer（CLI / hook / validator / routable lookup）→ 不算 runtime integration。doc-only trial 必須明寫「不接入 runtime」+ 未來接入 phase / entry condition + graduation deadline，且 doc-only trial 本身**不能聲稱已完成 runtime integration** — 詳見 [`governance/lifecycle/system-upgrade-governance.md`](../governance/lifecycle/system-upgrade-governance.md) §3 規則 8 與 [`system-upgrade-governance.yaml`](../governance/lifecycle/system-upgrade-governance.yaml) §`define_runtime_trigger_flow` |
| `Deferred Runtime Projection`（外放→收斂宣告） | plan 建立 `runtime/*.yaml` 但**不立即 project 到 `runtime.db`** 時 | 預設規則：`runtime/` 下 YAML 必須含 `runtime_projection.enabled: true` 並 project；外放是例外。若 plan 採「外放→後續收斂」模式（doc-only trial、schema 演進中、multi-session evidence 需要），必須在 §Decision Rationale 或 §Runtime Execution Path 明寫 (a) **不 project 的 reason**、(b) **預定 project 的 phase / 條件**。沒寫明 → reviewer 必擋。Phase 6 of [`bootstrap-contract-yaml-migration`](archived/2026-05-25-2200-bootstrap-contract-yaml-migration.md) 會落地 `validateRuntimeYamlProjects` commit-msg validator 機械強制此規則。 |
| `Open Questions` | 架構/流程 plan | 列出 completed 前需釐清的問題 |
| `完成條件` | 全部 plan | 依 [`governance/lifecycle/system-upgrade-governance.md`](../governance/lifecycle/system-upgrade-governance.md) §2 checklist 子集 |
| `Phase 0` Pre-Build Interrogation | 架構/流程/跨層 plan | per §Architecture Compatibility Preflight。**Phase 0 必須以 §Phase 0 公版開頭片段（Open Questions 核對）貼上的 checklist 開頭**，確保盤點時逐條核對 Open Questions |
| `Phase 1-N` | 全部 plan | 實作步驟 + 各 phase 完成條件 |
| `Stakeholder 同意項目` | 架構/流程 plan | 列出 sign-off 項 |
| `Per-surface consumer 表` | 任何 plan 新增 `route.*` entry 或 `runtime_projection.target_key` 或 commit-msg validator | 列出每個新 surface 與其 named consumer（discovery signal / Go validator / routable lookup / `manual_activation` annotation）。表格欄位：`Generated surface key` / `Named consumer(s)` / `Consumer 類型`。**Reviewer 必擋條件**：surface 列出但無對應 consumer（除非顯式 `manual_activation: { reason: ... }`）。本表是 §`define_runtime_trigger_flow` forbidden rules 的 plan-level 對應；commit 時由 `validateRuntimeTriggerWiring` 機械驗證 |
| `Glossary Impact row` | 全部 plan | 一行（在 §Decision Rationale 或 §完成條件 旁）：`Glossary Impact: <yes / no>`。若 yes，列出新引入的 framework vocabulary terms + 是否已在 `knowledge/glossary/ai-skill.md` 註冊。若 no，明寫 `no new framework vocabulary introduced`。Audit-time `ai-skill runtime audit` 的 glossary coverage warning 會 cross-check 此宣告 |
| `Watch-Out List citation` | 跨層 / Gen 4 forward plan | 引用對應 Gen 4 vision §Watch-Out List 的 wall（防 scope drift、防 over-engineering）。Path: [`architecture/ai-native-cognitive-ecosystem-system.md`](../architecture/ai-native-cognitive-ecosystem-system.md) §Watch-Out List |
| `與其他 plans 的關係` | 全部 plan | Cross-reference |

### `Decision Rationale` section 內容規範

當 plan 含 Decision Rationale 時，子章節須涵蓋：

```markdown
## Decision Rationale

### Problem & Why Now
（現狀問題 + 為什麼這個時間點要動）

### Decision
（具體決定做什麼）

### Alternatives Considered
- A. 維持現狀：reject because ...
- B. 完全重寫 X：reject because ...
- C. 漸進改造（accept）

### Why Not an ADR Yet
（為什麼這個階段不適合寫 ADR：未驗證、scope 還會調、open questions 未解、可能有更輕 promotion target）

### ADR Promotion Criteria（completed 時驗證）
- [ ] foundational + cross-session + cross-project + expensive-to-reverse + explains-why 全中
- [ ] Plan 結果證實 decision 可行
- [ ] Open Questions 全解
- [ ] 沒有更輕的 promotion target 適用（per ADR-007）
- [ ] 系統真實使用此 contract（具體 evidence 指標）

### Consequences（預期）
#### 正面
#### 負面
#### 風險
```

### Phase 0 公版開頭片段（Open Questions 核對）

每個 plan 的 Phase 0 **開頭**貼上下列公版 checklist，把「盤點順手解掉 Open Question 卻忘了回寫」變成結構性步驟。完成盤點後逐條回填，並把已解項目同步勾選 / 附註於 §Open Questions：

```markdown
### Phase 0.0 — Open Questions 核對（公版，必填）

逐條核對本 plan §Open Questions，標記處置並回寫：

- [ ] 已讀本 plan §Open Questions 全部條目
- [ ] 對每條標記 `resolved`（附 Phase 0 證據）/ `still-open` / `deferred`（附原因）
- [ ] `resolved` 的條目已同步勾選 / 附註於 §Open Questions
- [ ] 若盤點新發現問題，已加入 §Open Questions

| Open Question | 處置 | 證據 / 原因 |
|---|---|---|
| <Q1 摘要> | resolved / still-open / deferred | <Phase 0 盤點證據或延後原因> |
```

## Plan 執行前架構相容性檢查（Architecture Compatibility Preflight）

開始執行任何 `active/` plan 前，agent **必須**先確認 plan 與現行架構相容。此檢查是 blocking gate；未完成前不得進入 implementation phase。

在架構相容性檢查前，若 plan 會導向 code、workflow、governance、runtime、validation、schema、generated artifact 或 tool adapter 改動，必須先完成 [`workflow/software-delivery/requirements/pre-build-interrogation.md`](../workflow/software-delivery/requirements/pre-build-interrogation.md)。此 gate 用來確認 goal、scope、non-goals、acceptance、validation target、framework source-of-truth、duplication risk 與 blocker questions，避免 stale plan 直接進入實作。

### 檢查清單

| # | 檢查項目 | 說明 |
|---|---------|------|
| 1 | **Candidate files 存在性** | plan 列出的 source、generated surface、runtime table、workflow / metadata path 是否仍存在；缺檔需標 `not applicable` 或 `source missing` |
| 2 | **Source-of-truth 一致性** | 確認應修改的是 canonical source、SQLite canonical document、embedded source、compiler source 或 generated DB；不得只改不生效的 mirror / generated output |
| 3 | **Layer responsibility** | plan 是否把 policy、runtime state、workflow、metadata、analysis、intelligence 放在正確 layer |
| 4 | **Compiler / generated surface** | 涉及 `runtime/`、`knowledge/`、`metadata/`、`validation/` 時，確認 compiler / validator 會讀到該 source，並列出需要重新生成的 artifact |
| 5 | **Pre-build interrogation** | 若 plan 來自模糊需求或 framework 改動，確認已記錄需求拷問、source-of-truth discovery、duplication risk、open questions 與 assumptions |
| 6 | **Linked updates** | 依 [`enforcement/linked-updates.md`](../enforcement/linked-updates.md) 確認相關 README、metadata、contract inventory、routing registry、templates、runtime DB 或 validators 是否要同步 |
| 7 | **Open Questions 核對** | 逐條核對本 plan §Open Questions：Phase 0 / preflight 的盤點結果是否已回答、否定或細化任一 Open Question。對每條標記 `resolved`（附 Phase 0 證據）/ `still-open` / `deferred`（附原因），並把已解項目同步勾選或附註於 §Open Questions。**不得**只在工作筆記回答 Open Question 卻不回寫 plan |
| 8 | **Execution decision** | 若發現架構衝突、未解 blocker question 或 source-of-truth duplication risk，先暫停執行並更新 plan / 詢問使用者；不得邊實作邊假設 plan 仍正確 |

### 最低記錄格式

每次 preflight 至少要在工作筆記、plan Phase 0、或回覆中留下：

| 欄位 | 必填內容 |
| --- | --- |
| Trigger | 要開始執行哪個 plan / phase |
| Checked sources | 讀過哪些 current architecture sources |
| Conflicts | 無衝突，或列出 candidate path / source-of-truth / compiler / layer 衝突 |
| Interrogation | goal、scope、non-goals、acceptance、framework discovery、duplication risk、open questions / assumptions |
| Open Questions 核對 | 逐條列出本 plan 每個 Open Question 的處置：`resolved`（附 Phase 0 證據）/ `still-open` / `deferred`（附原因）；已解項目須回寫 §Open Questions |
| Decision | proceed / revise plan first / ask user / blocked |
| Validation | 用什麼方式確認（diff、runtime query、validator、link check、readback） |

### 強制執行規則

1. **任何 active plan 的 Phase 1 或 implementation phase 開始前，都必須先完成 Pre-build Interrogation 與 Architecture Compatibility Preflight。**
2. 若 plan 已有 Phase 0，Phase 0 必須包含此檢查；若沒有，agent 必須先補做 preflight，再決定是否需要更新 plan。
3. 若 preflight 發現 plan 與 current architecture 衝突、blocking question 未解，或會產生雙份 source-of-truth，必須先修正 plan 或取得使用者確認，不得直接繼續執行。
4. 涉及 `runtime.db`、generated reports、SQLite index 或 compiler outputs 時，preflight 必須確認「source 變更是否真的進入 generated surface」，以及舊 duplicate surface 是否已刪除、deprecate 或明確降級。
5. **Phase 0 必須以 Open Questions 核對 checklist 開頭**（見 §Plan 模板必填章節 `Phase 0` 列的公版片段）。Phase 0 的盤點若回答了任一 Open Question，必須在同一輪把該 Open Question 標記 `resolved` 並回寫 plan，不得讓盤點結果與 §Open Questions 狀態脫節。此公版片段是「盤點順手解掉 Open Question 卻忘了回寫」失效模式的結構性防呆。

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
| [`archived/2026-05-25-2100-runtime-cognitive-contract-v2.md`](archived/2026-05-25-2100-runtime-cognitive-contract-v2.md) | ✅ completed | Runtime Cognitive Contract v2：ADR-008 amendment；新增 validation_mode / derived cognitive_cost、compact/full adaptive disclosure、activation signal enforcement、high-risk capability snippet、inflated-reporting failure pattern 與 commit-msg validators。 |
| [`archived/2026-05-28-1200-gen3-runtime-trigger-audit-and-completion.md`](archived/2026-05-28-1200-gen3-runtime-trigger-audit-and-completion.md) | ✅ completed (auto-detected) | Gen 3 Runtime Trigger Audit & Completion — Phase 0–7 全部達成。落地 `ai-skill runtime audit` subcommand（md + `--json` 雙渲染）、3 個新 commit-msg validators（`validatePlanCheckboxSync` 第 16 / `validateRuntimeTriggerWiring` 第 17 / `validateEvidenceHierarchy` 第 18）、`pre_2026_05_28_doc_only_completion` grandfather flag（4 plans covered，sunset 2026-08-31 + conditional 2026-11-30）、warning-only glossary coverage guardrail（7 paths × backtick/snake_case heuristic）、plan template 4 個新 required sections。Audit baseline 242→237 orphan（5 wires 含 `enforcement.evidence_hierarchy.contract` / cognitive-state-evidence / memory-retrieval-activation / model-aware-routing / runtime-cognitive-modes）。Follow-up: [`active/2026-05-28-1830-plan-archival-audit-validator.md`](active/2026-05-28-1830-plan-archival-audit-validator.md)。 |
| [`active/2026-05-27-1557-tool-runtime-signal-economics-integration.md`](active/2026-05-27-1557-tool-runtime-signal-economics-integration.md) | draft | Tool Runtime Signal & Economics Integration：規劃把 `tools/` 從 document / routing index layer 升級為 runtime-readable signal source，並補上 execution economics layer 作為 Cognitive Mode discovery 的輸入訊號。**Sequencing：等 audit plan Phase 5 graduate（validateRuntimeTriggerWiring active）後再啟動**，才能享受新增 11 surfaces 自動保護。 |
| [`active/2026-05-28-1636-gen4-fitness-optimization-memory-interface-reservation.md`](active/2026-05-28-1636-gen4-fitness-optimization-memory-interface-reservation.md) | draft | Gen4 Fitness & Optimization Memory Interface Reservation：預留 positive optimization memory、rejected optimization memory、activation fitness 與 fitness placeholder schema；明確不做 autonomous optimizer / self-modifying governance / full telemetry DB。**Sequencing：排在 Gen3 audit 後，並以 economics plan 的 telemetry/economics primitives 作為 future input**。 |
| [`active/2026-05-28-1830-plan-archival-audit-validator.md`](active/2026-05-28-1830-plan-archival-audit-validator.md) | draft | Plan Archival Audit Validator — Gen3 audit plan §Phase 7 follow-up：將「archive 時所有 `- [ ]` 必須翻 `[x]` 或在 body 明文交代」這條 manual 規則機械化成第 19 個 commit-msg validator `validatePlanArchivalAudit`（block default；opt-out `[skip-plan-archival-audit]`）。3 scenarios + ≥ 5 fixture tests + dogfood archive。**Sequencing：可獨立排程；建議在 Gen3 audit plan 進 archived 之前 graduate 以保護該 archive commit**。 |
| [`active/2026-06-06-1800-sanitization-mechanical-enforcement.md`](active/2026-06-06-1800-sanitization-mechanical-enforcement.md) | in-progress | Sanitization Mechanical Enforcement：supersede `2026-05-31-2000` allowlist route，採 metadata-derived forbidden tokens + shared-layer topology + staged-content scanner，作為 `rule_classes[sanitization]` canonical executor；Phase 0 已完成 sibling supersede / parent reference sync，Phase 1 implementation pending。 |
| [`active/2026-06-09-1040-experience-validation-pipeline-evolution.md`](active/2026-06-09-1040-experience-validation-pipeline-evolution.md) | in-progress | Experience Validation Pipeline Evolution：記錄 responsive render-context governance 後續 open questions；Phase 1–3 已落地 Browser Evidence Metadata / Capture Envelope、Validation Coverage Model watch-list、Responsive domain downgrade criteria；Phase 4 Evidence Envelope 初步 spike 決定暫不 promotion，typed Context Taxonomy 仍等待 scenario pressure，不新增 runtime gate。 |
| [`archived/2026-06-10-0908-user-journey-validation-integration.md`](archived/2026-06-10-0908-user-journey-validation-integration.md) | ✅ completed (auto-detected) | User Journey Validation Integration：Phase 0–5 與 Vidoe-Test pilot 已完成；BDD 擁有 Journey Specification，validation workflow 負責 Journey Execution；criticality selection criteria、`validation_scope`、`expected_outcomes` / `observable_evidence`、workflow gates、artifact gates、glossary terms 與 4 個 validation scenarios 已落地。Pilot 確認 API success 不足以證明 journey pass，`membership_active` 需 DB readback，`playback_allowed` 需 protected resource readback。 |
| [`archived/2026-06-08-1544-evidence-acquisition-layer.md`](archived/2026-06-08-1544-evidence-acquisition-layer.md) | ✅ completed (auto-detected) | Evidence Acquisition Layer（Phase 1 of Validation Reasoning Taxonomy）：補上 `collection_method` 作為 evidence acquisition layer，讓 Browser Review、contract readback、static analysis、runtime trace 與 human observation 有明確位置；first landing 限於 UI-local workflow taxonomy + validation scenarios，並記錄未來 ownership graduation 到 shared validation-reasoning / finding taxonomy，不新增 `sd-browser-review`、runtime YAML 或 enforcement rule_class。 |
| [`archived/2026-06-08-1047-feedback-learning-report-obligation.md`](archived/2026-06-08-1047-feedback-learning-report-obligation.md) | ✅ completed (auto-detected) | Feedback / Learning Report Obligation：新增 final close-out learning disposition report，拆成 `feedback_decision` / `repo_context` / `writeback_status` 三維；runtime contract、stop hook schema validator、Cursor/Claude adapters、routing rules、validation scenarios 與 glossary entries 已落地。ADR promotion deferred until post-use spam check completes。 |
| [`archived/2026-06-08-1408-ui-governance-workflow.md`](archived/2026-06-08-1408-ui-governance-workflow.md) | ✅ completed (auto-detected) | UI Governance Workflow Integration：新增 `sd-ui-governance` workflow slice、software-delivery loading surface、artifact gates、focused evidence template、review checklist、advisory runtime-lite candidate signals 與 6 個 validation scenarios；first landing 保持 doc + scenario only，未 promotion 成 mechanical rule_class。 |
| [`archived/2026-05-26-1039-landing-page-positioning-refresh.md`](archived/2026-05-26-1039-landing-page-positioning-refresh.md) | ✅ completed | AI-native Cognitive Execution System Landing Page Refresh：根 README 已重構為 public-facing landing page，完整 overview 收斂到 `architecture/ai-native-cognitive-execution-system.md`；正式名稱使用 AI-native Cognitive Execution System，`Ai-skill` 僅作為尚未改名的 repo slug |
| [`archived/2026-05-25-1000-context-language-glossary-system.md`](archived/2026-05-25-1000-context-language-glossary-system.md) | ✅ completed | Context Language Glossary System：Phase 0–7 全部完成。`knowledge/glossary/README.md` schema spec + `ai-skill glossary validate` Go validator + 19 framework entries + SQLite projection（3 表）+ `route.knowledge.glossary` + `file_diff_glossary_touched` / `user_keyword_term_conflict` discovery signals + `validateGlossaryRetroOwn` commit-msg validator（gate.glossary.retro_own_required）全部 live。Gen 4 ecosystem-adaptation candidate terms 同步預收。 |
| [`archived/2026-05-22-1629-runtime-cognitive-modes-system.md`](archived/2026-05-22-1629-runtime-cognitive-modes-system.md) | ⚠️ completed (doc-only / pre-2026-strengthened) | Runtime Cognitive Modes System：把 `models/` 從 documentation layer 提升為 runtime activation；引入 4 維 cognitive mode primitive。`route.runtime.cognitive-modes` 已註冊但無 discovery signal 拉它（plan 主要被 commit-msg validators 直接消費，route 本身為 manual_activation candidate）。受 grandfather flag 保護；Phase 4 將決定升 manual_activation 或補 signal。依賴 [ADR-008](../constitution/ADR-008-runtime-cognitive-modes.md)。 |
| [`archived/2026-05-22-0855-executable-yaml-contract-migration.md`](archived/2026-05-22-0855-executable-yaml-contract-migration.md) | ✅ completed | Executable YAML Contract Migration：盤點哪些流程、gate、required reads、failure actions 應升級為 owner-layer YAML contract，並投影到 runtime generated surfaces，降低非 ChatGPT agent 漏跑流程的風險 |
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
| [`archived/2026-05-15-0920-runtime-execution-layer-upgrade-analysis.md`](archived/2026-05-15-0920-runtime-execution-layer-upgrade-analysis.md) | ✅ completed / archived | AI-native Cognitive Execution System 升級比對分析已完成；P0/P1/P2 execution runtime 缺口已由 `runtime/runtime.db`、SQLite canonical runtime documents、recovery、output governance、distributed runtime 與 cognitive governance plan 吸收，Agent VM 留作遠期方向 |
| [`archived/2026-05-15-0949-workflow-activation-contract-migration.md`](archived/2026-05-15-0949-workflow-activation-contract-migration.md) | ✅ superseded / archived | Per-workflow `activation-contract.yaml` 方案已被 ADR-006 registry-first workflow activation 取代；現行 source 是 activation #27、`route.workflow.*.activation_triggers` 與 `workflow/workflow-routing.md` |
| [`archived/2026-05-20-1039-runtime-recovery-escalation-system.md`](archived/2026-05-20-1039-runtime-recovery-escalation-system.md) | ✅ completed | Runtime Recovery & Escalation System — escalation policy、runtime guard、recovery procedure、metadata policy、workflow hooks 與 validation scenarios 全完成 |
| [`archived/2026-05-20-1307-ai-runtime-governance-five-step-integration.md`](archived/2026-05-20-1307-ai-runtime-governance-five-step-integration.md) | ✅ completed | AI Runtime Governance Five-Step Integration — Musk Five-Step source philosophy 與 AI runtime governance 轉譯層已完成 |
| [`archived/2026-05-20-1501-cognitive-state-evidence-governance.md`](archived/2026-05-20-1501-cognitive-state-evidence-governance.md) | ⚠️ completed (doc-only / pre-2026-strengthened) | Cognitive State & Evidence Governance — `route.governance.cognitive-state-evidence` 與 `enforcement.evidence_hierarchy.contract` 已 project 但無 discovery signal / commit-msg validator 消費；受 grandfather flag 保護，Phase 4 將補 wire `validateEvidenceHierarchy`。詳見 [`governance/lifecycle/system-upgrade-governance.yaml`](../governance/lifecycle/system-upgrade-governance.yaml) §`pre_2026_05_28_doc_only_completion`。 |
| [`archived/2026-05-20-1745-memory-retrieval-activation-governance.md`](archived/2026-05-20-1745-memory-retrieval-activation-governance.md) | ⚠️ completed (doc-only / pre-2026-strengthened) | Memory Retrieval & Activation Governance — `route.memory.retrieval-activation` 已註冊但 audit 仍判 orphan（無 discovery signal pull）；受 grandfather flag 保護。Phase 4 候選補 wire 對象。 |
| [`archived/2026-05-20-1802-model-aware-execution-routing.md`](archived/2026-05-20-1802-model-aware-execution-routing.md) | ⚠️ completed (doc-only / pre-2026-strengthened) | Model-Aware Execution Routing — `route.models.model-aware-routing` 與相關 generated surfaces 已 project 但無 commit-msg validator 引用；受 grandfather flag 保護。 |
| [`archived/2026-05-21-0834-cross-platform-go-script-runtime.md`](archived/2026-05-21-0834-cross-platform-go-script-runtime.md) | ✅ completed / archived | Cross-Platform Go Script Runtime — Windows、macOS、Linux repo-local binaries、native runtime refresh/validate/compile/query、CI artifacts、binary guards、mobile out-of-scope decision、legacy script disposition 已完成；持續生效 policy 轉由 `scripts/ai-skill-cli/docs/` 維護 |
| [`archived/2026-05-20-1601-ddd-intelligence-software-delivery-governance.md`](archived/2026-05-20-1601-ddd-intelligence-software-delivery-governance.md) | ✅ completed / archived | DDD Integration Plan — DDD domain intelligence、architecture selection、software-delivery architecture governance、metadata heuristics、validation scenarios、routing registry 與 generated runtime surfaces 已完成；DDD 維持 selectable architecture strategy，不 promotion 成 runtime invariant |
| [`archived/2026-05-20-1635-bdd-ddd-cognition-aligned-reframe.md`](archived/2026-05-20-1635-bdd-ddd-cognition-aligned-reframe.md) | ✅ completed / archived | BDD + DDD Cognition-Aligned Reframe — BDD 歸入 requirements cognition，DDD 歸入 domain architecture cognition，workflow 拆成 delivery stages，runtime 僅接收 metadata-only runtime-lite signal；routing、graphs、metadata、validation 與 generated runtime surfaces 已更新 |

## 誰會參考這裡（Inbound References）

- [`route.governance.durable-goal-boundary`](../knowledge/runtime/routing-registry.yaml) — candidate_sources 引用 `scripts/README.md`
- [`enforcement/conversation-goal-ledger.md`](../enforcement/conversation-goal-ledger.md) — 定義 active goal 與 durable planning 的邊界
- [`enforcement/linked-updates.md`](../enforcement/linked-updates.md) — 計畫完成後需執行連動更新檢查

## 與其他層的關係

- [`plans/archived/2026-05-11-1112-next-stage-upgrade-plan.md`](archived/2026-05-11-1112-next-stage-upgrade-plan.md) — 已完成的全局升級路線圖（所有 Phase 1-33 已執行完畢）
- [`governance/lifecycle/README.md`](../governance/lifecycle/README.md) — Skills Deprecation Timeline 等生命週期規則
- [`intelligence/engineering/agent-architecture/`](../intelligence/engineering/agent-architecture/) — 從已完成計畫中提煉的系統經驗結晶
