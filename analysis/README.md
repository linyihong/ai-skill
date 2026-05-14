# Analysis

`analysis/` 負責「如何觀察與拆解」。本層保存可重用的分析方法、觀察框架、問題拆解方式與 pattern extraction 路線，讓 agent 先知道如何取得證據，再決定是否沉澱成 workflow 或 intelligence。

## 目前入口

- [`apk/`](apk/README.md)：`apk-analysis` pilot 的分析方法候選目的地；目前仍 reference `skills/apk-analysis/`，不搬移大量內容。
- [`development-guidance/`](development-guidance/README.md)：開發指引的分析方法（風險翻譯、控制層選擇）。
- [`repo/`](repo/README.md)：Repository 分析與理解方法（結構觀察、依賴分析、心智模型建立、技術債評估）。
- [`production/`](production/README.md)：Production 問題分析與根因追蹤方法（incident 分類、觀測性資料判讀、RCA 流程、效能診斷）。
- [`issue/`](issue/README.md)：Issue 分類與優先級判斷方法（triage 流程、優先級計算、初步診斷、重複偵測）。

## 放什麼

- 授權情境下的技術觀察與拆解方法。
- Repo、APK、production issue 或其他系統的分析路線。
- 將 raw evidence 轉成可重用 pattern 的 extraction 方法。
- 可被 workflow 引用的分析檢查點與判讀順序。

## 不放什麼

- Agent 執行步驟與 orchestration flow；放到 `workflow/`。
- 工程 trade-off、架構 lesson、anti-pattern 結論；放到 `intelligence/`。
- 可執行 shared policy；放到 `shared-rules/`。
- 工具專屬路徑、hook 或 UI；放到 `ai-tools/`。

## 與既有層的關係

- `skills/` 目前仍是相容入口；本層只承接逐步抽出的分析方法。
- `workflow/` 可以引用本層，但不應複製完整分析知識。
- `intelligence/` 承接從分析結果萃取出的工程判斷。
- `shared-rules/` 仍負責授權、去敏、依賴讀取與 validation policy。

## 第一批候選遷移來源

- `skills/apk-analysis/techniques/`（已刪除，內容已遷移至 `analysis/apk/workflows/`）
- `skills/app-development-guidance/process/` 中偏分析與 discovery 的內容
- `architecture/next-stage-upgrade-plan.md` 中 `analysis/` 的分層說明
