# Contract Governance Heuristic（契約治理經驗法則）

> **Cross-Domain Promotion**: 文件優先順序階層模式已提取到 [`intelligence/engineering/heuristics/document-priority-hierarchy.md`](../../heuristics/document-priority-hierarchy.md)

**Status**: `candidate-intelligence`
**Source**: [`analysis/repo/contract-governance.md`](../../analysis/repo/contract-governance.md), [`skills/app-development-guidance/process/README.md`](../../skills/app-development-guidance/process/README.md)（已刪除）

## 原則

**If documents conflict, update the owning document instead of silently fixing only code.**

如果文件之間存在衝突，更新管轄該行為的文件，而不是默默地只修正程式碼。

## 為什麼

1. **程式碼是實作，不是規格** — 當程式碼與文件不一致時，程式碼可能是 bug，也可能是文件過時。只有更新管轄文件才能解決根本原因。
2. **文件優先順序防止無限迴圈** — 沒有優先順序時，每次衝突都需要重新協商，浪費時間。
3. **取消與延後的行為需要明確記錄** — 未記錄的取消會在未來被 agent 重新引入，因為沒有證據顯示它曾被考慮過。
4. **契約治理是 scalability 的前提** — 當多個 agent 或團隊維護同一個 repository 時，沒有治理規則會導致文件混亂。

## 何時適用

- Repository 有多份開發文件且需要定義優先順序。
- 文件與程式碼之間存在不一致。
- 多個 agent 或團隊共同維護同一個 repository。
- 需要建立文件變更的標準流程。

## 何時不適用

- 單一開發者維護的個人專案（但建議仍遵循，養成習慣）。
- Prototyping 階段，文件被明確視為非必要。
- 文件衝突已被確認是程式碼 bug（此時更新 tests，保持 docs 穩定）。

## 決策流程

```text
文件與程式碼不一致？
  ├── 分類衝突類型：
  │     ├── Product intent changed → 更新 product brief
  │     ├── BDD missing/stale → 從 evidence 回填 BDD
  │     ├── Contract stale → 更新 contract + consumers + tests
  │     ├── Implementation bug → 新增 regression tests，修正程式碼
  │     └── Test/fixture stale → 更新 tests/fixtures
  └── 按文件優先順序執行：
        1. Governance / framework contract
        2. Product plan / brief
        3. BDD behavior
        4. Domain / architecture / API contract
        5. Implementation
        6. Tests / fixtures
```

## 預設文件優先順序

| 優先級 | 文件類型 | 說明 |
|--------|----------|------|
| 1 | Governance / framework contract | Repository-wide invariants、required update rules |
| 2 | Product plan / accepted brief | Product intent、scope、non-goals |
| 3 | BDD behavior | Observable user/system behavior、acceptance criteria |
| 4 | Domain, architecture, API/interface contracts | 各層契約 |
| 5 | Implementation and generated clients | 實作與產生的客戶端 |
| 6 | Tests, fixtures, and examples | 測試與範例 |

## 常見誤用

| 誤用 | 正確 |
|------|------|
| 「程式碼是真相的唯一來源」 | 程式碼只顯示當前狀態，不顯示意圖；文件才是意圖的載體 |
| 「先修正程式碼，文件之後再補」 | 如果文件是管轄文件，應該先更新文件再修正程式碼 |
| 「取消的功能不需要記錄」 | 取消的功能必須明確記錄，否則未來可能被重新引入 |
| 「所有文件有同等權重」 | 不同文件有不同管轄範圍和優先順序 |

## Token Impact

避免因文件衝突導致的重複除錯循環。一個未解決的文件衝突可能導致 2-3 次錯誤的程式碼修改，每次耗費 15-30 分鐘。

---

← [回到 engineering/app-development-guidance/](README.md)
