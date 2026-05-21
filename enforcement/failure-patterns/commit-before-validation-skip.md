# Commit/Push Before Writeback Transaction Close（commit/push 前跳過 Ai-skill Writeback Transaction 關閉）

Status: candidate
Class: `validation-gap`

## Trigger

當 agent 在執行 `git add` → `git commit` → `git push` 前，沒有關閉 Ai-skill writeback transaction（依 [`dependency-reading.md`](../dependency-reading.md) §Ai-skill Writeback Transaction Guard），使用此 pattern。

具體觸發信號：
- 改動類型是 `plans/`、`architecture/`、`README.md` 等非核心層文件
- 改動規模小（1-2 個檔案）
- 前一次 commit 已經執行過 validator
- agent 說「改動很小，不需要完整驗證」

## Failure Mode

Agent 把 Ai-skill writeback transaction 的關閉條件視為「只有影響 runtime surface 時才需要執行」，而不是「每次寫入操作都必須關閉 transaction」。這導致：

1. 去敏檢查被跳過 — 可能意外提交含本機路徑或 project-specific evidence 的內容
2. `git diff` 檢查被跳過 — 提交前沒有確認實際變更內容
3. `git status --short --branch` 被跳過 — 可能遺漏未追蹤的檔案
4. Validator 被跳過 — 可能破壞 knowledge runtime 的完整性
5. Tool sync 被跳過 — 如果使用了 tool mirror 可能不同步
6. Push 後讀回被跳過 — 無法確認 promotion target 與索引已同步

## Risk

- 可重用文件中可能意外包含本機絕對路徑或 project-specific evidence
- Knowledge runtime 可能因未執行 validator 而處於不一致狀態
- 使用者需要額外提醒才能補做驗證，降低信任
- 同類錯誤會在下一輪或下一個 agent 重複發生

## Required Agent Action

1. 在執行 `git add` 之前，先確認已依 [`dependency-reading.md`](../dependency-reading.md) §Ai-skill Writeback Transaction Guard 開啟 writeback transaction
2. 關閉 transaction 前，完成交易關閉條件：
   - ✅ `git status --short --branch` 與 `git diff` 已檢查
   - ✅ 去敏檢查（依 [`sanitization.md`](../sanitization.md)）
   - ✅ 必要的 linked updates 已同步或明確寫出不適用理由
   - ✅ 若使用 tool mirror，執行 tool sync
   - ✅ 相關檔案已 `git add`、`git commit`、`git push`
   - ✅ Push 後已重新讀取更新過的入口、主要依賴、索引與 promotion target
   - ✅ 最後一次 `git status --short --branch` 顯示 clean，branch 沒有 ahead/behind
3. 不得因改動類型或規模而跳過上述任何項目
4. 若發現跳過，先補做驗證再繼續 commit

## Prevention Gate

當 agent 準備執行 `git add` 時，在繼續之前必須能回答：

- 我已經確認 writeback transaction 已開啟？
- 交易關閉條件中，哪幾項已完成、哪幾項尚未完成？
- 如果沒有全部完成，原因是什麼？（只有 reference-only 操作可豁免）
- 新增/修改的文件是否已檢查去敏？
- Validator 是否已執行且通過？

## 驗證

1. Commit log 中應有明確記錄「已關閉 Ai-skill writeback transaction」
2. 若發現 commit 前沒有 transaction 關閉記錄，視為違反本規則
3. 讀回新增或更新的 failure pattern、相關 index 與最終 `git status --short --branch`

## Linked Rules

- [`../dependency-reading.md`](../dependency-reading.md) — §Ai-skill Writeback Transaction Guard（交易關閉條件的唯一真相來源）
- [`governance/lifecycle/knowledge-update-flow.md`](../../governance/lifecycle/knowledge-update-flow.md) — Step 10 的驗證檢查清單（輔助參考，但交易關閉條件以 dependency-reading.md 為準）
- [`../sanitization.md`](../sanitization.md) — 去敏規則
- [`../linked-updates.md`](../linked-updates.md) — 連動更新檢查
- [`correction-loop-bypass.md`](correction-loop-bypass.md) — 相關的 validation-gap pattern

## Linked Validation Scenarios

- `validate_knowledge_runtime` — 執行 `ai-skill runtime validate` 確認 knowledge runtime 完整性
