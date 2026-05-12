# Context TTL & Pruning

`runtime/context/` 負責 **context 生命週期管理**。目標是防止 context 永久留在 context graph 導致 token 膨脹。

## 核心原則

1. **Default: session-only** — 多數 context 只活一個 task/session，session 結束後 prune。
2. **Explicit TTL** — 每條 context 明確定義 `session` / `task` / `conversation` 三種生命週期。
3. **Summary-first** — 先讀輕量 summary（300-500 tokens），需要才展開全文。
4. **Prune at boundary** — 在 task 完成、session 結束、或 token 接近門檻時 prune。

## TTL 類型

| TTL | 生命週期 | 範例 |
| --- | --- | --- |
| `session` | 1 task/session | Core Bootstrap、activation rules、architecture roadmap |
| `task` | N tasks | Skill workflow、feedback lesson、analysis method |
| `conversation` | 整個對話 | Routing registry、knowledge summary、TTL policy 本身 |

詳細定義：[`ttl-policy.yaml`](ttl-policy.yaml)

## Prune 策略

```
Trigger conditions:
  - Token threshold reached（>80% of session budget）
  - Task boundary（task completed）
  - Session boundary（new session）

Prune order:
  1. Deferred rules（never activated）
  2. Completed task context（TTL: task=1, expired）
  3. Old feedback lessons（TTL: task=1, expired）
  4. Architecture docs（TTL: session=1, can reload next session）

Preserve:
  - Core Bootstrap（always needed）
  - Agent goals（active task context）
  - Conversation goal ledger（session scope）
```

## 與既有層的關係

- `runtime/router/activation-rules.yaml`：決定哪些 rules 需要載入（影響 context 總量）
- `knowledge/summaries/`：輕量 summary 可頻繁載入而不需 prune
- `metadata/schema.md`：`context_cost.ttl` 欄位定義每個 atom 的 TTL
