# DDD Overengineering Risks

**Status**: `candidate-intelligence`

## 判斷原則

DDD 的錯誤不是「用了 pattern」，而是 pattern 的成本超過 business complexity。

## 風險訊號

- Aggregate 數量快速增加，但 invariant 不明確。
- CQRS 沒有 read/write divergence 或 scale 需求。
- Event sourcing 沒有 audit、temporal query、replay 或法規需求。
- Repository 只有單一 implementation，且不保護 aggregate boundary。
- Microservice split 沒有 deployment、ownership、scaling 或 compliance boundary。
- Abstraction count 超過 domain concept count。

## 必要行動

1. 做 architecture simplification review。
2. 評估是否回到 DDD Lite、feature module 或 simple service layer。
3. 移除 speculative abstraction。
4. 若保留複雜度，記錄 business complexity evidence。
