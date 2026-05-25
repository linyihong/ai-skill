# Correction Loop Bypass（修正閉環略過）

Status: validated
Class: `validation-gap`

## Trigger

當使用者指出 agent 的修正不完整、規劃漏掉 `.agent-goals`、沒有回饋 failure learning、只修了文字但沒有完成 linked updates / validation / commit / push / readback 時，使用此 pattern。

## Failure Mode

Agent 把使用者指出的問題視為單點文字修正，只修當下檔案或只口頭承認錯誤，卻沒有完成 correction close loop：分析漏掉原因、更新 goal ledger、沉澱 failure pattern、檢查 linked updates、驗證並關閉 repository writeback transaction。

## Risk

- 同類錯誤會在下一輪或下一個 agent 重複發生。
- 使用者以為系統規則已變強，但實際只有當下回覆被修正。
- `.agent-goals` 缺少 active work，導致長任務中斷後無法恢復。
- Repository 可能出現已修改但未驗證、未提交、未推送或未讀回的狀態。

## Required Agent Action

1. 先承認 correction loop gap，而不是只說「已修」。
2. 建立或更新 `.agent-goals/`，記錄目前修正目標、owner、parallelization、next action 與 completion criteria。
3. 依 `failure-learning-system.md` 分類失效模式。
4. 若可跨專案或跨 agent 重演，新增或更新 `enforcement/failure-patterns/`。
5. 依 `linked-updates.md` 檢查相關索引、入口、tool docs 或 workflow。
6. 執行可反查 validation：diff review、link/lint check、stale wording search、commit/push/readback、clean status。
7. 完成後刪除已驗證的 goal ledger entry。

## Prevention Gate

當使用者說「你沒有做閉環」「你沒有回饋失敗經驗」「你沒有建立 goal」「你只修當下」等語意時，agent 在繼續其他工作前必須能回答：

- 這次漏掉的 root cause 是什麼？
- 目前 `.agent-goals` 是否已有對應 active goal？
- 這是 `validation-gap`、`goal-ledger-miss`、`dependency-miss`、`source-mirror-drift` 或其他 class？
- Durable prevention 寫到哪裡？
- 哪些 linked updates 已檢查？
- 哪些 validation 已完成？

## Validation

讀回新增或更新的 failure pattern、相關 index、active/cleared goal ledger 與最終 `git status --short --branch`。若更新 canonical repository，必須完成 commit、push、readback 與 clean status。

## Linked Rules

- [`../failure-learning-system.md`](../failure-learning-system.md)
- [`../conversation-goal-ledger.md`](../conversation-goal-ledger.md)
- [`../goal-action-validation.md`](../goal-action-validation.md)
- [`../linked-updates.md`](../linked-updates.md)
- [`../dependency-reading.md`](../dependency-reading.md)

## Linked Validation Scenarios

- `validate_no_outdated_active_entrypoint` — 檢查 `workflow/`、`analysis/`、`runtime/onboarding/` 中是否仍有過時的 active entrypoint 參考，防止修正閉環略過後留下 stale 內容
