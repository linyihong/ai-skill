# Engineering Architecture Intelligence

放**架構思考模式**，不是架構教學。

## 核心

架構判斷力 — 不是理論介紹。

## 範例內容

- `modular-monolith-vs-microservices.md` — Microservices increase operational complexity rapidly. Prefer modular monolith unless independent scaling, separate deployment cadence, or team autonomy is required.
- `event-driven-tradeoffs.md` — Event-driven systems improve decoupling but increase debugging complexity and eventual consistency risk.
- `cqrs-when-not-to-use.md` — Avoid CQRS for CRUD-heavy low-complexity systems.
- [`architecture-selection/`](architecture-selection/README.md) — Architecture fit analysis、complexity evaluation、team topology、integration pressure 與 tradeoff matrix，用來避免 default architecture。

## 與其他層的關係

- DDD / 業務模型 tactical heuristics → `intelligence/engineering/domain/domain-driven-design/`
- Software-delivery 中的實際架構決策流程 → `workflow/software-delivery/architecture/`
- Architecture governance gate → `governance/ai-runtime-governance/software-delivery-architecture-governance.md`
- 架構事實與參考資料 → `knowledge/`
