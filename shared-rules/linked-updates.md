# 連動更新（全庫必須規則）

本檔是 **Ai-skill repository 全部文件與 skill 的連動更新規則**。若某項改動會影響其他文件、索引、同步流程、skill 入口、分類文件或範本，相關檔案 **必須**在同一次變更中更新或明確檢查，不可寫成「可選」或「之後再說」。

## Agent 必須做的事

1. 改任何 `shared-rules/`、`skills/`、根 `README.md`、同步腳本或模板前，先判斷是否有連動文件。
2. 若有連動文件，**必須**同步修改或明確寫下「已檢查，無需更新」的理由。
3. 若改動會影響 Cursor 可讀到的 skill 或 rules，**必須**執行 [`scripts/sync-cursor-bundle.sh`](../scripts/sync-cursor-bundle.sh)（除非使用者明講不要動本機 `~/.cursor`）。
4. 若改動 Ai-skill repo，除非使用者明講不要提交，**必須** `git add` → `commit` → `push`。
5. Commit/push 與必要的 bundle sync 完成後，**必須**依 [`dependency-reading.md`](dependency-reading.md) 重新讀取本次更新過的 skill/shared-rule 入口與主要依賴文件，確認 agent context 已載入最新版。
6. 回覆使用者時要說明已做了哪些連動更新；如果某些相關文件不用改，也要簡短說明原因。

## 常見連動關係

| 改動位置 | 必須同步更新或檢查 |
| --- | --- |
| `shared-rules/README.md` 或新增 shared rule | 根 `README.md`、相關 skill 的入口說明、`feedback_history` 模板引用。 |
| `shared-rules/dependency-reading.md` | `shared-rules/README.md`、`shared-rules/content-layering.md`、`shared-rules/linked-updates.md`、`.cursor/rules/` 的 always-apply agent rule、`skills/_template/SKILL.md`、`skills/ADDING_SKILLS.md`、所有現有 skill 的 `SKILL.md` 入口與根 `README.md`。 |
| `shared-rules/neutral-language.md` | `shared-rules/README.md`、`shared-rules/content-layering.md`、`shared-rules/feedback-lessons.md`、`skills/_template/SKILL.md`、`skills/ADDING_SKILLS.md`、所有現有 skill 的 `SKILL.md` 入口與根 `README.md`。 |
| `shared-rules/goal-action-validation.md` | `shared-rules/README.md`、`shared-rules/content-layering.md`、`shared-rules/feedback-lessons.md`、`skills/_template/SKILL.md`、`skills/ADDING_SKILLS.md`、所有現有 skill 的 `SKILL.md` 入口與根 `README.md`；若某 skill 有 `DOCUMENTATION.md` 或 `WORKFLOW.md` 的輸出格式，也需同步更新或明確檢查。 |
| `shared-rules/cross-skill-references.md` 或新增 cross-skill 關係 | referring skill 的 `SKILL.md` / `README.md` / `WORKFLOW.md` / `DOCUMENTATION.md`、target skill 的入口或接收格式、必要時 `skills/_template/SKILL.md` 與 `skills/ADDING_SKILLS.md`。 |
| `shared-rules/feedback-lessons.md` | 各 skill 的 `FEEDBACK.md`、`feedback_history/README.md`、新增 lesson 模板。 |
| `shared-rules/cursor-sync.md` 或 `scripts/sync-cursor-bundle.sh` | 根 `README.md`、`scripts/README.md`、Agents 必讀規則、實際執行同步。 |
| 新增 skill | 根 `README.md`、`skills/README.md`、`scripts/sync-cursor-bundle.sh` 實際同步結果。 |
| 修改 `skills/<name>/SKILL.md` 觸發條件或流程 | 該 skill 的 `README.md`、`RUNBOOK.md`、`WORKFLOW.md`、相關 cross-link。 |
| 新增 `feedback_history` lesson | 該 skill 的 `feedback_history/README.md`，必要時 `WORKFLOW.md`、`TOOLS.md`、`DOCUMENTATION.md` 或分類文件。 |
| 修改 `app-development-guidance/controls/` | 相關 `implementation/`、`platforms/`、`languages/`、`checklists/`。 |
| 修改 `app-development-guidance/process/` | 相關 `checklists/`、`templates/`、`implementation/`、`WORKFLOW.md`。 |
| 修改 `app-development-guidance/implementation/` | 相關 `controls/`、`platforms/`、`languages/`、`checklists/`。 |
| 修改 `app-development-guidance/templates/` | `templates/README.md`、`DOCUMENTATION.md`、`process/`、`CHECKLIST.md` 與引用該模板的文件。 |
| 修改 `app-development-guidance` 的 Product Brief / 企劃書驗證規則 | `process/README.md`、`WORKFLOW.md`、`CHECKLIST.md`、`templates/initial-development-docs.md`、`templates/README.md`、`DOCUMENTATION.md`、`SKILL.md`。 |
| 修改 `app-development-guidance/platforms/embedded/` 或 `implementation/embedded/` | `process/README.md`、`WORKFLOW.md`、`CHECKLIST.md`、`checklists/embedded-firmware-review.md`、`templates/initial-development-docs.md`。 |
| 修改 `app-development-guidance/implementation/backend/contract-codegen.md` 或 `vendor-integration.md` | `platforms/backend/api.md`、`checklists/api-security-review.md`、`process/README.md`、`WORKFLOW.md`、`CHECKLIST.md`、相關 `controls/`。 |
| 修改 `app-development-guidance/implementation/tooling/` | `process/README.md`、`WORKFLOW.md`、`CHECKLIST.md`、`checklists/contract-governance-review.md`、`templates/initial-development-docs.md`。 |
| 修改 `app-development-guidance` 的 implemented-first governance / traceability / BDD closure 規則 | `process/README.md`、`WORKFLOW.md`、`CHECKLIST.md`、`checklists/contract-governance-review.md`、`templates/initial-development-docs.md`、`SKILL.md`。 |

## 語氣規則

連動更新不是建議事項。若相關，就使用 **必須**、**需要同步**、**已檢查** 這類明確語氣；不要使用「可選」、「有空再補」、「建議之後」來描述必要的連動工作。

## 沒有連動更新時

可以不改其他檔，但要能說明理由，例如：

- 只修正 typo，沒有改變流程或規則。
- 只新增單一 lesson，尚未成熟到推廣進主流程；已更新 `feedback_history/README.md`。
- 某控制沒有平台或語言差異；已檢查相關資料夾，無需更新。

← [回到共用規則索引](README.md)
