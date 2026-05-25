# Software Delivery Validation Workflow

Validation stage 負責 proof acquisition：behavior correctness、business invariant correctness、execution correctness 分別需要對應 evidence。

## Test Ordering（test-first vs test-after）

| 變更類型 | 順序強制等級 | 出處 |
|---------|------------|------|
| Framework / runtime / governance / workflow / validation / scenario / metadata / compiler 改動 | **強制 test-first**（scenarios commit 早於實作 commit）| [`governance/lifecycle/system-upgrade-governance.md`](../../../governance/lifecycle/system-upgrade-governance.md) §3 規則 9 |
| 一般 production code（非 framework / runtime / governance）| 建議 test-first（test-driven design feedback）| [`intelligence/engineering/heuristics/test-driven-heuristic.md`](../../../intelligence/engineering/heuristics/test-driven-heuristic.md) |
| Doc-only trial / bug fix / typo / spike | 豁免（須明寫理由）| 規則 9 豁免清單 |

詳細原則：

- 通用 TDD（測試難寫 → 設計回饋）→ [`intelligence/engineering/heuristics/test-driven-heuristic.md`](../../../intelligence/engineering/heuristics/test-driven-heuristic.md)
- Framework 升級的順序強制（scenarios 必須先於實作）→ [`intelligence/engineering/development/test-first-framework-upgrade.md`](../../../intelligence/engineering/development/test-first-framework-upgrade.md)

兩者互補：通用 TDD 處理「測試難不難寫」（設計問題），test-first-framework-upgrade 處理「測試何時寫」（順序問題）。

## 進入本 stage 前

[`execution-flow.md`](../execution-flow.md) §4 測試策略定義 已確認測試策略；§7 驗證 執行 proof acquisition；§Test-First Ordering callout 確認 framework 升級已遵循順序。
