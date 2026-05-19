> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md)、[neutral-language](../../../../enforcement/neutral-language.md)、[goal-action-validation](../../../../enforcement/goal-action-validation.md) 與 [feedback-lessons](../../../../feedback/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-19 - H5 Bundle Parameter Codec And Runtime Readiness Gate

Status: candidate

#### One-line Summary

當 WebView/H5 API 只暴露加密的 `data` 參數時，先從前端 bundle 還原參數 codec 與 decoded shape；但「可產生同形 request」不等於「可在 SDK 外部成功重放」。

#### Human Explanation

嵌入式 H5 常把真實 query/body 包進單一加密欄位。只看動態 request metadata 會把它誤判為 opaque，進而阻塞後續分析；但反過來，只要在 bundle 中找到 encoder、key/iv 來源或 request interceptor，也不能立刻宣稱 SDK-ready。成功 schema 仍取決於 response wrapper、runtime token/session、gateway 條件、CORS/preflight 與錯誤語意。

比較穩定的分層是：先把 path-only evidence 升級成 request-params-decoded，再用 MITM、DevTools/CDP 或安全 response hook 取得 runtime 內的成功 response。外部 replay 若只得到業務錯誤碼，應記為 runtime-bound failure，而不是當成 API schema。

#### Trigger

- WebView / H5 request path 已被 `shouldInterceptRequest` 或代理看到，但 query/body 只有 `data`、`payload`、`p`、`params` 等單一加密欄位。
- H5 bundle 內可搜尋到 API path、HTTP client interceptor、CryptoJS/WebCrypto/AES helper、Base64 helper、request parameter rewrite logic。
- 外部 replay 以同形 request 打到 HTTP 200，但只返回 gateway/session/business error code。

#### Evidence

- Tool: H5 JavaScript bundle inspection, request interceptor reading, sanitized dynamic WebView metadata, local decode/encode smoke.
- Sanitized excerpt: A frontend bundle can show `JSON.stringify(params)` followed by symmetric encryption and Base64 wrapping; decoded shapes can reveal pagination, category, user/session-derived, or route-specific fields without storing raw values.
- Evidence path: Keep concrete hosts, endpoint paths, keys, tokens, user identifiers, ciphertext samples, and raw logs in project-private evidence or project docs; this lesson only records the generalized workflow and gates.

#### Generalized Lesson

1. **Search the bundle before calling `data` opaque**: look for API path strings, HTTP client interceptors, parameter rewrite functions, crypto helper names, and static key/iv assignments.
2. **Decode shape, not secrets**: document field names, value classes, enum-like values, pagination semantics, and whether fields are app/session/user-derived; do not persist raw ciphertext, tokens, ids, or bundle constants.
3. **Separate encoder readiness from API readiness**: being able to produce the encrypted parameter proves only the request codec. SDK readiness also needs a successful response schema and runtime authorization boundary.
4. **Treat HTTP 200 business errors carefully**: a replay that returns a business/session error proves transport reachability and wrapper shape, not success payload shape.
5. **Prefer runtime response capture for the next gate**: use MITM, WebView DevTools/CDP, or a safe response hook to capture the response inside the authorized WebView runtime before designing a client facade.

#### Agent Action

1. After path capture, fetch or inspect the exact H5 bundles loaded by the WebView.
2. Search for path strings and request interceptors; identify where params become the encrypted field.
3. Build a local decoder/encoder only in ephemeral analysis, then record sanitized decoded shape in project docs.
4. Mark status as `request-params-decoded`, not `schema-validated`, until runtime success responses are observed.
5. If external replay fails with a business/session error, document the error class and move the next action to runtime response capture rather than guessing SDK behavior.

#### Goal / Action / Validation

- Goal: Prevent premature SDK client implementation from request codec evidence while still extracting useful H5 API structure.
- Action: Promote path-only evidence to decoded request-shape evidence by reading the H5 bundle and validating decode/encode locally.
- Validation or reference source: Decoded request fields match multiple captured ciphertext samples or regenerated requests; success schema remains blocked until a runtime response capture provides wrapper keys and payload shape.

#### Applies When

- Embedded H5/WebView feature traffic uses encrypted single-field request parameters.
- The frontend bundle is available from WebView resource requests, static assets, APK assets, or cache.
- The work goal is API discovery, SDK readiness assessment, or capture planning.

#### Does Not Apply When

- The encrypted parameter is produced by native code, remote WebAssembly, or an unavailable runtime provider.
- The bundle constants or crypto material cannot be lawfully inspected or stored.
- The project only needs an H5 entry URL and will not implement or replay internal H5 APIs.

#### Validation

- Tracked docs contain no raw token, cookie, user id, ciphertext sample, full replay URL, private host, or static secret value.
- The decoded shape is derived from at least one dynamic request and one bundle-level encoder path.
- Any external replay result is labeled as success, transport-only, or runtime-bound failure; SDK readiness is not upgraded from a business-error response.

#### Promotion Target

- `analysis/apk/workflows/frida-hook-flow.md` for H5 bundle-to-request-shape workflow.
- `workflow/apk-analysis/execution-flow.md` for adding a `request-params-decoded` gate between `request-observed` and `schema-validated`.
- `workflow/software-delivery/execution-flow.md` only if this becomes a general SDK-readiness gate across multiple H5 integrations.

#### Required Linked Updates

- Project evidence should update the owning H5 API docs with sanitized decoded shapes and runtime-bound replay status.
- Existing short-window WebView capture lesson remains valid; this lesson refines the next gate after request path observation.
- `feedback/history/apk-analysis/README.md` category count should be updated with this lesson.
