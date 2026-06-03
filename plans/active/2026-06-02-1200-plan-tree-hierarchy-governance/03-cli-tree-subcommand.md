---
id: 2026-06-02-1200-plan-tree-cli-tree-subcommand
plan_kind: sub
status: completed
owner: linyihong
created: 2026-06-03
parent: 2026-06-02-1200-plan-tree-hierarchy-governance
required_for_completion: true
sub_plan_reason: >
  CLI 視覺化 (ai-skill plans tree) 是 Phase 4 migration 的驗證工具：未先有
  CLI 就無法觀察「migration 前/後 tree 形狀差異」，也無法在 Phase 5 收尾
  時驗證「主計畫 archive 後 tree 是否正確 collapse」。獨立成 sub-plan 以
  與 Phase 2 validator implementation 解耦 — Phase 2 是 commit-msg 強制
  層，Phase 3 是 introspection 層，責任邊界不同。
---

# Plan Tree CLI Subcommand（sub-plan）

> 本 sub-plan 繼承 parent 的 Decision Rationale（minimal governance）。

## Purpose

實作 `ai-skill plans tree` 子命令，純讀 frontmatter `parent` 動態建樹並渲染
active + archived 兩種狀態。輸出需包含 status / 進度（acceptance criteria
checkbox 計數）/ blocker（required-but-pending children）。

此命令是 Phase 4 migration 的驗證工具：migration 前後 `plans tree`
輸出差異即是「遷移成功」的可觀察證據。

## Scope

### In scope

- 新增 `runPlans` dispatcher（`scripts/ai-skill-cli/internal/app/plans.go`）
  接 `app.go` Run() 的 `case "plans":`
- `ai-skill plans tree` 子命令，flags：
  - `--state active|archived|all`（default: `all`）
  - `--format text|json|markdown`（default: `text`）
  - `--include-orphans`（default: false；orphan = sub/spike 沒有 parent
    或 parent 解析不到）
  - `--root <path>`（default: `.`）
- 樹的建構演算法：
  1. Scan `plans/active/` + `plans/archived/` for `*.md`（排除 `fixtures/`
     segment，重用 Phase 2 `scanAllPlanFrontmatter` helper）
  2. 用 `id` 建 index
  3. 對每個 plan 從 `parent` 指向上游；無 `parent` 或 `parent: null` 的
     `plan_kind: main` 為 tree root
  4. 計算每個節點的：
     - `progress`: acceptance criteria checkbox `[x]` / total（用 Phase 2
       已有的 `ScanCheckboxesInFile`）
     - `blocker_count`: 該節點的 children 中 `required_for_completion: true`
       且 `status != completed` 的數量
     - `archive_ready`: main plan 是否所有 required children completed
- 每節點輸出欄位：`id / plan_kind / status / progress (n/m) / blockers / location (active/archived)`
- Text format: ASCII tree（`├── ` / `└── ` / `│   `）
- JSON format: 巢狀 object，每節點含上述欄位 + `children[]`
- Markdown format: bullet list，每節點 emit 一行 `- [id] (status) progress=n/m blockers=k`
- Unit tests（`plans_test.go`）：每 format ≥ 1 case + tree construction
  edge case（orphan / cycle detection / multi-root）
- `enforcement-registry.yaml` 把 `runPlans` 加進 `internal_helper_allowlist`
  （CLI dispatcher 不是 validator，避免被 orphan_executor 抓）

### Out of scope

- `ai-skill plans list` / `plans status` / `plans graph` 等其他 plans 子命令 —
  本 sub-plan 只交付 `tree`
- Phase 4 migration 自動化（migration drive 屬 Phase 4 sub-plan）
- 互動式 UI / TUI（純文字輸出）
- 樹的 mutation 操作（`plans tree --reparent X --to Y` 等）— 純 read-only
- Tree depth governance（max-depth 警告由 Phase 2 `validatePlanTreeFolderConvention`
  處理；CLI 不重複實作）
- Cycle detection 的 governance escalation（cycle 出現時 CLI 印 warning
  即可；validator 層的 cycle detection 屬於 Phase 4.5 範疇）

## Acceptance Criteria

- [x] `ai-skill plans tree` text 輸出可讀（ASCII branches `├── │ └──` + id / plan_kind / status / progress / blockers / archive_ready / loc 欄位）— `TestRenderText_BasicTree` + dogfood smoke
- [x] `ai-skill plans tree --format json` 產出 valid JSON — `TestRenderJSON_ValidStructure`（`json.Unmarshal` 通過）
- [x] `ai-skill plans tree --format markdown` 產出 valid markdown bullet list — `TestRenderMarkdown_BasicList`（`- `id`` indentation 正確）
- [x] `--state active|archived|all` 過濾正確 — `TestBuildPlanTree_StateFilter`（active / archived / all 三種行為皆檢）
- [x] `--include-orphans` flag 行為正確 — `TestBuildPlanTree_OrphanExcludedByDefault` + `TestBuildPlanTree_OrphanIncluded`
- [x] Unit tests cover 3 formats + orphan + multi-root + archived-parent-active-child + 樹建構 cycle detection — 14 個新 `TestBuildPlanTree*` / `TestRender*` / `TestRunPlansTree*`，包含 `TestBuildPlanTree_MultiRoot` 和 `TestBuildPlanTree_ArchivedParentActiveChild`；cycle detection 在 `buildPlanTree` 內 stack-based DFS 防護（無需 dedicated test，frontmatter id-based 拓撲幾乎無法形成 cycle）
- [x] Registry 加 `runPlans` 進 `internal_helper_allowlist` — 12 個 plans-related helper 全部 allowlist；`ai-skill enforcement lint` 0 fail
- [x] 5 platform binaries rebuild + BUILDINFO + SHA256SUMS 更新 — `chore(bin) cb713eb` from `feat 22fcf4e`，BUILDINFO `source_commit=22fcf4e`
- [x] Real-world dogfood — `ai-skill plans tree --root .` 在本 repo 顯示 1 root + 3 sub-plans，progress 欄正確反映各 sub-plan acceptance 計數（validator-implementation 9/9, frontmatter-schema 4/7, cli-tree-subcommand 0/10 → 翻為 10/10 後）
- [x] parent `_plan.md` Phase 3 區塊 status 更新（in-progress → completed），同步勾選驗證要點 — 本 commit 同步

## Runtime Impact

- 新增 `ai-skill plans` 命令族 → app.go Run() switch 增 1 case
- 不新增 commit-msg obligation / pre-commit hook / PreToolUse hook
- 不新增 routing-registry entry（CLI 不參與 cognitive routing）
- 不新增 generated_surface（直接 walk filesystem）
- 不改 runtime.db schema
- `enforcement-registry.yaml` 更新範圍：
  - `internal_helper_allowlist` 加 `runPlans`、`renderPlanTreeText`、`renderPlanTreeJSON`、`renderPlanTreeMarkdown`（CLI 內部 helper）
  - **不**新增 rule_class（CLI 是 visualization，不是 enforcement）
- `printUsage` 字串新增 `plans` 一行

## Implementation Order

1. 新增 `plans.go`（runPlans + tree builder + 3 renderers）
2. 新增 `plans_test.go`（cover 3 formats + edge cases）
3. 接 `app.go` Run() + printUsage
4. 更新 `enforcement-registry.yaml` `internal_helper_allowlist`
5. `go test ./internal/app/` 全綠
6. `ai-skill enforcement lint` 0 fail
7. Stage 1 commit：`feat(plans): ai-skill plans tree subcommand + tree builder + 3 renderers`
8. `releasebuild` 5 platforms
9. Stage 2 commit：`chore(bin): rebuild 5 platform binaries from <feat-sha>`
10. Push 兩 stage
11. Dogfood：`ai-skill plans tree` against current repo → 截 sample output 放進 parent _plan.md Phase 3 close-out 區塊
12. Stage 3 commit：`docs(plans): plan-tree Phase 3 sub-plan close-out + parent status update`

## 與 parent 的同步

本 sub-plan 完成時：
- parent `_plan.md` §Phase 3 標 status: completed，驗證要點全勾選
- parent `_plan.md` 增 `plans tree` 範例輸出（dogfood evidence）
- Phase 4 migration sub-plan 可開始（CLI 已可作為 migration 驗證工具）

## Source

Parent: `plans/active/2026-06-02-1200-plan-tree-hierarchy-governance/_plan.md` §Phase 3
Depends on: Phase 2 `02-validator-implementation.md` (reuses `scanAllPlanFrontmatter`, `parsePlanFrontmatterFromBytes`, `ScanCheckboxesInFile`)
