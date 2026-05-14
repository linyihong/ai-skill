# Traceability Heuristic（文件追溯性經驗法則）

**Status**: `candidate-intelligence`
**Source**: [`analysis/repo/traceability-gate.md`](../../analysis/repo/traceability-gate.md), [`analysis/repo/contract-governance.md`](../../analysis/repo/contract-governance.md)

## 原則

**If a requirement cannot be traced to a test, it is not verified. If a test cannot be traced to a requirement, it is not justified.**

如果一個需求無法追溯到測試，它沒有被驗證。如果一個測試無法追溯到需求，它沒有被正當化。

## 為什麼

1. **追溯性防止「已實作但未驗證」的缺口** — 沒有追溯性的需求可能在程式碼中「存在」但從未被測試覆蓋。
2. **追溯性防止「已測試但無關」的浪費** — 沒有追溯性的測試可能是測試了不存在的需求或過時的行為。
3. **追溯性是變更影響分析的基礎** — 當需求改變時，追溯性能快速找出受影響的實作和測試。
4. **Stable IDs 是追溯性的基礎設施** — 沒有 stable IDs，追溯連結會在重構中斷裂。

## 何時適用

- 為已實作完成的 repository 建立文件追溯性。
- 確認每個需求都有對應的實作與測試。
- 從程式碼反向追溯回原始需求或規則。
- 評估變更的影響範圍時。

## 何時不適用

- 專案處於 prototyping 階段，需求仍在快速變動。
- 專案沒有 stable IDs 且短期內不會引入。
- 只需要高層次的架構理解，不需要逐項追溯。

## 決策流程

```text
需要建立追溯性？
  ├── 建立追溯連結：
  │     ├── Product/rule ID → BDD（哪個 behavior 證明該需求）
  │     ├── BDD → code refs（behavior 在哪裡實作）
  │     ├── BDD → test refs（behavior 如何被驗證）
  │     ├── Contract operation → fixture（provider/consumer 相容性）
  │     └── Generated client → source contract（防止 drift）
  ├── 使用 Stable IDs：
  │     ├── Feature IDs、Rule IDs、Operation IDs
  │     ├── Route names、Command names、Diagnostic codes
  │     ├── Event names、Scenario tags
  │     └── 避免使用行號、檔名（重構時會變）
  └── 處理未實作行為：
        ├── TBD → 待決定
        ├── noop → 無操作（intentionally empty）
        ├── not enforceable by tool → 工具無法強制執行
        ├── manual-only → 僅手動驗證
        └── out of scope → 明確排除在範圍外
```

## 常見誤用

| 誤用 | 正確 |
|------|------|
| 「測試全部通過，所以需求都實作了」 | 測試可能只覆蓋了部分場景，或測試本身有 gap |
| 「這個功能有測試，所以沒問題」 | 確認測試是否追溯到具體的需求 ID，還是「順便測的」 |
| 「用檔名作為追溯 ID」 | 檔名在重構中會改變；使用 stable IDs |
| 「未實作的行為不需要記錄」 | 未實作的行為必須標記原因，否則未來可能被誤認為 bug |

## Token Impact

避免因缺乏追溯性導致的回歸 bug。一個沒有追溯性的變更可能意外破壞未測試的行為，而發現時已經是 production incident。

---

← [回到 intelligence/engineering/analytical-reasoning/](README.md)
