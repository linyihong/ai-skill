# 共用規則（分類索引）

本目錄放**所有 skill 共用**的政策與約定，依主題分檔維護。**不要**在每一則 `feedback_history` lesson 裡重複貼上全文；條目頂部用相對路徑**連回此處**即可（模板與檔名規則見 [feedback-lessons.md](feedback-lessons.md)）。

## Agents（必讀）

1. **本檔是指索引，不是全文。**請先讀此 README，再依你將要做的行為，**打開並讀完下方表格中所有相關連結的全文**（不得只讀本頁摘要就當已遵守規則）。
2. **依賴讀取鐵則：**只要發現某個 skill、shared rule、tool-specific rule、模板或 feedback lesson 已更新、將被更新、或會影響目前任務，必須依 **[dependency-reading.md](dependency-reading.md)** 讀取相關依賴文件；包含 `shared-rules/` 本身，不可只讀單一 skill 檔。若使用者要求「重新讀 skill」或指出 shared rules 可能漏讀，必須先建立 dependency read ledger，列出已讀、缺檔不適用與仍阻塞項目。
3. **撰寫或新增 `feedback_history` lesson、或任何「回饋／沉澱技巧」行為前**，必須已讀過 **[feedback-lessons.md](feedback-lessons.md)**、**[reusable-guidance-boundary.md](reusable-guidance-boundary.md)**、**[sanitization.md](sanitization.md)**、**[neutral-language.md](neutral-language.md)** 與 **[goal-action-validation.md](goal-action-validation.md)**；若涉及授權邊界，另讀 **[authorization-scope.md](authorization-scope.md)**。
4. 索引與各分檔中的 **Markdown 連結**視為規則的一部分；請**循連結讀取**，直到該任務需要的條文都已載入為止。
5. **Git：**若你改動的是 **Ai-skill 這個 git repository**（`shared-rules/`、`skills/`、根目錄 README 等），**除非使用者明講不要提交**，否則**必須**在該 repo 根目錄完成 **`git add` → `git commit` → `git push`**；請自行申請所需工具權限。重新載入或重啟任一工具**不取代**這一步。第一次寫入 Ai-skill 或其工具同步路徑時，必須依 [dependency-reading.md](dependency-reading.md) 立刻開啟 writeback transaction；在切回專案分析或回覆完成前關閉 transaction。Commit/push 與必要 tool sync 完成後，還必須重新讀取本次更新過的 skill/shared-rule 入口與主要依賴文件，並用 `git status --short --branch` 確認工作樹乾淨且沒有 ahead/behind。若使用者未授權 push / merge 而 repo 仍有 pending commit、ahead/behind 或未合併狀態，最終回覆必須主動提醒使用者。
6. **工具同步：**若本機使用任何工具特定的 skill/rule 同步或 mirror，每次改動 **`shared-rules/`** 或 **`skills/`** 後**必須**執行已設定的同步流程；具體工具路徑與命令放在 [`ai-tools/`](../ai-tools/README.md)。**Agent** 在有權執行 shell 時必須代為執行已設定同步，除非使用者表示不要動本機工具設定。
7. **連動更新：**任何會影響其他文件、索引、skill 入口、同步流程或分類文件的改動，都**必須**依 [linked-updates.md](linked-updates.md) 同步更新或明確檢查；不得把必要連動說成「可選」。

| 分類 | 檔案 | 內容摘要 |
| --- | --- | --- |
| 授權與範圍 | [authorization-scope.md](authorization-scope.md) | 僅在授權範圍內分析、資料邊界。 |
| 去敏與占位符 | [sanitization.md](sanitization.md) | 什麼不可寫進可重用文件、占位符約定。 |
| 中性與低爭議用語 | [neutral-language.md](neutral-language.md) | 文件標題、檔名、摘要與正文避免高風險或爭議詞；用授權、合規、契約與風險控制語境描述。 |
| 工具中立文件 | [tool-neutral-documentation.md](tool-neutral-documentation.md) | 可重用文件預設保持工具中立；工具專屬路徑、hook、同步與 UI 步驟放到 `ai-tools/` 或工具設定檔。 |
| 文件 TODO | [document-todo-list.md](document-todo-list.md) | 文件有未完成、待決策、待補強或待驗證內容時，在前段放可掃描 TODO 表並連到相關章節、goal 或 issue。 |
| 目標、執行、驗證 | [goal-action-validation.md](goal-action-validation.md) | 每個重要工作單元要能反查目標、執行內容與驗證方式；純判斷題用參考來源與推論邊界。 |
| 對話目標閉環 | [conversation-goal-ledger.md](conversation-goal-ledger.md) | 使用專案本地 `.agent-goals/` 暫存帳本追蹤 active goals、優先權、拆解、轉移、multi-agent lock、完成驗證與刪除條件；多步驟、已有 todo、使用者要求繼續或看到 dirty files 時要先 status/init。 |
| 依賴文件讀取 | [dependency-reading.md](dependency-reading.md) | 發現 skill/shared rule/tool-specific rule/template/lesson 更新時，必須讀相關依賴文件，包含 shared-rules，並用 dependency read ledger 防止漏讀。 |
| 內容分層 | [content-layering.md](content-layering.md) | 共用規則／技巧／業務專案各自放哪。 |
| 可重用規則與專案證據邊界 | [reusable-guidance-boundary.md](reusable-guidance-boundary.md) | 技巧只沉澱通用原因、規則與驗證；專案 incident 證據留在專案文件；指出閉環不完整時必須分析原因並強化規則。 |
| 文件大小與拆分 | [document-sizing.md](document-sizing.md) | 文件變大時改成目錄、分類資料夾與多檔，避免單檔堆疊。 |
| Cross-skill references | [cross-skill-references.md](cross-skill-references.md) | 一個 skill 需要引用另一個 skill 的規範、模板、交接產物或驗證流程時怎麼寫。 |
| 連動更新 | [linked-updates.md](linked-updates.md) | 全庫必須連動更新規則：改一處影響多處時，相關文件必須同步更新或明確檢查。 |
| Feedback 與技巧條目 | [feedback-lessons.md](feedback-lessons.md) | **檔名規則、模板、agent 行為、索引**（唯一正文）；各 skill 目錄僅保留 `feedback_history/` 與可選的極短 `FEEDBACK.md` 入口。 |
| 工具同步 | [tool-neutral-documentation.md](tool-neutral-documentation.md) / [`ai-tools`](../ai-tools/README.md) | 通用規則保持工具中立；具體工具同步、hook、UI 與路徑放在 `ai-tools/` 或工具設定檔。 |

**單一真相來源：**只在本庫 **`shared-rules/`** 維護共用規則正文；佈署到專案時複製整個 `shared-rules/` 資料夾。
