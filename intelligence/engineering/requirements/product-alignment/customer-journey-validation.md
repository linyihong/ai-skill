# Customer Journey Validation

**Status**: `candidate-intelligence`

## 判斷原則

Customer Journey Map 從使用者真實路徑檢查需求。它回答「痛點在哪一段發生、使用者當下有什麼情緒或阻礙、功能是否真的命中那個時刻」。

## 檢查項

| Journey 欄位 | 問題 |
| --- | --- |
| Stage | 行為發生在 onboarding、activation、usage、retention、payment、support 或 recovery 哪一段？ |
| Trigger | 使用者為什麼進入這段旅程？ |
| Pain point | 哪個摩擦、疑惑、等待、錯誤或風險造成流失？ |
| Emotional low | 使用者情緒低點在哪裡？它是否影響 conversion、retention 或 trust？ |
| Blocker | 什麼阻擋使用者完成目標？ |
| Evidence | 來自觀察、support ticket、analytics、訪談、測試或明確 assumption？ |

## 風險訊號

- Feature 解的是內部便利，不是 journey pain。
- 使用者情緒低點與 feature investment 無關。
- 需求只描述 happy path，沒有對應 blocker / failure path。
- Journey stage 不清楚，導致 acceptance criteria 太泛。

## 行動

若 journey evidence 不足，將 claim 標記為 `assumption` 或 `open question`；若 feature 與 pain point 不對齊，回到 Impact Map 修正 `Who / How / What`。
