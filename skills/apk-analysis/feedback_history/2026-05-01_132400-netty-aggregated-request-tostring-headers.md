> 遵守 [共用規則索引](../../../shared-rules/README.md) 與 [feedback-lessons](../../../shared-rules/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-01 - Netty Aggregated Request toString Can Expose Headers

Status: validated

#### One-line Summary

Frida 直接讀 Netty `headers()` / `content()` 失敗時，`AggregatedFullHttpRequest.toString()` 仍可能輸出 request line 與 headers，可用來確認參數是否在 header 而非 body。

#### Human Explanation

在 hook Netty / local proxy handler 時，`FullHttpRequest` 物件有時可以讀到 method/path，卻無法透過 Frida 直接呼叫 `headers()` 或 `content()`；這可能來自 Java wrapper、interface dispatch 或方法解析差異。此時不要只判斷「沒有 header」。先檢查物件的 Java `toString()`，Netty 的 aggregated request 常會列出：

- request line,
- header names and values,
- content buffer metadata.

若 `toString()` 顯示 `content-length: 0`，但同時有加密／簽名類自訂 header，則可判斷參數可能在 header 中，下一步應找 header 組裝／解碼點，而不是繼續追空 body。

#### Trigger

- `FullHttpRequest` hook 已拿到 method/path。
- `headers()` / `content()` 直接呼叫或反射讀取失敗。
- Body hook 一直顯示 empty / no buffer，但 App 業務請求仍正常。

#### Evidence

- Tool: Frida hook on Netty `channelRead` / local proxy handler.
- Sanitized excerpt:
  - `toString()` shows `GET/POST /<api-path> HTTP/1.1`.
  - `toString()` shows custom encrypted/signature-like headers.
  - `toString()` shows `content-length: 0`.
- Evidence path: project-local `capture/` logs and sanitized notes; do not copy target header values, raw body, hostnames, tokens, cookies, or local paths into reusable skill docs.

#### Generalized Lesson

When Netty request accessors fail in Frida, use Java `toString()` on the actual request object as a fallback for route/header shape. Treat raw output as sensitive, because it may include full custom headers. Use it first to answer structural questions: method/path, header names, whether body is empty, and where parameters likely live.

#### Agent Action

When body/header capture appears empty:

1. Log the actual Java class name of the request object.
2. Cast to Netty interfaces and try method/path first.
3. If `headers()` or `content()` fail, call Java `toString()` on the request object.
4. Keep raw output in private `capture/`; publish only sanitized header names and structural facts.
5. If `content-length: 0` plus custom encrypted/signature headers are present, pivot to the header builder/decrypt/signing function.

#### Applies When

- APK uses Netty, local ProxyServer, loopback server, or HTTP object aggregator.
- Frida sees `AggregatedFullHttpRequest` or similar request classes.
- Goal is to determine request shape and parameter location.

#### Does Not Apply When

- Request body is binary/multipart and `toString()` omits content.
- The app uses a non-Netty HTTP stack.
- Raw header values are not authorized to be captured.

#### Validation

Validated when accessor-based header/body hooks fail, but Java `toString()` on the same request object reveals request line, headers, and content-length metadata sufficient to infer parameter location.

#### Promotion Target

- `TOOLS.md`
- `WORKFLOW.md`
