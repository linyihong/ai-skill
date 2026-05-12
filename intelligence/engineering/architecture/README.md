# Engineering Architecture Intelligence

放**架構思考模式**，不是架構教學。

## 核心

架構判斷力 — 不是理論介紹。

## 範例內容

- `modular-monolith-vs-microservices.md` — Microservices increase operational complexity rapidly. Prefer modular monolith unless independent scaling, separate deployment cadence, or team autonomy is required.
- `event-driven-tradeoffs.md` — Event-driven systems improve decoupling but increase debugging complexity and eventual consistency risk.
- `cqrs-when-not-to-use.md` — Avoid CQRS for CRUD-heavy low-complexity systems.

## 與其他層的關係

- 具體架構實作方式 → `skills/`
- 架構事實與參考資料 → `knowledge/`
