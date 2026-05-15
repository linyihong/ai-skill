> 遵守 [共用規則索引](../../../../enforcement/README.md) 與 [feedback-lessons](../../../../enforcement/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-15 - Dart AOT `string_refs.jsonl` PP_peep 可直接揭露硬編碼 AES 金鑰

#### One-line Summary

在 `unflutter` 產生的 `string_refs.jsonl` 中搜尋 `kind:"PP_peep"` 的條目，可直接找到函式 pool 中載入的硬編碼字串常數（如 AES key、secret、header name），無需反組譯或動態 hook。

#### Human Explanation

Dart AOT 編譯器將函式內使用的字串常數存放在該函式的 **pool area**（PC 相對定址區）。當函式需要某個字串時，會透過 `LDR X16, [X16, #offset]` 從 pool 載入。`unflutter` 的 `string_refs.jsonl` 會將這些 pool 載入記錄為 `kind:"PP_peep"`，並直接輸出字串的實際值。

這意味著：**不需要反組譯整個函式**，只需要在 `string_refs.jsonl` 中搜尋可疑的關鍵字（如 `encrypt`、`decrypt`、`key`、`secret`、header 名稱、base64 字串），就能直接找到硬編碼的密碼學金鑰或常數。

#### Trigger

- 已使用 `unflutter` 或類似工具產生 `string_refs.jsonl`。
- 正在分析 Dart AOT 函式中的加密/解密邏輯。
- 懷疑金鑰是硬編碼在函式內的常數字串（而非執行期動態產生）。

#### Evidence

- 工具：`unflutter` → `string_refs.jsonl`，`grep` / `jq` 過濾 `kind:"PP_peep"`。
- 去敏範例：在 `RequestInterceptor.onResponse` 的 `string_refs.jsonl` 條目中，`kind:"PP_peep"` 直接揭露了 AES-128 金鑰字串（16 字元 ASCII），該金鑰同時用於請求加密 header 與回應解密。
- 對應的反組譯確認：函式 disassembly 中確實有 `LDR X16, [X16, #offset]` 從 pool 載入該字串，然後傳遞給 `AES(AES._128048(...))` 建構子。
- 證據路徑：專案私有 capture 保留 raw `string_refs.jsonl`；本 lesson 只儲存通用結構。

#### Generalized Lesson

When analyzing Dart AOT functions for hardcoded cryptographic keys or string constants, search `string_refs.jsonl` for entries with `kind:"PP_peep"` (pool peep). These entries contain the actual string value loaded by the function, making them a fast and reliable source for finding:

- AES / HMAC / signing keys
- Secret strings and tokens
- Header names and API paths
- Base64-encoded constants

Cross-reference the found PC with the function's disassembly to confirm the load instruction (`LDR X16, [X16, #offset]`) and trace how the string is used (passed to constructor, compared, etc.).

#### Agent Action

1. After running `unflutter` (or equivalent AOT parser), locate `string_refs.jsonl`.
2. Search for `PP_peep` entries containing keywords like `encrypt`, `decrypt`, `key`, `secret`, `aes`, `hmac`, `sign`, or known header names.
3. For each candidate, note the `func` (function name), `pc` (program counter), and `value` (the actual string).
4. Read the function's disassembly around the candidate PC to confirm the load instruction and trace the string's usage.
5. If the string is a cryptographic key, verify by:
   - Checking if the same key appears in other functions (cross-function reuse).
   - Checking if the key length matches the expected algorithm (16 bytes = AES-128, 32 bytes = AES-256, etc.).
   - Confirming via dynamic Frida hook that the key parameter matches.
6. Document the key in project-private analysis files; do NOT include raw keys in reusable lessons (per sanitization rules).

#### Goal / Action / Validation

- **Goal**: Find hardcoded cryptographic keys in Dart AOT functions without full disassembly.
- **Action**: `grep '"PP_peep"' string_refs.jsonl | grep -i 'encrypt\|key\|secret\|aes'` or use `jq` to filter.
- **Validation**: Cross-reference the PC in the function's disassembly to confirm the `LDR` instruction and trace the key's usage path.

#### Applies When

- You have `string_refs.jsonl` from `unflutter` or a compatible AOT parser.
- The target function uses hardcoded string constants (not runtime-generated keys).
- You need to quickly identify cryptographic material across many functions.

#### Does Not Apply When

- Keys are derived at runtime (e.g., from key derivation functions, DH exchange, or external storage).
- The AOT parser did not produce `string_refs.jsonl` or the pool peep analysis is incomplete.
- You need to understand the full encryption flow, not just find the key.

#### Validation

- The string value from `PP_peep` matches the actual runtime value (confirmed via Frida `onEnter` parameter capture).
- The PC in `string_refs.jsonl` corresponds to a `LDR X16, [X16, #offset]` instruction in the disassembly.
- The same key appears in the expected algorithm context (e.g., passed to `AES._128048` constructor).

#### Promotion Target

`intelligence/engineering/analytical-reasoning/heuristics/` — this is a general-purpose heuristic for Dart AOT reverse engineering that applies across projects.

#### Required Linked Updates

- [`intelligence/engineering/analytical-reasoning/heuristics/README.md`](../../../../intelligence/engineering/analytical-reasoning/heuristics/README.md): Add entry for this heuristic.
- [`feedback/history/apk-analysis/README.md`](../README.md): Increment `flutter-dart-aot/` counter.
