# Dart AOT Pool Constant (PP_peep) 啟發式

`string_refs.jsonl` 中 `kind:"PP_peep"` 的條目可直接揭露硬編碼字串常數（如 AES key、secret），無需反組譯。

## 問題

在 Dart AOT 逆向分析中，尋找硬編碼的密碼學金鑰通常需要：
1. 反組譯整個函式
2. 追蹤 `LDR` 指令的資料流
3. 動態 hook 確認執行期值

這在分析大量函式時非常耗時。

## 根因

Dart AOT 編譯器將函式內使用的字串常數存放在該函式的 **pool area**（PC 相對定址區）。當函式需要某個字串時，會透過 `LDR X16, [X16, #offset]` 從 pool 載入。`unflutter` 的 `string_refs.jsonl` 會將這些 pool 載入記錄為 `kind:"PP_peep"`，並直接輸出字串的實際值。

這意味著：**不需要反組譯整個函式**，只需要在 `string_refs.jsonl` 中搜尋可疑的關鍵字（如 `encrypt`、`decrypt`、`key`、`secret`、header 名稱、base64 字串），就能直接找到硬編碼的密碼學金鑰或常數。

## 範例

在 `RequestInterceptor.onResponse` 的 `string_refs.jsonl` 條目中，`kind:"PP_peep"` 直接揭露了 AES-128 金鑰字串（16 字元 ASCII），該金鑰同時用於請求加密 header 與回應解密。

對應的反組譯確認：函式 disassembly 中確實有 `LDR X16, [X16, #offset]` 從 pool 載入該字串，然後傳遞給 `AES(AES._128048(...))` 建構子。

## Agent 操作步驟

1. 執行 `unflutter`（或相容的 AOT parser）後，找到 `string_refs.jsonl`。
2. 搜尋 `PP_peep` 條目中包含 `encrypt`、`decrypt`、`key`、`secret`、`aes`、`hmac`、`sign` 或已知 header 名稱的關鍵字。
3. 對每個候選條目，記錄 `func`（函式名稱）、`pc`（program counter）、`value`（實際字串值）。
4. 讀取候選 PC 附近的函式反組譯，確認 load instruction 並追蹤字串的使用路徑。
5. 如果字串是密碼學金鑰，透過以下方式驗證：
   - 檢查相同金鑰是否出現在其他函式中（跨函式重用）
   - 檢查金鑰長度是否符合預期演算法（16 bytes = AES-128, 32 bytes = AES-256 等）
   - 透過動態 Frida hook 確認金鑰參數是否匹配
6. 將金鑰記錄在專案私有的分析檔案中；**不要**在可重複使用的 lesson 中包含原始金鑰（依去敏規則）。

## 目標 / 行動 / 驗證

- **目標**：在不完整反組譯的情況下，找到 Dart AOT 函式中的硬編碼密碼學金鑰。
- **行動**：`grep '"PP_peep"' string_refs.jsonl | grep -i 'encrypt\|key\|secret\|aes'` 或使用 `jq` 過濾。
- **驗證**：交叉比對 PC 在函式反組譯中的位置，確認 `LDR` instruction 並追蹤金鑰的使用路徑。

## 適用時機

- 已有 `unflutter` 或相容 AOT parser 產生的 `string_refs.jsonl`。
- 目標函式使用硬編碼字串常數（而非執行期動態產生的金鑰）。
- 需要快速識別大量函式中的密碼學材料。

## 不適用時機

- 金鑰在執行期動態產生（例如來自金鑰衍生函式、DH 交換或外部儲存）。
- AOT parser 未產生 `string_refs.jsonl` 或 pool peep 分析不完整。
- 需要了解完整的加密流程，而不只是找到金鑰。

## 驗證

- `PP_peep` 的字串值與實際執行期值一致（透過 Frida `onEnter` 參數捕獲確認）。
- `string_refs.jsonl` 中的 PC 對應到反組譯中的 `LDR X16, [X16, #offset]` instruction。
- 相同金鑰出現在預期的演算法上下文中（例如傳遞給 `AES._128048` 建構子）。

## 來源

- `feedback/history/apk-analysis/flutter-dart-aot/2026-05-15_094700-dart-aot-pool-constant-pp-peep-hardcoded-key.md`
