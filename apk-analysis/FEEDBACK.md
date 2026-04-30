# Skill 回饋：規則與流程（本檔保持精簡）

本檔只說明 **如何**、**在哪裡** 撰寫與維護回饋；**具體 lesson 全文**請放在 **`feedback_history/`**，避免單一 `FEEDBACK.md` 無限膨脹。

成熟後再將規則整理進 `WORKFLOW.md`、`TOOLS.md` 或 `DOCUMENTATION.md`（見各條目中的 Promotion Target）。

## 條目放哪裡

| 內容 | 位置 |
| --- | --- |
| **撰寫規範、模板、索引說明** | **`FEEDBACK.md`**（本檔） |
| **每一條獨立 lesson**（候選／已驗證／已 promoted 備份） | **`feedback_history/YYYY-MM-DD_HHMMSS-<slug>.md`** |
| **條目總覽表**（可選，方便人類掃一眼） | **`feedback_history/README.md`** |

### 檔名規則（時間 + `<slug>`）

- 使用 **`YYYY-MM-DD_HHMMSS-<slug>.md`**：
  - **`YYYY-MM-DD`**：建立 lesson 的日期（本機）。
  - **`HHMMSS`**：**24 小時制**本機時間（6 位數字，例：`143052` = 14:30:52）。含時間可避免同日多檔碰撞、也方便依檔名排序。
- `<slug>` 建議 **短英文 kebab-case** 或 **有意義的英數縮寫**（例：`proxy-two-layer-tls`、`aapt-resolve-activity`）；中文標題可保留但不宜過長。
- **同一秒多條**：微調秒數或改 `<slug>`，勿覆寫既有檔。
- **修改既有 lesson**：在原檔 **追加修訂說明**（簡短段落）或建新檔並在舊檔頂部標 `deprecated → 見 xxx.md`；不要默默刪除歷史。

## 回饋原則

- **Cursor agent：** 在授權 APK 分析過程中一旦得到可重用技巧／失敗模式／驗證規則，應**主動**在 **`feedback_history/` 新增一個** `.md` 檔（依下方模板），**同一輪對話內**完成，**不要**等使用者提醒「記得回饋」。完成後可在 **`feedback_history/README.md`** 補一行索引（若專案維護者有在用）。
- 只寫通用方法，不寫特定 App 的私有結論。
- 必須去敏。
- 不得寫入本機絕對路徑、使用者名稱、私有工作目錄、clone 位置；用 `<AI_SKILL_REPO>`、`<PROJECT_ROOT>`、`<WORKSPACE>` 等 placeholder。
- 必須說明證據與適用／不適用條件。
- 不確定的想法：`Status: experimental`（或標於標題旁）。
- **不要**把長篇 lesson 貼進 **`FEEDBACK.md`**；本檔僅維持規則與模板。

## 人類也能讀的寫法

每一條 lesson 都要讓沒參與當次分析的人也看得懂。建議結構見下方模板（One-line Summary / Human Explanation / Agent Action …）。

## 新 lesson 模板（複製到新檔 `feedback_history/YYYY-MM-DD_HHMMSS-<slug>.md`）

```markdown
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

- **索引**：維護者可定期更新 **`feedback_history/README.md`**（表格：檔名、Status、標題、一句話摘要）；agent 新增檔案後**可選**追加表格列。
- **Git**：`feedback_history/` 與 `FEEDBACK.md` 一併版本控制；不要提交含機密的原始 log。
- **已搬遷**：既有長篇 lesson 已拆至 **`feedback_history/*.md`**（歷史批次檔名時間戳為 **`120000`–`120010`** 僅供排序，新建請用**當下** `HHMMSS`）；請自此新增檔案而非往舊版單檔底部堆疊。
