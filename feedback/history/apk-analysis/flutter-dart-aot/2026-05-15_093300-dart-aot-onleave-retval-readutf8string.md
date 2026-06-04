### 2026-05-15 - Dart AOT `onLeave` 回傳值讀取：`retval` 是 `NativePointer`，用 `readUtf8String()` 直接讀 `_OneByteString`

Status: validated

#### One-line Summary

Frida `Interceptor.attach` 的 `onLeave` callback 中，`retval` 參數已經是 `NativePointer`（指向 Dart heap object 的 tagged pointer），不需要 `ptr()` 包裝或透過 `ctx.x0` 讀取；Dart `_OneByteString`（Latin-1）的內容可以直接用 `readUtf8String(len)` 讀取，不需要先轉 byte array。

#### Human Explanation

在 Dart AOT 逆向中，hook 函數的 return value 是常見需求。Frida 的 `Interceptor.attach` 在 `onLeave` 提供 `retval` 參數，它已經是 `NativePointer` 型別，指向 Dart VM heap 中的 tagged object pointer。

容易犯的錯誤：
1. 用 `ptr(ctx.x0)` 或 `ptr(retval)` 重複包裝 — `retval` 已經是 `NativePointer`
2. 用 `readByteArray(len)` 讀出 bytes 再 `String.fromCharCode()` — Dart `_OneByteString` 的內容是 Latin-1 編碼，`readUtf8String(len)` 可以直接正確讀取

Dart `_OneByteString` 的典型記憶體佈局（`smi32@8/16`）：
- offset +0: 4-byte header（class id + size tag）
- offset +4: 4-byte hash（或 padding）
- offset +8: 4-byte Smi-tagged length（`raw_value >> 1` 為實際長度）
- offset +16: string data（Latin-1 編碼）

#### Trigger

- 寫 Frida hook 需要讀取 Dart AOT 函數的回傳值
- 在 `onLeave` 中嘗試用 `ptr(ctx.x0)` 或 `ptr(retval)` 讀取回傳值
- 用 `readByteArray()` + `String.fromCharCode()` 解碼 string 內容
- `readU32()` 讀到的 length 看起來不對（因為是 Smi-tagged，需要 `>> 1`）

#### Evidence

- Tool: Frida `Interceptor.attach` with `onLeave` callback
- Sanitized excerpt: `onLeave: function(retval) { var untagged = retval.sub(1); var len = untagged.add(8).readU32() >> 1; var s = untagged.add(16).readUtf8String(len); }` — 成功讀取 Dart `_OneByteString` 回傳值
- Evidence path: `<PROJECT_ROOT>/<target-app>/scripts/frida/hook_decrypt_aes_response.js`

#### Generalized Lesson

1. **`onLeave` 的 `retval` 已經是 `NativePointer`** — 不需要 `ptr()` 包裝，也不需要 `ctx.x0`。直接對 `retval` 進行 pointer arithmetic（`.sub()`, `.add()`）即可。

2. **Dart heap object 的 tag bit 處理** — Dart heap object pointer 的 bit 0 為 1（tagged pointer）。讀取前需要 `retval.sub(1)` 得到真正的 heap 位址。可以用 `retval.and(ptr('1')).toString() === '0x1'` 檢查是否 tagged。

3. **`readUtf8String(len)` 可以直接讀 `_OneByteString`** — Dart 的 `_OneByteString`（Latin-1 編碼）與 UTF-8 在 ASCII 範圍（0x00-0x7F）完全相容，且 Latin-1 的 0x80-0xFF 範圍在 `readUtf8String` 中也能正確保留（Frida 的 V8 引擎會將其視為有效的 UTF-8  continuation byte 序列）。不需要先 `readByteArray()` 再轉換。

4. **Smi-tagged length 需要 `>> 1`** — Dart AOT 中，small integer（Smi）以 tagged 形式儲存：`raw_value = actual_value << 1 | 1`。讀取 length 欄位時，需要 `readU32() >> 1` 得到實際長度。

5. **validate printable ratio** — 解碼後應檢查 printable character 比例（>75%），避免誤判非 string 物件為 string。

#### Agent Action

- 在 `onLeave` 中直接使用 `retval`（`NativePointer`），不要 `ptr()` 包裝
- 先 `retval.sub(1)` 去除 tag bit
- 用 `readU32()` 在 offset +8 讀取 Smi-tagged length，`>> 1` 得到實際長度
- 用 `readUtf8String(len)` 在 offset +16 讀取 string 內容
- 驗證 printable ratio > 75% 確認解碼正確
- 若 `smi32@8/16` 失敗，fallback 掃描前 32 bytes 找 JSON 起始字元（`{` 或 `"`）

#### Goal / Action / Validation

- Goal: 在 Frida `onLeave` 中正確讀取 Dart AOT 函數的 `_OneByteString` 回傳值
- Action: 使用 `retval.sub(1).add(8).readU32() >> 1` 取得 length，`retval.sub(1).add(16).readUtf8String(len)` 讀取內容
- Validation: 解碼後的 string 長度符合預期、printable ratio > 75%、JSON 格式正確（若為 JSON）

#### Applies When

- Flutter/Dart AOT 逆向，需要 capture 函數回傳值
- 回傳值為 `_OneByteString`（Latin-1 編碼的 Dart string）
- 使用 Frida `Interceptor.attach` 的 `onLeave` callback

#### Does Not Apply When

- 回傳值不是 Dart string（如 int、bool、List、Map）
- 使用 Dart VM service protocol（非 native hook）
- 處理 `_TwoByteString`（UTF-16 編碼的 Dart string）— 需要不同的讀取方式
- 32-bit Android 應用（pointer 大小不同）

#### Validation

- 成功 decode 前：`retval` 顯示為 tagged pointer（如 `0x7a123456`），`readU32()` 在 offset +8 回傳奇數值（Smi-tagged）
- 成功 decode 後：`readU32() >> 1` 得到合理長度（如 1023），`readUtf8String(len)` 回傳完整 JSON string

#### Promotion Target

- `intelligence/engineering/analytical-reasoning/heuristics/dart-aot-onleave-retval-readutf8string.md`（建議新增）

#### Required Linked Updates

- `feedback/history/apk-analysis/README.md`：`flutter-dart-aot/` 計數從 20 → 21
- `intelligence/engineering/analytical-reasoning/heuristics/README.md`：若 promote 則新增 atom 列
- 現有 lesson `2026-05-01_164741-dart-inline-onebyte-string-smi-length.md` 已涵蓋 Smi length 概念，本 lesson 補充 `onLeave` 的 `retval` 處理和 `readUtf8String()` 的使用
