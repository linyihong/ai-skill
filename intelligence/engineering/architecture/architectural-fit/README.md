# Architectural Fit Intelligence

本目錄保存 architecture fit cognition：根據 ambiguity、domain complexity、invariant density、workflow consistency、coordination burden 與 lifecycle pressure 動態選 strategy。

## 核心原則

不要 framework cosplay。先判斷 complexity，再選 CRUD、DDD Lite、Full DDD、event-driven 或 microservices。

## 目前條目

| 文件 | 用途 |
| --- | --- |
| [`architecture-fit-analysis.md`](architecture-fit-analysis.md) | 以 evidence 決定架構路線。 |
| [`ddd-fit-analysis.md`](ddd-fit-analysis.md) | 判斷是否需要 DDD / DDD Lite / Full DDD。 |
| [`crud-vs-domain-complexity.md`](crud-vs-domain-complexity.md) | 判斷 CRUD 何時足夠、何時 underfit。 |
| [`complexity-evaluation.md`](complexity-evaluation.md) | 評估 business / delivery / technical complexity。 |
| [`complexity-thresholds.md`](complexity-thresholds.md) | 將 complexity signal 分級。 |
| [`architecture-escalation.md`](architecture-escalation.md) | fit mismatch 何時升級。 |
| [`architecture-tradeoff-matrix.md`](architecture-tradeoff-matrix.md) | 架構策略取捨表。 |
| [`ddd-lite-vs-full-ddd.md`](ddd-lite-vs-full-ddd.md) | DDD Lite / Full DDD selection。 |
| [`team-topology-analysis.md`](team-topology-analysis.md) | ownership / team boundary。 |
| [`integration-pressure-analysis.md`](integration-pressure-analysis.md) | external model / coordination pressure。 |
| [`lifecycle-evaluation.md`](lifecycle-evaluation.md) | lifespan and evolution pressure。 |

## 輸出要求

Architecture recommendation 必須列出 chosen strategy、rejected lighter option、rejected heavier option、fit evidence、validation plan 與 upgrade/downgrade trigger。
