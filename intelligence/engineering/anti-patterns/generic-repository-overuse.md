# Generic Repository Overuse（泛型 Repository 過度使用）

**Status**: `candidate-intelligence`
**Source**: 通用後端開發經驗

## 原則

**Generic repositories often hide query intent and reduce performance visibility. Prefer explicit query methods over generic Find/GetAll patterns.**

泛型 Repository 經常隱藏查詢意圖並降低效能可視性。優先使用明確的查詢方法，而非泛型的 Find/GetAll 模式。

## 為什麼

1. **查詢意圖被隱藏**：`repository.FindAll()` 無法表達「我要找什麼」。呼叫端需要猜測這個查詢的條件、排序、關聯載入策略。
2. **效能問題難以發現**：Generic repository 通常使用 `IQueryable<T>` 或類似機制，讓呼叫端決定查詢範圍。這導致 N+1 query、不必要的欄位載入、缺少 index 的查詢。
3. **測試困難**：Generic repository 的 mock 需要模擬完整的查詢介面，測試程式碼複雜且脆弱。
4. **違反介面隔離原則**：一個 entity 可能只需要 2-3 種查詢方式，但 generic repository 暴露了所有 CRUD 操作。

## 何時適用

- **Prototype / MVP 階段**：快速開發優先，效能與可維護性次要。
- **CRUD-only entity**：entity 只有基本的增刪改查，沒有複雜的查詢邏輯。
- **Admin / internal tool**：效能要求不高，開發速度優先。

## 何時不適用

- **複雜的業務查詢**：查詢需要多個條件、排序、分頁、關聯載入。
- **效能敏感的 hot path**：每個 query 的效能都關鍵。
- **長期維護的專案**：Generic repository 在專案成長後會成為維護瓶頸。

## 替代方案

```text
❌ Generic Repository:
    repository.FindAll()  ← 不知道查詢什麼、載入什麼

✅ Explicit Query Method:
    orderRepository.FindPendingOrdersByCustomer(customerId, page, pageSize)
    ← 明確的查詢意圖、可針對性優化、可測試
```

## 決策流程

```text
使用 Generic Repository？
  ├── 這是 prototype / MVP？
  │     ├── 是 → 可以使用，但標記為 technical debt
  │     └── 否 → 繼續評估
  ├── entity 只有基本 CRUD？
  │     ├── 是 → 可以使用
  │     └── 否 → 需要明確的查詢方法
  ├── 查詢需要多個條件 / 排序 / 分頁？
  │     ├── 是 → 不要用 generic repository
  │     └── 否 → 繼續評估
  └── 這是 hot path？
        ├── 是 → 不要用 generic repository，使用明確的查詢方法 + 針對性優化
        └── 否 → 可以考慮，但需注意後續維護成本
```

## 常見誤用

| 誤用 | 正確 |
|------|------|
| 「Generic repository 減少重複程式碼」 | 減少的是 CRUD boilerplate，但增加了查詢意圖的模糊性。重複的 CRUD 程式碼可以用 code generation 解決 |
| 「用 IQueryable 讓呼叫端決定查詢」 | IQueryable 將查詢邏輯分散到多個呼叫端，難以追蹤效能問題。查詢邏輯應集中在 repository 中 |
| 「Generic repository 方便單元測試」 | Generic repository 的 mock 需要模擬完整查詢介面，反而增加測試複雜度 |

## Token Impact

Generic repository 在專案初期節省少量開發時間，但在專案成長後會持續消耗維護成本。每個 N+1 query 的診斷與修復成本約 1-3 小時。

---

← [回到 engineering/anti-patterns/](README.md)
