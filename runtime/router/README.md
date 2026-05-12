# Runtime Context Router

`runtime/router/` 負責 **context routing 決策邏輯**。Agent 在 session 啟動後，透過本層決定哪些知識需要載入、哪些可以 deferred。

## 路由流程

```
Session Start
  │
  ├─ 1. Load CORE_BOOTSTRAP.md（3 rules, ~800 tokens）
  │
  ├─ 2. Read README.md（OS layout, ~80 lines）
  │
  ├─ 3. Query skills-index.yaml → match task intent → find skill
  │
  ├─ 4. Check activation-rules.yaml → load lazy rules if triggered
  │
  ├─ 5. Read knowledge summary（300-500 tokens）
  │
  ├─ 6. Expand to full source only if needed
  │
  └─ 7. Apply TTL policy → prune context at task/session boundary
```

## 路由決策表

| 階段 | 輸入 | 輸出 | Token 成本 |
| --- | --- | --- | --- |
| 1. Bootstrap | Session start | Core rules (3) | ~800 |
| 2. Layout | Task intent | OS map | ~200 |
| 3. Skill routing | Task intent keywords | Skill ID + entrypoint | ~100 |
| 4. Rule activation | Task context | Lazy rules list | ~200 |
| 5. Summary | Skill ID | 300-500 token summary | ~400 |
| 6. Full source | Summary match | Full document | variable |

## 路由檔案

| 檔案 | 用途 |
| --- | --- |
| [`activation-rules.yaml`](activation-rules.yaml) | 定義 lazy-load rules 的觸發條件與優先權 |
| [`cost-budget.yaml`](cost-budget.yaml)（未來） | Session token budget 管理 |

## 與既有層的關係

- `knowledge/runtime/routing-registry.yaml`：machine-readable routing records（atom → source → summary → cost）
- `knowledge/indexes/README.md`：human-readable task intent routing table
- `skills-index.yaml`：skill-level routing index（triggers → entrypoint → summary）
- `runtime/context/ttl-policy.yaml`：context 生命週期管理
