# Software Delivery Architecture Governance

## Source Intelligence

source_intelligence:

- [`intelligence/engineering/architecture/architecture-selection/README.md`](../../intelligence/engineering/architecture/architecture-selection/README.md)
- [`intelligence/engineering/domain/domain-driven-design/README.md`](../../intelligence/engineering/domain/domain-driven-design/README.md)
- [`intelligence/engineering/anti-patterns/architecture-absolutism.md`](../../intelligence/engineering/anti-patterns/architecture-absolutism.md)
- [`workflow/software-delivery/architecture/README.md`](../../workflow/software-delivery/architecture/README.md)

本文件把 architecture fit analysis、DDD Lite / Full DDD selection、architecture minimality 與 overengineering detection 轉譯成 software-delivery governance。它不把 DDD 變成 runtime invariant，也不讓 architecture optimization 阻塞低風險 delivery。

## 觸發時機

在下列情況套用本治理：

- 使用者要求 architecture plan、design review、domain modeling、DDD、bounded context、CQRS、event sourcing 或 microservices。
- Agent 準備建議超過 simple service layer 的架構。
- 變更影響 domain invariant、contract boundary、long-lived maintainability、team ownership 或 integration boundary。
- 發現 architecture overengineering 或 architecture underfitting 的訊號。

## Governance Gate

| Gate | 通過條件 |
| --- | --- |
| Fit analysis | 已評估 domain complexity、invariant density、integration pressure、lifecycle、team scale 與 delivery priority。 |
| Minimality | 已說明更簡單方案為何不足，或選擇降級為簡單方案。 |
| Strategy selection | 已區分 CRUD / DDD Lite / Full DDD / event-driven / microservices。 |
| Rejected alternatives | 已列出 rejected lighter option 與 rejected heavier option。 |
| Overengineering review | 若引入 DDD/CQRS/event sourcing/microservices，已檢查 architecture inflation signal。 |
| Runtime boundary | 沒有把 DDD tactical pattern promotion 成 runtime primitive。 |

## 分層判斷

| 內容類型 | 目標層 |
| --- | --- |
| 為什麼這種架構適合 | `intelligence/engineering/architecture/architecture-selection/` |
| DDD tactical modeling 判斷 | `intelligence/engineering/domain/domain-driven-design/` |
| AI software-delivery governance gate | `governance/ai-runtime-governance/` |
| 實際執行步驟與 decision record | `workflow/software-delivery/architecture/` |
| 可機讀 fit signal | `metadata/architecture/` |
| 可測 failure mode | `validation/scenarios/architecture/` |

## Runtime Boundary

本治理不得直接寫入 `runtime.db` 或新增 phase blocking gate。只有下列壓縮訊號可作為未來 runtime-lite 候選：

- cross-context invariant violation。
- architecture drift causing delivery instability。
- boundary mismatch causing recovery loops。

即使 promotion，也只能是 reliability signal，不是 Full DDD enforcement。

## Validation Candidate

後續 scenario 應測：cargo-cult DDD、architecture fit mismatch、overengineering detection、bounded-context collapse 與 aggregate explosion。
