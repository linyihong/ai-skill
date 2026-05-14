# 開發指引分析方法

`analysis/development-guidance/` 負責「開發指引的分析方法」。本目錄保存從分析觀察到開發指引的轉換方法論，包括風險翻譯、控制層選擇、指引分類等分析技術。

> **遷移狀態**：此文件為新分層的 reference target，`skills/app-development-guidance/` 已不再作為 active entrypoint。新內容請直接寫入此文件。

## 範圍

本目錄涵蓋以下分析方法：

- **風險翻譯（Risk Translation）**：將分析觀察轉換為開發者視角的風險陳述。
- **擁有者層選擇（Owner Layer Selection）**：根據風險類型選擇最適合的控制層。
- **控制定義（Control Definition）**：定義控制措施的結構與驗證方法。
- **指引分類（Guidance Classification）**：將指引分類到正確的目錄位置。

## 參考 Catalog

領域分類索引已移至 [`metadata/development-guidance/`](../../metadata/development-guidance/)：

| Catalog | 說明 |
|---------|------|
| [`controls-catalog.md`](../../metadata/development-guidance/controls-catalog.md) | 跨平台安全控制類型與核心原則 |
| [`implementation-catalog.md`](../../metadata/development-guidance/implementation-catalog.md) | 實作模式分類與 contract-to-implementation 映射流程 |
| [`platforms-catalog.md`](../../metadata/development-guidance/platforms-catalog.md) | 平台特定指引分類 |
| [`languages-catalog.md`](../../metadata/development-guidance/languages-catalog.md) | 語言特定陷阱分類 |

## 當前來源參考

- `skills/app-development-guidance/WORKFLOW.md` §2-5（Translate To Risk → File The Guidance）
- `skills/app-development-guidance/DOCUMENTATION.md`（Reusable Note Structure）
- `skills/app-development-guidance/controls/`（Cross-platform security controls）
- `skills/app-development-guidance/implementation/`（Buildable implementation patterns）
- `skills/app-development-guidance/platforms/`（Platform-specific guidance）
- `skills/app-development-guidance/languages/`（Language-specific pitfalls）

## 閱讀順序

1. 先讀本 README 了解 scope。
2. 需要分析方法細節時讀 [`risk-translation.md`](risk-translation.md)。
3. 需要分類索引參考時讀 `metadata/development-guidance/` 下的 catalog。
4. 需要執行流程時參考 `workflow/software-delivery/execution-flow.md`。
5. 需要產出規範時參考 `workflow/software-delivery/artifact-gates.md`。

## 遷移說明

- 本目錄為 Phase 17 與 Phase 26 提取產物，內容來自 `skills/app-development-guidance/WORKFLOW.md` §2-4。
- 分類索引（catalogs）已移至 `metadata/development-guidance/`。
- 舊入口 `skills/app-development-guidance/` 已不再作為 active source of truth，新內容請直接寫入此目錄。
- 未來遷移完成條件：所有 analysis methods 完全提取、索引更新、舊入口保留 redirect reference。
