> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md)、[neutral-language](../../../../enforcement/neutral-language.md)、[goal-action-validation](../../../../enforcement/goal-action-validation.md)、[sanitization](../../../../enforcement/sanitization.md)、[reusable-guidance-boundary](../../../../enforcement/reusable-guidance-boundary.md) 與 [feedback-lessons](../../../feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-06-22 - Native decrypt mode N may reuse encrypt mode N key material

Status: candidate

#### One-line Summary

同一 native crypto 模組裡，**decrypt 的 mode index 常與 encrypt 共用 key/IV 表**。解出 encrypt 後應 hook key-setup 驗證 decrypt，而非假設 API blob 用另一套 key。

#### Human Explanation

`142700` 已記錄「opaque blob → decrypt → secondary compression」。若 encrypt path 的 mode `N` key 已確認，對 API 回傳的同名 `decryptStr(..., N)` **先 hook key material**；若相同，離線鏈多為 `base64(wire) → AES decrypt → inner base64 → zlib/json`。具體欄位名、腳本名、key 值留在專案。

#### Trigger

- Encrypt mode N solved; response blob still needs Frida decrypt
- Hook shows same 16B key/IV at crypto context init
- Wire blob length is multiple of 16 after base64 decode

#### Evidence

- Tool: key-setup hook + offline round-trip
- Sanitized excerpt: decrypt with known key == native output
- Evidence path: `<PROJECT_ROOT>/api/signing-re.md` §response blob decode

#### Generalized Lesson

```text
After encrypt RE for mode N:
  1. Hook decrypt entry with same mode N before full RE
  2. Compare key-setup inputs to encrypt path
  3. If match: reuse keys in offline decoder module
  4. Chain with secondary transform from 142700
```

#### Agent Action

1. 專案 evidence 記錄 wire format 與 decoder。
2. 交叉引用 `143200`（custom label ≠ non-standard cipher）。
3. Ai-skill 不寫 key bytes。

#### Goal / Action / Validation

- Goal: 避免對已知 mode 重複做 key RE。
- Validation: offline decode == native on same blob (project doc).

#### Applies When

- Paired encrypt/decrypt natives with mode integer
- Same `.so` / crypto context init function

#### Does Not Apply When

- Key hook shows different material per direction
- Blob not block-aligned after decode

#### Validation

- Project round-trip only; lesson stays name-free

#### Promotion Target

- `workflow/apk-analysis/execution-flow.md` §shared mode key table

#### Required Linked Updates

- `feedback/history/apk-analysis/README.md` 索引
- 已依 sanitization / reusable-guidance-boundary 自查
