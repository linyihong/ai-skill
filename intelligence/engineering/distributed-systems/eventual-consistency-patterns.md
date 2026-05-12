# Eventual Consistency Patterns（最終一致性模式）

**Status**: `candidate-intelligence`
**Source**: 分散式系統營運經驗

## 原則

**Eventual consistency requires idempotent consumers. If a consumer cannot safely process the same event twice, the system is not eventually consistent — it's eventually inconsistent.**

最終一致性需要冪等消費者。如果消費者無法安全地處理同一事件兩次，系統不是最終一致 — 而是最終不一致。

## 為什麼

1. **At-least-once delivery 是常態**：大部分 message queue（Kafka、RabbitMQ、SQS）保證 at-least-once delivery，不是 exactly-once。重複事件是正常現象。
2. **網路故障導致重試**：Consumer 處理成功但 ack 失敗時，事件會被重新投遞。
3. **重播是除錯工具**：需要重播事件來修復資料不一致時，consumer 必須能安全處理重複。
4. **非冪等操作導致資料錯誤**：如果 consumer 每次收到事件都 +1 餘額，重複事件會導致餘額錯誤。

## 何時需要冪等消費者

- **任何修改狀態的事件處理**：更新資料庫、呼叫外部 API、發送通知。
- **金錢相關操作**：付款、退款、餘額調整。
- **庫存調整**：入庫、出庫、預留。
- **任何「只應該發生一次」的操作**：建立帳戶、發送歡迎信。

## 何時不需要冪等消費者

- **純讀取操作**：不修改狀態的事件處理。
- **天然冪等的操作**：`SET status = 'confirmed' WHERE id = ?`（重複執行結果相同）。
- **日誌記錄**：重複的 log entry 可以接受。

## 實作模式

```text
1. 去重表（Deduplication Table）
   在 consumer 中維護 processed_event_ids 表
   收到事件 → 檢查 event_id 是否已處理 → 已處理則跳過

2. 樂觀鎖（Optimistic Locking）
   使用 version 或 updated_at 欄位
   UPDATE ... WHERE version = ? → 影響行數為 0 表示已處理

3. 冪等 Key（Idempotency Key）
   在 API 層使用 idempotency key
   相同 key 的請求只處理一次

4. 條件更新（Conditional Update）
   UPDATE ... WHERE status = 'pending'
   只有 pending 狀態的記錄才會被更新
```

## 決策流程

```text
設計 event consumer？
  ├── consumer 是否修改狀態？
  │     ├── 否 → 不需要冪等
  │     └── 是 → 需要冪等
  ├── 操作是否天然冪等？
  │     ├── 是 → 不需要額外處理
  │     └── 否 → 選擇冪等策略
  ├── 使用去重表？
  │     ├── 優點：簡單、通用
  │     └── 缺點：需要清理過期 event_id
  ├── 使用樂觀鎖？
  │     ├── 優點：不需要額外表格
  │     └── 缺點：需要 version 欄位
  └── 使用冪等 key？
        ├── 優點：API 層統一處理
        └── 缺點：需要 client 支援
```

## 常見誤用

| 誤用 | 正確 |
|------|------|
| 「Exactly-once delivery 可以解決」 | Exactly-once 在分散式系統中成本極高，且通常只是 at-least-once + 去重的組合 |
| 「重複事件很少發生，不用處理」 | 重複事件在生產環境中比想像中常見。不處理 = 資料不一致 |
| 「資料庫 unique constraint 就夠了」 | Unique constraint 只能防止完全相同的資料，無法處理「相同事件、不同時間」的場景 |

## Token Impact

非冪等 consumer 可能導致資料不一致，需要手動修復。每個不一致的修復成本約 1-4 小時（診斷 + 修復 + 驗證）。在生產環境中，這類問題通常是 P1 級 incident。

---

← [回到 engineering/distributed-systems/](README.md)
