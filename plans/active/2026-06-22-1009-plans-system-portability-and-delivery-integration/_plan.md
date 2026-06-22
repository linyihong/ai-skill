---
id: 2026-06-22-1009-plans-system-portability-and-delivery-integration
plan_kind: main
status: draft
owner: linyihong
created: 2026-06-22
parent: null
---

# Plans System Portability & Delivery Integration（plans 系統外部化與交付接入）

**Status**: `draft`
**Owner**: linyihong
**建立日期**: 2026-06-22
**Source**: 2026-06-22 對話 — 使用者觀察 plans 系統（plan-tree + 驗證機制 + 子計畫）已足夠成熟，希望 (1) 讓外部 repo 也能使用、(2) 把 software-delivery 接入「開發前先寫 plan」、(3) 用子計畫系統讓一個 sub-plan 可交給其他 agent 執行。
**Glossary Impact**: yes — 預期新引入 framework vocabulary：`plan_system_profile`（portable core vs governance overlay 邊界）、`agent_assignable` / `delegation_brief`（子計畫委派）。落地前須註冊到 `knowledge/glossary/ai-skill.md`；本 main plan 階段先宣告，sub-plan graduate 時才註冊。

> **Watch-Out List citation**：本 plan 對應 [`architecture/ai-native-cognitive-ecosystem-system.md`](../../../architecture/ai-native-cognitive-ecosystem-system.md) §Watch-Out List 的「process bloat」「premature abstraction」「over-engineering」防呆 — plan-first 採 advisory workflow ordering 而非機械 block（避免誤擋小修補）；委派只先做 schema + 雙路徑（人工 / agent），不建自動 orchestrator；外部化採共用 binary 而非整套治理搬遷，避免外部 repo 被 Ai-skill governance 綁架。

---

## Decision Rationale

### Problem & Why Now

plans 系統目前的價值（plan-tree 階層、frontmatter 單一 source-of-truth、5 個 commit-msg validators、archival audit / link integrity、pre-build interrogation → preflight → completion closure）都封在 Ai-skill repo 內部，且與 Ai-skill 專屬治理（`runtime.db`、`generated_surfaces`、routing-registry、Glossary、Cognitive Modes、ADR pipeline、intelligence atoms）混在同一份 plan 模板。三件事同時成熟到值得動手：

1. **外部 repo 想用**：`ai-skill plans tree --root <path>` 已支援任意 repo root（讀取面已跨 repo），但 commit-time 強制（validators）仍是本 repo 的 git hook，外部 repo 無法享受。
2. **software-delivery 缺 plan-first ordering**：現有 intake 有 pre-build-interrogation / Test-First Ordering，但沒有明文要求「實作前先有 plan」的順序，導致 plan 與 delivery 各走各的。
3. **子計畫已有獨立 acceptance / archive，但沒有委派出口**：sub-plan 可獨立 sign-off，卻沒有「自足 brief + 可交給其他 agent（人工或 Agent 工具）執行」的 schema。

### Decision

開一個 **plan tree（main + 3 sub-plan）**，dogfood plan-tree 本身，分階段 graduate：

- **Sub-plan 01 — External-repo plan system via shared binary**：外部 repo 透過**共用 ai-skill binary**（非 init-project 抽取安裝）使用 plans 系統。核心缺口是 commit-time 強制跨 repo：設計外部 repo 用共用 binary 跑 plan-tree / archival validators 的路徑（薄 git hook shim 或 `ai-skill plans validate --root` 子命令）。同時釐清 **portable core vs Ai-skill governance overlay** 邊界（`plan_system_profile`）。
- **Sub-plan 02 — software-delivery plan-first ordering**：把 plan-first 寫進 `workflow/software-delivery/` intake 段，接在 pre-build-interrogation / Test-First Ordering 之後，**advisory + review 檢查，不做機械 block**。
- **Sub-plan 03 — Sub-plan agent delegation**：在 sub-plan frontmatter 加 delegation 欄位（`agent_assignable` + 自足 `delegation_brief`），**同時支援人工派發與 Agent/Task 工具派發**（依專案需求二選一或並用），執行端兩條路都通。

### Alternatives Considered

- **A. 維持現狀（plans 只在 Ai-skill 內用）**：reject — 使用者明確要外部化，且讀取面已跨 repo，差最後一哩 commit-time 強制。
- **B. init-project 把整套 plan 模板抽取安裝進外部 repo**：reject（本輪）— 使用者選「共用 binary 指向外部 repo」；整套抽取會把 Ai-skill governance overlay 一起帶過去，違反 portability 初衷且工作量最大。保留為 01 的 future option。
- **C. plan-first 做成硬機械 gate（無 active plan 不准 commit code）**：reject（本輪）— 使用者選 workflow 層 ordering；硬 gate 易誤擋小修補、需複雜 opt-out。保留為 02 的後續升級候選（maturity ladder：ordering 觀察 → 視 evidence 再升級）。
- **D. 一次做完三件**：reject — 三件可分開 graduate，分階段降低 blast radius，且 plan-tree 本身就是用來表達這種 main↔sub 拆分。

### Why Not an ADR Yet

scope 仍會隨 Phase 0 盤點調整（特別是 01 的跨 repo 強制機制、03 的 delegation schema 形狀）；多個 Open Questions 未解；可能有更輕的 promotion target（例如 02 可能只需 workflow doc 更新，不需任何 runtime contract）。三條 sub-plan 各自 graduate 後，若浮現需跨 session / 跨 project 固化的決策，再評估 ADR。

### ADR Promotion Criteria（completed 時驗證）

- [ ] foundational + cross-session + cross-project + expensive-to-reverse + explains-why 全中
- [ ] 三條 sub-plan 結果證實 portable core 邊界、plan-first ordering、delegation schema 三者可行且被真實使用
- [ ] Open Questions 全解
- [ ] 沒有更輕的 promotion target 適用（per ADR-007）
- [ ] 至少一個外部 repo 真實透過共用 binary 跑過 plan validate（具體 evidence）

### Consequences（預期）

#### 正面
- 外部 repo 不需 fork Ai-skill governance 就能用成熟的 plan-tree + 驗證閉環。
- software-delivery 與 plans 系統對齊，「先規劃後實作」成為可見的 intake 順序。
- 子計畫成為可委派的執行單元，支援多 agent / 人工混合交付。

#### 負面
- 共用 binary 路徑讓外部 repo 對 Ai-skill binary 版本產生依賴（需版本相容策略）。
- delegation schema 增加 sub-plan frontmatter 表面積。

#### 風險
- portable core 與 governance overlay 邊界若沒切乾淨，外部 repo 仍會踩到 Ai-skill 專屬 validator（如 runtime trigger wiring）→ 01 Phase 0 必須先把「哪些 validator 屬 portable core」列清楚。
- plan-first ordering 若沒接好既有 pre-build-interrogation，會變成重複 gate（process bloat）。

---

## Runtime Execution Path

本 main plan 為 **規劃容器**，自身不直接接入 runtime；具體 runtime / workflow / validation 改動由各 sub-plan 宣告：

- **01** 可能新增 `ai-skill plans validate --root` 子命令（CLI consumer）與外部 repo git hook shim；若新增 commit-msg validator 行為或 `route.*`，須在 01 自己的 Runtime Execution Path + Per-surface consumer 表宣告，並由 `validateRuntimeTriggerWiring` 機械驗證。
- **02** 預期僅改 `workflow/software-delivery/` 文件（intake ordering）+ 可能新增 validation scenario；若不接 runtime，須在 02 明寫 doc-only + 未來接入條件。
- **03** 預期改 plan-tree frontmatter schema + 可能擴充 `validatePlanTreeFrontmatter`（新增 optional delegation 欄位驗證）；若擴充 validator 須宣告 trigger flow。

**doc-only trial 宣告**：本 main plan `status: draft` 階段不接入 runtime；接入由 sub-plan graduate 時各自落地。

---

## Open Questions

| # | Question | 處置 | 歸屬 sub-plan |
|---|----------|------|---------------|
| 1 | 外部 repo 跑 plan-tree validators 的最薄機制是什麼（git hook shim vs `plans validate` 子命令 vs 兩者）？ | still-open | 01 |
| 2 | portable core 到底包含哪幾個 validator？哪些是 Ai-skill governance overlay 必須排除？ | still-open | 01 |
| 3 | 共用 binary 的版本相容策略（外部 repo pin 哪個 binary、frontmatter schema 版本怎麼宣告）？ | still-open | 01 |
| 4 | plan-first ordering 與既有 pre-build-interrogation / Architecture Compatibility Preflight 如何不重複？ | still-open | 02 |
| 5 | delegation schema 欄位最小集合（`agent_assignable` / `delegation_brief` / context pack 指標）？人工與 agent 兩路徑共用同一份 brief 是否可行？ | still-open | 03 |
| 6 | agent 派發是否綁定特定工具（Task/Agent / worktree isolation），或保持 tool-neutral 只定義 brief 契約？ | still-open | 03 |

---

## 完成條件

- [ ] 三條 sub-plan（01 / 02 / 03）皆 `status: completed`
- [ ] 各 sub-plan 的 acceptance 達成並通過其宣告的 validation
- [ ] Open Questions 全部標記 `resolved` / `deferred`（附原因）並回寫
- [ ] Glossary Impact 落實：新 vocabulary 已註冊或明確不註冊
- [ ] 執行 Plan Completion Closure（含 `ai-skill runtime refresh` 若涉 knowledge/validation 層）
- [ ] plan tree 通過 `ai-skill plans tree` 檢視（main + 3 sub 階層正確）

## Phase 0 — Pre-Build Interrogation

### Phase 0.0 — Open Questions 核對（公版，必填）

逐條核對本 plan §Open Questions，標記處置並回寫：

- [ ] 已讀本 plan §Open Questions 全部條目
- [ ] 對每條標記 `resolved`（附 Phase 0 證據）/ `still-open` / `deferred`（附原因）
- [ ] `resolved` 的條目已同步勾選 / 附註於 §Open Questions
- [ ] 若盤點新發現問題，已加入 §Open Questions

| Open Question | 處置 | 證據 / 原因 |
|---|---|---|
| Q1-Q6 | still-open | 待各 sub-plan Phase 0 盤點解決 |

### Phase 0.1 — 架構相容性 preflight（main 層）

- [ ] 確認 `plans/README.md` Plan Tree Hierarchy 規則與本 tree frontmatter 相容
- [ ] 確認 `scripts/ai-skill-cli/internal/app/plan_tree.go` / `plans.go` 現行 schema（避免 03 改 frontmatter 撞既有 validator）
- [ ] 確認 `workflow/software-delivery/intake.md` 現行 intake 順序（02 接入點）
- [ ] 確認 `governance/lifecycle/plan-tree-hierarchy.md` 治理規則

## Phase 1-N

實作步驟拆進三條 sub-plan，各自有 Phase 0 / Phase 1-N / 完成條件：

| Sub-plan | 主題 | required_for_completion | 建議 sequencing |
|---|---|---|---|
| [`01-external-repo-plan-system-shared-binary.md`](01-external-repo-plan-system-shared-binary.md) | 外部 repo 經共用 binary 使用 plans 系統 + portable core 邊界 | true | 先做（解 Q2 portable core 邊界，後兩條依賴此邊界） |
| [`02-software-delivery-plan-first-ordering.md`](02-software-delivery-plan-first-ordering.md) | software-delivery plan-first workflow ordering（advisory） | true | 可與 01 並行（純 workflow doc） |
| [`03-subplan-agent-delegation.md`](03-subplan-agent-delegation.md) | sub-plan 委派 schema（人工 + agent 雙路徑） | true | 排在 01 之後（依賴 frontmatter schema 共識） |

## Stakeholder 同意項目

- [ ] 外部化採「共用 binary」而非 init-project 抽取（使用者 2026-06-22 已選）
- [ ] plan-first 採 workflow ordering advisory 而非機械 block（使用者 2026-06-22 已選）
- [ ] 委派同時支援人工與 agent 雙路徑（使用者 2026-06-22 已選）
- [ ] 分階段交付，dogfood plan tree（使用者 2026-06-22 已選）

## Glossary Impact

Glossary Impact: yes — 預期新增 `plan_system_profile`、`agent_assignable`、`delegation_brief`；本 main plan 階段僅宣告，實際註冊由 sub-plan graduate 時落地（避免提前註冊未定稿術語）。

## 與其他 plans 的關係

- [`archived/2026-06-02-1200-plan-tree-hierarchy-governance/_plan.md`](../../archived/2026-06-02-1200-plan-tree-hierarchy-governance/_plan.md) — 本 tree 直接建立於其 frontmatter schema + 5 validators 之上；03 將擴充其 frontmatter。
- [`archived/2026-05-28-1830-plan-archival-audit-validator.md`](../../archived/2026-05-28-1830-plan-archival-audit-validator.md) — archival audit validator，01 portable core 候選成員。
- [`workflow/software-delivery/intake.md`](../../../workflow/software-delivery/intake.md) — 02 的接入點。
