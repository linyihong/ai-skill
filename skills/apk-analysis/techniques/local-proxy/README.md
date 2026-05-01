# Local Proxy / Loopback Techniques

Use this category when app traffic first goes through a local loopback bridge, local proxy server, Netty handler, embedded TUN/VPN route, or similar in-app routing layer.

## When To Use

- pcap or logcat shows `127.0.0.1:<port>` or loopback HTTP.
- Wi-Fi MITM does not see business CONNECT, but the app still works and pcap shows external TLS.
- Logs mention `ProxyServer`, Netty handler, local bridge, TUN, sing-box, or embedded VPN-like routing.
- Hooking OkHttp/WebView does not expose the business upstream, but local handler objects do.

## Core Guidance

- Do not treat missing Wi-Fi MITM traffic as pinning until routing is proven.
- Separate local client-to-loopback requests from loopback-to-upstream requests.
- Hook high-semantic handler objects before raw sockets.
- For Netty requests, cast to `HttpRequest` / `FullHttpRequest` before reading method/URI.
- If headers/body accessors fail, use `toString()` only for private capture; public docs should keep header names and structure, not raw values.

## Handler Hook Flow

1. Confirm loopback evidence with pcap/logcat/MITM timing, such as `127.0.0.1:<port>`, `ProxyServer`, or Netty handler logs.
2. Identify handler methods and request object types rather than adding broad socket hooks.
3. Hook the handler argument that represents the local request, such as `FullHttpRequest` plus any resolved `URI` argument.
4. If `FullHttpRequest` does not expose method/path directly in Frida, cast to `io.netty.handler.codec.http.HttpRequest` or `FullHttpRequest`, then read `method/getMethod` and `uri/getUri`.
5. Redact query values by default; keep path shape, method, and header names.
6. If method/path are visible but `headers()` or `content()` fail, call the actual Java request object's `toString()` only in private capture. Public docs should summarize structure and `content-length`, not raw sensitive values.

## Attribution

Use local proxy metadata to map local UI/client calls to upstream APIs. If the upstream request is built by Flutter/Dart after the local handler, route into `../flutter-dart-aot/` only after evidence points there.

## Related Lessons

- `../../feedback_history/local-proxy/2026-04-30_120010-本機-loopback-proxyserver-轉發會讓-wi-fi-http-mitm-看不到業務-connect.md`
- `../../feedback_history/local-proxy/2026-05-01_114300-local-proxy-handler-uri-hook.md`
- `../../feedback_history/local-proxy/2026-05-01_131000-cast-netty-request-for-handler-route.md`
- `../../feedback_history/local-proxy/2026-05-01_132400-netty-aggregated-request-tostring-headers.md`
