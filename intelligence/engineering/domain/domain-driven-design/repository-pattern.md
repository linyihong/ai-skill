# Repository Pattern

**Status**: `candidate-intelligence`

## 判斷原則

Repository 應服務 aggregate persistence boundary，而不是為每個 entity 建立 generic CRUD wrapper。

相關 anti-pattern：[`generic-repository-overuse.md`](../../anti-patterns/generic-repository-overuse.md)。

## 適用訊號

- Aggregate 需要以 domain language 查找或保存。
- Persistence detail 會干擾 domain model。
- 查詢方法可以表達 use case intent。
- 測試需要替代 persistence，而不是模擬所有 query operator。

## 不適用訊號

- 單一 implementation、單純 CRUD、無 domain invariant。
- Repository 只轉呼叫 ORM，沒有語意價值。
- 泛型介面暴露 `FindAll`、`Query` 或任意 predicate。

## 決策規則

先問「這個 repository 是否保護 aggregate boundary 或表達 domain query？」如果答案是否，使用較簡單的資料存取方式。
