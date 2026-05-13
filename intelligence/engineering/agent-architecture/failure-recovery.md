# Failure Recovery（失效恢復）

**Status**: `candidate-intelligence`
**Source**: 本系統實際運作觀察

## 原則

**Agent's first recovery attempt after failure is usually the most reliable; repeated retries without strategy change degrade output quality.**

Agent 在失效後的第一次恢復嘗試通常是最可靠的；沒有策略變更的重複重試會降低輸出品質。

## 為什麼

1. **第一次恢復有最多上下文** — Agent 在剛發現錯誤時，對錯誤的上下文理解最完整。隨著重試次數增加，上下文被污染，判斷力下降。
2. **重複重試產生噪音** — 每次失敗的重試都會在上下文中留下錯誤的 tool call 記錄，這些記錄會干擾後續推理。
3. **Agent 難以自我修正策略** — Agent 傾向於用相同的方法重試，而不是停下來重新分析策略。這是因為「繼續嘗試」比「重新規劃」更符合 agent 的預設行為模式。
4. **疲勞效應** — 多次失敗後，agent 可能開始接受不完美的解決方案，或做出草率的決策。

## 症狀

| 症狀 | 說明 | 可信度 |
|------|------|--------|
| **相同錯誤重複** | Agent 用相同的方法重試同一操作 3+ 次，每次得到相同錯誤 | 高 |
| **恢復後品質下降** | 恢復成功後的輸出比正常狀態下的輸出品質明顯較低 | 高 |
| **忽略錯誤訊息** | Agent 重複嘗試時不再仔細閱讀錯誤訊息 | 中 |
| **降級策略** | Agent 從「最佳實踐」降級到「只要能動就好」 | 中 |

## 恢復策略

```text
發現錯誤？
  ├── 這是新類型的錯誤？
  │     ├── 是 → 分析錯誤訊息，制定新策略
  │     └── 否 → 查詢 failure-patterns/ 是否有對應 pattern
  │
  ├── 第一次重試？
  │     ├── 是 → 可以重試，但要確認策略正確
  │     └── 否 → 停下來，重新分析
  │
  ├── 已重試 2+ 次？
  │     ├── 是 → 必須停下來：
  │     │     ├── 讀取相關文件確認正確方法
  │     │     ├── 查詢 failure-patterns/
  │     │     ├── 考慮使用不同的工具或方法
  │     │     └── 如果仍然不確定，向 user 請求指引
  │     └── 否 → 繼續
  │
  └── 錯誤是工具層面的？
        ├── 是 → 檢查工具配置、版本、權限
        └── 否 → 檢查邏輯錯誤
```

## 預防方式

1. **設定重試上限** — 同一操作最多重試 2 次，之後必須停下來重新分析
2. **記錄失敗** — 將失敗記錄到 `failure-patterns/` 或 `feedback_history/`
3. **建立 validation scenario** — 如果錯誤是 stateless 可重現的，建立 validation scenario
4. **恢復後驗證** — 恢復成功後，執行額外的驗證步驟確認品質沒有下降
5. **策略變更** — 每次重試必須改變策略，不能使用完全相同的方法

## 不建議的做法

| 不建議 | 原因 |
|--------|------|
| 不分析錯誤就重試 | 重複相同的錯誤操作只會浪費 token |
| 重試超過 3 次 | 超過 3 次後，上下文已被嚴重污染 |
| 恢復後跳過驗證 | 恢復狀態下的輸出品質可能低於正常狀態 |

## 相關 atoms

- [`cognitive-boundaries.md`](cognitive-boundaries.md) — 認知邊界
- [`context-collapse.md`](context-collapse.md) — 上下文崩塌
- [`task-routing.md`](task-routing.md) — 任務路由

## Token Impact

無策略的重複重試是 token 浪費的典型模式。3 次無效重試可能消耗 5K-15K token，而且每次重試都會污染上下文，降低後續輸出的品質。

---

← [回到 agent-architecture/](README.md)
