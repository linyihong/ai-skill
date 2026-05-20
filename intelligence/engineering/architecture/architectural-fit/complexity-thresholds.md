# Complexity Thresholds

**Status**: `candidate-intelligence`

## 分級

| Level | Signal | Strategy |
| --- | --- | --- |
| Low | CRUD、低 invariant、短 lifecycle | CRUD / vertical slice / simple service layer |
| Medium | 有 shared language、部分 invariant、模組邊界 | DDD Lite / modular monolith / selective aggregate |
| High | 多 bounded context、高 invariant、長 lifecycle、多 team / integration | Full DDD / ACL / domain events / consistency boundaries |

## 規則

Architecture complexity 應與 domain complexity、lifecycle length、coordination burden 成比例。
