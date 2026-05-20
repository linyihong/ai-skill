# Aggregate Design

**Status**: `candidate-intelligence`

## 判斷原則

Aggregate 是一致性邊界，不是資料組合，也不是 ORM object graph。

既有核心 atom：[`aggregate-boundary-heuristics.md`](../aggregate-boundary-heuristics.md)。本文件補充 DDD integration 的 adoption boundary。

## 適用訊號

- 有明確 business invariant 需要在同一 transaction 內維持。
- use case 需要防止非法狀態轉移。
- 多個 entity 的更新不能被任意拆散。
- 併發衝突代表真實業務衝突，而不是技術實作不佳。

## 不適用訊號

- 只是想把 entity 包起來。
- 每張表都建立 aggregate。
- 只為了配合 repository pattern。
- read-heavy CRUD 頁面沒有關鍵不變量。

## 健康檢查

| 檢查 | 健康訊號 | 風險訊號 |
| --- | --- | --- |
| invariant | 一句話能說明 aggregate 保護什麼 | 只能列 entity 名稱 |
| transaction | 同一 aggregate 內完成 | 經常跨多個 aggregate ACID 更新 |
| size | 小且聚焦 | 為方便查詢持續膨脹 |
| event | 表達已發生的業務事實 | 用 event 取代同步一致性需求 |

## 最小化規則

若 invariant 不明確，先不要設計 aggregate；先回到 ubiquitous language 與 bounded context 分析。
