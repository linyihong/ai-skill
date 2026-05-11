# 內容分層

| 內容 | 放哪裡 |
| --- | --- |
| Repo-level **架構方向、roadmap、遷移階段與移除條件** | **`architecture/`**；例如 [`../architecture/ai-native-knowledge-operating-system.md`](../architecture/ai-native-knowledge-operating-system.md)。這類文件不是 executable shared rule，不放 `shared-rules/` 正文。 |
| 全 skill **共用政策**（授權、去敏、路徑、feedback 原則） | **`shared-rules/`** 各分類檔（本目錄） |
| 全庫**連動更新規則**（改一處會影響多處時要同步改哪些） | **[`linked-updates.md`](linked-updates.md)**（全庫唯一正文） |
| 全庫**文件用語規則**（中性、低爭議、避免 AI/搜尋誤判的標題、slug、摘要與正文） | **[`neutral-language.md`](neutral-language.md)**（全庫唯一正文） |
| 全庫**工具中立文件規則**（可重用 docs 不綁單一 IDE/agent；工具路徑、hook、UI、同步步驟放工具文件） | **[`tool-neutral-documentation.md`](tool-neutral-documentation.md)**（全庫唯一正文）；具體工具操作放 `ai-tools/<tool>.md` 或工具設定檔 |
| 全庫**決策效率與 context 控制規則**（先界定未知、用最高收益路線決策、避免無關讀取與 token 浪費） | **[`decision-efficiency.md`](decision-efficiency.md)**（全庫唯一正文）；skill 可保留領域專用路由表並連回此規則 |
| 全庫**失效學習系統**（將 agent 失效分類、記錄、推廣成防呆規則與 validation gate） | **[`failure-learning-system.md`](failure-learning-system.md)**（流程正文）；可重用跨 skill 失效模式放 **[`failure-patterns/`](failure-patterns/README.md)** |
| 全庫**文件 TODO 規則**（文件前段列出未完成、待決策、待補強、待驗證項並連回章節/goal） | **[`document-todo-list.md`](document-todo-list.md)**（全庫唯一正文）；跨文件或跨對話目標再連到 `.agent-goals/` |
| 全庫**工作驗證流程**（目標、執行、驗證；純判斷題用參考來源） | **[`goal-action-validation.md`](goal-action-validation.md)**（全庫唯一正文） |
| 全工具**對話目標閉環規則**（active goals、優先權、parallelization mode、owner/lock 決策、plan/todo links、missing/decision/strengthen、拆解、轉移、multi-agent lock、完成後刪除；多步驟 / todo / dirty files / 繼續前任務時先 status/init） | **[`conversation-goal-ledger.md`](conversation-goal-ledger.md)**（全庫唯一正文）；每個專案的暫存狀態放 `<PROJECT_ROOT>/.agent-goals/`，不進 git |
| 全庫**依賴讀取鐵則、dependency read ledger 與 writeback transaction gate**（發現 skill/rule/template/lesson 更新時必須讀相關依賴，重讀 skill 時列出已讀/不適用/阻塞項，寫入 Ai-skill 時必須先定位 canonical repo，不把工具 mirror 當 source，並完成 sync/commit/push/readback/clean status） | **[`dependency-reading.md`](dependency-reading.md)**（全庫唯一正文） |
| 全庫**可重用規則與專案證據邊界**（incident 只能抽象成通用原因、規則與驗證；具體證據留專案） | **[`reusable-guidance-boundary.md`](reusable-guidance-boundary.md)**（全庫唯一正文） |
| 可重用的**單一技巧、lesson 全文** | 未分類 skill 放 `skills/<skill-name>/feedback_history/`；已分類 skill 放 `skills/<skill-name>/feedback_history/<category>/` 或 `common/`；成熟後可整理進該 skill 的 `WORKFLOW.md` / `TOOLS.md` / `DOCUMENTATION.md` / 分類資料夾 |
| 某 skill 的**工具策略 adapter**（同一 skill 在某 AI 工具上的執行差異） | `skills/<skill-name>/tool-adapters/<tool>.md`；只寫該工具差異並連回核心 workflow，工具全域設定仍放 `ai-tools/<tool>.md` |
| **如何**下筆、命名、模板 | **[`feedback-lessons.md`](feedback-lessons.md)**（全庫唯一）；各 skill 可選保留一行入口 `FEEDBACK.md` 指向該檔 |
| **某 skill 需要引用另一個 skill 的規範或交接產物** | 在 referring skill 寫短 cross-link 與觸發條件；規則見 **[`cross-skill-references.md`](cross-skill-references.md)**；不要複製 target skill 的全文 |
| 特定 App 的 host、一次性結論、專屬實作、live run 結果、class/test 名稱、sample ID | **業務專案**文件，不進可重用 skill；若要沉澱，依 [`reusable-guidance-boundary.md`](reusable-guidance-boundary.md) 先抽象化 |

## 文件變大時

若任何 skill、技巧分類、shared rule、模板或寫作規範開始變大，不要把所有內容繼續堆在同一個 Markdown 檔。依 [`document-sizing.md`](document-sizing.md) 拆成「資料夾 + `README.md` 目錄 + 多個子檔」，讓 agent 只讀與任務相關的子文件。

← [回到共用規則索引](README.md)
