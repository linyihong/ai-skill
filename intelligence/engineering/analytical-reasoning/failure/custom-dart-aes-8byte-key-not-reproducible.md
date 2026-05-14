# Custom Dart AES 8-byte Key Not Reproducible（Dart 自訂 AES 8-byte 金鑰無法以標準 AES 重現）

> **⚠️ 2026-05-14 更新：此問題已解決。** 8-byte 金鑰可以透過 **AES/CTR/NoPadding + PKCS7 padding** 搭配 **zero-padding 至 16 bytes** 的方式以標準 Java AES 重現。詳見下方「## 已解決」章節。

## 問題

Frida hook 攔截到 Dart AOT 的自訂 `aesEncryptFallback` 函式使用 8-byte AES 金鑰時，無法以標準 Java `javax.crypto.Cipher` 或 Bouncy Castle `AESEngine` 重現加密結果。標準 AES 實作要求金鑰長度為 16/24/32 bytes，而 Dart 自訂實作使用非標準的金鑰衍生方式。

## 症狀

- Frida capture 顯示 AES 金鑰長度為 8 bytes（或其他非標準長度）
- 嘗試以 `Cipher.getInstance("AES/GCM/NoPadding")` 或 `Cipher.getInstance("AES/CBC/PKCS5Padding")` 加密時，`SecretKeySpec` 需要 16/24/32 bytes
- 嘗試以 Bouncy Castle `AESEngine` 或 `CBCBlockCipher` 加密時，金鑰或 IV 長度不符合要求
- 即使將金鑰填充至 16 bytes（重複、zero-pad、hash 等），輸出仍與 Dart 不一致
- IV 參數可能與金鑰相同（另一個非標準行為的信號）
- 所有捕獲的 ciphertext 共享相同的固定前綴，但變數區域無法匹配

## 診斷方式

1. **檢查金鑰長度**：Frida 顯示的 key 參數長度是否為 8 bytes？
2. **檢查 IV 是否等於金鑰**：如果 IV 參數與金鑰相同，表示實作使用非標準 nonce/IV 建構
3. **檢查輸出長度**：44 chars Base64 = 32 decoded bytes = 29 bytes 明文 + 3 bytes PKCS7 padding → 這是 AES-CBC 的行為，不是 GCM
4. **嘗試標準填充方式**：重複填充、zero-pad、MD5、SHA-256、串接、XOR、反轉——如果全部失敗，表示金鑰衍生方式非標準
5. **反組譯 Dart AOT 函式**：使用 `unflutter` 或 `blutter` 反組譯目標 offset，找出實際的金鑰擴展邏輯

## 解法（原始 — 已過時）

> 以下解法在 2026-05-14 之前有效。現在已知正確解法為 **AES/CTR/NoPadding + PKCS7 padding + zero-padded key/IV**。

當 8-byte 金鑰無法以標準 AES 重現時：

1. **不要浪費時間猜測金鑰衍生方式**——標準填充/雜湊方法幾乎不可能猜中
2. **反組譯 Dart AOT 函式**——使用 `unflutter` 或 `blutter` 反組譯目標 offset，找出實際的金鑰擴展邏輯
3. **或新增 Frida hook 追蹤金鑰擴展**——在 `aesEncryptFallback` 內部 hook `Uint8List` 寫入操作，觀察金鑰如何被擴展
4. **考慮混合方案**——如果金鑰衍生無法短期解決，可以捕獲固定前綴，只自訂產生變數區域
5. **記錄所有測試過的衍生方式**——避免未來重複測試相同的組合

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

## 受影響的檔案

- Java 實作中的 `Cipher.getInstance()` 參數選擇和 `SecretKeySpec` 金鑰規格
- Frida hook scripts 中對金鑰參數的解讀

## 相關 atoms

- `intelligence/engineering/analytical-reasoning/heuristics/hook-selection.md`
- `intelligence/engineering/analytical-reasoning/failure/processBlock-count-ambiguity.md`
- `feedback/history/apk-analysis/flutter-dart-aot/2026-05-14_000000-dart-custom-aes-8byte-key-not-reproducible.md`

## Token 影響

低。此 atom 在 Dart 自訂 AES 實作分析 session 中 lazy-load，約 200-300 tokens。
