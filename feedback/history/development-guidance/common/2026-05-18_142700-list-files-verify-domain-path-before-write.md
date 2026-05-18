> 遵守 [共用規則索引](../../../enforcement/README.md)、[dependency-reading](../../../enforcement/dependency-reading.md)、[neutral-language](../../../enforcement/neutral-language.md)、[goal-action-validation](../../../enforcement/goal-action-validation.md) 與 [feedback-lessons](../../feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-18 — 寫入 Feedback Lesson 前必須先 `list_files` 確認目標目錄存在

Status: validated

#### One-line Summary

當要寫入 feedback lesson 到 `feedback/history/<domain>/` 時，**必須先 `list_files` 確認目標 domain 目錄確實存在**，不能依賴內部記憶中的路徑名稱。

#### Human Explanation

Agent 在長時間對話中容易累積過時的內部記憶（例如某個 domain 曾經叫 `app-development-guidance`，但已改名為 `development-guidance`）。如果直接依賴記憶中的路徑寫入，會把 lesson 放到不存在的目錄下，或者更糟 — 在錯誤的位置建立新目錄，造成後續維護混亂。

正確的做法是：在決定 domain 歸屬後，先用 `list_files` 查看 `feedback/history/` 下實際有哪些 domain 目錄，確認目標存在，再寫入。

#### Trigger

寫入 JSON substring matching 的 feedback lesson 時，腦中浮現舊路徑 `app-development-guidance`（已改名為 `development-guidance`），沒有先確認目錄結構就直接用了錯誤的路徑。

#### Evidence

- Tool: `list_files`（確認目錄存在）
- Sanitized excerpt: Agent 依賴內部記憶中的 `app-development-guidance`，但實際目錄是 `development-guidance`
- Evidence path: `<AI_SKILL_REPO>/feedback/history/` 下的 domain 列表

#### Generalized Lesson

**寫入任何檔案到 `feedback/history/<domain>/` 前，必須先用 `list_files` 確認目標 domain 目錄確實存在。** 不要依賴內部記憶或之前 session 的上下文。如果目標 domain 不存在，應建立新目錄（而非放到錯誤的既有目錄下）。

#### Agent Action

1. 決定 lesson 的 domain 歸屬後，立即執行 `list_files("feedback/history/")` 查看實際 domain 列表
2. 確認目標 domain 存在後，再執行 `list_files("feedback/history/<domain>/")` 查看分類目錄
3. 若 domain 不存在，判斷是否需要建立（而非放到其他 domain 下）
4. 確認路徑正確後再寫入

#### Goal / Action / Validation

- Goal: 確保 feedback lesson 永遠寫入正確的 domain 路徑
- Action: 寫入前先 `list_files` 確認目錄存在
- Validation: 寫入後再次 `list_files` 確認檔案出現在預期位置

#### Applies When

- 建立新的 feedback lesson 時
- 任何需要寫入 `feedback/history/<domain>/` 的操作
- Agent 依賴內部記憶而非檔案系統狀態時

#### Does Not Apply When

- 已在同一輪對話中剛執行過 `list_files` 確認過目錄結構（但若經過多輪對話，仍建議重新確認）
- 修改已存在的檔案（路徑已確認過）

#### Validation

寫入後執行 `list_files("feedback/history/<domain>/<category>/")` 確認新檔案出現在預期位置。

#### Promotion Target

- `feedback/feedback-lessons.md`（在判斷流程中補強「確認目錄存在」步驟）

#### Required Linked Updates

- 已依 [`linked-updates.md`](../../../enforcement/linked-updates.md) 檢查：
  - `feedback/feedback-lessons.md`：需在判斷流程 Step 2 後加入「確認目錄存在」的子步驟
  - `runtime/generated/README.md`：可加入寫 feedback 前的提醒查詢
