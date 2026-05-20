# Aggregate Explosion

**Status**: `candidate-intelligence`

## 反模式

因為想「更 DDD」而建立過多 aggregate，導致一致性、transaction、event coordination 和理解成本上升。

## 訊號

- Aggregate 無法說明自己保護的不變量。
- 一個 use case 需要同時修改大量 aggregate。
- Domain event 只是為了串接過度切分的模型。
- Aggregate 數量成長快於業務概念數量。

## 修正

合併沒有獨立 invariant 的 aggregate。以 business transaction 和 lifecycle 重新檢查邊界。
