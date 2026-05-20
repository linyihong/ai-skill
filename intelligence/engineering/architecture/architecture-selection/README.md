# Architecture Selection Intelligence

本目錄保存 architecture selection 的判斷智慧。它回答「何時使用 CRUD、vertical slice、DDD Lite、Full DDD、modular monolith、event-driven 或 microservices」，而不是把任一架構當預設。

## 核心原則

Architecture complexity must match business complexity.

Agent 產出架構建議前，必須先說明：

- business complexity。
- invariant density。
- lifecycle length。
- team / ownership boundary。
- integration pressure。
- delivery speed priority。
- 被拒絕的更複雜或更簡單方案。

## 目前條目

| 文件 | 用途 |
| --- | --- |
| [`architecture-fit-analysis.md`](architecture-fit-analysis.md) | 以 evidence 決定架構路線。 |
| [`complexity-evaluation.md`](complexity-evaluation.md) | 評估 business / delivery / technical complexity。 |
| [`lifecycle-evaluation.md`](lifecycle-evaluation.md) | 判斷 MVP、成長期、長期平台的不同架構成本。 |
| [`team-topology-analysis.md`](team-topology-analysis.md) | 以團隊 ownership 判斷拆分是否合理。 |
| [`integration-pressure-analysis.md`](integration-pressure-analysis.md) | 判斷外部系統壓力與 ACL / event coordination 需求。 |
| [`architecture-tradeoff-matrix.md`](architecture-tradeoff-matrix.md) | CRUD / DDD Lite / Full DDD / event-driven / microservices 的取捨表。 |

## 與 DDD 的關係

DDD 是候選策略之一。若 fit analysis 不支持，agent 應主動推薦更簡單的架構。
