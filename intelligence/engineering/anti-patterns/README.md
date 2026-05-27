# Engineering Anti-Patterns

放**常見錯誤設計**。

## 核心

AI 自動避雷。

## 範例內容

- [`generic-repository-overuse.md`](generic-repository-overuse.md) — Generic repositories often hide query intent and reduce performance visibility.
- `shared-database-microservices.md` — Shared DB across services creates hidden coupling.（候選，尚未撰寫）
- `god-service-pattern.md` — Large service classes signal missing domain boundaries.（候選，尚未撰寫）
- [`cargo-cult-ddd.md`](cargo-cult-ddd.md) — DDD pattern 被套用到低 complexity 專案，但沒有 business complexity evidence。
- [`premature-cqrs.md`](premature-cqrs.md) — 沒有 read/write divergence 或 scale requirement 就導入 CQRS。
- [`aggregate-explosion.md`](aggregate-explosion.md) — Aggregate 數量超過 invariant / domain concept，造成一致性與理解成本膨脹。
- [`repository-abstraction-overuse.md`](repository-abstraction-overuse.md) — Repository abstraction 沒有 aggregate boundary 或 query intent 價值。
- [`architecture-absolutism.md`](architecture-absolutism.md) — 把單一架構風格當 universal default。
- [`migration-seeder-anti-patterns.md`](migration-seeder-anti-patterns.md) — 把大量業務資料以巨型 INSERT 包進 schema migration，使資料 lifecycle 被綁定到 schema lifecycle。
- [`session-scoped-implicit-state.md`](session-scoped-implicit-state.md) — 將 user / phase / tenant 等執行上下文塞進 session-scoped service 或 ambient context，對 agent runtime / replay / 並發場景是隱形依賴災難。
- [`migration-feature-bundling.md`](migration-feature-bundling.md) — 大型 migration / rewrite 時把搬遷與新功能綁進同一階段交付，使驗證失去 ground truth；應採 Parity-First Migration（Phase 1 等價、Phase 2 新功能）。
- [`shotgun-debugging.md`](shotgun-debugging.md) — Troubleshooting / fix 階段同時改多個 unrelated 變數，沒 baseline / 沒單變數變更 / 沒 after 對照；應走 measure → process → verify 三段式。AI 輔助下被放大十倍。

## 與根目錄 `anti-patterns/` 的差異

| 位置 | 內容 |
|------|------|
| [`anti-patterns/`](../../../anti-patterns/) | Agent 操作層的 anti-patterns（context explosion、recursive tool loop、hallucination loop 等） |
| `intelligence/engineering/anti-patterns/` | 工程領域的 anti-patterns（設計錯誤、架構錯誤等） |

## 與其他層的關係

- Agent 操作層的 anti-patterns → [`anti-patterns/`](../../../anti-patterns/)
- 具體的 failure pattern policy → `enforcement/failure-patterns/`
