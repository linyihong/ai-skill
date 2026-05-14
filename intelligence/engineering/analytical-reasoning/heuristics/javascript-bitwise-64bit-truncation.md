# JavaScript Bitwise 64-bit Truncation & Dart AOT Pointer Field Reading Heuristic（Frida Hook 指標運算啟發式）

> 如果遇到 `undecoded` 或 `hexdump` access violation，先看 [`failure/javascript-bitwise-64bit-truncation.md`](../failure/javascript-bitwise-64bit-truncation.md) 確認症狀。

## 問題

在 64-bit Android 上寫 Frida hook 讀取 Dart AOT heap 物件時，有兩個常見的 64-bit 指標截斷問題：

### 問題 A：JavaScript 位元運算子截斷

JavaScript 的位元運算子（`&`, `|`, `~`, `^`）會將數字轉為 32-bit 有號整數，導致指標運算靜默截斷 64-bit 位址。

### 問題 B：`readU32()` 讀取指標欄位

Dart AOT 物件的某些欄位包含**完整的 64-bit 指標**（8 bytes），必須用 Frida 的 `NativePointer.readPointer()` 讀取。使用 `readU32()` 只會讀到低 32 bits，造成位址截斷和 access violation。

## 原則

1. **絕對不要在 64-bit 值上使用 JavaScript 位元運算子**——它們永遠會截斷成 32-bit
2. Dart AOT 指標運算在 Frida 中應使用 `NativePointer` 方法：`ptr.and()`, `ptr.shr()`, `ptr.add()` 等
3. 如果必須用純 JS，使用 `BigInt`（`BigInt(ptr) & ~BigInt(7)`）或手動解析 hex 字串
4. 當所有 `tryReadDartString` 都回傳 `undecoded` 時，先懷疑指標運算錯誤，而不是字串讀取邏輯
5. 務必用 `hexdump` 或 `console.log(ptr.toString())` 確認位址正確
6. **Dart AOT 物件的指標欄位永遠用 `readPointer()`** — 它回傳 `NativePointer`（完整 64-bit），不是 number
7. **`readU32()` 只適用於 Smi（small integer）欄位**，不適用於指向其他物件的指標欄位
8. 當 `isLikelyPtr()` 對預期是指標的欄位值回傳 `false` 時，先檢查是否用了 `readU32()` 而不是 `readPointer()`
9. 一個安全的啟發式方法：如果欄位值在 `0x00000000`–`0x0000ffff` 範圍內，可能是 Smi；如果值很大（如 `0x7a6b027f4b`），則是指標

## 決策表

| 情境 | 建議做法 | 判斷信號 |
|------|---------|---------|
| 需要清除指標低 3 位元（untag） | 字串操作或 `ptr.and()` | `ptr & ~7` 在 JS 中截斷 |
| 需要設定/清除特定位元 | `NativePointer.and()`, `.or()`, `.shr()` | 位址缺少高 32 bits |
| 所有 `tryReadDartString` 回傳 `undecoded` | 先檢查 `untagHeapPtr` 實作 | `hexdump` 顯示 access violation |
| 需要跨 32-bit 邊界的位元遮罩 | `BigInt` 運算 | 值 > `0x7FFFFFFF` |
| 讀取 Dart 物件欄位，預期是指標 | 用 `readPointer()` 而非 `readU32()` | `hexdump` 顯示 access violation，位址缺少高 bits |
| 讀取 Dart 物件欄位，預期是 Smi | 用 `readU32()` | 值在 `0x00000000`–`0x0000ffff` 範圍內 |

## 不建議的做法

- 不要用 `ptr & ~7` 做 untag — 改用字串操作或 `ptr.and()`
- 不要假設 JavaScript 數字可以安全容納 64-bit 指標
- 不要忽略 `hexdump` 顯示的 access violation 位址
- **不要用 `readU32()` 讀取 Dart 物件的指標欄位** — 改用 `readPointer()`

## 驗證方式

### 問題 A：JavaScript 位元運算子

```javascript
// ❌ 錯誤：32-bit 截斷
function untagHeapPtr(p, ctx) {
  return ptr(p & ~7);  // 0x6b02303719 → 0x02303718
}

// ✅ 正確：字串操作
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

// ✅ 正確：NativePointer API
function untagHeapPtr(p, ctx) {
  return p.and(ptr(0xFFFFFFFFFFFFFFF8));
}
```

### 問題 B：`readU32()` vs `readPointer()`

```javascript
// ❌ 錯誤：只讀 4 bytes，截斷 64-bit 指標
var fieldVal = ptr(objAddr).add(0x38).readU32();
// 0x7a6b027f4b → 0x6b027f4b（缺少高 bits）

// ✅ 正確：讀取完整 8 bytes 指標
var fieldVal = ptr(objAddr).add(0x38).readPointer();
// 0x7a6b027f4b → 0x7a6b027f4b（完整位址）
```

## 來源

- Feedback lesson A: [`feedback/history/apk-analysis/flutter-dart-aot/2026-05-13_060300-javascript-bitwise-32bit-truncation-dart-64bit-ptr.md`](../../../../feedback/history/apk-analysis/flutter-dart-aot/2026-05-13_060300-javascript-bitwise-32bit-truncation-dart-64bit-ptr.md)
- Feedback lesson B: [`feedback/history/apk-analysis/flutter-dart-aot/2026-05-14_181900-dart-aot-object-field-readpointer-vs-readu32.md`](../../../../feedback/history/apk-analysis/flutter-dart-aot/2026-05-14_181900-dart-aot-object-field-readpointer-vs-readu32.md)
