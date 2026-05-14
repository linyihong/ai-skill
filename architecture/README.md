# Architecture

`architecture/` 存放永久性架構文件，定義系統的長期設計原則與演化方向。

## 用途

- 定義 OS 層級架構與相容性策略
- 記錄跨 session 的設計原則（不變的真理來源）
- 不包含執行計畫（執行計畫請見 [`plans/`](../plans/README.md)）

## 目前文件

| 文件 | 說明 |
|------|------|
| [`ai-native-knowledge-operating-system.md`](ai-native-knowledge-operating-system.md) | AI-native Knowledge Operating System 架構定義、層級說明、Reference-first 載入策略、相容層與 Deprecation 流程 |

## 誰會參考這裡（Inbound References）

變更本層內容時，需要一併檢查以下依賴方：

| 來源 | 關係 |
|------|------|
| [`route.architecture.permanent-docs`](../knowledge/runtime/routing-registry.yaml) | Routing registry record，agent 依此找到 architecture/ |
| [`README.md`](../README.md) | 根目錄 OS Layout 表格列出 architecture/ 層 |
| [`decisions/`](../decisions/README.md) | ADR 可能引用 architecture 架構原則 |
| [`plans/`](../plans/README.md) | 執行計畫可能引用 architecture 設計原則 |

## 與既有層的關係

- [`plans/`](../plans/README.md)：執行計畫存放處，完成後歸檔至 `plans/archived/`
- [`decisions/`](../decisions/README.md)：架構決策記錄（ADR），記錄為什麼做出某個架構選擇
- [`README.md`](../README.md)：根目錄 README 已列出 architecture 作為 OS Layout 的一層
