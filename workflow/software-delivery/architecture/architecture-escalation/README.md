# Architecture Escalation

## 目的

當 architecture recommendation 與 complexity evidence 不匹配時，停止直接實作並重做 fit analysis。

## 觸發

- Low complexity project 被推到 Full DDD / CQRS / event sourcing。
- High invariant domain 被壓成 CRUD-only。
- Requirement ambiguity 尚未解決。
- Architecture decision 無 rejected alternatives。
