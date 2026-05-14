# Platform Guidance Catalog

本文件列出各平台/應用類型的開發指引分類與索引。原承接 [`skills/app-development-guidance/platforms/`](../../skills/app-development-guidance/platforms/) 的內容（已刪除），提取為 tool-neutral 的分析參考。

> **遷移狀態**：`skills/app-development-guidance/platforms/` 已刪除。此文件為 canonical source，新內容請直接寫入此文件。

## 平台分類

| 平台 | 範圍 | 原始來源 |
|------|------|----------|
| **Mobile** | Android, iOS, Flutter, React Native, and mobile release concerns | [`platforms/mobile/`](../../skills/app-development-guidance/platforms/mobile/)（已刪除） |
| **Web** | Browser/frontend app hardening | [`platforms/web/`](../../skills/app-development-guidance/platforms/web/)（已刪除） |
| **Backend** | API/server controls that apps depend on | [`platforms/backend/`](../../skills/app-development-guidance/platforms/backend/)（已刪除） |
| **Embedded** | Firmware, sensors, hardware context, protocols, board bring-up, and hardware-in-loop validation | [`platforms/embedded/`](../../skills/app-development-guidance/platforms/embedded/)（已刪除） |

## 使用原則

1. **Prefer `controls/` for cross-platform principles** — 跨平台原則放在 controls，然後連結到這裡看平台差異。
2. **Implementation patterns live in `implementation/`** — 可建置步驟在 implementation，必須與 platform docs 保持連結。
3. **Platform docs 記錄該平台特有的風險與緩解** — 例如 Flutter 的 platform channel 邊界、Android 的 storage 選擇。

## 與其他層的關係

- `analysis/development-guidance/controls-catalog.md` 提供跨平台控制原則。
- `analysis/development-guidance/implementation-catalog.md` 提供可建置實作模式。
- `skills/app-development-guidance/platforms/` 是原始來源，已刪除。內容已由本文件承接。
