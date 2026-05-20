# Architecture Escalation

**Status**: `candidate-intelligence`

## 何時升級

Architecture escalation 只在 fit mismatch 會影響 behavior correctness、business invariant correctness 或 long-lived delivery correctness 時使用。

## 升級訊號

- Low complexity 專案被推向 Full DDD / CQRS / event sourcing。
- High invariant domain 被壓成 CRUD-only。
- Requirement ambiguity 尚未解決就開始 architecture design。
- Bounded context mismatch 導致 contract / recovery loop。

## 行動

先回 requirements cognition 或 domain-modeling cognition，不要直接堆更多架構 pattern。
