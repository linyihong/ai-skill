# Feedback 與技巧條目（共用規則）

「怎麼寫回饋、檔名怎麼取、模板長怎樣」**全部**在本檔維護；各 skill **不再**另存一份 `FEEDBACK.md` 正文（`apk-analysis` 目錄下僅保留極短入口檔，指向本檔）。

**每一條 lesson 全文**放在對應 skill 的 **`feedback_history/`**（例如 `skills/apk-analysis/feedback_history/`），**不要**把長篇堆進任何說明檔。

## 原則

- **不要**在每條 lesson 裡重複貼上 [sanitization.md](sanitization.md)、[authorization-scope.md](authorization-scope.md) 等全文；條目頂部用一行**引用** [README.md](README.md) 或**本檔**即可。
- **Cursor agent：** 在授權分析過程中一旦得到可重用技巧／失敗模式／驗證規則，應**主動**在同一輪對話內於該 skill 的 **`feedback_history/`** 新增檔案（依下方**檔名規則**與**模板**），**不要**等使用者提醒。
- 只寫**通用方法**，不寫特定 App 的私有結論；必須去敏；必須說明證據與適用／不適用條件；不確定標 `experimental`。
- 不得寫入本機絕對路徑、使用者名稱、私有工作目錄、clone 位置；用 `<AI_SKILL_REPO>`、`<PROJECT_ROOT>`、`<WORKSPACE>` 等 placeholder。

## 條目放哪裡

| 內容 | 位置 |
| --- | --- |
| **共用政策（全庫）** | [`shared-rules/README.md`](README.md) |
| **本檔** | 命名規則、模板、索引與 Git 約定（**唯一正文**） |
| **每一條獨立 lesson** | **`<skill>/feedback_history/YYYY-MM-DD_HHMMSS-<slug>.md`** |
| **條目總覽表**（可選） | **`<skill>/feedback_history/README.md`** |

範例：`skills/apk-analysis/feedback_history/`。

成熟後可將 lesson 整理進該 skill 的 `WORKFLOW.md`、`TOOLS.md` 或 `DOCUMENTATION.md`（見模板中 **Promotion Target**）。

## 檔名規則（時間 + `<slug>`）

- 使用 **`YYYY-MM-DD_HHMMSS-<slug>.md`**：
  - **`YYYY-MM-DD`**：建立 lesson 的日期（本機）。
  - **`HHMMSS`**：**24 小時制**本機時間（6 位數字，例：`143052` = 14:30:52）。含時間可避免同日多檔碰撞、也方便依檔名排序。
- `<slug>` 建議 **短英文 kebab-case** 或 **有意義的英數縮寫**（例：`proxy-two-layer-tls`、`aapt-resolve-activity`）；中文標題可保留但不宜過長。
- **同一秒多條**：微調秒數或改 `<slug>`，勿覆寫既有檔。
- **修改既有 lesson**：在原檔 **追加修訂說明**（簡短段落）或建新檔並在舊檔頂部標 `deprecated → 見 xxx.md`；不要默默刪除歷史。

## 新 lesson 模板

複製到新檔 **`<skill>/feedback_history/YYYY-MM-DD_HHMMSS-<slug>.md`**（以下引用路徑以檔案位於 `feedback_history/` 內為準）：

```markdown
> 遵守 [共用規則索引](../../../shared-rules/README.md) 與 [feedback-lessons](../../../shared-rules/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

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
```

## 同步與索引

- **索引**：維護者可定期更新 **`<skill>/feedback_history/README.md`**（表格：檔名、Status、標題、一句話摘要）；agent 新增檔案後**可選**追加表格列。
- **Git**：**`feedback_history/`** 版本控制；不要提交含機密的原始 log。
- **歷史**：既有長篇 lesson 應已拆至各 `feedback_history/*.md`（若見批次時間戳如 `120000`–`120010` 僅供排序，新建請用**當下** `HHMMSS`）；請自此新增檔案而非往舊版單檔底部堆疊。

← [回到共用規則索引](README.md)
