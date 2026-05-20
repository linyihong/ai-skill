# Engineering Domain Intelligence

放 **DDD / 業務模型智慧**。

## 核心

業務建模經驗。

## 範例內容

- `aggregate-boundary-heuristics.md` — If transaction scope crosses multiple business invariants, aggregate boundary may be incorrect.
- `inventory-domain-patterns.md` — Inventory systems require idempotent adjustments to survive retry scenarios.
- [`domain-driven-design/`](domain-driven-design/README.md) — DDD 作為 selectable architecture strategy 的業務模型判斷智慧，包含 bounded context、aggregate、ubiquitous language、DDD Lite 與 overengineering boundary。

## 與其他層的關係

- 特定業務領域的執行流程 → `workflow/`
- 業務領域的事實與術語 → `knowledge/`
- Architecture selection 與架構取捨 → `intelligence/engineering/architecture/architecture-selection/`
- Software-delivery 中的實際架構決策流程 → `workflow/software-delivery/architecture/`
