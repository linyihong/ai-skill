# Recursive Tool Loop

## 症狀

- 同一 tool 反覆呼叫（search → search → search）
- 每次呼叫都沒有新結果
- Tool call 次數快速增加
- Agent 陷入「搜尋 → 沒找到 → 再搜尋」循環

## 根本原因

1. **搜尋條件不精準**：第一次搜尋沒找到，擴大範圍再搜，持續擴大。
2. **無結果處理策略**：沒找到時不知道該停止或換方法。
3. **無呼叫次數限制**：沒有 circuit breaker。

## 影響

- Token 大量浪費
- 工作無進展
- Context 被無效 tool output 污染

## 預防

1. 設定 tool call 次數限制（`runtime/guards/circuit-breaker.yaml`）。
2. 搜尋前先確認關鍵詞與範圍。
3. 沒找到時先問使用者，而不是盲目擴大搜尋。
4. 使用 `tools/routing/` 的 explosion detection。

## 檢測

- 同一 tool 連續呼叫 > 3 次且無新結果
- 5 分鐘內 tool call > 15 次
- 搜尋結果為空但仍繼續搜尋

## 恢復

1. 停止目前 tool chain。
2. 詢問使用者正確的搜尋方向。
3. 必要時建立新 task 重新開始。

## 相關 Guards

- `runtime/guards/circuit-breaker.yaml`
- `tools/routing/README.md`
