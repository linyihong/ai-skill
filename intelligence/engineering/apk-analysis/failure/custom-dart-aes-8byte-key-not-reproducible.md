# Custom Dart AES 8-byte Key Not Reproducible（Dart 自訂 AES 8-byte 金鑰無法以標準 AES 重現）

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

## 解法

當 8-byte 金鑰無法以標準 AES 重現時：

1. **不要浪費時間猜測金鑰衍生方式**——標準填充/雜湊方法幾乎不可能猜中
2. **反組譯 Dart AOT 函式**——使用 `unflutter` 或 `blutter` 反組譯目標 offset，找出實際的金鑰擴展邏輯
3. **或新增 Frida hook 追蹤金鑰擴展**——在 `aesEncryptFallback` 內部 hook `Uint8List` 寫入操作，觀察金鑰如何被擴展
4. **考慮混合方案**——如果金鑰衍生無法短期解決，可以捕獲固定前綴，只自訂產生變數區域
5. **記錄所有測試過的衍生方式**——避免未來重複測試相同的組合

## 受影響的檔案

- Java 實作中的 `Cipher.getInstance()` 參數選擇和 `SecretKeySpec` 金鑰規格
- Frida hook scripts 中對金鑰參數的解讀

## 相關 atoms

- `intelligence/engineering/apk-analysis/heuristics/hook-selection.md`
- `intelligence/engineering/apk-analysis/failure/processBlock-count-ambiguity.md`
- `feedback/history/apk-analysis/flutter-dart-aot/2026-05-14_000000-dart-custom-aes-8byte-key-not-reproducible.md`

## Token 影響

低。此 atom 在 Dart 自訂 AES 實作分析 session 中 lazy-load，約 200-300 tokens。
