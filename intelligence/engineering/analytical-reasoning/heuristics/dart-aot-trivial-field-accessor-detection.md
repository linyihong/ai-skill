# Dart AOT 極短欄位存取子辨識

`intelligence/engineering/analytical-reasoning/heuristics/dart-aot-trivial-field-accessor-detection.md`

## 問題

逆向 Dart AOT 時，`flutter_meta.json` 或 `blutter` 產生的函數名稱可能誤導分析方向。短函數（≤12 bytes ARM64）幾乎總是 trivial field accessor，不是預期的業務 getter。

## 啟發式

| 徵象 | 解讀 | 行動 |
|--------|---------------|--------|
| 函數大小 ≤ 12 bytes（ARM64） | Trivial field accessor（`LDR X0, [X0, #offset]` + `RET`） | 不要分析此函數，hook 它的呼叫者 |
| 回傳值 class tag = `0x010d011c` | 空的 Dart container object | 表示 container 已分配但未填入資料 |
| 回傳值包含 4 個 `81 80 00 00` | 4 個 null fields | 確認 container 未初始化 |
| 無論何時呼叫，回傳值都相同 | 函數沒有副作用，只是讀取欄位 | 這是純粹的 field accessor |
| 同一個 offset 在不同來源有不同名稱 | 名稱衝突，需要進一步調查 | 用 Frida hook 確認實際行為 |

## 決策表

| 條件 | 結論 |
|-----------|-----------|
| 函數 ≤ 12 bytes + 回傳 empty container | Trivial field accessor，跳過 |
| 函數 ≤ 12 bytes + 回傳非空值 | Field accessor，但欄位有值，可追蹤誰寫入該欄位 |
| 函數 > 12 bytes + 名稱暗示特定用途 | 可能是真正的業務函數，需要分析 |
| 名稱衝突（不同來源不同命名） | 需要 Frida hook 確認真實行為 |

## 來源

- `feedback/history/apk-analysis/flutter-dart-aot/2026-05-14_081500-iv-f5447c-trivial-field-accessor-not-iv-getter.md`
