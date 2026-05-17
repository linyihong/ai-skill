> 遵守 [共用規則索引](../../../enforcement/README.md)、[dependency-reading](../../../enforcement/dependency-reading.md)、[neutral-language](../../../enforcement/neutral-language.md)、[goal-action-validation](../../../enforcement/goal-action-validation.md) 與 [feedback-lessons](../../../enforcement/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-15 - write_to_file 大檔案截斷風險與防護

Status: candidate

#### One-line Summary

當使用 `write_to_file` 寫入超過 ~700 行的 Java 原始檔時，內容可能被截斷（truncated），導致編譯錯誤。應改用 `apply_diff` 分段寫入，或寫完後立即編譯驗證完整性。

#### Human Explanation

在修改大型原始檔（如 Java live test，>700 行）時，`write_to_file` 工具可能因為內容過長而被截斷，只寫入檔案的前半部分。這會導致：
- 檔案不完整（缺少後半部的 method 實作）
- 編譯錯誤（reference 到未定義的方法）
- 難以察覺（工具回報 success，但檔案實際上不完整）

這個問題在以下情況特別容易發生：
- 檔案行數 > 700
- 同時修改多個 section（main method + 多個 helper methods）
- 使用 `write_to_file` 做完整重寫而非 `apply_diff` 做局部修改

#### Trigger

`write_to_file` 回報 success，但編譯時出現大量「cannot find symbol」錯誤，或檔案在 IDE 中顯示不完整（結尾在中間某個 statement）。

#### Evidence

- Tool: `write_to_file` with ~1000-line Java file
- Sanitized excerpt: 檔案在 `if (ehValue !=` 處截斷（line 742 of 1000+），後續約 260 行遺失
- Evidence path: `<PROJECT_ROOT>/apk-analysis-sdk/tata-sdk/src/test/java/com/tata/sdk/shortdrama/ShortDramaLiveTest.java`

#### Generalized Lesson

當需要修改大型檔案時：

1. **優先使用 `apply_diff`** 做局部修改，而非 `write_to_file` 完整重寫。
2. **如果必須完整重寫**，將檔案分成多個較小的區塊，用多次 `apply_diff` 或分段 `write_to_file` 完成。
3. **寫入後立即編譯驗證** — 不要假設寫入成功就代表內容完整。
4. **監控檔案大小** — 如果檔案超過 500 行，特別注意截斷風險。
5. **考慮模組化** — 如果檔案持續增長，考慮拆分成多個類別。

#### Agent Action

1. 對於大型檔案（>500 行）的修改，優先使用 `apply_diff` 做精確的局部修改。
2. 如果必須完整重寫，先確認檔案行數，考慮分段寫入。
3. 寫入完成後，立即執行編譯命令驗證完整性。
4. 如果編譯失敗且錯誤指向「找不到符號」，先懷疑檔案被截斷，讀取檔案結尾確認。

#### Goal / Action / Validation

- Goal: 避免 `write_to_file` 截斷導致檔案不完整
- Action: 對大型檔案使用 `apply_diff` 分段修改，寫入後編譯驗證
- Validation or reference source: 編譯成功且所有 method 都存在

#### Applies When

- 修改的檔案行數 > 500
- 使用 `write_to_file` 完整重寫現有檔案
- 檔案包含多個 method 定義

#### Does Not Apply When

- 檔案很小（< 100 行）
- 使用 `apply_diff` 做局部修改
- 建立全新的小檔案

#### Validation

寫入後讀取檔案最後 10 行，確認內容完整；執行編譯確認無語法錯誤。

#### Promotion Target

- `enforcement/README.md`（若提升為全庫規則）

#### Required Linked Updates

- 無需連動更新。
