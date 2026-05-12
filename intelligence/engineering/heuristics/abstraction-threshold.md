# Abstraction Threshold Heuristic（抽象化門檻經驗法則）

**Status**: `candidate-intelligence`
**Source**: 通用軟體工程經驗

## 原則

**If abstraction removes more clarity than duplication, do not abstract.**

如果抽象化移除的 clarity 少於它引入的複雜度，就不要抽象化。

## 為什麼

1. 抽象化的目的是隱藏複雜度，但如果抽象化本身比原始問題更難理解，它就失敗了。
2. Duplication 是程式碼的壞味道，但 premature abstraction 是更糟的壞味道 — 它隱藏了真實的差異。
3. 三次出現（Rule of Three）之前就抽象化，往往會選錯 abstraction boundary。

## 何時適用

- 同一段邏輯在不同地方出現三次以上，且每次的使用方式幾乎相同。
- 抽象化後的 interface 比原始實作更簡單、更直覺。
- 抽象化能讓新的使用案例變得 trivial（例如加入新的 repository implementation）。

## 何時不適用

- 只有兩處重複，且兩處的 context 或行為有細微差異。
- 抽象化後的 interface 需要大量參數或 config 來處理邊界情況。
- 抽象化引入新的概念（例如 Repository pattern、Factory pattern），而團隊需要花時間學習這些概念才能理解程式碼。

## 決策流程

```text
看到重複程式碼？
  ├── 出現次數 >= 3？
  │     ├── 是 → 考慮抽象化
  │     └── 否 → 繼續觀察
  │
  └── 考慮抽象化？
        ├── 抽象化後的 interface 比原始程式碼更簡單？
        │     ├── 是 → 可以抽象化
        │     └── 否 → 不要抽象化
        │
        └── 抽象化引入新概念？
              ├── 是，且團隊不熟悉 → 先不要，或加註解說明
              └── 否 → 可以抽象化
```

## 常見誤用

| 誤用 | 正確 |
|------|------|
| 兩個 repository 有相似的 CRUD 方法，就用 Generic Repository | Generic Repository 隱藏 query intent，讓 performance optimization 變得困難 |
| 三個 service 都 call 同一支 API，就包成一個共用 client | 如果每個 service 對 response 的處理方式不同，共用 client 會變成 if-else 地獄 |
| 到處都用 `Either` 或 `Result` 型別來處理錯誤 | 如果大部分函數不會失敗，強制 caller 處理 error type 反而降低 readability |

## Token Impact

避免 premature abstraction 產生的技術債。一個錯誤的抽象化可能需要 2-3 次重構才能修正，每次重構成本遠高於保留 duplication。

---

← [回到 engineering/heuristics/](README.md)
