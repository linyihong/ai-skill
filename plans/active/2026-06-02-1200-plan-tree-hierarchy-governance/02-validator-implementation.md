---
id: 2026-06-02-1200-plan-tree-validator-implementation
plan_kind: sub
status: in-progress
owner: linyihong
created: 2026-06-03
parent: 2026-06-02-1200-plan-tree-hierarchy-governance
required_for_completion: true
sub_plan_reason: >
  5 個 plan-tree validator 是本 governance plan 從「文件化規則」轉成「機械化
  enforcement」的關鍵階段。獨立成 sub-plan 以便 (a) hooks.go 改動範圍清楚
  bounded、(b) 與 Phase 1 schema doc 並行安全（schema 先 freeze、validator 後寫）、
  (c) 與 mechanical-enforcement-registry 的 commit-msg dispatcher 在同一個層級
  整合，risk 集中在這個 sub-plan 內部。Phase 2 啟動條件是 registry plan archive
  完成（避免 hooks.go 同時改動造成衝突），該條件於 2026-06-03 commit 5b3e089
  滿足。
---

# Plan Tree Validator Implementation（sub-plan）

> 本 sub-plan 繼承 parent 的 Decision Rationale（minimal governance）。

## Purpose

實作 parent `_plan.md` §Phase 2 列出的 5 個 plan-tree validator，使
plan tree hierarchy 從「靠 prose 自律」升級為「commit-msg + compile-time
機械擋」。實作完成後，sub-plan 漏 `parent` / `sub_plan_reason` /
`required_for_completion` 任一欄、parent 指向不存在 id、duplicate id 與
main archive 時 required child 未 completed 都會直接 reject commit。

本 sub-plan 同時補完 Phase 1 留下的 negative fixture 集合，作為 validator
unit test 的 testdata 錨點。

## Scope

### In scope

- 5 個 Go validator：
  - `validatePlanTreeFrontmatter`（block）— sub/spike 缺 `parent`、`sub_plan_reason` 空字串、缺 `required_for_completion` → 任一觸發 reject
  - `validatePlanTreeArchiveOrder`（block）— `git mv` 把 `plans/active/<main>/_plan.md` 搬到 `plans/archived/` 時，所有 `parent == <main>.id` 且 `required_for_completion: true` 的 sub-plan（不論 active / archived 位置）必須 `status: completed`
  - `validatePlanTreeParentReference`（block）— sub/spike `parent` 指向的 id 必須在全 repo plan 集合（active + archived）中可解析
  - `validatePlanTreeUniqueID`（block）— 全 repo `id` frontmatter 必須唯一；新增 / 改 id 撞既有 → reject
  - `validatePlanTreeFolderConvention`（warning only）— folder 缺 `_plan.md`、檔名不符 `^\d{2}-` 前綴、或 plan tree 深度 ≥ 3 → 印 warning，不擋
- Negative fixture 集合（在本 sub-plan folder 下 `fixtures/`，重用 Phase 1 happy-path 命名慣例）：
  - `sub-missing-parent.md` — sub 缺 `parent`
  - `sub-empty-reason.md` — `sub_plan_reason: ""`
  - `sub-missing-required.md` — sub 缺 `required_for_completion`
  - `parent-orphan.md` — parent 指向不存在的 id
  - `duplicate-id-a.md` + `duplicate-id-b.md` — 兩個 plan 共用同一 id
  - `archive-required-pending.md` — 模擬 main archive 但 required child 仍 in-progress
  - `folder-shape-violations/` — folder 缺 `_plan.md`、檔名違規、深度 3 範例
- Unit tests：每個 validator 至少 1 fail + 1 pass + 1 edge case，全部使用 `fixtures/` 下的 testdata
- `runtime/core-bootstrap.yaml` `per_commit_obligations[]` 新增 5 條 obligation：
  - `obligation.commit.plan_tree_frontmatter`（severity: block）
  - `obligation.commit.plan_tree_archive_order`（severity: block）
  - `obligation.commit.plan_tree_parent_reference`（severity: block）
  - `obligation.commit.plan_tree_unique_id`（severity: block）
  - `obligation.commit.plan_tree_folder_convention`（severity: warning）
- 每條 obligation 對應 opt-out trailer `[skip-plan-tree-<short>]`
- `commitMsgValidatorRegistry` map 註冊 5 條新 entry
- `defaultCommitMsgDispatchOrder` 同步加入 5 條（保留 fresh-clone fallback 行為）

### Out of scope

- `ai-skill plans tree` CLI subcommand — 屬於 Phase 3 `03-cli-tree-subcommand`
- `04-existing-plan-migration` — 屬於 Phase 4
- 舊 plan 全面 migration — Phase 2 validator land 後，舊 plan 仍是「無 frontmatter
  即不檢查」的兼容路徑（validator 對 frontmatter 缺失的 main plan 採 skip
  而非 reject，避免一次性 break 全 repo；migration drive 在 Phase 4 處理）
- ADR 升級 — 等 Phase 4 完成 ≥ 3 個 cluster migration 後評估
- Spike kind 特殊處理 — 目前 validator 對 spike 與 sub 採同邏輯（spike 預設
  `required_for_completion: false` 但仍要 `parent` / `sub_plan_reason`）

## Acceptance Criteria

- [ ] 5 個 Go validator 完成 + 全部 unit test PASS
- [ ] Negative fixture 集合落地（7+ 個 fixture file 或目錄）
- [ ] `runtime/core-bootstrap.yaml` 新增 5 條 `per_commit_obligations[]` 並通過 `ai-skill runtime compile`
- [ ] `commitMsgValidatorRegistry` 註冊 5 條，`defaultCommitMsgDispatchOrder` 同步
- [ ] `ai-skill runtime receipt` 輸出 `obligations=28 gates=25`（或對應的更新後總數）
- [ ] 既有 11+ commit-msg validator 全部仍 PASS（不破壞既存）
- [ ] 5 platform binaries rebuild + BUILDINFO + SHA256SUMS 更新
- [ ] parent `_plan.md` Phase 2 區塊 status 更新（in-progress → completed），同步勾選驗證要點
- [ ] 本 sub-plan archive（移到 `plans/archived/`）

## Runtime Impact

- **新增 5 條 commit-msg obligation** → `runtime.db.obligations` count 從 23 → 28
- **新增 5 條 commit-msg validator dispatch** → `commitMsgValidatorRegistry` map 與 `defaultCommitMsgDispatchOrder` slice 各加 5 entry
- 不新增 `routing-registry` route（plan-tree validator 不需 PreToolUse activation）
- 不新增 `generated_surfaces`（直接讀 frontmatter，無需 projection）
- 不新增 PreToolUse / PreCommit / Stop hook（純 commit-msg 階段執行）
- `enforcement-registry.yaml` 必須補：新增 rule_class `plan_tree_governance`
  coverage=`mechanical`，executors[] 列 5 條 symbol；同 commit 落地以滿足
  `validateEnforcementRuleRegistrySync`

## Implementation Order

1. 先寫 fixtures（讓 test data 落地）
2. 寫 Go validator + unit test（test-driven，紅→綠）
3. 接 `commitMsgValidatorRegistry` + `defaultCommitMsgDispatchOrder`
4. 更新 `runtime/core-bootstrap.yaml` `per_commit_obligations[]` + 跑 `ai-skill runtime compile`
5. 更新 `enforcement-registry.yaml` 新增 `plan_tree_governance` rule_class
6. 第一個 commit：`feat(plan-tree): implement 5 plan-tree validators + fixtures + obligations`
7. Rebuild 5 platform binaries
8. 第二個 commit：`chore(bin): rebuild 5 platform binaries from <feat-commit-sha>`
9. Push 兩個 commit
10. Readback：跑 `ai-skill runtime receipt`、`ai-skill enforcement coverage --self-check`、`ai-skill hooks commit-msg <test-message>` 驗證新 validator 有效
11. 更新 parent `_plan.md` Phase 2 status，archive 本 sub-plan
12. 第三個 commit：`docs(plans): plan-tree validator-implementation sub-plan close-out + archive`

## 與 parent 的同步

本 sub-plan 完成時：
- 在 parent `_plan.md` §Phase 2 標 status: completed，5 條 validator 驗證要點全部勾選
- 連動 `enforcement-registry.yaml`（新增 `plan_tree_governance` rule_class）
- 連動 `runtime/core-bootstrap.yaml`（5 條新 obligation）
- 連動 `governance/lifecycle/plan-tree-hierarchy.md` rule（若 Phase 1 已落地）
- Receipt 與 coverage CLI 數值更新（obligations / classes 計數）

## Source

Parent: `plans/active/2026-06-02-1200-plan-tree-hierarchy-governance/_plan.md` §Phase 2
Schema 依據: `plans/active/2026-06-02-1200-plan-tree-hierarchy-governance/01-frontmatter-schema.md`
Unblocker: 2026-06-03 commit 5b3e089（mechanical-enforcement-registry plan archive，解除 hooks.go 並行衝突）
