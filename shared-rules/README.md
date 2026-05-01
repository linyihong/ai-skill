# 共用規則（分類索引）

本目錄放**所有 skill 共用**的政策與約定，依主題分檔維護。**不要**在每一則 `feedback_history` lesson 裡重複貼上全文；條目頂部用相對路徑**連回此處**即可（模板與檔名規則見 [feedback-lessons.md](feedback-lessons.md)）。

## Agents（必讀）

1. **本檔是指索引，不是全文。**請先讀此 README，再依你將要做的行為，**打開並讀完下方表格中所有相關連結的全文**（不得只讀本頁摘要就當已遵守規則）。
2. **撰寫或新增 `feedback_history` lesson、或任何「回饋／沉澱技巧」行為前**，必須已讀過 **[feedback-lessons.md](feedback-lessons.md)** 與 **[sanitization.md](sanitization.md)**；若涉及授權邊界，另讀 **[authorization-scope.md](authorization-scope.md)**。
3. 索引與各分檔中的 **Markdown 連結**視為規則的一部分；請**循連結讀取**，直到該任務需要的條文都已載入為止。
4. **Git：**若你改動的是 **Ai-skill 這個 git repository**（`shared-rules/`、`skills/`、根目錄 README 等），**除非使用者明講不要提交**，否則應在該 repo 根目錄完成 **`git add` → `git commit` → `git push`**；請自行申請所需工具權限。重新載入 Cursor **不取代**這一步。

| 分類 | 檔案 | 內容摘要 |
| --- | --- | --- |
| 授權與範圍 | [authorization-scope.md](authorization-scope.md) | 僅在授權範圍內分析、資料邊界。 |
| 去敏與占位符 | [sanitization.md](sanitization.md) | 什麼不可寫進可重用文件、占位符約定。 |
| 內容分層 | [content-layering.md](content-layering.md) | 共用規則／技巧／業務專案各自放哪。 |
| Feedback 與技巧條目 | [feedback-lessons.md](feedback-lessons.md) | **檔名規則、模板、agent 行為、索引**（唯一正文）；各 skill 目錄僅保留 `feedback_history/` 與可選的極短 `FEEDBACK.md` 入口。 |
| 同步到 Cursor | [cursor-sync.md](cursor-sync.md) | 如何把 `shared-rules/` 與 `skills/` 佈署到 `.cursor`。 |

**單一真相來源：**只在本庫 **`shared-rules/`** 維護共用規則正文；佈署到專案時複製整個 `shared-rules/` 資料夾。
