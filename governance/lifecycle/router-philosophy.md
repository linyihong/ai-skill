# Runtime Context Router Philosophy

## Purpose

`runtime/README.md` 負責 **context routing 決策邏輯**。Agent 在 session 啟動後，透過本層決定哪些知識需要載入、哪些可以 deferred。

## Routing Flow

```
Session Start
  │
  ├─ 1. Load CORE_BOOTSTRAP.md（3 rules, ~800 tokens）
  │
  ├─ 2. Read README.md（OS layout, ~80 lines）
  │
  ├─ 3. Check contract-backed activation → load matching executable contracts
  │
  ├─ 4b. 若命中 activation #27 或 registry route.workflow.* activation_triggers：
  │       強制 routing-philosophy → 比對 activation_triggers → workflow-routing 歧義裁決 → execution-flow.md
  │
  ├─ 5. Read knowledge summary（300-500 tokens）
  │
  ├─ 6. Expand to full source only if needed
  │
  └─ 7. Apply TTL policy → prune context at task/session boundary
```

## Routing Decision Table

| 階段 | 輸入 | 輸出 | Token 成本 |
| --- | --- | --- | --- |
| 1. Bootstrap | Session start | Core rules (3) | ~800 |
| 2. Layout | Task intent | OS map | ~200 |
| 3. Skill routing | Task intent keywords | Skill ID + entrypoint | ~100 |
| 4. Contract activation | Task context | Matching executable contracts | ~200 |
| 5. Summary | Skill ID | 300-500 token summary | ~400 |
| 6. Full source | Summary match | Full document | variable |

## Contract Activation Usage

Owner-layer executable YAML contracts 是 rule activation 的 source-of-truth。Agent 應比對 contract 的 `activation`、routing registry 的 `activation_triggers`，以及 task intent、file change、user signal 與 validation gap；不要引用已移除的 Ruby activation engine 或舊 `activation_rules` tables。

Go CLI 若未來新增 activation query 命令，應以 owner-layer executable contracts、`generated_surfaces` 與 `knowledge/runtime/routing-registry.yaml` 為 source，並先補 validation fixture。

## 與既有文件的關係

- [`runtime/README.md`](../../runtime/README.md) — Runtime navigation entry point (contract projections, `activation-table.md`)
- [`runtime/runtime.db`](../../runtime/runtime.db) — Contract projections and runtime config tables
- [`runtime/router/activation-table.md`](../../runtime/router/activation-table.md) — Situation → Activated Rules table
- [`knowledge/runtime/routing-registry.yaml`](../../knowledge/runtime/routing-registry.yaml) — Machine-readable routing records
- [`knowledge/indexes/README.md`](../../knowledge/indexes/README.md) — Human-readable task intent routing table
- [`knowledge/runtime/routing-registry.yaml`](../../knowledge/runtime/routing-registry.yaml) — Machine-readable routing registry
- [`runtime/runtime.db`](../../runtime/runtime.db) — Context 生命週期管理
