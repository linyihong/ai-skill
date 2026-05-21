# Engineering Architecture Intelligence

放**架構思考模式**，不是架構教學。這一層現在採 cognition-first：domain-modeling、architectural-fit、delivery-alignment、system boundaries、modularity、coupling 與 consistency model，而不是以 DDD 作為 top-level 方法目錄。

## 核心

架構判斷力 — 不是理論介紹。

## 目前入口

| 子目錄 / 文件 | 用途 |
| --- | --- |
| [`domain-modeling/`](domain-modeling/README.md) | Domain architecture cognition：ubiquitous language、bounded context、aggregate boundary、invariant、consistency。 |
| [`architectural-fit/`](architectural-fit/README.md) | Architecture fit cognition：CRUD / DDD Lite / Full DDD / event-driven / microservices 選型與 tradeoff。 |
| [`delivery-alignment/`](delivery-alignment/README.md) | 確認 architecture decision 可落地到 delivery workflow、validation 與 governance。 |
| [`system-boundaries/`](system-boundaries/README.md) | Ownership、deployment、data、security 與 external integration boundary。 |
| [`modularity/`](modularity/README.md) | Module boundary、feature slice、package boundary 與 modular monolith。 |
| [`coupling-tradeoffs/`](coupling-tradeoffs/README.md) | Tight/loose coupling、coordination cost、abstraction cost。 |
| [`consistency-models/`](consistency-models/README.md) | ACID、eventual consistency、compensation、idempotency、ordering。 |
| [`modular-monolith-vs-microservices.md`](modular-monolith-vs-microservices.md) | Microservices increase operational complexity rapidly. |
| [`vendor-integration-architecture.md`](vendor-integration-architecture.md) | 多廠商整合的整合策略選型（Adapter / compile-time module / SPI / out-of-process / hybrid）；N ≥ 10 必須跳出 compile-time module per vendor。 |

## 與其他層的關係

- Requirements cognition → `intelligence/engineering/requirements/`。
- Software-delivery 中的實際架構決策流程 → `workflow/software-delivery/architecture/`。
- Architecture governance gate → `governance/ai-runtime-governance/software-delivery-architecture-governance.md`。
- 架構事實與參考資料 → `knowledge/`。
