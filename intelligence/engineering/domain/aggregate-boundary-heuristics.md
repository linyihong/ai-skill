# Aggregate Boundary Heuristics（聚合邊界經驗法則）

**Status**: `candidate-intelligence`
**Source**: DDD 實戰經驗

## 原則

**If transaction scope crosses multiple business invariants, aggregate boundary may be incorrect. An aggregate should be a consistency boundary, not a data grouping.**

如果交易範圍跨越了多個業務不變量，聚合邊界可能不正確。聚合應該是一致性邊界，而不是資料分組。

## 為什麼

1. **聚合的目的是保證不變量**：每個聚合負責保護一組業務不變量（business invariants）。如果一個交易需要修改多個聚合才能保證不變量，表示邊界劃分錯誤。
2. **過大的聚合導致效能問題**：一個聚合包含太多 entity 時，每次載入都需要讀取大量資料，且並發修改的衝突機率增加。
3. **過小的聚合導致一致性問題**：如果業務不變量需要跨多個聚合保證，最終一致性可能導致資料不一致。
4. **聚合邊界影響領域語言的清晰度**：好的聚合邊界讓領域語言自然對應到程式碼結構。

## 何時懷疑聚合邊界錯誤

- **交易需要鎖定多個聚合**：如果一個 use case 需要同時修改 2 個以上的聚合，且需要 ACID 保證，邊界可能錯了。
- **聚合包含太多 entity**：一個聚合包含 10+ 個 entity，且大部分 entity 很少一起修改。
- **聚合的 invariant 難以描述**：如果無法用一句話描述這個聚合保護什麼不變量，邊界可能模糊。
- **並發修改頻繁衝突**：多個使用者同時修改同一個聚合的不同部分，導致頻繁的樂觀鎖衝突。

## 何時不懷疑聚合邊界錯誤

- **聚合很小（1-3 個 entity）且 invariant 明確**：這是健康的聚合。
- **跨聚合的交易使用 eventual consistency 且業務可接受**：這是 DDD 的正確做法。
- **聚合的修改頻率低**：即使聚合較大，如果修改頻率低，也不是問題。

## 決策流程

```text
檢查聚合邊界？
  ├── 能否用一句話描述這個聚合的 invariant？
  │     ├── 是 → 繼續
  │     └── 否 → 邊界可能模糊，需要重新劃分
  ├── 交易是否需要修改多個聚合？
  │     ├── 是 → 檢查 invariant 是否真的跨聚合
  │     │     ├── 是 → 合併聚合或使用 saga
  │     │     └── 否 → 使用 eventual consistency
  │     └── 否 → 繼續
  ├── 聚合包含多少 entity？
  │     ├── > 10 → 考慮拆分
  │     ├── 4-10 → 檢查 entity 是否經常一起修改
  │     └── 1-3 → 健康
  └── 並發修改衝突頻率高？
        ├── 是 → 考慮拆分聚合或調整一致性模型
        └── 否 → 邊界可能正確
```

## 常見誤用

| 誤用 | 正確 |
|------|------|
| 「一個 Order 包含所有 OrderItem」 | Order 和 OrderItem 通常在同一聚合（Order 的 invariant 是總金額 = item 加總），但如果 OrderItem 需要獨立修改，應拆分 |
| 「User 和 Address 在同一聚合」 | User 和 Address 通常在不同聚合（Address 可以獨立修改，不影響 User 的 invariant） |
| 「為了效能把多個聚合合併」 | 聚合邊界是業務決策，不是效能決策。效能問題應透過 cache、read model 或 CQRS 解決 |

## Token Impact

錯誤的聚合邊界可能導致：
- 過大：每次載入浪費 2-5 倍不必要的資料讀取
- 過小：需要複雜的 saga 或 eventual consistency 邏輯
- 修正成本：重新劃分聚合邊界可能需要重構 3-10 個 use case

---

← [回到 engineering/domain/](README.md)
