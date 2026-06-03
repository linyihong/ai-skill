---
id: 2026-06-02-1200-plan-tree-existing-plan-migration
plan_kind: sub
status: in-progress
owner: linyihong
created: 2026-06-03
parent: 2026-06-02-1200-plan-tree-hierarchy-governance
required_for_completion: true
sub_plan_reason: >
  Phase 4 dogfood — 把現有 plans/active 與 plans/archived 內已有 implicit
  parent ↔ child 關係（純 prose 寫在檔頭）的 plan 群，遷移到新 frontmatter
  schema 並驗證 ai-skill plans tree 能正確渲染。獨立成 sub-plan 是因為
  遷移是一次性 batch 操作，跟「持續性 governance enforcement」(Phase 2) 與
  「visualization tool」(Phase 3) 責任邊界不同。
---

# Existing Plan Migration（sub-plan）

> 本 sub-plan 繼承 parent 的 Decision Rationale（minimal governance）。

## Purpose

把 `plans/active/` 與 `plans/archived/` 內既有的 implicit parent-child 群組
遷移為新的 frontmatter schema（`id` / `plan_kind` / `parent` /
`required_for_completion` / `sub_plan_reason`），讓 `ai-skill plans tree`
能正確顯示整個 repo 的 plan hierarchy。

驗證標準：遷移後 `ai-skill plans tree --root . --state all` 顯示 ≥ 1
parent-child cluster，且新 Phase 2 validators（particularly
`validatePlanTreeParentReference`）對遷移後狀態 PASS。

## Scope

### In scope（本 sub-plan 處理）

**Cluster 1 — Registry tree（首選 dogfood，per parent _plan.md §Phase 0 Inventory）**:
- Parent (archived): `plans/archived/2026-05-31-2100-mechanical-enforcement-registry.md`
  → `id: 2026-05-31-2100-mechanical-enforcement-registry`, `plan_kind: main`
- Child (active): `plans/active/2026-05-31-1900-workflow-activation-engine.md`
  → `parent: 2026-05-31-2100-...`, `plan_kind: sub`, `required_for_completion: true`
- Child (active): `plans/active/2026-05-31-2000-mechanical-sanitization-validator.md`
  → `parent: 2026-05-31-2100-...`, `plan_kind: sub`, `required_for_completion: true`
- Child (active): `plans/active/2026-06-01-0100-validation-scenario-governance-executor.md`
  → `parent: 2026-05-31-2100-...`, `plan_kind: sub`, `required_for_completion: false`
    （per Phase 3 round-2 R3 deferral；stub plan，非 parent acceptance gate）

**Standalones — 加最小 frontmatter 讓樹可見**:
- `plans/active/2026-05-27-1557-tool-runtime-signal-economics-integration.md`
  → `plan_kind: main`, `parent: null`, no required_for_completion
- `plans/active/2026-05-28-1636-gen4-fitness-optimization-memory-interface-reservation.md`
  → `plan_kind: main`, `parent: null`, no required_for_completion

**Other（plan-tree-hierarchy-governance cluster）**: 已 dogfood 完成，本 sub-plan 不重做。

### Out of scope

- 全 `plans/archived/` 歷史 plans 一次性遷移 — 只做 registry tree 一組以驗證
  pattern；其餘 archived plans 留作後續逐步遷移（plans 沒 frontmatter
  validators 會 silently skip，無功能性影響）
- 大規模 plan rename / id 標準化 — 維持既有 path 與檔名
- 已 archived 但有 implicit parent-child 結構的舊 cluster（若存在）— 暫不處理；
  active cluster 才是「會被 archive_ready 邏輯影響」的場景

## Acceptance Criteria

- [ ] Archived registry main plan 加 frontmatter（id + plan_kind: main）
- [ ] 3 個 active 子 plan 加 frontmatter（parent 全部指向 registry main）
- [ ] 2 個 standalone plan 加 frontmatter（plan_kind: main, parent: null）
- [ ] `ai-skill plans tree --root . --state all` 顯示 ≥ 3 個 root（registry main + 2 standalones + plan-tree-hierarchy-governance）
- [ ] `validatePlanTreeFrontmatter` / `validatePlanTreeParentReference` / `validatePlanTreeUniqueID` 在 commit 階段對遷移後 staged set 全 PASS（commit 過 commit-msg hook 即驗證）
- [ ] `ai-skill enforcement lint` 仍 0 fail
- [ ] parent `_plan.md` Phase 4 區塊 status 更新（待建 → completed），同步勾選驗證要點

## Runtime Impact

- 純 metadata 變更：6 個 plan 檔頭加 frontmatter，0 行 Go code
- 不重建 binaries
- 不改 runtime.db 或 generated_surfaces（frontmatter 是 plan-local metadata，walker 讀檔即可）
- 不新增 obligation / rule_class
- Phase 2 validators 既已 active，本 commit 是新 frontmatter 的 first
  large-scale dogfood — 若 schema 設計有問題會在此暴露

## Migration Order（嚴格）

**單一 commit 完成 atomic migration**：

1. 先 add frontmatter 到 archived parent（id 必須存在才 unblock subs 的 parent reference）
2. 同 commit add frontmatter 到 3 active subs（parent 指向 #1）
3. 同 commit add frontmatter 到 2 standalones（無 parent dependency）
4. 同 commit 完成後 push（避免 partial state）

如果分批 commit，會出現「subs 有 parent: X 但 X 尚未 migrate → validator block」的雞蛋問題。

## 與 parent 的同步

本 sub-plan 完成時：
- parent `_plan.md` §Phase 4 標 status: completed，驗證要點全勾選
- parent `_plan.md` 增 `ai-skill plans tree` 範例輸出（migration 後形狀 dogfood evidence）
- Phase 5 收尾可開始（CLI 已可用，migration 已 dogfood，governance/glossary/failure-pattern 收尾）

## Source

Parent: `plans/active/2026-06-02-1200-plan-tree-hierarchy-governance/_plan.md` §Phase 4
Cluster inventory: parent §Phase 0 Inventory 結果（Cluster 1 = Registry tree, 首選 dogfood）
Depends on: Phase 2（validators must be live to dogfood-validate the migration）+ Phase 3（CLI must be live to verify post-migration tree shape）
