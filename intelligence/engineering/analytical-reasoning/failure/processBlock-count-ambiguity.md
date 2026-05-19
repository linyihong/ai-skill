# processBlock Count Ambiguity（processBlock 呼叫次數歧義）

## 問題

Frida hook 攔截到 `PaddedBlockCipher.processBlock()` 被呼叫 43 次時，無法單獨從 block count 區分加密模式。43 blocks 可對應多種 mode：

| 模式 | 解釋 | 說明 |
|------|------|------|
| **CBC + PKCS7** | 42 data blocks + 1 padding block | 42×16=672 bytes 明文，PKCS7 補滿 1 block |
| **GCM** | 1 J0 block + 42 CTR mode blocks | GCM 內部先計算 J0（1 block），再用 CTR 模式加密 42 blocks |
| **CTR(SIC) + PKCS7** | 42 data blocks + 1 padding block | 與 CBC 相同 block count，但 mode 不同 |

## 症狀

- Frida capture 顯示 `PBC.ctor`（PaddedBlockCipher）被呼叫
- `processBlock` 被呼叫 43 次
- 無法從 block count 推斷 mode
- 如果 `GCMBlockCipher.ctor` 也在 capture 中出現，可能屬於**不同的 encryption group**（非 `eh` header 加密），不應混淆

## 診斷方式

1. **檢查 constructor chain**：確認 `PBC.ctor` 的參數是否包含 mode 資訊（Dart `encrypt` package 的 `AES()` constructor 接受 `AESMode` 參數）
2. **檢查 IV 長度**：CBC 用 16 bytes IV，GCM 用 12 bytes nonce
3. **執行 fixture 或授權 live proxy test**（最可靠）：用不同 mode 加密相同明文，比對 decrypt parity、response status class 或明確 pass/fail signal。不要把某次 App 的具體 status code 寫進 reusable atom。
4. **檢查 output 長度**：CBC output 是 16 的倍數（含 padding），CTR output 長度等於明文

## 解法

當 block count 無法確定 mode 時：

1. 先按 call group 分離不同 encryption 用途。
2. 依 [`../heuristics/dart-encrypt-package-mode-detection.md`](../heuristics/dart-encrypt-package-mode-detection.md) 檢查 mode、padding、IV/nonce 長度與 output length。
3. 用 fixture 或授權 live proxy test 驗證候選 mode。
4. 注意 `GCMBlockCipher.ctor` 可能屬於不同的 encryption group，不要混用。

## 相關 atoms

- `intelligence/engineering/analytical-reasoning/heuristics/hook-selection.md`
- `intelligence/engineering/analytical-reasoning/heuristics/dart-encrypt-package-mode-detection.md`
- `intelligence/engineering/analytical-reasoning/evidence-first-routing.md`
- `feedback/history/apk-analysis/flutter-dart-aot/2026-05-13_232600-dart-encrypt-package-aes-mode-cbc-not-ctr.md`

## Token 影響

低。此 atom 在 Dart `encrypt` package 分析 session 中 lazy-load，約 200-300 tokens。
