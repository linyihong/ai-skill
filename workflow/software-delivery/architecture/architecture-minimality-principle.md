# Architecture Minimality Principle

## 原則

選擇能保護目前已知 business complexity 的最小架構。不要為 speculative future complexity 預先付款。

## 必問問題

- 更簡單方案會破壞哪個 business invariant？
- 更重方案的成本是否有 evidence 支持？
- 這個 abstraction 現在是否有第二個 use case？
- 若未來需要升級，migration path 是否可接受？

## 降級規則

若找不到明確 invariant、lifecycle、integration 或 team evidence，降級為 DDD Lite、feature module 或 simple service layer。
