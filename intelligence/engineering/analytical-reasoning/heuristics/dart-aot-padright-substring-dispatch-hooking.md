# Dart AOT 字串操作假設與分派 Hooking

`intelligence/engineering/analytical-reasoning/heuristics/dart-aot-padright-substring-dispatch-hooking.md`

## 問題

逆向 Dart AOT 的字串操作時，容易誤判三個行為：
1. `padRight` 可能使用 null bytes（`\0`）而非 spaces（`0x20`）
2. `substring(0, N)` 在短字串上可能不 crash（自訂實作優雅處理越界）
3. vtable dispatch 比直接 hook 目標函數更穩健

## 啟發式

### 1. `padRight` Padding Detection

| 徵象 | 解讀 |
|--------|---------------|
| Hex dump 顯示 `\0`（0x00）而非 spaces（0x20） | 自訂實作使用 null bytes padding |
| 字串長度 < 預期但沒有 crash | 自訂實作不拋出越界異常 |
| 使用 `encrypt` 套件的 `padRight` | 可能使用標準 space padding |

### 2. `substring` Behavior

| 徵象 | 解讀 |
|--------|---------------|
| `substring(0, 32)` 在 16-char 字串上不 crash | 自訂實作優雅處理越界 |
| 回傳值長度 < 參數 N | 字串實際長度 < N，substring 回傳原字串 |

### 3. Dispatch Hooking Strategy

| 徵象 | 行動 |
|--------|--------|
| 不知道 vtable dispatch 的目標函數 | Hook `UBFX` + `LDR` + `BLR` 序列捕獲所有目標 |
| 直接 hook 目標函數失敗 | 改用 dispatch hooking 確認哪些函數被實際呼叫 |
| 需要捕獲完整流程 | 先用 dispatch hook，再收斂到具體函數 |

## 決策表

| 條件 | 建議作法 |
|-----------|---------------------|
| 需要確認 `padRight` 的填充字元 | Hex dump 輸出，不要只看字串表示 |
| 需要確認 `substring` 是否越界 | 檢查輸入字串的實際長度 |
| 不知道 vtable dispatch 的目標 | Hook dispatch 指令（`UBFX` + `LDR` + `BLR`） |
| 需要捕獲完整流程 | 先用 dispatch hook，再收斂 |

## 來源

- `feedback/history/apk-analysis/flutter-dart-aot/2026-05-14_073700-dart-aot-padright-null-bytes-substring-noop.md`
