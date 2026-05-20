# CRUD vs Domain Complexity

**Status**: `candidate-intelligence`

## 判斷原則

CRUD 是低 complexity 的好策略，不是落後策略。DDD 是 domain complexity 足夠高時的成本投資，不是預設。

## CRUD 足夠時

- Low invariant density。
- Low workflow orchestration。
- Single actor / simple state。
- Short lifecycle or MVP。
- Low integration pressure。

## CRUD underfit 時

- 多個 business invariant 需要保護。
- State transition 有高風險錯誤。
- 同一名詞在不同流程有不同語意。
- Cross-context workflow 常造成 contradiction。
