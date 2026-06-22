> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md)、[neutral-language](../../../../enforcement/neutral-language.md)、[goal-action-validation](../../../../enforcement/goal-action-validation.md)、[sanitization](../../../../enforcement/sanitization.md)、[reusable-guidance-boundary](../../../../enforcement/reusable-guidance-boundary.md) 與 [feedback-lessons](../../../feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-06-22 - Native encrypted-time getter may cache — first call rotates, consecutive calls return same blob

Status: candidate

#### One-line Summary

Encrypted-time 類 native getter 除 interceptor **雙呼叫**外，還可能有 **狀態快取**：連續兩次呼叫中 **第一次產生新 blob、第二次與之相同**；第三次若緊接第二次亦相同。RPC relay 仍須 mirror interceptor 的兩次呼叫，但不要用「每次 Java 呼叫都會變」假設做離線模型。

#### Human Explanation

`142100` 覆蓋「第二次輸出進 sign map」。實測還可能見到：`(call1) ≠ (call2) == (call3)`。表示 native 層在第一次呼叫時旋轉內部狀態/計數器，隨後短時間內回傳快取值。這與「每個 Java 呼叫都獨立隨機」不同，也與「兩次輸出完全相同」不同。離線 RE 須 hook 函數體內 mutex/全域狀態，而非僅記錄 Base64 樣本。**getter 名、wire 欄位名**留在專案。

#### Trigger

- Frida：手動連續呼叫同一 getter（相同 seed）三次，`a≠b` 且 `b==c`
- Interceptor 雙呼叫：第二次與第三次攔截間輸出可能相同
- 離線猜測 encrypted-time 完全隨機每次失敗

#### Evidence

- Tool: Frida Java 連續呼叫實驗 + spawn JNI log 成對相同 OUT
- Sanitized excerpt: triple-call `same=false` then `b==c=true`
- Evidence path: `<PROJECT_ROOT>/api/signing-re.md`（專案加密時間欄位 RE）

#### Generalized Lesson

```text
After confirming double-call (142100):
  1. Probe 3+ consecutive getter calls with same seed
  2. If a!=b and b==c: model = rotate-on-first + cache for burst
  3. RE target shifts to native global/mutex, not only blob algebra
  4. RPC relay unchanged: still call twice per signed request
```

#### Agent Action

1. 更新專案 signing-re 的 call-semantics 表。
2. 交叉引用 `142100`；Ai-skill 不寫 offset/key。

#### Goal / Action / Validation

- Goal: 避免錯誤隨機模型拖延 RE。
- Validation: 記錄 triple-call 行為 + interceptor 第二次仍可用於 sign。

#### Applies When

- Encrypted-time header via native getter
- Double-call pattern already confirmed

#### Does Not Apply When

- Every call returns unique value (true per-call RNG)
- Plain epoch header only

#### Validation

- Documented triple-call matrix in project evidence

#### Promotion Target

- `workflow/apk-analysis/execution-flow.md` §encrypted-time semantics

#### Required Linked Updates

- `feedback/history/apk-analysis/README.md` 索引追加
- `142100` Agent Action 交叉引用本條
- 已依 sanitization / reusable-guidance-boundary 自查
