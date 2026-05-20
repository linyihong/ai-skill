# Invariants

**Status**: `candidate-intelligence`

## 判斷原則

Business invariant 是 domain correctness 的核心。它描述「什麼必須永遠為真」，不是一般 validation rule 或 UI constraint。

## 來源

- Behavior scenarios 中反覆出現的 state transition。
- Product / domain contract 明確禁止的狀態。
- 合規、安全、金流、庫存、權限等 critical rule。
- Bug / incident 顯示曾被破壞的 business rule。

## 行動

每個 high-value invariant 都需要：owner context、觸發操作、非法狀態、validation target、recovery action。

## Runtime 邊界

Runtime 可以接收 `invariant_violation` signal，但不理解 aggregate purity 或 tactical DDD style。
