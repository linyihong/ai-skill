> 遵守 [共用規則索引](../../../../shared-rules/README.md)、[dependency-reading](../../../../shared-rules/dependency-reading.md)、[neutral-language](../../../../shared-rules/neutral-language.md)、[goal-action-validation](../../../../shared-rules/goal-action-validation.md) 與 [feedback-lessons](../../../../shared-rules/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-06 - Dart AOT Offset From ASM Address

Status: promoted

#### One-line Summary

用 Frida hook Dart AOT 時，offset 應以 asm 檔第一行函式地址為準，不要把輸出檔名 suffix 當成位址。

#### Human Explanation

部分 Dart AOT dump 檔名會帶一段短 hash / suffix，例如 `SomeFunction_4f914c.txt`。這段不一定是函式 entry offset。真正的 hook offset 應讀檔案第一行的函式地址，例如：

```text
0x00fd024c ... ; <SomeFunction_4f914c>
```

若誤用 `0x4f914c`，Frida 可能掛到 `r--` 資料區或其他非 code range，出現 `unable to intercept function`，且讀到的 bytes 與 asm prologue 不一致。

#### Trigger

媒體解密分析時，多個低位址函式（例如 decrypt/cache helper）使用檔名 suffix 作為 offset，Frida 全部拒絕 hook。經 offset/protection probe 後發現目標地址落在 `r--` range，且 bytes 不符合 asm。改用 asm 第一行完整地址後，函式可正常 hook 並命中。

#### Evidence

- Failure shape: `Interceptor.attach(base.add(0x4f914c))` -> `unable to intercept function`，range 為 `r--`，bytes 不符合函式 prologue。
- Correct shape: asm 第一行顯示 `0x00fd024c ... <FileDecryptLoader._loadAndDecryptImage...>`，改用 `base.add(0xfd024c)` 後可 hook 並命中。
- Evidence type: sanitized Frida offset probe log；不含 target host、token、URL 或 media bytes。

#### Generalized Lesson

在 blutter/unflutter/Dart AOT 產物中，檔名 suffix 只能當搜尋提示，不能直接當 hook offset。每次建立 Frida offset hook 前，先讀 asm 第一行地址，並用 runtime probe 檢查：

1. `base + offset` 所在 memory range 是否可執行（通常 `r-x`）。
2. runtime bytes 是否與 asm prologue 對得上。
3. 若 `Interceptor.attach` 失敗，先懷疑 offset / load bias / segment range，再懷疑 hook 內容。

#### Agent Action

下次做 Dart AOT offset hook：

1. 從 asm 第一行取完整十六進位地址。
2. 不使用檔名尾碼推導 offset，除非已與第一行地址一致。
3. 對失敗 hook 跑小型 probe：輸出 range protection、前 16-24 bytes、`+0/+4/+8...` attach 可行性。
4. 只有確認在 executable range 且 bytes 對齊後，才進入參數解析/物件 layout debugging。

#### Goal / Action / Validation

- Goal: 避免因錯誤 offset 把 Frida hook 掛到資料區，浪費時間 debug 假問題。
- Action: 使用 asm 第一行完整地址，並以 runtime range/prologue probe 驗證。
- Validation or reference source: `base + offset` 位於 executable range，runtime bytes 與 asm prologue 一致，hook 能命中。

#### Applies When

- 使用 unflutter/blutter/自製 dump 的 Dart AOT asm 檔建立 Frida native offset hook。
- 輸出檔名含短 hash / suffix，且與第一行地址不完全相同。
- `Interceptor.attach` 顯示 unable to intercept，但其他 offset hook 正常。

#### Does Not Apply When

- 使用符號表或 runtime API 直接解析到 function entry，且已驗證地址。
- hook 目標本來就是 callsite / inline stub，非函式 entry；這時仍需另行證明 callsite 可靠性。

#### Validation

- 至少一個更正後 offset 成功 hook 並命中。
- 文件或 hook 腳本中保留 offset 來源（asm 第一行地址或 probe 結果）。

#### Promotion Target

- `WORKFLOW.md` / Dart AOT offset hook checklist.

#### Required Linked Updates

- 若專案已有 private hook 腳本，修正為 asm 第一行完整 offset。
