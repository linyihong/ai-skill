# APK Analysis Workflow（APK 分析工作流程）

`workflow/apk-analysis/` 是 APK 分析執行流程的主要目錄，包含 tool-neutral 的執行步驟與產出規範。

## Scope（範圍）

此層擁有：

- 授權 APK 分析的 agent 執行序列。
- 任務分解、捕獲窗口、文件閉環和交接流程。
- 從工作流程步驟到 `analysis/apk/` 方法和 `intelligence/engineering/apk-analysis/` 課程的參考。
- 當 APK 發現成為應用程式/API/SDK 開發指引時的跨技能交接規則。

此層不擁有：

- 深度執行時期或流量分析技術內容；使用 `analysis/apk/`。
- 穩定的工程智慧和反模式；使用 `intelligence/engineering/apk-analysis/`。
- 工具特定的 UI、本地鏡像、hook 安裝或同步細節；使用 `ai-tools/` 或專案文件。
- 專案特定的發現、原始證據或私人服務細節。

## Current Source References（當前來源參考）

| 工作流程關注點 | 原始來源 | 提取狀態 |
| --- | --- | --- |
| Skill 觸發和授權邊界 | `../../skills/apk-analysis/SKILL.md` | 參考用，內容已遷移至本層 |
| 預設決策樹 | `../../skills/apk-analysis/WORKFLOW.md` | ✅ 已提取到 [`execution-flow.md`](execution-flow.md) |
| 捕獲窗口詳細規則 | `../../skills/apk-analysis/WORKFLOW.md` | ✅ 已提取到 [`execution-flow.md`](execution-flow.md) |
| 環境和工具準備 | `../../skills/apk-analysis/TOOLS.md` | 從 `analysis/apk/` 參考；不重複 |
| 文件和產出規範 | `../../skills/apk-analysis/DOCUMENTATION.md` | ✅ 已提取到 [`artifact-gates.md`](artifact-gates.md) |
| SDK 即時自我生成審計 | `../../skills/apk-analysis/DOCUMENTATION.md` | ✅ 已提取到 [`artifact-gates.md`](artifact-gates.md) |
| 身份資料自我生成審計 | `../../skills/apk-analysis/DOCUMENTATION.md` | ✅ 已提取到 [`artifact-gates.md`](artifact-gates.md) |
| UI 架構圖模板 | `../../skills/apk-analysis/DOCUMENTATION.md` | ✅ 已提取到 [`artifact-gates.md`](artifact-gates.md) |
| API 目錄詳細要求 | `../../skills/apk-analysis/DOCUMENTATION.md` | ✅ 已提取到 [`artifact-gates.md`](artifact-gates.md) |
| 清理規則 | `../../skills/apk-analysis/DOCUMENTATION.md` | ✅ 已提取到 [`artifact-gates.md`](artifact-gates.md) |
| 開發者指引備註 | `../../skills/apk-analysis/DOCUMENTATION.md` | ✅ 已提取到 [`artifact-gates.md`](artifact-gates.md) |
| 回饋課程撰寫技巧 | `../../skills/apk-analysis/DOCUMENTATION.md` | ✅ 已提取到 [`artifact-gates.md`](artifact-gates.md) |
| 回填規則 | `../../skills/apk-analysis/DOCUMENTATION.md` | ✅ 已提取到 [`artifact-gates.md`](artifact-gates.md) |
| 第一天專案 runbook | `../../skills/apk-analysis/RUNBOOK.md` | 候選用於 onboarding workflow |
| 開發指引交接 | `../../skills/app-development-guidance/` | 僅在觸發條件適用時參考 |

## Reference-First Workflow Shape（參考優先的工作流程形狀）

1. 確認授權、範圍、APK 版本、裝置/模擬器和允許的操作。
2. 透過 `analysis/apk/` 路由以識別流量/執行時期路徑。
3. 從 `analysis/apk/workflows/` 載入匹配的工作流程，並從 `intelligence/engineering/apk-analysis/` 載入智慧原子。
4. 產生已清理的專案產出：UI 地圖、操作到 API 矩陣、API 目錄、執行時期基線、fixtures 和適用的開放問題。
5. 如果發現必須成為應用程式/API/SDK 實作指引，使用已清理的 Feature Reconstruction Handoff 交接給 `workflow/app-development-guidance/`。
6. 如果出現可重複使用的課程，將課程保留在 `feedback/history/apk-analysis/` 中，直到提升規則將其移至智慧或回饋層。

## 與既有層的關係

- `workflow/apk-analysis/` 是 APK 分析執行流程的主要入口。所有 agent 應優先參考本目錄的內容。
- `analysis/apk/` 提供深度技術方法（traffic triage、Frida hook、proxy 架構等），被本 workflow 引用。
- `intelligence/engineering/apk-analysis/` 提供從 APK 分析中萃取的工程智慧（啟發式、反模式、失敗模式），被本 workflow 引用。
- `feedback/history/apk-analysis/` 儲存 APK 分析的具體課程記錄，可被提升至 intelligence 層。
- `skills/apk-analysis/` 是原始 skill 目錄，內容已逐步遷移至本層。新內容應直接寫入本層。
