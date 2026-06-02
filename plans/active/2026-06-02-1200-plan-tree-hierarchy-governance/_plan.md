---
id: 2026-06-02-1200-plan-tree-hierarchy-governance
plan_kind: main
status: draft
owner: linyihong
created: 2026-06-02
parent: null
---

# Plan Tree Hierarchy Governance（主計畫／子計畫樹狀治理）

**Status**: `draft`
**Owner**: linyihong
**建立日期**: 2026-06-02
**Source**: 2026-06-02 對話 — 使用者反映 `plans/active/` 橫向排列導致 main/sub plan 關係不直覺、難追蹤
**Glossary Impact**: yes — 新引入 framework vocabulary：`plan_kind` / `parent` / `required_for_completion` / `sub_plan_reason` / `plan tree`，須註冊到 `knowledge/glossary/ai-skill.md`

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

導入 **Minimal Governance plan tree**（2026-06-02 review 後 v2）：

- **A. Source of truth = frontmatter `id` + `parent`**（單一 truth）。Hierarchy 由 `parent` pointer 決定，CLI 從這裡建樹。
- **B. Folder + filename 是 UI convention，不是 source of truth**。主計畫 folder + `_plan.md` + `NN-<slug>.md` 前綴只是「人類眼睛友善」的擺放方式；validator 對 folder shape 只發 warning，不 block。
- **C. Frontmatter schema 最小集**：`id` / `plan_kind` / `parent` / `status` / `owner` / `created` / `required_for_completion`（sub only） / `sub_plan_reason`（sub only, free text）。
  - **不維護 `children:`** — runtime scan 推導。children 是衍生資料，parent 才是必要資料。
  - **不用 `completion_blocks_parent`** — 該詞描述了 archive/blocking 機制；改用 `required_for_completion: true|false` 描述業務意義（child 是否屬於 parent acceptance criteria），由 validator 推導出「未完成 → block parent archive」。
  - **不用 `sub_plan_trigger` enum** — 「為什麼拆 plan」是知識，不是 enum。改用 `sub_plan_reason: <free text>`，validator 只檢查非空。
- **D. 主計畫只列「sub-plan slug + 驗證要點摘要 + required_for_completion + 驗證方式」表**，細節留在 sub-plan。
- **E. `ai-skill plans tree` CLI subcommand**：純讀 frontmatter `parent` 動態建樹；即使 folder 放錯，tree 仍正確。
- **F. Depth 防呆 = `warning_at_depth: 3`**（非 block）— 不立法、讓真實案例決定上限。

### Alternatives Considered

- **A. 維持現狀（扁平 + 無 metadata）** — reject：使用者已明確反映「橫向看不直覺」，且 parent ↔ child 連結只能靠 prose。
- **B. 純 folder + 無 frontmatter** — reject：folder 放錯時整個 hierarchy 失效，沒有 robust source of truth。
- **C. 純 frontmatter + 扁平目錄（無 folder）** — reject as primary：解不了 `ls` 直觀性，仍需工具才能看到階層。但 frontmatter 是 truth，folder 是 UI，所以這方案可作為 graceful degradation（folder 缺失時 CLI 仍可建樹）。
- **D. Folder 結構 + 三軌 source（folder + frontmatter `children` + filename ordering）** — **rejected v1（2026-06-02 review）**：三個地方都在表達階層，會出現 `children: [01-schema]` 但實際 rename 成 `03-validator.md` 的同步問題。誰才是真實來源？
- **E. Folder 結構 + frontmatter `parent` 單一 truth + folder 為 UI convention（accept）** — `ls plans/active/<folder>/` 給人類友善視覺；frontmatter `parent` 給機器；folder shape validator 只發 warning 不 block；children 由 runtime scan 推導。
- **F. Hard enum trigger（`sub_plan_trigger: independent-signoff | ...`）** — **rejected v1**：enum 半年後會出現第 6 / 7 / 8 種，每次都要升 framework（enum + validator + glossary + docs）。改 free-text `sub_plan_reason: <非空>`，validator 只擋空字串。
- **G. `completion_blocks_parent: bool`（描述機制）** — **rejected v1**：重新發明 dependency graph；當 sub-plan 之間有 dependency（C 依賴 A）時不夠用。改 `required_for_completion: bool`（描述業務意義「是否屬 parent acceptance criteria」），validator 推導 archive block 邏輯。
- **H. `max_depth: 2`（硬限制）** — **rejected v1**：premature abstraction，目前根本沒出現 3 層需求。改 `warning_at_depth: 3`，讓真實案例決定。

### Why Not an ADR Yet

- 尚未跑過至少一輪「現有 plan 遷移」驗證 folder convention + frontmatter 雙軌是否真的解決追蹤痛點。
- `sub_plan_reason` free-text 模式是否真能在無 enum 強制下保持品質，仍需 implementation phase 證實。
- 與 `governance/lifecycle/system-upgrade-governance.yaml` §`define_runtime_trigger_flow` 的 sub-plan validator 整合方式尚未驗證。
- Open Questions 仍有 4 條未解（Q4 已 resolved as v2 minimal-governance）。

### ADR Promotion Criteria（completed 時驗證）

- [ ] foundational + cross-session + cross-project + expensive-to-reverse + explains-why 全中
- [ ] Plan 結果證實 folder 結構 + frontmatter 雙軌可長期維護
- [ ] Open Questions 全解
- [ ] 沒有更輕的 promotion target 適用（governance rule 本身可能就夠用，不需升 ADR）
- [ ] ≥ 3 個主計畫已用新格式跑完整生命週期（active → archived）

### Consequences

**正面**：
- `ls plans/active/<folder>/` 給人類直觀視覺，但 hierarchy 真實來源是 frontmatter `parent`，folder 放錯不會壞掉
- Sub-plan 細節與主計畫驗證要點分離，主計畫不再臃腫
- Validator 機械擋「frontmatter 缺欄位」「archive 時 required child 未 completed」「dangling parent pointer」「duplicate id」；folder shape 只發 warning，不擋 commit
- Tree CLI 從 `parent` pointer 動態建樹，給快速 status overview；referential integrity 由 validator 保證，不會出現 orphan node

**負面**：
- 一次性遷移成本（現有 6 個 active plan + N 個 archived plan）
- 新增 4 條 block validator + 1 條 warning validator，增加 hook 執行時間；其中 ParentReference / UniqueID 需要 scan 全 repo plan 集合，evaluator 需做檔案級 cache 避免 N² 成本

**風險**：
- 過度結構化 — 簡單 plan 被迫塞 folder（mitigation: `plan_kind: spike` 簡化模板）
- 巢狀深度失控（mitigation: `warning_at_depth: 3`，由真實案例決定上限，不立法硬擋）
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
| 新建 sub-plan commit | `commit-msg` hook | git commit staging plan with `plan_kind: sub` → `validatePlanTreeFrontmatter` | 本 plan §Frontmatter Schema | **Block** 若 `parent` 缺、`sub_plan_reason` 為空或缺、或 `required_for_completion` 缺 | unit test fixture `testdata/plan-tree/sub-missing-parent.md`、`sub-empty-reason.md` |
| 主計畫 archive commit | `commit-msg` hook | git commit moving `plans/active/<main>/` → `plans/archived/<main>/` → `validatePlanTreeArchiveOrder` | Runtime scan：所有 `parent == <main>` 且 `required_for_completion: true` 的 sub-plan `status` | **Block** 若有 `required_for_completion: true` 的 sub-plan `status != completed`。**只看 lifecycle status，不看 location**（sub-plan completed 但仍在 active/ 不阻擋 parent archive — archive 是儲存位置，completed 是生命週期狀態，不混） | unit test fixture `testdata/plan-tree/archive-with-required-pending.md` |
| Parent reference 檢查 | `commit-msg` hook | git commit staging sub-plan with `parent: <id>` → `validatePlanTreeParentReference` | 全 repo scan `plans/active/**/*.md` + `plans/archived/**/*.md` 收集 `id` 集合 | **Block** 若 `parent` 指向的 id 在 active + archived 都找不到（dangling pointer / orphan node） | unit test fixture `testdata/plan-tree/parent-orphan.md` |
| ID 唯一性檢查 | `commit-msg` hook | git commit staging any plan → `validatePlanTreeUniqueID` | 全 repo scan plan frontmatter `id` 欄位 | **Block** 若同一 `id` 出現在 ≥ 2 個檔案（含 active vs archived 跨目錄重複） | unit test fixture `testdata/plan-tree/duplicate-id.md` |
| Folder shape lint | `commit-msg` hook | git commit staging `plans/active/<x>/**` → `validatePlanTreeFolderConvention` | 本 plan §資料夾 convention | **Warning only**（不 block）— folder 缺 `_plan.md`、檔名不符 `NN-` 前綴、或深度 ≥ 3。輸出建議訊息，不擋 commit | unit test fixture `testdata/plan-tree/depth-3-warning.md` |
| Tree 渲染 | `ai-skill plans tree` CLI | 使用者執行 → 遞迴讀 `plans/active/` + `plans/archived/` → 解析 frontmatter `parent` → 動態建樹 | 本 plan §Frontmatter Schema | Print 樹狀 + status 進度 + warning（非 blocker）；即使 folder 放錯仍能建出正確 tree | golden test `testdata/plan-tree/tree-output.txt` |

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

## 資料夾 Convention（UI 層，非 source of truth）

> **重要**：以下是給人類眼睛友善的擺放方式。Hierarchy 的真實來源是 frontmatter `parent` pointer；folder 放錯時 CLI 仍能建出正確樹。Validator 對本節只發 warning，不 block commit。

```
plans/
  active/
    2026-06-02-1200-plan-tree-hierarchy-governance/   ← 主計畫資料夾，名建議 = main plan slug
      _plan.md                                        ← 主計畫本體（檔名建議 _plan.md）
      01-frontmatter-schema.md                        ← sub-plan，NN- 前綴給人類排序
      02-validator-implementation.md
      03-cli-tree-subcommand/                         ← sub-plan 自己也要拆時可包資料夾
        _plan.md
        01-renderer.md
        02-golden-tests.md
      04-existing-plan-migration.md
  archived/
    2026-06-02-1200-plan-tree-hierarchy-governance/   ← archive 時整個資料夾搬移（建議）
      _plan.md
      ...
```

**Convention（warning，非 block）**：
- 主計畫 folder 名建議 = main plan slug（含時間戳）。
- 主計畫本體檔名建議 `_plan.md`（底線開頭，自然排序在所有 `NN-` 前）。
- Sub-plan 檔名建議 `NN-<slug>.md`（NN 兩位數字），或 `NN-<slug>/_plan.md`。
- 巢狀深度 ≥ 3 層時發 warning，建議拆出獨立 main plan（不硬擋）。
- Archive 時建議整個 folder 一起搬。

**Block 規則四條**（在 §Runtime Execution Path）：
1. Sub-plan frontmatter 缺 `parent` / `sub_plan_reason` / `required_for_completion` → block。
2. 主計畫 archive 時，`required_for_completion: true` 的 sub-plan `status != completed` → block。**只看 status，不看 location**。
3. Sub-plan `parent` 指向的 id 不存在（dangling pointer） → block。
4. 同一 `id` 出現在 ≥ 2 個檔案（duplicate id） → block。

---

## Frontmatter Schema（Minimal Governance）

**主計畫**：

```yaml
---
id: <slug>
plan_kind: main
status: draft | in-progress | completed
owner: <name>
created: YYYY-MM-DD
parent: null
---
```

**Sub-plan**：

```yaml
---
id: <slug>                               # 全域唯一，不必包含 parent path
plan_kind: sub
status: draft | in-progress | completed
owner: <name>
created: YYYY-MM-DD
parent: <main-slug>                      # 指向 parent plan 的 id
required_for_completion: true | false    # 是否屬於 parent acceptance criteria
sub_plan_reason: >                       # 為什麼拆 plan（free text，非空）
  簡述拆分理由與 acceptance 邊界
---
```

**Spike**（簡化模板，只需 Goal + Acceptance + 結果回寫）：

```yaml
---
id: <slug>
plan_kind: spike
status: draft | completed
owner: <name>
created: YYYY-MM-DD
parent: <main-slug>
required_for_completion: false           # spike 預設不阻擋 parent archive
sub_plan_reason: >
  PoC / experiment 目的與時限
---
```

**設計原則**：
- `id` 是全域唯一 slug，**不再要求** `id: parent/sub-slug` 這種 path-like 格式（避免 rename 時連動成本）。
- **沒有 `children:` 欄位** — 由 runtime scan `parent` pointer 推導。
- **沒有 `completion_blocks_parent:`** — 用 `required_for_completion` 描述業務語意，archive blocker 由 validator 推導。
- **沒有 `sub_plan_trigger` enum** — 用 `sub_plan_reason` free text，validator 只擋空字串。
- 未來若需要表達 sub-plan 之間的依賴（C 依賴 A），可加 `depends_on: [<sub-id>]` 欄位（**不在本 plan 範圍**，留待真實案例驅動）。

---

## When to Open Sub-Plan（建議規則，非 enum）

**Validator 強制條件**：`sub_plan_reason` 非空字串。**不審內容**。

下列為**建議參考**（recommended triggers），寫進 `sub_plan_reason` 時可引用，但不強制：

| Recommended trigger | 觸發條件 | 範例 |
|---|---|---|
| Independent sign-off | 該支線需要獨立 stakeholder sign-off / acceptance 簽核 | DSL schema 設計獨立於 executor wiring |
| Multi-phase with own acceptance | 該支線跨 ≥ 3 個 phase 且有自己的 completion criteria | runtime trigger wiring 跨 schema + validator + readback |
| Independent runtime trigger | 該支線有自己的 runtime trigger flow / generated surface | 新增 `route.validation.executor` 需 wire discovery signal |
| Parallel owners | 兩個 parallel agent / session 需要同時推進，需 owner / lock 分隔 | child 拆給不同 owner 並行 |
| Independent archive (spike) | 該工作完成後可獨立 archive（主計畫仍 in-progress） | spike / experiment / 短期 PoC |

未來真實案例如出現第 6 / 7 種情境，**直接寫進 `sub_plan_reason` 即可，不需升 framework**。等該情境重複出現 ≥ 3 次再評估是否要 promote 為 recommended trigger 範例。

**不該開 sub-plan 的情境**（應留在主計畫加 phase 或 checkbox）：

- 單一 phase 內的 step 細分 → 用 checkbox。
- < 1 工作 session 可完成 → inline 寫進主計畫。
- 純文件補強、rename、typo → 直接 commit，不開 plan。
- 同一 acceptance criteria 底下的不同 angle → 同 plan 多 phase。

**設計原則**：「為什麼拆 plan」是知識，不是 enum。Framework 只擋「沒寫理由」，不擋「理由不在白名單」。

---

## 主計畫必填：Sub-Plan 驗證要點表

主計畫 §Phases 之後必填下列表，列出每個 sub-plan 的「驗證要點摘要」。細節留在 sub-plan 內。

| Sub-plan | 完成條件摘要 | required_for_completion | 驗證方式 |
|---|---|---|---|
| `01-frontmatter-schema` | schema YAML 文件化 + 範例 fixture（main / sub / spike） | true | unit test pass + 本 plan §Frontmatter Schema 對齊 |
| `02-validator-implementation` | 4 個 block validator + 1 個 warning validator 落地 + dispatch registry | true | `go test ./scripts/ai-skill-cli/internal/app/...` pass |
| `03-cli-tree-subcommand` | `ai-skill plans tree` 從 `parent` pointer 渲染 active + archived | false | golden test + 手動跑 CLI 驗證輸出 |
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
| Q1 深度防呆模式 | resolved | 2026-06-02 review v2：改 `warning_at_depth: 3`（非 block），讓真實案例決定上限 |
| Q2 Archive 順序 | resolved | 2026-06-02 v2.2：lifecycle / storage 分離；gate 只看 status，folder 搬法為操作自由 |
| Q3 Sub-plan Decision Rationale | resolved | 2026-06-02 v2.2：主計畫負責 Why、子計畫負責 How；sub-plan 預設繼承 parent rationale，必填僅 Purpose + Acceptance + Runtime Impact |
| Q4 Validator 強度 | resolved | 2026-06-02 review v2：parent / sub_plan_reason / required_for_completion 缺失為 block；folder shape 為 warning（minimal governance） |
| Q5 Spike 模板 | still-open | 待 Phase 0 評估最小章節集 |

### Phase 0 — Preflight checklist

- [x] 盤點 `plans/active/` 6 個 plan，標記隱性 parent-child（2026-06-02 完成）
- [ ] 盤點 `plans/archived/` 找出歷史 parent-child 範例（deferred — Phase 4 migration 啟動時做）
- [x] 確認 `scripts/ai-skill-cli/internal/app/hooks.go` validator registry 衝突風險（與 registry plan）— 已決議採並行安全路：Phase 2 延後至 registry plan archive
- [ ] 確認 `plans/README.md` 模板章節需要哪些連動更新（Phase 5 收尾再做）
- [ ] 確認 `enforcement/linked-updates.yaml` 是否需要新規則（Phase 5 收尾再做）
- [x] 驗證 folder 結構與既有 `validatePlanArchivalAudit` / `validatePlanCheckboxSync` / `validatePlanStatusSync` 不衝突（新 5 個 validator 為加法，不改既有）

### Phase 0 Inventory 結果

| Cluster | Main | Subs（隱性 parent ref 來源） | Migration target |
|---|---|---|---|
| **Cluster 1（Registry tree）** | `2026-05-31-2100-mechanical-enforcement-registry` (P1) | `2026-05-31-1900-workflow-activation-engine` (P2, "parent plan mechanical-enforcement-registry")<br>`2026-05-31-2000-mechanical-sanitization-validator` (P3, "parent meta-plan P1")<br>`2026-06-01-0100-validation-scenario-governance-executor` (stub, "Source: §Phase 3 Round-4 T1") | **Phase 4 dogfood（首選）** — 等該 cluster archive 後遷移成 folder 結構 |
| **Cluster 2（Plan-tree itself）** | `2026-06-02-1200-plan-tree-hierarchy-governance/_plan.md` | `01-frontmatter-schema.md`（Phase 1 已建）+ 03/04 待建 | Self-dogfood，已用新格式 |
| Standalone | `2026-05-27-1557-tool-runtime-signal-economics-integration` | — | 無需遷移 |
| Standalone | `2026-05-28-1636-gen4-fitness-optimization-memory-interface-reservation` | — | 無需遷移 |

Cluster 1 是天然 dogfood：3 個既有 plan 都已在檔頭 prose 明寫 parent 關係，遷移成本低（只需加 frontmatter + 建 folder）。

---

## Phase 1 — `01-frontmatter-schema`（sub-plan）

**Status**：in-progress（2026-06-02 啟動）

詳見 [`01-frontmatter-schema.md`](01-frontmatter-schema.md)（**已建**）。本主計畫驗證要點：schema 文件化 + ≥ 3 個 fixture（main / sub / spike）。

**Phase 1 已交付**：
- [x] `01-frontmatter-schema.md` sub-plan（dogfood new schema）
- [x] `fixtures/main-plan.md`
- [x] `fixtures/sub-plan.md`
- [x] `fixtures/spike-plan.md`
- [x] `governance/lifecycle/plan-tree-hierarchy.md` rule draft

**Phase 1 待完成**：
- [ ] Sub-plan 自己的 acceptance criteria 全 checked（最後 1 條：Phase 2 重用 fixtures 待 Phase 2 啟動時 verify）
- [ ] 將 Phase 1 mark completed 後，由 parent `_plan.md` Phase 5 收尾時連動 `plans/README.md`、glossary 註冊

---

## Phase 2 — `02-validator-implementation`（sub-plan）

詳見 [`02-validator-implementation.md`](02-validator-implementation.md)（待建）。本主計畫驗證要點：
- `validatePlanTreeFrontmatter`（**block**）— sub-plan 缺 `parent` / `sub_plan_reason`（空字串視為缺）/ `required_for_completion`
- `validatePlanTreeArchiveOrder`（**block**）— 主計畫 archive 時，所有 `parent == <main>` 且 `required_for_completion: true` 的 sub-plan 必須 `status: completed`（只看 status，不看 location）
- `validatePlanTreeParentReference`（**block**）— sub-plan `parent` 指向的 id 必須存在於全 repo plan 集合（active + archived）；防 orphan node
- `validatePlanTreeUniqueID`（**block**）— 全 repo plan `id` 必須唯一；防 parent pointer 指錯
- `validatePlanTreeFolderConvention`（**warning only**）— folder 缺 `_plan.md`、檔名不符 `NN-` 前綴、或深度 ≥ 3
- 5 個 validator 進 `hooks.go` registry，dispatch 順序與既有 11 個 validator 不衝突

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

1. **深度防呆模式** — RESOLVED（v2 2026-06-02）：改 `warning_at_depth: 3`，不立硬上限；超過時 CLI 發 warning，建議拆獨立 main plan。原 `max_depth: 2` 是 premature abstraction。
2. **Archive 順序** — RESOLVED（v2.2 2026-06-02）：lifecycle 與 storage 已明確分離 — `validatePlanTreeArchiveOrder` 只看 `status`，不看 location。所以「main + child 同 commit 一起搬 folder」「分次搬」都是合法的操作流程，框架不規範。Archive gate 只要求 required child `status: completed` 即可。
3. **Sub-plan Decision Rationale 可繼承否** — RESOLVED（v2.2 2026-06-02）：明文採「主計畫負責 Why、子計畫負責 How」原則。Sub-plan **不需** §Decision Rationale，**預設繼承 parent**。Sub-plan 必填章節最小集：(a) §Purpose（為什麼存在，簡短）、(b) §Acceptance Criteria、(c) §Runtime Impact（若有 runtime trigger / generated surface 改動）。避免 4 個 child 全部重寫 rationale 的維護負擔。
4. **Validator 強度** — RESOLVED（v2 2026-06-02）：minimal governance — `parent` / `sub_plan_reason`（空字串視為缺）/ `required_for_completion` 缺失為 block；不審 reason 內容；folder shape 全部為 warning。
5. **Spike 模板最小集** — `plan_kind: spike` 是否可只有 §Goal + §Acceptance + §結果回寫，免 Phase 0 公版？傾向是，但結果必須回寫主計畫對應 phase。
6. **Sub-plan dependency 表達**（v2 新增）— 當 sub-plan C 依賴 sub-plan A 完成時，是否需要 `depends_on: [<sub-id>]` 欄位？**Deferred — promotion gate：至少 3 個真實案例**（自然發生 C-depends-on-A 情境）再討論。理由：一旦加 `depends_on` 就會引入 DAG 而非 Tree，topological sort / cycle detection / graph validation 等複雜度成倍。本 plan 目前治理的是 Tree，不是 DAG；不為了通用而通用。

---

## 完成條件

- [ ] Phase 0 preflight 全部完成
- [ ] 4 個 sub-plan 全部 `status: completed` 且 archived
- [ ] Phase 5 收尾項目全 checked
- [ ] `governance/plan-tree-hierarchy.md` 落地
- [ ] 5 個 validator 進 `hooks.go` registry（4 block + 1 warning），unit test pass
- [ ] `ai-skill plans tree` CLI 可用
- [ ] 至少 1 組真實 parent-child plan 用新格式跑完（可用本 plan 自身作為 dogfood）
- [ ] Failure pattern 文件化
- [ ] Glossary terms 註冊
- [ ] 5 platform binaries rebuilt + pushed
- [ ] 本 plan 整 folder 搬 `plans/archived/2026-06-02-1200-plan-tree-hierarchy-governance/`

---

## Stakeholder 同意項目

- [x] linyihong: 目錄結構 = folder + `_plan.md` + `NN-` 前綴作為 UI convention（hierarchy 真實來源為 frontmatter `parent`）— **sign-off 2026-06-02**
- [x] linyihong: 「何時開 sub-plan」改用 free-text `sub_plan_reason`（非空 validator）+ recommended trigger 表（不強制）— **sign-off 2026-06-02 v2 review**
- [x] linyihong: 棄用 `children:` / `completion_blocks_parent` / `sub_plan_trigger enum` / `max_depth: 2`；改用 `parent` / `required_for_completion` / `sub_plan_reason` / `warning_at_depth: 3`（minimal governance）— **sign-off 2026-06-02 v2 review**
- [ ] linyihong: Open Questions Q2/Q3/Q5/Q6 解法（待 Phase 0 提案後 sign-off）
- [ ] linyihong: 是否將 plan tree 結構 promote 為 ADR（completed 後評估）

---

## Glossary Impact

**Glossary Impact: yes**

新引入 framework vocabulary：

| Term | 定義 | 註冊目標 |
|---|---|---|
| `plan_kind` | plan 的類型 enum：`main` / `sub` / `spike` | `knowledge/glossary/ai-skill.md` |
| `parent` | sub-plan 指向主計畫 id 的 frontmatter pointer，是 plan tree hierarchy 的單一 source of truth | `knowledge/glossary/ai-skill.md` |
| `required_for_completion` | sub-plan 是否屬於 parent 的 acceptance criteria 的 boolean；validator 由此推導 archive blocking | `knowledge/glossary/ai-skill.md` |
| `sub_plan_reason` | sub-plan 為什麼存在的 free-text 說明（非空，但不審內容） | `knowledge/glossary/ai-skill.md` |
| `plan tree` | 主計畫 + sub-plan 構成的階層結構，由 frontmatter `parent` pointer 動態建立 | `knowledge/glossary/ai-skill.md` |

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
