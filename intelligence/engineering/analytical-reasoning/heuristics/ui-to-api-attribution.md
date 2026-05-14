# UI-to-API Attribution Heuristic（UI 到 API 歸因啟發式）

## 問題

在 APK 分析中，UI 操作需要正確對應到 API 請求才能建立可靠的 feature-to-API mapping。但前景 package 驗證、feature context 確認、操作時間窗對齊等環節若處理不當，會導致 attribution 錯誤。

## 原則

- UI 截圖/XML 只有在前景 package 屬於目標 App 時，才能作為該 App 的操作證據。
- 同一 App 內也可能跑到錯誤頁面，重要 checkpoint 應再驗證目標 feature context。
- 操作時間窗（operation time window）是對齊 UI 操作與 API capture 的核心機制。
- Tab 預載、快取、背景同步可能干擾 timing-based attribution。
- Runtime probe 發現新 UI route 時，應回填到 UI 架構地圖以保留可重用的 route recipe。

## 決策表

| 情境 | 建議策略 | 原理 |
|------|---------|------|
| 前景 package 不是目標 App | 不取 UI 證據，檢查 replay 是否跳離目標 | 截圖/XML 不屬於目標 App |
| UI 操作跳到非目標 feature | 記錄觸發點與外部目的地類型 | 避免 API 對齊到錯誤 feature |
| 抓到 API 但不知道是哪個操作觸發 | 補 screenshot/UI hierarchy + operation id | 建立操作時間窗對齊 capture |
| 截圖 tab 與 API timing 對不上 | 標 trigger confidence low/medium | tab 預載、快取、背景同步可能干擾 |
| Runtime probe 發現新 UI route | 回填到 UI 架構地圖 | 保留可重用的 route recipe |

## 來源 lessons

- [`feedback/history/apk-analysis/http-api/2026-05-05_184700-screen-reachability-operation-recipes.md`](../../../feedback/history/apk-analysis/http-api/2026-05-05_184700-screen-reachability-operation-recipes.md)
- [`feedback/history/apk-analysis/common/2026-05-07_144300-ui-route-backfill-after-runtime-probe.md`](../../../feedback/history/apk-analysis/common/2026-05-07_144300-ui-route-backfill-after-runtime-probe.md)
- [`feedback/history/apk-analysis/common/2026-05-06_111558-foreground-package-validation.md`](../../../feedback/history/apk-analysis/common/2026-05-06_111558-foreground-package-validation.md)
- [`feedback/history/apk-analysis/common/2026-05-06_130300-feature-context-validation.md`](../../../feedback/history/apk-analysis/common/2026-05-06_130300-feature-context-validation.md)

## Token 影響

低。此 atom 在需要確認 UI-to-API attribution 正確性時 lazy-load，約 150-200 tokens。
