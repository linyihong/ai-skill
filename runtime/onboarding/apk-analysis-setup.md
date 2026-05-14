# APK Analysis Onboarding & Setup

本文件定義新 APK 分析專案的初始設定流程與提示詞模板。承接 [`workflow/apk-analysis/execution-flow.md`](../../workflow/apk-analysis/execution-flow.md) 與 [`workflow/apk-analysis/artifact-gates.md`](../../workflow/apk-analysis/artifact-gates.md) 的內容，提取為 tool-neutral 的 runtime onboarding 指引。

> **遷移狀態**：此文件為新分層的 reference target，舊 `skills/apk-analysis/` 已不再作為 active entrypoint。新內容請直接寫入此文件。

## 放置位置

有兩種用法：

### 1. 共用 skill repo（優先）

直接從本 repository 讀取：

```text
runtime/onboarding/apk-analysis-setup.md                      # 本文件：開場提示詞與放置方式
workflow/apk-analysis/execution-flow.md                       # 分析執行流程
workflow/apk-analysis/artifact-gates.md                       # 產出格式與完成定義
analysis/apk/workflows/                                       # 操作細節流程（依分類）
intelligence/engineering/analytical-reasoning/                        # 決策智慧（heuristics / anti-patterns / signals / failure）
```

新技巧回饋到：

```text
shared-rules/feedback-lessons.md                              # 檔名規則與 lesson 模板（全庫共用）
feedback/history/apk-analysis/<category>/                     # 每一條 lesson 獨立檔
```

驗證後同步更新：

```text
workflow/apk-analysis/execution-flow.md
workflow/apk-analysis/artifact-gates.md
analysis/apk/workflows/<workflow>.md
intelligence/engineering/analytical-reasoning/<atom-type>/<name>.md
```

### 2. 當作專案文件使用

使用時要明確告訴 AI：

```text
請先閱讀 runtime/onboarding/apk-analysis-setup.md，
並依照 workflow/apk-analysis/execution-flow.md 與 workflow/apk-analysis/artifact-gates.md 分析這個 APK。
如果過程中學到可重用技巧，請依 shared-rules/feedback-lessons.md 在 feedback/history/apk-analysis/<category>/ 新增檔案回饋；跨分類用 feedback/history/apk-analysis/common/。
```

### 3. 當作工具可讀 skill 使用

將 `workflow/apk-analysis/`、`analysis/apk/`、`intelligence/engineering/analytical-reasoning/` 與 `shared-rules/` 成對部署到你使用的 AI / agent 工具，或在工具提示中明確指定 repository 內的路徑。工具專屬部署方式請看 [`ai-tools/README.md`](../../ai-tools/README.md)，不要把工具專屬路徑寫進通用流程。

## 開場提示詞

新專案開始時，可以直接給 AI：

```text
請使用 apk-analysis skill 協助分析這個 APK。

目標：
- 先判斷核心流量走哪一層。
- 不要一開始假設是 pinning。
- 先做 localhost / pcap / proxy / Java hook / native or Flutter 判斷。
- 所有 token、device id、私密 host、個資都要去敏。
- 如果發現新技巧，請**主動**在 feedback/history/apk-analysis/<category>/ **新增 lesson 檔**（跨分類用 common/；格式見 shared-rules/feedback-lessons.md；不要等使用者說「記得回饋」），且要讓人類也看得懂；可同步更新 workflow/apk-analysis/execution-flow.md、workflow/apk-analysis/artifact-gates.md、analysis/apk/workflows/ 或 intelligence/engineering/analytical-reasoning/ 如已驗證。

目前材料：
- APK:
- version:
- package:
- device / emulator:
- root / Frida:
- proxy / MITM:
- 想先分析的功能流程:
```

## 讓 AI 自動回饋的提示

在新專案中可以加這句：

```text
每次你發現新的可重用 APK 分析技巧時，請不要只在對話裡說明；也請不要等使用者提醒「記得回饋」。
如果使用者或 reviewer 提出一個可能跨 APK 重用的操作改進、失敗模式或驗證規則，請先做泛化判斷；能泛化就立刻進 feedback loop。
請依 shared-rules/feedback-lessons.md，在同一輪對話於 feedback/history/apk-analysis/<category>/ 或 feedback/history/apk-analysis/common/ **新增一個** lesson 檔。
如果 lesson 尚未完全驗證，先標 candidate / experimental 並寫清楚 validation criteria；如果已經被驗證，也請同步更新 workflow/apk-analysis/execution-flow.md、workflow/apk-analysis/artifact-gates.md、analysis/apk/workflows/ 或 intelligence/engineering/analytical-reasoning/。
注意：lesson 內容要讓人類也能看懂，不只給 AI 看。
完成後請在 repository commit 並 push。
```

## 與其他層的關係

- `workflow/apk-analysis/execution-flow.md` 提供分析執行流程，本文件提供如何啟動分析的設定指引。
- `workflow/apk-analysis/artifact-gates.md` 提供產出格式與完成定義。
- `shared-rules/feedback-lessons.md` 提供 lesson 格式與回饋規則。
- `skills/apk-analysis/RUNBOOK.md` 是原始來源，已不再作為 active entrypoint（舊 `skills/` 結構已於 2026-05-13 標記為 deprecated）。
