# JavaScript Bitwise 64-bit Truncation Heuristic（Frida Hook 指標運算啟發式）

> 如果遇到 `undecoded` 或 `hexdump` access violation，先看 [`failure/javascript-bitwise-64bit-truncation.md`](../failure/javascript-bitwise-64bit-truncation.md) 確認症狀。

## 問題

在 64-bit Android 上寫 Frida hook 讀取 Dart AOT heap 物件時，JavaScript 的位元運算子（`&`, `|`, `~`, `^`）會將數字轉為 32-bit 有號整數，導致指標運算靜默截斷 64-bit 位址。

## 原則

1. **絕對不要在 64-bit 值上使用 JavaScript 位元運算子**——它們永遠會截斷成 32-bit
2. Dart AOT 指標運算在 Frida 中應使用 `NativePointer` 方法：`ptr.and()`, `ptr.shr()`, `ptr.add()` 等
3. 如果必須用純 JS，使用 `BigInt`（`BigInt(ptr) & ~BigInt(7)`）或手動解析 hex 字串
4. 當所有 `tryReadDartString` 都回傳 `undecoded` 時，先懷疑指標運算錯誤，而不是字串讀取邏輯
5. 務必用 `hexdump` 或 `console.log(ptr.toString())` 確認位址正確

## 決策表

| 情境 | 建議做法 | 判斷信號 |
|------|---------|---------|
| 需要清除指標低 3 位元（untag） | 字串操作或 `ptr.and()` | `ptr & ~7` 在 JS 中截斷 |
| 需要設定/清除特定位元 | `NativePointer.and()`, `.or()`, `.shr()` | 位址缺少高 32 bits |
| 所有 `tryReadDartString` 回傳 `undecoded` | 先檢查 `untagHeapPtr` 實作 | `hexdump` 顯示 access violation |
| 需要跨 32-bit 邊界的位元遮罩 | `BigInt` 運算 | 值 > `0x7FFFFFFF` |

## 不建議的做法

- 不要用 `ptr & ~7` 做 untag — 改用字串操作或 `ptr.and()`
- 不要假設 JavaScript 數字可以安全容納 64-bit 指標
- 不要忽略 `hexdump` 顯示的 access violation 位址

## 驗證方式

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

## 來源

- Feedback lesson: [`feedback/history/apk-analysis/flutter-dart-aot/2026-05-13_060300-javascript-bitwise-32bit-truncation-dart-64bit-ptr.md`](../../../../feedback/history/apk-analysis/flutter-dart-aot/2026-05-13_060300-javascript-bitwise-32bit-truncation-dart-64bit-ptr.md)
