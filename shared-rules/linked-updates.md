# 連動更新（全庫必須規則）

本檔是 **Ai-skill repository 全部文件與 skill 的連動更新規則**。若某項改動會影響其他文件、索引、同步流程、skill 入口、分類文件或範本，相關檔案 **必須**在同一次變更中更新或明確檢查，不可寫成「可選」或「之後再說」。

## Agent 必須做的事

1. 改任何 `shared-rules/`、`skills/`、根 `README.md`、同步腳本或模板前，先判斷是否有連動文件。
2. 若有連動文件，**必須**同步修改或明確寫下「已檢查，無需更新」的理由。
3. 第一次寫入 Ai-skill 或其工具同步路徑時，依 [`dependency-reading.md`](dependency-reading.md) 開啟 writeback transaction；連動更新、sync、commit、push、讀回與 clean status 都完成後才可關閉。
4. 若本輪明確使用或更新本機工具 mirror / symlink / copy snapshot，**必須**執行對應 tool sync；reference-only 策略不需要同步，具體工具命令放在 [`ai-tools/`](../ai-tools/README.md)。
5. 若改動 Ai-skill repo，除非使用者明講不要提交，**必須** `git add` → `commit` → `push`。
6. Commit/push 與必要 tool sync 完成後，**必須**依 [`dependency-reading.md`](dependency-reading.md) 重新讀取本次更新過的 skill/shared-rule 入口與主要依賴文件，確認 agent context 已載入最新版；reference-only 時 tool sync 記為不適用。
7. 最終回覆前必須再次執行 `git status --short --branch`；若仍有 modified/untracked/staged changes，或 branch ahead/behind remote，不得回覆「完成」，必須處理到乾淨或明確說明阻塞。
8. 回覆使用者時要說明已做了哪些連動更新；如果某些相關文件不用改，也要簡短說明原因。
9. 若發現 Ai-skill 有尚未推送、尚未合併、ahead/behind、或其他 pending commit 狀態，且使用者未明確要求 push / merge，最終回覆必須主動提醒目前狀態與下一步，不可只說「已更新」。

## 常見連動關係

| 改動位置 | 必須同步更新或檢查 |
| --- | --- |
| `shared-rules/README.md` 或新增 shared rule | 根 `README.md`、相關 skill 的入口說明（`skills/<name>/SKILL.md` 或 `workflow/<domain>/execution-flow.md`）、`feedback/history/` 模板引用。 |
| `shared-rules/reusable-guidance-boundary.md` | `shared-rules/README.md`、`shared-rules/content-layering.md`、`shared-rules/feedback-lessons.md`、`shared-rules/sanitization.md`、`shared-rules/goal-action-validation.md`、相關 skill 的入口（`skills/<name>/SKILL.md` 或 `workflow/<domain>/execution-flow.md`），以及已新增 lesson 的 promotion target 與 index。 |
| `shared-rules/dependency-reading.md` | `shared-rules/README.md`、`shared-rules/content-layering.md`、`shared-rules/linked-updates.md`、工具專用 always-apply agent rule、`skills/_template/SKILL.md`、`skills/ADDING_SKILLS.md`、所有現有 skill 的入口（`skills/<name>/SKILL.md` 或 `workflow/<domain>/execution-flow.md`）與根 `README.md`。若變更涉及重讀 / reload 防呆，這些入口必須提醒 agent 建立 dependency read ledger，並明列不存在的依賴檔為 `not applicable`。 |
| `shared-rules/neutral-language.md` | `shared-rules/README.md`、`shared-rules/content-layering.md`、`shared-rules/feedback-lessons.md`、`skills/_template/SKILL.md`、`skills/ADDING_SKILLS.md`、所有現有 skill 的入口（`skills/<name>/SKILL.md` 或 `workflow/<domain>/execution-flow.md`）與根 `README.md`。 |
| `shared-rules/tool-neutral-documentation.md` | `shared-rules/README.md`、`shared-rules/content-layering.md`、`shared-rules/linked-updates.md`、`shared-rules/dependency-reading.md`、根 `README.md`、`skills/README.md`、`skills/ADDING_SKILLS.md`、各 skill 的入口/README、`skills/*/tool-adapters/` 索引（或 `tools/adapters/`）、`ai-tools/README.md` 與受影響工具文件。 |
| `shared-rules/rule-weight.md` | `shared-rules/README.md`、`shared-rules/content-layering.md`、`shared-rules/dependency-reading.md`、`shared-rules/linked-updates.md`、`shared-rules/decision-efficiency.md`、`shared-rules/goal-action-validation.md`、工具專用 always-apply agent rule、`ai-tools/README.md`、`ai-tools/cursor.md`、`ai-tools/claude.md`。 |
| `shared-rules/decision-efficiency.md` | `shared-rules/README.md`、`shared-rules/content-layering.md`、`shared-rules/dependency-reading.md`、`shared-rules/linked-updates.md`、`governance/document-sizing.md`、有決策路由/context-loading 指引的 skill workflow 或 README（`workflow/<domain>/` 或 `skills/<name>/`）。 |
| `shared-rules/failure-learning-system.md` 或 `shared-rules/failure-patterns/` | `shared-rules/README.md`、`shared-rules/content-layering.md`、`shared-rules/dependency-reading.md`、`shared-rules/linked-updates.md`、`shared-rules/reusable-guidance-boundary.md`、`shared-rules/feedback-lessons.md`、`shared-rules/goal-action-validation.md`、受影響工具文件與被補強的 shared rule / skill workflow。 |
| `shared-rules/document-todo-list.md` | `shared-rules/README.md`、`shared-rules/content-layering.md`、`shared-rules/linked-updates.md`、`shared-rules/dependency-reading.md`、`shared-rules/conversation-goal-ledger.md`、`skills/ADDING_SKILLS.md`、相關模板與文件。 |
| `shared-rules/goal-action-validation.md` | `shared-rules/README.md`、`shared-rules/content-layering.md`、`shared-rules/feedback-lessons.md`、`skills/_template/SKILL.md`、`skills/ADDING_SKILLS.md`、所有現有 skill 的入口（`skills/<name>/SKILL.md` 或 `workflow/<domain>/execution-flow.md`）與根 `README.md`；若某 skill 有 `execution-flow.md` 或 `artifact-gates.md` 的輸出格式，也需同步更新或明確檢查。 |
| `shared-rules/conversation-goal-ledger.md` | `shared-rules/README.md`、`shared-rules/content-layering.md`、`shared-rules/dependency-reading.md`、`scripts/README.md`、相關 helper script、`ai-tools/` 中各工具整合文件；若 tool-specific hook / rule 實作變更，也需同步對應工具規則或 hook 文件。 |
| `shared-rules/cross-skill-references.md` 或新增 cross-skill 關係 | referring skill 的入口（`skills/<name>/SKILL.md` 或 `workflow/<domain>/execution-flow.md`）、target skill 的入口或接收格式、必要時 `skills/_template/SKILL.md` 與 `skills/ADDING_SKILLS.md`。 |
| `shared-rules/feedback-lessons.md` | `feedback/history/README.md`、各 domain 的 `feedback/history/<domain>/README.md`、新增 lesson 模板。舊結構 `skills/<name>/FEEDBACK.md` 與 `skills/<name>/feedback_history/README.md` 已於 2026-05-13 刪除。 |
| 工具同步文件或同步腳本 | 根 `README.md`、`scripts/README.md`、`ai-tools/` 對應工具文件、Agents 必讀規則、實際執行同步。 |
| 新增或修改 `tools/adapters/<tool>.md`（新分層）或 `skills/<name>/tool-adapters/<tool>.md`（舊結構，僅限尚未遷移的 skill；已遷移 skill 的舊 tool-adapters 路徑已被刪除） | 該 skill 的入口（`skills/<name>/README.md` 或 `workflow/<domain>/README.md`）、核心 `execution-flow.md` 中的 adapter link、`ai-tools/<tool>.md` 的 cross-link（若該工具已有集中說明）、必要 validation/checklist。 |
| `scripts/ai-skill-close-loop.sh` | `scripts/README.md`、根 `README.md`、`shared-rules/dependency-reading.md`、本檔；若改變 lock / commit / push 條件，也需同步相關 skill close-loop 說明。 |
| 新增 skill | 根 `README.md`、`skills/README.md`、`skills-index.yaml`、必要 tool sync 實際同步結果。 |
| 修改 `skills/<name>/SKILL.md` 或 `workflow/<domain>/execution-flow.md` 觸發條件或流程 | 該 skill 的 `README.md`、`runtime/onboarding/` 對應 quickstart、`skills-index.yaml`、相關 cross-link。 |
| 新增 `feedback/history/` lesson | `feedback/history/<domain>/README.md`，必要時 `feedback/history/<domain>/<category>/README.md`。 |
| 修改 `analysis/development-guidance/controls-catalog.md`（原 `skills/app-development-guidance/controls/`） | 相關 `analysis/development-guidance/implementation-catalog.md`、`analysis/development-guidance/platforms-catalog.md`、`analysis/development-guidance/languages-catalog.md`、`workflow/software-delivery/development-process.md`。 |
| 修改 `workflow/software-delivery/development-process.md`（原 `skills/app-development-guidance/process/`） | 相關 `workflow/software-delivery/artifact-gates.md`、`analysis/development-guidance/risk-translation.md`、`analysis/development-guidance/implementation-catalog.md`、`workflow/software-delivery/execution-flow.md`。 |
| 修改 `analysis/development-guidance/implementation-catalog.md`（原 `skills/app-development-guidance/implementation/`） | 相關 `analysis/development-guidance/controls-catalog.md`、`analysis/development-guidance/platforms-catalog.md`、`analysis/development-guidance/languages-catalog.md`。 |
| 修改 `analysis/development-guidance/` 的 templates（原 `skills/app-development-guidance/templates/`） | `workflow/software-delivery/artifact-gates.md`、`workflow/software-delivery/development-process.md`、`workflow/software-delivery/execution-flow.md` 與引用該模板的文件。 |
| 修改 `workflow/software-delivery/development-process.md` 的 Product Brief / 企劃書驗證規則 | `workflow/software-delivery/development-process.md`、`workflow/software-delivery/execution-flow.md`、`workflow/software-delivery/artifact-gates.md`、`analysis/development-guidance/README.md`。 |
| 修改 `analysis/development-guidance/platforms-catalog.md` 或 `analysis/development-guidance/implementation-catalog.md` 的 embedded 相關內容 | `workflow/software-delivery/development-process.md`、`workflow/software-delivery/execution-flow.md`、`workflow/software-delivery/artifact-gates.md`。 |
| 修改 `analysis/development-guidance/implementation-catalog.md` 的 backend 相關內容 | `analysis/development-guidance/controls-catalog.md`、`workflow/software-delivery/development-process.md`、`workflow/software-delivery/execution-flow.md`。 |
| 修改 `analysis/development-guidance/` 的 tooling 相關內容 | `workflow/software-delivery/development-process.md`、`workflow/software-delivery/execution-flow.md`、`workflow/software-delivery/artifact-gates.md`。 |
| 修改 `workflow/software-delivery/development-process.md` 的 implemented-first governance / traceability / BDD closure 規則 | `workflow/software-delivery/development-process.md`、`workflow/software-delivery/execution-flow.md`、`workflow/software-delivery/artifact-gates.md`、`analysis/development-guidance/README.md`。 |
| 修改 `knowledge/` 或 `validation/` 下的檔案 | 執行 `ruby scripts/refresh-knowledge-runtime.rb` 確認所有 validator 通過。 |
| **Plan 完成閉環**（plan 所有項目標記為完成） | 依 [`plans/README.md`](../plans/README.md#plan-完成閉環plan-completion-closure) 執行 7 項檢查清單：確認所有項目完成、執行 validator、檢查連動更新、更新 `plans/README.md` 狀態、搬移至 `archived/`（或標註例外原因）、commit & push、最終確認。 |
| **架構重構**（目錄重組、分層新增、路徑變更、命名變更） | 依 [`governance/lifecycle/intelligence-extraction-pipeline.md`](../governance/lifecycle/intelligence-extraction-pipeline.md) **Step 6a** 建立 validation scenario（至少一個，測試 AI 是否正確使用新路徑）。依 **Step 7a** 檢查以下 shared-rules 檔案：`governance/document-sizing.md`（範例路徑）、`linked-updates.md`（本檔表格）、`feedback-lessons.md`（Promotion Target 模板）、`shared-rules/README.md`（lazy-load 表格）、`skills/ADDING_SKILLS.md`（目錄結構建議）、`content-layering.md`（路徑描述）、`tool-neutral-documentation.md`（路徑描述）、`decision-efficiency.md`（Context Loading 步驟）、`cross-skill-references.md`（引用格式）。同時檢查 `failure-patterns/README.md` 索引與 `failure-learning-system.md` 的 Promotion Decision 表格。 |

## 閉環不完整時的強制補救

若使用者或 review 指出「更新閉環不完整」，agent 必須先分析漏掉的原因，再補規則與連動更新。常見原因與補救：

| 漏洞原因 | 必須補強 |
| --- | --- |
| 只修改了單一 skill / 單一檔案，沒有判斷是否屬於全庫規則 | 依 [`reusable-guidance-boundary.md`](reusable-guidance-boundary.md) 與本檔判斷正確層級；若是全庫行為，補 `shared-rules/`。 |
| 新增 lesson 後沒有推廣到 promotion target 或 checklist | 回到 lesson 的 **Promotion Target** / **Required Linked Updates**，逐一更新或明確寫不適用理由。 |
| 寫完 reusable docs 後沒有搜尋專案特例殘留 | 搜尋 skill/shared docs 中的 project name、host、endpoint、payload、class/test 名稱、sample ID、本機路徑等，移回 project docs。 |
| 只描述「改了什麼」，沒有描述「為什麼會漏」與「如何防止重犯」 | 補原因分析、決策規則、驗證步驟，並依 [`goal-action-validation.md`](goal-action-validation.md) 寫明目標、執行、驗證。 |
| 使用者指出 agent 反覆出現同類失效模式，但只修當下文件 | 依 [`failure-learning-system.md`](failure-learning-system.md) 分類失效模式、建立或更新 `shared-rules/failure-patterns/`，並補 prevention gate 與 validation。 |
| 更新完 skill / shared rules 後沒有提醒 repo 仍有 pending commit / ahead / behind 狀態 | 最終回覆必須列出 `git status --short --branch` 的關鍵狀態、哪些是本輪提交、哪些是既有 dirty changes，以及需要 push / merge / 清理的下一步。 |
| 有多個 dirty owner group 卻混成單一 commit | 沒有 active lock 時使用 `scripts/ai-skill-close-loop.sh --commit` 或手動依 shared-rules、scripts、各 skill owner 分開 commit；有 active lock 時停止並提醒，不得更新。 |

## 語氣規則

連動更新不是建議事項。若相關，就使用 **必須**、**需要同步**、**已檢查** 這類明確語氣；不要使用「可選」、「有空再補」、「建議之後」來描述必要的連動工作。

## 沒有連動更新時

可以不改其他檔，但要能說明理由，例如：

- 只修正 typo，沒有改變流程或規則。
- 只新增單一 lesson，尚未成熟到推廣進主流程；已更新 `feedback/history/<domain>/README.md`。
- 某控制沒有平台或語言差異；已檢查相關資料夾，無需更新。

← [回到共用規則索引](README.md)
