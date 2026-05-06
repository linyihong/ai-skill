# Feedback 與技巧條目（共用規則）

「怎麼寫回饋、檔名怎麼取、模板長怎樣」**全部**在本檔維護；各 skill **不再**另存一份 `FEEDBACK.md` 正文（`apk-analysis` 目錄下僅保留極短入口檔，指向本檔）。

**每一條 lesson 全文**放在對應 skill 的 **`feedback_history/`** 之下（未分類 skill 可直接放根層；已分類 skill 放 `<category>/` 或 `common/`），**不要**把長篇堆進任何說明檔。

## 原則

- **不要**在每條 lesson 裡重複貼上 [sanitization.md](sanitization.md)、[dependency-reading.md](dependency-reading.md)、[neutral-language.md](neutral-language.md)、[goal-action-validation.md](goal-action-validation.md)、[authorization-scope.md](authorization-scope.md) 等全文；條目頂部用一行**引用** [README.md](README.md) 或**本檔**即可。
- **Cursor agent：** 在授權分析過程中一旦得到可重用技巧／失敗模式／驗證規則，應**主動**在同一輪對話內於該 skill 的 **`feedback_history/`** 對應位置新增檔案（依下方**分類規則**、**檔名規則**與**模板**），**不要**等使用者提醒。
- 只寫**通用方法**，不寫特定 App 的私有結論；必須去敏；發現 skill/rule/template/lesson 更新時依 [dependency-reading.md](dependency-reading.md) 讀取依賴；標題、slug、摘要與正文必須依 [neutral-language.md](neutral-language.md) 使用中性低爭議用語；每個重要結論必須依 [goal-action-validation.md](goal-action-validation.md) 說明目標、執行、驗證或參考來源；必須說明證據與適用／不適用條件；不確定標 `experimental`。
- 不得寫入本機絕對路徑、使用者名稱、私有工作目錄、clone 位置；用 `<AI_SKILL_REPO>`、`<PROJECT_ROOT>`、`<WORKSPACE>` 等 placeholder。

## 條目放哪裡

| 內容 | 位置 |
| --- | --- |
| **共用政策（全庫）** | [`shared-rules/README.md`](README.md) |
| **本檔** | 命名規則、模板、索引與 Git 約定（**唯一正文**） |
| **每一條獨立 lesson（未分類 skill）** | **`<skill>/feedback_history/YYYY-MM-DD_HHMMSS-<slug>.md`** |
| **每一條獨立 lesson（已有分類的 skill）** | **`<skill>/feedback_history/<category>/YYYY-MM-DD_HHMMSS-<slug>.md`** |
| **條目總覽表**（可選） | **`<skill>/feedback_history/README.md`** 與必要的 **`<skill>/feedback_history/<category>/README.md`** |

範例：`skills/apk-analysis/feedback_history/`；若 skill 已有 `techniques/flutter-dart-aot/` 這類分類，對應 lesson 應放在 `skills/apk-analysis/feedback_history/flutter-dart-aot/`，跨分類或全域規則放 `skills/apk-analysis/feedback_history/common/`。

成熟後可將 lesson 整理進該 skill 的 `WORKFLOW.md`、`TOOLS.md` 或 `DOCUMENTATION.md`（見模板中 **Promotion Target**）。

## 分類規則

當某個 skill 內部已經開始按 runtime、platform、control、technique、checklist 等方式分類時，`feedback_history/` 也要跟著分類，避免所有 lesson 混在同一層：

- 新 lesson 優先放到 **`feedback_history/<category>/`**，其中 `<category>` 應對應該 skill 內既有分類名稱，例如 `flutter-dart-aot`、`http-api`、`controls`、`platforms`。
- 跨分類、全域適用或分類尚未明確的 lesson 放到 **`feedback_history/common/`**。
- 若一條 lesson 會 promote 到多個分類，放在主要分類，並在 lesson 的 **Promotion Target** / **Required Linked Updates** 寫出其他同步更新位置。
- `feedback_history/README.md` 應是總索引，列出 category folders；每個 category folder 可有自己的 `README.md` 表格。
- 既有歷史 lesson 若已被外部文件連結，可以先保留原路徑，並用 category README 索引清楚；若真的搬移，必須同一個 change 更新所有相對連結與索引。
- 不要為了分類而重複複製 lesson 內容；一條 lesson 只保留一份全文，其他地方用連結。

## 檔名規則（時間 + `<slug>`）

- 使用 **`YYYY-MM-DD_HHMMSS-<slug>.md`**：
  - **`YYYY-MM-DD`**：建立 lesson 的日期（本機）。
  - **`HHMMSS`**：**24 小時制**本機時間（6 位數字，例：`143052` = 14:30:52）。含時間可避免同日多檔碰撞、也方便依檔名排序。
- `<slug>` 建議 **短英文 kebab-case** 或 **有意義的英數縮寫**（例：`proxy-two-layer-tls`、`aapt-resolve-activity`）；中文標題可保留但不宜過長。
- **同一秒多條**：微調秒數或改 `<slug>`，勿覆寫既有檔。
- **修改既有 lesson**：在原檔 **追加修訂說明**（簡短段落）或建新檔並在舊檔頂部標 `deprecated → 見 xxx.md`；不要默默刪除歷史。

## 新 lesson 模板

複製到新檔：

- 未分類 skill：**`<skill>/feedback_history/YYYY-MM-DD_HHMMSS-<slug>.md`**
- 已分類 skill：**`<skill>/feedback_history/<category>/YYYY-MM-DD_HHMMSS-<slug>.md`**

注意：以下引用路徑以檔案位於 `feedback_history/` 內為準；若檔案在 `feedback_history/<category>/`，共用規則連結要多上一層，改成 `../../../../shared-rules/...`。

```markdown
> 遵守 [共用規則索引](../../../shared-rules/README.md)、[dependency-reading](../../../shared-rules/dependency-reading.md)、[neutral-language](../../../shared-rules/neutral-language.md)、[goal-action-validation](../../../shared-rules/goal-action-validation.md) 與 [feedback-lessons](../../../shared-rules/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### YYYY-MM-DD - [short title]

Status: candidate | validated | deprecated | promoted | experimental

#### One-line Summary

用一句人話說明這條 lesson。

#### Human Explanation

給人看的背景說明：為什麼重要、常見誤判是什麼、實務上怎麼判斷。

#### Trigger

遇到什麼現象或問題？

#### Evidence

- Tool:
- Sanitized excerpt:
- Evidence path:

#### Generalized Lesson

可重用的規則是什麼？

#### Agent Action

下次 agent 看到類似情境時，應該先做什麼、不要做什麼？

#### Goal / Action / Validation

- Goal:
- Action:
- Validation or reference source:

#### Applies When

- 條件 1

#### Does Not Apply When

- 條件 1

#### Validation

如何確認這條 lesson 是對的？

#### Promotion Target

- `WORKFLOW.md`
- `TOOLS.md`
- `DOCUMENTATION.md`
- `SKILL.md`

#### Required Linked Updates

- 依 [`linked-updates.md`](linked-updates.md) 列出必須同步更新或已檢查的相關文件；若無需連動更新，寫明原因。
```

## 同步與索引

- **索引**：維護者可定期整理 **`<skill>/feedback_history/README.md`**（表格：檔名、Status、標題、一句話摘要，或 category index）；若該 skill 已有索引，agent 新增 lesson 檔後**必須**追加表格列或明確說明為何暫不更新。已分類 skill 同時更新對應 **`feedback_history/<category>/README.md`**。
- **Git**：**`feedback_history/`** 版本控制；不要提交含機密的原始 log。
- **歷史**：既有長篇 lesson 應已拆至各 `feedback_history/*.md`（若見批次時間戳如 `120000`–`120010` 僅供排序，新建請用**當下** `HHMMSS`）；請自此新增檔案而非往舊版單檔底部堆疊。

← [回到共用規則索引](README.md)
