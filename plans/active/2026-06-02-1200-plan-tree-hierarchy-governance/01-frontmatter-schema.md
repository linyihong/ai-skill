---
id: 2026-06-02-1200-plan-tree-frontmatter-schema
plan_kind: sub
status: in-progress
owner: linyihong
created: 2026-06-02
parent: 2026-06-02-1200-plan-tree-hierarchy-governance
required_for_completion: true
sub_plan_reason: >
  Frontmatter schema 是其餘 4 條 validator 的依據（parent reference / unique
  id / archive order / folder convention 全部 query frontmatter）。schema
  必須先文件化 + fixture 落地，Phase 2 validator 才有測試錨點。獨立成 sub-plan
  以便 Phase 0/1 與 Phase 2 解耦（並行安全路 — registry plan archive 前不動 hooks.go）。
---

# Plan Tree Frontmatter Schema（sub-plan）

> 本 sub-plan 繼承 parent 的 Decision Rationale（minimal governance）。

## Purpose

文件化 plan tree frontmatter schema 並產出 3 個 fixture（main / sub / spike），讓 Phase 2 validator 有測試錨點，也讓使用者寫新 plan 時有可複製範例。

## Acceptance Criteria

- [x] 本 sub-plan 檔案落地（使用 new schema 自證）
- [ ] `governance/plan-tree-hierarchy.md` rule 文件落地
- [x] `fixtures/main-plan.md` 範例落地
- [x] `fixtures/sub-plan.md` 範例落地
- [x] `fixtures/spike-plan.md` 範例落地
- [ ] 本 sub-plan 已被 parent `_plan.md` §主計畫必填：Sub-Plan 驗證要點表 引用（已存在）
- [ ] Phase 2 validator 撰寫時直接重用本目錄 fixtures 作為 unit test testdata

## Runtime Impact

無 runtime trigger 改動；不新增 routing-registry entry，不新增 `generated_surfaces`。本 sub-plan 純文件化交付，Phase 2 才接 hooks.go。

## Schema 定義

詳見 parent `_plan.md` §Frontmatter Schema（Minimal Governance）。本 sub-plan 不重複定義，只補 fixture 與註解。

### Required fields（all plans）

| Field | Type | 說明 |
|---|---|---|
| `id` | string | 全域唯一 slug。建議格式：`YYYY-MM-DD-HHMM-<descriptor>` |
| `plan_kind` | enum | `main` / `sub` / `spike` |
| `status` | enum | `draft` / `in-progress` / `completed` |
| `owner` | string | accountable identity |
| `created` | date | ISO YYYY-MM-DD |
| `parent` | string\|null | main: `null`；sub/spike: parent main 的 `id` |

### Sub-plan 額外 required

| Field | Type | 說明 |
|---|---|---|
| `required_for_completion` | bool | 是否屬於 parent acceptance criteria（archive blocker 由此推導） |
| `sub_plan_reason` | string | 為什麼拆 plan（非空 free text；不審內容） |

### Spike 預設

`required_for_completion: false`（spike 預設不阻擋 parent archive），但 schema 不禁止 spike 設為 true。

## Validator 將檢查的條件（給 Phase 2 預先 freeze）

| 條件 | 嚴重度 | 對應 fixture |
|---|---|---|
| Sub-plan 缺 `parent` | block | `fixtures/sub-missing-parent.md`（Phase 2 補） |
| Sub-plan `sub_plan_reason` 空字串 | block | `fixtures/sub-empty-reason.md`（Phase 2 補） |
| Sub-plan 缺 `required_for_completion` | block | `fixtures/sub-missing-required.md`（Phase 2 補） |
| `parent` 指向不存在的 id | block | `fixtures/parent-orphan.md`（Phase 2 補） |
| 同一 `id` 出現 ≥ 2 次 | block | `fixtures/duplicate-id-{a,b}.md`（Phase 2 補） |
| `required_for_completion: true` 的 sub 在 parent archive 時未 completed | block | `fixtures/archive-required-pending.md`（Phase 2 補） |
| Folder 缺 `_plan.md` / 檔名不符 NN- / 深度 ≥ 3 | warning | `fixtures/folder-shape-violations/`（Phase 2 補） |

Phase 1 只交付 3 個 happy-path fixture；Phase 2 補完所有 negative fixture。

## 與 parent 的同步

本 sub-plan 完成時：
- 在 parent `_plan.md` §Phase 1 標 status: completed
- 把 `governance/plan-tree-hierarchy.md` rule 連動更新

## Source

Parent: `2026-06-02-1200-plan-tree-hierarchy-governance/_plan.md` §Phase 1
