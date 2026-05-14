# APK 分析 Skill 與人類指南

本包位於 **`skills/apk-analysis/`**（與其他 skill 並列於 [`skills/`](../README.md)）。**全庫共用政策與 feedback 寫作底線**見 [`../../shared-rules/README.md`](../../shared-rules/README.md)。

這個資料夾是可持續強化的 APK 分析知識包。它同時服務兩種讀者：

- AI / agent 工具：讀 `SKILL.md`，知道遇到 APK 分析、動態抓包、解密定位、Frida hook、Flutter/Dart AOT 分析時要怎麼做。
- 人類分析者：讀本 README、`WORKFLOW.md`、`TOOLS.md`、`DOCUMENTATION.md`，快速理解分析順序、必要工具、記錄格式與回饋方式。

## 下一階段 migration note

`skills/apk-analysis/` 目前仍是 active skill entrypoint。下一階段分層的 pilot map 見 [`../../plans/archived/apk-analysis-pilot-migration.md`](../../plans/archived/apk-analysis-pilot-migration.md)；候選 reference-first 目的地為 [`../../analysis/apk/`](../../analysis/apk/)、[`../../workflow/apk-analysis/`](../../workflow/apk-analysis/) 與 [`../../intelligence/engineering/analytical-reasoning/`](../../intelligence/engineering/analytical-reasoning/)。Pilot 期間不搬移大量內容，既有入口與連結維持可用。

## 目標

把一次次 APK 分析中真正可重用的「方法」沉澱下來，而不是保存某個目標 App 的私有結論。

這裡應該收：

- 如何判斷流量走哪一層。
- 什麼時候用 pcap、MITM、Frida、靜態分析或 Dart AOT。
- 如何找高語意 hook 點。
- 如何把動態樣本變成離線 fixture、schema 與測試。
- 如何把 API 整理成可維護的 API Catalog：總入口、分組索引、逐支 API 詳細文件、coverage/gap、UI 對照、SDK/client 欄位用途與驗證。
- **如何把「能連線拿到真實資料」的前提寫進專案：** Domain／執行環境基線（host family、session／opaque 參數、authorized identity material self-generation audit、分頁真實語意、簽章／gateway 前置，皆去敏）與 API Catalog **互鏈**，見 `DOCUMENTATION.md` § Domain／執行環境基線。
- 如何判斷 APK findings 是否已足夠開始 SDK/client/app tool 開發：若會連真實服務或跑 live integration，必須先通過 development readiness gate；若需要裝置、帳號、session、vendor attestation 或 server-issued material，還必須分析到每個 key group 能否自生成、怎麼生成或由誰提供；缺最小可跑因素時先補分析文件或列 blocker，不得從 API shape 直接開發 live-facing code。
- 如何撰寫可重現、可去敏、可回顧的分析文件。
- 如何把 APK 觀察整理成可交給 `app-development-guidance` 的功能重建交接規格。
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
| `techniques/` | 依 runtime / traffic family 拆分的技術分類；例如 Flutter/Dart AOT、HTTP API、local proxy、media/HLS。 |
| `FEEDBACK.md` | 入口：指向全庫共用的 [`feedback-lessons.md`](../../shared-rules/feedback-lessons.md)。 |
| `feedback_history/` | 每一條獨立 lesson 的 Markdown；必要時見 `README.md` 索引表。 |
| `RUNBOOK.md` | 新 APK 專案第一天如何套用 skill，以及如何把新經驗回饋回本包。 |

## 技術分類

`apk-analysis` 的主文件只保留共通流程與路由規則。特定 app/runtime/API family 的技巧放到 [`techniques/`](techniques/)：

| 分類 | 何時讀 |
| --- | --- |
| [`techniques/flutter-dart-aot/`](techniques/flutter-dart-aot/) | 證據指向 Flutter、Dart AOT、`libapp.so`、Dio/interceptor 或 Dart decode。 |
| [`techniques/http-api/`](techniques/http-api/) | 已觀測到 HTTP API，需要文件化 headers、request、response、fixture 或 replay。 |
| [`techniques/local-proxy/`](techniques/local-proxy/) | 證據指向 loopback、Netty/local proxy、embedded TUN/VPN 或本機轉發。 |
| [`techniques/media-hls/`](techniques/media-hls/) | 分析 HLS、key、segments、媒體下載、容器或圖片/音訊/影片格式。 |

原則：不要一開始讀完所有分類。先用 `WORKFLOW.md` 判斷流量與 runtime，再只讀相關分類；若是 A 類分析，就避免把 B 類文件放進上下文。

## 使用方式

開始分析前：

1. 確認授權範圍與目標 APK 版本。
2. 建立本次分析筆記位置。
3. 新專案第一天先讀 `RUNBOOK.md`。
4. 讀 `WORKFLOW.md`，先做流量路徑判斷，不要直接假設是 pinning 或加密。
5. 證據明確後，只讀 [`techniques/`](techniques/) 中對應分類。
6. 依 `TOOLS.md` 選最低干擾工具。
7. 依 `DOCUMENTATION.md` 記錄證據鏈。

分析完成後：

1. 把 target-specific API / schema / endpoint 結論寫回專案對應 API 文件。
2. 把可重用方法寫成 **`feedback_history/<category>/YYYY-MM-DD_HHMMSS-<slug>.md`**；跨分類用 `feedback_history/common/`（規則見 [`shared-rules/feedback-lessons.md`](../../shared-rules/feedback-lessons.md)）。
3. 如果方法已驗證，整理進 `WORKFLOW.md` 或 `TOOLS.md`。
4. 若專案有 SDK 或 client，將解碼規則補成 fixture / contract test。
5. 若得到的是「未來開發自家 App 時可用的設計、實作或防護 guidance」而非分析方法，寫入 [`app-development-guidance`](../app-development-guidance/)；本包只保留分析證據與方法。
6. 若目標是重新做出同等功能，專案分析文件必須包含 Feature Reconstruction Handoff，讓 [`app-development-guidance`](../app-development-guidance/) 可以接手產生 BDD、Domain Model Contract、API / Interface Contract、Error Handling Contract、implementation slices 與 tests。
7. 若分析文件要用來做 app 相關工具、SDK、client、mock API、fixture-driven implementation、contract test 或重建功能，agent 必須自動讀取並套用 [`app-development-guidance/SKILL.md`](../app-development-guidance/SKILL.md)，不要只在 `apk-analysis` 內寫實作計畫。

## 核心原則

- 先證明流量在哪一層，再選工具。
- 先判斷是否走代理，再談 CA / pinning。
- 優先 hook 高語意物件，例如 request options、response interceptor、decrypt function。
- 低層 socket / TLS hook 只作補證據或最後手段。
- 動態 hook 只是過渡；最終要沉澱成離線解碼器、fixture、schema 與測試。
- 文件要分離「方法」與「目標 App 的結論」。
- 分析文件要保留足夠的功能語意，不只保存 endpoint；重要 API 要能回到 capability、operation、domain concept、state/error behavior 與 fixture。
- API 列表不是只列 method/path；需要 API Catalog、分組、逐支 API 欄位語意、抓取來源、覆蓋率缺口、UI/API 對照與 SDK/client 欄位用途。
- App 開發 guidance 要分離到 [`app-development-guidance`](../app-development-guidance/)，避免把分析方法與產品開發 checklist 混在一起。
- 用 APK 分析文件產生工具、SDK、client、mock、contract test 或重建功能時，`app-development-guidance` 必須同輪接手開發文件與 blocker questions。
- 開始 live-facing SDK/client/app tool 開發前，必須先確認 project docs 已具備最小可跑因素；若 runtime 需要授權身分材料，還必須有 self-generation audit，回答「能否自生成、怎麼生成、不能時誰提供、reset/cooldown/error 如何驗證」。缺 runtime baseline 的工作只能落在離線 parser、fixture、mock 或文件補齊範圍。
- 新發現要回饋到 skill，但必須去敏、泛化、可驗證。

## 最小產出

一次合格的 APK 分析至少要留下：

- 分析環境：OS、device/emulator、APK version、arch、root/Frida 狀態。
- 流量路徑判斷：localhost、whole-device pcap、proxy/MITM、Java hook、native/Dart 路徑。
- 證據：pcap、hook log、MITM export、反編譯搜尋結果或 sanitized excerpt。
- 結論：哪些路徑有效、哪些被排除、下一步如何驗證。
- 功能重建交接：capability、screen/route/operation、domain concept candidates、API/interface contract、state/error handling、data lifecycle、fixtures、open questions。
- API Catalog：API 總入口、分組索引、逐支 API 詳細文件、coverage/gap、UI/API mapping、SDK/client 欄位用途與 validation/open questions。
- Development readiness gate：若要開始 live-facing SDK/client/app tool，專案文件已回答最小可跑因素，包括必要時的 authorized identity material self-generation audit，或把缺口列為 blocker / scoped out；若只具 skeleton baseline，僅允許離線 parser、fixture、mock 或文件工作。
- 去敏規則：哪些值被遮蔽，哪些文件不能提交。
- app-development-guidance handoff：若目標包含 app 工具、SDK、client 或重建功能，要記錄已啟用該 skill 與交接文件位置。
- 可重用 lesson：新增檔於 **`feedback_history/<category>/`** 或 **`feedback_history/common/`**。
- 可選的 Developer Guidance Notes：若對自家 App 開發有設計、實作或安全啟發，連到或回饋至 **`app-development-guidance/`**。
