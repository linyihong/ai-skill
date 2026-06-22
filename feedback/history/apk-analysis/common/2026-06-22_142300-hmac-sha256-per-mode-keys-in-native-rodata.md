> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md)、[neutral-language](../../../../enforcement/neutral-language.md)、[goal-action-validation](../../../../enforcement/goal-action-validation.md)、[sanitization](../../../../enforcement/sanitization.md)、[reusable-guidance-boundary](../../../../enforcement/reusable-guidance-boundary.md) 與 [feedback-lessons](../../../feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-06-22 - After plain SHA-256 fails, verify HMAC-SHA256 with per-mode keys from native rodata

Status: candidate

#### One-line Summary

`sign` 為 64 hex 且 plain `SHA256(canonical)` 不匹配時，下一個驗證分支是 **HMAC-SHA256(canonical, key_mode)**；`mode` 常為 Java 整數常數（如 normal/report/feedback）；key 常嵌在 native `.so` **`.rodata`**，可用 `strings` / 靜態常數名鄰近區段定位，**不得**把實際 key 寫入 Ai-skill，只寫提取與驗證流程。

#### Human Explanation

`sha256Encrypt(String, int)` 這類 API 在 plain hash 失敗後，高機率是 **HMAC-SHA256 + mode 分支 key**。DEX `<clinit>` 或 static field 常暴露 mode 整數與語意名（`*SignKey`）。native 庫 `.rodata` 常並排存放 32-byte ASCII key 與錯誤訊息／常數名。正確路徑：(1) Frida 捕獲 `(canonical, mode, sign)`；(2) 對每個 mode 用 Python `hmac.new(key, canonical, sha256).hexdigest()` 驗證；(3) key 來源記在 **project** `signing-re.md`，Ai-skill 只記「rodata + mode 表 + verify gate」。與 `141900` 銜接：141900 禁止在無證據時亂猜 key；本條是在 **已有 capture** 後的結構化下一步。

#### Trigger

- `hashlib.sha256(canonical)` ≠ captured 64-hex sign
- Native util 第二參數為 `int mode`；DEX 有 `*SignKey` 類 static int
- `strings lib*.so` 見多組 32 字元可列印 key 候選，鄰近 `SignKey` / `normal` / `report` 字串
- Frida JNI 已能列出 `(canonical, mode)` → `sign`

#### Evidence

- Tool: Frida JNI in/out + Python `hmac` one-liner per mode + `strings`/rodata xref on loaded `.so`
- Sanitized excerpt: mode `1|2|3` 各有一組 key；HMAC verify 與 JNI output 一致；method 名仍含 `sha256`
- Evidence path: `<PROJECT_ROOT>/api/signing-re.md`（key 與版本綁定，不進 Ai-skill）

#### Generalized Lesson

```text
After canonical + captured sign known, plain SHA256 ruled out (141900):
  1. For each observed mode int M: try HMAC-SHA256(key_M, canonical)
  2. key candidates: strings/rodata near *SignKey* labels in the JNI-loaded .so
  3. First mode where hmac == captured sign → algorithm confirmed for that path
  4. Document keys in PROJECT signing doc only; Ai-skill documents branch + verify
  5. Other natives (requestTime, DB pwd) may use different keys — do not reuse sign key
```

#### Agent Action

1. **每輪 sign RE 有進展時**：自問是否已寫此分支；若 HMAC 已驗證，更新 project `signing-re.md` + 離線 signer，並寫/更新本條或交叉引用。
2. Ai-skill 正文不含 host、package、offset、實際 key 字串。
3. 交叉引用 `141900`、`141800`、`142000`。

#### Goal / Action / Validation

- Goal: plain-hash 死胡同後有確定性的 HMAC + rodata 路徑，避免無限猜測。
- Action: mode 表來自 DEX static；key 來自 rodata；verify 用 capture 對照。
- Validation: ≥1 組 `(canonical, mode, sign)` Python HMAC 與 Frida JNI 一致。

#### Applies When

- Sign output 64 hex; plain SHA256 mismatch
- Native `sha256*` util with int mode parameter
- Protection/obfuscation library holds crypto

#### Does Not Apply When

- Plain SHA256 already matches capture
- Sign is RSA/ECDSA or non-hex encoding
- Keys only derivable server-side (no stable offline key in binary)

#### Validation

- Documented HMAC verify success per mode + project key location (not in Ai-skill)

#### Promotion Target

- `workflow/apk-analysis/execution-flow.md` §sign RE after plain-hash gate

#### Required Linked Updates

- `feedback/history/apk-analysis/README.md` 索引追加
- 交叉更新 `141900` Agent Action 指向本條為下一分支
- 已依 sanitization / reusable-guidance-boundary 自查
