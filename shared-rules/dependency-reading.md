# 依賴文件讀取鐵則

本規則適用於所有 agent 使用、修改或檢查 `shared-rules/`、`skills/`、`.cursor/rules/`、模板、feedback lessons、同步腳本與根索引時。目的不是增加形式流程，而是避免 agent 只讀單一文件，卻忽略已更新的依賴規則。

## 核心規則

只要發現某個 skill、shared rule、Cursor rule、模板或 feedback lesson 已更新、將被更新、或可能影響目前任務，agent 必須讀取它的相關依賴文件後才能下結論或繼續修改。

最低讀取範圍：

| 發現或修改的項目 | 必須讀取或明確檢查 |
| --- | --- |
| 任一 `skills/<name>/SKILL.md` | 該 skill 的 `README.md`、`WORKFLOW.md`、`DOCUMENTATION.md`、`CHECKLIST.md`、`FEEDBACK.md`、相關 `feedback_history/README.md`，以及 `shared-rules/README.md`。不存在的檔案可標記為不適用。 |
| 任一 skill 子文件 | 該 skill 的 `SKILL.md`、最近的目錄 `README.md`、相關 workflow/checklist/template、`shared-rules/linked-updates.md`。 |
| 任一 `shared-rules/*.md` | `shared-rules/README.md`、`shared-rules/content-layering.md`、`shared-rules/linked-updates.md`、`shared-rules/reusable-guidance-boundary.md`（若涉及 reusable guidance / incident / feedback）、受影響 skill 的 `SKILL.md` 或模板。 |
| 任一 `.cursor/rules/*.mdc` | 對應的 shared rule 正文、`shared-rules/README.md`、`shared-rules/cursor-sync.md`，以及受影響的 skill 入口。 |
| 任一 template | 模板目錄 `README.md`、引用該模板的 workflow/documentation/checklist、`shared-rules/linked-updates.md`。 |
| 任一 feedback lesson | 該分類 `README.md`、skill 的 `feedback_history/README.md`、`shared-rules/feedback-lessons.md`，以及 promotion target。 |

## Agent 行為

1. 先讀 `shared-rules/README.md`，再依任務讀相關 shared rule 全文。
2. 若任務碰到 skill，讀該 skill 入口與依賴文件；不要只依賴 `description` 或單一段落。
3. 若看到文件有 cross-link、promotion target、required linked updates、template reference、feedback index，或 reusable guidance / project incident 邊界，就循連結讀到任務所需的規則載入完成。
4. 若依賴文件不存在，記錄為 `not applicable`；若存在但未讀，不可宣稱已完成檢查。
5. 回覆或提交前，說明依賴讀取與連動更新的驗證方式。
6. 完成 `git commit`、`git push` 與必要的 bundle sync 後，必須重新讀取本次更新過的 skill/shared-rule 入口文件與主要依賴文件，確認目前 agent context 已載入最新版；不可只依賴提交前讀過的內容。
7. 最終回覆前必須執行 `git status --short --branch`。若 `Ai-skill` repo 仍有 modified/untracked/staged changes，或 branch 仍 ahead/behind remote，不得回覆「已完成」；必須先完成驗證、sync、commit、push、讀回，或明確說明被什麼阻塞。

## Ai-skill Writeback Transaction Guard

只要 agent 第一次寫入 `<AI_SKILL_REPO>/shared-rules/`、`<AI_SKILL_REPO>/skills/`、`.cursor/rules/`、模板、feedback lesson、同步腳本，或透過同步路徑（例如 `~/.cursor/shared-rules`、`~/.cursor/skills/<name>`）改到同一套內容，就必須立即把這件事視為一個尚未關閉的 **Ai-skill writeback transaction**，而不是等最終回覆才想起。

交易開始時必須：

1. 立刻建立或口頭標記一個明確的 close-loop 項目，例如 `Ai-skill close loop`。
2. 記錄 canonical repo 位置與本次 touched paths；若是從同步路徑讀寫，仍以 `<AI_SKILL_REPO>` 的 `git status --short --branch` 為準。
3. 在切回專案分析、長時間動態測試、或回覆「完成」前，先關閉這個 transaction。

交易關閉條件：

1. `git status --short --branch` 與 `git diff` 已檢查。
2. 必要的 linked updates 已同步或明確寫出不適用理由。
3. 必要的 bundle sync 已執行。
4. 相關檔案已 `git add`、`git commit`、`git push`。
5. Push 後已重新讀取更新過的入口、主要依賴、索引與 promotion target。
6. 最後一次 `git status --short --branch` 顯示 clean，且 branch 沒有 ahead/behind。

若 transaction 未關閉，agent 不得把注意力長時間切回業務專案，也不得把「已更新 skill」當作完成。可以繼續工作的唯一例外是：使用者明確要求暫停 Ai-skill close-loop；此時必須說明目前 dirty/ahead 狀態與下一步。

## Ai-skill 回寫完成門檻

只要 agent 在本庫回寫任何 `shared-rules/`、`skills/`、`.cursor/rules/`、模板、feedback lessons、README 或同步腳本，最終回覆前必須完成整個更新閉環：

1. `git status --short --branch` 檢查變更。
2. `git diff` 檢查將提交的內容，不得包含 secrets、raw tokens、私人 host、個資或本機絕對路徑。
3. 執行適用的 lints / Markdown link check / required linked updates 檢查。
4. 若影響 Cursor 可讀到的 skills/rules，執行 `./scripts/sync-cursor-bundle.sh`。
5. `git add` 相關檔案。
6. `git commit`。
7. `git push`。
8. Push 後重新讀取更新過的入口與主要依賴文件。
9. 再跑 `git status --short --branch`，必須看到沒有未提交變更，且 branch 不再 ahead/behind remote。

若第 9 步不乾淨，agent 必須回到第 1 步處理剩餘變更。不可在 dirty tree 或未 push 狀態下回覆「完成」。

## Commit / Push 後讀回 Gate

當本庫變更已完成 `git commit`、`git push`，且改動涉及 `shared-rules/`、`skills/`、`.cursor/rules/`、模板或 feedback lessons 時，agent 必須在最終回覆前做一次讀回：

| 更新類型 | Commit / push 後必須重新讀取 |
| --- | --- |
| `shared-rules/` | 更新過的 shared rule、`shared-rules/README.md`、`shared-rules/linked-updates.md`；若有 Cursor rule，也讀對應 `.cursor/rules/*.mdc`。 |
| `skills/<name>/` | 該 skill 的 `SKILL.md`，以及本次更新過的 workflow / documentation / checklist / template / feedback index。 |
| `.cursor/rules/` | 更新過的 `.mdc`，以及對應的 shared rule 正文。 |
| template 或 feedback lesson | 更新過的 template/lesson、索引 README、promotion target 或引用它的 workflow/documentation。 |

讀回目的：

- 確認提交後工作樹、bundle sync 與 agent context 沒有停在舊版本。
- 讓同一輪最終回覆能基於最新規則，而不是 commit 前暫存理解。
- 若下一個 agent 接手，最終回覆要明確說明已讀回哪些入口文件；若未能讀回，必須列為未完成驗證。

## 與連動更新的關係

本規則是「先讀依賴」；[`linked-updates.md`](linked-updates.md) 是「讀完後該同步更新或明確檢查哪些文件」。兩者都必須遵守：

- 沒有讀依賴，就不能可靠判斷是否需要連動更新。
- 已讀依賴但發現需要同步，就必須依 `linked-updates.md` 更新或說明無需更新的理由。
- 若改動會影響 Cursor 可讀到的 rules 或 skills，必須同步 bundle。

## 驗證

每次套用本規則時，至少要能回答：

| 欄位 | 必填內容 |
| --- | --- |
| 目標 | 這次要確認哪個 skill/rule/template 的依賴沒有漏讀。 |
| 執行 | 實際讀取或明確檢查了哪些依賴文件。 |
| 驗證 | `git diff`、Markdown link check、lints、required linked updates 檢查、bundle sync，或說明哪些文件不存在所以不適用。 |

← [回到共用規則索引](README.md)
