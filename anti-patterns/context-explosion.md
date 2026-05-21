# Context Explosion

## 症狀

- Context 持續成長，從未 pruning
- Token 使用率持續 > 80%
- Agent 回應變慢、品質下降
- 同一 session 中處理過多不相關任務

## 根本原因

1. **無 TTL 機制**：Context 一旦載入就永久保留。
2. **無 budget 管理**：沒有 token 預算概念。
3. **無 prune 策略**：不知道何時該清理 context。
4. **Task 堆疊**：在同一 session 中處理過多 task。

## 影響

- Token 成本暴增
- AI 推理品質下降（context 過大時）
- 需要手動建立新 session
- 工作進度中斷

## 預防

1. 實施 Context TTL（`runtime/runtime.db`）。
2. 設定 Token Budget（`runtime/runtime.db`）。
3. 在 task boundary 自動 prune。
4. 使用 summary-first routing 減少 context 載入量。

## 檢測

- `runtime/runtime.db` 的 `relevance` 與 `staleness` 維度
- `runtime/runtime.db` 的 `token_utilization` 信號
- Token 使用率 > 70% 時觸發警告

## 恢復

1. 執行 `runtime/runtime.db` 的 auto_archive。
2. 建立新 session。
3. 從 `memory/summary/` 載入必要的 session context。

## 相關 Guards

- `runtime/runtime.db`
- `runtime/runtime.db`
- `runtime/runtime.db`
- `runtime/runtime.db`
