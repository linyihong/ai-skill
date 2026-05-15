> 遵守 [共用規則索引](../../../enforcement/README.md)、[dependency-reading](../../../enforcement/dependency-reading.md)、[neutral-language](../../../enforcement/neutral-language.md)、[goal-action-validation](../../../enforcement/goal-action-validation.md) 與 [feedback-lessons](../../../enforcement/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-15 — 伺服器欄位驗證探針：先測試伺服器真正驗證哪些欄位，再逆向生成邏輯

Status: candidate

#### One-line Summary

在花時間逆向複雜的請求欄位生成邏輯之前，先系統性地測試伺服器是否真的驗證這些欄位的內容——伺服器可能根本不驗證。

#### Human Explanation

逆向專有 API 時，很容易假設每個請求欄位都必須「正確」生成。這會導致花大量精力逆向那些伺服器實際上不驗證的欄位。

正確的做法是：
1. 先用捕獲的值建立一個可工作的 baseline 請求
2. 然後逐一變化每個欄位（改成隨機值或垃圾值）
3. 如果伺服器仍然接受修改後的請求，表示該欄位不需要逆向——可以用合理的預設值自生成

這在以下場景特別重要：
- Identity material / 裝置指紋欄位
- 可能是可選的客戶端生成 token
- 伺服器可能只記錄但不驗證的欄位

#### Trigger

你有一個使用捕獲值可以正常工作的 API 請求，但某些欄位的生成邏輯非常複雜，需要大量逆向工程（Dart AOT 逆向、native code 分析等）。

#### Evidence

- Tool: Java JUnit 測試（`probeIdentityFieldValidation()`）
- Sanitized excerpt: 一個從最小欄位集合（app constants + random device_id）開始的測試，逐步加入每個欄位並使用修改後的值。全部 6 個探針測試都回傳 HTTP 200，證明伺服器不驗證 identity material 的內容。
- Evidence path: `<PROJECT_ROOT>/apk-analysis-sdk/tata-sdk-tests/src/test/java/com/tata/sdk/tests/live/LiveGuestLoginTest.java` — `probeIdentityFieldValidation()` 方法

#### Generalized Lesson

1. **先探針再逆向**：在投入時間逆向某個欄位的生成演算法之前，先測試伺服器是否真的驗證它。
2. **從最小開始**：從最少需要的欄位開始（app constants + 一個隨機識別碼）。如果這樣就能成功，你就省下了大量工作。
3. **一次只變一個欄位**：系統性地將每個欄位改成垃圾值/隨機值。如果伺服器仍然接受，表示該欄位不被驗證。
4. **記錄探針結果**：記錄哪些欄位被驗證、哪些不被驗證——這會指導自生成實作的方向。

#### Agent Action

遇到有很多請求欄位的複雜 API 時：

1. 先用捕獲值建立可工作的 baseline。
2. 建立一個探針測試，對每個欄位發送修改後的值。
3. 如果伺服器接受修改後的值，將該欄位標記為「可自生成，不需要逆向」。
4. 只逆向那些伺服器真正驗證的欄位。

#### Goal / Action / Validation

- Goal: 避免浪費精力逆向伺服器不驗證的欄位。
- Action: 建立一個系統性的探針測試，獨立變化每個欄位。
- Validation or reference source: 如果所有探針測試都通過（HTTP 200），表示這些欄位不被驗證。如果某個探針測試失敗（HTTP 400/403），表示該欄位需要逆向。

#### Applies When

- 逆向一個有許多請求欄位的專有 API
- 已知伺服器接受請求參數的某些變化
- 你有使用捕獲值的可工作 baseline 請求

#### Does Not Apply When

- 已知伺服器嚴格驗證所有欄位（例如加密簽名）
- 你沒有可工作的 baseline 請求
- 根據協定規範，這些欄位明顯是必需的

#### Validation

探針測試的結果是自我驗證的：如果用修改後的值仍然回傳 HTTP 200，表示該欄位不被驗證。建議執行多次探針測試以確認一致性。

#### Promotion Target

- `workflow/apk-analysis/execution-flow.md` — 加入步驟「在逆向之前先探測伺服器欄位驗證」
- `intelligence/engineering/analytical-reasoning/heuristics/` — 加入作為通用分析啟發式

#### Required Linked Updates

- 目前不需要連動更新。這是一條新的 lesson，不取代或使任何既有 lesson 失效。
