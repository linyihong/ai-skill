# PostgreSQL vs MongoDB（關聯式 vs 文件式資料庫取捨）

**Status**: `candidate-intelligence`
**Source**: 通用資料庫選型經驗

## 原則

**MongoDB improves schema flexibility but complicates transactional consistency. PostgreSQL is almost always the better default for general-purpose applications.**

MongoDB 提升 schema 靈活性，但增加交易一致性複雜度。PostgreSQL 幾乎總是通用型應用更好的預設選擇。

## 為什麼

1. **PostgreSQL 的 JSONB 已足夠靈活**：PostgreSQL 的 JSONB 欄位支援索引、查詢、部分更新，涵蓋大部分 MongoDB 的 use case，同時保留關聯式資料庫的 ACID 保證。
2. **交易一致性是硬需求**：大部分應用最終需要跨 document/table 的一致性。MongoDB 的 multi-document transaction 在 4.0 後才支援，且效能與 PostgreSQL 有差距。
3. **查詢複雜度**：PostgreSQL 的 JOIN、window functions、CTE、trigger、view 讓複雜查詢可以在資料庫層完成；MongoDB 的 aggregation pipeline 在複雜查詢時可讀性與效能都較差。
4. **生態成熟度**：PostgreSQL 的 ORM、migration tool、monitoring、backup 工具遠比 MongoDB 成熟。
5. **Schema 變更是雙面刃**：MongoDB 的 schema-less 在初期開發快速，但在後期維護時，隱含的 schema 約束反而比顯式 schema 更難管理。

## 何時適用 MongoDB

- **Prototype / MVP 階段**：schema 尚未穩定，需要快速迭代。
- **大量非結構化資料**：log、event、analytics 等 schema 變化頻繁的資料。
- **文件導向的資料模型**：資料天然是 document 結構（如 CMS content、product catalog），且跨 document 的關聯很少。
- **水平擴展需求明確**：寫入量超過單一 PostgreSQL instance 的處理能力，且 sharding 是必要架構。

## 何時不適用 MongoDB

- **需要跨 document/table 的交易一致性**：金融、庫存、訂單系統。
- **複雜的關聯查詢**：多個 entity 之間的 JOIN 是核心需求。
- **資料完整性要求高**：foreign key、unique constraint、check constraint 是必要功能。
- **團隊熟悉 SQL**：如果團隊已經熟悉 PostgreSQL，MongoDB 的學習曲線不值得。

## 決策流程

```text
選擇資料庫？
  ├── 資料模型是否天然是 document？
  │     ├── 否 → PostgreSQL
  │     └── 是 → 繼續評估
  ├── 需要跨 document 交易一致性？
  │     ├── 是 → PostgreSQL（或 MongoDB 4.0+，但需評估效能）
  │     └── 否 → 繼續評估
  ├── 查詢是否需要 JOIN / 複雜關聯？
  │     ├── 是 → PostgreSQL
  │     └── 否 → 繼續評估
  ├── 寫入量是否超過單一 instance 極限？
  │     ├── 否 → PostgreSQL
  │     └── 是 → MongoDB（或 PostgreSQL + read replica / sharding）
  └── 以上皆不符合？
        ├── 不確定 → PostgreSQL（預設選擇）
        └── 確定需要 document DB → MongoDB
```

## 常見誤用

| 誤用 | 正確 |
|------|------|
| 「MongoDB 比較快」 | 在正確的 use case 下 MongoDB 可能更快，但大部分 OLTP 場景 PostgreSQL 效能相當或更好 |
| 「Schema-less 比較靈活」 | Schema-less 在初期快速，但在後期維護時，隱含的 schema 約束比顯式 schema 更難管理 |
| 「先用 MongoDB，以後再遷移到 PostgreSQL」 | 資料遷移成本極高。如果預期最終需要關聯式資料庫，一開始就用 PostgreSQL |

## Token Impact

避免在專案初期選擇錯誤的資料庫，導致後期高昂的遷移成本。一個錯誤的資料庫選擇可能花費 2-6 週遷移，且過程中可能遺失資料或中斷服務。

---

← [回到 engineering/tradeoffs/](README.md)
