# Engineering Tradeoffs

放**技術取捨智慧**。

## 核心

「沒有銀彈」的理解。

## 目前內容

| 文件 | 用途 |
| --- | --- |
| [`ddd-cost-model.md`](ddd-cost-model.md) | 判斷 DDD 成本是否被 domain complexity 支付。 |
| [`bdd-cost-model.md`](bdd-cost-model.md) | 判斷 BDD-lite / full BDD runner 的成本邊界。 |
| [`architecture-overengineering.md`](architecture-overengineering.md) | 偵測 architecture / DDD overengineering。 |
| [`delivery-friction-analysis.md`](delivery-friction-analysis.md) | 分析 requirements、architecture、validation 與 recovery 的 delivery friction。 |

## 範例內容

- `postgres-vs-mongodb.md` — MongoDB improves schema flexibility but complicates transactional consistency.
- `websocket-vs-polling.md` — WebSocket reduces latency but increases connection management complexity.
- `sqlite-vs-postgres.md` — SQLite is ideal for local-first AI runtime systems with low operational overhead.

## 與其他層的關係

- 特定技術的執行方式 → `workflow/` 或 `analysis/`。
- 技術事實與參考資料 → `knowledge/`。
