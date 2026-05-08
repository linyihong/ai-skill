# 如何新增 Skill（方法與模板）

目標：在本 repository 的 **`skills/<skill-name>/`** 新增一個可被 agent 工具讀取的技巧包，並與 **`shared-rules/`** 對齊（不重複維護共用政策）。

## 1. 要不要開新 skill？

| 情況 | 建議 |
| --- | --- |
| 內容仍是 APK／流量／Frida／Flutter AOT 等同家族 | 擴充現有 **`apk-analysis/`**，用 `feedback_history/<category>/` 與 WORKFLOW／TOOLS／techniques 收斂 |
| 新領域（例：iOS 靜態流程、另一種合規測試方法論）且會長期累積 | 新建 **`skills/<name>/`** |

Skill **資料夾名稱**建議 **kebab-case**、簡短、穩定（例：`apk-analysis`，之後如 `ios-ipa-analysis`）。

## 2. 目錄結構（建議）

最小可運作組合：

```text
skills/<skill-name>/
  SKILL.md              # 必填：YAML frontmatter + 給 Agent 的入口正文
  feedback_history/     # 強烈建議：每條 lesson 一檔；若 skill 有分類，使用 feedback_history/<category>/（見 shared-rules/feedback-lessons.md）
```

常見擴充（依需求增量建立）：

```text
  README.md             # 給人類的導讀
  RUNBOOK.md            # 新專案第一天提示詞與路徑約定
  WORKFLOW.md           # 決策流程
  TOOLS.md              # 工具與環境
  DOCUMENTATION.md      # 產出格式
  tool-adapters/        # 可選：只有在某 AI 工具有 skill-specific 執行差異時建立
  FEEDBACK.md           # 可選：極短入口，連到 shared-rules/feedback-lessons.md
```

若單一文件開始過大，或一句規則展開成多個步驟、例外、模板、範例，依 [`shared-rules/document-sizing.md`](../shared-rules/document-sizing.md) 改成資料夾包裝：父層 `README.md` 做目錄與讀取路由，子檔保存具體內容。不要把不同技巧分類、寫作規範、工具教學和 feedback lesson 全堆在同一檔。

可直接複製範本（在 `<AI_SKILL_REPO>` 根目錄執行，將 `my-skill` 改成你的名稱）：

```bash
mkdir -p "skills/my-skill/feedback_history"
cp "skills/_template/SKILL.md" "skills/_template/FEEDBACK.md" "skills/my-skill/"
# 再編輯 skills/my-skill/SKILL.md 與 FEEDBACK.md，替換 <…> 占位符
# 可選：建立 skills/my-skill/README.md、RUNBOOK.md 等
# 可選：touch skills/my-skill/feedback_history/README.md 做索引表
```

## 3. `SKILL.md` 必填欄位（Agent）

檔案**開頭**使用 YAML frontmatter（與 `apk-analysis` 相同風格）：

| 欄位 | 說明 |
| --- | --- |
| `name` | 與資料夾名一致或至少可對應；建議 kebab-case |
| `description` | **英文**一段話，描述何時應套用此 skill（觸發詞、場景、邊界）。寫給路由／檢索用，避免特定客戶名與機密 |

正文建議包含：**何時用／不用**、**Quick Start**、**如何回饋**（指向 `shared-rules/feedback-lessons.md`）、**預設讀哪些檔**。

完整占位稿見 **[`_template/SKILL.md`](_template/SKILL.md)**。

## 4. 與共用規則的關係

- **授權、去敏、中性低爭議用語、目標/執行/驗證流程、依賴讀取鐵則、路徑占位、feedback 檔名與模板**：一律只維護在 **`shared-rules/`**，各 skill **不要**複製長文。
- **Cross-skill reference**：若某 skill 需要引用另一個 skill 的規範、模板、檢查清單或交接產物，依 **[`shared-rules/cross-skill-references.md`](../shared-rules/cross-skill-references.md)** 寫短引用，包含 target skill、觸發條件、交接 artifact、ownership boundary、去敏邊界與 linked updates；不要複製 target skill 的全文。
- **連動更新規則**：一律只維護在 **[`shared-rules/linked-updates.md`](../shared-rules/linked-updates.md)**；新增 skill 或修改 skill 結構時，受影響的索引、入口、同步文件、分類文件**必須**同步更新或明確檢查。
- **文件大小與拆分規則**：一律只維護在 **[`shared-rules/document-sizing.md`](../shared-rules/document-sizing.md)**；skill、技巧分類與寫作規範變大時，用資料夾與 `README.md` 目錄拆分。
- **中性與低爭議用語**：一律只維護在 **[`shared-rules/neutral-language.md`](../shared-rules/neutral-language.md)**；新增 skill 的標題、description、檔名、slug、索引與摘要都要避免高風險或容易造成 AI/搜尋誤判的詞，改用授權、合規、契約、風險控制等中性語境。
- **工具中立文件**：一律只維護在 **[`shared-rules/tool-neutral-documentation.md`](../shared-rules/tool-neutral-documentation.md)**；新增 skill 的 README / SKILL / workflow / template 預設不寫特定工具路徑、hook、UI 或同步細節。工具全域做法放到 [`ai-tools/`](../ai-tools/README.md)；若某 skill 對某工具有必要的執行差異，用 Strategy-style adapter 放 `skills/<skill>/tool-adapters/<tool>.md`，只寫差異並連回核心 workflow。
- **文件 TODO**：一律只維護在 **[`shared-rules/document-todo-list.md`](../shared-rules/document-todo-list.md)**；若新增 skill 文件仍有未完成、待決策、待補強或待驗證項目，在文件前段放 `Document TODO` 表並連到相關章節或 goal。
- **目標、執行、驗證流程**：一律只維護在 **[`shared-rules/goal-action-validation.md`](../shared-rules/goal-action-validation.md)**；新增 skill 的輸出格式、workflow、documentation 規則要能讓重要結論反查目標、執行、驗證，純判斷題則附參考來源與推論邊界。
- **依賴文件讀取鐵則**：一律只維護在 **[`shared-rules/dependency-reading.md`](../shared-rules/dependency-reading.md)**；新增、修改或重新讀取 skill 時，必須建立 dependency read ledger，讀 skill 入口、相關 README/workflow/checklist/template、shared-rules 與 linked updates，並把不存在的檔案標成 `not applicable`，不能只讀單一檔案。
- 各 skill 的 **`FEEDBACK.md`**（若需要）：維持與 [`apk-analysis/FEEDBACK.md`](apk-analysis/FEEDBACK.md) 相同模式——**幾行連結**到 [`shared-rules/feedback-lessons.md`](../shared-rules/feedback-lessons.md)。
- 每一則 lesson 頂部引用 `shared-rules`（路徑依檔案深度調整 `../../../shared-rules/...`）。

## 5. 在本庫登記與提交

1. 編輯根目錄 [**`README.md`**](../README.md) 的「現有 Skills」表格，加一行 `skills/<name>/` 與簡短說明。
2. 編輯 [**`skills/README.md`**](README.md) 的表格，同步列出。
3. 依 [`shared-rules/linked-updates.md`](../shared-rules/linked-updates.md) 檢查是否還需要同步更新 `RUNBOOK.md`、`WORKFLOW.md`、`DOCUMENTATION.md`、同步腳本或 cross-link。
4. 依 `shared-rules/dependency-reading.md` 關閉 Ai-skill writeback transaction：檢查 diff / linked updates，執行必要 sync，`git add` → `commit` → `push`，push 後讀回並確認 clean。

## 6. 同步到本機工具（可選）

**`shared-rules/`** 與各 **skill** 要成對出現：任何工具若只看到 skill 而看不到 shared rules，就會缺少授權、去敏、依賴讀取與 linked updates 底線。

具體工具部署、symlink、bundle、hook、reload 或設定方式放在 [`ai-tools/`](../ai-tools/README.md)。新增或修改 skill 後，依你使用的工具文件執行必要同步。

## 7. 工具 Strategy adapter（可選）

只有在某個 AI 工具對這個 skill 有真實執行差異時才建立 adapter，例如 tool event、hook、prompt injection、上下文載入順序、輸出限制或失敗模式不同。

建議結構：

```text
skills/<skill-name>/tool-adapters/
  README.md
  <tool>.md
```

寫法：

- 核心 `README.md` / `WORKFLOW.md` 保持工具中立，像 strategy interface。
- Adapter 只寫該工具的 execution strategy，不複製核心 workflow。
- Adapter 必須連回核心 workflow / tools / documentation 的相關章節。
- 工具全域設定仍放 `ai-tools/<tool>.md`；adapter 只放「這個 skill 對該工具」的差異。
- 若多個 adapter 重複同一段內容，抽回核心 skill 或 shared rule。

## 8. 檢查清單（新建完成前）

- [ ] `SKILL.md` 有合法 `name` / `description` frontmatter
- [ ] 正文有連到 `shared-rules` 與 `feedback-lessons`
- [ ] 標題、description、檔名、slug、索引與摘要已依 `shared-rules/neutral-language.md` 使用中性低爭議用語
- [ ] 文件已依 `shared-rules/tool-neutral-documentation.md` 保持工具中立；工具專屬路徑、hook、UI、同步步驟已放 `ai-tools/`
- [ ] 若建立 `tool-adapters/<tool>.md`，內容只包含 skill-specific 工具差異，並已連回核心 workflow；工具全域設定仍留在 `ai-tools/<tool>.md`
- [ ] 若文件尚有未完成、待決策、待補強或待驗證項目，已依 `shared-rules/document-todo-list.md` 在前段加入 TODO 表並連到相關章節/goal
- [ ] 輸出格式已依 `shared-rules/goal-action-validation.md` 要求重要工作單元包含目標、執行、驗證或參考來源
- [ ] 已依 `shared-rules/dependency-reading.md` 建立 dependency read ledger，讀取或明確檢查相關依賴文件，並標示缺檔不適用項
- [ ] 若引用其他 skill，已依 `shared-rules/cross-skill-references.md` 寫明 trigger、artifact、ownership boundary 與 linked updates
- [ ] 若文件開始變大，已依 `shared-rules/document-sizing.md` 拆成資料夾、目錄與子檔
- [ ] 已建立 `feedback_history/`（可先要 `README.md` 索引；若 skill 有分類，同步建立 `feedback_history/<category>/README.md`）
- [ ] 根目錄 `README.md` 與 `skills/README.md` 已更新
- [ ] 已依 `shared-rules/linked-updates.md` 完成或明確檢查必要連動更新
- [ ] 無真實本機絕對路徑、無機密寫入將 commit 的檔案

← [回到 skills 索引](README.md)
