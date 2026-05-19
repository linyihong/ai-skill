# Dart Encrypt Package Mode Detection（Dart 加密模式判斷啟發式）

## 問題

Dart AOT 分析中，constructor chain、`processBlock` 次數或短 key/IV 參數常不足以單獨判斷實際加密模式。若只靠單一信號，容易把 CBC、CTR/SIC、GCM 或參數 normalization 誤判成自訂 crypto。

## 原則

- Constructor 名稱是信號，不是結論；同一 capture 可能有多個 encryption group。
- `processBlock` count 可能同時符合 CBC + padding、CTR/SIC + padding 或 GCM internal blocks。
- Frida hook 到的 key/IV 可能是 normalization 前的 material，不一定是 cipher 最終輸入。
- 模式判斷必須用 output length、padding、IV/nonce 長度、call group 與 fixture / live proxy test 交叉驗證。

## 決策表

| 信號 | 判斷 | 下一步 |
| --- | --- | --- |
| `PaddedBlockCipher` constructor 出現，但只有 block count | 不足以區分 CBC / CTR(SIC) + padding | 檢查 inner cipher、IV 長度與 output 長度 |
| `GCMBlockCipher` constructor 出現 | 可能是 GCM，也可能是不同 encryption group | 用 call stack、參數來源與輸出 sink 分組 |
| Output 是 block size 倍數，但候選是 CTR/SIC | 可能是先 PKCS7 pad 再 CTR encrypt | 比對 plaintext length 與 padding bytes |
| Frida 顯示短 key/IV | 不直接判定為自訂 AES | 追蹤 zero-pad、copy、truncate 或 Uint8List normalization |
| CBC/GCM fixture 失敗 | 不代表標準 AES 不可重現 | 加入 CTR/SIC + normalization 候選 |
| 多個 mode 都可能成立 | 靜態信號不足 | 用 fixture 或授權 live proxy test 驗證 response status / decrypt parity |

## 驗證順序

1. 先按 call group 分離不同 encryption 用途，不混用 constructor 信號。
2. 記錄 key/IV material 的長度與 normalization 前後狀態，值本身留在專案 evidence。
3. 用 output 長度與 padding pattern 篩掉不可能的 mode。
4. 依序測試 CBC、CTR/SIC、GCM 等候選組合；每個候選只記錄 mode、padding、normalization 與 pass/fail。
5. 若 fixture 可重現，再回填到專案 decoder；若不可重現，回到 Dart AOT disassembly 或更高語意 hook。

## 不建議的做法

- 不要只因為短 key/IV 就宣稱「自訂 AES」。
- 不要只因為 `processBlock` 次數相同就判定 mode。
- 不要把 raw key、IV、ciphertext 或完整 response 寫進 reusable atom。
- 不要把某次 App 的 HTTP status 或 payload 當成全域規則；只記錄驗證方法。

## 相關 atoms

- [`../failure/processBlock-count-ambiguity.md`](../failure/processBlock-count-ambiguity.md)
- [`../failure/custom-dart-aes-8byte-key-not-reproducible.md`](../failure/custom-dart-aes-8byte-key-not-reproducible.md)
- [`hook-selection.md`](hook-selection.md)

## Token 影響

低。此 atom 在 Dart 加密模式判斷時 lazy-load，約 150-250 tokens。
