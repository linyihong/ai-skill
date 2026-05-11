# Claude 使用說明

本檔說明在 Claude 類工具中如何讓 Ai-skill 技巧穩定生效。核心原則與其他工具相同：中央庫是唯一真相來源，Claude 端只讀取、參照或同步快照。

## 自動配置

本專案在 `.claude/settings.json` 配置了 Claude Code 的規則加載入口。Claude Code 啟動時會自動讀取配置並指向：

- **Default Bootstrap 規則**: 見 [`shared-rules/README.md`](../shared-rules/README.md)（12-13 個必讀規則檔案）
- **架構層級導航**: 見 [`knowledge/indexes/README.md`](../knowledge/indexes/README.md)（任務路由索引）
- **Git 開發規則**: 自動指向開發分支 `claude/review-dev-rules-ndnTH`

Claude 不需要每次都被手動提示；只要打開此專案，配置會自動生效。

---

## 配置思想與邊界定義

### `.claude/settings.json` 的角色

`.claude/settings.json` 是一個**薄配置層**，只保留：
- 中央庫位置（`ai-skill_repo`）
- 指向規則的引用（rules、bootstrap、architecture）
- Git 操作配置（branch、push format）

**不在這裡存放**：
- Detailed bootstrap 清單（由 `shared-rules/README.md` 管理）
- 發現機制詳情（見本檔案下方）
- 工作流偏好（見 `.claude/settings.json` 的 `workflow_preferences`）

### 責任邊界

| 層級 | 位置 | 內容 |
|------|------|------|
| **通用 Claude 規則** | `ai-tools/claude.md` | 使用說明、操作注意、配置思想 |
| **Claude Code 配置入口** | `.claude/settings.json` | 配置項 + 指向規則的引用 |
| **Default Bootstrap** | `shared-rules/README.md` | 12-13 個必讀規則清單 |
| **架構層級導航** | `knowledge/indexes/README.md` | 任務意圖路由表 |

---

## 🔍 規則發現機制

### 核心概念

AI 在不同情境下應該讀什麼文件？

```
情境 → 查詢發現規則 → 找到對應文件 → 讀取並執行
```

### 發現規則（Situational Router）

| 情境 | 查詢 | 讀取 |
|------|------|------|
| **啟動/接手任務** | 什麼是這個知識庫？ | `README.md`、`shared-rules/README.md` |
| **規劃新任務** | 任務應該讀什麼？ | `knowledge/indexes/README.md`（任務路由）|
| **定義新規則** | 規則之間的依賴？ | `shared-rules/linked-updates.md` |
| **評估推廣/淘汰** | 生命週期如何管理？ | `governance/lifecycle/README.md` |
| **做架構決策** | 工程智慧與 trade-off？ | `intelligence/README.md` |
| **規劃複雜任務** | 流程如何分解？ | `workflow/README.md` |
| **提取教訓** | 失效如何學習？ | `shared-rules/failure-learning-system.md` |
| **不知道該怎辦** | 有沒有類似情況？ | `shared-rules/failure-patterns/` |

### 架構層級快速參考

| 層級 | 檔案 | 用途 |
|------|------|------|
| **knowledge** | `knowledge/indexes/README.md` | 任務意圖 → 文件對應 |
| **metadata** | `metadata/schema.md` | 知識原子結構與驗證 |
| **governance** | `governance/lifecycle/README.md` | 知識生命週期管理 |
| **intelligence** | `intelligence/README.md` | 工程決策與智慧 |
| **workflow** | `workflow/README.md` | 任務分解與流程 |
| **feedback** | `feedback/promotion/README.md` | 教訓推廣管道 |
| **memory** | `memory/README.md` | 長期記憶設計 |

## 何時使用

當你希望 Claude 使用本庫的某個 skill 時，不要只說「照技巧做」。請明確給 Claude：

- `<AI_SKILL_REPO>` 位置。
- 要使用的 `skills/<name>/SKILL.md`。
- 需要讀取的 workflow / tools / documentation / checklist。
- 必須先讀 `shared-rules/README.md`，再依任務讀相關 shared rules。
- 若任務跨多輪或可能中斷，使用 `<PROJECT_ROOT>/.agent-goals/` 追蹤目標閉環。

## 建議提示詞

通用提示：

```text
請使用 <skill-name> skill。

中央庫位置：
<AI_SKILL_REPO>

請先閱讀：
<AI_SKILL_REPO>/shared-rules/README.md
<AI_SKILL_REPO>/shared-rules/dependency-reading.md
<AI_SKILL_REPO>/shared-rules/linked-updates.md
<AI_SKILL_REPO>/shared-rules/conversation-goal-ledger.md
<AI_SKILL_REPO>/shared-rules/tool-neutral-documentation.md
<AI_SKILL_REPO>/shared-rules/rule-weight.md
<AI_SKILL_REPO>/shared-rules/decision-efficiency.md
<AI_SKILL_REPO>/shared-rules/failure-learning-system.md
<AI_SKILL_REPO>/shared-rules/document-todo-list.md
<AI_SKILL_REPO>/shared-rules/document-sizing.md
<AI_SKILL_REPO>/shared-rules/goal-action-validation.md
<AI_SKILL_REPO>/shared-rules/neutral-language.md
<AI_SKILL_REPO>/skills/<skill-name>/SKILL.md

再依該 skill 的 README / WORKFLOW / TOOLS / DOCUMENTATION / CHECKLIST 讀取必要依賴。
若某依賴檔不存在，請明確標示 not applicable，不要說已讀。

若本輪工作可能中斷或有多個目標，請依 conversation-goal-ledger 規則在 <PROJECT_ROOT>/.agent-goals/ 維護暫存目標，並明列 priority、parallelization mode、owner/lock、plan/todo links、missing/decision/strengthen。
若使用者指出你反覆失誤、更新到 tool mirror / project copy、漏讀依賴或漏做驗證，請依 failure-learning-system 分類失效模式，將可重用 lesson 推廣到 shared rule / failure pattern / skill feedback，並驗證防呆規則可被下次 agent 讀到。
完成 Ai-skill repo 變更時，請依 dependency-reading 完成 diff review、linked updates、tool sync、commit、push、讀回與 clean status。
```

APK 分析範例：

```text
請使用 apk-analysis skill 協助分析這個 APK。

中央庫位置：
<AI_SKILL_REPO>

請先讀：
<AI_SKILL_REPO>/shared-rules/README.md
<AI_SKILL_REPO>/shared-rules/dependency-reading.md
<AI_SKILL_REPO>/skills/apk-analysis/SKILL.md
<AI_SKILL_REPO>/skills/apk-analysis/RUNBOOK.md

再依 WORKFLOW.md / TOOLS.md / DOCUMENTATION.md 執行。
若得到可重用 lesson，依 shared-rules/feedback-lessons.md 寫入 feedback_history/，並完成 Ai-skill writeback transaction。
```

## Claude 操作注意

- Claude 不一定會自動探索整個 repo。提示時請列出入口檔與必要依賴，而不是只貼資料夾名。
- Claude 若只能看到單一專案，請同時提供 `<AI_SKILL_REPO>` 的可讀路徑，或把必要 skill/shared-rules 以工具支援的方式同步成可讀上下文。
- 如果 Claude 已經長時間對話，請先要求它讀 `<PROJECT_ROOT>/.agent-goals/README.md`，確認未完成項、待決策、優先順序、parallelization mode、owner/lock 與待補強內容。
- 如果 `.agent-goals/README.md` 顯示重疊目標已有其他 owner 或 active lock，要求 Claude 停止修改並提示使用者決定：等待、接手、拆成子目標或另開非重疊 goal。
- 如果 goal 標示 `single-owner` 或 `non-parallelizable`，不要讓 Claude 和其他 agent 分工同一流程；先取得使用者確認。
- 若 Claude 產生修改計畫，讓它把計畫 TODO 連到 `.agent-goals` 或文件前段的 `Document TODO`。
- 若 Claude 要改本庫，提醒它不要只更新文件；還要跑驗證、tool sync、commit、push、讀回和 clean status。
- 若 Claude 完成 `.agent-goals/` goal 後仍留下長期 roadmap、phase、migration、promotion、deprecation 或治理狀態，要求它先回寫到 durable planning 文件，再刪除 active goal；不要把 completed goal row 當長期記憶。

## 與 Tool Adapter 的關係

這裡只放 Claude 的通用使用方式。若某個 skill 針對 Claude 有特殊執行策略，例如上下文載入順序、prompt chunking、工具輸出限制或 Claude 專用 failure mode，放在：

```text
skills/<skill-name>/tool-adapters/claude.md
```

該 adapter 只寫 skill-specific 差異，並連回核心 `WORKFLOW.md` / `TOOLS.md` / `DOCUMENTATION.md`。不要把整個 skill 流程複製到 adapter。

## 驗證

使用 Claude 完成任務時，最後要求它回報：

- 讀了哪些 shared rules 與 skill 依賴。
- Default Bootstrap 是否已讀，以及哪些任務專屬規則後續補讀。
- 哪些依賴不存在，所以是 `not applicable`。
- 目標是否完成，還有哪些 `.agent-goals` 或 Document TODO 未完成。
- 驗證方法：diff review、lints、link check、測試、source check、tool sync、commit/push/readback/clean status。

← [回到 AI 工具索引](README.md)
