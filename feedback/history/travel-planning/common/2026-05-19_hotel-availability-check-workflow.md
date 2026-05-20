> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md)、[neutral-language](../../../../enforcement/neutral-language.md)、[goal-action-validation](../../../../enforcement/goal-action-validation.md) 與 [feedback-lessons](../../../feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

# Hotel Availability Check Workflow

### 2026-05-19 - 查詢飯店空房的工具優先順序與限制

Status: validated

#### One-line Summary

查飯店空房時，優先用 Chrome 瀏覽器工具操作訂房網站；WebFetch 抓不到 JS 動態載入的訂房系統，只能取得靜態資訊（房型、容量、聯絡方式）。

#### Human Explanation

使用者要求幫忙查特定住宿的空房時，訂房系統（Booking.com、じゃらん、楽天トラベル、Jalan、各旅館官方預約系統）幾乎全部使用 JavaScript 動態渲染，WebFetch 只能取得靜態 HTML，拿不到即時空位資訊。

Chrome browser extension（`mcp__Claude_in_Chrome`）可直接操作瀏覽器，能填入日期、人數並讀取搜尋結果，是唯一能查到即時空位的工具路徑。

若 Chrome extension 不可用，應改用 WebSearch 取得：
- 旅館官方電話 / Email
- 可接受人數上限（房型容量）
- 訂房平台連結（讓使用者自行查詢）

並明確告知使用者「無法查到即時空位，需自行前往訂房網站或致電確認」。

#### Trigger

- 使用者詢問「XX 飯店/旅館 有沒有空房」「幫我查空位」「可不可以預訂」
- 使用者提供住宿名稱 + 入住日期 + 人數，要求確認可訂性

#### Evidence

- 任務：查詢日本特定地區多間旅館的多人團體空房（紅葉旺季）
- Chrome extension 未連線，改用 WebFetch 嘗試訂房系統
- WebFetch 對某旅館官方訂房系統（SPA 架構）回傳空白內容
- WebFetch 對另一旅館訂房系統回傳 403
- 最終只能透過 WebSearch 取得靜態資訊（房型容量、電話、訂房平台連結）
- 詳細任務記錄留於 `<PROJECT_ROOT>/itinerary.md`

#### Generalized Lesson

查詢飯店空房的標準工作流程：

**Step 1：確認 Chrome extension 是否可用**
```
mcp__Claude_in_Chrome__tabs_context_mcp → createIfEmpty: true
```
- 若可用 → 進入 Step 2（瀏覽器操作）
- 若不可用 → 進入 Step 3（降級路線）

**Step 2：Chrome 操作（可取得即時空位）**
1. 前往訂房平台（Booking.com / じゃらん / Agoda 等）
2. 填入旅館名稱、入住日、退房日、人數
3. 讀取搜尋結果（`read_page` 或 `get_page_text`）
4. 回報可用房型與價格

**Step 3：降級路線（WebSearch + WebFetch，無法取得即時空位）**
1. `WebSearch`：搜尋旅館名稱，取得官網、訂房平台連結、電話
2. `WebFetch`：抓取旅館官網，取得：
   - 房型與最大入住人數
   - 電話 / Email（讓使用者直接致電或發信詢問）
   - 訂房平台連結清單
3. 明確告知使用者：「訂房系統為 JS 動態載入，WebFetch 無法取得即時空位，請直接前往以下連結或致電確認」

**關鍵限制說明（必須對使用者說清楚）**
- 旅館官方訂房系統與主流預訂平台幾乎全為 SPA / JS-heavy，WebFetch 抓不到即時空位
- WebFetch 只讀靜態 HTML，對動態系統結果為空白或 403
- 即使 Chrome extension 可用，部分系統有 CAPTCHA 或登入牆，可能仍需使用者自行操作

#### Anti-patterns

- ❌ 用 WebFetch 直接抓訂房系統 URL，誤以為能查到空位
- ❌ 未告知使用者查詢限制，讓使用者誤以為已確認空房
- ❌ Chrome extension 連線失敗後未切換至降級路線，直接放棄

#### Agent Action

- 先執行 `mcp__Claude_in_Chrome__tabs_context_mcp` 確認 Chrome 連線狀態，**再**選擇工具路徑
- 降級路線結束後，**必須**明確告知使用者「無法取得即時空位」及替代聯絡方式
- 不得沉默跳過空位查詢，也不得只提供訂房平台連結而不說明限制

#### Goal / Action / Validation

- Goal: 在 Chrome 可用時取得即時空位；不可用時提供靜態資訊 + 明確告知限制
- Action: 依三步驟工作流程執行（確認連線 → Chrome 操作 或 降級路線）
- Validation or reference source: 使用者能根據回報的資訊（即時空位 或 電話/連結）自行完成預訂決策

#### Applies When

- 使用者要求查詢特定住宿的空位或可訂性
- 任務涉及日本訂房平台（Booking.com、じゃらん、楽天、Agoda、一休、Trip.com）或旅館官方預約系統

#### Does Not Apply When

- 使用者只要求旅館資訊（地址、設施、價格範圍），不要求確認空位
- 任務是規劃行程候選清單，尚未到確認空位階段

#### Validation

- 執行 Step 1（確認 Chrome 連線）後，工具路徑是否與連線狀態一致
- 降級路線結束後，是否已回報限制說明與替代聯絡方式

#### Promotion Target

- `workflow/travel-planning/execution-flow.md`（Step 8 Lodging Planning 工具選擇段）
- `analysis/travel/sources-and-tools.md`（訂房平台工具欄位）

#### Required Linked Updates

- `workflow/travel-planning/execution-flow.md`：在 Lodging Planning 步驟加入 Chrome-first / WebFetch-降級的工具選擇規則
- `analysis/travel/sources-and-tools.md`：在訂房平台來源行補充「JS 動態系統限制，WebFetch 僅取靜態資訊」
- `knowledge/summaries/travel-planning.md`：確認是否需補充工具限制摘要（若未提及則新增）
- 已依 [reusable-guidance-boundary.md](../../../../enforcement/reusable-guidance-boundary.md) 確認：具體旅館名稱與旅行日期已從 Evidence 移除，僅保留通用規則

#### Related

- 工具：`mcp__Claude_in_Chrome__navigate`、`mcp__Claude_in_Chrome__form_input`、`mcp__Claude_in_Chrome__get_page_text`
- 訂房平台：Booking.com、Agoda、じゃらんnet、楽天トラベル、一休.com、Trip.com
