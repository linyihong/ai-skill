# Implementation Pattern Catalog

本文件列出可建置的實作模式分類與索引。原承接 [`skills/app-development-guidance/implementation/`](../../skills/app-development-guidance/implementation/) 的內容（已刪除），提取為 tool-neutral 的分析參考。

> **遷移狀態**：`skills/app-development-guidance/implementation/` 已刪除。此文件為 canonical source，新內容請直接寫入此文件。

## 分類

| 類別 | 範圍 | 原始來源 |
|------|------|----------|
| **Backend** | Server/API implementation patterns that mobile and web clients depend on | [`implementation/backend/`](../../skills/app-development-guidance/implementation/backend/)（已刪除） |
| **Mobile** | Android, iOS, Flutter, React Native implementation patterns | [`implementation/mobile/`](../../skills/app-development-guidance/implementation/mobile/)（已刪除） |
| **Embedded** | Firmware, sensor/protocol, hardware context, driver/service/application, and bring-up implementation patterns | [`implementation/embedded/`](../../skills/app-development-guidance/implementation/embedded/)（已刪除） |
| **Tooling** | IDE extensions, CLIs, linters, static analyzers, code generators, and internal automation | [`implementation/tooling/`](../../skills/app-development-guidance/implementation/tooling/)（已刪除） |
| **Examples** | Cross-cutting implementation patterns and snippets in pseudocode | [`implementation/examples/`](../../skills/app-development-guidance/implementation/examples/)（已刪除） |

## 使用流程

當從 contract-first process 開始工作時，使用 implementation docs 將 contracts 轉換為 build slices：

1. Map each Domain Model invariant to provider-side code and unit tests.
2. Map each API, event, command, or public interface contract to provider/consumer fixtures, mocks, or schema checks.
3. Map each Error Handling Contract entry to implementation behavior, logging redaction, and tests.
4. Keep implementation slices linked to the latest contract before teams or agents build in parallel.

## 與其他層的關係

- `analysis/development-guidance/controls-catalog.md` 提供跨平台控制原則，implementation 提供具體實作。
- `analysis/development-guidance/risk-translation.md` 提供從觀察到控制的流程。
- `skills/app-development-guidance/implementation/` 是原始來源，已刪除。內容已由本文件承接。
