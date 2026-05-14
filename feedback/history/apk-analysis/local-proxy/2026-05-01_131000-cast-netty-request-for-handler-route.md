> 遵守 [共用規則索引](../../../../enforcement/README.md) 與 [feedback-lessons](../../../../enforcement/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。
# Extracted — See [`analysis/apk/tools-and-failures.md`](../../../../analysis/apk/tools-and-failures.md)

### 2026-05-01 - Cast Netty Request Interfaces For Handler Routes

Status: validated

#### One-line Summary

Frida hook 到 Netty `FullHttpRequest` 參數時，直接呼叫 `req.method()` / `req.uri()` 可能失敗；先 cast 到 `HttpRequest` / `FullHttpRequest` 介面再讀 method/URI 更可靠。

#### Human Explanation

在 Android Java/Kotlin APK 中 hook Netty handler 時，Frida 拿到的參數不一定直接暴露所有 interface 方法。即使方法簽名顯示參數是 `FullHttpRequest`，直接 `req.method()`、`req.uri()` 可能回傳失敗或被 JavaScript property resolution 影響。

實務上可先 `Java.cast(req, Java.use('io.netty.handler.codec.http.HttpRequest'))`，再嘗試 `method()` / `getMethod()` 與 `uri()` / `getUri()`。這能在不讀 body、不列敏感 header 的前提下，取到 request method 和 path，配合 handler 解析出的 upstream `URI` 還原本機 proxy 的 route mapping。

#### Trigger

- Hook 到 Netty / local proxy handler 的 `FullHttpRequest` 參數。
- Handler 的 upstream `URI` 已可見，但 local request method/path 顯示 `<method?>` 或 `<uri?>`。
- 直接呼叫 request object 方法失敗，但 handler 仍正常轉發。

#### Evidence

- Tool: Frida targeted hook on a local proxy / Netty handler.
- Sanitized excerpt:
  - Before cast: handler hook logs upstream URI but request method/path are unknown.
  - After cast: handler hook logs `GET /<path>` and `POST /<path>` route metadata, query redacted.
- Evidence path: project-local `capture/` logs and sanitized notes; do not copy target hosts, raw queries, body, tokens, cookies, or local paths into reusable skill docs.

#### Generalized Lesson

For Netty request objects in Frida, treat Java interface casting as part of the hook. Cast to `HttpRequest` or `FullHttpRequest`, then read both old and new Netty method names (`getMethod`/`method`, `getUri`/`uri`). Keep output to method, path, and redacted route metadata unless a private capture explicitly requires more.

#### Agent Action

When a Netty handler hook receives a request object:

1. Try `Java.cast` to `io.netty.handler.codec.http.HttpRequest`.
2. Try `Java.cast` to `io.netty.handler.codec.http.FullHttpRequest`.
3. Read method via `method()` then `getMethod()`.
4. Read path via `uri()` then `getUri()`.
5. Redact query values and do not log body or sensitive headers by default.
6. Pair the local method/path with any resolved upstream `URI` parameter to document route mapping.

#### Applies When

- APK uses Netty, local ProxyServer, loopback handler, or embedded HTTP server classes.
- Frida hook signatures expose `FullHttpRequest`, `HttpRequest`, `ChannelHandlerContext`, or upstream `URI`.
- Goal is route/path confirmation, not full request body capture.

#### Does Not Apply When

- The handler uses a different HTTP framework without Netty request interfaces.
- The route data is only available in encrypted body/header fields.
- The analysis is already at a higher semantic layer such as app request options or Dart interceptor objects.

#### Validation

Validated when the same handler hook changes from unknown request metadata to visible sanitized `METHOD /path` values after interface casting, while the app continues to run and upstream URI logging still works.

#### Promotion Target

- `TOOLS.md`
- `WORKFLOW.md`
