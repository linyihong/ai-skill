# JavaScript Bitwise 64-bit Truncation（JavaScript 位元運算子截斷 64-bit 指標）

## 問題

在 64-bit Android 上寫 Frida hook 讀取 Dart AOT heap 物件時，`untagHeapPtr` 函數使用 JavaScript 位元運算子（`& ~7`）清除指標低 3 位元，但 JavaScript 的位元運算子會將數字轉為 32-bit 有號整數，導致 64-bit 位址的高 32 bits 被靜默丟棄。

## 症狀

| 症狀 | 發生時機 | 頻率 |
|------|---------|------|
| 所有 `tryReadDartString` 回傳 `undecoded` | 每次 hook 觸發 | 100% |
| `hexdump` 報 `access violation`，位址缺少高 32 bits | 嘗試讀取記憶體時 | 100% |
| 指標看起來有效（`isLikelyPtr` 回傳 true）但讀取失敗 | 每次指標運算 | 100% |

## 診斷方式

```javascript
// 檢查 untagHeapPtr 回傳值是否被截斷
var p = ptr("0x6b02303719");
var broken = ptr(p & ~7);        // 錯誤：0x02303718（缺少 0x6b0 前綴）
console.log("Broken:", broken.toString());
```

如果位址缺少高 32 bits，就是這個問題。

## 解法

見 [`intelligence/engineering/analytical-reasoning/heuristics/javascript-bitwise-64bit-truncation.md`](../heuristics/javascript-bitwise-64bit-truncation.md) — 包含正確的 `untagHeapPtr` 實作、`NativePointer` API 用法、以及決策表。

## 受影響的檔案

所有使用 JavaScript 位元運算子做指標運算的 Frida hook script 都應檢查：

- `hook_eh_generation_v7_extract_key_iv.js`（已修正）
- `hook_dart_request_interceptor.js`（可能有相同問題）
- `hook_self_generation_phase1.js`（可能有相同問題）
- `hook_self_generation_phase1_v10.js`（可能有相同問題）
- `hook_skyshield_getSignProxy.js`（可能有相同問題）

## 相關 atoms

- [`intelligence/engineering/analytical-reasoning/heuristics/javascript-bitwise-64bit-truncation.md`](../heuristics/javascript-bitwise-64bit-truncation.md)（解法與預防規則）
- `feedback/history/apk-analysis/flutter-dart-aot/2026-05-13_060300-javascript-bitwise-32bit-truncation-dart-64bit-ptr.md`（原始 lesson）

## Token 影響

低。此 atom 在 `undecoded` 或 access violation 發生時 lazy-load，約 100-150 tokens。
