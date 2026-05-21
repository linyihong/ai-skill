# Mandatory Step Blocker Bypass（強制步驟遇環境阻斷時靜默跳過）

Status: validated
Class: `process-gap` / `validation-gap`

## Trigger

當 agent 在執行強制步驟（如 knowledge-update-flow Step 9/10 的 runtime 驗證腳本）時遇到環境阻斷（如 Ruby 未安裝、CLI 工具缺失、服務未啟動），卻自行判斷「這是環境限制」後靜默跳過並繼續後續步驟，使用此 pattern。

具體觸發信號：
- 執行腳本或 CLI 工具時收到「command not found」或「not installed」錯誤
- 強制步驟的完成被 agent 備註為「環境限制，跳過」
- Agent 沒有把錯誤訊息顯示給使用者，也沒有等待使用者指示
- Step 流程繼續執行到下一步，使用者不知道某強制步驟被略過

## Failure Mode

Agent 把「環境阻斷」解釋為「我有理由跳過此強制步驟」，而不是「我必須立即停下來通知使用者」。這導致：

1. 使用者不知道某個強制步驟從未完成
2. Workflow 表面上「完成」，但實質上遺漏了關鍵驗證或更新
3. 錯誤被包裝在「備註」中，使用者需要主動審查才能發現
4. 同樣的環境問題在後續 session 繼續被靜默跳過，累積技術債

## Risk

- 強制步驟（如 runtime.db 更新、validator 執行）被架空，knowledge runtime 可能處於不一致狀態
- 使用者誤以為整個更新流程已正確執行完畢
- 環境問題從未被解決，因為 agent 從未明確要求使用者介入
- 信任損失：使用者期望 agent 遇到問題會主動告知，而非自行處理後沉默

## Required Agent Action

遇到強制步驟的環境阻斷時：

1. **立即停止**：不執行後續步驟
2. **顯示確切錯誤**：把錯誤訊息原文呈現給使用者（不得摘要或美化）
3. **說明影響**：解釋此步驟為強制步驟，跳過會造成什麼後果
4. **提出選項**：讓使用者選擇處理方式，例如：
   - 安裝缺少的工具後繼續
   - 改用手動驗證替代
   - 在流程文件中加入無工具替代程序後繼續
   - 暫緩此步驟並記錄為 open issue
5. **等待使用者明確指示**：不得自行選擇任何選項後繼續

## Prevention Gate

當 agent 準備跳過任何步驟前，必須能回答：

- 這個步驟在 workflow 中被標記為強制（`[強制]`）還是條件（`[條件]`）？
- 如果是強制步驟，使用者是否已在聊天介面明確授權跳過？
- 錯誤訊息是否已完整顯示給使用者（不是摘要，是原文）？
- 使用者選擇的處理方式是什麼？

若以上任何一項無法確認，**不得繼續執行後續步驟**。

## 驗證

1. Step 完成後，使用者可以在聊天記錄中找到確切的環境錯誤訊息
2. 流程中沒有任何強制步驟被標記為「環境限制，跳過」且無使用者確認
3. 若使用者選擇替代方案（如手動驗證），替代方案的執行結果已記錄

## Linked Rules

- [`../dependency-reading.md`](../dependency-reading.md) — §Ai-skill Writeback Transaction Guard（強制步驟定義的唯一來源）
- [`../../governance/lifecycle/knowledge-update-flow.md`](../../governance/lifecycle/knowledge-update-flow.md) — Step 9/10 強制執行條件
- [`commit-before-validation-skip.md`](commit-before-validation-skip.md) — 相關的 validation-gap（validation 被跳過的後果）
- [`correction-loop-bypass.md`](correction-loop-bypass.md) — 相關的 validation-gap pattern

## Linked Validation Scenarios

- Step 9 完成後確認：知識更新 session 中所有「環境阻斷」事件均有對應的使用者確認記錄
- Step 10 diff review 時確認：沒有強制步驟被標記為「跳過（環境限制）」且無使用者授權
