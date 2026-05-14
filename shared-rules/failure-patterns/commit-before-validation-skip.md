# Commit/Push Before Step 10 Validation Skip（commit/push 前跳過 knowledge-update-flow Step 10 驗證）

Status: candidate
Class: `validation-gap`

## Trigger

當 agent 在執行 `git add` → `git commit` → `git push` 前，沒有完整執行 knowledge-update-flow Step 10 的 6 項檢查清單時，使用此 pattern。

具體觸發信號：
- 改動類型是 `plans/`、`architecture/`、`README.md` 等非核心層文件
- 改動規模小（1-2 個檔案）
- 前一次 commit 已經執行過 validator
- agent 說「改動很小，不需要完整驗證」

## Failure Mode

Agent 把 knowledge-update-flow Step 10 的驗證檢查清單視為「只有影響 runtime surface 時才需要執行」，而不是「每次寫入操作都必須執行」。這導致：

1. 去敏檢查被跳過 — 可能意外提交含本機路徑或 project-specific evidence 的內容
2. `git diff` 檢查被跳過 — 提交前沒有確認實際變更內容
3. `git status --short --branch` 被跳過 — 可能遺漏未追蹤的檔案
4. Validator 被跳過 — 可能破壞 knowledge runtime 的完整性
5. Tool sync 被跳過 — 如果使用了 tool mirror 可能不同步

## Risk

- 可重用文件中可能意外包含本機絕對路徑或 project-specific evidence
- Knowledge runtime 可能因未執行 validator 而處於不一致狀態
- 使用者需要額外提醒才能補做驗證，降低信任
- 同類錯誤會在下一輪或下一個 agent 重複發生

## Required Agent Action

1. 在執行 `git add` 之前，先完整執行 knowledge-update-flow Step 10 的 6 項檢查清單：
   - ✅ `git status --short --branch` — 確認工作樹狀態
   - ✅ 去敏檢查 — 依 `sanitization.md` 檢查所有新增/修改的可重用文件
   - ✅ `git diff` — 檢查將提交的內容
   - ✅ 執行 validator / lint / link check
   - ✅ 若使用 tool mirror，執行 tool sync
   - ✅ 若有多個 owner group，使用 `scripts/ai-skill-close-loop.sh --commit`
2. 不得因改動類型或規模而跳過上述任何項目
3. 若發現跳過，先補做驗證再繼續 commit

## Prevention Gate

當 agent 準備執行 `git add` 時，在繼續之前必須能回答：

- 我已經執行 Step 10 的哪幾項檢查？
- 如果沒有全部執行，原因是什麼？（只有 reference-only 操作可豁免）
- 新增/修改的文件是否已檢查去敏？
- Validator 是否已執行且通過？

## 驗證

1. Commit log 中應有明確記錄「已執行 Step 10 檢查清單」
2. 若發現 commit 前沒有 Step 10 的執行記錄，視為違反本規則
3. 讀回新增或更新的 failure pattern、相關 index 與最終 `git status --short --branch`

## Linked Rules

- [`governance/lifecycle/knowledge-update-flow.md`](../../governance/lifecycle/knowledge-update-flow.md) — Step 10 的完整檢查清單
- [`../dependency-reading.md`](../dependency-reading.md) — §Ai-skill 回寫完成門檻、§Commit/Push 後讀回 Gate
- [`../sanitization.md`](../sanitization.md) — 去敏規則
- [`../linked-updates.md`](../linked-updates.md) — 連動更新檢查
- [`correction-loop-bypass.md`](correction-loop-bypass.md) — 相關的 validation-gap pattern

## Linked Validation Scenarios

- `validate_knowledge_runtime` — 執行 `ruby scripts/validate-knowledge-runtime.rb` 確認 knowledge runtime 完整性
