# UI Operation Stability Heuristic（UI 操作穩定性啟發式）

## 問題

在 APK 分析中，UI 操作是觸發 API 請求的主要手段，但不同操作策略（scroll、click、replay script）在不同情境下的穩定性差異很大。何時該用 bounded scroll、何時該固化為 operation script、何時該用 API-first replay？

## 原則

- UI 操作策略應根據 **API attribution 需求** 與 **UI 行為穩定性** 選擇，而非一律用同一種方式。
- Bounded scroll（top/mid/bottom）比無限滑動更可控，可避免混入背景請求。
- 可重放的 operation script 可減少每輪重新推理 UI 操作的成本。
- 每個 script 只做一個 flow，避免多 action 混雜 attribution。
- 若已有 API boundary 可保留 session/signing，API-first replay 比長 UI scroll 更快取得證據。

## 決策表

| 情境 | 建議策略 | 原理 |
|------|---------|------|
| 需要區分 initial load vs pagination vs tap-triggered API | Bounded scroll top/mid/bottom + 操作時間窗 | 避免無限制滑動混入背景請求 |
| 同一 feature 需要反覆 capture | 固化為可重放 operation script | 減少每輪重新推理 UI 操作 |
| 操作腳本結果不穩 | 每個 script 只做一個 flow + 輸出 timestamp | 避免多 action 混雜 attribution |
| 已有 API boundary 可保留 session/signing | API-first replay 取代長 UI scroll | 更快取得證據，UI 只留 attribution |
| UI behavior 本身就是問題 | UI capture + bounded gesture + package/feature context guard | 無法用 API replay 取代 |

## 來源 lessons

- [`feedback/history/apk-analysis/http-api/2026-05-01_171650-ui-architecture-map-from-screenshots.md`](../../../feedback/history/apk-analysis/http-api/2026-05-01_171650-ui-architecture-map-from-screenshots.md)
- [`feedback/history/apk-analysis/http-api/2026-05-05_183400-ui-automation-operation-scripts-for-api-capture.md`](../../../feedback/history/apk-analysis/http-api/2026-05-05_183400-ui-automation-operation-scripts-for-api-capture.md)
- [`feedback/history/apk-analysis/http-api/2026-05-05_183700-scrollable-clickable-screen-mapping.md`](../../../feedback/history/apk-analysis/http-api/2026-05-05_183700-scrollable-clickable-screen-mapping.md)
- [`feedback/history/apk-analysis/common/2026-05-07_124100-ui-fast-path-bounded-scroll.md`](../../../feedback/history/apk-analysis/common/2026-05-07_124100-ui-fast-path-bounded-scroll.md)

## Token 影響

低。此 atom 在需要決定 UI 操作策略時 lazy-load，約 150-200 tokens。
