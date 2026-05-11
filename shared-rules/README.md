# 共用規則（分類索引）

本目錄放**所有 skill 共用**的政策與約定，依主題分檔維護。**不要**在每一則 `feedback_history` lesson 裡重複貼上全文；條目頂部用相對路徑**連回此處**即可（模板與檔名規則見 [feedback-lessons.md](feedback-lessons.md)）。

## Agents（必讀）

1. **本檔是指索引，不是全文。**請先讀此 README，再依你將要做的行為，**打開並讀完下方表格中所有相關連結的全文**（不得只讀本頁摘要就當已遵守規則）。
2. **預設 bootstrap：**每次開啟 agent / 新 session / 接手長對話時，先讀本索引與下方 **[Default Bootstrap](#default-bootstrap)**；再依任務循連結讀取相關全文。不要一次把所有 shared rules 當作已讀。若規則看似衝突，先依 [rule-weight.md](rule-weight.md) 判斷權重。
3. **依賴讀取鐵則：**只要發現某個 skill、shared rule、tool-specific rule、模板或 feedback lesson 已更新、將被更新、或會影響目前任務，必須依 **[dependency-reading.md](dependency-reading.md)** 讀取相關依賴文件；包含 `shared-rules/` 本身，不可只讀單一 skill 檔。若使用者要求「重新讀 skill」或指出 shared rules 可能漏讀，必須先建立 dependency read ledger，列出已讀、缺檔不適用與仍阻塞項目。
4. **撰寫或新增 `feedback_history` lesson、或任何「回饋／沉澱技巧」行為前**，必須已讀過 **[feedback-lessons.md](feedback-lessons.md)**、**[reusable-guidance-boundary.md](reusable-guidance-boundary.md)**、**[sanitization.md](sanitization.md)**、**[neutral-language.md](neutral-language.md)** 與 **[goal-action-validation.md](goal-action-validation.md)**；若涉及授權邊界，另讀 **[authorization-scope.md](authorization-scope.md)**。
5. 索引與各分檔中的 **Markdown 連結**視為規則的一部分；請**循連結讀取**，直到該任務需要的條文都已載入為止。
6. **Git：**若你改動的是 **Ai-skill 這個 git repository**（`shared-rules/`、`skills/`、根目錄 README 等），**除非使用者明講不要提交**，否則**必須**在該 repo 根目錄完成 **`git add` → `git commit` → `git push`**；請自行申請所需工具權限。重新載入或重啟任一工具**不取代**這一步。第一次寫入 Ai-skill 或其工具同步路徑時，必須依 [dependency-reading.md](dependency-reading.md) 先定位 canonical `<AI_SKILL_REPO>` git repo，再開啟 writeback transaction；`~/.cursor/skills*`、bundles、專案 `.cursor/` 等工具部署路徑不可當成 source repo 完成回寫。在切回專案分析或回覆完成前關閉 transaction。Commit/push 與必要 tool sync 完成後，還必須重新讀取本次更新過的 skill/shared-rule 入口與主要依賴文件，並用 `git status --short --branch` 確認工作樹乾淨且沒有 ahead/behind。若使用者未授權 push / merge 而 repo 仍有 pending commit、ahead/behind 或未合併狀態，最終回覆必須主動提醒使用者。
7. **工具同步：**預設 reference-first 時不需要同步工具 mirror。只有本機明確使用工具特定的 symlink / bundle / copy mirror，且希望該 mirror 立即跟上時，才執行對應同步流程；具體工具路徑與命令放在 [`ai-tools/`](../ai-tools/README.md)。**Agent** 不得因為改了 `shared-rules/` 或 `skills/` 就自動跑工具同步，除非使用者要求、環境變數啟用，或目前任務就是維護該 mirror。
8. **連動更新：**任何會影響其他文件、索引、skill 入口、同步流程或分類文件的改動，都**必須**依 [linked-updates.md](linked-updates.md) 同步更新或明確檢查；不得把必要連動說成「可選」。

## Default Bootstrap

每次 agent / session 啟動、接手長對話、或使用者要求「先讀規則」時，預設載入下列最小集合：

| 順序 | 檔案 | 用途 |
| --- | --- | --- |
| 1 | [README.md](README.md) | shared rules 索引與讀取路由。 |
| 2 | [dependency-reading.md](dependency-reading.md) | 依賴讀取、dependency read ledger、Ai-skill writeback transaction。 |
| 3 | [linked-updates.md](linked-updates.md) | 判斷改一處時哪些文件必須同步。 |
| 4 | [conversation-goal-ledger.md](conversation-goal-ledger.md) | 對話目標閉環、owner/lock、parallelization mode、完成刪除。 |
| 5 | [tool-neutral-documentation.md](tool-neutral-documentation.md) | 可重用文件保持工具中立，工具差異放 adapter/tool docs。 |
| 6 | [rule-weight.md](rule-weight.md) | 規則權重與衝突優先序：安全/source/validation/user-goal 高於工具相容與效率偏好。 |
| 7 | [decision-efficiency.md](decision-efficiency.md) | 以最高收益路線做決策、降低無關讀取與 token/context 成本。 |
| 8 | [failure-learning-system.md](failure-learning-system.md) | 將 agent 失效模式轉成可累積的 failure pattern、防呆規則與驗證 gate。 |
| 9 | [document-todo-list.md](document-todo-list.md) | 文件前段 TODO 表與 goal/章節連結。 |
| 10 | [document-sizing.md](document-sizing.md) | 文件大小、拆分時機、資料夾與索引結構。 |
| 11 | [goal-action-validation.md](goal-action-validation.md) | 重要工作單元的目標、執行、驗證。 |
| 12 | [neutral-language.md](neutral-language.md) | 文件標題、摘要與正文的中性低爭議用語。 |

Bootstrap 不是「全部規則已讀」。完成 bootstrap 後，仍必須依任務讀相關全文，例如 feedback、sanitization、authorization、cross-skill references、skill-specific workflow / tools / documentation。

| 分類 | 檔案 | 內容摘要 |
| --- | --- | --- |
| 授權與範圍 | [authorization-scope.md](authorization-scope.md) | 僅在授權範圍內分析、資料邊界。 |
| 去敏與占位符 | [sanitization.md](sanitization.md) | 什麼不可寫進可重用文件、占位符約定。 |
| 中性與低爭議用語 | [neutral-language.md](neutral-language.md) | 文件標題、檔名、摘要與正文避免高風險或爭議詞；用授權、合規、契約與風險控制語境描述。 |
| 工具中立文件 | [tool-neutral-documentation.md](tool-neutral-documentation.md) | 可重用文件預設保持工具中立；工具專屬路徑、hook、同步與 UI 步驟放到 `ai-tools/` 或工具設定檔。 |
| 規則權重與衝突優先序 | [rule-weight.md](rule-weight.md) | 當 shared rules、skill workflow、tool adapter、使用者目標或效率規則看似衝突時，依安全/source/validation/user-goal/tool adapter/效率的權重排序處理。 |
| 決策效率 | [decision-efficiency.md](decision-efficiency.md) | 先界定未知，再依 time-to-evidence、語意距離、風險、驗證信號與 context 成本選下一步。 |
| 失效學習系統 | [failure-learning-system.md](failure-learning-system.md) / [failure-patterns](failure-patterns/README.md) | 使用 failure taxonomy、pattern records、promotion targets 與 validation gate，將重複 agent 失效模式沉澱成可重用防呆規則；例如 source/mirror 寫入漂移。 |
| 文件 TODO | [document-todo-list.md](document-todo-list.md) | 文件有未完成、待決策、待補強或待驗證內容時，在前段放可掃描 TODO 表並連到相關章節、goal 或 issue。 |
| 目標、執行、驗證 | [goal-action-validation.md](goal-action-validation.md) | 每個重要工作單元要能反查目標、執行內容與驗證方式；純判斷題用參考來源與推論邊界。 |
| 對話目標閉環 | [conversation-goal-ledger.md](conversation-goal-ledger.md) | 使用專案本地 `.agent-goals/` 暫存帳本追蹤 active goals、優先權、parallelization mode、owner/lock 決策、plan/todo links、missing/decision/strengthen、拆解、轉移、完成驗證與刪除條件；多步驟、已有 todo、使用者要求繼續或看到 dirty files 時要先 status/init。 |
| 依賴文件讀取 | [dependency-reading.md](dependency-reading.md) | 發現 skill/shared rule/tool-specific rule/template/lesson 更新時，必須讀相關依賴文件，包含 shared-rules，並用 dependency read ledger 防止漏讀。 |
| 內容分層 | [content-layering.md](content-layering.md) | 共用規則／技巧／業務專案各自放哪。 |
| 可重用規則與專案證據邊界 | [reusable-guidance-boundary.md](reusable-guidance-boundary.md) | 技巧只沉澱通用原因、規則與驗證；專案 incident 證據留在專案文件；指出閉環不完整時必須分析原因並強化規則。 |
| 文件大小與拆分 | [document-sizing.md](document-sizing.md) | 文件變大時改成目錄、分類資料夾與多檔，避免單檔堆疊。 |
| Cross-skill references | [cross-skill-references.md](cross-skill-references.md) | 一個 skill 需要引用另一個 skill 的規範、模板、交接產物或驗證流程時怎麼寫。 |
| 連動更新 | [linked-updates.md](linked-updates.md) | 全庫必須連動更新規則：改一處影響多處時，相關文件必須同步更新或明確檢查。 |
| Feedback 與技巧條目 | [feedback-lessons.md](feedback-lessons.md) | **檔名規則、模板、agent 行為、索引**（唯一正文）；各 skill 目錄僅保留 `feedback_history/` 與可選的極短 `FEEDBACK.md` 入口。 |
| 工具同步 / 參照 | [tool-neutral-documentation.md](tool-neutral-documentation.md) / [`ai-tools`](../ai-tools/README.md) | 通用規則保持工具中立；具體工具參照、symlink、同步、hook、UI 與路徑放在 `ai-tools/` 或工具設定檔。 |

**單一真相來源：**只在本庫 **`shared-rules/`** 維護共用規則正文；部署到工具或專案時優先參照中央庫或 symlink，需要離線快照時才複製整個 `shared-rules/` 資料夾。
