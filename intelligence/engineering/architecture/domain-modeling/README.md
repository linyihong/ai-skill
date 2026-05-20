# Domain Modeling Architecture Intelligence

本目錄保存 domain architecture cognition：domain complexity、ubiquitous language、bounded context、invariant、transactional consistency 與 domain events。這裡不是 DDD 方法目錄，而是處理 domain architecture failure 的知識層。

## 核心

DDD 是 domain complexity and consistency boundary management system。

## 目前條目

| 文件 | 用途 |
| --- | --- |
| [`ubiquitous-language.md`](ubiquitous-language.md) | 需求語言穩定後，建立 domain language。 |
| [`bounded-context.md`](bounded-context.md) | 判斷語言與模型邊界。 |
| [`aggregate-boundary.md`](aggregate-boundary.md) | 判斷 aggregate 是否保護真正的不變量。 |
| [`invariants.md`](invariants.md) | 定義 business invariant 與 validation target。 |
| [`transactional-consistency.md`](transactional-consistency.md) | 判斷同步一致性與 eventual consistency。 |
| [`domain-events.md`](domain-events.md) | 判斷 domain event 是否代表業務事實。 |
| [`anti-corruption-layer.md`](anti-corruption-layer.md) | 防止外部模型污染 internal domain language。 |
| [`domain-services.md`](domain-services.md) | 判斷跨 aggregate 的 domain rule ownership。 |
| [`repository-pattern.md`](repository-pattern.md) | 判斷 persistence boundary 是否有 domain value。 |
| [`tactical-vs-strategic-design.md`](tactical-vs-strategic-design.md) | 避免 tactical pattern 先於 strategic boundary。 |
| [`event-storming.md`](event-storming.md) | 用 planning method 發現 language / flow / boundary。 |

## 上下游

- 上游：`intelligence/engineering/requirements/behavior-modeling/` 先穩定 behavior language。
- 下游：`intelligence/engineering/architecture/architectural-fit/` 決定 CRUD / DDD Lite / Full DDD / event-driven 等 strategy。
