# Commit/Push Before Writeback Transaction Close（commit/push 前跳過 Ai-skill Writeback Transaction 關閉）

Status: candidate

#### One-line Summary

Agent 在執行 commit/push 前，跳過了 Ai-skill writeback transaction 的關閉條件（依 [`dependency-reading.md`](../../shared-rules/dependency-reading.md) §Ai-skill Writeback Transaction Guard），導致去敏檢查、diff review、validator、push 後讀回等步驟被省略。

#### Human Explanation

[`dependency-reading.md`](../../shared-rules/dependency-reading.md) §Ai-skill Writeback Transaction Guard 定義了任何寫入操作（shared-rules/、skills/、feedback lessons、plans/ 等）的完整交易關閉條件，共 7 項：git status/diff 檢查、去敏檢查、linked updates、tool sync、git add/commit/push、push 後讀回、clean status。但 agent 在實際執行時，經常只做 validator 就跳過其他項目直接 commit。這個問題在改動類型不涉及 `shared-rules/` 或 `skills/` 時更容易發生，因為 agent 會主觀判斷「這次改動不需要完整驗證」。

注意：`knowledge-update-flow.md` Step 10 的檢查清單與 Writeback Transaction Guard 的關閉條件有重疊，但後者是所有寫入操作的唯一真相來源。Step 10 是 knowledge-update-flow 專用的驗證步驟，不應被當作通用 commit/push 檢查清單。

#### Trigger

- 改動類型是 `plans/`、`architecture/` 等非核心層文件時
- 改動規模小（1-2 個檔案）時
- 前一次 commit 已經執行過 validator 時
- 沒有明確的 transaction close checklist 提醒時

#### Evidence

- Tool: Roo Code
- Sanitized excerpt: 在 commit `afe6702`（feat(plans): add enforcement-layer-enhancement.md plan）中，agent 直接執行 `git add` → `git commit` → `git push`，跳過了 writeback transaction 的關閉條件（git status 檢查、去敏檢查、git diff 檢查、validator 執行）。直到使用者提醒後才補做。
- Evidence path: 無 project-specific evidence，此為通用 agent 行為模式。

#### Generalized Lesson

Agent 在執行任何寫入操作（commit/push）前，**必須**關閉 Ai-skill writeback transaction（依 [`dependency-reading.md`](../../shared-rules/dependency-reading.md) §Ai-skill Writeback Transaction Guard 的 7 項關閉條件），不得因改動類型或規模而跳過。交易關閉條件的適用範圍是「任何寫入操作」，不是「只有影響 runtime surface 的改動」。

#### Agent Action

下次 agent 看到需要 commit/push 的情境時：

1. **先確認 writeback transaction 已開啟**（依 [`dependency-reading.md`](../../shared-rules/dependency-reading.md) §Ai-skill Writeback Transaction Guard）
2. 關閉 transaction 前，完成交易關閉條件：
   - ✅ `git status --short --branch` 與 `git diff` 已檢查
   - ✅ 去敏檢查（依 [`sanitization.md`](../../shared-rules/sanitization.md)）
   - ✅ 必要的 linked updates 已同步或明確寫出不適用理由
   - ✅ 若使用 tool mirror，執行 tool sync
   - ✅ 相關檔案已 `git add`、`git commit`、`git push`
   - ✅ Push 後已重新讀取更新過的入口、主要依賴、索引與 promotion target
   - ✅ 最後一次 `git status --short --branch` 顯示 clean，branch 沒有 ahead/behind
3. 即使改動只涉及 `plans/`、`README.md` 等非核心層，也必須關閉 transaction

#### Goal / Action / Validation

- Goal: 每次 commit/push 前都關閉 Ai-skill writeback transaction
- Action: 在執行 `git add` 之前，先依 [`dependency-reading.md`](../../shared-rules/dependency-reading.md) §Ai-skill Writeback Transaction Guard 逐項檢查關閉條件
- Validation or reference source: [`shared-rules/dependency-reading.md`](../../shared-rules/dependency-reading.md) §Ai-skill Writeback Transaction Guard

#### Applies When

- 任何需要 commit/push 的情境
- 任何改動類型（包括 plans/、README.md、architecture/ 等非核心層）

#### Does Not Apply When

- 僅讀取操作（reference-only），無任何寫入
- 使用 `--amend` 或 `--fixup` 修正前一次 commit 的 typo（但仍需檢查去敏）

#### Validation

1. 檢查 commit log 中是否有明確記錄「已關閉 Ai-skill writeback transaction」
2. 若發現 commit 前沒有 transaction 關閉記錄，視為違反本規則

#### Promotion Target

- `shared-rules/failure-patterns/` — 此為跨 skill 的 failure class，已建立 failure pattern 文件

#### Required Linked Updates

- 依 [`linked-updates.md`](../../shared-rules/linked-updates.md) 列出必須同步更新或已檢查的相關文件：
  - [`shared-rules/dependency-reading.md`](../../shared-rules/dependency-reading.md) — §Ai-skill Writeback Transaction Guard（交易關閉條件的唯一真相來源）
  - [`shared-rules/failure-patterns/commit-before-validation-skip.md`](../../shared-rules/failure-patterns/commit-before-validation-skip.md) — 已更新為指向 Writeback Transaction Guard
  - [`shared-rules/failure-patterns/README.md`](../../shared-rules/failure-patterns/README.md) — 已更新索引描述
- 已依 [`reusable-guidance-boundary.md`](../../shared-rules/reusable-guidance-boundary.md) 檢查：本 lesson 只包含 generalized agent 行為規則，無 project-specific evidence。
