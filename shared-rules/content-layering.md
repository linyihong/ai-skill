# 內容分層

| 內容 | 放哪裡 |
| --- | --- |
| Repo-level **架構方向、roadmap、遷移階段與移除條件** | **`architecture/`**；例如 [`../architecture/ai-native-knowledge-operating-system.md`](../architecture/ai-native-knowledge-operating-system.md)。這類文件不是 executable shared rule，不放 `shared-rules/` 正文。 |
| **長期目標、phase、未完成能力、migration 狀態、promotion / deprecation 決策** | **Durable planning 文件**：repo-level 放 `architecture/` 或相關 layer README；知識生命週期放 `governance/`；routing / atom 方向放 `knowledge/`、`metadata/`；業務專案放正式 project docs / issue tracker。**不要**用 `.agent-goals/` 保存長期目標或 completed archive。 |
| **本輪可中斷、需接手的 active conversation goal / implementation task** | **`<PROJECT_ROOT>/.agent-goals/`**；只追蹤目前對話的可恢復工作。完成驗證後刪除；若仍有長期方向，刪除前先回寫 durable planning 文件。 |
| 全 skill **共用政策**（授權、去敏、路徑、feedback 原則） | **`shared-rules/`** 各分類檔（本目錄） |
| 全庫**連動更新規則**（改一處會影響多處時要同步改哪些） | **[`linked-updates.md`](linked-updates.md)**（全庫唯一正文） |
| 全庫**文件用語與語言一致性規則**（正文預設繁體中文；中性、低爭議、避免 AI/搜尋誤判的標題、slug、摘要與正文） | **[`neutral-language.md`](neutral-language.md)**（全庫唯一正文） |
| 全庫**工具中立文件規則**（可重用 docs 不綁單一 IDE/agent；工具路徑、hook、UI、同步步驟放工具文件） | **[`tool-neutral-documentation.md`](tool-neutral-documentation.md)**（全庫唯一正文）；具體工具操作放 `ai-tools/<tool>.md` 或工具設定檔 |
| 全庫**規則權重與衝突優先序**（安全/source/validation/user-goal/tool adapter/效率規則彼此衝突時如何排序） | **[`rule-weight.md`](rule-weight.md)**（全庫唯一正文） |
| 全庫**決策效率與 context 控制規則**（先界定未知、用最高收益路線決策、避免無關讀取與 token 浪費） | **[`decision-efficiency.md`](decision-efficiency.md)**（全庫唯一正文）；skill 可保留領域專用路由表並連回此規則 |
| 全庫**失效學習系統**（將 agent 失效分類、記錄、推廣成防呆規則與 validation gate） | **[`failure-learning-system.md`](failure-learning-system.md)**（流程正文）；可重用跨 skill 失效模式放 **[`failure-patterns/`](failure-patterns/README.md)** |
| 全庫**文件 TODO 規則**（文件前段列出未完成、待決策、待補強、待驗證項並連回章節/goal） | **[`document-todo-list.md`](document-todo-list.md)**（全庫唯一正文）；跨文件或跨對話目標再連到 `.agent-goals/` |
| 全庫**工作驗證流程**（目標、執行、驗證；純判斷題用參考來源） | **[`goal-action-validation.md`](goal-action-validation.md)**（全庫唯一正文） |
| 全工具**對話目標閉環規則**（active goals、優先權、parallelization mode、owner/lock 決策、plan/todo links、missing/decision/strengthen、拆解、轉移、multi-agent lock、完成後刪除；多步驟 / todo / dirty files / 繼續前任務時先 status/init；長期目標需落到 durable planning 文件） | **[`conversation-goal-ledger.md`](conversation-goal-ledger.md)**（全庫唯一正文）；每個專案的暫存狀態放 `<PROJECT_ROOT>/.agent-goals/`，不進 git |
| 全庫**依賴讀取鐵則、dependency read ledger 與 writeback transaction gate**（發現 skill/rule/template/lesson 更新時必須讀相關依賴，重讀 skill 時列出已讀/不適用/阻塞項，寫入 Ai-skill 時必須先定位 canonical repo，不把工具 mirror 當 source，並完成 sync/commit/push/readback/clean status） | **[`dependency-reading.md`](dependency-reading.md)**（全庫唯一正文） |
| 全庫**可重用規則與專案證據邊界**（incident 只能抽象成通用原因、規則與驗證；具體證據留專案） | **[`reusable-guidance-boundary.md`](reusable-guidance-boundary.md)**（全庫唯一正文） |
| 可重用的**單一技巧、lesson 全文** | **`feedback/history/<domain>/`**（統一目標路徑）；舊結構 `skills/<skill-name>/feedback_history/` 已於 2026-05-13 刪除，所有 lesson 已搬遷至 `feedback/history/<domain>/`。成熟 lesson 可 promotion 到 `intelligence/<domain>/` 或 `workflow/<domain>/` |
| **intelligence 內部：entry vs solution 分層** | 見下方 [Intelligence Entry/Solution 分層](#intelligence-entrysolution-分層) |
| 某 skill 的**工具策略 adapter**（同一 skill 在某 AI 工具上的執行差異） | 新分層：`tools/adapters/<tool>.md`；舊結構保留：`skills/<skill-name>/tool-adapters/<tool>.md`（向後相容，僅限尚未遷移的 skill；已遷移 skill 的舊 tool-adapters 路徑已被刪除）。只寫該工具差異並連回核心 workflow，工具全域設定仍放 `ai-tools/<tool>.md` |
| **如何**下筆、命名、模板 | **[`feedback-lessons.md`](feedback-lessons.md)**（全庫唯一）；各 skill 可選保留一行入口 `FEEDBACK.md` 指向該檔 |
| **某 skill 需要引用另一個 skill 的規範或交接產物** | 在 referring skill 寫短 cross-link 與觸發條件；規則見 **[`cross-skill-references.md`](cross-skill-references.md)**；不要複製 target skill 的全文 |
| 特定 App 的 host、一次性結論、專屬實作、live run 結果、class/test 名稱、sample ID | **業務專案**文件，不進可重用 skill；若要沉澱，依 [`reusable-guidance-boundary.md`](reusable-guidance-boundary.md) 先抽象化 |

## 文件變大時

若任何 skill、技巧分類、shared rule、模板或寫作規範開始變大，不要把所有內容繼續堆在同一個 Markdown 檔。依 [`document-sizing.md`](document-sizing.md) 拆成「資料夾 + `README.md` 目錄 + 多個子檔」，讓 agent 只讀與任務相關的子文件。

### 跨專案適用說明

上述分層結構**不限於 AI agent 知識庫**。任何專案文件都可套用相同原則：

| 專案類型 | 範例結構 |
|----------|----------|
| API 規格書 | `api/README.md` + `api/authentication.md` + `api/endpoints/` |
| 系統設計文件 | `design/README.md` + `design/architecture.md` + `design/data-model.md` |
| 專案 wiki | `wiki/README.md` + `wiki/onboarding.md` + `wiki/development-setup.md` |
| 技術文件集 | `docs/README.md` + `docs/guides/` + `docs/references/` + `docs/tutorials/` |

## Intelligence Entry/Solution 分層

當 `intelligence/<domain>/` 下的知識點同時涉及「診斷/偵測」與「解法/預防」時，應採用 **entry/solution 分層**，避免同一知識點在多處重複：

### 原則

1. **一個知識點只放一份解法**，但可以有多個入口指向它。
2. **入口（Entry）**：症狀驅動、問題驅動的分類。當 agent 遇到具體錯誤或現象時，從這裡切入。例如：
   - `failure/` — 「我出了什麼錯？」（錯誤症狀、診斷方式）
   - `signals/` — 「我看到什麼現象？」（觀察到的信號）
   - `anti-patterns/` — 「我是不是在走彎路？」（可預防的錯誤模式）
3. **解法（Solution）**：規則驅動、預防驅動的分類。當 agent 需要知道「正確做法是什麼」時，從這裡找。例如：
   - `heuristics/` — 「我應該怎麼做？」（啟發式規則、決策表）
   - `patterns/` — 「正確的模式是什麼？」（設計模式、最佳實踐）
4. **入口檔案**只放症狀描述與診斷方式，解法直接指向 solution 分類的對應檔案，不重複解法內容。
5. **解法檔案**頂部可反向引用入口（「如果遇到症狀，先看 entry 確認」），但本體只放預防規則與驗證方式。

### 範例

```
intelligence/engineering/analysis/
  failure/
    javascript-bitwise-64bit-truncation.md    # 入口：症狀（undecoded、access violation）+ 診斷方式
                                                # → 解法見 heuristics/javascript-bitwise-64bit-truncation.md
  heuristics/
    javascript-bitwise-64bit-truncation.md    # 解法：預防規則、決策表、驗證碼
                                                # → 頂部引用 failure 入口
```

### 適用時機

- 同一個知識點同時有「診斷/偵測」和「解法/預防」兩個面向
- 多個入口（failure、signals、anti-patterns）指向同一個解法
- 解法需要在多處被引用（避免重複）

### 不適用時機

- 純粹的參考資料（如 API 規格、資料格式），沒有「診斷 vs 解法」的區分
- 只有單一面向的知識點（只有解法沒有診斷，或只有診斷沒有解法）

← [回到共用規則索引](README.md)
