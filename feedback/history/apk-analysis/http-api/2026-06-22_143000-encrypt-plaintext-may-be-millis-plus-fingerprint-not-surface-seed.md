> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md)、[neutral-language](../../../../enforcement/neutral-language.md)、[goal-action-validation](../../../../enforcement/goal-action-validation.md)、[sanitization](../../../../enforcement/sanitization.md)、[reusable-guidance-boundary](../../../../enforcement/reusable-guidance-boundary.md) 與 [feedback-lessons](../../../feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-06-22 - Encrypt plaintext may be millis + fingerprint, not the surface seed argument

Status: candidate

#### One-line Summary

Native encrypt core 的 **plaintext 往往不是 Java 傳入的表面 seed**（如 install path），而是 **`{millis}_{stable_device_fingerprint}`**（fingerprint 常為 colon 分隔 20-byte hex）。離線 RE 應 hook `strcpy`/encrypt 入口抓 plaintext，而非假設 seed 字串直接進 cipher。

#### Human Explanation

Java 層 getter 可能接收 path 類參數，但 native slow-path 會先從 install artifact 推導穩定指紋並快取；每次生成再 prepend 當前毫秒時間戳後送入 mode-N encrypt。指紋 suffix **同 install 穩定**；prefix millis **隨生成變化**。這解釋了為何 seed 參數會影響輸出，但 Frida hook encrypt 時看到的是 timestamp 字串而非原始 path。**欄位名、getter 名、offset、key** 留在專案 evidence。

#### Trigger

- Hook native encrypt core 卻讀不到表面 seed 明文
- `strcpy` 在 encrypt 函數內被呼叫且輸入像 `1782…_2F:D8:…`
- 離線用 path 直接當 AES plaintext 驗證失敗

#### Evidence

- Tool: Frida `strcpy` filter on encrypt core return address range + JNI seed log
- Sanitized excerpt: `PLAINTEXT 1782110705178_2F:D8:33:…:CE`（suffix stable across calls）
- Evidence path: `<PROJECT_ROOT>/api/signing-re.md`（專案加密時間欄位 RE）

#### Generalized Lesson

```text
When surface seed is a path but encrypt input is opaque:
  1. Hook strcpy/memcpy inside native encrypt (not only JNI jstring)
  2. Parse plaintext shape: millis + delimiter + fingerprint
  3. RE fingerprint derivation separately (slow-path file read / global cache)
  4. Offline model: encrypt(f"{ms}_{fingerprint}") -> base64, not encrypt(path)
```

#### Agent Action

1. 更新專案 signing-re 的 plaintext 表。
2. 交叉引用 `142100`（雙呼叫）、`142800`（快取語意）。
3. Ai-skill 不寫 SO offset / AES key / wire 欄位名。

#### Goal / Action / Validation

- Goal: 把離線 RE 從「path 雜湊」轉向「fingerprint + timestamp + mode-N encrypt」。
- Validation: 同一 session 多次 hook 見相同 fingerprint suffix、不同 millis prefix。

#### Applies When

- Native encrypt wrapper with path-like seed but opaque ciphertext
- Custom encrypted-time blob ~80 bytes after base64 decode

#### Does Not Apply When

- Plain epoch seconds in header
- Path string hashed directly without timestamp component

#### Validation

- Documented plaintext format in project evidence + Frida capture script

#### Promotion Target

- `workflow/apk-analysis/execution-flow.md` §encrypted-time plaintext discovery

#### Required Linked Updates

- `feedback/history/apk-analysis/README.md` 索引追加
- `142800` / `142100` cross-ref
- 已依 sanitization / reusable-guidance-boundary 自查
