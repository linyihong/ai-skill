# Domain Services

**Status**: `candidate-intelligence`

## 判斷原則

Domain service 承載不自然屬於單一 entity / value object / aggregate 的業務操作。它不是 application service、transaction script 或任意 helper 的新名字。

## 適用訊號

- 操作跨多個 aggregate，但核心規則仍是 domain rule。
- 規則需要 domain language 描述，且不屬於任一單一 object。
- 將規則塞入某個 aggregate 會造成錯誤 ownership。

## 不適用訊號

- 只是 orchestration、I/O、email、queue、HTTP、database transaction。
- 只是避免在 entity 裡寫方法。
- 只是 procedural service layer 的包裝。

## 邊界

Application service 負責 use case orchestration；domain service 負責業務規則。若 service 只呼叫 repository 和 adapter，它通常不是 domain service。
