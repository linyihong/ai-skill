> 遵守 [共用規則索引](../../../enforcement/README.md)、[dependency-reading](../../../enforcement/dependency-reading.md)、[neutral-language](../../../enforcement/neutral-language.md)、[goal-action-validation](../../../enforcement/goal-action-validation.md) 與 [feedback-lessons](../../../enforcement/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-14 - Dart 自訂 AES 實作：8-byte 金鑰無法以標準 Java/BC AES 重現

> **⚠️ 2026-05-14 更新：此問題已解決。** 8-byte 金鑰可以透過 **AES/CTR/NoPadding + PKCS7 padding** 搭配 **zero-padding 至 16 bytes** 的方式以標準 Java AES 重現。詳見下方「## 已解決」章節。

Status: resolved ✅

#### One-line Summary

逆向 Dart AOT 中的自訂 AES 實作時，若 Frida 顯示 AES 金鑰僅 8 bytes，則該實作使用非標準的金鑰衍生方式（key derivation），無法以標準 Java `javax.crypto.Cipher` 或 Bouncy Castle `AESEngine` 重現——兩者皆要求 16/24/32 bytes 金鑰。需透過 Frida 深入 hook 金鑰擴展（key expansion）階段或反組譯 Dart AOT 函式來找出實際的金鑰衍生邏輯。

#### Human Explanation

在分析某 App 的 header 生成時，Frida capture 顯示一個自訂的 `aesEncryptFallback` 函式：

1. 接收 8-byte 金鑰（x2）和 8-byte IV（x4，與金鑰相同）
2. 輸出為 44 字元 Base64（= 32 decoded bytes），符合 29 bytes 明文 + 3 bytes PKCS7 padding
3. 所有捕獲的 ciphertext 共享相同的 20-byte 固定前綴

嘗試以標準 Java AES 重現時，所有方法都失敗：

| 方法 | 結果 | 原因 |
|------|------|------|
| `javax.crypto.Cipher` AES/GCM | ❌ 失敗 | 需要 16-byte 金鑰、12-byte IV |
| `javax.crypto.Cipher` AES/CBC | ❌ 失敗 | 需要 16-byte 金鑰、16-byte IV |
| BC `AESEngine` | ❌ 失敗 | 需要 16/24/32-byte 金鑰 |
| BC `CBCBlockCipher` | ❌ 失敗 | 需要 16-byte IV |
| 重複填充金鑰至 16 bytes | ❌ 不匹配 | 輸出與 Dart 不同 |
| Zero-pad 金鑰至 16 bytes | ❌ 不匹配 | 輸出與 Dart 不同 |
| MD5(key) 作為金鑰 | ❌ 不匹配 | 輸出與 Dart 不同 |
| SHA-256(key) 前 16 bytes | ❌ 不匹配 | 輸出與 Dart 不同 |
| key+seed 串接 | ❌ 不匹配 | 輸出與 Dart 不同 |
| seed+key 串接 | ❌ 不匹配 | 輸出與 Dart 不同 |
| 反轉金鑰 | ❌ 不匹配 | 輸出與 Dart 不同 |
| XOR 0xFF | ❌ 不匹配 | 輸出與 Dart 不同 |

**關鍵教訓**：Dart 的自訂 AES 實作（非 `encrypt` 套件，而是直接操作 `Uint8List` 的函式）使用非標準的金鑰衍生方式。標準 AES 實作不接受 8-byte 金鑰，因此無法直接重現。需要反組譯 Dart AOT 函式或新增 Frida hook 追蹤金鑰擴展階段。

#### Trigger

- Frida capture 顯示 AES 金鑰長度為 8 bytes（或其他非 16/24/32 的長度）
- 嘗試以標準 Java `Cipher.getInstance("AES/GCM/NoPadding")` 或 `Cipher.getInstance("AES/CBC/PKCS5Padding")` 加密時，金鑰規格（`SecretKeySpec`）需要 16/24/32 bytes
- 嘗試以 Bouncy Castle `AESEngine` 或 `CBCBlockCipher` 加密時，金鑰或 IV 長度不符合要求
- 即使將金鑰填充至 16 bytes（重複、zero-pad、hash 等），輸出仍與 Dart 不一致
- 所有捕獲的 ciphertext 共享相同的固定前綴，但變數區域無法匹配

#### Evidence

- Tool: Frida capture script hooking `aesEncryptFallback` function
- Sanitized excerpt:
  - `aesEncryptFallback:enter | x2="<8byte_key>" (8 bytes) | x3="<8byte_seed>" (8 bytes) | x4="<8byte_key>" (8 bytes, IV=key) | x5="<29_byte_plaintext>" (29 bytes)`
  - `aesEncryptFallback:return | ret="<44_char_base64>" (44 chars Base64)`
- Test files: Multiple Java test files testing various key derivations (repeat-pad, zero-pad, MD5, SHA-256, concat, XOR, reversed) — all failed to match expected output
- Evidence path: `<PROJECT_ROOT>/capture/<frida_capture_log>.log`
- Dart AOT offset: Available from `unflutter`/`blutter` analysis

#### Generalized Lesson

1. **Dart 自訂 AES 實作可能使用非標準金鑰長度**——不同於 `encrypt` 套件（使用標準 AES），自訂實作（直接操作 `Uint8List`）可以接受任意長度的金鑰，並在內部進行非標準的金鑰衍生。
2. **8-byte 金鑰是紅旗信號**——標準 AES 需要 16/24/32 bytes 金鑰。看到 8-byte 金鑰時，表示 Dart 實作有自訂的金鑰擴展邏輯。
3. **標準 Java/BC AES 無法重現自訂金鑰衍生**——不要浪費時間嘗試各種填充/雜湊組合。直接反組譯 Dart AOT 函式或新增 Frida hook 追蹤金鑰擴展階段。
4. **IV = 金鑰是另一個紅旗**——當 IV 參數與金鑰相同時，表示實作可能使用非標準的 nonce/IV 建構方式。
5. **固定前綴可作為 fallback**——如果無法重現加密，可以捕獲一次固定前綴，只自訂產生變數區域，組合成有效 token。

#### Agent Action

分析 Dart AOT 中的自訂 AES 實作時：

1. **檢查金鑰長度**——如果 Frida 顯示金鑰長度不是 16/24/32 bytes，立即警覺這是自訂實作
2. **不要嘗試猜測金鑰衍生方式**——標準填充/雜湊方法（重複、zero-pad、MD5、SHA-256、串接、XOR）幾乎不可能猜中
3. **改為反組譯 Dart AOT 函式**——使用 `unflutter` 或 `blutter` 反組譯目標 offset，找出實際的金鑰擴展邏輯
4. **或新增 Frida hook 追蹤金鑰擴展**——在 `aesEncryptFallback` 內部 hook `Uint8List` 寫入操作，觀察金鑰如何被擴展
5. **考慮使用 captured prefix + self-generated variable 的混合方案**——如果金鑰衍生無法短期解決，可以捕獲固定前綴，只自訂產生變數區域
6. **記錄所有測試過的衍生方式**——避免未來重複測試相同的組合

#### Goal / Action / Validation

- Goal: 重現 Dart 自訂 AES 實作的加密輸出
- Action: 反組譯 Dart AOT 函式找出金鑰衍生邏輯，或新增 Frida hook 追蹤金鑰擴展階段
- Validation or reference source: 加密輸出與 Frida capture 的 ciphertext 完全一致

#### Applies When

- 從 Frida AOT capture 逆向 Dart 自訂 AES 實作（非 `encrypt` 套件）
- Frida 顯示 AES 金鑰長度為 8 bytes（或其他非標準長度）
- 標準 Java `javax.crypto.Cipher` 或 Bouncy Castle `AESEngine` 拒絕該金鑰長度
- 嘗試多種金鑰填充/衍生方式後仍無法匹配 Dart 輸出
- 有 Dart AOT 函式的 offset 可用於反組譯

#### Does Not Apply When

- Dart App 使用標準 `encrypt` 套件（金鑰長度為 16/24/32 bytes）
- 使用標準 Java/BC AES 即可成功重現
- 沒有 Dart AOT 函式的 offset 或無法反組譯
- 只需要 captured token 重放（不需要自訂產生）

#### Validation

- 已透過多個獨立 Java 測試檔案驗證：所有標準金鑰衍生方式（重複填充、zero-pad、MD5、SHA-256、串接、反轉、XOR）都無法匹配 Dart 輸出
- Frida capture 確認金鑰為 8 bytes、IV 與金鑰相同、輸出為 44 chars Base64
- Dart AOT offset 已識別但尚未反組譯

#### Promotion Target

- `intelligence/engineering/analytical-reasoning/failure/` — 新增 failure atom：「custom-dart-aes-8byte-key-not-reproducible.md」
- `intelligence/engineering/analytical-reasoning/heuristics/` — 更新 `hook-selection.md` 加入 8-byte key 的判斷規則

#### Required Linked Updates

- `<PROJECT_ROOT>/api/API列表/public/guest_login.md` — 更新 Encryption Parameters 表格（IV 與金鑰相同）；更新 Open Questions（新增 8-byte key derivation 問題）
- `<PROJECT_ROOT>/apk-analysis-sdk/.../SkyShieldXAspnetVersionProvider.java` — 更新 Javadoc 說明 8-byte key 無法以標準 AES 重現的事實；修正預設 IV

---

## ✅ 已解決（2026-05-14）

### 正確解法

SkyShield 的 `aesEncryptFallback` 使用 **標準 AES-128 CTR 模式**，但金鑰和 IV 在傳入前已被呼叫者 **zero-padded 至 16 bytes**。Frida hook 攔截到的是 padding 前的原始 8-byte 值，導致誤判為「自訂金鑰衍生」。

**正確參數**：

| 參數 | 值 | 說明 |
|------|-----|------|
| 加密模式 | `AES/CTR/NoPadding` | 即 SIC 模式，Java 中為 CTR |
| 金鑰 | `"l65tvNcw"` → zero-pad 至 16 bytes | `6c 36 35 74 76 4e 63 77 00 00 00 00 00 00 00 00` |
| IV/Nonce | `"5jZyks1r"` → zero-pad 至 16 bytes | `35 6a 5a 79 6b 73 31 72 00 00 00 00 00 00 00 00` |
| 填充 | PKCS7 padding（先 pad 再 CTR encrypt） | 29 bytes → 32 bytes (0x03 0x03 0x03) |
| 輸出 | Base64 encoded ciphertext | 32 bytes → 44 chars Base64 |

### Java 實作

```java
byte[] key16 = padZero("l65tvNcw".getBytes("UTF-8"), 16);
byte[] seed16 = padZero("5jZyks1r".getBytes("UTF-8"), 16);

Cipher cipher = Cipher.getInstance("AES/CTR/NoPadding");
cipher.init(Cipher.ENCRYPT_MODE, new SecretKeySpec(key16, "AES"), new IvParameterSpec(seed16));

// PKCS7 pad plaintext to 16-byte boundary first
byte[] paddedPt = pkcs7Pad(pt, 16);
byte[] result = cipher.doFinal(paddedPt);
String b64 = Base64.getEncoder().encodeToString(result);
```

### 驗證結果

- **28/28 測試案例全部通過** 🎉
- 所有 Frida capture 的 ciphertext 均與 Java 輸出完全一致
- 第一個 block 永遠相同（相同 key + IV），第二個 block 隨 random suffix 變化

### 為何之前失敗

| 錯誤嘗試 | 失敗原因 |
|----------|----------|
| AES/GCM | 錯誤的模式（GCM 需要 12-byte IV，且輸出包含 auth tag） |
| AES/CBC | 錯誤的模式（CBC 需要 16-byte IV，且輸出結構不同） |
| BC SICBlockCipher | 正確的模式但 `processBlock` 需要完整 16-byte block |
| 重複填充金鑰 | 正確的 key 是 zero-pad，不是 repeat-pad |
| 未先 PKCS7 pad 就 CTR encrypt | CTR 模式本身不需要 padding，但 Dart 實作先 pad 再 CTR |

### 關鍵教訓

1. **Frida hook 的參數不一定是加密函式收到的實際值**——呼叫者可能在傳入前已修改參數
2. **8-byte key 不一定是自訂金鑰衍生**——可能是呼叫者先 zero-pad 再傳入
3. **CTR 模式 + PKCS7 padding 是罕見但有效的組合**——通常 CTR 不需要 padding，但 Dart 實作確實先 pad 再 CTR
4. **不要只測試 CBC/GCM**——如果 CBC/GCM 都失敗，試試 CTR 模式
5. **第一個 block 匹配是關鍵信號**——如果第一個 block 的 keystream 匹配，表示 key 和 IV 正確，問題在於模式或 padding
