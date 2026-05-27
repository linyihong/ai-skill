# Analysis

`analysis/` 負責「如何觀察、拆解與取得證據」。本層保存可重用的分析方法、觀察框架、問題拆解方式與 pattern extraction 路線，讓 agent 先知道如何取得證據，再決定是否沉澱成 workflow 或 intelligence。

本層不是 raw case archive。單一 App、專案或一次 live run 的 raw logs、pcap、Frida output、class dump、host、endpoint、token、payload fragment 與執行證據，應留在業務專案的受控文件或 evidence 位置；若要進入本庫，必須先去敏並抽象成可重用的方法、lesson 或 decision atom。

## 目前入口

- [`apk/`](apk/README.md)：APK 分析的可重用分析方法。已從舊 `skills/apk-analysis/` 遷移至本層。
- [`development-guidance/`](development-guidance/README.md)：開發指引的分析方法（風險翻譯、控制層選擇）。
- [`repo/`](repo/README.md)：Repository 分析與理解方法（結構觀察、依賴分析、心智模型建立、技術債評估）。
- [`production/`](production/README.md)：Production 問題分析與根因追蹤方法（incident 分類、觀測性資料判讀、RCA 流程、效能診斷）。
- [`issue/`](issue/README.md)：Issue 分類與優先級判斷方法（triage 流程、優先級計算、初步診斷、重複偵測）。
- [`web/`](web/README.md)：Web Scraping 分析方法（網站結構觀察、反爬機制判讀、動態內容分析、MCP 工具設計評估）。
- [`security/`](security/README.md)：安全相關分析方法（token 流向、簽章一致性、雙簽章/雙加密 token 並存的審計）。
- [`ai-augmented-delivery/`](ai-augmented-delivery/README.md)：AI codegen 工具大幅進入開發流程後的觀察方法、量化資料與問題解剖（缺陷分布、perf test bottleneck anatomy）。

## 放什麼

- 授權情境下的技術觀察與拆解方法。
- Repo、APK、production issue 或其他系統的分析路線。
- 將 raw evidence 轉成可重用 pattern 的 extraction 方法。
- 可被 workflow 引用的分析檢查點與判讀順序。
- 可重複使用的證據取得與 triage 方法，例如「如何判斷流量走哪一層」。

## 不放什麼

- Raw logs、未去敏 traces、pcap、Frida output、class dump 或 command history。
- 某個 App / 專案專屬的 findings、host、endpoint、sample、payload 或 live run 結果。
- Agent 執行步驟與 orchestration flow；放到 `workflow/`。
- 工程 trade-off、架構 lesson、anti-pattern 結論；放到 `intelligence/`。
- 可執行 shared policy；放到 `enforcement/`。
- 工具專屬路徑、hook 或 UI；放到 `ai-tools/`。

## 快速判斷

新增內容前先問：

| 問題 | 放置位置 |
| --- | --- |
| 主要是在回答「如何取得與拆解證據？」 | `analysis/` |
| 主要是在回答「先做 A 再做 B？」 | `workflow/` |
| 主要是在回答「未來遇到類似情況如何判斷？」 | `intelligence/` |
| 主要是在保存「這次實際發生了什麼？」且脫離案例就失效 | 業務專案 evidence / 去敏後的 `feedback/history/` |

## 誰會參考這裡（Inbound References）

- [`route.workflow.apk-analysis`](../knowledge/runtime/routing-registry.yaml) — candidate_sources 引用 `analysis/apk/README.md`
- [`route.workflow.software-delivery`](../knowledge/runtime/routing-registry.yaml) — candidate_sources 引用 `analysis/development-guidance/README.md`
- [`route.workflow.travel-planning`](../knowledge/runtime/routing-registry.yaml) — candidate_sources 引用 `analysis/travel/README.md`
- [`route.intelligence.apk-analysis.atoms`](../knowledge/runtime/routing-registry.yaml) — required_dependencies 引用 `analysis/apk/README.md`
- [`route.analysis.apk.workflows`](../knowledge/runtime/routing-registry.yaml) — required_dependencies 引用 `analysis/apk/README.md`

## 與既有層的關係

- `skills/` 已 deprecated，`apk-analysis` 與 `travel-planning` 已遷移至本層與 `workflow/`。
- `workflow/` 可以引用本層，但不應複製完整分析知識。
- `intelligence/` 承接從分析結果萃取出的工程判斷。
- `enforcement/` 仍負責授權、去敏、依賴讀取與 validation policy。

## 遷移狀態

- `skills/apk-analysis/techniques/` — ✅ 已遷移至 `analysis/apk/workflows/`
- `skills/apk-analysis/TOOLS.md` — ✅ 已遷移至 `analysis/apk/tools-and-failures.md`
- `skills/apk-analysis/WORKFLOW.md`（分析部分） — ✅ 已遷移至 `analysis/apk/traffic-triage.md`
- `skills/travel-planning/TOOLS.md` — ✅ 已遷移至 `analysis/travel/sources-and-tools.md`
- `skills/travel-planning/README.md` — ✅ 已遷移至 `analysis/travel/README.md`
- `skills/app-development-guidance/` — ✅ 已遷移至新分層，舊目錄已刪除
