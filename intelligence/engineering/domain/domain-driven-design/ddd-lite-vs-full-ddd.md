# DDD Lite vs Full DDD

**Status**: `candidate-intelligence`

## 判斷原則

DDD Lite 保留 domain language、boundary awareness、critical invariant protection 與 explicit domain modeling；Full DDD 才引入完整 bounded contexts、aggregate consistency boundaries、ACL 與 event-driven coordination。

## 選擇矩陣

| 複雜度 | 適合 | 避免 |
| --- | --- | --- |
| 低 | CRUD、feature modules、simple service layer | full DDD、CQRS、event sourcing |
| 中 | DDD Lite、selective aggregate、modular monolith | microservices by default |
| 高 | Full DDD、bounded contexts、ACL、domain events | single global model |

## DDD Lite 最小集

- Ubiquitous language。
- Domain boundary awareness。
- Critical invariant protection。
- Minimal architecture decision record。

## Full DDD 進入條件

只有在 domain complexity、invariant density、integration pressure、lifecycle length 或 team scale 足以支付成本時才使用。
