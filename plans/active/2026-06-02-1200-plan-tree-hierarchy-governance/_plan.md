---
id: 2026-06-02-1200-plan-tree-hierarchy-governance
plan_kind: main
status: draft
owner: linyihong
created: 2026-06-02
parent: null
children: []
---

# Plan Tree Hierarchy Governance（主計畫／子計畫樹狀治理）

**Status**: `draft`
**Owner**: linyihong
**建立日期**: 2026-06-02
**Source**: 2026-06-02 對話 — 使用者反映 `plans/active/` 橫向排列導致 main/sub plan 關係不直覺、難追蹤
**Glossary Impact**: yes — 新引入 framework vocabulary：`plan_kind` / `sub_plan_trigger` / `completion_blocks_parent` / `plan tree`，須註冊到 `knowledge/glossary/ai-skill.md`

> **Watch-Out List citation**：本 plan 對應 [`architecture/ai-native-cognitive-ecosystem-system.md`](../../../architecture/ai-native-cognitive-ecosystem-system.md) §Watch-Out List 的「process bloat」與「premature abstraction」防呆 — 不為了結構而結構，只有當 plan 實際出現拆分需求時才開 sub-plan。

---

## Decision Rationale

### Problem & Why Now

`plans/active/` 目前是**扁平目錄**（橫向 `ls` list），但實際執行常出現：

1. 主計畫執行中發現要拆 sub-scope，agent 只能 (a) 在主計畫塞 Phase 6/7/8 直到無法管理，或 (b) 另開一個新 plan 用「與其他 plans 的關係」段落手寫 parent ↔ child 連結。
2. 從 `ls plans/active/` 看不出哪些是主軸、哪些是衍生支線。
3. 主計畫完成條件常被 sub-plan blocker 卡住，但表頭 `Status` 不會自動反映。
4. 已 archived 的 sub-plan 回頭追溯時找不到 parent context。

具體案例：`2026-06-01-0100-validation-scenario-governance-executor.md` 是 `2026-05-31-2100-mechanical-enforcement-registry.md` 的 round-5 U2 child plan，但兩者在 `plans/active/` 完全平鋪，只能靠檔頭 prose 連結。

### Decision

導入 **plan tree 三件套**：

- **A. 目錄結構** — 每個主計畫一個 folder，folder 名即 main plan slug；主計畫本體為 `_plan.md`；sub-plan 用 `NN-<slug>.md` 數字前綴排序；孫計畫多時可再包 `NN-<slug>/` 巢狀資料夾。
- **B. Frontmatter schema** — `plan_kind` / `parent` / `children` / `sub_plan_trigger` / `completion_blocks_parent`。
- **C. Hard rule + commit-msg validator** — sub-plan 必須在 frontmatter 標明 trigger（對應 §When to Open Sub-Plan 的 #1-#5），validator 機械擋。
- **D. 主計畫必填 sub-plan 驗證要點表** — 細節留在 sub-plan，主計畫只列「sub-plan slug / 完成條件摘要 / completion_blocks_parent / 驗證方式」。
- **E. `ai-skill plans tree` CLI subcommand** — 渲染樹狀視圖。

### Alternatives Considered

- **A. 維持現狀（扁平 + frontmatter parent line）** — reject：使用者已明確反映「橫向看不直覺」，純 frontmatter 仍需工具才能看到階層。
- **B. 純 frontmatter + tree CLI（不動目錄）** — reject as primary：向後相容好但解不了 `ls` 直觀性。可作為 Phase 過渡相容。
- **C. 改寫成單一 monolithic plan + 內部 sub-section** — reject：sub-plan 需獨立 archive / sign-off / owner / parallelization 時撐不住，違反 `plans/README.md` 對「跨 phase + 獨立 acceptance」的拆分原則。
- **D. 採用 folder + `_plan.md`（accept）** — `ls plans/active/<folder>/` 即看到主計畫與 sub-plan 共存；前綴排序給人類；frontmatter parent/children 給機器；hard validator 防 trigger 缺失。

### Why Not an ADR Yet

- 尚未跑過至少一輪「現有 plan 遷移」驗證 folder 結構是否真的解決追蹤痛點。
- `sub_plan_trigger` 5 條規則是否完整、會不會在實作中發現第 6 條，仍需 implementation phase 證實。
- 與 `governance/lifecycle/system-upgrade-governance.yaml` §`define_runtime_trigger_flow` 的 sub-plan validator 整合方式尚未驗證。
- Open Questions 仍有 5 條未解。

### ADR Promotion Criteria（completed 時驗證）

- [ ] foundational + cross-session + cross-project + expensive-to-reverse + explains-why 全中
- [ ] Plan 結果證實 folder 結構 + frontmatter 雙軌可長期維護
- [ ] Open Questions 全解
- [ ] 沒有更輕的 promotion target 適用（governance rule 本身可能就夠用，不需升 ADR）
- [ ] ≥ 3 個主計畫已用新格式跑完整生命週期（active → archived）

### Consequences

**正面**：
- `ls plans/active/<folder>/` 直接看到階層
- Sub-plan 細節與主計畫驗證要點分離，主計畫不再臃腫
- Validator 機械擋住「未標 trigger 就開 sub-plan」與「主計畫提前 archive 卻有 child blocker」
- Tree CLI 給快速 status overview

**負面**：
- 一次性遷移成本（現有 6 個 active plan + N 個 archived plan）
- Folder + frontmatter 雙重 source of truth，需 validator 保一致
- 新增 3 條 commit-msg validator，增加 hook 執行時間

**風險**：
- 過度結構化 — 簡單 plan 被迫塞 folder（mitigation: `plan_kind: spike` 簡化模板）
- 巢狀資料夾深度失控（mitigation: validator 限制 ≤ 2 層巢狀）
- 遷移期 archived plan 仍是扁平結構，CLI tree 需相容兩種格式

---

## Deferred Runtime Projection

本 plan 不新增 `runtime/*.yaml`。Frontmatter schema 直接由 commit-msg validator（Go 程式碼）解析；plan tree 渲染由 CLI subcommand 即時讀檔案系統，不 project 到 `runtime.db`。理由：

- (a) plan tree state 屬於檔案系統 reality，不需 runtime cache。
- (b) `runtime.db` 已有 `plan_status_sync` validator 對應，新增 projection 會雙份 source-of-truth。

若未來發現需要 query「all sub-plans of plan X」這類跨檔操作，再評估 project 到 `runtime.db generated_surfaces[plans.tree]`。

---

## Runtime Execution Path

| 階段 | Runtime owner | Trigger flow | Loaded contract | Runtime action / blocker | Evidence |
|------|---------------|--------------|-----------------|--------------------------|----------|
| 新建 sub-plan commit | `commit-msg` hook | git commit staging `plans/active/<main>/NN-*.md` 或 `plans/active/<main>/NN-*/_plan.md` → `validatePlanTreeFrontmatter` | 本 plan §Frontmatter Schema + §When to Open Sub-Plan | Block 若 `parent` 缺、`sub_plan_trigger` 缺、或 trigger 值不在 #1-#5 enum | unit test fixture `testdata/plan-tree/sub-without-trigger.md` |
| 主計畫 archive commit | `commit-msg` hook | git commit moving `plans/active/<main>/_plan.md` → `plans/archived/<main>/_plan.md` → `validatePlanTreeArchiveOrder` | 主計畫 frontmatter `children[]` + 各 child `completion_blocks_parent` | Block 若有 `completion_blocks_parent: true` 的 child 仍在 `plans/active/` | unit test fixture `testdata/plan-tree/archive-with-pending-child.md` |
| Folder shape lint | `commit-msg` hook | git commit staging `plans/active/<x>/**` → `validatePlanTreeFolderShape` | 本 plan §資料夾結構 | Block 若 folder 缺 `_plan.md`、檔名不符 `NN-` 前綴、或巢狀 > 2 層 | unit test fixture `testdata/plan-tree/missing-_plan.md` |
| Tree 渲染 | `ai-skill plans tree` CLI | 使用者執行 → 遞迴讀 `plans/active/` + `plans/archived/` → 解析 frontmatter | 本 plan §Frontmatter Schema | Print 樹狀 + status 進度（非 blocker） | golden test `testdata/plan-tree/tree-output.txt` |

**Forbidden 自我檢查**（per `define_runtime_trigger_flow`）：
- (a) 本 plan 不加 `route.*` 到 routing-registry。Validator 直接由 hooks.go registry dispatch，符合「discovery signal / commit-msg validator」要求。
- (b) 本 plan 不新增 `generated_surfaces`，故不需宣告 consumer。

---

## Per-surface consumer 表

| Generated surface key | Named consumer(s) | Consumer 類型 |
|---|---|---|
| （無新增 surface） | n/a | n/a |

本 plan 僅新增 3 個 Go validator 與 1 個 CLI subcommand，不產生 generated surface；故表為空，符合 §`define_runtime_trigger_flow` forbidden rule (b) 例外（無 surface = 無 consumer 義務）。

---

## 資料夾結構（規範）

```
plans/
  active/
    2026-06-02-1200-plan-tree-hierarchy-governance/   ← 主計畫資料夾，名 = main plan slug
      _plan.md                                        ← 主計畫本體（檔名固定 _plan.md）
      01-frontmatter-schema.md                        ← sub-plan，NN- 前綴排序
      02-validator-implementation.md
      03-cli-tree-subcommand/                         ← 孫計畫資料夾（當 sub-plan 自己也要拆時）
        _plan.md
        01-renderer.md
        02-golden-tests.md
      04-existing-plan-migration.md
  archived/
    2026-06-02-1200-plan-tree-hierarchy-governance/   ← archive 時整個資料夾搬移
      _plan.md
      01-frontmatter-schema.md
      ...
```

**規則**：
- 主計畫 folder 名 = main plan slug（含時間戳）。
- 主計畫本體檔名固定 `_plan.md`（底線開頭，排序在所有 `NN-` 前面）。
- Sub-plan 檔名格式：`NN-<slug>.md`（NN 為兩位數字，建議 01/02/03…），或當 sub-plan 自己也要拆時用 `NN-<slug>/_plan.md`。
- 巢狀深度上限 **2 層**（main → sub → grand-sub），超過時應考慮拆出獨立 main plan。
- Archive 時整個 folder 從 `plans/active/<slug>/` 搬到 `plans/archived/<slug>/`，folder 名不變。

---

## Frontmatter Schema

**主計畫**：

```yaml
---
id: <slug>
plan_kind: main
status: draft | in-progress | completed
owner: <name>
created: YYYY-MM-DD
parent: null
children:
  - 01-frontmatter-schema
  - 02-validator-implementation
  - 03-cli-tree-subcommand
---
```

**Sub-plan**：

```yaml
---
id: <main-slug>/NN-<sub-slug>
plan_kind: sub
status: draft | in-progress | completed
owner: <name>
created: YYYY-MM-DD
parent: <main-slug>
sub_plan_trigger: independent-signoff | multi-phase-with-own-acceptance | independent-runtime-trigger | parallel-owners | independent-archive-spike
completion_blocks_parent: true | false
---
```

**Spike**（簡化模板，只需 Phase 0 + Acceptance + 結果回寫）：

```yaml
---
id: <main-slug>/NN-<sub-slug>
plan_kind: spike
status: draft | completed
owner: <name>
created: YYYY-MM-DD
parent: <main-slug>
sub_plan_trigger: independent-archive-spike
completion_blocks_parent: false
---
```

---

## When to Open Sub-Plan（hard rule，validator 強制）

開新 sub-plan 的**強制條件**（`sub_plan_trigger` 必填，必須對應其中一條）：

| Trigger enum 值 | 觸發條件 | 範例 |
|---|---|---|
| `independent-signoff` | 該支線需要獨立 stakeholder sign-off / acceptance 簽核 | DSL schema 設計獨立於 executor wiring |
| `multi-phase-with-own-acceptance` | 該支線跨 ≥ 3 個 phase 且有自己的 completion criteria | runtime trigger wiring 跨 schema + validator + readback |
| `independent-runtime-trigger` | 該支線有自己的 runtime trigger flow / generated surface | 新增 `route.validation.executor` 需 wire discovery signal |
| `parallel-owners` | 兩個 parallel agent / session 需要同時推進，需 owner / lock 分隔 | child 拆給不同 owner 並行 |
| `independent-archive-spike` | 該工作完成後可獨立 archive（主計畫仍 in-progress） | spike / experiment / 短期 PoC |

**不該開 sub-plan 的情境**（應留在主計畫加 phase 或 checkbox）：

- 單一 phase 內的 step 細分 → 用 checkbox。
- < 1 工作 session 可完成 → inline 寫進主計畫。
- 純文件補強、rename、typo → 直接 commit，不開 plan。
- 同一 acceptance criteria 底下的不同 angle → 同 plan 多 phase。

---

## 主計畫必填：Sub-Plan 驗證要點表

主計畫 §Phases 之後必填下列表，列出每個 sub-plan 的「驗證要點摘要」。細節留在 sub-plan 內。

| Sub-plan | 完成條件摘要 | completion_blocks_parent | 驗證方式 |
|---|---|---|---|
| `01-frontmatter-schema` | schema YAML 文件化 + 範例 fixture | true | unit test pass + 本 plan §Frontmatter Schema 對齊 |
| `02-validator-implementation` | 3 個 Go validator 落地 + dispatch registry | true | `go test ./scripts/ai-skill-cli/internal/app/...` pass |
| `03-cli-tree-subcommand` | `ai-skill plans tree` 渲染 active + archived | false | golden test + 手動跑 CLI 驗證輸出 |
| `04-existing-plan-migration` | 盤點 6 個 active plan + 識別 parent-child 並遷移 | true | 遷移後 `ai-skill plans tree` 顯示正確階層 |

---

## Phase 0 — Pre-Build Interrogation

### Phase 0.0 — Open Questions 核對（公版，必填）

逐條核對本 plan §Open Questions：

- [ ] 已讀本 plan §Open Questions 全部條目
- [ ] 對每條標記 `resolved`（附 Phase 0 證據）/ `still-open` / `deferred`（附原因）
- [ ] `resolved` 的條目已同步勾選 / 附註於 §Open Questions
- [ ] 若盤點新發現問題，已加入 §Open Questions

| Open Question | 處置 | 證據 / 原因 |
|---|---|---|
| Q1 巢狀深度上限 | still-open | 待 Phase 0 盤點現有 plans 最深巢狀需求 |
| Q2 Archive 順序 | still-open | 待 Phase 0 確認搬移腳本 |
| Q3 Sub-plan Decision Rationale | still-open | 待 Phase 0 評估模板可繼承 |
| Q4 Validator 強度 | resolved | 使用者明確選 hard rule + validator |
| Q5 Spike 模板 | still-open | 待 Phase 0 評估最小章節集 |

### Phase 0 — Preflight checklist

- [ ] 盤點 `plans/active/` 6 個 plan，標記哪些有隱性 parent-child 關係
- [ ] 盤點 `plans/archived/` 找出歷史 parent-child 範例（特別是 `bootstrap-contract-yaml-migration` 系列）
- [ ] 確認 `scripts/ai-skill-cli/internal/app/hooks.go` validator registry 加入新 validator 的相容性
- [ ] 確認 `plans/README.md` 模板章節需要哪些連動更新
- [ ] 確認 `enforcement/linked-updates.yaml` 是否需要新規則
- [ ] 驗證 folder 結構與既有 `validatePlanArchivalAudit`、`validatePlanCheckboxSync`、`validatePlanStatusSync` 不衝突

---

## Phase 1 — `01-frontmatter-schema`（sub-plan）

詳見 [`01-frontmatter-schema.md`](01-frontmatter-schema.md)（待建）。本主計畫驗證要點：schema 文件化 + ≥ 3 個 fixture（main / sub / spike）。

---

## Phase 2 — `02-validator-implementation`（sub-plan）

詳見 [`02-validator-implementation.md`](02-validator-implementation.md)（待建）。本主計畫驗證要點：
- `validatePlanTreeFrontmatter` block sub-plan 缺 `sub_plan_trigger`
- `validatePlanTreeArchiveOrder` block 主計畫 archive 時 child 未 archive
- `validatePlanTreeFolderShape` block folder 缺 `_plan.md` 或檔名違規
- 3 個 validator 進 `hooks.go` registry，dispatch 順序與既有 11 個 validator 不衝突

---

## Phase 3 — `03-cli-tree-subcommand`（sub-plan）

詳見 [`03-cli-tree-subcommand/_plan.md`](03-cli-tree-subcommand/_plan.md)（待建）。本主計畫驗證要點：`ai-skill plans tree` 可渲染 active + archived 兩種狀態，輸出包含 status / 進度 / blocker。

---

## Phase 4 — `04-existing-plan-migration`（sub-plan）

詳見 [`04-existing-plan-migration.md`](04-existing-plan-migration.md)（待建）。本主計畫驗證要點：6 個 active plan 全部評估完畢；至少 1 組 parent-child 已用新格式遷移；遷移後 `ai-skill plans tree` 顯示正確階層。

---

## Phase 5 — 收尾（在主計畫直接做，不開 sub-plan）

- [ ] 更新 `plans/README.md`：新增 §Plan Tree Hierarchy 章節，連結本 plan
- [ ] 更新 `governance/` 或 `enforcement/`：新增 `plan-tree-hierarchy.md` rule
- [ ] 註冊 glossary terms 到 `knowledge/glossary/ai-skill.md`
- [ ] 寫 failure pattern `enforcement/failure-patterns/plan-tree-flat-ambiguity.md`
- [ ] 5 platform binaries rebuild + push
- [ ] 本 plan 整 folder 搬 `plans/archived/`

---

## Open Questions

1. **巢狀深度上限** — 目前提案 2 層。是否該允許 3 層（main → sub → grand-sub → great-grand-sub）？傾向不允許，超過時拆出獨立 main plan。
2. **Archive 順序** — 主計畫 + sub-plan 同 commit 一起搬 archived，還是分次？傾向同 commit（atomic），由 validator 確保順序。
3. **Sub-plan Decision Rationale 可繼承否** — sub-plan 是否需自己的 §Decision Rationale，或可在 frontmatter 標 `inherits_rationale: parent`？傾向 sub-plan 不需重複 Decision Rationale，但需有 §為什麼存在（簡短）+ §Acceptance。
4. **Validator 強度** — RESOLVED：使用者選 hard rule + commit-msg validator block。
5. **Spike 模板最小集** — `plan_kind: spike` 是否可只有 §Goal + §Acceptance + §結果回寫，免 Phase 0 公版？傾向是，但結果必須回寫主計畫對應 phase。

---

## 完成條件

- [ ] Phase 0 preflight 全部完成
- [ ] 4 個 sub-plan 全部 `status: completed` 且 archived
- [ ] Phase 5 收尾項目全 checked
- [ ] `governance/plan-tree-hierarchy.md` 落地
- [ ] 3 個新 validator 進 `hooks.go` registry，unit test pass
- [ ] `ai-skill plans tree` CLI 可用
- [ ] 至少 1 組真實 parent-child plan 用新格式跑完（可用本 plan 自身作為 dogfood）
- [ ] Failure pattern 文件化
- [ ] Glossary terms 註冊
- [ ] 5 platform binaries rebuilt + pushed
- [ ] 本 plan 整 folder 搬 `plans/archived/2026-06-02-1200-plan-tree-hierarchy-governance/`

---

## Stakeholder 同意項目

- [ ] linyihong: 目錄結構 = folder + `_plan.md` + `NN-` 前綴（**已 sign-off 2026-06-02**）
- [ ] linyihong: 「何時開 sub-plan」用 hard rule + validator（**已 sign-off 2026-06-02**）
- [ ] linyihong: Open Questions Q1/Q2/Q3/Q5 解法（待 Phase 0 提案後 sign-off）
- [ ] linyihong: 是否將 plan tree 結構 promote 為 ADR（completed 後評估）

---

## Glossary Impact

**Glossary Impact: yes**

新引入 framework vocabulary：

| Term | 定義 | 註冊目標 |
|---|---|---|
| `plan_kind` | plan 的類型 enum：`main` / `sub` / `spike` | `knowledge/glossary/ai-skill.md` |
| `sub_plan_trigger` | sub-plan 必填的開啟原因 enum，5 個值見 §When to Open Sub-Plan | `knowledge/glossary/ai-skill.md` |
| `completion_blocks_parent` | sub-plan 是否阻擋主計畫 archive 的 boolean flag | `knowledge/glossary/ai-skill.md` |
| `plan tree` | 主計畫 + sub-plan 構成的階層結構 | `knowledge/glossary/ai-skill.md` |

註冊將在 Phase 5 進行。

---

## 與其他 plans 的關係

- **Parent**: 無（本 plan 自身為 main plan）。
- **相關 plan**:
  - `plans/active/2026-05-31-2100-mechanical-enforcement-registry.md` — registry executors 註冊新 validator 的目標位置（Phase 2 整合）。
  - `plans/active/2026-06-01-0100-validation-scenario-governance-executor.md` — 既有的隱性 parent-child 範例（Phase 4 遷移 target）。
  - `plans/archived/2026-05-25-2200-bootstrap-contract-yaml-migration.md` — 歷史 multi-phase plan，可作為 Phase 4 「歷史 parent-child」 reference。
- **不衝突**:
  - 既有 `validatePlanArchivalAudit` / `validatePlanCheckboxSync` / `validatePlanStatusSync` 保持運作，新 validator 為加法。

---

## Dependency Read Ledger

| 欄位 | 內容 |
|---|---|
| Trigger | 2026-06-02 使用者反映 plans/active/ 橫向難追蹤、要求設計 main/sub plan 體系 + 開 sub-plan 規範 |
| Required set | `plans/README.md`、`enforcement/rule-weight.md`、`enforcement/dependency-reading.md`、`enforcement/conversation-goal-ledger.md`、`runtime/core-bootstrap.yaml`、`governance/lifecycle/system-upgrade-governance.yaml` §`define_runtime_trigger_flow`、現有 plan 範例 |
| Read | CORE_BOOTSTRAP.md、runtime/core-bootstrap.yaml、enforcement/rule-weight.md、enforcement/dependency-reading.md、enforcement/conversation-goal-ledger.md、plans/README.md（前 150 行）、`2026-06-01-0100-validation-scenario-governance-executor.md` |
| Not applicable | `validation/scenarios/` 結構（本 plan 不新增 scenario） |
| Deferred | `scripts/ai-skill-cli/internal/app/hooks.go` 完整 registry 結構（Phase 0 preflight 補讀） |
| Validation | 本 plan draft 提交後由使用者 review；Phase 0 完成 architecture compatibility preflight；implementation phase 跑 commit-msg validator unit tests |
