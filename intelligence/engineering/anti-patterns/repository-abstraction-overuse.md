# Repository Abstraction Overuse

**Status**: `candidate-intelligence`

## 反模式

為了抽象而建立 repository，卻沒有 aggregate boundary、查詢語意或 persistence 交換需求。

## 訊號

- Repository 只有單一 implementation。
- 方法是 `GetAll`、`Find`、`Save` 這類泛型 CRUD。
- 呼叫端仍需知道 ORM / query detail。
- 測試 mock 比真實查詢更難理解。

## 修正

只為 aggregate persistence boundary 或明確 domain query 建 repository。簡單 CRUD 可直接使用資料存取層或 framework convention。
