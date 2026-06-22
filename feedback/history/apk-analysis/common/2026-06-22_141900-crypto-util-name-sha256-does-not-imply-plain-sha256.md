> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md)、[neutral-language](../../../../enforcement/neutral-language.md)、[goal-action-validation](../../../../enforcement/goal-action-validation.md)、[sanitization](../../../../enforcement/sanitization.md)、[reusable-guidance-boundary](../../../../enforcement/reusable-guidance-boundary.md) 與 [feedback-lessons](../../../feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-06-22 - Crypto util named sha256* does not imply plain SHA-256 of canonical

Status: candidate

#### One-line Summary

DEX 中 `sha256Encrypt(String, int)` 或 native `.so` 內含 `SHA256` 符號，**不能**推論 `sign = SHA256(canonical)`；必須用 Frida **入參 canonical + 出參 hex** 對照，Python `hashlib.sha256` 不匹配時改走 native RE 或 in-app RPC relay。

#### Human Explanation

保護庫常內嵌標準 SHA-256 實作（round constants `sha256_k` 等），但對外 API 可能在 hash 前加 **salt/key/mode 分支/二次編碼**。方法名 `sha256Encrypt` 與 64 hex 輸出會誤導分析者直接寫 Python signer。正確 gate：拿到一組 `(canonical, sign)` 後立刻驗證 plain hash；失敗則停止猜測，改 hook native 或 `Java.use` RPC 呼叫同一 util。

#### Trigger

- Util method 名含 `sha256`、`encrypt`、`sign`
- `nm -D` 或 strings 在 protection `.so` 見 `SHA256::` / `_Z6sha256`
- Frida 已捕獲 canonical 字串與 64 hex sign
- `hashlib.sha256(canonical.encode()).hexdigest()` ≠ captured sign

#### Evidence

- Tool: Frida hook util `sha256Encrypt` + Python one-liner verify
- Sanitized excerpt: mode 整數參數（如 `1`）；plain SHA256 mismatch；native library 有 SHA256 class 但輸出不等于 plain hash
- Evidence path: `<PROJECT_ROOT>/api/signing-re.md`

#### Generalized Lesson

```text
After canonical string is known:
  1. Verify: hashlib.sha256(canonical) == captured sign ?
  2. If NO: do NOT brute-force HMAC with guessed keys as first step
  3. Prefer: Frida RPC call same util(canonical, mode) OR native trace
  4. Protection .so may embed SHA256 for other paths (requestTime encrypt, etc.)
```

#### Agent Action

1. Project docs 標註「plain SHA256 ruled out」與 verify 方法。
2. Ai-skill 不寫 crack key 步驟；只寫 verification gate。
3. 交叉引用 `sha256-hash-verify-python-not-shell`（用 Python 驗證，非 shell）。

#### Goal / Action / Validation

- Goal: 避免在錯誤假設上實作離線 signer。
- Action: canonical builder 與 crypto 分離；crypto 層標 blocking until RE/relay。
- Validation: 至少一組 capture 證明 plain hash 不匹配；relay 可 reproduce sign。

#### Applies When

- Native crypto util with sha256-like naming
- Sign output is 64 hex chars

#### Does Not Apply When

- Plain SHA256(canonical) already verified against capture
- Sign is HMAC with documented key in DEX strings

#### Validation

- Documented mismatch + successful Frida RPC reproduce for same canonical

#### Promotion Target

- `workflow/apk-analysis/execution-flow.md` §sign RE verification gate

#### Required Linked Updates

- `feedback/history/apk-analysis/README.md` 索引追加
- 已依 sanitization / reusable-guidance-boundary 自查
