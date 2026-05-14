# Document Priority Hierarchy Heuristic（文件優先順序階層經驗法則）

**Status**: `candidate-intelligence`
**Source**: 從 [`intelligence/engineering/development/contract-governance-heuristic.md`](../app-development-guidance/contract-governance-heuristic.md) 提取的跨領域通用部分

## 原則

**If documents conflict, update the owning document instead of silently fixing only one source.**

如果文件之間存在衝突，更新管轄該行為的文件，而不是默默地只修正其中一個來源。

## 為什麼

1. **任何單一來源都只顯示當前狀態，不顯示意圖** — 程式碼只顯示「當前的實作」，文件只顯示「最後一次更新時的認知」。只有管轄文件才能解決根本原因。
2. **沒有優先順序時，每次衝突都需要重新協商** — 浪費時間且結果不一致。
3. **取消或放棄的決定需要明確記錄** — 未記錄的取消會在未來被重新引入，因為沒有證據顯示它曾被考慮過。
4. **優先順序階層是 scalability 的前提** — 當多個 agent 或團隊維護同一個知識庫時，沒有治理規則會導致混亂。

## 何時適用

- 知識庫或 repository 有多份文件且需要定義優先順序。
- 文件之間存在不一致。
- 多個 agent 或團隊共同維護同一個知識庫。
- 需要建立文件變更的標準流程。

## 何時不適用

- 單一開發者維護的個人專案（但建議仍遵循，養成習慣）。
- Prototyping 階段，文件被明確視為非必要。
- 衝突已被確認是實作錯誤（此時更新 tests/validation，保持上層文件穩定）。

## 通用決策流程

```text
文件不一致？
  ├── 分類衝突類型：
  │     ├── 上層意圖改變 → 更新管轄該意圖的文件
  │     ├── 規格遺漏/過時 → 從 evidence 回填規格
  │     ├── 契約過時 → 更新契約 + consumers + validation
  │     ├── 實作錯誤 → 新增 regression tests，修正實作
  │     └── Test/fixture 過時 → 更新 tests/fixtures
  └── 按文件優先順序執行：
        1. Governance / framework contract（不變規則）
        2. Intent / plan / brief（意圖與範圍）
        3. Behavior specification（可觀察行為）
        4. Interface / API / architecture contracts（各層契約）
        5. Implementation（實作）
        6. Tests / fixtures / examples（驗證）
```

## 通用優先順序階層

| 優先級 | 層級 | 說明 | 範例 |
|--------|------|------|------|
| 1 | Governance / framework | 不變規則、required update rules | Repository invariants、lifecycle policy |
| 2 | Intent / plan | 意圖、範圍、non-goals | Product brief、project plan、task intent |
| 3 | Behavior specification | 可觀察行為、acceptance criteria | BDD scenarios、user stories、signal tables |
| 4 | Interface / architecture | 各層契約與邊界 | API contracts、architecture docs、routing registry |
| 5 | Implementation | 實作細節 | Source code、workflow steps、analysis methods |
| 6 | Validation | 驗證與範例 | Tests、fixtures、checklists、validation gates |

## 常見誤用

| 誤用 | 正確 |
|------|------|
| 「程式碼是真相的唯一來源」 | 程式碼只顯示當前狀態，不顯示意圖；上層文件才是意圖的載體 |
| 「先修正實作，文件之後再補」 | 如果文件是管轄文件，應該先更新文件再修正實作 |
| 「取消的功能不需要記錄」 | 取消的功能必須明確記錄，否則未來可能被重新引入 |
| 「所有文件有同等權重」 | 不同文件有不同管轄範圍和優先順序；Governance 層永遠最高 |

## Token Impact

避免因文件衝突導致的重複除錯循環。一個未解決的文件衝突可能導致 2-3 次錯誤的修改，每次耗費 15-30 分鐘。建立明確的優先順序階層後，衝突解決時間可降至 5 分鐘內。

---

← [回到 engineering/heuristics/](README.md)
