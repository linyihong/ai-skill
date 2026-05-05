# Ai-skill

這個 repository 是**專業共用**的 AI skill 知識庫，一般透過 **Git 與遠端同步**（與團隊或跨機器共用流程一致）。

未來不同專案需要 APK 分析、抓包、Frida、Proxyman、Dart AOT、解密或分析方法沉澱時，優先從這裡讀 skill，並把新的可重用技巧回饋回這裡。

**路徑約定：**在本機上，中央庫位置就是你心裡的 `<AI_SKILL_REPO>`（clone 所在目錄）。**請勿把真實本機絕對路徑寫進本庫任何會 commit 的文件**（含 `feedback_history`、`shared-rules`、規則範本）；對外一律用占位符 `<AI_SKILL_REPO>`、`<PROJECT_ROOT>`。

## 現有 Skills

| Skill | 用途 |
| --- | --- |
| `skills/apk-analysis/` | 授權 APK 流量分析、動態抓包、Flutter/Dart AOT、response 解碼、文件化與 `feedback_history/` 技巧條目。 |
| `skills/app-development-guidance/` | 將授權 App/API/embedded/firmware/hardware-product 觀察轉成開發 guidance、實作模式、控制項、檢查清單與驗證方法。 |

**目錄約定：**  
- **`skills/`**：各 skill 技巧包；之後新增 skill 放在 `skills/<name>/`，步驟見 [`skills/ADDING_SKILLS.md`](skills/ADDING_SKILLS.md)。  
- **`shared-rules/`**：**共用規則**（依主題分檔：授權、去敏、中性低爭議文件用語、內容分層、文件大小與拆分、cross-skill reference、**feedback 檔名／模板／agent 行為**、Cursor 同步等）；索引為 [`shared-rules/README.md`](shared-rules/README.md)，feedback 流程與模板集中在 [`feedback-lessons.md`](shared-rules/feedback-lessons.md)。各 skill 目錄下 **`FEEDBACK.md`** 若存在，僅為**一行入口**，不必重複維護正文。  
- **`ai-tools/`**：各 AI coding / agent 工具如何讀取、參照或同步本知識庫；Cursor 詳見 [`ai-tools/cursor.md`](ai-tools/cursor.md)，其他工具日後在此補檔。  
- **連動更新**：若改動會影響其他文件、索引、skill 入口、同步流程或分類文件，相關檔案**必須**依 [`shared-rules/linked-updates.md`](shared-rules/linked-updates.md) 同步更新或明確檢查，不得說成「可選」。  
- **每一則 `feedback_history`**：**不要**重複貼上共用規則全文，頂部引用 `shared-rules/` 即可。**同步到 `.cursor`** 時：先複製 **`shared-rules/`** 整包，再同步需要的 **`skills/<name>/`**（見 [`shared-rules/cursor-sync.md`](shared-rules/cursor-sync.md)）。

## AI 工具使用說明

詳細文件放在 [`ai-tools/`](ai-tools/)，根 README 只保留工具索引。新增其他 AI 工具時，在該目錄新增子檔並更新索引。

| 工具 | 文件 | 用途 |
| --- | --- | --- |
| Cursor | [`ai-tools/cursor.md`](ai-tools/cursor.md) | 啟用 apk-analysis、同步 `.cursor`、維持中央庫一致性。 |

## 回饋規則

可以回饋：

- 新的分析流程。
- 新的工具使用技巧。
- 新的失敗判讀方式。
- 新的去敏與 fixture 沉澱方式。
- 通用的媒體、解密、session、proxy、Frida 經驗。

不要回饋：

- 特定產品的完整 host / endpoint。
- token、secret、device id、帳號或個資。
- 未去敏 raw response。
- 本機絕對路徑、使用者名稱、私有工作目錄、clone 位置；請改用 `<AI_SKILL_REPO>`、`<PROJECT_ROOT>`、`<WORKSPACE>`。
- 只對單一產品有效、沒有泛化價值的結論。

## Git 規則

這個 repository 有自己的 git remote。更新 skill 後：

```bash
cd <AI_SKILL_REPO>
git status
git add .
git commit -m "Update apk analysis skill"
git push
```

不要把專案私有資料、抓包原始檔或未去敏樣本 commit 到這裡。

`.cursor` 的變更若屬於業務專案目錄，請勿誤把該專案才有的機密或絕對路徑抄進本庫；本庫維持泛化與 `<AI_SKILL_REPO>` / `<PROJECT_ROOT>` 等占位符。
