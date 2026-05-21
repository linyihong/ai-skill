# Cross-Layer Identifier Stability（跨層識別碼穩定性經驗法則）

**Status**: `candidate-intelligence`
**Source**: 跨層 refactor 痛苦樣本（多模組單體 Web App 三層對齊、governance runtime 跨 SQLite/YAML/Markdown drift 觀察）

## 原則

**當一個概念跨 ≥ 3 個 artifact 層存在時，stable opaque ID（數字、UUID、slug）優於 semantic naming。語意會 drift，ID 不會。**

可讀性 -1，可追蹤性 +5。

## 為什麼

1. Semantic name 反映「當下對概念的理解」，理解會演化（拆分、合併、改名）。每次語意演化都要跨層同步重命名，成本隨層數線性增加。
2. Opaque ID 與語意脫鉤，名字改了 ID 不變，跨層 trace 仍然可行。
3. 多層 drift 一旦發生，事後對齊成本通常是「一次性建立 ID 慣例」的 10 倍以上，且 drift 期間的 bug 無法歸因。

## 何時適用

- 一個概念出現在 ≥ 3 處：例如 DB schema、config file、code constant、文件 reference、UI label。
- 跨層命名容易被不同 reviewer 在不同時間獨立改動。
- 概念本身預期會演化（拆分 / 合併 / 改名）。
- 需要長期 traceability（audit、replay、incident 歸因）。

## 何時不適用

- 概念只活在單一檔案 / 單一層內。
- 命名已經是穩定行業術語（HTTP status code 名稱、SQL 關鍵字）。
- 對外 public API：ID 對使用者無意義，semantic name 才是契約。
- Throwaway prototype，3 個月內會丟掉。

## 判斷流程

```text
這個概念跨幾層 artifact？
  ├── 1 層 → semantic name 即可
  ├── 2 層 → semantic name + 一致性 lint
  └── ≥ 3 層 →
        ├── 概念會演化？
        │     ├── 是 → 加 stable opaque ID 當 join key，semantic name 當 display
        │     └── 否 → semantic name 可接受，但加 cross-layer reference lint
        └── 是否需要 audit / replay？
              ├── 是 → 強烈建議 opaque ID
              └── 否 → 看演化頻率
```

## 對 governance runtime 的應用

本專案的 `phase` / `obligation` / `gate` 同時出現在：

- `runtime/runtime.db` 多個 projection table
- `knowledge/runtime/routing-registry.yaml`
- 各層 README 文件
- workflow / intelligence atom 的 inbound reference

這已經滿足「≥ 3 層 + 概念會演化」的雙重條件。建議：

1. **每個 phase / obligation / gate 配一個 stable opaque ID**（例如 `OBL-0042`、`GATE-0017`），DB primary key + 文件 anchor 都用它。
2. **Semantic name 變成 display label**，可以改、可以多語、可以 deprecate。
3. **Cross-layer reference** 一律走 ID，不走 name。重命名只動 label，不動 join。

## 常見誤用

| 誤用 | 正確 |
|------|------|
| 用 `phase_name = "design"` 當 DB foreign key | 用 `phase_id = "PHS-0003"`，`name` 是 display 欄位 |
| 文件之間靠標題互引（`[請見「設計階段」]`） | 引 anchor ID（`[請見 PHS-0003]`） |
| 改名時跨層 grep-and-replace | 改 label，ID 不動，零跨層成本 |
| 對外 API 暴露 opaque ID 當 URL | 對外用 slug，內部用 opaque ID，兩者 mapping |

## 常被誤判為「過度設計」

「ID 看起來醜，可讀性差」是 valid 抱怨，但只在 1–2 層情境成立。在 ≥ 3 層情境下，可讀性損失被跨層 trace 收益完全覆蓋。判斷邊界是 **層數**，不是 **美感**。

## Token Impact

避免「概念改名造成的跨層 silent drift」。一次跨 3 層的概念演化，若靠 semantic name，平均需要 4–8 處同步修改 + 至少一次遺漏導致的事故；用 opaque ID 則為 1 處（label 表）。

---

← [回到 engineering/heuristics/](README.md)
