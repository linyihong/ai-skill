> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md)、[neutral-language](../../../../enforcement/neutral-language.md)、[goal-action-validation](../../../../enforcement/goal-action-validation.md)、[sanitization](../../../../enforcement/sanitization.md)、[reusable-guidance-boundary](../../../../enforcement/reusable-guidance-boundary.md) 與 [feedback-lessons](../../../feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-06-22 - R8-obfuscated OkHttp Response may block peekBody; hook converter layer

Status: candidate

#### One-line Summary

R8 混淆後 `Response` 可能只剩極少 public 方法（如 `okhttp3.u0`），在 `RealInterceptorChain.proceed` 後呼叫靜態 `u0.h(response, limit)` 或 `peekBody` 常 **0 條 response 日誌**；request URL 已命中時，response body 應改 hook **Gson/Moshi converter** 或 **OkHttp Interceptor** 解包後字串，而非只依賴 Response API。

#### Human Explanation

Request 側混淆 hook（`proceed(o0)`）成功後，agent 常在同一 hook 內讀 response body。混淆 Response 類可能移除或 rename `peekBody`/`body()`，靜態 helper 亦可能非 debug 用途。結果：REQ 豐富、RSP 為 0，易誤判「response 加密」而非 hook 點錯層。

#### Trigger

- Frida chain hook 有大量 `[REQ]` 業務 URL
- 同腳本 `[RSP]` / body preview **0 行**
- `getDeclaredMethods()` on obfuscated Response 僅 handful 方法，無 `peekBody`

#### Evidence

- Tool: Frida probe `okhttp3.u0` methods + chain hook with post-proceed body helper
- Sanitized excerpt: REQ > 100；RSP = 0；Response 公開方法 < 10
- Evidence path: `<PROJECT_ROOT>/scripts/frida/hook_capture_full.js`、`<PROJECT_ROOT>/capture/frida_w3_*.log`

#### Generalized Lesson

**Response capture 分層（request 已命中後）：**

```text
1. proceed 後 Response.toString() / 混淆 static helper — 快速探測
2. 若 0 輸出 → jadx 找 GsonConverterFactory / 自訂 Interceptor
3. Frida hook converter.fromJson 或 interceptor 返回 plaintext String
4. 備選：MITM on in-session CDN API host + trust bypass
```

Request chain hook 與 response converter hook **解耦**，各驗各的。

#### Agent Action

1. 報告分開寫「URL/schema 已確認」vs「body shape 待 converter hook」。
2. 勿因 RSP=0 就停止；下一版腳本只加 converter 一層（增量策略）。
3. 交叉引用 `common/2026-06-22_120000-okhttp-r8-obfuscated-request-realinterceptorchain-proceed.md`。

#### Goal / Action / Validation

- Goal: 縮短 REQ/RSP 不對稱時的误判時間。
- Action: `analysis/apk/traffic-triage.md` Java hook 小節補 response 層。
- Validation: converter hook 後至少 1 條 redacted JSON 片段（project capture only）。

#### Applies When

- Release APK + OkHttp3 + R8
- Request chain hook 已 stable

#### Does Not Apply When

- Debug build Response 仍有 `peekBody`
- 只需 URL/method/header（無 body 需求）

#### Validation

- 同一 capture：REQ 計數 > 0 且 RSP 計數 = 0 → 觸發 converter 路線
- converter hook 後 RSP > 0 或 MITM 明文 JSON

#### Promotion Target

- `workflow/apk-analysis/execution-flow.md` §response body
- 交叉引用 incremental frida hook versioning lesson

#### Required Linked Updates

- `feedback/history/apk-analysis/README.md` common 計數 +1
