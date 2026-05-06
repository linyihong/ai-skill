# 依賴文件讀取鐵則

本規則適用於所有 agent 使用、修改或檢查 `shared-rules/`、`skills/`、`.cursor/rules/`、模板、feedback lessons、同步腳本與根索引時。目的不是增加形式流程，而是避免 agent 只讀單一文件，卻忽略已更新的依賴規則。

## 核心規則

只要發現某個 skill、shared rule、Cursor rule、模板或 feedback lesson 已更新、將被更新、或可能影響目前任務，agent 必須讀取它的相關依賴文件後才能下結論或繼續修改。

最低讀取範圍：

| 發現或修改的項目 | 必須讀取或明確檢查 |
| --- | --- |
| 任一 `skills/<name>/SKILL.md` | 該 skill 的 `README.md`、`WORKFLOW.md`、`DOCUMENTATION.md`、`CHECKLIST.md`、`FEEDBACK.md`、相關 `feedback_history/README.md`，以及 `shared-rules/README.md`。不存在的檔案可標記為不適用。 |
| 任一 skill 子文件 | 該 skill 的 `SKILL.md`、最近的目錄 `README.md`、相關 workflow/checklist/template、`shared-rules/linked-updates.md`。 |
| 任一 `shared-rules/*.md` | `shared-rules/README.md`、`shared-rules/content-layering.md`、`shared-rules/linked-updates.md`、受影響 skill 的 `SKILL.md` 或模板。 |
| 任一 `.cursor/rules/*.mdc` | 對應的 shared rule 正文、`shared-rules/README.md`、`shared-rules/cursor-sync.md`，以及受影響的 skill 入口。 |
| 任一 template | 模板目錄 `README.md`、引用該模板的 workflow/documentation/checklist、`shared-rules/linked-updates.md`。 |
| 任一 feedback lesson | 該分類 `README.md`、skill 的 `feedback_history/README.md`、`shared-rules/feedback-lessons.md`，以及 promotion target。 |

## Agent 行為

1. 先讀 `shared-rules/README.md`，再依任務讀相關 shared rule 全文。
2. 若任務碰到 skill，讀該 skill 入口與依賴文件；不要只依賴 `description` 或單一段落。
3. 若看到文件有 cross-link、promotion target、required linked updates、template reference、feedback index，就循連結讀到任務所需的規則載入完成。
4. 若依賴文件不存在，記錄為 `not applicable`；若存在但未讀，不可宣稱已完成檢查。
5. 回覆或提交前，說明依賴讀取與連動更新的驗證方式。

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
