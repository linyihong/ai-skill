> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md)、[neutral-language](../../../../enforcement/neutral-language.md)、[goal-action-validation](../../../../enforcement/goal-action-validation.md)、[sanitization](../../../../enforcement/sanitization.md)、[reusable-guidance-boundary](../../../../enforcement/reusable-guidance-boundary.md) 與 [feedback-lessons](../../../feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-06-22 - Custom native F3AES may still be standard AES-128-CBC with rodata key/IV

Status: candidate

#### One-line Summary

Native 函數名/字串（如 `F3AES`）暗示自訂實作，但 **hook key-setup 讀入的 16B key/IV 常可直接用標準 AES-128-CBC + PKCS7** 還原 wire blob。不要因 custom label 跳過標準 cipher 驗證。

#### Human Explanation

Native encrypt core 可能標記自訂 label（如 `F3AES`），初試 rodata 字串當 key 用標準庫失敗（key/IV 順序或字串解碼錯誤）。Frida hook key-setup 得到實際 key/IV 後，標準 AES-128-CBC 與固定長度 wire blob **完全一致**。教訓：custom crypto label ≠ non-standard algorithm；先抓 runtime key material 再驗證標準 primitive。**offset、key 值、欄位名**留在專案。

#### Trigger

- rodata 可見 `AES`/`F3AES` 但 offline standard AES 初試失敗
- Wire output 固定 80 B（Base64 ~108 chars）after PKCS7 to 16-byte boundary
- Native mode integer selects key pair from rodata

#### Evidence

- Tool: Frida hook key-setup + PyCryptodome encrypt/decrypt round-trip
- Sanitized excerpt: `match True` on captured plaintext/ciphertext pair
- Evidence path: `<PROJECT_ROOT>/api/signing-re.md`（專案 encrypt RE）

#### Generalized Lesson

```text
After plaintext shape known (143000/143100):
  1. Hook key/IV setup (not only encrypt exit)
  2. Log exact 16-byte key and IV strings passed to context init
  3. Try AES-128-CBC PKCS7 before RE custom rounds
  4. Store keys in project evidence; Ai-skill omits key bytes in slug body if policy requires
```

#### Agent Action

1. 更新 project signing-re AES 表與 offline helper。
2. 交叉引用 `142300`（mode keys）、`143000`、`143100`。
3. Ai-skill 範式不寫具體 key 值（project doc 可寫）。

#### Goal / Action / Validation

- Goal: 避免在 custom AES RE 上過度投入。
- Validation: offline encrypt == Frida wire for same plaintext millis.

#### Applies When

- Opaque native encrypt with mode index + rodata strings
- Output length multiple of 16 after decode

#### Does Not Apply When

- Hook shows non-16 key length or stream cipher state
- Standard AES verify fails after correct runtime key capture

#### Validation

- Round-trip documented in project evidence

#### Promotion Target

- `workflow/apk-analysis/execution-flow.md` §verify standard cipher after key hook

#### Required Linked Updates

- `feedback/history/apk-analysis/README.md` 索引追加
- 已依 sanitization / reusable-guidance-boundary 自查
