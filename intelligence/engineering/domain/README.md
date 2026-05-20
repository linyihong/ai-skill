# Engineering Domain Intelligence

放業務領域智慧與 legacy domain entry。DDD tactical modeling 已重新歸入 `intelligence/engineering/architecture/domain-modeling/`，因為它在本系統中屬於 domain architecture cognition，而不是與 BDD 平行的方法論目錄。

## 核心

業務建模經驗與 domain-specific lessons。

## 範例內容

- [`aggregate-boundary-heuristics.md`](aggregate-boundary-heuristics.md) — If transaction scope crosses multiple business invariants, aggregate boundary may be incorrect.
- `inventory-domain-patterns.md` — Inventory systems require idempotent adjustments to survive retry scenarios.

## 與其他層的關係

- Requirements / behavior modeling → `intelligence/engineering/requirements/`。
- DDD / domain architecture modeling → `intelligence/engineering/architecture/domain-modeling/`。
- Architecture fit 與架構取捨 → `intelligence/engineering/architecture/architectural-fit/`。
- Software-delivery 中的實際架構決策流程 → `workflow/software-delivery/architecture/`。
- 特定業務領域的執行流程 → `workflow/`。
- 業務領域的事實與術語 → `knowledge/`。
