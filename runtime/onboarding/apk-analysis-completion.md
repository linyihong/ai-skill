# APK Analysis Completion Definition

本文件定義新 APK 分析專案的完成門檻。承接 [`workflow/apk-analysis/execution-flow.md`](../../workflow/apk-analysis/execution-flow.md) 與 [`workflow/apk-analysis/artifact-gates.md`](../../workflow/apk-analysis/artifact-gates.md) 的內容，提取為 tool-neutral 的 runtime completion gates。

> **遷移狀態**：此文件為新分層的 reference target，舊 `skills/apk-analysis/` 已不再作為 active entrypoint。新內容請直接寫入此文件。

## 初步分析完成定義

一次新 APK 初步分析完成時，至少應該得到：

| # | 項目 | 說明 |
|---|------|------|
| 1 | **流量路徑判斷結果** | 核心流量走哪一層（localhost / pcap / proxy / Java / native / Flutter） |
| 2 | **代理 / MITM 是否可用的證據** | 是否可成功代理流量，或存在 pinning / 繞過障礙 |
| 3 | **Java / native / Flutter stack 判斷** | APK 使用的技術棧判斷結果 |
| 4 | **初步 request metadata 或下一步 hook 計畫** | 至少知道 request shape 或下一步要 hook 什麼 |
| 5 | **response wrapper 或解密定位計畫** | 知道 response 如何包裝、是否需要解密 |
| 6 | **去敏規則與文件位置** | 已定義去敏規則並知道文件存放位置 |
| 7 | **Domain/runtime baseline 或 blocker 清單** | 若目標包含 SDK/client/replay/live integration：domain/runtime baseline 已回答最小可跑因素，或缺口已列 blocker / scoped out；僅 skeleton 時不得宣稱可開始 live-facing 開發 |
| 8 | **新 lesson 回饋** | 是否有新 lesson 回饋到 skill（若有發現新技巧） |

## 第一輪分析順序

1. **盤點 APK**：package name、version、architecture、permissions、native libraries、Flutter / React Native / Java/Kotlin / native 初步判斷。
2. **建立分析環境紀錄**：device / emulator、adb、Frida client/server、proxy tool、root status、allowed actions。
3. **流量路徑判斷**：localhost / loopback、whole-device pcap、system proxy / MITM、Java HTTP stack、native connect trace、Flutter / Dart AOT if applicable。
4. **找高語意 hook**：request options、response wrapper、response decoder / decryptor、token/session provider。
5. **分類路由**：先用共通流程判斷 runtime / traffic family。證據指向特定技術才讀對應 techniques 文件。不要在分類未明時一次讀完所有 technique folders。
6. **文件化**：失敗路徑也要寫、成功證據要可重現、API 結論和方法論要分開、去敏後才保存樣本。若下一步是 SDK/client/app tool/live integration 開發，先補 project-level domain/runtime baseline 的最小可跑因素；若需要 device/install/account/session/vendor/server-issued material，補 authorized identity material self-generation audit。
7. **回饋 skill**：新技巧寫入 `feedback/history/apk-analysis/<category>/` 或 `feedback/history/apk-analysis/common/`。已驗證技巧再同步進 `workflow/apk-analysis/execution-flow.md`、`workflow/apk-analysis/artifact-gates.md`、`analysis/apk/workflows/` 或 `intelligence/engineering/analytical-reasoning/`。

## 與其他層的關係

- `workflow/apk-analysis/execution-flow.md` 提供詳細的分析執行步驟，本文件定義何時算完成。
- `workflow/apk-analysis/artifact-gates.md` 提供產出規範與品質門檻。
- `skills/apk-analysis/RUNBOOK.md` 是原始來源，已不再作為 active entrypoint（舊 `skills/` 結構已於 2026-05-13 標記為 deprecated）。
