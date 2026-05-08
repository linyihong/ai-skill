# Claude 使用說明

本檔說明在 Claude 類工具中如何讓 Ai-skill 技巧穩定生效。核心原則與其他工具相同：中央庫是唯一真相來源，Claude 端只讀取、參照或同步快照。

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
<AI_SKILL_REPO>/shared-rules/decision-efficiency.md
<AI_SKILL_REPO>/shared-rules/document-todo-list.md
<AI_SKILL_REPO>/shared-rules/document-sizing.md
<AI_SKILL_REPO>/shared-rules/goal-action-validation.md
<AI_SKILL_REPO>/shared-rules/neutral-language.md
<AI_SKILL_REPO>/skills/<skill-name>/SKILL.md

再依該 skill 的 README / WORKFLOW / TOOLS / DOCUMENTATION / CHECKLIST 讀取必要依賴。
若某依賴檔不存在，請明確標示 not applicable，不要說已讀。

若本輪工作可能中斷或有多個目標，請依 conversation-goal-ledger 規則在 <PROJECT_ROOT>/.agent-goals/ 維護暫存目標，並明列 priority、parallelization mode、owner/lock、plan/todo links、missing/decision/strengthen。
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
