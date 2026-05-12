# Retry Smell Heuristic（重試壞味道經驗法則）

**Status**: `candidate-intelligence`
**Source**: 通用軟體工程經驗、分散式系統失效模式

## 原則

**More than 3 retries often indicates architectural failure, not transient instability.**

超過 3 次重試通常代表架構性失敗，而非短暫不穩定。

## 為什麼

1. 如果一個操作需要重試 3 次以上才能成功，問題通常不在於網路瞬斷，而是：
   - 下游服務不健康（degraded 而非 crashed）。
   - 請求本身有問題（例如 bad request、authentication failure）。
   - 架構設計有缺陷（例如缺少 circuit breaker、timeout 太短）。
2. 大量重試會放大問題：retry storm 可能讓已經不健康的下游服務完全崩潰。
3. 重試是暫時止痛藥，不是治療方案。如果重試「有效」，只是因為問題還不夠嚴重到被看見。

## 何時適用

- 看到 `retry(5)` 或 `while(true)` 搭配 sleep 的重試邏輯。
- 重試次數大於 3 且沒有 exponential backoff 或 jitter。
- 重試成功率高但 latency 不穩定（表示下游處於 degraded 狀態）。

## 何時不適用

- 只有 1-2 次重試，且有 exponential backoff + jitter。
- 重試是為了處理已知的、短暫的競爭條件（例如樂觀鎖衝突）。
- 重試是 client-side 的等冪性保證（例如 at-least-once delivery）。

## 決策流程

```text
看到重試邏輯？
  ├── 重試次數 > 3？
  │     ├── 是 → 這是 architectural smell，需要檢查下游健康度與請求正確性
  │     └── 否 → 繼續
  │
  └── 有 exponential backoff + jitter？
        ├── 是 → 合理，但監控重試率
        └── 否 → 加入 backoff / jitter，否則重試可能讓問題惡化
```

## 常見誤用

| 誤用 | 正確 |
|------|------|
| 對 4xx 錯誤重試（例如 401、403、400） | 4xx 是 client error，重試不會改變結果 |
| 對 timeout 重試 10 次 | 如果 timeout 是 1s，10 次重試 = 10s 的 blocking time，且可能讓下游更慢 |
| 沒有 circuit breaker 的重試 | 下游已經 degraded 時，重試只會加速崩潰 |

## Token Impact

避免 retry storm 造成的 cascading failure。一個沒有 backoff 的重試邏輯可以在 30 秒內讓一個 degraded 服務完全崩潰。

---

← [回到 engineering/heuristics/](README.md)
