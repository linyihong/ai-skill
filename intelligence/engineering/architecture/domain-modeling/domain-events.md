# Domain Events

**Status**: `candidate-intelligence`

## 判斷原則

Domain event 表示已發生且對業務有意義的事實，不是 service callback、資料同步技巧或所有狀態變更的預設包裝。

## 適用訊號

- 事件名稱能用業務語言描述已發生事實。
- 多個 bounded context 需要在不共享模型的情況下協調。
- 事件歷史本身有 audit、replay、通知或補償價值。
- 同步呼叫會造成 context coupling。

## 不適用訊號

- 單一 CRUD update 沒有下游業務意義。
- 只是為了套 event-driven architecture。
- 使用 event 逃避 transaction boundary 設計。
- 事件 payload 直接暴露內部 entity 結構。

## 決策規則

先確認事件是否代表 business fact，再決定 delivery mechanism。Domain event 是語意，message broker 是實作。
