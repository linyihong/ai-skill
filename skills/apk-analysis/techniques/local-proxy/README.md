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

## Related Lessons

- `../../feedback_history/2026-04-30_120010-本機-loopback-proxyserver-轉發會讓-wi-fi-http-mitm-看不到業務-connect.md`
- `../../feedback_history/2026-05-01_114300-local-proxy-handler-uri-hook.md`
- `../../feedback_history/2026-05-01_131000-cast-netty-request-for-handler-route.md`
- `../../feedback_history/2026-05-01_132400-netty-aggregated-request-tostring-headers.md`
