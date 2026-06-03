# Plan Tree Flat Ambiguity（扁平 plan 目錄使 main/sub 關係不可追蹤）

Status: validated
Class: `process-gap` / `traceability-gap`

## Trigger

當主計畫執行中需要拆出支線工作，agent 用下列任一方式處理 main ↔ sub 關係時，使用此 pattern：

- 把支線塞進主計畫的 Phase 6/7/8…，直到主計畫臃腫到無法管理。
- 另開一個新 plan，只在檔頭 prose 用「與其他 plans 的關係」段落手寫 parent ↔ child 連結。
- 在 `plans/active/` 平鋪所有 plan，靠人眼從 `ls` 推測哪些是主軸、哪些是衍生支線。

具體觸發訊號：

- `plans/active/` 出現兩個明顯有 parent-child 關係的 plan，但檔案橫向平鋪、無 frontmatter `parent`。
- Sub-plan 的 parent context 只存在於 prose（「Source: §Phase 3 Round-4 T1」「parent plan X」），不是機器可讀 pointer。
- 主計畫完成條件被某個 sub-plan blocker 卡住，但主計畫表頭 `Status` 不反映。
- 已 archived 的 sub-plan 回頭追溯時找不到 parent。

## Failure Mode

把「plan 之間的階層關係」當成 prose 連結或目錄擺放的隱性知識，導致：

1. **Hierarchy 不可機器追蹤**：沒有單一 source of truth，工具無法建樹、無法 cross-check referential integrity。
2. **主計畫臃腫**：所有支線細節塞回主計畫，phase 無限增長，sub-scope 與主軸驗證要點混雜。
3. **Orphan / dangling 連結**：prose 寫的 parent 名稱 rename 或 archive 後失效，無人偵測。
4. **Archive 順序失控**：主計畫 archive 時無法機械確認 required sub-plan 是否都 completed。
5. **三軌不同步**：若同時用 folder + frontmatter `children:` + filename ordering 三處表達階層，rename 後三者不一致，無人知道誰是真實來源。

## Risk

- **Traceability silent drift**：每多一個未登記的 parent-child，repo 的 plan 關係圖就偏離真實一點，累積後無人能還原全貌。
- **Premature abstraction 反向風險**：為了解決上述問題而引入 DAG（`depends_on` / `children` / enum trigger），在沒有真實案例時過度設計，framework 隨情境膨脹。
- **Completion 假象**：主計畫宣稱完成，但 required sub-plan 仍 in-progress，無 gate 攔截。

## Required Agent Action

拆 sub-plan 或盤點既有 plan 關係時：

1. **Frontmatter `parent` pointer 為唯一 source of truth**。Sub-plan 必帶 `parent: <main-id>` / `required_for_completion: bool` / `sub_plan_reason: <非空 free text>`；不維護 `children:`，runtime scan 推導。
2. **Folder + `_plan.md` + `NN-` 前綴只當 UI convention**。folder 放錯不該讓 hierarchy 失效；`ai-skill plans tree` 由 `parent` 動態建樹。
3. **Lifecycle 與 storage 分離**：archive gate 只看 `status: completed`，不看 `active|archived` location。
4. **不為通用而通用**：sub-plan 之間若無真實 `depends_on` 需求，維持 Tree 不升 DAG；enum 化「為什麼拆 plan」改用 free-text reason，validator 只擋空字串。
5. **單一 phase 內 step 細分、< 1 session 工作、純文件補強，不開 sub-plan**（用 checkbox / inline / 直接 commit）。

## Prevention Gate

Commit-msg hook 機械強制（`scripts/ai-skill-cli/internal/app/plan_tree.go`）：

| Validator | Severity | 防的失效 |
|---|---|---|
| `validatePlanTreeFrontmatter` | block | sub-plan 缺 `parent` / `sub_plan_reason` / `required_for_completion` |
| `validatePlanTreeArchiveOrder` | block | 主計畫 archive 時 required sub-plan 未 completed |
| `validatePlanTreeParentReference` | block | `parent` 指向不存在的 id（orphan / dangling） |
| `validatePlanTreeUniqueID` | block | 同一 `id` 出現 ≥ 2 次 |
| `validatePlanTreeFolderConvention` | warning | folder shape 違規（不 block） |

治理規則：[`governance/lifecycle/plan-tree-hierarchy.md`](../../governance/lifecycle/plan-tree-hierarchy.md)。
詞彙：[`knowledge/glossary/ai-skill.md`](../../knowledge/glossary/ai-skill.md)（`plan_tree` / `parent` / `plan_kind` / `required_for_completion` / `sub_plan_reason`）。

## Validation Method

- `ai-skill plans tree --state all` 渲染出的階層與真實 parent-child 一致，即使 folder 放錯。
- 對 sub-plan 缺欄位 / dangling parent / duplicate id 的 fixture，commit-msg validator 機械 block（unit test in `plan_tree_test.go`）。
- 主計畫 archive 時，required sub-plan 全 `completed` 才放行。

## 來源

落地 plan：[`plans/active/2026-06-02-1200-plan-tree-hierarchy-governance/_plan.md`](../../plans/active/2026-06-02-1200-plan-tree-hierarchy-governance/_plan.md)。
真實案例：registry cluster（1 parent + 3 children）原本在 `plans/active/` 平鋪、靠 prose 連結，遷移為 frontmatter `parent` 後可由 CLI 正確建樹。

← [回到失效模式索引](README.md)
