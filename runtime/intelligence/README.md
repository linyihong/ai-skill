# Intelligence Routing

`runtime/intelligence/` 定義 agent 如何根據目前 phase 與 task context，路由到 `intelligence/` 層的領域知識。這是 runtime 與 intelligence 之間的邊界層。

## 設計原則

1. **Intelligence 留在原處**：`intelligence/` 層的 prose 知識（heuristics、anti-patterns、tradeoffs、failure analysis）不搬進 runtime。Runtime 只存放 routing 規則與 generated YAML。
2. **Phase-aware routing**：不同 phase 需要不同的 intelligence domain。例如 `execution` phase 需要 analytical-reasoning heuristics，`validation` phase 需要 failure patterns。
3. **Deterministic-only routing**：Routing 規則只處理可明確分類的情境（domain match、phase match、task intent match）。無法分類的情境回退到 `intelligence/README.md` 全文搜尋。
4. **Lazy load**：Intelligence 知識只在對應 phase 或 task intent 觸發時才載入，不在 bootstrap 階段預載。

## Intelligence Domain 對照

| Domain | 路徑 | 適用 phase | 觸發條件 |
|--------|------|-----------|---------|
| Analytical Reasoning | `intelligence/engineering/analytical-reasoning/` | execution, validation | task intent 包含分析、逆向、偵測 |
| Agent Architecture | `intelligence/engineering/agent-architecture/` | checkpoint, finalize | task intent 包含架構、規則、設計 |
| Development | `intelligence/engineering/development/` | execution, validation | task intent 包含開發、BDD、測試 |
| Heuristics | `intelligence/engineering/heuristics/` | execution | 通用 heuristic 查詢 |
| Anti-patterns | `intelligence/engineering/anti-patterns/` | validation, finalize | 發現可疑模式或重複問題 |
| Tradeoffs | `intelligence/engineering/tradeoffs/` | execution, checkpoint | 需要技術決策 |
| Distributed Systems | `intelligence/engineering/distributed-systems/` | execution | task intent 包含分散式系統 |
| Domain | `intelligence/engineering/domain/` | execution | task intent 包含領域建模 |
| Failure | `intelligence/engineering/failure/` | validation, recovery | blocking gate 阻斷或驗證失敗 |
| Travel | `intelligence/travel/` | execution | task intent 包含旅行規劃 |
| Business | `intelligence/business/` | execution | task intent 包含商業分析 |
| IDE | `intelligence/ide/` | execution | task intent 包含 IDE 設定 |

## Routing 流程

```
Task Intent Identified
  │
  ├─ 1. Match domain via intelligence-routing.yaml
  │
  ├─ 2. Check current phase → filter applicable domains
  │
  ├─ 3. Load domain README.md (summary, ~200 tokens)
  │
  ├─ 4. If needed → load specific heuristic/pattern/tradeoff file
  │
  └─ 5. If no match → fallback to intelligence/README.md full search
```

## 與既有層的關係

- `runtime/intelligence/intelligence-routing.yaml`：machine-readable routing rules
- `intelligence/README.md`：intelligence layer overview
- `knowledge/runtime/routing-registry.yaml`：global routing registry（含 intelligence routes）
- `runtime/router/activation-rules.yaml`：lazy-load rule activation
- `runtime/phases/phase-machine.yaml`：phase definition（決定 intelligence loading 時機）
