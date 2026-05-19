# Dart AES Short Parameter Misread（短參數誤判為自訂 AES）

> 此 atom 是 failure 入口：記錄「Frida 看到短 key/IV 後，誤判為不可重現的自訂 AES」這類症狀。可重用解法與模式判斷放在 [`../heuristics/dart-encrypt-package-mode-detection.md`](../heuristics/dart-encrypt-package-mode-detection.md)。

## 問題

Frida hook 攔截到 Dart AOT 加密函式時，只看到短 key/IV 參數（例如 8 bytes），導致 agent 誤判為「自訂 AES」或「標準 AES 無法重現」。實際上，短參數可能只是呼叫鏈中間值；真正傳入 cipher 前可能已被 zero-pad、轉成固定長度 buffer，或搭配非直覺的 mode / padding 組合。

## 症狀

- Frida capture 顯示 AES key/IV 長度不是 16/24/32 bytes
- 嘗試以 `Cipher.getInstance("AES/GCM/NoPadding")` 或 `Cipher.getInstance("AES/CBC/PKCS5Padding")` 加密時，`SecretKeySpec` 需要 16/24/32 bytes
- 嘗試以 Bouncy Castle `AESEngine` 或 `CBCBlockCipher` 加密時，金鑰或 IV 長度不符合要求
- 初始 CBC/GCM 嘗試失敗，但 ciphertext block boundary 或前綴仍顯示可重現跡象
- IV / nonce 參數看似非標準，或與 hook 到的中間值關係不清楚

## 診斷方式

1. **確認 hook 點是否為最終 cipher input**：短 key/IV 可能是呼叫者傳入的原始 material，不是 cipher 收到的實際 buffer。
2. **檢查參數 normalization**：追蹤是否有 zero-pad、truncate、copyRange、Uint8List allocation 或 IV/nonce 建構。
3. **用輸出長度反推 mode / padding**：例如 output 是否固定為 block size 倍數、是否有 auth tag、是否需要先 pad 再 encrypt。
4. **用 live proxy / fixture test 驗證 mode**：不要只靠 block count 或 constructor 名稱下結論。
5. **若仍無法重現，再反組譯 Dart AOT 函式**：找出實際 key expansion、mode selection 與 padding path。

## 可重現模式

已觀察到的可重用模式：

| 誤判 | 實際可能性 | 驗證方式 |
| --- | --- | --- |
| 短 key 代表自訂 AES | 呼叫者先 zero-pad 至 16 bytes 再進 AES-128 | Hook normalization 後的 buffer 或比對 fixture |
| CTR 不需要 padding，所以不會有 block-size output | 實作可能先 PKCS7 pad，再用 CTR/SIC encrypt | 比對 plaintext length、ciphertext length 與 padding bytes |
| GCM constructor 出現就代表同一組 encryption group | GCM 可能屬於不同用途或不同 header group | 用 call stack、參數來源與輸出 sink 分組 |
| CBC/GCM 失敗就代表標準 AES 不可重現 | CTR/SIC + normalization 仍可能完全重現 | 依 heuristic 決策表測 CBC、CTR/SIC、GCM |

## 安全範例

```java
byte[] key16 = padZero("<key-material-redacted>".getBytes("UTF-8"), 16);
byte[] seed16 = padZero("<iv-material-redacted>".getBytes("UTF-8"), 16);

Cipher cipher = Cipher.getInstance("AES/CTR/NoPadding");
cipher.init(Cipher.ENCRYPT_MODE, new SecretKeySpec(key16, "AES"), new IvParameterSpec(seed16));

// PKCS7 pad plaintext to 16-byte boundary first
byte[] paddedPt = pkcs7Pad(pt, 16);
byte[] result = cipher.doFinal(paddedPt);
String b64 = Base64.getEncoder().encodeToString(result);
```

## 關鍵教訓

1. Frida hook 的參數不一定是 cipher 的最終輸入；呼叫者可能在傳入前後做 normalization。
2. 短 key/IV 不必然代表自訂 crypto；可能是標準 AES 前的 material preparation。
3. CTR/SIC + PKCS7 padding 雖不常見，但在 App 自訂流程中可能成立。
4. 不要只測 CBC/GCM；若輸出長度與 block pattern 支援，應把 CTR/SIC 納入候選。
5. 第一個 block 或 keystream prefix 匹配是重要信號，但仍需 fixture 或 live proxy test 驗證。

## 相關 atoms

- [`../heuristics/dart-encrypt-package-mode-detection.md`](../heuristics/dart-encrypt-package-mode-detection.md)
- [`processBlock-count-ambiguity.md`](processBlock-count-ambiguity.md)
- [`../heuristics/hook-selection.md`](../heuristics/hook-selection.md)
- `feedback/history/apk-analysis/flutter-dart-aot/2026-05-14_000000-dart-custom-aes-8byte-key-not-reproducible.md`

## Token 影響

低。此 atom 在 Dart 自訂 AES 實作分析 session 中 lazy-load，約 200-300 tokens。
