# Runtime Onboarding

## 放什麼

新專案或新任務的初始設定指引、開場提示詞模板、完成門檻定義。這些文件告訴使用者或 AI 如何開始使用某個 skill，以及何時算完成。

## 目前文件

| 文件 | 描述 |
|------|------|
| [`apk-analysis-setup.md`](apk-analysis-setup.md) | APK 分析專案的初始設定流程與提示詞模板 |
| [`apk-analysis-completion.md`](apk-analysis-completion.md) | APK 分析專案的完成門檻定義 |

## 與其他層的關係

- `runtime/` 根目錄提供 runtime 整體架構，本目錄提供具體的 onboarding 指引。
- `workflow/` 提供執行流程，本目錄提供如何啟動流程的設定指引。
- `skills/` 是原始來源，已不再作為 active entrypoint（舊 `skills/` 結構已於 2026-05-13 標記為 deprecated）。
