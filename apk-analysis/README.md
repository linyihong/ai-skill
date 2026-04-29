# APK 分析 Skill 與人類指南

這個資料夾是可持續強化的 APK 分析知識包。它同時服務兩種讀者：

- AI / Claude / Cursor agent：讀 `SKILL.md`，知道遇到 APK 分析、動態抓包、解密定位、Frida hook、Flutter/Dart AOT 分析時要怎麼做。
- 人類分析者：讀本 README、`WORKFLOW.md`、`TOOLS.md`、`DOCUMENTATION.md`，快速理解分析順序、必要工具、記錄格式與回饋方式。

## 目標

把一次次 APK 分析中真正可重用的「方法」沉澱下來，而不是保存某個目標 App 的私有結論。

這裡應該收：

- 如何判斷流量走哪一層。
- 什麼時候用 pcap、MITM、Frida、靜態分析或 Dart AOT。
- 如何找高語意 hook 點。
- 如何把動態樣本變成離線 fixture、schema 與測試。
- 如何撰寫可重現、可去敏、可回顧的分析文件。
- 新想法如何回饋到 skill，讓後續 agent 更強。

這裡不應該收：

- 特定 App 的完整 host、endpoint、token、secret、device id。
- 未去敏的 request / response。
- 使用者個資或帳號資料。
- 只對單一專案有效、沒有泛化價值的結論。

## 資料夾內容

| 文件 | 用途 |
| --- | --- |
| `SKILL.md` | 給 AI agent 的技能入口，包含觸發條件、工作原則與回饋規則。 |
| `WORKFLOW.md` | 人類與 agent 都可讀的分析決策流程。 |
| `TOOLS.md` | 常用工具、前置條件、適用情境與失敗判讀。 |
| `DOCUMENTATION.md` | 分析結果如何寫成可重現文件。 |
| `FEEDBACK.md` | 新技巧、新失敗模式、新驗證規則的回饋模板與待提升清單。 |
| `RUNBOOK.md` | 新 APK 專案第一天如何套用 skill，以及如何把新經驗回饋回本包。 |

## 使用方式

開始分析前：

1. 確認授權範圍與目標 APK 版本。
2. 建立本次分析筆記位置。
3. 新專案第一天先讀 `RUNBOOK.md`。
4. 讀 `WORKFLOW.md`，先做流量路徑判斷，不要直接假設是 pinning 或加密。
5. 依 `TOOLS.md` 選最低干擾工具。
6. 依 `DOCUMENTATION.md` 記錄證據鏈。

分析完成後：

1. 把 target-specific API / schema / endpoint 結論寫回專案對應 API 文件。
2. 把可重用方法寫入 `FEEDBACK.md`。
3. 如果方法已驗證，整理進 `WORKFLOW.md` 或 `TOOLS.md`。
4. 若專案有 SDK 或 client，將解碼規則補成 fixture / contract test。

## 核心原則

- 先證明流量在哪一層，再選工具。
- 先判斷是否走代理，再談 CA / pinning。
- 優先 hook 高語意物件，例如 request options、response interceptor、decrypt function。
- 低層 socket / TLS hook 只作補證據或最後手段。
- 動態 hook 只是過渡；最終要沉澱成離線解碼器、fixture、schema 與測試。
- 文件要分離「方法」與「目標 App 的結論」。
- 新發現要回饋到 skill，但必須去敏、泛化、可驗證。

## 最小產出

一次合格的 APK 分析至少要留下：

- 分析環境：OS、device/emulator、APK version、arch、root/Frida 狀態。
- 流量路徑判斷：localhost、whole-device pcap、proxy/MITM、Java hook、native/Dart 路徑。
- 證據：pcap、hook log、MITM export、反編譯搜尋結果或 sanitized excerpt。
- 結論：哪些路徑有效、哪些被排除、下一步如何驗證。
- 去敏規則：哪些值被遮蔽，哪些文件不能提交。
- 可重用 lesson：寫入 `FEEDBACK.md`。
