> 遵守 [共用規則索引](../../../../shared-rules/README.md) 與 [feedback-lessons](../../../../shared-rules/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-01 - Hook Local Proxy Handler URI, Not Just OkHttp

Status: validated

#### One-line Summary

App 內建 loopback ProxyServer 會讓 PC MITM 與 Java OkHttp hook 都只看到局部；hook handler 的 `FullHttpRequest` + `URI` 參數可直接確認上游業務 host。

#### Human Explanation

有些 APK 會在 App 內啟動本機 HTTP proxy / handler。前端或 SDK 先打 `127.0.0.1:<local-port>`，再由 handler 選擇真實上游 API host。這種情況下：

- PC 端 MITM 綁 Android global/Wi-Fi proxy，不一定在 handler 後面的外連路徑上。
- OkHttp `newCall` / `RealCall` hook 可能只看到 loopback、配置、校時或第三方。
- 只看 native `connect` 只能看到外部 IP，仍缺少上游 host/path。

若反射或方法探測發現本機 proxy handler 有 `FullHttpRequest`、`URI`、Netty `ChannelHandlerContext` 等參數，應優先 hook 這些 handler 方法。`URI` 參數通常已是 handler 解析後的上游目標，比 socket bytes 更高語意。

#### Trigger

- MITM 看不到業務 host，但 App 功能正常。
- OkHttp hook 看到 `127.0.0.1:<local-port>`、`/check`、`/op` 或本機 proxy health route。
- 類名、logcat 或反射結果出現 `ProxyServer`、`ProxyServerHandler`、`FullHttpRequest`、Netty handler。

#### Evidence

- Tool: Frida method signature probe + targeted handler hook.
- Sanitized excerpt:
  - `ProxyServerHandler` declared methods include overloads accepting `FullHttpRequest` and `java.net.URI`.
  - Hooking the overload logs upstream `https://<api-host>` URI while MITM still lacks that business route.
- Evidence path: project-local `capture/` logs and sanitized project notes; do not copy target hosts, tokens, header values, or private paths into reusable skill docs.

#### Generalized Lesson

When an APK uses a local proxy / Netty handler, the best hook point may be the handler method that receives both the local request object and the resolved upstream `URI`. This confirms route mapping without reconstructing TLS/socket bytes and without depending on PC MITM visibility.

#### Agent Action

When local proxy signs appear:

1. Probe proxy classes and declared methods first.
2. Look for overloads containing `FullHttpRequest`, `URI`, `HttpHeaders`, `ChannelHandlerContext`, or similar request/route objects.
3. Hook only method, sanitized URI, and non-sensitive route metadata first.
4. Avoid logging body, authorization, cookies, signature headers, or full query values unless explicitly required and kept in private `capture/`.
5. Use the result to distinguish:
   - local proxy is only health/config,
   - local proxy maps to business upstream,
   - route still needs Dart/native/pcap follow-up.

#### Applies When

- Android APK contains local HTTP proxy / Netty / loopback server behavior.
- MITM sees loopback or proxy-aware noise but not business upstream.
- Java OkHttp hooks are inconclusive or mostly show local URLs.

#### Does Not Apply When

- No local proxy/server/handler classes are present.
- Business host already appears in PC MITM and the remaining issue is TLS trust/pinning.
- The app uses pure Dart/native sockets without a Java-visible handler layer.

#### Validation

Validated when the same cold-start window shows:

- proxy handler method is hooked successfully,
- handler receives local request object and resolved upstream `URI`,
- upstream URI matches business-candidate hosts from DNS/SkyShield/pcap,
- PC MITM or OkHttp alone did not expose the same route.

#### Promotion Target

- `WORKFLOW.md`
- `TOOLS.md`
