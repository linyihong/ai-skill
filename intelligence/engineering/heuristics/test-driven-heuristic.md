# Test-Driven Heuristic（測試驅動經驗法則）

**Status**: `candidate-intelligence`
**Source**: 通用軟體工程經驗、TDD 實務

## 原則

**If writing a test for a function is difficult, the function's design is likely wrong.**

如果為一個函數寫測試很困難，這個函數的設計很可能有問題。

## 為什麼

1. 測試困難通常不是測試本身的問題，而是被測試的程式碼有設計問題：
   - 太多依賴（需要 mock 太多東西）。
   - 副作用隱藏在函數內部（I/O、DB、network call 混在邏輯中）。
   - 職責不清（一個函數做太多事）。
2. 測試是設計的回饋機制 — 如果測試難寫，代表設計需要改善。
3. 強迫先寫測試（TDD）自然會導向更好的設計，因為測試需求會強制你思考 interface 與職責邊界。

## 何時適用

- 寫 unit test 時需要 mock 5 個以上的 dependencies。
- 測試 setup 程式碼比測試本身還長。
- 需要用到 `PowerMock` 或類似的 bytecode manipulation 工具才能測試。
- 測試需要檢查 private method 的內部狀態。

## 何時不適用

- Integration test 或 end-to-end test — 這些本來就需要較多的 setup。
- Legacy code 的 characterization test — 先寫測試捕捉行為，再重構。
- 效能測試或 benchmark — 這些是另一種測試類型，不適用同樣的設計回饋原則。

## 決策流程

```text
寫測試遇到困難？
  ├── 需要 mock > 3 個 dependencies？
  │     ├── 是 → 檢查是否違反 Dependency Inversion Principle
  │     └── 否 → 繼續
  │
  └── 測試 setup 比測試本身長？
        ├── 是 → 函數可能違反 Single Responsibility Principle
        └── 否 → 繼續

如果以上任一為「是」：
  1. 先重構被測試的程式碼，再寫測試
  2. 不要為了讓測試通過而使用 PowerMock 或反射
```

## 常見誤用

| 誤用 | 正確 |
|------|------|
| 用 PowerMock mock static method | Static method 是設計問題，不是測試問題。考慮改用 dependency injection |
| 測試需要知道 private method 的實作細節 | 測試應該只測試 public interface；private method 是實作細節 |
| 測試需要建立複雜的物件圖 | 考慮用 Test Data Builder 或 Factory，但同時檢查是否 class 職責過多 |

## Token Impact

避免測試維護成本爆炸。一個需要 50 行 setup 的測試，維護成本是 5 行 setup 測試的 10 倍。更重要的是，難寫的測試會讓開發者放棄寫測試，導致 regression 風險上升。

---

← [回到 engineering/heuristics/](README.md)
