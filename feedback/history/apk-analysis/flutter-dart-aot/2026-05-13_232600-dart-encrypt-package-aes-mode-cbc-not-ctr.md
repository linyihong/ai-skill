> 遵守 [共用規則索引](../../../shared-rules/README.md)、[dependency-reading](../../../shared-rules/dependency-reading.md)、[neutral-language](../../../shared-rules/neutral-language.md)、[goal-action-validation](../../../shared-rules/goal-action-validation.md) 與 [feedback-lessons](../../../shared-rules/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-13 - Dart `encrypt` 套件 AES 模式：CBC 非 CTR/SIC

Status: validated

#### One-line Summary

逆向 Dart `encrypt` 套件的 AES 加密時，不要假設使用預設的 `AESMode.sic`（CTR）——App 可能明確傳入 `AESMode.cbc`。`PaddedBlockCipher` 包裹 `CBCBlockCipher` 是關鍵指標。透過 device proxy 實測 `AES/CBC/PKCS5Padding` vs `AES/CTR/NoPadding` 是確認模式的最終手段。

#### Human Explanation

Dart `encrypt` 套件 v5.0.3 的 `AES()` 建構子預設模式是 `AESMode.sic`（SIC = CTR）。但 App 可以明確傳入 `AESMode.cbc`。分析 Frida capture 時：

1. **`PaddedBlockCipher` 在 CBC 和 SIC 兩種模式下都會觸發**——兩種模式都可以被 `PaddedBlockCipher` 包裹以支援 PKCS7 padding。看到 `PBC.ctor` 無法區分兩者。
2. **`CBCBlockCipher.ctor` vs `SICBlockCipher.ctor`**——直接 hook 這兩個建構子可以區分，但 offset 可能錯誤或難以找到。
3. **43 次 `processBlock` 呼叫具有歧義性**——CBC（42 data + 1 padding）、GCM（1 J0 + 42 CTR）、CTR+PKCS7（42 data + 1 padding）都會產生 43 個 blocks。僅靠 block 數量無法區分。
4. **`GCMBlockCipher.ctor` 可能在不相關的加密群組中觸發**——在 v7 capture 中，`GCMBlockCipher.ctor` 在第二個加密群組（events 9962-9988）觸發，但這**不是** `eh` header 的加密。第一個群組（events 30-86）才是真正的 `eh` header 加密。

**最終確認方法**：透過 device proxy 實測比較兩種模式：
- `AES/CBC/PKCS5Padding` + prefix-as-IV → HTTP 200（SUCCESS）
- `AES/CTR/NoPadding` + 手動 PKCS7 padding → HTTP 551（rejected）

這證明了模式是 CBC，不是 CTR/SIC。

#### Trigger

- Frida capture 顯示 `PaddedBlockCipher.ctor` 但無法區分內部是 `CBCBlockCipher` 還是 `SICBlockCipher`
- Block 數量分析（43 blocks）在 CBC、GCM、CTR+PKCS7 之間具有歧義性
- `GCMBlockCipher.ctor` 在某些加密群組觸發，但不是目標群組
- Java 實作使用 `AES/CTR/NoPadding` 回傳 HTTP 551（PerimeterX 阻擋）

#### Evidence

- Tool: Live proxy test（`/tmp/CompareModesLiveTest.java`）比較 `AES/CBC/PKCS5Padding` vs `AES/CTR/NoPadding`
- Sanitized excerpt:
  - `AES/CBC/PKCS5Padding` → HTTP 200，含 Set-Cookie（PHPSESSID, sl-session）
  - `AES/CTR/NoPadding` → HTTP 551（PerimeterX）
- Evidence path: `<PROJECT_ROOT>/TATA/docs/eh-header-generation-status-2026-05-13.md`

#### Generalized Lesson

1. **Dart `encrypt` 套件的預設模式（`AESMode.sic`）不可靠**——務必檢查 App 是否明確傳入不同模式。`AES.AES.ctor` hook 的 `args[2]`（mode 參數）是以物件指標傳遞，不是字串，所以 `tryReadDartString` 會顯示 "unknown"。
2. **`PaddedBlockCipher` 同時包裹 CBC 和 SIC**——看到 `PBC.ctor` 無法區分兩者。必須直接 hook `CBCBlockCipher.ctor` 或 `SICBlockCipher.ctor`。
3. **43 次 `processBlock` 呼叫具有歧義性**——可能是 CBC（42 data + 1 padding）、GCM（1 J0 + 42 CTR）或 CTR+PKCS7（42 data + 1 padding）。僅靠 block 數量不夠。
4. **Live proxy test 是最終確認手段**——不確定時，執行獨立 Java 測試，透過 device proxy 以不同 cipher 配置發送請求。HTTP status code（200 vs 551）是 ground truth。
5. **Java `PKCS5Padding` 對 AES 而言等同於 PKCS7**——兩者使用相同的 padding 方案（名稱差異是歷史因素）。`Cipher.getInstance("AES/CBC/PKCS5Padding")` 匹配 Dart 的 `PaddedBlockCipher('AES/CBC/PKCS7')`。

#### Agent Action

分析 Dart `encrypt` 套件的 AES 加密時：

1. **不要假設使用預設模式**——檢查 App 是否明確傳入 `AESMode.cbc`
2. **同時 hook `CBCBlockCipher.ctor` 和 `SICBlockCipher.ctor`**——不只是 `PBC.ctor`
3. **驗證加密群組的隔離性**——確保分析的 Frida events 屬於目標加密（如 `eh` header），而不是不相關的加密群組
4. **不確定時，執行 live proxy test**——透過 device proxy 比較 `AES/CBC/PKCS5Padding` vs `AES/CTR/NoPadding`
5. **Java `AES/CBC/PKCS5Padding`** 是 Dart `PaddedBlockCipher('AES/CBC/PKCS7')` 的正確對應——不需要手動 padding

#### Goal / Action / Validation

- Goal: 正確識別 Dart `encrypt` 套件使用的 AES 模式
- Action: 執行 live proxy test 比較 CBC vs CTR 模式；在 Frida 中確認 `CBCBlockCipher.ctor` 的 offset
- Validation or reference source: HTTP 200 且含 Set-Cookie headers 確認模式正確

#### Applies When

- 從 Frida AOT capture 逆向 Dart `encrypt` 套件的 AES 加密
- `PaddedBlockCipher.ctor` 觸發但 `CBCBlockCipher.ctor` / `SICBlockCipher.ctor` 的 hook 缺失或失敗
- Block 數量分析具有歧義性（43 blocks 可能是 CBC、GCM 或 CTR+PKCS7）
- 有 device proxy（adb forward）可用於 live testing

#### Does Not Apply When

- Dart App 使用不同的加密函式庫（如直接使用 `pointycastle`，非 `encrypt` 套件）
- 沒有 device proxy 可用於 live testing
- 加密沒有被 `PaddedBlockCipher` 包裹（如直接使用 `AESEngine`）

#### Validation

- 已透過 live proxy test 確認：`AES/CBC/PKCS5Padding` → HTTP 200，`AES/CTR/NoPadding` → HTTP 551
- Dart `encrypt` 套件原始碼確認 `AESMode.cbc` 是有效的明確模式

#### Promotion Target

- `intelligence/engineering/analysis/heuristics/` — 新增 heuristic：「Dart encrypt package mode detection」
- `workflow/apk-analysis/execution-flow.md` — 新增步驟：「Frida capture 有歧義性時，透過 live proxy test 驗證 AES 模式」

#### Required Linked Updates

- 無需連動更新；這是新 lesson。
