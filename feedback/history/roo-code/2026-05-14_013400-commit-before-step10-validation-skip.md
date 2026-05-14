# Commit/Push Before Step 10 Validation Skip（commit/push 前跳過 knowledge-update-flow Step 10 驗證）

Status: candidate

#### One-line Summary

Agent 在執行 commit/push 前，跳過了 knowledge-update-flow Step 10（驗證）的完整檢查清單，導致去敏檢查、diff review、validator 等步驟被省略。

#### Human Explanation

knowledge-update-flow.md 定義了完整的 Step 1-11 流程，其中 Step 10（驗證）包含 6 項檢查：git status、去敏檢查、git diff、validator/lint、tool sync、close-loop script。但 agent 在實際執行時，經常只做 validator（第 4 項）就跳過其他項目直接 commit。這個問題在改動類型不涉及 `shared-rules/` 或 `skills/` 時更容易發生，因為 agent 會主觀判斷「這次改動不需要完整驗證」。

#### Trigger

- 改動類型是 `plans/`、`architecture/` 等非核心層文件時
- 改動規模小（1-2 個檔案）時
- 前一次 commit 已經執行過 validator 時
- 沒有明確的 checklist 提醒時

#### Evidence

- Tool: Roo Code
- Sanitized excerpt: 在 commit `afe6702`（feat(plans): add enforcement-layer-enhancement.md plan）中，agent 直接執行 `git add` → `git commit` → `git push`，跳過了 Step 10 的 git status 檢查、去敏檢查、git diff 檢查、validator 執行。直到使用者提醒後才補做。
- Evidence path: 無 project-specific evidence，此為通用 agent 行為模式。

#### Generalized Lesson

Agent 在執行任何 commit/push 前，**必須**完整執行 knowledge-update-flow Step 10 的 6 項檢查清單，不得因改動類型或規模而跳過。Step 10 的適用條件是「任何寫入操作」，不是「只有影響 runtime surface 的改動」。

#### Agent Action

下次 agent 看到需要 commit/push 的情境時：

1. **先執行 Step 10 檢查清單**，再執行 Step 11（Commit/Push/Readback）
2. 檢查清單不可省略：
   - ✅ `git status --short --branch` — 確認工作樹狀態
   - ✅ 去敏檢查 — 依 sanitization.md 檢查所有新增/修改的可重用文件
   - ✅ `git diff` — 檢查將提交的內容
   - ✅ 執行 validator / lint / link check
   - ✅ 若使用 tool mirror，執行 tool sync
   - ✅ 若有多個 owner group，使用 close-loop script
3. 即使改動只涉及 `plans/`、`README.md` 等非核心層，也必須執行 Step 10

#### Goal / Action / Validation

- Goal: 每次 commit/push 前都完整執行 Step 10 檢查清單
- Action: 在執行 `git add` 之前，先逐項檢查 Step 10 的 6 個項目
- Validation or reference source: `governance/lifecycle/knowledge-update-flow.md` §Step 10

#### Applies When

- 任何需要 commit/push 的情境
- 任何改動類型（包括 plans/、README.md、architecture/ 等非核心層）

#### Does Not Apply When

- 僅讀取操作（reference-only），無任何寫入
- 使用 `--amend` 或 `--fixup` 修正前一次 commit 的 typo（但仍需檢查去敏）

#### Validation

1. 檢查 commit log 中是否有明確記錄「已執行 Step 10 檢查清單」
2. 若發現 commit 前沒有 Step 10 的執行記錄，視為違反本規則

#### Promotion Target

- `shared-rules/failure-patterns/` — 此為跨 skill 的 failure class，可建立 failure pattern 文件

#### Required Linked Updates

- 依 [`linked-updates.md`](../../shared-rules/linked-updates.md) 列出必須同步更新或已檢查的相關文件：
  - `governance/lifecycle/knowledge-update-flow.md` — 已確認 Step 10 檢查清單完整，無需修改
  - `shared-rules/failure-patterns/` — 若 promotion 後需建立 failure pattern 文件
- 已依 [`reusable-guidance-boundary.md`](../../shared-rules/reusable-guidance-boundary.md) 檢查：本 lesson 只包含 generalized agent 行為規則，無 project-specific evidence。
