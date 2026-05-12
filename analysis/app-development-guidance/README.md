# App Development Guidance Analysis Methods

`analysis/app-development-guidance/` 負責「開發指引的分析方法」。本目錄保存從分析觀察到開發指引的轉換方法論，包括風險翻譯、控制層選擇、指引分類等分析技術。

> **相容性規則**：`skills/app-development-guidance/WORKFLOW.md` 仍為 active skill entrypoint。本文件為 reference target，兩者應保持同步。

## Scope

本目錄涵蓋以下分析方法：

- **Risk Translation**：將分析觀察轉換為開發者視角的風險陳述。
- **Owner Layer Selection**：根據風險類型選擇最適合的控制層。
- **Control Definition**：定義控制措施的結構與驗證方法。
- **Guidance Classification**：將指引分類到正確的目錄位置。

## Current Source References

- `skills/app-development-guidance/WORKFLOW.md` §2-5（Translate To Risk → File The Guidance）
- `skills/app-development-guidance/DOCUMENTATION.md`（Reusable Note Structure）

## Read Order

1. 先讀本 README 了解 scope。
2. 需要執行流程時參考 `workflow/app-development-guidance/execution-flow.md`。
3. 需要產出規範時參考 `workflow/app-development-guidance/artifact-gates.md`。

## Migration Notes

- 本目錄為 Phase 17 提取產物，內容來自 `skills/app-development-guidance/WORKFLOW.md` §2-4。
- 舊入口 `skills/app-development-guidance/WORKFLOW.md` 仍為 active source of truth。
- 未來遷移完成條件：所有 analysis methods 完全提取、索引更新、舊入口保留 redirect reference。
