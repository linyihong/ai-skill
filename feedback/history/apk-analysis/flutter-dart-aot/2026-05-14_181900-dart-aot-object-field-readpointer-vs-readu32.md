### 2026-05-14 - Dart AOT 物件欄位讀取：`readPointer()` vs `readU32()` — 64-bit 指標不可用 32-bit 讀取

#### One-line Summary
Dart AOT 物件的某些欄位包含**完整的 64-bit 指標**（8 bytes），必須用 Frida 的 `NativePointer.readPointer()` 讀取；使用 `readU32()` 只會讀到低 32 bits，造成位址截斷和 access violation。

#### Human Explanation
在 64-bit Android 上，Dart AOT 物件的欄位可能包含指向其他 heap 物件的指標。這些指標是**完整的 64-bit 位址**（8 bytes），不是 compressed 32-bit 指標。如果使用 `readU32()`（只讀 4 bytes）來讀取這樣的欄位，只會得到位址的低 32 bits，高 32 bits 被靜默丟棄。

例如，一個 Dart 物件的欄位在 byte offset 0x38 處包含指標 `0x0000007a6b027f4b`：
- `readU32()` 回傳 `0x6b027f4b`（number，缺少 `0x0000007a` 前綴）
- `readPointer()` 回傳 `0x7a6b027f4b`（NativePointer，完整位址）

使用截斷的位址進行記憶體讀取會觸發 access violation，因為 `0x6b027f4b` 不是有效的對映位址。

這與 JavaScript 位元運算子截斷（`&`, `|`, `~` 造成 32-bit 截斷）是**不同的問題**——這裡是 Frida 的 `NativePointer` API 選擇問題，不是 JS 語言特性問題。

#### Trigger
- 寫 Frida hook 讀取 Dart AOT heap 物件的欄位
- 使用 `readU32()` 讀取欄位值，然後嘗試用該值作為指標進行 `hexdump` 或 `readByteArray`
- `hexdump` 回報 "access violation"，且位址看起來缺少高 bits
- 或者 `isLikelyPtr()` 對欄位值回傳 `false`，因為截斷後的位址落在使用者空間的低位區域

#### Evidence
- `hook_final_assembly.js` v2 中，`x3` 物件的欄位在 offset 0x38 處包含指標 `0x7a6b027f4b`
- `readU32()` 回傳 `0x6b027f4b`（截斷，缺少 `0x7a` 前綴）
- 使用 `0x6b027f4b` 進行 `hexdump` 觸發 "access violation"
- 修正為 `readPointer()` 後，正確讀取到完整位址 `0x7a6b027f4b`，後續資料讀取成功
- 同樣問題也發生在 return value 的欄位讀取（offsets 0x08, 0x0c, 0x10, 0x0f, 0x14, 0x18）

#### Generalized Lesson
1. **Dart AOT 物件的指標欄位永遠用 `readPointer()`** — 它回傳 `NativePointer`（完整 64-bit），不是 number
2. **`readU32()` 只適用於 Smi（small integer）欄位**，不適用於指向其他物件的指標欄位
3. 當 `isLikelyPtr()` 對預期是指標的欄位值回傳 `false` 時，先檢查是否用了 `readU32()` 而不是 `readPointer()`
4. 在 Dart AOT 中，物件欄位的類型無法從靜態分析得知——必須透過實驗判斷該欄位是 Smi 還是指標
5. 一個安全的啟發式方法：如果欄位值在 `0x00000000`–`0x0000ffff` 範圍內，可能是 Smi；如果值很大（如 `0x7a6b027f4b`），則是指標

#### Agent Action
- 發現 `x3` 物件的 offset 0x38 處的資料讀取失敗（access violation）
- 檢查 `readU32()` 的回傳值，發現它只回傳了 4 bytes（低 32 bits）
- 改用 `readPointer()` 後，正確讀取到完整的 64-bit 位址
- 修正所有 Dart 物件欄位讀取：將 `readU32()` 替換為 `readPointer()`，特別是在 offsets 0x08, 0x0c, 0x10, 0x0f, 0x14, 0x18, 0x38, 0x48, 0x58 等位置
- 修正後所有資料讀取都正確

#### Goal / Action / Validation
- **目標**：在 Frida 中正確讀取 Dart AOT 物件的指標欄位
- **行動**：將 `readU32()` 替換為 `readPointer()` 用於所有指標欄位讀取
- **驗證**：修正後 `hexdump` 成功顯示資料，不再有 access violation

#### Applies When
- 寫 64-bit Dart AOT 應用的 Frida hook
- 需要讀取 Dart 物件的欄位，且不確定該欄位是 Smi 還是指標
- `readU32()` 回傳的值看起來像截斷的位址（缺少高 bits）

#### Does Not Apply When
- 讀取 Smi（small integer）欄位（`readU32()` 是正確的）
- 使用 Dart VM 的 service protocol（不是直接記憶體讀取）
- 處理 32-bit Android 應用

#### Validation
- 修正前：`ptr(objAddr).add(0x38).readU32()` → `0x6b027f4b`（截斷，錯誤）
- 修正後：`ptr(objAddr).add(0x38).readPointer()` → `0x7a6b027f4b`（完整，正確）

#### Promotion Target
- `intelligence/engineering/analytical-reasoning/heuristics/dart-aot-field-read-pointer-vs-smi.md`（建議新增）

#### Required Linked Updates
- `feedback/history/apk-analysis/README.md`：`flutter-dart-aot/` 計數從 19 → 20
- `intelligence/engineering/analytical-reasoning/heuristics/README.md`：新增 atom 列（若 promote）
- 所有現有 hook script 中的 Dart 物件欄位讀取都應檢查是否正確使用 `readPointer()` vs `readU32()`
