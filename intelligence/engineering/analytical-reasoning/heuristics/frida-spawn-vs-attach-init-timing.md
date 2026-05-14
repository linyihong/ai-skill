# Frida spawn 與 attach 的初始化時機

`intelligence/engineering/analytical-reasoning/heuristics/frida-spawn-vs-attach-init-timing.md`

## 問題

Frida attach 模式（`frida -U <package>`）無法 hook App 啟動階段的初始化函數（static initializer、singleton constructor、`JNI_OnLoad`、library constructor）。必須使用 spawn 模式（`frida -U -f <package>`）。

此外，Frida 的 JavaScript runtime 沒有 Node.js 的 `Buffer` API，操作二進位資料需使用 `NativePointer` 方法。

## 啟發式

### Spawn 與 attach 的選擇

| 目標函數時機 | 建議模式 | 理由 |
|----------------------|-----------------|--------|
| App 啟動階段（static initializer、constructor） | **Spawn**（`-f`） | Attach 模式錯過初始化 |
| 使用者操作階段（button click、API call） | Attach | 更輕量，不重啟 App |
| 不確定執行時機 | **先試 spawn** | 如果 spawn 會觸發但 attach 不會，就是初始化函數 |
| 需要捕獲完整 lifecycle | Spawn | 從頭開始捕獲所有事件 |

### Frida JavaScript 限制

| Node.js API（對照用） | Frida 等效寫法 |
|------------|-----------------|
| `Buffer.from()` | `ptr.readByteArray(length)` |
| `Buffer.alloc()` | `Memory.alloc(size)` |
| `uint8[i]` | `ptr.add(i).readU8()` |
| `uint8[i] = val` | `ptr.add(i).writeU8(val)` |
| `require('buffer')` | ❌ 不可用 |

## 決策表

| 條件 | 行動 |
|-----------|--------|
| Hook 在 attach 模式下從未觸發 | 先用 spawn 模式確認 offset 是否正確 |
| 需要操作二進位資料 | 使用 `ptr.add(i).readU8()` 而非 `Buffer` |
| 需要捕獲初始化流程 | 使用 spawn 模式 |
| 只需要 hook 執行階段函數 | 使用 attach 模式（較輕量） |

## 來源

- `feedback/history/apk-analysis/common/2026-05-14_073700-frida-spawn-vs-attach-init-timing-no-buffer.md`
- `feedback/history/apk-analysis/flutter-dart-aot/2026-05-14_081500-initialize-hook-does-fire-in-spawn.md`
