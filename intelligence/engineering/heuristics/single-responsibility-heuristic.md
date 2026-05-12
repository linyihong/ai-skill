# Single Responsibility Heuristic（單一職責經驗法則）

**Status**: `candidate-intelligence`
**Source**: 通用軟體工程經驗、SOLID 原則實務

## 原則

**If you can't describe what a module does without using "and", it has too many responsibilities.**

如果你無法在不用「和」的情況下描述一個模組的職責，它就有太多責任。

## 為什麼

1. 單一職責不是「只做一件事」，而是「只有一個改變的理由」。
2. 當一個 class 或 function 有多個改變理由時，任何一個需求的變更都可能影響到其他不相關的行為。
3. 多職責的模組難以測試、難以重用、難以理解。

## 何時適用

- Code review 中看到 class 名稱包含 `And`、`Manager`、`Util`、`Helper`。
- 一個 method 超過 20 行且做了多件事（validation + transformation + persistence + notification）。
- 修改一個功能時，不小心 break 了另一個不相關的功能。

## 何時不適用

- 簡單的 glue code 或 script（例如 migration script、one-off data fix）。
- 效能關鍵的 hot path（有時合併職責是為了避免 function call overhead）。
- Framework 要求的 entry point（例如 Controller、Activity、Fragment），這些本來就是 framework 層級的職責聚合點。

## 判斷方法

### 命名測試

如果 class 名稱需要包含「和」才能描述它的職責，就代表需要拆分：

| 壞名稱 | 好名稱 |
|--------|--------|
| `OrderManager`（管理訂單和庫存和通知） | `OrderProcessor`、`InventoryReserver`、`NotificationSender` |
| `UserHelper`（驗證用戶和發送 email 和產生報表） | `UserValidator`、`EmailService`、`ReportGenerator` |
| `DataUtils`（解析 CSV 和壓縮檔案和計算 hash） | `CsvParser`、`FileCompressor`、`HashCalculator` |

### 改變理由測試

問自己：「這個 class 有幾個改變理由？」

```text
class InvoiceService {
  // 改變理由 1: 發票格式變更
  generatePdf(invoice) { ... }

  // 改變理由 2: Email 模板變更
  sendEmail(invoice) { ... }

  // 改變理由 3: 資料庫 schema 變更
  saveToDatabase(invoice) { ... }
}
```

三個改變理由 → 需要拆分。

## Token Impact

避免 god class 造成的維護災難。一個 2000 行的 class 修改任何功能都需要讀完整個檔案，context 成本是拆分後的 5-10 倍。

---

← [回到 engineering/heuristics/](README.md)
