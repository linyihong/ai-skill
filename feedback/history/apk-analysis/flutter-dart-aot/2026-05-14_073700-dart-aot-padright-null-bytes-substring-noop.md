> 遵守 [共用規則索引](../../../enforcement/README.md)、[dependency-reading](../../../enforcement/dependency-reading.md)、[neutral-language](../../../enforcement/neutral-language.md)、[goal-action-validation](../../../enforcement/goal-action-validation.md) 與 [feedback-lessons](../../../enforcement/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-14 - Dart AOT 逆向：`padRight` 使用 null bytes、`substring` 越界不 crash、dispatch hooking 技巧

Status: validated

#### One-line Summary

逆向 Dart AOT 的 `_OneByteString.padRight` 時，不要假設使用標準 Dart 的 space padding——自訂實作可能使用 null bytes（`\0`）。`substring(0, 32)` 在 16-char 字串上可能不 crash（自訂實作優雅處理越界）。hook vtable dispatch 指令（`UBFX` + `LDR`）比 hook 目標函數更能捕獲完整流程。

#### Human Explanation

在逆向 Dart AOT 編譯的 App 時，有三個容易誤判的行為：

1. **`padRight` 使用 null bytes 而非 spaces**：標準 Dart `String.padRight()` 使用空格（0x20）填充。但自訂實作（`_OneByteString.padRight`）可能使用 null bytes（0x00）。這會影響後續的字串處理——null bytes 不是可列印字元，在 hex dump 中看起來像截斷，但 Dart 的 `_OneByteString` 可以正常包含它們。

2. **`substring(0, 32)` 在 16-char 字串上不 crash**：標準 Dart `String.substring()` 在 start > length 時會拋出 `RangeError`。但自訂實作可能優雅處理越界情況，直接回傳原字串而不 crash。這意味著看到 `substring(0, 32)` 被呼叫不代表字串真的有 32 個字元。

3. **Dispatch hooking 技巧**：Dart AOT 的虛擬函式呼叫透過 vtable dispatch 實現。hook dispatch 指令（`UBFX X0, X0, #12, #20` 提取 class ID → `LDR X16, [X21, X0, LSL #3]` 載入 vtable 條目 → `BLR X16` 間接呼叫）可以捕獲所有可能的目標函數，而不需要事先知道具體是哪個函數被呼叫。

#### Trigger

- Frida capture 顯示 `padRight` 回傳的字串包含 `\0` bytes 而非 spaces
- `substring(0, 32)` 被呼叫在 16-char 字串上但沒有 crash
- 嘗試 hook 特定目標函數失敗（offset 錯誤或函數未被呼叫），改用 dispatch 指令 hook 成功捕獲完整流程

#### Evidence

- Tool: Frida AOT hook（`hook_dispatch_first_call.js`）
- Sanitized excerpt:
  - `padRight("l65tvNcw", 16)` → `"l65tvNcw\0\0\0\0\0\0\0\0"`（null bytes，非 spaces）
  - `substring(0, 32)` on 16-char string → returns same 16-char string（no crash）
  - Dispatch at `UBFX X0, X0, #12, #20` + `LDR X16, [X21, X0, LSL #3]` + `BLR X16` successfully captures all vtable targets
- Evidence path: `<PROJECT_ROOT>/TATA/scripts/frida/hook_dispatch_first_call.js`

#### Generalized Lesson

1. **Dart AOT `padRight` 可能使用 null bytes**——不要假設使用標準 Dart space padding。檢查 hex dump 確認填充字元。如果看到 `\0`，這是自訂實作，不是 bug。
2. **Dart AOT `substring` 可能不拋出越界異常**——不要假設 `substring(0, N)` 代表字串長度 >= N。檢查實際字串長度。
3. **Dispatch hooking 比目標 hook 更穩健**——hook vtable dispatch 指令（`UBFX` + `LDR` + `BLR`）可以捕獲所有可能的目標，不需要事先知道具體函數。這在逆向未知流程時特別有用。
4. **先 hook dispatch，再收斂到具體函數**——先用 dispatch hook 捕獲完整流程，確認哪些函數被實際呼叫後，再針對這些函數做精確 hook。

#### Agent Action

逆向 Dart AOT 的字串操作時：

1. **hook `padRight` 時同時輸出 hex dump**——確認填充字元是 spaces（0x20）還是 null bytes（0x00）或其他
2. **hook `substring` 時檢查實際字串長度**——不要假設 `substring(0, N)` 的參數 N 代表字串長度
3. **優先使用 dispatch hooking**——hook `UBFX` + `LDR` + `BLR` 序列捕獲 vtable dispatch，而不是直接 hook 猜測的目標函數
4. **逐步收斂**——先用 dispatch hook 了解完整流程，再用精確 hook 捕獲細節

#### Goal / Action / Validation

- Goal: 正確逆向 Dart AOT 的字串操作和虛擬函式呼叫
- Action: 使用 dispatch hooking 捕獲完整流程；hex dump 確認填充字元；檢查 substring 的實際行為
- Validation or reference source: Frida capture 顯示 padRight 回傳 null bytes、substring 不 crash、dispatch hook 捕獲所有目標

#### Applies When

- 逆向 Dart AOT 編譯的 Flutter App
- 需要 hook `padRight`、`substring` 或其他 `_OneByteString` 方法
- 不確定 vtable dispatch 的目標函數
- 需要捕獲完整的函式呼叫流程

#### Does Not Apply When

- App 使用標準 Dart VM（非 AOT 編譯）
- 使用 Dart 原始碼級別的偵錯工具（如 DevTools）
- 字串操作使用標準 Dart 函式庫（無自訂覆寫）

#### Validation

- Frida capture 確認 `padRight` 回傳 null bytes（hex dump 顯示 `\0`）
- Frida capture 確認 `substring(0, 32)` 在 16-char 字串上回傳原字串
- Dispatch hook 成功捕獲所有 vtable 目標，無遺漏

#### Promotion Target

- `intelligence/engineering/analytical-reasoning/heuristics/` — 新增 heuristic：「Dart AOT string operation assumptions」
- `workflow/apk-analysis/execution-flow.md` — 新增步驟：「Dart AOT dispatch hooking 策略」

#### Required Linked Updates

- 無需連動更新；這是新 lesson。
