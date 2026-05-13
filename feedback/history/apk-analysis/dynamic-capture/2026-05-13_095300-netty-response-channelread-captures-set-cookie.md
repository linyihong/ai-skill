> 遵守 [共用規則索引](../../../shared-rules/README.md)、[dependency-reading](../../../shared-rules/dependency-reading.md)、[neutral-language](../../../shared-rules/neutral-language.md)、[goal-action-validation](../../../shared-rules/goal-action-validation.md) 與 [feedback-lessons](../../../shared-rules/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-13 - Netty Response channelRead Captures Set-Cookie Headers at Proxy Level

Status: candidate

#### One-line Summary

當 App 使用 Netty 作為 proxy 層時，hook `channelRead` on `HttpResponseDecoder` / `HttpObjectAggregator` 可以捕獲完整的 HTTP response headers（包括 `Set-Cookie`），即使 response body 是空的或無法直接讀取。

#### Human Explanation

在分析 App 的 session 管理時，需要從 HTTP response 中提取 `Set-Cookie` header 來取得 session token。如果 App 使用 Netty 作為 proxy 層，標準的 Java Frida hook（如 `HttpURLConnection`、`OkHttp`）可能無法捕獲 response，因為流量經過 Netty pipeline 而非標準 Java HTTP 客戶端。

**實務案例**：某 App 的 guest login response 包含多個 `Set-Cookie` headers，但這些 cookie 只在 Netty proxy 層的 response path 中出現。透過 hook `channelRead` on `HttpResponseDecoder` / `HttpObjectAggregator`，可以：

1. 使用 `Java.cast()` 到 `HttpResponse` interface 來讀取 status code 和 headers
2. 使用 `headers().getAll('Set-Cookie')` 來提取所有 cookie 值
3. 使用 `toString()` 作為 fallback（當 `headers()` 無法直接呼叫時）

**關鍵教訓**：不要只 hook request path。Response path 的 `channelRead` 是捕獲 session cookie 的最佳位置。

#### Trigger

- App 使用 Netty 作為 proxy 或 HTTP 伺服器
- 需要從 HTTP response 中提取 `Set-Cookie` 或 session token
- Java-level `HttpURLConnection` / `OkHttp` hook 沒有捕獲到 response
- Frida capture 看到 request 但看不到 response headers

#### Evidence

- Tool: Frida（hook Netty `channelRead` on response decoder classes）
- Sanitized excerpt: `channelRead` on `HttpResponseDecoder` / `HttpObjectAggregator` 成功捕獲 guest login response 的 `Set-Cookie` headers
- Evidence path: `<PROJECT_ROOT>/capture/`（cookie capture logs）、`<PROJECT_ROOT>/capture/self_gen_phase1_v10_analysis.md`

#### Generalized Lesson

當需要從 Netty-based App 捕獲 HTTP response headers（特別是 session cookie）時：

1. **Hook response path**：`channelRead` on `HttpResponseDecoder` / `HttpObjectAggregator` — 這是 Netty pipeline 中 response 經過的位置
2. **使用 Java.cast() 到 HttpResponse**：`Java.cast(msg, Java.use('io.netty.handler.codec.http.HttpResponse'))` 可以存取 status code 和 headers
3. **使用 headers().getAll('Set-Cookie')**：Netty 的 `HttpHeaders.getAll()` 可以返回所有同名 header 值（多個 `Set-Cookie` header）
4. **toString() 作為 fallback**：如果 `headers()` 無法直接呼叫，`toString()` 仍可能輸出完整的 header 列表
5. **不要 truncate**：Cookie 值可能很長，確保 log 輸出沒有長度限制

#### Agent Action

下次 agent 需要從 Netty-based App 捕獲 session cookie 時：

1. 註冊 `channelRead` hook on `HttpResponseDecoder` / `HttpObjectAggregator`
2. 在 hook 中檢查 msg 是否為 `DefaultHttpResponse` / `DefaultFullHttpResponse`
3. 嘗試 `Java.cast()` 到 `HttpResponse` interface
4. 使用 `headers().getAll('Set-Cookie')` 提取 cookie 值
5. 如果 cast 失敗，使用 `toString()` 作為 fallback（不 truncate）
6. 同時 hook `ProxyServerHandler.channelRead` 來捕獲 proxy 層的 response

#### Goal / Action / Validation

- Goal: 讓 agent 能從 Netty-based App 的 response path 捕獲 session cookie
- Action: hook `channelRead` on response decoder classes，使用 `Java.cast()` + `headers().getAll('Set-Cookie')`
- Validation or reference source: 實測驗證 — Netty response channelRead hook 成功捕獲 `Set-Cookie` headers

#### Applies When

- App 使用 Netty 作為 proxy 或 HTTP 伺服器
- 需要從 HTTP response 中提取 session cookie 或 header
- Java-level HTTP hook 無法捕獲 response
- 目標是分析 session 管理或 token 生命週期

#### Does Not Apply When

- App 不使用 Netty（使用 OkHttp / HttpURLConnection 直接連線）
- 只需要 request 資訊（不需要 response headers）
- Session cookie 在非 HTTP 層（如 WebSocket、自訂 protocol）

#### Validation

- 註冊 `channelRead` hook on `HttpResponseDecoder` / `HttpObjectAggregator`
- 觸發目標 API 請求（如 guest login）
- 確認 log 中有 `Set-Cookie` header 輸出
- 如果沒有，檢查 msg class name 是否為預期的 response 類型

#### Promotion Target

- `analysis/apk/workflows/dynamic-capture-flow.md`（新分層）

#### Required Linked Updates

- 無需連動更新；此為新 lesson，尚未 promote 到正式文件
