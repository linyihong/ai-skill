# Context TTL & Pruning Philosophy

## Core Principles

1. **Default: session-only** — 多數 context 只活一個 task/session，session 結束後 prune。
2. **Explicit TTL** — 每條 context 明確定義 `session` / `task` / `conversation` 三種生命週期。
3. **Summary-first** — 先讀輕量 summary（300-500 tokens），需要才展開全文。
4. **Prune at boundary** — 在 task 完成、session 結束、或 token 接近門檻時 prune。

## TTL Types

| TTL | 生命週期 | 範例 |
| --- | --- | --- |
| `session` | 1 task/session | Core Bootstrap、activation rules、architecture roadmap |
| `task` | N tasks | Skill workflow、feedback lesson、analysis method |
| `conversation` | 整個對話 | Routing registry、knowledge summary、TTL policy 本身 |

## Prune Strategy

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

## 與既有文件的關係

- [`governance/ai-runtime-governance/context-attention-governance.md`](../ai-runtime-governance/context-attention-governance.md) — 將 context collapse / attention budgeting intelligence 轉譯成 context loading、recap 與 prune gate
- [`runtime/context/`](../../runtime/context/) — Runtime navigation entry point (data file: `ttl-policy.yaml`)
- [`runtime/runtime.db`](../../runtime/runtime.db) — TTL policy definitions
- [`runtime/runtime.db`](../../runtime/runtime.db) — 決定哪些 rules 需要載入（影響 context 總量）
- [`knowledge/summaries/`](../../knowledge/summaries/) — 輕量 summary 可頻繁載入而不需 prune
- [`metadata/schema.md`](../../metadata/schema.md) — `context_cost.ttl` 欄位定義每個 atom 的 TTL
