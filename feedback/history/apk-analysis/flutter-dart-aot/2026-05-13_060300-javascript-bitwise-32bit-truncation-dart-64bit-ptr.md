### 2026-05-13 - JavaScript 位元運算子在 Frida Hook 中截斷 64-bit 指標

#### 一句話總結
JavaScript 的位元運算子（`&`, `|`, `~`, `^`）會將數字轉為 32-bit 有號整數，導致 Dart AOT 的指標運算（如 `untagHeapPtr`）靜默截斷 64-bit 位址，造成記憶體讀取錯誤。

#### 人類可讀說明
在 64-bit Android 上寫 Frida hook 讀取 Dart AOT heap 物件時，常需要清除指標的低 3 位元（untag），例如 `ptr & ~7`。但在 JavaScript 中，`0x6b02303719 & ~7` 並不會得到 `0x6b02303718`——JavaScript 會先把兩個運算元轉成 32-bit 有號整數，結果變成 `0x02303718`（高 32 bits 被靜默丟棄）。這導致後續所有記憶體讀取都存取到錯誤位址，回傳垃圾資料或觸發 access violation。

修正方式：避免使用 JavaScript 位元運算子，改用字串操作或 Frida 的 `NativePointer` API（`ptr.and()`, `ptr.shr()` 等）。

#### 觸發條件
- 寫 Frida hook 讀取 Dart AOT heap 物件
- `untagHeapPtr` 函數使用 JavaScript 位元運算子（如 `& ~7`）
- 所有 `tryReadDartString` 都回傳 `undecoded`，即使指標看起來有效
- `hexdump` 顯示 "access violation"，位址缺少高 32 bits

#### 證據
- `hexdump failed: Error: access violation accessing 0x2303718` — 位址 `0x2303718` 缺少 `0x6b0` 前綴，正確應為 `0x6b02303718`
- `untagHeapPtr(0x6b02303719, ctx)` 回傳 `0x02303718` 而非 `0x6b02303718`
- 根因：`0x6b02303719 & ~7` 在 JavaScript 中 = `0x02303718`（32-bit 截斷）
- 修正方式：改用字串操作清除低位元

#### 通用教訓
1. **絕對不要在 64-bit 值上使用 JavaScript 位元運算子**——它們永遠會截斷成 32-bit
2. Dart AOT 指標運算在 Frida 中應使用 `NativePointer` 方法：`ptr.and()`, `ptr.shr()`, `ptr.add()` 等
3. 如果必須用純 JS，使用 `BigInt`（`BigInt(ptr) & ~BigInt(7)`）或手動解析 hex 字串
4. 當所有 `tryReadDartString` 都回傳 `undecoded` 時，先懷疑指標運算錯誤，而不是字串讀取邏輯
5. 務必用 `hexdump` 或 `console.log(ptr.toString())` 確認位址正確

#### Agent 行動
- 發現所有 `tryReadDartString` 都回傳 `undecoded`，即使指標看起來有效
- 用 `hexdump` 檢查實際讀取的記憶體——發現 access violation，位址被截斷
- 追蹤到 `untagHeapPtr` 使用 JavaScript `&` 運算子造成截斷
- 修正：重寫 `untagHeapPtr` 使用字串操作：
  ```javascript
  function untagHeapPtr(p, ctx) {
    var s = p.toString(16);
    var len = s.length;
    var lastChar = s[len - 1];
    var lastDigit = parseInt(lastChar, 16);
    var cleared = lastDigit & ~7;
    var newLast = cleared.toString(16);
    var newS = s.substring(0, len - 1) + newLast;
    return new NativePointer('0x' + newS);
  }
  ```
- 修正後所有 Dart 字串讀取都正確

#### 目標 / 行動 / 驗證
- **目標**：在 Frida 中正確讀取 Dart AOT heap 字串
- **行動**：用字串操作取代 `untagHeapPtr` 中的位元運算子
- **驗證**：修正後 `tryReadDartString` 正確回傳 `[api/public/?service=Plaza.heartbeat&l=zh-cn]` 而非 `undecoded`

#### 適用時機
- 寫 64-bit Dart AOT 應用的 Frida hook
- 任何需要清除/設定特定位元的指標運算
- 在可能 > 32-bit 的值上使用 `&`, `|`, `~`, `^` 運算子

#### 不適用時機
- 僅處理 32-bit 位址（現代 Android 很少見）
- 使用 Frida 的 `NativePointer` 方法（`ptr.and()`, `ptr.or()` 等，它們正確處理 64-bit）
- 使用 BigInt 運算

#### 驗證方式
- 修正前：`untagHeapPtr(0x6b02303719)` → `0x02303718`（錯誤，32-bit 截斷）
- 修正後：`untagHeapPtr(0x6b02303719)` → `0x6b02303718`（正確）

#### Promotion Target
- `intelligence/engineering/analytical-reasoning/heuristics/javascript-bitwise-64bit-truncation.md`（已 promote）

#### Required Linked Updates
- `feedback/history/apk-analysis/README.md`：`flutter-dart-aot/` 計數從 13 → 14（已更新）
- `intelligence/engineering/analytical-reasoning/heuristics/README.md`：新增 atom 列（已更新）
- 所有現有 hook script 中的 `untagHeapPtr` 實作都應檢查此 bug
- 受影響的 script：`hook_eh_generation_v7_extract_key_iv.js`（已修正）、`hook_dart_request_interceptor.js`（可能有相同問題）、`hook_self_generation_phase1.js`（可能有相同問題）
