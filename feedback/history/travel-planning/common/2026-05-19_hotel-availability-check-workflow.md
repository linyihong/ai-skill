> 遵守 [共用規則索引](../../../../enforcement/README.md) 與 [feedback-lessons](../../../../enforcement/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

# Hotel Availability Check Workflow

### 2026-05-19 - 查詢飯店空房的工具優先順序與限制

Status: active

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
- 絕大多數日本旅館訂房系統（489ban、yadosys、jalan、rakuten travel 等）均為 SPA / JS-heavy
- WebFetch 只讀靜態 HTML，對這類系統結果為空白或 403
- 即使 Chrome extension 可用，部分系統有 CAPTCHA 或登入牆，可能仍需使用者自行操作

#### Anti-patterns

- ❌ 用 WebFetch 直接抓訂房系統 URL，誤以為能查到空位
- ❌ 未告知使用者查詢限制，讓使用者誤以為已確認空房
- ❌ Chrome extension 連線失敗後未切換至降級路線，直接放棄

#### Related

- 工具：`mcp__Claude_in_Chrome__navigate`、`mcp__Claude_in_Chrome__form_input`、`mcp__Claude_in_Chrome__get_page_text`
- 訂房平台：Booking.com、Agoda、じゃらんnet、楽天トラベル、一休.com、Trip.com
