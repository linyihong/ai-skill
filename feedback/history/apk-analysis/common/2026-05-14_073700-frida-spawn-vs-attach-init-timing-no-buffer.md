> 遵守 [共用規則索引](../../../enforcement/README.md)、[dependency-reading](../../../enforcement/dependency-reading.md)、[neutral-language](../../../enforcement/neutral-language.md)、[goal-action-validation](../../../enforcement/goal-action-validation.md) 與 [feedback-lessons](../../../enforcement/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-14 - Frida 實戰：spawn 模式 vs attach 模式的初始化時機、JS 無 Buffer API

Status: validated

#### One-line Summary

Frida attach 模式無法 hook App 啟動階段的初始化函數（如 Dart AOT 的 static initializer），必須使用 spawn 模式（`-f` flag）。Frida 的 JavaScript runtime 沒有 Node.js 的 `Buffer` API，操作二進位資料需使用 `p.add(i).readU8()` 模式。

#### Human Explanation

在 Frida 逆向中有兩個常見陷阱：

1. **Spawn vs Attach 初始化時機**：Frida 的 attach 模式（`frida -U com.example.app`）是在 App 已經啟動後才附加 hook。這意味著任何在 App 啟動階段執行的初始化函數（如 Dart AOT 的 static initializer、native library 的 constructor、`JNI_OnLoad` 等）都已經執行完畢，hook 永遠不會觸發。使用 spawn 模式（`frida -U -f com.example.app`）可以在 App 啟動前就載入 hook 腳本，捕獲完整的初始化流程。

2. **Frida JS 無 `Buffer` API**：Frida 的 JavaScript runtime 是基於 Duktape 或 QuickJS，不是 Node.js。這意味著 `Buffer`、`Uint8Array` 等 Node.js 內建模組不可用。操作二進位資料時，需要使用 Frida 原生的 `NativePointer` 方法：`ptr.add(offset).readU8()` 讀取單個 byte、`ptr.add(offset).writeU8(value)` 寫入單個 byte、`ptr.readByteArray(length)` 讀取 byte array。

#### Trigger

- `Interceptor.attach` 在 attach 模式下從未觸發，但在 spawn 模式下正常觸發
- Frida 腳本報錯 `ReferenceError: 'Buffer' is not defined`
- 嘗試使用 `Buffer.from()` 或 `Uint8Array` 操作二進位資料失敗

#### Evidence

- Tool: Frida hook script comparison（attach vs spawn）
- Sanitized excerpt:
  - `hookInitialize` at offset `0xe5b458` never fires in attach mode
  - Same hook fires correctly in spawn mode（`frida -U -f <package> -l script.js`）
  - `Buffer` is not defined → use `ptr.add(i).readU8()` instead
- Evidence path: `<PROJECT_ROOT>/<target-app>/scripts/frida/hook_dispatch_first_call.js`

#### Generalized Lesson

1. **初始化函數必須用 spawn 模式**——如果目標 hook 在 attach 模式下從未觸發，嘗試 spawn 模式。這適用於所有在 App 啟動階段執行的程式碼。
2. **Frida JS 沒有 `Buffer`**——不要嘗試 `require('buffer')` 或 `Buffer.from()`。使用 Frida 原生的 `NativePointer` 方法操作二進位資料。
3. **Spawn 模式的缺點**：spawn 模式會重新啟動 App，可能觸發登入、session recovery 或 rate limit。如果只需要 hook 執行階段的函數，attach 模式更輕量。
4. **混合策略**：先用 spawn 模式確認初始化流程，再用 attach 模式進行後續的細部 hook。

#### Agent Action

部署 Frida hook 時：

1. **先判斷目標函數的執行時機**——如果是初始化階段（static initializer、constructor、library loading），使用 spawn 模式
2. **如果 attach 模式 hook 不觸發，改用 spawn 模式**——不要假設 offset 錯誤
3. **操作二進位資料時使用 `ptr.add(i).readU8()`**——不要嘗試使用 `Buffer` 或 `Uint8Array`
4. **記錄使用的模式**——在腳本或文件中標明使用 attach 還是 spawn 模式

#### Goal / Action / Validation

- Goal: 正確選擇 Frida 的附加模式，避免初始化 hook 遺漏
- Action: 對初始化函數使用 spawn 模式；對執行階段函數使用 attach 模式
- Validation or reference source: 同一 hook 在 spawn 模式下觸發但在 attach 模式下不觸發

#### Applies When

- 需要 hook App 啟動階段的初始化函數
- Frida attach 模式的 hook 從未觸發
- 編寫 Frida JS 腳本需要操作二進位資料
- 需要捕獲 Dart AOT 的 static initializer 或 native library constructor

#### Does Not Apply When

- 只需要 hook 使用者操作階段的函數（attach 模式即可）
- 使用 Frida Gadget（嵌入模式，初始化時機不同）
- 使用其他語言（Python、Swift）的 Frida binding

#### Validation

- Spawn 模式成功觸發 `hookInitialize`，attach 模式不觸發
- `ptr.add(i).readU8()` 成功讀取二進位資料，`Buffer` 報錯

#### Promotion Target

- `intelligence/engineering/analytical-reasoning/heuristics/` — 新增 heuristic：「Frida spawn vs attach init timing」
- `workflow/apk-analysis/execution-flow.md` — 新增步驟：「Frida 部署策略：spawn vs attach」

#### Required Linked Updates

- 無需連動更新；這是新 lesson。
