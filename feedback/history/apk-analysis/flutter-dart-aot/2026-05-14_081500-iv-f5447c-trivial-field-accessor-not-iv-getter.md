> 遵守 [共用規則索引](../../../enforcement/README.md)、[dependency-reading](../../../enforcement/dependency-reading.md)、[neutral-language](../../../enforcement/neutral-language.md)、[goal-action-validation](../../../enforcement/goal-action-validation.md) 與 [feedback-lessons](../../../enforcement/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-14 - Dart AOT 逆向：短函數（≤12 bytes）很可能是 trivial field accessor，不是預期的 getter

Status: validated

#### One-line Summary

逆向 Dart AOT 時，若一個函數只有 **≤12 bytes**（如 `0x1a2b57c` 僅 12 bytes），且回傳值是一個 **空的 Dart container object**（class tag `0x010d011c`，4 個 null fields），則該函數是 **trivial field accessor**（直接從物件偏移讀取欄位回傳），不是預期的業務 getter（如 IV getter）。不要浪費時間分析它——真正的邏輯在呼叫者或另一個函數中。

#### Human Explanation

在分析某 App 的 `x-aspnet-version` token 生成時，`flutter_meta.json` 中有一個名為 `IV_f5447c` 的函數（offset `0x1a2b57c`），名稱暗示它是 IV getter。但 Frida hook 顯示：

1. 函數只有 **12 bytes**（ARM64 的 `LDR X0, [X0, #offset]` + `RET`）
2. 回傳值是一個 **空的 Dart container object**（class tag `0x010d011c`，4 個欄位都是 `81 80 00 00` = null）
3. 無論何時呼叫，回傳值都相同——永遠是空的 container

這表示 `IV_f5447c` 只是一個 **trivial field accessor**——它從某個物件讀取一個欄位並回傳。該欄位在初始化時被設為一個空的 container，且**從未被填入資料**。真正的 IV 資料在另一個地方被使用。

**關鍵教訓**：`flutter_meta.json` 的函數名稱是從 Dart 編譯符號推導的，不一定反映實際用途。短函數（≤12 bytes）幾乎總是 trivial field accessor，不是業務邏輯。

#### Trigger

- `flutter_meta.json` 中的函數名稱暗示某個用途（如 `IV_*`），但 Frida hook 顯示函數極短（≤12 bytes）
- 回傳值是一個空的 Dart container object（class tag `0x010d011c`）
- 無論何時呼叫，回傳值都相同（空的 container）
- 同一個 offset 在另一個 hook script 中被命名為不同的名稱（如 `keyExpansion2`）

#### Evidence

- Tool: Frida AOT hook（`hook_capture_iv.js`）
- Sanitized excerpt:
  - `IV:enter | x1=0x6b02bf88f9` → `IV:return | ret=0x6b02bf8939 raw32=1c 01 0d 01 00 00 00 00 81 80 00 00 81 80 00 00 81 80 00 00 81 80 00 00 00 00 00 00 00 00 00 00`
  - Class tag `0x010d011c` = empty container, 4 null fields (`81 80 00 00`)
  - Function size: 12 bytes (ARM64 `LDR X0, [X0, #0x10]` + `RET`)
- Evidence path: `<PROJECT_ROOT>/capture/frida_capture_iv_20260514.log`

#### Generalized Lesson

1. **`flutter_meta.json` 的函數名稱不可靠**——名稱是從 Dart 編譯符號推導的，可能誤導。短名稱（如 `IV_*`）不保證是預期的 getter。
2. **短函數（≤12 bytes）幾乎總是 trivial field accessor**——ARM64 上 12 bytes 只能容納 `LDR X0, [X0, #offset]` + `RET`。這種函數直接從物件偏移讀取欄位，不做任何運算。
3. **空的 Dart container object 是初始化未完成的信號**——class tag `0x010d011c` + 4 null fields 表示 container 已分配但從未被填入資料。真正的資料在另一個地方被設定。
4. **交叉比對不同 hook script 的命名**——同一個 offset 在不同 script 中可能被命名為不同名稱（如 `IV_f5447c` vs `keyExpansion2`）。名稱衝突是 red flag，需要進一步調查。
5. **不要浪費時間分析 trivial field accessor**——直接 hook 它的呼叫者或尋找真正的資料來源。

#### Agent Action

分析 Dart AOT 函數時：

1. **先檢查函數大小**——如果函數只有 ≤12 bytes，它幾乎肯定是 trivial field accessor
2. **檢查回傳值的 class tag**——`0x010d011c` = empty container，`0x010d011c` + null fields = 未初始化的 container
3. **交叉比對名稱**——同一個 offset 在不同來源（`flutter_meta.json`、手動命名、其他 hook script）的名稱是否一致
4. **如果函數是 trivial field accessor，hook 它的呼叫者**——真正的邏輯在上一層
5. **記錄函數的真實行為**——不要只依賴名稱，用 Frida hook 確認實際行為

#### Goal / Action / Validation

- Goal: 正確識別 Dart AOT 函數的真實用途，避免被名稱誤導
- Action: 檢查函數大小、回傳值 class tag、交叉比對命名
- Validation or reference source: Frida capture 顯示函數只有 12 bytes，回傳值永遠是空的 container

#### Applies When

- 分析 `flutter_meta.json` 或 `blutter` 產生的 Dart AOT 函數列表
- 函數名稱暗示特定用途（如 `IV_*`、`get*`、`*Provider`）
- 函數大小 ≤12 bytes（ARM64）
- 回傳值看起來像空的 container object

#### Does Not Apply When

- 函數大小 > 12 bytes（有實際邏輯）
- 使用標準 Dart VM（非 AOT 編譯）
- 函數名稱來自人工命名（非自動推導）

#### Validation

- Frida capture 確認函數只有 12 bytes
- Frida capture 確認回傳值永遠是空的 container（class tag `0x010d011c`，4 null fields）
- 同一個 offset 在 `hook_skyshield_aes_key_expansion.js` 中被命名為 `keyExpansion2`

#### Promotion Target

- `intelligence/engineering/analytical-reasoning/heuristics/` — 新增 heuristic：「Dart AOT trivial field accessor detection」
- `workflow/apk-analysis/execution-flow.md` — 新增步驟：「分析 Dart AOT 函數前先檢查大小和回傳值」

#### Required Linked Updates

- 無需連動更新；這是新 lesson。
