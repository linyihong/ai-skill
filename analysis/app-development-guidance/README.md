# App Development Guidance Analysis Methods

`analysis/app-development-guidance/` 負責「開發指引的分析方法」。本目錄保存從分析觀察到開發指引的轉換方法論，包括風險翻譯、控制層選擇、指引分類等分析技術，以及跨平台安全控制、實作模式、平台指引、語言陷阱的 catalog 參考。

> **遷移狀態**：此文件為新分層的 reference target，`skills/app-development-guidance/` 已不再作為 active entrypoint。新內容請直接寫入此文件。

## Scope

本目錄涵蓋以下分析方法與參考 catalog：

### 分析方法

- **Risk Translation**：將分析觀察轉換為開發者視角的風險陳述。
- **Owner Layer Selection**：根據風險類型選擇最適合的控制層。
- **Control Definition**：定義控制措施的結構與驗證方法。
- **Guidance Classification**：將指引分類到正確的目錄位置。

### 參考 Catalog

| Catalog | 說明 | 原始來源 |
|---------|------|----------|
| [`controls-catalog.md`](controls-catalog.md) | 跨平台安全控制類型與核心原則（API Transport、Auth & Session、Local Storage、Logging & Telemetry、Anti-Tamper、Release Build） | `skills/app-development-guidance/controls/` |
| [`implementation-catalog.md`](implementation-catalog.md) | 實作模式分類（Backend、Mobile、Embedded、Tooling、Examples）與 contract-to-implementation 映射流程 | `skills/app-development-guidance/implementation/` |
| [`platforms-catalog.md`](platforms-catalog.md) | 平台特定指引分類（Mobile、Web、Backend、Embedded） | `skills/app-development-guidance/platforms/` |
| [`languages-catalog.md`](languages-catalog.md) | 語言特定陷阱分類（Dart、Kotlin/Java、Swift、TypeScript） | `skills/app-development-guidance/languages/` |

## Current Source References

- `skills/app-development-guidance/WORKFLOW.md` §2-5（Translate To Risk → File The Guidance）
- `skills/app-development-guidance/DOCUMENTATION.md`（Reusable Note Structure）
- `skills/app-development-guidance/controls/`（Cross-platform security controls）
- `skills/app-development-guidance/implementation/`（Buildable implementation patterns）
- `skills/app-development-guidance/platforms/`（Platform-specific guidance）
- `skills/app-development-guidance/languages/`（Language-specific pitfalls）

## Read Order

1. 先讀本 README 了解 scope 與可用 catalog。
2. 需要安全控制參考時讀 `controls-catalog.md`。
3. 需要實作模式參考時讀 `implementation-catalog.md`。
4. 需要平台或語言特定指引時讀 `platforms-catalog.md` 或 `languages-catalog.md`。
5. 需要執行流程時參考 `workflow/app-development-guidance/execution-flow.md`。
6. 需要產出規範時參考 `workflow/app-development-guidance/artifact-gates.md`。

## Migration Notes

- 本目錄為 Phase 17 與 Phase 26 提取產物，內容來自 `skills/app-development-guidance/WORKFLOW.md` §2-4 及 `skills/app-development-guidance/` 的 controls/、implementation/、platforms/、languages/ 子目錄。
- 舊入口 `skills/app-development-guidance/` 已不再作為 active source of truth，新內容請直接寫入此目錄。
- 未來遷移完成條件：所有 analysis methods 完全提取、索引更新、舊入口保留 redirect reference。
