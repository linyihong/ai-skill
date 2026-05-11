# AI-native Knowledge Operating System

這個 repository 是 AI-native Knowledge Operating System 的中央知識庫：以 Git 維護單一真相來源，讓 agent 直接 reference `<AI_SKILL_REPO>` 讀取 rules、skills、tool adapters、goal state guidance、failure learning 與 close-loop automation。

未來不同專案需要 APK 分析、app/API/embedded development guidance、旅行規劃或其他可重用 agent 能力時，優先從這裡讀 skill；新的可重用技巧、規則與架構方向也回寫到這裡。工具端的 symlink、bundle 或 copy snapshot 只是相容層，不取代本 repository。

**路徑約定：**在本機上，中央庫位置就是你心裡的 `<AI_SKILL_REPO>`（clone 所在目錄），且必須是本 repository 的 git root。`~/.cursor/skills*`、`~/.cursor/shared-rules`、`~/.cursor/bundles/*` 或專案 `.cursor/` 是工具部署 / mirror 路徑，不能取代 source repo。**請勿把真實本機絕對路徑寫進本庫任何會 commit 的文件**（含 `feedback_history`、`shared-rules`、規則範本）；對外一律用占位符 `<AI_SKILL_REPO>`、`<PROJECT_ROOT>`。

## 作業層

| 層級 | 路徑 | 用途 |
| --- | --- | --- |
| 架構方向 | [`architecture/`](architecture/) | 定義 AI-native Knowledge Operating System roadmap、相容層盤點與遷移條件。 |
| 共用作業規則 | [`shared-rules/`](shared-rules/README.md) | 跨 skill 的依賴讀取、連動更新、驗證、goal ledger、failure learning、語言用語與文件邊界規則。 |
| 能力模組 | [`skills/`](skills/README.md) | 可重用 agent 能力、workflow、checklist、文件模板、techniques 與 feedback lessons。 |
| 工具 adapter | [`ai-tools/`](ai-tools/README.md) | 工具專屬 reference、symlink、bundle、hook、UI 與 troubleshooting 指引。 |
| 閉環自動化 | [`scripts/`](scripts/README.md) | goal ledger helper、保守分組 commit/push automation，以及可選 tool sync bridge。 |
| 暫存目標狀態 | `.agent-goals/` | 專案本地 active goal ledger；不提交，完成並驗證後刪除。 |

## 下一階段正式分層

這些 top-level directories 先定義 AI-native Knowledge Operating System 的責任邊界；目前不代表既有 `skills/`、`shared-rules/` 或 `ai-tools/` 已完成遷移。

| 層級 | 路徑 | 目前用途 |
| --- | --- | --- |
| 分析方法 | [`analysis/`](analysis/README.md) | 觀察、拆解、pattern extraction 與分析路線。 |
| 工程智慧 | [`intelligence/`](intelligence/README.md) | Engineering decision、trade-off、anti-pattern 與 reusable domain knowledge。 |
| 執行流程 | [`workflow/`](workflow/README.md) | Planning、decomposition、review、orchestration 與 handoff flow。 |
| Runtime 設計 | [`runtime/`](runtime/README.md) | Dynamic loading、context routing、context pruning 與 coordination design。 |
| 長期記憶 | [`memory/`](memory/README.md) | Episodic/project/failure memory 的抽象化與 replay 邊界。 |
| Feedback 演化 | [`feedback/`](feedback/README.md) | Lesson extraction、refinement、promotion 與 feedback loop。 |
| Model 協作 | [`models/`](models/README.md) | Model capability profile、routing、compression 與 prompt adaptation。 |
| 知識治理 | [`governance/`](governance/README.md) | Lifecycle、cleanup、splitting、validation 與 dependency maintenance design。 |
| 知識導航 | [`knowledge/`](knowledge/README.md) | Knowledge atoms、indexes、summaries、graphs 與 runtime navigation。 |
| 控制資料 | [`metadata/`](metadata/README.md) | Knowledge Atom schema、ranking、confidence、compatibility 與 runtime metadata。 |

## Agent 作業流程

1. 讀取本 `README.md`，先理解 OS layout。
2. 從 [`shared-rules/README.md`](shared-rules/README.md) 載入 shared-rule bootstrap。
3. 只讀取任務相關的 skill 入口與依賴文件。
4. 任務跨多步驟時，在 `<PROJECT_ROOT>/.agent-goals/` 追蹤可恢復目標。
5. 可重用知識只編輯 `<AI_SKILL_REPO>` 的 canonical source，不改 tool mirror。
6. 驗證 linked updates，commit、push、讀回更新入口，並確認 clean status。

`reference-first` 是預設：agent 直接讀本 repository。`symlink`、`bundle`、`copy snapshot` 只是在工具無法穩定 reference 中央庫時使用的相容層。

## 新專案 Quickstart

新專案啟動 agent 時，可直接使用這段 prompt：

```text
Use the AI-native Knowledge Operating System.

Canonical repository:
<AI_SKILL_REPO>

Project:
<PROJECT_ROOT>

First read:
<AI_SKILL_REPO>/README.md
<AI_SKILL_REPO>/shared-rules/README.md

Load the Default Bootstrap from shared-rules/README.md, including:
- dependency-reading.md
- linked-updates.md
- conversation-goal-ledger.md
- tool-neutral-documentation.md
- rule-weight.md
- decision-efficiency.md
- failure-learning-system.md
- document-todo-list.md
- document-sizing.md
- goal-action-validation.md
- neutral-language.md

Then choose the task-relevant skill from <AI_SKILL_REPO>/skills/.
If the task spans multiple steps, first read or initialize <PROJECT_ROOT>/.agent-goals/.
For reusable knowledge updates, edit the canonical repository only, not .cursor, ~/.cursor, bundles, or copied snapshots.
For repository changes, validate linked updates, commit, push, reread changed entries, and confirm clean git status.
```

一般情境只需要 `reference-first`。只有目前工具無法穩定直接讀 `<AI_SKILL_REPO>` 時，才使用 `symlink`、`bundle` 或 `copy snapshot`；工具專屬設定放在 [`ai-tools/`](ai-tools/README.md)。

## 能力模組

| Skill | 用途 |
| --- | --- |
| `skills/apk-analysis/` | 授權 APK 流量分析、動態抓包、Flutter/Dart AOT、response 解碼、文件化與 `feedback_history/` 技巧條目。 |
| `skills/app-development-guidance/` | 將授權 App/API/embedded/firmware/hardware-product 觀察轉成開發 guidance、實作模式、控制項、檢查清單與驗證方法。 |
| `skills/travel-planning/` | 依目的地、日期、交通與玩法查證營業時間、精準 Google Maps 點位、旅行社/套裝行程參考、飛機/新幹線/巴士/渡輪/自駕等長距離交通比較、玩法/餐飲與當地評分工具篩選、行程時間壓力、天氣/道路風險、交通訂票與票價/自駕成本、加油/充電補給、住宿/民宿或車中泊安靜度、路線折返檢查、可加入行事曆/地圖/提醒 app 的欄位與備案。 |

**目錄約定：**  
- **`skills/`**：各 skill 技巧包；之後新增 skill 放在 `skills/<name>/`，步驟見 [`skills/ADDING_SKILLS.md`](skills/ADDING_SKILLS.md)。  
- **`shared-rules/`**：**共用規則**（依主題分檔：授權、去敏、中性低爭議文件用語、工具中立文件、目標/執行/驗證流程、對話目標閉環、失效學習系統、依賴文件讀取、內容分層、文件大小與拆分、cross-skill reference、**feedback 檔名／模板／agent 行為**、工具同步等）；索引為 [`shared-rules/README.md`](shared-rules/README.md)，feedback 流程與模板集中在 [`feedback-lessons.md`](shared-rules/feedback-lessons.md)。各 skill 目錄下 **`FEEDBACK.md`** 若存在，僅為**一行入口**，不必重複維護正文。
- **`architecture/`**：repo-level 架構與 roadmap，不是可執行 shared rule；目前方向見 [`architecture/ai-native-knowledge-operating-system.md`](architecture/ai-native-knowledge-operating-system.md)。
- **`ai-tools/`**：各 AI coding / agent 工具如何讀取、參照或同步本知識庫；工具專屬路徑、hook、UI 與同步細節都放在這裡，不寫進通用 skill / shared rule 正文。
- **連動更新**：若改動會影響其他文件、索引、skill 入口、同步流程或分類文件，相關檔案**必須**依 [`shared-rules/linked-updates.md`](shared-rules/linked-updates.md) 同步更新或明確檢查，不得說成「可選」。  
- **每一則 `feedback_history`**：**不要**重複貼上共用規則全文，頂部引用 `shared-rules/` 即可。工具端同步與部署方式請看 [`ai-tools/`](ai-tools/README.md)。

## AI 工具使用說明

詳細文件放在 [`ai-tools/`](ai-tools/)，根 README 只保留工具索引。新增其他 AI 工具時，在該目錄新增子檔並更新索引。

| 工具 | 文件 | 用途 |
| --- | --- | --- |
| 工具整合 | [`ai-tools/`](ai-tools/README.md) | 各工具如何讀取、同步與套用本知識庫；工具專屬內容集中在此。 |

## 預設 shared rules 載入

每次開啟 agent、新 session 或接手長對話時，先讀 [`shared-rules/README.md`](shared-rules/README.md) 的 **Default Bootstrap**。它只載入索引與必讀規則；後續仍依任務讀 skill-specific workflow、tools、documentation 與其他 shared rules。

## 架構 Roadmap

詳細 roadmap、相容層盤點與 copy/bundle sync 移除條件放在 [`architecture/ai-native-knowledge-operating-system.md`](architecture/ai-native-knowledge-operating-system.md)。可執行政策保留在 `shared-rules/`；工具專屬設定保留在 `ai-tools/`。

下一階段完整升級規劃放在 [`architecture/next-stage-upgrade-plan.md`](architecture/next-stage-upgrade-plan.md)，用來規劃 `analysis/`、`intelligence/`、`workflow/`、`runtime/`、`memory/`、`feedback/`、`models/`、`governance/`、`knowledge/`、`metadata/` 等正式分層。

## 對話目標閉環

若工作可能中斷、跨多輪、被多 agent 接手、拆成多個子目標、已建立 TodoWrite、使用者要求繼續前一個多步驟任務，或 agent 已看到 active project 有 modified / staged / untracked files，依 [`shared-rules/conversation-goal-ledger.md`](shared-rules/conversation-goal-ledger.md) 在業務專案本地先檢查或維護 `<PROJECT_ROOT>/.agent-goals/`。Goal 需要明確列出 priority、parallelization mode、owner/lock 決策、plan/todo links、missing/decision/strengthen、next action、completion criteria 與 validation。這是暫存狀態，不進 git；目標完成並驗證後刪除。可用 [`scripts/agent-goals.sh`](scripts/agent-goals.sh) 建立、更新、拆解、暫停與完成刪除 goal。

## 失效學習系統

若使用者指出 agent 反覆失誤、寫錯 source/mirror、漏讀依賴、忘記目標、漏做驗證或閉環不完整，依 [`shared-rules/failure-learning-system.md`](shared-rules/failure-learning-system.md) 分類失效模式並選擇 promotion target。可重用跨 skill 的失效模式放在 [`shared-rules/failure-patterns/`](shared-rules/failure-patterns/README.md)；skill-specific lesson 仍放對應 skill 的 `feedback_history/`；專案 incident 證據留在專案文件，不進 reusable docs。

## 回饋規則

可以回饋：

- 新的分析流程。
- 新的工具使用技巧。
- 新的失敗判讀方式。
- 新的去敏與 fixture 沉澱方式。
- 通用的媒體、解密、session、proxy、Frida 經驗。

不要回饋：

- 特定產品的完整 host / endpoint。
- token、secret、device id、帳號或個資。
- 未去敏 raw response。
- 本機絕對路徑、使用者名稱、私有工作目錄、clone 位置；請改用 `<AI_SKILL_REPO>`、`<PROJECT_ROOT>`、`<WORKSPACE>`。
- 只對單一產品有效、沒有泛化價值的結論。

## Git 規則

這個 repository 有自己的 git remote。更新 skill 後：

```bash
cd <AI_SKILL_REPO>
./scripts/ai-skill-close-loop.sh          # dry-run：檢查 lock 與分組
./scripts/ai-skill-close-loop.sh --commit # 沒有 active lock 時分組提交
./scripts/ai-skill-close-loop.sh --commit --push
```

Commit/push 與必要的工具同步完成後，agent 還要重新讀取本次更新過的 skill/shared-rule 入口與主要依賴文件，並依 `shared-rules/dependency-reading.md` 留下 dependency read ledger（已讀、缺檔不適用、阻塞項與驗證），避免使用 commit 前的舊上下文。

若偵測到其他 agent / user 的 active close-loop lock，停止自動 commit / push，回報目前 owner 與下一步；不要混合提交他人的變更。

不要把專案私有資料、抓包原始檔或未去敏樣本 commit 到這裡。

工具設定目錄的變更若屬於業務專案，請勿誤把該專案才有的機密或絕對路徑抄進本庫；本庫維持泛化與 `<AI_SKILL_REPO>` / `<PROJECT_ROOT>` 等占位符。
