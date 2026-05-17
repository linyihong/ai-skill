# Runtime Context Router Philosophy

## Purpose

`runtime/router/` 負責 **context routing 決策邏輯**。Agent 在 session 啟動後，透過本層決定哪些知識需要載入、哪些可以 deferred。

## Routing Flow

```
Session Start
  │
  ├─ 1. Load CORE_BOOTSTRAP.md（3 rules, ~800 tokens）
  │
  ├─ 2. Read README.md（OS layout, ~80 lines）
  │
  ├─ 3. Query skills-index.yaml → match task intent → find skill
  │
  ├─ 4. Run activation-engine.rb → load lazy rules if triggered
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
| 4. Rule activation | Task context | Lazy rules list | ~200 |
| 5. Summary | Skill ID | 300-500 token summary | ~400 |
| 6. Full source | Summary match | Full document | variable |

## Activation Engine Usage

[`activation-engine.rb`](../../runtime/router/activation-engine.rb) 是程式化的 activation 判斷工具，接受 task intent、file change、user signal 等輸入，輸出應該 activate 的 rule 列表。

```bash
# 指定 task intent
ruby runtime/router/activation-engine.rb --intent migration

# 指定 file changes
ruby runtime/router/activation-engine.rb --file-changed enforcement/rule-weight.md

# 複合條件
ruby runtime/router/activation-engine.rb --intent migration --file-changed "**/*.md"
```

## 與既有文件的關係

- [`runtime/router/`](../../runtime/router/) — Runtime navigation entry point (data files: `activation-rules.yaml`, `activation-engine.rb`, `activation-table.md`)
- [`runtime/router/activation-rules.yaml`](../../runtime/router/activation-rules.yaml) — Lazy-load rules with activation conditions
- [`runtime/router/activation-engine.rb`](../../runtime/router/activation-engine.rb) — Activation engine implementation
- [`runtime/router/activation-table.md`](../../runtime/router/activation-table.md) — Situation → Activated Rules table
- [`knowledge/runtime/routing-registry.yaml`](../../knowledge/runtime/routing-registry.yaml) — Machine-readable routing records
- [`knowledge/indexes/README.md`](../../knowledge/indexes/README.md) — Human-readable task intent routing table
- [`skills-index.yaml`](../../skills-index.yaml) — Skill-level routing index
- [`runtime/context/ttl-policy.yaml`](../../runtime/context/ttl-policy.yaml) — Context 生命週期管理
