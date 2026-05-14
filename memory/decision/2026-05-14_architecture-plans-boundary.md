# Session-level Decision: architecture/ vs plans/ Boundary Clarification

## Status
accepted

## Context
`architecture/` 目錄同時包含了永久架構文件（`ai-native-knowledge-operating-system.md`）和執行計畫（`context-cost-optimization-plan.md`、`apk-analysis-pilot-migration.md`、`next-stage-upgrade-plan.md`），與 `plans/` 目錄產生重疊。

使用者要求清理重疊，讓 `architecture/` 只保留永久文件。

## Decision
- `architecture/` 只保留永久架構定義：[`ai-native-knowledge-operating-system.md`](../../architecture/ai-native-knowledge-operating-system.md)
- 已完成的執行計畫搬到 `plans/archived/`：[`context-cost-optimization.md`](../../plans/archived/context-cost-optimization.md)、[`apk-analysis-pilot-migration.md`](../../plans/archived/apk-analysis-pilot-migration.md)
- 進行中的路線圖搬到 `plans/active/`：[`next-stage-upgrade-plan.md`](../../plans/archived/next-stage-upgrade-plan.md)
- 更新了 43 處交叉引用，validator 通過

## Consequences
- **正面**：`architecture/` 與 `plans/` 邊界清晰，不再重疊
- **正面**：所有執行計畫統一在 `plans/` 下（active/ 進行中、archived/ 已完成）
- **風險**：既有外部引用可能仍指向舊路徑（已全部更新並驗證）

## Related
- [`decisions/ADR-003-three-layer-architecture.md`](../../decisions/ADR-003-three-layer-architecture.md) — Three-Layer Architecture
- [`plans/README.md`](../../plans/README.md) — Plans 目錄規則與狀態
