---
id: 2026-06-22-1009-subplan-agent-delegation
plan_kind: sub
status: draft
owner: linyihong
created: 2026-06-22
parent: 2026-06-22-1009-plans-system-portability-and-delivery-integration
required_for_completion: true
sub_plan_reason: >
  委派 schema 會擴充 plan-tree frontmatter，依賴 01 對 portable core 邊界的
  共識（新欄位算 portable 還是 Ai-skill-only？是否進 plan_system_profile？），
  因此排在 01 之後。獨立成 sub-plan 以便 delegation schema + 人工/agent 雙路徑
  契約可獨立設計與 sign-off，並可先 reservation（只設計 schema 不實作自動派發）。
---

# Sub-plan Agent Delegation（sub-plan）

**Status**: `draft`
**Owner**: linyihong
**Parent**: [`_plan.md`](_plan.md)

## Source Request
用子計畫系統讓一個 sub-plan 可交給其他 agent 執行。使用者澄清：**人工派發與 agent 派發兩種都要支援，依專案需求選用或並用**。

## Scope
- **In**：sub-plan frontmatter 新增 delegation 欄位（`agent_assignable` + 自足 `delegation_brief` 契約）；定義 brief 必含內容（goal / acceptance / in-out scope / context pack 指標 / validation）；人工派發與 agent 派發共用同一份 brief 的契約。
- **Out**：自動 orchestrator（自動偵測 + 自動 spawn agent + 自動收斂結果）— 保留為 future；本輪採 reservation pattern（schema + 雙路徑契約，不建自動 orchestrator）。
- **Affected**：plan-tree frontmatter schema、`scripts/ai-skill-cli/internal/app/plan_tree.go`（`PlanFrontmatter` struct + 可選 `validatePlanTreeFrontmatter` 擴充）、`governance/lifecycle/plan-tree-hierarchy.md`、`plans/README.md`。

## Decision Rationale（sub 層）
sub-plan 已具獨立 acceptance / archive，是天然的委派單元，缺的是**自足 brief 契約**：執行者（人或 agent）不讀整個 main plan 也能獨立完成。雙路徑共用同一份 `delegation_brief`：
- **人工派發**：把 brief 貼給另一位開發 / 另一個 session。
- **agent 派發**：把 brief 餵給 Agent/Task 工具（可選 worktree isolation）。
brief 契約保持 **tool-neutral**（Q6）：只定義「自足執行所需資訊」，不綁特定 agent 工具。`agent_assignable: bool` 標示此 sub-plan 是否設計為可獨立委派（有完整 brief + 不依賴未完成的 sibling）。

### Alternatives
- A. 純文件慣例（sub-plan 標「可指派」但無 schema）：partial — 使用者要兩路徑且要可靠 brief，純慣例不足。
- B. 只做 agent 自動 orchestrator：reject — over-engineering，使用者要人工也能用。
- C. schema + 雙路徑契約 + reservation（accept）。

## Open Questions（本 sub）
- Q5（delegation schema 最小欄位集；人工/agent 共用同一 brief 是否可行）。
- Q6（agent 派發是否綁工具 vs tool-neutral brief 契約）。

## Phase 0 — Pre-Build Interrogation

### Phase 0.0 — Open Questions 核對（公版，必填）
- [ ] 已讀 main + 本 sub §Open Questions 全部條目
- [ ] 對每條標記 `resolved` / `still-open` / `deferred`
- [ ] `resolved` 條目回寫
- [ ] 新問題已加入 §Open Questions

| Open Question | 處置 | 證據 / 原因 |
|---|---|---|
| Q5 schema 最小欄位 | still-open | Phase 1 設計 brief 契約 |
| Q6 tool-neutral vs 綁工具 | still-open | Phase 1 決定 brief 不綁工具 |

### Phase 0.1 — 架構盤點（**依賴 01 portable core 邊界**）
- [ ] 讀 `plan_tree.go` `PlanFrontmatter` struct + `validatePlanTreeFrontmatter`，確認新增 optional 欄位不破壞既有 5 validators。
- [ ] 與 01 對齊：delegation 欄位是否進 `plan_system_profile`（外部 repo 也能用），還是 Ai-skill-only。
- [ ] 確認新欄位為 **optional**（不破壞既有 sub-plan，未宣告 = 不可委派 / 不變）。

## Phase 1 — Delegation schema + brief 契約設計
- [ ] 定義 frontmatter 新欄位：`agent_assignable: bool`（optional，default false）；`delegation_brief` 指標（指向本檔某 section 或獨立 brief 區塊）。
- [ ] 定義 `delegation_brief` 必含：goal / in-out scope / acceptance / context pack（需讀哪些檔）/ validation / 不依賴未完成 sibling 的聲明。
- [ ] tool-neutral：brief 不綁特定 agent 工具（Q6 → tool-neutral）。
- [ ] 文件化於 `governance/lifecycle/plan-tree-hierarchy.md` + `plans/README.md`。

## Phase 2 — Validator 擴充（optional 欄位）+ 雙路徑說明
- [ ] 擴充 `validatePlanTreeFrontmatter`：當 `agent_assignable: true` 時，驗 `delegation_brief` 必要子欄位非空（block）；未宣告則不變（向後相容）。
- [ ] 測試：tmp fixture（assignable+完整 brief pass / assignable+缺 brief fail / 未宣告 pass）。
- [ ] 文件化雙路徑：人工派發 SOP + agent 派發 SOP（後者可選 worktree isolation；保持 tool-neutral，工具細節放 `ai-tools/`）。
- [ ] **若擴充 validator 行為，補 Runtime Execution Path trigger flow**（commit-msg validator 已是既有 dispatch，新增子驗證須宣告）。

## Phase 3 — Dogfood
- [ ] 挑一個未來真實 sub-plan 標 `agent_assignable: true` + 寫完整 brief，驗證人工或 agent 能僅憑 brief 獨立執行（acceptance evidence）。

## 完成條件
- [ ] delegation frontmatter schema + brief 契約落地（Q5 resolved）
- [ ] 欄位 optional、向後相容（既有 sub-plan 不受影響）
- [ ] validator 擴充 + 測試通過
- [ ] 人工 + agent 雙路徑 SOP 落地（tool-neutral，Q6 resolved）
- [ ] 至少一次 dogfood brief 委派 evidence
- [ ] 與 01 `plan_system_profile` 邊界對齊（新欄位 portable 與否已決定）

## Glossary Impact
Glossary Impact: yes — 新增 `agent_assignable`、`delegation_brief`；Phase 1 落地時註冊到 `knowledge/glossary/ai-skill.md`。

## 與其他 plans 的關係
- 擴充 [`archived/2026-06-02-1200-plan-tree-hierarchy-governance/_plan.md`](../../archived/2026-06-02-1200-plan-tree-hierarchy-governance/_plan.md) 的 frontmatter schema 與 `validatePlanTreeFrontmatter`。
- 依賴 [`01-external-repo-plan-system-shared-binary.md`](01-external-repo-plan-system-shared-binary.md) 的 `plan_system_profile` 邊界決定新欄位是否 portable。
