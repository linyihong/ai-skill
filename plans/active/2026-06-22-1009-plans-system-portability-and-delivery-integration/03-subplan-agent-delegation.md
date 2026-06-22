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
  共識（新欄位算 portable 還是 Ai-skill-only？是否進 plan_profile？），
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
- **In**：sub-plan frontmatter 新增 **nested `delegation` 物件**（`enabled` / `modes` / `brief` / `constraints`）；定義 brief 必含內容（goal / acceptance / in-out scope / context pack 指標 / validation）；人工派發與 agent 派發共用同一份 brief 的契約。
- **Out**：自動 orchestrator（自動偵測 + 自動 spawn agent + 自動收斂結果）— 保留為 future；本輪採 reservation pattern（schema + 雙路徑契約，不建自動 orchestrator）。
- **Affected**：plan-tree frontmatter schema、`scripts/ai-skill-cli/internal/app/plan_tree.go`（`PlanFrontmatter` struct + 可選 `validatePlanTreeFrontmatter` 擴充）、`governance/lifecycle/plan-tree-hierarchy.md`、`plans/README.md`。

## Decision Rationale（sub 層）
sub-plan 已具獨立 acceptance / archive，是天然的委派單元，缺的是**自足 brief 契約**：執行者（人或 agent）不讀整個 main plan 也能獨立完成。

**Schema 修正（回應 review #5）：避免把「可委派」與「brief 存在」綁死。** 扁平的 `agent_assignable: bool` + `delegation_brief` 會讓「是否可派」「派給誰」「brief」混在一起，未來要支援 manual-only / agent-only / hybrid / forbidden 就得破 schema。改用 **nested `delegation` 物件**：

```yaml
delegation:
  enabled: true            # 是否開放委派（取代 agent_assignable）
  modes: [human, agent]    # 允許路徑：human / agent 任一或並列
  brief:                   # 自足 brief（tool-neutral）
    goal: ...
    scope: { in: [...], out: [...] }
    acceptance: ...
    context_pack: [...]    # 需讀哪些檔
    validation: ...
  constraints:             # 委派限制（worktree / 不可碰路徑 / 需 sign-off）
    - ...
```

雙路徑共用同一份 `brief`：人工派發把 brief 貼給另一開發 / session；agent 派發把 brief 餵給 Agent/Task 工具（可選 worktree isolation，細節歸 `constraints` 與 `ai-tools/`，brief 本身 tool-neutral，Q6）。

### Alternatives
- A. 純文件慣例（sub-plan 標「可指派」但無 schema）：partial — 使用者要兩路徑且要可靠 brief，純慣例不足。
- B. 扁平 `agent_assignable` + `delegation_brief`：reject — 隱藏耦合（assignable 綁 brief 存在），不支援 mode 變化，未來必破 schema（review #5）。
- C. 只做 agent 自動 orchestrator：reject — over-engineering，使用者要人工也能用。
- D. nested `delegation` 物件 + reservation（accept）。

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
| Q5 delegation 最小契約 | still-open | Phase 1 設計 nested `delegation`（enabled/modes/brief/constraints），避免 assignable↔brief 綁死 |
| Q6 tool-neutral vs 綁工具 | still-open | Phase 1 brief tool-neutral，工具細節歸 constraints/ai-tools |

### Phase 0.1 — 架構盤點（需與 01 對齊 frontmatter schema；**不依賴外部 repo 能力**）
- [ ] 讀 `plan_tree.go` `PlanFrontmatter` struct + `validatePlanTreeFrontmatter`，確認新增 optional 欄位不破壞既有 5 validators。
- [ ] 與 01 對齊：delegation 欄位是否進 `plan_profile`（外部 repo 也能用），還是 Ai-skill-only。
- [ ] 確認新欄位為 **optional**（不破壞既有 sub-plan，未宣告 = 不可委派 / 不變）。

## Phase 1 — Delegation schema + brief 契約設計
- [ ] 定義 nested `delegation` 物件（optional，未宣告 = 不可委派 / 行為不變）：`enabled` / `modes:[human|agent]` / `brief` / `constraints`。
- [ ] 定義 `brief` 必含：goal / in-out scope / acceptance / context_pack（需讀哪些檔）/ validation / 不依賴未完成 sibling 的聲明。
- [ ] tool-neutral：`brief` 不綁特定 agent 工具；工具/隔離細節歸 `constraints` + `ai-tools/`（Q6）。
- [ ] 文件化於 `governance/lifecycle/plan-tree-hierarchy.md` + `plans/README.md`。

## Phase 2 — Validator 擴充（optional 欄位）+ 雙路徑說明
- [ ] 擴充 `validatePlanTreeFrontmatter`：當 `delegation.enabled: true` 時，驗 `delegation.modes` 非空且值合法、`delegation.brief` 必要子欄位非空（block）；未宣告 `delegation` 則不變（向後相容）。
- [ ] 測試：tmp fixture（enabled+完整 brief pass / enabled+缺 brief fail / enabled+空 modes fail / 未宣告 pass）。
- [ ] 文件化雙路徑：human 派發 SOP + agent 派發 SOP（後者可選 worktree isolation；保持 tool-neutral，工具細節放 `ai-tools/`）。
- [ ] **若擴充 validator 行為，補 Runtime Execution Path trigger flow**（commit-msg validator 已是既有 dispatch，新增子驗證須宣告）。

## Phase 3 — Dogfood（回應 review #6：兩路徑各一次）
- [ ] 挑一個真實 sub-plan 設 `delegation.enabled: true` + 完整 brief。
- [ ] **human 路徑 evidence**：另一 session / 開發僅憑 brief 獨立完成一次。
- [ ] **agent 路徑 evidence**：以 Agent/Task 工具僅憑 brief 在 worktree 執行一次。
- [ ] 兩次皆記錄 brief 是否足夠自足（缺漏回饋修正 schema）。

## 完成條件
- [ ] nested `delegation` schema + brief 契約落地（Q5 resolved）
- [ ] 欄位 optional、向後相容（既有 sub-plan 不受影響）
- [ ] validator 擴充 + 測試通過（含 enabled/modes/brief 各 violation）
- [ ] human + agent 雙路徑 SOP 落地（tool-neutral，Q6 resolved）
- [ ] dogfood evidence：human 一次 + agent 一次
- [ ] 與 01 `plan_profile` 邊界對齊（`delegation` 欄位 portable 與否已決定）

## Glossary Impact
Glossary Impact: yes — 新增 `delegation`（nested 委派 schema：enabled/modes/brief/constraints）；Phase 1 落地時註冊到 `knowledge/glossary/ai-skill.md`。取代早期扁平 `agent_assignable` / `delegation_brief` 提案。

## 與其他 plans 的關係
- 擴充 [`archived/2026-06-02-1200-plan-tree-hierarchy-governance/_plan.md`](../../archived/2026-06-02-1200-plan-tree-hierarchy-governance/_plan.md) 的 frontmatter schema 與 `validatePlanTreeFrontmatter`。
- 依賴 [`01-external-repo-plan-system-shared-binary.md`](01-external-repo-plan-system-shared-binary.md) 的 `plan_profile` 邊界決定新欄位是否 portable。
