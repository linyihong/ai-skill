# Software Delivery Architecture Workflow

本目錄保存 software-delivery 中的 architecture selection 執行流程。它用來在 design review、planning、greenfield、large refactor 或 long-lived domain system 中做架構 fit analysis。

## 何時使用

- 使用者要求 architecture plan、design review、domain model、DDD、bounded context、CQRS、event sourcing、microservices 或 modular monolith。
- 變更可能影響 domain invariant、long-lived maintainability、team ownership、integration boundary 或 deployment boundary。
- Agent 準備提出比 simple service layer 更複雜的架構。

## 不必使用

- 小型 bug fix。
- 純 UI 文案、樣式、低風險 CRUD。
- 不改變 domain model、contract 或 deployment boundary 的局部 refactor。

## 流程

1. 先讀 [`architecture-fit-analysis.md`](architecture-fit-analysis.md)。
2. 若需要正式 decision，使用 [`architecture-decision-framework.md`](architecture-decision-framework.md)。
3. 若提到 DDD，讀 [`bounded-context-evaluation.md`](bounded-context-evaluation.md) 與 `intelligence/engineering/domain/domain-driven-design/README.md`。
4. 若架構複雜度增加，讀 [`architecture-overengineering-detection.md`](architecture-overengineering-detection.md)。
5. 若 fit mismatch 會阻塞 delivery，依 [`architecture-escalation-policy.md`](architecture-escalation-policy.md) 升級。
6. 所有建議都要通過 [`architecture-minimality-principle.md`](architecture-minimality-principle.md)。

## 輸出

Architecture recommendation 至少包含：chosen strategy、rejected lighter option、rejected heavier option、fit evidence、validation plan、upgrade/downgrade trigger。
